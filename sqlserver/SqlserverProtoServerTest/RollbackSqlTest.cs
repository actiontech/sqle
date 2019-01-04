using System;
using Xunit;
using SqlserverProtoServer;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Collections.Generic;
using System.IO;
using SqlserverProto;
using NLog;

namespace SqlServerProtoServerTest {
    public class RollbackSqlTest {
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

        private void rollbackCreateTable() {
            var text = "CREATE TABLE table1(a INT)";
            var statementList = ParseStatementList(text);
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseName = "database1";
            context.ExpectSchemaName = "schema1";
            context.ExpectTableName = "table1";
            foreach (var statement in statementList.Statements) {
                var isDDL = false;
                var isDML = false;
                var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                Assert.True(isDDL == true);
                Assert.True(isDML == false);
                Assert.Equal("DROP TABLE database1.schema1.table1;", rollbackSql);
            }
        }

        private void rollbackDropTable() {
            var context = new SqlserverContext(new SqlserverMeta()) {
                IsTest = true,
                ExpectDatabaseName = "database1",
                ExpectSchemaName = "schema1",
                ExpectTableName = "table1",
            };

            {
                StatementList initStatementList = ParseStatementList("CREATE TABLE database1.schema1.table1(" +
                                                                     "col1 INT NOT NULL, " +
                                                                     "col2 INT NOT NULL, " +
                                                                     "CONSTRAINT PK_1 PRIMARY KEY (col1), " +
                                                                     "CONSTRAINT UN_1 UNIQUE (col2), " +
                                                                     "INDEX IX_1 (col2))");
                CreateTableStatement initStatement = initStatementList.Statements[0] as CreateTableStatement;
                context.UpdateContext(LogManager.GetCurrentClassLogger(), initStatement);
            }
            var text = "DROP TABLE database1.schema1.table1";
            var statementList = ParseStatementList(text);
            foreach (var statement in statementList.Statements) {
                var isDDL = false;
                var isDML = false;
                var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                Assert.True(isDDL == true);
                Assert.True(isDML == false);
                Assert.Equal("CREATE TABLE database1.schema1.table1 (col1 INT NOT NULL,col2 INT NOT NULL,CONSTRAINT PK_1 PRIMARY KEY (col1),CONSTRAINT UN_1 UNIQUE (col2),INDEX IX_1 (col2));", rollbackSql);
            }
        }

