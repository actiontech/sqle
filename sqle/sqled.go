package sqled

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/actiontech/sqle/sqle/utils"

	"github.com/actiontech/sqle/sqle/api"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/auditplan"

	"github.com/facebookgo/grace/gracenet"
)

func Run(config *config.Config) error {
	// init logger
	log.InitLogger(config.Server.SqleCnf.LogPath)
	defer log.ExitLogger()

	log.Logger().Infoln("starting sqled server")

	secretKey := config.Server.SqleCnf.SecretKey
	if secretKey != "" {
		if err := utils.SetSecretKey([]byte(secretKey)); err != nil {
			return fmt.Errorf("set secret key error, %v, check your secret key in config file", err)
		}
	}

	if err := driver.InitPlugins(config.Server.SqleCnf.PluginPath); err != nil {
		return fmt.Errorf("init plugins error: %v", err)
	}

	dbConfig := config.Server.DBCnf.MysqlCnf

	dbPassword := dbConfig.Password
	// Support using secret mysql password in sqled config, read secret_mysql_password first,
	// but you can still use mysql_password to be compatible with older versions.
	secretPassword := dbConfig.SecretPassword
	if secretPassword != "" {
		password, err := utils.AesDecrypt(secretPassword)
		if err != nil {
			return fmt.Errorf("read db info from config file error, %d", err)
		}
		dbPassword = password
	}

	s, err := model.NewStorage(dbConfig.User, dbPassword,
		dbConfig.Host, dbConfig.Port, dbConfig.Schema, config.Server.SqleCnf.DebugLog)
	if err != nil {
		return fmt.Errorf("get new storage failed: %v", err)
	}
	model.InitStorage(s)

	if config.Server.SqleCnf.AutoMigrateTable {
		if err := s.AutoMigrate(); err != nil {
			return fmt.Errorf("auto migrate table failed: %v", err)
		}
		if err := s.CreateRulesIfNotExist(driver.AllRules()); err != nil {
			return fmt.Errorf("create rules failed while auto migrating table: %v", err)
		}
		if err := s.CreateDefaultTemplate(driver.AllRules()); err != nil {
			return fmt.Errorf("create default template failed while auto migrating table: %v", err)
		}
		if err := s.CreateAdminUser(); err != nil {
			return fmt.Errorf("create default admin user failed while auto migrating table: %v", err)
		}
		if err := s.CreateDefaultWorkflowTemplate(); err != nil {
			return fmt.Errorf("create default workflow template failed while auto migrating table: %v", err)
		}
	}
	exitChan := make(chan struct{})
	server.InitSqled(exitChan)
	auditPlanMgrQuitCh := auditplan.InitManager(model.GetStorage())

	net := &gracenet.Net{}
	go api.StartApi(net, exitChan, config.Server.SqleCnf)

	killChan := make(chan os.Signal, 1)
	// os.Kill is like kill -9 which kills a process immediately, can't be caught
	signal.Notify(killChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2 /*graceful-shutdown*/)
	select {
	case <-exitChan:
		auditPlanMgrQuitCh <- struct{}{}
		log.Logger().Infoln("sqled server will exit")
	case sig := <-killChan:
		switch sig {
		case syscall.SIGUSR2:
			if pid, err := net.StartProcess(); nil != err {
				log.Logger().Infof("Graceful restarted by signal SIGUSR2, but failed: %v", err)
				return err
			} else {
				log.Logger().Infof("Graceful restarted, new pid is %v", pid)
			}
			log.Logger().Infof("old sqled exit")
		default:
			log.Logger().Infof("Exit by signal %v", sig)
		}
	}
	log.Logger().Info("stop sqled server")
	return nil
}
