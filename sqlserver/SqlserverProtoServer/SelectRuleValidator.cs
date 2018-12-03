using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class SelectWhereRuleValidator : RuleValidator {
        public bool WhereClauseHasColumn(BooleanExpression booleanExpression) {
            Console.WriteLine("booleanExpression:{0}", booleanExpression);
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

        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            if (statement is SelectStatement) {
                var select = statement as SelectStatement;
                var querySpec = select.QueryExpression as QuerySpecification;
                if (querySpec.WhereClause == null || !WhereClauseHasColumn(querySpec.WhereClause.SearchCondition)) {
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                }
            }
        }

        public SelectWhereRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }


    public class SelectAllRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            if (statement is SelectStatement) {
                var select = statement as SelectStatement;
                var querySpec = select.QueryExpression as QuerySpecification;
                foreach (var selectElement in querySpec.SelectElements) {
                    if (selectElement is SelectStarExpression) {
                        context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                    }
                }
            }
        }

        public SelectAllRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
