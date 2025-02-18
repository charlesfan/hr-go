package employee_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/controller/route/employee"
	"github.com/charlesfan/hr-go/repository/db/daos"
	employeeModel "github.com/charlesfan/hr-go/repository/db/models/employee"
	"github.com/charlesfan/hr-go/service"
	"github.com/charlesfan/hr-go/test"
	"github.com/charlesfan/hr-go/test/fixture"
)

type EmployeeTestCaseSuite struct {
	env *test.Env
	c   *gin.Engine
}

func setupEmployeeTestCaseSuite(t *testing.T) (EmployeeTestCaseSuite, func(t *testing.T)) {
	s := EmployeeTestCaseSuite{
		env: test.SetupEnv(t),
		c:   gin.New(),
	}
	// prepare gin
	s.c.Use(gin.Recovery())
	s.c.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"content-type, set-sookie, *"},
	}))

	// prepare service
	a := service.NewAuthenticationService(service.AuthenticationServiceConfig{
		TokenExpired: time.Hour * 24 * 1, // one week
		Key:          config.APP_SECRET,
	}, fixture.NewCacheMock())

	c := s.env.DBConfig()
	err := daos.Init(c)
	assert.Nil(t, err)
	dbFactory := daos.NewDBRepoFactory()
	r := dbFactory.EmployeeRepo()

	// overwrite rsi UsersService
	service.EmployeeService = service.NewEmployeeService(r, a)

	employee.ConfigRouterGroup(s.c.Group("/v1"))

	return s, func(t *testing.T) {
		defer s.env.Close()
	}
}

func TestLoginHandler(t *testing.T) {
	s, teardownTestCase := setupEmployeeTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name         string
		route        string
		method       string
		body         string
		responseCode int
		setupSubTest test.SetupSubTest
	}{
		{
			name:         "Login by email",
			method:       "POST",
			route:        "/v1/employee/login",
			body:         `{"email":"em@test.com","password":"foobar"}`,
			responseCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.env.DB.Migrator().DropTable(&employeeModel.Employee{})
				s.env.DB.AutoMigrate(&employeeModel.Employee{})
				s.env.DB.Create(fixture.GetPrepareEmployee())

				return func(t *testing.T) {
				}
			},
		},
		/***
		{
			name:         "Login by error email",
			method:       "POST",
			route:        "/v1/login/email",
			body:         `{"email":"mba@liteon.com","password":"qw123456789"}`,
			responseCode: http.StatusForbidden,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.env.DB.DropTable(&user.User{})
				s.env.DB.AutoMigrate(&user.User{})

				return func(t *testing.T) {
				}
			},
		},
		{
			name:         "Login by error password",
			method:       "POST",
			route:        "/v1/login/email",
			body:         `{"email":"mba@liteon.com","password":"error_pwd"}`,
			responseCode: http.StatusForbidden,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.env.DB.DropTable(&user.User{})
				s.env.DB.AutoMigrate(&user.User{})
				s.env.DB.Create(fixture.GetPrepareUser1())
				s.env.DB.Create(fixture.GetPrepareUser2())

				return func(t *testing.T) {
				}
			},
		},
		***/
	}

	for _, tc := range tt {
		teardownSubTest := tc.setupSubTest(t)
		defer teardownSubTest(t)

		req := httptest.NewRequest(tc.method, tc.route, strings.NewReader(tc.body))
		req.Header.Set("Content-Type", gin.MIMEJSON)
		rec := httptest.NewRecorder()
		s.c.ServeHTTP(rec, req)

		assert.Equal(t, tc.responseCode, rec.Code)
	}
}
