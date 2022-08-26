module github.com/actiontech/sqle

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/actiontech/mybatis-mapper-2-sql v0.1.1-0.20220728081924-8483c9ff0a98
	github.com/agiledragon/gomonkey v2.0.2+incompatible
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alibabacloud-go/darabonba-openapi v0.1.18
	github.com/alibabacloud-go/rds-20140815/v2 v2.1.0
	github.com/alibabacloud-go/tea v1.1.19
	github.com/alibabacloud-go/tea-utils v1.4.3
	github.com/clbanning/mxj/v2 v2.5.6 // indirect
	github.com/chanxuehong/util v0.0.0-20200304121633-ca8141845b13 // indirect
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
	github.com/go-ldap/ldap/v3 v3.4.1
	github.com/go-openapi/jsonreference v0.19.4 // indirect
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gogf/gf/v2 v2.1.2
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-plugin v1.4.2
	github.com/jackc/pgx/v4 v4.13.0
	github.com/jinzhu/gorm v1.9.15
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/jmoiron/sqlx v1.3.3
	github.com/labstack/echo/v4 v4.6.1
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
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.1
	github.com/swaggo/echo-swagger v1.0.0
	github.com/swaggo/swag v1.6.7
	github.com/timtadh/data-structures v0.5.3 // indirect
	github.com/timtadh/lexmachine v0.2.2
	github.com/ungerik/go-dry v0.0.0-20210209114055-a3e162a9e62e
	github.com/urfave/cli/v2 v2.1.1
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a
	google.golang.org/grpc v1.39.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/chanxuehong/wechat.v1 v1.0.0-20171118020122-aad7e298d1e7
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	vitess.io/vitess v0.12.0
)

replace (
	github.com/pingcap/log => github.com/pingcap/log v0.0.0-20191012051959-b742a5d432e9
	github.com/pingcap/parser => github.com/sjjian/parser v0.0.0-20220614062700-e3219e3d6833
	google.golang.org/grpc => google.golang.org/grpc v1.29.0
)
