package errorcode

import "fmt"

type DbError struct {
	DbName	string
	Err 	error
}

func (dbErr DbError) Error() string {
	return fmt.Sprintf("Database Name: %s Error %v", dbErr.DbName, dbErr.Err)
}

type HttpError struct {
	HttpStatusCode 	int
	Err 			error
}

func (httpErr HttpError) Error() string {
	return fmt.Sprintf("Http Status Code: %d Error %v", httpErr.HttpStatusCode, httpErr.Err)
}