package service

import (
	"time"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/repository/cache"
	"github.com/charlesfan/hr-go/repository/db/daos"
	"github.com/charlesfan/hr-go/repository/db/models/employee"
)

var (
	// === Repository ===
	employeeRepo employee.Repository
	// === Service ===
	AuthenticationService AuthenticationServicer
	EmployeeService       EmployeeServicer
)

func Init(r daos.DBRepoFactory) {
	// === Repository ===
	employeeRepo = r.EmployeeRepo()
	// === Service ===
	// authentication service
	AuthenticationService = NewAuthenticationService(AuthenticationServiceConfig{
		TokenExpired: time.Hour * 24 * 7, // one week
		Key:          config.APP_SECRET,
	}, cache.NewRedis())
	// employee service
	EmployeeService = NewEmployeeService(employeeRepo, AuthenticationService)
}
