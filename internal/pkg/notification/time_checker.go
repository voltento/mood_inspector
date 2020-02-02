package notification

import (
	"github.com/pkg/errors"
	"github.com/voltento/mood_inspector/pkg/daytime"
	"github.com/voltento/mood_inspector/pkg/errorswrp"
	"math"
	"math/rand"
	"time"
)

const timeFault = time.Second * 10

type TimeChecker interface {
	Check(time time.Time) bool
}

func NewTimeChecker(cfg *NotificationCfg) (TimeChecker, error) {
	var timeProvider TimeChecker
	if len(cfg.CertainTime) > 0 {
		timeProvider = &dailyCertainTimeChecker{certainTimes: cfg.CertainTime[0:len(cfg.CertainTime)]}
	}

	if cfg.RandomTime != nil {
		if timeProvider != nil {
			return nil, errors.New("several time types provided")
		}

		if p, err := newDailyRandomTimeChecker(cfg); err != nil {
			return nil, errorswrp.Wrap(err, "can not build random timer")
		} else {
			timeProvider = p
		}
	}

	if timeProvider == nil {
		return nil, errors.New("no time provided")
	}

	return timeProvider, nil
}

type dailyRandomTimeChecker struct {
	doubleCallChecker
	lastProcessedTime daytime.DayTime
	nextCallTime      daytime.DayTime
	from              time.Duration
	to                time.Duration
	period            time.Duration
	extraPeriod       time.Duration
}

func (d *dailyRandomTimeChecker) Check(t time.Time) bool {
	if d.isDoubleCall(t) {
		return false
	}

	testTime := daytime.NewDayTimeFromTime(t)

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

func (d *dailyRandomTimeChecker) Equal(r *dailyRandomTimeChecker) bool {
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

func (d *dailyRandomTimeChecker) buildNextCallTime(newNextCallTime daytime.DayTime) daytime.DayTime {
	nextRand := func(from daytime.DayTime) daytime.DayTime {
		nextTime := from.Add(d.period)
		if d.extraPeriod != 0 {
			nextTime = nextTime.Add(time.Duration(rand.Uint64()) % d.extraPeriod)
		}
		return nextTime
	}

	for i := 0; i < 1000; i += 1 {
		newNextCallTime = nextRand(newNextCallTime)
		if d.inPeriodSendPeriod(newNextCallTime) {
			return newNextCallTime
		}
	}

	return daytime.NewDayTimeFromTime(time.Now())
}

func (d *dailyRandomTimeChecker) inPeriodSendPeriod(t daytime.DayTime) bool {
	if t.Get().Nanoseconds() > d.from.Nanoseconds() && t.Get().Nanoseconds() < d.to.Nanoseconds() {
		return true
	}

	return false
}

func newDailyRandomTimeChecker(config *NotificationCfg) (*dailyRandomTimeChecker, error) {
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
	timeCheker := &dailyRandomTimeChecker{
		lastProcessedTime: daytime.NewDayTimeFromTime(lastProcessedTime),
		from:              from,
		to:                to,
		period:            config.RandomTime.Period,
		extraPeriod:       config.RandomTime.ExtraPeriod,
	}
	timeCheker.nextCallTime = timeCheker.buildNextCallTime(daytime.NewDayTimeFromTime(time.Now()))
	return timeCheker, nil
}

// This type of timeChecker does not worry about date, only about time of a day
type dailyCertainTimeChecker struct {
	doubleCallChecker
	certainTimes      []time.Time
	lastProcessedTime *time.Time
}

func (c *dailyCertainTimeChecker) Check(t time.Time) bool {
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
