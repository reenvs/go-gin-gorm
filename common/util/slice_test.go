package util

import (
	"reflect"
	"testing"
)

func Test_DeleteSliceItem(t *testing.T) {
	type args struct {
		src    []uint32
		target []uint32
	}
	tests := []struct {
		name string
		args args
		want []uint32
	}{
		{
			name: "正常测试",
			args: args{
				src:    []uint32{1, 2, 3, 4, 5},
				target: []uint32{1, 2, 5},
			},
			want: []uint32{3, 4},
		},
		{
			name: "要删除的元素乱序排列",
			args: args{
				src:    []uint32{1, 2, 3, 4, 5},
				target: []uint32{4, 1},
			},
			want: []uint32{2, 3, 5},
		},
		{
			name: "删除的元素有的存在有的不存在",
			args: args{
				src:    []uint32{1, 2, 3, 4, 5},
				target: []uint32{1, 3, 5, 6},
			},
			want: []uint32{2, 4},
		},
		{
			name: "删除的元素不存在",
			args: args{
				src:    []uint32{1, 2, 3, 4, 5},
				target: []uint32{6},
			},
			want: []uint32{1, 2, 3, 4, 5},
		},
		{
			name: "删除的元素有重复",
			args: args{
				src:    []uint32{1, 2, 3, 4, 5},
				target: []uint32{1, 2, 3, 3},
			},
			want: []uint32{4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteSliceItem(tt.args.src, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteSliceItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_RmDuplicate(t *testing.T) {
	type args struct {
		list []uint32
	}
	tests := []struct {
		name string
		args args
		want []uint32
	}{
		// TODO: Add test cases.
		{
			name: "有重复(刚开始正序)",
			args: args{
				list: []uint32{1, 2, 3, 4, 5, 5, 3, 1},
			},
			want: []uint32{1, 2, 3, 4, 5},
		},
		{
			name: "有重复(刚开始逆序)",
			args: args{
				list: []uint32{5, 4, 3, 2, 1, 1, 2, 3, 4, 5},
			},
			want: []uint32{5, 4, 3, 2, 1},
		},
		{
			name: "无重复",
			args: args{
				list: []uint32{1, 2, 3, 4, 5},
			},
			want: []uint32{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RmDuplicate(tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rmDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
