package config

import (
	"fmt"
	"os"
)

// Figure out if we're self-hosted or on RDS, as well as what ID we can use - Heroku is treated separately
func identifySystem(config ServerConfig) (systemID string, systemType string, systemScope string, systemIDFallback string, systemTypeFallback string, systemScopeFallback string) {
	// Allow overrides from config or env variables
	systemID = config.SystemID
	systemType = config.SystemType
	systemScope = config.SystemScope
	systemIDFallback = config.SystemIDFallback
	systemTypeFallback = config.SystemTypeFallback
	systemScopeFallback = config.SystemScopeFallback

	if config.AwsDbInstanceID != "" || config.AwsDbClusterID != "" || systemType == "amazon_rds" {
		systemType = "amazon_rds"
		if systemScope == "" {
			if config.AwsAccountID != "" {
				clusterPrefix := ""
				if config.AwsDbInstanceID == "" && config.AwsDbClusterID != "" {
					if config.AwsDbClusterReadonly {
						clusterPrefix = "cluster-ro-"
					} else {
						clusterPrefix = "cluster-"
					}
				}
				systemScope = config.AwsRegion + "/" + clusterPrefix + config.AwsAccountID
				if systemScopeFallback == "" {
					systemScopeFallback = config.AwsRegion
				}
			} else {
				systemScope = config.AwsRegion
			}
		}
		if systemID == "" {
			if config.AwsDbInstanceID != "" {
				systemID = config.AwsDbInstanceID
			} else if config.AwsDbClusterID != "" {
				systemID = config.AwsDbClusterID
			}
		}
	} else if config.AzureDbServerName != "" || systemType == "azure_database" {
		systemType = "azure_database"
		if systemID == "" {
			systemID = config.AzureDbServerName
		}
	} else if (config.GcpProjectID != "" && config.GcpCloudSQLInstanceID != "") || (config.GcpProjectID != "" && config.GcpAlloyDBClusterID != "" && config.GcpAlloyDBInstanceID != "") || systemType == "google_cloudsql" {
		systemType = "google_cloudsql"
		if systemScope == "" {
			systemScope = config.GcpProjectID
		}
		if systemID == "" {
			if config.GcpCloudSQLInstanceID != "" {
				systemID = config.GcpCloudSQLInstanceID
			} else if config.GcpAlloyDBClusterID != "" && config.GcpAlloyDBInstanceID != "" {
				systemID = config.GcpAlloyDBClusterID + ":" + config.GcpAlloyDBInstanceID
			}
		}
	} else if (config.CrunchyBridgeClusterID != "") || systemType == "crunchy_bridge" {
		systemType = "crunchy_bridge"
		if systemID == "" {
			systemID = config.CrunchyBridgeClusterID
		}
	} else if (config.AivenProjectID != "" && config.AivenServiceID != "") || systemType == "aiven" {
		systemType = "aiven"
		if systemID == "" {
			systemID = config.AivenServiceID
		}
		if systemScope == "" {
			systemScope = config.AivenProjectID
		}

		if systemTypeFallback == "" {
			systemTypeFallback = "self_hosted"
		}
		if systemIDFallback == "" {
			systemIDFallback = selfManagedSystemID(config)
		}
		if systemScopeFallback == "" {
			systemScopeFallback = selfManagedSystemScope(config)
		}
	} else {
		systemType = "self_hosted"
		if systemID == "" {
			systemID = selfManagedSystemID(config)
			if systemScope == "" {
				systemScope = selfManagedSystemScope(config)
			}
		}
	}
	return
}

func selfManagedSystemID(config ServerConfig) string {
	hostname := config.GetDbHost()
	if hostname == "" || hostname == "localhost" || hostname == "127.0.0.1" {
		hostname, _ = os.Hostname()
	}
	return hostname
}

func selfManagedSystemScope(config ServerConfig) string {
	scope := fmt.Sprintf("%d/%s", config.GetDbPort(), config.GetDbName())
	if config.DbAllNames {
		scope += "*"
	}
	return scope
}
