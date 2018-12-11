using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Data.SqlClient;
using System.Collections.Generic;

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
                    break;

                case AlterTableStatement alterTableStatement:
                    isDDL = true;
                    break;

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
            return String.Format("DROP DATABASE {0}", statement.DatabaseName.Value);
        }

        public void GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(SqlserverContext context, SchemaObjectName schemaObjectName, out String databaseName, out String schemaName, out String tableName) {
            var databaseIdentifier = schemaObjectName.DatabaseIdentifier;
            if (databaseIdentifier != null) {
                databaseName = databaseIdentifier.Value;
            } else {
                databaseName = context.GetCurrentDatabase();
            }

            var schemaIdentifier = schemaObjectName.SchemaIdentifier;
            if (schemaIdentifier != null) {
                schemaName = schemaIdentifier.Value;
            } else {
                schemaName = context.GetCurrentSchema();
            }

            tableName = schemaObjectName.BaseIdentifier.Value;
            return;
        }

        public String GenerateCreateTableRollbackSql(SqlserverContext context, CreateTableStatement statement) {
            String databaseName, schemaName, tableName;
            GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(context, statement.SchemaObjectName, out databaseName, out schemaName, out tableName);
            return String.Format("DROP TABLE {0}.{1}.{2}", databaseName, schemaName, tableName);
        }

        public String GenerateDropTableRollbackSql(SqlserverContext context, DropTableStatement statement) {
            //var rollbackSql = "";
            foreach (var tableObject in statement.Objects) {
                // todo
            }
            return "";
        }

        public String GenerateAlterTableRollbackSql(SqlserverContext context, AlterTableStatement statement) {
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
                return String.Format("CREATE {0} {1} INDEX {2} ON {3}.{4}.{5} ({6})", unique, type, indexName, databaseName, schemaName, tableName, String.Join(',', columns));
            }

            return "";
        }

        public String GenerateCreateIndexRollbackSql(SqlserverContext context, CreateIndexStatement statement) {
            var indexName = statement.Name.Value;
            String databaseName, schemaName, tableName;
            GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(context, statement.OnName, out databaseName, out schemaName, out tableName);

            // if create index using DropExisting option, just rollback it using original create index sql
            bool dropExisting = false;
            foreach(var option in statement.IndexOptions) {
                if (option is IndexStateOption) {
                    var indexStateOption = option as IndexStateOption;
                    if (indexStateOption.OptionKind == IndexOptionKind.DropExisting) {
                        dropExisting = true;
                    }
                }
            }
            if (dropExisting) {
                return String.Format("{0} WITH (DROP_EXISTING=ON)", GetCreateIndexSql(context, indexName, databaseName, schemaName, tableName));
            }

            return String.Format("DROP INDEX {0} ON {1}.{2}.{3}", indexName, databaseName, schemaName, tableName);
        }

        public String GenerateDropIndexRollbackSql(SqlserverContext context, DropIndexStatement statement) {
            var rollbackSql = "";
            foreach (var dropIndexClauseBase in statement.DropIndexClauses) {
                if (dropIndexClauseBase is DropIndexClause) {
                    var dropIndexClause = dropIndexClauseBase as DropIndexClause;
                    var indexName = dropIndexClause.Index.Value;
                    String databaseName, schemaName, tableName;
                    GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(context, dropIndexClause.Object, out databaseName, out schemaName, out tableName);

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