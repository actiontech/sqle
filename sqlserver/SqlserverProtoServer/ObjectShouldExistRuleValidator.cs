using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class ObjectShouldExistRuleValidator : RuleValidator {
        public List<String> DatabaseNames;
        public List<String> Schemas;
        public List<String> TableNames;

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            List<SchemaObjectName> schemaObjectNames = new List<SchemaObjectName>();
            switch (statement) {
                // USE database
                case UseStatement useStatement:
                    DatabaseNames.Add(useStatement.DatabaseName.Value);
                    break;

                //  CREATE TABLE database.schema.table(col1 INT)
                case CreateTableStatement createTableStatement:
                    DatabaseNames = AddDatabaseName(DatabaseNames, context, createTableStatement.SchemaObjectName);
                    Schemas = AddSchemaName(Schemas, createTableStatement.SchemaObjectName);
                    TableNames = AddTableName(TableNames, context, createTableStatement.SchemaObjectName);
                    break;

                // ALTER TABLE database.schema.table ALTER COLUMN col1 INT NOT NULL
                case AlterTableStatement alertTableStatemet:
                    DatabaseNames = AddDatabaseName(DatabaseNames, context, alertTableStatemet.SchemaObjectName);
                    Schemas = AddSchemaName(Schemas, alertTableStatemet.SchemaObjectName);
                    TableNames = AddTableName(TableNames, context, alertTableStatemet.SchemaObjectName);
                    break;

                case SelectStatement selectStatement:
                    if (selectStatement.QueryExpression is QuerySpecification) {
                        schemaObjectNames = AddSchemaObjectNameFromQuerySpecification(schemaObjectNames, selectStatement.QueryExpression as QuerySpecification);
                        GetDatabaseAndSchemaAndTableNames(schemaObjectNames, context, DatabaseNames, Schemas, TableNames);
                    }
                    break;

                case InsertStatement insertStatement:
                    var insertSpec = insertStatement.InsertSpecification;
                    // INSERT INTO table1 VALUES(1)
                    schemaObjectNames = AddSchemaObjectNameFromTableReference(schemaObjectNames, insertSpec.Target);
                    // INSERT INTO table1 SELECT...
                    if (insertSpec.InsertSource is SelectInsertSource) {
                        var selectInsertSource = insertSpec.InsertSource as SelectInsertSource;
                        if (selectInsertSource.Select is QuerySpecification) {
                            schemaObjectNames = AddSchemaObjectNameFromQuerySpecification(schemaObjectNames, selectInsertSource.Select as QuerySpecification);
                        }
                    }
                    GetDatabaseAndSchemaAndTableNames(schemaObjectNames, context, DatabaseNames, Schemas, TableNames);
                    break;

                case DeleteStatement deleteStatement:
                    var deleteSpec = deleteStatement.DeleteSpecification;
                    // DELETE FROM table1;
                    schemaObjectNames = AddSchemaObjectNameFromTableReference(schemaObjectNames, deleteSpec.Target);
                    // DELETE FROM schema1.table1 WHERE col1 IN (SELECT tbl2.col1 FROM schema2.table2 AS tbl2 INNER JOIN table3 AS tbl3 ON tbl2.col2=tbl3.col2)
                    WhereClause whereClause = deleteSpec.WhereClause;
                    if (whereClause != null) {
                        if (whereClause.SearchCondition is InPredicate) {
                            InPredicate inPredicate = whereClause.SearchCondition as InPredicate;
                            if (inPredicate.Subquery != null) {
                                schemaObjectNames = AddSchemaObjectNameFromQuerySpecification(schemaObjectNames, inPredicate.Subquery.QueryExpression as QuerySpecification);
                            }
                        }
                    }
                    GetDatabaseAndSchemaAndTableNames(schemaObjectNames, context, DatabaseNames, Schemas, TableNames);
                    break;

                case UpdateStatement updateStatement:
                    var updateSpec = updateStatement.UpdateSpecification;
                    // UPDATE schema1.table1 SET schema1.table1.col1 = 1
                    schemaObjectNames = AddSchemaObjectNameFromTableReference(schemaObjectNames, updateSpec.Target);
                    // UPDATE schema1.table2 SET schema1.table2.col2 = schema1.table2.col2 + schema1.table1.col2 FROM table2 INNER JOIN table1 ON (table2.col1 = table1.col1)
                    schemaObjectNames = AddSchemaObjectNameFromFromClause(schemaObjectNames, updateSpec.FromClause);

                    GetDatabaseAndSchemaAndTableNames(schemaObjectNames, context, DatabaseNames, Schemas, TableNames);
                    break;
            }
        }

        public void reset() {
            DatabaseNames = new List<string>();
            Schemas = new List<string>();
            TableNames = new List<string>();
        }

        public ObjectShouldExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) {
            DatabaseNames = new List<string>();
            Schemas = new List<string>();
            TableNames = new List<string>();
        }
    }

    public class DatabaseShouldExistRuleValidator : ObjectShouldExistRuleValidator {
        public override void Check(SqlserverContext context, TSqlStatement statement) {
            base.Check(context, statement);

            List<String> notExistDatabaseNames = new List<String>();
            foreach (var databaseName in DatabaseNames) {
                var databaseExisted = DatabaseExists(context, databaseName);
                if (!databaseExisted) {
                    notExistDatabaseNames.Add(databaseName);
                }
            }
            if (notExistDatabaseNames.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(String.Join(',', notExistDatabaseNames)));
            }

            List<String> notExistSchemas = new List<String>();
            foreach (var schema in Schemas) {
                var schemaExited = SchemaExists(context, schema);
                if (!schemaExited) {
                    notExistSchemas.Add(schema);
                }
            }
            if (notExistSchemas.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(String.Join(',', notExistSchemas)));
            }

            reset();
        }

        public DatabaseShouldExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class TableShouldExistRuleValidator : ObjectShouldExistRuleValidator {
        public override void Check(SqlserverContext context, TSqlStatement statement) {
            base.Check(context, statement);

            List<String> notExistTableNames = new List<String>();
            foreach (var tableName in TableNames) {
                var tableIdentifier = tableName.Split('.');
                if (tableIdentifier.Length != 3) {
                    continue;
                }
                var exist = TableExists(context, tableIdentifier[0], tableIdentifier[1], tableIdentifier[2]);
                if (!exist) {
                    notExistTableNames.Add(tableName);
                }
            }

            if (notExistTableNames.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(String.Join(',', notExistTableNames)));
            }

            reset();
        }

        public TableShouldExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
