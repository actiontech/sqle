package config

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Config struct {
	Servers []ServerConfig
}

// ServerIdentifier -
//
//	Unique identity of each configured server, for deduplication inside the collector.
//
//	Note we intentionally don't include the Fallback variables in the identifier, since that is mostly intended
//	to help transition systems when their "identity" is altered due to collector changes - in the collector we rely
//	on the non-Fallback values only.
type ServerIdentifier struct {
	APIKey      string
	APIBaseURL  string
	SystemID    string
	SystemType  string
	SystemScope string
}

// ServerConfig -
//
//	Contains the information how to connect to a Postgres instance,
//	with optional AWS credentials to get metrics
//	from AWS CloudWatch as well as RDS logfiles
type ServerConfig struct {
	APIKey     string `ini:"api_key"`
	APIBaseURL string `ini:"api_base_url"`

	ErrorCallback   string `ini:"error_callback"`
	SuccessCallback string `ini:"success_callback"`

	EnableReports    bool `ini:"enable_reports"`
	DisableLogs      bool `ini:"disable_logs"`
	DisableActivity  bool `ini:"disable_activity"`
	EnableLogExplain bool `ini:"enable_log_explain"`

	DbURL                 string `ini:"db_url"`
	DbName                string `ini:"db_name"`
	DbUsername            string `ini:"db_username"`
	DbPassword            string `ini:"db_password"`
	DbHost                string `ini:"db_host"`
	DbPort                int    `ini:"db_port"`
	DbSslMode             string `ini:"db_sslmode"`
	DbSslRootCert         string `ini:"db_sslrootcert"`
	DbSslRootCertContents string `ini:"db_sslrootcert_contents"`
	DbSslCert             string `ini:"db_sslcert"`
	DbSslCertContents     string `ini:"db_sslcert_contents"`
	DbSslKey              string `ini:"db_sslkey"`
	DbSslKeyContents      string `ini:"db_sslkey_contents"`
	DbUseIamAuth          bool   `ini:"db_use_iam_auth"`

	// We have to do some tricks to support sslmode=prefer, namely we have to
	// first try an SSL connection (= require), and if that fails change the
	// sslmode to none
	DbSslModePreferFailed bool

	DbExtraNames []string // Additional databases that should be fetched (determined by additional databases in db_name)
	DbAllNames   bool     // All databases except template databases should be fetched (determined by * in the db_name list)

	AwsRegion               string `ini:"aws_region"`
	AwsAccountID            string `ini:"aws_account_id"`
	AwsDbInstanceID         string `ini:"aws_db_instance_id"`
	AwsDbClusterID          string `ini:"aws_db_cluster_id"`
	AwsDbClusterReadonly    bool   `ini:"aws_db_cluster_readonly"`
	AwsAccessKeyID          string `ini:"aws_access_key_id"`
	AwsSecretAccessKey      string `ini:"aws_secret_access_key"`
	AwsAssumeRole           string `ini:"aws_assume_role"`
	AwsWebIdentityTokenFile string `ini:"aws_web_identity_token_file"`
	AwsRoleArn              string `ini:"aws_role_arn"`

	// Support for custom AWS endpoints
	// See https://docs.aws.amazon.com/sdk-for-go/api/aws/endpoints/
	AwsEndpointSigningRegion       string `ini:"aws_endpoint_signing_region"`
	AwsEndpointSigningRegionLegacy string `ini:"aws_endpoint_rds_signing_region"`
	AwsEndpointRdsURL              string `ini:"aws_endpoint_rds_url"`
	AwsEndpointEc2URL              string `ini:"aws_endpoint_ec2_url"`
	AwsEndpointCloudwatchURL       string `ini:"aws_endpoint_cloudwatch_url"`
	AwsEndpointCloudwatchLogsURL   string `ini:"aws_endpoint_cloudwatch_logs_url"`

	AzureDbServerName          string `ini:"azure_db_server_name"`
	AzureEventhubNamespace     string `ini:"azure_eventhub_namespace"`
	AzureEventhubName          string `ini:"azure_eventhub_name"`
	AzureADTenantID            string `ini:"azure_ad_tenant_id"`
	AzureADClientID            string `ini:"azure_ad_client_id"`
	AzureADClientSecret        string `ini:"azure_ad_client_secret"`
	AzureADCertificatePath     string `ini:"azure_ad_certificate_path"`
	AzureADCertificatePassword string `ini:"azure_ad_certificate_password"`

	GcpProjectID          string `ini:"gcp_project_id"` // Optional for CloudSQL (you can pass the full "Connection name" as the instance ID)
	GcpCloudSQLInstanceID string `ini:"gcp_cloudsql_instance_id"`
	GcpAlloyDBClusterID   string `ini:"gcp_alloydb_cluster_id"`
	GcpAlloyDBInstanceID  string `ini:"gcp_alloydb_instance_id"`
	GcpPubsubSubscription string `ini:"gcp_pubsub_subscription"`
	GcpCredentialsFile    string `ini:"gcp_credentials_file"`

	CrunchyBridgeClusterID string `ini:"crunchy_bridge_cluster_id"`

	AivenProjectID string `ini:"aiven_project_id"`
	AivenServiceID string `ini:"aiven_service_id"`

	SectionName string
	Identifier  ServerIdentifier

	SystemID            string `ini:"api_system_id"`
	SystemType          string `ini:"api_system_type"`
	SystemScope         string `ini:"api_system_scope"`
	SystemIDFallback    string `ini:"api_system_id_fallback"`
	SystemTypeFallback  string `ini:"api_system_type_fallback"`
	SystemScopeFallback string `ini:"api_system_scope_fallback"`

	AlwaysCollectSystemData bool `ini:"always_collect_system_data"`
	DisableCitusSchemaStats bool `ini:"disable_citus_schema_stats"`

	// Configures the location where logfiles are - this can either be a directory,
	// or a file - needs to readable by the regular pganalyze user
	LogLocation string `ini:"db_log_location"`

	// Configures the collector to tail a local docker container using
	// "docker logs -t" - this is currently experimental and mostly intended for
	// development and debugging. The value needs to be the name of the container.
	LogDockerTail string `ini:"db_log_docker_tail"`

	// Configures the collector to start a built-in syslog server that listens
	// on the specifed "hostname:port" for Postgres log messages
	LogSyslogServer string `ini:"db_log_syslog_server"`

	// Configures the collector to use the "pg_read_file" (superuser) or
	// "pganalyze.read_log_file" (helper) function to retrieve log data
	// directly over the Postgres connection. This only works when superuser
	// access to the server is possible, either directly, or via the helper
	// function. Used by default for Crunchy Bridge.
	LogPgReadFile bool `ini:"db_log_pg_read_file"`

	// Specifies a table pattern to ignore - no statistics will be collected for
	// tables that match the name. This uses Golang's filepath.Match function for
	// comparison, so you can e.g. use "*" for wildcard matching.
	//
	// Deprecated: Please use ignore_schema_regexp instead, since that uses an
	// optimized code path in the collector and can avoid long-running queries.
	IgnoreTablePattern string `ini:"ignore_table_pattern"`

	// Specifies a regular expression to ignore - no statistics will be collected for
	// tables, views, functions, or schemas that match the name. Note that the match
	// is applied to the '.'-joined concantenation of schema name and object name.
	// E.g., to ignore tables that start with "ignored_", set this to "^ignored_". To
	// ignore table "foo" only in the public schema, set to "^public\.foo$" (N.B.: you
	// should escape the dot since that has special meaning in a regexp).
	IgnoreSchemaRegexp string `ini:"ignore_schema_regexp"`

	// Specifies the frequency of query statistics collection in seconds
	//
	// Currently supported values: 600 (10 minutes), 60 (1 minute)
	//
	// Defaults to once per minute (60)
	QueryStatsInterval int `ini:"query_stats_interval"`

	// Maximum connections allowed to the database with the collector
	// application_name, in order to protect against accidental connection leaks
	// in the collector
	//
	// This defaults to 10 connections, but you may want to raise this when running
	// the collector multiple times against the same database server
	MaxCollectorConnections int `ini:"max_collector_connections"`

	// Do not monitor this server while it is a replica (according to pg_is_in_recovery),
	// but keep checking on standard snapshot intervals and automatically start monitoring
	// once the server is promoted
	SkipIfReplica bool `ini:"skip_if_replica"`

	// Configuration for PII filtering
	FilterLogSecret   string `ini:"filter_log_secret"`   // none/all/credential/parsing_error/statement_text/statement_parameter/table_data/ops/unidentified (comma separated)
	FilterQuerySample string `ini:"filter_query_sample"` // none/normalize/all (defaults to "none")
	FilterQueryText   string `ini:"filter_query_text"`   // none/unparsable (defaults to "unparsable")

	// HTTP proxy overrides
	HTTPProxy  string `ini:"http_proxy"`
	HTTPSProxy string `ini:"https_proxy"`
	NoProxy    string `ini:"no_proxy"`

	// HTTP clients to be used for API connections
	HTTPClient          *http.Client
	HTTPClientWithRetry *http.Client
}

