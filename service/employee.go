package service

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/charlesfan/hr-go/repository/db"
	"github.com/charlesfan/hr-go/repository/db/models/employee"
	"github.com/charlesfan/hr-go/utils/log"
)

type Employee struct {
	Id       string
	Name     string
	Email    string
	password string
}

type employeeService struct {
	employeeRepo employee.Repository
	authService  AuthenticationServicer
}

func (s *employeeService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *employeeService) LoginByEmailPassword(email, password string) (string, int64, int) {
	cnd := db.NewCnd().Eq(employee.EmailColumn, email)
	users, _, err := s.employeeRepo.Find(cnd)
	if err != nil {
		log.Error(err)
		return "", 0, ErrorCodeDatabaseFail
	}

	if len(users) <= 0 {
		log.Errorf("email: %s not found", email)
		return "", 0, ErrorCodeDataGetFail
	}

	u := users[0]
	if err := s.checkPasswordHash(password, u.Password); !err {
		return "", 0, ErrorCodeForbidden
	}

	tokenStr, expiredAt, errCode := s.authService.CreateToken(JwtConfig{
		Name: u.Name,
	})
	if errCode != ErrorCodeSuccess {
		return "", 0, ErrorCodeTokenCreateFail
	}

	return tokenStr, expiredAt, ErrorCodeSuccess
}

func NewEmployeeService(r employee.Repository, a AuthenticationServicer) EmployeeServicer {
	return &employeeService{
		employeeRepo: r,
		authService:  a,
	}
}
