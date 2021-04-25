package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/universe/ucommon/v4/util"
)

type MycatConfig struct {
	AlgorithmSchemas map[string]*AlgorithmSchema `json:"schema_list"`
	DataHosts        map[string]*DataHost        `json:"data_host_list"`
}

type AlgorithmSchema struct {
	AlgorithmTables map[string]*AlgorithmTable `json:"table_list"`
	DataNode        *DataNode                  `json:"data_node"`
}

type AlgorithmTable struct {
	name           string
	ShardingColumn string      `json:"sharding_columns"`
	DataNodes      []*DataNode `json:"data_node_list"`
}

type DataNode struct {
	DataHostName string `json:"data_host"`
	Database     string `json:"database"`
}

type DataHost struct {
	User     string        `json:"user"`
	Host     string        `json:"host"`
	Port     string        `json:"port"`
	Password util.Password `json:"password"`
}

func (m *MycatConfig) IsShardingSchema(schemaName string) (bool, error) {
	if m.AlgorithmSchemas != nil {
		if schema, ok := m.AlgorithmSchemas[schemaName]; ok {
			if schema.AlgorithmTables != nil {
				return true, nil
			} else {
				return false, nil
			}
		}
	}
	return false, fmt.Errorf("schema %s not found in mycat config", schemaName)
}

func (m *MycatConfig) GetShardingColumn(schemaName, tableName string) (string, error) {
	ok, err := m.IsShardingSchema(schemaName)
	if err != nil || !ok {
		return "", err
	}
	table, ok := m.AlgorithmSchemas[schemaName].AlgorithmTables[tableName]
	if !ok {
		return "", fmt.Errorf("table %s not found in mycat config", tableName)
	}
	return table.ShardingColumn, nil
}

func LoadMycatServerFromXML(serverText, schemasText, rulesText []byte) (*Instance, error) {
	var err error
	serverXML := &ServerXML{}
	err = xml.Unmarshal(serverText, serverXML)
	if err != nil {
		return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR, err)
	}
	schemasXML := &SchemasXML{}
	err = xml.Unmarshal(schemasText, schemasXML)
	if err != nil {
		return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR, err)
	}
	rulesXML := &RulesXML{}
	err = xml.Unmarshal(rulesText, rulesXML)
	if err != nil {
		return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR, err)
	}
	var instance = &Instance{
		DbType: DB_TYPE_MYCAT,
	}
	instance.MycatConfig = &MycatConfig{
		AlgorithmSchemas: map[string]*AlgorithmSchema{},
		DataHosts:        map[string]*DataHost{},
	}

	// load all dataHost from schema.xml
	var allDataHosts = map[string]*DataHost{}
	for _, hosts := range schemasXML.DataHosts {
		// just get first writeHost, if writeHost exists
		if len(hosts.WriteHosts) >= 1 {
			host := hosts.WriteHosts[0]
			ip, port := unmarshalUrl(host.Url)
			allDataHosts[hosts.Name] = &DataHost{
				User: host.User,
				Port: port,
				Host: ip,
			}
		}
	}

	// load all dataNode from schema.xml
	// name, database, dataHost can use ',', '$', '-' to configure multi nodes
	// but the (database size) * (dataHost size) must equal the size of name
	// every dataHost has all database in its tag
	// eg: <dataNode name="dn1$0-750" dataHost="host$1-10" database="db$0-75" />
	// host1 has database of dn1$0-75 (name is dn$1-75)
	// host2 has database of dn1$0-75 (name is dn$76-151)
	var allDataNodes = map[string]*DataNode{}
	for _, dataNode := range schemasXML.DataNodes {
		names := splitMultiNodes(dataNode.Name)
		databases := splitMultiNodes(dataNode.Database)
		dataHosts := splitMultiNodes(dataNode.DataHostName)
		if len(names) != len(databases)*len(dataHosts) {
			return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR,
				fmt.Errorf("parser dataNode %s failed", dataNode.Name))
		}
		for n, name := range names {
			allDataNodes[name] = &DataNode{
				DataHostName: dataHosts[n/len(databases)],
				Database:     databases[n%len(databases)],
			}
		}
	}

	// load all schema form schema.xml
	var AllAlgorithmSchemas = map[string]*AlgorithmSchema{}
	for _, schema := range schemasXML.Schemas {
		as := &AlgorithmSchema{}
		if schema.Tables != nil {
			as.AlgorithmTables = map[string]*AlgorithmTable{}
			for _, table := range schema.Tables {
				t := &AlgorithmTable{
					name:      table.Name,
					DataNodes: []*DataNode{},
				}
				nodeNames := splitMultiNodes(table.DataNodeList)
				for _, nodeName := range nodeNames {
					node, ok := allDataNodes[nodeName]
					if !ok {
						err = fmt.Errorf("dataNode %s not found in schema.xml", nodeName)
						return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR, err)
					}
					t.DataNodes = append(t.DataNodes, node)
				}

				if table.RuleName != "" {
					rule, exist := rulesXML.getRuleByName(table.RuleName)
					if !exist {
						err = fmt.Errorf("rule %s not found in rule.xml", table.RuleName)
						return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR, err)
					}
					t.ShardingColumn = rule.ShardingColumn
				}

				as.AlgorithmTables[table.Name] = t
			}
		}
		if schema.DataNodeName != "" {
			node, ok := allDataNodes[schema.DataNodeName]
			if !ok {
				err = fmt.Errorf("dataNode %s not found in schema.xml", schema.DataNodeName)
				return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR, err)
			}
			as.DataNode = node
		}
		AllAlgorithmSchemas[schema.Name] = as
	}

	instance.User = serverXML.User.Name

	schemas := []string{}
	for _, property := range serverXML.User.PropertyList {
		if property.Name == "schemas" {
			schemas = strings.Split(property.Value, ",")
			break
		}
	}
	if len(schemas) <= 0 {
		return instance, nil
	}

	for _, schema := range schemas {
		s, ok := AllAlgorithmSchemas[schema]
		if !ok {
			err = fmt.Errorf("schema %s not found in schema.xml", schema)
			return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR, err)
		}
		instance.MycatConfig.AlgorithmSchemas[schema] = s
		if s.DataNode != nil {
			instance.MycatConfig.DataHosts[s.DataNode.DataHostName] = allDataHosts[s.DataNode.DataHostName]
		}
		for _, table := range s.AlgorithmTables {
			for _, node := range table.DataNodes {
				_, ok := allDataHosts[node.DataHostName]
				if !ok {
					err = fmt.Errorf("dataHost %s not found in schema.xml", node.DataHostName)
					return nil, errors.New(errors.PARSER_MYCAT_CONFIG_ERROR, err)
				}
				instance.MycatConfig.DataHosts[node.DataHostName] = allDataHosts[node.DataHostName]
			}
		}
	}
	return instance, nil
}

