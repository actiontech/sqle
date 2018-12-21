using System;
using System.Collections.Generic;
using SqlserverProto;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Data.SqlClient;

namespace SqlserverProtoServer {
    public enum RULE_LEVEL {
        NORMAL, NOTICE, WARN, ERROR
    }

    public class RULE_LEVEL_STRING {
        public static Dictionary<RULE_LEVEL, String> RuleLevels = new Dictionary<RULE_LEVEL, String> {
            {RULE_LEVEL.NORMAL, "normal"},
            {RULE_LEVEL.NOTICE, "notice"},
            {RULE_LEVEL.WARN, "warn"},
            {RULE_LEVEL.ERROR, "error"}
        };

        public static String GetRuleLevelString(RULE_LEVEL level) {
            return RuleLevels[level];
        }
    }

    public static class DefaultRules {
        // rule names
        // This SCHEMA is DATABASE which comes from MySQL
        public const String SCHEMA_NOT_EXIST = "schema_not_exist";
        public const String SCHEMA_EXIST = "schema_exist";
        public const String TABLE_NOT_EXIST = "table_not_exist";
        public const String TABLE_EXIST = "table_exist";
        public const String DDL_CREATE_TABLE_NOT_EXIST = "ddl_create_table_not_exist";
        public const String DDL_CHECK_OBJECT_NAME_LENGTH = "ddl_check_object_name_length";
        public const String DDL_CHECK_PRIMARY_KEY_EXIST = "ddl_check_primary_key_exist";
        public const String DDL_CHECK_PRIMARY_KEY_TYPE = "ddl_check_primary_key_type";
        public const String DDL_DISABLE_VARCHAR_MAX = "ddl_disable_varchar_max";
        public const String DDL_CHECK_TYPE_CHAR_LENGTH = "ddl_check_type_char_length";
        public const String DDL_DISABLE_FOREIGN_KEY = "ddl_disable_foreign_key";
        public const String DDL_CHECK_INDEX_COUNT = "ddl_check_index_count";
        public const String DDL_CHECK_COMPOSITE_INDEX_MAX = "ddl_check_composite_index_max";
        public const String DDL_DISABLE_USING_KEYWORD = "ddl_disable_using_keyword";
        private const string DDL_TABLE_USING_INNODB_UTF8MB4 = "ddl_create_table_using_innodb";
        public const String DDL_DISABLE_INDEX_DATA_TYPE_BLOB = "ddl_disable_index_column_blob";
        public const String DDL_CHECK_ALTER_TABLE_NEED_MERGE = "ddl_check_alter_table_need_merge";
        public const String DDL_DISABLE_DROP_STATEMENT = "ddl_disable_drop_statement";
        public const String DML_CHECK_INVALID_WHERE_CONDITION = "ddl_check_invalid_where_condition";
        public const String DML_DISABE_SELECT_ALL_COLUMN = "dml_disable_select_all_column";

