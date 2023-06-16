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

// rds.go - the rds APIs definition supported by the RDS service
package rds

import (
	"fmt"
	"strconv"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/http"
)

// CreateRds - create a RDS with the specific parameters
//
// PARAMS:
//     - args: the arguments to create a rds
// RETURNS:
//     - *InstanceIds: the result of create RDS, contains new RDS's instanceIds
//     - error: nil if success otherwise the specific error
func (c *Client) CreateRds(args *CreateRdsArgs) (*CreateResult, error) {
	if args == nil {
		return nil, fmt.Errorf("unset args")
	}

	if args.Engine == "" {
		return nil, fmt.Errorf("unset Engine")
	}

	if args.EngineVersion == "" {
		return nil, fmt.Errorf("unset EngineVersion")
	}

	if args.Billing.PaymentTiming == "" {
		return nil, fmt.Errorf("unset PaymentTiming")
	}

	result := &CreateResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.POST).
		WithURL(getRdsUri()).
		WithQueryParamFilter("clientToken", args.ClientToken).
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		WithResult(result).
		Do()

	return result, err
}

// CreateReadReplica - create a readReplica RDS with the specific parameters
//
// PARAMS:
//     - args: the arguments to create a readReplica rds
// RETURNS:
//     - *InstanceIds: the result of create a readReplica RDS, contains the readReplica RDS's instanceIds
//     - error: nil if success otherwise the specific error
func (c *Client) CreateReadReplica(args *CreateReadReplicaArgs) (*CreateResult, error) {
	if args == nil {
		return nil, fmt.Errorf("unset args")
	}

	if args.SourceInstanceId == "" {
		return nil, fmt.Errorf("unset SourceInstanceId")
	}

	if args.Billing.PaymentTiming == "" {
		return nil, fmt.Errorf("unset PaymentTiming")
	}

	result := &CreateResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.POST).
		WithURL(getRdsUri()).
		WithQueryParamFilter("clientToken", args.ClientToken).
		WithQueryParam("readReplica", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		WithResult(result).
		Do()

	return result, err
}

// CreateRdsProxy - create a proxy RDS with the specific parameters
//
// PARAMS:
//     - args: the arguments to create a readReplica rds
// RETURNS:
//     - *InstanceIds: the result of create a readReplica RDS, contains the readReplica RDS's instanceIds
//     - error: nil if success otherwise the specific error
func (c *Client) CreateRdsProxy(args *CreateRdsProxyArgs) (*CreateResult, error) {
	if args == nil {
		return nil, fmt.Errorf("unset args")
	}

	if args.SourceInstanceId == "" {
		return nil, fmt.Errorf("unset SourceInstanceId")
	}

	if args.Billing.PaymentTiming == "" {
		return nil, fmt.Errorf("unset PaymentTiming")
	}

	result := &CreateResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.POST).
		WithURL(getRdsUri()).
		WithQueryParamFilter("clientToken", args.ClientToken).
		WithQueryParam("rdsproxy", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		WithResult(result).
		Do()

	return result, err
}

// ListRds - list all RDS with the specific parameters
//
// PARAMS:
//     - args: the arguments to list all RDS
// RETURNS:
//     - *ListRdsResult: the result of list all RDS, contains all rds' meta
//     - error: nil if success otherwise the specific error
func (c *Client) ListRds(args *ListRdsArgs) (*ListRdsResult, error) {
	if args == nil {
		args = &ListRdsArgs{}
	}

	if args.MaxKeys <= 0 || args.MaxKeys > 1000 {
		args.MaxKeys = 1000
	}

	result := &ListRdsResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUri()).
		WithQueryParamFilter("marker", args.Marker).
		WithQueryParamFilter("maxKeys", strconv.Itoa(args.MaxKeys)).
		WithResult(result).
		Do()

	return result, err
}

// GetDetail - get a specific rds Instance's detail
//
// PARAMS:
//     - instanceId: the specific rds Instance's ID
// RETURNS:
//     - *Instance: the specific rdsInstance's detail
//     - error: nil if success otherwise the specific error
func (c *Client) GetDetail(instanceId string) (*Instance, error) {
	result := &Instance{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId)).
		WithResult(result).
		Do()

	return result, err
}

// DeleteRds - delete a rds
//
// PARAMS:
//     - instanceIds: the specific instanceIds
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) DeleteRds(instanceIds string) error {
	return bce.NewRequestBuilder(c).
		WithMethod(http.DELETE).
		WithURL(getRdsUri()).
		WithQueryParamFilter("instanceIds", instanceIds).
		Do()
}

