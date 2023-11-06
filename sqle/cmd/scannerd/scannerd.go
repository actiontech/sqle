package main

import (
	"context"
	"fmt"
	"os"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/cmd"

	"github.com/fatih/color"
)

var version string

func main() {
	var code int
	ctx := context.WithValue(context.Background(), cmd.VersionKey, version)

	err := cmd.Execute(ctx)
	if err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		code = 1
	}

	color.Unset()
	if code != 0 {
		os.Exit(code)
	}
}
