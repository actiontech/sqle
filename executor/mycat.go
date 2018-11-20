package executor

import (
	"database/sql/driver"
	"sqle/model"
)

type MycatConn struct {
	Mycat    *BaseConn
	DataHost map[string]*BaseConn
	config   *model.MycatConfig
}

func newMycatConn(instance *model.Instance, schema string) (Db, error) {
	var mc = &MycatConn{
		DataHost: map[string]*BaseConn{},
		config:   instance.MycatConfig,
	}
	conn, err := newConn(instance, schema)
	if err != nil {
		return nil, err
	}
	mc.Mycat = conn
	for name, dataHost := range mc.config.DataHosts {
		conn, err := newConn(&model.Instance{
			DbType:   model.DB_TYPE_MYSQL,
			Host:     dataHost.Host,
			Port:     dataHost.Port,
			User:     dataHost.User,
			Password: dataHost.Password,
		}, schema)
		if err != nil {
			return nil, err
		}
		mc.DataHost[name] = conn
	}
	return mc, nil
}

func (mc *MycatConn) Close() {
	mc.Mycat.Close()
	for _, conn := range mc.DataHost {
		conn.Close()
	}
}

func (mc *MycatConn) Ping() error {
	if err := mc.Mycat.Ping(); err != nil {
		return err
	}
	//TODO: ping using goroutine
	for _, conn := range mc.DataHost {
		if err := conn.Ping(); err != nil {
			return err
		}
	}
	return nil
}

func (mc *MycatConn) Exec(query string) (driver.Result, error) {
	return nil, nil
}

func (mc *MycatConn) Query(query string, args ...interface{}) ([]map[string]string, error) {
	return nil, nil
}
