# GoHive
[![Build Status](https://travis-ci.com/beltran/gohive.svg?branch=master)](https://travis-ci.com/beltran/gohive) [![Coverage Status](https://coveralls.io/repos/github/beltran/gohive/badge.svg?branch=master)](https://coveralls.io/github/beltran/gohive?branch=master)


GoHive is a driver for Hive and the [Spark Distributed SQL Engine](https://spark.apache.org/docs/latest/sql-distributed-sql-engine.html) in go that supports connection mechanisms KERBEROS(Gssapi Sasl), NONE(Plain Sasl), LDAP, CUSTOM and NOSASL, both for binary and http transport, with and without SSL. The kerberos mechanism will pick a different authentication level depending on `hive.server2.thrift.sasl.qop`.

## Installation
Gohive can be installed with:
```
go get github.com/beltran/gohive
```

To add kerberos support gohive requires header files to build against the GSSAPI C library. They can be installed with:
- Ubuntu: `sudo apt-get install libkrb5-dev`
- MacOS: `brew install homebrew/dupes/heimdal --without-x11`
- Debian: `yum install -y krb5-devel`

Then:
```
go get -tags kerberos github.com/beltran/gohive
```

## Quickstart

```go
    connection, errConn := gohive.Connect("hs2.example.com", 10000, "KERBEROS", configuration)
    if errConn != nil {
        log.Fatal(errConn)
    }
    cursor := connection.Cursor()

    cursor.Exec(ctx, "INSERT INTO myTable VALUES(1, '1'), (2, '2'), (3, '3'), (4, '4')")
    if cursor.Err != nil {
        log.Fatal(cursor.Err)
    }

    cursor.Exec(ctx, "SELECT * FROM myTable")
    if cursor.Err != nil {
        log.Fatal(cursor.Err)
    }

    var i int32
    var s string
    for cursor.HasMore(ctx) {
        cursor.FetchOne(ctx, &i, &s)
        if cursor.Err != nil {
            log.Fatal(cursor.Err)
        }
        log.Println(i, s)
    }

    cursor.Close()
    connection.Close()
```

`cursor.HasMore` may query hive for more rows if not all of them have been received. Once the row is
read is discarded from memory so as long as the fetch size is not too big there's no limit to how much
data can be queried.

## Supported connections
### Connect with Sasl kerberos:
``` go
configuration := NewConnectConfiguration()
configuration.Service = "hive"
// Previously kinit should have done: kinit -kt ./secret.keytab hive/hs2.example.com@EXAMPLE.COM
connection, errConn := Connect("hs2.example.com", 10000, "KERBEROS", configuration)
```
This implies setting in hive-site.xml:
- `hive.server2.authentication = KERBEROS`
- `hive.server2.authentication.kerberos.principal = hive/_HOST@EXAMPLE.COM`
- `hive.server2.authentication.kerberos.keytab = path/to/keytab.keytab`

### Connnect using Plain Sasl:
``` go
configuration := NewConnectConfiguration()
// If it's not set it will be picked up from the logged user
configuration.Username = "myUsername"
// This may not be necessary
configuration.Password = "myPassword"
connection, errConn := Connect("hs2.example.com", 10000, "NONE", configuration)
```
This implies setting in hive-site.xml:

- `hive.server2.authentication = NONE`

### Connnect using No Sasl:
``` go
connection, errConn := Connect("hs2.example.com", 10000, "NOSASL", NewConnectConfiguration())
```
This implies setting in hive-site.xml:

- `hive.server2.authentication = NOSASL`

### Connect using Http transport mode
Binary transport mode is supported for this three options(PLAIN, KERBEROS and NOSASL). Http transport is supported for PLAIN and KERBEROS:
``` go
configuration := NewConnectConfiguration()
configuration.HttpPath = "cliservice" // this is the default path in hive configuration.
configuration.TransportMode = "http"
configuration.Service = "hive"

connection, errConn := Connect("hs2.example.com", 10000, "KERBEROS", configuration)
```
This implies setting in hive-site.xml:

- `hive.server2.authentication = KERBEROS`, or `NONE`
- `hive.server2.transport.mode = http`
- `hive.server2.thrift.http.port = 10001`

## Zookeeper
A connection can be made using zookeeper:

```go
connection, errConn := ConnectZookeeper("zk1.example.com:2181,zk2.example.com:2181", "NONE", configuration)
```
The last two parameters determine how the connection to hive will be made once the hive hosts are retrieved from zookeeper.

## NULL values
For example if a `NULL` value is in a row, the following operations would put `0` into `i`:
```
var i int32
cursor.FetchOne(context.Background(), &i)
```
To differentiate between these two values (`NULL` and `0`) the following will set `i` to `nil` or `*i` to `0`:
```
var i *int32 = new(int32)
cursor.FetchOne(context.Background(), &i)
```
which will produce the same result as:
```
var i *int32
cursor.FetchOne(context.Background(), &i)
```
Alternatively, using the rowmap API, `m := cursor.RowMap(context.Background())`,
 `m` would be `map[string]interface{}{"table_name.column_name": nil}` for a `NULL` value. It will return a map
where the keys are `table_name.column_name`. This works fine with hive but using [Spark Thirft SQL server](https://spark.apache.org/docs/latest/sql-distributed-sql-engine.html) `table_name` is not present and the keys are `column_name` and it can [lead to problems](https://github.com/beltran/gohive/issues/120) if two tables have the same column name so the `FetchOne` API should be used in this case.

## Running tests
Tests can be run with:
```
./scripts/integration
```
This uses [dhive](https://github.com/beltran/dhive) and it will start two docker instances with hive and kerberos. `kinit`, `klist`, `kdestroy` have to be installed locally. `hs2.example.com` will have to be an alias for 127.0.0.1 in `/etc/hosts`. The krb5 configuration file should be created with `bash scripts/create_krbconf.sh`. Overall the [steps used in the travis CI](https://github.com/beltran/gohive/blob/ec69b5601829296a56ca0558693ed30c11180a94/.travis.yml#L24-L46) can be followed.