        // rules
        public static Dictionary<String, RuleValidator> RuleValidators = new Dictionary<String, RuleValidator> {
            {
                SCHEMA_NOT_EXIST,
                new DatabaseShouldExistRuleValidator(
                    SCHEMA_NOT_EXIST,
                    "操作数据库时，数据库必须存在",
                    "database或者schema {0} 不存在",
                    RULE_LEVEL.ERROR
                )
            },
            {
                SCHEMA_EXIST,
                new DatabaseShouldNotExistRuleValidator(
                    SCHEMA_EXIST,
                    "创建数据库时，数据库不能存在",
                    "database {0} 已存在",
                    RULE_LEVEL.ERROR
                )
            },
            {
                TABLE_NOT_EXIST,
                new TableShouldExistRuleValidator(
                    TABLE_NOT_EXIST,
                    "操作表时，表必须存在",
                    "表 {0} 不存在",
                    RULE_LEVEL.ERROR
                )
            },
            {
                TABLE_EXIST,
                new TableShouldNotExistRuleValidator(
                    TABLE_EXIST,
                    "创建表时，表不能存在",
                    "表 {0} 已存在",
                    RULE_LEVEL.ERROR
                )
            },
            {
                // There is no CREATE TABLE IF NOT EXISTS statement
                DDL_CREATE_TABLE_NOT_EXIST,
                new FakerRuleValidator(
                    DDL_CREATE_TABLE_NOT_EXIST,
                    "新建表必须加入if not exists create，保证重复执行不报错",
                    "新建表必须加入if not exists create，保证重复执行不报错",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_OBJECT_NAME_LENGTH,
                new ObjectNameMaxLengthRuleValidator(
                    DDL_CHECK_OBJECT_NAME_LENGTH,
                    "表名、列名、索引名的长度不能大于64字节",
                    "表名、列名、索引名的长度不能大于64字节",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_PRIMARY_KEY_EXIST,
                new PrimaryKeyShouldExistRuleValidator(
                    DDL_CHECK_PRIMARY_KEY_EXIST,
                    "表必须有主键",
                    "表必须有主键",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_PRIMARY_KEY_TYPE,
                new PrimaryKeyAutoIncrementRuleValidator(
                    DDL_CHECK_PRIMARY_KEY_TYPE,
                    "主键建议使用自增",
                    "主键建议使用自增",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_DISABLE_VARCHAR_MAX,
                new StringTypeShouldNoVarcharMaxRuleValidator(
                    DDL_DISABLE_VARCHAR_MAX,
                    "禁止使用 varchar(max)",
                    "禁止使用 varchar(max)",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_TYPE_CHAR_LENGTH,
                new StringTypeShouldNotExceedMaxLengthRuleValidator(
                    DDL_CHECK_TYPE_CHAR_LENGTH,
                    "char长度大于20时，必须使用varchar类型",
                    "char长度大于20时，必须使用varchar类型",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_DISABLE_FOREIGN_KEY,
                new ForeignKeyRuleValidator(
                    DDL_DISABLE_FOREIGN_KEY,
                    "禁止使用外键",
                    "禁止使用外键",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_INDEX_COUNT,
                new NumberOfIndexesShouldNotExceedMaxRuleValidator(
                    DDL_CHECK_INDEX_COUNT,
                    "索引个数建议不超过5个",
                    "索引个数建议不超过5个",
                    RULE_LEVEL.NOTICE
                )
            },
            {
                DDL_CHECK_COMPOSITE_INDEX_MAX,
                new NumberOfCompsiteIndexColumnsShouldNotExceedMaxRuleValidator(
                    DDL_CHECK_COMPOSITE_INDEX_MAX,
                    "复合索引的列数量不建议超过5个",
                    "复合索引的列数量不建议超过5个",
                    RULE_LEVEL.NOTICE
                )
            },
            {
                DDL_DISABLE_USING_KEYWORD,
                new ObjectNameRuleValidator(
                    DDL_DISABLE_USING_KEYWORD,
                    "数据库对象命名禁止使用关键字",
                    "数据库对象命名禁止使用关键字 %s",
                    RULE_LEVEL.ERROR
                )
            },

              {
               DDL_TABLE_USING_INNODB_UTF8MB4,
               new FakerRuleValidator(
                   DDL_TABLE_USING_INNODB_UTF8MB4,
                   "建议使用Innodb引擎,utf8mb4字符集",
                   "建议使用Innodb引擎,utf8mb4字符集",
                   RULE_LEVEL.NOTICE
               )
              },

            {
                DDL_DISABLE_INDEX_DATA_TYPE_BLOB,
                new DisableAddIndexForColumnsTypeBlob(
                    DDL_DISABLE_INDEX_DATA_TYPE_BLOB,
                    "禁止将blob类型的列加入索引",
                    "禁止将blob类型的列加入索引",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_ALTER_TABLE_NEED_MERGE,
                new MergeAlterTableRuleValidator(
                    DDL_CHECK_ALTER_TABLE_NEED_MERGE,
                    "存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
                    "已存在对该表的修改语句，建议合并成一个ALTER语句",
                    RULE_LEVEL.NOTICE
                )
            },
            {
                DDL_DISABLE_DROP_STATEMENT,
                new DisableDropRuleValidator(
                    DDL_DISABLE_DROP_STATEMENT,
                    "禁止除索引外的drop操作",
                    "禁止除索引外的drop操作",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DML_CHECK_INVALID_WHERE_CONDITION,
                new SelectWhereRuleValidator(
                    DML_CHECK_INVALID_WHERE_CONDITION,
                    "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
                    "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
                    RULE_LEVEL.ERROR)
            },
            {
                DML_DISABE_SELECT_ALL_COLUMN,
                new SelectAllRuleValidator(
                    DML_DISABE_SELECT_ALL_COLUMN,
                    "不建议使用select *",
                    "不建议使用select *",
                    RULE_LEVEL.NOTICE
                )
            }
        };
    }

    /// <summary>
    /// Advise result context represents context of all advise rule results for one sql.
    /// </summary>
    public class AdviseResultContext {
        public RULE_LEVEL Level;
        public String Message;

        public AdviseResultContext() {
            Level = RULE_LEVEL.NORMAL;
            Message = "";
        }

        public void AddAdviseResult(RULE_LEVEL level, String message) {
            if (Level < level) {
                Level = level;
            }

            var formatMsg = String.Format("[{0}]{1}", RULE_LEVEL_STRING.GetRuleLevelString(level), message);
            if (String.IsNullOrEmpty(Message)) {
                Message = formatMsg;
            } else {
                Message += "\n" + formatMsg;
            }
        }

        public void ResetAdviseResult() {
            Level = RULE_LEVEL.NORMAL;
            Message = "";
        }

        public AdviseResult GetAdviseResult() {
            AdviseResult adviseResult = new AdviseResult();
            adviseResult.AdviseLevel = GetLevel();
            adviseResult.AdviseResultMessage = GetMessage();
            return adviseResult;
        }

        public String GetLevel() {
            return RULE_LEVEL_STRING.GetRuleLevelString(Level);
        }

        public String GetMessage() {
            return Message;
        }
    }


    /// <summary>
    /// Rule validator context represens context of sqls.
    /// </summary>
    public class SqlserverContext {
        public SqlserverMeta SqlserverMeta;
        public Config Config;

        // advise context
        public Dictionary<String/*database*/, bool> AllDatabases;
        public Dictionary<String/*schema*/, bool> AllSchemas;
        public Dictionary<String/*schema.table*/, bool> AllTables;
        public bool databaseHasLoad;
        public bool schemaHasLoad;
        public bool tableHasLoad;

        public class DDLAction {
            public const String ADD_DATABASE = "add_database";
            public const String ADD_SCHEMA = "add_schema";
            public const String ADD_TABLE = "add_table";
            public const String REMOVE_DATABASE = "remove_database";
            public const String REMOVE_SCHEMA = "remove_schema";
            public const String REMOVE_TABLE = "remove_table";

            public String ID;
            public String Action;
        }

        public List<DDLAction> DDLActions;
        public Dictionary<String/*table*/, List<AlterTableStatement>> AlterTableStmts;
        public bool IsDDL;
        public bool IsDML;

        public AdviseResultContext AdviseResultContext;

        // rollback context
        public Dictionary<String/*database.schema.table*/, Dictionary<String/*column*/, String/*column definition*/>> TableColumnDefinitions;
        public Dictionary<String/*database.schema.table*/, Dictionary<String/*constraint*/, String/*constraint definition*/>> TableConstraintDefinitions;
        public Dictionary<String/*database.schema.table*/, Dictionary<String/*index*/, String/*index definition*/>> TableIndexDefinitions;

        public String GetConnectionString() {
            return String.Format(
                "Server=tcp:{0},{1};" +
                "Database=master;" +
                "User ID={2};" +
                "Password={3};",
                SqlserverMeta.Host, SqlserverMeta.Port,
                SqlserverMeta.User,
                SqlserverMeta.Password);
        }

        public void LoadDatabasesFromDB() {
            if (databaseHasLoad) {
                return;
            }

            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand("SELECT name FROM sys.databases", connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        AllDatabases[reader["name"] as String] = true;
                    }
                    databaseHasLoad = true;
                } finally {
                    reader.Close();
                }
            }
        }

        public void LoadSchemasFromDB() {
            if (schemaHasLoad) {
                return;
            }

            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand("SELECT name FROM sys.schemas", connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        AllSchemas[reader["name"] as String] = true;
                    }
                    schemaHasLoad = true;
                } finally {
                    reader.Close();
                }
            }
        }

        public void LoadTablesFromDB() {
            if (tableHasLoad) {
                return;
            }

            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand("SELECT TABLE_CATALOG, TABLE_SCHEMA, TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE'", connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        String database = reader["TABLE_CATALOG"] as String;
                        String schema = reader["TABLE_SCHEMA"] as String;
                        String table = reader["TABLE_NAME"] as String;
                        AllTables[String.Format("{0}.{1}.{2}", database, schema, table)] = true;
                    }
                    tableHasLoad = true;
                } finally {
                    reader.Close();
                }
            }
        }

        public void LoadCurrentDatabase() {
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand("SELECT DB_NAME() AS Database_name", connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        SqlserverMeta.CurrentDatabase = reader["Database_name"] as String;
                    }
                } finally {
                    reader.Close();
                }
            }
        }

        public void LoadCurrentSchema() {
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand("SELECT SCHEMA_NAME() AS Schema_name", connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        SqlserverMeta.CurrentSchema = reader["Schema_name"] as String;
                    }
                } finally {
                    reader.Close();
                }
            }
        }

