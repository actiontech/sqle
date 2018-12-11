using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class DisableDropRuleValidator : RuleValidator {
        public override void Check(SqlserverContext context, TSqlStatement statement) {
            if(statement is DropDatabaseStatement) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
            if (statement is DropTableStatement) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public DisableDropRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
