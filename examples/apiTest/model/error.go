package model

import "errors"

type MyError struct {
	Code int
	Msg  string
}

func NewMyError(code int, msg string) *MyError {
	return &MyError{Code: code, Msg: msg}

}

func (e *MyError) Error() string {
	return e.Msg
}

func GetErrorResponse(code int, err error, resp *Response) {
	var e *MyError
	switch {
	case errors.As(err, &e):
		resp.Code = e.Code
		resp.Msg = e.Msg
	}
	resp.Code = code
	resp.Msg = err.Error()
}
