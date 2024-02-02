package v1

import base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"

// swagger:parameters RegisterDMSProxyTarget
type RegisterDMSProxyTargetReq struct {
	// register dms proxy
	// in:body
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
