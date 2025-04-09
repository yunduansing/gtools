package model

import "errors"

type MyError struct {
	Code int64
	Msg  string
}

func NewMyError(code int64, msg string) *MyError {
	return &MyError{Code: code, Msg: msg}

}

func (e *MyError) Error() string {
	return e.Msg
}

func GetErrorResponse(code int64, err error, resp *Response) {
	var e *MyError
	switch {
	case errors.As(err, &e):
		resp.Code = e.Code
		resp.Msg = e.Msg
	}
	resp.Code = code
	resp.Msg = err.Error()
}
