package test

import (
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/repository/db/dialects"
	"github.com/charlesfan/hr-go/utils/tmpfile"
)

type Env struct {
	T     *testing.T
	B     *testing.B
	DB    *gorm.DB
	dbf   *os.File
	fs    []*os.File
	Debug bool
}

func (e *Env) Close() {
	sqldb, _ := e.DB.DB()
	sqldb.Close()
	for _, f := range e.fs {
		if f != nil {
			f.Close()
			os.Remove(f.Name())
		}
	}
}

func (e *Env) checkDebug() {
	e.Debug = false
	if os.Getenv("DEBUG") == "true" {
		e.Debug = true
	}
}

func (e *Env) setupDB() error {
	f, err := tempfile.TempFileWithSuffix(os.TempDir(), "gorm", ".db")
	if f == nil || err != nil {
		return err
	}
	e.fs = append(e.fs, f)

	db, err := gorm.Open(sqlite.Open(f.Name()), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		e.T.Errorf("gorm open fail => %v", err)
		return err
	}
	sqldb, err := db.DB()
	if err != nil {
		return err
	}
	sqldb.SetMaxIdleConns(10)

	e.DB = db
	e.dbf = f

	return nil
}

func (e *Env) DBConfig() config.Config {
	c := config.Config{
		Env: "development",
	}

	if e.dbf != nil {
		c.Database = &config.Database{
			Dialect: dialects.Sqlite.String(),
			Host:    e.dbf.Name(),
		}
	}

	return c
}

func (e *Env) runMigration() {
	values := []interface{}{}

	for _, value := range values {
		e.DB.Migrator().DropTable(value)
	}

	if err := e.DB.AutoMigrate(values...).Error; err != nil {
		panic(err)
	}
}

func SetupEnv(t *testing.T) *Env {
	t.Log("prepare database sqlite3")

	env := new(Env)
	env.T = t

	env.checkDebug()
	env.setupDB()
	//env.runMigration()

	return env
}
