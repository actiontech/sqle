package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ConfigurationSummaryForCreate 参数模板信息。
type ConfigurationSummaryForCreate struct {

	// 参数组ID。
	Id string `json:"id"`

	// 参数组名称。
	Name string `json:"name"`

	// 参数组描述。
	Description *string `json:"description,omitempty"`

	// 引擎版本。
	DatastoreVersionName string `json:"datastore_version_name"`

	// 引擎名。
	DatastoreName ConfigurationSummaryForCreateDatastoreName `json:"datastore_name"`

	// 创建时间，格式为\"yyyy-MM-ddTHH:mm:ssZ\"。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	Created string `json:"created"`

	// 更新时间，格式为\"yyyy-MM-ddTHH:mm:ssZ\"。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	Updated string `json:"updated"`
}

func (o ConfigurationSummaryForCreate) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigurationSummaryForCreate struct{}"
	}

	return strings.Join([]string{"ConfigurationSummaryForCreate", string(data)}, " ")
}

type ConfigurationSummaryForCreateDatastoreName struct {
	value string
}

type ConfigurationSummaryForCreateDatastoreNameEnum struct {
	MYSQL      ConfigurationSummaryForCreateDatastoreName
	POSTGRESQL ConfigurationSummaryForCreateDatastoreName
	SQLSERVER  ConfigurationSummaryForCreateDatastoreName
	MARIADB    ConfigurationSummaryForCreateDatastoreName
}

func GetConfigurationSummaryForCreateDatastoreNameEnum() ConfigurationSummaryForCreateDatastoreNameEnum {
	return ConfigurationSummaryForCreateDatastoreNameEnum{
		MYSQL: ConfigurationSummaryForCreateDatastoreName{
			value: "mysql",
		},
		POSTGRESQL: ConfigurationSummaryForCreateDatastoreName{
			value: "postgresql",
		},
		SQLSERVER: ConfigurationSummaryForCreateDatastoreName{
			value: "sqlserver",
		},
		MARIADB: ConfigurationSummaryForCreateDatastoreName{
			value: "mariadb",
		},
	}
}

func (c ConfigurationSummaryForCreateDatastoreName) Value() string {
	return c.value
}

func (c ConfigurationSummaryForCreateDatastoreName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ConfigurationSummaryForCreateDatastoreName) UnmarshalJSON(b []byte) error {
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