// SupportsLogDownload - Determines whether the specified config can download logs
func (config ServerConfig) SupportsLogDownload() bool {
	return config.AwsDbInstanceID != "" || config.AwsDbClusterID != "" || config.LogPgReadFile
}

// GetPqOpenString - Gets the database configuration as a string that can be passed to lib/pq for connecting
func (config ServerConfig) GetPqOpenString(dbNameOverride string, passwordOverride string) (string, error) {
	var dbUsername, dbPassword, dbName, dbHost, dbSslMode, dbSslRootCert, dbSslCert, dbSslKey string
	var dbPort int

	if config.DbURL != "" {
		u, err := url.Parse(config.DbURL)
		if err != nil {
			return "", fmt.Errorf("Failed to parse database URL: %w", err)
		}

		if u.User != nil {
			dbUsername = u.User.Username()
			dbPassword, _ = u.User.Password()
		}

		if u.Path != "" {
			dbName = u.Path[1:len(u.Path)]
		}

		hostSplits := strings.SplitN(u.Host, ":", 2)
		dbHost = hostSplits[0]
		if len(hostSplits) > 1 {
			dbPort, _ = strconv.Atoi(hostSplits[1])
		}

		querySplits := strings.Split(u.RawQuery, "&")
		for _, querySplit := range querySplits {
			keyValue := strings.SplitN(querySplit, "=", 2)
			switch keyValue[0] {
			case "sslmode":
				dbSslMode = keyValue[1]
			case "sslrootcert":
				dbSslRootCert = keyValue[1]
			case "sslcert":
				dbSslCert = keyValue[1]
			case "sslkey":
				dbSslKey = keyValue[1]
			}
		}
	}

	dbinfo := []string{}

	if config.DbUsername != "" {
		dbUsername = config.DbUsername
	}
	if passwordOverride != "" {
		dbPassword = passwordOverride
	} else if config.DbPassword != "" {
		dbPassword = config.DbPassword
	}
	if dbNameOverride != "" {
		dbName = dbNameOverride
	} else if config.DbName != "" {
		dbName = config.DbName
	}
	if config.DbHost != "" {
		dbHost = config.DbHost
	}
	if config.DbPort != 0 {
		dbPort = config.DbPort
	}
	if config.DbSslMode != "" {
		dbSslMode = config.DbSslMode
	}
	if config.DbSslRootCert != "" {
		dbSslRootCert = config.DbSslRootCert
	}
	if config.DbSslCert != "" {
		dbSslCert = config.DbSslCert
	}
	if config.DbSslKey != "" {
		dbSslKey = config.DbSslKey
	}

	// Defaults if nothing is set
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == 0 {
		dbPort = 5432
	}
	if dbSslMode == "" {
		dbSslMode = "prefer"
	}

	// Handle SSL mode prefer
	if dbSslMode == "prefer" {
		if config.DbSslModePreferFailed {
			dbSslMode = "disable"
		} else {
			dbSslMode = "require"
		}
	}

	// Handle SSL certificates shipped with the collector
	//
	// Note: "rds-ca-2019-root" is a legacy cert expiring in 2024 that is part of the global CA set
	if dbSslRootCert == "rds-ca-2019-root" || dbSslRootCert == "rds-ca-global" {
		dbSslRootCert = "/usr/share/pganalyze-collector/sslrootcert/rds-ca-global.pem"
	}

	// Generate the actual string
	if dbUsername != "" {
		dbinfo = append(dbinfo, fmt.Sprintf("user='%s'", strings.Replace(dbUsername, "'", "\\'", -1)))
	}
	if dbPassword != "" {
		dbinfo = append(dbinfo, fmt.Sprintf("password='%s'", strings.Replace(dbPassword, "'", "\\'", -1)))
	}
	if dbName != "" {
		dbinfo = append(dbinfo, fmt.Sprintf("dbname='%s'", strings.Replace(dbName, "'", "\\'", -1)))
	}
	if dbHost != "" {
		dbinfo = append(dbinfo, fmt.Sprintf("host='%s'", strings.Replace(dbHost, "'", "\\'", -1)))
	}
	if dbPort != 0 {
		dbinfo = append(dbinfo, fmt.Sprintf("port=%d", dbPort))
	}
	if dbSslMode != "" {
		dbinfo = append(dbinfo, fmt.Sprintf("sslmode=%s", dbSslMode))
	}
	if dbSslRootCert != "" {
		dbinfo = append(dbinfo, fmt.Sprintf("sslrootcert='%s'", strings.Replace(dbSslRootCert, "'", "\\'", -1)))
	}
	if dbSslCert != "" {
		dbinfo = append(dbinfo, fmt.Sprintf("sslcert='%s'", strings.Replace(dbSslCert, "'", "\\'", -1)))
	}
	if dbSslKey != "" {
		dbinfo = append(dbinfo, fmt.Sprintf("sslkey='%s'", strings.Replace(dbSslKey, "'", "\\'", -1)))
	}
	dbinfo = append(dbinfo, "connect_timeout=10")

	return strings.Join(dbinfo, " "), nil
}

