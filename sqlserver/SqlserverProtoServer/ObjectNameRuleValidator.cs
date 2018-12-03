using System;
using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class ObjectNameRuleValidator : RuleValidator {
        public List<String> Names;

        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            switch (statement) {
                case CreateDatabaseStatement createDatabaseStatement:
                    Names.Add(createDatabaseStatement.DatabaseName.Value);
                    break;
                case CreateTableStatement createTableStatement:
                    // table name
                    Names.Add(createTableStatement.SchemaObjectName.BaseIdentifier.Value);
                    // column names
                    var columnDefinitions = createTableStatement.Definition.ColumnDefinitions;
                    foreach (var columnDefinition in columnDefinitions) {
                        Names.Add(columnDefinition.ColumnIdentifier.Value);
                    }
                    // index names
                    var tableConstraints = createTableStatement.Definition.TableConstraints;
                    foreach (var tableConstraint in tableConstraints) {
                        Names.Add(tableConstraint.ConstraintIdentifier.Value);
                    }
                    break;

                case AlterTableStatement alterTableStatement:
                    switch (alterTableStatement) {
                        // add column (with or without constraint)
                        case AlterTableAddTableElementStatement alterTableAddTableElementStatement:
                            foreach (var columnDefinition in alterTableAddTableElementStatement.Definition.ColumnDefinitions) {
                                Names.Add(columnDefinition.ColumnIdentifier.Value);
                                foreach (var constraint in columnDefinition.Constraints) {
                                    if (constraint.ConstraintIdentifier != null) {
                                        Names.Add(constraint.ConstraintIdentifier.Value);
                                    }
                                }
                            }
                            break;
                    }

                    break;

                // sp_rename [ @objname = ] 'object_name' , [ @newname = ] 'new_name' [ , [@objtype = ] 'object_type' ]
                // rename table: EXEC sp_rename 'schema1.table_old', 'table_new'
                // rename column: EXEC sp_rename 'schema1.table1.col_old', 'col_new', 'COLUMN'
                // rename index: EXEC sp_rename N'schema1.table1.index_old', N'index_new', N'INDEX'
                // rename constaint: EXEC sp_rename 'schema1.PK_constraint1', 'PK_constraint2'
                case ExecuteStatement executeStatement:
                    ExecutableEntity executableEntity = executeStatement.ExecuteSpecification.ExecutableEntity;
                    if (executableEntity is ExecutableProcedureReference) {
                        ExecutableProcedureReference executableProcedureReference = executableEntity as ExecutableProcedureReference;
                        if (executableProcedureReference.ProcedureReference != null) {
                            var procedureReference = executableProcedureReference.ProcedureReference;
                            if (procedureReference.ProcedureReference.Name.BaseIdentifier.Value == "sp_rename") {
                                var parameters = executableProcedureReference.Parameters;
                                if (parameters.Count >= 2) {
                                    var newName = ((StringLiteral)parameters[1].ParameterValue).Value;
                                    Names.Add(newName);
                                }
                            }
                        }
                    }

                    break;

                // add index
                case CreateIndexStatement createIndexStatement:
                    Names.Add(createIndexStatement.Name.Value);
                    break;
                default:
                    Console.WriteLine("type:{0}", statement.GetType());
                    break;
            }

            Show(Names, "name:{0}");
            reset();
        }

        public void reset() {
            Names = new List<String>();
        }

        public ObjectNameRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) {
            Names = new List<String>();
        }
    }

    public class ObjectNameMaxLengthRuleValidator: ObjectNameRuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            base.Check(context, statement);

            foreach (var name in Names) {
                if (name.Length > 64) {
                    context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage());
                }
            }

            reset();
        }

        public ObjectNameMaxLengthRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) {

        }
    }

    public class ObjectNameShouldNotContainsKeywordRuleValidator : ObjectNameRuleValidator {
        public override void Check(RuleValidatorContext context, TSqlStatement statement) {
            base.Check(context, statement);

            var invalidNames = new List<String>();
            foreach (var name in Names) {
                if (IsReserverdKeyword(name)) {
                    invalidNames.Add(name);
                }
            }

            if (invalidNames.Count > 0) {
                context.AdviseResultContext.AddAdviseResult(GetLevel(), GetMessage(String.Join(',', invalidNames)));
            }

            reset();
        }

        public ObjectNameShouldNotContainsKeywordRuleValidator(String name, String desc, String msg, RULE_LEVEL level) : base(name, desc, msg, level) {

        }

        public bool IsReserverdKeyword(String name) {
            return Keywords.Contains(name.ToUpper()) ? true : false;
        }

        public static List<String> Keywords = new List<string> {
            "ADD", "ALL", "ALTER", "AND", "ANY", "AS", "ASC", "AUTHORIZATION" , "BACKUP", "BEGIN", "BETWEEN", "BREAK", "BROWSE", "BULK", "BY" , "CASCADE", "CASE", "CHECK", "CHECKPOINT", "CLOSE", "CLUSTERED" , "COALESCE", "COLLATE", "COLUMN", "COMMIT", "COMPUTE", "CONSTRAINT" , "CONTAINS", "CONTAINSTABLE", "CONTINUE", "CONVERT", "CREATE" , "CROSS", "CURRENT", "CURRENT_DATE", "CURRENT_TIME" , "CURRENT_TIMESTAMP", "CURRENT_USER", "CURSOR", "DATABASE", "DBCC" , "DEALLOCATE", "DECLARE", "DEFAULT", "DELETE", "DENY", "DESC" , "DISK", "DISTINCT", "DISTRIBUTED", "DOUBLE", "DROP", "DUMMY" , "DUMP", "ELSE", "END", "ERRLVL", "ESCAPE", "EXCEPT", "EXEC" , "EXECUTE", "EXISTS", "EXIT", "FETCH", "FILE", "FILLFACTOR", "FOR" , "FOREIGN", "FREETEXT", "FREETEXTTABLE", "FROM", "FULL", "FUNCTION" , "GOTO", "GRANT", "GROUP", "HAVING", "HOLDLOCK", "IDENTITY" , "IDENTITY_INSERT", "IDENTITYCOL", "IF", "IN", "INDEX", "INNER" , "INSERT", "INTERSECT", "INTO", "IS", "JOIN", "KEY", "KILL", "LEFT" , "LIKE", "LINENO", "LOAD", "NATIONAL", "NOCHECK", "NONCLUSTERED" , "NOT", "NULL", "NULLIF", "OF", "OFF", "OFFSETS", "ON", "OPEN" , "OPENDATASOURCE", "OPENQUERY", "OPENROWSET", "OPENXML", "OPTION" , "OR", "ORDER", "OUTER", "OVER", "PERCENT", "PLAN", "PRECISION" , "PRIMARY", "PRINT", "PROC", "PROCEDURE", "PUBLIC", "RAISERROR" , "READ", "READTEXT", "RECONFIGURE", "REFERENCES", "REPLICATION" , "RESTORE", "RESTRICT", "RETURN", "REVOKE", "RIGHT", "ROLLBACK" , "ROWCOUNT", "ROWGUIDCOL", "RULE", "SAVE", "SCHEMA", "SELECT" , "SESSION_USER", "SET", "SETUSER", "SHUTDOWN", "SOME", "STATISTICS" , "SYSTEM_USER", "TABLE", "TEXTSIZE", "THEN", "TO", "TOP", "TRANSACTION" , "TRIGGER", "TRUNCATE", "TSEQUAL", "UNION", "UNIQUE", "UPDATE" , "UPDATETEXT", "USE", "USER", "VALUES", "VARYING", "VIEW" , "WAITFOR", "WHEN", "WHERE", "WHILE", "WITH", "WRITETEXT"
        };
    }
}
