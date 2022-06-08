// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package errs

import "github.com/pingcap/errors"

const (
	// NotLeaderErr indicates the the non-leader member received the requests which should be received by leader.
	NotLeaderErr = "is not leader"
	// MismatchLeaderErr indicates the the non-leader member received the requests which should be received by leader.
	MismatchLeaderErr = "mismatch leader id"
)

// common error in multiple packages
var (
	ErrGetSourceStore      = errors.Normalize("failed to get the source store", errors.RFCCodeText("PD:common:ErrGetSourceStore"))
	ErrIncorrectSystemTime = errors.Normalize("incorrect system time", errors.RFCCodeText("PD:common:ErrIncorrectSystemTime"))
)

// tso errors
var (
	ErrSetLocalTSOConfig  = errors.Normalize("set local tso config failed, %s", errors.RFCCodeText("PD:tso:ErrSetLocalTSOConfig"))
	ErrGetAllocator       = errors.Normalize("get allocator failed, %s", errors.RFCCodeText("PD:tso:ErrGetAllocator"))
	ErrGetLocalAllocator  = errors.Normalize("get local allocator failed, %s", errors.RFCCodeText("PD:tso:ErrGetLocalAllocator"))
	ErrSyncMaxTS          = errors.Normalize("sync max ts failed, %s", errors.RFCCodeText("PD:tso:ErrSyncMaxTS"))
	ErrResetUserTimestamp = errors.Normalize("reset user timestamp failed, %s", errors.RFCCodeText("PD:tso:ErrResetUserTimestamp"))
	ErrGenerateTimestamp  = errors.Normalize("generate timestamp failed, %s", errors.RFCCodeText("PD:tso:ErrGenerateTimestamp"))
	ErrInvalidTimestamp   = errors.Normalize("invalid timestamp", errors.RFCCodeText("PD:tso:ErrInvalidTimestamp"))
	ErrLogicOverflow      = errors.Normalize("logic part overflow", errors.RFCCodeText("PD:tso:ErrLogicOverflow"))
)

// member errors
var (
	ErrEtcdLeaderNotFound = errors.Normalize("etcd leader not found", errors.RFCCodeText("PD:member:ErrEtcdLeaderNotFound"))
	ErrMarshalLeader      = errors.Normalize("marshal leader failed", errors.RFCCodeText("PD:member:ErrMarshalLeader"))
)

// core errors
var (
	ErrWrongRangeKeys      = errors.Normalize("wrong range keys", errors.RFCCodeText("PD:core:ErrWrongRangeKeys"))
	ErrStoreNotFound       = errors.Normalize("store %v not found", errors.RFCCodeText("PD:core:ErrStoreNotFound"))
	ErrPauseLeaderTransfer = errors.Normalize("store %v is paused for leader transfer", errors.RFCCodeText("PD:core:ErrPauseLeaderTransfer"))
	ErrStoreTombstone      = errors.Normalize("store %v has been removed", errors.RFCCodeText("PD:core:ErrStoreTombstone"))
	ErrStoreDestroyed      = errors.Normalize("store %v has been physically destroyed", errors.RFCCodeText("PD:core:ErrStoreDestroyed"))
	ErrStoreUnhealthy      = errors.Normalize("store %v is unhealthy", errors.RFCCodeText("PD:core:ErrStoreUnhealthy"))
	ErrSlowStoreEvicted    = errors.Normalize("store %v is evited as a slow store", errors.RFCCodeText("PD:core:ErrSlowStoreEvicted"))
)

// client errors
var (
	ErrClientCreateTSOStream = errors.Normalize("create TSO stream failed", errors.RFCCodeText("PD:client:ErrClientCreateTSOStream"))
	ErrClientGetTSOTimeout   = errors.Normalize("get TSO timeout", errors.RFCCodeText("PD:client:ErrClientGetTSOTimeout"))
	ErrClientGetTSO          = errors.Normalize("get TSO failed, %v", errors.RFCCodeText("PD:client:ErrClientGetTSO"))
	ErrClientGetLeader       = errors.Normalize("get leader from %v error", errors.RFCCodeText("PD:client:ErrClientGetLeader"))
	ErrClientGetMember       = errors.Normalize("get member failed", errors.RFCCodeText("PD:client:ErrClientGetMember"))
)

