package inspector

import (
	"sqle/model"
	"strconv"
	"sync"
)

const (
	CONFIG_DML_ROLLBACK_MAX_ROWS = "dml_rollback_max_rows"
	CONFIG_DDL_OSC_SIZE_LIMIT    = "ddl_osc_size_limit"
)

var configMap = map[string]*model.Config{
	CONFIG_DML_ROLLBACK_MAX_ROWS: &model.Config{
		Name:    CONFIG_DML_ROLLBACK_MAX_ROWS,
		Value:   "",
		Default: "1000",
		Desc:    "在 DML 语句中预计影响行数超过指定值则不回滚",
	},
	CONFIG_DDL_OSC_SIZE_LIMIT: &model.Config{
		Name:    CONFIG_DDL_OSC_SIZE_LIMIT,
		Value:   "",
		Default: "16",
		Desc:    "改表时，表空间超过指定大小(MB)审核时输出osc改写建议",
	},
}
var configMutex sync.Mutex

func UpdateConfig(name, value string) {
	configMutex.Lock()
	if config, ok := configMap[name]; ok {
		config.Value = value
	}
	configMutex.Unlock()
}

func GetConfigInt(name string) int64 {
	configMutex.Lock()
	defer configMutex.Unlock()
	config, ok := configMap[name]
	if !ok {
		return 0
	}
	value, err := strconv.ParseInt(config.Value, 10, 64)
	if err == nil {
		return value
	}
	value, err = strconv.ParseInt(config.Default, 10, 64)
	if err == nil {
		return value
	}
	return 0
}

func GetAllConfig() []model.Config {
	configMutex.Lock()
	configs := make([]model.Config, 0, len(configMap))
	for _, config := range configMap {
		configs = append(configs, *config)
	}
	configMutex.Unlock()
	return configs
}
