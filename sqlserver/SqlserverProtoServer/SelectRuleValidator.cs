using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;

namespace SqlserverProtoServer {
    public class SelectWhereRuleValidator : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public bool WhereClauseHasColumn(BooleanExpression booleanExpression) {
            switch (booleanExpression) {
                case BooleanComparisonExpression comparisonExpression:
                    if (comparisonExpression.FirstExpression is ColumnReferenceExpression || comparisonExpression.SecondExpression is ColumnReferenceExpression) {
                        return true;
                    }
                    break;

                case BooleanNotExpression notExpression:
                    return WhereClauseHasColumn(notExpression.Expression);

                case BooleanParenthesisExpression parenthesisExpression:
                    return WhereClauseHasColumn(parenthesisExpression.Expression);

                case BooleanBinaryExpression binaryExpression:
                    return WhereClauseHasColumn(binaryExpression.FirstExpression) && WhereClauseHasColumn(binaryExpression.SecondExpression);

                case InPredicate inPredicate:
                    if (inPredicate.Expression is ColumnReferenceExpression) {
                        return true;
                    }
                    break;

                case LikePredicate likePredicate:
                    if (likePredicate.FirstExpression is ColumnReferenceExpression) {
                        return true;
                    }
                    break;

                case BooleanTernaryExpression ternaryExpression:
                    if (ternaryExpression.FirstExpression is ColumnReferenceExpression) {
                        return true;
                    }
                    break;

                case BooleanIsNullExpression isNullExpression:
                    if (isNullExpression.Expression is ColumnReferenceExpression) {
                        return true;
                    }
                    break;
            }

            return false;
        }

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            if (statement is SelectStatement) {
                var select = statement as SelectStatement;
                var querySpec = select.QueryExpression as QuerySpecification;
                if (querySpec.WhereClause == null || !WhereClauseHasColumn(querySpec.WhereClause.SearchCondition)) {
                    logger.Debug("There is no effective where clause");
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                }
            }
        }

        public SelectWhereRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }


    public class SelectAllRuleValidator : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            if (statement is SelectStatement) {
                var select = statement as SelectStatement;
                var querySpec = select.QueryExpression as QuerySpecification;
                foreach (var selectElement in querySpec.SelectElements) {
                    if (selectElement is SelectStarExpression) {
                        logger.Debug("There is select all expression");
                        context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                    }
                }
            }
        }

        public SelectAllRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
