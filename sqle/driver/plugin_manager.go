package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	v1 "github.com/actiontech/sqle/sqle/driver/v1"
	v2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"google.golang.org/grpc"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
)

var ErrPluginNotFound = errors.New("plugin not found")

var BuiltInPluginBoots = map[string] /*plugin name*/ PluginBoot{}

type pluginManager struct {
	pluginNames []string
	metas       map[string]v2.DriverMetas
	driverBoots map[string]PluginBoot
}

var PluginManager = &pluginManager{
	pluginNames: []string{},
	metas:       map[string]v2.DriverMetas{},
	driverBoots: map[string]PluginBoot{},
}

func GetPluginManager() *pluginManager {
	return PluginManager
}

func (pm *pluginManager) GetAllRules() map[string][]*v2.Rule {
	rules := map[string][]*v2.Rule{}
	for _, p := range pm.pluginNames {
		meta := pm.metas[p]
		rules[p] = meta.Rules
	}
	return rules
}

func (pm *pluginManager) AllDrivers() []string {
	return pm.pluginNames
}

func (pm *pluginManager) AllAdditionalParams() map[string] /*driver name*/ params.Params {
	newParams := map[string]params.Params{}
	for k, v := range pm.metas {
		newParams[k] = v.DatabaseAdditionalParams.Copy()
	}
	return newParams
}

func (pm *pluginManager) register(boot PluginBoot) error {
	meta, err := boot.Register()
	if err != nil {
		return err
	}
	if _, ok := pm.metas[meta.PluginName]; ok {
		return fmt.Errorf("duplicated driver name %s", meta.PluginName)
	}
	pm.pluginNames = append(pm.pluginNames, meta.PluginName)
	pm.metas[meta.PluginName] = *meta
	pm.driverBoots[meta.PluginName] = boot
	return nil
}

var SQLEGRPCDialOptions = []grpc.DialOption{}

func getClientConfig(path string) *goPlugin.ClientConfig {
	return &goPlugin.ClientConfig{
		HandshakeConfig: v2.HandshakeConfig,
		VersionedPlugins: map[int]goPlugin.PluginSet{
			v1.ProtocolVersion: v1.PluginSet,
			v2.ProtocolVersion: v2.PluginSet,
		},
		Cmd:              exec.Command(path),
		AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
		GRPCDialOptions:  SQLEGRPCDialOptions,
	}
}

func (pm *pluginManager) Start(pluginDir string) error {
	// register built-in plugin, now is MySQL.
	for name, b := range BuiltInPluginBoots {
		err := pm.register(b)
		if err != nil {
			return fmt.Errorf("start built-in %s plugin failed, error: %v", name, err)
		}
	}

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
		path := filepath.Join(pluginDir, p.Name())

		client := goPlugin.NewClient(getClientConfig(path))
		_, err := client.Client()
		if err != nil {
			return err
		}

		var boot PluginBoot
		switch client.NegotiatedVersion() {
		case v1.ProtocolVersion:
			boot = &PluginBootV1{cfg: getClientConfig, path: path, client: client}
		case v2.ProtocolVersion:
			boot = &PluginBootV2{client: client}
		}
		err = pm.register(boot)
		if err != nil {
			return fmt.Errorf("unable to load plugin: %v, error: %v", path, err)
		}
	}
	return nil
}

func (pm *pluginManager) Stop() {
	for name, b := range pm.driverBoots {
		err := b.Stop()
		if err != nil {
			log.NewEntry().Warnf("stop %s plugin failed, error: %v", name, err)
		}
	}
}

func (pm *pluginManager) isPluginExists(pluginName string) bool {
	if _, ok := pm.metas[pluginName]; ok {
		return true
	}
	return false
}

func (pm *pluginManager) OpenPlugin(pluginName string, cfg *v2.Config) (Plugin, error) {
	if !pm.isPluginExists(pluginName) {
		return nil, ErrPluginNotFound
	}
	return pm.driverBoots[pluginName].Open(cfg)
}
