package notification

import (
	"github.com/voltento/mood_inspector/pkg/daytime"
	"testing"
	"time"
)

func Test_timeToDurationFromStartOfDay(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "Hours",
			args: args{
				t: time.Date(1, 1, 2, 3, 0, 0, 0, time.Local),
			},
			want: time.Hour * 3,
		},
		{
			name: "Minutes",
			args: args{
				t: time.Date(2020, 1, 2, 0, 10, 0, 0, time.Local),
			},
			want: time.Minute * 10,
		},
		{
			name: "Seconds",
			args: args{
				t: time.Date(1, 1, 3, 0, 0, 100, 0, time.Local),
			},
			want: time.Second * 100,
		},
		{
			name: "Nanoseconds",
			args: args{
				t: time.Date(1, 1, 3, 0, 0, 0, 9, time.Local),
			},
			want: time.Nanosecond * 9,
		},
		{
			name: "Hour + minutes + seconds + nanoseconds",
			args: args{
				t: time.Date(1, 1, 3, 1, 2, 3, 4, time.Local),
			},
			want: time.Hour*1 + time.Minute*2 + time.Second*3 + time.Nanosecond*4,
		},
		{
			name: "america local hour + minutes + seconds + nanoseconds",
			args: args{
				t: time.Date(1, 1, 3, 1, 2, 3, 4, loadLocationOrPanic("America/New_York")),
			},
			want: time.Hour*1 + time.Minute*2 + time.Second*3 + time.Nanosecond*4,
		},
		{
			name: "gmt local hour + minutes + seconds + nanoseconds",
			args: args{
				t: time.Date(1, 1, 3, 1, 2, 3, 4, loadLocationOrPanic("GMT")),
			},
			want: time.Hour*1 + time.Minute*2 + time.Second*3 + time.Nanosecond*4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeToDurationFromStartOfDay(tt.args.t); got != tt.want {
				t.Errorf("timeToDurationFromStartOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func loadLocationOrPanic(local string) *time.Location {
	l, err := time.LoadLocation(local)
	if err != nil {
		panic(err.Error())
	}
	return l
}

func Test_durationDiffAbs(t *testing.T) {
	type args struct {
		l time.Duration
		r time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "hour l > r",
			args: args{
				l: time.Hour * 2,
				r: time.Hour * 1,
			},
			want: time.Hour * 1,
		},
		{
			name: "hour l < r",
			args: args{
				l: time.Hour * 2,
				r: time.Hour * 3,
			},
			want: time.Hour * 1,
		},
		{
			name: "minutes l > r",
			args: args{
				l: time.Minute * 2,
				r: time.Minute * 1,
			},
			want: time.Minute * 1,
		},
		{
			name: "minutes l < r",
			args: args{
				l: time.Minute * 5,
				r: time.Minute * 1,
			},
			want: time.Minute * 4,
		},
		{
			name: "seconds l > r",
			args: args{
				l: time.Second * 5,
				r: time.Second * 3,
			},
			want: time.Second * 2,
		},
		{
			name: "seconds l < r",
			args: args{
				l: time.Second * 0,
				r: time.Second * 3,
			},
			want: time.Second * 3,
		},
		{
			name: "milliseconds l > r",
			args: args{
				l: time.Millisecond * 100,
				r: time.Millisecond * 1,
			},
			want: time.Millisecond * 99,
		},
		{
			name: "milliseconds l < r",
			args: args{
				l: time.Millisecond * 3,
				r: time.Millisecond * 103,
			},
			want: time.Millisecond * 100,
		},
		{
			name: "nanoseconds l > r",
			args: args{
				l: time.Nanosecond * 100,
				r: time.Nanosecond * 1,
			},
			want: time.Nanosecond * 99,
		},
		{
			name: "nanoseconds l < r",
			args: args{
				l: time.Nanosecond * 3,
				r: time.Nanosecond * 103,
			},
			want: time.Nanosecond * 100,
		},
		{
			name: "hour + minutes + seconds l < r",
			args: args{
				l: time.Hour + time.Minute*59 + time.Second*60,
				r: time.Hour * 1,
			},
			want: time.Minute * 60,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := durationDiffAbs(tt.args.l, tt.args.r); got != tt.want {
				t.Errorf("durationDiffAbs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_certainTime_CanSendNow(t *testing.T) {
	type fields struct {
		certainTimes      []time.Time
		lastProcessedTime *time.Time
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "one value ok",
			fields: fields{
				certainTimes: []time.Time{time.Date(0, 0, 0, 1, 1, 10, 0, time.Local)},
			},
			args: args{
				t: time.Date(0, 0, 0, 1, 1, 11, 0, time.Local),
			},
			want: true,
		},
		{
			name: "one value field: t > expected for 10 s",
			fields: fields{
				certainTimes: []time.Time{time.Date(0, 0, 0, 1, 1, 10, 0, time.Local)},
			},
			args: args{
				t: time.Date(0, 0, 0, 1, 1, 21, 0, time.Local),
			},
			want: false,
		},
		{
			name: "one value field: t < expected for 10 s",
			fields: fields{
				certainTimes: []time.Time{time.Date(2019, 0, 0, 21, 1, 10, 0, time.Local)},
			},
			args: args{
				t: time.Date(0, 0, 5, 1, 1, 10, 0, time.Local),
			},
			want: false,
		},
		{
			name: "one value field: t < expected for 1 s",
			fields: fields{
				certainTimes: []time.Time{time.Date(1, 0, 0, 21, 1, 21, 0, time.Local)},
			},
			args: args{
				t: time.Date(0, 0, 0, 21, 1, 20, 0, time.Local),
			},
			want: false,
		},
		{
			name: "one value field: t < expected for 1 h",
			fields: fields{
				certainTimes: []time.Time{time.Date(1, 0, 0, 2, 1, 20, 0, time.Local)},
			},
			args: args{
				t: time.Date(0, 0, 0, 1, 1, 20, 0, time.Local),
			},
			want: false,
		},
		{
			name: "two values ok: second ok",
			fields: fields{
				certainTimes: []time.Time{
					time.Date(0, 0, 0, 22, 1, 10, 0, time.Local),
					time.Date(0, 0, 0, 1, 1, 10, 0, time.Local),
				},
			},
			args: args{
				t: time.Date(0, 3, 0, 1, 1, 20, 0, time.Local),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &dailyCertainTimeChecker{
				certainTimes:      tt.fields.certainTimes,
				lastProcessedTime: tt.fields.lastProcessedTime,
			}
			if got := c.Check(tt.args.t); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}

			if resend := c.Check(tt.args.t); resend {
				t.Errorf("expected Check() returns false on second call")
			}
		})
	}
}

func Test_dailyRandomTime_inPeriodSendPeriod(t *testing.T) {
	type fields struct {
		from time.Duration
		to   time.Duration
	}
	type args struct {
		t daytime.DayTime
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "in",
			fields: fields{
				from: time.Hour,
				to:   time.Hour + time.Minute*30,
			},
			args: args{daytime.NewDayTime(time.Hour + time.Minute*5)},
			want: true,
		},
		{
			name: "out",
			fields: fields{
				from: time.Hour,
				to:   time.Hour + time.Minute*30,
			},
			args: args{daytime.NewDayTime(time.Hour + time.Minute*31)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dailyRandomTimeChecker{
				from: tt.fields.from,
				to:   tt.fields.to,
			}
			if got := d.inPeriodSendPeriod(tt.args.t); got != tt.want {
				t.Errorf("inPeriodSendPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dailyRandomTime_buildNextCallTime(t *testing.T) {
	type fields struct {
		lastProcessedTime daytime.DayTime
		nextCallTime      daytime.DayTime
		from              time.Duration
		to                time.Duration
		period            time.Duration
		extraPeriod       time.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "ok",
			fields: fields{
				lastProcessedTime: daytime.NewDayTimeFromTime(time.Now()),
				nextCallTime:      daytime.NewDayTimeFromTime(time.Now().Add(time.Hour * -3)),
				from:              time.Hour * 1,
				to:                time.Hour * 2,
				period:            time.Minute * 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dailyRandomTimeChecker{
				lastProcessedTime: tt.fields.lastProcessedTime,
				nextCallTime:      tt.fields.nextCallTime,
				from:              tt.fields.from,
				to:                tt.fields.to,
				period:            tt.fields.period,
				extraPeriod:       tt.fields.extraPeriod,
			}

			for i := 0; i < 100; i += 1 {
				got := d.buildNextCallTime(daytime.NewDayTimeFromTime(time.Now()))
				if in := d.inPeriodSendPeriod(got); !in {
					t.Errorf("buildNextCallTime() produced time not in the range. Time: %v range %v:%v", got, d.from, d.to)
				}
			}
		})
	}
}

func Test_newDailyRandomTime(t *testing.T) {
	type args struct {
		config *NotificationCfg
	}
	tests := []struct {
		name    string
		args    args
		want    *dailyRandomTimeChecker
		wantErr bool
	}{
		{
			name: "ok",
			args: args{config: &NotificationCfg{
				Name:    "a",
				Message: "msg",
				RandomTime: &RandomTimeCfg{
					From:        buildTime("1:15AM"),
					To:          buildTime("7:00PM"),
					Period:      time.Minute,
					ExtraPeriod: time.Minute,
				},
			}},
			want: &dailyRandomTimeChecker{
				from:        time.Hour + time.Minute*15,
				to:          time.Hour * 19,
				period:      time.Minute,
				extraPeriod: time.Minute,
			},
			wantErr: false,
		},
		{
			name: "from > to",
			args: args{config: &NotificationCfg{
				RandomTime: &RandomTimeCfg{
					From: buildTime("7:00PM"),
					To:   buildTime("1:15AM"),
				},
			}},
			wantErr: true,
		},
		{
			name: "period is to big",
			args: args{config: &NotificationCfg{
				RandomTime: &RandomTimeCfg{
					From:   buildTime("7:00PM"),
					To:     buildTime("8:15PM"),
					Period: time.Hour + time.Minute*16,
				},
			}},
			wantErr: true,
		},
		{
			name: "no from",
			args: args{config: &NotificationCfg{
				RandomTime: &RandomTimeCfg{
					To:     buildTime("8:15PM"),
					Period: time.Hour + time.Minute*16,
				},
			}},
			wantErr: true,
		},
		{
			name: "no to",
			args: args{config: &NotificationCfg{
				RandomTime: &RandomTimeCfg{
					From:   buildTime("8:15PM"),
					Period: time.Hour + time.Minute*16,
				},
			}},
			wantErr: true,
		},
		{
			name: "no period",
			args: args{config: &NotificationCfg{
				RandomTime: &RandomTimeCfg{
					From: buildTime("8:15PM"),
					To:   buildTime("10:15PM"),
				},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newDailyRandomTimeChecker(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("newDailyRandomTimeChecker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if !tt.want.Equal(got) {
				t.Errorf("newDailyRandomTimeChecker() \ngot  = %v, \nwant %v", got, tt.want)
			}
		})
	}
}

func Test_dailyRandomTime_CanSendNow(t *testing.T) {
	type fields struct {
		lastProcessedTime daytime.DayTime
		nextCallTime      daytime.DayTime
		from              time.Duration
		to                time.Duration
		period            time.Duration
		extraPeriod       time.Duration
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "ok",
			fields: fields{
				lastProcessedTime: fromNow(-time.Second * 20),
				nextCallTime:      fromNow(-time.Second * 1),
				from:              fromNowDuration(-time.Hour * 1),
				to:                fromNowDuration(time.Hour * 1),
				period:            time.Minute,
			},
			args: args{time.Now()},
			want: true,
		},
		{
			name: "to early",
			fields: fields{
				lastProcessedTime: fromNow(-time.Second * 20),
				nextCallTime:      fromNow(-time.Second * 1),
				from:              fromNowDuration(-time.Hour * 1),
				to:                fromNowDuration(time.Hour * 1),
				period:            time.Minute,
			},
			args: args{time.Now().Add(-time.Second * 12)},
			want: false,
		},
		{
			name: "to late",
			fields: fields{
				lastProcessedTime: fromNow(-time.Second * 20),
				nextCallTime:      fromNow(-time.Second * 1),
				from:              fromNowDuration(-time.Hour * 1),
				to:                fromNowDuration(time.Hour * 1),
				period:            time.Minute,
			},
			args: args{time.Now().Add(time.Second * 12)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dailyRandomTimeChecker{
				lastProcessedTime: tt.fields.lastProcessedTime,
				nextCallTime:      tt.fields.nextCallTime,
				from:              tt.fields.from,
				to:                tt.fields.to,
				period:            tt.fields.period,
				extraPeriod:       tt.fields.extraPeriod,
			}
			if got := d.Check(tt.args.t); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
			if d.Check(tt.args.t) {
				t.Errorf("Check() returned true on double call ")
			}
		})
	}
}

func fromNow(d time.Duration) daytime.DayTime {
	return daytime.NewDayTimeFromTime(time.Now()).Add(d)
}

func fromNowDuration(d time.Duration) time.Duration {
	t := fromNow(d)
	return t.Get()
}
