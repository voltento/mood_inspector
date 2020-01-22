package notification

import (
	"errors"
	"time"
)

const timeFault = time.Second * 10

type TimeChecker interface {
	CanSendNow(time time.Time) bool
}

func NewTimeProvider(cfg *NotificationCfg) (TimeChecker, error) {
	var timeProvider TimeChecker
	if len(cfg.CertainTime) > 0 {
		timeProvider = &certainTime{certainTimes: cfg.CertainTime[0:len(cfg.CertainTime)]}
	}

	if timeProvider == nil {
		return nil, errors.New("no time provided")
	}

	return timeProvider, nil
}

type certainTime struct {
	certainTimes      []time.Time
	lastProcessedTime time.Time
}

func (c *certainTime) CanSendNow(time time.Time) bool {
	panic("implement me")
}

func timeToDurationFromStartOfDay(t time.Time) time.Duration {
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.Sub(today)
}
