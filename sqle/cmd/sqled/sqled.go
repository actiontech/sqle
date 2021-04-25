package main

import (
	"actiontech.cloud/sqle/sqle/sqle"
	"actiontech.cloud/sqle/sqle/sqle/config"
	"actiontech.cloud/sqle/sqle/sqle/utils"
	"fmt"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var version string
var port int
var user string
var mysqlUser string
var mysqlPass string
var mysqlHost string
var mysqlPort string
var mysqlSchema string
var configPath string
var pidFile string
var debug bool
var autoMigrateTable bool
var logPath = "./logs"
var sqlServerParserServerHost = "127.0.0.1"
var sqlServerParserServerPort = "10001"

var httpsEnable bool
var certFilePath string
var keyFilePath string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sqle",
		Short: "SQLe",
		Long:  "SQLe\n\nVersion:\n  " + version,
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args); nil != err {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 10000, "http server port")
	rootCmd.PersistentFlags().StringVarP(&mysqlUser, "mysql-user", "", "sqle", "mysql user")
	rootCmd.PersistentFlags().StringVarP(&mysqlPass, "mysql-password", "", "sqle", "mysql password")
	rootCmd.PersistentFlags().StringVarP(&mysqlHost, "mysql-host", "", "localhost", "mysql host")
	rootCmd.PersistentFlags().StringVarP(&mysqlPort, "mysql-port", "", "3306", "mysql port")
	rootCmd.PersistentFlags().StringVarP(&mysqlSchema, "mysql-schema", "", "sqle", "mysql schema")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&pidFile, "pidfile", "", "", "pid file path")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "debug mode, print more log")
	rootCmd.PersistentFlags().BoolVarP(&autoMigrateTable, "auto-migrate-table", "", false, "auto migrate table if table model has changed")
	rootCmd.PersistentFlags().BoolVarP(&httpsEnable, "enable-https", "", false, "enable https")
	rootCmd.PersistentFlags().StringVarP(&certFilePath, "cert-file-path", "", "", "https cert file path")
	rootCmd.PersistentFlags().StringVarP(&keyFilePath, "key-file-path", "", "", "https key file path")
	rootCmd.Execute()
}

func run(cmd *cobra.Command, _ []string) error {
	var cfg = &config.Config{}

	// read config from file first, then read from cmd args.
	if configPath != "" {
		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("load config path: %s failed error :%v", configPath, err)
		}
		err = yaml.Unmarshal(b, cfg)
		if err != nil {
			return fmt.Errorf("unmarshal config file error %v", err)
		}
	} else {
		mysqlPass, err := utils.DecodeString(mysqlPass)
		if err != nil {
			return fmt.Errorf("decode mysql password to string error : %v", err)
		}
		cfg = &config.Config{
			Server: config.Server{
				SqleCnf: config.SqleConfig{
					SqleServerPort:   port,
					AutoMigrateTable: autoMigrateTable,
					DebugLog:         debug,
					LogPath:          logPath,
					EnableHttps:      httpsEnable,
					CertFilePath:     certFilePath,
					KeyFilePath:      keyFilePath,
				},
				DBCnf: config.DatabaseConfig{
					MysqlCnf: config.MysqlConfig{
						Host:     mysqlHost,
						Port:     mysqlPort,
						User:     mysqlUser,
						Password: mysqlPass,
						Schema:   mysqlSchema,
					},
				},
			},
		}
	}
	if pidFile != "" {
		f, err := os.Create(pidFile)
		if err != nil {
			return err
		}
		fmt.Fprintf(f, "%d\n", os.Getpid())
		f.Close()
		defer func() {
			os.Remove(pidFile)
		}()
	}
	return sqled.Run(cfg)
}
