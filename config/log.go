package config

type Log struct {
	File  string `development:"LOG_FILE" production:"LOG_FILE"`
	Level string `development:"LOG_LEVEL" production:"LOG_LEVEL"`
}
