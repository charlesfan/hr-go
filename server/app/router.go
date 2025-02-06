package app

import (
	"github.com/gin-gonic/gin"

	"github.com/charlesfan/hr-go/config"
)

type RequestHeader struct {
	Authorization string `json:"authorization"`
}

type Router struct {
	addr   string
	router *gin.Engine
}

func NewRouter(addr string) *Router {
	return &Router{
		addr:   addr,
		router: gin.Default(),
	}
}

func (r *Router) Config(c config.Config) {
	r.router.MaxMultipartMemory = 8 << 20 // 8 MiB

	v := r.router.Group("/hr")
	v.Any("", func(c *gin.Context) {
		c.String(200, "power by Charles")
	})
}

func (r *Router) Run() {
	r.router.Run(r.addr)
}
