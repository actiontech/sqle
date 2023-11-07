/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"dm/i18n"
	"sync"
)

// 发版标记
var version = "8.1.2.174"
var build_date = "2022.11.21"
var svn = "13574"

var globalDmDriver = newDmDriver()

func init() {
	sql.Register("dm", globalDmDriver)
}

func driverInit(svcConfPath string) {
	load(svcConfPath)
	if GlobalProperties != nil && GlobalProperties.Len() > 0 {
		setDriverAttributes(GlobalProperties)
	}
	globalDmDriver.createFilterChain(nil, GlobalProperties)

	switch Locale {
	case 0:
		i18n.InitConfig(i18n.Messages_zh_CN)
	case 1:
		i18n.InitConfig(i18n.Messages_en_US)
	case 2:
		i18n.InitConfig(i18n.Messages_zh_TW)
	}
}

type DmDriver struct {
	filterable
	readPropMutex sync.Mutex
}

func newDmDriver() *DmDriver {
	d := new(DmDriver)
	d.idGenerator = dmDriverIDGenerator
	return d
}

/*************************************************************
 ** PUBLIC METHODS AND FUNCTIONS
 *************************************************************/
func (d *DmDriver) Open(dsn string) (driver.Conn, error) {
	return d.open(dsn)
}

func (d *DmDriver) OpenConnector(dsn string) (driver.Connector, error) {
	return d.openConnector(dsn)
}

func (d *DmDriver) open(dsn string) (*DmConnection, error) {
	c, err := d.openConnector(dsn)
	if err != nil {
		return nil, err
	}
	return c.connect(context.Background())
}

func (d *DmDriver) openConnector(dsn string) (*DmConnector, error) {
	connector := new(DmConnector).init()
	connector.url = dsn
	connector.dmDriver = d
	d.readPropMutex.Lock()
	err := connector.mergeConfigs(dsn)
	d.readPropMutex.Unlock()
	if err != nil {
		return nil, err
	}
	connector.createFilterChain(connector, nil)
	return connector, nil
}
