package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/charlesfan/hr-go/controller/common"
	"github.com/charlesfan/hr-go/controller/resp"
	"github.com/charlesfan/hr-go/service"
	"github.com/charlesfan/hr-go/utils/log"
)

type RequestHeader struct {
	Authorization string `json:"authorization"`
}

func TokenVerify(c *gin.Context) {
	rh := &RequestHeader{}
	if err := c.ShouldBindHeader(rh); err != nil {
		resp.WriteResponse(c, service.ErrorCodeForbidden, nil)
		c.Abort()
		return
	}

	header := rh.Authorization
	auth := strings.SplitN(header, " ", 2)
	if len(auth) == 2 && auth[0] == "Bearer" {
		token := auth[1]
		claims, errCode := service.AuthenticationService.Verify(token)
		if errCode != service.ErrorCodeSuccess {
			resp.WriteResponse(c, errCode, nil)
			c.Abort()
			return
		}
		log.Debugf("%+v", claims)
		c.Set(common.TokenClaims, *claims)
	} else {
		resp.WriteResponse(c, service.ErrorCodeForbidden, nil)
		c.Abort()
		return
	}
}
