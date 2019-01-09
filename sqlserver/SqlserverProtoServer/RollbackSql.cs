using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Data.SqlClient;
using System.Collections.Generic;
using System.Diagnostics.Contracts;
using NLog;

namespace SqlserverProtoServer {
    public class RollbackSql {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public String GetRollbackSql(SqlserverContext context, TSqlStatement statement, out bool isDDL, out bool isDML) {
            isDDL = false;
            isDML = false;

            switch (statement) {
                // ddl
                case CreateTableStatement createTableStatement:
                    isDDL = true;
                    return GenerateCreateTableRollbackSql(context, createTableStatement);

                case DropTableStatement dropTableStatement:
                    isDDL = true;
                    return GenerateDropTableRollbackSql(context, dropTableStatement);

                case AlterTableStatement alterTableStatement:
                    isDDL = true;
                    return GenerateAlterTableRollbackSql(context, alterTableStatement);

                case ExecuteStatement executeStatement:
                    isDDL = true;
                    return GenerateRenameRollbackSql(context, executeStatement);

                case CreateIndexStatement createIndexStatement:
                    isDDL = true;
                    return GenerateCreateIndexRollbackSql(context, createIndexStatement);

                case DropIndexStatement dropIndexStatement:
                    isDDL = true;
                    return GenerateDropIndexRollbackSql(context, dropIndexStatement);

                // dml
                case InsertStatement insertStatement:
                    isDML = true;
                    return GenerateInsertRollbackSql(context, insertStatement);

                case DeleteStatement deleteStatement:
                    isDML = true;
                    return GenerateDeleteRollbackSql(context, deleteStatement);

                case UpdateStatement updateStatement:
                    isDML = true;
                    return GenerateUpdateRollbackSql(context, updateStatement);
            }

            return "";
        }

        public String GenerateCreateTableRollbackSql(SqlserverContext context, CreateTableStatement statement) {
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(statement.SchemaObjectName, out String databaseName, out String schemaName, out String tableName);
            return String.Format("DROP TABLE {0}.{1}.{2};", databaseName, schemaName, tableName);
        }

        public String GetCreateTableSql(SqlserverContext context, String databaseName, String schemaName, String tableName) {
            var sqlLines = new List<String>();

            // columns
            var columnNameToDefinition = context.GetTableColumnDefinitions(logger, databaseName, schemaName, tableName);
            foreach (var columnsPair in columnNameToDefinition) {
                String columnDefinition = columnsPair.Value;
                sqlLines.Add(columnDefinition);
            }

            // primary key & unique contraint
            var constraintNameToDefinition = context.GetTableConstraintDefinitions(logger, databaseName, schemaName, tableName);
            foreach (var constraintPair in constraintNameToDefinition) {
                String constraintDefinition = constraintPair.Value;
                sqlLines.Add(constraintDefinition);
            }

            // index
            var indexNameToDefinition = context.GetTableIndexDefinitions(logger, databaseName, schemaName, tableName);
            foreach (var indexPair in indexNameToDefinition) {
                String indexDefinition = indexPair.Value;
                sqlLines.Add(indexDefinition);
            }

            if (sqlLines.Count > 0) {
                return String.Format("CREATE TABLE {0}.{1}.{2} ({3});", databaseName, schemaName, tableName, String.Join(',', sqlLines));
            }
            logger.Info("table {0}.{1}.{2} definition not found", databaseName, schemaName, tableName);
            return "";
        }

        public String GenerateDropTableRollbackSql(SqlserverContext context, DropTableStatement statement) {
            var rollbackSql = "";
            foreach (var tableObject in statement.Objects) {
                String databaseName, schemaName, tableName;
                if (tableObject.DatabaseIdentifier != null) {
                    databaseName = tableObject.DatabaseIdentifier.Value;
                } else {
                    databaseName = context.GetCurrentDatabase();
                }
                if (tableObject.SchemaIdentifier != null) {
                    schemaName = tableObject.SchemaIdentifier.Value;
                } else {
                    schemaName = context.GetCurrentSchema();
                }
                tableName = tableObject.BaseIdentifier.Value;

                rollbackSql += GetCreateTableSql(context, databaseName, schemaName, tableName);
            }
            return rollbackSql;
        }

