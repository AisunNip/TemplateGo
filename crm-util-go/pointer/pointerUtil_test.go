package pointer

import (
	"fmt"
	"testing"
	"time"
)

func TestPointer(t *testing.T) {
	// & means "address of"
	// * means "value of"

	// declare pointer with null value
	var stringPointer *string
	var f64Pointer *float64
	var f32Pointer *float32
	var i64Pointer *int64
	var i32Pointer *int32
	var i16Pointer *int16
	var i8Pointer *int8
	var iPointer *int
	var boolPointer *bool
	var timePointer *time.Time

	fmt.Println("String Val:", GetStringValue(stringPointer))
	fmt.Println("Float64 Val:", GetFloat64Value(f64Pointer))
	fmt.Println("Float32 Val:", GetFloat32Value(f32Pointer))
	fmt.Println("Int64 Val:", GetInt64Value(i64Pointer))
	fmt.Println("Int32 Val:", GetInt32Value(i32Pointer))
	fmt.Println("Int16 Val:", GetInt16Value(i16Pointer))
	fmt.Println("Int8 Val:", GetInt8Value(i8Pointer))
	fmt.Println("Int Val:", GetIntValue(iPointer))
	fmt.Println("Boolean Val:", GetBooleanValue(boolPointer))
	fmt.Println("Time Val:", GetTimeValue(timePointer))

	str1 := NewString("abc")
	str2 := "Abc"
	fmt.Println("EqualsString:", EqualsString(str1, str2))
	fmt.Println("EqualsIgnoreCaseString:", EqualsIgnoreCaseString(str1, str2))
}
