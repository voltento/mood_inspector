package notification

import (
	"errors"
	"math/rand"
)

type MessageProvider interface {
	Message() string
}

func NewMessageProvider(cfg *NotificationCfg) (MessageProvider, error) {
	var messageProvider MessageProvider
	if len(cfg.Message) > 0 {
		messageProvider = &simpleMessageProvider{message: cfg.Message}
	}

	if len(cfg.RandomMessage) > 0 {
		if messageProvider != nil {
			return nil, errors.New("several message types are defined simultaneously")
		}

		messageProvider = &randomMessageProvider{messages: cfg.RandomMessage[0:len(cfg.RandomMessage)]}
	}

	if messageProvider == nil {
		return nil, errors.New("no message provided")
	}

	return messageProvider, nil
}

type simpleMessageProvider struct {
	message string
}

func (s *simpleMessageProvider) Message() string {
	return s.message
}

type randomMessageProvider struct {
	messages []string
}

func (r *randomMessageProvider) Message() string {
	i := rand.Int() % len(r.messages)
	return r.messages[i]
}
