package compare

import (
	"context"

	"github.com/actiontech/sqle/sqle/common"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
)

type Compared struct {
	BaseInstance     *model.Instance
	ComparedInstance *model.Instance
}

type SchemaObject struct {
	SchemaName          string
	ComparedResult      string
	DatabaseDiffObjects []*DatabaseDiffObject
	InconsistentNum     int
}

type DatabaseDiffObject struct {
	InconsistentNum    int
	ObjectType         string
	ObjectsDiffResults []*ObjectDiffResult
}

type ObjectDiffResult struct {
	ComparedResult string
	ObjectName     string
}

func (c *Compared) ExecDatabaseCompare(baseSchemaName, comparedSchemaName *string) (*SchemaObject, error) {
	basePlugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), c.BaseInstance, "")
	if err != nil {
		return nil, err
	}
	defer basePlugin.Close(context.TODO())
	var baseInfos []*driverV2.DatabasSchemaInfo
	if baseSchemaName != nil && comparedSchemaName != nil {
		baseInfos = append(baseInfos, &driverV2.DatabasSchemaInfo{
			ScheamName: *baseSchemaName,
		})
	}
	baseRes, err := basePlugin.GetDatabaseObjectDDL(context.TODO(), baseInfos)
	if err != nil {
		return nil, err
	}
	comparedPlugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), c.ComparedInstance, "")
	if err != nil {
		return nil, err
	}
	defer comparedPlugin.Close(context.TODO())
	var infos []*driverV2.DatabasSchemaInfo
	if baseSchemaName != nil && comparedSchemaName != nil {
		infos = append(infos, &driverV2.DatabasSchemaInfo{
			ScheamName: *comparedSchemaName,
		})
	}
	comparedRes, err := basePlugin.GetDatabaseObjectDDL(context.TODO(), infos)
	if err != nil {
		return nil, err
	}
	for _, base := range baseRes {
		for _, compared := range comparedRes {
			if compared.SchemaName == base.SchemaName {
				// TODO 对比逻辑

			}
		}
	}
	return nil, err
}

func (c *Compared) GetDatabaseDiffModifySQLs(objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	basePlugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), c.BaseInstance, "")
	if err != nil {
		return nil, err
	}
	defer basePlugin.Close(context.TODO())
	calibratedDSN := &driverV2.DSN{
		Host:             c.ComparedInstance.Host,
		Port:             c.ComparedInstance.Port,
		User:             c.ComparedInstance.User,
		Password:         c.ComparedInstance.Password,
		AdditionalParams: c.ComparedInstance.AdditionalParams,
	}
	modifySQLs, err := basePlugin.GetDatabaseDiffModifySQL(context.TODO(), calibratedDSN, objInfos)
	if err != nil {
		return nil, err
	}
	dbDiffSQLs := make([]*driverV2.DatabaseDiffModifySQLResult, len(modifySQLs))
	for i, schemaDiff := range modifySQLs {
		dbDiffSQLs[i] = &driverV2.DatabaseDiffModifySQLResult{
			SchemaName: schemaDiff.SchemaName,
			ModifySQLs: schemaDiff.ModifySQLs,
		}
	}
	return dbDiffSQLs, nil
}
