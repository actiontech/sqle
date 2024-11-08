package oracle

import (
	"fmt"
	"testing"
)

func Test_Name(t *testing.T) {
	_, err := NewDB(&DSN{
		Host:        "10.186.62.16",
		Port:        "1521",
		User:        "sys",
		Password:    "COEdhd/umg4=1",
		ServiceName: "ORCLCDB",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
