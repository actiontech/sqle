package driver

import (
	"github.com/actiontech/sqle/sqle/pkg/params"
	goPlugin "github.com/hashicorp/go-plugin"
	"os/exec"
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

type PluginClient interface {
	Kill()
	Client() (goPlugin.ClientProtocol, error)
	RegisterPlugin(c PluginClient) error
}

type pluginClientOld struct {
	path string
	c    *goPlugin.Client
}

func newOldClientFromFile(path string) *pluginClientOld {
	return &pluginClientOld{
		path: path,
	}
}

func (p *pluginClientOld) Kill() {
	p.c.Kill()
}

func (p *pluginClientOld) Client() (goPlugin.ClientProtocol, error) {
	p.resetClient()
	return p.c.Client()
}

func (p *pluginClientOld) RegisterPlugin(c PluginClient) error {
	_, err := registerAuditDriver(c)
	return err
}

func (p *pluginClientOld) resetClient() {
	if p.c != nil {
		p.c.Kill()
	}
	p.c = goPlugin.NewClient(&goPlugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: goPlugin.PluginSet{
			PluginNameDriver: &auditDriverPlugin{},
		},
		Cmd:              exec.Command(p.path),
		AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
	})
}

type pluginClient struct {
	path string
	c    *goPlugin.Client
}

func newClientFromFile(path string) *pluginClient {
	return &pluginClient{
		path: path,
	}
}

func (p *pluginClient) Kill() {
	p.c.Kill()
}

func (p *pluginClient) Client() (goPlugin.ClientProtocol, error) {
	p.resetClient()
	return p.c.Client()
}

func (p *pluginClient) RegisterPlugin(c PluginClient) error {
	pluginName, err := registerAuditDriver(c)
	if err != nil {
		return err
	}
	return registerPlugin(pluginName, c)
}

func (p *pluginClient) resetClient() {
	if p.c != nil {
		p.c.Kill()
	}
	p.c = goPlugin.NewClient(&goPlugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		VersionedPlugins: defaultPluginSet,
		Cmd:              exec.Command(p.path),
		AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
	})
}

func testConnClient(client PluginClient) bool {
	c, err := client.Client()
	if err != nil {
		return false
	}
	defer client.Kill()
	err = c.Ping()
	return err == nil
}

func RegisterDriverFromClient(client PluginClient) error {
	return client.RegisterPlugin(client)
}
