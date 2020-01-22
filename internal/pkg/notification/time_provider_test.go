package notification

import (
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
