package inventorypb

//go-sumtype:decl Service

// Service is a common interface for all types of Services.
type Service interface {
	sealedService() //nolint:unused
}

// in order of ServiceType enum

func (*MySQLService) sealedService()      {}
func (*MongoDBService) sealedService()    {}
func (*PostgreSQLService) sealedService() {}
func (*ProxySQLService) sealedService()   {}
func (*HAProxyService) sealedService()    {}
func (*ExternalService) sealedService()   {}
