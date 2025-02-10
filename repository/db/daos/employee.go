package daos

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/charlesfan/hr-go/repository/db"
	"github.com/charlesfan/hr-go/repository/db/models"
	"github.com/charlesfan/hr-go/repository/db/models/employee"
	"github.com/charlesfan/hr-go/utils/log"
)

type employeeRepo struct {
	db *gorm.DB
}

func (r *employeeRepo) Create(f *employee.Employee) error {
	d := r.db.Create(f)
	if err := d.Error; err != nil {
		log.Error("employeeRepo Create fail => ", err)
		return db.ErrChecking(err)
	}
	return nil
}

func (r *employeeRepo) Get(pk models.PrimaryKey[string]) (employee.Employee, error) {
	f := employee.Employee{}

	if !pk.Validate() {
		err := errors.New("primary key is nil")
		log.Error("employeeRepo Get fail => ", err)
		return f, err
	}

	cnd := db.NewCnd()
	for k, v := range pk {
		cnd.Eq(k, v)
	}

	cnd.Find(r.db, &f)
	return f, nil
}

func (r *employeeRepo) Delete(pk models.PrimaryKey[string]) (int64, error) {
	data, err := r.Get(pk)
	if err != nil {
		return 0, err
	}

	if data.ID == "" {
		return 0, fmt.Errorf("employee not found: %s", data.ID)
	}

	result := r.db.Delete(data)
	row := result.RowsAffected
	if err = result.Error; err != nil {
		log.Error("employeeRepo delete fail => ", err)
	}
	return row, err
}

func (r *employeeRepo) Find(cnd *db.Cnd) (list []employee.Employee, paging *db.Paging, err error) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &employee.Employee{})
	if cnd.Paging != nil {
		paging = &db.Paging{
			Page:  cnd.Paging.Page,
			Limit: cnd.Paging.Limit,
			Total: count,
		}
	}
	return
}

func (r *employeeRepo) Query(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}

func newEmployeeRepo(db *gorm.DB) employee.Repository {
	return &employeeRepo{
		db: db,
	}
}
