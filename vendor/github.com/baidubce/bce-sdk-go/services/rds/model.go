/*
 * Copyright 2020 Baidu, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the
 * License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions
 * and limitations under the License.
 */

// model.go - definitions of the request arguments and results data structure model

package rds

import (
	"github.com/baidubce/bce-sdk-go/model"
)

type CreateRdsArgs struct {
	ClientToken       string           `json:"-"`
	Billing           Billing          `json:"billing"`
	PurchaseCount     int              `json:"purchaseCount,omitempty"`
	InstanceName      string           `json:"instanceName,omitempty"`
	Engine            string           `json:"engine"`
	EngineVersion     string           `json:"engineVersion"`
	Category          string           `json:"category,omitempty"`
	CpuCount          int              `json:"cpuCount"`
	MemoryCapacity    float64          `json:"memoryCapacity"`
	VolumeCapacity    int              `json:"volumeCapacity"`
	DiskIoType        string           `json:"diskIoType"`
	ZoneNames         []string         `json:"zoneNames,omitempty"`
	VpcId             string           `json:"vpcId,omitempty"`
	IsDirectPay       bool             `json:"isDirectPay,omitempty"`
	Subnets           []SubnetMap      `json:"subnets,omitempty"`
	Tags              []model.TagModel `json:"tags,omitempty"`
	AutoRenewTimeUnit string           `json:"autoRenewTimeUnit,omitempty"`
	AutoRenewTime     int              `json:"autoRenewTime,omitempty"`
	BgwGroupId        string           `json:"bgwGroupId,omitempty"`
}

type Billing struct {
	PaymentTiming string      `json:"paymentTiming"`
	Reservation   Reservation `json:"reservation,omitempty"`
}

type Reservation struct {
	ReservationLength   int    `json:"reservationLength,omitempty"`
	ReservationTimeUnit string `json:"reservationTimeUnit,omitempty"`
}

type SubnetMap struct {
	ZoneName string `json:"zoneName"`
	SubnetId string `json:"subnetId"`
}

type CreateResult struct {
	InstanceIds []string `json:"instanceIds"`
}

type CreateReadReplicaArgs struct {
	ClientToken      string           `json:"-"`
	Billing          Billing          `json:"billing"`
	PurchaseCount    int              `json:"purchaseCount,omitempty"`
	SourceInstanceId string           `json:"sourceInstanceId"`
	InstanceName     string           `json:"instanceName,omitempty"`
	CpuCount         int              `json:"cpuCount"`
	MemoryCapacity   float64          `json:"memoryCapacity"`
	VolumeCapacity   int              `json:"volumeCapacity"`
	ZoneNames        []string         `json:"zoneNames,omitempty"`
	VpcId            string           `json:"vpcId,omitempty"`
	IsDirectPay      bool             `json:"isDirectPay,omitempty"`
	Subnets          []SubnetMap      `json:"subnets,omitempty"`
	Tags             []model.TagModel `json:"tags,omitempty"`
}

type CreateRdsProxyArgs struct {
	ClientToken      string           `json:"-"`
	Billing          Billing          `json:"billing"`
	SourceInstanceId string           `json:"sourceInstanceId"`
	InstanceName     string           `json:"instanceName,omitempty"`
	NodeAmount       int              `json:"nodeAmount"`
	ZoneNames        []string         `json:"zoneNames,omitempty"`
	VpcId            string           `json:"vpcId,omitempty"`
	IsDirectPay      bool             `json:"isDirectPay,omitempty"`
	Subnets          []SubnetMap      `json:"subnets,omitempty"`
	Tags             []model.TagModel `json:"tags,omitempty"`
}

type ListRdsArgs struct {
	Marker  string
	MaxKeys int
}

type Instance struct {
	InstanceId         string       `json:"instanceId"`
	InstanceName       string       `json:"instanceName"`
	Engine             string       `json:"engine"`
	EngineVersion      string       `json:"engineVersion"`
	Category           string       `json:"category"`
	InstanceStatus     string       `json:"instanceStatus"`
	CpuCount           int          `json:"cpuCount"`
	MemoryCapacity     float64      `json:"memoryCapacity"`
	VolumeCapacity     int          `json:"volumeCapacity"`
	NodeAmount         int          `json:"nodeAmount"`
	UsedStorage        float64      `json:"usedStorage"`
	PublicAccessStatus string       `json:"publicAccessStatus"`
	InstanceCreateTime string       `json:"instanceCreateTime"`
	InstanceExpireTime string       `json:"instanceExpireTime"`
	Endpoint           Endpoint     `json:"endpoint"`
	SyncMode           string       `json:"syncMode"`
	BackupPolicy       BackupPolicy `json:"backupPolicy"`
	Region             string       `json:"region"`
	InstanceType       string       `json:"instanceType"`
	SourceInstanceId   string       `json:"sourceInstanceId"`
	SourceRegion       string       `json:"sourceRegion"`
	ZoneNames          []string     `json:"zoneNames"`
	VpcId              string       `json:"vpcId"`
	Subnets            []Subnet     `json:"subnets"`
	Topology           Topology     `json:"topology"`
	Task               string       `json:"task"`
	PaymentTiming      string       `json:"paymentTiming"`
	BgwGroupId         string       `json:"bgwGroupId"`
}

type ListRdsResult struct {
	Marker      string     `json:"marker"`
	MaxKeys     int        `json:"maxKeys"`
	IsTruncated bool       `json:"isTruncated"`
	NextMarker  string     `json:"nextMarker"`
	Instances   []Instance `json:"instances"`
}

