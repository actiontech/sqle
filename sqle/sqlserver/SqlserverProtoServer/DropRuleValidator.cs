using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;

namespace SqlserverProtoServer {
    public class DisableDropRuleValidator : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            if(statement is DropDatabaseStatement) {
                logger.Debug("There exists drop database statement");
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
            if (statement is DropTableStatement) {
                logger.Debug("There exists drop table statement");
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public DisableDropRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
