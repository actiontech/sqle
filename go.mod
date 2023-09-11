module github.com/actiontech/sqle

go 1.16

require (
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/99designs/gqlgen v0.17.20
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/actiontech/dms v0.0.0-20230907030052-c3b5a7505949
	github.com/actiontech/mybatis-mapper-2-sql v0.3.0
	github.com/agiledragon/gomonkey v2.0.2+incompatible
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alibabacloud-go/darabonba-openapi v0.1.18
	github.com/alibabacloud-go/darabonba-openapi/v2 v2.0.2
	github.com/alibabacloud-go/dingtalk v1.4.88
	github.com/alibabacloud-go/rds-20140815/v2 v2.1.0
	github.com/alibabacloud-go/tea v1.1.19
	github.com/alibabacloud-go/tea-utils v1.4.3
	github.com/alibabacloud-go/tea-utils/v2 v2.0.1
	github.com/baidubce/bce-sdk-go v0.9.151
	github.com/bwmarrin/snowflake v0.3.0
	github.com/clbanning/mxj/v2 v2.5.6 // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548
	github.com/cznic/parser v0.0.0-20181122101858-d773202d5b1f
	github.com/cznic/sortutil v0.0.0-20181122101858-f5f958428db8
	github.com/cznic/strutil v0.0.0-20181122101858-275e90344537
	github.com/cznic/y v0.0.0-20181122101901-b05e8c2e8d7b
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/freeport v0.0.0-20150612182905-d4adf43b75b9 // indirect
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/fatih/color v1.13.0
	github.com/github/gh-ost v1.1.3-0.20210727153850-e484824bbd68
	github.com/go-ini/ini v1.63.2
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.14.1
	github.com/go-sql-driver/mysql v1.7.0
	github.com/gogf/gf/v2 v2.1.2
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-plugin v1.4.2
	github.com/jackc/pgx/v4 v4.13.0
	github.com/jinzhu/gorm v1.9.15
	github.com/jmoiron/sqlx v1.3.3
	github.com/labstack/echo/v4 v4.10.2
	github.com/larksuite/oapi-sdk-go/v3 v3.0.23
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/moby/sys/mountinfo v0.6.2
	github.com/nxadm/tail v1.4.8
	github.com/openark/golib v0.0.0-20210531070646-355f37940af8
	github.com/percona/go-mysql v0.0.0-20210427141028-73d29c6da78c
	github.com/percona/pmm-agent v2.15.1+incompatible
	github.com/pingcap/parser v3.0.12+incompatible
	github.com/pingcap/tidb v1.1.0-beta.0.20200630082100-328b6d0a955c
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/sijms/go-ora/v2 v2.2.15
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.8.2
	github.com/swaggo/echo-swagger v1.0.0
	github.com/swaggo/swag v1.6.7
	github.com/timtadh/data-structures v0.5.3 // indirect
	github.com/timtadh/lexmachine v0.2.2
	github.com/ungerik/go-dry v0.0.0-20210209114055-a3e162a9e62e
	github.com/urfave/cli/v2 v2.8.1
	github.com/vektah/gqlparser/v2 v2.5.1
	golang.org/x/net v0.11.0
	golang.org/x/sys v0.9.0
	google.golang.org/grpc v1.50.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	vitess.io/vitess v0.12.0
)

replace (
	cloud.google.com/go/compute/metadata => cloud.google.com/go/compute/metadata v0.1.0
	github.com/labstack/echo/v4 => github.com/labstack/echo/v4 v4.6.1
	github.com/percona/pmm-agent => github.com/taolx0/pmm-agent v0.0.0-20230614092412-936a5cff4635
	github.com/pingcap/log => github.com/pingcap/log v0.0.0-20191012051959-b742a5d432e9
	github.com/pingcap/parser => github.com/sjjian/parser v0.0.0-20220614062700-e3219e3d6833
	golang.org/x/net => golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/sys => golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
	google.golang.org/grpc => google.golang.org/grpc v1.29.0
)
