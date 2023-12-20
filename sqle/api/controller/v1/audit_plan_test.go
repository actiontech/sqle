package v1_test

import (
	"testing"

	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/model"
)

func TestIsSqlInBlackList(t *testing.T) {
	filter := v1.ConvertToBlackFilter([]*model.BlackListAuditPlanSQL{
		{
			FilterContent: "SELECT",
			FilterType:    "SQL",
		}, {
			FilterContent: "table_1",
			FilterType:    "SQL",
		},
	})

	matchSqls := []string{
		"SELECT * FROM users",
		"DELETE From tAble_1",
		"SELECT COUNT(*) FROM table_2",
	}
	for _, matchSql := range matchSqls {
		if !filter.IsSqlInBlackList(matchSql) {
			t.Error("Expected SQL to match blacklist")
		}
	}
	notMatchSqls := []string{
		"INSERT INTO users VALUES (1, 'John')",
		"DELETE  From schools",
		"SHOW CREATE TABLE table_2",
	}
	for _, notMatchSql := range notMatchSqls {
		if filter.IsSqlInBlackList(notMatchSql) {
			t.Error("Did not expect SQL to match blacklist")
		}
	}
}

func TestIsIpInBlackList(t *testing.T) {
	filter := v1.ConvertToBlackFilter([]*model.BlackListAuditPlanSQL{
		{
			FilterContent: "192.168.1.23",
			FilterType:    "IP",
		}, {
			FilterContent: "10.0.5.67",
			FilterType:    "IP",
		},
	})

	matchIps := []string{
		"10.0.5.67",
		"192.168.1.23",
	}
	for _, matchIp := range matchIps {
		if !filter.HasEndpointInBlackList([]string{matchIp}) {
			t.Error("Expected Ip to match blacklist")
		}
	}

	notMatchIps := []string{
		"172.16.254.89",
		"134.12.45.78",
		"50.67.89.12",
	}
	for _, notMatchIp := range notMatchIps {
		if filter.HasEndpointInBlackList([]string{notMatchIp}) {
			t.Error("Did not expect Ip to match blacklist")
		}
	}
}

func TestIsCidrInBlackList(t *testing.T) {
	filter := v1.ConvertToBlackFilter([]*model.BlackListAuditPlanSQL{
		{
			FilterContent: "192.168.0.0/24",
			FilterType:    "CIDR",
		}, {
			FilterContent: "10.100.0.0/16",
			FilterType:    "CIDR",
		},
	})

	matchIps := []string{
		"10.100.1.2",
		"10.100.25.45",
		"192.168.0.2",
		"192.168.0.45",
	}
	for _, matchIp := range matchIps {
		if !filter.HasEndpointInBlackList([]string{matchIp}) {
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
		if filter.HasEndpointInBlackList([]string{notMatchIp}) {
			t.Error("Did not expect CIDR to match blacklist")
		}
	}
}

func TestIsHostInBlackList(t *testing.T) {
	filter := v1.ConvertToBlackFilter([]*model.BlackListAuditPlanSQL{
		{
			FilterContent: "host",
			FilterType:    "HOST",
		}, {
			FilterContent: "some_site",
			FilterType:    "HOST",
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
		if !filter.HasEndpointInBlackList([]string{matchHost}) {
			t.Error("Expected HOST to match blacklist")
		}
	}

	notMatchHosts := []string{
		"other_site/home",
		"any_other_site/local",
	}
	for _, noMatchHost := range notMatchHosts {
		if filter.HasEndpointInBlackList([]string{noMatchHost}) {
			t.Error("Did not expect HOST to match blacklist")
		}
	}
}
