package driver

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/actiontech/sqle/sqle/driver/proto"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pkg/errors"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// DSN provide necessary information to connect to database.
type DSN struct {
	Host             string
	Port             string
	User             string
	Password         string
	AdditionalParams params.Params

	// DatabaseName is the default database to connect.
	DatabaseName string
}

// Config define the configuration for driver.
type Config struct {
	DSN   *DSN
	Rules []*Rule
}

type PluginClient struct {
	path string
	c    *goPlugin.Client
}

func TestPluginDriver(path string) bool {
	c := newClientFromFile(path)
	cp, err := c.Client()
	if err != nil {
		log.NewEntry().Errorf("test conn plugin failed: %v", err)
		return false
	}
	defer c.Kill()

	err = cp.Ping()
	if err != nil {
		log.NewEntry().Errorf("test conn plugin failed: %v", err)
		return false
	}
	return true
}

func RegisterDrivers(path string) (pluginName string, err error) {
	c := newClientFromFile(path)
	cp, err := c.Client()
	if err != nil {
		return "", err
	}
	defer c.Kill()

	pluginName, version, drvClient, err := registerAuditDriver(cp)
	if err != nil {
		return "", err
	}

	if err := registerQueryDriver(pluginName, cp); err != nil {
		log.Logger().WithFields(logrus.Fields{
			"plugin_name": pluginName,
			"plugin_type": PluginNameQueryDriver,
		}).Infof("plugin not exist or failed to load. err: %v", err)
	}

	if err := registerAnalysisDriver(pluginName, cp); err != nil {
		log.Logger().WithFields(logrus.Fields{
			"plugin_name": pluginName,
			"plugin_type": PluginNameAnalysisDriver,
		}).Infof("plugin not exist or failed to load. err: %v", err)
	}

	// to be compatible with old plugins
	// the old plugin will panic if it call close() here
	if version >= DefaultPluginVersion {
		_, err = drvClient.Close(context.TODO(), &proto.Empty{})
		if err != nil {
			log.Logger().Errorf("gracefully close plugins failed, will force kill the sub progress. err: %v", err)
		}
	}
	return pluginName, nil
}

func newClientFromFile(path string) *PluginClient {
	client := goPlugin.NewClient(&goPlugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		VersionedPlugins: defaultPluginSet,
		Cmd:              exec.Command(path),
		AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
		GRPCDialOptions:  SQLEGRPCDialOptions,
	})
	return &PluginClient{
		path: path,
		c:    client,
	}
}

func (p *PluginClient) Kill() {
	if p.c != nil {
		p.c.Kill()
	}
}

func (p *PluginClient) Client() (goPlugin.ClientProtocol, error) {
	return p.c.Client()
}

var SQLEGRPCDialOptions = []grpc.DialOption{}

var driverManagerMu = &sync.RWMutex{}
var driverManagers = make(map[string]driverManagerHandler)

type driverManagerHandler struct {
	newDriverManagerFunc newDriverManagerHandler
}

type newDriverManagerHandler func(log *logrus.Entry, dbType string, config *Config) (DriverManager, error)

func RegisterDriverFromClient(path string) error {
	pluginName, err := RegisterDrivers(path)
	if err != nil {
		return fmt.Errorf("register plugin failed: %v", err)
	}

	handler := func(log *logrus.Entry, dbType string, config *Config) (DriverManager, error) {
		client := newClientFromFile(path)
		gRPCClient, err := client.Client()
		if err != nil {
			return nil, err
		}
		closeCh := make(chan struct{})
		go func() {
			<-closeCh
			client.Kill()
		}()

		drvMgr := &PluginDriverManager{
			grpcClient:    gRPCClient,
			pluginCloseCh: closeCh,
			config:        config,
			dbType:        dbType,
			log:           log,
		}

		if err = drvMgr.initAuditDriver(); err != nil {
			return nil, err
		}
		if err = drvMgr.initSQLQueryDriver(); err != nil {
			return nil, err
		}
		if err = drvMgr.initAnalysisDriver(); err != nil {
			return nil, err
		}

		return drvMgr, nil
	}

	RegisterDriverManger(pluginName, handler)
	return nil
}

