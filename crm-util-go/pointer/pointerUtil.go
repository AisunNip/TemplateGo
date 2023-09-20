package pointer

import (
	"strings"
	"time"
)

func NewString(s string) *string {
	return &s
}

func GetStringValue(s *string) string {
	var val string

	if s != nil {
		val = *s
	}

	return val
}

func EqualsString(str1 *string, str2 string) bool {
	isEqual := false

	str1Val := GetStringValue(str1)
	if str1Val == str2 {
		isEqual = true
	}

	return isEqual
}

func EqualsIgnoreCaseString(str1 *string, str2 string) bool {
	str1Val := GetStringValue(str1)
	return strings.EqualFold(str1Val, str2)
}

func NewFloat64(f float64) *float64 {
	return &f
}

func GetFloat64Value(f *float64) float64 {
	var val float64

	if f != nil {
		val = *f
	}

	return val
}

func NewFloat32(f float32) *float32 {
	return &f
}

func GetFloat32Value(f *float32) float32 {
	var val float32

	if f != nil {
		val = *f
	}

	return val
}

func NewInt64(i int64) *int64 {
	return &i
}

func GetInt64Value(i *int64) int64 {
	var val int64

	if i != nil {
		val = *i
	}

	return val
}

func NewInt32(i int32) *int32 {
	return &i
}

func GetInt32Value(i *int32) int32 {
	var val int32

	if i != nil {
		val = *i
	}

	return val
}

func NewInt16(i int16) *int16 {
	return &i
}

func GetInt16Value(i *int16) int16 {
	var val int16

	if i != nil {
		val = *i
	}

	return val
}

func NewInt8(i int8) *int8 {
	return &i
}

func GetInt8Value(i *int8) int8 {
	var val int8

	if i != nil {
		val = *i
	}

	return val
}

func NewInt(i int) *int {
	return &i
}

func GetIntValue(i *int) int {
	var val int

	if i != nil {
		val = *i
	}

	return val
}

func NewBoolean(b bool) *bool {
	return &b
}

func GetBooleanValue(b *bool) bool {
	var val bool

	if b != nil {
		val = *b
	}

	return val
}

func NewTime(t time.Time) *time.Time {
	return &t
}

func GetTimeValue(t *time.Time) time.Time {
	var val time.Time

	if t != nil {
		val = *t
	}

	return val
}
