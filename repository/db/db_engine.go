package db

import (
	"errors"

	"gorm.io/gorm"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/utils/log"
)

type Engine struct {
	sqlDB
	sqlConfig *config.Database
	gormDB    *gorm.DB
}

func (e *Engine) Run() {
	err := e.sqlDB.New(e.sqlConfig.Dialect)
	if err != nil {
		log.Error(errors.Unwrap(err))
	}
}

func (e *Engine) Migration(tables []interface{}) {
	if err := e.gormDB.AutoMigrate(tables...); err != nil {
		log.Error(err)
	}
}

func (e *Engine) GormDB() *gorm.DB {
	return e.gormDB
}

func New(c config.Config) *Engine {
	e := &Engine{
		sqlConfig: c.Database,
	}

	e.sqlDB.engine = e

	return e
}
