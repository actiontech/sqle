package oracle

import (
	"fmt"
	"testing"
)

func Test_Name(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDB(&DSN{
				Host:        "10.186.62.16",
				Port:        "1521",
				User:        "system",
				Password:    "COEdhd/umg4=1",
				ServiceName: "ORCLCDB",
			})
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	}

	_, err := NewDB(&DSN{
		Host:        "10.186.62.16",
		Port:        "1521",
		User:        "system",
		Password:    "COEdhd/umg4=1",
		ServiceName: "ORCLCDB",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
