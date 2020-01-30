package daytime

import "time"

type DayTime struct {
	t time.Duration
}

func NewDayTimeFromTime(t time.Time) DayTime {
	d := time.Hour * time.Duration(t.Hour())
	d += time.Minute * time.Duration(t.Minute())
	d += time.Second * time.Duration(t.Second())
	d += time.Nanosecond * time.Duration(t.Nanosecond())

	return NewDayTime(d)
}

func NewDayTime(t time.Duration) DayTime {
	r := DayTime{}
	r.Set(t)
	return r
}

func (t *DayTime) cut() {
	if diff := t.t - time.Hour*24; diff > 0 {
		t.t = diff
	}
}

func (t *DayTime) Set(v time.Duration) {
	t.t = v
	t.cut()
}

func (t *DayTime) Get() time.Duration {
	return t.t
}

func (t DayTime) Add(v time.Duration) DayTime {
	t.Set(t.Get() + v)
	return t
}
