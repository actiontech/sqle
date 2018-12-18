using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Data.SqlClient;
using System.Collections.Generic;
using System.Diagnostics.Contracts;

namespace SqlserverProtoServer {
    public class RollbackSql {
        public String GetRollbackSql(SqlserverContext context, TSqlStatement statement, out bool isDDL, out bool isDML) {
            isDDL = false;
            isDML = false;

            switch (statement) {
                // ddl
                case CreateDatabaseStatement createDatabaseStatement:
                    isDDL = true;
                    return GenerateCreateDatabaseRollbackSql(context, createDatabaseStatement);

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
                    break;

                case DeleteStatement deleteStatement:
                    isDML = true;
                    break;

                case UpdateStatement updateStatement:
                    isDML = true;
                    break;
            }

            return "";
        }

        public String GenerateCreateDatabaseRollbackSql(SqlserverContext context, CreateDatabaseStatement statement) {
            return String.Format("DROP DATABASE {0};", statement.DatabaseName.Value);
        }

        public String GenerateCreateTableRollbackSql(SqlserverContext context, CreateTableStatement statement) {
            String databaseName, schemaName, tableName;
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(statement.SchemaObjectName, out databaseName, out schemaName, out tableName);
            return String.Format("DROP TABLE {0}.{1}.{2};", databaseName, schemaName, tableName);
        }

        public String GetCreateTableSql(SqlserverContext context, String databaseName, String schemaName, String tableName) {
            String connectionString = context.GetConnectionString();
            var sqlLines = new List<String>();

            // columns
            var columnNameToDefinition = context.GetTableColumnDefinitions(databaseName, schemaName, tableName);
            foreach (var columnsPair in columnNameToDefinition) {
                String columnDefinition = columnsPair.Value;
                sqlLines.Add(columnDefinition);
            }

            // primary key & unique contraint
            var constraintNameToDefinition = context.GetTableConstraintDefinitions(databaseName, schemaName, tableName);
            foreach (var constraintPair in constraintNameToDefinition) {
                String constraintDefinition = constraintPair.Value;
                sqlLines.Add(constraintDefinition);
            }

            // index
            var indexNameToDefinition = context.GetTableIndexDefinitions(databaseName, schemaName, tableName);
            foreach (var indexPair in indexNameToDefinition) {
                String indexDefinition = indexPair.Value;
                sqlLines.Add(indexDefinition);
            }

            if (sqlLines.Count > 0) {
                return String.Format("CREATE TABLE {0}.{1}.{2} ({3});", databaseName, schemaName, tableName, String.Join(',', sqlLines));
            }

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
            String databaseName, schemaName, tableName;
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(statement.SchemaObjectName, out databaseName, out schemaName, out tableName);
            var key = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            var rollbackPrefix = String.Format("ALTER TABLE {0}.{1}.{2}", databaseName, schemaName, tableName);
            var rollbackActions = new List<String>();
            switch (statement) {
                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    var tableDefinition = alterTableAddTableElementStatement.Definition;

                    var addConstraints = new List<String>();
                    foreach (var tableConstraint in tableDefinition.TableConstraints) {
                        addConstraints.Add(tableConstraint.ConstraintIdentifier.Value);
                    }
                    if (addConstraints.Count > 0) {
                        rollbackActions.Add(String.Format("{0} DROP CONSTRAINT {1}", rollbackPrefix, String.Join(',', addConstraints)));
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
                            foreach (var tableColumnDefinitionPair in context.GetTableColumnDefinitions(databaseName, schemaName, tableName)) {
                                if (tableColumnDefinitionPair.Key == elemName) {
                                    dropColumnDefinitions.Add(tableColumnDefinitionPair.Value);
                                }
                            }
                            lastType = TableElementType.Column;
                        }

                        if (elemType == TableElementType.Constraint) {
                            foreach (var tableConstraintDefinitionPair in context.GetTableConstraintDefinitions(databaseName, schemaName, tableName)) {
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
                    foreach (var tableColumnDefinitionPair in context.GetTableColumnDefinitions(databaseName, schemaName, tableName)) {
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
                SqlCommand command = new SqlCommand(String.Format("SELECT a.type_desc AS Index_type, a.is_unique AS Index_unique, COL_NAME(b.object_id, b.column_id) AS Col_name FROM sys.indexes a JOIN sys.index_columns b ON a.object_id=b.object_id AND a.index_id =b.index_id WHERE a.object_id=OBJECT_ID('{0}') AND a.name='{1}';", tableName, indexName), connection);
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
            String databaseName, schemaName, tableName;
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(statement.OnName, out databaseName, out schemaName, out tableName);

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
                    String databaseName, schemaName, tableName;
                    context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(dropIndexClause.Object, out databaseName, out schemaName, out tableName);

                    var rollbackIndexSql = GetCreateIndexSql(context, indexName, databaseName, schemaName, tableName);
                    if (rollbackIndexSql != "") {
                        rollbackSql += String.Format("{0};", rollbackIndexSql);
                    }
                }
            }
            return rollbackSql;
        }

        public String GenerateInsertRollbackSql(SqlserverContext context, InsertStatement statement) {
            return "";
        }

        public String GenerateDeleteRollbackSql(SqlserverContext context, DeleteStatement statement) {
            return "";
        }

        public String GenerateUpdateRollbackSql(SqlserverContext context, UpdateStatement statement) {
            return "";
        }
    }
}