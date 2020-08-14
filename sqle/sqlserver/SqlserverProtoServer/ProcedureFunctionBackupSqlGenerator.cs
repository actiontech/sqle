using System.Collections.Generic;
using Microsoft.SqlServer.TransactSql.ScriptDom;
using NLog;
using SqlserverProto;
using System;
using System.Data.SqlClient;

namespace SqlserverProtoServer {
    public class ProcedureFunctionBackupSqlGenerator {
        public SqlserverMeta SqlserverMeta;

        public ProcedureFunctionBackupSqlGenerator(SqlserverMeta sqlserverMeta) {
            this.SqlserverMeta = sqlserverMeta;
        }

        public String GetConnectionString() {
            return String.Format(
                "Server=tcp:{0},{1};" +
                "Database=master;" +
                "User ID={2};" +
                "Password={3};",
                SqlserverMeta.Host, SqlserverMeta.Port,
                SqlserverMeta.User,
                SqlserverMeta.Password);
        }

        public string GetObjectDefinition(Logger logger, string objectname) {
            var ret = "";
            String connectionString = GetConnectionString();
            using (SqlConnection connection = new SqlConnection(connectionString)) {
                var commandStr = String.Format("SELECT OBJECT_DEFINITION (OBJECT_ID('{0}')) AS ObjectDefinition", objectname);
                logger.Info("sql query: {0}", commandStr);
                SqlCommand command = new SqlCommand(commandStr, connection);
                connection.Open();
                SqlDataReader reader = command.ExecuteReader();
                try {
                    while (reader.Read()) {
                        ret = (string)reader["ObjectDefinition"];
                    }
                } finally {
                    reader.Close();
                }
            }

            return ret;
        }

        public string GetName(SchemaObjectName schemaObjectName) {
            var identifiers = new List<string>();
            if (schemaObjectName.SchemaIdentifier != null) {
                identifiers.Add(schemaObjectName.SchemaIdentifier.Value);
            }
            if (schemaObjectName.BaseIdentifier != null) {
                identifiers.Add(schemaObjectName.BaseIdentifier.Value);
            }
            return string.Join('.', identifiers);
        }

        public string getBackupSql(Logger logger, string objectName) {
            var backupSql = "";
            try {
                var objectDefinition = GetObjectDefinition(logger, objectName);
                var newObjectName = String.Format("{0}_{1}", objectName, DateTime.Now.ToString("yyyy_MM_dd_HH_mm"));
                backupSql = objectDefinition.Replace(objectName, newObjectName);
            } catch (Exception e) {
                logger.Fatal("GetBackupSql for {0} error, message: {1}", objectName, e.Message);
                logger.Fatal("GetBackupSql for {0} error, stacktrace: {1}", objectName, e.StackTrace);
            }

            return backupSql;
        }

        public List<string> GetbackupSqlsForObejcts(Logger logger, IList<SchemaObjectName> objects) {
            var backupSqls = new List<string>();
            foreach (var obj in objects) {
                var objName = GetName(obj);
                var backupSql = getBackupSql(logger, objName);
                if (backupSql != "") {
                    logger.Info("backupSql for {0} is {1}", objName, backupSql);
                    backupSqls.Add(backupSql);
                }
            }
            return backupSqls;
        }
    }
}
