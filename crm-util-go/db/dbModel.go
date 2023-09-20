package db

import (
	"crm-util-go/logging"
	"database/sql"
	"encoding/json"
	"github.com/gocql/gocql"
	"reflect"
	"strings"
	"time"
)

const (
	CrmFormatDate  string = "2006-01-02T15:04:05.000Z0700"
	DbMaxFailTimes int    = 20
)

type DbAuthen struct {
	Username string
	Password string
}

type CassandraPool struct {
	Hosts              []string
	Port               int
	ConnectTimeout     time.Duration
	Keyspace           string
	ConsistencyLevel   gocql.Consistency
	DbAuthen           DbAuthen
	CQLVersion         string
	NativeProtoVersion int
	Logger             *logging.PatternLogger
}

/*
	MaxLifetime := 5 * time.Minute
*/
type DBPool struct {
	DataSourceName string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxLifetime    time.Duration
	Logger         *logging.PatternLogger
	AppName        string
}

/*
	NodeName := "galera1"
	DataSourceName := "user:password@tcp(dbhost1:3306)/db_name"
	"ccbcdv/ccbcdv#234@172.19.190.148:1555,172.19.190.157:1555/CRMOLPRD?poolMaxSessions=100"
	"crmapp:crmapp2020@tcp(172.19.208.111:3306)/CRMX2?charset=utf8&checkConnLiveness=true&timeout=5s&readTimeout=60s&writeTimeout=60s&parseTime=true&loc=Asia%2FBangkok"
*/
type DBNode struct {
	NodeName       string
	DataSourceName string
}

type DBPoolCluster struct {
	DBNode       []DBNode
	RegisterName string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	Logger       *logging.PatternLogger
}

type CrmDateTime time.Time

func (ct CrmDateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)

	if t.IsZero() {
		return []byte("null"), nil
	}

	val := "\"" + t.Format(CrmFormatDate) + "\""
	return []byte(val), nil
}

func (ct *CrmDateTime) UnmarshalJSON(b []byte) error {

	s := strings.Trim(string(b), "\"")

	if s == "null" {
		*ct = CrmDateTime(time.Time{})
		return nil
	}

	t, err := time.Parse(CrmFormatDate, s)

	if err != nil {
		return err
	}

	*ct = CrmDateTime(t)

	return err
}

func (ct *CrmDateTime) String() string {
	var output string

	t := time.Time(*ct)

	if !t.IsZero() {
		output = t.Format(CrmFormatDate)
	}

	return output
}

func (ct *CrmDateTime) Time() time.Time {
	return time.Time(*ct)
}

type Int64 sql.NullInt64

func (ni *Int64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*ni = Int64{i.Int64, false}
	} else {
		*ni = Int64{i.Int64, true}
	}
	return nil
}

func (ni *Int64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *Int64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

type Float64 sql.NullFloat64

func (nf *Float64) Scan(value interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*nf = Float64{f.Float64, false}
	} else {
		*nf = Float64{f.Float64, true}
	}
	return nil
}

func (nf *Float64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

func (nf *Float64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}

type String sql.NullString

func (ns *String) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*ns = String{s.String, false}
	} else {
		*ns = String{s.String, true}
	}

	return nil
}

func (ns *String) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *String) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}
