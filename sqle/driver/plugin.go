package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pingcap/errors"
)

const (
	DefaultPluginVersion = 1
)

var defaultPluginSet = map[int]goPlugin.PluginSet{
	DefaultPluginVersion: goPlugin.PluginSet{
		PluginNameDriver: &auditDriverPlugin{},
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
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}

func (p *PluginServer) AddDriverPlugin(plugin goPlugin.Plugin) {
	p.AddPlugin(PluginNameDriver, DefaultPluginVersion, plugin)
}

func (p *PluginServer) AddPlugin(pluginName string, pluginVersion int, plugin goPlugin.Plugin) {
	p.mutex.Lock()
	if _, ok := p.plugins[pluginVersion]; !ok {
		p.plugins[pluginVersion] = goPlugin.PluginSet{}
	}
	p.plugins[pluginVersion][pluginName] = plugin
	p.mutex.Unlock()
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
		var client PluginClient
		client = newClientFromFile(binaryPath)
		if !testConnClient(client) {
			client = newOldClientFromFile(binaryPath)
			if !testConnClient(client) {
				return fmt.Errorf("unable to load plugin: %v", binaryPath)
			}
		}
		if err := RegisterDriverFromClient(client); err != nil {
			return err
		}

	}
	return nil
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
