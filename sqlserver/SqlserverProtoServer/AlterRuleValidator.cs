using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class MergeAlterTableRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            if (statement is AlterTableStatement) {
                AlterTableStatement alterTableStatement = statement as AlterTableStatement;
                String tableName = alterTableStatement.SchemaObjectName.BaseIdentifier.Value;
                List<AlterTableStatement> alterTableStatements;
                if (context.AlterTableStmts.ContainsKey(tableName)) {
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
