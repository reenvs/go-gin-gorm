package util

import "testing"

func TestRandom(t *testing.T) {
	type args struct {
		ids []uint32
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "test", args: args{ids: []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}},
		{name: "test", args: args{ids: []uint32{}}},
		{name: "test", args: args{ids: []uint32{11, 12, 13, 14, 15, 16, 17, 18, 19, 110}}},
		{name: "test", args: args{ids: []uint32{12, 22, 32, 42, 52, 62, 72, 82, 92, 102}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Random(tt.args.ids)
			t.Logf("res:%v", tt.args.ids)
		})
	}
}
