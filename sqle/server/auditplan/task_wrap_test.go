package auditplan

import (
	"testing"

	"github.com/actiontech/sqle/sqle/model"
)

func TestIsSqlInBlackList(t *testing.T) {
	filter := ConvertToBlackFilter([]*model.BlackListAuditPlanSQL{
		{
			FilterContent: "SELECT",
			FilterType:    "sql",
		}, {
			FilterContent: "table_1",
			FilterType:    "sql",
		}, {
			FilterContent: "ignored_service",
			FilterType:    "sql",
		},
	})

	matchSqls := []string{
		"SELECT * FROM users",
		"DELETE From tAble_1",
		"SELECT COUNT(*) FROM table_2",
		`/* this is a comment, Service: ignored_service */ 
		select * from table_ignored where id < 123;`,
		`/* this is a comment, Service: ignored_service */ update * from table_ignored where id < 123;`,
	}
	for _, matchSql := range matchSqls {
		if _, isSqlInBlackList := filter.IsSqlInBlackList(matchSql); !isSqlInBlackList {
			t.Error("Expected SQL to match blacklist")
		}
	}
	notMatchSqls := []string{
		"INSERT INTO users VALUES (1, 'John')",
		"DELETE  From schools",
		"SHOW CREATE TABLE table_2",
		`/* this is a comment, Service: ignored_
		service */ update * from table_ignored where id < 123;`,
	}
	for _, notMatchSql := range notMatchSqls {
		if _, isSqlInBlackList := filter.IsSqlInBlackList(notMatchSql); isSqlInBlackList {
			t.Error("Did not expect SQL to match blacklist")
		}
	}
}

func TestIsIpInBlackList(t *testing.T) {
	filter := ConvertToBlackFilter([]*model.BlackListAuditPlanSQL{
		{
			FilterContent: "192.168.1.23",
			FilterType:    "ip",
		}, {
			FilterContent: "10.0.5.67",
			FilterType:    "ip",
		},
	})

	matchIps := []string{
		"10.0.5.67",
		"192.168.1.23",
	}
	for _, matchIp := range matchIps {
		if _, hasEndpointInBlackList := filter.HasEndpointInBlackList([]string{matchIp}); !hasEndpointInBlackList {
			t.Error("Expected Ip to match blacklist")
		}
	}

	notMatchIps := []string{
		"172.16.254.89",
		"134.12.45.78",
		"50.67.89.12",
	}
	for _, notMatchIp := range notMatchIps {
		if _, hasEndpointInBlackList := filter.HasEndpointInBlackList([]string{notMatchIp}); hasEndpointInBlackList {
			t.Error("Did not expect Ip to match blacklist")
		}
	}
}

func TestIsCidrInBlackList(t *testing.T) {
	filter := ConvertToBlackFilter([]*model.BlackListAuditPlanSQL{
		{
			FilterContent: "192.168.0.0/24",
			FilterType:    "cidr",
		}, {
			FilterContent: "10.100.0.0/16",
			FilterType:    "cidr",
		},
	})

	matchIps := []string{
		"10.100.1.2",
		"10.100.25.45",
		"192.168.0.2",
		"192.168.0.45",
	}
	for _, matchIp := range matchIps {
		if _, hasEndpointInBlackList := filter.HasEndpointInBlackList([]string{matchIp}); !hasEndpointInBlackList {
			t.Error("Expected CIDR to match blacklist")
		}
	}

	notMatchIps := []string{
		"172.16.254.89",
		"134.12.45.78",
		"50.67.89.12",
		"172.30.1.2",
		"172.30.30.45",
	}
	for _, notMatchIp := range notMatchIps {
		if _, hasEndpointInBlackList := filter.HasEndpointInBlackList([]string{notMatchIp}); hasEndpointInBlackList {
			t.Error("Did not expect CIDR to match blacklist")
		}
	}
}

func TestIsHostInBlackList(t *testing.T) {
	filter := ConvertToBlackFilter([]*model.BlackListAuditPlanSQL{
		{
			FilterContent: "host",
			FilterType:    "host",
		}, {
			FilterContent: "some_site",
			FilterType:    "host",
		},
	})

	matchHosts := []string{
		"local_host",
		"local_Host.com",
		"any_Host.io",
		"some_Site.org/home/",
		"Some_site.cn/mysql",
	}

	for _, matchHost := range matchHosts {
		if _, hasEndpointInBlackList := filter.HasEndpointInBlackList([]string{matchHost}); !hasEndpointInBlackList {
			t.Error("Expected HOST to match blacklist")
		}
	}

	notMatchHosts := []string{
		"other_site/home",
		"any_other_site/local",
	}
	for _, noMatchHost := range notMatchHosts {
		if _, hasEndpointInBlackList := filter.HasEndpointInBlackList([]string{noMatchHost}); hasEndpointInBlackList {
			t.Error("Did not expect HOST to match blacklist")
		}
	}
}
