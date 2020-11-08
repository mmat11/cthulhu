package bot

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log/level"

	"cthulhu/store"
	"cthulhu/telegram"
)

const (
	counterCommand = "counter"
	countersKey    = "customCounters"
)

type (
	CounterArgs struct {
		Op    string
		Match string
	}
	CountersMap map[string]struct{}
)

func parseCounterArgs(args []string) (*CounterArgs, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("usage: /counter get,add,remove match")
	}

	ca := &CounterArgs{Op: args[0]}
	if len(args) == 2 {
		ca.Match = strings.ToLower(strings.Trim(args[1], "\""))
	}

	switch ca.Op {
	case "add", "remove":
		if ca.Match == "" {
			return nil, errors.New("missing argument: match")
		}
	}

	return ca, nil
}

func (s *service) handleCounter(ctx context.Context, updateReq *telegram.Update) error {
	var (
		chatID   = updateReq.Message.Chat.ID
		authorID = updateReq.Message.From.ID
		args     = telegram.CommandArgumentsSlice(updateReq.Message.CommandArguments())
	)

	if !s.Config.isMod(authorID) {
		return nil
	}

	counterArgs, err := parseCounterArgs(args)
	if err != nil {
		level.Error(s.Logger).Log(
			"msg", "error parsing args",
			"err", err,
		)
		s.Telegram.SendMessage(ctx, chatID, err.Error())
		return nil
	}

	level.Info(s.Logger).Log(
		"msg", "parsed args",
		"args", counterArgs,
	)

	counters, err := s.readCounters(ctx)
	if err != nil {
		return err
	}

	switch counterArgs.Op {
	case "get":
		clst := make([]string, 0)
		for c := range counters {
			clst = append(clst, c)
		}
		message := fmt.Sprintf("counters: %s", strings.Join(clst, ", "))
		s.Telegram.SendMessage(ctx, chatID, message)
		return nil
	case "add":
		if _, ok := counters[counterArgs.Match]; ok {
			message := fmt.Sprintf("counter already exists: %s", counterArgs.Match)
			s.Telegram.SendMessage(ctx, chatID, message)
			return nil
		}
		counters[counterArgs.Match] = struct{}{}
	case "remove":
		if _, ok := counters[counterArgs.Match]; !ok {
			message := fmt.Sprintf("counter does not exist: %s", counterArgs.Match)
			s.Telegram.SendMessage(ctx, chatID, message)
			return nil
		}
		delete(counters, counterArgs.Match)
	default:
		s.Telegram.SendMessage(ctx, chatID, fmt.Sprintf("invalid operation: %s", counterArgs.Op))
		return nil
	}

	if err := s.Store.Update(ctx, countersKey, MarshalCountersMap(counters)); err != nil {
		level.Error(s.Logger).Log(
			"msg", "error updating counters",
			"counters", counters,
			"err", err,
		)
	}

	return s.Telegram.SendMessage(ctx, chatID, fmt.Sprintf("counter updated: %s", counterArgs.Match))
}

func (s *service) readCounters(ctx context.Context) (CountersMap, error) {
	var counters = CountersMap{}
	countersB, err := s.Store.Read(ctx, countersKey)
	if err != nil {
		if err != store.ErrKeyNotFound {
			level.Error(s.Logger).Log("msg", "error reading counters", "err", err)
			return nil, err
		}
		countersB = MarshalCountersMap(counters)
		if cErr := s.Store.Create(ctx, countersKey, countersB); err != nil {
			level.Error(s.Logger).Log(
				"msg", "error creating counters",
				"counters", counters,
				"err", cErr,
			)
		}
	}
	counters, err = UnmarshalCountersMap(countersB)
	if err != nil {
		level.Error(s.Logger).Log("msg", "error unmarshaling counters")
		return nil, err
	}
	return counters, nil
}
