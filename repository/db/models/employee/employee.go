package employee

import (
	"time"

	"gorm.io/gorm"

	"github.com/charlesfan/hr-go/repository/db"
	"github.com/charlesfan/hr-go/repository/db/models"
)

const (
	NameColumn       = "name"
	EmailColumn      = "email"
	DepartmentColumn = "department"
)

type Employee struct {
	ID             string `gorm:"column:id;unique;primary_key"`
	Name           string `gorm:"column:name;not null"`
	Email          string `gorm:"column:email;unique;not null"`
	Password       string `gorm:"column:password;not null"`
	Department     string `gorm:"column:department;not null"`
	RegularizeDate int64  `gorm:"column:regularize_date;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeleteAt       time.Time
}

func (Employee) TableName() string {
	return "hr_employee"
}

func (f *Employee) PK() models.PrimaryKey[string] {
	m := make(models.PrimaryKey[string])
	m["id"] = f.ID
	return m
}

func (f *Employee) SetPK(pk models.PrimaryKey[string]) *Employee {
	f.ID = pk["id"]
	return f
}

type Repository interface {
	Get(models.PrimaryKey[string]) (Employee, error)
	Create(*Employee) error
	Delete(models.PrimaryKey[string]) (int64, error)
	Find(*db.Cnd) ([]Employee, *db.Paging, error)
	Query(query interface{}, args ...interface{}) *gorm.DB
}
