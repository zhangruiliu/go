package timex

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"github.com/jinzhu/now"
)

var timeFormaterr = errors.New("not a valid time string")

type Time time.Time

func (this Time) Before(other Time) bool {
	t := time.Time(this)
	tt := time.Time(other)
	return t.Before(tt)
}

func Parse(s1 string, s2 string) (t Time, err error) {
	t1, err := time.Parse(s1, s2)
	if err != nil {
		panic(err)
	}

	return Time(t1), err
}
func (this Time) Equal(other Time) bool {
	t := time.Time(this)
	tt := time.Time(other)
	return t.Year() == tt.Year() && t.Month() == tt.Month() &&
		t.Day() == tt.Day() && t.Hour() == tt.Hour() && t.Minute() == tt.Minute() &&
		t.Second() == tt.Second()
}

func (this Time) ToString() string {
	t := time.Time(this)
	str := fmt.Sprintf(`%04d-%02d-%02d %02d:%02d:%02d`, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return str
}

func (this *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.Token()
		if err == io.EOF {
			// found end of element
			break
		}
		if err != nil {
			return err
		}
		if char, ok := t.(xml.CharData); ok {
			s := string([]byte(char))
			var tim time.Time
			var flag bool = true //大于19位为true  10~19 为false
			timeRune := []rune(string(s))
			if len(timeRune) >= 19 {
				s = string(timeRune[0:19])
				flag = true
			} else if len(timeRune) > 10 && len(timeRune) < 19 {
				s = string(timeRune[0:10])
				flag = false
			} else {
				return timeFormaterr
			}

			timeRune = []rune(string(s))
			year := timeRune[0:4]
			month := timeRune[5:7]
			day := timeRune[8:10]
			yearInt, err := strconv.Atoi(string(year))
			if err != nil {
				return timeFormaterr
			}

			monthInt, err := strconv.Atoi(string(month))
			if err != nil {
				return timeFormaterr
			}

			dayInt, err := strconv.Atoi(string(day))
			if err != nil {
				return timeFormaterr
			}

			tim = time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)

			if !flag {
				*this = Time(tim)
				d.Token()
				return nil
			}

			hour := timeRune[11:13]
			minute := timeRune[14:16]
			second := timeRune[17:19]
			hourInt, err := strconv.Atoi(string(hour))
			if err != nil {
				return timeFormaterr
			}
			minuteInt, err := strconv.Atoi(string(minute))
			if err != nil {
				return timeFormaterr
			}
			secondInt, err := strconv.Atoi(string(second))
			if err != nil {
				return timeFormaterr
			}
			tim = time.Date(yearInt, time.Month(monthInt), dayInt, hourInt, minuteInt, secondInt, 0, time.Local)
			*this = Time(tim)
			d.Token()
			return nil

		}
	}
	return nil
}

func (m Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(start)
	t := time.Time(m)
	str := fmt.Sprintf(`%04d-%02d-%02d %02d:%02d:%02d`, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	e.EncodeToken(xml.CharData([]byte(str)))
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

func (this Time) MarshalJSON() ([]byte, error) {
	t := time.Time(this)
	str := fmt.Sprintf(`"%04d-%02d-%02d %02d:%02d:%02d"`, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return []byte(str), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if data == nil || string(data) == `null` || len(data) < 1 {
		ti := Time(time.Time{})
		t = &ti
		return
	}

	s := strings.Trim(string(data), `"`)

	var tim time.Time
	var flag bool = true //大于19位为true  10~19 为false
	timeRune := []rune(string(s))
	if len(timeRune) >= 19 {
		s = string(timeRune[0:19])
		flag = true
	} else if len(timeRune) > 10 && len(timeRune) < 19 {
		s = string(timeRune[0:10])
		flag = false
	} else {
		return fmt.Errorf(`not valid time, ` + string(data) + ` length ` + strconv.Itoa(len(timeRune)))
	}

	timeRune = []rune(string(s))
	year := timeRune[0:4]
	month := timeRune[5:7]
	day := timeRune[8:10]
	yearInt, err := strconv.Atoi(string(year))
	if err != nil {
		return timeFormaterr
	}

	monthInt, err := strconv.Atoi(string(month))
	if err != nil {
		return timeFormaterr
	}

	dayInt, err := strconv.Atoi(string(day))
	if err != nil {
		return timeFormaterr
	}

	tim = time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.Local)

	if !flag {
		*t = Time(tim)
		return nil
	}

	hour := timeRune[11:13]
	minute := timeRune[14:16]
	second := timeRune[17:19]
	hourInt, err := strconv.Atoi(string(hour))
	if err != nil {
		return timeFormaterr
	}
	minuteInt, err := strconv.Atoi(string(minute))
	if err != nil {
		return timeFormaterr
	}
	secondInt, err := strconv.Atoi(string(second))
	if err != nil {
		return timeFormaterr
	}
	tim = time.Date(yearInt, time.Month(monthInt), dayInt, hourInt, minuteInt, secondInt, 0, time.Local)
	*t = Time(tim)
	return nil

}

const (
	YYYYMMDD = "yyyymmdd"
	Format1 = "yyyy/mm/dd"
	Format2 = "yyyy-mm-dd"
	UplodaConsignmentFormat = "yyyymmddhhmisi"
	SystemFormat = "yyyy-mm-dd hh:mi:si"
)

func NewTime(str string, format string) (ret Time, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(`not valid time string `)
		}
	}()
	switch format {
	case YYYYMMDD:
		f4 := str[0:4]
		f6 := str[4:6]
		f8 := str[6:8]
		y, err := strconv.ParseInt(f4, 10, 64)
		if err != nil {
			panic(err)
		}
		m, err := strconv.ParseInt(f6, 10, 64)
		if err != nil {
			panic(err)
		}
		r, err := strconv.ParseInt(f8, 10, 64)
		if err != nil {
			panic(err)
		}

		return Time(time.Date(int(y), time.Month(m), int(r), 0, 0, 0, 0, time.Local)), nil
		break
	case Format1, Format2:
		//2012/08/07
		f4 := str[0:4]
		f6 := str[5:7]
		f8 := str[8:10]
		y, err := strconv.ParseInt(f4, 10, 64)
		if err != nil {
			panic(err)
		}
		m, err := strconv.ParseInt(f6, 10, 64)
		if err != nil {
			panic(err)
		}
		r, err := strconv.ParseInt(f8, 10, 64)
		if err != nil {
			panic(err)
		}

		return Time(time.Date(int(y), time.Month(m), int(r), 0, 0, 0, 0, time.Local)), nil
	case SystemFormat:

		t,err := now.Parse(str)
		return Time(t),err
	default:
		panic(`not support format` + format)
	}
	return
}

var defaultTime = time.Time{}

func (t Time)Format(format string) string {
	switch format {
	case UplodaConsignmentFormat:
		if defaultTime.Equal(time.Time(t)) {
			return ""
		}
		y := time.Time(t).Year()
		m := time.Time(t).Month()
		d := time.Time(t).Day()
		h := time.Time(t).Hour()
		mi := time.Time(t).Minute()
		sec := time.Time(t).Second()
		return fmt.Sprintf("%04d%02d%02d%02d%02d%02d", y, m, d, h, mi, sec)
	}
	return ""

}
