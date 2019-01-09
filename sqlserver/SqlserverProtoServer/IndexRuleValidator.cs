using System;
using System.Collections.Generic;
using System.Data.SqlClient;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;

namespace SqlserverProtoServer {
    public class NumberOfCompsiteIndexColumnsShouldNotExceedMaxRuleValidator : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public const int COMPOSITE_INDEX_MAX = 5;

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            int compositeIndexMax = 0;
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    if (createTableStatement.Definition.Indexes == null) {
                        break;
                    }
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
            if (compositeIndexMax > COMPOSITE_INDEX_MAX) {
                logger.Debug("composite index exceed max {0}", COMPOSITE_INDEX_MAX);
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public NumberOfCompsiteIndexColumnsShouldNotExceedMaxRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class NumberOfIndexesShouldNotExceedMaxRuleValidator : RuleValidator {
        public bool IsTest;
        public int ExpectIndexNumber;

        protected Logger logger = LogManager.GetCurrentClassLogger();

        public int GetIndexCounterFromTableDefinition(TableDefinition tableDefinition) {
            int indexCounter = tableDefinition.Indexes.Count;
            indexCounter += GetIndexCounterFromColumnDefinitions(tableDefinition.ColumnDefinitions);
            indexCounter += GetIndexCounterFromTableConstraints(tableDefinition.TableConstraints);
            return indexCounter;
        }

        public int GetIndexCounterFromColumnDefinitions(IList<ColumnDefinition> columnDefinitions) {
            if (columnDefinitions == null) {
                return 0;
            }

            int indexCounter = 0;
            foreach (var columnDefinition in columnDefinitions) {
                if (columnDefinition.Constraints == null) {
                    continue;
                }
                foreach (var constraint in columnDefinition.Constraints) {
                    if (constraint is UniqueConstraintDefinition) {
                        indexCounter++;
                    }
                }
            }
            return indexCounter;
        }

        public int GetIndexCounterFromTableConstraints(IList<ConstraintDefinition> constraintDefinitions) {
            if (constraintDefinitions == null) {
                return 0;
            }

            var indexCounter = 0;
            foreach (var constraint in constraintDefinitions) {
                if (constraint is UniqueConstraintDefinition) {
                    indexCounter++;
                }
            }
            return indexCounter;
        }

        public int GetNumberOfIndexesOnTable(SqlserverContext context, String tableName) {
            if (IsTest) {
                return ExpectIndexNumber;
            }

            int indexNumber = 0;
            String connectionString = context.GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                String queryString = String.Format("SELECT COUNT(*) AS Index_number FROM sys.indexes WHERE object_id=OBJECT_ID('{0}')", tableName);
                SqlCommand command = new SqlCommand(queryString, connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        indexNumber = (int)reader["Index_number"];
                    }
                } finally {
                    reader.Close();
                }
            }
            return indexNumber;
        }

