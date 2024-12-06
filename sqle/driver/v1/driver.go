package driverV1

import (
	"context"
	"fmt"
	"sync"

	"github.com/actiontech/sqle/sqle/driver/v1/proto"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
)

var handshakeConfig = goPlugin.HandshakeConfig{
	ProtocolVersion:  ProtocolVersion,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

const (
	ProtocolVersion = 1
)

const (
	PluginNameAuditDriver    = "audit-driver"
	PluginNameQueryDriver    = "query-driver"
	PluginNameAnalysisDriver = "analysis-driver"
)

var PluginSet = goPlugin.PluginSet{
	PluginNameAuditDriver:    &auditDriverPlugin{},
	PluginNameAnalysisDriver: &analysisDriverPlugin{},
	PluginNameQueryDriver:    &queryDriverPlugin{},
}

var (
	// auditDrivers store instantiate handlers for MySQL or gRPC plugin.
	auditDrivers = make(map[string]struct{})
	driversMu    sync.RWMutex
)

var queryDriverMu = &sync.RWMutex{}
var queryDrivers = make(map[string]struct{})

var analysisDriverMu = &sync.RWMutex{}
var analysisDrivers = make(map[string]struct{})

func checkQueryDriver(cp goPlugin.ClientProtocol) error {
	rawI, err := cp.Dispense(PluginNameQueryDriver)
	if err != nil {
		return err //todo
	}
	// srv can only be proto.QueryDriverClient
	//nolint:forcetypeassert
	s := rawI.(proto.QueryDriverClient)

	// The test target plugin implements the QueryDriver plugin
	_, err = s.Init(context.TODO(), &proto.InitRequest{})
	if err != nil {
		return err //todo
	}
	return nil
}

func checkAnalysisDriver(cp goPlugin.ClientProtocol) error {
	rawI, err := cp.Dispense(PluginNameAnalysisDriver)
	if err != nil {
		return err
	}
	//nolint:forcetypeassert
	a := rawI.(proto.AnalysisDriverClient)

	// The test target plugin implements the AnalysisDriver plugin
	_, err = a.Init(context.TODO(), &proto.AnalysisDriverInitRequest{})
	if err != nil {
		return err
	}
	return nil
}

func RegisterDrivers(c *goPlugin.Client, clientCfg func(cmdBase string, cmdArgs []string) *goPlugin.ClientConfig,
	cmdBase string, cmdArgs []string) (pluginName string, rules []*Rule, additionalParams params.Params,
	enableQuery, enableSQLAnalysis bool, err error) {
	cp, err := c.Client()
	if err != nil {
		log.NewEntry().Errorf("test conn plugin failed: %v", err)
		return "", nil, nil, false, false, err
	}
	defer c.Kill()

	err = cp.Ping()
	if err != nil {
		log.NewEntry().Errorf("test conn plugin failed: %v", err)
		return "", nil, nil, false, false, err
	}

	rawI, err := cp.Dispense(PluginNameAuditDriver)
	if err != nil {
		return "", nil, nil, false, false, err
	}
	// client can only be proto.DriverClient
	//nolint:forcetypeassert
	client := rawI.(proto.DriverClient)

	pluginMeta, err := client.Metas(context.TODO(), &proto.Empty{})
	if err != nil {
		return "", nil, nil, false, false, err
	}
	l := log.Logger().WithField("plugin_name", pluginMeta.Name)

	_, exist := auditDrivers[pluginMeta.Name]
	if exist {
		panic("duplicated driver name")
	}

	// init audit driver, so that we can use Close to inform all plugins with the same progress to recycle resource
	_, err = client.Init(context.TODO(), &proto.InitRequest{})
	if err != nil {
		return "", nil, nil, false, false, err
	}
	l.WithField("plugin_type", PluginNameAuditDriver).Infoln("plugin inited")

	driversMu.Lock()
	auditDrivers[pluginMeta.Name] = struct{}{}
	driversMu.Unlock()

	rules = make([]*Rule, 0, len(pluginMeta.Rules))
	for _, rule := range pluginMeta.Rules {
		rules = append(rules, convertRuleFromProtoToDriver(rule))
	}

	// check and init query driver
	err = checkQueryDriver(cp)
	if err != nil {
		l.WithField("plugin_type", PluginNameQueryDriver).Infof("plugin not exist or failed to load. err: %v", err)
	} else {
		enableQuery = true
		queryDriverMu.Lock()
		queryDrivers[pluginMeta.Name] = struct{}{}
		queryDriverMu.Unlock()
		l.WithField("plugin_type", PluginNameQueryDriver).Infoln("plugin inited")
	}

	// check and init analysis driver
	err = checkAnalysisDriver(cp)
	if err != nil {
		l.WithField("plugin_type", PluginNameAnalysisDriver).Infof("plugin not exist or failed to load. err: %v", err)
	} else {
		enableSQLAnalysis = true
		analysisDriverMu.Lock()
		analysisDrivers[pluginMeta.Name] = struct{}{}
		analysisDriverMu.Unlock()
		l.WithField("plugin_type", PluginNameAnalysisDriver).Infoln("plugin inited")
	}

	// to be compatible with old plugins
	// the old plugin will panic if it call close() here
	if pluginMeta.Version >= DefaultPluginVersion {
		_, err = client.Close(context.TODO(), &proto.Empty{})
		if err != nil {
			log.Logger().Errorf("gracefully close plugins failed, will force kill the sub progress. err: %v", err)
		}
	}

	handler := func(log *logrus.Entry, dbType string, config *Config) (DriverManager, error) {

		client := goPlugin.NewClient(clientCfg(cmdBase, cmdArgs))
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

	RegisterDriverManger(pluginMeta.Name, handler)

	return pluginMeta.Name, rules, proto.ConvertProtoParamToParam(pluginMeta.GetAdditionalParams()), enableQuery, enableSQLAnalysis, nil
}

// DSN provide necessary information to connect to database.
type DSN struct {
	Host                    string
	Port                    string
	User                    string
	Password                string
	AdditionalParams        params.Params
	SQLAllowQueryAuditLevel string
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

func (p *PluginClient) Kill() {
	if p.c != nil {
		p.c.Kill()
	}
}

func (p *PluginClient) Client() (goPlugin.ClientProtocol, error) {
	return p.c.Client()
}

var driverManagerMu = &sync.RWMutex{}
var driverManagers = make(map[string]driverManagerHandler)

type driverManagerHandler struct {
	newDriverManagerFunc newDriverManagerHandler
}

type newDriverManagerHandler func(log *logrus.Entry, dbType string, config *Config) (DriverManager, error)

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

func NewDriverManger(log *logrus.Entry, dbType string, config *Config) (DriverManager, error) {
	driverManagerMu.RLock()
	defer driverManagerMu.RUnlock()
	h, exist := driverManagers[dbType]
	if !exist {
		return nil, fmt.Errorf("driver type %v is not supported", dbType)
	}

	return h.newDriverManagerFunc(log, dbType, config)
}

func (d *PluginDriverManager) GetAuditDriver() (Driver, error) {
	if d.auditPluginClient == nil {
		return nil, fmt.Errorf("audit driver type %v is not supported", d.dbType)
	}
	return &driverImpl{d.auditPluginClient, d.pluginCloseCh}, nil
}

func (d *PluginDriverManager) GetSQLQueryDriver() (SQLQueryDriver, error) {
	if d.queryPluginClient == nil {
		return nil, fmt.Errorf("SQL query driver type %v is not supported", d.dbType)
	}
	return &queryDriverImpl{d.queryPluginClient}, nil
}

func (d *PluginDriverManager) GetAnalysisDriver() (AnalysisDriver, error) {
	if d.analysisPluginClient == nil {
		return nil, fmt.Errorf("analysis driver type %v is not supported", d.dbType)
	}
	return &analysisDriverImpl{d.analysisPluginClient}, nil
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
