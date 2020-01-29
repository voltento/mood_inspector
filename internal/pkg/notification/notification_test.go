package notification

import (
	"github.com/voltento/mood_inspector/internal/pkg/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_notification_SendIfNeed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		timeChecker TimeChecker
		msgProvider MessageProvider
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "message_provider: simple time_checker:dailyCertainTime",
			fields: fields{
				timeChecker: &dailyCertainTime{
					certainTimes: []time.Time{time.Date(1, 1, 2, 3, 0, 0, 0, time.Local)},
				},
				msgProvider: &simpleMessageProvider{"foo"},
			},
			args: args{t: time.Date(1, 1, 2, 3, 0, 1, 0, time.Local)},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := mocks.NewMockSender(mockCtrl)
			sender.EXPECT().Send(tt.want).Times(1)

			n := &notification{
				timeChecker: tt.fields.timeChecker,
				msgProvider: tt.fields.msgProvider,
			}
			n.SendIfNeed(tt.args.t, sender)
			n.SendIfNeed(tt.args.t, sender)
		})
	}
}
