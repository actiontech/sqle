package v1

import (
	"sync"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

var ErrSQLIsNotSupported = errors.New("SQL is not supported")

const (
	// grpc error code
	grpcErrSQLIsNotSupported = 1000
)

const (
	DefaultPluginVersion = 1
)

var defaultPluginSet = map[int]goPlugin.PluginSet{
	DefaultPluginVersion: goPlugin.PluginSet{
		PluginNameAuditDriver: &auditDriverPlugin{},
	},
}

type PluginServer struct {
	plugins map[int]goPlugin.PluginSet
	mutex   *sync.Mutex
}

func NewPlugin() *PluginServer {
	return &PluginServer{
		plugins: map[int]goPlugin.PluginSet{},
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
		GRPCServer: SQLEGrpcServer,
	})
}

var SQLEGrpcServer = goPlugin.DefaultGRPCServer

func (p *PluginServer) AddDriverPlugin(plugin goPlugin.Plugin) {
	p.AddPlugin(PluginNameAuditDriver, DefaultPluginVersion, plugin)
}

func (p *PluginServer) AddQueryDriverPlugin(plugin goPlugin.Plugin) {
	p.AddPlugin(PluginNameQueryDriver, DefaultPluginVersion, plugin)
}

func (p *PluginServer) AddAnalysisDriverPlugin(plugin goPlugin.Plugin) {
	p.AddPlugin(PluginNameAnalysisDriver, DefaultPluginVersion, plugin)
}

func (p *PluginServer) AddPlugin(pluginName string, pluginVersion int, plugin goPlugin.Plugin) {
	p.mutex.Lock()
	if _, ok := p.plugins[pluginVersion]; !ok {
		p.plugins[pluginVersion] = goPlugin.PluginSet{}
	}
	p.plugins[pluginVersion][pluginName] = plugin
	p.mutex.Unlock()
}

var handshakeConfig = goPlugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// ServePlugin start plugin process service. It should be called on plugin process.
// Deprecated: Use PluginServer.AddDriverPlugin and PluginServer.Serve instead.
func ServePlugin(r Registerer, newDriver func(cfg *Config) Driver) {
	name := r.Name()
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: handshakeConfig,

		Plugins: goPlugin.PluginSet{
			name: &auditDriverPlugin{Srv: &auditDriverGRPCServer{r: r, newDriver: newDriver}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}