        public const int INDEX_MAX_NUMBER = 5;
        public override void Check(SqlserverContext context, TSqlStatement statement) {
            int indexCounter = 0;
            String tableName = "";

            switch (statement) {
                case CreateTableStatement createTableStatement:
                    indexCounter = GetIndexCounterFromTableDefinition(createTableStatement.Definition);
                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    tableName = alterTableAddTableElementStatement.SchemaObjectName.BaseIdentifier.Value;
                    indexCounter = GetNumberOfIndexesOnTable(context, tableName);
                    indexCounter += GetIndexCounterFromColumnDefinitions(alterTableAddTableElementStatement.Definition.ColumnDefinitions);
                    break;

                case CreateIndexStatement createIndexStatement:
                    tableName = createIndexStatement.OnName.BaseIdentifier.Value;
                    indexCounter = GetNumberOfIndexesOnTable(context, tableName);
                    indexCounter++;
                    break;
            }

            if (indexCounter > INDEX_MAX_NUMBER) {
                logger.Debug("number of {0}'s indexes exceed max {1}", tableName, INDEX_MAX_NUMBER);
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public NumberOfIndexesShouldNotExceedMaxRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class DisableAddIndexForColumnsTypeBlob : RuleValidator {
        public bool IsTest;
        public Dictionary<String, String> ExpectColumnAndType;
        public Dictionary<String, bool> ExpectIndexColumns;

        protected Logger logger = LogManager.GetCurrentClassLogger();

        public bool hasIndexForColumnsTypeBlobInTableDefinition(TableDefinition tableDefinition) {
            if (tableDefinition == null) {
                return false;
            }

            // unique index
            foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                if (columnDefinition.Constraints == null) {
                    continue;
                }

                foreach (var constriant in columnDefinition.Constraints) {
                    if (constriant is UniqueConstraintDefinition) {
                        if (IsBlobType(columnDefinition.DataType)) {
                            return true;
                        }
                    }
                }
            }


            if (tableDefinition.TableConstraints != null) {
                foreach (var constraint in tableDefinition.TableConstraints) {
                    if (constraint is UniqueConstraintDefinition) {
                        UniqueConstraintDefinition uniqueConstraintDefinition = constraint as UniqueConstraintDefinition;
                        if (hasIndexForColumnsTypeBlob(tableDefinition, uniqueConstraintDefinition.Columns)) {
                            return true;
                        }
                    }
                }
            }

            // indexes
            if (tableDefinition.Indexes != null) {
                foreach (var index in tableDefinition.Indexes) {
                    if (hasIndexForColumnsTypeBlob(tableDefinition, index.Columns)) {
                        return true;
                    }
                }
            }

            return false;
        }

        public bool hasIndexForColumnsTypeBlob(TableDefinition tableDefinition, IList<ColumnWithSortOrder> columns) {
            if (tableDefinition == null || columns == null) {
                return false;
            }

            foreach (var column in columns) {
                ColumnReferenceExpression columnReferenceExpression = column.Column;
                foreach (var identifier in columnReferenceExpression.MultiPartIdentifier.Identifiers) {
                    foreach (var columnDefinition in tableDefinition.ColumnDefinitions) {
                        if (identifier.Value == columnDefinition.ColumnIdentifier.Value && IsBlobType(columnDefinition.DataType)) {
                            return true;
                        }
                    }
                }
            }
            return false;
        }

        public Dictionary<String, bool> GetIndexColumnsOnTable(SqlserverContext context, String tableName) {
            if (IsTest) {
                return ExpectIndexColumns;
            }

            Dictionary<String, bool> indexedColumns = new Dictionary<string, bool>();
            String connectionString = context.GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                String queryString = String.Format("SELECT COL_NAME(object_id, column_id) AS Column_name FROM sys.index_columns WHERE object_id=OBJECT_ID('{0}')", tableName);
                SqlCommand command = new SqlCommand(queryString, connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        indexedColumns[reader["Column_name"] as String] = true;
                    }
                } finally {
                    reader.Close();
                }
            }
            return indexedColumns;
        }