        private void rollbackAlterTable() {
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseName = "database1";
            context.ExpectSchemaName = "dbo";
            context.ExpectTableName = "test";
            var logger = LogManager.GetCurrentClassLogger();
            var index = 0;
            foreach (var text in new String[]{
                "CREATE TABLE dbo.test(column_b INT, column_c INT, column_d INT, CONSTRAINT my_constraint UNIQUE (column_c), CONSTRAINT my_pk_constraint UNIQUE (column_d));",

                "ALTER TABLE dbo.test ADD AddDate smalldatetime NULL CONSTRAINT AddDateDflt DEFAULT GETDATE() WITH VALUES;",
                "ALTER TABLE dbo.test ADD column_b INT IDENTITY CONSTRAINT column_b_pk PRIMARY KEY, " +
                    "column_c INT NULL CONSTRAINT column_c_fk REFERENCES test1(column_a), " +
                    "column_d VARCHAR(16) NULL CONSTRAINT column_d_chk CHECK (column_d LIKE '[0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]' OR column_d LIKE '([0-9][0-9][0-9]) [0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]'), " +
                    "column_e DECIMAL(3,3) CONSTRAINT column_e_default DEFAULT .081;",
                "ALTER TABLE dbo.test DROP COLUMN column_c, column_d;",

                "ALTER TABLE dbo.test WITH NOCHECK ADD CONSTRAINT exd_check CHECK (column_a > 1);",
                "ALTER TABLE dbo.test DROP CONSTRAINT my_constraint, my_pk_constraint, COLUMN column_b;",

                "ALTER TABLE dbo.test ALTER COLUMN column_b DECIMAL(5,2);",
                "ALTER TABLE test ALTER COLUMN column_d varchar(50) ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = [CEK1], ENCRYPTION_TYPE=Randomized, ALGORITHM='AEAD_AES_256_CBC_HMAC_SHA_256') NULL;",

                // rename table schem.table
                "EXEC sp_rename 'dbo.test', 'test1';",
                // rename column table.column|schema.table.column
                "EXEC sp_rename 'dbo.test.column_b', 'column_b1', 'COLUMN';"
            }) {
                var statementList = ParseStatementList(text);
                foreach (var statement in statementList.Statements) {
                    var isDDL = false;
                    var isDML = false;
                    Console.WriteLine("index:{0}, text:{1}", index, text);
                    var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                    if (index == 1) {
                        Console.WriteLine("rollback 1:{0}", rollbackSql);
                        Assert.Equal("ALTER TABLE database1.dbo.test DROP COLUMN AddDate", rollbackSql);
                    }
                    if (index == 2) {
                        Console.WriteLine("rollback 2:{0}", rollbackSql);
                        Assert.Equal("ALTER TABLE database1.dbo.test DROP COLUMN column_b,column_c,column_d,column_e", rollbackSql);
                    }
                    if (index == 3) {
                        Console.WriteLine("rollback 3:{0}", rollbackSql);
                        Assert.Equal("ALTER TABLE database1.dbo.test ADD column_c INT,column_d INT", rollbackSql);
                    }
                    if (index == 4) {
                        Console.WriteLine("rollback 4:{0}", rollbackSql);
                        Assert.Equal("ALTER TABLE database1.dbo.test DROP CONSTRAINT exd_check", rollbackSql);
                    }
                    if (index == 5) {
                        Console.WriteLine("rollback 5:{0}", rollbackSql);
                        Assert.Equal("ALTER TABLE database1.dbo.test ADD column_b INT;ALTER TABLE database1.dbo.test ADD CONSTRAINT my_constraint UNIQUE (column_c),CONSTRAINT my_pk_constraint UNIQUE (column_d)", rollbackSql);
                    }
                    if (index == 6) {
                        Console.WriteLine("rollback 6:{0}", rollbackSql);
                        Assert.Equal("ALTER TABLE database1.dbo.test ALTER COLUMN column_b INT", rollbackSql);
                    }
                    if (index == 7) {
                        Console.WriteLine("rollback 7:{0}", rollbackSql);
                        Assert.Equal("ALTER TABLE database1.dbo.test ALTER COLUMN column_d INT", rollbackSql);
                    }
                    if (index == 8) {
                        Console.WriteLine("rollback 8:{0}", rollbackSql);
                        Assert.Equal("EXEC sp_rename 'dbo.test1', 'test'", rollbackSql);
                    }
                    if (index == 9) {
                        Console.WriteLine("rollback 9:{0}", rollbackSql);
                        Assert.Equal("EXEC sp_rename 'dbo.test.column_b1', 'column_b', 'COLUMN'", rollbackSql);
                    }
                    context.UpdateContext(logger, statement);
                    index += 1;
                    Console.WriteLine("=====================================================");
                }
            }
        }

        private void rollbackCreateIndex() {
            var text = "CREATE UNIQUE INDEX IX2 ON table1 (col1, col2)";
            var statementList = ParseStatementList(text);
            var context = new SqlserverContext(new SqlserverMeta());
            context.IsTest = true;
            context.ExpectDatabaseName = "database1";
            context.ExpectSchemaName = "schema1";
            context.ExpectTableName = "table1";
            foreach (var statement in statementList.Statements) {
                var isDDL = false;
                var isDML = false;
                var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                Assert.True(isDDL == true);
                Assert.True(isDML == false);
                Assert.Equal("DROP INDEX IX2 ON database1.schema1.table1;", rollbackSql);
            }
        }

        [Fact]
        public void RollbackDDLSqlTest() {
            rollbackCreateTable();
            rollbackCreateIndex();
            rollbackDropTable();
            rollbackAlterTable();
        }

        private void rollbackInsert() {
            var context = new SqlserverContext(new SqlserverMeta(), new Config() {
                DMLRollbackMaxRows = 10000
            }) {
                IsTest = true,
                ExpectDatabaseName = "database1",
                ExpectSchemaName = "schema1",
                ExpectTableName = "tbl1"
            };
            var primaryKeys = new Dictionary<String, List<String>>();
            primaryKeys["database1.schema1.tbl1"] = new List<String>() {
                "col1"
            };
            context.PrimaryKeys = primaryKeys;
            context.ExpectColumns = new List<String>() {
                "col1", "col2"
            };

            var index = 0;
            foreach (var text in new String[]{
                "INSERT INTO tbl1(col1, col2) VALUES (1, 2), (3, 4);",
                "INSERT INTO tbl1 VALUES (5, 6, 7);"
            }) {
                Console.WriteLine("index:{0}, text:{1}", index, text);
                var statementList = ParseStatementList(text);
                foreach (var statement in statementList.Statements) {
                    var isDDL = false;
                    var isDML = false;
                    var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                    Assert.True(isDDL == false);
                    Assert.True(isDML == true);
                    if (index == 0) {
                        Console.WriteLine("rollback 0: {0}", rollbackSql);
                        Assert.Equal("DELETE FROM database1.schema1.tbl1 WHERE col1 = '1';\nDELETE FROM database1.schema1.tbl1 WHERE col1 = '3';", rollbackSql);
                        context.ExpectColumns = new List<String>() {
                            "col1", "col2", "col3"
                        };
                    }
                    if (index == 1) {
                        Console.WriteLine("rollback 1: {0}", rollbackSql);
                        Assert.Equal("DELETE FROM database1.schema1.tbl1 WHERE col1 = '5';", rollbackSql);
                    }
                    Console.WriteLine("=====================================================");
                    index +=1;
                }
            }
        }