// GetDbHost - Gets the database hostname from the given configuration
func (config ServerConfig) GetDbHost() string {
	if config.DbURL != "" {
		u, err := url.Parse(config.DbURL)
		if err != nil {
			return ""
		}
		parts := strings.Split(u.Host, ":")
		return parts[0]
	}

	return config.DbHost
}

func (config ServerConfig) GetDbURLRedacted() string {
	if config.DbURL == "" {
		return ""
	}

	u, err := url.Parse(config.DbURL)
	if err != nil {
		return "<unparsable>"
	}

	u.User = url.User(u.User.Username())
	return u.String()
}

// GetDbPort - Gets the database port from the given configuration
func (config ServerConfig) GetDbPort() int {
	if config.DbURL != "" {
		u, err := url.Parse(config.DbURL)
		if err != nil {
			return 5432
		}
		parts := strings.Split(u.Host, ":")

		if len(parts) == 2 {
			port, _ := strconv.Atoi(parts[1])
			return port
		}

		return 5432
	}

	return config.DbPort
}

// GetDbUsername - Gets the database hostname from the given configuration
func (config ServerConfig) GetDbUsername() string {
	if config.DbURL != "" {
		u, err := url.Parse(config.DbURL)
		if err != nil {
			return ""
		}
		if u != nil && u.User != nil {
			return u.User.Username()
		}
	}

	return config.DbUsername
}

// GetDbName - Gets the database name from the given configuration
func (config ServerConfig) GetDbName() string {
	if config.DbURL != "" {
		u, err := url.Parse(config.DbURL)
		if err != nil {
			return ""
		}
		if len(u.Path) > 0 {
			return u.Path[1:len(u.Path)]
		}
	}

	return config.DbName
}
