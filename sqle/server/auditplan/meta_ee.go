//go:build enterprise
// +build enterprise

package auditplan

import "github.com/actiontech/sqle/sqle/pkg/params"

const (
	TypeOceanBaseForMySQLMybatis = "ocean_base_for_mysql_mybatis"
	TypeOceanBaseForMySQLTopSQL  = "ocean_base_for_mysql_top_sql"
	TypeDB2TopSQL                = "db2_top_sql"
)

const (
	InstanceTypeOceanBaseForMySQL = "OceanBase For MySQL"
	InstanceTypeDB2               = "DB2"
)

const (
	paramKeyIndicator = "indicator"
	paramKeyTopN      = "top_n"
)

var EEMetas = []Meta{
	{
		Type:         TypeOceanBaseForMySQLMybatis,
		Desc:         "Mybatis 扫描",
		InstanceType: InstanceTypeOceanBaseForMySQL,
		CreateTask:   NewDefaultTask,
	},
	{
		Type:         TypeOceanBaseForMySQLTopSQL,
		Desc:         "Top SQL",
		InstanceType: InstanceTypeOceanBaseForMySQL,
		CreateTask:   NewOBMySQLTopSQLTask,
		Params: []*params.Param{
			{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyTopN,
				Desc:  "Top N",
				Value: "3",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyIndicator,
				Desc:  "关注指标",
				Value: OBMySQLIndicatorElapsedTime,
				Type:  params.ParamTypeString,
			},
		},
	},
	{
		Type:         TypeDB2TopSQL,
		Desc:         "DB2 Top SQL",
		InstanceType: InstanceTypeDB2,
		CreateTask:   NewDB2TopSQLTask,
		Params: []*params.Param{
			{
				Key:   paramKeyCollectIntervalMinute,
				Desc:  "采集周期（分钟）",
				Value: "60",
				Type:  params.ParamTypeInt,
			},
			{
				Key:   paramKeyTopN,
				Desc:  "Top N",
				Value: "3",
				Type:  params.ParamTypeInt,
			},
		},
	},
}

func init() {
	for _, meta := range EEMetas {
		Metas = append(Metas, meta)
		MetaMap[meta.Type] = meta
	}
}
