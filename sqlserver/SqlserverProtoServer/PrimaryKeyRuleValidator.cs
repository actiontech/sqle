using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class PrimaryKeyShouldExistRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            if (statement is CreateTableStatement) {
                bool hasPrimaryKey = false;
                CreateTableStatement createTableStatement = statement as CreateTableStatement;
                TableDefinition tableDefinition = createTableStatement.Definition;

                /*
                    CREATE TABLE schema1.table1(
                        col1 INT NOT NULL PRIMARY KEY CLUSTERED)
                */
                foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                    foreach (var constraint in columnDefinition.Constraints) {
                        if (constraint is UniqueConstraintDefinition) {
                            UniqueConstraintDefinition uniqueConstraintDefinition = constraint as UniqueConstraintDefinition;
                            if (uniqueConstraintDefinition.IsPrimaryKey) {
                                hasPrimaryKey = true;
                            }
                        }
                    }
                }

                /*
                    CREATE TABLE schema1.table1(
                        col1 INT NOT NULL,
                        col2 INT NOT NULL,
                        CONSTRAINT PK_constraint PRIMARY KEY CLUSTERED(col1, col2) WITH (IGNORE_DUP_KEY = OFF))
                */
                foreach (var tableConstraint in tableDefinition.TableConstraints) {
                    if (tableConstraint is UniqueConstraintDefinition) {
                        UniqueConstraintDefinition uniqueConstraintDefinition = tableConstraint as UniqueConstraintDefinition;
                        if (uniqueConstraintDefinition.IsPrimaryKey) {
                            hasPrimaryKey = true;
                        }
                    }
                }

                if (!hasPrimaryKey) {
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                }
            }
        }

        public PrimaryKeyShouldExistRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class PrimaryKeyAutoIncrementRuleValidator : RuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            CreateTableStatement createTableStatement = statement as CreateTableStatement;
            TableDefinition tableDefinition = createTableStatement.Definition;
            bool isPrimaryKeyAutoIncrement = false;

            /*
                    CREATE TABLE schema1.table1(
                        col1 INT IDENTITY(1,1) PRIMARY KEY CLUSTERED)
                */
            foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                bool isPrimaryColumn = false;
                foreach (var constraint in columnDefinition.Constraints) {
                    if (constraint is UniqueConstraintDefinition) {
                        UniqueConstraintDefinition uniqueConstraintDefinition = constraint as UniqueConstraintDefinition;
                        if (uniqueConstraintDefinition.IsPrimaryKey) {
                            isPrimaryColumn = true;
                        }
                    }
                }
                if (isPrimaryColumn && columnDefinition.IdentityOptions != null) {
                    isPrimaryKeyAutoIncrement =  true;
                }
            }

            /*
                    CREATE TABLE schema1.table1(
                        col1 INT IDENTITY(1, 1),
                        col2 INT NOT NULL,
                        CONSTRAINT PK_constraint PRIMARY KEY(col1, col2))
                */
            foreach (var tableConstraint in tableDefinition.TableConstraints) {
                if (tableConstraint is UniqueConstraintDefinition) {
                    UniqueConstraintDefinition uniqueConstraintDefinition = tableConstraint as UniqueConstraintDefinition;
                    if (uniqueConstraintDefinition.IsPrimaryKey) {
                        foreach (var primaryColumn in uniqueConstraintDefinition.Columns) {
                            ColumnReferenceExpression columnReferenceExpression = primaryColumn.Column;
                            foreach (var identifier in columnReferenceExpression.MultiPartIdentifier.Identifiers) {
                                foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                                    if (identifier.Value == columnDefinition.ColumnIdentifier.Value && columnDefinition.IdentityOptions != null) {
                                        isPrimaryKeyAutoIncrement = true;
                                    }
                                }
                            }
                        }
                    }
                }
            }

            if (!isPrimaryKeyAutoIncrement) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public PrimaryKeyAutoIncrementRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
