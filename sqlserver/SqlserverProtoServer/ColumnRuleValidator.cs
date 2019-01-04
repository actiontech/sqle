using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;
using System.Collections.Generic;

namespace SqlserverProtoServer {
    public class CheckColumnWithoutDefault : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public bool HasDefaultValueForNoneColumnBlob(IList<ColumnDefinition> columnDefinitions) {
            foreach (var columnDefinition in columnDefinitions) {
                if (IsBlobType(columnDefinition.DataType) || columnDefinition.IdentityOptions != null) {
                    continue;
                }

                if (columnDefinition.DefaultConstraint == null) {
                    return false;
                }
            }

            return true;
        }

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            IList<ColumnDefinition> columnDefinitions = null;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    columnDefinitions = createTableStatement.Definition.ColumnDefinitions;
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    columnDefinitions = alterTableAddTableElementStatement.Definition.ColumnDefinitions;
                    break;
            }

            if (!HasDefaultValueForNoneColumnBlob(columnDefinitions)) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public CheckColumnWithoutDefault(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class CheckColumnTimestampWithoutDefaut : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public bool HasDefaultValueForColumnTimestamp(IList<ColumnDefinition> columnDefinitions) {
            var timeTypes = new List<String>() {
                "DATE", "DATETIME", "DATETIME2", "DATETIMEOFFSET", "SMALLDATETIME", "TIME",
            };
            foreach (var columnDefinition in columnDefinitions) {
                var typeName = columnDefinition.DataType.Name.BaseIdentifier.Value;
                if (timeTypes.Contains(typeName) && columnDefinition.DefaultConstraint == null) {
                    logger.Debug("column {0} of time type should contain default value", columnDefinition.ColumnIdentifier.Value);
                    return false;
                }
            }

            return true;
        }

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            IList<ColumnDefinition> columnDefinitions = null;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    columnDefinitions = createTableStatement.Definition.ColumnDefinitions;
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    columnDefinitions = alterTableAddTableElementStatement.Definition.ColumnDefinitions;
                    break;
            }

            if (!HasDefaultValueForColumnTimestamp(columnDefinitions)) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public CheckColumnTimestampWithoutDefaut(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class CheckColumnBlobNotNull : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public bool NullableForColumnBlob(IList<ColumnDefinition> columnDefinitions) {
            var nullable = true;
            foreach (var columnDefinition in columnDefinitions) {
                if (IsBlobType(columnDefinition.DataType)) {
                    foreach (var constraint in columnDefinition.Constraints) {
                        if (constraint is NullableConstraintDefinition) {
                            var nullableConstraint = constraint as NullableConstraintDefinition;
                            if (!nullableConstraint.Nullable) {
                                nullable = false;
                            }
                        }
                    }
                }
            }

            return nullable;
        }

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            IList<ColumnDefinition> columnDefinitions = null;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    columnDefinitions = createTableStatement.Definition.ColumnDefinitions;
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    columnDefinitions = alterTableAddTableElementStatement.Definition.ColumnDefinitions;
                    break;
            }

            if (!NullableForColumnBlob(columnDefinitions)) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public CheckColumnBlobNotNull(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class CheckColumnBlobDefaultNotNull : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public bool DefaultIsNullForColumnBlob(IList<ColumnDefinition> columnDefinitions) {
            var defaultIsNull = true;
            foreach (var columnDefinition in columnDefinitions) {
                if (IsBlobType(columnDefinition.DataType) && columnDefinition.DefaultConstraint != null) {
                    if (!(columnDefinition.DefaultConstraint.Expression is NullLiteral)) {
                        defaultIsNull = false;
                    }
                }
            }

            return defaultIsNull;
        }
        public override void Check(SqlserverContext context, TSqlStatement statement) {
            IList<ColumnDefinition> columnDefinitions = null;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    columnDefinitions = createTableStatement.Definition.ColumnDefinitions;
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    columnDefinitions = alterTableAddTableElementStatement.Definition.ColumnDefinitions;
                    break;
            }

            if (!DefaultIsNullForColumnBlob(columnDefinitions)) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public CheckColumnBlobDefaultNotNull(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
