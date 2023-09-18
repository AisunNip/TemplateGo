package validate

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

/*
https://golang.org/pkg/regexp/syntax/

"^" matches the start of input (i.e. beginning of your string)
"$" matches the end of input (i.e. the end of your string)
"()" are grouping operators

\d             digits (== [0-9])
\D             not digits (== [^0-9])
\s             whitespace (== [\t\n\f\r ])
\S             not whitespace (== [^\t\n\f\r ])
\w             word characters (== [0-9A-Za-z_])
\W             not word characters (== [^0-9A-Za-z_])

x*             zero or more x, prefer more
x+             one or more x, prefer more
x?             zero or one x, prefer one
x{n,m}         n or n+1 or ... or m x, prefer more
x{n,}          n or more x, prefer more
x{n}           exactly n x
x*?            zero or more x, prefer fewer
x+?            one or more x, prefer fewer
x??            zero or one x, prefer zero
x{n,m}?        n or n+1 or ... or m x, prefer fewer
x{n,}?         n or more x, prefer fewer
x{n}?          exactly n x
*/

var digitRegExp = regexp.MustCompile("^-?\\d+$")
var idNoRegExp = regexp.MustCompile("^\\d{13}$")
var emailRegExp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var mobileRegExp = regexp.MustCompile("^(\\+66|66|0)\\d{9}$")
var homePhoneRegExp = regexp.MustCompile("^(\\+66|66|0)\\d{8}$")

func IsDigits(digits string) bool {
	return digitRegExp.MatchString(digits)
}

func IsBoolean(data string) bool {
	_, err := strconv.ParseBool(data)

	if err == nil {
		return true
	} else {
		return false
	}
}

func IsThaiIDNo(idNo string) bool {
	return idNoRegExp.MatchString(idNo)
}

func IsEmail(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	return emailRegExp.MatchString(email)
}

func IsMobileNo(mobileNo string) bool {
	return mobileRegExp.MatchString(mobileNo)
}

func IsHomePhoneNo(phoneNo string) bool {
	return homePhoneRegExp.MatchString(phoneNo)
}

func HasStringValue(data string) bool {
	if len(data) > 0 {
		return true
	} else {
		return false
	}
}

func HasIntValue(data int) bool {
	if data > 0 {
		return true
	} else {
		return false
	}
}

func HasDateTime(data time.Time) bool {
	if data.IsZero() {
		return false
	} else {
		return true
	}
}

func CheckRequestHeader(context echo.Context) error {
	employeeid := context.Request().Header.Get("employeeid")
	if employeeid == "" {
		return errors.New("employeeid")
	}

	return nil
}
func CheckRequestBody(body map[string]string, key ...string) error {
	for _, k := range key {
		if body[k] == "" {
			return errors.New(k)
		}
	}

	return nil
}
