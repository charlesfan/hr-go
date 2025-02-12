package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/repository/db/daos"
	"github.com/charlesfan/hr-go/repository/db/models/employee"
	"github.com/charlesfan/hr-go/service"
	"github.com/charlesfan/hr-go/test"
	"github.com/charlesfan/hr-go/test/fixture"
)

type EmployeeServiceTestCaseSuite struct {
	env     *test.Env
	service service.EmployeeServicer
}

func setupEmployeeServiceTestCase(t *testing.T) (EmployeeServiceTestCaseSuite, func(t *testing.T)) {
	s := EmployeeServiceTestCaseSuite{
		env: test.SetupEnv(t),
	}

	c := s.env.DBConfig()
	err := daos.Init(c)
	assert.Nil(t, err)
	dbFactory := daos.NewDBRepoFactory()

	a := service.NewAuthenticationService(service.AuthenticationServiceConfig{
		TokenExpired: time.Hour * 24 * 1, // one week
		Key:          config.APP_SECRET,
	})

	r := dbFactory.EmployeeRepo()

	s.service = service.NewEmployeeService(r, a)

	return s, func(t *testing.T) {
		defer s.env.Close()
	}
}

func TestEmployeeService_LoginByEmailPassword(t *testing.T) {
	s, teardownTestCase := setupEmployeeServiceTestCase(t)
	defer teardownTestCase(t)

	tt := []struct {
		name             string
		wantResponseCode int
		givenEmail       string
		givenPassword    string
		setupSubTest     test.SetupSubTest
	}{
		{
			name:             "Login Success",
			wantResponseCode: service.ErrorCodeSuccess,
			givenEmail:       fixture.GetPrepareEmployee().Email,
			givenPassword:    "foobar",
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.env.DB.Migrator().DropTable(&employee.Employee{})
				s.env.DB.AutoMigrate(&employee.Employee{})
				s.env.DB.Create(fixture.GetPrepareEmployee())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:             "Error Password",
			wantResponseCode: service.ErrorCodeForbidden,
			givenEmail:       fixture.GetPrepareEmployee().Email,
			givenPassword:    "error_password",
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.env.DB.Migrator().DropTable(&employee.Employee{})
				s.env.DB.AutoMigrate(&employee.Employee{})
				s.env.DB.Create(fixture.GetPrepareEmployee())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:             "no user",
			wantResponseCode: service.ErrorCodeDataGetFail,
			givenEmail:       fixture.GetPrepareEmployee().Email,
			givenPassword:    "foobar",
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.env.DB.Migrator().DropTable(&employee.Employee{})
				s.env.DB.AutoMigrate(&employee.Employee{})

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			tokenStr, _, errCode := s.service.LoginByEmailPassword(tc.givenEmail, tc.givenPassword)
			assert.Equal(t, errCode, tc.wantResponseCode)
			if errCode == service.ErrorCodeSuccess {
				assert.NotZero(t, tokenStr)
			}
		})
	}
}
