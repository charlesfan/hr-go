package repo

import (
	"sync"

	"gorm.io/gorm"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/repository/db"
)

var (
	once          sync.Once
	dbEngine      *db.Engine
	dbrepoFactory *dbRepoFactory
)

type DBRepoFactory interface {
}

func Init(c config.Config) error {
	dbEngine = db.New(c)
	dbEngine.Run()

	return nil
}

type dbRepoFactory struct {
	gormDB *gorm.DB
}
