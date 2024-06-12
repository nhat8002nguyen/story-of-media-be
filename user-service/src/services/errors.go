package services

import "fmt"

type CodeType string

var ERROR_NOT_FOUND = CodeType("404")
var ERROR_INTERNAL_SERVER = CodeType("500")
var ERROR_UNAUTHORIZED = CodeType("401")

type CustomError struct {
	Code    CodeType
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("error %s with code: %s", e.Message, e.Code)
}
