package main

import (
	"fmt"
	"io/ioutil"
	"os"

	dmsCommonConf "github.com/actiontech/dms/pkg/dms-common/conf"
	sqled "github.com/actiontech/sqle/sqle"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var version string
var port int

// var user string
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
var logMaxSizeMB int
var logMaxBackupNumber int
var httpsEnable bool
var certFilePath string
var keyFilePath string
var pluginPath string

func init() {
	config.Version = version

}

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
	rootCmd.Flags().IntVarP(&port, "port", "p", 10000, "http server port")
	rootCmd.Flags().StringVarP(&mysqlUser, "mysql-user", "", "sqle", "mysql user")
	rootCmd.Flags().StringVarP(&mysqlPass, "mysql-password", "", "sqle", "mysql password")
	rootCmd.Flags().StringVarP(&mysqlHost, "mysql-host", "", "localhost", "mysql host")
	rootCmd.Flags().StringVarP(&mysqlPort, "mysql-port", "", "3306", "mysql port")
	rootCmd.Flags().StringVarP(&mysqlSchema, "mysql-schema", "", "sqle", "mysql schema")
	rootCmd.Flags().StringVarP(&configPath, "config", "", "", "config file path")
	rootCmd.Flags().StringVarP(&pidFile, "pidfile", "", "", "pid file path")
	rootCmd.Flags().BoolVarP(&debug, "debug", "", false, "debug mode, print more log")
	rootCmd.Flags().BoolVarP(&autoMigrateTable, "auto-migrate-table", "", false, "auto migrate table if table model has changed")
	rootCmd.Flags().IntVarP(&logMaxSizeMB, "log-max-size-mb", "", 1024, "log max size (MB)")
	rootCmd.Flags().IntVarP(&logMaxBackupNumber, "log-max-backup-number", "", 2, "log max backup number")
	rootCmd.Flags().BoolVarP(&httpsEnable, "enable-https", "", false, "enable https")
	rootCmd.Flags().StringVarP(&certFilePath, "cert-file-path", "", "", "https cert file path")
	rootCmd.Flags().StringVarP(&keyFilePath, "key-file-path", "", "", "https key file path")
	rootCmd.Flags().StringVarP(&pluginPath, "plugin-path", "", "", "plugin path")

	rootCmd.AddCommand(genSecretPasswordCmd())
	if err := rootCmd.Execute(); err != nil {
		log.NewEntry().Error("sqle abnormal termination:", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, _ []string) error {
	var cfg = &config.Options{}

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
		cfg = &config.Options{
			SqleOptions: config.SqleOptions{
				BaseOptions: dmsCommonConf.BaseOptions{
					APIServiceOpts: &dmsCommonConf.APIServerOpts{
						Port:         port,
						EnableHttps:  httpsEnable,
						CertFilePath: certFilePath,
						KeyFilePath:  keyFilePath,
					},
				},
				Service: config.SeviceOpts{
					AutoMigrateTable:   autoMigrateTable,
					DebugLog:           debug,
					LogPath:            logPath,
					LogMaxSizeMB:       logMaxSizeMB,
					LogMaxBackupNumber: logMaxBackupNumber,
					PluginPath:         pluginPath,
					Database: config.Database{
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
	return sqled.Run(&cfg.SqleOptions)
}
