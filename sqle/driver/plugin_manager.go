package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/hashicorp/go-hclog"

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

func (pm *pluginManager) GetDriverMetasOfPlugin(pluginName string) *driverV2.DriverMetas {
	if dm, exist := pm.metas[pluginName]; exist {
		return &dm
	}
	return nil
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
			PluginName:            meta.PluginName,
			DatabaseDefaultPort:   meta.DatabaseDefaultPort,
			Logo:                  meta.Logo,
			EnabledOptionalModule: meta.EnabledOptionalModule,
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
	return meta.IsOptionalModuleEnabled(expectModule)
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
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   "plugin-client",
			Output: log.Logger().Out,
			Level:  hclog.Trace,
		}),
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

	// kill plugins process residual and remove pidfile
	var wg sync.WaitGroup
	dir := GetPluginPidDirPath(pluginDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".pid") {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = KillResidualPluginsProcess(path)
				if err != nil {
					log.NewEntry().Warnf("stop residual plugin %s error: %v", path, err)
				}
			}()
		}
		return nil
	}); err != nil {
		log.NewEntry().Warnf("stop residual plugin file path walk error: %v", err)
	}
	wg.Wait()

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

		if len(cmdArgs) == 0 && strings.HasSuffix(p.Name(), ".jar") {
			javaPluginCmd := fmt.Sprintf("java -jar %s", cmdBase)
			cmdBase = "sh"
			cmdArgs = append(cmdArgs, "-c", javaPluginCmd)
		}

		client := goPlugin.NewClient(getClientConfig(cmdBase, cmdArgs))
		_, err := client.Client()
		if err != nil {
			return fmt.Errorf("plugin %v failed to start, error: %v Please check the sqled.log for more details", p.Name(), err)
		}

		pluginPidFilePath := GetPluginPidFilePath(pluginDir, p.Name())
		err = WritePidFile(pluginPidFilePath, int64(client.ReattachConfig().Pid))
		if err != nil {
			return fmt.Errorf("write plugin %s pid file failed, error: %v", pluginPidFilePath, err)
		}
		var pp PluginProcessor
		switch client.NegotiatedVersion() {
		case driverV1.ProtocolVersion:
			pp = &PluginProcessorV1{cfg: getClientConfig, cmdBase: cmdBase, cmdArgs: cmdArgs, client: client}
		case driverV2.ProtocolVersion:
			pp = &PluginProcessorV2{cfg: getClientConfig, cmdBase: cmdBase, cmdArgs: cmdArgs, client: client, pluginPidFilePath: pluginPidFilePath}
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

func KillResidualPluginsProcess(pidFile string) error {
	process, err := GetProcessByPidFile(pidFile)
	if err != nil {
		return fmt.Errorf("get plugin %s process failed, error: %v", pidFile, err)
	}
	if process != nil {
		err = StopProcess(process)
		if err != nil {
			return fmt.Errorf("stop plugin process [%v] failed, error: %v", process.Pid, err)
		}
	}
	err = os.Remove(pidFile)
	if err != nil {
		return fmt.Errorf("remove pid file %s error: %v", pidFile, err)
	}
	return nil
}

// 根据pid文件获取进程信息
func GetProcessByPidFile(pluginPidFile string) (*os.Process, error) {
	if _, err := os.Stat(pluginPidFile); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		pidContent, err := os.ReadFile(pluginPidFile)
		if err != nil {
			return nil, err
		}
		if len(pidContent) == 0 {
			return nil, nil
		}
		pid, err := strconv.Atoi(string(pidContent))
		if err != nil {
			return nil, err
		}
		// 获取进程
		process, err := GetProcessByPid(pid)
		if err != nil {
			return nil, err
		}
		return process, nil
	}
	return nil, nil
}

// 根据pid获取进程信息，若进程已退出则返回nil
func GetProcessByPid(pid int) (*os.Process, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}
	// 检查进程是否存在的方式
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return nil, nil
		}
		return nil, err
	}
	return process, nil
}

// 退出进程
func StopProcess(process *os.Process) error {
	doneChan := time.NewTicker(2 * time.Second)
	defer doneChan.Stop()
	for {
		select {
		case <-doneChan.C:
			log.NewEntry().Warnf("stop plugin process [%v] failed, just kill it ", process.Pid)
			err := process.Kill()
			if err != nil {
				return err
			}
			return nil
		default:
			err := process.Signal(syscall.SIGTERM)
			if errors.Is(err, os.ErrProcessDone) {
				return nil
			}
		}
	}
}

func GetPluginPidDirPath(pluginDir string) string {
	return filepath.Join(pluginDir, "pidfile")
}

func GetPluginPidFilePath(pluginDir string, pluginName string) string {
	return filepath.Join(GetPluginPidDirPath(pluginDir), pluginName+".pid")
}

func WritePidFile(pidFilePath string, pid int64) error {
	if err := os.MkdirAll(filepath.Dir(pidFilePath), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(pidFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d", pid)
	if err != nil {
		return err
	}
	return nil
}
