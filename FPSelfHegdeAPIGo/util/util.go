package util

import (
	"FPSelfHegdeAPIGo/logging"
	"strings"
	"time"
)

type Util struct {
	ConnectTimeout time.Duration
	Logger         *logging.PatternLogger
}

func (util *Util) ConvertMsisdn(msisdn string) string {

	// m := *msisdn
	var m string

	m = strings.ReplaceAll(msisdn, "+", "")

	if strings.HasPrefix(m, "66") {
		m = strings.Replace(m, "66", "0", 1)
	}

	return m
}

// GetDBDataString2Date : function convert SqlString to format date(string)
func (util *Util) GetDBDataString2Date(calDate string) *string {

	var response string
	then, err := time.Parse("2006-01-02T15:04:05Z+07:00", calDate)
	if err != nil {

		then, err = time.Parse("2006-01-02 15:04:05", calDate)
		if err != nil {
			// fmt.Println("Parse " + calDate + " to " + calDate)
			then, err = time.Parse("2006-01-02T15:04:05+07:00", calDate)
			if err != nil {
				then, err = time.Parse("2006-01-02T15:04:05Z", calDate)
				if err != nil {
					return nil
				}
			}

		}
	}

	// t := time.Now()
	// loc, _ := time.LoadLocation(t.Location().String())
	// then = then.In(loc)
	// fmt.Println(then)
	response = then.Format("2006-01-02T15:04:05")
	response = response + "+07:00"
	// fmt.Println("Parse " + calDate + " to " + response)
	return &response
}

func (util *Util) DiffDate(t1 time.Time, t2 time.Time) int {
	days := t2.Sub(t1).Hours() / 24
	return int(days)
}
