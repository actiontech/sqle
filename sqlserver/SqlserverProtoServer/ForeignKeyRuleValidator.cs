using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class ForeignKeyRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            Console.WriteLine("statement type:{0}", statement);
            bool hasForeignKey = false;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    foreach (var columnDefinition in createTableStatement.Definition.ColumnDefinitions) {
                        foreach (var constraint in columnDefinition.Constraints) {
                            if (constraint is ForeignKeyConstraintDefinition) {
                                hasForeignKey = true;
                            }
                        }
                    }
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    TableDefinition tableDefinition = alterTableAddTableElementStatement.Definition;
                    foreach (var tableConstaint in tableDefinition.TableConstraints) {
                        if (tableConstaint is ForeignKeyConstraintDefinition) {
                            hasForeignKey = true;
                        }
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
