package utils

import (
	//"fmt"
	"time"
)

/*
func GetDate(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04")
}
func GetDateMH(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format("01-02 03:04")
}*/
func GetDate(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02")
}

func GetDateMH(timestamp time.Time) string {
	return timestamp.Format("2006-01-02 15:04")
}

func GetTimeParse(times string) int64 {
	if "" == times {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04", times, loc)
	return parse.Unix()
}

func GetCurrentTime() string {

	return time.Now().Format("2006-01-02 15:04")
}

func GetDateParse(dates string) int64 {
	if "" == dates {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02", dates, loc)
	return parse.Unix()
}
