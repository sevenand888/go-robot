package utils

import "time"

// TodayDateString 获取当前日期字符串 (格式: YYYYMMDD)
func TodayDateString() string {
	return time.Now().Format("20060102")
}
