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
            var context = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            });
            foreach (var statement in statementList.Statements) {
                var isDDL = false;
                var isDML = false;
                var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                Assert.True(isDDL == true);
                Assert.True(isDML == false);
                Assert.Equal("DROP TABLE master..table1;", rollbackSql);
            }
        }

        private void rollbackAlterTable() {
            var context = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            });
            var logger = LogManager.GetCurrentClassLogger();
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
                    Console.WriteLine("{0}", text);
                    new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                    context.UpdateContext(logger, statement);
                    Console.WriteLine("=====================================================");
                }
            }
        }

        private void rollbackCreateIndex() {
            var text = "CREATE UNIQUE INDEX IX2 ON table1 (col1, col2) WITH (DROP_EXISTING=ON)";
            var statementList = ParseStatementList(text);
            var context = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            });
            foreach (var statement in statementList.Statements) {
                var isDDL = false;
                var isDML = false;
                var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                Assert.True(isDDL == true);
                Assert.True(isDML == false);
                Assert.Equal("DROP INDEX IX2 ON master..table1;", rollbackSql);
            }
        }

        private void rollbackDropIndex() {
            var text = "DROP INDEX IX1 ON tbl6";
            var statementList = ParseStatementList(text);
            var context = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            });
            foreach (var statement in statementList.Statements) {
                var isDDL = false;
                var isDML = false;
                var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                Assert.True(isDDL == true);
                Assert.True(isDML == false);
                Assert.Equal("CREATE  NONCLUSTERED INDEX IX1 ON master..table1 (c);", rollbackSql);
            }
        }

        private void rollbackDropTable() {
            //var text = "DROP TABLE dbo.WorkOut";
            var text = "DROP TABLE dbo.tbl7";
            var statementList = ParseStatementList(text);
            var context = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            });
            foreach (var statement in statementList.Statements) {
                var isDDL = false;
                var isDML = false;
                var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                Assert.True(isDDL == true);
                Assert.True(isDML == false);
                Console.WriteLine("{0}", rollbackSql);
                //Assert.Equal("CREATE  NONCLUSTERED INDEX IX1 ON master..table1 (c);", rollbackSql);
            }
        }

        [Fact]
        public void RollbackDDLSqlTest() {
            //rollbackCreateTable();
            //rollbackCreateIndex();
            //rollbackDropIndex();
            //rollbackDropTable();
            //rollbackAlterTable();
        }

        private void rollbackInsert() {
            var context = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            }, new Config() {
                DMLRollbackMaxRows = 10000
            });
            foreach (var text in new String[]{
                "INSERT INTO tbl1(col1, col2) VALUES (1, 2), (3, 4);",
                "INSERT INTO tbl1 VALUES (5, 6, 7);"
            }) {
                var statementList = ParseStatementList(text);
                foreach (var statement in statementList.Statements) {
                    var isDDL = false;
                    var isDML = false;
                    var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                    Assert.True(isDDL == false);
                    Assert.True(isDML == true);
                }
            }
        }

        private void rollbackDelete() {
            var context = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            }, new Config() {
                DMLRollbackMaxRows = 10000
            });
            foreach (var text in new String[]{
                "DELETE FROM tbl1 WHERE col1 = 2 AND col2 = 3;",
                "DELETE FROM tbl1;",
                "DELETE TOP(1) FROM tbl1;"
            }) {
                var statementList = ParseStatementList(text);
                foreach (var statement in statementList.Statements) {
                    var isDDL = false;
                    var isDML = false;
                    var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                    Assert.True(isDDL == false);
                    Assert.True(isDML == true);
                    Console.WriteLine("{0}", rollbackSql);
                }
            }
        }

        private void rollbackUpdate() {
            var context = new SqlserverContext(new SqlserverMeta() {
                Host = "10.186.62.15",
                Port = "1433",
                User = "sa",
                Password = "123456aB"
            }, new Config() {
                DMLRollbackMaxRows = 10000
            });
            foreach (var text in new String[]{
                "UPDATE tbl3 SET col1=\"dddd\", col2=2 WHERE col1='aa';",
            }) {
                var statementList = ParseStatementList(text);
                foreach (var statement in statementList.Statements) {
                    var isDDL = false;
                    var isDML = false;
                    var rollbackSql = new RollbackSql().GetRollbackSql(context, statement, out isDDL, out isDML);
                    Assert.True(isDDL == false);
                    Assert.True(isDML == true);
                    Console.WriteLine("{0}", rollbackSql);
                }
            }
        }

        [Fact]
        public void RollbackDMLSqlTest() {
            //rollbackInsert();
            //rollbackDelete();
            rollbackUpdate();
        }
    }
}