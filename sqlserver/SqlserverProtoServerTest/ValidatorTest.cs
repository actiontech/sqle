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
            var ruleValidatorContext = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            });
            ruleValidatorContext.SqlserverMeta.CurrentDatabase = "master";

            Console.WriteLine("ruleName:{0}, text:{1}", ruleName, text);
            var statementList = ParseStatementList(text);
            foreach (var statement in statementList.Statements) {
                validator.Check(ruleValidatorContext, statement);
                return ruleValidatorContext.AdviseResultContext.GetAdviseResult();
            }
            return null;
        }

        private void MyAssert(String ruleName, String text, String expectLevel, String expectMsg) {
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
            {
                try {
                    foreach (var text in new String[]{
                        "CREATE DATABASE master"
                    }) {
                        MyAssert(DefaultRules.SCHEMA_EXIST,
                           text,
                           RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                           "[error]database master 已存在");
                    }

                    foreach (var text in new String[]{
                        "CREATE DATABASE database1"
                    }) {
                        MyAssert(DefaultRules.SCHEMA_EXIST,
                                text,
                                 RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                                 "");
                    }

                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*TABLE_NOT_EXIST*/
            {
                try {
                    foreach (var text in new String[]{

                        "ALTER TABLE database1.schema1.table1 ALTER COLUMN col1 INT NOT NULL",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 INT NOT NULL",
                        "SELECT col1 FROM database1.schema1.table1",
                        "SELECT col1 FROM database1.schema1.table1 AS table2",
                        "INSERT INTO schema1.table1 VALUES (1)",
                        "INSERT INTO database1.schema1.table1 VALUES (1)",
                        "DELETE FROM schema1.table1 WHERE col1 IN ('a')",
                        "UPDATE schema1.table1 SET schema1.table1.col1 = 1"

                    }) {
                        MyAssert(DefaultRules.TABLE_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]表 schema1.table1 不存在");
                    }

                    foreach (var text in new String[]{

                        "ALTER TABLE table1 ALTER COLUMN col1 INT NOT NULL",
                        "SELECT col1 FROM table1",
                        "SELECT col1 FROM table1 AS table2",
                        "INSERT INTO table1 VALUES (1)",
                        "DELETE FROM table1"

                    }) {
                        MyAssert(DefaultRules.TABLE_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]表 table1 不存在");
                    }

                    foreach (var text in new String[]{
                        "SELECT col1 FROM table1 JOIN table2 ON table1.col2=table2.col2",
                        "UPDATE table1 SET table1.col2 = table1.col2 + table2.col2 FROM table2 INNER JOIN table1 ON (table2.col1 = table1.col1)"
                    }) {
                        MyAssert(DefaultRules.TABLE_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]表 table1,table2 不存在");
                    }

                    foreach (var text in new String[]{
                        "SELECT col1 FROM database1.schema1.table1 JOIN (SELECT col2 FROM schema2.table2 JOIN schema3.table3 ON table2.col3=table3.col3) as table4 ON tabl4.col2=table1.col2",
                        "INSERT INTO schema1.table1 SELECT 'SELECT', tbl2.col1, tbl3.col2, tbl2.col2 FROM schema2.table2 AS tbl2 INNER JOIN schema3.table3 AS tbl3 ON tbl2.col1=tbl3.col1 WHERE tbl2.col1 LIKE '2%' ORDER BY tbl2.col1, tlb3.col2",
                        "DELETE FROM schema1.table1 WHERE col1 IN (SELECT tbl2.col1 FROM schema2.table2 AS tbl2 INNER JOIN schema3.table3 AS tbl3 ON tbl2.col2=tbl3.col2)",
                    }) {
                        MyAssert(DefaultRules.TABLE_NOT_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]表 schema1.table1,schema2.table2,schema3.table3 不存在");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /* FakerRuleValidator
             * DDL_CREATE_TABLE_NOT_EXIST
             * DDL_TABLE_USING_INNODB_UTF8MB4
            */
            {
            }

            /*DDL_CHECK_OBJECT_NAME_LENGTH*/
            {
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
                        MyAssert(DefaultRules.DDL_CHECK_OBJECT_NAME_LENGTH,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]表名、列名、索引名的长度不能大于64字节");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_CHECK_PRIMARY_KEY_EXIST*/
            {
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
                        MyAssert(DefaultRules.DDL_CHECK_PRIMARY_KEY_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                             "");
                    }

                    foreach (var text in new String[]{
                        "CREATE TABLE schema1.table1(" +
                            "col1 INT" +
                        ")"
                    }) {
                        MyAssert(DefaultRules.DDL_CHECK_PRIMARY_KEY_EXIST,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]表必须有主键");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_CHECK_PRIMARY_KEY_TYPE*/
            {
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
                        MyAssert(DefaultRules.DDL_CHECK_PRIMARY_KEY_TYPE,
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
                        MyAssert(DefaultRules.DDL_CHECK_PRIMARY_KEY_TYPE,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]主键建议使用自增");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_DISABLE_VARCHAR_MAX*/
            {
                try {
                    foreach (var text in new String[]{

                        "CREATE TABLE schema1.table1(" +
                            "col1 VARCHAR(MAX) PRIMARY KEY CLUSTERED" +
                        ")",
                        "ALTER TABLE schema1.table1 ADD col1 VARCHAR(MAX)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 VARCHAR(MAX)"

                    }) {
                        MyAssert(DefaultRules.DDL_DISABLE_VARCHAR_MAX,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]禁止使用 varchar(max)");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_CHECK_TYPE_CHAR_LENGTH*/
            {
                try {
                    foreach (var text in new String[]{

                        "CREATE TABLE schema1.table1(" +
                            "col1 CHAR(60) PRIMARY KEY CLUSTERED" +
                        ")",
                        "ALTER TABLE schema1.table1 ADD col1 CHAR(60)",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 CHAR(60)"

                    }) {
                        MyAssert(DefaultRules.DDL_CHECK_TYPE_CHAR_LENGTH,
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
                        MyAssert(DefaultRules.DDL_CHECK_TYPE_CHAR_LENGTH,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                             "");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_DISABLE_FOREIGN_KEY*/
            {
                try {
                    foreach (var text in new String[] {

                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL REFERENCES schema1.table2(col11)" +
                        ")",

                        "ALTER TABLE schema1.table1 ADD CONSTRAINT FK_fk1 FOREIGN KEY(col1) REFERENCES schema1.table2(col11)"

                    }) {
                        MyAssert(DefaultRules.DDL_DISABLE_FOREIGN_KEY,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                             "[error]禁止使用外键");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_CHECK_INDEX_COUNT*/
            {
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
                        MyAssert(DefaultRules.DDL_CHECK_INDEX_COUNT,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                             "");
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

                        // There must a table1 table in sql server when executing those samples
                        /*
                        "ALTER TABLE table1 ADD col10 VARCHAR(20) NULL CONSTRAINT constraint1 UNIQUE",
                        "CREATE INDEX IX_index1 ON table1(col1)"
                        */

                    }) {
                        MyAssert(DefaultRules.DDL_CHECK_INDEX_COUNT,
                             text,
                             RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NOTICE),
                             "[notice]索引个数建议不超过5个");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_CHECK_COMPOSITE_INDEX_MAX*/
            {
                try {
                    foreach (var text in new String[] {

                        "CREATE TABLE schema1.table1(" +
                            "col1 INT NOT NULL UNIQUE," +
                            "col2 INT NOT NULL," +
                            "INDEX IX_index2 (col1, col2)" +
                        ")",
                        "CREATE INDEX IX_index1 ON table1(col1, col2)"

                    }) {
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
                        MyAssert(DefaultRules.DDL_CHECK_COMPOSITE_INDEX_MAX,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NOTICE),
                            "[notice]复合索引的列数量不建议超过5个");
                    }
                } catch (Exception e) {
                    Console.WriteLine("{0}", e.Message);
                }
            }

            /*DDL_DISABLE_INDEX_DATA_TYPE_BLOB*/
            {
                try {
                    foreach (var text in new String[] {
                        "CREATE TABLE table1(" +
                            "col1 IMAGE" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 XML" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 TEXT" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 VARBINARY(MAX)" +
                        ")",

                        "ALTER TABLE table1 ADD col1 IMAGE",
                        "ALTER TABLE table1 ADD col1 XML",
                        "ALTER TABLE table1 ADD col1 TEXT",
                        "ALTER TABLE table1 ADD col1 VARBINARY(MAX)",

                        // There must a table1 table in sql server when executing those two samples
                        /*
                        "ALTER TABLE table2 ALTER COLUMN col2 IMAGE",
                        "ALTER TABLE table2 ALTER COLUMN col2 XML",
                        "ALTER TABLE table2 ALTER COLUMN col2 TEXT",
                        "ALTER TABLE table2 ALTER COLUMN col2 VARBINARY(MAX)",

                        "CREATE INDEX IX_index1 ON table1(col1)"
                        */
                    }) {
                        MyAssert(DefaultRules.DDL_DISABLE_INDEX_DATA_TYPE_BLOB,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                            "");
                    }

                    foreach (var text in new String[]{
                        "CREATE TABLE table1(" +
                            "col1 IMAGE PRIMARY KEY" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 XML NOT NULL UNIQUE" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 TEXT," +
                            "CONSTRAINT IX_index1 UNIQUE(col1)" +
                        ")",
                        "CREATE TABLE table1(" +
                            "col1 VARBINARY(MAX)," +
                            "INDEX IX_index1 (col1)" +
                        ")",

                        "ALTER TABLE table1 ADD col1 IMAGE CONSTRAINT constraint1 UNIQUE",
                        "ALTER TABLE table1 ADD col1 XML CONSTRAINT constraint1 UNIQUE",
                        "ALTER TABLE table1 ADD col1 TEXT CONSTRAINT constraint1 UNIQUE",
                        "ALTER TABLE table1 ADD col1 VARBINARY(MAX) CONSTRAINT constraint1 UNIQUE",
                        /*
                        // There must a table1 table in sql server when executing those two samples
                        "ALTER TABLE table1 ALTER COLUMN col1 IMAGE",
                        "ALTER TABLE table1 ALTER COLUMN col1 XML",
                        "ALTER TABLE table1 ALTER COLUMN col1 TEXT",
                        "ALTER TABLE table1 ALTER COLUMN col1 VARBINARY(MAX)",

                        "CREATE INDEX IX_index1 ON table4(col1)"
                        */
                    }) {
                        MyAssert(DefaultRules.DDL_DISABLE_INDEX_DATA_TYPE_BLOB,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]禁止将blob类型的列加入索引");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            /*DDL_CHECK_ALTER_TABLE_NEED_MERGE*/
            {
                try {
                    var ruleValidatorContext = new SqlserverContext(new SqlserverMeta() {
                        Host = "10.186.62.15",
                        Port = "1433",
                        User = "sa",
                        Password = "123456aB"
                    });
                    var ruleValidator = DefaultRules.RuleValidators[DefaultRules.DDL_CHECK_ALTER_TABLE_NEED_MERGE];
                    var sqlIndex = 0;
                    foreach (var sql in new String[] {

                        "ALTER TABLE schema1.table1 ALTER COLUMN col1 INT",
                        "ALTER TABLE schema1.table1 ALTER COLUMN col2 INT"

                    }) {
                        Console.WriteLine("ruleName:{0}, text:{1}", DefaultRules.DDL_CHECK_ALTER_TABLE_NEED_MERGE, sql);
                        var statementList = ParseStatementList(sql);
                        foreach (var statement in statementList.Statements) {
                            ruleValidator.Check(ruleValidatorContext, statement);
                            ruleValidatorContext.UpdateContext(statement);
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


            /*DDL_DISABLE_DROP_STATEMENT*/
            {
                try {
                    foreach (var text in new String[] {

                        "DROP DATABASE db1",
                        "DROP TABLE schema1.table1"

                    }) {
                        MyAssert(DefaultRules.DDL_DISABLE_DROP_STATEMENT,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.ERROR),
                            "[error]禁止除索引外的drop操作");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            /*DML_CHECK_INVALID_WHERE_CONDITION*/
            {
                try {
                    foreach (var text in new String[] {

                        "SELECT col1 FROM table1",
                        "SELECT col1 FROM table1 WHERE 1=1",
                        "SELECT col1 FROM table1 WHERE 'a' IN ('a', 'b', 'c')",
                        "SELECT col1 FROM table1 WHERE NULL IS NULL"

                    }) {
                        MyAssert(DefaultRules.DML_CHECK_INVALID_WHERE_CONDITION,
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
                        MyAssert(DefaultRules.DML_CHECK_INVALID_WHERE_CONDITION,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                            "");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }

            /*DML_DISABE_SELECT_ALL_COLUMN*/
            {
                try {
                    foreach (var text in new String[] {

                        "select * from table1"

                    }) {
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
                        MyAssert(DefaultRules.DML_DISABE_SELECT_ALL_COLUMN,
                            text,
                            RULE_LEVEL_STRING.GetRuleLevelString(RULE_LEVEL.NORMAL),
                            "");
                    }
                } catch (Exception e) {
                    Console.WriteLine("exception:{0}", e.Message);
                }
            }
        }
    }
}
