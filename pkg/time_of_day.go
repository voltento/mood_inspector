package pkg

import "time"

type TimeOfDay struct {
	t time.Duration
}

func NewTimeOfDay(t time.Duration) TimeOfDay {
	r := TimeOfDay{}
	r.Set(t)
	return r
}

func (t *TimeOfDay) cut() {
	if diff := t.t - time.Hour*24; diff > 0 {
		t.t = diff
	}
}

func (t *TimeOfDay) Set(v time.Duration) {
	t.t = v
	t.cut()
}

func (t *TimeOfDay) Get() time.Duration {
	return t.t
}

func (t TimeOfDay) Add(v time.Duration) TimeOfDay {
	t.Set(t.Get() + v)
	return t
}
