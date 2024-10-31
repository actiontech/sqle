//go:build enterprise
// +build enterprise

package differ

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/jmoiron/sqlx"
)

// Schemas returns a slice of schemas on the instance visible to the user. If
// called with no args, all non-system schemas will be returned. Or pass one or
// more schema names as args to filter the result to just those schemas.
// Note that the ordering of the resulting slice is not guaranteed.
func Schemas(ignoreSysDatabase bool, conn *executor.Executor, schemaInfos []*driverV2.DatabasSchemaInfo) ([]*Schema, error) {

	type rawSchema struct {
		Name      string `json:"schema_name"`
		CharSet   string `json:"default_character_set_name"`
		Collation string `json:"default_collation_name"`
	}

	var query string

	schemaNames := make([]string, len(schemaInfos))
	for i, schemaName := range schemaInfos {
		schemaNames[i] = schemaName.SchemaName
	}

	// Note on these queries: MySQL 8.0 changes information_schema column names to
	// come back from queries in all caps, so we need to explicitly use AS clauses
	// in order to get them back as lowercase and have sqlx Select() work
	if len(schemaInfos) == 0 {
		query = `
			SELECT schema_name AS schema_name, default_character_set_name AS default_character_set_name,
			       default_collation_name AS default_collation_name
			FROM   information_schema.schemata
			WHERE  schema_name NOT IN ('information_schema', 'performance_schema', 'mysql', 'test', 'sys')`
	} else {
		// If instance is using lower_case_table_names=2, apply an explicit collation
		// to ensure the schema name comes back with its original lettercasing. See
		// https://dev.mysql.com/doc/refman/8.0/en/charset-collation-information-schema.html
		var lctn2Collation string
		if ignoreSysDatabase {
			lctn2Collation = " COLLATE utf8_general_ci"
		}
		query = fmt.Sprintf(`
			SELECT schema_name AS schema_name, default_character_set_name AS default_character_set_name,
			       default_collation_name AS default_collation_name
			FROM   information_schema.schemata
			WHERE  schema_name%s IN (?)`, lctn2Collation)
	}
	query, args, err := sqlx.In(query, schemaNames)
	if err != nil {
		return nil, err
	}
	results, err := conn.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	ret := make([]*rawSchema, len(results))
	for i, record := range results {
		ret[i] = &rawSchema{
			Name:      record["schema_name"].String,
			CharSet:   record["default_character_set_name"].String,
			Collation: record["default_collation_name"].String,
		}
	}

	schemas := make([]*Schema, len(ret))
	for n, rawSchema := range ret {
		schemas[n] = &Schema{
			Name:      rawSchema.Name,
			CharSet:   rawSchema.CharSet,
			Collation: rawSchema.Collation,
		}
		// Create a non-cached connection pool with this schema as the default
		// database. The instance.querySchemaX calls below can establish a lot of
		// connections, so we will explicitly close the pool afterwards, to avoid
		// keeping a very large number of conns open. (Although idle conns eventually
		// get closed automatically, this may take too long.)
		results, err := conn.Db.Query("SELECT VERSION()")
		if err != nil {
			return nil, err
		}
		// 当前仅获取mysql
		flavor := ParseFlavor(fmt.Sprintf("%s %s", "mysql:", results[0]["VERSION()"].String))
		tableNames := getObjetNamesBySchema(rawSchema.Name, driverV2.ObjectType_TABLE, schemaInfos)
		if len(tableNames) > 0 {
			schemas[n].Tables, err = querySchemaTables(conn, rawSchema.Name, tableNames, flavor)
			if err != nil {
				return nil, err
			}
		}
		procedureNames := getObjetNamesBySchema(rawSchema.Name, driverV2.ObjectType_PROCEDURE, schemaInfos)
		if len(procedureNames) > 0 {
			schemas[n].Routines, err = querySchemaRoutines(conn, rawSchema.Name, procedureNames, driverV2.ObjectType_PROCEDURE, flavor)
			if err != nil {
				return nil, err
			}
		}
		functionNames := getObjetNamesBySchema(rawSchema.Name, driverV2.ObjectType_FUNCTION, schemaInfos)
		if len(functionNames) > 0 {
			schemas[n].Routines, err = querySchemaRoutines(conn, rawSchema.Name, functionNames, driverV2.ObjectType_FUNCTION, flavor)
			if err != nil {
				return nil, err
			}
		}

	}
	return schemas, nil
}

func getObjetNamesBySchema(schemaName string, objectType string, schemaInfos []*driverV2.DatabasSchemaInfo) []string {
	objectNames := make([]string, 0)
	for _, schemaInfo := range schemaInfos {
		if schemaInfo.SchemaName == schemaName {
			for _, obj := range schemaInfo.DatabaseObjects {
				if obj.ObjectType == objectType {
					objectNames = append(objectNames, obj.ObjectName)
				}
			}

		}
	}
	return objectNames
}