// schedule errors
var (
	ErrUnexpectedOperatorStatus = errors.Normalize("operator with unexpected status", errors.RFCCodeText("PD:schedule:ErrUnexpectedOperatorStatus"))
	ErrUnknownOperatorStep      = errors.Normalize("unknown operator step found", errors.RFCCodeText("PD:schedule:ErrUnknownOperatorStep"))
	ErrMergeOperator            = errors.Normalize("merge operator error, %s", errors.RFCCodeText("PD:schedule:ErrMergeOperator"))
	ErrCreateOperator           = errors.Normalize("unable to create operator, %s", errors.RFCCodeText("PD:schedule:ErrCreateOperator"))
)

// scheduler errors
var (
	ErrSchedulerExisted                 = errors.Normalize("scheduler existed", errors.RFCCodeText("PD:scheduler:ErrSchedulerExisted"))
	ErrSchedulerDuplicated              = errors.Normalize("scheduler duplicated", errors.RFCCodeText("PD:scheduler:ErrSchedulerDuplicated"))
	ErrSchedulerNotFound                = errors.Normalize("scheduler not found", errors.RFCCodeText("PD:scheduler:ErrSchedulerNotFound"))
	ErrScheduleConfigNotExist           = errors.Normalize("the config does not exist", errors.RFCCodeText("PD:scheduler:ErrScheduleConfigNotExist"))
	ErrSchedulerConfig                  = errors.Normalize("wrong scheduler config %s", errors.RFCCodeText("PD:scheduler:ErrSchedulerConfig"))
	ErrCacheOverflow                    = errors.Normalize("cache overflow", errors.RFCCodeText("PD:scheduler:ErrCacheOverflow"))
	ErrInternalGrowth                   = errors.Normalize("unknown interval growth type error", errors.RFCCodeText("PD:scheduler:ErrInternalGrowth"))
	ErrSchedulerCreateFuncNotRegistered = errors.Normalize("create func of %v is not registered", errors.RFCCodeText("PD:scheduler:ErrSchedulerCreateFuncNotRegistered"))
)

// placement errors
var (
	ErrRuleContent   = errors.Normalize("invalid rule content, %s", errors.RFCCodeText("PD:placement:ErrRuleContent"))
	ErrLoadRule      = errors.Normalize("load rule failed", errors.RFCCodeText("PD:placement:ErrLoadRule"))
	ErrLoadRuleGroup = errors.Normalize("load rule group failed", errors.RFCCodeText("PD:placement:ErrLoadRuleGroup"))
	ErrBuildRuleList = errors.Normalize("build rule list failed, %s", errors.RFCCodeText("PD:placement:ErrBuildRuleList"))
)

// cluster errors
var (
	ErrNotBootstrapped = errors.Normalize("TiKV cluster not bootstrapped, please start TiKV first", errors.RFCCodeText("PD:cluster:ErrNotBootstrapped"))
	ErrStoreIsUp       = errors.Normalize("store is still up, please remove store gracefully", errors.RFCCodeText("PD:cluster:ErrStoreIsUp"))
)

// versioninfo errors
var (
	ErrFeatureNotExisted = errors.Normalize("feature not existed", errors.RFCCodeText("PD:versioninfo:ErrFeatureNotExisted"))
)

// autoscaling errors
var (
	ErrUnsupportedMetricsType   = errors.Normalize("unsupported metrics type %v", errors.RFCCodeText("PD:autoscaling:ErrUnsupportedMetricsType"))
	ErrUnsupportedComponentType = errors.Normalize("unsupported component type %v", errors.RFCCodeText("PD:autoscaling:ErrUnsupportedComponentType"))
	ErrUnexpectedType           = errors.Normalize("unexpected type %v", errors.RFCCodeText("PD:autoscaling:ErrUnexpectedType"))
	ErrTypeConversion           = errors.Normalize("type conversion error", errors.RFCCodeText("PD:autoscaling:ErrTypeConversion"))
	ErrEmptyMetricsResponse     = errors.Normalize("metrics response from Prometheus is empty", errors.RFCCodeText("PD:autoscaling:ErrEmptyMetricsResponse"))
	ErrEmptyMetricsResult       = errors.Normalize("result from Prometheus is empty, %s", errors.RFCCodeText("PD:autoscaling:ErrEmptyMetricsResult"))
)

// apiutil errors
var (
	ErrRedirect = errors.Normalize("redirect failed", errors.RFCCodeText("PD:apiutil:ErrRedirect"))
)

