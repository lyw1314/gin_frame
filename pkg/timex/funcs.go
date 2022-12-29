// author by lipengfei5 @2022-05-18

package timex

import (
	"strconv"
	"time"
)

// 今天
func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}

// 现在
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 明天
func GetTomorrowDate() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

// 昨天
func GetYesterdayDate() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}

// 上周同日
func GetLastWeekDate() string {
	return time.Now().AddDate(0, 0, -7).Format("2006-01-02")
}

// 上月同日
func GetLastMonthDate() string {
	return time.Now().AddDate(0, -1, 0).Format("2006-01-02")
}

func GetTimerLine(sourceDate string) (start, end string) {
	var s, e time.Time
	if len(sourceDate) <= 0 {
		s, _ = time.Parse("2006-01-02 15:04:05", GetTodayDate()+" 00:05:00")
		e, _ = time.Parse("2006-01-02 15:04:05", GetTodayDate()+" 23:59:59")
	} else {
		s, _ = time.Parse("2006-01-02 15:04:05", sourceDate+" 00:05:00")
		e, _ = time.Parse("2006-01-02 15:04:05", sourceDate+" 23:59:59")
	}
	return s.Format("2006-01-02 15:04:05"), e.Format("2006-01-02 15:04:05")
}

// 获取上个季度的日期
func GetLastSeasonDate(currentDate string) (start, end string) {
	var s, e time.Time
	curr, _ := time.Parse("2006-01-02", currentDate)
	month := int(curr.Month())
	switch month {
	case 1, 2, 3:
		s, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.AddDate(-1, 0, 0).Year())+"-10-01", time.Local)
		e, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.AddDate(-1, 0, 0).Year())+"-12-31", time.Local)
	case 4, 5, 6:
		s, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-01-01", time.Local)
		e, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-03-31", time.Local)
	case 7, 8, 9:
		s, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-04-01", time.Local)
		e, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-06-30", time.Local)
	case 10, 11, 12:
		s, _ = time.Parse("2006-01-02", strconv.Itoa(curr.Year())+"-07-01")
		e, _ = time.Parse("2006-01-02", strconv.Itoa(curr.Year())+"-09-30")
	default:
		break
	}
	return s.Format("2006-01-02"), e.Format("2006-01-02")
}

// 获取当前季度的日期
func GetCurrentSeasonDate(currentDate string) (start, end string) {
	var s, e time.Time
	curr, _ := time.Parse("2006-01-02", currentDate)
	month := int(curr.Month())
	switch month {
	case 1, 2, 3:
		s, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-01-01", time.Local)
		e, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-03-31", time.Local)
	case 4, 5, 6:
		s, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-04-01", time.Local)
		e, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-06-30", time.Local)
	case 7, 8, 9:
		s, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-07-01", time.Local)
		e, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-09-30", time.Local)
	case 10, 11, 12:
		s, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-10-01", time.Local)
		e, _ = time.ParseInLocation("2006-01-02", strconv.Itoa(curr.Year())+"-12-31", time.Local)
	default:
		break
	}
	return s.Format("2006-01-02"), e.Format("2006-01-02")
}

// 获取区间的天数
func GetDateRange(startdate, enddate string) (daterange []string) {
	starttime, _ := time.ParseInLocation("2006-01-02", startdate, time.Local)
	endtime, _ := time.ParseInLocation("2006-01-02", enddate, time.Local)

	days := (endtime.Unix()-starttime.Unix())/86400 + 1
	var i int64
	for i = 0; i < days; i++ {
		daterange = append(daterange, time.Unix(starttime.Unix()+i*86400, 0).Format("2006-01-02"))
	}
	return daterange
}

func GetDiffDays(startdate, enddate string) int64 {
	starttime, _ := time.ParseInLocation("2006-01-02", startdate, time.Local)
	endtime, _ := time.ParseInLocation("2006-01-02", enddate, time.Local)

	if starttime.Unix() < endtime.Unix() {
		starttime, endtime = endtime, starttime
	}
	return (starttime.Unix() - endtime.Unix()) / 86400
}

func GetOffsetDate(currdate string, offset int64) string {
	currTime, _ := time.ParseInLocation("2006-01-02", currdate, time.Local)
	return time.Unix(currTime.Unix()+86400*offset, 0).Format("2006-01-02")
}