type Subnet struct {
	Name     string `json:"name"`
	SubnetId string `json:"subnetId"`
	ZoneName string `json:"zoneName"`
	Cidr     string `json:"cidr"`
	VpcId    string `json:"vpcId"`
}

type Endpoint struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	VnetIp  string `json:"vnetIp"`
	InetIp  string `json:"inetIp"`
}

type BackupPolicy struct {
	BackupDays    string `json:"backupDays"`
	BackupTime    string `json:"backupTime"`
	Persistent    bool   `json:"persistent"`
	ExpireInDays  int    `json:"expireInDays"`
	FreeSpaceInGB int    `json:"freeSpaceInGb"`
}

type Topology struct {
	Rdsproxy    []string `json:"rdsproxy"`
	Master      []string `json:"master"`
	ReadReplica []string `json:"readReplica"`
}

type ResizeRdsArgs struct {
	CpuCount       int     `json:"cpuCount"`
	MemoryCapacity float64 `json:"memoryCapacity"`
	VolumeCapacity int     `json:"volumeCapacity"`
	NodeAmount     int     `json:"nodeAmount,omitempty"`
	IsDirectPay    bool    `json:"isDirectPay,omitempty"`
}

type CreateAccountArgs struct {
	ClientToken        string              `json:"-"`
	AccountName        string              `json:"accountName"`
	Password           string              `json:"password"`
	AccountType        string              `json:"accountType,omitempty"`
	DatabasePrivileges []DatabasePrivilege `json:"databasePrivileges,omitempty"`
	Desc               string              `json:"desc,omitempty"`
	Type               string              `json:"type,omitempty"`
}

type DatabasePrivilege struct {
	DbName   string `json:"dbName"`
	AuthType string `json:"authType"`
}

type Account struct {
	AccountName        string              `json:"accountName"`
	Remark             string              `json:"remark"`
	Status             string              `json:"status"`
	Type               string              `json:"type"`
	AccountType        string              `json:"accountType"`
	DatabasePrivileges []DatabasePrivilege `json:"databasePrivileges"`
	Desc               string              `json:"desc"`
}

type ListAccountResult struct {
	Accounts []Account `json:"accounts"`
}

type UpdateInstanceNameArgs struct {
	InstanceName string `json:"instanceName"`
}

type ModifySyncModeArgs struct {
	SyncMode string `json:"syncMode"`
}

type ModifyEndpointArgs struct {
	Address string `json:"address"`
}

type ModifyPublicAccessArgs struct {
	PublicAccess bool `json:"publicAccess"`
}

type ModifyBackupPolicyArgs struct {
	BackupDays   string `json:"backupDays"`
	BackupTime   string `json:"backupTime"`
	Persistent   bool   `json:"persistent"`
	ExpireInDays int    `json:"expireInDays"`
}

type GetBackupListArgs struct {
	Marker  string
	MaxKeys int
}

type Snapshot struct {
	SnapshotId          string `json:"backupId"`
	SnapshotSizeInBytes int64  `json:"backupSize"`
	SnapshotType        string `json:"backupType"`
	SnapshotStatus      string `json:"backupStatus"`
	SnapshotStartTime   string `json:"backupStartTime"`
	SnapshotEndTime     string `json:"backupEndTime"`
	DownloadUrl         string `json:"downloadUrl"`
	DownloadExpires     string `json:"downloadExpires"`
}

type GetBackupListResult struct {
	Marker      string     `json:"marker"`
	MaxKeys     int        `json:"maxKeys"`
	IsTruncated bool       `json:"isTruncated"`
	NextMarker  string     `json:"nextMarker"`
	Backups     []Snapshot `json:"backups"`
}

type GetZoneListResult struct {
	Zones []ZoneName `json:"zones"`
}

type ZoneName struct {
	ZoneNames []string `json:"zoneNames"`
}

type ListSubnetsArgs struct {
	VpcId    string `json:"vpcId"`
	ZoneName string `json:"zoneName"`
}

type ListSubnetsResult struct {
	Subnets []Subnet `json:"subnets"`
}

type GetSecurityIpsResult struct {
	Etag        string   `json:"etag"`
	SecurityIps []string `json:"securityIps"`
}

type UpdateSecurityIpsArgs struct {
	SecurityIps []string `json:"securityIps"`
}

type ListParametersResult struct {
	Etag       string      `json:"etag"`
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Name          string `json:"name"`
	DefaultValue  string `json:"defaultValue"`
	Value         string `json:"value"`
	PendingValue  string `json:"pendingValue"`
	Type          string `json:"type"`
	Dynamic       string `json:"dynamic"`
	Modifiable    string `json:"modifiable"`
	AllowedValues string `json:"allowedValues"`
	Desc          string `json:"desc"`
}

type UpdateParameterArgs struct {
	Parameters []KVParameter `json:"parameters"`
}

type KVParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AutoRenewArgs struct {
	InstanceIds       []string `json:"instanceIds"`
	AutoRenewTimeUnit string   `json:"autoRenewTimeUnit"`
	AutoRenewTime     int      `json:"autoRenewTime"`
}

type SlowLogDownloadTaskListResult struct {
	Slowlogs []Slowlog `json:"slowlogs"`
}

type SlowLogDownloadDetail struct {
	Slowlogs []SlowlogDetail `json:"slowlogs"`
}
type Slowlog struct {
	SlowlogId          string `json:"slowlogId"`
	SlowlogSizeInBytes int    `json:"slowlogSizeInBytes"`
	SlowlogStartTime   string `json:"slowlogStartTime"`
	SlowlogEndTime     string `json:"slowlogEndTime"`
}

type SlowlogDetail struct {
	Url             string `json:"url"`
	DownloadExpires string `json:"downloadExpires"`
}