// ResizeRds - resize an RDS with the specific parameters
//
// PARAMS:
//     - instanceId: the specific instanceId
//     - args: the arguments to resize an RDS
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) ResizeRds(instanceId string, args *ResizeRdsArgs) error {
	if args == nil {
		return fmt.Errorf("unset args")
	}

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)).
		WithQueryParam("resize", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// CreateAccount - create a account with the specific parameters
//
// PARAMS:
//     - instanceId: the specific instanceId
//     - args: the arguments to create a account
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) CreateAccount(instanceId string, args *CreateAccountArgs) error {
	if args == nil {
		return fmt.Errorf("unset args")
	}

	if args.AccountName == "" {
		return fmt.Errorf("unset AccountName")
	}

	if args.Password == "" {
		return fmt.Errorf("unset Password")
	}

	cryptedPass, err := Aes128EncryptUseSecreteKey(c.Config.Credentials.SecretAccessKey, args.Password)
	if err != nil {
		return err
	}
	args.Password = cryptedPass

	return bce.NewRequestBuilder(c).
		WithMethod(http.POST).
		WithURL(getRdsUriWithInstanceId(instanceId)+"/account").
		WithQueryParamFilter("clientToken", args.ClientToken).
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// ListAccount - list all account of a RDS instance with the specific parameters
//
// PARAMS:
//     - instanceId: the specific rds Instance's ID
// RETURNS:
//     - *ListAccountResult: the result of list all account, contains all accounts' meta
//     - error: nil if success otherwise the specific error
func (c *Client) ListAccount(instanceId string) (*ListAccountResult, error) {
	result := &ListAccountResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId) + "/account").
		WithResult(result).
		Do()

	return result, err
}

// GetAccount - get an account of a RDS instance with the specific parameters
//
// PARAMS:
//     - instanceId: the specific rds Instance's ID
//     - accountName: the specific account's name
// RETURNS:
//     - *Account: the account's meta
//     - error: nil if success otherwise the specific error
func (c *Client) GetAccount(instanceId, accountName string) (*Account, error) {
	result := &Account{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId) + "/account/" + accountName).
		WithResult(result).
		Do()

	return result, err
}

// DeleteAccount - delete an account of a RDS instance
//
// PARAMS:
//     - instanceIds: the specific instanceIds
//     - accountName: the specific account's name
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) DeleteAccount(instanceId, accountName string) error {
	return bce.NewRequestBuilder(c).
		WithMethod(http.DELETE).
		WithURL(getRdsUriWithInstanceId(instanceId) + "/account/" + accountName).
		Do()
}

// RebootInstance - reboot a specified instance
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - instanceId: id of the instance to be rebooted
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) RebootInstance(instanceId string) error {

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)).
		WithQueryParam("reboot", "").
		Do()
}

// UpdateInstanceName - update name of a specified instance
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - instanceId: id of the instance
//     - args: the arguments to update instanceName
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) UpdateInstanceName(instanceId string, args *UpdateInstanceNameArgs) error {

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)).
		WithQueryParam("rename", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// UpdateSyncMode - update sync mode of a specified instance
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - instanceId: id of the instance
//     - args: the arguments to update syncMode
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) ModifySyncMode(instanceId string, args *ModifySyncModeArgs) error {

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)).
		WithQueryParam("modifySyncMode", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// ModifyEndpoint - modify the prefix of endpoint
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - instanceId: id of the instance
//     - args: the arguments to modify endpoint
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) ModifyEndpoint(instanceId string, args *ModifyEndpointArgs) error {

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)).
		WithQueryParam("modifyEndpoint", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// ModifyPublicAccess - modify public access
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - instanceId: id of the instance
//     - args: the arguments to modify public access
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) ModifyPublicAccess(instanceId string, args *ModifyPublicAccessArgs) error {

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)).
		WithQueryParam("modifyPublicAccess", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// ModifyBackupPolicy - modify backup policy
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - instanceId: id of the instance
//     - args: the arguments to modify public access
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) ModifyBackupPolicy(instanceId string, args *ModifyBackupPolicyArgs) error {

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)).
		WithQueryParam("modifyBackupPolicy", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// GetBackupList - get backup list of the instance
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - instanceId: id of the instance
// RETURNS:
//     - *GetBackupListResult: result of the backup list
//     - error: nil if success otherwise the specific error
func (c *Client) GetBackupList(instanceId string, args *GetBackupListArgs) (*GetBackupListResult, error) {

	if args == nil {
		args = &GetBackupListArgs{}
	}

	if args.MaxKeys <= 0 || args.MaxKeys > 1000 {
		args.MaxKeys = 1000
	}

	result := &GetBackupListResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId)+"/backup").
		WithQueryParamFilter("marker", args.Marker).
		WithQueryParamFilter("maxKeys", strconv.Itoa(args.MaxKeys)).
		WithResult(result).
		Do()

	return result, err
}

// GetBackupDetail - get backup detail of the instance's backup
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - instanceId: id of the instance
//     - backupId: id of the backup
// RETURNS:
//     - *Snapshot: result of the backup detail
//     - error: nil if success otherwise the specific error
func (c *Client) GetBackupDetail(instanceId string, backupId string) (*Snapshot, error) {
	result := &Snapshot{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId) + "/backup/" + backupId).
		WithResult(result).
		Do()

	return result, err
}

