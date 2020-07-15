package task

import (
	"bytes"
	"context"
	"strings"
	"text/template"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/yaml.v2"

	"cthulhu/bot"
	"cthulhu/metrics"
	"cthulhu/store"
	"cthulhu/telegram"
)

func init() {
	Register(welcomeMessageTaskName, WelcomeMessageTask)
}

type (
	taskConfig map[string]string
	tmplData   struct {
		NewUsers string
	}
	UserSet                   map[string]struct{}
	WelcomeMessageTaskContext struct {
		UsersProcessed UserSet
	}
)

const welcomeMessageTaskName string = "WelcomeMessage"

func WelcomeMessageTask(
	ctx context.Context,
	logger log.Logger,
	config bot.Config,
	store store.Service,
	tg telegram.Service,
	metrics metrics.Service,
	args bot.TaskArgs) func() {
	return func() {
		level.Info(logger).Log("msg", "running task")

		var (
			startTime                         = time.Now()
			newUsers    UserSet               = make(UserSet)
			taskConfigs map[string]taskConfig = make(map[string]taskConfig)
			tpl         string
			tplData     tmplData
		)

		for _, arg := range args {
			taskCfg, err := readArgValue(arg.Arg.Value)
			if err != nil {
				level.Error(logger).Log("msg", "could not read task config", "err", err)
				return
			}
			taskConfigs[arg.Arg.Name] = *taskCfg
		}

		for _, cfg := range taskConfigs {
			groupID, ok := cfg["group_id"]
			if !ok {
				level.Error(logger).Log("msg", "missing group ID")
				continue
			}

			newUsers = getNewUsers(ctx, store, groupID)
			usersAlreadyProcessed, err := getUsersProcessed(ctx, store, groupID)
			if err != nil {
				level.Error(logger).Log("msg", "error while getting processed users", "err", err)
				continue
			}

			for k := range newUsers {
				if _, ok := usersAlreadyProcessed[k]; ok {
					delete(newUsers, k)
				}
			}

			tplData = tmplData{NewUsers: newUsers.String()}
			switch len(newUsers) {
			case 0:
				level.Info(logger).Log("msg", "0 new users since the last run", "group_id", groupID)
				continue
			case 1:
				tpl, ok = cfg["message_template_single"]
			default:
				tpl, ok = cfg["message_template_multiple"]
			}
			if !ok {
				level.Error(logger).Log("msg", "missing message template", "group_id", groupID)
				continue
			}
			message, err := processTemplate(tpl, tplData)
			if err != nil {
				level.Error(logger).Log("msg", "failed to process template", "err", err)
				continue
			}
			tg.SendMessage(ctx, int64FromStr(groupID), message)
			if err := setUsersProcessed(ctx, store, groupID, newUsers); err != nil {
				level.Error(logger).Log("msg", "error while setting users processed", "err", err)
			}
		}

		metrics.ObserveTasksDuration(
			welcomeMessageTaskName,
			float64(time.Since(startTime).Seconds()),
		)
	}
}

func getUsersProcessed(ctx context.Context, store store.Service, groupID string) (UserSet, error) {
	var processed UserSet = make(UserSet)

	key := welcomeMessageTaskName + groupID

	b, err := store.Read(ctx, key)
	if err != nil {
		// no context
		return processed, nil
	}

	c, err := UnmarshalWMTaskContext(b)
	if err != nil {
		return nil, err
	}

	processed = c.UsersProcessed
	return processed, nil
}

func setUsersProcessed(ctx context.Context, store store.Service, groupID string, newUsers UserSet) error {
	key := welcomeMessageTaskName + groupID

	b, err := store.Read(ctx, key)
	if err != nil {
		// context needs to be created
		return store.Create(ctx, key, MarshalWMTaskContext(&WelcomeMessageTaskContext{UsersProcessed: newUsers}))
	}

	c, err := UnmarshalWMTaskContext(b)
	if err != nil {
		return err
	}

	// merge new UserSet with existing UserSet in context
	c.UsersProcessed.Merge(newUsers)
	return store.Update(ctx, key, MarshalWMTaskContext(c))
}

func getNewUsers(ctx context.Context, store store.Service, groupID string) UserSet {
	var newUsers UserSet = make(UserSet)

	for _, v := range store.GetAll(ctx) {
		u, err := telegram.UnmarshalUpdate(v)
		if err == nil {
			if u.Message.Chat.ID == int64FromStr(groupID) {
				if u.Message.NewChatMembers != nil {
					for _, user := range *u.Message.NewChatMembers {
						newUsers[telegram.GetUserName(user)] = struct{}{}
					}
				}
			}
		}
	}
	return newUsers
}

func processTemplate(tplString string, data tmplData) (string, error) {
	var (
		t         *template.Template
		tmplBytes bytes.Buffer
	)

	t, err := template.New("tmpl").Parse(tplString)
	if err != nil {
		return "", err
	}
	if err := t.Execute(&tmplBytes, data); err != nil {
		return "", err
	}

	return tmplBytes.String(), nil
}

func readArgValue(data string) (*taskConfig, error) {
	var cfg taskConfig
	err := yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (us UserSet) Merge(b UserSet) {
	for k, v := range b {
		us[k] = v
	}
}

func (us UserSet) String() string {
	var s strings.Builder
	for k := range us {
		s.WriteString(k)
		s.WriteString(", ")
	}
	return strings.TrimRight(s.String(), ", ")
}
