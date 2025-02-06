package config

type Server struct {
	Schema     string `development:"SERVER_SCHEMA" production:"SERVER_SCHEMA"`
	Host       string `development:"SERVER_HOST" production:"SERVER_HOST"`
	Port       string `development:"SERVER_PORT" production:"SERVER_PORT"`
	MaxWorkers int    `development:"SERVER_MAX_WORKERS" production:"SERVER_MAX_WORKERS"`
	AppSecret  string `development:"SERVER_APP_SECRET" production:"SERVER_APP_SECRET"`
}
