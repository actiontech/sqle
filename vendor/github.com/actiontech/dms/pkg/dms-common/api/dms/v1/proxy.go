package v1

import (
	"fmt"

	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
)

// swagger:model
type RegisterDMSProxyTargetReq struct {
	DMSProxyTarget *DMSProxyTarget `json:"dms_proxy_target" validate:"required"`
}

// A dms proxy target
type DMSProxyTarget struct {
	// target name
	// Required: true
	Name string `json:"name" validate:"required"`
	// target addr, eg: http://10.1.2.1:5432
	// Required: true
	Addr string `json:"addr" validate:"required,url"`
	// version number
	// Required: true
	Version string `json:"version" validate:"required"`
	// url prefix that need to be proxy, eg: /v1/user
	// Required: true
	ProxyUrlPrefixs []string `json:"proxy_url_prefixs" validate:"required"`
	// the scenario is used to differentiate scenarios
	Scenario ProxyScenario `json:"scenario"`
}

func (s *DMSProxyTarget) String() string {
	return fmt.Sprintf("{name: %v, addr: %v, version: %v, Scenario %v}", s.Name, s.Addr, s.Version, s.Scenario)
}

// swagger:enum ProxyScenario
type ProxyScenario string

const (
	ProxyScenarioInternalService     ProxyScenario = "internal_service"
	ProxyScenarioThirdPartyIntegrate ProxyScenario = "thrid_party_integrate"
)

// swagger:model RegisterDMSProxyTargetReply
type RegisterDMSProxyTargetReply struct {
	// Generic reply
	base.GenericResp
}
