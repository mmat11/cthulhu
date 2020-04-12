package task

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"text/template"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/yaml.v2"

	"cthulhu/bot"
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
	welcomeMessageTaskContext struct {
		UsersProcessed []string
	}
)

const welcomeMessageTaskName string = "WelcomeMessage"

func WelcomeMessageTask(
	ctx context.Context,
	logger log.Logger,
	config bot.Config,
	store store.Service,
	tg telegram.Service,
	args bot.TaskArgs) func() {
	return func() {
		level.Info(logger).Log("msg", "running task")

		var (
			newUsersFromUpdates, newUsers []string
			taskConfigs                   map[string]taskConfig = make(map[string]taskConfig)
			tpl                           string
			tplData                       tmplData
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

			newUsersFromUpdates = getNewUsers(ctx, store, groupID)
			usersAlreadyProcessed, err := getUsersProcessed(ctx, store, groupID)
			if err != nil {
				level.Error(logger).Log("msg", "error while getting processed users", "err", err)
				continue
			}

			for _, u := range newUsersFromUpdates {
				if _, ok := usersAlreadyProcessed[u]; !ok {
					newUsers = append(newUsers, u)
				}
			}

			tplData = tmplData{NewUsers: strings.Join(newUsers, ", ")}
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
			setUsersProcessed(ctx, store, groupID, newUsers)
			newUsers = []string{}
		}
	}
}

func getUsersProcessed(ctx context.Context, store store.Service, groupID string) (map[string]struct{}, error) {
	var processed map[string]struct{} = make(map[string]struct{})

	key := welcomeMessageTaskName + groupID

	val, err := store.Read(ctx, key)
	if err != nil {
		// no context
		return processed, nil
	}

	if taskContext, ok := val.(welcomeMessageTaskContext); ok {
		for _, u := range taskContext.UsersProcessed {
			processed[u] = struct{}{}
		}
		return processed, nil
	}
	return nil, errors.New("unknown context format")
}

func setUsersProcessed(ctx context.Context, store store.Service, groupID string, newUsers []string) {
	var taskContext welcomeMessageTaskContext

	key := welcomeMessageTaskName + groupID

	val, err := store.Read(ctx, key)
	if err != nil {
		// context needs to be created
		taskContext = welcomeMessageTaskContext{UsersProcessed: newUsers}
		store.Create(ctx, key, taskContext)
		return
	}
	if taskContext, ok := val.(welcomeMessageTaskContext); ok {
		taskContext.UsersProcessed = append(taskContext.UsersProcessed, newUsers...)
		store.Update(ctx, key, taskContext)
	}
}

func getNewUsers(ctx context.Context, store store.Service, groupID string) []string {
	var newUsers []string

	for _, v := range store.GetAll(ctx) {
		if u, ok := v.(*telegram.Update); ok {
			if u.Message.Chat.ID == int64FromStr(groupID) {
				if u.Message.NewChatMembers != nil {
					for _, user := range *u.Message.NewChatMembers {
						newUsers = append(newUsers, telegram.GetUserName(user))
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
