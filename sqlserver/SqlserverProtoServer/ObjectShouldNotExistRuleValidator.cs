using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;

namespace SqlserverProtoServer {
    public class ObjectShouldNotExistRuleValidator : RuleValidator {
        public List<String> DatabaseNames;
        public List<String> TableNames;

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    TableNames = AddTableName(TableNames, context, createTableStatement.SchemaObjectName);
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
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            base.Check(context, statement);

            foreach(var databaseName in DatabaseNames) {
                var exist = DatabaseExists(logger, context, databaseName);
                if (exist) {
                    logger.Debug("database {0} should not exist", databaseName);
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(databaseName));
                }
            }

            reset();
        }

        public DatabaseShouldNotExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class TableShouldNotExistRuleValidator : ObjectShouldNotExistRuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            base.Check(context, statement);

            foreach(var tableName in TableNames) {
                var tableIdentifier = tableName.Split(".");
                if (tableIdentifier.Length != 3) {
                    continue;
                }
                var exist = TableExists(logger, context, tableIdentifier[0], tableIdentifier[1], tableIdentifier[2]);
                if (exist) {
                    logger.Debug("table {0} should not exist", tableIdentifier[2]);
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(tableName));
                }
            }

            reset();
        }

        public TableShouldNotExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) {
        }
    }
}