        public String GenerateAlterTableRollbackSql(SqlserverContext context, AlterTableStatement statement) {
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(statement.SchemaObjectName, out String databaseName, out String schemaName, out String tableName);
            var key = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            var rollbackPrefix = String.Format("ALTER TABLE {0}.{1}.{2}", databaseName, schemaName, tableName);
            var rollbackActions = new List<String>();
            switch (statement) {
                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    var tableDefinition = alterTableAddTableElementStatement.Definition;

                    var addConstraints = new List<String>();
                    if (tableDefinition.TableConstraints != null) {
                        foreach (var tableConstraint in tableDefinition.TableConstraints) {
                            if (tableConstraint.ConstraintIdentifier == null) {
                                continue;
                            }
                            addConstraints.Add(tableConstraint.ConstraintIdentifier.Value);
                        }
                        if (addConstraints.Count > 0) {
                            rollbackActions.Add(String.Format("{0} DROP CONSTRAINT {1}", rollbackPrefix, String.Join(',', addConstraints)));
                        }
                    }

                    var addColumns = new List<String>();
                    foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                        addColumns.Add(columnDefinition.ColumnIdentifier.Value);
                    }
                    if (addColumns.Count > 0) {
                        rollbackActions.Add(String.Format("{0} DROP COLUMN {1}", rollbackPrefix, String.Join(',', addColumns)));
                    }
                    break;

                case AlterTableDropTableElementStatement alterTableDropTableElementStatement:
                    var alterTableDropTableElements = alterTableDropTableElementStatement.AlterTableDropTableElements;
                    var lastType = TableElementType.NotSpecified;
                    var dropColumnDefinitions = new List<String>();
                    var dropConstaintDefinitions = new List<String>();
                    foreach (var elem in alterTableDropTableElements) {
                        var elemName = elem.Name.Value;
                        var elemType = elem.TableElementType;
                        // in "ALTER TABLE dbo.test DROP CONSTRAINT my_constraint, my_pk_constraint, COLUMN column_b",
                        // my_pk_constraint's type will be NotSpecified
                        if (elemType == TableElementType.NotSpecified) {
                            elemType = lastType;
                        }

                        if (elemType == TableElementType.Column) {
                            foreach (var tableColumnDefinitionPair in context.GetTableColumnDefinitions(logger, databaseName, schemaName, tableName)) {
                                if (tableColumnDefinitionPair.Key == elemName) {
                                    dropColumnDefinitions.Add(tableColumnDefinitionPair.Value);
                                }
                            }
                            lastType = TableElementType.Column;
                        }

                        if (elemType == TableElementType.Constraint) {
                            foreach (var tableConstraintDefinitionPair in context.GetTableConstraintDefinitions(logger, databaseName, schemaName, tableName)) {
                                if (tableConstraintDefinitionPair.Key == elemName) {
                                    dropConstaintDefinitions.Add(tableConstraintDefinitionPair.Value);
                                }
                            }
                            lastType = TableElementType.Constraint;
                        }
                    }

                    if (dropColumnDefinitions.Count > 0) {
                        rollbackActions.Add(String.Format("{0} ADD {1}", rollbackPrefix, String.Join(',', dropColumnDefinitions)));
                    }
                    if (dropConstaintDefinitions.Count > 0) {
                        rollbackActions.Add(String.Format("{0} ADD {1}", rollbackPrefix, String.Join(',', dropConstaintDefinitions)));
                    }
                    break;

                case AlterTableAlterColumnStatement alterTableAlterColumnStatement:
                    var alterColumn = alterTableAlterColumnStatement.ColumnIdentifier.Value;
                    foreach (var tableColumnDefinitionPair in context.GetTableColumnDefinitions(logger, databaseName, schemaName, tableName)) {
                        if (tableColumnDefinitionPair.Key == alterColumn) {
                            rollbackActions.Add(String.Format("{0} ALTER COLUMN {1}", rollbackPrefix, tableColumnDefinitionPair.Value));
                        }
                    }
                    break;
            }
            return String.Join(';', rollbackActions);
        }

