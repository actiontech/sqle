package sqled

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"actiontech.cloud/sqle/sqle/sqle/api"
	"actiontech.cloud/sqle/sqle/sqle/config"
	"actiontech.cloud/sqle/sqle/sqle/driver"
	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"actiontech.cloud/sqle/sqle/sqle/server"
	"actiontech.cloud/sqle/sqle/sqle/server/auditplan"

	"github.com/facebookgo/grace/gracenet"
)

func Run(config *config.Config) error {
	// init logger
	log.InitLogger(config.Server.SqleCnf.LogPath)
	defer log.ExitLogger()

	log.Logger().Infoln("starting sqled server")

	if err := driver.InitPlugins(config.Server.SqleCnf.PluginPath); err != nil {
		return fmt.Errorf("init plugins error: %v", err)
	}

	dbConfig := config.Server.DBCnf.MysqlCnf
	s, err := model.NewStorage(dbConfig.User, dbConfig.Password,
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
			return fmt.Errorf("create default workflow template failed while auto migrateing table: %v", err)
		}
	}

	exitChan := make(chan struct{}, 0)
	server.InitSqled(exitChan)
	auditPlanMgrQuitCh := auditplan.InitManager(model.GetStorage())

	net := &gracenet.Net{}
	go api.StartApi(net, exitChan, config.Server.SqleCnf)

	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR2 /*graceful-shutdown*/)
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
