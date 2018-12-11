using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class StringTypeShouldNotExceedMaxLengthRuleValidator : RuleValidator {
        public const int CHAR_MAX_LENGTH = 20;

        public bool isCharLengthExceedMaxLengthInDefinitions(IList<ColumnDefinition> columnDefinitions) {
            foreach (var columnDefinition in columnDefinitions) {
                if (isCharLengthExceedMaxLengthInDefinitionsInDataType(columnDefinition.DataType)) {
                    return true;
                }
            }

            return false;
        }

        public bool isCharLengthExceedMaxLengthInDefinitionsInDataType(DataTypeReference dataTypeReference) {
            if (dataTypeReference is ParameterizedDataTypeReference) {
                ParameterizedDataTypeReference parameterizedDataTypeReference = dataTypeReference as ParameterizedDataTypeReference;
                if (parameterizedDataTypeReference.Name.BaseIdentifier.Value.ToUpper() == "CHAR") {
                    foreach (var parameter in parameterizedDataTypeReference.Parameters) {
                        if (parameter is IntegerLiteral) {
                            IntegerLiteral integerLiteral = parameter as IntegerLiteral;
                            int varcharLength = Int32.Parse(integerLiteral.Value);
                            if (varcharLength > CHAR_MAX_LENGTH) {
                                return true;
                            }
                        }
                    }
                }
            }
            return false;
        }

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    TableDefinition tableDefinition = createTableStatement.Definition;
                    if (isCharLengthExceedMaxLengthInDefinitions(tableDefinition.ColumnDefinitions)) {
                        context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                        return;
                    }
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    if (isCharLengthExceedMaxLengthInDefinitions(alterTableAddTableElementStatement.Definition.ColumnDefinitions)) {
                        context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                        return;
                    }
                    break;

                case AlterTableAlterColumnStatement alterTableAlterColumnStatement:
                    if (isCharLengthExceedMaxLengthInDefinitionsInDataType(alterTableAlterColumnStatement.DataType)) {
                        context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                        return;
                    }
                    break;
            }
        }

        public StringTypeShouldNotExceedMaxLengthRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class StringTypeShouldNoVarcharMaxRuleValidator : RuleValidator {
        public bool isVarcharMaxInDefinitions(IList<ColumnDefinition> columnDefinitions) {
            foreach (var columnDefinition in columnDefinitions) {
                if (isVarcharMaxInDataType(columnDefinition.DataType)) {
                    return true;
                }
            }
            return false;
        }

        public bool isVarcharMaxInDataType(DataTypeReference dataTypeReference) {
            if (dataTypeReference is ParameterizedDataTypeReference) {
                ParameterizedDataTypeReference parameterizedDataTypeReference = dataTypeReference as ParameterizedDataTypeReference;
                if (parameterizedDataTypeReference.Name.BaseIdentifier.Value.ToUpper() == "VARCHAR") {
                    foreach (var parameter in parameterizedDataTypeReference.Parameters) {
                        if (parameter is MaxLiteral) {
                            if (parameter.Value.ToUpper() == "MAX") {
                                return true;
                            }
                        }
                    }
                }
            }

            return false;
        }

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    TableDefinition tableDefinition = createTableStatement.Definition;
                    if (isVarcharMaxInDefinitions(tableDefinition.ColumnDefinitions)) {
                        context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                        return;
                    }
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    if (isVarcharMaxInDefinitions(alterTableAddTableElementStatement.Definition.ColumnDefinitions)) {
                        context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                        return;
                    }
                    break;

                case AlterTableAlterColumnStatement alterTableAlterColumnStatement:
                    if (isVarcharMaxInDataType(alterTableAlterColumnStatement.DataType)) {
                        context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                        return;
                    }
                    break;
            }
        }

        public StringTypeShouldNoVarcharMaxRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
