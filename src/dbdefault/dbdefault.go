package dbdefault

import (
	"timex"
	"time"
)

const (
	DbIntNull int64 = -2147483648
	DbStringNull string = "\r"
//DbTimeNull timex.Time = timex.Time(time.Time{})
	DbFloat64Null float64 = -2147483648
)

var (
	DbTimeNull timex.Time = timex.Time(time.Time{})
	SYSDATE timex.Time = timex.Time(time.Time{}.Add(time.Second))
)

func IsIntNull(a int64) bool {
	return a == DbIntNull
}

func IsStringNull(a string) bool {
	return a == DbStringNull
}

func IsFloat64Null(a float64) bool {
	return a == DbFloat64Null
}

func IsTimeNull(t timex.Time) bool {
	return t.Equal(DbTimeNull)
}

func IsStringNullOrEmpty(a string) bool {
	return IsStringNull(a) || a == ""
}