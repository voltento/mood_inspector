package notification

import (
	"github.com/voltento/mood_inspector/pkg/errorswrp"
	"time"
)

type notification struct {
	timeChecker TimeChecker
	msgProvider MessageProvider
}

func (n *notification) SendIfNeed(t time.Time, s Sender) {
	if n.timeChecker.Check(t) {
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
		return nil, errorswrp.Wrap(err, "can not create time checker")
	}

	if msgProvider, err = NewMessageProvider(cfg); err != nil {
		return nil, errorswrp.Wrap(err, "can not create message provider")
	}

	return &notification{timeChecker: timeChecker, msgProvider: msgProvider}, nil
}
