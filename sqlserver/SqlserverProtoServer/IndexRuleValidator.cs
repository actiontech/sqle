using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class NumberOfCompsiteIndexColumnsShouldNotExceedMaxRuleValidator : RuleValidator {
        public const int COMPOSITE_INDEX_MAX = 5;
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            int compositeIndexMax = 0;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    foreach (var index in createTableStatement.Definition.Indexes) {
                        if (index.Columns.Count > compositeIndexMax) {
                            compositeIndexMax = index.Columns.Count;
                        }
                    }
                    break;

                case CreateIndexStatement createIndexStatement:
                    compositeIndexMax = createIndexStatement.Columns.Count;
                    break;
            }
            Console.WriteLine(compositeIndexMax);
            if (compositeIndexMax > COMPOSITE_INDEX_MAX) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public NumberOfCompsiteIndexColumnsShouldNotExceedMaxRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class NumberOfIndexesShouldNotExceedMaxRuleValidator : RuleValidator {
        public int GetIndexCounterFromTableDefinition(TableDefinition tableDefinition) {
            int indexCounter = tableDefinition.Indexes.Count;
            indexCounter += GetIndexCounterFromColumnDefinitions(tableDefinition.ColumnDefinitions);
            indexCounter += GetIndexCounterFromTableConstraints(tableDefinition.TableConstraints);
            return indexCounter;
        }

        public int GetIndexCounterFromColumnDefinitions(IList<ColumnDefinition> columnDefinitions) {
            int indexCounter = 0;
            foreach (var columnDefinition in columnDefinitions) {
                foreach (var constraint in columnDefinition.Constraints) {
                    if (constraint is UniqueConstraintDefinition) {
                        indexCounter++;
                    }
                }
            }
            return indexCounter;
        }

        public int GetIndexCounterFromTableConstraints(IList<ConstraintDefinition> constraintDefinitions) {
            var indexCounter = 0;
            foreach (var constraint in constraintDefinitions) {
                if (constraint is UniqueConstraintDefinition) {
                    indexCounter++;
                }
            }
            return indexCounter;
        }

        public const int INDEX_MAX_NUMBER = 5;
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            int indexCounter = 0;
            TableDefinition tableDefinition = null;
            String tableName = "";

            switch (statement) {
                case CreateTableStatement createTableStatement:
                    indexCounter = GetIndexCounterFromTableDefinition(createTableStatement.Definition);
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    tableName = alterTableAddTableElementStatement.SchemaObjectName.BaseIdentifier.Value;
                    tableDefinition = GetTableDefinition(tableName);
                    indexCounter = GetIndexCounterFromTableDefinition(tableDefinition);
                    indexCounter += GetIndexCounterFromColumnDefinitions(alterTableAddTableElementStatement.Definition.ColumnDefinitions);
                    break;

                case CreateIndexStatement createIndexStatement:
                    tableDefinition = GetTableDefinition(createIndexStatement.OnName.BaseIdentifier.Value);
                    indexCounter = GetIndexCounterFromTableDefinition(tableDefinition);
                    indexCounter++;
                    break;
            }

            if (indexCounter > INDEX_MAX_NUMBER) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public NumberOfIndexesShouldNotExceedMaxRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class DisableAddIndexForColumnsTypeBlob : RuleValidator {
        public bool hasIndexForColumnsTypeBlobInTableDefinition(TableDefinition tableDefinition) {
            // unique index
            foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                foreach (var constriant in columnDefinition.Constraints) {
                    if (constriant is UniqueConstraintDefinition) {
                        if (columnDefinition.DataType.Name.BaseIdentifier.Value.ToUpper() == "BLOB") {
                            return true;
                        }
                    }
                }
            }

            foreach (var constraint in tableDefinition.TableConstraints) {
                if (constraint is UniqueConstraintDefinition) {
                    UniqueConstraintDefinition uniqueConstraintDefinition = constraint as UniqueConstraintDefinition;
                    if (hasIndexForColumnsTypeBlob(tableDefinition, uniqueConstraintDefinition.Columns)) {
                        return true;
                    }
                }
            }
            // indexes
            foreach (var index in tableDefinition.Indexes) {
                if (hasIndexForColumnsTypeBlob(tableDefinition, index.Columns)) {
                    return true;
                }
            }
            return false;
        }

        public bool hasIndexForColumnsTypeBlob(TableDefinition tableDefinition, IList<ColumnWithSortOrder> columns) {
            foreach (var column in columns) {
                ColumnReferenceExpression columnReferenceExpression = column.Column;
                foreach (var identifier in columnReferenceExpression.MultiPartIdentifier.Identifiers) {
                    foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                        if (identifier.Value == columnDefinition.ColumnIdentifier.Value && columnDefinition.DataType.Name.BaseIdentifier.Value.ToUpper() == "BLOB") {
                            return true;
                        }
                    }
                }
            }
            return false;
        }

        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            bool indexDataTypeIsBlob = false;
            TableDefinition tableDefinition;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    tableDefinition = createTableStatement.Definition;
                    if (hasIndexForColumnsTypeBlobInTableDefinition(tableDefinition)) {
                        indexDataTypeIsBlob = true;
                    }
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    tableDefinition = alterTableAddTableElementStatement.Definition;
                    if (hasIndexForColumnsTypeBlobInTableDefinition(tableDefinition)) {
                        indexDataTypeIsBlob = true;
                    }
                    break;

                case AlterTableAlterColumnStatement alterTableAlterColumnStatement:
                    if (alterTableAlterColumnStatement.DataType.Name.BaseIdentifier.Value.ToUpper() == "BLOB") {
                        String columnName = alterTableAlterColumnStatement.ColumnIdentifier.Value;
                        tableDefinition = GetTableDefinition(alterTableAlterColumnStatement.SchemaObjectName.BaseIdentifier.Value);
                        foreach (var index in tableDefinition.Indexes) {
                            foreach (var column in index.Columns) {
                                ColumnReferenceExpression columnReferenceExpression = column.Column;
                                foreach (var indentifier in columnReferenceExpression.MultiPartIdentifier.Identifiers) {
                                    if (indentifier.Value == columnName) {
                                        indexDataTypeIsBlob = true;
                                    }
                                }
                            }
                        }
                    }
                    break;

                case CreateIndexStatement createIndexStatement:
                    tableDefinition = GetTableDefinition(createIndexStatement.OnName.BaseIdentifier.Value);
                    if (hasIndexForColumnsTypeBlob(tableDefinition, createIndexStatement.Columns)) {
                        indexDataTypeIsBlob = true;
                    }
                    break;
            }
            if (indexDataTypeIsBlob) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public DisableAddIndexForColumnsTypeBlob(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
