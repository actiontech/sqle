package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/actiontech/sqle/sqle/config"

	"github.com/actiontech/sqle/sqle/driver/common"
	driverV1 "github.com/actiontech/sqle/sqle/driver/v1"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var ErrPluginNotFound = errors.New("plugin not found")

func NewErrPluginAPINotImplement(m driverV2.OptionalModule) error {
	return fmt.Errorf("plugin not implement api %s", m)
}

var BuiltInPluginProcessors = map[string] /*plugin name*/ PluginProcessor{}

type pluginManager struct {
	pluginNames      []string
	metas            map[string]driverV2.DriverMetas
	pluginProcessors map[string]PluginProcessor
}

var PluginManager = &pluginManager{
	pluginNames:      []string{},
	metas:            map[string]driverV2.DriverMetas{},
	pluginProcessors: map[string]PluginProcessor{},
}

func GetPluginManager() *pluginManager {
	return PluginManager
}

func (pm *pluginManager) GetAllRules() map[string][]*driverV2.Rule {
	rules := map[string][]*driverV2.Rule{}
	for _, p := range pm.pluginNames {
		meta := pm.metas[p]
		rules[p] = meta.Rules
	}
	return rules
}

func (pm *pluginManager) AllDrivers() []string {
	return pm.pluginNames
}

func (pm *pluginManager) AllDriverMetas() []*driverV2.DriverMetas {
	metas := make([]*driverV2.DriverMetas, len(pm.metas))

	for i := range pm.pluginNames {
		pluginName := pm.pluginNames[i]
		meta := pm.metas[pluginName]
		metas[i] = &driverV2.DriverMetas{
			PluginName:          meta.PluginName,
			DatabaseDefaultPort: meta.DatabaseDefaultPort,
			Logo:                meta.Logo,
		}
	}

	return metas
}

func (pm *pluginManager) AllLogo() map[string][]byte {
	logoMap := map[string][]byte{}
	for _, pluginName := range pm.pluginNames {
		meta := pm.metas[pluginName]
		if meta.Logo != nil {
			logoMap[pluginName] = meta.Logo
		}
	}
	return logoMap
}

func (pm *pluginManager) AllAdditionalParams() map[string] /*driver name*/ params.Params {
	newParams := map[string]params.Params{}
	for k, v := range pm.metas {
		newParams[k] = v.DatabaseAdditionalParams.Copy()
	}
	return newParams
}

func (pm *pluginManager) IsOptionalModuleEnabled(pluginName string, expectModule driverV2.OptionalModule) bool {
	meta, ok := pm.metas[pluginName]
	if !ok {
		return false
	}
	for _, m := range meta.EnabledOptionalModule {
		if m == expectModule {
			return true
		}
	}
	return false
}

func (pm *pluginManager) register(pp PluginProcessor) error {
	meta, err := pp.GetDriverMetas()
	if err != nil {
		return err
	}
	if _, ok := pm.metas[meta.PluginName]; ok {
		return fmt.Errorf("duplicated driver name %s", meta.PluginName)
	}
	pm.pluginNames = append(pm.pluginNames, meta.PluginName)
	pm.metas[meta.PluginName] = *meta
	pm.pluginProcessors[meta.PluginName] = pp
	return nil
}

func getClientConfig(cmdBase string, cmdArgs []string) *goPlugin.ClientConfig {
	return &goPlugin.ClientConfig{
		HandshakeConfig: driverV2.HandshakeConfig,
		VersionedPlugins: map[int]goPlugin.PluginSet{
			driverV1.ProtocolVersion: driverV1.PluginSet,
			driverV2.ProtocolVersion: driverV2.PluginSet,
		},
		Cmd:              exec.Command(cmdBase, cmdArgs...),
		AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
		GRPCDialOptions:  common.GRPCDialOptions,
		StartTimeout:     10 * time.Minute,
	}
}

func (pm *pluginManager) Start(pluginDir string, pluginConfigList []config.PluginConfig) error {
	// register built-in plugin, now is MySQL.
	for name, b := range BuiltInPluginProcessors {
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
		cmdBase := filepath.Join(pluginDir, p.Name())
		cmdArgs := make([]string, 0)

		for _, pluginConfig := range pluginConfigList {
			if p.Name() == pluginConfig.PluginName {
				cmdBase = "sh"
				cmdArgs = append(cmdArgs, "-c", pluginConfig.CMD)
				break
			}
		}

		client := goPlugin.NewClient(getClientConfig(cmdBase, cmdArgs))
		_, err := client.Client()
		if err != nil {
			return err
		}

		var pp PluginProcessor
		switch client.NegotiatedVersion() {
		case driverV1.ProtocolVersion:
			pp = &PluginProcessorV1{cfg: getClientConfig, cmdBase: cmdBase, cmdArgs: cmdArgs, client: client}
		case driverV2.ProtocolVersion:
			pp = &PluginProcessorV2{cfg: getClientConfig, cmdBase: cmdBase, cmdArgs: cmdArgs, client: client}
		}
		if err := pm.register(pp); err != nil {
			stopErr := pp.Stop()
			if stopErr != nil {
				log.NewEntry().Warnf("stop plugin %s failed, error: %v", p.Name(), stopErr)
			}
			return fmt.Errorf("unable to load plugin: %v, error: %v", p.Name(), err)
		}

	}
	return nil
}

func (pm *pluginManager) Stop() {
	for name, pp := range pm.pluginProcessors {
		err := pp.Stop()
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

func (pm *pluginManager) OpenPlugin(l *logrus.Entry, pluginName string, cfg *driverV2.Config) (Plugin, error) {
	if !pm.isPluginExists(pluginName) {
		return nil, ErrPluginNotFound
	}
	return pm.pluginProcessors[pluginName].Open(l, cfg)
}
