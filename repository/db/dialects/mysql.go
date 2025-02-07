package dialects

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/utils/log"
)

func MySqlDB(c *config.Database) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password,
		c.Host, c.Port,
		c.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("failed to connect database, %+v", err)
		return nil
	}

	return db
}
