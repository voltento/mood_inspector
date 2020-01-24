package notification

import "time"

type doubleCallChecker struct {
	lastProcessedTime *time.Time
}

func (d *doubleCallChecker) setLastCallTime(t time.Time) {
	d.lastProcessedTime = &t
}

func (d *doubleCallChecker) isDoubleCall(t time.Time) bool {
	testD := timeToDurationFromStartOfDay(t)
	if d.lastProcessedTime != nil {
		diff := durationDiffAbs(timeToDurationFromStartOfDay(*d.lastProcessedTime), testD)
		if diff <= timeFault {
			return true
		}
	}

	return false
}
