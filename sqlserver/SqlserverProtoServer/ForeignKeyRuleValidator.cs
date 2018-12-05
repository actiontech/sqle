using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class ForeignKeyRuleValidator : RuleValidator {
        public bool hasForeignKeyConstraint(IList<ConstraintDefinition> constraints) {
            foreach (var constrait in constraints) {
                if (constrait is ForeignKeyConstraintDefinition) {
                    return true;
                }
            }
            return false;
        }
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            bool hasForeignKey = false;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    foreach (var columnDefinition in createTableStatement.Definition.ColumnDefinitions) {
                        if (hasForeignKeyConstraint(columnDefinition.Constraints)) {
                            hasForeignKey = true;
                        }
                    }
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    if (hasForeignKeyConstraint(alterTableAddTableElementStatement.Definition.TableConstraints)) {
                        hasForeignKey = true;
                    }
                    break;
            }

            if (hasForeignKey) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public ForeignKeyRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
