package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Proxy struct {

	// Proxy实例ID。
	PoolId string `json:"pool_id"`

	// Proxy实例开启状态，取值范围如下。 - open：打开。 - closed：关闭。 - frozen：已冻结。 - opening：打开中。 - closing：关闭中。 - freezing：冻结中。 - unfreezing：解冻中。
	Status string `json:"status"`

	// Proxy读写分离地址。
	Address string `json:"address"`

	// elb模式的虚拟IP信息。
	ElbVip string `json:"elb_vip"`

	// 弹性公网IP信息。
	Eip string `json:"eip"`

	// Proxy端口信息。
	Port int32 `json:"port"`

	// Proxy实例状态。 - abnormal：异常。 - normal：正常。 - creating：创建中。 - deleted：已删除。
	PoolStatus string `json:"pool_status"`

	// 延时阈值（单位：KB）。
	DelayThresholdInKilobytes int32 `json:"delay_threshold_in_kilobytes"`

	// Proxy实例规格的CPU数量。
	Cpu string `json:"cpu"`

	// Proxy实例规格的内存数量。
	Mem string `json:"mem"`

	// Proxy节点个数。
	NodeNum int32 `json:"node_num"`

	// Proxy节点信息。
	Nodes []ProxyNode `json:"nodes"`

	// Proxy主备模式，取值范围：Ha。
	Mode string `json:"mode"`
}

func (o Proxy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Proxy struct{}"
	}

	return strings.Join([]string{"Proxy", string(data)}, " ")
}
