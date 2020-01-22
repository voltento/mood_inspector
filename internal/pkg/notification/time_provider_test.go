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
