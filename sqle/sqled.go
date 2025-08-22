package sqled

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"github.com/actiontech/dms/pkg/dms-common/pkg/http"
	"github.com/actiontech/sqle/sqle/api"
	"github.com/actiontech/sqle/sqle/dms"
	knowledge_base "github.com/actiontech/sqle/sqle/server/knowledge_base"
	optimizationRule "github.com/actiontech/sqle/sqle/server/optimization/rule"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/server/cluster"

	"github.com/facebookgo/grace/gracenet"
)

//go:embed docs/swagger.yaml
var sqleSwaggerYaml []byte

func Run(options *config.SqleOptions) error {
	// init logger
	sqleCnf := options.Service
	log.InitLogger(sqleCnf.LogPath, sqleCnf.LogMaxSizeMB, sqleCnf.LogMaxBackupNumber, sqleCnf.DebugLog)
	defer log.ExitLogger()

	log.Logger().Infoln("starting sqled server")
	defer log.Logger().Info("stop sqled server")

	// validate config
	err := validateConfig(options)
	if err != nil {
		return err
	}

	secretKey := options.SecretKey
	if secretKey != "" {
		// reset jwt singing key, default dms token
		if err := http.ResetJWTSigningKeyAndDefaultToken(secretKey); err != nil {
			return err
		}

		// reset aes secret key
		if err := dmsCommonAes.ResetAesSecretKey(secretKey); err != nil {
			return err
		}
	}

	// nofify singal
	exitChan := make(chan struct{})
	net := &gracenet.Net{}
	go NotifySignal(exitChan, net)

	// init plugins
	{
		defer driver.GetPluginManager().Stop()
		if err := driver.GetPluginManager().Start(sqleCnf.PluginPath, options.Service.PluginConfig); err != nil {
			return fmt.Errorf("init plugins error: %v", err)
		}
	}

	// init storage
	{
		dbConfig := options.Service.Database

		dbPassword := dbConfig.Password
		// Support using secret mysql password in sqled config, read secret_mysql_password first,
		// but you can still use mysql_password to be compatible with older versions.
		secretPassword := dbConfig.SecretPassword
		if secretPassword != "" {
			password, err := dmsCommonAes.AesDecrypt(secretPassword)
			if err != nil {
				return fmt.Errorf("read db info from config file error, %v", err)
			}
			dbPassword = password
		}

		s, err := model.NewStorage(dbConfig.User, dbPassword,
			dbConfig.Host, dbConfig.Port, dbConfig.Schema, sqleCnf.DebugLog)
		if err != nil {
			return fmt.Errorf("get new storage failed: %v", err)
		}
		model.InitStorage(s)

		if sqleCnf.AutoMigrateTable {
			if err := s.AutoMigrate(); err != nil {
				return fmt.Errorf("auto migrate table failed: %v", err)
			}
			err := s.CreateDefaultWorkflowTemplateIfNotExist()
			if err != nil {
				return fmt.Errorf("create workflow template failed: %v", err)
			}
			if err := s.CreateRuleCategoriesRelated(); err != nil {
				return fmt.Errorf("create rule categories related failed while auto migrating table: %v", err)
			}
			rules := model.MergeOptimizationRules(driver.GetPluginManager().GetAllRules(), optimizationRule.OptimizationRuleMap)
			if err := s.CreateRulesIfNotExist(rules); err != nil {
				return fmt.Errorf("create rules failed while auto migrating table: %v", err)
			}
			if err := s.DeleteRulesIfNotExist(rules); err != nil {
				return fmt.Errorf("delete rules failed while auto migrating table: %v", err)
			}

			if err := s.CreateDefaultTemplateIfNotExist(model.ProjectIdForGlobalRuleTemplate, driver.GetPluginManager().GetAllRules()); err != nil {
				return fmt.Errorf("create default template failed while auto migrating table: %v", err)
			}
			if err := knowledge_base.LoadKnowledge(rules); err != nil {
				return fmt.Errorf("LoadKnowledge failed : %v", err)
			}
		}
		{
			if err := s.CreateDefaultReportPushConfigIfNotExist(model.DefaultProjectUid); err != nil {
				return fmt.Errorf("create default report push config failed: %v", err)
			}
		}
	}

	err = dms.RegisterAsDMSTarget(options)
	if err != nil {
		return fmt.Errorf("register to dms failed :%v", err)
	}

	if options.OptimizationConfig.OptimizationKey != "" && options.OptimizationConfig.OptimizationURL != "" {
		// todo flash del optimize rules
		//optimizationRule.InitOptimizationRule()
	}

	server.InitSqled(exitChan)

	var node cluster.Node
	if sqleCnf.EnableClusterMode {
		cluster.IsClusterMode = true
		log.Logger().Infoln("running sqled server on cluster mode")
		node = cluster.DefaultNode
		node.Join(fmt.Sprintf("%v", options.ID))
		defer node.Leave()
	} else {
		node = &cluster.NoClusterNode{}
	}

	jm := server.NewServerJobManger(node)
	jm.Start()
	defer jm.Stop()

	go api.StartApi(net, exitChan, options, sqleSwaggerYaml)

	// Wait for exit signal from NotifySignal goroutine
	<-exitChan
	log.Logger().Infoln("sqled server will exit")
	return nil
}

func validateConfig(options *config.SqleOptions) error {
	sqleCnf := options.Service
	if sqleCnf.EnableClusterMode {
		if options.ID == 0 {
			return fmt.Errorf("server id is required on cluster mode")
		}
		if options.ReportHost == "" {
			return fmt.Errorf("report host is required on cluster mode")
		}
	}
	return nil
}

func NotifySignal(exitChan chan struct{}, net *gracenet.Net) {
	killChan := make(chan os.Signal, 1)
	// os.Kill is like kill -9 which kills a process immediately, can't be caught
	signal.Notify(killChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2 /*graceful-shutdown*/)
	sig := <-killChan
	switch sig {
	case syscall.SIGUSR2:
		if pid, err := net.StartProcess(); nil != err {
			log.Logger().Infof("Graceful restarted by signal SIGUSR2, but failed: %v", err)
		} else {
			log.Logger().Infof("Graceful restarted, new pid is %v", pid)
		}
		log.Logger().Infof("old sqled exit")
	default:
		log.Logger().Infof("Exit by signal %v", sig)
	}

	close(exitChan)
}
