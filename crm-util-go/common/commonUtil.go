package common

import (
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/big"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func NewUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

/*
output: darwin, freebsd, linux, windows
*/
func GetOperatingSystem() string {
	return runtime.GOOS
}

func GetReflect(data interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(data), reflect.ValueOf(data)
}

func RandomNumber(min int64, max int64) (int64, error) {
	randomNo, err := rand.Int(rand.Reader, big.NewInt(max-min))

	if err != nil {
		return 0, err
	}

	return randomNo.Int64() + min, nil
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

func StringToInt64(data string) (int64, error) {
	data = strings.Trim(data, " ")
	return strconv.ParseInt(data, 10, 64)
}

func StringToInt(data string) (int, error) {
	data = strings.Trim(data, " ")
	i64, err := strconv.ParseInt(data, 10, 0)
	return int(i64), err
}

func IntToString(data int) string {
	return strconv.FormatInt(int64(data), 10)
}

func Int64ToString(data int64) string {
	return strconv.FormatInt(data, 10)
}

func StringToBool(data string) (bool, error) {
	data = strings.Trim(data, " ")
	return strconv.ParseBool(data)
}

func ByteToString(data []byte) string {
	return string(data)
}

func GenerateUniqueKey(prefix string) (string, error) {
	key := ""
	currTime := time.Now().UnixNano()
	randomNo, err := RandomNumber(1, 9999999999)

	if err != nil {
		return key, err
	}

	if len(prefix) > 0 {
		key = prefix + "-"
	}

	key = key + strconv.FormatInt(currTime, 36) + "-" + strconv.FormatInt(randomNo, 36)
	return key, nil
}
