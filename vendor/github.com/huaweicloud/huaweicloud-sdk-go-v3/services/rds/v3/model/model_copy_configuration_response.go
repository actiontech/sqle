package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CopyConfigurationResponse Response Object
type CopyConfigurationResponse struct {

	// 参数模板ID。
	Id *string `json:"id,omitempty"`

	// 参数模板名称。
	Name *string `json:"name,omitempty"`

	// 描述。
	Description *string `json:"description,omitempty"`

	// 数据库版本名称
	DatastoreVersionName *string `json:"datastore_version_name,omitempty"`

	// 数据库名称。
	DatastoreName *string `json:"datastore_name,omitempty"`

	// 创建时间，格式为\"yyyy-MM-ddTHH:mm:ssZ\"。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	CreateTime *string `json:"create_time,omitempty"`

	// 更新时间，格式为\"yyyy-MM-ddTHH:mm:ssZ\"。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	UpdateTime     *string `json:"update_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CopyConfigurationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CopyConfigurationResponse struct{}"
	}

	return strings.Join([]string{"CopyConfigurationResponse", string(data)}, " ")
}
