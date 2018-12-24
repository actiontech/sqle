using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;

namespace SqlserverProtoServer {
    public class MergeAlterTableRuleValidator : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            if (statement is AlterTableStatement) {
                AlterTableStatement alterTableStatement = statement as AlterTableStatement;
                String tableName = alterTableStatement.SchemaObjectName.BaseIdentifier.Value;
                List<AlterTableStatement> alterTableStatements;
                if (context.AlterTableStmts.ContainsKey(tableName)) {
                    logger.Debug("There exists multiple alter table statements for table:{0}", tableName);
                    alterTableStatements = context.AlterTableStmts[tableName];
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                } else {
                    alterTableStatements = new List<AlterTableStatement>();
                    context.AlterTableStmts[tableName] = alterTableStatements;
                }
                alterTableStatements.Add(alterTableStatement);
            }
        }

        public MergeAlterTableRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
