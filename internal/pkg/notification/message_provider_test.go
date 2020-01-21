package notification

import (
	"reflect"
	"testing"
)

func TestNewMessageProvider(t *testing.T) {
	type args struct {
		cfg *NotificationCfg
	}
	tests := []struct {
		name    string
		args    args
		want    MessageProvider
		wantErr bool
	}{
		{
			name: "simpleMessageProvider",
			args: args{
				cfg: &NotificationCfg{
					Message: "foo",
				},
			},
			want:    &simpleMessageProvider{message: "foo"},
			wantErr: false,
		},
		{
			name: "randomMessageProvider",
			args: args{
				cfg: &NotificationCfg{
					RandomMessage: []string{"a1", "a2"},
				},
			},
			want:    &randomMessageProvider{messages: []string{"a1", "a2"}},
			wantErr: false,
		},
		{
			name: "double message provider defined field",
			args: args{
				cfg: &NotificationCfg{
					Message:       "foo",
					RandomMessage: []string{"a1", "a2"},
				},
			},
			wantErr: true,
		},
		{
			name: "no messages defined field",
			args: args{
				cfg: &NotificationCfg{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMessageProvider(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMessageProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageProvider() got = %v, want %v", got, tt.want)
			}
		})
	}
}
