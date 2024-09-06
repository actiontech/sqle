# pg_query_go [![GoDoc](https://godoc.org/github.com/pganalyze/pg_query_go/v4?status.svg)](https://godoc.org/github.com/pganalyze/pg_query_go/v4)

Go version of https://github.com/pganalyze/pg_query

This Go library and its cgo extension use the actual PostgreSQL server source to parse SQL queries and return the internal PostgreSQL parse tree.

You can find further background to why a query's parse tree is useful here: https://pganalyze.com/blog/parse-postgresql-queries-in-ruby.html


## Installation

```
go get github.com/pganalyze/pg_query_go/v4@latest
```

Due to compiling parts of PostgreSQL, the first time you build against this library it will take a bit longer.

Expect up to 3 minutes. You can use `go build -x` to see the progress.

## Usage with Go modules

When integrating this library using Go modules, and using a vendor/ directory,
you will need to explicitly copy over some of the C build files, since Go does
not copy files in subfolders without .go files whilst vendoring.

The best way to do so is to use [modvendor](https://github.com/goware/modvendor),
and vendor your modules like this:

```
go mod vendor
go get -u github.com/goware/modvendor
modvendor -copy="**/*.c **/*.h **/*.proto" -v
```

## Usage

### Parsing a query into JSON

Put the following in a new Go package, after having installed pg_query as above:

```go
package main

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v4"
)

func main() {
	tree, err := pg_query.ParseToJSON("SELECT 1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", tree)
}
```

Running will output the query's parse tree as JSON:

```json
{"version":150001,"stmts":[{"stmt":{"SelectStmt":{"targetList":[{"ResTarget":{"val":{"A_Const":{"ival":{"ival":1},"location":7}},"location":7}}],"limitOption":"LIMIT_OPTION_DEFAULT","op":"SETOP_NONE"}}}]}
```

### Parsing a query into Go structs

When working with the query information inside Go its recommended you use the `Parse()` method which returns Go structs:

```go
package main

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v4"
)

func main() {
	result, err := pg_query.Parse("SELECT 42")
	if err != nil {
		panic(err)
	}

	// This will output "42"
	fmt.Printf("%d\n", result.Stmts[0].Stmt.GetSelectStmt().GetTargetList()[0].GetResTarget().GetVal().GetAConst().GetIval().Ival)
}
```

You can find all the node types in the `pg_query.pb.go` Protobuf definition.

### Deparsing a parse tree back into a SQL statement

In order to go back from a parse tree to a SQL statement, you can use the deparsing functionality:

```go
package main

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v4"
)

func main() {
	result, err := pg_query.Parse("SELECT 42")
	if err != nil {
		panic(err)
	}

	result.Stmts[0].Stmt.GetSelectStmt().GetTargetList()[0].GetResTarget().Val = pg_query.MakeAConstStrNode("Hello World", -1)

	stmt, err := pg_query.Deparse(result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", stmt)
}
```

This will output the following:

```
SELECT 'Hello World'
```

Note that it is currently not recommended to pass unsanitized input to the deparser, as it may lead to crashes.

### Parsing a PL/pgSQL function into JSON (Experimental)

Put the following in a new Go package, after having installed pg_query as above:

```go
package main

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v4"
)

func main() {
	tree, err := pg_query.ParsePlPgSqlToJSON(
		`CREATE OR REPLACE FUNCTION cs_fmt_browser_version(v_name varchar, v_version varchar)
  			RETURNS varchar AS $$
  			BEGIN
  			    IF v_version IS NULL THEN
  			        RETURN v_name;
  			    END IF;
  			    RETURN v_name || '/' || v_version;
  			END;
  			$$ LANGUAGE plpgsql;`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", tree)
}
```

Running will output the functions's parse tree as JSON:

```json
$ go run main.go
[
{"PLpgSQL_function":{"datums":[{"PLpgSQL_var":{"refname":"v_name","datatype":{"PLpgSQL_type":{"typname":"UNKNOWN"}}}},{"PLpgSQL_var":{"refname":"v_version","datatype":{"PLpgSQL_type":{"typname":"UNKNOWN"}}}},{"PLpgSQL_var":{"refname":"found","datatype":{"PLpgSQL_type":{"typname":"UNKNOWN"}}}}],"action":{"PLpgSQL_stmt_block":{"lineno":2,"body":[{"PLpgSQL_stmt_if":{"lineno":3,"cond":{"PLpgSQL_expr":{"query":"v_version IS NULL"}},"then_body":[{"PLpgSQL_stmt_return":{"lineno":4,"expr":{"PLpgSQL_expr":{"query":"v_name"}}}}]}},{"PLpgSQL_stmt_return":{"lineno":6,"expr":{"PLpgSQL_expr":{"query":"v_name || '/' || v_version"}}}}]}}}}
]
```

## Benchmarks

```
$ make benchmark
go build -a
go test -test.bench=. -test.run=XXX -test.benchtime 10s -test.benchmem -test.cpu=4
goos: darwin
goarch: arm64
pkg: github.com/pganalyze/pg_query_go/v4
BenchmarkParseSelect1-4                  	 3230398	      3656 ns/op	    1104 B/op	      20 allocs/op
BenchmarkParseSelect2-4                  	  927363	     12739 ns/op	    2896 B/op	      59 allocs/op
BenchmarkParseCreateTable-4              	  399819	     30080 ns/op	    8432 B/op	     151 allocs/op
BenchmarkParseSelect1Parallel-4          	10951803	      1094 ns/op	    1104 B/op	      20 allocs/op
BenchmarkParseSelect2Parallel-4          	 3255471	      3675 ns/op	    2896 B/op	      59 allocs/op
BenchmarkParseCreateTableParallel-4      	 1341716	      8919 ns/op	    8432 B/op	     151 allocs/op
BenchmarkRawParseSelect1-4               	 4275111	      2795 ns/op	     192 B/op	       5 allocs/op
BenchmarkRawParseSelect2-4               	 1252704	      9534 ns/op	     352 B/op	       5 allocs/op
BenchmarkRawParseCreateTable-4           	  503385	     23168 ns/op	    1120 B/op	       5 allocs/op
BenchmarkRawParseSelect1Parallel-4       	15446528	       780.1 ns/op	     192 B/op	       5 allocs/op
BenchmarkRawParseSelect2Parallel-4       	 4638837	      2595 ns/op	     352 B/op	       5 allocs/op
BenchmarkRawParseCreateTableParallel-4   	 1932316	      6197 ns/op	    1120 B/op	       5 allocs/op
BenchmarkFingerprintSelect1-4            	 6583965	      1795 ns/op	     112 B/op	       4 allocs/op
BenchmarkFingerprintSelect2-4            	 2962663	      4015 ns/op	     112 B/op	       4 allocs/op
BenchmarkFingerprintCreateTable-4        	 1796041	      6831 ns/op	     112 B/op	       4 allocs/op
BenchmarkNormalizeSelect1-4              	10073278	      1171 ns/op	      72 B/op	       4 allocs/op
BenchmarkNormalizeSelect2-4              	 6029834	      1932 ns/op	     104 B/op	       4 allocs/op
BenchmarkNormalizeCreateTable-4          	 4703816	      2490 ns/op	     184 B/op	       4 allocs/op
PASS
ok  	github.com/pganalyze/pg_query_go/v4	273.449s
```

Note that allocation counts exclude the cgo portion, so they are higher than shown here.

See `benchmark_test.go` for details on the benchmarks.


## Authors

- [Lukas Fittl](mailto:lukas@fittl.com)


## License

Copyright (c) 2015, Lukas Fittl <lukas@fittl.com><br>
pg_query_go is licensed under the 3-clause BSD license, see LICENSE file for details.

This project includes code derived from the [PostgreSQL project](http://www.postgresql.org/),
see LICENSE.POSTGRESQL for details.
