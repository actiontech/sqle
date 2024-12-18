//go:build enterprise
// +build enterprise

package executor

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
)

// When using keywords as schema names, you need to pay attention to wrapping them in quotation marks
func (c *Executor) ShowCreateSchema(schemaName string) (string, error) {
	query := fmt.Sprintf("SHOW CREATE DATABASE %s", schemaName)

	result, err := c.Db.Query(query)
	if err != nil {
		return "", err
	}

	if len(result) != 1 {
		err := fmt.Errorf("show create database error, result is %v", result)
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	return result[0]["Create Database"].String, nil
}

func (c *Executor) GetTableEngine(schemaName string, tableName string) (string, error) {
	query := "SELECT TABLE_NAME, ENGINE FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"

	result, err := c.Db.Query(query, schemaName, tableName)
	if err != nil {
		return "", err
	}

	if len(result) != 1 {
		err := fmt.Errorf("get table engine, result is %v", result)
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	return result[0]["ENGINE"].String, nil
}

func (c *Executor) ShowSchemaProcedures(schema string) ([]string, error) {
	query := fmt.Sprintf("select ROUTINE_NAME from information_schema.routines where routine_schema = '%s' AND routine_type = 'PROCEDURE'", schema)

	if c.IsLowerCaseTableNames() {
		schema = strings.ToLower(schema)

		query = fmt.Sprintf(
			"select ROUTINE_NAME from information_schema.routines where lower(routine_schema)='%s' and routine_type = 'PROCEDURE'", schema)
	}

	result, err := c.Db.Query(query)
	if err != nil {
		return nil, err
	}
	procedures := make([]string, len(result))
	for n, v := range result {
		if len(v) < 1 {
			err := fmt.Errorf("show procedures error, result not match")
			c.Db.Logger().Error(err)
			return procedures, errors.New(errors.ConnectRemoteDatabaseError, err)
		}
		for key, table := range v {
			if key != "ROUTINE_NAME" {
				continue
			}
			procedures[n] = table.String
			break
		}
	}
	return procedures, nil
}

// When using keywords as procedure names, you need to pay attention to wrapping them in quotation marks
func (c *Executor) ShowCreateProcedure(procedureName string) (string, error) {
	result, err := c.Db.Query(fmt.Sprintf("show create procedure %s", procedureName))
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		err := fmt.Errorf("show create procedure error, result is %v", result)
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	if query, ok := result[0]["Create Procedure"]; !ok {
		err := fmt.Errorf("show create procedure error, column \"Create Procedure\" not found")
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	} else {
		return query.String, nil
	}
}

func (c *Executor) ShowSchemaFunctions(schema string) ([]string, error) {
	query := fmt.Sprintf("select ROUTINE_NAME from information_schema.routines where routine_schema = '%s' AND routine_type = 'FUNCTION'", schema)

	if c.IsLowerCaseTableNames() {
		schema = strings.ToLower(schema)

		query = fmt.Sprintf(
			"select ROUTINE_NAME from information_schema.routines where lower(routine_schema)='%s' and routine_type = 'FUNCTION'", schema)
	}

	result, err := c.Db.Query(query)
	if err != nil {
		return nil, err
	}
	procedures := make([]string, len(result))
	for n, v := range result {
		if len(v) < 1 {
			err := fmt.Errorf("show procedures error, result not match")
			c.Db.Logger().Error(err)
			return procedures, errors.New(errors.ConnectRemoteDatabaseError, err)
		}
		for key, procedure := range v {
			if key != "ROUTINE_NAME" {
				continue
			}
			procedures[n] = procedure.String
			break
		}
	}
	return procedures, nil
}

// When using keywords as function names, you need to pay attention to wrapping them in quotation marks
func (c *Executor) ShowCreateFunction(functionName string) (string, error) {
	result, err := c.Db.Query(fmt.Sprintf("show create function %s", functionName))
	if err != nil {
		return "", err
	}
	if len(result) != 1 {
		err := fmt.Errorf("show create function error, result is %v", result)
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	}
	if query, ok := result[0]["Create Function"]; !ok {
		err := fmt.Errorf("show create function error, column \"Create Function\" not found")
		c.Db.Logger().Error(err)
		return "", errors.New(errors.ConnectRemoteDatabaseError, err)
	} else {
		return query.String, nil
	}
}

// TODO 数据库结构对比目前视图、事件、触发器还未支持，若后续支持放开此处注释的代码以获取对象的名称（代码须经验证）

// func (c *Executor) ShowSchemaTriggers(schema string) ([]string, error) {
// 	query := fmt.Sprintf("select TRIGGER_NAME from information_schema.triggers where trigger_schema = '%s'", schema)
//
// 	if c.IsLowerCaseTableNames() {
// 		schema = strings.ToLower(schema)
//
// 		query = fmt.Sprintf(
// 			"select TRIGGER_NAME from information_schema.triggers where lower(trigger_schema)='%s'", schema)
// 	}
//
// 	result, err := c.Db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	triggers := make([]string, len(result))
// 	for n, v := range result {
// 		if len(v) < 1 {
// 			err := fmt.Errorf("show triggers error, result not match")
// 			c.Db.Logger().Error(err)
// 			return triggers, errors.New(errors.ConnectRemoteDatabaseError, err)
// 		}
// 		for key, trigger := range v {
// 			if key != "TRIGGER_NAME" {
// 				continue
// 			}
// 			triggers[n] = trigger.String
// 			break
// 		}
// 	}
// 	return triggers, nil
// }
//
// func (c *Executor) ShowSchemaEvents(schema string) ([]string, error) {
// 	query := fmt.Sprintf("select EVENT_NAME from information_schema.events where event_schema = '%s'", schema)
//
// 	if c.IsLowerCaseTableNames() {
// 		schema = strings.ToLower(schema)
//
// 		query = fmt.Sprintf(
// 			"select EVENT_NAME from information_schema.events where lower(event_schema)='%s'", schema)
// 	}
//
// 	result, err := c.Db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	events := make([]string, len(result))
// 	for n, v := range result {
// 		if len(v) < 1 {
// 			err := fmt.Errorf("show events error, result not match")
// 			c.Db.Logger().Error(err)
// 			return events, errors.New(errors.ConnectRemoteDatabaseError, err)
// 		}
// 		for key, event := range v {
// 			if key != "EVENT_NAME" {
// 				continue
// 			}
// 			events[n] = event.String
// 			break
// 		}
// 	}
// 	return events, nil
// }
