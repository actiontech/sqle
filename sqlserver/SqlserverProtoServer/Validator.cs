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
                    "database或者schema {0} 已存在",
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
    public class RuleValidatorContext {
        public SqlserverMeta SqlserverMeta;

        public Dictionary<String/*database*/, bool> AllDatabases;
        public Dictionary<String/*schema*/, bool> AllSchemas;
        public Dictionary<String/*schema*/, Dictionary<String/*table*/, bool>> AllTables;
        public Dictionary<String/*table*/, List<AlterTableStatement>> AlterTableStmts;
        public bool hasLoadFromDB;
        public bool IsDDL;
        public bool IsDML;

        public AdviseResultContext AdviseResultContext;

        public RuleValidatorContext(SqlserverMeta sqlserverMeta) {
            this.SqlserverMeta = sqlserverMeta;
            AllDatabases = new Dictionary<String, bool>();
            AllSchemas = new Dictionary<string, bool>();
            AllTables = new Dictionary<String, Dictionary<String, bool>>();
            AlterTableStmts = new Dictionary<string, List<AlterTableStatement>>();
            AdviseResultContext = new AdviseResultContext();
            hasLoadFromDB = false;
            IsDDL = false;
            IsDML = false;
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
        public abstract void Check(RuleValidatorContext context, TSqlStatement statement);

        void LoadFromDB(RuleValidatorContext context) {
            if (context.hasLoadFromDB) {
                return;
            }

            // for test
            context.SqlserverMeta.Host = "10.186.62.15";
            context.SqlserverMeta.Port = "1433";
            context.SqlserverMeta.User = "sa";
            context.SqlserverMeta.Password = "123456aB";

            String connectionString = String.Format(
                "Server=tcp:{0},{1};" +
                "Database=master;" +
                "User ID={2};" +
                "Password={3};",
                context.SqlserverMeta.Host, context.SqlserverMeta.Port,
                context.SqlserverMeta.User,
                context.SqlserverMeta.Password);

            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand("SELECT name FROM sys.databases", connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while(reader.Read()) {
                        context.AllDatabases[reader["name"] as String] = true;
                    }
                } finally {
                    reader.Close();
                }
            }

            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand("SELECT name FROM sys.schemas", connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        context.AllSchemas[reader["name"] as String] = true;
                    }
                } finally {
                    reader.Close();
                }
            }

            using (SqlConnection connection = new SqlConnection(connectionString)) {
                SqlCommand command = new SqlCommand("SELECT TABLE_SCHEMA, TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE'", connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        String schema = reader["TABLE_SCHEMA"] as String;
                        String table = reader["TABLE_NAME"] as String;
                        if (context.AllTables.ContainsKey(schema)) {
                            context.AllTables[schema][table] = true;
                        } else {
                            Dictionary<String, bool> schemaTables = new Dictionary<String, bool>();
                            schemaTables[table] = true;
                            context.AllTables[schema] = schemaTables;
                        }
                    }
                } finally {
                    reader.Close();
                }
            }
        }

        public bool DatabaseExists(RuleValidatorContext context, String databaseName) {
            LoadFromDB(context);
            return context.AllDatabases.ContainsKey(databaseName);
        }

        public bool SchemaExists(RuleValidatorContext context, String schema) {
            LoadFromDB(context);
            return context.AllSchemas.ContainsKey(schema);
        }

        public bool TableExists(RuleValidatorContext context, String table) {
            foreach (var pair in context.AllTables) {
                Dictionary<String, bool> valuePairs = pair.Value;
                if (valuePairs.ContainsKey(table)) {
                    return true;
                }
            }
            return false;
        }

        public TableDefinition GetTableDefinition(String tableName) {
            return new TableDefinition();
        }

        protected RuleValidator(String name, String desc, String msg, RULE_LEVEL level) {
            Name = name;
            Desc = desc;
            Message = msg;
            Level = level;
        }

        public List<String> AddDatabaseName(List<String> databaseNames, RuleValidatorContext context, SchemaObjectName schemaObjectName) {
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

        public List<string> AddTableName(List<String> tableNames, SchemaObjectName schemaObjectName) {
            var baseIndentifier = schemaObjectName.BaseIdentifier;
            if (!tableNames.Contains(baseIndentifier.Value)) {
                tableNames.Add(baseIndentifier.Value);
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

        public void GetDatabaseAndSchemaAndTableNames(List<SchemaObjectName> schemaObjectNames, RuleValidatorContext context, List<String> databaseNames, List<String> schemaNames, List<String> tableNames) {
            foreach (var schemaObject in schemaObjectNames) {
                databaseNames = AddDatabaseName(databaseNames, context, schemaObject);
                schemaNames = AddSchemaName(schemaNames, schemaObject);
                tableNames = AddTableName(tableNames, schemaObject);
            }
        }

        public void Show(List<String> list, String format) {
            foreach (var item in list) {
                Console.WriteLine(format, item);
            }
        }
    }

    // FakeRuleValidator implements rule validator which do nothing.
    public class FakerRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) { }

        public FakerRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
