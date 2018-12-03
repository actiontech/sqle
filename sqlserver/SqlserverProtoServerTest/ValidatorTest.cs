using System;
using Xunit;
using SqlserverProtoServer;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Collections.Generic;
using System.IO;
using SqlserverProto;

namespace SqlServerProtoServerTest {
    public class RuleVaidatorTest {
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

        public AdviseResult validate(String ruleName, String text) {
            var validator = DefaultRules.RuleValidators[ruleName];
            var ruleValidatorContext = new RuleValidatorContext(new SqlserverMeta());
            ruleValidatorContext.SqlserverMeta.CurrentDatabase = "master";

            Console.WriteLine("ruleName:{0}, text:{1}", ruleName, text);
            var statementList = ParseStatementList(text);
            foreach (var statement in statementList.Statements) {
                validator.Check(ruleValidatorContext, statement);
                return ruleValidatorContext.AdviseResultContext.GetAdviseResult();
            }
            return null;
        }

        public void MyAssert(String ruleName, String text, String expectLevel, String expectMsg) {
            AdviseResult adviseResult = validate(ruleName, text);
            //Console.WriteLine("{0}, {1}", adviseResult.AdviseLevel, adviseResult.AdviseResultMessage);
            Assert.Equal(expectLevel, adviseResult.AdviseLevel);
            Assert.Equal(expectMsg, adviseResult.AdviseResultMessage);
        }

