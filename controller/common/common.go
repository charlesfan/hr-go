package common

import (
	"github.com/gin-gonic/gin"

	"github.com/charlesfan/hr-go/service"
)

const TokenClaims = "claims"

func UserInfo(c *gin.Context) (re service.CustomClaims, ok bool) {
	claims, ok := c.Get(TokenClaims)
	if !ok {
		return
	}

	re, ok = claims.(service.CustomClaims)
	return
}
