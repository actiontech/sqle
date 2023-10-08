package config

import dmsCommonConf "github.com/actiontech/dms/pkg/dms-common/conf"

type Options struct {
	SqleOptions SqleOptions `yaml:"sqle"`
}
type SqleOptions struct {
	dmsCommonConf.BaseOptions `yaml:",inline"`
	DMSServerAddress          string     `yaml:"dms_server_address"`
	Service                   SeviceOpts `yaml:"service"`
}

type SeviceOpts struct {
	EnableClusterMode  bool           `yaml:"enable_cluster_mode"`
	AutoMigrateTable   bool           `yaml:"auto_migrate_table"`
	DebugLog           bool           `yaml:"debug_log"`
	LogPath            string         `yaml:"log_path"`
	LogMaxSizeMB       int            `yaml:"log_max_size_mb"`
	LogMaxBackupNumber int            `yaml:"log_max_backup_number"`
	PluginPath         string         `yaml:"plugin_path"`
	Database           Database       `yaml:"database"`
	PluginConfig       []PluginConfig `yaml:"plugin_config"`
}

type Database struct {
	Host           string `yaml:"mysql_host"`
	Port           string `yaml:"mysql_port"`
	User           string `yaml:"mysql_user"`
	Password       string `yaml:"mysql_password,omitempty"`
	SecretPassword string `yaml:"secret_mysql_password,omitempty"`
	Schema         string `yaml:"mysql_schema"`
}

type PluginConfig struct {
	PluginName string `yaml:"plugin_name"`
	CMD        string `yaml:"cmd"`
}
