package notification

import (
	"github.com/spf13/viper"
	"strings"
	"testing"
	"time"
)

func buildTime(s string) time.Time {
	t, err := time.Parse(time.Kitchen, s)
	if err != nil {
		panic(err.Error())
	}
	return t
}

func Test_loadReminderConfig(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name               string
		args               args
		want               NotificationCfg
		hasReadConfigError bool
		hasUnmarshalError  bool
	}{
		{
			name: "simple test",
			args: args{
				data: `{
      "name": "simple_reminder",
      "message": "notifications message",
      "random_message": ["message 1", "message 2"],
      "random_time": {
        "from": "9:00AM",
        "to": "10:00AM",
        "period": "1h",
        "extra_period": "30m"
      },
      "certain_time": ["10:00AM"]
    }`,
			},
			want: NotificationCfg{
				Name:          "simple_reminder",
				Message:       "notifications message",
				RandomMessage: []string{"message 1", "message 2"},
				RandomTime: RandomTimeCfg{
					From:        buildTime("9:00AM"),
					To:          buildTime("10:00AM"),
					Period:      time.Hour,
					ExtraPeriod: time.Minute * 30,
				},
				CertainTime: []time.Time{buildTime("10:00AM")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configurator := viper.New()
			configurator.SetConfigType("json")
			err := configurator.ReadConfig(strings.NewReader(tt.args.data))
			if hasError := err != nil; hasError != tt.hasReadConfigError {
				if err != nil {
					t.Errorf("got unexpected read error: %v", err.Error())
				} else {
					t.Errorf("expected read error")
				}
			}

			result := NotificationCfg{}
			err = configurator.Unmarshal(&result, DecoderConfigOptions())
			if hasError := err != nil; hasError != tt.hasUnmarshalError {
				if err != nil {
					t.Errorf("got unexpected unmarshal error: %v", err.Error())
				} else {
					t.Errorf("expected unmarshal error")
				}
			}

			if !result.Equal(&tt.want) {
				t.Errorf("loadReminderConfig() = %v, want %v", result, tt.want)
			}
		})
	}
}
