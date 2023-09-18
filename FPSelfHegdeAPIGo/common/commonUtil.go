package common

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

/*
output darwin, freebsd, linux, windows
*/
func GetOperatingSystem() string {
	return runtime.GOOS
}

func GetMsisdnFormat(s string) string {
	if s != "" && strings.HasPrefix(s, "0") {
		return strings.Replace(s, "0", "66", 1)
	} else {
		return s
	}
}

func GetReflect(data interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(data), reflect.ValueOf(data)
}

func RandomNumber(min int, max int) int {
	return rand.Intn(max-min) + min
}

func Round(data float64) float64 {
	return math.Round(data*100) / 100
}

func Ceil(data float64) float64 {
	return math.Ceil(data*100) / 100
}

func StringToFloat64(data string) (float64, error) {
	data = strings.Trim(data, " ")
	return strconv.ParseFloat(data, 64)
}

func Float64ToString(data float64) string {
	return fmt.Sprintf("%f", data)
}

func StringToInt(data string) (int64, error) {
	data = strings.Trim(data, " ")
	return strconv.ParseInt(data, 10, 64)
}

func IntToString(data int) string {
	return strconv.FormatInt(int64(data), 10)
}

func StringToBool(data string) (bool, error) {
	data = strings.Trim(data, " ")
	return strconv.ParseBool(data)
}

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

func Node(name string) string {
	var builder strings.Builder
	builder.WriteString("<")
	builder.WriteString(name)
	builder.WriteString(">")
	return builder.String()
}

func NodeValue(name string, value string) string {
	var builder strings.Builder
	if len(value) > 0 {
		if "%CLEAR%" == value {
			value = ""
		}
		builder.WriteString("<")
		builder.WriteString(name)
		builder.WriteString(">")
		builder.WriteString(escapeXml(value))
		builder.WriteString("</")
		builder.WriteString(name)
		builder.WriteString(">")
	}
	return builder.String()
}

func escapeXml(value string) string {
	if len(value) > 0 {
		value = strings.ReplaceAll(value, "&", "&amp;")
		value = strings.ReplaceAll(value, ">", "&gt;")
		value = strings.ReplaceAll(value, "<", "&lt;")
		value = strings.ReplaceAll(value, "\"", "&quot;")
		value = strings.ReplaceAll(value, "'", "&apos;")
	}
	return value
}

func GetLowerFirstVariable(input string) (output string) {
	output = string(unicode.ToLower([]rune(input)[0])) + input[1:]
	return
}
