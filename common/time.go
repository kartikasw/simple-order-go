package common

import (
	"time"
)

const layout = "2006-01-02T15:04:05-07:00"

func ParseStringToTime(arg string) (time.Time, error) {
	t, err := time.Parse(layout, arg)
	return t, err
}

func ParseTimeToString(t time.Time) string {
	result := t.Format(layout)
	return result
}
