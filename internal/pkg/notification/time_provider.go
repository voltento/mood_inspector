package notification

import (
	"errors"
	"math"
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
	lastProcessedTime *time.Time
}

func (c *certainTime) CanSendNow(t time.Time) bool {
	//if c.lastProcessedTime != nil {
	//
	//	diff := timeToDurationFromStartOfDay(*c.lastProcessedTime) - timeToDurationFromStartOfDay(t)
	//	if diff
	//}
	panic("not implemented")
}

func timeToDurationFromStartOfDay(t time.Time) time.Duration {
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.Sub(today)
}

func durationDiffAbs(l time.Duration, r time.Duration) time.Duration {
	const x = time.Nanosecond
	lMs := int64(l / x)
	rMs := int64(r / x)
	return time.Duration(math.Abs(float64(lMs-rMs))) * x
}
