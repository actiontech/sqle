package main

import (
	"fmt"

	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/spf13/cobra"
)

func backfillSQLManageFirstAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backfill-sql-manage-first-audit",
		Short: "Backfill SQL manage first audit result from latest audit result",
		RunE: func(cmd *cobra.Command, args []string) error {
			storage, err := newStorageFromCommandConfig()
			if err != nil {
				return err
			}
			model.InitStorage(storage)

			affectedRows, err := storage.BackfillSQLManageFirstAuditResult()
			if err != nil {
				return err
			}
			fmt.Printf("backfilled sql_manage_records first audit result rows: %d\n", affectedRows)
			return nil
		},
	}
	cmd.Flags().StringVarP(&configPath, "config", "", "", "config file path")
	cmd.Flags().StringVarP(&mysqlUser, "mysql-user", "", "sqle", "mysql user")
	cmd.Flags().StringVarP(&mysqlPass, "mysql-password", "", "sqle", "mysql password")
	cmd.Flags().StringVarP(&mysqlHost, "mysql-host", "", "localhost", "mysql host")
	cmd.Flags().StringVarP(&mysqlPort, "mysql-port", "", "3306", "mysql port")
	cmd.Flags().StringVarP(&mysqlSchema, "mysql-schema", "", "sqle", "mysql schema")
	cmd.Flags().BoolVarP(&debug, "debug", "", false, "debug mode, print more log")
	return cmd
}

func newStorageFromCommandConfig() (*model.Storage, error) {
	if configPath != "" {
		config.ParseConfigFile(configPath)
		dbConfig := config.GetOptions().SqleOptions.Service.Database
		dbPassword := dbConfig.Password
		if dbConfig.SecretPassword != "" {
			password, err := dmsCommonAes.AesDecrypt(dbConfig.SecretPassword)
			if err != nil {
				return nil, fmt.Errorf("read db info from config file error, %d", err)
			}
			dbPassword = password
		}
		return model.NewStorage(dbConfig.User, dbPassword, dbConfig.Host, dbConfig.Port, dbConfig.Schema, debug)
	}

	plainPassword, err := utils.DecodeString(mysqlPass)
	if err != nil {
		return nil, fmt.Errorf("decode mysql password to string error : %v", err)
	}
	return model.NewStorage(mysqlUser, plainPassword, mysqlHost, mysqlPort, mysqlSchema, debug)
}
