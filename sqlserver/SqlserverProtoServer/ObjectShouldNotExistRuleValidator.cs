using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;

namespace SqlserverProtoServer {
    public class DatabaseShouldNotExistRuleValidator : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            if (statement is CreateDatabaseStatement) {
                var createDatabaseStatement = statement as CreateDatabaseStatement;
                var databaseName = createDatabaseStatement.DatabaseName.Value;
                if (DatabaseExists(logger, context, databaseName)) {
                    logger.Debug("database {0} should not exist", databaseName);
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(databaseName));
                }
            }
        }

        public DatabaseShouldNotExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class TableShouldNotExistRuleValidator : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            if (statement is CreateTableStatement) {
                var createTableStatement = statement as CreateTableStatement;
                var tableNames = new List<String>();
                tableNames = AddTableName(tableNames, context, createTableStatement.SchemaObjectName);

                foreach (var tableName in tableNames) {
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
            }
        }

        public TableShouldNotExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
