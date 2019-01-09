using System;
using System.Collections.Generic;
using SqlserverProto;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Data.SqlClient;
using NLog;

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
        // error message
        public static String DATABASE_NOT_EXIST_MSG = "database {0} 不存在";
        public static String DATABASE_EXIST_MSG = "database {0} 已存在";
        public static String SCHEMA_NOT_EXIST_MSG = "schema {0} 不存在";
        public static String TABLE_NOT_EXIST_MSG = "表 {0} 不存在";
        public static String TABLE_EXIST_MSG = "表 {0} 已存在";
        public static String COLUMN_NOT_EXIST_MSG = "字段 {0} 不存在";
        public static String COLUMN_EXIST_MSG = "字段 {0} 已存在";
        public static String INDEX_NOT_EXIST_MSG = "索引 {0} 不存在";
        public static String INDEX_EXIST_MSG = "索引 {0} 已存在";
        public static String CONSTRAINT_NOT_EXIST_MSG = "约束 {0} 不存在";
        public static String CONSTRAINT_EXIST_MSG = "约束 {0} 已存在";
        public static String DUPLICATE_COLUMN_ERROR_MSG = "字段名 {0} 重复";
        public static String DUPLICATE_INDEX_ERROR_MSG = "索引名 {0} 重复";
        public static String DUPLICATE_CONSTAINT_ERROR_MSG = "约束名 {0} 重复";
        public static String PRIMARY_KEY_MULTI_ERROR_MSG = "主键只能设置一个";
        public static String PRIMARY_KEY_EXIST_MSG = "已经存在主键，不能再添加";
        public static String KEY_COLUMN_NOT_EXIST_MSG = "索引字段 {0} 不存在";
        public static String CONSTRAINT_COLUMN_NOT_EXIST_MSG = "约束字段 {0} 不存在";
        public static String NOT_MATCH_VALUES_AND_COLUMNS = "指定的值列数与字段列数不匹配";

        // rule names
        public const String DDL_CHECK_OBJECT_NAME_LENGTH = "ddl_check_object_name_length";
        public const String DDL_CHECK_PK_NOT_EXIST = "ddl_check_pk_not_exist";
        public const String DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT = "ddl_check_pk_without_auto_increment";
        public const String DDL_CHECK_COLUMN_VARCHAR_MAX = "ddl_check_column_varchar_max";
        public const String DDL_CHECK_COLUMN_CHAR_LENGTH = "ddl_check_column_char_length";
        public const String DDL_DISABLE_FK = "ddl_disable_fk";
        public const String DDL_CHECK_INDEX_COUNT = "ddl_check_index_count";
        public const String DDL_CHECK_COMPOSITE_INDEX_MAX = "ddl_check_composite_index_max";
        public const String DDL_CHECK_OBJECT_NAME_USING_KEYWORD = "ddl_check_object_name_using_keyword";
        public const String DDL_CHECK_INDEX_COLUMN_WITH_BLOB = "ddl_check_index_column_with_blob";
        public const String DDL_CHECK_ALTER_TABLE_NEED_MERGE = "ddl_check_alter_table_need_merge";
        public const String DDL_DISABLE_DROP_STATEMENT = "ddl_disable_drop_statement";
        public const String ALL_CHECK_WHERE_IS_INVALID = "all_check_where_is_invalid";
        public const String DML_DISABE_SELECT_ALL_COLUMN = "dml_disable_select_all_column";
        public const String DDL_CHECK_INDEX_PREFIX = "ddl_check_index_prefix";
        public const String DDL_CHECK_UNIQUE_INDEX_PREFIX = "ddl_check_unique_index_prefix";
        public const String DDL_CHECK_COLUMN_WITHOUT_DEFAULT = "ddl_check_column_without_default";
        public const String DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT = "ddl_check_column_timestamp_without_default";
        public const String DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL = "ddl_check_column_blob_with_not_null";
        public const String DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL = "ddl_check_column_blob_default_is_not_null";
        public const String DML_CHECK_WITH_LIMIT = "dml_check_with_limit";

        // rules
        public static Dictionary<String, RuleValidator> RuleValidators = new Dictionary<String, RuleValidator> {
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
                DDL_CHECK_PK_NOT_EXIST,
                new PrimaryKeyShouldExistRuleValidator(
                    DDL_CHECK_PK_NOT_EXIST,
                    "表必须有主键",
                    "表必须有主键",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT,
                new PrimaryKeyAutoIncrementRuleValidator(
                    DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT,
                    "主键建议使用自增",
                    "主键建议使用自增",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_COLUMN_VARCHAR_MAX,
                new StringTypeShouldNoVarcharMaxRuleValidator(
                    DDL_CHECK_COLUMN_VARCHAR_MAX,
                    "禁止使用 varchar(max)",
                    "禁止使用 varchar(max)",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_COLUMN_CHAR_LENGTH,
                new StringTypeShouldNotExceedMaxLengthRuleValidator(
                    DDL_CHECK_COLUMN_CHAR_LENGTH,
                    "char长度大于20时，必须使用varchar类型",
                    "char长度大于20时，必须使用varchar类型",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_DISABLE_FK,
                new ForeignKeyRuleValidator(
                    DDL_DISABLE_FK,
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
                DDL_CHECK_OBJECT_NAME_USING_KEYWORD,
                new ObjectNameShouldNotContainsKeywordRuleValidator(
                    DDL_CHECK_OBJECT_NAME_USING_KEYWORD,
                    "数据库对象命名禁止使用关键字",
                    "数据库对象命名禁止使用关键字 %s",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_INDEX_COLUMN_WITH_BLOB,
                new DisableAddIndexForColumnsTypeBlob(
                    DDL_CHECK_INDEX_COLUMN_WITH_BLOB,
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
                ALL_CHECK_WHERE_IS_INVALID,
                new SelectWhereRuleValidator(
                    ALL_CHECK_WHERE_IS_INVALID,
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
            },
            {
                DDL_CHECK_INDEX_PREFIX,
                new CheckNormalIndexPrefix(
                    DDL_CHECK_INDEX_PREFIX,
                    "普通索引必须要以 \"idx_\" 为前缀",
                    "普通索引必须要以 \"idx_\" 为前缀",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_UNIQUE_INDEX_PREFIX,
                new CheckUniqueIndexPrefix(
                    DDL_CHECK_UNIQUE_INDEX_PREFIX,
                    "unique索引必须要以 \"uniq_\" 为前缀",
                    "unique索引必须要以 \"uniq_\" 为前缀",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_COLUMN_WITHOUT_DEFAULT,
                new CheckColumnWithoutDefault(
                    DDL_CHECK_COLUMN_WITHOUT_DEFAULT,
                    "除了自增列及大字段列之外，每个列都必须添加默认值",
                    "除了自增列及大字段列之外，每个列都必须添加默认值",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT,
                new CheckColumnTimestampWithoutDefaut(
                    DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT,
                    "timestamp 类型的列必须添加默认值",
                    "timestamp 类型的列必须添加默认值",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL,
                new CheckColumnBlobNotNull(
                    DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL,
                    "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
                    "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL,
                new CheckColumnBlobDefaultNotNull(
                    DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL,
                    "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
                    "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DML_CHECK_WITH_LIMIT,
                new TopConditionRuleValidator(
                    DML_CHECK_WITH_LIMIT,
                    "delete/update 语句不能有limit/top条件",
                    "delete/update 语句不能有top条件",
                    RULE_LEVEL.ERROR
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
        public bool BaseRuleFailed;
        public static bool BASE_RULE_OK = true;
        public static bool BASE_RULE_FAILED = false;

        public AdviseResultContext() {
            Level = RULE_LEVEL.NORMAL;
            Message = "";
            BaseRuleFailed = BASE_RULE_OK;
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
            BaseRuleFailed = BASE_RULE_OK;
        }

        public void SetBaseRuleStatus(bool failed) {
            BaseRuleFailed = failed;
        }

        public bool GetBaseRuleStatus() {
            return BaseRuleFailed;
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
        public Dictionary<String/*database.schema.table*/, bool> AllTables;
        public Dictionary<String/*database.schema.table*/, List<String>> PrimaryKeys;
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

        // test items
        public bool IsTest;
        public bool ExpectDatabaseExist;
        public bool ExpectSchemaExist;
        public bool ExpectTableExist;
        public String ExpectCurrentDatabase;
        public String ExpectCurrentSchema;
        public String ExpectDatabaseName;
        public String ExpectSchemaName;
        public String ExpectTableName;
        public List<String> ExpectColumns;
        public List<Dictionary<String, String>> ExpectRecords;
        public int ExpectRecordsCount;

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
            if (IsTest) {
                return ExpectCurrentDatabase;
            }
            if (SqlserverMeta.CurrentDatabase != "") {
                return SqlserverMeta.CurrentDatabase;
            }

            LoadCurrentDatabase();

            return SqlserverMeta.CurrentDatabase;
        }

        public String GetCurrentSchema() {
            if (IsTest) {
                return ExpectCurrentSchema;
            }
            if (SqlserverMeta.CurrentSchema != "") {
                return SqlserverMeta.CurrentSchema;
            }

            LoadCurrentSchema();

            return SqlserverMeta.CurrentSchema;
        }

        public List<String> GetPrimaryKeys(String databaseName, String schemaName, String tableName) {
            var key = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (PrimaryKeys.ContainsKey(key)) {
                return PrimaryKeys[key];
            }

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
            PrimaryKeys[key] = ret;
            return ret;
        }

        public List<String> GetColumns(String databaseName, String schemaName, String tableName) {
            if (IsTest) {
                return ExpectColumns;
            }
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
            if (IsTest) {
                return ExpectRecords;
            }

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
            if (IsTest) {
                return ExpectRecordsCount;
            }

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

        public Dictionary<String, String> GetTableColumnDefinitions(Logger logger, String databaseName, String schemaName, String tableName) {
            var columnDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableColumnDefinitions.ContainsKey(columnDefinitionKey)) {
                return TableColumnDefinitions[columnDefinitionKey];
            }

            if (!TableExists(logger, databaseName, schemaName, tableName)) {
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
            if (definition == null || definition.TableConstraints == null) {
                return;
            }
            var constraintDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            foreach (var tableConstraint in definition.TableConstraints) {
                if (tableConstraint.ConstraintIdentifier == null) {
                    continue;
                }
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

        public Dictionary<String, String> GetTableConstraintDefinitions(Logger logger, String databaseName, String schemaName, String tableName) {
            var constraintDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableConstraintDefinitions.ContainsKey(constraintDefinitionKey)) {
                return TableConstraintDefinitions[constraintDefinitionKey];
            }

            if (!TableExists(logger, databaseName, schemaName, tableName)) {
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
            if (definition == null || definition.Indexes == null) {
                return;
            }

            var indexDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            foreach (var index in definition.Indexes) {
                var indexName = index.Name.Value;
                var indexString = "";
                for (int i = index.FirstTokenIndex; i <= index.LastTokenIndex; i++) {
                    indexString += index.ScriptTokenStream[i].Text;
                }

                if (!TableIndexDefinitions.ContainsKey(indexDefinitionKey)) {
                    TableIndexDefinitions[indexDefinitionKey] = new Dictionary<String, String>();
                }
                TableIndexDefinitions[indexDefinitionKey][indexName] = indexString;
            }
        }

        public Dictionary<String, String> GetTableIndexDefinitions(Logger logger, String databaseName, String schemaName, String tableName) {
            var indexDefinitionKey = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
            if (TableIndexDefinitions.ContainsKey(indexDefinitionKey)) {
                return TableIndexDefinitions[indexDefinitionKey];
            }

            if (!TableExists(logger, databaseName, schemaName, tableName)) {
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
            PrimaryKeys = new Dictionary<string, List<string>>();
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

        public bool DatabaseExists(Logger logger, String databaseName) {
            if (IsTest) {
                return ExpectDatabaseExist;
            }
            bool notBeDroped = true;
            foreach (var action in DDLActions) {
                if (action.ID == databaseName && action.Action == DDLAction.ADD_DATABASE) {
                    logger.Info("ADD_DATABASE:{0}", databaseName);
                    notBeDroped = true;
                }
                if (action.ID == databaseName && action.Action == DDLAction.REMOVE_DATABASE) {
                    logger.Info("REMOVE_DATABASE:{0}", databaseName);
                    notBeDroped = false;
                }
            }

            if (!notBeDroped) {
                logger.Info("database {0} be dropped", databaseName);
                return false;
            }

            var allDatabases = GetAllDatabases();
            logger.Info("current databases:{0}", String.Join(",", allDatabases.Keys));
            return allDatabases.ContainsKey(databaseName);
        }

        public bool SchemaExists(Logger logger, String schema) {
            if (IsTest) {
                return ExpectSchemaExist;
            }
            bool notBeDroped = true;
            foreach (var action in DDLActions) {
                if (action.ID == schema && action.Action == DDLAction.ADD_SCHEMA) {
                    logger.Info("ADD_SCHEMA:{0}", schema);
                    notBeDroped = true;
                }
                if (action.ID == schema && action.Action == DDLAction.REMOVE_SCHEMA) {
                    logger.Info("REMOVE_SCHEMA:{0}", schema);
                    notBeDroped = false;
                }
            }

            if (!notBeDroped) {
                logger.Info("schema {0} be dropped", schema);
                return false;
            }

            var allschemas = GetAllSchemas();
            logger.Info("current schemas:{0}", String.Join(",", allschemas.Keys));
            return allschemas.ContainsKey(schema);
        }

        public bool TableExists(Logger logger, String databaseName, String schema, String tableName) {
            if (IsTest) {
                return ExpectTableExist;
            }
            if (schema == "") {
                schema = GetCurrentSchema();
            }
            String id = String.Format("{0}.{1}.{2}", databaseName, schema, tableName);
            logger.Info("table key:{0}", id);
            bool databaseBeDropped = false;
            bool tableBeDropped = false;
            foreach (var action in DDLActions) {
                if (action.ID == databaseName && action.Action == DDLAction.REMOVE_DATABASE) {
                    logger.Info("REMOVE_DATABASE:{0}", databaseName);
                    databaseBeDropped = true;
                }
                if (action.ID == databaseName && action.Action == DDLAction.ADD_DATABASE) {
                    logger.Info("ADD_DATABASE:{0}", databaseName);
                    databaseBeDropped = false;
                }
                if (action.ID == tableName && action.Action == DDLAction.REMOVE_TABLE) {
                    logger.Info("REMOVE_TABLE:{0}", tableName);
                    tableBeDropped = true;
                }
                if (action.ID == tableName && action.Action == DDLAction.ADD_TABLE) {
                    logger.Info("ADD_TABLE:{0}", tableName);
                    tableBeDropped = false;
                }
            }

            if (databaseBeDropped || tableBeDropped) {
                return false;
            }

            var allTables = GetAllTables();
            logger.Info("current tables:{0}", String.Join(",", allTables.Keys));
            return allTables.ContainsKey(id);
        }

        public void GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(SchemaObjectName schemaObjectName, out String databaseName, out String schemaName, out String tableName) {
            if (IsTest) {
                databaseName = ExpectDatabaseName;
                schemaName = ExpectSchemaName;
                tableName = ExpectTableName;
                return;
            }

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

        public void UpdateContext(Logger logger, TSqlStatement sqlStatement/*, bool needUpdateDefinition*/) {
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

                    var primaryKeys = new List<String>();
                    foreach (var columnDefinition in createTableStatement.Definition.ColumnDefinitions) {
                        if (columnDefinition.Constraints == null) {
                            continue;
                        }
                        foreach (var columnConstraint in columnDefinition.Constraints) {
                            if (columnConstraint is UniqueConstraintDefinition) {
                                if ((columnConstraint as UniqueConstraintDefinition).IsPrimaryKey) {
                                    primaryKeys.Add(columnDefinition.ColumnIdentifier.Value);
                                }
                            }
                        }
                    }
                    if (createTableStatement.Definition.TableConstraints != null) {
                        foreach (var tableConstraint in createTableStatement.Definition.TableConstraints) {
                            if (tableConstraint is UniqueConstraintDefinition) {
                                var uniqueConstraint = tableConstraint as UniqueConstraintDefinition;
                                if (uniqueConstraint.IsPrimaryKey) {
                                    foreach (var column in uniqueConstraint.Columns) {
                                        foreach (var identifier in column.Column.MultiPartIdentifier.Identifiers) {
                                            primaryKeys.Add(identifier.Value);
                                        }
                                    }
                                }
                            }
                        }
                    }

                    PrimaryKeys[String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName)] = primaryKeys;

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

                        var newPrimaryKeys = new Dictionary<String, List<String>>();
                        foreach (var primaryKeyPair in PrimaryKeys) {
                            var tableIdentifier = primaryKeyPair.Key.Split(".");
                            if (tableIdentifier.Length == 3 && tableIdentifier[0] == database.Value) {
                                continue;
                            }
                            newPrimaryKeys[primaryKeyPair.Key] = primaryKeyPair.Value;
                        }
                        PrimaryKeys = newPrimaryKeys;
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
                        var key = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
                        DDLAction dropTableAction = new DDLAction {
                            ID = key,
                            Action = DDLAction.REMOVE_TABLE,
                        };
                        DDLActions.Add(dropTableAction);

                        ResetTableColumnDefinitions(databaseName, schemaName, tableName);
                        ResetTableConstraintDefinitions(databaseName, schemaName, tableName);
                        ResetTableIndexDefinitions(databaseName, schemaName, tableName);

                        if (PrimaryKeys.ContainsKey(key)) {
                            PrimaryKeys.Remove(key);
                        }
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

                                        // update tableDefinition
                                        var tableColumnDefinitions = GetTableColumnDefinitions(logger, databaseName, schemaName, tableName);
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
                                        var key = String.Format("{0}.{1}.{2}", databaseName, schemaName, tableName);
                                        TableColumnDefinitions[key] = newColumnDefinitions;

                                        // update Primarykeys
                                        if (PrimaryKeys.ContainsKey(key)) {
                                            var primaryColumns = PrimaryKeys[key];
                                            var newPrimaryColumns = new List<String>();
                                            foreach (var primaryColumn in primaryColumns) {
                                                if (primaryColumn == oldColumnName) {
                                                    newPrimaryColumns.Add(newColumnName);
                                                } else {
                                                    newPrimaryColumns.Add(primaryColumn);
                                                }
                                            }
                                            PrimaryKeys[key] = newPrimaryColumns;
                                        }
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

                                        var tableColumnDefinitions = GetTableColumnDefinitions(logger, databaseName, schemaName, tableName);
                                        if (tableColumnDefinitions.Count > 0) {
                                            TableColumnDefinitions.Remove(key);
                                            TableColumnDefinitions[newKey] = tableColumnDefinitions;
                                        }

                                        var tableIndexDefinitions = GetTableIndexDefinitions(logger, databaseName, schemaName, tableName);
                                        if (tableIndexDefinitions.Count > 0) {
                                            TableIndexDefinitions.Remove(key);
                                            TableIndexDefinitions[newKey] = tableIndexDefinitions;
                                        }

                                        var tableConstaintDefinitions = GetTableConstraintDefinitions(logger, databaseName, schemaName, tableName);
                                        if (tableConstaintDefinitions.Count > 0) {
                                            TableConstraintDefinitions.Remove(key);
                                            TableConstraintDefinitions[newKey] = tableConstaintDefinitions;
                                        }

                                        if (PrimaryKeys.ContainsKey(key)) {
                                            var primaryColumns = PrimaryKeys[key];
                                            PrimaryKeys.Remove(key);
                                            PrimaryKeys[newKey] = primaryColumns;
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

        public bool DatabaseExists(Logger logger, SqlserverContext context, String databaseName) {
            return context.DatabaseExists(logger, databaseName);
        }

        public bool SchemaExists(Logger logger, SqlserverContext context, String schema) {
            return context.SchemaExists(logger, schema);
        }

        public bool TableExists(Logger logger, SqlserverContext context, String databaseName, String schema, String table) {
            return context.TableExists(logger, databaseName, schema, table);
        }

        protected RuleValidator(String name, String desc, String msg, RULE_LEVEL level) {
            Name = name;
            Desc = desc;
            Message = msg;
            Level = level;
        }

        protected RuleValidator() {}

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

        public bool IsBlobType(DataTypeReference dataTypeReference) {
            switch (dataTypeReference) {
                case SqlDataTypeReference sqlDataTypeReference:
                    return IsBlobTypeString(sqlDataTypeReference.Name.BaseIdentifier.Value, sqlDataTypeReference.Parameters);

                case XmlDataTypeReference xmlDataTypeReference:
                    return IsBlobTypeString(xmlDataTypeReference.Name.BaseIdentifier.Value, null);
            }
            return false;
        }

        public bool IsBlobTypeString(String type, IList<Literal> parameters) {
            switch (type.ToLower()) {
                case "image":
                case "text":
                case "xml":
                case "varbinary(max)":
                    return true;
                case "varbinary":
                    foreach (var param in parameters) {
                        if (param.Value.ToLower() == "max") {
                            return true;
                        }
                    }
                    break;
            }
            return false;
        }
    }
}
