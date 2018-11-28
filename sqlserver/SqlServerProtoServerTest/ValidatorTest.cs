using System;
using Xunit;
using SqlserverProtoServer;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using System.Collections.Generic;
using System.IO;

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

        [Fact]
        public void SelectAllRuleValidatorTest() {
            var validator = DefaultRules.RuleValidators[DefaultRules.DML_DISABE_SELECT_ALL_COLUMN];
            var ruleValidatorContext = new RuleValidatorContext();

            // not select *
            var statementList = ParseStatementList("select hello from world");
            foreach (var statement in statementList.Statements) {
                validator.Check(ruleValidatorContext, statement);
                var auditResult = ruleValidatorContext.AuditResultContext.GetAuditResult();
                Assert.Equal(AuditResultContext.RuleLevels[RULE_LEVEL.NORMAL], auditResult.AuditLevel);
                Assert.Equal("", auditResult.AuditResultMessage);
            }

            // with select *
            ruleValidatorContext.AuditResultContext.ResetAuditResult();
            statementList = ParseStatementList("select * from world");
            foreach (var statement in statementList.Statements) {
                validator.Check(ruleValidatorContext, statement);
                var auditResult = ruleValidatorContext.AuditResultContext.GetAuditResult();
                Assert.Equal(AuditResultContext.RuleLevels[RULE_LEVEL.NOTICE], auditResult.AuditLevel);
                Assert.Equal("[notice]不建议使用select *", auditResult.AuditResultMessage);
            }
        }
    }
}