// grpcutil errors
var (
	ErrSecurityConfig = errors.Normalize("security config error: %s", errors.RFCCodeText("PD:grpcutil:ErrSecurityConfig"))
)

// server errors
var (
	ErrServiceRegistered     = errors.Normalize("service with path [%s] already registered", errors.RFCCodeText("PD:server:ErrServiceRegistered"))
	ErrAPIInformationInvalid = errors.Normalize("invalid api information, group %s version %s", errors.RFCCodeText("PD:server:ErrAPIInformationInvalid"))
	ErrClientURLEmpty        = errors.Normalize("client url empty", errors.RFCCodeText("PD:server:ErrClientEmpty"))
	ErrLeaderNil             = errors.Normalize("leader is nil", errors.RFCCodeText("PD:server:ErrLeaderNil"))
	ErrCancelStartEtcd       = errors.Normalize("etcd start canceled", errors.RFCCodeText("PD:server:ErrCancelStartEtcd"))
	ErrConfigItem            = errors.Normalize("cannot set invalid configuration", errors.RFCCodeText("PD:server:ErrConfiguration"))
)

// logutil errors
var (
	ErrInitFileLog = errors.Normalize("init file log error, %s", errors.RFCCodeText("PD:logutil:ErrInitFileLog"))
)

// typeutil errors
var (
	ErrBytesToUint64 = errors.Normalize("invalid data, must 8 bytes, but %d", errors.RFCCodeText("PD:typeutil:ErrBytesToUint64"))
)

// The third-party project error.
// url errors
var (
	ErrURLParse      = errors.Normalize("parse url error", errors.RFCCodeText("PD:url:ErrURLParse"))
	ErrQueryUnescape = errors.Normalize("inverse transformation of QueryEscape error", errors.RFCCodeText("PD:url:ErrQueryUnescape"))
)

// grpc errors
var (
	ErrGRPCDial         = errors.Normalize("dial error", errors.RFCCodeText("PD:grpc:ErrGRPCDial"))
	ErrCloseGRPCConn    = errors.Normalize("close gRPC connection failed", errors.RFCCodeText("PD:grpc:ErrCloseGRPCConn"))
	ErrGRPCSend         = errors.Normalize("send request error", errors.RFCCodeText("PD:grpc:ErrGRPCSend"))
	ErrGRPCRecv         = errors.Normalize("receive response error", errors.RFCCodeText("PD:grpc:ErrGRPCRecv"))
	ErrGRPCCloseSend    = errors.Normalize("close send error", errors.RFCCodeText("PD:grpc:ErrGRPCCloseSend"))
	ErrGRPCCreateStream = errors.Normalize("create stream error", errors.RFCCodeText("PD:grpc:ErrGRPCCreateStream"))
)

// proto errors
var (
	ErrProtoUnmarshal = errors.Normalize("failed to unmarshal proto", errors.RFCCodeText("PD:proto:ErrProtoUnmarshal"))
	ErrProtoMarshal   = errors.Normalize("failed to marshal proto", errors.RFCCodeText("PD:proto:ErrProtoMarshal"))
)

// etcd errors
var (
	ErrNewEtcdClient     = errors.Normalize("new etcd client failed", errors.RFCCodeText("PD:etcd:ErrNewEtcdClient"))
	ErrStartEtcd         = errors.Normalize("start etcd failed", errors.RFCCodeText("PD:etcd:ErrStartEtcd"))
	ErrEtcdURLMap        = errors.Normalize("etcd url map error", errors.RFCCodeText("PD:etcd:ErrEtcdURLMap"))
	ErrEtcdGrantLease    = errors.Normalize("etcd lease failed", errors.RFCCodeText("PD:etcd:ErrEtcdGrantLease"))
	ErrEtcdTxnInternal   = errors.Normalize("internal etcd transaction error occurred", errors.RFCCodeText("PD:etcd:ErrEtcdTxnInternal"))
	ErrEtcdTxnConflict   = errors.Normalize("etcd transaction failed, conflicted and rolled back", errors.RFCCodeText("PD:etcd:ErrEtcdTxnConflict"))
	ErrEtcdKVPut         = errors.Normalize("etcd KV put failed", errors.RFCCodeText("PD:etcd:ErrEtcdKVPut"))
	ErrEtcdKVDelete      = errors.Normalize("etcd KV delete failed", errors.RFCCodeText("PD:etcd:ErrEtcdKVDelete"))
	ErrEtcdKVGet         = errors.Normalize("etcd KV get failed", errors.RFCCodeText("PD:etcd:ErrEtcdKVGet"))
	ErrEtcdKVGetResponse = errors.Normalize("etcd invalid get value response %v, must only one", errors.RFCCodeText("PD:etcd:ErrEtcdKVGetResponse"))
	ErrEtcdGetCluster    = errors.Normalize("etcd get cluster from remote peer failed", errors.RFCCodeText("PD:etcd:ErrEtcdGetCluster"))
	ErrEtcdMoveLeader    = errors.Normalize("etcd move leader error", errors.RFCCodeText("PD:etcd:ErrEtcdMoveLeader"))
	ErrEtcdTLSConfig     = errors.Normalize("etcd TLS config error", errors.RFCCodeText("PD:etcd:ErrEtcdTLSConfig"))
	ErrEtcdWatcherCancel = errors.Normalize("watcher canceled", errors.RFCCodeText("PD:etcd:ErrEtcdWatcherCancel"))
	ErrCloseEtcdClient   = errors.Normalize("close etcd client failed", errors.RFCCodeText("PD:etcd:ErrCloseEtcdClient"))
	ErrEtcdMemberList    = errors.Normalize("etcd member list failed", errors.RFCCodeText("PD:etcd:ErrEtcdMemberList"))
)

