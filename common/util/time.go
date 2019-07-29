package util

import (
	"time"
)

func IntDate(t time.Time) uint32 {
	year, month, day := t.Date()
	return uint32(year*10000 + int(month)*100 + day)
}
