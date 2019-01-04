using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;

namespace SqlserverProtoServer {
    public class TopConditionRuleValidator : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            TopRowFilter topRowFilter = null;
            switch (statement) {
                case DeleteStatement deleteStatement:
                    topRowFilter = deleteStatement.DeleteSpecification.TopRowFilter;
                    break;

                case UpdateStatement updateStatement:
                    topRowFilter = updateStatement.UpdateSpecification.TopRowFilter;
                    break;
            }

            if (topRowFilter != null) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public TopConditionRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
