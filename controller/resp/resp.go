package resp

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func WriteResponse(c *gin.Context, errCode int, data interface{}) {
	coder := ParseError(errCode)
	c.JSON(coder.HTTPStatus(), &Response{
		Code: coder.Code(),
		Msg:  coder.Text(),
		Data: data,
	})

}
