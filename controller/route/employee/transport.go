package employee

import (
	"github.com/gin-gonic/gin"
)

func ConfigRouterGroup(group *gin.RouterGroup) {
	c := NewEmployeeController()

	employeeGroup := group.Group("/employee")
	{
		employeeGroup.POST("login", c.Login)
	}
}
