package config

type MySQLConfig struct {
	User string `env: "MSSQL_USER"`
	Pass string `env: "MSSQL_PASS"`
	Host string `env: "MSSQL_HOST"`
	Port string `env: "MSSQL_PORT" envDefault:"3306"`
}
