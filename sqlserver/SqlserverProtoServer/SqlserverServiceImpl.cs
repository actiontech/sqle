using System;
using System.Threading.Tasks;
using System.Collections.Generic;
using System.IO;
using Grpc.Core;
using SqlserverProto;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace SqlserverProtoServer {
    public class SqlServerServiceImpl : SqlserverService.SqlserverServiceBase {
        /*
         * The following compatibility level values can be configured (not all versions supports all of the above listed compatibility level):
         * Product                          Database Engine Version         Compatibility Level Designation         Supported Compatibility Level Values
         * SQL Server 2019 preview                     15                              150                                     150,140,130,120,110,100
         * SQL Server 2017 (14.x)                      14                              140                                     140,130,120,110,100
         * Azure SQL Database logical server           12                              130                                     150,140,130,120,110,100
         * Azure SQL Database Managed Instance         12                              130                                     150,140,130,120,110,100
         * SQL Server 2016 (13.x)                      13                              130                                     130,120,110,100
         * SQL Server 2014 (12.x)                      12                              120                                     120,110,100
         * SQL Server 2012 (11.x)                      11                              110                                     110,100,90
         * SQL Server 2008 R2                          10.5                            100                                     100,90,80
         * SQL Server 2008                             10                              100                                     100,90,80
         * SQL Server 2005 (9.x)                       9                               90                                      90,80
         * SQL Server 2000                             8                               80                                      80
         * 
         * more information: https://docs.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql-compatibility-level?view=sql-server-2017
         */
        private const String SQL80 = "80";
        private const String SQL90 = "90";
        private const String SQL100 = "100";
        private const String SQL110 = "110";
        private const String SQL120 = "120";
        private const String SQL130 = "130";

        // sql server parser
        private readonly Dictionary<String, TSqlParser> SqlParsers;

        private TSqlParser GetParser(String version) {
            // set default sql parser version to SQL100
            if (version == "") {
                version = SQL130;
            }

            var parser = SqlParsers[version];
            if (parser == null) {
                throw new ArgumentException(String.Format("unsupported TSqlParser version:{0}", version));
            }

            return parser;
        }

        // construct function
        public SqlServerServiceImpl() {
            SqlParsers = new Dictionary<String, TSqlParser> {
                {SQL80, new TSql80Parser(false)},
                {SQL90, new TSql90Parser(false)},
                {SQL100, new TSql100Parser(false)},
                {SQL110, new TSql110Parser(false)},
                {SQL120, new TSql120Parser(false)},
                {SQL130, new TSql130Parser(false)},
            };
        }

        // parse sqls
        private StatementList ParseStatementList(String version, String text) {
            // get parser
            var parser = GetParser(version);

            // parse sqls
            var reader = new StringReader(text);
            IList<ParseError> errorList;
            var statementList = parser.ParseStatementList(reader, out errorList);
            if (errorList.Count > 0) {
                throw new ArgumentException(String.Format("parse sql {0} error: {1}", text, errorList.ToString()));
            }

            return statementList;
        }

        // Splite sqls
        public override Task<SplitSqlsOutput> GetSplitSqls(SplitSqlsInput request, ServerCallContext context) {
            var output = new SplitSqlsOutput();
            var version = request.Version;
            var sqls = request.Sqls;
            var statementList = ParseStatementList(version, sqls);

            foreach (var statement in statementList.Statements) {
                var sql = "";
                for (int index = statement.FirstTokenIndex; index <= statement.LastTokenIndex; index++) {
                    sql += statement.ScriptTokenStream[index].Text;
                }

                var splitSql = new Sql();
                splitSql.Sql_ = sql;
                splitSql.IsDDL = IsDDL(statement);
                splitSql.IsDML = IsDML(statement);
                output.SplitSqls.Add(splitSql);
            }

            return Task.FromResult(output);
        }

        public bool IsDDL(TSqlStatement statement) {
            if (statement is CreateDatabaseStatement) {
                return true;
            } else if (statement is CreateTableStatement) {
                return true;
            } else if (statement is AlterTableStatement) {
                return true;
            } else if (statement is CreateIndexStatement) {
                return true;
            } else if (statement is DropIndexStatement) {
                return true;
            } else if (statement is DropTableStatement) {
                return true;
            }
            return false;
        }

        public bool IsDML(TSqlStatement statement) {
            if (statement is InsertStatement) {
                return true;
            } else if (statement is UpdateStatement) {
                return true;
            } else if (statement is DeleteStatement) {
                return true;
            }
            return false;
        }

        // Audit implement
        public override Task<AdviseOutput> Advise(AdviseInput request, ServerCallContext context) {
            var output = new AdviseOutput();
            var version = request.Version;
            var sqls = request.Sqls;
            var ruleNames = request.RuleNames;
            var ruleValidatorContext = new RuleValidatorContext(request.SqlserverMeta);

            foreach (var sql in sqls) {
                var statementList = ParseStatementList(version, sql);
                foreach (var statement in statementList.Statements) {
                    foreach (var ruleName in ruleNames) {
                        var ruleValidator = DefaultRules.RuleValidators[ruleName];
                        if (ruleValidator == null) {
                            continue;
                        }
                        ruleValidator.Check(ruleValidatorContext, statement);
                    }

                    ruleValidatorContext.UpdateContext(statement);

                    output.AdviseResults.Add(ruleValidatorContext.AdviseResultContext.GetAdviseResult());
                    ruleValidatorContext.AdviseResultContext.ResetAdviseResult();
                }
            }

            return Task.FromResult(output);
        }

        // GetRollbackSqls implement
        public override Task<GetRollbackSqlsOutput> GetRollbackSqls(GetRollbackSqlsInput request, ServerCallContext context) {
            return base.GetRollbackSqls(request, context);
        }
    }
}