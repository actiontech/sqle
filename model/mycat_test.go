package model

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInspector(t *testing.T) {
	server := `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mycat:server SYSTEM "server.dtd">
<mycat:server xmlns:mycat="http://io.mycat/">
	<user name="root">
		<property name="password">asd2010</property>
		<property name="schemas">masterdb,singledb</property>
	</user>

</mycat:server>
`
	schema := `
<?xml version="1.0"?>
<!DOCTYPE mycat:schema SYSTEM "schema.dtd">
<mycat:schema xmlns:mycat="http://io.mycat/">
	<schema name="masterdb" checkSQLschema="false" sqlMaxLimit="100">
		<table name="tb1" dataNode="dn1,dn2,dn3" rule="sharding-by-intfile"/>
		<table name="tb2" dataNode="dn1,dn2,dn3" rule="sharding-by-intfile"/>
	</schema>
	<schema name="singledb" checkSQLschema="false" sqlMaxLimit="100" dataNode="dn4"/>	
	<dataNode name="dn1" dataHost="host1" database="masterdb"/>
	<dataNode name="dn2" dataHost="host2" database="masterdb"/>
	<dataNode name="dn3" dataHost="host3" database="masterdb"/>
	<dataNode name="dn4" dataHost="host1" database="singledb"/>

	<dataHost name="host1" maxCon="1000" minCon="10" balance="1" writeType="0" dbType="mysql" dbDriver="native" switchType="-1" slaveThreshold="100">
		<heartbeat>select user()</heartbeat>
		<writeHost host="m1" url="172.20.130.2:3306" user="root" password="m1test"></writeHost>
	</dataHost>

	<dataHost name="host2" maxCon="1000" minCon="10" balance="1" writeType="0" dbType="mysql" dbDriver="native" switchType="-1" slaveThreshold="100">
                <heartbeat>select user()</heartbeat>
                <writeHost host="s1" url="172.20.130.3:3306" user="root" password="s1test"></writeHost>
        </dataHost>

	<dataHost name="host3" maxCon="1000" minCon="10" balance="1" writeType="0" dbType="mysql" dbDriver="native" switchType="-1" slaveThreshold="100">
                <heartbeat>select user()</heartbeat>
                <writeHost host="s2" url="172.20.130.4:3306" user="root" password="s2test"></writeHost>
        </dataHost>

</mycat:schema>
`
	rule := `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mycat:rule SYSTEM "rule.dtd">
<mycat:rule xmlns:mycat="http://io.mycat/">
	<tableRule name="rule1">
		<rule>
			<columns>id</columns>
			<algorithm>func1</algorithm>
		</rule>
	</tableRule>

	<tableRule name="rule2">
		<rule>
			<columns>user_id</columns>
			<algorithm>func1</algorithm>
		</rule>
	</tableRule>

	<tableRule name="sharding-by-intfile">
		<rule>
			<columns>sharding_id</columns>
			<algorithm>hash-int</algorithm>
		</rule>
	</tableRule>
</mycat:rule>
`
	serverXML := &ServerXML{}
	err := xml.Unmarshal([]byte(server), serverXML)
	if err != nil {
		t.Error(err)
		return
	}
	schemaXML := &SchemasXML{}
	err = xml.Unmarshal([]byte(schema), schemaXML)
	if err != nil {
		t.Error(err)
		return
	}
	rulesXML := &RulesXML{}
	err = xml.Unmarshal([]byte(rule), rulesXML)
	if err != nil {
		t.Error(err)
		return
	}
	instance, err := LoadMycatServerFromXML(serverXML, schemaXML, rulesXML)
	if err != nil {
		t.Error(err)
		return
	}
	mycatConfig := instance.MycatConfig
	assert.Len(t, mycatConfig.DataHosts, 3)
	assert.Len(t, mycatConfig.AlgorithmSchemas, 2)
	assert.Len(t, mycatConfig.AlgorithmSchemas["masterdb"].AlgorithmTables, 2)
	assert.Nil(t, mycatConfig.AlgorithmSchemas["masterdb"].DataNode)
	assert.Nil(t, mycatConfig.AlgorithmSchemas["singledb"].AlgorithmTables)
	assert.NotNil(t, mycatConfig.AlgorithmSchemas["singledb"].DataNode)
}

func TestSplitMultiNodes(t *testing.T) {
	assert.Equal(t, []string{"abc1", "abc2", "abc3", "abc4", "abc5"},
		splitMultiNodes("abc$1-5"))
	assert.Equal(t, []string{"abc11", "abc12", "abc13", "abc14", "abc15", "abc21", "abc22"},
		splitMultiNodes("abc1$1-5,abc2$1-2"))
	assert.Equal(t, []string{"abc"}, splitMultiNodes("abc"))
	assert.Equal(t, []string{"abc1", "abc2"}, splitMultiNodes("abc1,abc2"))
	assert.Equal(t, []string{"abc11", "abc12", "abc2"}, splitMultiNodes("abc1$1-2,abc2"))
	assert.Equal(t, []string{"abc1", "abc21", "abc22"}, splitMultiNodes("abc1,abc2$1-2"))
}
