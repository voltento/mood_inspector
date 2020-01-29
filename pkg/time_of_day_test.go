package pkg

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTimeOfDay(t *testing.T) {
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "simple",
			args: args{t: time.Hour + time.Second},
			want: time.Hour + time.Second,
		},
		{
			name: "overflow hour",
			args: args{t: time.Hour*24 + time.Second + time.Hour},
			want: time.Hour + time.Second,
		},
		{
			name: "overflow minute",
			args: args{t: time.Hour*24 + time.Minute},
			want: time.Minute,
		},
		{
			name: "overflow second",
			args: args{t: time.Hour*24 + time.Second},
			want: time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeOfDay := NewTimeOfDay(tt.args.t)
			if got := timeOfDay.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTimeOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