func RegisterDriverManger(pluginName string, handler newDriverManagerHandler) {
	driverManagerMu.RLock()
	_, exist := driverManagers[pluginName]
	driverManagerMu.RUnlock()
	if exist {
		panic(fmt.Sprintf("duplicated driver name [%v]", pluginName))
	}

	driverManagerMu.Lock()
	driverManagers[pluginName] = driverManagerHandler{
		newDriverManagerFunc: handler,
	}
	driverManagerMu.Unlock()
}

type DriverManager interface {
	GetAuditDriver() (Driver, error)
	GetSQLQueryDriver() (SQLQueryDriver, error)
	GetAnalysisDriver() (AnalysisDriver, error)
	// Close invoke grpc.Close of audit plugin to inform plugin process to recycle their resource
	// resource of all drivers should be recycle in this function
	Close(ctx context.Context)
}

type PluginDriverManager struct {
	grpcClient           goPlugin.ClientProtocol
	pluginCloseCh        chan struct{}
	dbType               string
	log                  *logrus.Entry
	config               *Config
	auditPluginClient    proto.DriverClient
	queryPluginClient    proto.QueryDriverClient
	analysisPluginClient proto.AnalysisDriverClient
}

func (d *PluginDriverManager) GetAuditDriver() (Driver, error) {
	if d.auditPluginClient == nil {
		return nil, fmt.Errorf("audit driver type %v is not supported", d.dbType)
	}
	return &driverImpl{d.auditPluginClient, d.pluginCloseCh}, nil
}

func (d *PluginDriverManager) initAuditDriver() error {
	_, exist := auditDrivers[d.dbType]
	if !exist {
		return nil
	}

	if d.auditPluginClient != nil {
		return nil
	}

	rawI, err := d.grpcClient.Dispense(PluginNameAuditDriver)
	if err != nil {
		return fmt.Errorf("dispense audit driver failed: %v", err)
	}
	// pluginInst can only be proto.QueryDriverClient
	//nolint:forcetypeassert
	pluginInst := rawI.(proto.DriverClient)

	// protoRules send to plugin for Audit.
	protoRules := make([]*proto.Rule, len(d.config.Rules))
	for i, rule := range d.config.Rules {
		protoRules[i] = convertRuleFromDriverToProto(rule)
	}

	initRequest := &proto.InitRequest{
		Rules: protoRules,
	}
	if d.config != nil && d.config.DSN != nil {
		initRequest.Dsn = &proto.DSN{
			Host:             d.config.DSN.Host,
			Port:             d.config.DSN.Port,
			User:             d.config.DSN.User,
			Password:         d.config.DSN.Password,
			AdditionalParams: proto.ConvertParamToProtoParam(d.config.DSN.AdditionalParams),

			// database is to open.
			Database: d.config.DSN.DatabaseName,
		}
	}

	_, err = pluginInst.Init(context.TODO(), initRequest)
	if err != nil {
		return fmt.Errorf("init audit driver failed: %v", err)
	}
	d.auditPluginClient = pluginInst
	return nil
}

func (d *PluginDriverManager) GetSQLQueryDriver() (SQLQueryDriver, error) {
	if d.queryPluginClient == nil {
		return nil, fmt.Errorf("SQL query driver type %v is not supported", d.dbType)
	}
	return &queryDriverImpl{d.queryPluginClient}, nil
}

