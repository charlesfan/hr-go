package fixture

import "github.com/charlesfan/hr-go/repository/db/models/employee"

func GetPrepareEmployee() *employee.Employee {
	return &employee.Employee{
		ID:             "aeb8ada1-d8f9-46c6-8bf2-b1fbbc5530d7",
		Name:           "employee01",
		Email:          "em@test.com",
		Department:     "hr",
		RegularizeDate: 1739253769000,
	}
}
