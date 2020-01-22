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
		timeProvider = &dailyCertainTime{certainTimes: cfg.CertainTime[0:len(cfg.CertainTime)]}
	}

	if timeProvider == nil {
		return nil, errors.New("no time provided")
	}

	return timeProvider, nil
}

// This type of timeChecker does not worry about date, only about time of a day
type dailyCertainTime struct {
	certainTimes      []time.Time
	lastProcessedTime *time.Time
}

func (c *dailyCertainTime) CanSendNow(t time.Time) bool {
	testD := timeToDurationFromStartOfDay(t)
	if c.lastProcessedTime != nil {
		diff := durationDiffAbs(timeToDurationFromStartOfDay(*c.lastProcessedTime), testD)
		if diff <= timeFault {
			return false
		}
	}

	for _, suggestedTime := range c.certainTimes {
		suggestedD := timeToDurationFromStartOfDay(suggestedTime)
		if suggestedD.Seconds() > testD.Seconds() {
			continue
		}
		if durationDiffAbs(suggestedD, testD) > timeFault {
			continue
		}

		c.lastProcessedTime = &t
		return true
	}

	return false
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
