package pkg

import (
	"reflect"
	"testing"
)

func buildMap(ids ...ID) map[ID]struct{} {
	result := make(map[ID]struct{})
	for _, id := range ids {
		result[id] = struct{}{}
	}
	return result
}

func Test_chats_AddChat(t *testing.T) {
	type expect struct {
		ids map[ID]struct{}
	}
	type args struct {
		ids []ID
	}
	tests := []struct {
		name   string
		expect expect
		args   args
	}{
		{
			name:   "Add value",
			expect: expect{ids: buildMap(1)},
			args:   args{ids: []ID{1}},
		},
		{
			name:   "Add two values",
			expect: expect{ids: buildMap(1, 2)},
			args:   args{ids: []ID{1, 2}},
		},
		{
			name:   "Add three values",
			expect: expect{ids: buildMap(1, 2, 3)},
			args:   args{ids: []ID{1, 2, 3}},
		},
		{
			name:   "Add two double values",
			expect: expect{ids: buildMap(1)},
			args:   args{ids: []ID{1, 1}},
		},
		{
			name:   "Empty",
			expect: expect{ids: buildMap()},
			args:   args{ids: []ID{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := chats{ids: make(map[ID]struct{})}
			for _, id := range tt.args.ids {
				ch.AddChat(id)
			}

			if !reflect.DeepEqual(tt.expect.ids, ch.Get()) {
				t.Errorf("Values are not equal: %v %v", tt.expect.ids, ch.Get())
			}
		})
	}
}
