package executor

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/sirupsen/logrus"
	"sqle/errors"
	"sqle/model"
)

type MycatConn struct {
	log    *logrus.Entry
	Conn   *BaseConn
	config *model.MycatConfig
}

func newMycatConn(entry *logrus.Entry, instance *model.Instance, schema string) (Db, error) {
	var mc = &MycatConn{
		log:    entry,
		config: instance.MycatConfig,
	}
	conn, err := newConn(entry, instance, schema)
	if err != nil {
		entry.Error(err)
		return nil, err
	}
	mc.Conn = conn
	return mc, nil
}

func (mc *MycatConn) openDataHostConn(name, schema string) (Db, error) {
	dataHost, ok := mc.config.DataHosts[name]
	if !ok {
		msg := fmt.Errorf("data host %s not found", name)
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, msg)
	}
	conn, err := newConn(mc.log, &model.Instance{
		DbType:   model.DB_TYPE_MYSQL,
		Host:     dataHost.Host,
		Port:     dataHost.Port,
		User:     dataHost.User,
		Password: string(dataHost.Password),
	}, schema)
	if err != nil {
		mc.Logger().Error("connect mycat data host failed")
		return nil, err
	}
	return conn, nil
}

func (mc *MycatConn) Close() {
	mc.Conn.Close()
}

func (mc *MycatConn) Ping() error {
	if err := mc.Conn.Ping(); err != nil {
		return err
	}
	return nil
}

func (mc *MycatConn) Exec(query string) (driver.Result, error) {
	return mc.Conn.Exec(query)
}

func (mc *MycatConn) Transact(qs ...string) error {
	return mc.Conn.Transact(qs...)
}

func (mc *MycatConn) ExecDDL(query, schema, table string) error {
	as, ok := mc.config.AlgorithmSchemas[schema]
	if !ok {
		msg := fmt.Errorf("schema %s not found in mycat algorithm schemas", schema)
		mc.log.Error(msg)
		return errors.New(errors.CONNECT_REMOTE_DB_ERROR, msg)
	}
	if as.AlgorithmTables == nil {
		if as.DataNode != nil {
			conn, err := mc.openDataHostConn(as.DataNode.DataHostName, as.DataNode.Database)
			if err != nil {
				return errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
			}
			defer conn.Close()
			return conn.ExecDDL(query, "", "")
		}
	} else {
		at, ok := as.AlgorithmTables[table]
		if !ok {
			msg := fmt.Errorf("table %s not found in mycat algorithm schema %s", table, schema)
			mc.log.Error(msg)
			return errors.New(errors.CONNECT_REMOTE_DB_ERROR, msg)
		}
		conns := []Db{}
		defer func() {
			for _, conn := range conns {
				conn.Close()
			}
		}()

		for _, node := range at.DataNodes {
			conn, err := mc.openDataHostConn(node.DataHostName, node.Database)
			if err != nil {
				return errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
			}
			conns = append(conns, conn)
		}
		for _, conn := range conns {
			err := conn.ExecDDL(query, "", "")
			if err != nil {
				return errors.New(errors.CONNECT_REMOTE_DB_ERROR, err)
			}
		}
	}
	return nil
}

func (mc *MycatConn) Query(query string, args ...interface{}) ([]map[string]sql.NullString, error) {
	return mc.Conn.Query(query, args...)
}

func (mc *MycatConn) Logger() *logrus.Entry {
	return mc.log
}
