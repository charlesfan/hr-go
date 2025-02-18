package db

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/charlesfan/hr-go/repository/db/dialects"
)

type ISqlDB interface {
	IsConnected() bool
	Close() bool
	New(d string) error
}

type sqlDB struct {
	engine *Engine
}

func (w *sqlDB) IsConnected() bool {
	if w.engine.gormDB == nil {
		return false
	}
	return true
}

func (w *sqlDB) New(d string) error {
	var gdb *gorm.DB

	c := w.engine.sqlConfig

	switch dialects.Dialect(d) {
	case dialects.MySql:
		gdb = dialects.MySqlDB(c)
	case dialects.Sqlite:
		gdb = dialects.SqliteDB(c)
	default:
		return fmt.Errorf("Database not support: %q", d)
	}

	w.engine.gormDB = gdb
	return nil
}

func Initsql(gdb *gorm.DB) error {
	if InitSqlStr != "" {
		err := gdb.Exec(InitSqlStr).Error
		if err != nil {
			return err
		}
	}

	return nil
}