func (d *PluginDriverManager) initSQLQueryDriver() error {
	_, exist := queryDrivers[d.dbType]
	if !exist {
		return nil
	}

	if d.queryPluginClient != nil {
		return nil
	}

	rawI, err := d.grpcClient.Dispense(PluginNameQueryDriver)
	if err != nil {
		return fmt.Errorf("dispense SQL query driver failed: %v", err)
	}
	// pluginInst can only be proto.QueryDriverClient
	//nolint:forcetypeassert
	pluginInst := rawI.(proto.QueryDriverClient)

	initRequest := &proto.InitRequest{
		Rules: []*proto.Rule{},
	}
	if d.config != nil && d.config.DSN != nil {
		initRequest.Dsn = &proto.DSN{
			Host:             d.config.DSN.Host,
			Port:             d.config.DSN.Port,
			User:             d.config.DSN.User,
			Password:         d.config.DSN.Password,
			AdditionalParams: proto.ConvertParamToProtoParam(d.config.DSN.AdditionalParams),

			// database is to open.
			Database: d.config.DSN.DatabaseName,
		}
	}
	_, err = pluginInst.Init(context.TODO(), initRequest)
	if err != nil {
		return fmt.Errorf("init SQL query driver failed: %v", err)
	}
	d.queryPluginClient = pluginInst
	return nil
}

func (d *PluginDriverManager) GetAnalysisDriver() (AnalysisDriver, error) {
	if d.analysisPluginClient == nil {
		return nil, fmt.Errorf("analysis driver type %v is not supported", d.dbType)
	}
	return &analysisDriverImpl{d.analysisPluginClient}, nil
}

func (d *PluginDriverManager) initAnalysisDriver() error {
	_, exist := analysisDrivers[d.dbType]
	if !exist {
		return nil
	}

	if d.analysisPluginClient != nil {
		return nil
	}

	rawI, err := d.grpcClient.Dispense(PluginNameAnalysisDriver)
	if err != nil {
		return fmt.Errorf("dispense analysis driver failed: %v", err)
	}
	//nolint:forcetypeassert
	pluginInst := rawI.(proto.AnalysisDriverClient)

	initRequest := &proto.AnalysisDriverInitRequest{}
	if d.config != nil && d.config.DSN != nil {
		initRequest.Dsn = &proto.DSN{
			Host:             d.config.DSN.Host,
			Port:             d.config.DSN.Port,
			User:             d.config.DSN.User,
			Password:         d.config.DSN.Password,
			AdditionalParams: proto.ConvertParamToProtoParam(d.config.DSN.AdditionalParams),

			// database is to open.
			Database: d.config.DSN.DatabaseName,
		}
	}
	_, err = pluginInst.Init(context.TODO(), initRequest)
	if err != nil {
		return fmt.Errorf("init analysis driver failed: %v", err)
	}
	d.analysisPluginClient = pluginInst
	return nil
}

func (d *PluginDriverManager) Close(ctx context.Context) {
	impl := &driverImpl{d.auditPluginClient, d.pluginCloseCh}
	impl.Close(ctx)
}

func NewDriverManger(log *logrus.Entry, dbType string, config *Config) (DriverManager, error) {
	driverManagerMu.RLock()
	defer driverManagerMu.RUnlock()
	h, exist := driverManagers[dbType]
	if !exist {
		return nil, fmt.Errorf("driver type %v is not supported", dbType)
	}

	return h.newDriverManagerFunc(log, dbType, config)
}

// InitPlugins init plugins at plugins directory. It should be called on host process.
func InitPlugins(pluginDir string) error {
	if pluginDir == "" {
		return nil
	}

	// read plugin file
	var plugins []os.FileInfo
	if err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "init plugin")
		}

		if info.IsDir() || info.Mode()&0111 == 0 {
			return nil
		}
		plugins = append(plugins, info)
		return nil
	}); err != nil {
		return err
	}

	// register plugin
	for _, p := range plugins {
		binaryPath := filepath.Join(pluginDir, p.Name())

		// check plugin
		if !TestPluginDriver(binaryPath) {
			return fmt.Errorf("unable to load plugin: %v", binaryPath)
		}

		if err := RegisterDriverFromClient(binaryPath); err != nil {
			return err
		}

	}
	return nil
}
