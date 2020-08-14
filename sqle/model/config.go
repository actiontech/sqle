

package model

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	SqleCnf SqleConfig     `yaml:"sqle_config"`
	DBCnf   DatabaseConfig `yaml:"db_config"`
}

type SqleConfig struct {
	SqleServerPort   int    `yaml:"server_port"`
	AutoMigrateTable bool   `yaml:"auto_migrate_table"`
	DebugLog         bool   `yaml:"debug_log"`
	LogPath          string `yaml:"log_path"`
}

type DatabaseConfig struct {
	MysqlCnf     MysqlConfig     `yaml:"mysql_cnf"`
	SqlServerCnf SqlServerConfig `yaml:"sql_server_cnf"`
}

type MysqlConfig struct {
	Host     string `yaml:"mysql_host"`
	Port     string `yaml:"mysql_port"`
	User     string `yaml:"mysql_user"`
	Password string `yaml:"mysql_password"`
	Schema   string `yaml:"mysql_schema"`
}

type SqlServerConfig struct {
	Host string `yaml:"sql_server_host"`
	Port string `yaml:"sql_server_port"`
}
