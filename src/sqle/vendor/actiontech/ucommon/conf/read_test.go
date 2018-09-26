package conf

import (
	"fmt"
	"testing"
)

var (
	byteConfig = `[default1]
host = something.com
port = 443
active = true  
compression = off

[service-2]
port = 444
skip_lock
host = something.com

[service-1]#comment1
compression = on
`
)

func TestReadConfigBytes_Normal(t *testing.T) {
	t.Logf("conf is like before :\n%v\n", byteConfig)
	conf, err := ReadConfigBytes([]byte(byteConfig))
	assertNoError(t, err)
	options, err := conf.GetOptions("default1")
	assertNoError(t, err)
	if len(options) != 4 {
		for i, opt := range options {
			fmt.Printf("%v : %v,len=%v\n", i, opt, len(opt))
		}
		t.Fatalf("read conf uncorrect , expect num of default1.option is 4, result is %v", len(options))
	}

	host, err := conf.GetString("default1", "host")
	assertNoError(t, err)
	if "something.com" != host {
		t.Fatalf("expect default1.host=something.com, result default1.host=%v", host)
	}

	port, err := conf.GetInt("default1", "port")
	assertNoError(t, err)
	if 443 != port {
		t.Fatalf("expect default1.port=, result default1.port=%v", port)
	}

	active, err := conf.GetBool("default1", "active")
	assertNoError(t, err)
	if !active {
		t.Fatalf("expect default1.active=true, result default1.active=%v", active)
	}

	compression, err := conf.GetBool("default1", "compression")
	assertNoError(t, err)
	if compression {
		t.Fatalf("expect default1.compression=false, result default1.compression=%v", compression)
	}
}

// could have comments at section line
// eg: [setion1] #comments
func TestReadConfigBytes_SectionWithComment(t *testing.T) {
	t.Logf("conf is like before :\n%v\n", byteConfig)
	conf, err := ReadConfigBytes([]byte(byteConfig))
	assertNoError(t, err)
	if !conf.HasSection("service-1") {
		t.Fatalf("read conf error, section(service-1) with comment be ignored")
	}
}

func assertNoError(t *testing.T, err error) {
	if nil != err {
		t.Fatalf(err.Error())
	}
}

// if option don't have value,should give this option defautl value(true)
func TestReadConfigBytes_OptionWithoutValue(t *testing.T) {
	t.Logf("conf is like before :\n%v\n", byteConfig)
	conf, err := ReadConfigBytes([]byte(byteConfig))
	assertNoError(t, err)

	opt, err := conf.GetBool("service-2", "skip_lock")
	assertNoError(t, err)
	if !opt {
		t.Fatalf("expect service-2.skip_lock=true, result is %v", opt)
	}
	formalOpt, err := conf.GetInt("service-2", "port")
	assertNoError(t, err)
	if 444 != formalOpt {
		t.Fatalf("expect service-2.port=444, result is %v", formalOpt)
	}
}
