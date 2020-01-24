package pkg

import (
	"github.com/voltento/mood_inspector/internal/pkg/notification"
	"gitlab.mobbtech.com/iqbus/iqbus_go_client/errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_buildConfigFromFile(t *testing.T) {
	type args struct {
		fileContent string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{

		{
			name: "reminder name",
			args: args{
				fileContent: `{
                                "notifications": [
                                  {
                                    "name": "simple_reminder"
                                  }
                                ]
                               }`,
			},
			want: &Config{Notifications: []notification.NotificationCfg{
				{
					Name: "simple_reminder",
				},
			}},
		},

		{
			name: "certain time",
			args: args{
				fileContent: `{
                                "notifications": [
                                  {
                	                "certain_time": ["10:00AM", "9:00PM"]
								  }
  							    ]
							   }`,
			},
			want: &Config{Notifications: []notification.NotificationCfg{
				{
					CertainTime: []time.Time{loadKitckenTimeOrPanic("10:00AM"), loadKitckenTimeOrPanic("9:00PM")},
				},
			}},
		},

		{
			name: "message",
			args: args{
				fileContent: `{
                                "notifications": [
                                  {
                                    "message": "notifications message"
                                  }
                                ]
                              }`,
			},
			want: &Config{Notifications: []notification.NotificationCfg{
				{
					Message: "notifications message",
				},
			}},
		},

		{
			name: "Certain reminder ok",
			want: &Config{Notifications: []notification.NotificationCfg{{
				Name:        "simple_reminder",
				Message:     "notifications message",
				CertainTime: []time.Time{loadKitckenTimeOrPanic("10:00AM"), loadKitckenTimeOrPanic("9:00PM")},
			}},
			},
			args: args{
				fileContent: `{
 								 "notifications": [
                                   {
                                      "name": "simple_reminder",
                                      "message": "notifications message",
                                      "certain_time": ["10:00AM", "9:00PM"]
                                   }
                                   ]
         				      }`,
			},
		},

		{
			name: "random_message",
			want: &Config{Notifications: []notification.NotificationCfg{{
				RandomMessage: []string{"a", "b"},
			}},
			},
			args: args{
				fileContent: `{
 								 "notifications": [
                                   {
                                      "random_message": ["a", "b"]
                                   }
                                   ]
         				      }`,
			},
		},

		{
			name: "random_time",
			want: &Config{Notifications: []notification.NotificationCfg{{
				RandomTime: &notification.RandomTimeCfg{
					From:        loadKitckenTimeOrPanic("8:29AM"),
					To:          loadKitckenTimeOrPanic("1:39PM"),
					Period:      time.Hour + time.Minute*15,
					ExtraPeriod: time.Minute * 13,
				},
			}},
			},
			args: args{
				fileContent: `{
								 "notifications": [
 								 	{	
										"random_time":  {
											"from": "8:29AM",
								    		"to": "1:39PM",
											"period": "1h15m",
											"extra_period": "13m"
                                 		 }
								 	}
								 ]
         				      }`,
			},
		},

		{
			name: "bot",
			want: &Config{TelegramBot: TelegramBot{
				Id:     "a",
				Secret: "b",
			}},
			args: args{
				fileContent: `{
								 "telegram_bot": {
									"id": "a",
									"secret": "b"
							     }
         				      }`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := ioutil.TempFile("./", "test*")
			if err != nil {
				panic(errors.Wrap(err, "can not create a temp file: %v").Error())
			}
			defer os.Remove(f.Name())

			if _, err := f.WriteString(tt.args.fileContent); err != nil {
				panic(errors.Wrap(err, "can not write to the temp file: %v").Error())
			}

			got, err := buildConfigFromFile(f.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("buildConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildConfigFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func loadKitckenTimeOrPanic(s string) time.Time {
	t, err := time.Parse(time.Kitchen, s)
	if err != nil {
		panic(err.Error())
	}
	return t
}
