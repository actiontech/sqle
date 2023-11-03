module github.com/actiontech/sqle

go 1.19

require (
	github.com/99designs/gqlgen v0.17.20
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/actiontech/java-sql-extractor v0.0.0-20231019025645-418760361a3b
	github.com/actiontech/mybatis-mapper-2-sql v0.4.0
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
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548
	github.com/cznic/parser v0.0.0-20181122101858-d773202d5b1f
	github.com/cznic/sortutil v0.0.0-20181122101858-f5f958428db8
	github.com/cznic/strutil v0.0.0-20181122101858-275e90344537
	github.com/cznic/y v0.0.0-20181122101901-b05e8c2e8d7b
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/facebookgo/grace v0.0.0-20180706040059-75cf19382434
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/fatih/color v1.13.0
	github.com/github/gh-ost v1.1.3-0.20210727153850-e484824bbd68
	github.com/go-git/go-git/v5 v5.9.0
	github.com/go-ini/ini v1.63.2
	github.com/go-ldap/ldap/v3 v3.4.1
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.5 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gogf/gf/v2 v2.1.2
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.3
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-plugin v1.4.2
	github.com/jackc/pgx/v4 v4.13.0
	github.com/jinzhu/gorm v1.9.15
	github.com/jmoiron/sqlx v1.3.3
	github.com/labstack/echo/v4 v4.6.1
	github.com/larksuite/oapi-sdk-go/v3 v3.0.18
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/moby/sys/mountinfo v0.6.2
	github.com/nxadm/tail v1.4.11
	github.com/openark/golib v0.0.0-20210531070646-355f37940af8
	github.com/percona/go-mysql v0.0.0-20210427141028-73d29c6da78c
	github.com/percona/pmm-agent v2.15.1+incompatible
	github.com/pingcap/parser v3.0.12+incompatible
	github.com/pingcap/tidb v1.1.0-beta.0.20200630082100-328b6d0a955c
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/sijms/go-ora/v2 v2.2.15
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.1
	github.com/swaggo/echo-swagger v1.0.0
	github.com/swaggo/swag v1.6.7
	github.com/timtadh/data-structures v0.5.3 // indirect
	github.com/timtadh/lexmachine v0.2.2
	github.com/ungerik/go-dry v0.0.0-20210209114055-a3e162a9e62e
	github.com/urfave/cli/v2 v2.8.1
	github.com/vektah/gqlparser/v2 v2.5.1
	golang.org/x/net v0.15.0
	google.golang.org/grpc v1.46.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/chanxuehong/wechat.v1 v1.0.0-20171118020122-aad7e298d1e7
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	vitess.io/vitess v0.12.0
)

require (
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/sys v0.12.0
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20200615164410-66371956d46c // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230828082145-3c4c8a2d2371 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/acomagu/bufpipe v1.0.4 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/alibabacloud-go/alibabacloud-gateway-spi v0.0.4 // indirect
	github.com/alibabacloud-go/debug v0.0.0-20190504072949-9472017b5c68 // indirect
	github.com/alibabacloud-go/endpoint-util v1.1.0 // indirect
	github.com/alibabacloud-go/openapi-util v0.1.0 // indirect
	github.com/alibabacloud-go/tea-xml v1.1.2 // indirect
	github.com/aliyun/credentials-go v1.1.2 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/chanxuehong/util v0.0.0-20200304121633-ca8141845b13 // indirect
	github.com/clbanning/mxj/v2 v2.5.6 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/cznic/golex v0.0.0-20181122101858-9c343928389c // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/freeport v0.0.0-20150612182905-d4adf43b75b9 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/go-mysql-org/go-mysql v1.3.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mitchellh/go-testing-interface v1.14.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pingcap/errors v0.11.5-0.20201126102027-b0a155152ca3 // indirect
	github.com/pingcap/log v0.0.0-20210317133921-96f4fcab92a4 // indirect
	github.com/pingcap/tipb v0.0.0-20200522051215-f31a15d98fce // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20190728182440-6a916e37a237 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shirou/gopsutil v2.19.10+incompatible // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726 // indirect
	github.com/siddontang/go-log v0.0.0-20180807004314-8d05993dda07 // indirect
	github.com/skeema/knownhosts v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14 // indirect
	github.com/tjfoc/gmsm v1.3.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/crypto v0.13.0 // indirect
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.13.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220414192740-2d67ff6cf2b4 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/percona/pmm-agent => github.com/taolx0/pmm-agent v0.0.0-20230614092412-936a5cff4635
	github.com/pingcap/log => github.com/pingcap/log v0.0.0-20191012051959-b742a5d432e9
	github.com/pingcap/parser => github.com/sjjian/parser v0.0.0-20231020015929-c5d7ca486d80
	google.golang.org/grpc => google.golang.org/grpc v1.29.0
)
