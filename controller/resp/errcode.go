package resp

import (
	"net/http"

	"github.com/charlesfan/hr-go/service"
)

type Coder interface {
	Code() int
	Text() string
	SetText(s string)
	HTTPStatus() int
}

type xcode struct {
	code       int
	text       string
	httpStatus int
}

func (x xcode) Code() int {
	return x.code
}

func (x xcode) Text() string {
	return x.text
}

func (x *xcode) SetText(s string) {
	x.text = s
}

func (x xcode) HTTPStatus() int {
	return x.httpStatus
}

func (x xcode) Error() string {
	return x.text
}

func ParseError(code int) Coder {
	if code == 0 {
		return nil
	}

	xc := &xcode{
		code:       -1,
		httpStatus: http.StatusInternalServerError,
	}

	if msg := service.ErrorMsg(code); msg != "" {
		xc.code = code
		xc.text = msg
		xc.httpStatus = service.ErrorStatusCode(code)
	}

	return xc
}
