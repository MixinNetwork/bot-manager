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
	local, _ := time.LoadLocation("Asia/Chongqing")
	_time = _time.In(local)
	return _time.Format("2006/01/02 15:04:05.999")
}
