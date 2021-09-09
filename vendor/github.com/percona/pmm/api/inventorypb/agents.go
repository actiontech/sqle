package inventorypb

//go-sumtype:decl Agent

// Agent is a common interface for all types of Agents.
type Agent interface {
	sealedAgent() //nolint:unused
}

// in order of AgentType enum

func (*PMMAgent) sealedAgent()                        {}
func (*VMAgent) sealedAgent()                         {}
func (*NodeExporter) sealedAgent()                    {}
func (*MySQLdExporter) sealedAgent()                  {}
func (*MongoDBExporter) sealedAgent()                 {}
func (*PostgresExporter) sealedAgent()                {}
func (*ProxySQLExporter) sealedAgent()                {}
func (*QANMySQLPerfSchemaAgent) sealedAgent()         {}
func (*QANMySQLSlowlogAgent) sealedAgent()            {}
func (*QANMongoDBProfilerAgent) sealedAgent()         {}
func (*QANPostgreSQLPgStatementsAgent) sealedAgent()  {}
func (*QANPostgreSQLPgStatMonitorAgent) sealedAgent() {}
func (*RDSExporter) sealedAgent()                     {}
func (*ExternalExporter) sealedAgent()                {}
func (*AzureDatabaseExporter) sealedAgent()           {}
