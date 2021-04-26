package utils

import (
	"time"
)

func GetDate(num int) string {
	var today = time.Now()
	if num == 0 {
		return today.Format("2006-01-02")
	} else {
		oldDay := today.AddDate(0, 0, num)
		return oldDay.Format("2006-01-02")
	}
}

func FormatTime(_time time.Time) string {
	return _time.Format(time.RFC3339)
}
