package timeUtil_test

import (
	"crm-util-go/timeUtil"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTimeToString(t *testing.T) {
	currDT := time.Now()
	layout := "yyyy-mm-ddThh:mi:ss.SSSZ0700"
	timeStr := timeUtil.TimeToString(layout, currDT)
	fmt.Println(timeStr)

	if len(timeStr) == 0 {
		t.Errorf("expected result yyyy-mm-ddThh:mi:ss.SSS+0700 but empty result")
		fmt.Println(timeStr)
	}
}

func TestStringToTime(t *testing.T) {
	layout := "yyyy-mm-ddThh:mi:ss.SSSZ0700"
	dateTime, err := timeUtil.StringToTime(layout, "2022-01-31T23:59:59.020+0700")

	if err != nil {
		t.Errorf("StringToTime error: " + err.Error())
	} else {
		fmt.Println("DateTime: ", dateTime)
	}
}

func TestIsPastDate(t *testing.T) {
	date := time.Now().AddDate(0, 0, -1)
	fmt.Println("date:", date)

	isPast := timeUtil.IsPastDate(date)

	if !isPast {
		t.Errorf("expected result to true but %s", strconv.FormatBool(isPast))
	}
}

func TestIsFutureDate(t *testing.T) {
	date := time.Now().AddDate(0, 0, 1)
	fmt.Println("date:", date)

	isFuture := timeUtil.IsFutureDate(date)

	if !isFuture {
		t.Errorf("expected result to true but %s", strconv.FormatBool(isFuture))
	}
}

func TestIsLeapYear(t *testing.T) {
	fmt.Println("2020 is leap year: ", timeUtil.IsLeapYear(2020))
	fmt.Println("2021 is leap year: ", timeUtil.IsLeapYear(2021))
	fmt.Println("2024 is leap year: ", timeUtil.IsLeapYear(2024))
}

func TestIsUnixZero(t *testing.T) {
	var dateTime time.Time
	fmt.Println(timeUtil.IsUnixZero(dateTime))
	fmt.Println(dateTime.IsZero())
}

func TestGetBeginningOfMonth(t *testing.T) {
	fmt.Println(timeUtil.GetBeginningOfMonth())
}

func TestGetEndOfMonth(t *testing.T) {
	fmt.Println(timeUtil.GetEndOfMonth())
}

func TestSetStartDay(t *testing.T) {
	startDT := timeUtil.SetStartDay(time.Now())
	startDTStr := timeUtil.TimeToString("yyyy-mm-dd hh:mi:ss.SSS", startDT)
	fmt.Println(startDTStr)

	if !strings.HasSuffix(startDTStr, "00:00:00.000") {
		t.Errorf("expected result to 00:00:00.000 but %s", startDTStr)
	}
}

func TestSetEndDay(t *testing.T) {
	endDT := timeUtil.SetEndDay(time.Now())
	endDTStr := timeUtil.TimeToString("yyyy-mm-dd hh:mi:ss.SSS", endDT)
	fmt.Println(endDTStr)

	if !strings.HasSuffix(endDTStr, "23:59:59.999") {
		t.Errorf("expected result to 23:59:59.999 but %s", endDTStr)
	}
}

func TestAddDays(t *testing.T) {
	currDT := time.Now()
	resultDT := timeUtil.AddDays(currDT, 1)
	fmt.Printf("Current Date: %s, Result Date: %s", currDT.String(), resultDT.String())

	duration := timeUtil.DiffTime(resultDT, currDT)
	durationStr := duration.String()
	fmt.Println(durationStr)

	if durationStr != "24h0m0s" {
		t.Errorf("expected result to 24h0m0s but %s", durationStr)
	}
}

func TestAddMonths(t *testing.T) {
	currDT := time.Now()
	resultDT := timeUtil.AddMonths(currDT, 1)
	fmt.Printf("Current Date: %s, Result Date: %s", currDT.String(), resultDT.String())
}

func TestAddYears(t *testing.T) {
	currDT := time.Now()
	resultDT := timeUtil.AddYears(currDT, 1)
	fmt.Printf("Current Date: %s, Result Date: %s", currDT.String(), resultDT.String())
}

func TestAddHours(t *testing.T) {
	currDT := time.Now()
	resultDT := timeUtil.AddHours(currDT, 1)
	fmt.Printf("Current Date: %s, Result Date: %s", currDT.String(), resultDT.String())
}

func TestAddMinutes(t *testing.T) {
	currDT := time.Now()
	resultDT := timeUtil.AddMinutes(currDT, 1)
	fmt.Printf("Current Date: %s, Result Date: %s", currDT.String(), resultDT.String())
}

func TestAddSeconds(t *testing.T) {
	currDT := time.Now()
	resultDT := timeUtil.AddSeconds(currDT, 1)
	fmt.Printf("Current Date: %s, Result Date: %s", currDT.String(), resultDT.String())
}
