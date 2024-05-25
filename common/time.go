package common

import "time"

func ParseStringToTime(arg string) (t time.Time, err error) {
	t, err = time.Parse("2006-01-02T15:04:05-07:00", arg)
	return
}
