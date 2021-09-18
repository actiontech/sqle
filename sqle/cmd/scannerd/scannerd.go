package main

import (
	"fmt"
	"os"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/cmd"

	"github.com/fatih/color"
)

func main() {
	var code int
	err := cmd.Execute()
	if err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		code = 1
	}

	color.Unset()
	if code != 0 {
		os.Exit(code)
	}
}
