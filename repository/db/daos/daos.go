package daos

import (
	"sync"

	"gorm.io/gorm"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/repository/db"
	"github.com/charlesfan/hr-go/repository/db/models/employee"
)

var (
	once          sync.Once
	dbEngine      *db.Engine
	dbrepoFactory *dbRepoFactory
)

type DBRepoFactory interface {
	EmployeeRepo() employee.Repository
}

func Init(c config.Config) error {
	dbEngine = db.New(c)
	dbEngine.Run()

	return nil
}

type onceRepo struct {
	once sync.Once
	repo interface{}
}

type dbRepoFactory struct {
	gormDB *gorm.DB

	employeeRepository onceRepo
}

func (s *dbRepoFactory) EmployeeRepo() employee.Repository {
	s.employeeRepository.once.Do(func() {
		dbEngine.Migration([]interface{}{&employee.Employee{}})
		s.employeeRepository.repo = newEmployeeRepo(s.gormDB)
	})

	return s.employeeRepository.repo.(employee.Repository)
}

func NewDBRepoFactory() DBRepoFactory {
	once.Do(func() {
		dbrepoFactory = &dbRepoFactory{
			gormDB: dbEngine.GormDB(),
		}
	})
	return dbrepoFactory
}
