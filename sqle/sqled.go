package sqled

import (
	"actiontech.cloud/universe/sqle/v4/sqle/api"
	"actiontech.cloud/universe/sqle/v4/sqle/config"
	"actiontech.cloud/universe/sqle/v4/sqle/inspector"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/sqle/v4/sqle/model"
	"actiontech.cloud/universe/sqle/v4/sqle/server"
	"actiontech.cloud/universe/ucommon/v4/ubootstrap"
	"fmt"
	"github.com/facebookgo/grace/gracenet"
	"syscall"
)

func Run(config *config.Config) error {
	// init logger
	log.InitLogger(config.Server.SqleCnf.LogPath)
	defer log.ExitLogger()

	log.Logger().Infoln("starting sqled server")

	err := inspector.LoadPtTemplateFromFile("./scripts/pt-online-schema-change.template")
	if err != nil {
		return fmt.Errorf("load './scripts/pt-online-schema-change.template/' failed: %v", err)
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
		if err := s.CreateRulesIfNotExist(inspector.InitRules); err != nil {
			return fmt.Errorf("create rules failed while auto migrating table: %v", err)
		}
		if err := s.CreateDefaultTemplate(inspector.DefaultTemplateRules); err != nil {
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

	net := &gracenet.Net{}
	go api.StartApi(net, exitChan, config.Server.SqleCnf)

	killChan := ubootstrap.ListenKillSignal()
	select {
	case <-exitChan:
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
