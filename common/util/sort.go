package util

import (
	"bytes"
	"fmt"
	"sort"
)

type Uint32Slice []uint32

func (c Uint32Slice) Len() int           { return len(c) }
func (c Uint32Slice) Less(i, j int) bool { return c[i] < c[j] }
func (c Uint32Slice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type Int64Slice []int64

func (c Int64Slice) Len() int           { return len(c) }
func (c Int64Slice) Less(i, j int) bool { return c[i] < c[j] }
func (c Int64Slice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type stringSorter []string

func (s stringSorter) Len() int {
	return len(s)
}

func (s stringSorter) Less(i, j int) bool {
	if s[i] < s[j] {
		return true
	} else {
		return false
	}
}

func (s stringSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func SortString(s []string) []string {
	ss := stringSorter(s)
	sort.Sort(ss)
	return ss
}

func ToOrderedString(us []uint32) string {
	var is []int
	for _, i := range us {
		is = append(is, int(i))
	}
	sort.Ints(is)

	var buffer bytes.Buffer
	for _, i := range is {
		buffer.WriteString(fmt.Sprint(i))
	}
	return buffer.String()
}
