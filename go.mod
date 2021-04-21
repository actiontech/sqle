module actiontech.cloud/universe/sqle/v4

go 1.14

require (
	actiontech.cloud/universe/ucommon/v4 v4.2101.1-0.20210204023918-a44296eb4ef0
	actiontech.cloud/universe/ucore-common/v4 v4.0.0-20210222085528-4892e77d2aac
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cznic/golex v0.0.0-20181122101858-9c343928389c // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548
	github.com/cznic/parser v0.0.0-20181122101858-d773202d5b1f
	github.com/cznic/sortutil v0.0.0-20150617083342-4c7342852e65
	github.com/cznic/strutil v0.0.0-20181122101858-275e90344537
	github.com/cznic/y v0.0.0-20181122101901-b05e8c2e8d7b
	github.com/denisenkom/go-mssqldb v0.0.0-20200620013148-b91950f658ec
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.3.5
	github.com/jinzhu/gorm v1.9.15
	github.com/jmoiron/sqlx v1.3.1
	github.com/labstack/echo/v4 v4.0.0
	github.com/labstack/gommon v0.2.8
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pingcap/parser v3.0.12+incompatible
	github.com/pingcap/tidb v0.0.0-20200312110807-8c4696b3f340 // v3.0.12
	github.com/sirupsen/logrus v1.2.0
	github.com/sjjian/mybatis-mapper-2-sql v0.1.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/cobra v0.0.7
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/swaggo/echo-swagger v1.0.0
	github.com/swaggo/swag v1.6.7
	github.com/urfave/cli/v2 v2.1.1
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	google.golang.org/grpc v1.28.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/ini.v1 v1.57.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/jinzhu/gorm => actiontech.cloud/universe/gorm v0.0.0-20190520085104-6d6ea8fa4ec5
