using System;
using Xunit;
using SqlserverProtoServer;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Collections.Generic;
using System.IO;
using SqlserverProto;
using NLog;

namespace SqlServerProtoServerTest {
    public class ValidatorTest {
        public StatementList ParseStatementList(string text) {
            var parser = new TSql130Parser(false);
            var reader = new StringReader(text);
            IList<ParseError> errorList;
            var statementList = parser.ParseStatementList(reader, out errorList);
            if (errorList.Count > 0) {
                throw new ArgumentException(String.Format("parse sql {0} error: {1}", text, errorList.ToString()));
            }

            return statementList;
        }

        [Fact]
        public void ParseStatementListTestWithGOTest() {
            try {
                StatementList statementList;
                var sqlserverImpl = new SqlServerServiceImpl();
                statementList = sqlserverImpl.ParseStatementList(LogManager.GetCurrentClassLogger(), "", 
                                                                 @"CREATE TABLE table1(col1 INT)
                                                                   GO
                                                                   INSERT INTO table1 VALUES(1)");
                Assert.True(statementList.Statements.Count > 0);
                statementList = sqlserverImpl.ParseStatementList(LogManager.GetCurrentClassLogger(), "", 
                                                                       @"CREATE TABLE table1(col1 INT);
                                                                       GO;
                                                                       INSERT INTO table1 VALUES(1);");

            } catch (Exception e) {
                Console.WriteLine(e.Message);
                Console.WriteLine(e.StackTrace);
            }
        }

        private void IsInvalidUse() {
            // database not exist
            StatementList statementList = ParseStatementList("USE database1");
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseExist = false;
            Console.WriteLine("IsInvalidUse");
            foreach (var statment in statementList.Statements) {
                var invalid = new BaseRuleValidator().IsInvalidUse(LogManager.GetCurrentClassLogger(), context, statment as UseStatement);
                Assert.True(invalid == true);
            }
            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
            Assert.Equal("[error]database database1 不存在", context.AdviseResultContext.GetMessage());
        }

        private void IsInvalidCreateTable() {
            {
                StatementList statementList = ParseStatementList("CREATE TABLE database1.schema1.tbl1 (" +
                                                                 "col1 INT PRIMARY KEY," +
                                                                 "col1 INT," +
                                                                 "col2 INT," +
                                                                 "INDEX IX1 (col1)," +
                                                                 "INDEX IX1 (col3)," +
                                                                 "CONSTRAINT PK_Constaint PRIMARY KEY (col1)," +
                                                                 "CONSTRAINT UN_1 UNIQUE (col1)," +
                                                                 "CONSTRAINT UN_1 UNIQUE (col4)" +
                                                                 ");");
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseExist = false;
                context.ExpectSchemaExist = false;
                context.ExpectTableExist = true;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "tbl1";
                Console.WriteLine();
                Console.WriteLine("IsInvalidCreateTable");

                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidCreateTableStatement(LogManager.GetCurrentClassLogger(), context, statment as CreateTableStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]database database1 不存在\n[error]schema schema1 不存在\n[error]表 tbl1 已存在\n[error]字段名 col1 重复\n[error]索引名 IX1 重复\n[error]索引字段 col3 不存在\n[error]约束名 UN_1 重复\n[error]约束字段 col4 不存在\n[error]主键只能设置一个", context.AdviseResultContext.GetMessage());
            }
        }

        private void IsInvalidAlterTable() {
            {
                StatementList statementList = ParseStatementList("ALTER TABLE database1.schema1.table1 ADD col1 INT NOT NULL");
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseExist = false;
                context.ExpectSchemaExist = false;
                context.ExpectTableExist = false;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";
                Console.WriteLine();
                Console.WriteLine("IsInvalidAlterTable");

                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidAlterTableStatement(LogManager.GetCurrentClassLogger(), context, statment as AlterTableStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]database database1 不存在\n[error]schema schema1 不存在\n[error]表 table1 不存在", context.AdviseResultContext.GetMessage());
            }

            {
                // init data
                StatementList initStatementList = ParseStatementList("CREATE TABLE database1.schema1.table1(" +
                                                                     "col1 INT NOT NULL, " +
                                                                     "col2 INT NOT NULL, " +
                                                                     "CONSTRAINT PK_1 PRIMARY KEY (col1)," +
                                                                     "CONSTRAINT UN_1 UNIQUE (col2)," +
                                                                     "INDEX IX_1 (col2))");
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseExist = true;
                context.ExpectSchemaExist = true;
                context.ExpectTableExist = true;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";
                context.UpdateContext(LogManager.GetCurrentClassLogger(), initStatementList.Statements[0]);

                // add column
                {
                    StatementList statementList = ParseStatementList("ALTER TABLE database1.schema1.table1 ADD col1 INT NOT NULL");
                    foreach (var statment in statementList.Statements) {
                        var invalid = new BaseRuleValidator().IsInvalidAlterTableStatement(LogManager.GetCurrentClassLogger(), context, statment as AlterTableStatement);
                        Assert.True(invalid == true);
                    }
                    Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                    Assert.Equal("[error]字段 col1 已存在", context.AdviseResultContext.GetMessage());
                    context.AdviseResultContext.ResetAdviseResult();
                }

                // add index
                {
                    StatementList statementList = ParseStatementList("ALTER TABLE database1.schema1.table1 ADD INDEX IX_1 (col3)");
                    foreach (var statment in statementList.Statements) {
                        var invalid = new BaseRuleValidator().IsInvalidAlterTableStatement(LogManager.GetCurrentClassLogger(), context, statment as AlterTableStatement);
                        Assert.True(invalid == true);
                    }
                    Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                    Assert.Equal("[error]字段 col3 不存在\n[error]索引 IX_1 已存在", context.AdviseResultContext.GetMessage());
                    context.AdviseResultContext.ResetAdviseResult();
                }

                // add constraint
                {
                    StatementList statementList = ParseStatementList("ALTER TABLE database1.schema1.table1 ADD CONSTRAINT UN_1 UNIQUE (col3)");
                    foreach (var statment in statementList.Statements) {
                        var invalid = new BaseRuleValidator().IsInvalidAlterTableStatement(LogManager.GetCurrentClassLogger(), context, statment as AlterTableStatement);
                        Assert.True(invalid == true);
                    }
                    Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                    Assert.Equal("[error]字段 col3 不存在\n[error]约束 UN_1 已存在", context.AdviseResultContext.GetMessage());
                    context.AdviseResultContext.ResetAdviseResult();
                }

                // add primary key constraint
                {
                    StatementList statementList = ParseStatementList("ALTER TABLE database1.schema1.table1 ADD CONSTRAINT PK_1 PRIMARY KEY (col3)");
                    foreach (var statment in statementList.Statements) {
                        var invalid = new BaseRuleValidator().IsInvalidAlterTableStatement(LogManager.GetCurrentClassLogger(), context, statment as AlterTableStatement);
                        Assert.True(invalid == true);
                    }
                    Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                    Assert.Equal("[error]已经存在主键，不能再添加\n[error]字段 col3 不存在\n[error]约束 PK_1 已存在", context.AdviseResultContext.GetMessage());
                    context.AdviseResultContext.ResetAdviseResult();
                }
            }
        }

        private void IsInvalidDropTable() {
            StatementList statementList = ParseStatementList("DROP TABLE database1.schema1.table1;");
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseExist = false;
            context.ExpectSchemaExist = false;
            context.ExpectTableExist = false;
            context.ExpectDatabaseName = "database1";
            context.ExpectSchemaName = "schema1";
            context.ExpectTableName = "table1";
            Console.WriteLine();
            Console.WriteLine("IsInvalidDropTable");

            foreach (var statment in statementList.Statements) {
                var invalid = new BaseRuleValidator().IsInvalidDropTableStatement(LogManager.GetCurrentClassLogger(), context, statment as DropTableStatement);
                Assert.True(invalid == true);
            }
            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
            Assert.Equal("[error]database database1 不存在\n[error]schema schema1 不存在\n[error]表 table1 不存在", context.AdviseResultContext.GetMessage());
        }

        private void IsInvalidCreateDatabase() {
            // database not exist
            StatementList statementList = ParseStatementList("CREATE DATABASE database1");
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseExist = true;
            Console.WriteLine();
            Console.WriteLine("IsInvalidCreateDatabase");

            foreach (var statment in statementList.Statements) {
                var invalid = new BaseRuleValidator().IsInvalidCreateDatabaseStatement(LogManager.GetCurrentClassLogger(), context, statment as CreateDatabaseStatement);
                Assert.True(invalid == true);
            }
            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
            Assert.Equal("[error]database database1 已存在", context.AdviseResultContext.GetMessage());
        }

        private void IsInvalidDropDatabase() {
            // database exist
            StatementList statementList = ParseStatementList("DROP DATABASE database1");
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseExist = false;
            Console.WriteLine();
            Console.WriteLine("IsInvalidDropDatabase");

            foreach (var statment in statementList.Statements) {
                var invalid = new BaseRuleValidator().IsInvalidDropDatabaseStatement(LogManager.GetCurrentClassLogger(), context, statment as DropDatabaseStatement);
                Assert.True(invalid == true);
            }
            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
            Assert.Equal("[error]database database1 不存在", context.AdviseResultContext.GetMessage());
        }

        private void IsInvalidCreateIndex() {
            Console.WriteLine();
            Console.WriteLine("IsInvalidCreateIndex");

            {
                StatementList statementList = ParseStatementList("CREATE INDEX IX_1 ON database1.schema1.table1 (col1);");
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseExist = false;
                context.ExpectSchemaExist = false;
                context.ExpectTableExist = false;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidCreateIndexStatement(LogManager.GetCurrentClassLogger(), context, statment as CreateIndexStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]database database1 不存在\n[error]schema schema1 不存在\n[error]表 table1 不存在", context.AdviseResultContext.GetMessage());
            }

            {
                // init data
                StatementList initStatementList = ParseStatementList("CREATE TABLE database1.schema1.table1(" +
                                                                     "col1 INT NOT NULL, " +
                                                                     "col2 INT NOT NULL, " +
                                                                     "CONSTRAINT PK_1 PRIMARY KEY (col1)," +
                                                                     "CONSTRAINT UN_1 UNIQUE (col2)," +
                                                                     "INDEX IX_1 (col2))");
                CreateTableStatement initStatement = initStatementList.Statements[0] as CreateTableStatement;
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseExist = true;
                context.ExpectSchemaExist = true;
                context.ExpectTableExist = true;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";
                context.UpdateContext(LogManager.GetCurrentClassLogger(), initStatementList.Statements[0]);

                {
                    StatementList statementList = ParseStatementList("CREATE INDEX IX_1 ON database1.schema1.table1 (col3);");
                    foreach (var statment in statementList.Statements) {
                        var invalid = new BaseRuleValidator().IsInvalidCreateIndexStatement(LogManager.GetCurrentClassLogger(), context, statment as CreateIndexStatement);
                        Assert.True(invalid == true);
                    }
                    Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                    Assert.Equal("[error]索引 IX_1 已存在\n[error]字段 col3 不存在", context.AdviseResultContext.GetMessage());
                }
            }
        }

        private void IsInvalidDropIndex() {
            Console.WriteLine();
            Console.WriteLine("IsInvalidDropIndex");

            // init data
            StatementList initStatementList = ParseStatementList("CREATE TABLE database1.schema1.table1(" +
                                                                 "col1 INT NOT NULL, " +
                                                                 "col2 INT NOT NULL, " +
                                                                 "CONSTRAINT PK_1 PRIMARY KEY (col1)," +
                                                                 "CONSTRAINT UN_1 UNIQUE (col2)," +
                                                                 "INDEX IX_1 (col2))");
            CreateTableStatement initStatement = initStatementList.Statements[0] as CreateTableStatement;
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseExist = true;
            context.ExpectSchemaExist = true;
            context.ExpectTableExist = true;
            context.ExpectDatabaseName = "database1";
            context.ExpectSchemaName = "schema1";
            context.ExpectTableName = "table1";
            context.UpdateContext(LogManager.GetCurrentClassLogger(), initStatementList.Statements[0]);

            {
                StatementList statementList = ParseStatementList("DROP INDEX IX_2 ON schema1.table1;");
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidDropIndexStatement(LogManager.GetCurrentClassLogger(), context, statment as DropIndexStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]索引 IX_2 不存在", context.AdviseResultContext.GetMessage());
            }
        }

        private void IsInvalidInsert() {
            Console.WriteLine("IsInvalidInsert");

            {
                StatementList statementList = ParseStatementList("INSERT INTO databse1.schema1.table1 VALUES(1);");
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";


                context.ExpectDatabaseExist = false;
                context.ExpectSchemaExist = true;
                context.ExpectTableExist = true;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidInsertStatement(LogManager.GetCurrentClassLogger(), context, statment as InsertStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]database database1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();


                context.ExpectDatabaseExist = true;
                context.ExpectSchemaExist = false;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidInsertStatement(LogManager.GetCurrentClassLogger(), context, statment as InsertStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]schema schema1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();

                context.ExpectSchemaExist = true;
                context.ExpectTableExist = false;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidInsertStatement(LogManager.GetCurrentClassLogger(), context, statment as InsertStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]表 table1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();
            }

            {
                // init data
                StatementList initStatementList = ParseStatementList("CREATE TABLE database1.schema1.table1(" +
                                                                     "col1 INT NOT NULL, " +
                                                                     "col2 INT NOT NULL, " +
                                                                     "CONSTRAINT PK_1 PRIMARY KEY (col1)," +
                                                                     "CONSTRAINT UN_1 UNIQUE (col2)," +
                                                                     "INDEX IX_1 (col2))");
                CreateTableStatement initStatement = initStatementList.Statements[0] as CreateTableStatement;
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseExist = true;
                context.ExpectSchemaExist = true;
                context.ExpectTableExist = true;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";
                context.UpdateContext(LogManager.GetCurrentClassLogger(), initStatementList.Statements[0]);

                {
                    StatementList statementList = ParseStatementList("INSERT INTO database1.schema1.table1(col1, col2, col2, col3) VALUES(2, 3);");
                    foreach (var statment in statementList.Statements) {
                        var invalid = new BaseRuleValidator().IsInvalidInsertStatement(LogManager.GetCurrentClassLogger(), context, statment as InsertStatement);
                        Assert.True(invalid == true);
                    }
                    Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                    Assert.Equal("[error]字段 col3 不存在\n[error]字段名 col2 重复\n[error]指定的值列数与字段列数不匹配", context.AdviseResultContext.GetMessage());
                }
            }
        }

        private void IsInvalidUpdate() {
            Console.WriteLine();
            Console.WriteLine("IsInvalidUpdate");

            {
                StatementList statementList = ParseStatementList("UPDATE databse1.schema1.table1 SET col1=1;");
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";


                context.ExpectDatabaseExist = false;
                context.ExpectSchemaExist = true;
                context.ExpectTableExist = true;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidUpdateStatement(LogManager.GetCurrentClassLogger(), context, statment as UpdateStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]database database1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();


                context.ExpectDatabaseExist = true;
                context.ExpectSchemaExist = false;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidUpdateStatement(LogManager.GetCurrentClassLogger(), context, statment as UpdateStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]schema schema1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();

                context.ExpectSchemaExist = true;
                context.ExpectTableExist = false;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidUpdateStatement(LogManager.GetCurrentClassLogger(), context, statment as UpdateStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]表 table1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();
            }

            {
                // init data
                StatementList initStatementList = ParseStatementList("CREATE TABLE database1.schema1.table1(" +
                                                                     "col1 INT NOT NULL, " +
                                                                     "col2 INT NOT NULL, " +
                                                                     "CONSTRAINT PK_1 PRIMARY KEY (col1)," +
                                                                     "CONSTRAINT UN_1 UNIQUE (col2)," +
                                                                     "INDEX IX_1 (col2))");
                CreateTableStatement initStatement = initStatementList.Statements[0] as CreateTableStatement;
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseExist = true;
                context.ExpectSchemaExist = true;
                context.ExpectTableExist = true;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";
                context.UpdateContext(LogManager.GetCurrentClassLogger(), initStatementList.Statements[0]);

                {
                    StatementList statementList = ParseStatementList("UPDATE database1.schema1.table1 SET col1=1, col1=2, col2=2, col3=3;");
                    foreach (var statment in statementList.Statements) {
                        var invalid = new BaseRuleValidator().IsInvalidUpdateStatement(LogManager.GetCurrentClassLogger(), context, statment as UpdateStatement);
                        Assert.True(invalid == true);
                    }
                    Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                    Assert.Equal("[error]字段 col3 不存在\n[error]字段名 col1 重复", context.AdviseResultContext.GetMessage());
                }
            }
        }

        private void IsInvalidDelete() {
            Console.WriteLine("IsInvalidDelete");

            {
                StatementList statementList = ParseStatementList("DELETE FROM databse1.schema1.table1;");
                var context = new SqlserverContext(new SqlserverMeta());
                context.IsTest = true;
                context.ExpectDatabaseName = "database1";
                context.ExpectSchemaName = "schema1";
                context.ExpectTableName = "table1";


                context.ExpectDatabaseExist = false;
                context.ExpectSchemaExist = true;
                context.ExpectTableExist = true;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidDeleteStatement(LogManager.GetCurrentClassLogger(), context, statment as DeleteStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]database database1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();


                context.ExpectDatabaseExist = true;
                context.ExpectSchemaExist = false;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidDeleteStatement(LogManager.GetCurrentClassLogger(), context, statment as DeleteStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]schema schema1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();

                context.ExpectSchemaExist = true;
                context.ExpectTableExist = false;
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidDeleteStatement(LogManager.GetCurrentClassLogger(), context, statment as DeleteStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]表 table1 不存在", context.AdviseResultContext.GetMessage());
                context.AdviseResultContext.ResetAdviseResult();
            }
        }

        private void IsInvalidSelect() {
            Console.WriteLine();
            Console.WriteLine("IsInvalidSelect");

            // init data
            StatementList initStatementList = ParseStatementList("CREATE TABLE database1.schema1.table1(" +
                                                                 "col1 INT NOT NULL, " +
                                                                 "col2 INT NOT NULL, " +
                                                                 "CONSTRAINT PK_1 PRIMARY KEY (col1)," +
                                                                 "CONSTRAINT UN_1 UNIQUE (col2)," +
                                                                 "INDEX IX_1 (col2))");
            CreateTableStatement initStatement = initStatementList.Statements[0] as CreateTableStatement;
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseExist = true;
            context.ExpectSchemaExist = true;
            context.ExpectTableExist = false;
            context.ExpectDatabaseName = "database1";
            context.ExpectSchemaName = "schema1";
            context.ExpectTableName = "table1";
            context.UpdateContext(LogManager.GetCurrentClassLogger(), initStatementList.Statements[0]);

            {
                StatementList statementList = ParseStatementList("SELECT col1 FROM database1.schema1.table1;");
                foreach (var statment in statementList.Statements) {
                    var invalid = new BaseRuleValidator().IsInvalidSelectStatement(LogManager.GetCurrentClassLogger(), context, statment as SelectStatement);
                    Assert.True(invalid == true);
                }
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                Assert.Equal("[error]表 table1 不存在", context.AdviseResultContext.GetMessage());
            }
        }

        [Fact]
        public void BaseValidatorTest() {
            IsInvalidUse();
            IsInvalidCreateTable();
            IsInvalidAlterTable();
            IsInvalidDropTable();
            IsInvalidCreateDatabase();
            IsInvalidDropDatabase();
            IsInvalidCreateIndex();
            IsInvalidDropIndex();
            IsInvalidInsert();
            IsInvalidUpdate();
            IsInvalidDelete();
            IsInvalidSelect();
        }

        private void MyAssert(String ruleName, String text, String expectRuleLevel, String expectErrMsg) {
            var statementList = ParseStatementList(text);
            foreach (var statment in statementList.Statements) {
                var context = new SqlserverContext(new SqlserverMeta());
                var validator = DefaultRules.RuleValidators[ruleName];
                validator.Check(context, statment);
                Assert.Equal(expectRuleLevel, context.AdviseResultContext.GetLevel());
                Assert.Equal(expectErrMsg, context.AdviseResultContext.GetMessage());
            }
        }

        [Fact]
        public void RuleValidatorTest() {
            //DDL_CHECK_OBJECT_NAME_LENGTH
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_OBJECT_NAME_LENGTH");
                try {
                    foreach (var text in new String[]{
                        "CREATE DATABASE database1",
                        "CREATE TABLE table1 (col1 INT, INDEX IX_col1 NONCLUSTERED(col1))",
                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(20) NULL, col2 INT NULL",
                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(20) NULL CONSTRAINT constraint1 UNIQUE",
                        "ALTER TABLE schema1.table1 WITH NOCHECK ADD CONSTRAINT constraint1 CHECK (col1 > 1)",
                        "EXEC sp_rename @objectname='schema1.table_old', @newname='table_new'",
                        "EXEC sp_rename 'schema1.table1.col_old', 'col_new', 'COLUMN'",
                        "EXEC sp_rename N'schema1.table1.index_old', N'index_new', N'INDEX'",
                        "CREATE UNIQUE INDEX index1 ON schema1.table1(col1)"
                    }) {
                        Console.WriteLine("text:{0}", text);
                        MyAssert(DefaultRules.DDL_CHECK_OBJECT_NAME_LENGTH,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                             "");
                    }

                    foreach (var text in new String[]{
                        "CREATE DATABASE database123456789012345678901234567890123456789012345678901234567890",
                        "CREATE TABLE table123456789012345678901234567890123456789012345678901234567890 (col1 INT, INDEX IX_col1 NONCLUSTERED(col1))",
                        "CREATE TABLE table1 (col12345678901234567890123456789012345678901234567890123456789012 INT, INDEX IX_col1 NONCLUSTERED(col123456789012345678901234567890123456789012345678901234567890))",
                        "CREATE TABLE table1 (col1 INT, INDEX IX_col123456789012345678901234567890123456789012345678901234567890 NONCLUSTERED(col1))",
                        "ALTER TABLE schema1.table1 ADD col12345678901234567890123456789012345678901234567890123456789012 VARCHAR(20) NULL, col2 INT NULL",
                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(20) NULL CONSTRAINT constraint123456789012345678901234567890123456789012345678901234567890 UNIQUE",
                        "ALTER TABLE schema1.table1 WITH NOCHECK ADD CONSTRAINT constraint123456789012345678901234567890123456789012345678901234567890 CHECK (col1 > 1)",
                        "EXEC sp_rename @objectname='schema1.table_old', @newname='table_new123456789012345678901234567890123456789012345678901234567890'",
                        "EXEC sp_rename 'schema1.table1.col_old', 'col_new123456789012345678901234567890123456789012345678901234567890', 'COLUMN'",
                        "EXEC sp_rename N'schema1.table1.index_old', N'index_new123456789012345678901234567890123456789012345678901234567890', N'INDEX'",
                        "CREATE UNIQUE INDEX index123456789012345678901234567890123456789012345678901234567890 ON schema1.table1(col1)"

                    }) {
                        Console.WriteLine("text:{0}", text);
                        MyAssert(DefaultRules.DDL_CHECK_OBJECT_NAME_LENGTH,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]表名、列名、索引名的长度不能大于64字节");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            //DDL_CHECK_PK_NOT_EXIST
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_PK_NOT_EXIST");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT PRIMARY KEY CLUSTERED" +
                        ")",
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT IDENTITY(1, 1)," +
                            "col2 INT NOT NULL," +
                            "CONSTRAINT PK_constraint PRIMARY KEY(col1, col2)" +
                        ")"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_PK_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                             "");
                    }

                    foreach (var text in new String[]{
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT" +
                        ")"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_PK_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]表必须有主键");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            //DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT IDENTITY(1,1) PRIMARY KEY CLUSTERED)",

                        "CREATE TABLE schema1.table1(" +
                            "col1 INT IDENTITY(1, 1)," +
                            "col2 INT NOT NULL," +
                            "CONSTRAINT PK_constraint PRIMARY KEY(col1, col2)" +
                        ")"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                             "");
                    }

                    foreach (var text in new String[]{
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT" +
                        ")",

                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL," +
                            "col2 INT NOT NULL," +
                            "CONSTRAINT PK_constraint PRIMARY KEY(col1, col2)" +
                        ")"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]主键建议使用自增");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            //DDL_CHECK_COLUMN_VARCHAR_MAX
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_COLUMN_VARCHAR_MAX");
                try {
                    foreach (var text in new String[]{

                        "CREATE TABLE schema1.table1(" +
                            "col1 VARCHAR(MAX) PRIMARY KEY CLUSTERED" +
                        ")",
                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(MAX)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 VARCHAR(MAX)"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COLUMN_VARCHAR_MAX,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]禁止使用 varchar(max)");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            //DDL_CHECK_COLUMN_CHAR_LENGTH
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_COLUMN_CHAR_LENGTH");
                try {
                    foreach (var text in new String[]{

                        "CREATE TABLE schema1.table1(" +
                            "col1 CHAR(60) PRIMARY KEY CLUSTERED" +
                        ")",
                        "ALTER TABLE schema1.table1 ADD col1 CHAR(60)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 CHAR(60)"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COLUMN_CHAR_LENGTH,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]char长度大于20时，必须使用varchar类型");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }

                try {
                    foreach (var text in new String[]{

                        "CREATE TABLE schema1.table1(" +
                            "col1 CHAR(20) PRIMARY KEY CLUSTERED" +
                        ")",
                        "ALTER TABLE schema1.table1 ADD col1 CHAR(20)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 CHAR(20)"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COLUMN_CHAR_LENGTH,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                             "");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            //DDL_DISABLE_FK
            {
                Console.WriteLine();
                Console.WriteLine("DDL_DISABLE_FK");
                try {
                    foreach (var text in new String[] {

                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL REFERENCES schema1.table2(col11)" +
                        ")",

                        "create table marks(s_id int not null,test_no int not null,marks int not null default(0),primary key(s_id,test_no),foreign key(s_id) references student(s_id))",

                        "ALTER TABLE schema1.table1 ADD CONSTRAINT FK_fk1 FOREIGN KEY(col1) REFERENCES schema1.table2(col11)"

                    }) {
                        Console.WriteLine("{0}", text);
                        MyAssert(DefaultRules.DDL_DISABLE_FK,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]禁止使用外键");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            //DDL_CHECK_INDEX_COUNT
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_INDEX_COUNT");
                try {
                    foreach (var text in new String[] {
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL UNIQUE," +
                            "CONSTRAINT IX_index1 UNIQUE(col1)," +
                            "INDEX IX_index2 (col1)" +
                        ")",

                        "ALTER TABLE schema1.table100 ADD col1 VARCHAR(20) NULL CONSTRAINT constraint1 UNIQUE",

                        "CREATE INDEX IX_index1 ON table100(col1)"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        var statementList = ParseStatementList(text);
                        foreach (var statment in statementList.Statements) {
                            var context = new SqlserverContext(new SqlserverMeta());
                            var validator = DefaultRules.RuleValidators[DefaultRules.DDL_CHECK_INDEX_COUNT];
                            var checkIndexCountValidator = validator as NumberOfIndexesShouldNotExceedMaxRuleValidator;
                            checkIndexCountValidator.IsTest = true;
                            checkIndexCountValidator.ExpectIndexNumber = 0;

                            validator.Check(context, statment);
                            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL), context.AdviseResultContext.GetLevel());
                            Assert.Equal("", context.AdviseResultContext.GetMessage());
                        }
                    }

                    foreach (var text in new String[] {
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL PRIMARY KEY," +
                            "col2 INT NOT NULL," +
                            "col3 INT NOT NULL," +
                            "col4 INT NOT NULL," +
                            "CONSTRAINT IX_index1 UNIQUE(col1)," +
                            "INDEX IX_index2 (col1)," +
                            "INDEX IX_index3 (col2)," +
                            "INDEX IX_index4 (col3)," +
                            "INDEX IX_index5 (col4)" +
                        ")",
                        "ALTER TABLE table1 ADD col10 VARCHAR(20) NULL CONSTRAINT constraint1 UNIQUE",
                        "CREATE INDEX IX_index1 ON table1(col1)"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        var statementList = ParseStatementList(text);
                        foreach (var statment in statementList.Statements) {
                            var context = new SqlserverContext(new SqlserverMeta());
                            var validator = DefaultRules.RuleValidators[DefaultRules.DDL_CHECK_INDEX_COUNT];
                            var checkIndexCountValidator = validator as NumberOfIndexesShouldNotExceedMaxRuleValidator;
                            checkIndexCountValidator.IsTest = true;
                            checkIndexCountValidator.ExpectIndexNumber = 10;

                            validator.Check(context, statment);
                            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NOTICE), context.AdviseResultContext.GetLevel());
                            Assert.Equal("[notice]索引个数建议不超过5个", context.AdviseResultContext.GetMessage());
                        }
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            //DDL_CHECK_COMPOSITE_INDEX_MAX
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_COMPOSITE_INDEX_MAX");
                try {
                    foreach (var text in new String[] {

                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL UNIQUE," +
                            "col2 INT NOT NULL," +
                            "INDEX IX_index2 (col1, col2)" +
                        ")",
                        "CREATE INDEX IX_index1 ON table1(col1, col2)"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COMPOSITE_INDEX_MAX,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                            "");
                    }

                    foreach (var text in new String[] {
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL," +
                            "col2 INT NOT NULL," +
                            "col3 INT NOT NULL," +
                            "col4 INT NOT NULL," +
                            "col5 INT NOT NULL," +
                            "col6 INT NOT NULL," +
                            "INDEX IX_index2 (col1, col2, col3, col4, col5, col6)" +
                        ")",
                        "CREATE INDEX IX_index1 ON table1(col1, col2, col3, col4, col5, col6)"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COMPOSITE_INDEX_MAX,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NOTICE),
                            "[notice]复合索引的列数量不建议超过5个");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            // DDL_CHECK_OBJECT_NAME_USING_KEYWORD
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_OBJECT_NAME_USING_KEYWORD");
                try {
                    foreach (var text in new String[] {
                        "CREATE DATABASE database1",
                        "CREATE TABLE table1 (table1 INT, INDEX index1 NONCLUSTERED(col1))",
                        "ALTER TABLE schema1.table1 ADD table1 VARCHAR(20) NULL",
                        "ALTER TABLE schema1.table1 ADD table1 VARCHAR(20) NULL CONSTRAINT constraint1 UNIQUE",
                        "ALTER TABLE schema1.table1 WITH NOCHECK ADD CONSTRAINT constraint1 CHECK (column1 > 1)",
                        "EXEC sp_rename @objectname='schema1.table_old', @newname='table1'",
                        "EXEC sp_rename 'schema1.table1.col_old', 'COLUMN1', 'COLUMN'",
                        "EXEC sp_rename N'schema1.table1.index_old', N'INDEX1', N'INDEX'",
                        "CREATE UNIQUE INDEX index1 ON schema1.table1(table1)"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_OBJECT_NAME_USING_KEYWORD,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                            "");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            //DDL_CHECK_INDEX_COLUMN_WITH_BLOB
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_INDEX_COLUMN_WITH_BLOB");
                try {
                    foreach (var text in new String[] {
                        "CREATE TABLE table1(" +
                            "col1 IMAGE," +
                            "INDEX IX_1 (col1)" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 XML," +
                            "INDEX IX_1 (col1)" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 TEXT," +
                            "INDEX IX_1 (col1)" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 VARBINARY(MAX)," +
                            "INDEX IX_1 (col1)" +
                        ")",
                        "ALTER TABLE table1 ADD col1 IMAGE UNIQUE",
                        "ALTER TABLE table1 ADD col1 XML CONSTRAINT constraint1 UNIQUE",
                        "ALTER TABLE table1 ADD col1 TEXT CONSTRAINT constraint1 UNIQUE",
                        "ALTER TABLE table1 ADD col1 VARBINARY(MAX) CONSTRAINT constraint1 UNIQUE",
                        "CREATE INDEX IX_1 ON table1(col1)",
                        "ALTER TABLE table1 ALTER COLUMN col1 IMAGE"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        var statementList = ParseStatementList(text);
                        foreach (var statment in statementList.Statements) {
                            var context = new SqlserverContext(new SqlserverMeta());
                            var validator = DefaultRules.RuleValidators[DefaultRules.DDL_CHECK_INDEX_COLUMN_WITH_BLOB];
                            var checkIndexCountValidator = validator as DisableAddIndexForColumnsTypeBlob;
                            checkIndexCountValidator.IsTest = true;
                            checkIndexCountValidator.ExpectIndexColumns = new Dictionary<string, bool>(){
                                {"col1", true},
                            };
                            checkIndexCountValidator.ExpectColumnAndType = new Dictionary<string, string>(){
                                {"col1", "image"},
                            };

                            validator.Check(context, statment);
                            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR), context.AdviseResultContext.GetLevel());
                            Assert.Equal("[error]禁止将blob类型的列加入索引", context.AdviseResultContext.GetMessage());
                        }
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DDL_CHECK_ALTER_TABLE_NEED_MERGE
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_ALTER_TABLE_NEED_MERGE");
                try {
                    var ruleValidatorContext = new SqlserverContext(new SqlserverMeta());
                    var ruleValidator = DefaultRules.RuleValidators[DefaultRules.DDL_CHECK_ALTER_TABLE_NEED_MERGE];
                    var sqlIndex = 0;
                    foreach (var sql in new String[] {
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 INT",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col2 INT"
                    }) {
                        Console.WriteLine("text:{0}", sql);

                        var statementList = ParseStatementList(sql);
                        foreach (var statement in statementList.Statements) {
                            ruleValidator.Check(ruleValidatorContext, statement);
                            ruleValidatorContext.UpdateContext(LogManager.GetCurrentClassLogger(), statement);
                        }

                        var result = ruleValidatorContext.AdviseResultContext.GetAdviseResult();
                        if (sqlIndex == 0) {
                            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL), result.AdviseLevel);
                            Assert.Equal("", result.AdviseResultMessage);
                        }
                        if (sqlIndex == 1) {
                            Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NOTICE), result.AdviseLevel);
                            Assert.Equal("[notice]已存在对该表的修改语句，建议合并成一个ALTER语句", result.AdviseResultMessage);
                        }
                        sqlIndex++;
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }


            //DDL_DISABLE_DROP_STATEMENT
            {
                Console.WriteLine();
                Console.WriteLine("DDL_DISABLE_DROP_STATEMENT");
                try {
                    foreach (var text in new String[] {

                        "DROP DATABASE db1",
                        "DROP TABLE schema1.table1"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_DISABLE_DROP_STATEMENT,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]禁止除索引外的drop操作");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //ALL_CHECK_WHERE_IS_INVALID
            {
                Console.WriteLine();
                Console.WriteLine("ALL_CHECK_WHERE_IS_INVALID");
                try {
                    foreach (var text in new String[] {

                        "SELECT col1 FROM table1",
                        "SELECT col1 FROM table1 WHERE 1=1",
                        "SELECT col1 FROM table1 WHERE 'a' IN ('a', 'b', 'c')",
                        "SELECT col1 FROM table1 WHERE NULL IS NULL",

                        "UPDATE table1 SET col1=1",
                        "UPDATE table1 SET col1=1 WHERE 1=1",
                        "UPDATE table1 SET col1=1 WHERE 'a' IN ('a', 'b', 'c')",
                        "UPDATE table1 SET col1=1 WHERE NULL IS NULL",

                        "DELETE FROM table1",
                        "DELETE FROM table1 WHERE 1=1",
                        "DELETE FROM table1 WHERE 'a' IN ('a', 'b', 'c')",
                        "DELETE FROM table1 WHERE NULL IS NULL"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.ALL_CHECK_WHERE_IS_INVALID,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }

                try {
                    foreach (var text in new String[] {

                        "SELECT col1 FROM table1 WHERE col1=1",
                        "SELECT col1 FROM table1 WHERE NOT a > 100",
                        "SELECT col1 FROM table1 WHERE (col1=1)",
                        "SELECT col1 FROM table1 WHERE col1 LIKE ('%a%')",
                        "SELECT col1 FROM table1 WHERE col1 = 1 OR col1 = 2",
                        "SELECT col1 FROM table1 WHERE col1 IN ('a', 'b', 'c')",
                        "SELECT col1 FROM table1 WHERE col1 BETWEEN 100 AND 200",
                        "SELECT col1 FROM table1 WHERE col1 IS NULL",
                        "SELECT col1 FROM table1 WHERE 1 < col1"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.ALL_CHECK_WHERE_IS_INVALID,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                            "");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DML_DISABE_SELECT_ALL_COLUMN
            {
                Console.WriteLine();
                Console.WriteLine("DML_DISABE_SELECT_ALL_COLUMN");
                try {
                    foreach (var text in new String[] {

                        "select * from table1"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DML_DISABE_SELECT_ALL_COLUMN,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NOTICE),
                            "[notice]不建议使用select *");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }

                try {
                    foreach (var text in new String[] {

                        "select col1 from table1"

                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DML_DISABE_SELECT_ALL_COLUMN,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                            "");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DDL_CHECK_INDEX_PREFIX
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_INDEX_PREFIX");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE table1(col1 INT, INDEX index1(col1));",
                        "CREATE INDEX index1 ON table1(col1);"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_INDEX_PREFIX,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]普通索引必须要以 \"idx_\" 为前缀");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DDL_CHECK_UNIQUE_INDEX_PREFIX
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_UNIQUE_INDEX_PREFIX");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE table1(col1 INT, CONSTRAINT constraint1 UNIQUE(col1));",
                        "CREATE TABLE table1(col1 INT, INDEX index1 UNIQUE(col1));",
                        "CREATE UNIQUE INDEX index1 ON table1(col1);",
                        "ALTER TABLE table1 ADD CONSTRAINT constraint1 UNIQUE (col1);"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_UNIQUE_INDEX_PREFIX,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]unique索引必须要以 \"uniq_\" 为前缀");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DDL_CHECK_COLUMN_WITHOUT_DEFAULT
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_COLUMN_WITHOUT_DEFAULT");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE table1(col1 INT, col2 INT DEFAULT 0);",
                        "ALTER TABLE table1 ADD col1 VARCHAR(20) DEFAULT 0, col2 INT;",
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COLUMN_WITHOUT_DEFAULT,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]除了自增列及大字段列之外，每个列都必须添加默认值");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE table1(col1 DATE)",
                        "CREATE TABLE table1(col1 DATETIME)",
                        "CREATE TABLE table1(col1 DATETIME2)",
                        "CREATE TABLE table1(col1 DATETIMEOFFSET)",
                        "CREATE TABLE table1(col1 SMALLDATETIME)",
                        "CREATE TABLE table1(col1 TIME)",

                        "ALTER TABLE table1 ADD col1 DATE",
                        "ALTER TABLE table1 ADD col1 DATETIME",
                        "ALTER TABLE table1 ADD col1 DATETIME2",
                        "ALTER TABLE table1 ADD col1 DATETIMEOFFSET",
                        "ALTER TABLE table1 ADD col1 SMALLDATETIME",
                        "ALTER TABLE table1 ADD col1 TIME",
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]timestamp 类型的列必须添加默认值");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE table1(col1 TEXT NOT NULL)",

                        "ALTER TABLE table1 ADD col1 TEXT NOT NULL",
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE table1(col1 TEXT NOT NULL)",

                        "ALTER TABLE table1 ADD col1 TEXT NOT NULL",
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL
            {
                Console.WriteLine();
                Console.WriteLine("DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL");
                try {
                    foreach (var text in new String[]{
                        "CREATE TABLE table1(col1 TEXT DEFAULT '123')",

                        "ALTER TABLE table1 ADD col1 TEXT DEFAULT '123'"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            //DML_CHECK_WITH_LIMIT
            {
                Console.WriteLine();
                Console.WriteLine("DML_CHECK_WITH_LIMIT");
                try {
                    foreach (var text in new String[]{
                        "DELETE TOP(100) FROM table1;",
                        "UPDATE TOP(100) table1 SET col1=1;"
                    }) {
                        Console.WriteLine("text:{0}", text);

                        MyAssert(DefaultRules.DML_CHECK_WITH_LIMIT,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]delete/update 语句不能有top条件");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }
        }
    }
}
