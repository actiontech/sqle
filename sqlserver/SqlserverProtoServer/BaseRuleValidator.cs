using System;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;
using System.Collections.Generic;

namespace SqlserverProtoServer {
    public class BaseRuleValidator : RuleValidator {
        public override void Check(SqlserverContext context, TSqlStatement statement) {
            var isInvalid = false;
            var logger = LogManager.GetCurrentClassLogger();
            switch (statement) {
                case UseStatement useStatement:
                    isInvalid = IsInvalidUse(logger, context, useStatement);
                    break;

                case CreateTableStatement createTableStatement:
                    isInvalid = IsInvalidCreateTableStatement(logger, context, createTableStatement);
                    break;

                case AlterTableStatement alterTableStatement:
                    isInvalid = IsInvalidAlterTableStatement(logger, context, alterTableStatement);
                    break;

                case DropTableStatement dropTableStatement:
                    isInvalid = IsInvalidDropTableStatement(logger, context, dropTableStatement);
                    break;

                case CreateDatabaseStatement createDatabaseStatement:
                    isInvalid = IsInvalidCreateDatabaseStatement(logger, context, createDatabaseStatement);
                    break;

                case DropDatabaseStatement dropDatabaseStatement:
                    isInvalid = IsInvalidDropDatabaseStatement(logger, context, dropDatabaseStatement);
                    break;

                case CreateIndexStatement createIndexStatement:
                    isInvalid = IsInvalidCreateIndexStatement(logger, context, createIndexStatement);
                    break;

                case DropIndexStatement dropIndexStatement:
                    isInvalid = IsInvalidDropIndexStatement(logger, context, dropIndexStatement);
                    break;

                case InsertStatement insertStatement:
                    isInvalid = IsInvalidInsertStatement(logger, context, insertStatement);
                    break;

                case UpdateStatement updateStatement:
                    isInvalid = IsInvalidUpdateStatement(logger, context, updateStatement);
                    break;

                case DeleteStatement deleteStatement:
                    isInvalid = IsInvalidDeleteStatement(logger, context, deleteStatement);
                    break;

                case SelectStatement selectStatement:
                    isInvalid = IsInvalidSelectStatement(logger, context, selectStatement);
                    break;
            }

            if (isInvalid) {
                context.AdviseResultContext.SetBaseRuleStatus(AdviseResultContext.BASE_RULE_FAILED);
                return;
            }
        }