// dashboard errors
var (
	ErrDashboardStart = errors.Normalize("start dashboard failed", errors.RFCCodeText("PD:dashboard:ErrDashboardStart"))
	ErrDashboardStop  = errors.Normalize("stop dashboard failed", errors.RFCCodeText("PD:dashboard:ErrDashboardStop"))
)

// strconv errors
var (
	ErrStrconvParseInt   = errors.Normalize("parse int error", errors.RFCCodeText("PD:strconv:ErrStrconvParseInt"))
	ErrStrconvParseUint  = errors.Normalize("parse uint error", errors.RFCCodeText("PD:strconv:ErrStrconvParseUint"))
	ErrStrconvParseFloat = errors.Normalize("parse float error", errors.RFCCodeText("PD:strconv:ErrStrconvParseFloat"))
)

// prometheus errors
var (
	ErrPrometheusPushMetrics  = errors.Normalize("push metrics to gateway failed", errors.RFCCodeText("PD:prometheus:ErrPrometheusPushMetrics"))
	ErrPrometheusCreateClient = errors.Normalize("create client error", errors.RFCCodeText("PD:prometheus:ErrPrometheusCreateClient"))
	ErrPrometheusQuery        = errors.Normalize("query error", errors.RFCCodeText("PD:prometheus:ErrPrometheusQuery"))
)

// http errors
var (
	ErrSendRequest    = errors.Normalize("send HTTP request failed", errors.RFCCodeText("PD:http:ErrSendRequest"))
	ErrWriteHTTPBody  = errors.Normalize("write HTTP body failed", errors.RFCCodeText("PD:http:ErrWriteHTTPBody"))
	ErrNewHTTPRequest = errors.Normalize("new HTTP request failed", errors.RFCCodeText("PD:http:ErrNewHTTPRequest"))
)

// ioutil error
var (
	ErrIORead = errors.Normalize("IO read error", errors.RFCCodeText("PD:ioutil:ErrIORead"))
)

// os error
var (
	ErrOSOpen = errors.Normalize("open error", errors.RFCCodeText("PD:os:ErrOSOpen"))
)

// dir error
var (
	ErrReadDirName = errors.Normalize("read dir name error", errors.RFCCodeText("PD:dir:ErrReadDirName"))
)

// netstat error
var (
	ErrNetstatTCPSocks = errors.Normalize("TCP socks error", errors.RFCCodeText("PD:netstat:ErrNetstatTCPSocks"))
)

// hex error
var (
	ErrHexDecodingString = errors.Normalize("decode string %s error", errors.RFCCodeText("PD:hex:ErrHexDecodingString"))
)

// filepath errors
var (
	ErrFilePathAbs = errors.Normalize("failed to convert a path to absolute path", errors.RFCCodeText("PD:filepath:ErrFilePathAbs"))
)

// plugin errors
var (
	ErrLoadPlugin       = errors.Normalize("failed to load plugin", errors.RFCCodeText("PD:plugin:ErrLoadPlugin"))
	ErrLookupPluginFunc = errors.Normalize("failed to lookup plugin function", errors.RFCCodeText("PD:plugin:ErrLookupPluginFunc"))
)

