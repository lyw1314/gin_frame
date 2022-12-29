package timex

import (
	"database/sql/driver"
	"time"
)

// type Day time.Time
type Time struct {
	time.Time
}

const localTimeFormat = "2006-01-02 15:04:05"

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+localTimeFormat+`"`, string(data), time.Local)
	*t = Time{now}
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(localTimeFormat)+2)
	b = append(b, '"')
	b = append(b, []byte(t.String())...)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	if t.IsZero() {
		return "0000-00-00 00:00:00"
	}

	return t.Format(localTimeFormat)
}

func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return "0000-00-00 00:00:00", nil
	}
	return t.Time, nil
}

func (t *Time) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*t = Time{vt}
	case string:
		tTime, _ := time.Parse("2006/01/02 15:04:05", vt)
		*t = Time{tTime}
	default:
		return nil
	}
	return nil
}

func Now() Time {
	return Time{Time: time.Now()}
}
