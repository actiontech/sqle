package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowInstanceConfigurationResponse Response Object
type ShowInstanceConfigurationResponse struct {

	// 引擎版本。
	DatastoreVersionName *string `json:"datastore_version_name,omitempty"`

	// 引擎名。
	DatastoreName *ShowInstanceConfigurationResponseDatastoreName `json:"datastore_name,omitempty"`

	// 创建时间，格式为\"yyyy-MM-ddTHH:mm:ssZ\"。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	Created *string `json:"created,omitempty"`

	// 更新时间，格式为\"yyyy-MM-ddTHH:mm:ssZ\"。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	Updated *string `json:"updated,omitempty"`

	// 参数对象，用户基于默认参数模板自定义的参数配置。
	ConfigurationParameters *[]ConfigurationParameter `json:"configuration_parameters,omitempty"`
	HttpStatusCode          int                       `json:"-"`
}

func (o ShowInstanceConfigurationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInstanceConfigurationResponse struct{}"
	}

	return strings.Join([]string{"ShowInstanceConfigurationResponse", string(data)}, " ")
}

type ShowInstanceConfigurationResponseDatastoreName struct {
	value string
}

type ShowInstanceConfigurationResponseDatastoreNameEnum struct {
	MYSQL      ShowInstanceConfigurationResponseDatastoreName
	POSTGRESQL ShowInstanceConfigurationResponseDatastoreName
	SQLSERVER  ShowInstanceConfigurationResponseDatastoreName
	MARIADB    ShowInstanceConfigurationResponseDatastoreName
}

func GetShowInstanceConfigurationResponseDatastoreNameEnum() ShowInstanceConfigurationResponseDatastoreNameEnum {
	return ShowInstanceConfigurationResponseDatastoreNameEnum{
		MYSQL: ShowInstanceConfigurationResponseDatastoreName{
			value: "mysql",
		},
		POSTGRESQL: ShowInstanceConfigurationResponseDatastoreName{
			value: "postgresql",
		},
		SQLSERVER: ShowInstanceConfigurationResponseDatastoreName{
			value: "sqlserver",
		},
		MARIADB: ShowInstanceConfigurationResponseDatastoreName{
			value: "mariadb",
		},
	}
}

func (c ShowInstanceConfigurationResponseDatastoreName) Value() string {
	return c.value
}

func (c ShowInstanceConfigurationResponseDatastoreName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowInstanceConfigurationResponseDatastoreName) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
