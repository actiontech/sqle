using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class ObjectShouldNotExistRuleValidator : RuleValidator {
        public List<String> DatabaseNames;
        public List<String> TableNames;

        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    DatabaseNames = AddDatabaseName(DatabaseNames, context, createTableStatement.SchemaObjectName);
                    TableNames = AddTableName(TableNames, createTableStatement.SchemaObjectName);
                    break;

                case CreateDatabaseStatement createDatabaseStatement:
                    DatabaseNames.Add(createDatabaseStatement.DatabaseName.Value);
                    break;
            }
        }

        public  void reset() {
            DatabaseNames = new List<String>();
            TableNames = new List<String>();
        }

        public ObjectShouldNotExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) {
            DatabaseNames = new List<String>();
            TableNames = new List<String>();
        }
    }

    public class DatabaseShouldNotExistRuleValidator : ObjectShouldNotExistRuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            base.Check(context, statement);

            foreach(var databaseName in DatabaseNames) {
                var exist = DatabaseExists(context, databaseName);
                if (exist) {
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(databaseName));
                }
            }

            reset();
        }

        public DatabaseShouldNotExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class TableShouldNotExistRuleValidator : ObjectShouldNotExistRuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            base.Check(context, statement);

            foreach(var tableName in TableNames) {
                var exist = TableExists(context, tableName);
                if (exist) {
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(tableName));
                }
            }

            reset();
        }

        public TableShouldNotExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) {
        }
    }
}
