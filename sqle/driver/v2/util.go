package driverV2

import (
	"fmt"
	"math/rand"
	"time"

	protoV2 "github.com/actiontech/sqle/sqle/driver/v2/proto"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

const (
	// grpc error code
	GrpcErrSQLIsNotSupported = 1000
)

const (
	SQLTypeDML = "dml"
	SQLTypeDDL = "ddl"
	SQLTypeDQL = "dql"
)

const (
	DriverTypeMySQL          = "MySQL"
	DriverTypePostgreSQL     = "PostgreSQL"
	DriverTypeTiDB           = "TiDB"
	DriverTypeSQLServer      = "SQL Server"
	DriverTypeOracle         = "Oracle"
	DriverTypeDB2            = "DB2"
	DriverTypeOceanBase      = "OceanBase For MySQL"
	DriverTypeTDSQLForInnoDB = "TDSQL For InnoDB"
)

type DriverNotSupportedError struct {
	DriverTyp string
}

func (e *DriverNotSupportedError) Error() string {
	return fmt.Sprintf("driver type %v is not supported", e.DriverTyp)
}

type OptionalModule uint32

const (
	OptionalModuleGenRollbackSQL = iota
	OptionalModuleQuery
	OptionalModuleExplain
	OptionalModuleGetTableMeta
	OptionalModuleExtractTableFromSQL
	OptionalModuleEstimateSQLAffectRows
	OptionalModuleKillProcess
)

func (m OptionalModule) String() string {
	switch m {
	case OptionalModuleGenRollbackSQL:
		return "GenRollbackSQL"
	case OptionalModuleQuery:
		return "Query"
	case OptionalModuleExplain:
		return "Explain"
	case OptionalModuleGetTableMeta:
		return "GetTableMeta"
	case OptionalModuleExtractTableFromSQL:
		return "ExtractTableFromSQL"
	case OptionalModuleEstimateSQLAffectRows:
		return "EstimateSQLAffectRows"
	case OptionalModuleKillProcess:
		return "KillProcess"
	default:
		return "Unknown"
	}
}

type DriverMetas struct {
	PluginName               string
	DatabaseDefaultPort      int64
	Logo                     []byte
	DatabaseAdditionalParams params.Params
	Rules                    []*Rule
	EnabledOptionalModule    []OptionalModule
}

func ConvertRuleFromProtoToDriver(rule *protoV2.Rule) *Rule {
	var ps = make(params.Params, 0, len(rule.Params))
	for _, p := range rule.Params {
		ps = append(ps, &params.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  params.ParamType(p.Type),
		})
	}
	return &Rule{
		Name:       rule.Name,
		Category:   rule.Category,
		Desc:       rule.Desc,
		Annotation: rule.Annotation,
		Level:      RuleLevel(rule.Level),
		Params:     ps,
		Knowledge:  RuleKnowledge{Content: rule.Knowledge.GetContent()},
	}
}

func ConvertRuleFromDriverToProto(rule *Rule) *protoV2.Rule {
	var params = make([]*protoV2.Param, 0, len(rule.Params))
	for _, p := range rule.Params {
		params = append(params, &protoV2.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  string(p.Type),
		})
	}
	return &protoV2.Rule{
		Name:       rule.Name,
		Desc:       rule.Desc,
		Annotation: rule.Annotation,
		Level:      string(rule.Level),
		Category:   rule.Category,
		Params:     params,
		Knowledge: &protoV2.Knowledge{
			Content: rule.Knowledge.Content,
		},
	}
}

func ConvertParamToProtoParam(p params.Params) []*protoV2.Param {
	pp := make([]*protoV2.Param, len(p))
	for i, v := range p {
		if v == nil {
			continue
		}
		pp[i] = &protoV2.Param{
			Key:   v.Key,
			Value: v.Value,
			Desc:  v.Desc,
			Type:  string(v.Type),
		}
	}
	return pp
}

func ConvertProtoParamToParam(p []*protoV2.Param) params.Params {
	pp := make(params.Params, len(p))
	for i, v := range p {
		if v == nil {
			continue
		}
		pp[i] = &params.Param{
			Key:   v.Key,
			Value: v.Value,
			Desc:  v.Desc,
			Type:  params.ParamType(v.Type),
		}
	}
	return pp
}

func ConvertTabularDataToProto(td TabularData) *protoV2.TabularData {
	columns := make([]*protoV2.TabularDataHead, 0, len(td.Columns))
	for _, c := range td.Columns {
		columns = append(columns, &protoV2.TabularDataHead{
			Name: c.Name,
			Desc: c.Desc,
		})
	}

	rows := make([]*protoV2.TabularDataRows, 0, len(td.Rows))
	for _, r := range td.Rows {
		rows = append(rows, &protoV2.TabularDataRows{
			Items: r,
		})
	}

	return &protoV2.TabularData{
		Columns: columns,
		Rows:    rows,
	}
}

func ConvertProtoTabularDataToDriver(pTd *protoV2.TabularData) TabularData {
	columns := make([]TabularDataHead, 0, len(pTd.Columns))
	for _, c := range pTd.Columns {
		columns = append(columns, TabularDataHead{
			Name: c.Name,
			Desc: c.Desc,
		})
	}

	rows := make([][]string, 0, len(pTd.Rows))
	for _, r := range pTd.Rows {
		rows = append(rows, r.Items)
	}

	return TabularData{
		Columns: columns,
		Rows:    rows,
	}
}

func ConvertTableMetaToProto(meta *TableMeta) *protoV2.TableMeta {
	return &protoV2.TableMeta{
		ColumnsInfo:    &protoV2.ColumnsInfo{Data: ConvertTabularDataToProto(meta.ColumnsInfo.TabularData)},
		IndexesInfo:    &protoV2.IndexesInfo{Data: ConvertTabularDataToProto(meta.IndexesInfo.TabularData)},
		CreateTableSQL: meta.CreateTableSQL,
	}
}

func ConvertProtoTableMetaToDriver(meta *protoV2.TableMeta) *TableMeta {
	return &TableMeta{
		ColumnsInfo:    ColumnsInfo{TabularData: ConvertProtoTabularDataToDriver(meta.ColumnsInfo.Data)},
		IndexesInfo:    IndexesInfo{TabularData: ConvertProtoTabularDataToDriver(meta.IndexesInfo.Data)},
		CreateTableSQL: meta.CreateTableSQL,
	}
}

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
