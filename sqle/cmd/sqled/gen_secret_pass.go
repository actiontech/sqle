package main

import (
	"fmt"
	"io/ioutil"
	"os"

	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func genSecretPasswordCmd() *cobra.Command {
	run := func() error {
		var cfg = &config.SqleOptions{}
		if configPath != "" {
			b, err := ioutil.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("load config path: %s failed error :%v", configPath, err)
			}
			err = yaml.Unmarshal(b, cfg)
			if err != nil {
				return fmt.Errorf("unmarshal config file error %v", err)
			}

			secretKey := cfg.SecretKey
			if secretKey != "" {
				if err = dmsCommonAes.ResetAesSecretKey(secretKey); err != nil {
					return fmt.Errorf("set secret key error, %v, check your secret key in config file", err)
				}
			}

			password := cfg.Service.Database.Password
			if password == "" {
				return fmt.Errorf("mysql_password is empty")
			}
			secretPassword, err := dmsCommonAes.AesEncrypt(password)
			if err != nil {
				return fmt.Errorf("gen secret password error, %d", err)
			}
			cfg.Service.Database.SecretPassword = secretPassword
			cfg.Service.Database.Password = ""
		} else {
			return fmt.Errorf("--config is required")
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshal sqle config error %v", err)
		}
		err = ioutil.WriteFile(configPath, data, 0666)
		if err != nil {
			return fmt.Errorf("%v write sqle config file error %v", configPath, err)
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:   "gen_secret_pass",
		Short: "generate secret mysql password in sqled config file",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(); nil != err {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVarP(&configPath, "config", "", "", "config file path")
	return cmd
}
