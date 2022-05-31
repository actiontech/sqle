//go:build !enterprise
// +build !enterprise

package driver

func registerPlugin(pluginName string, c PluginClient) error {
	return nil
}
