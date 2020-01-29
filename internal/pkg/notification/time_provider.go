package notification

import (
	"github.com/pkg/errors"
	"github.com/voltento/mood_inspector/pkg"
	"math"
	"math/rand"
	"time"
)

const timeFault = time.Second * 10

type TimeChecker interface {
	CanSendNow(time time.Time) bool
}

func NewTimeChecker(cfg *NotificationCfg) (TimeChecker, error) {
	var timeProvider TimeChecker
	if len(cfg.CertainTime) > 0 {
		timeProvider = &dailyCertainTime{certainTimes: cfg.CertainTime[0:len(cfg.CertainTime)]}
	}

	if cfg.RandomTime != nil {
		if timeProvider != nil {
			return nil, errors.New("several time types provided")
		}

		if p, err := newDailyRandomTime(cfg); err != nil {
			return nil, errors.Wrap(err, "can not build random timer")
		} else {
			timeProvider = p
		}
	}

	if timeProvider == nil {
		return nil, errors.New("no time provided")
	}

	return timeProvider, nil
}

type dailyRandomTime struct {
	doubleCallChecker
	lastProcessedTime pkg.TimeOfDay
	nextCallTime      pkg.TimeOfDay
	from              time.Duration
	to                time.Duration
	period            time.Duration
	extraPeriod       time.Duration
}

func (d *dailyRandomTime) CanSendNow(t time.Time) bool {
	if d.isDoubleCall(t) {
		return false
	}

	testTime := pkg.NewTimeOfDayFromTime(t)

	if d.nextCallTime.Get() > testTime.Get() {
		return false
	}

	diff := d.nextCallTime.Get().Nanoseconds() - testTime.Get().Nanoseconds()

	if diffSeconds := time.Duration(math.Abs(float64(diff))); diffSeconds > timeFault {
		return false
	}

	d.setLastCallTime(t)
	d.nextCallTime = d.buildNextCallTime(testTime)

	return true
}

func (d *dailyRandomTime) Equal(r *dailyRandomTime) bool {
	if d.from != r.from {
		return false
	}

	if d.to != r.to {
		return false
	}

	if d.period != r.period {
		return false
	}

	if d.extraPeriod != r.extraPeriod {
		return false
	}

	return true
}

func (d *dailyRandomTime) buildNextCallTime(newNextCallTime pkg.TimeOfDay) pkg.TimeOfDay {
	nextRand := func(from pkg.TimeOfDay) pkg.TimeOfDay {
		nextTime := from.Add(d.period)
		if d.extraPeriod != 0 {
			nextTime = nextTime.Add(d.extraPeriod % time.Duration(rand.Uint64()))
		}
		return nextTime
	}

	for i := 0; i < 1000; i += 1 {
		newNextCallTime = nextRand(newNextCallTime)
		if d.inPeriodSendPeriod(newNextCallTime) {
			return newNextCallTime
		}
	}

	return pkg.NewTimeOfDayFromTime(time.Now())
}

func (d *dailyRandomTime) inPeriodSendPeriod(t pkg.TimeOfDay) bool {
	if t.Get().Nanoseconds() > d.from.Nanoseconds() && t.Get().Nanoseconds() < d.to.Nanoseconds() {
		return true
	}

	return false
}

func newDailyRandomTime(config *NotificationCfg) (*dailyRandomTime, error) {
	randomTimeConfig := *config.RandomTime
	emptyTime := time.Time{}
	if randomTimeConfig.From == emptyTime {
		return nil, errors.New("'from' was not set")
	}

	if randomTimeConfig.To == emptyTime {
		return nil, errors.New("'to' was not set")
	}

	if randomTimeConfig.Period == 0 {
		return nil, errors.New("'period' was not set")
	}

	from := timeToDurationFromStartOfDay(config.RandomTime.From)
	to := timeToDurationFromStartOfDay(config.RandomTime.To)

	if from > to {
		return nil, errors.New("wrong arguments in random time provider: from < to")
	}

	if to-from <= config.RandomTime.Period {
		return nil, errors.New("wrong arguments in random time provider: period is too big")
	}

	lastProcessedTime := time.Now().Add(-config.RandomTime.ExtraPeriod)
	timeCheker := &dailyRandomTime{
		lastProcessedTime: pkg.NewTimeOfDayFromTime(lastProcessedTime),
		from:              from,
		to:                to,
		period:            config.RandomTime.Period,
		extraPeriod:       config.RandomTime.ExtraPeriod,
	}
	timeCheker.nextCallTime = timeCheker.buildNextCallTime(pkg.NewTimeOfDayFromTime(time.Now()))
	return timeCheker, nil
}

// This type of timeChecker does not worry about date, only about time of a day
type dailyCertainTime struct {
	doubleCallChecker
	certainTimes      []time.Time
	lastProcessedTime *time.Time
}

func (c *dailyCertainTime) CanSendNow(t time.Time) bool {
	if c.isDoubleCall(t) {
		return false
	}

	for _, suggestedTime := range c.certainTimes {
		testD := timeToDurationFromStartOfDay(t)
		suggestedD := timeToDurationFromStartOfDay(suggestedTime)
		if suggestedD.Seconds() > testD.Seconds() {
			continue
		}
		if durationDiffAbs(suggestedD, testD) > timeFault {
			continue
		}

		c.setLastCallTime(t)
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
