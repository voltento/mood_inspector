package notification

import (
	"github.com/pkg/errors"
	"time"
)

type Sender interface {
	Send(msg string)
}

type Notification interface {
	SendIfNeed(t time.Time, s Sender)
}

type notification struct {
	timeChecker TimeChecker
	msgProvider MessageProvider
}

func (n *notification) SendIfNeed(t time.Time, s Sender) {
	if n.timeChecker.CanSendNow(t) {
		s.Send(n.msgProvider.Message())
	}
}

func NewNotification(cfg *NotificationCfg) (Notification, error) {
	var (
		timeChecker TimeChecker
		err         error
		msgProvider MessageProvider
	)

	if timeChecker, err = NewTimeChecker(cfg); err != nil {
		return nil, errors.Wrap(err, "can not create time checker")
	}

	if msgProvider, err = NewMessageProvider(cfg); err != nil {
		return nil, errors.Wrap(err, "can not create message provider")
	}

	return &notification{timeChecker: timeChecker, msgProvider: msgProvider}, nil
}
