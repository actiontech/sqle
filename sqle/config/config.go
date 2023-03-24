package config

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	SqleCnf        SqleConfig     `yaml:"sqle_config"`
	DBCnf          DatabaseConfig `yaml:"db_config"`
	SQLQueryConfig SQLQueryConfig `yaml:"sql_query_config"`
	PluginConfig   []PluginConfig `yaml:"plugin_config"`
}

type SqleConfig struct {
	ServerId           string `yaml:"server_id"`
	EnableClusterMode  bool   `yaml:"enable_cluster_mode"`
	SqleServerPort     int    `yaml:"server_port"`
	EnableHttps        bool   `yaml:"enable_https"`
	CertFilePath       string `yaml:"cert_file_path"`
	KeyFilePath        string `yaml:"key_file_path"`
	AutoMigrateTable   bool   `yaml:"auto_migrate_table"`
	DebugLog           bool   `yaml:"debug_log"`
	LogPath            string `yaml:"log_path"`
	LogMaxSizeMB       int    `yaml:"log_max_size_mb"`
	LogMaxBackupNumber int    `yaml:"log_max_backup_number"`
	PluginPath         string `yaml:"plugin_path"`
	SecretKey          string `yaml:"secret_key"`
}

type DatabaseConfig struct {
	MysqlCnf MysqlConfig `yaml:"mysql_cnf"`
}

type MysqlConfig struct {
	Host           string `yaml:"mysql_host"`
	Port           string `yaml:"mysql_port"`
	User           string `yaml:"mysql_user"`
	Password       string `yaml:"mysql_password,omitempty"`
	SecretPassword string `yaml:"secret_mysql_password,omitempty"`
	Schema         string `yaml:"mysql_schema"`
}

type SQLQueryConfig struct {
	EnableHttps              bool   `yaml:"enable_https"`
	CloudBeaverHost          string `yaml:"cloud_beaver_host"`
	CloudBeaverPort          string `yaml:"cloud_beaver_port"`
	CloudBeaverAdminUser     string `yaml:"cloud_beaver_admin_user"`
	CloudBeaverAdminPassword string `yaml:"cloud_beaver_admin_password"`
}

type PluginConfig struct {
	PluginName string `yaml:"plugin_name"`
	CMD        string `yaml:"cmd"`
}
