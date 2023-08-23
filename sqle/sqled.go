package sqled

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api"
	"github.com/actiontech/sqle/sqle/api/controller"

	// "github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/service"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/cluster"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/facebookgo/grace/gracenet"
)

func Run(config *config.Config) error {
	// init logger
	sqleCnf := config.Server.SqleCnf
	log.InitLogger(sqleCnf.LogPath, sqleCnf.LogMaxSizeMB, sqleCnf.LogMaxBackupNumber)
	defer log.ExitLogger()

	log.Logger().Infoln("starting sqled server")
	defer log.Logger().Info("stop sqled server")

	if sqleCnf.EnableClusterMode && sqleCnf.ServerId == "" {
		return fmt.Errorf("server id is required on cluster mode")
	}

	secretKey := sqleCnf.SecretKey
	if secretKey != "" {
		if err := utils.SetSecretKey([]byte(secretKey)); err != nil {
			return fmt.Errorf("set secret key error, %v, check your secret key in config file", err)
		}
	}

	defer driver.GetPluginManager().Stop()
	if err := driver.GetPluginManager().Start(sqleCnf.PluginPath, config.Server.PluginConfig); err != nil {
		return fmt.Errorf("init plugins error: %v", err)
	}

	// service.InitSQLQueryConfig(sqleCnf.SqleServerPort, sqleCnf.EnableHttps, config.Server.SQLQueryConfig)

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
		dbConfig.Host, dbConfig.Port, dbConfig.Schema, sqleCnf.DebugLog)
	if err != nil {
		return fmt.Errorf("get new storage failed: %v", err)
	}
	model.InitStorage(s)

	err = api.RegisterAsDMSTarget(config.Server.SqleCnf)
	if err != nil {
		return fmt.Errorf("register to dms failed :%v", err)
	}

	if sqleCnf.AutoMigrateTable {
		if err := s.AutoMigrate(); err != nil {
			return fmt.Errorf("auto migrate table failed: %v", err)
		}
		if err := s.CreateRulesIfNotExist(driver.GetPluginManager().GetAllRules()); err != nil {
			return fmt.Errorf("create rules failed while auto migrating table: %v", err)
		}
		namespaces, _, err := dmsobject.ListNamespaces(context.Background(), controller.GetDMSServerAddress(), dmsV1.ListNamespaceReq{
			PageSize:  999,
			PageIndex: 1,
		})
		if err != nil {
			return err
		}
		for _, namespace := range namespaces {
			if err := s.CreateDefaultTemplateIfNotExist(model.ProjectUID(namespace.NamespaceUid), driver.GetPluginManager().GetAllRules()); err != nil {
				return fmt.Errorf("create default template failed while auto migrating table: %v", err)
			}
		}
		// if err := s.CreateDefaultProject(); err != nil {
		// 	return fmt.Errorf("create default project failed while auto migrating table: %v", err)
		// }
		// if err := s.CreateDefaultRole(); err != nil {
		// 	return fmt.Errorf("create default rule failed while auto migrating table: %v", err)
		// }
	}
	exitChan := make(chan struct{})
	server.InitSqled(exitChan)

	var node cluster.Node
	if sqleCnf.EnableClusterMode {
		cluster.IsClusterMode = true
		log.Logger().Infoln("running sqled server on cluster mode")
		node = cluster.DefaultNode
		node.Join(sqleCnf.ServerId)
		defer node.Leave()
	} else {
		node = &cluster.NoClusterNode{}
	}

	jm := server.NewServerJobManger(node)
	jm.Start()
	defer jm.Stop()

	net := &gracenet.Net{}
	go api.StartApi(net, exitChan, sqleCnf)

	killChan := make(chan os.Signal, 1)
	// os.Kill is like kill -9 which kills a process immediately, can't be caught
	signal.Notify(killChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2 /*graceful-shutdown*/)
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
	return nil
}
