package daos

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/charlesfan/hr-go/repository/db/models"
	"github.com/charlesfan/hr-go/repository/db/models/employee"
	"github.com/charlesfan/hr-go/test"
	"github.com/charlesfan/hr-go/test/fixture"
)

type EmployeeTestCaseSuite struct {
	env          *test.Env
	employeeRepo employee.Repository
}

func setupEmployeeTestCase(t *testing.T) (EmployeeTestCaseSuite, func(t *testing.T)) {
	s := EmployeeTestCaseSuite{
		env: test.SetupEnv(t),
	}
	s.employeeRepo = newEmployeeRepo(s.env.DB)

	return s, func(t *testing.T) {
		defer s.env.Close()
	}
}

func TestEmployeeRepository_Get(t *testing.T) {
	s, teardownTestCase := setupEmployeeTestCase(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testPK        models.PrimaryKey[string]
		wantData      *employee.Employee
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:     "success",
			wantData: fixture.GetPrepareEmployee(),
			testPK:   fixture.GetPrepareEmployee().PK(),
			err:      nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.env.DB.Migrator().DropTable(&employee.Employee{})
				s.env.DB.AutoMigrate(&employee.Employee{})
				s.env.DB.Create(fixture.GetPrepareEmployee())
				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			p, err := s.employeeRepo.Get(tc.testPK)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				assert.Equal(t, tc.wantData.ID, p.ID)
				assert.Equal(t, tc.wantData.Name, p.Name)
				assert.Equal(t, tc.wantData.Email, p.Email)
				assert.Equal(t, tc.wantData.Department, p.Department)
				assert.Equal(t, tc.wantData.RegularizeDate, p.RegularizeDate)
			}
		})
	}
}

func TestEmployeeRepository_Create(t *testing.T) {
	s, teardownTestCase := setupEmployeeTestCase(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *employee.Employee
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:     "success",
			testData: fixture.GetPrepareEmployee(),
			err:      nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.env.DB.Migrator().DropTable(&employee.Employee{})
				s.env.DB.AutoMigrate(&employee.Employee{})
				return func(t *testing.T) {
				}
			},
		},
		{
			name:     "sameIdFail",
			testData: fixture.GetPrepareEmployee(),
			err:      errors.New("UNIQUE constraint failed: hr_employee.id"),
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.env.DB.Migrator().DropTable(&employee.Employee{})
				s.env.DB.AutoMigrate(&employee.Employee{})
				d := fixture.GetPrepareEmployee()
				d.Email = "new@test.com"
				s.env.DB.Create(d)

				return func(t *testing.T) {

				}
			},
		},
		{
			name:     "sameEmailFail",
			testData: fixture.GetPrepareEmployee(),
			err:      errors.New("UNIQUE constraint failed: hr_employee.email"),
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.env.DB.Migrator().DropTable(&employee.Employee{})
				s.env.DB.AutoMigrate(&employee.Employee{})
				d := fixture.GetPrepareEmployee()
				d.ID = "newId"
				s.env.DB.Create(d)

				return func(t *testing.T) {

				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			err := s.employeeRepo.Create(tc.testData)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				p, _ := s.employeeRepo.Get(tc.testData.PK())
				assert.Equal(t, tc.testData.ID, p.ID)
				assert.Equal(t, tc.testData.Name, p.Name)
				assert.Equal(t, tc.testData.Email, p.Email)
				assert.Equal(t, tc.testData.Department, p.Department)
				assert.Equal(t, tc.testData.RegularizeDate, p.RegularizeDate)
			}
		})
	}
}