// json errors
var (
	ErrJSONMarshal   = errors.Normalize("failed to marshal json", errors.RFCCodeText("PD:json:ErrJSONMarshal"))
	ErrJSONUnmarshal = errors.Normalize("failed to unmarshal json", errors.RFCCodeText("PD:json:ErrJSONUnmarshal"))
)

// leveldb errors
var (
	ErrLevelDBClose = errors.Normalize("close leveldb error", errors.RFCCodeText("PD:leveldb:ErrLevelDBClose"))
	ErrLevelDBWrite = errors.Normalize("leveldb write error", errors.RFCCodeText("PD:leveldb:ErrLevelDBWrite"))
	ErrLevelDBOpen  = errors.Normalize("leveldb open file error", errors.RFCCodeText("PD:leveldb:ErrLevelDBOpen"))
)

// semver
var (
	ErrSemverNewVersion = errors.Normalize("new version error", errors.RFCCodeText("PD:semver:ErrSemverNewVersion"))
)

// log
var (
	ErrInitLogger = errors.Normalize("init logger error", errors.RFCCodeText("PD:log:ErrInitLogger"))
)

// encryption
var (
	ErrEncryptionInvalidMethod      = errors.Normalize("invalid encryption method", errors.RFCCodeText("PD:encryption:ErrEncryptionInvalidMethod"))
	ErrEncryptionInvalidConfig      = errors.Normalize("invalid config", errors.RFCCodeText("PD:encryption:ErrEncryptionInvalidConfig"))
	ErrEncryptionGenerateIV         = errors.Normalize("fail to generate iv", errors.RFCCodeText("PD:encryption:ErrEncryptionGenerateIV"))
	ErrEncryptionGCMEncrypt         = errors.Normalize("GCM encryption fail", errors.RFCCodeText("PD:encryption:ErrEncryptionGCMEncrypt"))
	ErrEncryptionGCMDecrypt         = errors.Normalize("GCM decryption fail", errors.RFCCodeText("PD:encryption:ErrEncryptionGCMDecrypt"))
	ErrEncryptionCTREncrypt         = errors.Normalize("CTR encryption fail", errors.RFCCodeText("PD:encryption:ErrEncryptionCTREncrypt"))
	ErrEncryptionCTRDecrypt         = errors.Normalize("CTR decryption fail", errors.RFCCodeText("PD:encryption:ErrEncryptionCTRDecrypt"))
	ErrEncryptionEncryptRegion      = errors.Normalize("encrypt region fail", errors.RFCCodeText("PD:encryption:ErrEncryptionEncryptRegion"))
	ErrEncryptionDecryptRegion      = errors.Normalize("decrypt region fail", errors.RFCCodeText("PD:encryption:ErrEncryptionDecryptRegion"))
	ErrEncryptionNewDataKey         = errors.Normalize("fail to generate data key", errors.RFCCodeText("PD:encryption:ErrEncryptionNewDataKey"))
	ErrEncryptionNewMasterKey       = errors.Normalize("fail to get master key", errors.RFCCodeText("PD:encryption:ErrEncryptionNewMasterKey"))
	ErrEncryptionCurrentKeyNotFound = errors.Normalize("current data key not found", errors.RFCCodeText("PD:encryption:ErrEncryptionCurrentKeyNotFound"))
	ErrEncryptionKeyNotFound        = errors.Normalize("data key not found", errors.RFCCodeText("PD:encryption:ErrEncryptionKeyNotFound"))
	ErrEncryptionKeysWatcher        = errors.Normalize("data key watcher error", errors.RFCCodeText("PD:encryption:ErrEncryptionKeysWatcher"))
	ErrEncryptionLoadKeys           = errors.Normalize("load data keys error", errors.RFCCodeText("PD:encryption:ErrEncryptionLoadKeys"))
	ErrEncryptionRotateDataKey      = errors.Normalize("failed to rotate data key", errors.RFCCodeText("PD:encryption:ErrEncryptionRotateDataKey"))
	ErrEncryptionSaveDataKeys       = errors.Normalize("failed to save data keys", errors.RFCCodeText("PD:encryption:ErrEncryptionSaveDataKeys"))
	ErrEncryptionKMS                = errors.Normalize("KMS error", errors.RFCCodeText("PD:ErrEncryptionKMS"))
)
