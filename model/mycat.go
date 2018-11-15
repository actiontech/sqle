package model

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type MycatServer struct {
	Host             string
	User             string
	Port             string
	Password         string
	AlgorithmSchemas map[string]AlgorithmSchema
	DataHosts        map[string]Instance
}

type AlgorithmSchema struct {
	AlgorithmTables map[string]AlgorithmTable
	DataNodes       string
}

type AlgorithmTable struct {
	name           string
	ShardingColumn string
	DataNodes      map[string]ShardingSchema
}

type ShardingSchema struct {
	DataHostName string
	Database     string
}

func LoadMycatServerConfig(server, schema, rule []byte) (*MycatServer, error) {
	var mycat = &MycatServer{
		AlgorithmSchemas: map[string]AlgorithmSchema{},
		DataHosts:        map[string]Instance{},
	}

	var rulesXML = new(RulesXML)
	var schemasXML = new(SchemasXML)
	var serverXML = new(ServerXML)
	var err error
	err = xml.Unmarshal(server, serverXML)
	if err != nil {
		goto ERROR
	}
	err = xml.Unmarshal(schema, schemasXML)
	if err != nil {
		goto ERROR
	}
	err = xml.Unmarshal(rule, rulesXML)
	if err != nil {
		goto ERROR
	}

	for _, hosts := range schemasXML.DataHosts {
		// just get first writeHost, if writeHost exists
		if len(hosts.WriteHosts) >= 1 {
			host := hosts.WriteHosts[0]
			ip := ""
			port := ""
			url := strings.Split(host.Url, ":")
			if len(url) >= 2 {
				ip = url[len(url)-2]
				port = url[len(url)-1]
			}
			mycat.DataHosts[hosts.Name] = Instance{
				User:   host.User,
				Port:   port,
				Host:   ip,
				DbType: "mysql",
			}
		}
	}

	for _, schema := range schemasXML.Schemas {
		as := AlgorithmSchema{}
		if schema.Tables != nil {
			for _, table := range schema.Tables {
				t := &AlgorithmTable{
					name: table.Name,
				}
				if table.RuleName != "" {
					rule, exist := rulesXML.getRuleByName(table.RuleName)
					if !exist {
						err = fmt.Errorf("rule %s not found in rule.xml", table.RuleName)
						goto ERROR
					}
					t.ShardingColumn = rule.ShardingColumn
				}
				node, err := schemasXML.getDataNode(table.DataNodeName)
				if err != nil {
					goto ERROR
				}
			}
		}
		as.DataNodes = schema.DataNodeName

		mycat.AlgorithmSchemas[schema.Name] = as
	}

ERROR:
	return nil, err
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
	DataNodeName string `xml:"dataNode,attr"`
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

func (r *SchemasXML) getDataNode(name string) (*DataNodeXML, error) {
	var node *DataNodeXML
	for _, n := range r.DataNodes {
		if n.DataHostName == name {
			node = n
			break
		}
	}
	if node == nil {
		return nil, fmt.Errorf("dataNode %s not found in schema.xml", name)
	}
	count := 0
	hosts := strings.Split(node.DataHostName, ",")
	for _, host := range hosts {
		for _, dataHost := range r.DataHosts {
			if dataHost.Name == host {
				count += 1
				continue
			}
		}
	}
	if len(hosts) > count {
		return nil, fmt.Errorf("dataNode ")
	}
	return node, nil
}
