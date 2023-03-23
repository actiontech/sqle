package license

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_GetAuditPlansByReq(t *testing.T) {
	sss, err := EncodeLicense(&LicensePermission{
		WorkDurationDay: 10,
		Version:         "999",
		UserCount:       20,
		NumberOfInstanceOfEachType: map[string]LimitOfType{
			"MySQL": {
				DBType: "MySQL",
				Count:  10,
			},
		},
	}, "EODbn3MOmOPYDsmsmdffo3Dbl8PcmsmYDeEOns5MEOqbo39uK4MCDB6PHN0bmOTbm3Tbq8U2q8Te68TcmSciAs5PEhiLm3iwn3iwn3Lwm84wm8Hw68jy75IMEg6di30doOqdoP4toOmuoP4eoPqeK4MCDB6PHN0rq8TrmeTbneUgn3Tv6OTsmSciAs5PEhiL6OiwnPDwq9iwo30wqOmwo4i=")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sss)
}

func TestCheckHardware(t *testing.T) {
	l := &License{
		LicenseContent: LicenseContent{
			HardwareSign: "test1",
		},
	}
	assert.NoError(t, l.CheckHardwareSignIsMatch("test1"))
	assert.Error(t, l.CheckHardwareSignIsMatch("test2"))
	assert.Error(t, l.CheckHardwareSignIsMatch(""))

	l = &License{
		LicenseContent: LicenseContent{
			HardwareSign: "test1",
			ClusterHardwareSign: map[string]string{
				"node3": "test3",
				"node4": "test4",
			},
		},
	}
	assert.NoError(t, l.CheckHardwareSignIsMatch("test1"))
	assert.NoError(t, l.CheckHardwareSignIsMatch("test3"))
	assert.NoError(t, l.CheckHardwareSignIsMatch("test4"))

	assert.Error(t, l.CheckHardwareSignIsMatch("test2"))
	assert.Error(t, l.CheckHardwareSignIsMatch("test5"))
	assert.Error(t, l.CheckHardwareSignIsMatch(""))
}

func TestCheckLicenseExol(t *testing.T) {
	l := &License{
		LicenseContent: LicenseContent{
			HardwareSign: "xxxxx1",
			Permission: LicensePermission{
				WorkDurationDay: 10,
			},
		},
		LicenseStatus: LicenseStatus{
			WorkDurationHour: 0,
		},
	}
	assert.NoError(t, l.CheckLicenseNotExpired())
	l.WorkDurationHour = 1
	assert.NoError(t, l.CheckLicenseNotExpired())
	l.WorkDurationHour = 239
	assert.NoError(t, l.CheckLicenseNotExpired())
	l.WorkDurationHour = 240
	assert.Error(t, l.CheckLicenseNotExpired())
	l.WorkDurationHour = 241
	assert.Error(t, l.CheckLicenseNotExpired())

	l.Permission.WorkDurationDay = 0
	l.WorkDurationHour = 0
	assert.Error(t, l.CheckLicenseNotExpired())
	l.WorkDurationHour = 1
	assert.Error(t, l.CheckLicenseNotExpired())
}

func TestCheckCreateUser(t *testing.T) {
	l := &License{
		LicenseContent: LicenseContent{
			Permission: LicensePermission{
				UserCount: 10,
			},
		},
	}

	assert.NoError(t, l.CheckCanCreateUser(0))
	assert.NoError(t, l.CheckCanCreateUser(9))

	assert.Error(t, l.CheckCanCreateUser(10))
	assert.Error(t, l.CheckCanCreateUser(11))
	assert.Error(t, l.CheckCanCreateUser(12))
	assert.Error(t, l.CheckCanCreateUser(1000))

	l = &License{
		LicenseContent: LicenseContent{
			Permission: LicensePermission{
				UserCount: 0,
			},
		},
	}
	assert.Error(t, l.CheckCanCreateUser(0))
	assert.Error(t, l.CheckCanCreateUser(1))
	assert.Error(t, l.CheckCanCreateUser(1000))
}

func TestCheckCreateInstance(t *testing.T) {
	//
	type testCase struct {
		total  LimitOfEachType
		usage  LimitOfEachType
		dbType string
		pass   bool
	}
	cases := []testCase{}

	// 无实例容量
	cases = append(cases, []testCase{
		{
			total: LimitOfEachType{},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 0}},
			dbType: "MySQL",
			pass:   false,
		},
		{
			total: LimitOfEachType{},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 1}},
			dbType: "MySQL",
			pass:   false,
		},
	}...)

	// 单数据库类型
	cases = append(cases, []testCase{
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
			},
			usage:  LimitOfEachType{},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 0}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 1}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 9}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10}},
			dbType: "MySQL",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 11}},
			dbType: "MySQL",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 100}},
			dbType: "MySQL",
			pass:   false,
		},
	}...)

	// 单数据库类型带custom
	cases = append(cases, []testCase{
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 0}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 1}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 9}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 11}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 19}},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 20}},
			dbType: "MySQL",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 21}},
			dbType: "MySQL",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 100}},
			dbType: "MySQL",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10}},
			dbType: "TiDB",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10}},
			dbType: "TiDB",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 19}},
			dbType: "TiDB",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 19}},
			dbType: "TiDB",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 20}},
			dbType: "TiDB",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL":  LimitOfType{DBType: "MySQL", Count: 10},
				"custom": LimitOfType{DBType: "custom", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 21}},
			dbType: "TiDB",
			pass:   false,
		},
	}...)

	// 多数据库类型
	cases = append(cases, []testCase{
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 0},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 9},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			dbType: "MySQL",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			dbType: "MySQL",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 11},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			dbType: "MySQL",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 0},
			},
			dbType: "TiDB",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 9},
			},
			dbType: "TiDB",
			pass:   true,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			dbType: "TiDB",
			pass:   false,
		},
		{
			total: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 10},
			},
			usage: LimitOfEachType{
				"MySQL": LimitOfType{DBType: "MySQL", Count: 10},
				"TiDB":  LimitOfType{DBType: "TiDB", Count: 11},
			},
			dbType: "TiDB",
			pass:   false,
		},
	}...)
	for _, c := range cases {
		l := &License{
			LicenseContent: LicenseContent{
				Permission: LicensePermission{
					NumberOfInstanceOfEachType: c.total,
				},
			},
		}
		err := l.CheckCanCreateInstance(c.dbType, c.usage)
		if c.pass {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}

}