// ServerXML is the unmarshal struct object for server.xml
type ServerXML struct {
	XMLName xml.Name `xml:"server"`
	User    *UserXML `xml:"user"`
}

type UserXML struct {
	Name         string         `xml:"name,attr"`
	PropertyList []*PropertyXML `xml:"property"`
}

type PropertyXML struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",innerxml"`
}

// SchemasXML is the unmarshal struct object for schema.xml
type SchemasXML struct {
	XMLName   xml.Name       `xml:"schema"`
	Schemas   []*SchemaXML   `xml:"schema"`
	DataNodes []*DataNodeXML `xml:"dataNode"`
	DataHosts []*DataHostXML `xml:"dataHost"`
}

type SchemaXML struct {
	Name         string      `xml:"name,attr"`
	DataNodeName string      `xml:"dataNode,attr"`
	Tables       []*TableXML `xml:"table"`
}

type TableXML struct {
	Name         string `xml:"name,attr"`
	DataNodeList string `xml:"dataNode,attr"` // "dn1,dn2"
	RuleName     string `xml:"rule,attr"`
}

type DataNodeXML struct {
	Name         string `xml:"name,attr"`
	DataHostName string `xml:"dataHost,attr"`
	Database     string `xml:"database,attr"`
}

type DataHostXML struct {
	Name       string          `xml:"name,attr"`
	WriteHosts []*WriteHostXML `xml:"writeHost"`
}

type WriteHostXML struct {
	Url  string `xml:"url,attr"`
	User string `xml:"user,attr"`
}

// RulesXML is the unmarshal struct object for rule.xml
type RulesXML struct {
	XMLName xml.Name   `xml:"rule"`
	Rules   []*RuleXML `xml:"tableRule"`
}

type RuleXML struct {
	Name           string `xml:"name,attr"`
	ShardingColumn string `xml:"rule>columns"`
}

func (r *RulesXML) getRuleByName(name string) (*RuleXML, bool) {
	if r == nil {
		return nil, false
	}
	for _, rule := range r.Rules {
		if rule.Name == name {
			return rule, true
		}
	}
	return nil, false
}

func (s *SchemasXML) getDataNode(name string) (*DataNodeXML, error) {
	var node *DataNodeXML
	for _, n := range s.DataNodes {
		if name == n.Name {
			node = n
			break
		}
	}
	if node == nil {
		return nil, fmt.Errorf("dataNode %s not found in schema.xml", name)
	}
	return node, nil
}

func unmarshalUrl(url string) (host, port string) {
	u := strings.Split(url, ":")
	if len(url) >= 2 {
		host = u[len(u)-2]
		port = u[len(u)-1]
	}
	return
}

// splitMultiNodes split string to list, by using mycat config rule
// eg: input: node1$1-5,node2$1-2 output: [node11, node12, node13, node14, node15, node21, node22]
func splitMultiNodes(nodes string) []string {
	result := []string{}
	nodeList := strings.Split(nodes, ",")
	for _, node := range nodeList {
		s := strings.Split(node, "$")
		if len(s) != 2 {
			result = append(result, node)
			continue
		}
		prefix := s[0]
		interval := strings.Split(s[1], "-")
		if len(interval) != 2 {
			result = append(result, node)
			continue
		}
		min, err := strconv.ParseInt(interval[0], 10, 64)
		if err != nil {
			fmt.Println(err)
			result = append(result, node)
			continue
		}
		max, err := strconv.ParseInt(interval[1], 10, 64)
		if err != nil {
			result = append(result, node)
			continue
		}
		if min > max {
			result = append(result, node)
			continue
		}
		for i := min; i <= max; i++ {
			result = append(result, fmt.Sprintf("%s%d", prefix, i))
		}
	}
	return result
}
