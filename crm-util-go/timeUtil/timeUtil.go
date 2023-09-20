package timeUtil

import (
	"strings"
	"time"
)

func changeTimeLayout(layout string) string {
	layout = strings.ReplaceAll(layout, "yyyy", "2006")
	layout = strings.ReplaceAll(layout, "mm", "01")
	layout = strings.ReplaceAll(layout, "dd", "02")
	layout = strings.ReplaceAll(layout, "hh", "15")
	layout = strings.ReplaceAll(layout, "mi", "04")
	layout = strings.ReplaceAll(layout, "ss", "05")
	layout = strings.ReplaceAll(layout, "SSS", "000")

	return layout
}

func DiffTime(endTime time.Time, startTime time.Time) time.Duration {
	return endTime.Sub(startTime)
}

/*
input layout:
	yyyy-mm-dd hh:mi:ss
	yyyy-mm-dd hh:mi:ss.SSS
	yyyy-mm-ddThh:mi:ss.SSS
	yyyy-mm-ddThh:mi:ss.SSSZ0700
	yyyy-mm-ddThh:mi:ssZ07:00
*/
func StringToTime(layout string, dateTime string) (time.Time, error) {
	layout = changeTimeLayout(layout)
	return time.Parse(layout, dateTime)
}

/*
input layout:
	yyyy-mm-dd hh:mi:ss
	yyyy-mm-dd hh:mi:ss.SSS
	yyyy-mm-ddThh:mi:ss.SSS
	yyyy-mm-ddThh:mi:ss.SSSZ0700
	yyyy-mm-ddThh:mi:ssZ07:00
*/
func TimeToString(layout string, dateTime time.Time) string {
	var dt string

	if !dateTime.IsZero() {
		layout = changeTimeLayout(layout)
		dt = dateTime.Format(layout)
	}

	return dt
}

func SetStartDay(dateTime time.Time) time.Time {
	var startDT time.Time

	if !dateTime.IsZero() {
		startDT = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(),
			0, 0, 0, 0, dateTime.Location())
	}

	return startDT
}

func SetEndDay(dateTime time.Time) time.Time {
	var endDT time.Time

	if !dateTime.IsZero() {
		endDT = time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(),
			23, 59, 59, 999999999, dateTime.Location())
	}

	return endDT
}

func IsLeapYear(year int) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}

func IsPastDate(dateTime time.Time) bool {
	currDT := time.Now()

	currDTLoc := time.Date(currDT.Year(), currDT.Month(), currDT.Day(),
		currDT.Hour(), currDT.Minute(), currDT.Second(),
		currDT.Nanosecond(), dateTime.Location())

	return currDTLoc.After(dateTime)
}

func IsFutureDate(dateTime time.Time) bool {
	currDT := time.Now()

	currDTLoc := time.Date(currDT.Year(), currDT.Month(), currDT.Day(),
		currDT.Hour(), currDT.Minute(), currDT.Second(),
		currDT.Nanosecond(), dateTime.Location())

	return currDTLoc.Before(dateTime)
}

func IsUnixZero(dateTime time.Time) bool {
	unixZero := time.Unix(0, 0)
	return dateTime.Equal(unixZero)
}

func GetBeginningOfMonth() time.Time {
	currDT := time.Now()
	currentYear, currentMonth, _ := currDT.Date()
	currentLocation := currDT.Location()

	return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
}

func GetEndOfMonth() time.Time {
	beginningOfMonth := GetBeginningOfMonth()
	return beginningOfMonth.AddDate(0, 1, -1)
}

func AddDays(dateTime time.Time, days int) time.Time {
	return dateTime.AddDate(0, 0, days)
}

func AddMonths(dateTime time.Time, months int) time.Time {
	return dateTime.AddDate(0, months, 0)
}

func AddYears(dateTime time.Time, years int) time.Time {
	return dateTime.AddDate(years, 0, 0)
}

func AddHours(dateTime time.Time, hours int) time.Time {
	return dateTime.Add(time.Hour * time.Duration(hours))
}

func AddMinutes(dateTime time.Time, minutes int) time.Time {
	return dateTime.Add(time.Minute * time.Duration(minutes))
}

func AddSeconds(dateTime time.Time, seconds int) time.Time {
	return dateTime.Add(time.Second * time.Duration(seconds))
}