        public bool IsInvalidUse(Logger logger, SqlserverContext context, UseStatement statement) {
            var databaseName = statement.DatabaseName.Value;
            if (!context.DatabaseExists(logger, databaseName)) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, databaseName));
                logger.Info("database {0} should exist", databaseName);
                return true;
            }
            return false;
        }

        public bool IsInvalidCreateTableStatement(Logger logger, SqlserverContext context, CreateTableStatement statement) {
            var isInvalid = false;
            var schemaObject = statement.SchemaObjectName;
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(statement.SchemaObjectName, out String databaseName, out String schemaName, out String tableName);
            // database should exist
            {
                if (!context.DatabaseExists(logger, databaseName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, databaseName));
                    isInvalid = true;
                    logger.Info("database {0} should exist", databaseName);
                }
            }

            // schema should exist
            {
                if (!context.SchemaExists(logger, schemaName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, schemaName));
                    isInvalid = true;
                    logger.Info("schema {0} should exist", schemaName);
                }
            }

            // table should not exist
            {
                if (context.TableExists(logger, databaseName, schemaName, tableName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_EXIST_MSG, tableName));
                    isInvalid = true;
                    logger.Info("table {0} should not exist", tableName);
                }
            }

            {
                // no duplicate column
                var columnNames = new List<String>();
                var pkCounter = 0;
                foreach (var columnDefinition in statement.Definition.ColumnDefinitions) {
                    columnNames.Add(columnDefinition.ColumnIdentifier.Value);
                    if (columnDefinition.Constraints == null) {
                        continue;
                    }
                    foreach (var columnConstraint in columnDefinition.Constraints) {
                        if (columnConstraint is UniqueConstraintDefinition) {
                            var uniqueConstraintDefinition = columnConstraint as UniqueConstraintDefinition;
                            if (uniqueConstraintDefinition.IsPrimaryKey) {
                                pkCounter += 1;
                            }
                        }
                    }
                }
                logger.Info("create table columns:{0}", String.Join(",", columnNames));
                var duplicatedColumns = GetDuplicatedNames(columnNames);
                if (duplicatedColumns.Count > 0) {
                    logger.Info("table {0} has duplicated columns: {1}", tableName, String.Join(",", duplicatedColumns));
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DUPLICATE_COLUMN_ERROR_MSG, String.Join(",", duplicatedColumns)));
                    isInvalid = true;
                }

                // no duplicate index
                var indexNames = new List<String>();
                var indexColumnNames = new List<String>();
                if (statement.Definition.Indexes != null) {
                    foreach (var indexDefinition in statement.Definition.Indexes) {
                        indexNames.Add(indexDefinition.Name.Value);
                        foreach (var column in indexDefinition.Columns) {
                            ColumnReferenceExpression columnReferenceExpression = column.Column;
                            var identifiers = columnReferenceExpression.MultiPartIdentifier.Identifiers;
                            if (identifiers.Count > 0) {
                                var identifier = identifiers[identifiers.Count - 1];
                                indexColumnNames.Add(identifier.Value);
                            }
                        }
                    }
                }
                logger.Info("create table indexes:{0}", String.Join(",", indexNames));
                var duplicatedIndexes = GetDuplicatedNames(indexNames);
                if (duplicatedIndexes.Count > 0) {
                    logger.Info("table {0} has duplicated index: {1}", tableName, String.Join(",", duplicatedIndexes));
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DUPLICATE_INDEX_ERROR_MSG, String.Join(",", duplicatedIndexes)));
                    isInvalid = true;
                }

                // index column should be table column
                var notExistKeyColumns = new List<String>();
                foreach (var indexColumnName in indexColumnNames) {
                    if (!columnNames.Contains(indexColumnName)) {
                        notExistKeyColumns.Add(indexColumnName);
                    }
                }
                if (notExistKeyColumns.Count > 0) {
                    logger.Info("key columns {0} not exist in table {1}", String.Join(",", notExistKeyColumns, tableName));
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.KEY_COLUMN_NOT_EXIST_MSG, String.Join(",", notExistKeyColumns)));
                    isInvalid = true;
                }

                // no duplicate constaint
                var constraintNames = new List<String>();
                var constraintColumnNames = new List<String>();
                if (statement.Definition.TableConstraints != null) {
                    foreach (var constraintDefinition in statement.Definition.TableConstraints) {
                        if (constraintDefinition.ConstraintIdentifier == null) {
                            continue;
                        }
                        constraintNames.Add(constraintDefinition.ConstraintIdentifier.Value);
                        if (constraintDefinition is UniqueConstraintDefinition) {
                            var uniqueConstaintDefinition = constraintDefinition as UniqueConstraintDefinition;
                            if (uniqueConstaintDefinition.IsPrimaryKey) {
                                pkCounter += 1;
                            }
                            foreach (var column in uniqueConstaintDefinition.Columns) {
                                ColumnReferenceExpression columnReferenceExpression = column.Column;
                                var identifiers = columnReferenceExpression.MultiPartIdentifier.Identifiers;
                                if (identifiers.Count > 0) {
                                    var identifier = identifiers[identifiers.Count - 1];
                                    constraintColumnNames.Add(identifier.Value);
                                }
                            }
                        }
                    }
                }
                logger.Info("create table constraints:{0}, constraint columns:{1}", String.Join(",", constraintNames), String.Join(",", constraintColumnNames));
                var duplicatedConstaints = GetDuplicatedNames(constraintNames);
                if (duplicatedConstaints.Count > 0) {
                    logger.Info("table {0} has duplicated constaint: {1}", tableName, String.Join(",", duplicatedConstaints));
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DUPLICATE_CONSTAINT_ERROR_MSG, String.Join(",", duplicatedConstaints)));
                    isInvalid = true;
                } else {
                    var existedConstraints = ExistedConstraints(context, databaseName, constraintNames);
                    if (existedConstraints.Count > 0) {
                        logger.Info("existed constaints: {0}", String.Join(",", existedConstraints));
                        context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DUPLICATE_CONSTAINT_ERROR_MSG, String.Join(",", existedConstraints)));
                        isInvalid = true;
                    }
                }

                // constraint column should be table column
                var notExsitConstaintColumns = new List<String>();
                foreach (var constaintColumnName in constraintColumnNames) {
                    if (!columnNames.Contains(constaintColumnName)) {
                        notExsitConstaintColumns.Add(constaintColumnName);
                    }
                }
                if (notExsitConstaintColumns.Count > 0) {
                    logger.Info("constraint columns {0} not exist in table {1}", String.Join(",", notExsitConstaintColumns, tableName));
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.CONSTRAINT_COLUMN_NOT_EXIST_MSG, String.Join(",", notExsitConstaintColumns)));
                    isInvalid = true;
                }

                // no multiple primary key constraint
                if (pkCounter > 1) {
                    logger.Info("table {0} has {1} PRIMARY KEY definition", tableName, pkCounter);
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, DefaultRules.PRIMARY_KEY_MULTI_ERROR_MSG);
                    isInvalid = true;
                }
            }

            return isInvalid;
        }

        public bool IsInvalidAlterTableStatement(Logger logger, SqlserverContext context, AlterTableStatement statement) {
            var isInvalid = false;
            var schemaObject = statement.SchemaObjectName;
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(schemaObject, out String databaseName, out String schemaName, out String tableName);
            // database should exist
            {
                if (!context.DatabaseExists(logger, databaseName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, databaseName));
                    isInvalid = true;
                    logger.Info("database {0} should exist", databaseName);
                }
            }

            // schema should exist
            {
                if (!context.SchemaExists(logger, schemaName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, schemaName));
                    isInvalid = true;
                    logger.Info("schema {0} should exist", schemaName);
                }
            }

            // table should exist
            {
                if (!context.TableExists(logger, databaseName, schemaName, tableName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_NOT_EXIST_MSG, tableName));
                    isInvalid = true;
                    logger.Info("table {0} should exist", tableName);
                    return isInvalid;
                }
            }

            var needExistsColumnName = new List<String>();
            var needNotExistsColumnName = new List<String>();
            var needExistsIndexName = new List<String>();
            var needNotExistsIndexName = new List<String>();
            var needExistsConstraintName = new List<String>();
            var needNotExistsConstaintName = new List<String>();
            var isAddPrimaryKey = false;

            switch (statement) {
                case AlterTableAddTableElementStatement addTableElementStatement:
                    var definition = addTableElementStatement.Definition;
                    var addColumnDefinitions = context.GetTableColumnDefinitions(logger, databaseName, schemaName, tableName);

                    if (definition.ColumnDefinitions.Count > 0) {
                        foreach (var columnDefinition in definition.ColumnDefinitions) {
                            var addColumnName = columnDefinition.ColumnIdentifier.Value;
                            logger.Info("addColumnName:{0}", addColumnName);
                            if (addColumnDefinitions.ContainsKey(addColumnName)) {
                                needNotExistsColumnName.Add(addColumnName);
                            }
                        }
                    }

                    if (definition.Indexes.Count > 0) {
                        var addIndexDefinitions = context.GetTableIndexDefinitions(logger, databaseName, schemaName, tableName);
                        foreach (var index in definition.Indexes) {
                            var addIndexName = index.Name.Value;
                            logger.Info("addIndexName:{0}", addIndexName);
                            if (addIndexDefinitions.ContainsKey(addIndexName)) {
                                needNotExistsIndexName.Add(addIndexName);
                            }

                            foreach (var column in index.Columns) {
                                var identifiers = column.Column.MultiPartIdentifier.Identifiers;
                                if (identifiers.Count > 0) {
                                    var identifier = identifiers[identifiers.Count - 1];
                                    if (!addColumnDefinitions.ContainsKey(identifier.Value)) {
                                        needExistsColumnName.Add(identifier.Value);
                                    }
                                }
                            }
                        }
                    }

                    if (definition.TableConstraints.Count > 0) {
                        var constraintDefinitions = context.GetTableConstraintDefinitions(logger, databaseName, schemaName, tableName);
                        foreach (var tableConstaint in definition.TableConstraints) {
                            if (tableConstaint.ConstraintIdentifier == null) {
                                continue;
                            }
                            var constaintName = tableConstaint.ConstraintIdentifier.Value;
                            logger.Info("constraint name:{0}", constaintName);
                            if (constraintDefinitions.ContainsKey(constaintName)) {
                                needNotExistsConstaintName.Add(constaintName);
                            }

                            if (tableConstaint is UniqueConstraintDefinition) {
                                var uniqueConstaint = tableConstaint as UniqueConstraintDefinition;
                                if (uniqueConstaint.IsPrimaryKey) {
                                    isAddPrimaryKey = true;
                                }

                                foreach (var column in uniqueConstaint.Columns) {
                                    var identifiers = column.Column.MultiPartIdentifier.Identifiers;
                                    if (identifiers.Count > 0) {
                                        var identifier = identifiers[identifiers.Count - 1];
                                        if (!addColumnDefinitions.ContainsKey(identifier.Value)) {
                                            needExistsColumnName.Add(identifier.Value);
                                        }
                                    }
                                }
                            }
                        }
                    }

                    if (isAddPrimaryKey) {
                        var primaryKeys = context.GetPrimaryKeys(databaseName, schemaName, tableName);
                        if (primaryKeys.Count > 0) {
                            context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, DefaultRules.PRIMARY_KEY_EXIST_MSG);
                            isInvalid = true;
                            logger.Info("table {0} has primary key", tableName);
                        }
                    }

                    break;

                case AlterTableAlterColumnStatement alterColumnStatement:
                    var columnName = alterColumnStatement.ColumnIdentifier.Value;
                    var alterColumnDefinitions = context.GetTableColumnDefinitions(logger, databaseName, schemaName, tableName);
                    logger.Info("alert column:{0}", columnName);
                    if (!alterColumnDefinitions.ContainsKey(columnName)) {
                        needNotExistsColumnName.Add(columnName);
                    }

                    break;

                case AlterTableAlterIndexStatement alterIndexStatement:
                    var indexName = alterIndexStatement.IndexIdentifier.Value;
                    logger.Info("alert index:{0}", indexName);
                    var alterIndexDefinitions = context.GetTableIndexDefinitions(logger, databaseName, schemaName, tableName);
                    if (!alterIndexDefinitions.ContainsKey(indexName)) {
                        needExistsIndexName.Add(indexName);
                    }

                    break;

                case AlterTableDropTableElementStatement dropTableElementStatement:
                    foreach (var dropTableElement in dropTableElementStatement.AlterTableDropTableElements) {
                        if (!dropTableElement.IsIfExists) {
                            switch (dropTableElement.TableElementType) {
                                case TableElementType.Column:
                                    logger.Info("drop column:{0}", dropTableElement.Name.Value);
                                    needExistsColumnName.Add(dropTableElement.Name.Value);
                                    break;

                                case TableElementType.Constraint:
                                    logger.Info("drop constraint:{0}", dropTableElement.Name.Value);
                                    needExistsConstraintName.Add(dropTableElement.Name.Value);
                                    break;

                                case TableElementType.Index:
                                    logger.Info("drop index:{0}", dropTableElement.Name.Value);
                                    needExistsIndexName.Add(dropTableElement.Name.Value);
                                    break;
                            }
                        }
                    }
                    break;
            }

            if (needExistsColumnName.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.COLUMN_NOT_EXIST_MSG, String.Join(",", needExistsColumnName)));
                isInvalid = true;
                logger.Info("column {0} should exist", String.Join(",", needExistsColumnName));
            }

            if (needNotExistsColumnName.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.COLUMN_EXIST_MSG, String.Join(",", needNotExistsColumnName)));
                isInvalid = true;
                logger.Info("column {0} should not exist", String.Join(",", needNotExistsColumnName));
            }

            if (needExistsIndexName.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.INDEX_NOT_EXIST_MSG, String.Join(",", needExistsIndexName)));
                isInvalid = true;
                logger.Info("index {0} should exist", String.Join(",", needExistsIndexName));
            }

            if (needNotExistsIndexName.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.INDEX_EXIST_MSG, String.Join(",", needNotExistsIndexName)));
                isInvalid = true;
                logger.Info("index {0} should not exist", String.Join(",", needNotExistsIndexName));
            }

            if (needExistsConstraintName.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.CONSTRAINT_NOT_EXIST_MSG, String.Join(",", needExistsConstraintName)));
                isInvalid = true;
                logger.Info("constraint {0} should exist", String.Join(",", needExistsConstraintName));
            }

            if (needNotExistsConstaintName.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.CONSTRAINT_EXIST_MSG, String.Join(",", needNotExistsConstaintName)));
                isInvalid = true;
                logger.Info("constaint {0} should not exist", String.Join(",", needNotExistsConstaintName));
            }

            return isInvalid;
        }

        public bool IsInvalidDropTableStatement(Logger logger, SqlserverContext context, DropTableStatement statement) {
            if (statement.IsIfExists) {
                return false;
            }

            var isInvalid = false;
            var needExistsDatabaseNames = new List<String>();
            var needExistsSchemaNames = new List<String>();
            var needExistsTableNames = new List<String>();
            foreach (var droppedObject in statement.Objects) {
                context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(droppedObject, out String databaseName, out String schemaName, out String tableName);
                if (!context.DatabaseExists(logger, databaseName)) {
                    needExistsDatabaseNames.Add(databaseName);
                }
                if (!context.SchemaExists(logger, schemaName)) {
                    needExistsSchemaNames.Add(schemaName);
                }
                if (!context.TableExists(logger, databaseName, schemaName, tableName)) {
                    needExistsTableNames.Add(tableName);
                }
            }

            // database should exist
            if (needExistsDatabaseNames.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, String.Join(",", needExistsDatabaseNames)));
                isInvalid = true;
                logger.Info("database {0} should exist", String.Join(",", needExistsDatabaseNames));
            }
            // schema should exist
            if (needExistsSchemaNames.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, String.Join(",", needExistsSchemaNames)));
                isInvalid = true;
                logger.Info("schema {0} should exist", String.Join(",", needExistsSchemaNames));
            }
            // table should exist
            if (needExistsTableNames.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_NOT_EXIST_MSG, String.Join(",", needExistsTableNames)));
                isInvalid = true;
                logger.Info("table {0} should not exist", String.Join(",", needExistsTableNames));
            }


            return isInvalid;
        }

        public bool IsInvalidCreateDatabaseStatement(Logger logger, SqlserverContext context, CreateDatabaseStatement statement) {
            var databaseName = statement.DatabaseName.Value;
            if (context.DatabaseExists(logger, databaseName)) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_EXIST_MSG, databaseName));
                logger.Info("database {0} should not exist", databaseName);
                return true;
            }
            return false;
        }

        public bool IsInvalidDropDatabaseStatement(Logger logger, SqlserverContext context, DropDatabaseStatement statement) {
            if (statement.IsIfExists) {
                return false;
            }

            var notExistDatabaseNames = new List<String>();
            foreach (var database in statement.Databases) {
                var databaseName = database.Value;
                if (!context.DatabaseExists(logger, databaseName)) {
                    notExistDatabaseNames.Add(databaseName);
                }
            }

            if (notExistDatabaseNames.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, String.Join(",", notExistDatabaseNames)));
                logger.Info("database {0} should exist", String.Join(",", notExistDatabaseNames));
                return true;
            }

            return false;
        }

        public bool IsInvalidCreateIndexStatement(Logger logger, SqlserverContext context, CreateIndexStatement statement) {
            var isInvalid = false;
            var schemaObject = statement.OnName;
            context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(schemaObject, out String databaseName, out String schemaName, out String tableName);
            // database should exist
            {
                if (!context.DatabaseExists(logger, databaseName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, databaseName));
                    isInvalid = true;
                    logger.Info("database {0} should exist", databaseName);
                }
            }

            // schema should exist
            {
                if (!context.SchemaExists(logger, schemaName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, schemaName));
                    isInvalid = true;
                    logger.Info("schema {0} should exist", schemaName);
                }
            }

            // table should exist
            {
                if (!context.TableExists(logger, databaseName, schemaName, tableName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_NOT_EXIST_MSG, tableName));
                    isInvalid = true;
                    logger.Info("table {0} should exist", tableName);
                    return isInvalid;
                }
            }

            var indexDefinitions = context.GetTableIndexDefinitions(logger, databaseName, schemaName, tableName);
            if (indexDefinitions.ContainsKey(statement.Name.Value)) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.INDEX_EXIST_MSG, statement.Name.Value));
                isInvalid = true;
                logger.Info("index {0} should not exist", statement.Name.Value);
            }
            logger.Info("create index:{0}", statement.Name.Value);

            var columnDefinitions = context.GetTableColumnDefinitions(logger, databaseName, schemaName, tableName);
            var needExistsColumns = new List<String>();
            foreach (var column in statement.Columns) {
                var identifiers = column.Column.MultiPartIdentifier.Identifiers;
                if (identifiers.Count > 0) {
                    var identifier = identifiers[identifiers.Count - 1];
                    if (!columnDefinitions.ContainsKey(identifier.Value)) {
                        needExistsColumns.Add(identifier.Value);
                    }
                }
            }
            if (needExistsColumns.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.COLUMN_NOT_EXIST_MSG, String.Join(",", needExistsColumns)));
                isInvalid = true;
                logger.Info("column {0} should exist", String.Join(",", needExistsColumns));
            }

            return isInvalid;
        }

        public bool IsInvalidDropIndexStatement(Logger logger, SqlserverContext context, DropIndexStatement statement) {
            if (statement.IsIfExists) {
                return false;
            }

            var isInvalid = false;
            var needExistsIndex = new List<String>();
            foreach (var dropIndexClauseBase in statement.DropIndexClauses) {
                if (dropIndexClauseBase is DropIndexClause) {
                    var dropIndexCaluse = dropIndexClauseBase as DropIndexClause;
                    logger.Info("drop index:{0}", dropIndexCaluse.Index.Value);
                    var schemaObject = dropIndexCaluse.Object;
                    context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(schemaObject, out String databaseName, out String schemaName, out String tableName);
                    // database should exist
                    {
                        if (!context.DatabaseExists(logger, databaseName)) {
                            context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, databaseName));
                            isInvalid = true;
                            logger.Info("database {0} should exist", databaseName);
                        }
                    }

                    // schema should exist
                    {
                        if (!context.SchemaExists(logger, schemaName)) {
                            context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, schemaName));
                            isInvalid = true;
                            logger.Info("schema {0} should exist", schemaName);
                        }
                    }

                    // table should exist
                    {
                        if (!context.TableExists(logger, databaseName, schemaName, tableName)) {
                            context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_NOT_EXIST_MSG, tableName));
                            isInvalid = true;
                            logger.Info("table {0} should exist", tableName);
                            continue;
                        }
                    }

                    var indexDefinitions = context.GetTableIndexDefinitions(logger, databaseName, schemaName, tableName);
                    if (!indexDefinitions.ContainsKey(dropIndexCaluse.Index.Value)) {
                        needExistsIndex.Add(dropIndexCaluse.Index.Value);
                    }
                }
            }

            if (needExistsIndex.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.INDEX_NOT_EXIST_MSG, String.Join(",", needExistsIndex)));
                isInvalid = true;
                logger.Info("index {0} should exist", String.Join(",", needExistsIndex));
            }

            return isInvalid;
        }

        public bool IsInvalidInsertStatement(Logger logger, SqlserverContext context, InsertStatement statement) {
            var isInvalid = false;
            var tableReference = statement.InsertSpecification.Target;
            if (tableReference is NamedTableReference) {
                var namedTableReference = tableReference as NamedTableReference;
                var schemaObjectName = namedTableReference.SchemaObject;
                context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(schemaObjectName, out String databaseName, out String schemaName, out String tableName);
                if (!context.DatabaseExists(logger, databaseName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, databaseName));
                    logger.Info("databases {0} should exists", databaseName);
                    return true;
                }
                if (!context.SchemaExists(logger, schemaName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, schemaName));
                    logger.Info("schemas {0} shold exists", schemaName);
                    return true;
                }
                if (!context.TableExists(logger, databaseName, schemaName, tableName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_NOT_EXIST_MSG, tableName));
                    logger.Info("tables {0} should exists", tableName);
                    return true;
                }

                if (statement.InsertSpecification.Columns != null) {
                    var columnDefinitions = context.GetTableColumnDefinitions(logger, databaseName, schemaName, tableName);
                    var needExistsColumns = new List<String>();
                    var insertColumns = new List<String>();

                    foreach (var column in statement.InsertSpecification.Columns) {
                        var identifiers = column.MultiPartIdentifier.Identifiers;
                        if (identifiers.Count > 0) {
                            var identifier = identifiers[identifiers.Count - 1];
                            insertColumns.Add(identifier.Value);
                            if (!columnDefinitions.ContainsKey(identifier.Value)) {
                                needExistsColumns.Add(identifier.Value);
                            }
                        }
                    }
                    if (insertColumns.Count == 0) {
                        foreach (var columnDefinition in columnDefinitions) {
                            insertColumns.Add(columnDefinition.Key);
                        }
                    }
                    logger.Info("insert columns:{0}", String.Join(",", insertColumns));

                    if (needExistsColumns.Count > 0) {
                        context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.COLUMN_NOT_EXIST_MSG, String.Join(",", needExistsColumns)));
                        logger.Info("columns {0} should exists", String.Join(",", needExistsColumns));
                        isInvalid = true;
                    }

                    var duplicatedColumns = GetDuplicatedNames(insertColumns);
                    if (duplicatedColumns.Count > 0) {
                        context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DUPLICATE_COLUMN_ERROR_MSG, String.Join(",", duplicatedColumns)));
                        logger.Info("duplicated columns {0} in insert statement", String.Join(",", duplicatedColumns));
                        isInvalid = true;
                    }

                    if (statement.InsertSpecification.InsertSource is ValuesInsertSource) {
                        var valuesInsertSource = statement.InsertSpecification.InsertSource as ValuesInsertSource;
                        foreach (var rowValue in valuesInsertSource.RowValues) {
                            if (rowValue.ColumnValues.Count != insertColumns.Count) {
                                context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, DefaultRules.NOT_MATCH_VALUES_AND_COLUMNS);
                                logger.Info("column not match values for {0}", tableName);
                                isInvalid = true;
                                break;
                            }
                        }
                    }
                }
            }

            return isInvalid;
        }

        public bool IsInvalidUpdateStatement(Logger logger, SqlserverContext context, UpdateStatement statement) {
            var isInvalid = false;
            TableReference tableReference = null;
            if (statement.UpdateSpecification.FromClause != null && statement.UpdateSpecification.FromClause.TableReferences != null) {
                tableReference = statement.UpdateSpecification.FromClause.TableReferences[0];
            } else {
                tableReference = statement.UpdateSpecification.Target;
            }
            if (tableReference != null && tableReference is NamedTableReference) {
                var namedTableReference = tableReference as NamedTableReference;
                var schemaObjectName = namedTableReference.SchemaObject;
                context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(schemaObjectName, out String databaseName, out String schemaName, out String tableName);
                if (!context.DatabaseExists(logger, databaseName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, databaseName));
                    logger.Info("databases {0} should exists", databaseName);
                    return true;
                }
                if (!context.SchemaExists(logger, schemaName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, schemaName));
                    logger.Info("schemas {0} shold exists", schemaName);
                    return true;
                }
                if (!context.TableExists(logger, databaseName, schemaName, tableName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_NOT_EXIST_MSG, tableName));
                    logger.Info("tables {0} should exists", tableName);
                    return true;
                }

                var columnDefinitions = context.GetTableColumnDefinitions(logger, databaseName, schemaName, tableName);
                var needExistsColumns = new List<String>();
                var updateColumns = new List<String>();
                foreach (var setClause in statement.UpdateSpecification.SetClauses) {
                    if (setClause is AssignmentSetClause) {
                        var assignmentSetClause = setClause as AssignmentSetClause;
                        var identifiers = assignmentSetClause.Column.MultiPartIdentifier.Identifiers;
                        if (identifiers.Count > 0) {
                            var identifier = identifiers[identifiers.Count - 1];
                            updateColumns.Add(identifier.Value);
                            if (!columnDefinitions.ContainsKey(identifier.Value)) {
                                needExistsColumns.Add(identifier.Value);
                            }
                        }
                    }
                }
                logger.Info("update columns:{0}", String.Join(",", updateColumns));

                if (needExistsColumns.Count > 0) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.COLUMN_NOT_EXIST_MSG, String.Join(",", needExistsColumns)));
                    logger.Info("columns {0} should exists", String.Join(",", needExistsColumns));
                    isInvalid = true;
                }

                var duplicatedColumns = GetDuplicatedNames(updateColumns);
                if (duplicatedColumns.Count > 0) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DUPLICATE_COLUMN_ERROR_MSG, String.Join(",", duplicatedColumns)));
                    logger.Info("duplicated columns {0} in insert statement", String.Join(",", duplicatedColumns));
                    isInvalid = true;
                }
            }

            return isInvalid;
        }

        public bool IsInvalidDeleteStatement(Logger logger, SqlserverContext context, DeleteStatement statement) {
            TableReference tableReference = null;
            if (statement.DeleteSpecification.FromClause != null && statement.DeleteSpecification.FromClause.TableReferences != null) {
                tableReference = statement.DeleteSpecification.FromClause.TableReferences[0];
            } else {
                tableReference = statement.DeleteSpecification.Target;
            }
            if (tableReference != null && tableReference is NamedTableReference) {
                var namedTableReference = tableReference as NamedTableReference;
                var schemaObjectName = namedTableReference.SchemaObject;
                context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(schemaObjectName, out String databaseName, out String schemaName, out String tableName);
                if (!context.DatabaseExists(logger, databaseName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, databaseName));
                    logger.Info("databases {0} should exists", databaseName);
                    return true;
                }
                if (!context.SchemaExists(logger, schemaName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, schemaName));
                    logger.Info("schemas {0} shold exists", schemaName);
                    return true;
                }
                if (!context.TableExists(logger, databaseName, schemaName, tableName)) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_NOT_EXIST_MSG, tableName));
                    logger.Info("tables {0} should exists", tableName);
                    return true;
                }
            }
            return false;
        }

        public bool IsInvalidSelectStatement(Logger logger, SqlserverContext context, SelectStatement statement) {
            var isInvalid = false;
            if (statement.QueryExpression is QuerySpecification) {
                var querySpecfication = statement.QueryExpression as QuerySpecification;
                var needExistsDatabaseNames = new List<String>();
                var needExistsSchemaNames = new List<String>();
                var needExistsTableNames = new List<String>();
                foreach (var tableReference in querySpecfication.FromClause.TableReferences) {
                    if (tableReference is NamedTableReference) {
                        var namedTableReference = tableReference as NamedTableReference;
                        var schemaObjectName = namedTableReference.SchemaObject;
                        context.GetDatabaseNameAndSchemaNameAndTableNameFromSchemaObjectName(schemaObjectName, out String databaseName, out String schemaName, out String tableName);
                        if (!context.DatabaseExists(logger, databaseName)) {
                            needExistsDatabaseNames.Add(databaseName);
                            continue;
                        }
                        if (!context.SchemaExists(logger, schemaName)) {
                            needExistsSchemaNames.Add(schemaName);
                            continue;
                        }
                        if (!context.TableExists(logger, databaseName, schemaName, tableName)) {
                            needExistsTableNames.Add(tableName);
                            continue;
                        }
                    }
                }

                if (needExistsDatabaseNames.Count > 0) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.DATABASE_NOT_EXIST_MSG, String.Join(",", needExistsDatabaseNames)));
                    logger.Info("databases {0} should exists", String.Join(",", needExistsDatabaseNames));
                    isInvalid = true;
                }

                if (needExistsSchemaNames.Count > 0) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.SCHEMA_NOT_EXIST_MSG, String.Join(",", needExistsSchemaNames)));
                    logger.Info("schemas {0} shold exists", String.Join(",", needExistsSchemaNames));
                    isInvalid = true;
                }

                if (needExistsTableNames.Count > 0) {
                    context.AdviseResultContext.AddAdviseResult(RULE_LEVEL.ERROR, String.Format(DefaultRules.TABLE_NOT_EXIST_MSG, String.Join(",", needExistsTableNames)));
                    logger.Info("tables {0} should exists", String.Join(",", needExistsTableNames));
                    isInvalid = true;
                }
            }
            return isInvalid;
        }

        public List<String> GetDuplicatedNames(List<String> names) {
            var ret = new List<String>();
            var nameMap = new Dictionary<String, int>();
            foreach (var name in names) {
                if (nameMap.ContainsKey(name)) {
                    nameMap[name] += 1;
                } else {
                    nameMap[name] = 1;
                }
            }

            foreach (var namePair in nameMap) {
                if (namePair.Value > 1) {
                    ret.Add(namePair.Key);
                }
            }

            return ret;
        }

        public List<String> ExistedConstraints(SqlserverContext context, String databaseName, List<String> constraintNames) {
            var ret = new List<String>();
            var existConstraints = context.GetConstraintNames(databaseName);
            foreach (var constraintName in constraintNames) {
                foreach (var existConstraint in existConstraints) {
                    if (constraintName == existConstraint) {
                        ret.Add(constraintName);
                        break;
                    }
                }
            }

            return ret;
        }

        public BaseRuleValidator() {}
    }
}