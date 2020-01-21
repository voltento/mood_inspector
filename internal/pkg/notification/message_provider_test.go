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

func Test_simpleMessageProvider_Message(t *testing.T) {
	type fields struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "simple",
			fields: fields{message: "foo"},
			want:   "foo",
		},
		{
			name:   "empty",
			fields: fields{message: ""},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := simpleMessageProvider{
				message: tt.fields.message,
			}
			if got := s.Message(); got != tt.want {
				t.Errorf("Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randomMessageProvider_Message(t *testing.T) {
	type fields struct {
		messages []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "simple",
			fields: fields{messages: []string{"foo"}},
			want:   "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := randomMessageProvider{
				messages: tt.fields.messages,
			}
			if got := r.Message(); got != tt.want {
				t.Errorf("Message() = %v, want %v", got, tt.want)
			}
		})
	}
}
