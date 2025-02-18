package employee

import (
	"github.com/gin-gonic/gin"

	"github.com/charlesfan/hr-go/controller/resp"
	"github.com/charlesfan/hr-go/service"
	"github.com/charlesfan/hr-go/utils/log"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

type EmployeeController struct {
	employeeSrv service.EmployeeServicer
}

func NewEmployeeController() *EmployeeController {
	return &EmployeeController{
		employeeSrv: service.EmployeeService,
	}
}

func (a *EmployeeController) Login(c *gin.Context) {
	var r LoginRequest

	if err := c.Bind(&r); err != nil {
		log.Error("user binding error: ", err)
		resp.WriteResponse(c, service.ErrorCodeBadRequest, nil)

		return
	}

	token, _, errCode := a.employeeSrv.LoginByEmailPassword(r.Email, r.Password)
	resp.WriteResponse(c, errCode, token)
}
