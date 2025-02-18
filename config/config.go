package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/charlesfan/hr-go/utils/log"
)

type (
	EnvType  string
	LogLevel string
)

var content Config

const (
	Dev EnvType = "development"
	Pro EnvType = "production"

	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
)

func (e EnvType) String() string { return string(e) }

func (l LogLevel) String() string { return string(l) }

type Config struct {
	Env      string
	Server   *Server
	Database *Database
	Log      *Log
	Redis    *Redis
}

func NewConfig() Config {
	return content
}

func Init() {
	// init config, default value
	c := Config{
		Env: "development",
		Log: &Log{
			Level: DebugLevel.String(),
		},
		Server: &Server{
			Schema:     "http",
			Host:       "0.0.0.0",
			Port:       "8080",
			MaxWorkers: 128,
			AppSecret:  "FFPGR3tAh1WIT1cMYPxgIlnd6CbPlc0b",
		},
		Database: &Database{
			Dialect:  "mysql",
			Host:     "db",
			Port:     "3306",
			Database: "hr_go",
			User:     "root",
			Password: "dev123",
			Schema:   "user",
		},
		Redis: &Redis{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	}

	if err := c.loadConfig(); err != nil {
		log.Error("load config failed")
		os.Exit(1)
	}

	content = c
}

func (c *Config) loadConfig() error {
	cfgPath := viper.GetString("config")
	viper.SetEnvPrefix("HR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file (%s), %v", cfgPath, err)
			os.Exit(1)
		}
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		return err
	}

	b, _ := json.MarshalIndent(&c, "", "  ")
	fmt.Println(string(b))

	return nil
}
