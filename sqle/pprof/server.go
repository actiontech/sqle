package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/actiontech/sqle/sqle/log"
)

// StartServer 启动独立的 pprof HTTP 服务器
// port: pprof 服务器监听端口，如果为 0 则不启动
func StartServer(port int) error {
	if port <= 0 {
		log.Logger().Infof("pprof server disabled (port: %d)", port)
		return nil
	}

	address := fmt.Sprintf("0.0.0.0:%d", port)
	log.Logger().Infof("starting pprof server on %s", address)

	// pprof 包在导入时会自动注册路由到 http.DefaultServeMux
	// 只需要启动一个 HTTP 服务器即可
	if err := http.ListenAndServe(address, nil); err != nil {
		return fmt.Errorf("pprof server failed: %v", err)
	}

	return nil
}

// StartServerAsync 异步启动独立的 pprof HTTP 服务器
func StartServerAsync(port int) {
	if port <= 0 {
		log.Logger().Infof("pprof server disabled (port: %d)", port)
		return
	}

	go func() {
		if err := StartServer(port); err != nil {
			log.Logger().Errorf("pprof server error: %v", err)
		}
	}()
}
