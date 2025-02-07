package dialects

import (
	glog "log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/utils/log"
)

func SqliteDB(c *config.Database) *gorm.DB {
	newLogger := logger.New(
		glog.New(os.Stdout, "\r\n", glog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open(c.Host), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Errorf("gorm open fail => %+v", err)
		return nil
	}

	return db
}
