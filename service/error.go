package service

/**

custom error code

http status code + category + serial number

xxx 00 xx General
xxx 01 xx Employee

*/

import (
	"fmt"
	"reflect"
	"strconv"
)

type ErrorCode int

const ErrorCodeSuccess = 2000000

// 400 00
const (
	ErrorCodeBadRequest = iota + 4000000
	ErrorCodeParameterInvalid
)

// 401 00
const (
	ErrorCodeTokenInvalid = iota + 4010000
	ErrorCodeTokenExpired
)

// 401 01

// 403 00
const (
	ErrorCodeForbidden = iota + 4030000
)

// 404 00
const (
	ErrorCodeUserNotFound = iota + 4040000
)

// 404 01
const (
	ErrorCodeDataGetFail = iota + 4040100
)

// 500 00
const (
	ErrorCodeDatabaseFail = iota + 5000000
	ErrorCodeTokenCreateFail
)

var errorMsg = map[int]string{
	ErrorCodeSuccess:          "success",
	ErrorCodeBadRequest:       "bad request",
	ErrorCodeParameterInvalid: "parameter invalid",
	ErrorCodeTokenInvalid:     "token invalid",
	ErrorCodeTokenCreateFail:  "token create fail",
	ErrorCodeForbidden:        "Forbidden",
	ErrorCodeTokenExpired:     "token expired",
	ErrorCodeUserNotFound:     "user is not exist",
	ErrorCodeDatabaseFail:     "database failure",
	ErrorCodeDataGetFail:      "fail to get data",
}

func ErrorMsg(code int) string {
	return errorMsg[code]
}

func ErrorStatusCode(code int) int {
	i, _ := strconv.Atoi(fmt.Sprintf("%d", code)[:3])
	return i
}

func NewErrors(s ...interface{}) *errors {
	r := &errors{
		[]interface{}{},
	}

	for _, v := range s {
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			v1 := reflect.ValueOf(v)
			for i := 0; i < v1.Len(); i++ {
				r.s = append(r.s, v1.Index(i).Interface())
			}
		default:
			r.s = append(r.s, v)
		}
	}

	return r
}

type errors struct {
	s []interface{}
}

func (m *errors) Add(s interface{}) *errors {
	m.s = append(m.s, s)
	return m
}

func (m *errors) Error() []interface{} {
	return m.s
}

func (m *errors) NotEmpty() bool {
	return len(m.s) > 0
}
