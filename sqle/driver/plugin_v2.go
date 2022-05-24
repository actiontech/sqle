package driver

import (
	"context"
	"github.com/actiontech/sqle/sqle/driver/proto"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/sirupsen/logrus"
	"os/exec"
	"sync"

	goPlugin "github.com/hashicorp/go-plugin"
)

const (
	PluginNameDriver = "driver"
)

const (
	DefaultPluginVersion = 1
)

var defaultPluginSet = map[int]goPlugin.PluginSet{
	DefaultPluginVersion: goPlugin.PluginSet{
		PluginNameDriver: &driverPlugin{},
	},
}

type PluginServer struct {
	plugins map[int]goPlugin.PluginSet
	mutex   *sync.Mutex
}

func NewPlugin() *PluginServer {
	return &PluginServer{
		plugins: defaultPluginSet,
		mutex:   &sync.Mutex{},
	}
}

func (p *PluginServer) Serve() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig:  handshakeConfig,
		VersionedPlugins: p.plugins,
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}

func (p *PluginServer) AddPlugin(pluginName string, pluginVersion int, plugin goPlugin.Plugin) {
	p.mutex.Lock()
	if _, ok := p.plugins[pluginVersion]; !ok {
		p.plugins[pluginVersion] = goPlugin.PluginSet{}
	}
	p.plugins[pluginVersion][pluginName] = plugin
	p.mutex.Unlock()
}

const (
	ClientVersionOne = 1
	ClientVersionTwo = 2
)

type PluginClient struct {
	path          string
	c             *goPlugin.Client
	ClientVersion int
}

func (p *PluginClient) Kill() {
	p.c.Kill()
}

func (p *PluginClient) Client() (goPlugin.ClientProtocol, error) {
	p.resetClient()
	return p.c.Client()
}

func (p *PluginClient) resetClient() {
	if p.c != nil {
		p.c.Kill()
	}
	switch p.ClientVersion {
	case ClientVersionOne:
		p.c = goPlugin.NewClient(&goPlugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins: goPlugin.PluginSet{
				PluginNameDriver: &driverPlugin{},
			},
			Cmd:              exec.Command(p.path),
			AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
		})
	case ClientVersionTwo:
		p.c = goPlugin.NewClient(&goPlugin.ClientConfig{
			HandshakeConfig:  handshakeConfig,
			VersionedPlugins: defaultPluginSet,
			Cmd:              exec.Command(p.path),
			AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
		})
	}
}

func testConnClient(client *PluginClient) bool {
	c, err := client.Client()
	if err != nil {
		return false
	}
	defer client.Kill()
	err = c.Ping()
	return err == nil
}

func RegisterDriverFromClient(client *PluginClient) error {
	closeCh := make(chan struct{})
	srv, err := getServerHandle(client, closeCh)
	if err != nil {
		return err
	}
	pluginMeta, err := srv.Metas(context.TODO(), &proto.Empty{})
	if err != nil {
		return err
	}
	close(closeCh)

	// driverRules get from plugin when plugin initialize.
	var driverRules = make([]*Rule, 0, len(pluginMeta.Rules))
	for _, rule := range pluginMeta.Rules {
		driverRules = append(driverRules, convertRuleFromProtoToDriver(rule))
	}

	handler := func(log *logrus.Entry, config *Config) (Driver, error) {
		pluginCloseCh := make(chan struct{})
		srv, err := getServerHandle(client, pluginCloseCh)
		if err != nil {
			return nil, err
		}

		// protoRules send to plugin for Audit.
		var protoRules []*proto.Rule
		for _, rule := range config.Rules {
			protoRules = append(protoRules, convertRuleFromDriverToProto(rule))
		}

		initRequest := &proto.InitRequest{
			Rules: protoRules,
		}
		if config.DSN != nil {
			initRequest.Dsn = &proto.DSN{
				Host:             config.DSN.Host,
				Port:             config.DSN.Port,
				User:             config.DSN.User,
				Password:         config.DSN.Password,
				AdditionalParams: proto.ConvertParamToProtoParam(config.DSN.AdditionalParams),

				// database is to open.
				Database: config.DSN.DatabaseName,
			}
		}

		_, err = srv.Init(context.TODO(), initRequest)
		if err != nil {
			return nil, err
		}

		return &driverPluginClient{srv, pluginCloseCh}, nil

	}

	Register(pluginMeta.Name, handler, driverRules, proto.ConvertProtoParamToParam(pluginMeta.GetAdditionalParams()))

	switch client.ClientVersion {
	case ClientVersionTwo:
		// register of private plugins, general plugins do not need to be placed in exclusiveRegisterPlugin
		err = exclusiveRegisterPlugin(pluginMeta.Name, client)
	}

	if err != nil {
		return err
	}

	log.Logger().WithFields(logrus.Fields{
		"plugin_name":    pluginMeta.Name,
		"plugin_version": client.ClientVersion,
	}).Infoln("plugin inited")
	return nil
}

func getServerHandle(client *PluginClient, closeCh <-chan struct{}) (proto.DriverClient, error) {
	gRPCClient, err := client.Client()
	if err != nil {
		return nil, err
	}
	go func() {
		<-closeCh
		client.Kill()
	}()
	rawI, err := gRPCClient.Dispense(PluginNameDriver)
	if err != nil {
		return nil, err
	}
	// srv can only be proto.DriverClient
	//nolint:forcetypeassert
	srv := rawI.(proto.DriverClient)

	return srv, nil
}

func newClientV1FromFile(path string) *PluginClient {
	return &PluginClient{
		ClientVersion: ClientVersionOne,
		path:          path,
	}
}

func newClientV2FromFile(path string) *PluginClient {
	return &PluginClient{
		ClientVersion: ClientVersionTwo,
		path:          path,
	}
}
