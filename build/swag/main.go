package main

import (
	"fmt"
	"log"
	"os"

	"github.com/swaggo/swag"
	"github.com/swaggo/swag/gen"
	"github.com/urfave/cli/v2"
)

const (
	searchDirFlag        = "dir"
	excludeFlag          = "exclude"
	generalInfoFlag      = "generalInfo"
	propertyStrategyFlag = "propertyStrategy"
	outputFlag           = "output"
	parseVendorFlag      = "parseVendor"
	parseDependencyFlag  = "parseDependency"
	markdownFilesFlag    = "markdownFiles"
	parseInternal        = "parseInternal"
	generatedTimeFlag    = "generatedTime"
)

var initFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    generalInfoFlag,
		Aliases: []string{"g"},
		Value:   "main.go",
		Usage:   "Go file path in which 'swagger general API Info' is written",
	},
	&cli.StringFlag{
		Name:    searchDirFlag,
		Aliases: []string{"d"},
		Value:   "./",
		Usage:   "Directory you want to parse",
	},
	&cli.StringFlag{
		Name:  excludeFlag,
		Usage: "exclude directories and files when searching, comma separated",
	},
	&cli.StringFlag{
		Name:    propertyStrategyFlag,
		Aliases: []string{"p"},
		Value:   "camelcase",
		Usage:   "Property Naming Strategy like snakecase,camelcase,pascalcase",
	},
	&cli.StringFlag{
		Name:    outputFlag,
		Aliases: []string{"o"},
		Value:   "./docs",
		Usage:   "Output directory for all the generated files(swagger.json, swagger.yaml and doc.go)",
	},
	&cli.BoolFlag{
		Name:  parseVendorFlag,
		Usage: "Parse go files in 'vendor' folder, disabled by default",
	},
	&cli.BoolFlag{
		Name:  parseDependencyFlag,
		Usage: "ParseDependencies whether swag should be parse outside dependency folder: 0 none, 1 models, 2 operations, 3 all. 0 by default",
	},
	&cli.StringFlag{
		Name:    markdownFilesFlag,
		Aliases: []string{"md"},
		Value:   "",
		Usage:   "Parse folder containing markdown files to use as description, disabled by default",
	},
	&cli.BoolFlag{
		Name:  "parseInternal",
		Usage: "Parse go files in internal packages, disabled by default",
	},
	&cli.BoolFlag{
		Name:  "generatedTime",
		Usage: "Generate timestamp at the top of docs.go, true by default",
	},
}

func initAction(c *cli.Context) error {
	strategy := c.String(propertyStrategyFlag)

	switch strategy {
	case swag.CamelCase, swag.SnakeCase, swag.PascalCase:
	default:
		return fmt.Errorf("not supported %s propertyStrategy", strategy)
	}

	return gen.New().Build(&gen.Config{
		SearchDir:          c.String(searchDirFlag),
		Excludes:           c.String(excludeFlag),
		MainAPIFile:        c.String(generalInfoFlag),
		PropNamingStrategy: strategy,
		OutputDir:          c.String(outputFlag),
		ParseVendor:        c.Bool(parseVendorFlag),
		ParseDependency:    c.Bool(parseDependencyFlag),
		MarkdownFilesDir:   c.String(markdownFilesFlag),
		ParseInternal:      c.Bool(parseInternal),
		GeneratedTime:      c.Bool(generatedTimeFlag),
	})
}

func main() {
	app := cli.NewApp()
	app.Version = swag.Version
	app.Usage = "Automatically generate RESTful API documentation with Swagger 2.0 for Go."
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Create docs.go",
			Action:  initAction,
			Flags:   initFlags,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
