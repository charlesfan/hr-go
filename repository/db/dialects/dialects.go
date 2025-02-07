package dialects

type Dialect string

func (d Dialect) String() string {
	return string(d)
}

type ConnInfo struct {
	Dialect  Dialect
	Host     string
	Port     string
	User     string
	Database string
	Password string
}

const (
	MySql  Dialect = "mysql"
	Sqlite Dialect = "sqlite"
)