        public Dictionary<String, String> GetColumnAndTypeOnTable(SqlserverContext context, String tableName) {
            if (IsTest) {
                return ExpectColumnAndType;
            }

            Dictionary<String, String> columnTypes = new Dictionary<String, String>();
            String connectionString = context.GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                String queryString = String.Format("SELECT COLUMN_NAME AS Column_name, DATA_TYPE AS Data_type, CHARACTER_MAXIMUM_LENGTH as Max_length FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='{0}'", tableName);
                SqlCommand command = new SqlCommand(queryString, connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        var columnName = reader["Column_name"] as String;
                        var columnType = reader["Data_type"] as String;
                        if (columnType == "varbinary") {
                            if ((int)reader["Max_length"] == -1) {
                                columnType += "(max)";
                            }
                        }
                        columnTypes[columnName] = columnType;
                    }
                } finally {
                    reader.Close();
                }
            }
            return columnTypes;
        }

        public override void Check(SqlserverContext context, TSqlStatement statement) {
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
                    if (IsBlobType(alterTableAlterColumnStatement.DataType)) {
                        String columnName = alterTableAlterColumnStatement.ColumnIdentifier.Value;
                        Dictionary<String, bool> indexedColumns = GetIndexColumnsOnTable(context, alterTableAlterColumnStatement.SchemaObjectName.BaseIdentifier.Value);
                        if (indexedColumns.ContainsKey(columnName)) {
                            indexDataTypeIsBlob = true;
                        }
                    }
                    break;

                case CreateIndexStatement createIndexStatement:
                    Dictionary<String, String> columnTypes = GetColumnAndTypeOnTable(context, createIndexStatement.OnName.BaseIdentifier.Value);
                    foreach (var column in createIndexStatement.Columns) {
                        ColumnReferenceExpression columnReferenceExpression = column.Column;
                        foreach (var identifier in columnReferenceExpression.MultiPartIdentifier.Identifiers) {
                            if (!columnTypes.ContainsKey(identifier.Value)) {
                                continue;
                            }
                            if (columnTypes.ContainsKey(identifier.Value) && IsBlobTypeString(columnTypes[identifier.Value], null)) {
                                indexDataTypeIsBlob = true;
                            }
                        }
                    }
                    break;
            }
            if (indexDataTypeIsBlob) {
                logger.Debug("There is index for blob type column");
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
            }
        }

        public DisableAddIndexForColumnsTypeBlob(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class CheckNormalIndexPrefix : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            var indexes = new List<String>();
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    if (createTableStatement.Definition.Indexes == null) {
                        break;
                    }

                    foreach (var index in createTableStatement.Definition.Indexes) {
                        if (!index.Unique) {
                            indexes.Add(index.Name.Value);
                        }
                    }
                    break;

                case CreateIndexStatement createIndexStatement:
                    if (!createIndexStatement.Unique) {
                        indexes.Add(createIndexStatement.Name.Value);
                    }
                    break;
            }

            foreach (var index in indexes) {
                if (!index.StartsWith("idx_")) {
                    logger.Debug("index {0} should start with idx_", index);
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                    return;
                }
            }
        }

        public CheckNormalIndexPrefix(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }

    public class CheckUniqueIndexPrefix : RuleValidator {
        protected Logger logger = LogManager.GetCurrentClassLogger();

        public override void Check(SqlserverContext context, TSqlStatement statement) {
            var indexes = new List<String>();
            switch (statement) {
                case CreateTableStatement createTableStatement:
                    if (createTableStatement.Definition.TableConstraints != null) {
                        foreach (var constaint in createTableStatement.Definition.TableConstraints) {
                            if (constaint is UniqueConstraintDefinition) {
                                indexes.Add(constaint.ConstraintIdentifier.Value);
                            }
                        }
                    }

                    if (createTableStatement.Definition.Indexes != null) {
                        foreach (var index in createTableStatement.Definition.Indexes) {
                            if (index.Unique) {
                                indexes.Add(index.Name.Value);
                            }
                        }
                    }

                    break;

                case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                    if (alterTableAddTableElementStatement.Definition.TableConstraints != null) {
                        foreach (var constraint in alterTableAddTableElementStatement.Definition.TableConstraints) {
                            if (constraint is UniqueConstraintDefinition) {
                                indexes.Add(constraint.ConstraintIdentifier.Value);
                            }
                        }
                    }

                    break;

                case CreateIndexStatement createIndexStatement:
                    if (createIndexStatement.Unique) {
                        indexes.Add(createIndexStatement.Name.Value);
                    }
                    break;
            }

            foreach (var index in indexes) {
                if (!index.StartsWith("uniq_")) {
                    logger.Debug("index {0} should start with uniq_", index);
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                    return;
                }
            }
        }

        public CheckUniqueIndexPrefix(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) { }
    }
}
