package log

import (
	"fmt"
	"strings"
	"testing"
)

func TestFilterPassword(t *testing.T) {
	line := "abc"
	result := defaultFilterPassword(line)
	if "abc" != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = "password="
	result = defaultFilterPassword(line)
	if `password=******` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = "password=1234qwer"
	result = defaultFilterPassword(line)
	if `password=******` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = "password=1234qwer abcd"
	result = defaultFilterPassword(line)
	if `password=****** abcd` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = "password:1234qwer"
	result = defaultFilterPassword(line)
	if `password:******` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = "identified by 1234qwer"
	result = defaultFilterPassword(line)
	if `identified by ******` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = "-P 123456 -F"
	result = defaultFilterPassword(line)
	if `-P ****** -F` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = "-P 123456 -f"
	result = defaultFilterPassword(line)
	if `-P ****** -f` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = "-P 123456"
	result = defaultFilterPassword(line)
	if `-P ******` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = `"password":123 abcd`
	result = defaultFilterPassword(line)
	if `"password":****** abcd` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = `"password":123,abcd`
	result = defaultFilterPassword(line)
	if `"password":******,abcd` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = `"replPassword":"123456"`
	result = defaultFilterPassword(line)
	if `"replPassword":"******"` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = `"replPassword":123456`
	result = defaultFilterPassword(line)
	if `"replPassword":******` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = `"backupPassword":"123456"`
	result = defaultFilterPassword(line)
	if `"backupPassword":"******"` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = `Password(123456)`
	result = defaultFilterPassword(line)
	if `Password(******)` != result {
		t.Errorf("'%v' not filter, result is  %v", line, result)
	}

	line = `request for mysql-oz91mc, {"common":{"instID":"mysql-oz91mc","reason":"rpc call promote mysql-oz91mc","consulIndex":885},"logApplierUser":"uguard_apply","logApplierPassword":"123456","logApplierPassword":"123456","scsiLevel":5,"sip":"172.17.5.100","ignoreFetchError":true,"srcDbMeta":{"user":"uguard_op","password":"123456","host":"127.0.0.1","port":"3306","execSQLTimeout":10},"srcScsi":{},"srcAddr":"172.17.5.3","srcAgentPort":"5710","srcLogbin":"/opt/mysql/binlog/3306/mysql-bin","localDbMeta":{"user":"uguard_op","password":"123456","host":"127.0.0.1","port":"3306","execSQLTimeout":10},"localScsi":{},"runUser":"gxna","mycnfPath":"/opt/mysql/etc/3306/my.cnf"}`
	result = defaultFilterPassword(line)
	if `request for mysql-oz91mc, {"common":{"instID":"mysql-oz91mc","reason":"rpc call promote mysql-oz91mc","consulIndex":885},"logApplierUser":"uguard_apply","logApplierPassword":"******","logApplierPassword":"******","scsiLevel":5,"sip":"172.17.5.100","ignoreFetchError":true,"srcDbMeta":{"user":"uguard_op","password":"******","host":"127.0.0.1","port":"3306","execSQLTimeout":10},"srcScsi":{},"srcAddr":"172.17.5.3","srcAgentPort":"5710","srcLogbin":"/opt/mysql/binlog/3306/mysql-bin","localDbMeta":{"user":"uguard_op","password":"******","host":"127.0.0.1","port":"3306","execSQLTimeout":10},"localScsi":{},"runUser":"gxna","mycnfPath":"/opt/mysql/etc/3306/my.cnf"}` != result {
		t.Errorf("'%v' not filter,\n result is  %v", line, result)
	}

	line = `<?xml version="1.0"?>
<!DOCTYPE mycat:schema SYSTEM "schema.dtd">
<mycat:schema xmlns:mycat="http://io.mycat/">
    <dataHost name="jiedian1" dbType="mysql" dbDriver="native" maxCon="100" minCon="10" balance="0" switchType="-1" slaveThreshold="-1" tempReadHostAvailable="0">
        <heartbeat>show slave status</heartbeat>
        <writeHost host="mysql-agogo2" url="172.20.30.11:3306" user="root" password="123">
            <readHost host="mysql-wtp047" url="172.20.30.13:3306" user="root" password="123"></readHost>
        </writeHost>
    </dataHost>
</mycat:schema>`
	result = defaultFilterPassword(line)
	if `<?xml version="1.0"?>
<!DOCTYPE mycat:schema SYSTEM "schema.dtd">
<mycat:schema xmlns:mycat="http://io.mycat/">
    <dataHost name="jiedian1" dbType="mysql" dbDriver="native" maxCon="100" minCon="10" balance="0" switchType="-1" slaveThreshold="-1" tempReadHostAvailable="0">
        <heartbeat>show slave status</heartbeat>
        <writeHost host="mysql-agogo2" url="172.20.30.11:3306" user="root" password="******">
            <readHost host="mysql-wtp047" url="172.20.30.13:3306" user="root" password="******"></readHost>
        </writeHost>
    </dataHost>
</mycat:schema>` != result {
		t.Errorf("%v\n=====================\n%v\n", line, result)
	}
}

func TestDone(t *testing.T) {
	instance = NewTestLogger()
	instance.setLevelAbility(detail, true)
	st := NewStage()
	st.Enter("test")
	Write(st).Detail("detail").Brief("brief").Done()
	Write(st).Detail("detail2").Done()
	Detail(st, "only detail")
	Detail(st, "identified by 1234qwer")
}

func TestBuilder(t *testing.T) {
	testLogger := NewTestLogger()
	instance = testLogger
	instance.setLevelAbility(detail, true)
	stage := NewStage()
	stage.Enter("TestBuilder")

	instance.setLevelAbility(brief, false)
	instance.setLevelAbility(key, true)
	Write(stage).UserError("user error").Brief("brief error").Detail("detail error").Done()
	UserInfo(stage, "userinfo")
	fmt.Println("userinfo")
	userlog := strings.Join(testLogger.logContents[user], "")
	fmt.Println(userlog)

	Key(stage, "key log")
	fmt.Println("key log")
	keylog := strings.Join(testLogger.logContents[key], "")
	fmt.Println(keylog)

	Brief(stage, "brief log")
	fmt.Println("brief log")
	brieflog := strings.Join(testLogger.logContents[brief], "")
	fmt.Println(brieflog)

	Detail(stage, "detail log")
	fmt.Println("detail log")
	detaillog := strings.Join(testLogger.logContents[detail], "")
	fmt.Println(detaillog)

	// Detail(stage, "detail log")
	fmt.Println("trace log")
}