        public String GenerateRenameRollbackSql(SqlserverContext context, ExecuteStatement statement) {
            var rollbackSql = "";
            var entity = statement.ExecuteSpecification.ExecutableEntity;
            if (entity is ExecutableProcedureReference) {
                var procedure = entity as ExecutableProcedureReference;
                if (procedure.ProcedureReference is ProcedureReferenceName) {
                    var procedureName = procedure.ProcedureReference as ProcedureReferenceName;
                    var baseName = procedureName.ProcedureReference.Name.BaseIdentifier.Value;
                    if (baseName.ToLower() == "sp_rename") {
                        // EXEC sp_rename N'dbo.test.IX_index0', N'IX_index1', N'INDEX' => EXEC sp_renam N'dbo.test.IX_index1', N'IX_index0', N'INDEX' 
                        if (procedure.Parameters.Count >= 2) {
                            var objName = procedure.Parameters[0].ParameterValue as StringLiteral;
                            var newName = procedure.Parameters[1].ParameterValue as StringLiteral;

                            var rollbackNewName = "";
                            var objNameArray = objName.Value.Split(".");
                            if (newName.IsNational) {
                                rollbackNewName = String.Format("N'{0}'", objNameArray[objNameArray.Length - 1]);
                            } else {
                                rollbackNewName = String.Format("'{0}'", objNameArray[objNameArray.Length - 1]);
                            }

                            var rollbackObjName = "";
                            objNameArray[objNameArray.Length - 1] = newName.Value;
                            if (objName.IsNational) {
                                rollbackObjName = String.Format("N'{0}'", String.Join('.', objNameArray));
                            } else {
                                rollbackObjName = String.Format("'{0}'", String.Join('.', objNameArray));
                            }


                            rollbackSql = String.Format("EXEC sp_rename {0}, {1}", rollbackObjName, rollbackNewName);
                            for (int i = 2; i < procedure.Parameters.Count; i++) {
                                var param = procedure.Parameters[i].ParameterValue as StringLiteral;
                                if (param.IsNational) {
                                    rollbackSql += String.Format(", N'{0}'", param.Value);
                                } else {
                                    rollbackSql += String.Format(", '{0}'", param.Value);
                                }
                            }
                        }
                    }
                }
            }

            return rollbackSql;
        }