        public Dictionary<String, bool> GetAllDatabases() {
            LoadDatabasesFromDB();
            return AllDatabases;
        }

        public Dictionary<String, bool> GetAllSchemas() {
            LoadSchemasFromDB();
            return AllSchemas;
        }

        public Dictionary<String, bool> GetAllTables() {
            LoadTablesFromDB();
            return AllTables;
        }

        public String GetCurrentDatabase() {
            if (SqlserverMeta.CurrentDatabase != "") {
                return SqlserverMeta.CurrentDatabase;
            }

            LoadCurrentDatabase();

            return SqlserverMeta.CurrentDatabase;
        }

        public String GetCurrentSchema() {
            if (SqlserverMeta.CurrentSchema != "") {
                return SqlserverMeta.CurrentSchema;
            }

            LoadCurrentSchema();

            return SqlserverMeta.CurrentSchema;
        }

        public List<String> GetPrimaryKeys(String databaseName, String schemaName, String tableName) {
            var ret = new List<String>();
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand(String.Format("SELECT " +
                                                                  "c.name AS Primary_key " +
                                                                  "FROM sys.indexes ix JOIN sys.index_columns ic ON ix.object_id=ic.object_id AND ix.index_id=ic.index_id JOIN sys.columns c ON ic.object_id=c.object_id AND ic.column_id=c.column_id " +
                                                                  "WHERE ix.object_id=OBJECT_ID('{0}.{1}.{2}', 'U') AND ix.is_primary_key = 1", databaseName, schemaName, tableName), connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        ret.Add((String)reader["Primary_key"]);
                    }
                } finally {
                    reader.Close();
                }
            }
            return ret;
        }

        public List<String> GetColumns(String databaseName, String schemaName, String tableName) {
            var ret = new List<String>();
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand(String.Format("SELECT " +
                                                                  "c.name AS Column_name " +
                                                                  "FROM sys.columns c WHERE c.object_id=OBJECT_ID('{0}.{1}.{2}', 'U')", databaseName, schemaName, tableName), connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        ret.Add((String)reader["Column_name"]);
                    }
                } finally {
                    reader.Close();
                }
            }
            return ret;
        }

        public List<Dictionary<String, String>> GetRecords(String databaseName, String schemaName, String tableName, String where) {
            var ret = new List<Dictionary<String, String>>();
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                String query = String.Format("SELECT * FROM {0}.{1}.{2} {3}", databaseName, schemaName, tableName, where);
                Console.WriteLine("query:{0}", query);
                SqlCommand command = new SqlCommand(query, connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        var row = new Dictionary<String, String>();
                        for (int index = 0; index < reader.FieldCount; index++) {
                            row[reader.GetName(index)] = reader.GetValue(index).ToString();
                        }
                        ret.Add(row);
                    }
                } finally {
                    reader.Close();
                }
            }
            return ret;
        }

        public int GetRecordsCount(String databaseName, String schemaName, String tableName, String where) {
            var ret = 0;
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                String query = String.Format("SELECT COUNT(*) AS counter FROM {0}.{1}.{2} {3}", databaseName, schemaName, tableName, where);
                Console.WriteLine("query:{0}", query);
                SqlCommand command = new SqlCommand(String.Format("SELECT COUNT(*) AS counter FROM {0}.{1}.{2} {3}", databaseName, schemaName, tableName, where), connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        ret = (int)reader["counter"];
                    }
                } finally {
                    reader.Close();
                }
            }
            return ret;
        }

        public void ResetTableColumnDefinitions(String databaseName, String schemaName, String tableName) {
            var columnDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableColumnDefinitions.ContainsKey(columnDefinitionKey)) {
                TableColumnDefinitions.Remove(columnDefinitionKey);
            }
        }

        public void SetTableColumnDefinitions(TableDefinition definition, String databaseName, String schemaName, String tableName) {
            var columnDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            foreach (var columnDefinition in definition.ColumnDefinitions) {
                var columnName = columnDefinition.ColumnIdentifier.Value;
                var columnString = "";
                for (int index = columnDefinition.FirstTokenIndex; index <= columnDefinition.LastTokenIndex; index++) {
                    columnString += columnDefinition.ScriptTokenStream[index].Text;
                }

                if (!TableColumnDefinitions.ContainsKey(columnDefinitionKey)) {
                    TableColumnDefinitions[columnDefinitionKey] = new Dictionary<String, String>();
                }
                TableColumnDefinitions[columnDefinitionKey][columnName] = columnString;
            }
        }

        public Dictionary<String, String> GetTableColumnDefinitions(String databaseName, String schemaName, String tableName) {
            var columnDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableColumnDefinitions.ContainsKey(columnDefinitionKey)) {
                return TableColumnDefinitions[columnDefinitionKey];
            }

            if (!TableExists(databaseName, schemaName, tableName)) {
                return new Dictionary<String, String>();
            }

            var result = new Dictionary<String, String>();
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand(String.Format("SELECT " +
                                                                  "c.name AS Column_name, " +
                                                                  "tp.name AS Type_name, " +
                                                                  "c.is_computed AS Is_computed, " +
                                                                  "OBJECT_DEFINITION(c.object_id, c.column_id) AS Column_definition, " +
                                                                  "c.system_type_id AS System_type_id," +
                                                                  "c.user_type_id AS User_type_id, " +
                                                                  "SCHEMA_NAME(tp.schema_id) AS Schema_name, " +
                                                                  "c.max_length AS Max_length, " +
                                                                  "c.precision AS Precesion, " +
                                                                  "c.scale AS Scale, " +
                                                                  "c.collation_name AS Collation_name, " +
                                                                  "c.is_nullable AS Is_nullable," +
                                                                  "OBJECT_NAME(c.default_object_id) AS Default_constraint_name, " +
                                                                  "OBJECT_DEFINITION(c.default_object_id) AS Default_definition, " +
                                                                  "cc.name AS Check_constraint_name, " +
                                                                  "cc.definition AS Check_definition, " +
                                                                  "c.is_identity AS Is_identity, " +
                                                                  "CAST(IDENTITYPROPERTY(c.object_id, 'SeedValue') AS VARCHAR(5)) AS Identity_base, " +
                                                                  "CAST(IDENTITYPROPERTY(c.object_id, 'IncrementValue') AS VARCHAR(5)) AS Identity_incr " +
                                                                  "FROM sys.columns c JOIN sys.types tp ON c.user_type_id=tp.user_type_id LEFT JOIN sys.check_constraints cc ON c.object_id=cc.parent_object_id AND cc.parent_column_id=c.column_id " +
                                                                  "WHERE c.object_id=OBJECT_ID('{0}.{1}.{2}', 'U')", databaseName, schemaName, tableName), connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        var sqlLine = "";
                        var colName = (String)reader["Column_name"];
                        sqlLine += reader["Column_name"];

                        bool isComputed = (bool)reader["Is_computed"];
                        if (isComputed) {
                            var columnDefinition = (String)reader["Column_definition"];
                            sqlLine += String.Format(" AS {0}", columnDefinition);
                        } else {
                            var systemTypeID = (Byte)reader["System_type_id"];
                            var userTypeID = (Int32)reader["User_type_id"];
                            var typeName = ((String)reader["Type_name"]).ToUpper();
                            if (systemTypeID != userTypeID) {
                                var schema = (String)reader["Schema_name"];
                                sqlLine += String.Format(" {0}.{1}", schema, typeName);
                            } else {
                                sqlLine += String.Format(" {0}", typeName.ToUpper());
                            }

                            if (typeName == "VARCHAR" || typeName == "CHAR" || typeName == "VARBINARY" || typeName == "BINARY") {
                                var maxLen = (Int16)reader["Max_length"];
                                sqlLine += String.Format("({0})", maxLen == -1 ? "MAX" : maxLen.ToString());
                            }
                            if (typeName == "NVARCHAR" || typeName == "NCHAR") {
                                var maxLen = (Int16)reader["Max_length"];
                                sqlLine += String.Format("({0})", maxLen == -1 ? "MAX" : (maxLen / 2).ToString());
                            }
                            if (typeName == "DATETIME2" || typeName == "TIME2" || typeName == "DATETIMEOFFSET") {
                                var scale = (Int32)reader["Scale"];
                                sqlLine += String.Format("({0})", scale);
                            }
                            if (typeName == "DECIMAL") {
                                var precesion = (Int32)reader["Precesion"];
                                var scale = (Int32)reader["Scale"];
                                sqlLine += String.Format("({0},{1})", precesion, scale);
                            }

                            var collationName = reader["Collation_name"];
                            if (systemTypeID == userTypeID && collationName.ToString() != "") {
                                sqlLine += String.Format(" COLLATE {0}", collationName.ToString());
                            }

                            var isNullable = (bool)reader["Is_nullable"];
                            if (!isNullable) {
                                sqlLine += " NOT NULL";
                            }

                            var defaultConstraintName = reader["Default_constraint_name"];
                            var defaultDefinition = reader["Default_definition"];
                            if (defaultConstraintName.ToString() != "" && defaultDefinition.ToString() != "") {
                                sqlLine += String.Format(" CONSTRAINT {0} DEFAULT {1}", defaultConstraintName, defaultDefinition);
                            }

                            var checkConstraintName = reader["Check_constraint_name"];
                            var checkDefinition = reader["Check_definition"];
                            if (checkConstraintName.ToString() != "" && checkDefinition.ToString() != "") {
                                sqlLine += String.Format(" CONSTRAINT {0} CHECK {1}", checkConstraintName, checkDefinition);
                            }

                            var isIdentity = (bool)reader["Is_identity"];
                            if (isIdentity) {
                                var identityBase = reader["Identity_base"];
                                var identityIncr = reader["Identity_incr"];
                                sqlLine += String.Format(" IDENTITY({0}, {1})", identityBase, identityIncr);
                            }
                        }

                        result[colName] = sqlLine;
                    }
                } finally {
                    reader.Close();
                }
            }

            TableColumnDefinitions[columnDefinitionKey] = result;
            return result;
        }

        public void ResetTableConstraintDefinitions(String databaseName, String schemaName, String tableName) {
            var constraintDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableConstraintDefinitions.ContainsKey(constraintDefinitionKey)) {
                TableConstraintDefinitions.Remove(constraintDefinitionKey);
            }
        }

        public void SetTableConstraintDefinitions(TableDefinition definition, String databaseName, String schemaName, String tableName) {
            var constraintDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            foreach (var tableConstraint in definition.TableConstraints) {
                var constraintName = tableConstraint.ConstraintIdentifier.Value;
                var constraintString = "";
                for (int index = tableConstraint.FirstTokenIndex; index <= tableConstraint.LastTokenIndex; index++) {
                    constraintString += tableConstraint.ScriptTokenStream[index].Text;
                }

                if (!TableConstraintDefinitions.ContainsKey(constraintDefinitionKey)) {
                    TableConstraintDefinitions[constraintDefinitionKey] = new Dictionary<String, String>();
                }
                TableConstraintDefinitions[constraintDefinitionKey][constraintName] = constraintString;
            }
        }

        public Dictionary<String, String> GetTableConstraintDefinitions(String databaseName, String schemaName, String tableName) {
            var constraintDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableConstraintDefinitions.ContainsKey(constraintDefinitionKey)) {
                return TableConstraintDefinitions[constraintDefinitionKey];
            }

            if (!TableExists(databaseName, schemaName, tableName)) {
                return new Dictionary<String, String>();
            }

            var result = new Dictionary<String, String>();
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand(String.Format("SELECT " +
                                                                  "kc.name AS Key_name, " +
                                                                  "ic.index_id AS Index_id, " +
                                                                  "c.name AS Column_name, " +
                                                                  "kc.type AS Type, " +
                                                                  "ic.is_descending_key AS Is_descending_key " +
                                                                  "FROM sys.key_constraints kc JOIN sys.index_columns ic ON kc.parent_object_id=ic.object_id AND kc.unique_index_id=ic.index_id JOIN sys.columns c ON ic.object_id=c.object_id AND ic.column_id=c.column_id " +
                                                                  "WHERE kc.parent_object_id=OBJECT_ID('{0}.{1}.{2}', 'U') AND (kc.type='PK' OR kc.type='UQ')", databaseName, schemaName, tableName), connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    var primaryKeyCols = new Dictionary<String, List<String>>();
                    var uniqueKeyCols = new Dictionary<String, List<String>>();
                    var ifClusteredCols = new Dictionary<String, bool>();
                    while (reader.Read()) {
                        var keyName = (String)reader["Key_name"];
                        var indexID = (Int32)reader["Index_id"];
                        var colName = (String)reader["Column_name"];
                        if (indexID == 1) {
                            ifClusteredCols[keyName] = true;
                        } else {
                            ifClusteredCols[keyName] = false;
                        }
                        var isDescendingKey = (bool)reader["Is_descending_key"];
                        var col = isDescendingKey ? colName + " DESC" : colName;


                        var type = (String)reader["Type"];
                        if (type == "PK") {
                            if (!primaryKeyCols.ContainsKey(keyName)) {
                                primaryKeyCols[keyName] = new List<string>();
                            }
                            primaryKeyCols[keyName].Add(col);
                        }
                        if (type == "UQ") {
                            if (!uniqueKeyCols.ContainsKey(keyName)) {
                                uniqueKeyCols[keyName] = new List<string>();
                            }
                            uniqueKeyCols[keyName].Add(col);
                        }
                    }

                    foreach (var primaryKeyColsPair in primaryKeyCols) {
                        var keyName = primaryKeyColsPair.Key;
                        var cols = primaryKeyColsPair.Value;
                        var primaryKeyConstraint = String.Format("CONSTRAINT {0} PRIMARY KEY {1} ({2})", keyName, ifClusteredCols[keyName] ? "CLUSTERED" : "NONCLUSTERED", String.Join(',', cols));
                        result[keyName] = primaryKeyConstraint;
                    }

                    foreach (var uniqueKeyColsPair in uniqueKeyCols) {
                        var keyName = uniqueKeyColsPair.Key;
                        var cols = uniqueKeyColsPair.Value;
                        var uniqueKeyConstraint = String.Format("CONSTRAINT {0} UNIQUE {1} ({2})", keyName, ifClusteredCols[keyName] ? "CLUSTERED" : "NONCLUSTERED", String.Join(',', cols));
                        result[keyName] = uniqueKeyConstraint;
                    }
                } finally {
                    reader.Close();
                }
            }

            TableConstraintDefinitions[constraintDefinitionKey] = result;
            return result;
        }

        public void ResetTableIndexDefinitions(String databaseName, String schemaName, String tableName) {
            var indexDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableIndexDefinitions.ContainsKey(indexDefinitionKey)) {
                TableIndexDefinitions.Remove(indexDefinitionKey);
            }
        }

        public void SetTableIndexDefinitions(TableDefinition definition, String databaseName, String schemaName, String tableName) {
            var indexDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            foreach (var index in definition.Indexes) {
                var indexName = index.Name.Value;
                var indexString = "";
                for (int i = index.FirstTokenIndex; i <= index.LastTokenIndex; i++) {
                    indexString += index.ScriptTokenStream[i];
                }

                if (!TableIndexDefinitions.ContainsKey(indexName)) {
                    TableIndexDefinitions[indexName] = new Dictionary<String, String>();
                }
                TableIndexDefinitions[indexDefinitionKey][indexName] = indexString;
            }
        }

        public Dictionary<String, String> GetTableIndexDefinitions(String databaseName, String schemaName, String tableName) {
            var indexDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableIndexDefinitions.ContainsKey(indexDefinitionKey)) {
                return TableIndexDefinitions[indexDefinitionKey];
            }

            if (!TableExists(databaseName, schemaName, tableName)) {
                return new Dictionary<String, String>();
            }

            var result = new Dictionary<String, String>();
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand(String.Format("SELECT " +
                                                                  "ix.name AS Index_name, " +
                                                                  "ix.type_desc AS Type_desc, " +
                                                                  "c.name AS Column_name, " +
                                                                  "ic.is_descending_key AS Is_descending_key " +
                                                                  "FROM sys.indexes ix JOIN sys.index_columns ic ON ix.object_id=ic.object_id AND ix.index_id=ic.index_id JOIN sys.columns c ON ic.object_id=c.object_id AND ic.column_id=c.column_id " +
                                                                  "WHERE ix.object_id=OBJECT_ID('{0}.{1}.{2}', 'U') AND ix.is_primary_key !=1 AND ix.is_unique_constraint !=1 AND ix.auto_created != 1", databaseName, schemaName, tableName), connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    var indexCols = new Dictionary<String, List<String>>();
                    var indexTypeDesc = new Dictionary<String, String>();
                    while (reader.Read()) {
                        var indexName = (String)reader["Index_name"];

                        var typeDesc = (String)reader["Type_desc"];
                        indexTypeDesc[indexName] = typeDesc;

                        var colName = (String)reader["Column_name"];
                        var isDescendingKey = (bool)reader["Is_descending_key"];
                        var col = isDescendingKey ? colName + " DESC" : colName;
                        if (!indexCols.ContainsKey(indexName)) {
                            indexCols[indexName] = new List<string>();
                        }
                        indexCols[indexName].Add(col);
                    }

                    foreach (var indexColsPair in indexCols) {
                        var indexName = indexColsPair.Key;
                        var cols = indexColsPair.Value;
                        var indexDefinition = String.Format("INDEX {0} {1} ({2})", indexName, indexTypeDesc[indexName], String.Join(',', cols));
                        result[indexName] = indexDefinition;
                    }
                } finally {
                    reader.Close();
                }
            }

            TableIndexDefinitions[indexDefinitionKey] = result;
            return result;
        }

        public SqlserverContext(SqlserverMeta sqlserverMeta, Config config): this(sqlserverMeta) {
            this.Config = config;
        }

        public SqlserverContext(SqlserverMeta sqlserverMeta) {
            this.SqlserverMeta = sqlserverMeta;
            AllDatabases = new Dictionary<String, bool>();
            AllSchemas = new Dictionary<string, bool>();
            AllTables = new Dictionary<String, bool>();
            DDLActions = new List<DDLAction>();
            AlterTableStmts = new Dictionary<string, List<AlterTableStatement>>();
            AdviseResultContext = new AdviseResultContext();
            tableHasLoad = false;
            schemaHasLoad = false;
            databaseHasLoad = false;
            IsDDL = false;
            IsDML = false;
            TableColumnDefinitions = new Dictionary<String, Dictionary<String, String>>();
            TableConstraintDefinitions = new Dictionary<String, Dictionary<String, String>>();
            TableIndexDefinitions = new Dictionary<String, Dictionary<String, String>>();
        }

        public bool DatabaseExists(String databaseName) {
            bool notBeDroped = true;
            foreach (var action in DDLActions) {
                if (action.ID == databaseName && action.Action == DDLAction.ADD_DATABASE) {
                    notBeDroped = true;
                }
                if (action.ID == databaseName && action.Action == DDLAction.REMOVE_DATABASE) {
                    notBeDroped = false;
                }
            }

            if (!notBeDroped) {
                return false;
            }

            var allDatabases = GetAllDatabases();
            return allDatabases.ContainsKey(databaseName);
        }

        public bool SchemaExists(String schema) {
            bool notBeDroped = true;
            foreach (var action in DDLActions) {
                if (action.ID == schema && action.Action == DDLAction.ADD_SCHEMA) {
                    notBeDroped = true;
                }
                if (action.ID == schema && action.Action == DDLAction.REMOVE_SCHEMA) {
                    notBeDroped = false;
                }
            }

            if (!notBeDroped) {
                return false;
            }

            var allschemas = GetAllSchemas();
            return allschemas.ContainsKey(schema);
        }

        public bool TableExists(String databaseName, String schema, String tableName) {
            if (schema == "") {
                schema = GetCurrentSchema();
            }
            String id = String.Format("{0}.{1}.{2}", databaseName, schema, tableName);
            bool databaseBeDropped = false;
            bool tableBeDropped = false;
            foreach (var action in DDLActions) {
                if (action.ID == id && action.Action == DDLAction.REMOVE_DATABASE) {
                    databaseBeDropped = true;
                }
                if (action.ID == id && action.Action == DDLAction.ADD_DATABASE) {
                    databaseBeDropped = false;
                }
                if (action.ID == id && action.Action == DDLAction.REMOVE_TABLE) {
                    tableBeDropped = true;
                }
                if (action.ID == id && action.Action == DDLAction.ADD_TABLE) {
                    tableBeDropped = false;
                }
            }

            if (databaseBeDropped && tableBeDropped) {
                return false;
            }

            var allTables = GetAllTables();
            return allTables.ContainsKey(id);
        }

        public void GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(SchemaObjectName schemaObjectName, out String databaseName, out String schemaName, out String tableName) {
            if (schemaObjectName == null) {
                databaseName = GetCurrentDatabase();
                schemaName = GetCurrentSchema();
                tableName = "";
                return;
            }

            var databaseIdentifier = schemaObjectName.DatabaseIdentifier;
            if (databaseIdentifier != null) {
                databaseName = databaseIdentifier.Value;
            } else {
                databaseName = GetCurrentDatabase();
            }

            var schemaIdentifier = schemaObjectName.SchemaIdentifier;
            if (schemaIdentifier != null) {
                schemaName = schemaIdentifier.Value;
            } else {
                schemaName = GetCurrentSchema();
            }

            tableName = schemaObjectName.BaseIdentifier.Value;
            return;
        }

        public void UpdateContext(TSqlStatement sqlStatement/*, bool needUpdateDefinition*/) {
            String databaseName, schemaName, tableName;
            switch (sqlStatement) {
                case UseStatement useStatement:
                    SqlserverMeta.CurrentDatabase = useStatement.DatabaseName.Value;
                    break;

                case CreateDatabaseStatement createDatabaseStatement:
                    DDLAction addDatabaseAction = new DDLAction {
                        ID = createDatabaseStatement.DatabaseName.Value,
                        Action = DDLAction.ADD_DATABASE,
                    };
                    DDLActions.Add(addDatabaseAction);
                    break;

                case CreateSchemaStatement createSchemaStatement:
                    DDLAction addSchemaAction = new DDLAction {
                        ID = createSchemaStatement.Name.Value,
                        Action = DDLAction.ADD_SCHEMA,
                    };
                    DDLActions.Add(addSchemaAction);
                    break;

                case CreateTableStatement createTableStatement:
                    GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(createTableStatement.SchemaObjectName, out databaseName, out schemaName, out tableName);
                    DDLAction addTableAction = new DDLAction {
                        ID = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName),
                        Action = DDLAction.ADD_TABLE,
                    };
                    DDLActions.Add(addTableAction);

                    SetTableColumnDefinitions(createTableStatement.Definition, databaseName, schemaName, tableName);
                    SetTableConstraintDefinitions(createTableStatement.Definition, databaseName, schemaName, tableName);
                    SetTableIndexDefinitions(createTableStatement.Definition, databaseName, schemaName, tableName);

                    break;

                case DropDatabaseStatement dropDatabaseStatement:
                    foreach (var database in dropDatabaseStatement.Databases) {
                        DDLAction removeDatabaseAction = new DDLAction {
                            ID = database.Value,
                            Action = DDLAction.REMOVE_DATABASE,
                        };
                        DDLActions.Add(removeDatabaseAction);

                        var newTableColumnDefinitions = new Dictionary<String, Dictionary<String, String>>();
                        foreach (var tableColumnDefinitionPair in TableColumnDefinitions) {
                            var tableIdentifier = tableColumnDefinitionPair.Key.Split(".");
                            if (tableIdentifier.Length == 3 && tableIdentifier[0] == database.Value) {
                                continue;
                            }
                            newTableColumnDefinitions[tableColumnDefinitionPair.Key] = tableColumnDefinitionPair.Value;
                        }
                        TableColumnDefinitions = newTableColumnDefinitions;

                        var newTableConstraintDefinitions = new Dictionary<String, Dictionary<String, String>>();
                        foreach (var tableConstraintDefinitionPair in TableConstraintDefinitions) {
                            var tableIdentifier = tableConstraintDefinitionPair.Key.Split(".");
                            if (tableIdentifier.Length == 3 && tableIdentifier[0] == database.Value) {
                                continue;
                            }
                            newTableConstraintDefinitions[tableConstraintDefinitionPair.Key] = tableConstraintDefinitionPair.Value;
                        }
                        TableConstraintDefinitions = newTableConstraintDefinitions;

                        var newTableIndexDefinitions = new Dictionary<String, Dictionary<String, String>>();
                        foreach (var tableIndexDefinitionPair in TableIndexDefinitions) {
                            var tableIdentifier = tableIndexDefinitionPair.Key.Split(".");
                            if (tableIdentifier.Length == 3 && tableIdentifier[0] == database.Value) {
                                continue;
                            }
                            newTableIndexDefinitions[tableIndexDefinitionPair.Key] = tableIndexDefinitionPair.Value;
                        }
                        TableIndexDefinitions = newTableIndexDefinitions;
                    }
                    break;

                case DropSchemaStatement dropSchemaStatement:
                    DDLAction removeSchemaAction = new DDLAction {
                        ID = dropSchemaStatement.Schema.BaseIdentifier.Value,
                        Action = DDLAction.REMOVE_SCHEMA,
                    };
                    DDLActions.Add(removeSchemaAction);
                    break;

                case DropTableStatement dropTableStatement:
                    foreach (var schemaObject in dropTableStatement.Objects) {
                        GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(schemaObject, out databaseName, out schemaName, out tableName);
                        DDLAction dropTableAction = new DDLAction {
                            ID = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName),
                            Action = DDLAction.REMOVE_TABLE,
                        };
                        DDLActions.Add(dropTableAction);

                        ResetTableColumnDefinitions(databaseName, schemaName, tableName);
                        ResetTableConstraintDefinitions(databaseName, schemaName, tableName);
                        ResetTableIndexDefinitions(databaseName, schemaName, tableName);
                    }
                    break;

                case ExecuteStatement executeStatement:
                    var entity = executeStatement.ExecuteSpecification.ExecutableEntity;
                    if (entity is ExecutableProcedureReference) {
                        var procedure = entity as ExecutableProcedureReference;
                        if (procedure.ProcedureReference is ProcedureReferenceName) {
                            var procedureName = procedure.ProcedureReference as ProcedureReferenceName;
                            var baseName = procedureName.ProcedureReference.Name.BaseIdentifier.Value;
                            if (baseName.ToLower() == "sp_rename") {
                                GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(null, out databaseName, out schemaName, out tableName);
                                if (procedure.Parameters.Count > 2) {
                                    // rename column
                                    var type = procedure.Parameters[2].ParameterValue as StringLiteral;
                                    if (type.Value.ToUpper() == "COLUMN") {
                                        var oldColumnName = "";
                                        var objNameArray = (procedure.Parameters[0].ParameterValue as StringLiteral).Value.Split(".");
                                        if (objNameArray.Length == 2) {
                                            tableName = objNameArray[0];
                                            oldColumnName = objNameArray[1];
                                        }
                                        if (objNameArray.Length == 3) {
                                            schemaName = objNameArray[0];
                                            tableName = objNameArray[1];
                                            oldColumnName = objNameArray[2];
                                        }
                                        var newColumnName = (procedure.Parameters[1].ParameterValue as StringLiteral).Value.Split(".")[0];

                                        var tableColumnDefinitions = GetTableColumnDefinitions(databaseName, schemaName, tableName);
                                        if (tableColumnDefinitions.Count == 0) {
                                            break;
                                        }
                                        var newColumnDefinitions = new Dictionary<String, String>();
                                        foreach (var columnDefinitionPair in tableColumnDefinitions) {
                                            if (columnDefinitionPair.Key == newColumnName) {
                                                newColumnDefinitions[newColumnName] = columnDefinitionPair.Value;
                                            } else {
                                                newColumnDefinitions[columnDefinitionPair.Key] = columnDefinitionPair.Value;
                                            }
                                        }
                                        TableColumnDefinitions[String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName)] = newColumnDefinitions;
                                    }
                                } else if (procedure.Parameters.Count == 2){
                                    // rename table
                                    var objNameArray = (procedure.Parameters[0].ParameterValue as StringLiteral).Value.Split(".");
                                    if (objNameArray.Length == 2) {
                                        schemaName = objNameArray[0];
                                        tableName = objNameArray[1];
                                        var newTableName = (procedure.Parameters[1].ParameterValue as StringLiteral).Value.Split(".")[0];
                                        String key = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
                                        String newKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, newTableName);

                                        var tableColumnDefinitions = GetTableColumnDefinitions(databaseName, schemaName, tableName);
                                        if (tableColumnDefinitions.Count > 0) {
                                            TableColumnDefinitions.Remove(key);
                                            TableColumnDefinitions[newKey] = tableColumnDefinitions;
                                        }

                                        var tableIndexDefinitions = GetTableIndexDefinitions(databaseName, schemaName, tableName);
                                        if (tableIndexDefinitions.Count > 0) {
                                            TableIndexDefinitions.Remove(key);
                                            TableIndexDefinitions[newKey] = tableIndexDefinitions;
                                        }

                                        var tableConstaintDefinitions = GetTableConstraintDefinitions(databaseName, schemaName, tableName);
                                        if (tableConstaintDefinitions.Count > 0) {
                                            TableConstraintDefinitions.Remove(key);
                                            TableConstraintDefinitions[newKey] = tableConstaintDefinitions;
                                        }
                                    }
                                }
                            }
                        }
                    }
                    break;
            }
        }
    }

    public abstract class RuleValidator {
        protected String Name;
        protected String Desc;
        protected String Message;
        protected RULE_LEVEL Level;

        // return validator name
        public String GetName() {
            return Name;
        }

        // return validator description
        public String GetDescription() {
            return Desc;
        }

        // return validator message
        public String GetMessage(params String[] paras) {
            return String.Format(Message, paras);
        }

        // return validator level
        public RULE_LEVEL GetLevel() {
            return Level;
        }

        public String GetLevelString() {
            return RULE_LEVEL_STRING.GetRuleLevelString(Level);
        }

        // if check failed, it will throw exception
        public abstract void Check(SqlserverContext context, TSqlStatement statement);

        public bool DatabaseExists(SqlserverContext context, String databaseName) {
            return context.DatabaseExists(databaseName);
        }

        public bool SchemaExists(SqlserverContext context, String schema) {
            return context.SchemaExists(schema);
        }

        public bool TableExists(SqlserverContext context, String databaseName, String schema, String table) {
            return context.TableExists(databaseName, schema, table);
        }

        protected RuleValidator(String name, String desc, String msg, RULE_LEVEL level) {
            Name = name;
            Desc = desc;
            Message = msg;
            Level = level;
        }

        public List<String> AddDatabaseName(List<String> databaseNames, SqlserverContext context, SchemaObjectName schemaObjectName) {
            var databaseIndentifier = schemaObjectName.DatabaseIdentifier;
            if (databaseIndentifier != null && databaseIndentifier.Value != "" && !databaseNames.Contains(databaseIndentifier.Value)) {
                databaseNames.Add(databaseIndentifier.Value);
            } else if (!databaseNames.Contains(context.SqlserverMeta.CurrentDatabase)) {
                databaseNames.Add(context.SqlserverMeta.CurrentDatabase);
            }

            return databaseNames;
        }

        public List<String> AddSchemaName(List<String> schemaNames, SchemaObjectName schemaObjectName) {
            var schemaIndentifier = schemaObjectName.SchemaIdentifier;
            if (schemaIndentifier != null && schemaIndentifier.Value != "" && !schemaNames.Contains(schemaIndentifier.Value)) {
                schemaNames.Add(schemaIndentifier.Value);
            }
            return schemaNames;
        }

        public List<string> AddTableName(List<String> tableNames, SqlserverContext context, SchemaObjectName schemaObjectName) {
            var databaseIdentifier = schemaObjectName.DatabaseIdentifier;
            var schemaIdentifier = schemaObjectName.SchemaIdentifier;
            var baseIndentifier = schemaObjectName.BaseIdentifier;
            var databaseName = "";
            var schemaName = "";
            var tableName = "";
            if (databaseIdentifier != null) {
                databaseName = databaseIdentifier.Value;
            } else {
                databaseName = context.GetCurrentDatabase();
            }
            if (schemaIdentifier != null) {
                schemaName = schemaIdentifier.Value;
            } else {
                schemaName = context.GetCurrentSchema();
            }
            tableName = baseIndentifier.Value;

            var key = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (!tableNames.Contains(key)) {
                tableNames.Add(key);
            }
            return tableNames;
        }

        public List<SchemaObjectName> AddSchemaObjectNameFromTableReference(List<SchemaObjectName> schemaObjectNames, TableReference tableReference) {
            switch (tableReference) {
                // SELECT col1 FROM table1 JOIN (SELECT col2 FROM table2 JOIN table3 ON table2.col3=table3.col3) as table4 ON tabl4.col2=table1.col2
                case QualifiedJoin qualifiedJoin:
                    if (qualifiedJoin.FirstTableReference != null) {
                        schemaObjectNames = AddSchemaObjectNameFromTableReference(schemaObjectNames, qualifiedJoin.FirstTableReference);
                    }
                    if (qualifiedJoin.SecondTableReference != null) {
                        schemaObjectNames = AddSchemaObjectNameFromTableReference(schemaObjectNames, qualifiedJoin.SecondTableReference);
                    }
                    break;
                case QueryDerivedTable queryDerivedTable:
                    if (queryDerivedTable.QueryExpression is QuerySpecification) {
                        var querySpec = queryDerivedTable.QueryExpression as QuerySpecification;
                        schemaObjectNames = AddSchemaObjectNameFromFromClause(schemaObjectNames, querySpec.FromClause);
                    }
                    break;

                // SELECT col1 FROM database1.schema1.table1 AS table2
                case NamedTableReference namedTableReference:
                    schemaObjectNames.Add(namedTableReference.SchemaObject);
                    break;
                default:
                    return schemaObjectNames;
            }
            return schemaObjectNames;
        }

        public List<SchemaObjectName> AddSchemaObjectNameFromQuerySpecification(List<SchemaObjectName> schemaObjectNames, QuerySpecification querySpecification) {
            // FromClause could be null such as SELECT @@IDENTIFY
            return AddSchemaObjectNameFromFromClause(schemaObjectNames, querySpecification.FromClause);
        }

        public List<SchemaObjectName> AddSchemaObjectNameFromFromClause(List<SchemaObjectName> schemaObjectNames, FromClause fromClause) {
            if (fromClause != null) {
                var tableReferences = fromClause.TableReferences;
                foreach (var tableReference in tableReferences) {
                    schemaObjectNames = AddSchemaObjectNameFromTableReference(schemaObjectNames, tableReference);
                }
            }
            return schemaObjectNames;
        }

        public void GetDatabaseAndSchemaAndTableNames(List<SchemaObjectName> schemaObjectNames, SqlserverContext context, List<String> databaseNames, List<String> schemaNames, List<String> tableNames) {
            foreach (var schemaObject in schemaObjectNames) {
                databaseNames = AddDatabaseName(databaseNames, context, schemaObject);
                schemaNames = AddSchemaName(schemaNames, schemaObject);
                tableNames = AddTableName(tableNames, context, schemaObject);
            }
        }
    }

    // FakeRuleValidator implements rule validator which do nothing.
    public class FakerRuleValidator : RuleValidator {
        public override void Check(SqlserverContext context, TSqlStatement statement) { }

        public FakerRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
