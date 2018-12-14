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
            Console.WriteLine("{0}", statement);
            String databaseName, schemaName, tableName;
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(statement.SchemaObjectName, out databaseName, out schemaName, out tableName);
            var rollbackActions = new List<String>();
            switch (statement) {
                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    var tableDefinition = alterTableAddTableElementStatement.Definition;

                    var addConstraints = new List<String>();
                    foreach (var tableConstraint in tableDefinition.TableConstraints) {
                        addConstraints.Add(tableConstraint.ConstraintIdentifier.Value);
                    }
                    if (addConstraints.Count > 0) {
                        rollbackActions.Add("DROP CONSTRAINT " + String.Join(',', addConstraints));
                    }

                    var addColumns= new List<String>();
                    foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                        addColumns.Add(columnDefinition.ColumnIdentifier.Value);
                    }
                    if (addColumns.Count > 0) {
                        rollbackActions.Add("DROP COLUMN " + String.Join(',', addColumns));
                    }

                    Console.WriteLine("{0}", String.Join(';', rollbackActions));
                    break;

                case AlterTableDropTableElementStatement alterTableDropTableElementStatement:
                    var dropColumns = new List<String>();
                    var dropConstraints = new List<String>();
                    var alterTableDropTableElements = alterTableDropTableElementStatement.AlterTableDropTableElements;
                    foreach (var elem in alterTableDropTableElements) {
                        if (elem.TableElementType == TableElementType.Column) {
                            dropColumns.Add(elem.Name.Value);
                        }
                        if (elem.TableElementType == TableElementType.Constraint) {
                            dropConstraints.Add(elem.Name.Value);
                        }
                    }

                    break;
            }
            return "";
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
                        rollbackSql += String.Format("{0};", GetCreateIndexSql(context, indexName, databaseName, schemaName, tableName));
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