        [Fact]
        public void RuleValidatorTest() {
            /*SCHEMA_NOT_EXIST*/
            {
                try {
                    foreach (var text in new String[]{
                        "USE database1",
                    }) {
                        MyAssert(DefaultRules.SCHEMA_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]database或者schema database1 不存在");
                    }

                    foreach (var text in new String[]{
                        "CREATE TABLE schema1.table1(col1 INT)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 INT NOT NULL",
                        "SELECT col1 FROM schema1.table1",
                        "INSERT INTO schema1.table1 VALUES (1)",
                        "DELETE FROM schema1.table1",
                        "UPDATE schema1.table1 SET schema1.table1.col1=1"
                    }) {
                        MyAssert(DefaultRules.SCHEMA_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]database或者schema schema1 不存在");
                    }

                    foreach (var text in new String[]{
                        "CREATE TABLE database1.schema1.table1(col1 INT)",
                        "ALTER TABLE database1.schema1.table1 ALTER COLUMN col1 INT NOT NULL",
                        "SELECT col1 FROM database1.schema1.table1",
                        "INSERT INTO database1.schema1.table1 VALUES (1)",
                        "DELETE FROM database1.schema1.table1",
                        "UPDATE database1.schema1.table1 SET schema1.table1.col1=1"
                    }) {
                        MyAssert(DefaultRules.SCHEMA_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]database或者schema database1 不存在\n[error]database或者schema schema1 不存在");
                    }

                    foreach (var text in new String[]{
                        "USE master",
                        "CREATE TABLE master.dbo.table1(col1 INT)",
                        "ALTER TABLE master.dbo.table1 ALTER COLUMN col1 INT NOT NULL",
                        "SELECT col1 FROM master.dbo.table1",
                        "INSERT INTO master.dbo.table1 VALUES(1)",
                        "UPDATE master.dbo.table1 SET dbo.table1.col1=1"
                    }) {
                        MyAssert(DefaultRules.SCHEMA_NOT_EXIST,
                           text,
                           RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                           "");
                    }
                   
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*SCHEMA_EXIST*/
            {/*
                try {
                    foreach (var text in new String[]{
                        "USE database1",
                    }) {
                        MyAssert(DefaultRules.SCHEMA_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                             "");
                    }

                    foreach (var text in new String[]{
                        "CREATE TABLE schema1.table1(col1 INT)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 INT NOT NULL",
                        "SELECT col1 FROM schema1.table1",
                        "INSERT INTO schema1.table1 VALUES (1)",
                        "DELETE FROM schema1.table1",
                        "UPDATE schema1.table1 SET schema1.table1.col1=1"
                    }) {
                        MyAssert(DefaultRules.SCHEMA_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]database或者schema schema1 不存在");
                    }

                    foreach (var text in new String[]{
                        "CREATE TABLE database1.schema1.table1(col1 INT)",
                        "ALTER TABLE database1.schema1.table1 ALTER COLUMN col1 INT NOT NULL",
                        "SELECT col1 FROM database1.schema1.table1",
                        "INSERT INTO database1.schema1.table1 VALUES (1)",
                        "DELETE FROM database1.schema1.table1",
                        "UPDATE database1.schema1.table1 SET schema1.table1.col1=1"
                    }) {
                        MyAssert(DefaultRules.SCHEMA_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]database或者schema database1 不存在\n[error]database或者schema schema1 不存在");
                    }

                    foreach (var text in new String[]{
                        "USE master",
                        "CREATE TABLE master.dbo.table1(col1 INT)",
                        "ALTER TABLE master.dbo.table1 ALTER COLUMN col1 INT NOT NULL",
                        "SELECT col1 FROM master.dbo.table1",
                        "INSERT INTO master.dbo.table1 VALUES(1)",
                        "UPDATE master.dbo.table1 SET dbo.table1.col1=1"
                    }) {
                        MyAssert(DefaultRules.SCHEMA_NOT_EXIST,
                           text,
                           RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                           "");
                    }

                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
                */
            }

            /*TABLE_NOT_EXIST*/
            {
                try {
                    MyAssert(DefaultRules.TABLE_NOT_EXIST,
                             "CREATE TABLE database1.schema1.table1(col1 INT)",
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]表 table1 不存在");
                    /*
                    "CREATE TABLE database1.schema1.table1(col1 INT)",
                    "CREATE TABLE database1.table1(col1 INT)",
                    "CREATE TABLE schema1.table1(col1 INT)",
                    "CREATE TABLE table1(col1 INT)",
                    "ALTER TABLE database1.schema1.table1 ALTER COLUMN col1 INT NOT NULL",
                    "ALTER TABLE database1.table1 ALTER COLUMN col1 INT NOT NULL",
                    "ALTER TABLE schema1.table1 ALTER COLUMN col1 INT NOT NULL",
                    "ALTER TABLE table1 ALTER COLUMN col1 INT NOT NULL",
                    "SELECT col1 FROM table1",
                    "SELECT col1 FROM talbe1 AS table2",
                    "SELECT col1 FROM database1.schema1.table1",
                    "SELECT col1 FROM database1.schema1.table1 AS table2",
                    "SELECT col1 FROM table1 JOIN table2 ON table1.col2=table2.col2",
                    "SELECT col1 FROM database1.schema1.table1 JOIN (SELECT col2 FROM table2 JOIN table3 ON table2.col3=table3.col3) as table4 ON tabl4.col2=table1.col2"
                    "INSERT INTO table1 VALUES (1)",
                    "INSERT INTO schema1.table1 VALUES (1)",
                    "INSERT INTO database1.schema1.table1 VALUES (1)",
                    "INSERT INTO schema1.table1 SELECT 'SELECT', tbl2.col1, tbl3.col2, tbl2.col2 FROM schema2.table2 AS tbl2 INNER JOIN schema3.table3 AS tbl3 ON tbl2.col1=tbl3.col1 WHERE tbl2.col1 LIKE '2%' ORDER BY tbl2.col1, tlb3.col2",
                    "DELETE FROM table1",
                    "DELETE FROM schema1.table1 WHERE col1 IN (SELECT tbl2.col1 FROM schema2.table2 AS tbl2 INNER JOIN table3 AS tbl3 ON tbl2.col2=tbl3.col2)",
                    "DELETE FROM schema1.table1 WHERE col1 IN ('a')"
                    "UPDATE schema1.table1 SET schema1.table1.col1 = 1",
                    "UPDATE schema1.table2 SET schema1.table2.col2 = schema1.table2.col2 + schema1.table1.col2 FROM table2 INNER JOIN table1 ON (table2.col1 = table1.col1)"
                    */

                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_CHECK_OBJECT_NAME_LENGTH*/
            {
                try {
                    foreach (var text in new String[]{
                        /*
                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(20) NULL, col2 INT NULL",
                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(20) NULL CONSTRAINT constraint1 UNIQUE",
                        "ALTER TABLE schema1.table1 WITH NOCHECK ADD CONSTRAINT constraint1 CHECK (col1 > 1)",
                        "EXEC sp_rename @objectname='schema1.table_old', @newname='table_new'",
                        "EXEC sp_rename 'schema1.table1.col_old', 'col_new', 'COLUMN'",
                        "EXEC sp_rename N'schema1.table1.index_old', N'index_new', N'INDEX'",
                        "CREATE UNIQUE INDEX index1 ON schema1.table1(col1)"
                        */
                    }) {
                        AdviseResult adviseResult = validate(DefaultRules.DDL_CHECK_OBJECT_NAME_LENGTH, text);
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_CHECK_PRIMARY_KEY_EXIST*/
            {
                try {
                    foreach (var text in new String[]{
                        /*
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT PRIMARY KEY CLUSTERED" +
                        ")",
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT IDENTITY(1,1) PRIMARY KEY CLUSTERED" +
                        ")",

                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL," +
                            "col2 INT NOT NULL," +
                            "CONSTRAINT PK_constraint PRIMARY KEY(col1)" +
                        ")",
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT IDENTITY(1, 1)," +
                            "col2 INT NOT NULL," +
                            "CONSTRAINT PK_constraint PRIMARY KEY(col1, col2)" +
                        ")"
                        */
                    }) {
                        //ObjectShouldExistRuleValidator(DefaultRules.DDL_CHECK_PRIMARY_KEY_EXIST, text);
                        validate(DefaultRules.DDL_CHECK_PRIMARY_KEY_TYPE, text);
                    }
                } catch (Exception e) {
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DDL_CHECK_TYPE_CHAR_LENGTH*/
            {
                try {
                    foreach (var text in new String[]{
                        /*
                        "CREATE TABLE schema1.table1(" +
                            "col1 CHAR(60) PRIMARY KEY CLUSTERED" +
                        ")",
                        "ALTER TABLE schema1.table1 ADD col1 CHAR(60)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 CHAR(60)"
                        */
                    }) {
                        validate(DefaultRules.DDL_CHECK_TYPE_CHAR_LENGTH, text);
                    }
                } catch (Exception e) {
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DDL_DISABLE_VARCHAR_MAX*/
            {
                try {
                    foreach (var text in new String[]{
                        /*
                        "CREATE TABLE schema1.table1(" +
                            "col1 VARCHAR(MAX) PRIMARY KEY CLUSTERED" +
                        ")",
                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(MAX)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 VARCHAR(MAX)"
                        */
                    }) {
                        validate(DefaultRules.DDL_DISABLE_VARCHAR_MAX, text);
                    }
                } catch (Exception e) {
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DDL_DISABLE_FOREIGN_KEY*/
            {
                try {
                    foreach (var text in new String[] {
                        /*
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL REFERENCES schema1.table2(col11)" +
                        ")",

                        "ALTER TABLE schema1.table1 ADD CONSTRAINT FK_fk1 FOREIGN KEY(col1) REFERENCES schema1.table2(col11)"
                        */
                    }) {
                        validate(DefaultRules.DDL_DISABLE_FOREIGN_KEY, text);
                    }
                } catch (Exception e) {
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DDL_CHECK_INDEX_COUNT*/
            {
                try {
                    foreach (var text in new String[] {
                        /*
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL UNIQUE," +
                        "CONSTRAINT IX_index1 UNIQUE(col1)," +
                        "INDEX IX_index2 (col1)" +
                        ")",

                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(20) NULL CONSTRAINT constraint1 UNIQUE",

                        "CREATE INDEX IX_index1 ON table1(col1)"
                        */
                    }) {
                        validate(DefaultRules.DDL_CHECK_INDEX_COUNT, text);
                    }
                } catch (Exception e) {
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DDL_CHECK_COMPOSITE_INDEX_MAX*/
            {
                try {
                    foreach (var text in new String[] {
                        /*
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL UNIQUE," +
                            "col2 INT NOT NULL," +
                            "INDEX IX_index2 (col1, col2)" +
                        ")",
                        
                        "CREATE INDEX IX_index1 ON table1(col1, col2)"
                        */
                    }) {
                        validate(DefaultRules.DDL_CHECK_COMPOSITE_INDEX_MAX, text);
                    }
                } catch (Exception e) {
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DDL_DISABLE_INDEX_DATA_TYPE_BLOB*/
            {
                try {
                    foreach (var text in new String[]{
                        /*
                        "CREATE TABLE schema1.table1(" +
                            "a INT NOT NULL," +
                            "b BLOB," +
                            "c INT NOT NULL," +
                            "d INT NOT NULL," +
                            "e BLOB," +
                            "CONSTRAINT PK_table1 PRIMARY KEY CLUSTERED(a)," +
                            "CONSTRAINT IX_index1 UNIQUE(b, c)," +
                            "INDEX IX_index2 NONCLUSTERED(d, e))",

                        "ALTER TABLE schema1.table1 ADD col1 BLOB CONSTRAINT IX_index1 UNIQUE",

                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 BLOB",

                        "CREATE INDEX IX_index1 ON table1(col1)"
                        */
                    }) {
                        validate(DefaultRules.DDL_DISABLE_INDEX_DATA_TYPE_BLOB, text);
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DDL_DISABLE_DROP_STATEMENT*/
            {
                try {
                    foreach (var text in new String[] {
                        /*
                        "DROP DATABASE db1",
                        "DROP TABLE schema1.table1"
                        */
                    }) {
                        validate(DefaultRules.DDL_DISABLE_DROP_STATEMENT, text);
                    }
                } catch (Exception e) {
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DDL_CHECK_ALTER_TABLE_NEED_MERGE*/
            {
                try {
                    foreach (var text in new String[] {
                        /*
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 INT; ALTER TABLE schema1.table1 ALTER COLUMN col2 INT;"
                        */
                    }) {
                        validate(DefaultRules.DDL_CHECK_ALTER_TABLE_NEED_MERGE, text);
                    }
                } catch (Exception e) {
                    Assert.IsType<NotImplementedException>(e);
                }
            }

            /*DML_CHECK_INVALID_WHERE_CONDITION*/
            {
                try {
                    foreach (var text in new String[] {
                        /*
                        "SELECT col1 FROM table1",
                        "SELECT col1 FROM table1 WHERE 1=1",
                        "SELECT col1 FROM table1 WHERE col1=1",
                        "SELECT col1 FROM table1 WHERE NOT a > 100",
                        "SELECT col1 FROM table1 WHERE (col1=1)",
                        "SELECT col1 FROM table1 WHERE col1 LIKE ('%a%')",
                        "SELECT col1 FROM table1 WHERE col1 = 1 OR col1 = 2",
                        "SELECT col1 FROM table1 WHERE col1 IN ('a', 'b', 'c')",
                        "SELECT col1 FROM table1 WHERE 'a' IN ('a', 'b', 'c')",
                        "SELECT col1 FROM table1 WHERE col1 BETWEEN 100 AND 200",
                        "SELECT col1 FROM table1 WHERE col1 IS NULL",
                        "SELECT col1 FROM table1 WHERE NULL IS NULL",
                        "SELECT col1 FROM table1 WHERE 1 < col1",
                        "SELECT col1 FROM table1 WHERE EXISTS (SELECT col2 FROM table2)"
                        */
                    }) {
                        validate(DefaultRules.DML_CHECK_INVALID_WHERE_CONDITION, text);
                    }
                } catch (Exception e) {
                    Console.WriteLine(e.Message);
                    Assert.IsType<NotImplementedException>(e);
                }
            }
        }

        [Fact]
        public void SelectAllRuleValidatorTest() {
            var validator = DefaultRules.RuleValidators[DefaultRules.DML_DISABE_SELECT_ALL_COLUMN];
            var ruleValidatorContext = new RuleValidatorContext(new SqlserverMeta());

            // without select *
            var statementList = ParseStatementList("select hello from world");
            foreach (var statement in statementList.Statements) {
                validator.Check(ruleValidatorContext, statement);
                var adviseResult = ruleValidatorContext.AdviseResultContext.GetAdviseResult();
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL), adviseResult.AdviseLevel);
                Assert.Equal("", adviseResult.AdviseResultMessage);
            }

            // with select *
            ruleValidatorContext.AdviseResultContext.ResetAdviseResult();
            statementList = ParseStatementList("select * from world");
            foreach (var statement in statementList.Statements) {
                validator.Check(ruleValidatorContext, statement);
                var adviseResult = ruleValidatorContext.AdviseResultContext.GetAdviseResult();
                Assert.Equal(RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NOTICE), adviseResult.AdviseLevel);
                Assert.Equal("[notice]不建议使用select *", adviseResult.AdviseResultMessage);
            }
        }
    }
}
