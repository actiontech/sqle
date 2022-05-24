//go:build !enterprise
// +build !enterprise

package driver

func exclusiveRegisterPlugin(pluginName string, c *PluginClient) error {
	return nil
}
