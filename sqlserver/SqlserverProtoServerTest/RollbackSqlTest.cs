using System;
using Xunit;
using SqlserverProtoServer;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Collections.Generic;
using System.IO;
using SqlserverProto;

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

        public void RollbackCreateDatabase() {
            var text = "CREATE DATABASE db1";
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
                Assert.Equal("DROP DATABASE db1", rollbackSql);
            }
        }

        public void RollbackCreateTable() {
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
                Assert.Equal("DROP TABLE master..table1", rollbackSql);
            }
        }

        public void RollbackCreateIndex() {
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
                Assert.Equal("CREATE UNIQUE NONCLUSTERED INDEX IX2 ON master..table1 (a,b) WITH (DROP_EXISTING=ON)", rollbackSql);
            }
        }

        public void RollbackDropIndex() {
            var text = "DROP INDEX IX1 ON table1";
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

        [Fact]
        public void RollbackDDLSqlTest() {
            // create database db1;
            RollbackCreateDatabase();
            RollbackCreateTable();
            RollbackCreateIndex();
            RollbackDropIndex();
        }

        [Fact]
        public void RollbackDMLSqlTest() {

        }
    }
}