        private void rollbackDelete() {
            var context = new SqlserverContext(new SqlserverMeta(), new Config() {
                DMLRollbackMaxRows = 10000
            }) {
                IsTest = true,
                ExpectDatabaseName = "database1",
                ExpectSchemaName = "schema1",
                ExpectTableName = "tbl1",
            };
            var columns = new List<String>() {
                "col1", "col2"
            };
            context.ExpectColumns = columns;
            var records = new List<Dictionary<String, String>>();
            var record = new Dictionary<String, String>();
            record["col1"] = "2";
            record["col2"] = "3";
            records.Add(record);
            context.ExpectRecords = records;

            var index = 0;
            foreach (var text in new String[]{
                "DELETE FROM tbl1 WHERE col1 = 2 AND col2 = 3;",
                "DELETE FROM tbl1;",
                "DELETE TOP(1) FROM tbl1;"
            }) {
                Console.WriteLine("index:{0}, text:{1}", index, text);
                var statementList = ParseStatementList(text);
                foreach (var statement in statementList.Statements) {
                    var isDDL = false;
                    var isDML = false;
                    var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                    Assert.True(isDDL == false);
                    Assert.True(isDML == true);
                    if (index == 0) {
                        Console.WriteLine("rollback 0: {0}", rollbackSql);
                        Assert.Equal("INSERT INTO database1.schema1.tbl1 (col1, col2) VALUES ('2', '3')", rollbackSql);
                    }
                    if (index == 1) {
                        Console.WriteLine("rollback 1: {0}", rollbackSql);
                        Assert.Equal("INSERT INTO database1.schema1.tbl1 (col1, col2) VALUES ('2', '3')", rollbackSql);
                    }
                    if (index == 2) {
                        Console.WriteLine("rollback 2: {0}", rollbackSql);
                        Assert.Equal("INSERT INTO database1.schema1.tbl1 (col1, col2) VALUES ('2', '3')", rollbackSql);
                    }
                    Console.WriteLine("=====================================================");
                    index += 1;
                }
            }
        }

        private void rollbackUpdate() {
            var context = new SqlserverContext(new SqlserverMeta(), new Config() {
                DMLRollbackMaxRows = 10000
            }) {
                IsTest = true,
                ExpectDatabaseName = "database1",
                ExpectSchemaName = "schema1",
                ExpectTableName = "tbl3"
            };
            var primaryKeys = new Dictionary<String, List<String>>();
            primaryKeys["database1.schema1.tbl3"] = new List<String>() {
                "col1"
            };
            context.PrimaryKeys = primaryKeys;
            var columns = new List<String>() {
                "col1", "col2"
            };
            context.ExpectColumns = columns;
            var records = new List<Dictionary<String, String>>();
            var record = new Dictionary<String, String>();
            record["col1"] = "aa";
            record["col2"] = "3";
            records.Add(record);
            context.ExpectRecords = records;

            foreach (var text in new String[]{
                "UPDATE tbl3 SET col1=\"dddd\", col2=2 WHERE col1='aa';",
            }) {
                Console.WriteLine("text:{0}", text);
                var statementList = ParseStatementList(text);
                foreach (var statement in statementList.Statements) {
                    var isDDL = false;
                    var isDML = false;
                    var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                    Assert.True(isDDL == false);
                    Assert.True(isDML == true);
                    Console.WriteLine("rollback: {0}", rollbackSql);
                    Assert.Equal("UPDATE database1.schema1.tbl3 SET col1 = 'aa', col2 = '3' WHERE col1 = \"dddd\";", rollbackSql);
                }
            }
        }

        [Fact]
        public void RollbackDMLSqlTest() {
            rollbackInsert();
            rollbackDelete();
            rollbackUpdate();
        }
    }
}