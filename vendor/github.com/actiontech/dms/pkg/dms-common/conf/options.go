package conf

import "fmt"

type BaseOptions struct {
	ID                int64          `yaml:"id" validate:"required"`
	APIServiceOpts    *APIServerOpts `yaml:"api"`
	SecretKey         string         `yaml:"secret_key"`
	ServerId          string         `yaml:"server_id"`
	ReportHost        string         `yaml:"report_host"` //the host name or IP address of the cluster node
	EnableClusterMode bool           `yaml:"enable_cluster_mode"`
}

type APIServerOpts struct {
	Addr         string `yaml:"addr" validate:"required"`
	Port         int    `yaml:"port" validate:"required"`
	EnableHttps  bool   `yaml:"enable_https"`
	CertFilePath string `yaml:"cert_file_path"`
	KeyFilePath  string `yaml:"key_file_path"`
}

func (o *BaseOptions) GetAPIServer() *APIServerOpts {
	return o.APIServiceOpts
}

func (api *APIServerOpts) GetHTTPAddr() string {
	return fmt.Sprintf("%v:%v", api.Addr, api.Port)
}
