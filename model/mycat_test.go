package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testServer = `
<?xml version="1.0" encoding="UTF-8"?>
<mycat:server xmlns:mycat="http://io.mycat/">
	<user name="root">
		<property name="password">asd2010</property>
		<property name="schemas">masterdb,singledb</property>
	</user>

</mycat:server>
`
var testRule = `
<?xml version="1.0" encoding="UTF-8"?>
<mycat:rule xmlns:mycat="http://io.mycat/">
	<tableRule name="sharding-by-intfile">
		<rule>
			<columns>sharding_id</columns>
			<algorithm>hash-int</algorithm>
		</rule>
	</tableRule>
</mycat:rule>
`

func TestNewParserMycatConfig_Normal(t *testing.T) {
	schema := `
<?xml version="1.0"?>
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

	<dataHost name="host1">
		<heartbeat>select user()</heartbeat>
		<writeHost host="m1" url="172.20.130.2:3306" user="root" password="m1test"></writeHost>
	</dataHost>

	<dataHost name="host2">
		<heartbeat>select user()</heartbeat>
		<writeHost host="s1" url="172.20.130.3:3306" user="root" password="s1test"></writeHost>
	</dataHost>

	<dataHost name="host3">
		<heartbeat>select user()</heartbeat>
		<writeHost host="s2" url="172.20.130.4:3306" user="root" password="s2test"></writeHost>
	</dataHost>
</mycat:schema>
`
	instance, err := LoadMycatServerFromXML([]byte(testServer), []byte(schema), []byte(testRule))
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

func TestNewParserMycatConfig_Wildcard(t *testing.T) {
	schema := `
<?xml version="1.0"?>
<mycat:schema xmlns:mycat="http://io.mycat/">
	<schema name="masterdb" checkSQLschema="false" sqlMaxLimit="100">
		<table name="tb1" dataNode="dn$1-15" rule="sharding-by-intfile"/>
	</schema>
	<schema name="singledb" checkSQLschema="false" sqlMaxLimit="100" dataNode="dn2"/>	
	<dataNode name="dn$1-15" dataHost="host$1-3" database="db$1-5"/>
	<dataNode name="dn2" dataHost="host1" database="db2"/>
	<dataHost name="host1">
		<heartbeat>select user()</heartbeat>
		<writeHost host="m1" url="172.20.130.2:3306" user="root" password="m1test"></writeHost>
	</dataHost>

	<dataHost name="host2">
		<heartbeat>select user()</heartbeat>
		<writeHost host="s1" url="172.20.130.3:3306" user="root" password="s1test"></writeHost>
	</dataHost>

	<dataHost name="host3">
		<heartbeat>select user()</heartbeat>
		<writeHost host="s2" url="172.20.130.4:3306" user="root" password="s2test"></writeHost>
	</dataHost>
	<dataHost name="host4">
		<heartbeat>select user()</heartbeat>
		<writeHost host="s2" url="172.20.130.4:3306" user="root" password="s2test"></writeHost>
	</dataHost>
</mycat:schema>
`
	instance, err := LoadMycatServerFromXML([]byte(testServer), []byte(schema), []byte(testRule))
	if err != nil {
		t.Error(err)
		return
	}
	mycatConfig := instance.MycatConfig
	assert.Len(t, mycatConfig.DataHosts, 3)
	assert.NotNil(t, mycatConfig.DataHosts["host1"])
	assert.NotNil(t, mycatConfig.DataHosts["host2"])
	assert.NotNil(t, mycatConfig.DataHosts["host3"])
	assert.Nil(t, mycatConfig.DataHosts["host4"])
	assert.Len(t, mycatConfig.AlgorithmSchemas, 2)
	assert.Len(t, mycatConfig.AlgorithmSchemas["masterdb"].AlgorithmTables, 1)
	assert.Nil(t, mycatConfig.AlgorithmSchemas["masterdb"].DataNode)
	table := mycatConfig.AlgorithmSchemas["masterdb"].AlgorithmTables["tb1"]
	assert.NotNil(t, table)
	nodes := []string{}
	for _, node := range table.DataNodes {
		nodes = append(nodes, fmt.Sprintf("%s:%s", node.DataHostName, node.Database))
	}
	expect := []string{"host1:db1", "host1:db2", "host1:db3", "host1:db4", "host1:db5",
		"host2:db1", "host2:db2", "host2:db3", "host2:db4", "host2:db5",
		"host3:db1", "host3:db2", "host3:db3", "host3:db4", "host3:db5"}
	assert.Equal(t, expect, nodes)
}

func TestSplitMultiNodes(t *testing.T) {
	assert.Equal(t, []string{"abc1", "abc2", "abc3", "abc4", "abc5"},
		splitMultiNodes("abc$1-5"))
	assert.Equal(t, []string{"abc11", "abc12", "abc13", "abc14", "abc15", "abc21", "abc22"},
		splitMultiNodes("abc1$1-5,abc2$1-2"))
	assert.Equal(t, []string{"abc"}, splitMultiNodes("abc"))
	assert.Equal(t, []string{"abc1", "abc2"}, splitMultiNodes("abc1,abc2"))
	assert.Equal(t, []string{"abc11", "abc12", "abc13", "abc14", "abc15", "abc2"}, splitMultiNodes("abc1$1-5,abc2"))
	assert.Equal(t, []string{"abc1", "abc21", "abc22"}, splitMultiNodes("abc1,abc2$1-2"))
}