// GetZoneList - list all zone
//
// PARAMS:
//     - cli: the client agent which can perform sending request
// RETURNS:
//     - *GetZoneListResult: result of the zone list
//     - error: nil if success otherwise the specific error
func (c *Client) GetZoneList() (*GetZoneListResult, error) {
	result := &GetZoneListResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(URI_PREFIX + "/zone").
		WithResult(result).
		Do()

	return result, err
}

// ListsSubnet - list all Subnets
//
// PARAMS:
//     - cli: the client agent which can perform sending request
//     - args: the arguments to list all subnets, not necessary
// RETURNS:
//     - *ListSubnetsResult: result of the subnet list
//     - error: nil if success otherwise the specific error
func (c *Client) ListSubnets(args *ListSubnetsArgs) (*ListSubnetsResult, error) {
	if args == nil {
		args = &ListSubnetsArgs{}
	}

	result := &ListSubnetsResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(URI_PREFIX+"/subnet").
		WithQueryParamFilter("vpcId", args.VpcId).
		WithQueryParamFilter("zoneName", args.ZoneName).
		WithResult(result).
		Do()

	return result, err
}

// GetSecurityIps - get all SecurityIps
//
// PARAMS:
//     - instanceId: the specific rds Instance's ID
// RETURNS:
//     - *GetSecurityIpsResult: all security IP
//     - error: nil if success otherwise the specific error
func (c *Client) GetSecurityIps(instanceId string) (*GetSecurityIpsResult, error) {
	result := &GetSecurityIpsResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId) + "/securityIp").
		WithResult(result).
		Do()

	return result, err
}

// UpdateSecurityIps - update SecurityIps
//
// PARAMS:
//     - instanceId: the specific rds Instance's ID
//     - Etag: get latest etag by GetSecurityIps
//     - Args: all SecurityIps
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) UpdateSecurityIps(instanceId, Etag string, args *UpdateSecurityIpsArgs) error {

	headers := map[string]string{"x-bce-if-match": Etag}

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)+"/securityIp").
		WithHeaders(headers).
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// ListParameters - list all parameters of a RDS instance
//
// PARAMS:
//     - instanceId: the specific rds Instance's ID
// RETURNS:
//     - *ListParametersResult: the result of list all parameters
//     - error: nil if success otherwise the specific error
func (c *Client) ListParameters(instanceId string) (*ListParametersResult, error) {
	result := &ListParametersResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId) + "/parameter").
		WithResult(result).
		Do()

	return result, err
}

// UpdateParameter - update Parameter
//
// PARAMS:
//     - instanceId: the specific rds Instance's ID
//     - Etag: get latest etag by ListParameters
//     - Args: *UpdateParameterArgs
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) UpdateParameter(instanceId, Etag string, args *UpdateParameterArgs) error {

	headers := map[string]string{"x-bce-if-match": Etag}

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUriWithInstanceId(instanceId)+"/parameter").
		WithHeaders(headers).
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// autoRenew - create autoRenew
//
// PARAMS:
//     - Args: *autoRenewArgs
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) AutoRenew(args *AutoRenewArgs) error {

	return bce.NewRequestBuilder(c).
		WithMethod(http.PUT).
		WithURL(getRdsUri()).
		WithQueryParam("autoRenew", "").
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE).
		WithBody(args).
		Do()
}

// getSlowLogDownloadTaskList
//
// PARAMS:
//     - instanceId: the specific rds Instance's ID
//     - datetime: the log time. range(datetime, datetime + 24 hours)
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) GetSlowLogDownloadTaskList(instanceId, datetime string) (*SlowLogDownloadTaskListResult, error) {
	fmt.Println(getRdsUriWithInstanceId(instanceId) + "/slowlogs/logList/" + datetime)
	result := &SlowLogDownloadTaskListResult{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId) + "/slowlogs/logList/" + datetime).
		WithResult(result).
		Do()
	fmt.Println(result, err)
	return result, err
}

// getSlowLogDownloadDetail
//
// PARAMS:
//     - Args: *slowLogDownloadTaskListArgs
// RETURNS:
//     - error: nil if success otherwise the specific error
func (c *Client) GetSlowLogDownloadDetail(instanceId, logId, downloadValidTimeInSec string) (*SlowLogDownloadDetail, error) {
	result := &SlowLogDownloadDetail{}
	err := bce.NewRequestBuilder(c).
		WithMethod(http.GET).
		WithURL(getRdsUriWithInstanceId(instanceId) + "/slowlogs/download_url/" + logId + "/" + downloadValidTimeInSec).
		WithResult(result).
		Do()
	return result, err
}

func (c *Client) Request(method, uri string, body interface{}) (interface{}, error) {
	res := struct{}{}
	req := bce.NewRequestBuilder(c).
		WithMethod(method).
		WithURL(uri).
		WithHeader(http.CONTENT_TYPE, bce.DEFAULT_CONTENT_TYPE)
	var err error
	if body != nil {
		err = req.
			WithBody(body).
			WithResult(res).
			Do()
	} else {
		err = req.
			WithResult(res).
			Do()
	}

	return res, err
}
