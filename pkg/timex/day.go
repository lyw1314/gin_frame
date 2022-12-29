// author by lipengfei5 @2022-05-18

package timex

import (
	"database/sql/driver"
	"time"
)

type Day struct {
	time.Time
}

const (
	dayFormat = "2006-01-02"
)

func (t *Day) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+dayFormat+`"`, string(data), time.Local)
	*t = Day{now}
	return
}

func (t Day) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(dayFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, dayFormat)
	b = append(b, '"')
	return b, nil
}
func (t Day) String() string {
	return t.Format(dayFormat)
}

// 使用string，防止mysql的date类型查询时，数据映射错误的问题
func (t Day) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t *Day) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*t = Day{vt}
	case string:
		tTime, _ := time.Parse("2006/01/02 15:04:05", vt)
		*t = Day{tTime}
	}
	return nil
}