        public String GetCreateIndexSql(SqlserverContext context, String indexName, String databaseName, String schemaName, String tableName) {
            String type = "";
            String unique = "";
            List<String> columns = new List<string>();
            String connectionString = context.GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand(String.Format("SELECT a.type_desc AS Index_type, a.is_unique AS Index_unique, COL_NAME(b.object_id, b.column_id) AS Col_name FROM {0}.sys.indexes a JOIN {0}.sys.index_columns b ON a.object_id=b.object_id AND a.index_id =b.index_id WHERE a.object_id=OBJECT_ID('{1}') AND a.name='{2}';", databaseName, tableName, indexName), connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        type = (String)reader["Index_type"];
                        unique = ((bool)reader["Index_unique"] ? "UNIQUE" : "");
                        columns.Add((String)reader["Col_name"]);
                    }
                } finally {
                    reader.Close();
                }
            }

            if (columns.Count > 0) {
                return String.Format("CREATE {0} {1} INDEX {2} ON {3}.{4}.{5} ({6});", unique, type, indexName, databaseName, schemaName, tableName, String.Join(',', columns));
            }

            return "";
        }

        public String GenerateCreateIndexRollbackSql(SqlserverContext context, CreateIndexStatement statement) {
            var indexName = statement.Name.Value;
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(statement.OnName, out String databaseName, out String schemaName, out String tableName);

            // if create index using DropExisting option, just rollback it using original create index sql
            bool dropExisting = false;
            foreach (var option in statement.IndexOptions) {
                if (option is IndexStateOption) {
                    var indexStateOption = option as IndexStateOption;
                    if (indexStateOption.OptionKind == IndexOptionKind.DropExisting) {
                        dropExisting = true;
                    }
                }
            }
            var rollbackSql = "";
            if (dropExisting) {
                rollbackSql = GetCreateIndexSql(context, indexName, databaseName, schemaName, tableName);
                if (rollbackSql != "") {
                    return String.Format("{0} WITH (DROP_EXISTING=ON);", rollbackSql);
                }
            }

            return String.Format("DROP INDEX {0} ON {1}.{2}.{3};", indexName, databaseName, schemaName, tableName);
        }

        public String GenerateDropIndexRollbackSql(SqlserverContext context, DropIndexStatement statement) {
            var rollbackSql = "";
            foreach (var dropIndexClauseBase in statement.DropIndexClauses) {
                if (dropIndexClauseBase is DropIndexClause) {
                    var dropIndexClause = dropIndexClauseBase as DropIndexClause;
                    var indexName = dropIndexClause.Index.Value;
                    context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(dropIndexClause.Object, out String databaseName, out String schemaName, out String tableName);

                    var rollbackIndexSql = GetCreateIndexSql(context, indexName, databaseName, schemaName, tableName);
                    if (rollbackIndexSql != "") {
                        rollbackSql += String.Format("{0};\n", rollbackIndexSql);
                    }
                }
            }
            return rollbackSql;
        }

        public String GenerateInsertRollbackSql(SqlserverContext context, InsertStatement statement) {
            var rollbackSqls = new List<String>();
            var insertSpecification = statement.InsertSpecification;
            if (insertSpecification.InsertSource is ValuesInsertSource) {
                var tableReference = insertSpecification.Target as NamedTableReference;
                context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(tableReference.SchemaObject, out String databaseName, out String schemaName, out String tableName);

                var primaryKeys = context.GetPrimaryKeys(databaseName, schemaName, tableName);
                if (primaryKeys.Count == 0) {
                    return "";
                }
                var columns = new List<String>();
                if (insertSpecification.Columns.Count > 0) {
                    foreach (var column in insertSpecification.Columns) {
                        foreach (var identifier in column.MultiPartIdentifier.Identifiers) {
                            columns.Add(identifier.Value);
                        }
                    }
                } else {
                    columns = context.GetColumns(databaseName, schemaName, tableName);
                }

                var insertSource = insertSpecification.InsertSource as ValuesInsertSource;
                if (insertSource.RowValues.Count > context.Config.DMLRollbackMaxRows) {
                    return "";
                }
                foreach (var rowValue in insertSource.RowValues) {
                    if (rowValue.ColumnValues.Count != columns.Count) {
                        return "";
                    }

                    var whereCondition = new List<String>();
                    for (int index = 0; index < columns.Count; index++) {
                        var columnName = columns[index];
                        foreach (var primaryKey in primaryKeys) {
                            if (columnName == primaryKey) {
                                var value = "";
                                switch (rowValue.ColumnValues[index]) {
                                    case IntegerLiteral integerLiteral:
                                        value = integerLiteral.Value;
                                        break;
                                    case StringLiteral stringLiteral:
                                        value = stringLiteral.Value;
                                        break;
                                    case NumericLiteral numericLiteral:
                                        value = numericLiteral.Value;
                                        break;
                                }
                                if (value != "") {
                                    whereCondition.Add(String.Format("{0} = '{1}'", columnName, value));
                                }
                            }
                        }
                    }
                    if (whereCondition.Count != primaryKeys.Count) {
                        return "";
                    }
                    rollbackSqls.Add(String.Format("DELETE FROM {0}.{1}.{2} WHERE {3};", databaseName, schemaName, tableName, String.Join(" AND ", whereCondition)));
                }
            }

            return String.Join('\n', rollbackSqls);
        }

        public String GenerateDeleteRollbackSql(SqlserverContext context, DeleteStatement statement) {
            var rollbackSql = "";
            var deleteSpecification = statement.DeleteSpecification;
            if (deleteSpecification.Target is NamedTableReference) {
                var tableReference = deleteSpecification.Target as NamedTableReference;
                context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(tableReference.SchemaObject, out String databaseName, out String schemaName, out String tableName);

                var where = "";
                if (deleteSpecification.WhereClause != null) {
                    var whereClause = deleteSpecification.WhereClause;
                    for (int index = whereClause.FirstTokenIndex; index <= whereClause.LastTokenIndex; index++) {
                        where += whereClause.ScriptTokenStream[index].Text;
                    }
                }

                var recordsCount = context.GetRecordsCount(databaseName, schemaName, tableName, where);
                if (recordsCount > context.Config.DMLRollbackMaxRows) {
                    return "";
                }
                var records = context.GetRecords(databaseName, schemaName, tableName, where);
                var columns = context.GetColumns(databaseName, schemaName, tableName);

                var values = new List<String>();
                foreach (var record in records) {
                    if (record.Count != columns.Count) {
                        return "";
                    }

                    var recordValues = new List<String>();
                    foreach (var columnName in columns) {
                        var value = "NULL";
                        if (record.ContainsKey(columnName) && record[columnName] != "") {
                            value = String.Format("'{0}'", record[columnName]);
                        }
                        recordValues.Add(value);
                    }
                    values.Add(String.Format("({0})", String.Join(", ", recordValues)));
                }

                if (values.Count > 0) {
                    rollbackSql = String.Format("INSERT INTO {0}.{1}.{2} ({3}) VALUES {4}", databaseName, schemaName, tableName, String.Join(", ", columns), String.Join(", ", values));
                }
            }
            return rollbackSql;
        }

        public String GenerateUpdateRollbackSql(SqlserverContext context, UpdateStatement statement) {
            var rollbacksql = "";
            var updateSpecification = statement.UpdateSpecification;
            if (updateSpecification.Target is NamedTableReference) {
                var tableReference = updateSpecification.Target as NamedTableReference;
                context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(tableReference.SchemaObject, out String databaseName, out String schemaName, out String tableName);

                var primaryKeys = context.GetPrimaryKeys(databaseName, schemaName, tableName);
                if (primaryKeys.Count == 0) {
                    return "";
                }

                var where = "";
                if (updateSpecification.WhereClause != null) {
                    var whereClause = updateSpecification.WhereClause;
                    for (int index = whereClause.FirstTokenIndex; index <= whereClause.LastTokenIndex; index++) {
                        where += whereClause.ScriptTokenStream[index].Text;
                    }
                }
                var recordsCount = context.GetRecordsCount(databaseName, schemaName, tableName, where);
                if (recordsCount > context.Config.DMLRollbackMaxRows) {
                    return "";
                }
                var records = context.GetRecords(databaseName, schemaName, tableName, where);
                var columns = context.GetColumns(databaseName, schemaName, tableName);

                foreach (var record in records) {
                    if (record.Count != columns.Count) {
                        return "";
                    }

                    var whereConditions = new List<String>();
                    var value = new List<String>();
                    foreach (var columnName in columns) {
                        var colChanged = false;
                        var isPkChanged = false;
                        var isPk = false;
                        foreach (var pkColumnName in primaryKeys) {
                            if (pkColumnName == columnName) {
                                isPk = true;
                                break;
                            }
                        }

                        var pkValue = "";
                        foreach (var setClause in updateSpecification.SetClauses) {
                            var updatedColumn = setClause.ScriptTokenStream[setClause.FirstTokenIndex].Text;
                            if (updatedColumn == columnName) {
                                colChanged = true;
                                if (isPk) {
                                    isPkChanged = true;
                                    pkValue = setClause.ScriptTokenStream[setClause.LastTokenIndex].Text;
                                }
                            }
                        }

                        var v = "NULL";
                        if (record.ContainsKey(columnName)) {
                            v = String.Format("'{0}'", record[columnName]);
                        }

                        if (colChanged) {
                            value.Add(String.Format("{0} = {1}", columnName, v));
                        }
                        if (isPk) {
                            if (isPkChanged) {
                                whereConditions.Add(String.Format("{0} = {1}", columnName, pkValue));
                            } else {
                                whereConditions.Add(String.Format("{0} = {1}", columnName, v));
                            }
                        }

                    }
                    rollbacksql += String.Format("UPDATE {0}.{1}.{2} SET {3} WHERE {4};", databaseName, schemaName, tableName, String.Join(", ", value), String.Join(" AND ", whereConditions));
                }
            }
            return rollbacksql;
        }
    }
}