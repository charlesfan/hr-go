package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

type Database struct {
	Dialect  string `development:"DATABASE_DIALECT" production:"DATABASE_DIALECT"`
	Host     string `development:"DATABASE_HOST" production:"DATABASE_HOST"`
	Port     string `development:"DATABASE_PORT" production:"DATAB_PORT"`
	Database string `development:"DATABASE_DATABASE" production:"DATABASE_DATABASE"`
	User     string `development:"DATABASE_USER" production:"DATABASE_USER"`
	Password string `development:"DATABASE_PASSWORD" production:"DATABASE_PASSWORD"`
	Schema   string `development:"DATABASE_SCHEMA" production:"DATABASE_SCHEMA"`
	Level    string `development:"DATABASE_LEVEL" production:"DATABASE_LEVEl"`
}

func (d *Database) GetUser() string {
	return d.User
}

func (d *Database) GetPassword() string {
	return d.Password
}

func (d *Database) GetDB() string {
	return d.Database
}

func (d *Database) GetSchema() string {
	return d.Schema
}

func (d *Database) GetHost() string {
	return d.Host
}

func (d *Database) GetPort() int {
	i, err := strconv.Atoi(d.Port)
	if err != nil {
		panic(fmt.Errorf("Port error: %v\n", err))
	}
	return i
}

func (d *Database) GetDebug() bool {
	return viper.GetString("ENV") == Dev.String()
}

func (d *Database) GetLogLevel() logger.LogLevel {
	logLevel := logger.Info
	switch strings.ToLower(d.Level) {
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	}
	return logLevel
}
