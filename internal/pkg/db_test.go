package pkg

import (
	"reflect"
	"testing"
)

func Test_chats_AddChat(t *testing.T) {
	type expect struct {
		ids []ID
	}
	type args struct {
		ids []ID
	}
	tests := []struct {
		name   string
		fields expect
		args   args
	}{
		{
			name:   "Add value",
			fields: expect{ids: []ID{1}},
			args:   args{ids: []ID{1}},
		},
		{
			name:   "Add two values",
			fields: expect{ids: []ID{1, 2}},
			args:   args{ids: []ID{1, 2}},
		},
		{
			name:   "Add three values",
			fields: expect{ids: []ID{1, 2, 3}},
			args:   args{ids: []ID{1, 2, 3}},
		},
		{
			name:   "Add two double values",
			fields: expect{ids: []ID{1}},
			args:   args{ids: []ID{1, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := chats{ids: make(map[ID]struct{})}
			for _, id := range tt.args.ids {
				ch.AddChat(id)
			}

			if !reflect.DeepEqual(tt.fields.ids, ch.Get()) {
				t.Errorf("Values are not equal: %v %v", tt.fields.ids, ch.Get())
			}
		})
	}
}
