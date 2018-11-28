using System;
using System.Collections.Generic;
using SqlserverProto;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public enum RULE_LEVEL {
        NORMAL, NOTICE, WARN, ERROR
    }

    public static class DefaultRules {
        // rule names
        public const string SCHEMA_NOT_EXIST = "schema_not_exist";
        public const string SCHEMA_EXIST = "schema_exist";
        public const string TABLE_NOT_EXIST = "table_not_exist";
        public const string TABLE_EXIST = "table_exist";
        public const string DDL_CREATE_TABLE_NOT_EXIST = "ddl_create_table_not_exist";
        public const string DDL_CHECK_OBJECT_NAME_LENGTH = "ddl_check_object_name_length";
        public const string DDL_CHECK_PRIMARY_KEY_EXIST = "ddl_check_primary_key_exist";
        public const string DDL_CHECK_PRIMARY_KEY_TYPE = "ddl_check_primary_key_type";
        public const string DDL_DISABLE_VARCHAR_MAX = "ddl_disable_varchar_max";
        public const string DDL_CHECK_TYPE_CHAR_LENGTH = "ddl_check_type_char_length";
        public const string DDL_DISABLE_FOREIGN_KEY = "ddl_disable_foreign_key";
        public const string DDL_CHECK_INDEX_COUNT = "ddl_check_index_count";
        public const string DDL_CHECK_COMPOSITE_INDEX_MAX = "ddl_check_composite_index_max";
        public const string DDL_DISABLE_USING_KEYWORD = "ddl_disable_using_keyword";
        //private const string DDL_TABLE_USING_INNODB_UTF8MB4 = "ddl_create_table_using_innodb";
        public const string DDL_DISABLE_INDEX_DATA_TYPE_BLOB = "ddl_disable_index_column_blob";
        public const string DDL_CHECK_ALTER_TABLE_NEED_MERGE = "ddl_check_alter_table_need_merge";
        public const string DDL_DISABLE_DROP_STATEMENT = "ddl_disable_drop_statement";
        public const string DML_CHECK_INVALID_WHERE_CONDITION = "ddl_check_invalid_where_condition";
        public const string DML_DISABE_SELECT_ALL_COLUMN = "dml_disable_select_all_column";

        // rules
        public static Dictionary<string, RuleValidator> RuleValidators = new Dictionary<string, RuleValidator> {
            {
                SCHEMA_NOT_EXIST,
                new ObjectNotExistRuleValidator(
                    SCHEMA_NOT_EXIST,
                    "操作数据库时，数据库必须存在",
                    "schema {0} 不存在",
                    RULE_LEVEL.ERROR
                )
            },
            {
                SCHEMA_EXIST,
                new ObjectExistRuleValidator(
                    SCHEMA_EXIST,
                    "创建数据库时，数据库不能存在",
                    "schema {0} 已存在",
                    RULE_LEVEL.ERROR
                )
            },
            {
                TABLE_NOT_EXIST,
                new ObjectNotExistRuleValidator(
                    TABLE_NOT_EXIST,
                    "操作表时，表必须存在",
                    "表 {0} 不存在",
                    RULE_LEVEL.ERROR
                )
            },
            {
                TABLE_EXIST,
                new ObjectExistRuleValidator(
                    TABLE_EXIST,
                    "创建表时，表不能存在",
                    "表 {0} 已存在",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CREATE_TABLE_NOT_EXIST,
                new IfNotExistRuleValidator(
                    DDL_CREATE_TABLE_NOT_EXIST,
                    "新建表必须加入if not exists create，保证重复执行不报错",
                    "新建表必须加入if not exists create，保证重复执行不报错",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_OBJECT_NAME_LENGTH,
                new NewObjectNameRuleValidator(
                    DDL_CHECK_OBJECT_NAME_LENGTH,
                    "表名、列名、索引名的长度不能大于64字节",
                    "表名、列名、索引名的长度不能大于64字节",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_PRIMARY_KEY_EXIST,
                new PrimaryKeyRuleValidator(
                    DDL_CHECK_PRIMARY_KEY_EXIST,
                    "表必须有主键",
                    "表必须有主键",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_PRIMARY_KEY_TYPE,
                new PrimaryKeyRuleValidator(
                    DDL_CHECK_PRIMARY_KEY_TYPE,
                    "主键建议使用自增，且为bigint无符号类型，即bigint unsigned",
                    "主键建议使用自增，且为bigint无符号类型，即bigint unsigned",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_DISABLE_VARCHAR_MAX,
                new DisableVarcharMaxRuleValidator(
                    DDL_DISABLE_VARCHAR_MAX,
                    "禁止使用 varchar(max)",
                    "禁止使用 varchar(max)",
                    RULE_LEVEL.ERROR
                )
            },
            {
                DDL_CHECK_TYPE_CHAR_LENGTH,
                new StringTypeRuleValidator(
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
                new IndexRuleValidator(
                    DDL_CHECK_INDEX_COUNT,
                    "索引个数建议不超过5个",
                    "索引个数建议不超过5个",
                    RULE_LEVEL.NOTICE
                )
            },
            {
                DDL_CHECK_COMPOSITE_INDEX_MAX,
                new IndexRuleValidator(
                    DDL_CHECK_COMPOSITE_INDEX_MAX,
                    "复合索引的列数量不建议超过5个",
                    "复合索引的列数量不建议超过5个",
                    RULE_LEVEL.NOTICE
                )
            },
            {
                DDL_DISABLE_USING_KEYWORD,
                new NewObjectNameRuleValidator(
                    DDL_DISABLE_USING_KEYWORD,
                    "数据库对象命名禁止使用关键字",
                    "数据库对象命名禁止使用关键字 %s",
                    RULE_LEVEL.ERROR
                )
            },
            /*
             * {
             *  DDL_TABLE_USING_INNODB_UTF8MB4,
             *  new EngineAndCharacterSetRuleValidator(
             *      DDL_TABLE_USING_INNODB_UTF8MB4,
             *      "建议使用Innodb引擎,utf8mb4字符集",
             *      "建议使用Innodb引擎,utf8mb4字符集",
             *      RULE_LEVEL.NOTICE
             *  )
             * },
            */
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
    /// Rule validator context represens context of sqls.
    /// </summary>
    public class RuleValidatorContext {
        public string CurrentSchema;
        public Dictionary<string/*schema*/, bool> AllSchemas;
        public bool IsSchemaLoaded;
        public Dictionary<string/*schema*/, Dictionary<string/*table*/, bool>> AllTables;
        public bool IsDDL;
        public bool IsDML;

        public AuditResultContext AuditResultContext;

        public RuleValidatorContext() {
            AllSchemas = new Dictionary<string, bool>();
            AllTables = new Dictionary<string, Dictionary<string, bool>>();
            AuditResultContext = new AuditResultContext();
                   }
    }

    /// <summary>
    /// Audit result context represents context of all audit rule results for one sql.
    /// </summary>
    public class AuditResultContext {
        public RULE_LEVEL Level;
        public string Message;
        public static Dictionary<RULE_LEVEL, string> RuleLevels = new Dictionary<RULE_LEVEL, string> {
            {RULE_LEVEL.NORMAL, "normal"},
            {RULE_LEVEL.NOTICE, "notice"},
            {RULE_LEVEL.WARN, "warn"},
            {RULE_LEVEL.ERROR, "error"}
        };

        public AuditResultContext() {
            Level = RULE_LEVEL.NORMAL;
            Message = "";
        }

        public void AddAuditResult(RULE_LEVEL level, string message) {
            if (Level < level) {
                Level = level;
            }

            var formatMsg = String.Format("[{0}]{1}", RuleLevels[level], message);
            if (String.IsNullOrEmpty(Message)) {
                Message = formatMsg;
            } else {
                Message += "\n" + formatMsg;
            }
        }

        public void ResetAuditResult() {
            Level = RULE_LEVEL.NORMAL;
            Message = "";
        }

        public AuditResult GetAuditResult() {
            AuditResult auditResult = new AuditResult();
            auditResult.AuditLevel = GetLevel();
            auditResult.AuditResultMessage = GetMessage();
            return auditResult;
        }

        public string GetLevel() {
            return RuleLevels[Level];
        }

        public string GetMessage() {
            return Message;
        }
    }

    public abstract class RuleValidator {
        protected string Name;
        protected string Desc;
        protected string Message;
        protected RULE_LEVEL Level;

        // if check failed, it will throw exception
        public abstract void Check(RuleValidatorContext context, TSqlStatement statement);
        // return validator name
        public string GetName() {
            return Name;
        }
        // return validator description
        public string GetDescription() {
            return Desc;
        }
        // return validator message
        public string GetMessage(params string[] paras) {
            return String.Format(Message, paras);
        }
        // return validator level
        public RULE_LEVEL GetLevel() {
            return Level;
        }

        protected RuleValidator(string name, string desc, string msg, RULE_LEVEL level) {
            Name = name;
            Desc = desc;
            Message = msg;
            Level = level;
        }
    }

    public class ObjectNotExistRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public ObjectNotExistRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class ObjectExistRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public ObjectExistRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class IfNotExistRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public IfNotExistRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class NewObjectNameRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public NewObjectNameRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class PrimaryKeyRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public PrimaryKeyRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class DisableVarcharMaxRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public DisableVarcharMaxRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class StringTypeRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public StringTypeRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class ForeignKeyRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public ForeignKeyRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class IndexRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public IndexRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
    /*
    public class EngineAndCharacterSetRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public EngineAndCharacterSetRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
    */

    public class DisableAddIndexForColumnsTypeBlob : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public DisableAddIndexForColumnsTypeBlob(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class MergeAlterTableRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public MergeAlterTableRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class DisableDropRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public DisableDropRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class SelectWhereRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            throw new NotImplementedException();
        }

        public SelectWhereRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class SelectAllRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            if (statement is SelectStatement) {
                var select = statement as SelectStatement;
                var querySpec = select.QueryExpression as QuerySpecification;
                foreach (var selectElement in querySpec.SelectElements) {
                    if (selectElement is SelectStarExpression) {
                        context.AuditResultContext.AddAuditResult(GetLevel(), GetMessage());
                    }
                }
            }
        }

        public SelectAllRuleValidator(string name, string desc, string msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
