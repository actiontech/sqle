// This file is auto-generated, don't edit it. Thanks.
/**
 *
 */
package client

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	endpointutil "github.com/alibabacloud-go/endpoint-util/service"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

type CancelScheduleTasksRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	TaskId               *string `json:"TaskId,omitempty" xml:"TaskId,omitempty"`
}

func (s CancelScheduleTasksRequest) String() string {
	return tea.Prettify(s)
}

func (s CancelScheduleTasksRequest) GoString() string {
	return s.String()
}

func (s *CancelScheduleTasksRequest) SetDBClusterId(v string) *CancelScheduleTasksRequest {
	s.DBClusterId = &v
	return s
}

func (s *CancelScheduleTasksRequest) SetOwnerAccount(v string) *CancelScheduleTasksRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CancelScheduleTasksRequest) SetOwnerId(v int64) *CancelScheduleTasksRequest {
	s.OwnerId = &v
	return s
}

func (s *CancelScheduleTasksRequest) SetResourceOwnerAccount(v string) *CancelScheduleTasksRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CancelScheduleTasksRequest) SetResourceOwnerId(v int64) *CancelScheduleTasksRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *CancelScheduleTasksRequest) SetTaskId(v string) *CancelScheduleTasksRequest {
	s.TaskId = &v
	return s
}

type CancelScheduleTasksResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	Success   *bool   `json:"Success,omitempty" xml:"Success,omitempty"`
}

func (s CancelScheduleTasksResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CancelScheduleTasksResponseBody) GoString() string {
	return s.String()
}

func (s *CancelScheduleTasksResponseBody) SetRequestId(v string) *CancelScheduleTasksResponseBody {
	s.RequestId = &v
	return s
}

func (s *CancelScheduleTasksResponseBody) SetSuccess(v bool) *CancelScheduleTasksResponseBody {
	s.Success = &v
	return s
}

type CancelScheduleTasksResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CancelScheduleTasksResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CancelScheduleTasksResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelScheduleTasksResponse) GoString() string {
	return s.String()
}

func (s *CancelScheduleTasksResponse) SetHeaders(v map[string]*string) *CancelScheduleTasksResponse {
	s.Headers = v
	return s
}

func (s *CancelScheduleTasksResponse) SetStatusCode(v int32) *CancelScheduleTasksResponse {
	s.StatusCode = &v
	return s
}

func (s *CancelScheduleTasksResponse) SetBody(v *CancelScheduleTasksResponseBody) *CancelScheduleTasksResponse {
	s.Body = v
	return s
}

type CheckAccountNameRequest struct {
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CheckAccountNameRequest) String() string {
	return tea.Prettify(s)
}

func (s CheckAccountNameRequest) GoString() string {
	return s.String()
}

func (s *CheckAccountNameRequest) SetAccountName(v string) *CheckAccountNameRequest {
	s.AccountName = &v
	return s
}

func (s *CheckAccountNameRequest) SetDBClusterId(v string) *CheckAccountNameRequest {
	s.DBClusterId = &v
	return s
}

func (s *CheckAccountNameRequest) SetOwnerAccount(v string) *CheckAccountNameRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CheckAccountNameRequest) SetOwnerId(v int64) *CheckAccountNameRequest {
	s.OwnerId = &v
	return s
}

func (s *CheckAccountNameRequest) SetResourceOwnerAccount(v string) *CheckAccountNameRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CheckAccountNameRequest) SetResourceOwnerId(v int64) *CheckAccountNameRequest {
	s.ResourceOwnerId = &v
	return s
}

type CheckAccountNameResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CheckAccountNameResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CheckAccountNameResponseBody) GoString() string {
	return s.String()
}

func (s *CheckAccountNameResponseBody) SetRequestId(v string) *CheckAccountNameResponseBody {
	s.RequestId = &v
	return s
}

type CheckAccountNameResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CheckAccountNameResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CheckAccountNameResponse) String() string {
	return tea.Prettify(s)
}

func (s CheckAccountNameResponse) GoString() string {
	return s.String()
}

func (s *CheckAccountNameResponse) SetHeaders(v map[string]*string) *CheckAccountNameResponse {
	s.Headers = v
	return s
}

func (s *CheckAccountNameResponse) SetStatusCode(v int32) *CheckAccountNameResponse {
	s.StatusCode = &v
	return s
}

func (s *CheckAccountNameResponse) SetBody(v *CheckAccountNameResponseBody) *CheckAccountNameResponse {
	s.Body = v
	return s
}

type CheckDBNameRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CheckDBNameRequest) String() string {
	return tea.Prettify(s)
}

func (s CheckDBNameRequest) GoString() string {
	return s.String()
}

func (s *CheckDBNameRequest) SetDBClusterId(v string) *CheckDBNameRequest {
	s.DBClusterId = &v
	return s
}

func (s *CheckDBNameRequest) SetDBName(v string) *CheckDBNameRequest {
	s.DBName = &v
	return s
}

func (s *CheckDBNameRequest) SetOwnerAccount(v string) *CheckDBNameRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CheckDBNameRequest) SetOwnerId(v int64) *CheckDBNameRequest {
	s.OwnerId = &v
	return s
}

func (s *CheckDBNameRequest) SetResourceOwnerAccount(v string) *CheckDBNameRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CheckDBNameRequest) SetResourceOwnerId(v int64) *CheckDBNameRequest {
	s.ResourceOwnerId = &v
	return s
}

type CheckDBNameResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CheckDBNameResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CheckDBNameResponseBody) GoString() string {
	return s.String()
}

func (s *CheckDBNameResponseBody) SetRequestId(v string) *CheckDBNameResponseBody {
	s.RequestId = &v
	return s
}

type CheckDBNameResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CheckDBNameResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CheckDBNameResponse) String() string {
	return tea.Prettify(s)
}

func (s CheckDBNameResponse) GoString() string {
	return s.String()
}

func (s *CheckDBNameResponse) SetHeaders(v map[string]*string) *CheckDBNameResponse {
	s.Headers = v
	return s
}

func (s *CheckDBNameResponse) SetStatusCode(v int32) *CheckDBNameResponse {
	s.StatusCode = &v
	return s
}

func (s *CheckDBNameResponse) SetBody(v *CheckDBNameResponseBody) *CheckDBNameResponse {
	s.Body = v
	return s
}

type CloseAITaskRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CloseAITaskRequest) String() string {
	return tea.Prettify(s)
}

func (s CloseAITaskRequest) GoString() string {
	return s.String()
}

func (s *CloseAITaskRequest) SetDBClusterId(v string) *CloseAITaskRequest {
	s.DBClusterId = &v
	return s
}

func (s *CloseAITaskRequest) SetOwnerAccount(v string) *CloseAITaskRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CloseAITaskRequest) SetOwnerId(v int64) *CloseAITaskRequest {
	s.OwnerId = &v
	return s
}

func (s *CloseAITaskRequest) SetRegionId(v string) *CloseAITaskRequest {
	s.RegionId = &v
	return s
}

func (s *CloseAITaskRequest) SetResourceOwnerAccount(v string) *CloseAITaskRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CloseAITaskRequest) SetResourceOwnerId(v int64) *CloseAITaskRequest {
	s.ResourceOwnerId = &v
	return s
}

type CloseAITaskResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TaskId    *string `json:"TaskId,omitempty" xml:"TaskId,omitempty"`
}

func (s CloseAITaskResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CloseAITaskResponseBody) GoString() string {
	return s.String()
}

func (s *CloseAITaskResponseBody) SetRequestId(v string) *CloseAITaskResponseBody {
	s.RequestId = &v
	return s
}

func (s *CloseAITaskResponseBody) SetTaskId(v string) *CloseAITaskResponseBody {
	s.TaskId = &v
	return s
}

type CloseAITaskResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CloseAITaskResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CloseAITaskResponse) String() string {
	return tea.Prettify(s)
}

func (s CloseAITaskResponse) GoString() string {
	return s.String()
}

func (s *CloseAITaskResponse) SetHeaders(v map[string]*string) *CloseAITaskResponse {
	s.Headers = v
	return s
}

func (s *CloseAITaskResponse) SetStatusCode(v int32) *CloseAITaskResponse {
	s.StatusCode = &v
	return s
}

func (s *CloseAITaskResponse) SetBody(v *CloseAITaskResponseBody) *CloseAITaskResponse {
	s.Body = v
	return s
}

type CloseDBClusterMigrationRequest struct {
	ContinueEnableBinlog *bool   `json:"ContinueEnableBinlog,omitempty" xml:"ContinueEnableBinlog,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CloseDBClusterMigrationRequest) String() string {
	return tea.Prettify(s)
}

func (s CloseDBClusterMigrationRequest) GoString() string {
	return s.String()
}

func (s *CloseDBClusterMigrationRequest) SetContinueEnableBinlog(v bool) *CloseDBClusterMigrationRequest {
	s.ContinueEnableBinlog = &v
	return s
}

func (s *CloseDBClusterMigrationRequest) SetDBClusterId(v string) *CloseDBClusterMigrationRequest {
	s.DBClusterId = &v
	return s
}

func (s *CloseDBClusterMigrationRequest) SetOwnerAccount(v string) *CloseDBClusterMigrationRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CloseDBClusterMigrationRequest) SetOwnerId(v int64) *CloseDBClusterMigrationRequest {
	s.OwnerId = &v
	return s
}

func (s *CloseDBClusterMigrationRequest) SetResourceOwnerAccount(v string) *CloseDBClusterMigrationRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CloseDBClusterMigrationRequest) SetResourceOwnerId(v int64) *CloseDBClusterMigrationRequest {
	s.ResourceOwnerId = &v
	return s
}

type CloseDBClusterMigrationResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CloseDBClusterMigrationResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CloseDBClusterMigrationResponseBody) GoString() string {
	return s.String()
}

func (s *CloseDBClusterMigrationResponseBody) SetRequestId(v string) *CloseDBClusterMigrationResponseBody {
	s.RequestId = &v
	return s
}

type CloseDBClusterMigrationResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CloseDBClusterMigrationResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CloseDBClusterMigrationResponse) String() string {
	return tea.Prettify(s)
}

func (s CloseDBClusterMigrationResponse) GoString() string {
	return s.String()
}

func (s *CloseDBClusterMigrationResponse) SetHeaders(v map[string]*string) *CloseDBClusterMigrationResponse {
	s.Headers = v
	return s
}

func (s *CloseDBClusterMigrationResponse) SetStatusCode(v int32) *CloseDBClusterMigrationResponse {
	s.StatusCode = &v
	return s
}

func (s *CloseDBClusterMigrationResponse) SetBody(v *CloseDBClusterMigrationResponseBody) *CloseDBClusterMigrationResponse {
	s.Body = v
	return s
}

type CreateAccountRequest struct {
	AccountDescription   *string `json:"AccountDescription,omitempty" xml:"AccountDescription,omitempty"`
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	AccountPassword      *string `json:"AccountPassword,omitempty" xml:"AccountPassword,omitempty"`
	AccountPrivilege     *string `json:"AccountPrivilege,omitempty" xml:"AccountPrivilege,omitempty"`
	AccountType          *string `json:"AccountType,omitempty" xml:"AccountType,omitempty"`
	ClientToken          *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CreateAccountRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAccountRequest) GoString() string {
	return s.String()
}

func (s *CreateAccountRequest) SetAccountDescription(v string) *CreateAccountRequest {
	s.AccountDescription = &v
	return s
}

func (s *CreateAccountRequest) SetAccountName(v string) *CreateAccountRequest {
	s.AccountName = &v
	return s
}

func (s *CreateAccountRequest) SetAccountPassword(v string) *CreateAccountRequest {
	s.AccountPassword = &v
	return s
}

func (s *CreateAccountRequest) SetAccountPrivilege(v string) *CreateAccountRequest {
	s.AccountPrivilege = &v
	return s
}

func (s *CreateAccountRequest) SetAccountType(v string) *CreateAccountRequest {
	s.AccountType = &v
	return s
}

func (s *CreateAccountRequest) SetClientToken(v string) *CreateAccountRequest {
	s.ClientToken = &v
	return s
}

func (s *CreateAccountRequest) SetDBClusterId(v string) *CreateAccountRequest {
	s.DBClusterId = &v
	return s
}

func (s *CreateAccountRequest) SetDBName(v string) *CreateAccountRequest {
	s.DBName = &v
	return s
}

func (s *CreateAccountRequest) SetOwnerAccount(v string) *CreateAccountRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateAccountRequest) SetOwnerId(v int64) *CreateAccountRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateAccountRequest) SetResourceOwnerAccount(v string) *CreateAccountRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateAccountRequest) SetResourceOwnerId(v int64) *CreateAccountRequest {
	s.ResourceOwnerId = &v
	return s
}

type CreateAccountResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateAccountResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateAccountResponseBody) GoString() string {
	return s.String()
}

func (s *CreateAccountResponseBody) SetRequestId(v string) *CreateAccountResponseBody {
	s.RequestId = &v
	return s
}

type CreateAccountResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateAccountResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateAccountResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAccountResponse) GoString() string {
	return s.String()
}

func (s *CreateAccountResponse) SetHeaders(v map[string]*string) *CreateAccountResponse {
	s.Headers = v
	return s
}

func (s *CreateAccountResponse) SetStatusCode(v int32) *CreateAccountResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateAccountResponse) SetBody(v *CreateAccountResponseBody) *CreateAccountResponse {
	s.Body = v
	return s
}

type CreateBackupRequest struct {
	ClientToken          *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CreateBackupRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateBackupRequest) GoString() string {
	return s.String()
}

func (s *CreateBackupRequest) SetClientToken(v string) *CreateBackupRequest {
	s.ClientToken = &v
	return s
}

func (s *CreateBackupRequest) SetDBClusterId(v string) *CreateBackupRequest {
	s.DBClusterId = &v
	return s
}

func (s *CreateBackupRequest) SetOwnerAccount(v string) *CreateBackupRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateBackupRequest) SetOwnerId(v int64) *CreateBackupRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateBackupRequest) SetResourceOwnerAccount(v string) *CreateBackupRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateBackupRequest) SetResourceOwnerId(v int64) *CreateBackupRequest {
	s.ResourceOwnerId = &v
	return s
}

type CreateBackupResponseBody struct {
	BackupJobId *string `json:"BackupJobId,omitempty" xml:"BackupJobId,omitempty"`
	RequestId   *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateBackupResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateBackupResponseBody) GoString() string {
	return s.String()
}

func (s *CreateBackupResponseBody) SetBackupJobId(v string) *CreateBackupResponseBody {
	s.BackupJobId = &v
	return s
}

func (s *CreateBackupResponseBody) SetRequestId(v string) *CreateBackupResponseBody {
	s.RequestId = &v
	return s
}

type CreateBackupResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateBackupResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateBackupResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateBackupResponse) GoString() string {
	return s.String()
}

func (s *CreateBackupResponse) SetHeaders(v map[string]*string) *CreateBackupResponse {
	s.Headers = v
	return s
}

func (s *CreateBackupResponse) SetStatusCode(v int32) *CreateBackupResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateBackupResponse) SetBody(v *CreateBackupResponseBody) *CreateBackupResponse {
	s.Body = v
	return s
}

type CreateDBClusterRequest struct {
	AutoRenew                              *bool   `json:"AutoRenew,omitempty" xml:"AutoRenew,omitempty"`
	BackupRetentionPolicyOnClusterDeletion *string `json:"BackupRetentionPolicyOnClusterDeletion,omitempty" xml:"BackupRetentionPolicyOnClusterDeletion,omitempty"`
	ClientToken                            *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	CloneDataPoint                         *string `json:"CloneDataPoint,omitempty" xml:"CloneDataPoint,omitempty"`
	ClusterNetworkType                     *string `json:"ClusterNetworkType,omitempty" xml:"ClusterNetworkType,omitempty"`
	CreationCategory                       *string `json:"CreationCategory,omitempty" xml:"CreationCategory,omitempty"`
	CreationOption                         *string `json:"CreationOption,omitempty" xml:"CreationOption,omitempty"`
	DBClusterDescription                   *string `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBMinorVersion                         *string `json:"DBMinorVersion,omitempty" xml:"DBMinorVersion,omitempty"`
	DBNodeClass                            *string `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
	DBType                                 *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion                              *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	DefaultTimeZone                        *string `json:"DefaultTimeZone,omitempty" xml:"DefaultTimeZone,omitempty"`
	GDNId                                  *string `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	LowerCaseTableNames                    *string `json:"LowerCaseTableNames,omitempty" xml:"LowerCaseTableNames,omitempty"`
	OwnerAccount                           *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId                                *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ParameterGroupId                       *string `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	PayType                                *string `json:"PayType,omitempty" xml:"PayType,omitempty"`
	Period                                 *string `json:"Period,omitempty" xml:"Period,omitempty"`
	RegionId                               *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceGroupId                        *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	ResourceOwnerAccount                   *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId                        *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityIPList                         *string `json:"SecurityIPList,omitempty" xml:"SecurityIPList,omitempty"`
	SourceResourceId                       *string `json:"SourceResourceId,omitempty" xml:"SourceResourceId,omitempty"`
	TDEStatus                              *bool   `json:"TDEStatus,omitempty" xml:"TDEStatus,omitempty"`
	UsedTime                               *string `json:"UsedTime,omitempty" xml:"UsedTime,omitempty"`
	VPCId                                  *string `json:"VPCId,omitempty" xml:"VPCId,omitempty"`
	VSwitchId                              *string `json:"VSwitchId,omitempty" xml:"VSwitchId,omitempty"`
	ZoneId                                 *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s CreateDBClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateDBClusterRequest) GoString() string {
	return s.String()
}

func (s *CreateDBClusterRequest) SetAutoRenew(v bool) *CreateDBClusterRequest {
	s.AutoRenew = &v
	return s
}

func (s *CreateDBClusterRequest) SetBackupRetentionPolicyOnClusterDeletion(v string) *CreateDBClusterRequest {
	s.BackupRetentionPolicyOnClusterDeletion = &v
	return s
}

func (s *CreateDBClusterRequest) SetClientToken(v string) *CreateDBClusterRequest {
	s.ClientToken = &v
	return s
}

func (s *CreateDBClusterRequest) SetCloneDataPoint(v string) *CreateDBClusterRequest {
	s.CloneDataPoint = &v
	return s
}

func (s *CreateDBClusterRequest) SetClusterNetworkType(v string) *CreateDBClusterRequest {
	s.ClusterNetworkType = &v
	return s
}

func (s *CreateDBClusterRequest) SetCreationCategory(v string) *CreateDBClusterRequest {
	s.CreationCategory = &v
	return s
}

func (s *CreateDBClusterRequest) SetCreationOption(v string) *CreateDBClusterRequest {
	s.CreationOption = &v
	return s
}

func (s *CreateDBClusterRequest) SetDBClusterDescription(v string) *CreateDBClusterRequest {
	s.DBClusterDescription = &v
	return s
}

func (s *CreateDBClusterRequest) SetDBMinorVersion(v string) *CreateDBClusterRequest {
	s.DBMinorVersion = &v
	return s
}

func (s *CreateDBClusterRequest) SetDBNodeClass(v string) *CreateDBClusterRequest {
	s.DBNodeClass = &v
	return s
}

func (s *CreateDBClusterRequest) SetDBType(v string) *CreateDBClusterRequest {
	s.DBType = &v
	return s
}

func (s *CreateDBClusterRequest) SetDBVersion(v string) *CreateDBClusterRequest {
	s.DBVersion = &v
	return s
}

func (s *CreateDBClusterRequest) SetDefaultTimeZone(v string) *CreateDBClusterRequest {
	s.DefaultTimeZone = &v
	return s
}

func (s *CreateDBClusterRequest) SetGDNId(v string) *CreateDBClusterRequest {
	s.GDNId = &v
	return s
}

func (s *CreateDBClusterRequest) SetLowerCaseTableNames(v string) *CreateDBClusterRequest {
	s.LowerCaseTableNames = &v
	return s
}

func (s *CreateDBClusterRequest) SetOwnerAccount(v string) *CreateDBClusterRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateDBClusterRequest) SetOwnerId(v int64) *CreateDBClusterRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateDBClusterRequest) SetParameterGroupId(v string) *CreateDBClusterRequest {
	s.ParameterGroupId = &v
	return s
}

func (s *CreateDBClusterRequest) SetPayType(v string) *CreateDBClusterRequest {
	s.PayType = &v
	return s
}

func (s *CreateDBClusterRequest) SetPeriod(v string) *CreateDBClusterRequest {
	s.Period = &v
	return s
}

func (s *CreateDBClusterRequest) SetRegionId(v string) *CreateDBClusterRequest {
	s.RegionId = &v
	return s
}

func (s *CreateDBClusterRequest) SetResourceGroupId(v string) *CreateDBClusterRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *CreateDBClusterRequest) SetResourceOwnerAccount(v string) *CreateDBClusterRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateDBClusterRequest) SetResourceOwnerId(v int64) *CreateDBClusterRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *CreateDBClusterRequest) SetSecurityIPList(v string) *CreateDBClusterRequest {
	s.SecurityIPList = &v
	return s
}

func (s *CreateDBClusterRequest) SetSourceResourceId(v string) *CreateDBClusterRequest {
	s.SourceResourceId = &v
	return s
}

func (s *CreateDBClusterRequest) SetTDEStatus(v bool) *CreateDBClusterRequest {
	s.TDEStatus = &v
	return s
}

func (s *CreateDBClusterRequest) SetUsedTime(v string) *CreateDBClusterRequest {
	s.UsedTime = &v
	return s
}

func (s *CreateDBClusterRequest) SetVPCId(v string) *CreateDBClusterRequest {
	s.VPCId = &v
	return s
}

func (s *CreateDBClusterRequest) SetVSwitchId(v string) *CreateDBClusterRequest {
	s.VSwitchId = &v
	return s
}

func (s *CreateDBClusterRequest) SetZoneId(v string) *CreateDBClusterRequest {
	s.ZoneId = &v
	return s
}

type CreateDBClusterResponseBody struct {
	DBClusterId     *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OrderId         *string `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	RequestId       *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	ResourceGroupId *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
}

func (s CreateDBClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateDBClusterResponseBody) GoString() string {
	return s.String()
}

func (s *CreateDBClusterResponseBody) SetDBClusterId(v string) *CreateDBClusterResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *CreateDBClusterResponseBody) SetOrderId(v string) *CreateDBClusterResponseBody {
	s.OrderId = &v
	return s
}

func (s *CreateDBClusterResponseBody) SetRequestId(v string) *CreateDBClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *CreateDBClusterResponseBody) SetResourceGroupId(v string) *CreateDBClusterResponseBody {
	s.ResourceGroupId = &v
	return s
}

type CreateDBClusterResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateDBClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateDBClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateDBClusterResponse) GoString() string {
	return s.String()
}

func (s *CreateDBClusterResponse) SetHeaders(v map[string]*string) *CreateDBClusterResponse {
	s.Headers = v
	return s
}

func (s *CreateDBClusterResponse) SetStatusCode(v int32) *CreateDBClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateDBClusterResponse) SetBody(v *CreateDBClusterResponseBody) *CreateDBClusterResponse {
	s.Body = v
	return s
}

type CreateDBClusterEndpointRequest struct {
	AutoAddNewNodes       *string `json:"AutoAddNewNodes,omitempty" xml:"AutoAddNewNodes,omitempty"`
	ClientToken           *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId           *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBEndpointDescription *string `json:"DBEndpointDescription,omitempty" xml:"DBEndpointDescription,omitempty"`
	EndpointConfig        *string `json:"EndpointConfig,omitempty" xml:"EndpointConfig,omitempty"`
	EndpointType          *string `json:"EndpointType,omitempty" xml:"EndpointType,omitempty"`
	Nodes                 *string `json:"Nodes,omitempty" xml:"Nodes,omitempty"`
	OwnerAccount          *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId               *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ReadWriteMode         *string `json:"ReadWriteMode,omitempty" xml:"ReadWriteMode,omitempty"`
	ResourceOwnerAccount  *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId       *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CreateDBClusterEndpointRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateDBClusterEndpointRequest) GoString() string {
	return s.String()
}

func (s *CreateDBClusterEndpointRequest) SetAutoAddNewNodes(v string) *CreateDBClusterEndpointRequest {
	s.AutoAddNewNodes = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetClientToken(v string) *CreateDBClusterEndpointRequest {
	s.ClientToken = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetDBClusterId(v string) *CreateDBClusterEndpointRequest {
	s.DBClusterId = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetDBEndpointDescription(v string) *CreateDBClusterEndpointRequest {
	s.DBEndpointDescription = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetEndpointConfig(v string) *CreateDBClusterEndpointRequest {
	s.EndpointConfig = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetEndpointType(v string) *CreateDBClusterEndpointRequest {
	s.EndpointType = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetNodes(v string) *CreateDBClusterEndpointRequest {
	s.Nodes = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetOwnerAccount(v string) *CreateDBClusterEndpointRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetOwnerId(v int64) *CreateDBClusterEndpointRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetReadWriteMode(v string) *CreateDBClusterEndpointRequest {
	s.ReadWriteMode = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetResourceOwnerAccount(v string) *CreateDBClusterEndpointRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateDBClusterEndpointRequest) SetResourceOwnerId(v int64) *CreateDBClusterEndpointRequest {
	s.ResourceOwnerId = &v
	return s
}

type CreateDBClusterEndpointResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateDBClusterEndpointResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateDBClusterEndpointResponseBody) GoString() string {
	return s.String()
}

func (s *CreateDBClusterEndpointResponseBody) SetRequestId(v string) *CreateDBClusterEndpointResponseBody {
	s.RequestId = &v
	return s
}

type CreateDBClusterEndpointResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateDBClusterEndpointResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateDBClusterEndpointResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateDBClusterEndpointResponse) GoString() string {
	return s.String()
}

func (s *CreateDBClusterEndpointResponse) SetHeaders(v map[string]*string) *CreateDBClusterEndpointResponse {
	s.Headers = v
	return s
}

func (s *CreateDBClusterEndpointResponse) SetStatusCode(v int32) *CreateDBClusterEndpointResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateDBClusterEndpointResponse) SetBody(v *CreateDBClusterEndpointResponseBody) *CreateDBClusterEndpointResponse {
	s.Body = v
	return s
}

type CreateDBEndpointAddressRequest struct {
	ConnectionStringPrefix *string `json:"ConnectionStringPrefix,omitempty" xml:"ConnectionStringPrefix,omitempty"`
	DBClusterId            *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBEndpointId           *string `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	NetType                *string `json:"NetType,omitempty" xml:"NetType,omitempty"`
	OwnerAccount           *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId                *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount   *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId        *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CreateDBEndpointAddressRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateDBEndpointAddressRequest) GoString() string {
	return s.String()
}

func (s *CreateDBEndpointAddressRequest) SetConnectionStringPrefix(v string) *CreateDBEndpointAddressRequest {
	s.ConnectionStringPrefix = &v
	return s
}

func (s *CreateDBEndpointAddressRequest) SetDBClusterId(v string) *CreateDBEndpointAddressRequest {
	s.DBClusterId = &v
	return s
}

func (s *CreateDBEndpointAddressRequest) SetDBEndpointId(v string) *CreateDBEndpointAddressRequest {
	s.DBEndpointId = &v
	return s
}

func (s *CreateDBEndpointAddressRequest) SetNetType(v string) *CreateDBEndpointAddressRequest {
	s.NetType = &v
	return s
}

func (s *CreateDBEndpointAddressRequest) SetOwnerAccount(v string) *CreateDBEndpointAddressRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateDBEndpointAddressRequest) SetOwnerId(v int64) *CreateDBEndpointAddressRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateDBEndpointAddressRequest) SetResourceOwnerAccount(v string) *CreateDBEndpointAddressRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateDBEndpointAddressRequest) SetResourceOwnerId(v int64) *CreateDBEndpointAddressRequest {
	s.ResourceOwnerId = &v
	return s
}

type CreateDBEndpointAddressResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateDBEndpointAddressResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateDBEndpointAddressResponseBody) GoString() string {
	return s.String()
}

func (s *CreateDBEndpointAddressResponseBody) SetRequestId(v string) *CreateDBEndpointAddressResponseBody {
	s.RequestId = &v
	return s
}

type CreateDBEndpointAddressResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateDBEndpointAddressResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateDBEndpointAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateDBEndpointAddressResponse) GoString() string {
	return s.String()
}

func (s *CreateDBEndpointAddressResponse) SetHeaders(v map[string]*string) *CreateDBEndpointAddressResponse {
	s.Headers = v
	return s
}

func (s *CreateDBEndpointAddressResponse) SetStatusCode(v int32) *CreateDBEndpointAddressResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateDBEndpointAddressResponse) SetBody(v *CreateDBEndpointAddressResponseBody) *CreateDBEndpointAddressResponse {
	s.Body = v
	return s
}

type CreateDBLinkRequest struct {
	ClientToken          *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBLinkName           *string `json:"DBLinkName,omitempty" xml:"DBLinkName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SourceDBName         *string `json:"SourceDBName,omitempty" xml:"SourceDBName,omitempty"`
	TargetDBAccount      *string `json:"TargetDBAccount,omitempty" xml:"TargetDBAccount,omitempty"`
	TargetDBInstanceName *string `json:"TargetDBInstanceName,omitempty" xml:"TargetDBInstanceName,omitempty"`
	TargetDBName         *string `json:"TargetDBName,omitempty" xml:"TargetDBName,omitempty"`
	TargetDBPasswd       *string `json:"TargetDBPasswd,omitempty" xml:"TargetDBPasswd,omitempty"`
	TargetIp             *string `json:"TargetIp,omitempty" xml:"TargetIp,omitempty"`
	TargetPort           *string `json:"TargetPort,omitempty" xml:"TargetPort,omitempty"`
	VpcId                *string `json:"VpcId,omitempty" xml:"VpcId,omitempty"`
}

func (s CreateDBLinkRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateDBLinkRequest) GoString() string {
	return s.String()
}

func (s *CreateDBLinkRequest) SetClientToken(v string) *CreateDBLinkRequest {
	s.ClientToken = &v
	return s
}

func (s *CreateDBLinkRequest) SetDBClusterId(v string) *CreateDBLinkRequest {
	s.DBClusterId = &v
	return s
}

func (s *CreateDBLinkRequest) SetDBLinkName(v string) *CreateDBLinkRequest {
	s.DBLinkName = &v
	return s
}

func (s *CreateDBLinkRequest) SetOwnerAccount(v string) *CreateDBLinkRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateDBLinkRequest) SetOwnerId(v int64) *CreateDBLinkRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateDBLinkRequest) SetRegionId(v string) *CreateDBLinkRequest {
	s.RegionId = &v
	return s
}

func (s *CreateDBLinkRequest) SetResourceOwnerAccount(v string) *CreateDBLinkRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateDBLinkRequest) SetResourceOwnerId(v int64) *CreateDBLinkRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *CreateDBLinkRequest) SetSourceDBName(v string) *CreateDBLinkRequest {
	s.SourceDBName = &v
	return s
}

func (s *CreateDBLinkRequest) SetTargetDBAccount(v string) *CreateDBLinkRequest {
	s.TargetDBAccount = &v
	return s
}

func (s *CreateDBLinkRequest) SetTargetDBInstanceName(v string) *CreateDBLinkRequest {
	s.TargetDBInstanceName = &v
	return s
}

func (s *CreateDBLinkRequest) SetTargetDBName(v string) *CreateDBLinkRequest {
	s.TargetDBName = &v
	return s
}

func (s *CreateDBLinkRequest) SetTargetDBPasswd(v string) *CreateDBLinkRequest {
	s.TargetDBPasswd = &v
	return s
}

func (s *CreateDBLinkRequest) SetTargetIp(v string) *CreateDBLinkRequest {
	s.TargetIp = &v
	return s
}

func (s *CreateDBLinkRequest) SetTargetPort(v string) *CreateDBLinkRequest {
	s.TargetPort = &v
	return s
}

func (s *CreateDBLinkRequest) SetVpcId(v string) *CreateDBLinkRequest {
	s.VpcId = &v
	return s
}

type CreateDBLinkResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateDBLinkResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateDBLinkResponseBody) GoString() string {
	return s.String()
}

func (s *CreateDBLinkResponseBody) SetRequestId(v string) *CreateDBLinkResponseBody {
	s.RequestId = &v
	return s
}

type CreateDBLinkResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateDBLinkResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateDBLinkResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateDBLinkResponse) GoString() string {
	return s.String()
}

func (s *CreateDBLinkResponse) SetHeaders(v map[string]*string) *CreateDBLinkResponse {
	s.Headers = v
	return s
}

func (s *CreateDBLinkResponse) SetStatusCode(v int32) *CreateDBLinkResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateDBLinkResponse) SetBody(v *CreateDBLinkResponseBody) *CreateDBLinkResponse {
	s.Body = v
	return s
}

type CreateDBNodesRequest struct {
	ClientToken          *string                       `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string                       `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNode               []*CreateDBNodesRequestDBNode `json:"DBNode,omitempty" xml:"DBNode,omitempty" type:"Repeated"`
	EndpointBindList     *string                       `json:"EndpointBindList,omitempty" xml:"EndpointBindList,omitempty"`
	ImciSwitch           *string                       `json:"ImciSwitch,omitempty" xml:"ImciSwitch,omitempty"`
	OwnerAccount         *string                       `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64                        `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PlannedEndTime       *string                       `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string                       `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string                       `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64                        `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CreateDBNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateDBNodesRequest) GoString() string {
	return s.String()
}

func (s *CreateDBNodesRequest) SetClientToken(v string) *CreateDBNodesRequest {
	s.ClientToken = &v
	return s
}

func (s *CreateDBNodesRequest) SetDBClusterId(v string) *CreateDBNodesRequest {
	s.DBClusterId = &v
	return s
}

func (s *CreateDBNodesRequest) SetDBNode(v []*CreateDBNodesRequestDBNode) *CreateDBNodesRequest {
	s.DBNode = v
	return s
}

func (s *CreateDBNodesRequest) SetEndpointBindList(v string) *CreateDBNodesRequest {
	s.EndpointBindList = &v
	return s
}

func (s *CreateDBNodesRequest) SetImciSwitch(v string) *CreateDBNodesRequest {
	s.ImciSwitch = &v
	return s
}

func (s *CreateDBNodesRequest) SetOwnerAccount(v string) *CreateDBNodesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateDBNodesRequest) SetOwnerId(v int64) *CreateDBNodesRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateDBNodesRequest) SetPlannedEndTime(v string) *CreateDBNodesRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *CreateDBNodesRequest) SetPlannedStartTime(v string) *CreateDBNodesRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *CreateDBNodesRequest) SetResourceOwnerAccount(v string) *CreateDBNodesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateDBNodesRequest) SetResourceOwnerId(v int64) *CreateDBNodesRequest {
	s.ResourceOwnerId = &v
	return s
}

type CreateDBNodesRequestDBNode struct {
	TargetClass *string `json:"TargetClass,omitempty" xml:"TargetClass,omitempty"`
	ZoneId      *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s CreateDBNodesRequestDBNode) String() string {
	return tea.Prettify(s)
}

func (s CreateDBNodesRequestDBNode) GoString() string {
	return s.String()
}

func (s *CreateDBNodesRequestDBNode) SetTargetClass(v string) *CreateDBNodesRequestDBNode {
	s.TargetClass = &v
	return s
}

func (s *CreateDBNodesRequestDBNode) SetZoneId(v string) *CreateDBNodesRequestDBNode {
	s.ZoneId = &v
	return s
}

type CreateDBNodesResponseBody struct {
	DBClusterId *string                             `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeIds   *CreateDBNodesResponseBodyDBNodeIds `json:"DBNodeIds,omitempty" xml:"DBNodeIds,omitempty" type:"Struct"`
	OrderId     *string                             `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	RequestId   *string                             `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateDBNodesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateDBNodesResponseBody) GoString() string {
	return s.String()
}

func (s *CreateDBNodesResponseBody) SetDBClusterId(v string) *CreateDBNodesResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *CreateDBNodesResponseBody) SetDBNodeIds(v *CreateDBNodesResponseBodyDBNodeIds) *CreateDBNodesResponseBody {
	s.DBNodeIds = v
	return s
}

func (s *CreateDBNodesResponseBody) SetOrderId(v string) *CreateDBNodesResponseBody {
	s.OrderId = &v
	return s
}

func (s *CreateDBNodesResponseBody) SetRequestId(v string) *CreateDBNodesResponseBody {
	s.RequestId = &v
	return s
}

type CreateDBNodesResponseBodyDBNodeIds struct {
	DBNodeId []*string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty" type:"Repeated"`
}

func (s CreateDBNodesResponseBodyDBNodeIds) String() string {
	return tea.Prettify(s)
}

func (s CreateDBNodesResponseBodyDBNodeIds) GoString() string {
	return s.String()
}

func (s *CreateDBNodesResponseBodyDBNodeIds) SetDBNodeId(v []*string) *CreateDBNodesResponseBodyDBNodeIds {
	s.DBNodeId = v
	return s
}

type CreateDBNodesResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateDBNodesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateDBNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateDBNodesResponse) GoString() string {
	return s.String()
}

func (s *CreateDBNodesResponse) SetHeaders(v map[string]*string) *CreateDBNodesResponse {
	s.Headers = v
	return s
}

func (s *CreateDBNodesResponse) SetStatusCode(v int32) *CreateDBNodesResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateDBNodesResponse) SetBody(v *CreateDBNodesResponseBody) *CreateDBNodesResponse {
	s.Body = v
	return s
}

type CreateDatabaseRequest struct {
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	AccountPrivilege     *string `json:"AccountPrivilege,omitempty" xml:"AccountPrivilege,omitempty"`
	CharacterSetName     *string `json:"CharacterSetName,omitempty" xml:"CharacterSetName,omitempty"`
	Collate              *string `json:"Collate,omitempty" xml:"Collate,omitempty"`
	Ctype                *string `json:"Ctype,omitempty" xml:"Ctype,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBDescription        *string `json:"DBDescription,omitempty" xml:"DBDescription,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CreateDatabaseRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateDatabaseRequest) GoString() string {
	return s.String()
}

func (s *CreateDatabaseRequest) SetAccountName(v string) *CreateDatabaseRequest {
	s.AccountName = &v
	return s
}

func (s *CreateDatabaseRequest) SetAccountPrivilege(v string) *CreateDatabaseRequest {
	s.AccountPrivilege = &v
	return s
}

func (s *CreateDatabaseRequest) SetCharacterSetName(v string) *CreateDatabaseRequest {
	s.CharacterSetName = &v
	return s
}

func (s *CreateDatabaseRequest) SetCollate(v string) *CreateDatabaseRequest {
	s.Collate = &v
	return s
}

func (s *CreateDatabaseRequest) SetCtype(v string) *CreateDatabaseRequest {
	s.Ctype = &v
	return s
}

func (s *CreateDatabaseRequest) SetDBClusterId(v string) *CreateDatabaseRequest {
	s.DBClusterId = &v
	return s
}

func (s *CreateDatabaseRequest) SetDBDescription(v string) *CreateDatabaseRequest {
	s.DBDescription = &v
	return s
}

func (s *CreateDatabaseRequest) SetDBName(v string) *CreateDatabaseRequest {
	s.DBName = &v
	return s
}

func (s *CreateDatabaseRequest) SetOwnerAccount(v string) *CreateDatabaseRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateDatabaseRequest) SetOwnerId(v int64) *CreateDatabaseRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateDatabaseRequest) SetResourceOwnerAccount(v string) *CreateDatabaseRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateDatabaseRequest) SetResourceOwnerId(v int64) *CreateDatabaseRequest {
	s.ResourceOwnerId = &v
	return s
}

type CreateDatabaseResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateDatabaseResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateDatabaseResponseBody) GoString() string {
	return s.String()
}

func (s *CreateDatabaseResponseBody) SetRequestId(v string) *CreateDatabaseResponseBody {
	s.RequestId = &v
	return s
}

type CreateDatabaseResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateDatabaseResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateDatabaseResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateDatabaseResponse) GoString() string {
	return s.String()
}

func (s *CreateDatabaseResponse) SetHeaders(v map[string]*string) *CreateDatabaseResponse {
	s.Headers = v
	return s
}

func (s *CreateDatabaseResponse) SetStatusCode(v int32) *CreateDatabaseResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateDatabaseResponse) SetBody(v *CreateDatabaseResponseBody) *CreateDatabaseResponse {
	s.Body = v
	return s
}

type CreateGlobalDatabaseNetworkRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	GDNDescription       *string `json:"GDNDescription,omitempty" xml:"GDNDescription,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s CreateGlobalDatabaseNetworkRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateGlobalDatabaseNetworkRequest) GoString() string {
	return s.String()
}

func (s *CreateGlobalDatabaseNetworkRequest) SetDBClusterId(v string) *CreateGlobalDatabaseNetworkRequest {
	s.DBClusterId = &v
	return s
}

func (s *CreateGlobalDatabaseNetworkRequest) SetGDNDescription(v string) *CreateGlobalDatabaseNetworkRequest {
	s.GDNDescription = &v
	return s
}

func (s *CreateGlobalDatabaseNetworkRequest) SetOwnerAccount(v string) *CreateGlobalDatabaseNetworkRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateGlobalDatabaseNetworkRequest) SetOwnerId(v int64) *CreateGlobalDatabaseNetworkRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateGlobalDatabaseNetworkRequest) SetResourceOwnerAccount(v string) *CreateGlobalDatabaseNetworkRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateGlobalDatabaseNetworkRequest) SetResourceOwnerId(v int64) *CreateGlobalDatabaseNetworkRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *CreateGlobalDatabaseNetworkRequest) SetSecurityToken(v string) *CreateGlobalDatabaseNetworkRequest {
	s.SecurityToken = &v
	return s
}

type CreateGlobalDatabaseNetworkResponseBody struct {
	GDNId     *string `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateGlobalDatabaseNetworkResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateGlobalDatabaseNetworkResponseBody) GoString() string {
	return s.String()
}

func (s *CreateGlobalDatabaseNetworkResponseBody) SetGDNId(v string) *CreateGlobalDatabaseNetworkResponseBody {
	s.GDNId = &v
	return s
}

func (s *CreateGlobalDatabaseNetworkResponseBody) SetRequestId(v string) *CreateGlobalDatabaseNetworkResponseBody {
	s.RequestId = &v
	return s
}

type CreateGlobalDatabaseNetworkResponse struct {
	Headers    map[string]*string                       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateGlobalDatabaseNetworkResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateGlobalDatabaseNetworkResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateGlobalDatabaseNetworkResponse) GoString() string {
	return s.String()
}

func (s *CreateGlobalDatabaseNetworkResponse) SetHeaders(v map[string]*string) *CreateGlobalDatabaseNetworkResponse {
	s.Headers = v
	return s
}

func (s *CreateGlobalDatabaseNetworkResponse) SetStatusCode(v int32) *CreateGlobalDatabaseNetworkResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateGlobalDatabaseNetworkResponse) SetBody(v *CreateGlobalDatabaseNetworkResponseBody) *CreateGlobalDatabaseNetworkResponse {
	s.Body = v
	return s
}

type CreateParameterGroupRequest struct {
	DBType               *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion            *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ParameterGroupDesc   *string `json:"ParameterGroupDesc,omitempty" xml:"ParameterGroupDesc,omitempty"`
	ParameterGroupName   *string `json:"ParameterGroupName,omitempty" xml:"ParameterGroupName,omitempty"`
	Parameters           *string `json:"Parameters,omitempty" xml:"Parameters,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s CreateParameterGroupRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateParameterGroupRequest) GoString() string {
	return s.String()
}

func (s *CreateParameterGroupRequest) SetDBType(v string) *CreateParameterGroupRequest {
	s.DBType = &v
	return s
}

func (s *CreateParameterGroupRequest) SetDBVersion(v string) *CreateParameterGroupRequest {
	s.DBVersion = &v
	return s
}

func (s *CreateParameterGroupRequest) SetOwnerAccount(v string) *CreateParameterGroupRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateParameterGroupRequest) SetOwnerId(v int64) *CreateParameterGroupRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateParameterGroupRequest) SetParameterGroupDesc(v string) *CreateParameterGroupRequest {
	s.ParameterGroupDesc = &v
	return s
}

func (s *CreateParameterGroupRequest) SetParameterGroupName(v string) *CreateParameterGroupRequest {
	s.ParameterGroupName = &v
	return s
}

func (s *CreateParameterGroupRequest) SetParameters(v string) *CreateParameterGroupRequest {
	s.Parameters = &v
	return s
}

func (s *CreateParameterGroupRequest) SetRegionId(v string) *CreateParameterGroupRequest {
	s.RegionId = &v
	return s
}

func (s *CreateParameterGroupRequest) SetResourceOwnerAccount(v string) *CreateParameterGroupRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateParameterGroupRequest) SetResourceOwnerId(v int64) *CreateParameterGroupRequest {
	s.ResourceOwnerId = &v
	return s
}

type CreateParameterGroupResponseBody struct {
	ParameterGroupId *string `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	RequestId        *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateParameterGroupResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateParameterGroupResponseBody) GoString() string {
	return s.String()
}

func (s *CreateParameterGroupResponseBody) SetParameterGroupId(v string) *CreateParameterGroupResponseBody {
	s.ParameterGroupId = &v
	return s
}

func (s *CreateParameterGroupResponseBody) SetRequestId(v string) *CreateParameterGroupResponseBody {
	s.RequestId = &v
	return s
}

type CreateParameterGroupResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateParameterGroupResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateParameterGroupResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateParameterGroupResponse) GoString() string {
	return s.String()
}

func (s *CreateParameterGroupResponse) SetHeaders(v map[string]*string) *CreateParameterGroupResponse {
	s.Headers = v
	return s
}

func (s *CreateParameterGroupResponse) SetStatusCode(v int32) *CreateParameterGroupResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateParameterGroupResponse) SetBody(v *CreateParameterGroupResponseBody) *CreateParameterGroupResponse {
	s.Body = v
	return s
}

type CreateStoragePlanRequest struct {
	ClientToken          *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	Period               *string `json:"Period,omitempty" xml:"Period,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	StorageClass         *string `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
	StorageType          *string `json:"StorageType,omitempty" xml:"StorageType,omitempty"`
	UsedTime             *string `json:"UsedTime,omitempty" xml:"UsedTime,omitempty"`
}

func (s CreateStoragePlanRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateStoragePlanRequest) GoString() string {
	return s.String()
}

func (s *CreateStoragePlanRequest) SetClientToken(v string) *CreateStoragePlanRequest {
	s.ClientToken = &v
	return s
}

func (s *CreateStoragePlanRequest) SetOwnerAccount(v string) *CreateStoragePlanRequest {
	s.OwnerAccount = &v
	return s
}

func (s *CreateStoragePlanRequest) SetOwnerId(v int64) *CreateStoragePlanRequest {
	s.OwnerId = &v
	return s
}

func (s *CreateStoragePlanRequest) SetPeriod(v string) *CreateStoragePlanRequest {
	s.Period = &v
	return s
}

func (s *CreateStoragePlanRequest) SetResourceOwnerAccount(v string) *CreateStoragePlanRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *CreateStoragePlanRequest) SetResourceOwnerId(v int64) *CreateStoragePlanRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *CreateStoragePlanRequest) SetStorageClass(v string) *CreateStoragePlanRequest {
	s.StorageClass = &v
	return s
}

func (s *CreateStoragePlanRequest) SetStorageType(v string) *CreateStoragePlanRequest {
	s.StorageType = &v
	return s
}

func (s *CreateStoragePlanRequest) SetUsedTime(v string) *CreateStoragePlanRequest {
	s.UsedTime = &v
	return s
}

type CreateStoragePlanResponseBody struct {
	DBInstanceId *string `json:"DBInstanceId,omitempty" xml:"DBInstanceId,omitempty"`
	OrderId      *string `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	RequestId    *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s CreateStoragePlanResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateStoragePlanResponseBody) GoString() string {
	return s.String()
}

func (s *CreateStoragePlanResponseBody) SetDBInstanceId(v string) *CreateStoragePlanResponseBody {
	s.DBInstanceId = &v
	return s
}

func (s *CreateStoragePlanResponseBody) SetOrderId(v string) *CreateStoragePlanResponseBody {
	s.OrderId = &v
	return s
}

func (s *CreateStoragePlanResponseBody) SetRequestId(v string) *CreateStoragePlanResponseBody {
	s.RequestId = &v
	return s
}

type CreateStoragePlanResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateStoragePlanResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateStoragePlanResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateStoragePlanResponse) GoString() string {
	return s.String()
}

func (s *CreateStoragePlanResponse) SetHeaders(v map[string]*string) *CreateStoragePlanResponse {
	s.Headers = v
	return s
}

func (s *CreateStoragePlanResponse) SetStatusCode(v int32) *CreateStoragePlanResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateStoragePlanResponse) SetBody(v *CreateStoragePlanResponseBody) *CreateStoragePlanResponse {
	s.Body = v
	return s
}

type DeleteAccountRequest struct {
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteAccountRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccountRequest) GoString() string {
	return s.String()
}

func (s *DeleteAccountRequest) SetAccountName(v string) *DeleteAccountRequest {
	s.AccountName = &v
	return s
}

func (s *DeleteAccountRequest) SetDBClusterId(v string) *DeleteAccountRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteAccountRequest) SetOwnerAccount(v string) *DeleteAccountRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteAccountRequest) SetOwnerId(v int64) *DeleteAccountRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteAccountRequest) SetResourceOwnerAccount(v string) *DeleteAccountRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteAccountRequest) SetResourceOwnerId(v int64) *DeleteAccountRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteAccountResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteAccountResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccountResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteAccountResponseBody) SetRequestId(v string) *DeleteAccountResponseBody {
	s.RequestId = &v
	return s
}

type DeleteAccountResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteAccountResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteAccountResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccountResponse) GoString() string {
	return s.String()
}

func (s *DeleteAccountResponse) SetHeaders(v map[string]*string) *DeleteAccountResponse {
	s.Headers = v
	return s
}

func (s *DeleteAccountResponse) SetStatusCode(v int32) *DeleteAccountResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteAccountResponse) SetBody(v *DeleteAccountResponseBody) *DeleteAccountResponse {
	s.Body = v
	return s
}

type DeleteBackupRequest struct {
	BackupId             *string `json:"BackupId,omitempty" xml:"BackupId,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteBackupRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteBackupRequest) GoString() string {
	return s.String()
}

func (s *DeleteBackupRequest) SetBackupId(v string) *DeleteBackupRequest {
	s.BackupId = &v
	return s
}

func (s *DeleteBackupRequest) SetDBClusterId(v string) *DeleteBackupRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteBackupRequest) SetOwnerAccount(v string) *DeleteBackupRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteBackupRequest) SetOwnerId(v int64) *DeleteBackupRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteBackupRequest) SetResourceOwnerAccount(v string) *DeleteBackupRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteBackupRequest) SetResourceOwnerId(v int64) *DeleteBackupRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteBackupResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteBackupResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteBackupResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteBackupResponseBody) SetRequestId(v string) *DeleteBackupResponseBody {
	s.RequestId = &v
	return s
}

type DeleteBackupResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteBackupResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteBackupResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBackupResponse) GoString() string {
	return s.String()
}

func (s *DeleteBackupResponse) SetHeaders(v map[string]*string) *DeleteBackupResponse {
	s.Headers = v
	return s
}

func (s *DeleteBackupResponse) SetStatusCode(v int32) *DeleteBackupResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteBackupResponse) SetBody(v *DeleteBackupResponseBody) *DeleteBackupResponse {
	s.Body = v
	return s
}

type DeleteDBClusterRequest struct {
	BackupRetentionPolicyOnClusterDeletion *string `json:"BackupRetentionPolicyOnClusterDeletion,omitempty" xml:"BackupRetentionPolicyOnClusterDeletion,omitempty"`
	DBClusterId                            *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount                           *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId                                *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount                   *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId                        *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteDBClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBClusterRequest) GoString() string {
	return s.String()
}

func (s *DeleteDBClusterRequest) SetBackupRetentionPolicyOnClusterDeletion(v string) *DeleteDBClusterRequest {
	s.BackupRetentionPolicyOnClusterDeletion = &v
	return s
}

func (s *DeleteDBClusterRequest) SetDBClusterId(v string) *DeleteDBClusterRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteDBClusterRequest) SetOwnerAccount(v string) *DeleteDBClusterRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteDBClusterRequest) SetOwnerId(v int64) *DeleteDBClusterRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteDBClusterRequest) SetResourceOwnerAccount(v string) *DeleteDBClusterRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteDBClusterRequest) SetResourceOwnerId(v int64) *DeleteDBClusterRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteDBClusterResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteDBClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBClusterResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteDBClusterResponseBody) SetRequestId(v string) *DeleteDBClusterResponseBody {
	s.RequestId = &v
	return s
}

type DeleteDBClusterResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteDBClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteDBClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBClusterResponse) GoString() string {
	return s.String()
}

func (s *DeleteDBClusterResponse) SetHeaders(v map[string]*string) *DeleteDBClusterResponse {
	s.Headers = v
	return s
}

func (s *DeleteDBClusterResponse) SetStatusCode(v int32) *DeleteDBClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteDBClusterResponse) SetBody(v *DeleteDBClusterResponseBody) *DeleteDBClusterResponse {
	s.Body = v
	return s
}

type DeleteDBClusterEndpointRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBEndpointId         *string `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteDBClusterEndpointRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBClusterEndpointRequest) GoString() string {
	return s.String()
}

func (s *DeleteDBClusterEndpointRequest) SetDBClusterId(v string) *DeleteDBClusterEndpointRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteDBClusterEndpointRequest) SetDBEndpointId(v string) *DeleteDBClusterEndpointRequest {
	s.DBEndpointId = &v
	return s
}

func (s *DeleteDBClusterEndpointRequest) SetOwnerAccount(v string) *DeleteDBClusterEndpointRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteDBClusterEndpointRequest) SetOwnerId(v int64) *DeleteDBClusterEndpointRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteDBClusterEndpointRequest) SetResourceOwnerAccount(v string) *DeleteDBClusterEndpointRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteDBClusterEndpointRequest) SetResourceOwnerId(v int64) *DeleteDBClusterEndpointRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteDBClusterEndpointResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteDBClusterEndpointResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBClusterEndpointResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteDBClusterEndpointResponseBody) SetRequestId(v string) *DeleteDBClusterEndpointResponseBody {
	s.RequestId = &v
	return s
}

type DeleteDBClusterEndpointResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteDBClusterEndpointResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteDBClusterEndpointResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBClusterEndpointResponse) GoString() string {
	return s.String()
}

func (s *DeleteDBClusterEndpointResponse) SetHeaders(v map[string]*string) *DeleteDBClusterEndpointResponse {
	s.Headers = v
	return s
}

func (s *DeleteDBClusterEndpointResponse) SetStatusCode(v int32) *DeleteDBClusterEndpointResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteDBClusterEndpointResponse) SetBody(v *DeleteDBClusterEndpointResponseBody) *DeleteDBClusterEndpointResponse {
	s.Body = v
	return s
}

type DeleteDBEndpointAddressRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBEndpointId         *string `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	NetType              *string `json:"NetType,omitempty" xml:"NetType,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteDBEndpointAddressRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBEndpointAddressRequest) GoString() string {
	return s.String()
}

func (s *DeleteDBEndpointAddressRequest) SetDBClusterId(v string) *DeleteDBEndpointAddressRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteDBEndpointAddressRequest) SetDBEndpointId(v string) *DeleteDBEndpointAddressRequest {
	s.DBEndpointId = &v
	return s
}

func (s *DeleteDBEndpointAddressRequest) SetNetType(v string) *DeleteDBEndpointAddressRequest {
	s.NetType = &v
	return s
}

func (s *DeleteDBEndpointAddressRequest) SetOwnerAccount(v string) *DeleteDBEndpointAddressRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteDBEndpointAddressRequest) SetOwnerId(v int64) *DeleteDBEndpointAddressRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteDBEndpointAddressRequest) SetResourceOwnerAccount(v string) *DeleteDBEndpointAddressRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteDBEndpointAddressRequest) SetResourceOwnerId(v int64) *DeleteDBEndpointAddressRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteDBEndpointAddressResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteDBEndpointAddressResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBEndpointAddressResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteDBEndpointAddressResponseBody) SetRequestId(v string) *DeleteDBEndpointAddressResponseBody {
	s.RequestId = &v
	return s
}

type DeleteDBEndpointAddressResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteDBEndpointAddressResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteDBEndpointAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBEndpointAddressResponse) GoString() string {
	return s.String()
}

func (s *DeleteDBEndpointAddressResponse) SetHeaders(v map[string]*string) *DeleteDBEndpointAddressResponse {
	s.Headers = v
	return s
}

func (s *DeleteDBEndpointAddressResponse) SetStatusCode(v int32) *DeleteDBEndpointAddressResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteDBEndpointAddressResponse) SetBody(v *DeleteDBEndpointAddressResponseBody) *DeleteDBEndpointAddressResponse {
	s.Body = v
	return s
}

type DeleteDBLinkRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBLinkName           *string `json:"DBLinkName,omitempty" xml:"DBLinkName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteDBLinkRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBLinkRequest) GoString() string {
	return s.String()
}

func (s *DeleteDBLinkRequest) SetDBClusterId(v string) *DeleteDBLinkRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteDBLinkRequest) SetDBLinkName(v string) *DeleteDBLinkRequest {
	s.DBLinkName = &v
	return s
}

func (s *DeleteDBLinkRequest) SetOwnerAccount(v string) *DeleteDBLinkRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteDBLinkRequest) SetOwnerId(v int64) *DeleteDBLinkRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteDBLinkRequest) SetResourceOwnerAccount(v string) *DeleteDBLinkRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteDBLinkRequest) SetResourceOwnerId(v int64) *DeleteDBLinkRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteDBLinkResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteDBLinkResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBLinkResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteDBLinkResponseBody) SetRequestId(v string) *DeleteDBLinkResponseBody {
	s.RequestId = &v
	return s
}

type DeleteDBLinkResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteDBLinkResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteDBLinkResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBLinkResponse) GoString() string {
	return s.String()
}

func (s *DeleteDBLinkResponse) SetHeaders(v map[string]*string) *DeleteDBLinkResponse {
	s.Headers = v
	return s
}

func (s *DeleteDBLinkResponse) SetStatusCode(v int32) *DeleteDBLinkResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteDBLinkResponse) SetBody(v *DeleteDBLinkResponseBody) *DeleteDBLinkResponse {
	s.Body = v
	return s
}

type DeleteDBNodesRequest struct {
	ClientToken          *string   `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string   `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeId             []*string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty" type:"Repeated"`
	OwnerAccount         *string   `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64    `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string   `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64    `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteDBNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBNodesRequest) GoString() string {
	return s.String()
}

func (s *DeleteDBNodesRequest) SetClientToken(v string) *DeleteDBNodesRequest {
	s.ClientToken = &v
	return s
}

func (s *DeleteDBNodesRequest) SetDBClusterId(v string) *DeleteDBNodesRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteDBNodesRequest) SetDBNodeId(v []*string) *DeleteDBNodesRequest {
	s.DBNodeId = v
	return s
}

func (s *DeleteDBNodesRequest) SetOwnerAccount(v string) *DeleteDBNodesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteDBNodesRequest) SetOwnerId(v int64) *DeleteDBNodesRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteDBNodesRequest) SetResourceOwnerAccount(v string) *DeleteDBNodesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteDBNodesRequest) SetResourceOwnerId(v int64) *DeleteDBNodesRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteDBNodesResponseBody struct {
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OrderId     *string `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	RequestId   *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteDBNodesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBNodesResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteDBNodesResponseBody) SetDBClusterId(v string) *DeleteDBNodesResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DeleteDBNodesResponseBody) SetOrderId(v string) *DeleteDBNodesResponseBody {
	s.OrderId = &v
	return s
}

func (s *DeleteDBNodesResponseBody) SetRequestId(v string) *DeleteDBNodesResponseBody {
	s.RequestId = &v
	return s
}

type DeleteDBNodesResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteDBNodesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteDBNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteDBNodesResponse) GoString() string {
	return s.String()
}

func (s *DeleteDBNodesResponse) SetHeaders(v map[string]*string) *DeleteDBNodesResponse {
	s.Headers = v
	return s
}

func (s *DeleteDBNodesResponse) SetStatusCode(v int32) *DeleteDBNodesResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteDBNodesResponse) SetBody(v *DeleteDBNodesResponseBody) *DeleteDBNodesResponse {
	s.Body = v
	return s
}

type DeleteDatabaseRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteDatabaseRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteDatabaseRequest) GoString() string {
	return s.String()
}

func (s *DeleteDatabaseRequest) SetDBClusterId(v string) *DeleteDatabaseRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteDatabaseRequest) SetDBName(v string) *DeleteDatabaseRequest {
	s.DBName = &v
	return s
}

func (s *DeleteDatabaseRequest) SetOwnerAccount(v string) *DeleteDatabaseRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteDatabaseRequest) SetOwnerId(v int64) *DeleteDatabaseRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteDatabaseRequest) SetResourceOwnerAccount(v string) *DeleteDatabaseRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteDatabaseRequest) SetResourceOwnerId(v int64) *DeleteDatabaseRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteDatabaseResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteDatabaseResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteDatabaseResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteDatabaseResponseBody) SetRequestId(v string) *DeleteDatabaseResponseBody {
	s.RequestId = &v
	return s
}

type DeleteDatabaseResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteDatabaseResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteDatabaseResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteDatabaseResponse) GoString() string {
	return s.String()
}

func (s *DeleteDatabaseResponse) SetHeaders(v map[string]*string) *DeleteDatabaseResponse {
	s.Headers = v
	return s
}

func (s *DeleteDatabaseResponse) SetStatusCode(v int32) *DeleteDatabaseResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteDatabaseResponse) SetBody(v *DeleteDatabaseResponseBody) *DeleteDatabaseResponse {
	s.Body = v
	return s
}

type DeleteGlobalDatabaseNetworkRequest struct {
	GDNId                *string `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s DeleteGlobalDatabaseNetworkRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteGlobalDatabaseNetworkRequest) GoString() string {
	return s.String()
}

func (s *DeleteGlobalDatabaseNetworkRequest) SetGDNId(v string) *DeleteGlobalDatabaseNetworkRequest {
	s.GDNId = &v
	return s
}

func (s *DeleteGlobalDatabaseNetworkRequest) SetOwnerAccount(v string) *DeleteGlobalDatabaseNetworkRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteGlobalDatabaseNetworkRequest) SetOwnerId(v int64) *DeleteGlobalDatabaseNetworkRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteGlobalDatabaseNetworkRequest) SetResourceOwnerAccount(v string) *DeleteGlobalDatabaseNetworkRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteGlobalDatabaseNetworkRequest) SetResourceOwnerId(v int64) *DeleteGlobalDatabaseNetworkRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DeleteGlobalDatabaseNetworkRequest) SetSecurityToken(v string) *DeleteGlobalDatabaseNetworkRequest {
	s.SecurityToken = &v
	return s
}

type DeleteGlobalDatabaseNetworkResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteGlobalDatabaseNetworkResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteGlobalDatabaseNetworkResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteGlobalDatabaseNetworkResponseBody) SetRequestId(v string) *DeleteGlobalDatabaseNetworkResponseBody {
	s.RequestId = &v
	return s
}

type DeleteGlobalDatabaseNetworkResponse struct {
	Headers    map[string]*string                       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteGlobalDatabaseNetworkResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteGlobalDatabaseNetworkResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteGlobalDatabaseNetworkResponse) GoString() string {
	return s.String()
}

func (s *DeleteGlobalDatabaseNetworkResponse) SetHeaders(v map[string]*string) *DeleteGlobalDatabaseNetworkResponse {
	s.Headers = v
	return s
}

func (s *DeleteGlobalDatabaseNetworkResponse) SetStatusCode(v int32) *DeleteGlobalDatabaseNetworkResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteGlobalDatabaseNetworkResponse) SetBody(v *DeleteGlobalDatabaseNetworkResponseBody) *DeleteGlobalDatabaseNetworkResponse {
	s.Body = v
	return s
}

type DeleteMaskingRulesRequest struct {
	DBClusterId  *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	RuleNameList *string `json:"RuleNameList,omitempty" xml:"RuleNameList,omitempty"`
}

func (s DeleteMaskingRulesRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteMaskingRulesRequest) GoString() string {
	return s.String()
}

func (s *DeleteMaskingRulesRequest) SetDBClusterId(v string) *DeleteMaskingRulesRequest {
	s.DBClusterId = &v
	return s
}

func (s *DeleteMaskingRulesRequest) SetRuleNameList(v string) *DeleteMaskingRulesRequest {
	s.RuleNameList = &v
	return s
}

type DeleteMaskingRulesResponseBody struct {
	Message   *string `json:"Message,omitempty" xml:"Message,omitempty"`
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	Success   *bool   `json:"Success,omitempty" xml:"Success,omitempty"`
}

func (s DeleteMaskingRulesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteMaskingRulesResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteMaskingRulesResponseBody) SetMessage(v string) *DeleteMaskingRulesResponseBody {
	s.Message = &v
	return s
}

func (s *DeleteMaskingRulesResponseBody) SetRequestId(v string) *DeleteMaskingRulesResponseBody {
	s.RequestId = &v
	return s
}

func (s *DeleteMaskingRulesResponseBody) SetSuccess(v bool) *DeleteMaskingRulesResponseBody {
	s.Success = &v
	return s
}

type DeleteMaskingRulesResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteMaskingRulesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteMaskingRulesResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteMaskingRulesResponse) GoString() string {
	return s.String()
}

func (s *DeleteMaskingRulesResponse) SetHeaders(v map[string]*string) *DeleteMaskingRulesResponse {
	s.Headers = v
	return s
}

func (s *DeleteMaskingRulesResponse) SetStatusCode(v int32) *DeleteMaskingRulesResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteMaskingRulesResponse) SetBody(v *DeleteMaskingRulesResponseBody) *DeleteMaskingRulesResponse {
	s.Body = v
	return s
}

type DeleteParameterGroupRequest struct {
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ParameterGroupId     *string `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DeleteParameterGroupRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteParameterGroupRequest) GoString() string {
	return s.String()
}

func (s *DeleteParameterGroupRequest) SetOwnerAccount(v string) *DeleteParameterGroupRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DeleteParameterGroupRequest) SetOwnerId(v int64) *DeleteParameterGroupRequest {
	s.OwnerId = &v
	return s
}

func (s *DeleteParameterGroupRequest) SetParameterGroupId(v string) *DeleteParameterGroupRequest {
	s.ParameterGroupId = &v
	return s
}

func (s *DeleteParameterGroupRequest) SetRegionId(v string) *DeleteParameterGroupRequest {
	s.RegionId = &v
	return s
}

func (s *DeleteParameterGroupRequest) SetResourceOwnerAccount(v string) *DeleteParameterGroupRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DeleteParameterGroupRequest) SetResourceOwnerId(v int64) *DeleteParameterGroupRequest {
	s.ResourceOwnerId = &v
	return s
}

type DeleteParameterGroupResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DeleteParameterGroupResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteParameterGroupResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteParameterGroupResponseBody) SetRequestId(v string) *DeleteParameterGroupResponseBody {
	s.RequestId = &v
	return s
}

type DeleteParameterGroupResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteParameterGroupResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteParameterGroupResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteParameterGroupResponse) GoString() string {
	return s.String()
}

func (s *DeleteParameterGroupResponse) SetHeaders(v map[string]*string) *DeleteParameterGroupResponse {
	s.Headers = v
	return s
}

func (s *DeleteParameterGroupResponse) SetStatusCode(v int32) *DeleteParameterGroupResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteParameterGroupResponse) SetBody(v *DeleteParameterGroupResponseBody) *DeleteParameterGroupResponse {
	s.Body = v
	return s
}

type DescribeAITaskStatusRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeAITaskStatusRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeAITaskStatusRequest) GoString() string {
	return s.String()
}

func (s *DescribeAITaskStatusRequest) SetDBClusterId(v string) *DescribeAITaskStatusRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeAITaskStatusRequest) SetOwnerAccount(v string) *DescribeAITaskStatusRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeAITaskStatusRequest) SetOwnerId(v int64) *DescribeAITaskStatusRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeAITaskStatusRequest) SetRegionId(v string) *DescribeAITaskStatusRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeAITaskStatusRequest) SetResourceOwnerAccount(v string) *DescribeAITaskStatusRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeAITaskStatusRequest) SetResourceOwnerId(v int64) *DescribeAITaskStatusRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeAITaskStatusResponseBody struct {
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	RequestId   *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	Status      *string `json:"Status,omitempty" xml:"Status,omitempty"`
	StatusName  *string `json:"StatusName,omitempty" xml:"StatusName,omitempty"`
}

func (s DescribeAITaskStatusResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeAITaskStatusResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeAITaskStatusResponseBody) SetDBClusterId(v string) *DescribeAITaskStatusResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeAITaskStatusResponseBody) SetRequestId(v string) *DescribeAITaskStatusResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeAITaskStatusResponseBody) SetStatus(v string) *DescribeAITaskStatusResponseBody {
	s.Status = &v
	return s
}

func (s *DescribeAITaskStatusResponseBody) SetStatusName(v string) *DescribeAITaskStatusResponseBody {
	s.StatusName = &v
	return s
}

type DescribeAITaskStatusResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeAITaskStatusResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeAITaskStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeAITaskStatusResponse) GoString() string {
	return s.String()
}

func (s *DescribeAITaskStatusResponse) SetHeaders(v map[string]*string) *DescribeAITaskStatusResponse {
	s.Headers = v
	return s
}

func (s *DescribeAITaskStatusResponse) SetStatusCode(v int32) *DescribeAITaskStatusResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeAITaskStatusResponse) SetBody(v *DescribeAITaskStatusResponseBody) *DescribeAITaskStatusResponse {
	s.Body = v
	return s
}

type DescribeAccountsRequest struct {
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeAccountsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeAccountsRequest) GoString() string {
	return s.String()
}

func (s *DescribeAccountsRequest) SetAccountName(v string) *DescribeAccountsRequest {
	s.AccountName = &v
	return s
}

func (s *DescribeAccountsRequest) SetDBClusterId(v string) *DescribeAccountsRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeAccountsRequest) SetOwnerAccount(v string) *DescribeAccountsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeAccountsRequest) SetOwnerId(v int64) *DescribeAccountsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeAccountsRequest) SetPageNumber(v int32) *DescribeAccountsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeAccountsRequest) SetPageSize(v int32) *DescribeAccountsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeAccountsRequest) SetResourceOwnerAccount(v string) *DescribeAccountsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeAccountsRequest) SetResourceOwnerId(v int64) *DescribeAccountsRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeAccountsResponseBody struct {
	Accounts        []*DescribeAccountsResponseBodyAccounts `json:"Accounts,omitempty" xml:"Accounts,omitempty" type:"Repeated"`
	PageNumber      *int32                                  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount *int32                                  `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId       *string                                 `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeAccountsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeAccountsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeAccountsResponseBody) SetAccounts(v []*DescribeAccountsResponseBodyAccounts) *DescribeAccountsResponseBody {
	s.Accounts = v
	return s
}

func (s *DescribeAccountsResponseBody) SetPageNumber(v int32) *DescribeAccountsResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeAccountsResponseBody) SetPageRecordCount(v int32) *DescribeAccountsResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeAccountsResponseBody) SetRequestId(v string) *DescribeAccountsResponseBody {
	s.RequestId = &v
	return s
}

type DescribeAccountsResponseBodyAccounts struct {
	AccountDescription       *string                                                   `json:"AccountDescription,omitempty" xml:"AccountDescription,omitempty"`
	AccountLockState         *string                                                   `json:"AccountLockState,omitempty" xml:"AccountLockState,omitempty"`
	AccountName              *string                                                   `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	AccountPasswordValidTime *string                                                   `json:"AccountPasswordValidTime,omitempty" xml:"AccountPasswordValidTime,omitempty"`
	AccountStatus            *string                                                   `json:"AccountStatus,omitempty" xml:"AccountStatus,omitempty"`
	AccountType              *string                                                   `json:"AccountType,omitempty" xml:"AccountType,omitempty"`
	DatabasePrivileges       []*DescribeAccountsResponseBodyAccountsDatabasePrivileges `json:"DatabasePrivileges,omitempty" xml:"DatabasePrivileges,omitempty" type:"Repeated"`
}

func (s DescribeAccountsResponseBodyAccounts) String() string {
	return tea.Prettify(s)
}

func (s DescribeAccountsResponseBodyAccounts) GoString() string {
	return s.String()
}

func (s *DescribeAccountsResponseBodyAccounts) SetAccountDescription(v string) *DescribeAccountsResponseBodyAccounts {
	s.AccountDescription = &v
	return s
}

func (s *DescribeAccountsResponseBodyAccounts) SetAccountLockState(v string) *DescribeAccountsResponseBodyAccounts {
	s.AccountLockState = &v
	return s
}

func (s *DescribeAccountsResponseBodyAccounts) SetAccountName(v string) *DescribeAccountsResponseBodyAccounts {
	s.AccountName = &v
	return s
}

func (s *DescribeAccountsResponseBodyAccounts) SetAccountPasswordValidTime(v string) *DescribeAccountsResponseBodyAccounts {
	s.AccountPasswordValidTime = &v
	return s
}

func (s *DescribeAccountsResponseBodyAccounts) SetAccountStatus(v string) *DescribeAccountsResponseBodyAccounts {
	s.AccountStatus = &v
	return s
}

func (s *DescribeAccountsResponseBodyAccounts) SetAccountType(v string) *DescribeAccountsResponseBodyAccounts {
	s.AccountType = &v
	return s
}

func (s *DescribeAccountsResponseBodyAccounts) SetDatabasePrivileges(v []*DescribeAccountsResponseBodyAccountsDatabasePrivileges) *DescribeAccountsResponseBodyAccounts {
	s.DatabasePrivileges = v
	return s
}

type DescribeAccountsResponseBodyAccountsDatabasePrivileges struct {
	AccountPrivilege *string `json:"AccountPrivilege,omitempty" xml:"AccountPrivilege,omitempty"`
	DBName           *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
}

func (s DescribeAccountsResponseBodyAccountsDatabasePrivileges) String() string {
	return tea.Prettify(s)
}

func (s DescribeAccountsResponseBodyAccountsDatabasePrivileges) GoString() string {
	return s.String()
}

func (s *DescribeAccountsResponseBodyAccountsDatabasePrivileges) SetAccountPrivilege(v string) *DescribeAccountsResponseBodyAccountsDatabasePrivileges {
	s.AccountPrivilege = &v
	return s
}

func (s *DescribeAccountsResponseBodyAccountsDatabasePrivileges) SetDBName(v string) *DescribeAccountsResponseBodyAccountsDatabasePrivileges {
	s.DBName = &v
	return s
}

type DescribeAccountsResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeAccountsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeAccountsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeAccountsResponse) GoString() string {
	return s.String()
}

func (s *DescribeAccountsResponse) SetHeaders(v map[string]*string) *DescribeAccountsResponse {
	s.Headers = v
	return s
}

func (s *DescribeAccountsResponse) SetStatusCode(v int32) *DescribeAccountsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeAccountsResponse) SetBody(v *DescribeAccountsResponseBody) *DescribeAccountsResponse {
	s.Body = v
	return s
}

type DescribeAutoRenewAttributeRequest struct {
	DBClusterIds         *string `json:"DBClusterIds,omitempty" xml:"DBClusterIds,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceGroupId      *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeAutoRenewAttributeRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeAutoRenewAttributeRequest) GoString() string {
	return s.String()
}

func (s *DescribeAutoRenewAttributeRequest) SetDBClusterIds(v string) *DescribeAutoRenewAttributeRequest {
	s.DBClusterIds = &v
	return s
}

func (s *DescribeAutoRenewAttributeRequest) SetOwnerAccount(v string) *DescribeAutoRenewAttributeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeAutoRenewAttributeRequest) SetOwnerId(v int64) *DescribeAutoRenewAttributeRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeAutoRenewAttributeRequest) SetPageNumber(v int32) *DescribeAutoRenewAttributeRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeAutoRenewAttributeRequest) SetPageSize(v int32) *DescribeAutoRenewAttributeRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeAutoRenewAttributeRequest) SetRegionId(v string) *DescribeAutoRenewAttributeRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeAutoRenewAttributeRequest) SetResourceGroupId(v string) *DescribeAutoRenewAttributeRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeAutoRenewAttributeRequest) SetResourceOwnerAccount(v string) *DescribeAutoRenewAttributeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeAutoRenewAttributeRequest) SetResourceOwnerId(v int64) *DescribeAutoRenewAttributeRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeAutoRenewAttributeResponseBody struct {
	Items            *DescribeAutoRenewAttributeResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	PageNumber       *int32                                       `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *int32                                       `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                                      `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *int32                                       `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeAutoRenewAttributeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeAutoRenewAttributeResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeAutoRenewAttributeResponseBody) SetItems(v *DescribeAutoRenewAttributeResponseBodyItems) *DescribeAutoRenewAttributeResponseBody {
	s.Items = v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBody) SetPageNumber(v int32) *DescribeAutoRenewAttributeResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBody) SetPageRecordCount(v int32) *DescribeAutoRenewAttributeResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBody) SetRequestId(v string) *DescribeAutoRenewAttributeResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBody) SetTotalRecordCount(v int32) *DescribeAutoRenewAttributeResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeAutoRenewAttributeResponseBodyItems struct {
	AutoRenewAttribute []*DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute `json:"AutoRenewAttribute,omitempty" xml:"AutoRenewAttribute,omitempty" type:"Repeated"`
}

func (s DescribeAutoRenewAttributeResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeAutoRenewAttributeResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeAutoRenewAttributeResponseBodyItems) SetAutoRenewAttribute(v []*DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) *DescribeAutoRenewAttributeResponseBodyItems {
	s.AutoRenewAttribute = v
	return s
}

type DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute struct {
	AutoRenewEnabled *bool   `json:"AutoRenewEnabled,omitempty" xml:"AutoRenewEnabled,omitempty"`
	DBClusterId      *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	Duration         *int32  `json:"Duration,omitempty" xml:"Duration,omitempty"`
	PeriodUnit       *string `json:"PeriodUnit,omitempty" xml:"PeriodUnit,omitempty"`
	RegionId         *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	RenewalStatus    *string `json:"RenewalStatus,omitempty" xml:"RenewalStatus,omitempty"`
}

func (s DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) String() string {
	return tea.Prettify(s)
}

func (s DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) GoString() string {
	return s.String()
}

func (s *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) SetAutoRenewEnabled(v bool) *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute {
	s.AutoRenewEnabled = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) SetDBClusterId(v string) *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute {
	s.DBClusterId = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) SetDuration(v int32) *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute {
	s.Duration = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) SetPeriodUnit(v string) *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute {
	s.PeriodUnit = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) SetRegionId(v string) *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute {
	s.RegionId = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute) SetRenewalStatus(v string) *DescribeAutoRenewAttributeResponseBodyItemsAutoRenewAttribute {
	s.RenewalStatus = &v
	return s
}

type DescribeAutoRenewAttributeResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeAutoRenewAttributeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeAutoRenewAttributeResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeAutoRenewAttributeResponse) GoString() string {
	return s.String()
}

func (s *DescribeAutoRenewAttributeResponse) SetHeaders(v map[string]*string) *DescribeAutoRenewAttributeResponse {
	s.Headers = v
	return s
}

func (s *DescribeAutoRenewAttributeResponse) SetStatusCode(v int32) *DescribeAutoRenewAttributeResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeAutoRenewAttributeResponse) SetBody(v *DescribeAutoRenewAttributeResponseBody) *DescribeAutoRenewAttributeResponse {
	s.Body = v
	return s
}

type DescribeBackupLogsRequest struct {
	BackupRegion         *string `json:"BackupRegion,omitempty" xml:"BackupRegion,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EndTime              *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	StartTime            *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeBackupLogsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupLogsRequest) GoString() string {
	return s.String()
}

func (s *DescribeBackupLogsRequest) SetBackupRegion(v string) *DescribeBackupLogsRequest {
	s.BackupRegion = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetDBClusterId(v string) *DescribeBackupLogsRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetEndTime(v string) *DescribeBackupLogsRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetOwnerAccount(v string) *DescribeBackupLogsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetOwnerId(v int64) *DescribeBackupLogsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetPageNumber(v int32) *DescribeBackupLogsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetPageSize(v int32) *DescribeBackupLogsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetResourceOwnerAccount(v string) *DescribeBackupLogsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetResourceOwnerId(v int64) *DescribeBackupLogsRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeBackupLogsRequest) SetStartTime(v string) *DescribeBackupLogsRequest {
	s.StartTime = &v
	return s
}

type DescribeBackupLogsResponseBody struct {
	Items            *DescribeBackupLogsResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	PageNumber       *string                              `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *string                              `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                              `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *string                              `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeBackupLogsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupLogsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeBackupLogsResponseBody) SetItems(v *DescribeBackupLogsResponseBodyItems) *DescribeBackupLogsResponseBody {
	s.Items = v
	return s
}

func (s *DescribeBackupLogsResponseBody) SetPageNumber(v string) *DescribeBackupLogsResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeBackupLogsResponseBody) SetPageRecordCount(v string) *DescribeBackupLogsResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeBackupLogsResponseBody) SetRequestId(v string) *DescribeBackupLogsResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeBackupLogsResponseBody) SetTotalRecordCount(v string) *DescribeBackupLogsResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeBackupLogsResponseBodyItems struct {
	BackupLog []*DescribeBackupLogsResponseBodyItemsBackupLog `json:"BackupLog,omitempty" xml:"BackupLog,omitempty" type:"Repeated"`
}

func (s DescribeBackupLogsResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupLogsResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeBackupLogsResponseBodyItems) SetBackupLog(v []*DescribeBackupLogsResponseBodyItemsBackupLog) *DescribeBackupLogsResponseBodyItems {
	s.BackupLog = v
	return s
}

type DescribeBackupLogsResponseBodyItemsBackupLog struct {
	BackupLogEndTime     *string `json:"BackupLogEndTime,omitempty" xml:"BackupLogEndTime,omitempty"`
	BackupLogId          *string `json:"BackupLogId,omitempty" xml:"BackupLogId,omitempty"`
	BackupLogName        *string `json:"BackupLogName,omitempty" xml:"BackupLogName,omitempty"`
	BackupLogSize        *string `json:"BackupLogSize,omitempty" xml:"BackupLogSize,omitempty"`
	BackupLogStartTime   *string `json:"BackupLogStartTime,omitempty" xml:"BackupLogStartTime,omitempty"`
	DownloadLink         *string `json:"DownloadLink,omitempty" xml:"DownloadLink,omitempty"`
	IntranetDownloadLink *string `json:"IntranetDownloadLink,omitempty" xml:"IntranetDownloadLink,omitempty"`
	LinkExpiredTime      *string `json:"LinkExpiredTime,omitempty" xml:"LinkExpiredTime,omitempty"`
}

func (s DescribeBackupLogsResponseBodyItemsBackupLog) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupLogsResponseBodyItemsBackupLog) GoString() string {
	return s.String()
}

func (s *DescribeBackupLogsResponseBodyItemsBackupLog) SetBackupLogEndTime(v string) *DescribeBackupLogsResponseBodyItemsBackupLog {
	s.BackupLogEndTime = &v
	return s
}

func (s *DescribeBackupLogsResponseBodyItemsBackupLog) SetBackupLogId(v string) *DescribeBackupLogsResponseBodyItemsBackupLog {
	s.BackupLogId = &v
	return s
}

func (s *DescribeBackupLogsResponseBodyItemsBackupLog) SetBackupLogName(v string) *DescribeBackupLogsResponseBodyItemsBackupLog {
	s.BackupLogName = &v
	return s
}

func (s *DescribeBackupLogsResponseBodyItemsBackupLog) SetBackupLogSize(v string) *DescribeBackupLogsResponseBodyItemsBackupLog {
	s.BackupLogSize = &v
	return s
}

func (s *DescribeBackupLogsResponseBodyItemsBackupLog) SetBackupLogStartTime(v string) *DescribeBackupLogsResponseBodyItemsBackupLog {
	s.BackupLogStartTime = &v
	return s
}

func (s *DescribeBackupLogsResponseBodyItemsBackupLog) SetDownloadLink(v string) *DescribeBackupLogsResponseBodyItemsBackupLog {
	s.DownloadLink = &v
	return s
}

func (s *DescribeBackupLogsResponseBodyItemsBackupLog) SetIntranetDownloadLink(v string) *DescribeBackupLogsResponseBodyItemsBackupLog {
	s.IntranetDownloadLink = &v
	return s
}

func (s *DescribeBackupLogsResponseBodyItemsBackupLog) SetLinkExpiredTime(v string) *DescribeBackupLogsResponseBodyItemsBackupLog {
	s.LinkExpiredTime = &v
	return s
}

type DescribeBackupLogsResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeBackupLogsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeBackupLogsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupLogsResponse) GoString() string {
	return s.String()
}

func (s *DescribeBackupLogsResponse) SetHeaders(v map[string]*string) *DescribeBackupLogsResponse {
	s.Headers = v
	return s
}

func (s *DescribeBackupLogsResponse) SetStatusCode(v int32) *DescribeBackupLogsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeBackupLogsResponse) SetBody(v *DescribeBackupLogsResponseBody) *DescribeBackupLogsResponse {
	s.Body = v
	return s
}

type DescribeBackupPolicyRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeBackupPolicyRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupPolicyRequest) GoString() string {
	return s.String()
}

func (s *DescribeBackupPolicyRequest) SetDBClusterId(v string) *DescribeBackupPolicyRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeBackupPolicyRequest) SetOwnerAccount(v string) *DescribeBackupPolicyRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeBackupPolicyRequest) SetOwnerId(v int64) *DescribeBackupPolicyRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeBackupPolicyRequest) SetResourceOwnerAccount(v string) *DescribeBackupPolicyRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeBackupPolicyRequest) SetResourceOwnerId(v int64) *DescribeBackupPolicyRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeBackupPolicyResponseBody struct {
	BackupFrequency                              *string `json:"BackupFrequency,omitempty" xml:"BackupFrequency,omitempty"`
	BackupRetentionPolicyOnClusterDeletion       *string `json:"BackupRetentionPolicyOnClusterDeletion,omitempty" xml:"BackupRetentionPolicyOnClusterDeletion,omitempty"`
	DataLevel1BackupFrequency                    *string `json:"DataLevel1BackupFrequency,omitempty" xml:"DataLevel1BackupFrequency,omitempty"`
	DataLevel1BackupPeriod                       *string `json:"DataLevel1BackupPeriod,omitempty" xml:"DataLevel1BackupPeriod,omitempty"`
	DataLevel1BackupRetentionPeriod              *string `json:"DataLevel1BackupRetentionPeriod,omitempty" xml:"DataLevel1BackupRetentionPeriod,omitempty"`
	DataLevel1BackupTime                         *string `json:"DataLevel1BackupTime,omitempty" xml:"DataLevel1BackupTime,omitempty"`
	DataLevel2BackupAnotherRegionRegion          *string `json:"DataLevel2BackupAnotherRegionRegion,omitempty" xml:"DataLevel2BackupAnotherRegionRegion,omitempty"`
	DataLevel2BackupAnotherRegionRetentionPeriod *string `json:"DataLevel2BackupAnotherRegionRetentionPeriod,omitempty" xml:"DataLevel2BackupAnotherRegionRetentionPeriod,omitempty"`
	DataLevel2BackupPeriod                       *string `json:"DataLevel2BackupPeriod,omitempty" xml:"DataLevel2BackupPeriod,omitempty"`
	DataLevel2BackupRetentionPeriod              *string `json:"DataLevel2BackupRetentionPeriod,omitempty" xml:"DataLevel2BackupRetentionPeriod,omitempty"`
	PreferredBackupPeriod                        *string `json:"PreferredBackupPeriod,omitempty" xml:"PreferredBackupPeriod,omitempty"`
	PreferredBackupTime                          *string `json:"PreferredBackupTime,omitempty" xml:"PreferredBackupTime,omitempty"`
	PreferredNextBackupTime                      *string `json:"PreferredNextBackupTime,omitempty" xml:"PreferredNextBackupTime,omitempty"`
	RequestId                                    *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeBackupPolicyResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupPolicyResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeBackupPolicyResponseBody) SetBackupFrequency(v string) *DescribeBackupPolicyResponseBody {
	s.BackupFrequency = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetBackupRetentionPolicyOnClusterDeletion(v string) *DescribeBackupPolicyResponseBody {
	s.BackupRetentionPolicyOnClusterDeletion = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetDataLevel1BackupFrequency(v string) *DescribeBackupPolicyResponseBody {
	s.DataLevel1BackupFrequency = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetDataLevel1BackupPeriod(v string) *DescribeBackupPolicyResponseBody {
	s.DataLevel1BackupPeriod = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetDataLevel1BackupRetentionPeriod(v string) *DescribeBackupPolicyResponseBody {
	s.DataLevel1BackupRetentionPeriod = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetDataLevel1BackupTime(v string) *DescribeBackupPolicyResponseBody {
	s.DataLevel1BackupTime = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetDataLevel2BackupAnotherRegionRegion(v string) *DescribeBackupPolicyResponseBody {
	s.DataLevel2BackupAnotherRegionRegion = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetDataLevel2BackupAnotherRegionRetentionPeriod(v string) *DescribeBackupPolicyResponseBody {
	s.DataLevel2BackupAnotherRegionRetentionPeriod = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetDataLevel2BackupPeriod(v string) *DescribeBackupPolicyResponseBody {
	s.DataLevel2BackupPeriod = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetDataLevel2BackupRetentionPeriod(v string) *DescribeBackupPolicyResponseBody {
	s.DataLevel2BackupRetentionPeriod = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetPreferredBackupPeriod(v string) *DescribeBackupPolicyResponseBody {
	s.PreferredBackupPeriod = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetPreferredBackupTime(v string) *DescribeBackupPolicyResponseBody {
	s.PreferredBackupTime = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetPreferredNextBackupTime(v string) *DescribeBackupPolicyResponseBody {
	s.PreferredNextBackupTime = &v
	return s
}

func (s *DescribeBackupPolicyResponseBody) SetRequestId(v string) *DescribeBackupPolicyResponseBody {
	s.RequestId = &v
	return s
}

type DescribeBackupPolicyResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeBackupPolicyResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeBackupPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupPolicyResponse) GoString() string {
	return s.String()
}

func (s *DescribeBackupPolicyResponse) SetHeaders(v map[string]*string) *DescribeBackupPolicyResponse {
	s.Headers = v
	return s
}

func (s *DescribeBackupPolicyResponse) SetStatusCode(v int32) *DescribeBackupPolicyResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeBackupPolicyResponse) SetBody(v *DescribeBackupPolicyResponseBody) *DescribeBackupPolicyResponse {
	s.Body = v
	return s
}

type DescribeBackupTasksRequest struct {
	BackupJobId          *string `json:"BackupJobId,omitempty" xml:"BackupJobId,omitempty"`
	BackupMode           *string `json:"BackupMode,omitempty" xml:"BackupMode,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeBackupTasksRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupTasksRequest) GoString() string {
	return s.String()
}

func (s *DescribeBackupTasksRequest) SetBackupJobId(v string) *DescribeBackupTasksRequest {
	s.BackupJobId = &v
	return s
}

func (s *DescribeBackupTasksRequest) SetBackupMode(v string) *DescribeBackupTasksRequest {
	s.BackupMode = &v
	return s
}

func (s *DescribeBackupTasksRequest) SetDBClusterId(v string) *DescribeBackupTasksRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeBackupTasksRequest) SetOwnerAccount(v string) *DescribeBackupTasksRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeBackupTasksRequest) SetOwnerId(v int64) *DescribeBackupTasksRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeBackupTasksRequest) SetResourceOwnerAccount(v string) *DescribeBackupTasksRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeBackupTasksRequest) SetResourceOwnerId(v int64) *DescribeBackupTasksRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeBackupTasksResponseBody struct {
	Items     *DescribeBackupTasksResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	RequestId *string                               `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeBackupTasksResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupTasksResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeBackupTasksResponseBody) SetItems(v *DescribeBackupTasksResponseBodyItems) *DescribeBackupTasksResponseBody {
	s.Items = v
	return s
}

func (s *DescribeBackupTasksResponseBody) SetRequestId(v string) *DescribeBackupTasksResponseBody {
	s.RequestId = &v
	return s
}

type DescribeBackupTasksResponseBodyItems struct {
	BackupJob []*DescribeBackupTasksResponseBodyItemsBackupJob `json:"BackupJob,omitempty" xml:"BackupJob,omitempty" type:"Repeated"`
}

func (s DescribeBackupTasksResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupTasksResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeBackupTasksResponseBodyItems) SetBackupJob(v []*DescribeBackupTasksResponseBodyItemsBackupJob) *DescribeBackupTasksResponseBodyItems {
	s.BackupJob = v
	return s
}

type DescribeBackupTasksResponseBodyItemsBackupJob struct {
	BackupJobId          *string `json:"BackupJobId,omitempty" xml:"BackupJobId,omitempty"`
	BackupProgressStatus *string `json:"BackupProgressStatus,omitempty" xml:"BackupProgressStatus,omitempty"`
	JobMode              *string `json:"JobMode,omitempty" xml:"JobMode,omitempty"`
	Process              *string `json:"Process,omitempty" xml:"Process,omitempty"`
	StartTime            *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
	TaskAction           *string `json:"TaskAction,omitempty" xml:"TaskAction,omitempty"`
}

func (s DescribeBackupTasksResponseBodyItemsBackupJob) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupTasksResponseBodyItemsBackupJob) GoString() string {
	return s.String()
}

func (s *DescribeBackupTasksResponseBodyItemsBackupJob) SetBackupJobId(v string) *DescribeBackupTasksResponseBodyItemsBackupJob {
	s.BackupJobId = &v
	return s
}

func (s *DescribeBackupTasksResponseBodyItemsBackupJob) SetBackupProgressStatus(v string) *DescribeBackupTasksResponseBodyItemsBackupJob {
	s.BackupProgressStatus = &v
	return s
}

func (s *DescribeBackupTasksResponseBodyItemsBackupJob) SetJobMode(v string) *DescribeBackupTasksResponseBodyItemsBackupJob {
	s.JobMode = &v
	return s
}

func (s *DescribeBackupTasksResponseBodyItemsBackupJob) SetProcess(v string) *DescribeBackupTasksResponseBodyItemsBackupJob {
	s.Process = &v
	return s
}

func (s *DescribeBackupTasksResponseBodyItemsBackupJob) SetStartTime(v string) *DescribeBackupTasksResponseBodyItemsBackupJob {
	s.StartTime = &v
	return s
}

func (s *DescribeBackupTasksResponseBodyItemsBackupJob) SetTaskAction(v string) *DescribeBackupTasksResponseBodyItemsBackupJob {
	s.TaskAction = &v
	return s
}

type DescribeBackupTasksResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeBackupTasksResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeBackupTasksResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupTasksResponse) GoString() string {
	return s.String()
}

func (s *DescribeBackupTasksResponse) SetHeaders(v map[string]*string) *DescribeBackupTasksResponse {
	s.Headers = v
	return s
}

func (s *DescribeBackupTasksResponse) SetStatusCode(v int32) *DescribeBackupTasksResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeBackupTasksResponse) SetBody(v *DescribeBackupTasksResponseBody) *DescribeBackupTasksResponse {
	s.Body = v
	return s
}

type DescribeBackupsRequest struct {
	BackupId             *string `json:"BackupId,omitempty" xml:"BackupId,omitempty"`
	BackupMode           *string `json:"BackupMode,omitempty" xml:"BackupMode,omitempty"`
	BackupRegion         *string `json:"BackupRegion,omitempty" xml:"BackupRegion,omitempty"`
	BackupStatus         *string `json:"BackupStatus,omitempty" xml:"BackupStatus,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EndTime              *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	StartTime            *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeBackupsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupsRequest) GoString() string {
	return s.String()
}

func (s *DescribeBackupsRequest) SetBackupId(v string) *DescribeBackupsRequest {
	s.BackupId = &v
	return s
}

func (s *DescribeBackupsRequest) SetBackupMode(v string) *DescribeBackupsRequest {
	s.BackupMode = &v
	return s
}

func (s *DescribeBackupsRequest) SetBackupRegion(v string) *DescribeBackupsRequest {
	s.BackupRegion = &v
	return s
}

func (s *DescribeBackupsRequest) SetBackupStatus(v string) *DescribeBackupsRequest {
	s.BackupStatus = &v
	return s
}

func (s *DescribeBackupsRequest) SetDBClusterId(v string) *DescribeBackupsRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeBackupsRequest) SetEndTime(v string) *DescribeBackupsRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeBackupsRequest) SetOwnerAccount(v string) *DescribeBackupsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeBackupsRequest) SetOwnerId(v int64) *DescribeBackupsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeBackupsRequest) SetPageNumber(v int32) *DescribeBackupsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeBackupsRequest) SetPageSize(v int32) *DescribeBackupsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeBackupsRequest) SetResourceOwnerAccount(v string) *DescribeBackupsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeBackupsRequest) SetResourceOwnerId(v int64) *DescribeBackupsRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeBackupsRequest) SetStartTime(v string) *DescribeBackupsRequest {
	s.StartTime = &v
	return s
}

type DescribeBackupsResponseBody struct {
	Items            *DescribeBackupsResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	PageNumber       *string                           `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *string                           `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                           `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *string                           `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeBackupsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeBackupsResponseBody) SetItems(v *DescribeBackupsResponseBodyItems) *DescribeBackupsResponseBody {
	s.Items = v
	return s
}

func (s *DescribeBackupsResponseBody) SetPageNumber(v string) *DescribeBackupsResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeBackupsResponseBody) SetPageRecordCount(v string) *DescribeBackupsResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeBackupsResponseBody) SetRequestId(v string) *DescribeBackupsResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeBackupsResponseBody) SetTotalRecordCount(v string) *DescribeBackupsResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeBackupsResponseBodyItems struct {
	Backup []*DescribeBackupsResponseBodyItemsBackup `json:"Backup,omitempty" xml:"Backup,omitempty" type:"Repeated"`
}

func (s DescribeBackupsResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupsResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeBackupsResponseBodyItems) SetBackup(v []*DescribeBackupsResponseBodyItemsBackup) *DescribeBackupsResponseBodyItems {
	s.Backup = v
	return s
}

type DescribeBackupsResponseBodyItemsBackup struct {
	BackupEndTime   *string `json:"BackupEndTime,omitempty" xml:"BackupEndTime,omitempty"`
	BackupId        *string `json:"BackupId,omitempty" xml:"BackupId,omitempty"`
	BackupMethod    *string `json:"BackupMethod,omitempty" xml:"BackupMethod,omitempty"`
	BackupMode      *string `json:"BackupMode,omitempty" xml:"BackupMode,omitempty"`
	BackupSetSize   *string `json:"BackupSetSize,omitempty" xml:"BackupSetSize,omitempty"`
	BackupStartTime *string `json:"BackupStartTime,omitempty" xml:"BackupStartTime,omitempty"`
	BackupStatus    *string `json:"BackupStatus,omitempty" xml:"BackupStatus,omitempty"`
	BackupType      *string `json:"BackupType,omitempty" xml:"BackupType,omitempty"`
	BackupsLevel    *string `json:"BackupsLevel,omitempty" xml:"BackupsLevel,omitempty"`
	ConsistentTime  *string `json:"ConsistentTime,omitempty" xml:"ConsistentTime,omitempty"`
	DBClusterId     *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	IsAvail         *string `json:"IsAvail,omitempty" xml:"IsAvail,omitempty"`
}

func (s DescribeBackupsResponseBodyItemsBackup) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupsResponseBodyItemsBackup) GoString() string {
	return s.String()
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupEndTime(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupEndTime = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupId(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupId = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupMethod(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupMethod = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupMode(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupMode = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupSetSize(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupSetSize = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupStartTime(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupStartTime = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupStatus(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupStatus = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupType(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupType = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetBackupsLevel(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.BackupsLevel = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetConsistentTime(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.ConsistentTime = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetDBClusterId(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.DBClusterId = &v
	return s
}

func (s *DescribeBackupsResponseBodyItemsBackup) SetIsAvail(v string) *DescribeBackupsResponseBodyItemsBackup {
	s.IsAvail = &v
	return s
}

type DescribeBackupsResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeBackupsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeBackupsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeBackupsResponse) GoString() string {
	return s.String()
}

func (s *DescribeBackupsResponse) SetHeaders(v map[string]*string) *DescribeBackupsResponse {
	s.Headers = v
	return s
}

func (s *DescribeBackupsResponse) SetStatusCode(v int32) *DescribeBackupsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeBackupsResponse) SetBody(v *DescribeBackupsResponseBody) *DescribeBackupsResponse {
	s.Body = v
	return s
}

type DescribeCharacterSetNameRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeCharacterSetNameRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeCharacterSetNameRequest) GoString() string {
	return s.String()
}

func (s *DescribeCharacterSetNameRequest) SetDBClusterId(v string) *DescribeCharacterSetNameRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeCharacterSetNameRequest) SetOwnerAccount(v string) *DescribeCharacterSetNameRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeCharacterSetNameRequest) SetOwnerId(v int64) *DescribeCharacterSetNameRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeCharacterSetNameRequest) SetRegionId(v string) *DescribeCharacterSetNameRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeCharacterSetNameRequest) SetResourceOwnerAccount(v string) *DescribeCharacterSetNameRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeCharacterSetNameRequest) SetResourceOwnerId(v int64) *DescribeCharacterSetNameRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeCharacterSetNameResponseBody struct {
	CharacterSetNameItems *DescribeCharacterSetNameResponseBodyCharacterSetNameItems `json:"CharacterSetNameItems,omitempty" xml:"CharacterSetNameItems,omitempty" type:"Struct"`
	Engine                *string                                                    `json:"Engine,omitempty" xml:"Engine,omitempty"`
	RequestId             *string                                                    `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeCharacterSetNameResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeCharacterSetNameResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeCharacterSetNameResponseBody) SetCharacterSetNameItems(v *DescribeCharacterSetNameResponseBodyCharacterSetNameItems) *DescribeCharacterSetNameResponseBody {
	s.CharacterSetNameItems = v
	return s
}

func (s *DescribeCharacterSetNameResponseBody) SetEngine(v string) *DescribeCharacterSetNameResponseBody {
	s.Engine = &v
	return s
}

func (s *DescribeCharacterSetNameResponseBody) SetRequestId(v string) *DescribeCharacterSetNameResponseBody {
	s.RequestId = &v
	return s
}

type DescribeCharacterSetNameResponseBodyCharacterSetNameItems struct {
	CharacterSetName []*string `json:"CharacterSetName,omitempty" xml:"CharacterSetName,omitempty" type:"Repeated"`
}

func (s DescribeCharacterSetNameResponseBodyCharacterSetNameItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeCharacterSetNameResponseBodyCharacterSetNameItems) GoString() string {
	return s.String()
}

func (s *DescribeCharacterSetNameResponseBodyCharacterSetNameItems) SetCharacterSetName(v []*string) *DescribeCharacterSetNameResponseBodyCharacterSetNameItems {
	s.CharacterSetName = v
	return s
}

type DescribeCharacterSetNameResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeCharacterSetNameResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeCharacterSetNameResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeCharacterSetNameResponse) GoString() string {
	return s.String()
}

func (s *DescribeCharacterSetNameResponse) SetHeaders(v map[string]*string) *DescribeCharacterSetNameResponse {
	s.Headers = v
	return s
}

func (s *DescribeCharacterSetNameResponse) SetStatusCode(v int32) *DescribeCharacterSetNameResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeCharacterSetNameResponse) SetBody(v *DescribeCharacterSetNameResponseBody) *DescribeCharacterSetNameResponse {
	s.Body = v
	return s
}

type DescribeDBClusterAccessWhitelistRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterAccessWhitelistRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAccessWhitelistRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAccessWhitelistRequest) SetDBClusterId(v string) *DescribeDBClusterAccessWhitelistRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterAccessWhitelistRequest) SetOwnerAccount(v string) *DescribeDBClusterAccessWhitelistRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterAccessWhitelistRequest) SetOwnerId(v int64) *DescribeDBClusterAccessWhitelistRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterAccessWhitelistRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterAccessWhitelistRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterAccessWhitelistRequest) SetResourceOwnerId(v int64) *DescribeDBClusterAccessWhitelistRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterAccessWhitelistResponseBody struct {
	DBClusterSecurityGroups *DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroups `json:"DBClusterSecurityGroups,omitempty" xml:"DBClusterSecurityGroups,omitempty" type:"Struct"`
	Items                   *DescribeDBClusterAccessWhitelistResponseBodyItems                   `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	RequestId               *string                                                              `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDBClusterAccessWhitelistResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAccessWhitelistResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAccessWhitelistResponseBody) SetDBClusterSecurityGroups(v *DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroups) *DescribeDBClusterAccessWhitelistResponseBody {
	s.DBClusterSecurityGroups = v
	return s
}

func (s *DescribeDBClusterAccessWhitelistResponseBody) SetItems(v *DescribeDBClusterAccessWhitelistResponseBodyItems) *DescribeDBClusterAccessWhitelistResponseBody {
	s.Items = v
	return s
}

func (s *DescribeDBClusterAccessWhitelistResponseBody) SetRequestId(v string) *DescribeDBClusterAccessWhitelistResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroups struct {
	DBClusterSecurityGroup []*DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup `json:"DBClusterSecurityGroup,omitempty" xml:"DBClusterSecurityGroup,omitempty" type:"Repeated"`
}

func (s DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroups) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroups) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroups) SetDBClusterSecurityGroup(v []*DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup) *DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroups {
	s.DBClusterSecurityGroup = v
	return s
}

type DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup struct {
	SecurityGroupId   *string `json:"SecurityGroupId,omitempty" xml:"SecurityGroupId,omitempty"`
	SecurityGroupName *string `json:"SecurityGroupName,omitempty" xml:"SecurityGroupName,omitempty"`
}

func (s DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup) SetSecurityGroupId(v string) *DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup) SetSecurityGroupName(v string) *DescribeDBClusterAccessWhitelistResponseBodyDBClusterSecurityGroupsDBClusterSecurityGroup {
	s.SecurityGroupName = &v
	return s
}

type DescribeDBClusterAccessWhitelistResponseBodyItems struct {
	DBClusterIPArray []*DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray `json:"DBClusterIPArray,omitempty" xml:"DBClusterIPArray,omitempty" type:"Repeated"`
}

func (s DescribeDBClusterAccessWhitelistResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAccessWhitelistResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAccessWhitelistResponseBodyItems) SetDBClusterIPArray(v []*DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray) *DescribeDBClusterAccessWhitelistResponseBodyItems {
	s.DBClusterIPArray = v
	return s
}

type DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray struct {
	DBClusterIPArrayAttribute *string `json:"DBClusterIPArrayAttribute,omitempty" xml:"DBClusterIPArrayAttribute,omitempty"`
	DBClusterIPArrayName      *string `json:"DBClusterIPArrayName,omitempty" xml:"DBClusterIPArrayName,omitempty"`
	SecurityIps               *string `json:"SecurityIps,omitempty" xml:"SecurityIps,omitempty"`
}

func (s DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray) SetDBClusterIPArrayAttribute(v string) *DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray {
	s.DBClusterIPArrayAttribute = &v
	return s
}

func (s *DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray) SetDBClusterIPArrayName(v string) *DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray {
	s.DBClusterIPArrayName = &v
	return s
}

func (s *DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray) SetSecurityIps(v string) *DescribeDBClusterAccessWhitelistResponseBodyItemsDBClusterIPArray {
	s.SecurityIps = &v
	return s
}

type DescribeDBClusterAccessWhitelistResponse struct {
	Headers    map[string]*string                            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterAccessWhitelistResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterAccessWhitelistResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAccessWhitelistResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAccessWhitelistResponse) SetHeaders(v map[string]*string) *DescribeDBClusterAccessWhitelistResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterAccessWhitelistResponse) SetStatusCode(v int32) *DescribeDBClusterAccessWhitelistResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterAccessWhitelistResponse) SetBody(v *DescribeDBClusterAccessWhitelistResponseBody) *DescribeDBClusterAccessWhitelistResponse {
	s.Body = v
	return s
}

type DescribeDBClusterAttributeRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterAttributeRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAttributeRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAttributeRequest) SetDBClusterId(v string) *DescribeDBClusterAttributeRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterAttributeRequest) SetOwnerAccount(v string) *DescribeDBClusterAttributeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterAttributeRequest) SetOwnerId(v int64) *DescribeDBClusterAttributeRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterAttributeRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterAttributeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterAttributeRequest) SetResourceOwnerId(v int64) *DescribeDBClusterAttributeRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterAttributeResponseBody struct {
	BlktagTotal               *int64                                           `json:"BlktagTotal,omitempty" xml:"BlktagTotal,omitempty"`
	BlktagUsed                *int64                                           `json:"BlktagUsed,omitempty" xml:"BlktagUsed,omitempty"`
	Category                  *string                                          `json:"Category,omitempty" xml:"Category,omitempty"`
	CreationTime              *string                                          `json:"CreationTime,omitempty" xml:"CreationTime,omitempty"`
	DBClusterDescription      *string                                          `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBClusterId               *string                                          `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBClusterNetworkType      *string                                          `json:"DBClusterNetworkType,omitempty" xml:"DBClusterNetworkType,omitempty"`
	DBClusterStatus           *string                                          `json:"DBClusterStatus,omitempty" xml:"DBClusterStatus,omitempty"`
	DBNodes                   []*DescribeDBClusterAttributeResponseBodyDBNodes `json:"DBNodes,omitempty" xml:"DBNodes,omitempty" type:"Repeated"`
	DBType                    *string                                          `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion                 *string                                          `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	DBVersionStatus           *string                                          `json:"DBVersionStatus,omitempty" xml:"DBVersionStatus,omitempty"`
	DataLevel1BackupChainSize *int64                                           `json:"DataLevel1BackupChainSize,omitempty" xml:"DataLevel1BackupChainSize,omitempty"`
	DeletionLock              *int32                                           `json:"DeletionLock,omitempty" xml:"DeletionLock,omitempty"`
	Engine                    *string                                          `json:"Engine,omitempty" xml:"Engine,omitempty"`
	ExpireTime                *string                                          `json:"ExpireTime,omitempty" xml:"ExpireTime,omitempty"`
	Expired                   *string                                          `json:"Expired,omitempty" xml:"Expired,omitempty"`
	InodeTotal                *int64                                           `json:"InodeTotal,omitempty" xml:"InodeTotal,omitempty"`
	InodeUsed                 *int64                                           `json:"InodeUsed,omitempty" xml:"InodeUsed,omitempty"`
	IsLatestVersion           *bool                                            `json:"IsLatestVersion,omitempty" xml:"IsLatestVersion,omitempty"`
	IsProxyLatestVersion      *bool                                            `json:"IsProxyLatestVersion,omitempty" xml:"IsProxyLatestVersion,omitempty"`
	LockMode                  *string                                          `json:"LockMode,omitempty" xml:"LockMode,omitempty"`
	MaintainTime              *string                                          `json:"MaintainTime,omitempty" xml:"MaintainTime,omitempty"`
	PayType                   *string                                          `json:"PayType,omitempty" xml:"PayType,omitempty"`
	ProxyCpuCores             *string                                          `json:"ProxyCpuCores,omitempty" xml:"ProxyCpuCores,omitempty"`
	ProxyStandardCpuCores     *string                                          `json:"ProxyStandardCpuCores,omitempty" xml:"ProxyStandardCpuCores,omitempty"`
	ProxyStatus               *string                                          `json:"ProxyStatus,omitempty" xml:"ProxyStatus,omitempty"`
	ProxyType                 *string                                          `json:"ProxyType,omitempty" xml:"ProxyType,omitempty"`
	RegionId                  *string                                          `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	RequestId                 *string                                          `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	ResourceGroupId           *string                                          `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	SQLSize                   *int64                                           `json:"SQLSize,omitempty" xml:"SQLSize,omitempty"`
	StorageMax                *int64                                           `json:"StorageMax,omitempty" xml:"StorageMax,omitempty"`
	StoragePayType            *string                                          `json:"StoragePayType,omitempty" xml:"StoragePayType,omitempty"`
	StorageSpace              *int64                                           `json:"StorageSpace,omitempty" xml:"StorageSpace,omitempty"`
	StorageType               *string                                          `json:"StorageType,omitempty" xml:"StorageType,omitempty"`
	StorageUsed               *int64                                           `json:"StorageUsed,omitempty" xml:"StorageUsed,omitempty"`
	SubCategory               *string                                          `json:"SubCategory,omitempty" xml:"SubCategory,omitempty"`
	Tags                      []*DescribeDBClusterAttributeResponseBodyTags    `json:"Tags,omitempty" xml:"Tags,omitempty" type:"Repeated"`
	VPCId                     *string                                          `json:"VPCId,omitempty" xml:"VPCId,omitempty"`
	VSwitchId                 *string                                          `json:"VSwitchId,omitempty" xml:"VSwitchId,omitempty"`
	ZoneIds                   *string                                          `json:"ZoneIds,omitempty" xml:"ZoneIds,omitempty"`
}

func (s DescribeDBClusterAttributeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAttributeResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAttributeResponseBody) SetBlktagTotal(v int64) *DescribeDBClusterAttributeResponseBody {
	s.BlktagTotal = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetBlktagUsed(v int64) *DescribeDBClusterAttributeResponseBody {
	s.BlktagUsed = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetCategory(v string) *DescribeDBClusterAttributeResponseBody {
	s.Category = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetCreationTime(v string) *DescribeDBClusterAttributeResponseBody {
	s.CreationTime = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDBClusterDescription(v string) *DescribeDBClusterAttributeResponseBody {
	s.DBClusterDescription = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDBClusterId(v string) *DescribeDBClusterAttributeResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDBClusterNetworkType(v string) *DescribeDBClusterAttributeResponseBody {
	s.DBClusterNetworkType = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDBClusterStatus(v string) *DescribeDBClusterAttributeResponseBody {
	s.DBClusterStatus = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDBNodes(v []*DescribeDBClusterAttributeResponseBodyDBNodes) *DescribeDBClusterAttributeResponseBody {
	s.DBNodes = v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDBType(v string) *DescribeDBClusterAttributeResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDBVersion(v string) *DescribeDBClusterAttributeResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDBVersionStatus(v string) *DescribeDBClusterAttributeResponseBody {
	s.DBVersionStatus = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDataLevel1BackupChainSize(v int64) *DescribeDBClusterAttributeResponseBody {
	s.DataLevel1BackupChainSize = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetDeletionLock(v int32) *DescribeDBClusterAttributeResponseBody {
	s.DeletionLock = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetEngine(v string) *DescribeDBClusterAttributeResponseBody {
	s.Engine = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetExpireTime(v string) *DescribeDBClusterAttributeResponseBody {
	s.ExpireTime = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetExpired(v string) *DescribeDBClusterAttributeResponseBody {
	s.Expired = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetInodeTotal(v int64) *DescribeDBClusterAttributeResponseBody {
	s.InodeTotal = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetInodeUsed(v int64) *DescribeDBClusterAttributeResponseBody {
	s.InodeUsed = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetIsLatestVersion(v bool) *DescribeDBClusterAttributeResponseBody {
	s.IsLatestVersion = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetIsProxyLatestVersion(v bool) *DescribeDBClusterAttributeResponseBody {
	s.IsProxyLatestVersion = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetLockMode(v string) *DescribeDBClusterAttributeResponseBody {
	s.LockMode = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetMaintainTime(v string) *DescribeDBClusterAttributeResponseBody {
	s.MaintainTime = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetPayType(v string) *DescribeDBClusterAttributeResponseBody {
	s.PayType = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetProxyCpuCores(v string) *DescribeDBClusterAttributeResponseBody {
	s.ProxyCpuCores = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetProxyStandardCpuCores(v string) *DescribeDBClusterAttributeResponseBody {
	s.ProxyStandardCpuCores = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetProxyStatus(v string) *DescribeDBClusterAttributeResponseBody {
	s.ProxyStatus = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetProxyType(v string) *DescribeDBClusterAttributeResponseBody {
	s.ProxyType = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetRegionId(v string) *DescribeDBClusterAttributeResponseBody {
	s.RegionId = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetRequestId(v string) *DescribeDBClusterAttributeResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetResourceGroupId(v string) *DescribeDBClusterAttributeResponseBody {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetSQLSize(v int64) *DescribeDBClusterAttributeResponseBody {
	s.SQLSize = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetStorageMax(v int64) *DescribeDBClusterAttributeResponseBody {
	s.StorageMax = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetStoragePayType(v string) *DescribeDBClusterAttributeResponseBody {
	s.StoragePayType = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetStorageSpace(v int64) *DescribeDBClusterAttributeResponseBody {
	s.StorageSpace = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetStorageType(v string) *DescribeDBClusterAttributeResponseBody {
	s.StorageType = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetStorageUsed(v int64) *DescribeDBClusterAttributeResponseBody {
	s.StorageUsed = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetSubCategory(v string) *DescribeDBClusterAttributeResponseBody {
	s.SubCategory = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetTags(v []*DescribeDBClusterAttributeResponseBodyTags) *DescribeDBClusterAttributeResponseBody {
	s.Tags = v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetVPCId(v string) *DescribeDBClusterAttributeResponseBody {
	s.VPCId = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetVSwitchId(v string) *DescribeDBClusterAttributeResponseBody {
	s.VSwitchId = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBody) SetZoneIds(v string) *DescribeDBClusterAttributeResponseBody {
	s.ZoneIds = &v
	return s
}

type DescribeDBClusterAttributeResponseBodyDBNodes struct {
	AddedCpuCores    *string `json:"AddedCpuCores,omitempty" xml:"AddedCpuCores,omitempty"`
	CreationTime     *string `json:"CreationTime,omitempty" xml:"CreationTime,omitempty"`
	DBNodeClass      *string `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
	DBNodeId         *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	DBNodeRole       *string `json:"DBNodeRole,omitempty" xml:"DBNodeRole,omitempty"`
	DBNodeStatus     *string `json:"DBNodeStatus,omitempty" xml:"DBNodeStatus,omitempty"`
	FailoverPriority *int32  `json:"FailoverPriority,omitempty" xml:"FailoverPriority,omitempty"`
	HotReplicaMode   *string `json:"HotReplicaMode,omitempty" xml:"HotReplicaMode,omitempty"`
	ImciSwitch       *string `json:"ImciSwitch,omitempty" xml:"ImciSwitch,omitempty"`
	MasterId         *string `json:"MasterId,omitempty" xml:"MasterId,omitempty"`
	MaxConnections   *int32  `json:"MaxConnections,omitempty" xml:"MaxConnections,omitempty"`
	MaxIOPS          *int32  `json:"MaxIOPS,omitempty" xml:"MaxIOPS,omitempty"`
	ZoneId           *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s DescribeDBClusterAttributeResponseBodyDBNodes) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAttributeResponseBodyDBNodes) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetAddedCpuCores(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.AddedCpuCores = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetCreationTime(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.CreationTime = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetDBNodeClass(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.DBNodeClass = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetDBNodeId(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.DBNodeId = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetDBNodeRole(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.DBNodeRole = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetDBNodeStatus(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.DBNodeStatus = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetFailoverPriority(v int32) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.FailoverPriority = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetHotReplicaMode(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.HotReplicaMode = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetImciSwitch(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.ImciSwitch = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetMasterId(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.MasterId = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetMaxConnections(v int32) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.MaxConnections = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetMaxIOPS(v int32) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.MaxIOPS = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyDBNodes) SetZoneId(v string) *DescribeDBClusterAttributeResponseBodyDBNodes {
	s.ZoneId = &v
	return s
}

type DescribeDBClusterAttributeResponseBodyTags struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s DescribeDBClusterAttributeResponseBodyTags) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAttributeResponseBodyTags) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAttributeResponseBodyTags) SetKey(v string) *DescribeDBClusterAttributeResponseBodyTags {
	s.Key = &v
	return s
}

func (s *DescribeDBClusterAttributeResponseBodyTags) SetValue(v string) *DescribeDBClusterAttributeResponseBodyTags {
	s.Value = &v
	return s
}

type DescribeDBClusterAttributeResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterAttributeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterAttributeResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAttributeResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAttributeResponse) SetHeaders(v map[string]*string) *DescribeDBClusterAttributeResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterAttributeResponse) SetStatusCode(v int32) *DescribeDBClusterAttributeResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterAttributeResponse) SetBody(v *DescribeDBClusterAttributeResponseBody) *DescribeDBClusterAttributeResponse {
	s.Body = v
	return s
}

type DescribeDBClusterAuditLogCollectorRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterAuditLogCollectorRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAuditLogCollectorRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAuditLogCollectorRequest) SetDBClusterId(v string) *DescribeDBClusterAuditLogCollectorRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterAuditLogCollectorRequest) SetOwnerAccount(v string) *DescribeDBClusterAuditLogCollectorRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterAuditLogCollectorRequest) SetOwnerId(v int64) *DescribeDBClusterAuditLogCollectorRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterAuditLogCollectorRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterAuditLogCollectorRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterAuditLogCollectorRequest) SetResourceOwnerId(v int64) *DescribeDBClusterAuditLogCollectorRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterAuditLogCollectorResponseBody struct {
	CollectorStatus *string `json:"CollectorStatus,omitempty" xml:"CollectorStatus,omitempty"`
	RequestId       *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDBClusterAuditLogCollectorResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAuditLogCollectorResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAuditLogCollectorResponseBody) SetCollectorStatus(v string) *DescribeDBClusterAuditLogCollectorResponseBody {
	s.CollectorStatus = &v
	return s
}

func (s *DescribeDBClusterAuditLogCollectorResponseBody) SetRequestId(v string) *DescribeDBClusterAuditLogCollectorResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDBClusterAuditLogCollectorResponse struct {
	Headers    map[string]*string                              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterAuditLogCollectorResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterAuditLogCollectorResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAuditLogCollectorResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAuditLogCollectorResponse) SetHeaders(v map[string]*string) *DescribeDBClusterAuditLogCollectorResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterAuditLogCollectorResponse) SetStatusCode(v int32) *DescribeDBClusterAuditLogCollectorResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterAuditLogCollectorResponse) SetBody(v *DescribeDBClusterAuditLogCollectorResponseBody) *DescribeDBClusterAuditLogCollectorResponse {
	s.Body = v
	return s
}

type DescribeDBClusterAvailableResourcesRequest struct {
	DBNodeClass          *string `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
	DBType               *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion            *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PayType              *string `json:"PayType,omitempty" xml:"PayType,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	ZoneId               *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s DescribeDBClusterAvailableResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAvailableResourcesRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetDBNodeClass(v string) *DescribeDBClusterAvailableResourcesRequest {
	s.DBNodeClass = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetDBType(v string) *DescribeDBClusterAvailableResourcesRequest {
	s.DBType = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetDBVersion(v string) *DescribeDBClusterAvailableResourcesRequest {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetOwnerAccount(v string) *DescribeDBClusterAvailableResourcesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetOwnerId(v int64) *DescribeDBClusterAvailableResourcesRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetPayType(v string) *DescribeDBClusterAvailableResourcesRequest {
	s.PayType = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetRegionId(v string) *DescribeDBClusterAvailableResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterAvailableResourcesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetResourceOwnerId(v int64) *DescribeDBClusterAvailableResourcesRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesRequest) SetZoneId(v string) *DescribeDBClusterAvailableResourcesRequest {
	s.ZoneId = &v
	return s
}

type DescribeDBClusterAvailableResourcesResponseBody struct {
	AvailableZones []*DescribeDBClusterAvailableResourcesResponseBodyAvailableZones `json:"AvailableZones,omitempty" xml:"AvailableZones,omitempty" type:"Repeated"`
	RequestId      *string                                                          `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDBClusterAvailableResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAvailableResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAvailableResourcesResponseBody) SetAvailableZones(v []*DescribeDBClusterAvailableResourcesResponseBodyAvailableZones) *DescribeDBClusterAvailableResourcesResponseBody {
	s.AvailableZones = v
	return s
}

func (s *DescribeDBClusterAvailableResourcesResponseBody) SetRequestId(v string) *DescribeDBClusterAvailableResourcesResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDBClusterAvailableResourcesResponseBodyAvailableZones struct {
	RegionId         *string                                                                          `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	SupportedEngines []*DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines `json:"SupportedEngines,omitempty" xml:"SupportedEngines,omitempty" type:"Repeated"`
	ZoneId           *string                                                                          `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s DescribeDBClusterAvailableResourcesResponseBodyAvailableZones) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAvailableResourcesResponseBodyAvailableZones) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAvailableResourcesResponseBodyAvailableZones) SetRegionId(v string) *DescribeDBClusterAvailableResourcesResponseBodyAvailableZones {
	s.RegionId = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesResponseBodyAvailableZones) SetSupportedEngines(v []*DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines) *DescribeDBClusterAvailableResourcesResponseBodyAvailableZones {
	s.SupportedEngines = v
	return s
}

func (s *DescribeDBClusterAvailableResourcesResponseBodyAvailableZones) SetZoneId(v string) *DescribeDBClusterAvailableResourcesResponseBodyAvailableZones {
	s.ZoneId = &v
	return s
}

type DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines struct {
	AvailableResources []*DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources `json:"AvailableResources,omitempty" xml:"AvailableResources,omitempty" type:"Repeated"`
	Engine             *string                                                                                            `json:"Engine,omitempty" xml:"Engine,omitempty"`
}

func (s DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines) SetAvailableResources(v []*DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources) *DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines {
	s.AvailableResources = v
	return s
}

func (s *DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines) SetEngine(v string) *DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEngines {
	s.Engine = &v
	return s
}

type DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources struct {
	Category    *string `json:"Category,omitempty" xml:"Category,omitempty"`
	DBNodeClass *string `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
}

func (s DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources) SetCategory(v string) *DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources {
	s.Category = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources) SetDBNodeClass(v string) *DescribeDBClusterAvailableResourcesResponseBodyAvailableZonesSupportedEnginesAvailableResources {
	s.DBNodeClass = &v
	return s
}

type DescribeDBClusterAvailableResourcesResponse struct {
	Headers    map[string]*string                               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterAvailableResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterAvailableResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterAvailableResourcesResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterAvailableResourcesResponse) SetHeaders(v map[string]*string) *DescribeDBClusterAvailableResourcesResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterAvailableResourcesResponse) SetStatusCode(v int32) *DescribeDBClusterAvailableResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterAvailableResourcesResponse) SetBody(v *DescribeDBClusterAvailableResourcesResponseBody) *DescribeDBClusterAvailableResourcesResponse {
	s.Body = v
	return s
}

type DescribeDBClusterEndpointsRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBEndpointId         *string `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterEndpointsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterEndpointsRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterEndpointsRequest) SetDBClusterId(v string) *DescribeDBClusterEndpointsRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterEndpointsRequest) SetDBEndpointId(v string) *DescribeDBClusterEndpointsRequest {
	s.DBEndpointId = &v
	return s
}

func (s *DescribeDBClusterEndpointsRequest) SetOwnerAccount(v string) *DescribeDBClusterEndpointsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterEndpointsRequest) SetOwnerId(v int64) *DescribeDBClusterEndpointsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterEndpointsRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterEndpointsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterEndpointsRequest) SetResourceOwnerId(v int64) *DescribeDBClusterEndpointsRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterEndpointsResponseBody struct {
	Items     []*DescribeDBClusterEndpointsResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Repeated"`
	RequestId *string                                        `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDBClusterEndpointsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterEndpointsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterEndpointsResponseBody) SetItems(v []*DescribeDBClusterEndpointsResponseBodyItems) *DescribeDBClusterEndpointsResponseBody {
	s.Items = v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBody) SetRequestId(v string) *DescribeDBClusterEndpointsResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDBClusterEndpointsResponseBodyItems struct {
	AddressItems          []*DescribeDBClusterEndpointsResponseBodyItemsAddressItems `json:"AddressItems,omitempty" xml:"AddressItems,omitempty" type:"Repeated"`
	AutoAddNewNodes       *string                                                    `json:"AutoAddNewNodes,omitempty" xml:"AutoAddNewNodes,omitempty"`
	DBEndpointDescription *string                                                    `json:"DBEndpointDescription,omitempty" xml:"DBEndpointDescription,omitempty"`
	DBEndpointId          *string                                                    `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	EndpointConfig        *string                                                    `json:"EndpointConfig,omitempty" xml:"EndpointConfig,omitempty"`
	EndpointType          *string                                                    `json:"EndpointType,omitempty" xml:"EndpointType,omitempty"`
	NodeWithRoles         *string                                                    `json:"NodeWithRoles,omitempty" xml:"NodeWithRoles,omitempty"`
	Nodes                 *string                                                    `json:"Nodes,omitempty" xml:"Nodes,omitempty"`
	ReadWriteMode         *string                                                    `json:"ReadWriteMode,omitempty" xml:"ReadWriteMode,omitempty"`
}

func (s DescribeDBClusterEndpointsResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterEndpointsResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetAddressItems(v []*DescribeDBClusterEndpointsResponseBodyItemsAddressItems) *DescribeDBClusterEndpointsResponseBodyItems {
	s.AddressItems = v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetAutoAddNewNodes(v string) *DescribeDBClusterEndpointsResponseBodyItems {
	s.AutoAddNewNodes = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetDBEndpointDescription(v string) *DescribeDBClusterEndpointsResponseBodyItems {
	s.DBEndpointDescription = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetDBEndpointId(v string) *DescribeDBClusterEndpointsResponseBodyItems {
	s.DBEndpointId = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetEndpointConfig(v string) *DescribeDBClusterEndpointsResponseBodyItems {
	s.EndpointConfig = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetEndpointType(v string) *DescribeDBClusterEndpointsResponseBodyItems {
	s.EndpointType = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetNodeWithRoles(v string) *DescribeDBClusterEndpointsResponseBodyItems {
	s.NodeWithRoles = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetNodes(v string) *DescribeDBClusterEndpointsResponseBodyItems {
	s.Nodes = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItems) SetReadWriteMode(v string) *DescribeDBClusterEndpointsResponseBodyItems {
	s.ReadWriteMode = &v
	return s
}

type DescribeDBClusterEndpointsResponseBodyItemsAddressItems struct {
	ConnectionString            *string `json:"ConnectionString,omitempty" xml:"ConnectionString,omitempty"`
	IPAddress                   *string `json:"IPAddress,omitempty" xml:"IPAddress,omitempty"`
	NetType                     *string `json:"NetType,omitempty" xml:"NetType,omitempty"`
	Port                        *string `json:"Port,omitempty" xml:"Port,omitempty"`
	PrivateZoneConnectionString *string `json:"PrivateZoneConnectionString,omitempty" xml:"PrivateZoneConnectionString,omitempty"`
	VPCId                       *string `json:"VPCId,omitempty" xml:"VPCId,omitempty"`
	VSwitchId                   *string `json:"VSwitchId,omitempty" xml:"VSwitchId,omitempty"`
	VpcInstanceId               *string `json:"VpcInstanceId,omitempty" xml:"VpcInstanceId,omitempty"`
}

func (s DescribeDBClusterEndpointsResponseBodyItemsAddressItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterEndpointsResponseBodyItemsAddressItems) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterEndpointsResponseBodyItemsAddressItems) SetConnectionString(v string) *DescribeDBClusterEndpointsResponseBodyItemsAddressItems {
	s.ConnectionString = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItemsAddressItems) SetIPAddress(v string) *DescribeDBClusterEndpointsResponseBodyItemsAddressItems {
	s.IPAddress = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItemsAddressItems) SetNetType(v string) *DescribeDBClusterEndpointsResponseBodyItemsAddressItems {
	s.NetType = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItemsAddressItems) SetPort(v string) *DescribeDBClusterEndpointsResponseBodyItemsAddressItems {
	s.Port = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItemsAddressItems) SetPrivateZoneConnectionString(v string) *DescribeDBClusterEndpointsResponseBodyItemsAddressItems {
	s.PrivateZoneConnectionString = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItemsAddressItems) SetVPCId(v string) *DescribeDBClusterEndpointsResponseBodyItemsAddressItems {
	s.VPCId = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItemsAddressItems) SetVSwitchId(v string) *DescribeDBClusterEndpointsResponseBodyItemsAddressItems {
	s.VSwitchId = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponseBodyItemsAddressItems) SetVpcInstanceId(v string) *DescribeDBClusterEndpointsResponseBodyItemsAddressItems {
	s.VpcInstanceId = &v
	return s
}

type DescribeDBClusterEndpointsResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterEndpointsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterEndpointsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterEndpointsResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterEndpointsResponse) SetHeaders(v map[string]*string) *DescribeDBClusterEndpointsResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterEndpointsResponse) SetStatusCode(v int32) *DescribeDBClusterEndpointsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterEndpointsResponse) SetBody(v *DescribeDBClusterEndpointsResponseBody) *DescribeDBClusterEndpointsResponse {
	s.Body = v
	return s
}

type DescribeDBClusterMigrationRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterMigrationRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMigrationRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMigrationRequest) SetDBClusterId(v string) *DescribeDBClusterMigrationRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterMigrationRequest) SetOwnerAccount(v string) *DescribeDBClusterMigrationRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterMigrationRequest) SetOwnerId(v int64) *DescribeDBClusterMigrationRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterMigrationRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterMigrationRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterMigrationRequest) SetResourceOwnerId(v int64) *DescribeDBClusterMigrationRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterMigrationResponseBody struct {
	Comment                *string                                                        `json:"Comment,omitempty" xml:"Comment,omitempty"`
	DBClusterEndpointList  []*DescribeDBClusterMigrationResponseBodyDBClusterEndpointList `json:"DBClusterEndpointList,omitempty" xml:"DBClusterEndpointList,omitempty" type:"Repeated"`
	DBClusterId            *string                                                        `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBClusterReadWriteMode *string                                                        `json:"DBClusterReadWriteMode,omitempty" xml:"DBClusterReadWriteMode,omitempty"`
	DelayedSeconds         *int32                                                         `json:"DelayedSeconds,omitempty" xml:"DelayedSeconds,omitempty"`
	DtsInstanceId          *string                                                        `json:"DtsInstanceId,omitempty" xml:"DtsInstanceId,omitempty"`
	ExpiredTime            *string                                                        `json:"ExpiredTime,omitempty" xml:"ExpiredTime,omitempty"`
	MigrationStatus        *string                                                        `json:"MigrationStatus,omitempty" xml:"MigrationStatus,omitempty"`
	RdsEndpointList        []*DescribeDBClusterMigrationResponseBodyRdsEndpointList       `json:"RdsEndpointList,omitempty" xml:"RdsEndpointList,omitempty" type:"Repeated"`
	RdsReadWriteMode       *string                                                        `json:"RdsReadWriteMode,omitempty" xml:"RdsReadWriteMode,omitempty"`
	RequestId              *string                                                        `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	SourceRDSDBInstanceId  *string                                                        `json:"SourceRDSDBInstanceId,omitempty" xml:"SourceRDSDBInstanceId,omitempty"`
	Topologies             *string                                                        `json:"Topologies,omitempty" xml:"Topologies,omitempty"`
}

func (s DescribeDBClusterMigrationResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMigrationResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMigrationResponseBody) SetComment(v string) *DescribeDBClusterMigrationResponseBody {
	s.Comment = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetDBClusterEndpointList(v []*DescribeDBClusterMigrationResponseBodyDBClusterEndpointList) *DescribeDBClusterMigrationResponseBody {
	s.DBClusterEndpointList = v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetDBClusterId(v string) *DescribeDBClusterMigrationResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetDBClusterReadWriteMode(v string) *DescribeDBClusterMigrationResponseBody {
	s.DBClusterReadWriteMode = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetDelayedSeconds(v int32) *DescribeDBClusterMigrationResponseBody {
	s.DelayedSeconds = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetDtsInstanceId(v string) *DescribeDBClusterMigrationResponseBody {
	s.DtsInstanceId = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetExpiredTime(v string) *DescribeDBClusterMigrationResponseBody {
	s.ExpiredTime = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetMigrationStatus(v string) *DescribeDBClusterMigrationResponseBody {
	s.MigrationStatus = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetRdsEndpointList(v []*DescribeDBClusterMigrationResponseBodyRdsEndpointList) *DescribeDBClusterMigrationResponseBody {
	s.RdsEndpointList = v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetRdsReadWriteMode(v string) *DescribeDBClusterMigrationResponseBody {
	s.RdsReadWriteMode = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetRequestId(v string) *DescribeDBClusterMigrationResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetSourceRDSDBInstanceId(v string) *DescribeDBClusterMigrationResponseBody {
	s.SourceRDSDBInstanceId = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBody) SetTopologies(v string) *DescribeDBClusterMigrationResponseBody {
	s.Topologies = &v
	return s
}

type DescribeDBClusterMigrationResponseBodyDBClusterEndpointList struct {
	AddressItems []*DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems `json:"AddressItems,omitempty" xml:"AddressItems,omitempty" type:"Repeated"`
	DBEndpointId *string                                                                    `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	EndpointType *string                                                                    `json:"EndpointType,omitempty" xml:"EndpointType,omitempty"`
}

func (s DescribeDBClusterMigrationResponseBodyDBClusterEndpointList) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMigrationResponseBodyDBClusterEndpointList) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointList) SetAddressItems(v []*DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointList {
	s.AddressItems = v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointList) SetDBEndpointId(v string) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointList {
	s.DBEndpointId = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointList) SetEndpointType(v string) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointList {
	s.EndpointType = &v
	return s
}

type DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems struct {
	ConnectionString *string `json:"ConnectionString,omitempty" xml:"ConnectionString,omitempty"`
	IPAddress        *string `json:"IPAddress,omitempty" xml:"IPAddress,omitempty"`
	NetType          *string `json:"NetType,omitempty" xml:"NetType,omitempty"`
	Port             *string `json:"Port,omitempty" xml:"Port,omitempty"`
	VPCId            *string `json:"VPCId,omitempty" xml:"VPCId,omitempty"`
	VSwitchId        *string `json:"VSwitchId,omitempty" xml:"VSwitchId,omitempty"`
}

func (s DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) SetConnectionString(v string) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems {
	s.ConnectionString = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) SetIPAddress(v string) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems {
	s.IPAddress = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) SetNetType(v string) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems {
	s.NetType = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) SetPort(v string) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems {
	s.Port = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) SetVPCId(v string) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems {
	s.VPCId = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems) SetVSwitchId(v string) *DescribeDBClusterMigrationResponseBodyDBClusterEndpointListAddressItems {
	s.VSwitchId = &v
	return s
}

type DescribeDBClusterMigrationResponseBodyRdsEndpointList struct {
	AddressItems []*DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems `json:"AddressItems,omitempty" xml:"AddressItems,omitempty" type:"Repeated"`
	DBEndpointId *string                                                              `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	EndpointType *string                                                              `json:"EndpointType,omitempty" xml:"EndpointType,omitempty"`
}

func (s DescribeDBClusterMigrationResponseBodyRdsEndpointList) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMigrationResponseBodyRdsEndpointList) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointList) SetAddressItems(v []*DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) *DescribeDBClusterMigrationResponseBodyRdsEndpointList {
	s.AddressItems = v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointList) SetDBEndpointId(v string) *DescribeDBClusterMigrationResponseBodyRdsEndpointList {
	s.DBEndpointId = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointList) SetEndpointType(v string) *DescribeDBClusterMigrationResponseBodyRdsEndpointList {
	s.EndpointType = &v
	return s
}

type DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems struct {
	ConnectionString *string `json:"ConnectionString,omitempty" xml:"ConnectionString,omitempty"`
	IPAddress        *string `json:"IPAddress,omitempty" xml:"IPAddress,omitempty"`
	NetType          *string `json:"NetType,omitempty" xml:"NetType,omitempty"`
	Port             *string `json:"Port,omitempty" xml:"Port,omitempty"`
	VPCId            *string `json:"VPCId,omitempty" xml:"VPCId,omitempty"`
	VSwitchId        *string `json:"VSwitchId,omitempty" xml:"VSwitchId,omitempty"`
}

func (s DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) SetConnectionString(v string) *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems {
	s.ConnectionString = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) SetIPAddress(v string) *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems {
	s.IPAddress = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) SetNetType(v string) *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems {
	s.NetType = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) SetPort(v string) *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems {
	s.Port = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) SetVPCId(v string) *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems {
	s.VPCId = &v
	return s
}

func (s *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems) SetVSwitchId(v string) *DescribeDBClusterMigrationResponseBodyRdsEndpointListAddressItems {
	s.VSwitchId = &v
	return s
}

type DescribeDBClusterMigrationResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterMigrationResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterMigrationResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMigrationResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMigrationResponse) SetHeaders(v map[string]*string) *DescribeDBClusterMigrationResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterMigrationResponse) SetStatusCode(v int32) *DescribeDBClusterMigrationResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterMigrationResponse) SetBody(v *DescribeDBClusterMigrationResponseBody) *DescribeDBClusterMigrationResponse {
	s.Body = v
	return s
}

type DescribeDBClusterMonitorRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterMonitorRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMonitorRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMonitorRequest) SetDBClusterId(v string) *DescribeDBClusterMonitorRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterMonitorRequest) SetOwnerAccount(v string) *DescribeDBClusterMonitorRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterMonitorRequest) SetOwnerId(v int64) *DescribeDBClusterMonitorRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterMonitorRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterMonitorRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterMonitorRequest) SetResourceOwnerId(v int64) *DescribeDBClusterMonitorRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterMonitorResponseBody struct {
	Period    *string `json:"Period,omitempty" xml:"Period,omitempty"`
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDBClusterMonitorResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMonitorResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMonitorResponseBody) SetPeriod(v string) *DescribeDBClusterMonitorResponseBody {
	s.Period = &v
	return s
}

func (s *DescribeDBClusterMonitorResponseBody) SetRequestId(v string) *DescribeDBClusterMonitorResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDBClusterMonitorResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterMonitorResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterMonitorResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterMonitorResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterMonitorResponse) SetHeaders(v map[string]*string) *DescribeDBClusterMonitorResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterMonitorResponse) SetStatusCode(v int32) *DescribeDBClusterMonitorResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterMonitorResponse) SetBody(v *DescribeDBClusterMonitorResponseBody) *DescribeDBClusterMonitorResponse {
	s.Body = v
	return s
}

type DescribeDBClusterParametersRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterParametersRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterParametersRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterParametersRequest) SetDBClusterId(v string) *DescribeDBClusterParametersRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterParametersRequest) SetOwnerAccount(v string) *DescribeDBClusterParametersRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterParametersRequest) SetOwnerId(v int64) *DescribeDBClusterParametersRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterParametersRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterParametersRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterParametersRequest) SetResourceOwnerId(v int64) *DescribeDBClusterParametersRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterParametersResponseBody struct {
	DBType            *string                                                   `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion         *string                                                   `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	Engine            *string                                                   `json:"Engine,omitempty" xml:"Engine,omitempty"`
	RequestId         *string                                                   `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	RunningParameters *DescribeDBClusterParametersResponseBodyRunningParameters `json:"RunningParameters,omitempty" xml:"RunningParameters,omitempty" type:"Struct"`
}

func (s DescribeDBClusterParametersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterParametersResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterParametersResponseBody) SetDBType(v string) *DescribeDBClusterParametersResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBody) SetDBVersion(v string) *DescribeDBClusterParametersResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBody) SetEngine(v string) *DescribeDBClusterParametersResponseBody {
	s.Engine = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBody) SetRequestId(v string) *DescribeDBClusterParametersResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBody) SetRunningParameters(v *DescribeDBClusterParametersResponseBodyRunningParameters) *DescribeDBClusterParametersResponseBody {
	s.RunningParameters = v
	return s
}

type DescribeDBClusterParametersResponseBodyRunningParameters struct {
	Parameter []*DescribeDBClusterParametersResponseBodyRunningParametersParameter `json:"Parameter,omitempty" xml:"Parameter,omitempty" type:"Repeated"`
}

func (s DescribeDBClusterParametersResponseBodyRunningParameters) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterParametersResponseBodyRunningParameters) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterParametersResponseBodyRunningParameters) SetParameter(v []*DescribeDBClusterParametersResponseBodyRunningParametersParameter) *DescribeDBClusterParametersResponseBodyRunningParameters {
	s.Parameter = v
	return s
}

type DescribeDBClusterParametersResponseBodyRunningParametersParameter struct {
	CheckingCode          *string `json:"CheckingCode,omitempty" xml:"CheckingCode,omitempty"`
	DataType              *string `json:"DataType,omitempty" xml:"DataType,omitempty"`
	DefaultParameterValue *string `json:"DefaultParameterValue,omitempty" xml:"DefaultParameterValue,omitempty"`
	Factor                *string `json:"Factor,omitempty" xml:"Factor,omitempty"`
	ForceRestart          *bool   `json:"ForceRestart,omitempty" xml:"ForceRestart,omitempty"`
	IsModifiable          *bool   `json:"IsModifiable,omitempty" xml:"IsModifiable,omitempty"`
	IsNodeAvailable       *string `json:"IsNodeAvailable,omitempty" xml:"IsNodeAvailable,omitempty"`
	ParamRelyRule         *string `json:"ParamRelyRule,omitempty" xml:"ParamRelyRule,omitempty"`
	ParameterDescription  *string `json:"ParameterDescription,omitempty" xml:"ParameterDescription,omitempty"`
	ParameterName         *string `json:"ParameterName,omitempty" xml:"ParameterName,omitempty"`
	ParameterStatus       *string `json:"ParameterStatus,omitempty" xml:"ParameterStatus,omitempty"`
	ParameterValue        *string `json:"ParameterValue,omitempty" xml:"ParameterValue,omitempty"`
}

func (s DescribeDBClusterParametersResponseBodyRunningParametersParameter) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterParametersResponseBodyRunningParametersParameter) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetCheckingCode(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.CheckingCode = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetDataType(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.DataType = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetDefaultParameterValue(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.DefaultParameterValue = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetFactor(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.Factor = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetForceRestart(v bool) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.ForceRestart = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetIsModifiable(v bool) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.IsModifiable = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetIsNodeAvailable(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.IsNodeAvailable = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetParamRelyRule(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.ParamRelyRule = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetParameterDescription(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.ParameterDescription = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetParameterName(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.ParameterName = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetParameterStatus(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.ParameterStatus = &v
	return s
}

func (s *DescribeDBClusterParametersResponseBodyRunningParametersParameter) SetParameterValue(v string) *DescribeDBClusterParametersResponseBodyRunningParametersParameter {
	s.ParameterValue = &v
	return s
}

type DescribeDBClusterParametersResponse struct {
	Headers    map[string]*string                       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterParametersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterParametersResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterParametersResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterParametersResponse) SetHeaders(v map[string]*string) *DescribeDBClusterParametersResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterParametersResponse) SetStatusCode(v int32) *DescribeDBClusterParametersResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterParametersResponse) SetBody(v *DescribeDBClusterParametersResponseBody) *DescribeDBClusterParametersResponse {
	s.Body = v
	return s
}

type DescribeDBClusterPerformanceRequest struct {
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EndTime     *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	Key         *string `json:"Key,omitempty" xml:"Key,omitempty"`
	StartTime   *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeDBClusterPerformanceRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterPerformanceRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterPerformanceRequest) SetDBClusterId(v string) *DescribeDBClusterPerformanceRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterPerformanceRequest) SetEndTime(v string) *DescribeDBClusterPerformanceRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeDBClusterPerformanceRequest) SetKey(v string) *DescribeDBClusterPerformanceRequest {
	s.Key = &v
	return s
}

func (s *DescribeDBClusterPerformanceRequest) SetStartTime(v string) *DescribeDBClusterPerformanceRequest {
	s.StartTime = &v
	return s
}

type DescribeDBClusterPerformanceResponseBody struct {
	DBClusterId     *string                                                  `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBType          *string                                                  `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion       *string                                                  `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	EndTime         *string                                                  `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	PerformanceKeys *DescribeDBClusterPerformanceResponseBodyPerformanceKeys `json:"PerformanceKeys,omitempty" xml:"PerformanceKeys,omitempty" type:"Struct"`
	RequestId       *string                                                  `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	StartTime       *string                                                  `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeDBClusterPerformanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterPerformanceResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterPerformanceResponseBody) SetDBClusterId(v string) *DescribeDBClusterPerformanceResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBody) SetDBType(v string) *DescribeDBClusterPerformanceResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBody) SetDBVersion(v string) *DescribeDBClusterPerformanceResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBody) SetEndTime(v string) *DescribeDBClusterPerformanceResponseBody {
	s.EndTime = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBody) SetPerformanceKeys(v *DescribeDBClusterPerformanceResponseBodyPerformanceKeys) *DescribeDBClusterPerformanceResponseBody {
	s.PerformanceKeys = v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBody) SetRequestId(v string) *DescribeDBClusterPerformanceResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBody) SetStartTime(v string) *DescribeDBClusterPerformanceResponseBody {
	s.StartTime = &v
	return s
}

type DescribeDBClusterPerformanceResponseBodyPerformanceKeys struct {
	PerformanceItem []*DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem `json:"PerformanceItem,omitempty" xml:"PerformanceItem,omitempty" type:"Repeated"`
}

func (s DescribeDBClusterPerformanceResponseBodyPerformanceKeys) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterPerformanceResponseBodyPerformanceKeys) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterPerformanceResponseBodyPerformanceKeys) SetPerformanceItem(v []*DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem) *DescribeDBClusterPerformanceResponseBodyPerformanceKeys {
	s.PerformanceItem = v
	return s
}

type DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem struct {
	DBNodeId    *string                                                                       `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	Measurement *string                                                                       `json:"Measurement,omitempty" xml:"Measurement,omitempty"`
	MetricName  *string                                                                       `json:"MetricName,omitempty" xml:"MetricName,omitempty"`
	Points      *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPoints `json:"Points,omitempty" xml:"Points,omitempty" type:"Struct"`
}

func (s DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem) SetDBNodeId(v string) *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.DBNodeId = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem) SetMeasurement(v string) *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.Measurement = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem) SetMetricName(v string) *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.MetricName = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem) SetPoints(v *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPoints) *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.Points = v
	return s
}

type DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPoints struct {
	PerformanceItemValue []*DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue `json:"PerformanceItemValue,omitempty" xml:"PerformanceItemValue,omitempty" type:"Repeated"`
}

func (s DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPoints) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPoints) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPoints) SetPerformanceItemValue(v []*DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPoints {
	s.PerformanceItemValue = v
	return s
}

type DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue struct {
	Timestamp *int64  `json:"Timestamp,omitempty" xml:"Timestamp,omitempty"`
	Value     *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) SetTimestamp(v int64) *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue {
	s.Timestamp = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) SetValue(v string) *DescribeDBClusterPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue {
	s.Value = &v
	return s
}

type DescribeDBClusterPerformanceResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterPerformanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterPerformanceResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterPerformanceResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterPerformanceResponse) SetHeaders(v map[string]*string) *DescribeDBClusterPerformanceResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterPerformanceResponse) SetStatusCode(v int32) *DescribeDBClusterPerformanceResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterPerformanceResponse) SetBody(v *DescribeDBClusterPerformanceResponseBody) *DescribeDBClusterPerformanceResponse {
	s.Body = v
	return s
}

type DescribeDBClusterSSLRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterSSLRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterSSLRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterSSLRequest) SetDBClusterId(v string) *DescribeDBClusterSSLRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterSSLRequest) SetOwnerAccount(v string) *DescribeDBClusterSSLRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterSSLRequest) SetOwnerId(v int64) *DescribeDBClusterSSLRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterSSLRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterSSLRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterSSLRequest) SetResourceOwnerId(v int64) *DescribeDBClusterSSLRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterSSLResponseBody struct {
	Items         []*DescribeDBClusterSSLResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Repeated"`
	RequestId     *string                                  `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	SSLAutoRotate *string                                  `json:"SSLAutoRotate,omitempty" xml:"SSLAutoRotate,omitempty"`
}

func (s DescribeDBClusterSSLResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterSSLResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterSSLResponseBody) SetItems(v []*DescribeDBClusterSSLResponseBodyItems) *DescribeDBClusterSSLResponseBody {
	s.Items = v
	return s
}

func (s *DescribeDBClusterSSLResponseBody) SetRequestId(v string) *DescribeDBClusterSSLResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBClusterSSLResponseBody) SetSSLAutoRotate(v string) *DescribeDBClusterSSLResponseBody {
	s.SSLAutoRotate = &v
	return s
}

type DescribeDBClusterSSLResponseBodyItems struct {
	DBEndpointId        *string `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	SSLConnectionString *string `json:"SSLConnectionString,omitempty" xml:"SSLConnectionString,omitempty"`
	SSLEnabled          *string `json:"SSLEnabled,omitempty" xml:"SSLEnabled,omitempty"`
	SSLExpireTime       *string `json:"SSLExpireTime,omitempty" xml:"SSLExpireTime,omitempty"`
}

func (s DescribeDBClusterSSLResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterSSLResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterSSLResponseBodyItems) SetDBEndpointId(v string) *DescribeDBClusterSSLResponseBodyItems {
	s.DBEndpointId = &v
	return s
}

func (s *DescribeDBClusterSSLResponseBodyItems) SetSSLConnectionString(v string) *DescribeDBClusterSSLResponseBodyItems {
	s.SSLConnectionString = &v
	return s
}

func (s *DescribeDBClusterSSLResponseBodyItems) SetSSLEnabled(v string) *DescribeDBClusterSSLResponseBodyItems {
	s.SSLEnabled = &v
	return s
}

func (s *DescribeDBClusterSSLResponseBodyItems) SetSSLExpireTime(v string) *DescribeDBClusterSSLResponseBodyItems {
	s.SSLExpireTime = &v
	return s
}

type DescribeDBClusterSSLResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterSSLResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterSSLResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterSSLResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterSSLResponse) SetHeaders(v map[string]*string) *DescribeDBClusterSSLResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterSSLResponse) SetStatusCode(v int32) *DescribeDBClusterSSLResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterSSLResponse) SetBody(v *DescribeDBClusterSSLResponseBody) *DescribeDBClusterSSLResponse {
	s.Body = v
	return s
}

type DescribeDBClusterTDERequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterTDERequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterTDERequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterTDERequest) SetDBClusterId(v string) *DescribeDBClusterTDERequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterTDERequest) SetOwnerAccount(v string) *DescribeDBClusterTDERequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterTDERequest) SetOwnerId(v int64) *DescribeDBClusterTDERequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterTDERequest) SetResourceOwnerAccount(v string) *DescribeDBClusterTDERequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterTDERequest) SetResourceOwnerId(v int64) *DescribeDBClusterTDERequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterTDEResponseBody struct {
	DBClusterId      *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EncryptNewTables *string `json:"EncryptNewTables,omitempty" xml:"EncryptNewTables,omitempty"`
	EncryptionKey    *string `json:"EncryptionKey,omitempty" xml:"EncryptionKey,omitempty"`
	RequestId        *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TDEStatus        *string `json:"TDEStatus,omitempty" xml:"TDEStatus,omitempty"`
}

func (s DescribeDBClusterTDEResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterTDEResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterTDEResponseBody) SetDBClusterId(v string) *DescribeDBClusterTDEResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterTDEResponseBody) SetEncryptNewTables(v string) *DescribeDBClusterTDEResponseBody {
	s.EncryptNewTables = &v
	return s
}

func (s *DescribeDBClusterTDEResponseBody) SetEncryptionKey(v string) *DescribeDBClusterTDEResponseBody {
	s.EncryptionKey = &v
	return s
}

func (s *DescribeDBClusterTDEResponseBody) SetRequestId(v string) *DescribeDBClusterTDEResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBClusterTDEResponseBody) SetTDEStatus(v string) *DescribeDBClusterTDEResponseBody {
	s.TDEStatus = &v
	return s
}

type DescribeDBClusterTDEResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterTDEResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterTDEResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterTDEResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterTDEResponse) SetHeaders(v map[string]*string) *DescribeDBClusterTDEResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterTDEResponse) SetStatusCode(v int32) *DescribeDBClusterTDEResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterTDEResponse) SetBody(v *DescribeDBClusterTDEResponseBody) *DescribeDBClusterTDEResponse {
	s.Body = v
	return s
}

type DescribeDBClusterVersionRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClusterVersionRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterVersionRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterVersionRequest) SetDBClusterId(v string) *DescribeDBClusterVersionRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterVersionRequest) SetOwnerAccount(v string) *DescribeDBClusterVersionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClusterVersionRequest) SetOwnerId(v int64) *DescribeDBClusterVersionRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClusterVersionRequest) SetResourceOwnerAccount(v string) *DescribeDBClusterVersionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClusterVersionRequest) SetResourceOwnerId(v int64) *DescribeDBClusterVersionRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClusterVersionResponseBody struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBLatestVersion      *string `json:"DBLatestVersion,omitempty" xml:"DBLatestVersion,omitempty"`
	DBMinorVersion       *string `json:"DBMinorVersion,omitempty" xml:"DBMinorVersion,omitempty"`
	DBRevisionVersion    *string `json:"DBRevisionVersion,omitempty" xml:"DBRevisionVersion,omitempty"`
	DBVersion            *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	DBVersionStatus      *string `json:"DBVersionStatus,omitempty" xml:"DBVersionStatus,omitempty"`
	IsLatestVersion      *string `json:"IsLatestVersion,omitempty" xml:"IsLatestVersion,omitempty"`
	IsProxyLatestVersion *string `json:"IsProxyLatestVersion,omitempty" xml:"IsProxyLatestVersion,omitempty"`
	ProxyLatestVersion   *string `json:"ProxyLatestVersion,omitempty" xml:"ProxyLatestVersion,omitempty"`
	ProxyRevisionVersion *string `json:"ProxyRevisionVersion,omitempty" xml:"ProxyRevisionVersion,omitempty"`
	ProxyVersionStatus   *string `json:"ProxyVersionStatus,omitempty" xml:"ProxyVersionStatus,omitempty"`
	RequestId            *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDBClusterVersionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterVersionResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterVersionResponseBody) SetDBClusterId(v string) *DescribeDBClusterVersionResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetDBLatestVersion(v string) *DescribeDBClusterVersionResponseBody {
	s.DBLatestVersion = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetDBMinorVersion(v string) *DescribeDBClusterVersionResponseBody {
	s.DBMinorVersion = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetDBRevisionVersion(v string) *DescribeDBClusterVersionResponseBody {
	s.DBRevisionVersion = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetDBVersion(v string) *DescribeDBClusterVersionResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetDBVersionStatus(v string) *DescribeDBClusterVersionResponseBody {
	s.DBVersionStatus = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetIsLatestVersion(v string) *DescribeDBClusterVersionResponseBody {
	s.IsLatestVersion = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetIsProxyLatestVersion(v string) *DescribeDBClusterVersionResponseBody {
	s.IsProxyLatestVersion = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetProxyLatestVersion(v string) *DescribeDBClusterVersionResponseBody {
	s.ProxyLatestVersion = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetProxyRevisionVersion(v string) *DescribeDBClusterVersionResponseBody {
	s.ProxyRevisionVersion = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetProxyVersionStatus(v string) *DescribeDBClusterVersionResponseBody {
	s.ProxyVersionStatus = &v
	return s
}

func (s *DescribeDBClusterVersionResponseBody) SetRequestId(v string) *DescribeDBClusterVersionResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDBClusterVersionResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClusterVersionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClusterVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClusterVersionResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClusterVersionResponse) SetHeaders(v map[string]*string) *DescribeDBClusterVersionResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClusterVersionResponse) SetStatusCode(v int32) *DescribeDBClusterVersionResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClusterVersionResponse) SetBody(v *DescribeDBClusterVersionResponseBody) *DescribeDBClusterVersionResponse {
	s.Body = v
	return s
}

type DescribeDBClustersRequest struct {
	DBClusterDescription *string                         `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBClusterIds         *string                         `json:"DBClusterIds,omitempty" xml:"DBClusterIds,omitempty"`
	DBClusterStatus      *string                         `json:"DBClusterStatus,omitempty" xml:"DBClusterStatus,omitempty"`
	DBNodeIds            *string                         `json:"DBNodeIds,omitempty" xml:"DBNodeIds,omitempty"`
	DBType               *string                         `json:"DBType,omitempty" xml:"DBType,omitempty"`
	OwnerAccount         *string                         `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64                          `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32                          `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32                          `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	PayType              *string                         `json:"PayType,omitempty" xml:"PayType,omitempty"`
	RegionId             *string                         `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceGroupId      *string                         `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	ResourceOwnerAccount *string                         `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64                          `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	Tag                  []*DescribeDBClustersRequestTag `json:"Tag,omitempty" xml:"Tag,omitempty" type:"Repeated"`
}

func (s DescribeDBClustersRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersRequest) SetDBClusterDescription(v string) *DescribeDBClustersRequest {
	s.DBClusterDescription = &v
	return s
}

func (s *DescribeDBClustersRequest) SetDBClusterIds(v string) *DescribeDBClustersRequest {
	s.DBClusterIds = &v
	return s
}

func (s *DescribeDBClustersRequest) SetDBClusterStatus(v string) *DescribeDBClustersRequest {
	s.DBClusterStatus = &v
	return s
}

func (s *DescribeDBClustersRequest) SetDBNodeIds(v string) *DescribeDBClustersRequest {
	s.DBNodeIds = &v
	return s
}

func (s *DescribeDBClustersRequest) SetDBType(v string) *DescribeDBClustersRequest {
	s.DBType = &v
	return s
}

func (s *DescribeDBClustersRequest) SetOwnerAccount(v string) *DescribeDBClustersRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClustersRequest) SetOwnerId(v int64) *DescribeDBClustersRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClustersRequest) SetPageNumber(v int32) *DescribeDBClustersRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeDBClustersRequest) SetPageSize(v int32) *DescribeDBClustersRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeDBClustersRequest) SetPayType(v string) *DescribeDBClustersRequest {
	s.PayType = &v
	return s
}

func (s *DescribeDBClustersRequest) SetRegionId(v string) *DescribeDBClustersRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeDBClustersRequest) SetResourceGroupId(v string) *DescribeDBClustersRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeDBClustersRequest) SetResourceOwnerAccount(v string) *DescribeDBClustersRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClustersRequest) SetResourceOwnerId(v int64) *DescribeDBClustersRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeDBClustersRequest) SetTag(v []*DescribeDBClustersRequestTag) *DescribeDBClustersRequest {
	s.Tag = v
	return s
}

type DescribeDBClustersRequestTag struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s DescribeDBClustersRequestTag) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersRequestTag) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersRequestTag) SetKey(v string) *DescribeDBClustersRequestTag {
	s.Key = &v
	return s
}

func (s *DescribeDBClustersRequestTag) SetValue(v string) *DescribeDBClustersRequestTag {
	s.Value = &v
	return s
}

type DescribeDBClustersResponseBody struct {
	Items            *DescribeDBClustersResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	PageNumber       *int32                               `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *int32                               `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                              `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *int32                               `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeDBClustersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersResponseBody) SetItems(v *DescribeDBClustersResponseBodyItems) *DescribeDBClustersResponseBody {
	s.Items = v
	return s
}

func (s *DescribeDBClustersResponseBody) SetPageNumber(v int32) *DescribeDBClustersResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeDBClustersResponseBody) SetPageRecordCount(v int32) *DescribeDBClustersResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeDBClustersResponseBody) SetRequestId(v string) *DescribeDBClustersResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBClustersResponseBody) SetTotalRecordCount(v int32) *DescribeDBClustersResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeDBClustersResponseBodyItems struct {
	DBCluster []*DescribeDBClustersResponseBodyItemsDBCluster `json:"DBCluster,omitempty" xml:"DBCluster,omitempty" type:"Repeated"`
}

func (s DescribeDBClustersResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersResponseBodyItems) SetDBCluster(v []*DescribeDBClustersResponseBodyItemsDBCluster) *DescribeDBClustersResponseBodyItems {
	s.DBCluster = v
	return s
}

type DescribeDBClustersResponseBodyItemsDBCluster struct {
	Category             *string                                              `json:"Category,omitempty" xml:"Category,omitempty"`
	CreateTime           *string                                              `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	DBClusterDescription *string                                              `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBClusterId          *string                                              `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBClusterNetworkType *string                                              `json:"DBClusterNetworkType,omitempty" xml:"DBClusterNetworkType,omitempty"`
	DBClusterStatus      *string                                              `json:"DBClusterStatus,omitempty" xml:"DBClusterStatus,omitempty"`
	DBNodeClass          *string                                              `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
	DBNodeNumber         *int32                                               `json:"DBNodeNumber,omitempty" xml:"DBNodeNumber,omitempty"`
	DBNodes              *DescribeDBClustersResponseBodyItemsDBClusterDBNodes `json:"DBNodes,omitempty" xml:"DBNodes,omitempty" type:"Struct"`
	DBType               *string                                              `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion            *string                                              `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	DeletionLock         *int32                                               `json:"DeletionLock,omitempty" xml:"DeletionLock,omitempty"`
	Engine               *string                                              `json:"Engine,omitempty" xml:"Engine,omitempty"`
	ExpireTime           *string                                              `json:"ExpireTime,omitempty" xml:"ExpireTime,omitempty"`
	Expired              *string                                              `json:"Expired,omitempty" xml:"Expired,omitempty"`
	LockMode             *string                                              `json:"LockMode,omitempty" xml:"LockMode,omitempty"`
	PayType              *string                                              `json:"PayType,omitempty" xml:"PayType,omitempty"`
	RegionId             *string                                              `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceGroupId      *string                                              `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	StoragePayType       *string                                              `json:"StoragePayType,omitempty" xml:"StoragePayType,omitempty"`
	StorageSpace         *int64                                               `json:"StorageSpace,omitempty" xml:"StorageSpace,omitempty"`
	StorageUsed          *int64                                               `json:"StorageUsed,omitempty" xml:"StorageUsed,omitempty"`
	Tags                 *DescribeDBClustersResponseBodyItemsDBClusterTags    `json:"Tags,omitempty" xml:"Tags,omitempty" type:"Struct"`
	VpcId                *string                                              `json:"VpcId,omitempty" xml:"VpcId,omitempty"`
	ZoneId               *string                                              `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s DescribeDBClustersResponseBodyItemsDBCluster) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersResponseBodyItemsDBCluster) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetCategory(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.Category = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetCreateTime(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.CreateTime = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBClusterDescription(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBClusterDescription = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBClusterId(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBClusterNetworkType(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBClusterNetworkType = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBClusterStatus(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBClusterStatus = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBNodeClass(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBNodeClass = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBNodeNumber(v int32) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBNodeNumber = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBNodes(v *DescribeDBClustersResponseBodyItemsDBClusterDBNodes) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBNodes = v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBType(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBType = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDBVersion(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetDeletionLock(v int32) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.DeletionLock = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetEngine(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.Engine = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetExpireTime(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.ExpireTime = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetExpired(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.Expired = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetLockMode(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.LockMode = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetPayType(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.PayType = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetRegionId(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.RegionId = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetResourceGroupId(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetStoragePayType(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.StoragePayType = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetStorageSpace(v int64) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.StorageSpace = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetStorageUsed(v int64) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.StorageUsed = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetTags(v *DescribeDBClustersResponseBodyItemsDBClusterTags) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.Tags = v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetVpcId(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.VpcId = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBCluster) SetZoneId(v string) *DescribeDBClustersResponseBodyItemsDBCluster {
	s.ZoneId = &v
	return s
}

type DescribeDBClustersResponseBodyItemsDBClusterDBNodes struct {
	DBNode []*DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode `json:"DBNode,omitempty" xml:"DBNode,omitempty" type:"Repeated"`
}

func (s DescribeDBClustersResponseBodyItemsDBClusterDBNodes) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersResponseBodyItemsDBClusterDBNodes) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterDBNodes) SetDBNode(v []*DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode) *DescribeDBClustersResponseBodyItemsDBClusterDBNodes {
	s.DBNode = v
	return s
}

type DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode struct {
	DBNodeClass *string `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
	DBNodeId    *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	DBNodeRole  *string `json:"DBNodeRole,omitempty" xml:"DBNodeRole,omitempty"`
	RegionId    *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ZoneId      *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode) SetDBNodeClass(v string) *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode {
	s.DBNodeClass = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode) SetDBNodeId(v string) *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode {
	s.DBNodeId = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode) SetDBNodeRole(v string) *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode {
	s.DBNodeRole = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode) SetRegionId(v string) *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode {
	s.RegionId = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode) SetZoneId(v string) *DescribeDBClustersResponseBodyItemsDBClusterDBNodesDBNode {
	s.ZoneId = &v
	return s
}

type DescribeDBClustersResponseBodyItemsDBClusterTags struct {
	Tag []*DescribeDBClustersResponseBodyItemsDBClusterTagsTag `json:"Tag,omitempty" xml:"Tag,omitempty" type:"Repeated"`
}

func (s DescribeDBClustersResponseBodyItemsDBClusterTags) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersResponseBodyItemsDBClusterTags) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterTags) SetTag(v []*DescribeDBClustersResponseBodyItemsDBClusterTagsTag) *DescribeDBClustersResponseBodyItemsDBClusterTags {
	s.Tag = v
	return s
}

type DescribeDBClustersResponseBodyItemsDBClusterTagsTag struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s DescribeDBClustersResponseBodyItemsDBClusterTagsTag) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersResponseBodyItemsDBClusterTagsTag) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterTagsTag) SetKey(v string) *DescribeDBClustersResponseBodyItemsDBClusterTagsTag {
	s.Key = &v
	return s
}

func (s *DescribeDBClustersResponseBodyItemsDBClusterTagsTag) SetValue(v string) *DescribeDBClustersResponseBodyItemsDBClusterTagsTag {
	s.Value = &v
	return s
}

type DescribeDBClustersResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClustersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClustersResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersResponse) SetHeaders(v map[string]*string) *DescribeDBClustersResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClustersResponse) SetStatusCode(v int32) *DescribeDBClustersResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClustersResponse) SetBody(v *DescribeDBClustersResponseBody) *DescribeDBClustersResponse {
	s.Body = v
	return s
}

type DescribeDBClustersWithBackupsRequest struct {
	DBClusterDescription *string `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBClusterIds         *string `json:"DBClusterIds,omitempty" xml:"DBClusterIds,omitempty"`
	DBType               *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion            *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	IsDeleted            *int32  `json:"IsDeleted,omitempty" xml:"IsDeleted,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBClustersWithBackupsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersWithBackupsRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersWithBackupsRequest) SetDBClusterDescription(v string) *DescribeDBClustersWithBackupsRequest {
	s.DBClusterDescription = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetDBClusterIds(v string) *DescribeDBClustersWithBackupsRequest {
	s.DBClusterIds = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetDBType(v string) *DescribeDBClustersWithBackupsRequest {
	s.DBType = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetDBVersion(v string) *DescribeDBClustersWithBackupsRequest {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetIsDeleted(v int32) *DescribeDBClustersWithBackupsRequest {
	s.IsDeleted = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetOwnerAccount(v string) *DescribeDBClustersWithBackupsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetOwnerId(v int64) *DescribeDBClustersWithBackupsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetPageNumber(v int32) *DescribeDBClustersWithBackupsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetPageSize(v int32) *DescribeDBClustersWithBackupsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetRegionId(v string) *DescribeDBClustersWithBackupsRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetResourceOwnerAccount(v string) *DescribeDBClustersWithBackupsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBClustersWithBackupsRequest) SetResourceOwnerId(v int64) *DescribeDBClustersWithBackupsRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBClustersWithBackupsResponseBody struct {
	Items            *DescribeDBClustersWithBackupsResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	PageNumber       *int32                                          `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *int32                                          `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                                         `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *int32                                          `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeDBClustersWithBackupsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersWithBackupsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersWithBackupsResponseBody) SetItems(v *DescribeDBClustersWithBackupsResponseBodyItems) *DescribeDBClustersWithBackupsResponseBody {
	s.Items = v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBody) SetPageNumber(v int32) *DescribeDBClustersWithBackupsResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBody) SetPageRecordCount(v int32) *DescribeDBClustersWithBackupsResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBody) SetRequestId(v string) *DescribeDBClustersWithBackupsResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBody) SetTotalRecordCount(v int32) *DescribeDBClustersWithBackupsResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeDBClustersWithBackupsResponseBodyItems struct {
	DBCluster []*DescribeDBClustersWithBackupsResponseBodyItemsDBCluster `json:"DBCluster,omitempty" xml:"DBCluster,omitempty" type:"Repeated"`
}

func (s DescribeDBClustersWithBackupsResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersWithBackupsResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersWithBackupsResponseBodyItems) SetDBCluster(v []*DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) *DescribeDBClustersWithBackupsResponseBodyItems {
	s.DBCluster = v
	return s
}

type DescribeDBClustersWithBackupsResponseBodyItemsDBCluster struct {
	CreateTime           *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	DBClusterDescription *string `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBClusterNetworkType *string `json:"DBClusterNetworkType,omitempty" xml:"DBClusterNetworkType,omitempty"`
	DBClusterStatus      *string `json:"DBClusterStatus,omitempty" xml:"DBClusterStatus,omitempty"`
	DBNodeClass          *string `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
	DBType               *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion            *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	DeletedTime          *string `json:"DeletedTime,omitempty" xml:"DeletedTime,omitempty"`
	DeletionLock         *int32  `json:"DeletionLock,omitempty" xml:"DeletionLock,omitempty"`
	Engine               *string `json:"Engine,omitempty" xml:"Engine,omitempty"`
	ExpireTime           *string `json:"ExpireTime,omitempty" xml:"ExpireTime,omitempty"`
	Expired              *string `json:"Expired,omitempty" xml:"Expired,omitempty"`
	IsDeleted            *int32  `json:"IsDeleted,omitempty" xml:"IsDeleted,omitempty"`
	LockMode             *string `json:"LockMode,omitempty" xml:"LockMode,omitempty"`
	PayType              *string `json:"PayType,omitempty" xml:"PayType,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	VpcId                *string `json:"VpcId,omitempty" xml:"VpcId,omitempty"`
	ZoneId               *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetCreateTime(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.CreateTime = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDBClusterDescription(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DBClusterDescription = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDBClusterId(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDBClusterNetworkType(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DBClusterNetworkType = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDBClusterStatus(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DBClusterStatus = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDBNodeClass(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DBNodeClass = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDBType(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DBType = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDBVersion(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDeletedTime(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DeletedTime = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetDeletionLock(v int32) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.DeletionLock = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetEngine(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.Engine = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetExpireTime(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.ExpireTime = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetExpired(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.Expired = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetIsDeleted(v int32) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.IsDeleted = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetLockMode(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.LockMode = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetPayType(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.PayType = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetRegionId(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.RegionId = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetVpcId(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.VpcId = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster) SetZoneId(v string) *DescribeDBClustersWithBackupsResponseBodyItemsDBCluster {
	s.ZoneId = &v
	return s
}

type DescribeDBClustersWithBackupsResponse struct {
	Headers    map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBClustersWithBackupsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBClustersWithBackupsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBClustersWithBackupsResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBClustersWithBackupsResponse) SetHeaders(v map[string]*string) *DescribeDBClustersWithBackupsResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBClustersWithBackupsResponse) SetStatusCode(v int32) *DescribeDBClustersWithBackupsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBClustersWithBackupsResponse) SetBody(v *DescribeDBClustersWithBackupsResponseBody) *DescribeDBClustersWithBackupsResponse {
	s.Body = v
	return s
}

type DescribeDBInitializeVariableRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBInitializeVariableRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBInitializeVariableRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBInitializeVariableRequest) SetDBClusterId(v string) *DescribeDBInitializeVariableRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBInitializeVariableRequest) SetOwnerAccount(v string) *DescribeDBInitializeVariableRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBInitializeVariableRequest) SetOwnerId(v int64) *DescribeDBInitializeVariableRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBInitializeVariableRequest) SetResourceOwnerAccount(v string) *DescribeDBInitializeVariableRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBInitializeVariableRequest) SetResourceOwnerId(v int64) *DescribeDBInitializeVariableRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBInitializeVariableResponseBody struct {
	DBType    *string                                            `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion *string                                            `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	RequestId *string                                            `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	Variables *DescribeDBInitializeVariableResponseBodyVariables `json:"Variables,omitempty" xml:"Variables,omitempty" type:"Struct"`
}

func (s DescribeDBInitializeVariableResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBInitializeVariableResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBInitializeVariableResponseBody) SetDBType(v string) *DescribeDBInitializeVariableResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeDBInitializeVariableResponseBody) SetDBVersion(v string) *DescribeDBInitializeVariableResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBInitializeVariableResponseBody) SetRequestId(v string) *DescribeDBInitializeVariableResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBInitializeVariableResponseBody) SetVariables(v *DescribeDBInitializeVariableResponseBodyVariables) *DescribeDBInitializeVariableResponseBody {
	s.Variables = v
	return s
}

type DescribeDBInitializeVariableResponseBodyVariables struct {
	Variable []*DescribeDBInitializeVariableResponseBodyVariablesVariable `json:"Variable,omitempty" xml:"Variable,omitempty" type:"Repeated"`
}

func (s DescribeDBInitializeVariableResponseBodyVariables) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBInitializeVariableResponseBodyVariables) GoString() string {
	return s.String()
}

func (s *DescribeDBInitializeVariableResponseBodyVariables) SetVariable(v []*DescribeDBInitializeVariableResponseBodyVariablesVariable) *DescribeDBInitializeVariableResponseBodyVariables {
	s.Variable = v
	return s
}

type DescribeDBInitializeVariableResponseBodyVariablesVariable struct {
	Charset *string `json:"Charset,omitempty" xml:"Charset,omitempty"`
	Collate *string `json:"Collate,omitempty" xml:"Collate,omitempty"`
	Ctype   *string `json:"Ctype,omitempty" xml:"Ctype,omitempty"`
}

func (s DescribeDBInitializeVariableResponseBodyVariablesVariable) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBInitializeVariableResponseBodyVariablesVariable) GoString() string {
	return s.String()
}

func (s *DescribeDBInitializeVariableResponseBodyVariablesVariable) SetCharset(v string) *DescribeDBInitializeVariableResponseBodyVariablesVariable {
	s.Charset = &v
	return s
}

func (s *DescribeDBInitializeVariableResponseBodyVariablesVariable) SetCollate(v string) *DescribeDBInitializeVariableResponseBodyVariablesVariable {
	s.Collate = &v
	return s
}

func (s *DescribeDBInitializeVariableResponseBodyVariablesVariable) SetCtype(v string) *DescribeDBInitializeVariableResponseBodyVariablesVariable {
	s.Ctype = &v
	return s
}

type DescribeDBInitializeVariableResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBInitializeVariableResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBInitializeVariableResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBInitializeVariableResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBInitializeVariableResponse) SetHeaders(v map[string]*string) *DescribeDBInitializeVariableResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBInitializeVariableResponse) SetStatusCode(v int32) *DescribeDBInitializeVariableResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBInitializeVariableResponse) SetBody(v *DescribeDBInitializeVariableResponseBody) *DescribeDBInitializeVariableResponse {
	s.Body = v
	return s
}

type DescribeDBLinksRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBLinkName           *string `json:"DBLinkName,omitempty" xml:"DBLinkName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBLinksRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBLinksRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBLinksRequest) SetDBClusterId(v string) *DescribeDBLinksRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBLinksRequest) SetDBLinkName(v string) *DescribeDBLinksRequest {
	s.DBLinkName = &v
	return s
}

func (s *DescribeDBLinksRequest) SetOwnerAccount(v string) *DescribeDBLinksRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBLinksRequest) SetOwnerId(v int64) *DescribeDBLinksRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBLinksRequest) SetResourceOwnerAccount(v string) *DescribeDBLinksRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBLinksRequest) SetResourceOwnerId(v int64) *DescribeDBLinksRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBLinksResponseBody struct {
	DBInstanceName *string                                   `json:"DBInstanceName,omitempty" xml:"DBInstanceName,omitempty"`
	DBLinkInfos    []*DescribeDBLinksResponseBodyDBLinkInfos `json:"DBLinkInfos,omitempty" xml:"DBLinkInfos,omitempty" type:"Repeated"`
	RequestId      *string                                   `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDBLinksResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBLinksResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBLinksResponseBody) SetDBInstanceName(v string) *DescribeDBLinksResponseBody {
	s.DBInstanceName = &v
	return s
}

func (s *DescribeDBLinksResponseBody) SetDBLinkInfos(v []*DescribeDBLinksResponseBodyDBLinkInfos) *DescribeDBLinksResponseBody {
	s.DBLinkInfos = v
	return s
}

func (s *DescribeDBLinksResponseBody) SetRequestId(v string) *DescribeDBLinksResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDBLinksResponseBodyDBLinkInfos struct {
	DBInstanceName       *string `json:"DBInstanceName,omitempty" xml:"DBInstanceName,omitempty"`
	DBLinkName           *string `json:"DBLinkName,omitempty" xml:"DBLinkName,omitempty"`
	SourceDBName         *string `json:"SourceDBName,omitempty" xml:"SourceDBName,omitempty"`
	TargetAccount        *string `json:"TargetAccount,omitempty" xml:"TargetAccount,omitempty"`
	TargetDBInstanceName *string `json:"TargetDBInstanceName,omitempty" xml:"TargetDBInstanceName,omitempty"`
	TargetDBName         *string `json:"TargetDBName,omitempty" xml:"TargetDBName,omitempty"`
}

func (s DescribeDBLinksResponseBodyDBLinkInfos) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBLinksResponseBodyDBLinkInfos) GoString() string {
	return s.String()
}

func (s *DescribeDBLinksResponseBodyDBLinkInfos) SetDBInstanceName(v string) *DescribeDBLinksResponseBodyDBLinkInfos {
	s.DBInstanceName = &v
	return s
}

func (s *DescribeDBLinksResponseBodyDBLinkInfos) SetDBLinkName(v string) *DescribeDBLinksResponseBodyDBLinkInfos {
	s.DBLinkName = &v
	return s
}

func (s *DescribeDBLinksResponseBodyDBLinkInfos) SetSourceDBName(v string) *DescribeDBLinksResponseBodyDBLinkInfos {
	s.SourceDBName = &v
	return s
}

func (s *DescribeDBLinksResponseBodyDBLinkInfos) SetTargetAccount(v string) *DescribeDBLinksResponseBodyDBLinkInfos {
	s.TargetAccount = &v
	return s
}

func (s *DescribeDBLinksResponseBodyDBLinkInfos) SetTargetDBInstanceName(v string) *DescribeDBLinksResponseBodyDBLinkInfos {
	s.TargetDBInstanceName = &v
	return s
}

func (s *DescribeDBLinksResponseBodyDBLinkInfos) SetTargetDBName(v string) *DescribeDBLinksResponseBodyDBLinkInfos {
	s.TargetDBName = &v
	return s
}

type DescribeDBLinksResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBLinksResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBLinksResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBLinksResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBLinksResponse) SetHeaders(v map[string]*string) *DescribeDBLinksResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBLinksResponse) SetStatusCode(v int32) *DescribeDBLinksResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBLinksResponse) SetBody(v *DescribeDBLinksResponseBody) *DescribeDBLinksResponse {
	s.Body = v
	return s
}

type DescribeDBNodePerformanceRequest struct {
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeId    *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	EndTime     *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	Key         *string `json:"Key,omitempty" xml:"Key,omitempty"`
	StartTime   *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeDBNodePerformanceRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodePerformanceRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBNodePerformanceRequest) SetDBClusterId(v string) *DescribeDBNodePerformanceRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBNodePerformanceRequest) SetDBNodeId(v string) *DescribeDBNodePerformanceRequest {
	s.DBNodeId = &v
	return s
}

func (s *DescribeDBNodePerformanceRequest) SetEndTime(v string) *DescribeDBNodePerformanceRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeDBNodePerformanceRequest) SetKey(v string) *DescribeDBNodePerformanceRequest {
	s.Key = &v
	return s
}

func (s *DescribeDBNodePerformanceRequest) SetStartTime(v string) *DescribeDBNodePerformanceRequest {
	s.StartTime = &v
	return s
}

type DescribeDBNodePerformanceResponseBody struct {
	DBNodeId        *string                                               `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	DBType          *string                                               `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion       *string                                               `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	EndTime         *string                                               `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	PerformanceKeys *DescribeDBNodePerformanceResponseBodyPerformanceKeys `json:"PerformanceKeys,omitempty" xml:"PerformanceKeys,omitempty" type:"Struct"`
	RequestId       *string                                               `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	StartTime       *string                                               `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeDBNodePerformanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodePerformanceResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBNodePerformanceResponseBody) SetDBNodeId(v string) *DescribeDBNodePerformanceResponseBody {
	s.DBNodeId = &v
	return s
}

func (s *DescribeDBNodePerformanceResponseBody) SetDBType(v string) *DescribeDBNodePerformanceResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeDBNodePerformanceResponseBody) SetDBVersion(v string) *DescribeDBNodePerformanceResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBNodePerformanceResponseBody) SetEndTime(v string) *DescribeDBNodePerformanceResponseBody {
	s.EndTime = &v
	return s
}

func (s *DescribeDBNodePerformanceResponseBody) SetPerformanceKeys(v *DescribeDBNodePerformanceResponseBodyPerformanceKeys) *DescribeDBNodePerformanceResponseBody {
	s.PerformanceKeys = v
	return s
}

func (s *DescribeDBNodePerformanceResponseBody) SetRequestId(v string) *DescribeDBNodePerformanceResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBNodePerformanceResponseBody) SetStartTime(v string) *DescribeDBNodePerformanceResponseBody {
	s.StartTime = &v
	return s
}

type DescribeDBNodePerformanceResponseBodyPerformanceKeys struct {
	PerformanceItem []*DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem `json:"PerformanceItem,omitempty" xml:"PerformanceItem,omitempty" type:"Repeated"`
}

func (s DescribeDBNodePerformanceResponseBodyPerformanceKeys) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodePerformanceResponseBodyPerformanceKeys) GoString() string {
	return s.String()
}

func (s *DescribeDBNodePerformanceResponseBodyPerformanceKeys) SetPerformanceItem(v []*DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem) *DescribeDBNodePerformanceResponseBodyPerformanceKeys {
	s.PerformanceItem = v
	return s
}

type DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem struct {
	Measurement *string                                                                    `json:"Measurement,omitempty" xml:"Measurement,omitempty"`
	MetricName  *string                                                                    `json:"MetricName,omitempty" xml:"MetricName,omitempty"`
	Points      *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPoints `json:"Points,omitempty" xml:"Points,omitempty" type:"Struct"`
}

func (s DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem) GoString() string {
	return s.String()
}

func (s *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem) SetMeasurement(v string) *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.Measurement = &v
	return s
}

func (s *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem) SetMetricName(v string) *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.MetricName = &v
	return s
}

func (s *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem) SetPoints(v *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPoints) *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.Points = v
	return s
}

type DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPoints struct {
	PerformanceItemValue []*DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue `json:"PerformanceItemValue,omitempty" xml:"PerformanceItemValue,omitempty" type:"Repeated"`
}

func (s DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPoints) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPoints) GoString() string {
	return s.String()
}

func (s *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPoints) SetPerformanceItemValue(v []*DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPoints {
	s.PerformanceItemValue = v
	return s
}

type DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue struct {
	Timestamp *int64  `json:"Timestamp,omitempty" xml:"Timestamp,omitempty"`
	Value     *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) GoString() string {
	return s.String()
}

func (s *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) SetTimestamp(v int64) *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue {
	s.Timestamp = &v
	return s
}

func (s *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) SetValue(v string) *DescribeDBNodePerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue {
	s.Value = &v
	return s
}

type DescribeDBNodePerformanceResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBNodePerformanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBNodePerformanceResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodePerformanceResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBNodePerformanceResponse) SetHeaders(v map[string]*string) *DescribeDBNodePerformanceResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBNodePerformanceResponse) SetStatusCode(v int32) *DescribeDBNodePerformanceResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBNodePerformanceResponse) SetBody(v *DescribeDBNodePerformanceResponseBody) *DescribeDBNodePerformanceResponse {
	s.Body = v
	return s
}

type DescribeDBNodesParametersRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeIds            *string `json:"DBNodeIds,omitempty" xml:"DBNodeIds,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDBNodesParametersRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodesParametersRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBNodesParametersRequest) SetDBClusterId(v string) *DescribeDBNodesParametersRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBNodesParametersRequest) SetDBNodeIds(v string) *DescribeDBNodesParametersRequest {
	s.DBNodeIds = &v
	return s
}

func (s *DescribeDBNodesParametersRequest) SetOwnerAccount(v string) *DescribeDBNodesParametersRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDBNodesParametersRequest) SetOwnerId(v int64) *DescribeDBNodesParametersRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDBNodesParametersRequest) SetResourceOwnerAccount(v string) *DescribeDBNodesParametersRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDBNodesParametersRequest) SetResourceOwnerId(v int64) *DescribeDBNodesParametersRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDBNodesParametersResponseBody struct {
	DBNodeIds []*DescribeDBNodesParametersResponseBodyDBNodeIds `json:"DBNodeIds,omitempty" xml:"DBNodeIds,omitempty" type:"Repeated"`
	DBType    *string                                           `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion *string                                           `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	Engine    *string                                           `json:"Engine,omitempty" xml:"Engine,omitempty"`
	RequestId *string                                           `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDBNodesParametersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodesParametersResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBNodesParametersResponseBody) SetDBNodeIds(v []*DescribeDBNodesParametersResponseBodyDBNodeIds) *DescribeDBNodesParametersResponseBody {
	s.DBNodeIds = v
	return s
}

func (s *DescribeDBNodesParametersResponseBody) SetDBType(v string) *DescribeDBNodesParametersResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBody) SetDBVersion(v string) *DescribeDBNodesParametersResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBody) SetEngine(v string) *DescribeDBNodesParametersResponseBody {
	s.Engine = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBody) SetRequestId(v string) *DescribeDBNodesParametersResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDBNodesParametersResponseBodyDBNodeIds struct {
	DBNodeId          *string                                                            `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	RunningParameters []*DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters `json:"RunningParameters,omitempty" xml:"RunningParameters,omitempty" type:"Repeated"`
}

func (s DescribeDBNodesParametersResponseBodyDBNodeIds) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodesParametersResponseBodyDBNodeIds) GoString() string {
	return s.String()
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIds) SetDBNodeId(v string) *DescribeDBNodesParametersResponseBodyDBNodeIds {
	s.DBNodeId = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIds) SetRunningParameters(v []*DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) *DescribeDBNodesParametersResponseBodyDBNodeIds {
	s.RunningParameters = v
	return s
}

type DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters struct {
	CheckingCode          *string `json:"CheckingCode,omitempty" xml:"CheckingCode,omitempty"`
	DataType              *string `json:"DataType,omitempty" xml:"DataType,omitempty"`
	DefaultParameterValue *string `json:"DefaultParameterValue,omitempty" xml:"DefaultParameterValue,omitempty"`
	Factor                *string `json:"Factor,omitempty" xml:"Factor,omitempty"`
	ForceRestart          *bool   `json:"ForceRestart,omitempty" xml:"ForceRestart,omitempty"`
	IsModifiable          *bool   `json:"IsModifiable,omitempty" xml:"IsModifiable,omitempty"`
	IsNodeAvailable       *string `json:"IsNodeAvailable,omitempty" xml:"IsNodeAvailable,omitempty"`
	ParamRelyRule         *string `json:"ParamRelyRule,omitempty" xml:"ParamRelyRule,omitempty"`
	ParameterDescription  *string `json:"ParameterDescription,omitempty" xml:"ParameterDescription,omitempty"`
	ParameterName         *string `json:"ParameterName,omitempty" xml:"ParameterName,omitempty"`
	ParameterStatus       *string `json:"ParameterStatus,omitempty" xml:"ParameterStatus,omitempty"`
	ParameterValue        *string `json:"ParameterValue,omitempty" xml:"ParameterValue,omitempty"`
}

func (s DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) GoString() string {
	return s.String()
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetCheckingCode(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.CheckingCode = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetDataType(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.DataType = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetDefaultParameterValue(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.DefaultParameterValue = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetFactor(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.Factor = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetForceRestart(v bool) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.ForceRestart = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetIsModifiable(v bool) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.IsModifiable = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetIsNodeAvailable(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.IsNodeAvailable = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetParamRelyRule(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.ParamRelyRule = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetParameterDescription(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.ParameterDescription = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetParameterName(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.ParameterName = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetParameterStatus(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.ParameterStatus = &v
	return s
}

func (s *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters) SetParameterValue(v string) *DescribeDBNodesParametersResponseBodyDBNodeIdsRunningParameters {
	s.ParameterValue = &v
	return s
}

type DescribeDBNodesParametersResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBNodesParametersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBNodesParametersResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBNodesParametersResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBNodesParametersResponse) SetHeaders(v map[string]*string) *DescribeDBNodesParametersResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBNodesParametersResponse) SetStatusCode(v int32) *DescribeDBNodesParametersResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBNodesParametersResponse) SetBody(v *DescribeDBNodesParametersResponseBody) *DescribeDBNodesParametersResponse {
	s.Body = v
	return s
}

type DescribeDBProxyPerformanceRequest struct {
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EndTime     *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	Key         *string `json:"Key,omitempty" xml:"Key,omitempty"`
	StartTime   *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeDBProxyPerformanceRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBProxyPerformanceRequest) GoString() string {
	return s.String()
}

func (s *DescribeDBProxyPerformanceRequest) SetDBClusterId(v string) *DescribeDBProxyPerformanceRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBProxyPerformanceRequest) SetEndTime(v string) *DescribeDBProxyPerformanceRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeDBProxyPerformanceRequest) SetKey(v string) *DescribeDBProxyPerformanceRequest {
	s.Key = &v
	return s
}

func (s *DescribeDBProxyPerformanceRequest) SetStartTime(v string) *DescribeDBProxyPerformanceRequest {
	s.StartTime = &v
	return s
}

type DescribeDBProxyPerformanceResponseBody struct {
	DBClusterId     *string                                                `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBType          *string                                                `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion       *string                                                `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	EndTime         *string                                                `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	PerformanceKeys *DescribeDBProxyPerformanceResponseBodyPerformanceKeys `json:"PerformanceKeys,omitempty" xml:"PerformanceKeys,omitempty" type:"Struct"`
	RequestId       *string                                                `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	StartTime       *string                                                `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeDBProxyPerformanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBProxyPerformanceResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDBProxyPerformanceResponseBody) SetDBClusterId(v string) *DescribeDBProxyPerformanceResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBody) SetDBType(v string) *DescribeDBProxyPerformanceResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBody) SetDBVersion(v string) *DescribeDBProxyPerformanceResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBody) SetEndTime(v string) *DescribeDBProxyPerformanceResponseBody {
	s.EndTime = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBody) SetPerformanceKeys(v *DescribeDBProxyPerformanceResponseBodyPerformanceKeys) *DescribeDBProxyPerformanceResponseBody {
	s.PerformanceKeys = v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBody) SetRequestId(v string) *DescribeDBProxyPerformanceResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBody) SetStartTime(v string) *DescribeDBProxyPerformanceResponseBody {
	s.StartTime = &v
	return s
}

type DescribeDBProxyPerformanceResponseBodyPerformanceKeys struct {
	PerformanceItem []*DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem `json:"PerformanceItem,omitempty" xml:"PerformanceItem,omitempty" type:"Repeated"`
}

func (s DescribeDBProxyPerformanceResponseBodyPerformanceKeys) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBProxyPerformanceResponseBodyPerformanceKeys) GoString() string {
	return s.String()
}

func (s *DescribeDBProxyPerformanceResponseBodyPerformanceKeys) SetPerformanceItem(v []*DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem) *DescribeDBProxyPerformanceResponseBodyPerformanceKeys {
	s.PerformanceItem = v
	return s
}

type DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem struct {
	DBNodeId    *string                                                                     `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	Measurement *string                                                                     `json:"Measurement,omitempty" xml:"Measurement,omitempty"`
	MetricName  *string                                                                     `json:"MetricName,omitempty" xml:"MetricName,omitempty"`
	Points      *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPoints `json:"Points,omitempty" xml:"Points,omitempty" type:"Struct"`
}

func (s DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem) GoString() string {
	return s.String()
}

func (s *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem) SetDBNodeId(v string) *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.DBNodeId = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem) SetMeasurement(v string) *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.Measurement = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem) SetMetricName(v string) *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.MetricName = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem) SetPoints(v *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPoints) *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItem {
	s.Points = v
	return s
}

type DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPoints struct {
	PerformanceItemValue []*DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue `json:"PerformanceItemValue,omitempty" xml:"PerformanceItemValue,omitempty" type:"Repeated"`
}

func (s DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPoints) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPoints) GoString() string {
	return s.String()
}

func (s *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPoints) SetPerformanceItemValue(v []*DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPoints {
	s.PerformanceItemValue = v
	return s
}

type DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue struct {
	Timestamp *int64  `json:"Timestamp,omitempty" xml:"Timestamp,omitempty"`
	Value     *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) GoString() string {
	return s.String()
}

func (s *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) SetTimestamp(v int64) *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue {
	s.Timestamp = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue) SetValue(v string) *DescribeDBProxyPerformanceResponseBodyPerformanceKeysPerformanceItemPointsPerformanceItemValue {
	s.Value = &v
	return s
}

type DescribeDBProxyPerformanceResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDBProxyPerformanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDBProxyPerformanceResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDBProxyPerformanceResponse) GoString() string {
	return s.String()
}

func (s *DescribeDBProxyPerformanceResponse) SetHeaders(v map[string]*string) *DescribeDBProxyPerformanceResponse {
	s.Headers = v
	return s
}

func (s *DescribeDBProxyPerformanceResponse) SetStatusCode(v int32) *DescribeDBProxyPerformanceResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDBProxyPerformanceResponse) SetBody(v *DescribeDBProxyPerformanceResponseBody) *DescribeDBProxyPerformanceResponse {
	s.Body = v
	return s
}

type DescribeDatabasesRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeDatabasesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDatabasesRequest) GoString() string {
	return s.String()
}

func (s *DescribeDatabasesRequest) SetDBClusterId(v string) *DescribeDatabasesRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDatabasesRequest) SetDBName(v string) *DescribeDatabasesRequest {
	s.DBName = &v
	return s
}

func (s *DescribeDatabasesRequest) SetOwnerAccount(v string) *DescribeDatabasesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDatabasesRequest) SetOwnerId(v int64) *DescribeDatabasesRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDatabasesRequest) SetPageNumber(v int32) *DescribeDatabasesRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeDatabasesRequest) SetPageSize(v int32) *DescribeDatabasesRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeDatabasesRequest) SetResourceOwnerAccount(v string) *DescribeDatabasesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDatabasesRequest) SetResourceOwnerId(v int64) *DescribeDatabasesRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeDatabasesResponseBody struct {
	Databases       *DescribeDatabasesResponseBodyDatabases `json:"Databases,omitempty" xml:"Databases,omitempty" type:"Struct"`
	PageNumber      *int32                                  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount *int32                                  `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId       *string                                 `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeDatabasesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDatabasesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDatabasesResponseBody) SetDatabases(v *DescribeDatabasesResponseBodyDatabases) *DescribeDatabasesResponseBody {
	s.Databases = v
	return s
}

func (s *DescribeDatabasesResponseBody) SetPageNumber(v int32) *DescribeDatabasesResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeDatabasesResponseBody) SetPageRecordCount(v int32) *DescribeDatabasesResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeDatabasesResponseBody) SetRequestId(v string) *DescribeDatabasesResponseBody {
	s.RequestId = &v
	return s
}

type DescribeDatabasesResponseBodyDatabases struct {
	Database []*DescribeDatabasesResponseBodyDatabasesDatabase `json:"Database,omitempty" xml:"Database,omitempty" type:"Repeated"`
}

func (s DescribeDatabasesResponseBodyDatabases) String() string {
	return tea.Prettify(s)
}

func (s DescribeDatabasesResponseBodyDatabases) GoString() string {
	return s.String()
}

func (s *DescribeDatabasesResponseBodyDatabases) SetDatabase(v []*DescribeDatabasesResponseBodyDatabasesDatabase) *DescribeDatabasesResponseBodyDatabases {
	s.Database = v
	return s
}

type DescribeDatabasesResponseBodyDatabasesDatabase struct {
	Accounts         *DescribeDatabasesResponseBodyDatabasesDatabaseAccounts `json:"Accounts,omitempty" xml:"Accounts,omitempty" type:"Struct"`
	CharacterSetName *string                                                 `json:"CharacterSetName,omitempty" xml:"CharacterSetName,omitempty"`
	DBDescription    *string                                                 `json:"DBDescription,omitempty" xml:"DBDescription,omitempty"`
	DBName           *string                                                 `json:"DBName,omitempty" xml:"DBName,omitempty"`
	DBStatus         *string                                                 `json:"DBStatus,omitempty" xml:"DBStatus,omitempty"`
	Engine           *string                                                 `json:"Engine,omitempty" xml:"Engine,omitempty"`
}

func (s DescribeDatabasesResponseBodyDatabasesDatabase) String() string {
	return tea.Prettify(s)
}

func (s DescribeDatabasesResponseBodyDatabasesDatabase) GoString() string {
	return s.String()
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabase) SetAccounts(v *DescribeDatabasesResponseBodyDatabasesDatabaseAccounts) *DescribeDatabasesResponseBodyDatabasesDatabase {
	s.Accounts = v
	return s
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabase) SetCharacterSetName(v string) *DescribeDatabasesResponseBodyDatabasesDatabase {
	s.CharacterSetName = &v
	return s
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabase) SetDBDescription(v string) *DescribeDatabasesResponseBodyDatabasesDatabase {
	s.DBDescription = &v
	return s
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabase) SetDBName(v string) *DescribeDatabasesResponseBodyDatabasesDatabase {
	s.DBName = &v
	return s
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabase) SetDBStatus(v string) *DescribeDatabasesResponseBodyDatabasesDatabase {
	s.DBStatus = &v
	return s
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabase) SetEngine(v string) *DescribeDatabasesResponseBodyDatabasesDatabase {
	s.Engine = &v
	return s
}

type DescribeDatabasesResponseBodyDatabasesDatabaseAccounts struct {
	Account []*DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount `json:"Account,omitempty" xml:"Account,omitempty" type:"Repeated"`
}

func (s DescribeDatabasesResponseBodyDatabasesDatabaseAccounts) String() string {
	return tea.Prettify(s)
}

func (s DescribeDatabasesResponseBodyDatabasesDatabaseAccounts) GoString() string {
	return s.String()
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabaseAccounts) SetAccount(v []*DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount) *DescribeDatabasesResponseBodyDatabasesDatabaseAccounts {
	s.Account = v
	return s
}

type DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount struct {
	AccountName      *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	AccountPrivilege *string `json:"AccountPrivilege,omitempty" xml:"AccountPrivilege,omitempty"`
	AccountStatus    *string `json:"AccountStatus,omitempty" xml:"AccountStatus,omitempty"`
	PrivilegeStatus  *string `json:"PrivilegeStatus,omitempty" xml:"PrivilegeStatus,omitempty"`
}

func (s DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount) String() string {
	return tea.Prettify(s)
}

func (s DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount) GoString() string {
	return s.String()
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount) SetAccountName(v string) *DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount {
	s.AccountName = &v
	return s
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount) SetAccountPrivilege(v string) *DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount {
	s.AccountPrivilege = &v
	return s
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount) SetAccountStatus(v string) *DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount {
	s.AccountStatus = &v
	return s
}

func (s *DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount) SetPrivilegeStatus(v string) *DescribeDatabasesResponseBodyDatabasesDatabaseAccountsAccount {
	s.PrivilegeStatus = &v
	return s
}

type DescribeDatabasesResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDatabasesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDatabasesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDatabasesResponse) GoString() string {
	return s.String()
}

func (s *DescribeDatabasesResponse) SetHeaders(v map[string]*string) *DescribeDatabasesResponse {
	s.Headers = v
	return s
}

func (s *DescribeDatabasesResponse) SetStatusCode(v int32) *DescribeDatabasesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDatabasesResponse) SetBody(v *DescribeDatabasesResponseBody) *DescribeDatabasesResponse {
	s.Body = v
	return s
}

type DescribeDetachedBackupsRequest struct {
	BackupId             *string `json:"BackupId,omitempty" xml:"BackupId,omitempty"`
	BackupMode           *string `json:"BackupMode,omitempty" xml:"BackupMode,omitempty"`
	BackupRegion         *string `json:"BackupRegion,omitempty" xml:"BackupRegion,omitempty"`
	BackupStatus         *string `json:"BackupStatus,omitempty" xml:"BackupStatus,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EndTime              *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	StartTime            *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeDetachedBackupsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeDetachedBackupsRequest) GoString() string {
	return s.String()
}

func (s *DescribeDetachedBackupsRequest) SetBackupId(v string) *DescribeDetachedBackupsRequest {
	s.BackupId = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetBackupMode(v string) *DescribeDetachedBackupsRequest {
	s.BackupMode = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetBackupRegion(v string) *DescribeDetachedBackupsRequest {
	s.BackupRegion = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetBackupStatus(v string) *DescribeDetachedBackupsRequest {
	s.BackupStatus = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetDBClusterId(v string) *DescribeDetachedBackupsRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetEndTime(v string) *DescribeDetachedBackupsRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetOwnerAccount(v string) *DescribeDetachedBackupsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetOwnerId(v int64) *DescribeDetachedBackupsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetPageNumber(v int32) *DescribeDetachedBackupsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetPageSize(v int32) *DescribeDetachedBackupsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetResourceOwnerAccount(v string) *DescribeDetachedBackupsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetResourceOwnerId(v int64) *DescribeDetachedBackupsRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeDetachedBackupsRequest) SetStartTime(v string) *DescribeDetachedBackupsRequest {
	s.StartTime = &v
	return s
}

type DescribeDetachedBackupsResponseBody struct {
	Items            *DescribeDetachedBackupsResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	PageNumber       *string                                   `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *string                                   `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                                   `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *string                                   `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeDetachedBackupsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeDetachedBackupsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeDetachedBackupsResponseBody) SetItems(v *DescribeDetachedBackupsResponseBodyItems) *DescribeDetachedBackupsResponseBody {
	s.Items = v
	return s
}

func (s *DescribeDetachedBackupsResponseBody) SetPageNumber(v string) *DescribeDetachedBackupsResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBody) SetPageRecordCount(v string) *DescribeDetachedBackupsResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBody) SetRequestId(v string) *DescribeDetachedBackupsResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBody) SetTotalRecordCount(v string) *DescribeDetachedBackupsResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeDetachedBackupsResponseBodyItems struct {
	Backup []*DescribeDetachedBackupsResponseBodyItemsBackup `json:"Backup,omitempty" xml:"Backup,omitempty" type:"Repeated"`
}

func (s DescribeDetachedBackupsResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeDetachedBackupsResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeDetachedBackupsResponseBodyItems) SetBackup(v []*DescribeDetachedBackupsResponseBodyItemsBackup) *DescribeDetachedBackupsResponseBodyItems {
	s.Backup = v
	return s
}

type DescribeDetachedBackupsResponseBodyItemsBackup struct {
	BackupEndTime   *string `json:"BackupEndTime,omitempty" xml:"BackupEndTime,omitempty"`
	BackupId        *string `json:"BackupId,omitempty" xml:"BackupId,omitempty"`
	BackupMethod    *string `json:"BackupMethod,omitempty" xml:"BackupMethod,omitempty"`
	BackupMode      *string `json:"BackupMode,omitempty" xml:"BackupMode,omitempty"`
	BackupSetSize   *string `json:"BackupSetSize,omitempty" xml:"BackupSetSize,omitempty"`
	BackupStartTime *string `json:"BackupStartTime,omitempty" xml:"BackupStartTime,omitempty"`
	BackupStatus    *string `json:"BackupStatus,omitempty" xml:"BackupStatus,omitempty"`
	BackupType      *string `json:"BackupType,omitempty" xml:"BackupType,omitempty"`
	BackupsLevel    *string `json:"BackupsLevel,omitempty" xml:"BackupsLevel,omitempty"`
	ConsistentTime  *string `json:"ConsistentTime,omitempty" xml:"ConsistentTime,omitempty"`
	DBClusterId     *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	IsAvail         *string `json:"IsAvail,omitempty" xml:"IsAvail,omitempty"`
	StoreStatus     *string `json:"StoreStatus,omitempty" xml:"StoreStatus,omitempty"`
}

func (s DescribeDetachedBackupsResponseBodyItemsBackup) String() string {
	return tea.Prettify(s)
}

func (s DescribeDetachedBackupsResponseBodyItemsBackup) GoString() string {
	return s.String()
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupEndTime(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupEndTime = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupId(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupId = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupMethod(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupMethod = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupMode(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupMode = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupSetSize(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupSetSize = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupStartTime(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupStartTime = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupStatus(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupStatus = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupType(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupType = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetBackupsLevel(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.BackupsLevel = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetConsistentTime(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.ConsistentTime = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetDBClusterId(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.DBClusterId = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetIsAvail(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.IsAvail = &v
	return s
}

func (s *DescribeDetachedBackupsResponseBodyItemsBackup) SetStoreStatus(v string) *DescribeDetachedBackupsResponseBodyItemsBackup {
	s.StoreStatus = &v
	return s
}

type DescribeDetachedBackupsResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeDetachedBackupsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeDetachedBackupsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeDetachedBackupsResponse) GoString() string {
	return s.String()
}

func (s *DescribeDetachedBackupsResponse) SetHeaders(v map[string]*string) *DescribeDetachedBackupsResponse {
	s.Headers = v
	return s
}

func (s *DescribeDetachedBackupsResponse) SetStatusCode(v int32) *DescribeDetachedBackupsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeDetachedBackupsResponse) SetBody(v *DescribeDetachedBackupsResponseBody) *DescribeDetachedBackupsResponse {
	s.Body = v
	return s
}

type DescribeGlobalDatabaseNetworkRequest struct {
	GDNId                *string `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s DescribeGlobalDatabaseNetworkRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworkRequest) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworkRequest) SetGDNId(v string) *DescribeGlobalDatabaseNetworkRequest {
	s.GDNId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkRequest) SetOwnerAccount(v string) *DescribeGlobalDatabaseNetworkRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkRequest) SetOwnerId(v int64) *DescribeGlobalDatabaseNetworkRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkRequest) SetResourceOwnerAccount(v string) *DescribeGlobalDatabaseNetworkRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkRequest) SetResourceOwnerId(v int64) *DescribeGlobalDatabaseNetworkRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkRequest) SetSecurityToken(v string) *DescribeGlobalDatabaseNetworkRequest {
	s.SecurityToken = &v
	return s
}

type DescribeGlobalDatabaseNetworkResponseBody struct {
	Connections    []*DescribeGlobalDatabaseNetworkResponseBodyConnections `json:"Connections,omitempty" xml:"Connections,omitempty" type:"Repeated"`
	CreateTime     *string                                                 `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	DBClusters     []*DescribeGlobalDatabaseNetworkResponseBodyDBClusters  `json:"DBClusters,omitempty" xml:"DBClusters,omitempty" type:"Repeated"`
	DBType         *string                                                 `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion      *string                                                 `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	GDNDescription *string                                                 `json:"GDNDescription,omitempty" xml:"GDNDescription,omitempty"`
	GDNId          *string                                                 `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	GDNStatus      *string                                                 `json:"GDNStatus,omitempty" xml:"GDNStatus,omitempty"`
	RequestId      *string                                                 `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeGlobalDatabaseNetworkResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworkResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetConnections(v []*DescribeGlobalDatabaseNetworkResponseBodyConnections) *DescribeGlobalDatabaseNetworkResponseBody {
	s.Connections = v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetCreateTime(v string) *DescribeGlobalDatabaseNetworkResponseBody {
	s.CreateTime = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetDBClusters(v []*DescribeGlobalDatabaseNetworkResponseBodyDBClusters) *DescribeGlobalDatabaseNetworkResponseBody {
	s.DBClusters = v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetDBType(v string) *DescribeGlobalDatabaseNetworkResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetDBVersion(v string) *DescribeGlobalDatabaseNetworkResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetGDNDescription(v string) *DescribeGlobalDatabaseNetworkResponseBody {
	s.GDNDescription = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetGDNId(v string) *DescribeGlobalDatabaseNetworkResponseBody {
	s.GDNId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetGDNStatus(v string) *DescribeGlobalDatabaseNetworkResponseBody {
	s.GDNStatus = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBody) SetRequestId(v string) *DescribeGlobalDatabaseNetworkResponseBody {
	s.RequestId = &v
	return s
}

type DescribeGlobalDatabaseNetworkResponseBodyConnections struct {
	ConnectionString *string `json:"ConnectionString,omitempty" xml:"ConnectionString,omitempty"`
	NetType          *string `json:"NetType,omitempty" xml:"NetType,omitempty"`
	Port             *string `json:"Port,omitempty" xml:"Port,omitempty"`
}

func (s DescribeGlobalDatabaseNetworkResponseBodyConnections) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworkResponseBodyConnections) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyConnections) SetConnectionString(v string) *DescribeGlobalDatabaseNetworkResponseBodyConnections {
	s.ConnectionString = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyConnections) SetNetType(v string) *DescribeGlobalDatabaseNetworkResponseBodyConnections {
	s.NetType = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyConnections) SetPort(v string) *DescribeGlobalDatabaseNetworkResponseBodyConnections {
	s.Port = &v
	return s
}

type DescribeGlobalDatabaseNetworkResponseBodyDBClusters struct {
	DBClusterDescription *string                                                       `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBClusterId          *string                                                       `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBClusterStatus      *string                                                       `json:"DBClusterStatus,omitempty" xml:"DBClusterStatus,omitempty"`
	DBNodeClass          *string                                                       `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
	DBNodes              []*DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes `json:"DBNodes,omitempty" xml:"DBNodes,omitempty" type:"Repeated"`
	DBType               *string                                                       `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion            *string                                                       `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	ExpireTime           *string                                                       `json:"ExpireTime,omitempty" xml:"ExpireTime,omitempty"`
	Expired              *string                                                       `json:"Expired,omitempty" xml:"Expired,omitempty"`
	PayType              *string                                                       `json:"PayType,omitempty" xml:"PayType,omitempty"`
	RegionId             *string                                                       `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ReplicaLag           *string                                                       `json:"ReplicaLag,omitempty" xml:"ReplicaLag,omitempty"`
	Role                 *string                                                       `json:"Role,omitempty" xml:"Role,omitempty"`
	StorageUsed          *string                                                       `json:"StorageUsed,omitempty" xml:"StorageUsed,omitempty"`
}

func (s DescribeGlobalDatabaseNetworkResponseBodyDBClusters) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworkResponseBodyDBClusters) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetDBClusterDescription(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.DBClusterDescription = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetDBClusterId(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.DBClusterId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetDBClusterStatus(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.DBClusterStatus = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetDBNodeClass(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.DBNodeClass = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetDBNodes(v []*DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.DBNodes = v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetDBType(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.DBType = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetDBVersion(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.DBVersion = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetExpireTime(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.ExpireTime = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetExpired(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.Expired = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetPayType(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.PayType = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetRegionId(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.RegionId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetReplicaLag(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.ReplicaLag = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetRole(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.Role = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClusters) SetStorageUsed(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClusters {
	s.StorageUsed = &v
	return s
}

type DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes struct {
	CreationTime     *string `json:"CreationTime,omitempty" xml:"CreationTime,omitempty"`
	DBNodeClass      *string `json:"DBNodeClass,omitempty" xml:"DBNodeClass,omitempty"`
	DBNodeId         *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	DBNodeRole       *string `json:"DBNodeRole,omitempty" xml:"DBNodeRole,omitempty"`
	DBNodeStatus     *string `json:"DBNodeStatus,omitempty" xml:"DBNodeStatus,omitempty"`
	FailoverPriority *int32  `json:"FailoverPriority,omitempty" xml:"FailoverPriority,omitempty"`
	MaxConnections   *int32  `json:"MaxConnections,omitempty" xml:"MaxConnections,omitempty"`
	MaxIOPS          *int32  `json:"MaxIOPS,omitempty" xml:"MaxIOPS,omitempty"`
	ZoneId           *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetCreationTime(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.CreationTime = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetDBNodeClass(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.DBNodeClass = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetDBNodeId(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.DBNodeId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetDBNodeRole(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.DBNodeRole = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetDBNodeStatus(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.DBNodeStatus = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetFailoverPriority(v int32) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.FailoverPriority = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetMaxConnections(v int32) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.MaxConnections = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetMaxIOPS(v int32) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.MaxIOPS = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes) SetZoneId(v string) *DescribeGlobalDatabaseNetworkResponseBodyDBClustersDBNodes {
	s.ZoneId = &v
	return s
}

type DescribeGlobalDatabaseNetworkResponse struct {
	Headers    map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeGlobalDatabaseNetworkResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeGlobalDatabaseNetworkResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworkResponse) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworkResponse) SetHeaders(v map[string]*string) *DescribeGlobalDatabaseNetworkResponse {
	s.Headers = v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponse) SetStatusCode(v int32) *DescribeGlobalDatabaseNetworkResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworkResponse) SetBody(v *DescribeGlobalDatabaseNetworkResponseBody) *DescribeGlobalDatabaseNetworkResponse {
	s.Body = v
	return s
}

type DescribeGlobalDatabaseNetworksRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	GDNDescription       *string `json:"GDNDescription,omitempty" xml:"GDNDescription,omitempty"`
	GDNId                *string `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s DescribeGlobalDatabaseNetworksRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworksRequest) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetDBClusterId(v string) *DescribeGlobalDatabaseNetworksRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetGDNDescription(v string) *DescribeGlobalDatabaseNetworksRequest {
	s.GDNDescription = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetGDNId(v string) *DescribeGlobalDatabaseNetworksRequest {
	s.GDNId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetOwnerAccount(v string) *DescribeGlobalDatabaseNetworksRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetOwnerId(v int64) *DescribeGlobalDatabaseNetworksRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetPageNumber(v int32) *DescribeGlobalDatabaseNetworksRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetPageSize(v int32) *DescribeGlobalDatabaseNetworksRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetResourceOwnerAccount(v string) *DescribeGlobalDatabaseNetworksRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetResourceOwnerId(v int64) *DescribeGlobalDatabaseNetworksRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksRequest) SetSecurityToken(v string) *DescribeGlobalDatabaseNetworksRequest {
	s.SecurityToken = &v
	return s
}

type DescribeGlobalDatabaseNetworksResponseBody struct {
	Items            []*DescribeGlobalDatabaseNetworksResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Repeated"`
	PageNumber       *int32                                             `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *int32                                             `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                                            `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *int32                                             `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeGlobalDatabaseNetworksResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworksResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworksResponseBody) SetItems(v []*DescribeGlobalDatabaseNetworksResponseBodyItems) *DescribeGlobalDatabaseNetworksResponseBody {
	s.Items = v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBody) SetPageNumber(v int32) *DescribeGlobalDatabaseNetworksResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBody) SetPageRecordCount(v int32) *DescribeGlobalDatabaseNetworksResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBody) SetRequestId(v string) *DescribeGlobalDatabaseNetworksResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBody) SetTotalRecordCount(v int32) *DescribeGlobalDatabaseNetworksResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeGlobalDatabaseNetworksResponseBodyItems struct {
	CreateTime     *string                                                      `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	DBClusters     []*DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters `json:"DBClusters,omitempty" xml:"DBClusters,omitempty" type:"Repeated"`
	DBType         *string                                                      `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion      *string                                                      `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	GDNDescription *string                                                      `json:"GDNDescription,omitempty" xml:"GDNDescription,omitempty"`
	GDNId          *string                                                      `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	GDNStatus      *string                                                      `json:"GDNStatus,omitempty" xml:"GDNStatus,omitempty"`
}

func (s DescribeGlobalDatabaseNetworksResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworksResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItems) SetCreateTime(v string) *DescribeGlobalDatabaseNetworksResponseBodyItems {
	s.CreateTime = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItems) SetDBClusters(v []*DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters) *DescribeGlobalDatabaseNetworksResponseBodyItems {
	s.DBClusters = v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItems) SetDBType(v string) *DescribeGlobalDatabaseNetworksResponseBodyItems {
	s.DBType = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItems) SetDBVersion(v string) *DescribeGlobalDatabaseNetworksResponseBodyItems {
	s.DBVersion = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItems) SetGDNDescription(v string) *DescribeGlobalDatabaseNetworksResponseBodyItems {
	s.GDNDescription = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItems) SetGDNId(v string) *DescribeGlobalDatabaseNetworksResponseBodyItems {
	s.GDNId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItems) SetGDNStatus(v string) *DescribeGlobalDatabaseNetworksResponseBodyItems {
	s.GDNStatus = &v
	return s
}

type DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters struct {
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	RegionId    *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	Role        *string `json:"Role,omitempty" xml:"Role,omitempty"`
}

func (s DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters) SetDBClusterId(v string) *DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters {
	s.DBClusterId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters) SetRegionId(v string) *DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters {
	s.RegionId = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters) SetRole(v string) *DescribeGlobalDatabaseNetworksResponseBodyItemsDBClusters {
	s.Role = &v
	return s
}

type DescribeGlobalDatabaseNetworksResponse struct {
	Headers    map[string]*string                          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeGlobalDatabaseNetworksResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeGlobalDatabaseNetworksResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeGlobalDatabaseNetworksResponse) GoString() string {
	return s.String()
}

func (s *DescribeGlobalDatabaseNetworksResponse) SetHeaders(v map[string]*string) *DescribeGlobalDatabaseNetworksResponse {
	s.Headers = v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponse) SetStatusCode(v int32) *DescribeGlobalDatabaseNetworksResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeGlobalDatabaseNetworksResponse) SetBody(v *DescribeGlobalDatabaseNetworksResponseBody) *DescribeGlobalDatabaseNetworksResponse {
	s.Body = v
	return s
}

type DescribeLogBackupPolicyRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeLogBackupPolicyRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeLogBackupPolicyRequest) GoString() string {
	return s.String()
}

func (s *DescribeLogBackupPolicyRequest) SetDBClusterId(v string) *DescribeLogBackupPolicyRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeLogBackupPolicyRequest) SetOwnerAccount(v string) *DescribeLogBackupPolicyRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeLogBackupPolicyRequest) SetOwnerId(v int64) *DescribeLogBackupPolicyRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeLogBackupPolicyRequest) SetResourceOwnerAccount(v string) *DescribeLogBackupPolicyRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeLogBackupPolicyRequest) SetResourceOwnerId(v int64) *DescribeLogBackupPolicyRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeLogBackupPolicyResponseBody struct {
	EnableBackupLog                       *int32  `json:"EnableBackupLog,omitempty" xml:"EnableBackupLog,omitempty"`
	LogBackupAnotherRegionRegion          *string `json:"LogBackupAnotherRegionRegion,omitempty" xml:"LogBackupAnotherRegionRegion,omitempty"`
	LogBackupAnotherRegionRetentionPeriod *string `json:"LogBackupAnotherRegionRetentionPeriod,omitempty" xml:"LogBackupAnotherRegionRetentionPeriod,omitempty"`
	LogBackupRetentionPeriod              *int32  `json:"LogBackupRetentionPeriod,omitempty" xml:"LogBackupRetentionPeriod,omitempty"`
	RequestId                             *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeLogBackupPolicyResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeLogBackupPolicyResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeLogBackupPolicyResponseBody) SetEnableBackupLog(v int32) *DescribeLogBackupPolicyResponseBody {
	s.EnableBackupLog = &v
	return s
}

func (s *DescribeLogBackupPolicyResponseBody) SetLogBackupAnotherRegionRegion(v string) *DescribeLogBackupPolicyResponseBody {
	s.LogBackupAnotherRegionRegion = &v
	return s
}

func (s *DescribeLogBackupPolicyResponseBody) SetLogBackupAnotherRegionRetentionPeriod(v string) *DescribeLogBackupPolicyResponseBody {
	s.LogBackupAnotherRegionRetentionPeriod = &v
	return s
}

func (s *DescribeLogBackupPolicyResponseBody) SetLogBackupRetentionPeriod(v int32) *DescribeLogBackupPolicyResponseBody {
	s.LogBackupRetentionPeriod = &v
	return s
}

func (s *DescribeLogBackupPolicyResponseBody) SetRequestId(v string) *DescribeLogBackupPolicyResponseBody {
	s.RequestId = &v
	return s
}

type DescribeLogBackupPolicyResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeLogBackupPolicyResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeLogBackupPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeLogBackupPolicyResponse) GoString() string {
	return s.String()
}

func (s *DescribeLogBackupPolicyResponse) SetHeaders(v map[string]*string) *DescribeLogBackupPolicyResponse {
	s.Headers = v
	return s
}

func (s *DescribeLogBackupPolicyResponse) SetStatusCode(v int32) *DescribeLogBackupPolicyResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeLogBackupPolicyResponse) SetBody(v *DescribeLogBackupPolicyResponseBody) *DescribeLogBackupPolicyResponse {
	s.Body = v
	return s
}

type DescribeMaskingRulesRequest struct {
	DBClusterId  *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	RuleNameList *string `json:"RuleNameList,omitempty" xml:"RuleNameList,omitempty"`
}

func (s DescribeMaskingRulesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeMaskingRulesRequest) GoString() string {
	return s.String()
}

func (s *DescribeMaskingRulesRequest) SetDBClusterId(v string) *DescribeMaskingRulesRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeMaskingRulesRequest) SetRuleNameList(v string) *DescribeMaskingRulesRequest {
	s.RuleNameList = &v
	return s
}

type DescribeMaskingRulesResponseBody struct {
	DBClusterId *string                               `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	Data        *DescribeMaskingRulesResponseBodyData `json:"Data,omitempty" xml:"Data,omitempty" type:"Struct"`
	Message     *string                               `json:"Message,omitempty" xml:"Message,omitempty"`
	RequestId   *string                               `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	Success     *bool                                 `json:"Success,omitempty" xml:"Success,omitempty"`
}

func (s DescribeMaskingRulesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeMaskingRulesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeMaskingRulesResponseBody) SetDBClusterId(v string) *DescribeMaskingRulesResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeMaskingRulesResponseBody) SetData(v *DescribeMaskingRulesResponseBodyData) *DescribeMaskingRulesResponseBody {
	s.Data = v
	return s
}

func (s *DescribeMaskingRulesResponseBody) SetMessage(v string) *DescribeMaskingRulesResponseBody {
	s.Message = &v
	return s
}

func (s *DescribeMaskingRulesResponseBody) SetRequestId(v string) *DescribeMaskingRulesResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeMaskingRulesResponseBody) SetSuccess(v bool) *DescribeMaskingRulesResponseBody {
	s.Success = &v
	return s
}

type DescribeMaskingRulesResponseBodyData struct {
	RuleList []*string `json:"RuleList,omitempty" xml:"RuleList,omitempty" type:"Repeated"`
}

func (s DescribeMaskingRulesResponseBodyData) String() string {
	return tea.Prettify(s)
}

func (s DescribeMaskingRulesResponseBodyData) GoString() string {
	return s.String()
}

func (s *DescribeMaskingRulesResponseBodyData) SetRuleList(v []*string) *DescribeMaskingRulesResponseBodyData {
	s.RuleList = v
	return s
}

type DescribeMaskingRulesResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeMaskingRulesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeMaskingRulesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeMaskingRulesResponse) GoString() string {
	return s.String()
}

func (s *DescribeMaskingRulesResponse) SetHeaders(v map[string]*string) *DescribeMaskingRulesResponse {
	s.Headers = v
	return s
}

func (s *DescribeMaskingRulesResponse) SetStatusCode(v int32) *DescribeMaskingRulesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeMaskingRulesResponse) SetBody(v *DescribeMaskingRulesResponseBody) *DescribeMaskingRulesResponse {
	s.Body = v
	return s
}

type DescribeMetaListRequest struct {
	BackupId             *string `json:"BackupId,omitempty" xml:"BackupId,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	GetDbName            *string `json:"GetDbName,omitempty" xml:"GetDbName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	RestoreTime          *string `json:"RestoreTime,omitempty" xml:"RestoreTime,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s DescribeMetaListRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeMetaListRequest) GoString() string {
	return s.String()
}

func (s *DescribeMetaListRequest) SetBackupId(v string) *DescribeMetaListRequest {
	s.BackupId = &v
	return s
}

func (s *DescribeMetaListRequest) SetDBClusterId(v string) *DescribeMetaListRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeMetaListRequest) SetGetDbName(v string) *DescribeMetaListRequest {
	s.GetDbName = &v
	return s
}

func (s *DescribeMetaListRequest) SetOwnerAccount(v string) *DescribeMetaListRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeMetaListRequest) SetOwnerId(v int64) *DescribeMetaListRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeMetaListRequest) SetPageNumber(v int32) *DescribeMetaListRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeMetaListRequest) SetPageSize(v int32) *DescribeMetaListRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeMetaListRequest) SetResourceOwnerAccount(v string) *DescribeMetaListRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeMetaListRequest) SetResourceOwnerId(v int64) *DescribeMetaListRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeMetaListRequest) SetRestoreTime(v string) *DescribeMetaListRequest {
	s.RestoreTime = &v
	return s
}

func (s *DescribeMetaListRequest) SetSecurityToken(v string) *DescribeMetaListRequest {
	s.SecurityToken = &v
	return s
}

type DescribeMetaListResponseBody struct {
	DBClusterId      *string                              `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	Items            []*DescribeMetaListResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Repeated"`
	PageNumber       *string                              `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize         *string                              `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	RequestId        *string                              `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalPageCount   *string                              `json:"TotalPageCount,omitempty" xml:"TotalPageCount,omitempty"`
	TotalRecordCount *string                              `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeMetaListResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeMetaListResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeMetaListResponseBody) SetDBClusterId(v string) *DescribeMetaListResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeMetaListResponseBody) SetItems(v []*DescribeMetaListResponseBodyItems) *DescribeMetaListResponseBody {
	s.Items = v
	return s
}

func (s *DescribeMetaListResponseBody) SetPageNumber(v string) *DescribeMetaListResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeMetaListResponseBody) SetPageSize(v string) *DescribeMetaListResponseBody {
	s.PageSize = &v
	return s
}

func (s *DescribeMetaListResponseBody) SetRequestId(v string) *DescribeMetaListResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeMetaListResponseBody) SetTotalPageCount(v string) *DescribeMetaListResponseBody {
	s.TotalPageCount = &v
	return s
}

func (s *DescribeMetaListResponseBody) SetTotalRecordCount(v string) *DescribeMetaListResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeMetaListResponseBodyItems struct {
	Database *string   `json:"Database,omitempty" xml:"Database,omitempty"`
	Tables   []*string `json:"Tables,omitempty" xml:"Tables,omitempty" type:"Repeated"`
}

func (s DescribeMetaListResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeMetaListResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeMetaListResponseBodyItems) SetDatabase(v string) *DescribeMetaListResponseBodyItems {
	s.Database = &v
	return s
}

func (s *DescribeMetaListResponseBodyItems) SetTables(v []*string) *DescribeMetaListResponseBodyItems {
	s.Tables = v
	return s
}

type DescribeMetaListResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeMetaListResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeMetaListResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeMetaListResponse) GoString() string {
	return s.String()
}

func (s *DescribeMetaListResponse) SetHeaders(v map[string]*string) *DescribeMetaListResponse {
	s.Headers = v
	return s
}

func (s *DescribeMetaListResponse) SetStatusCode(v int32) *DescribeMetaListResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeMetaListResponse) SetBody(v *DescribeMetaListResponseBody) *DescribeMetaListResponse {
	s.Body = v
	return s
}

type DescribeParameterGroupRequest struct {
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ParameterGroupId     *string `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeParameterGroupRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupRequest) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupRequest) SetOwnerAccount(v string) *DescribeParameterGroupRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeParameterGroupRequest) SetOwnerId(v int64) *DescribeParameterGroupRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeParameterGroupRequest) SetParameterGroupId(v string) *DescribeParameterGroupRequest {
	s.ParameterGroupId = &v
	return s
}

func (s *DescribeParameterGroupRequest) SetRegionId(v string) *DescribeParameterGroupRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeParameterGroupRequest) SetResourceOwnerAccount(v string) *DescribeParameterGroupRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeParameterGroupRequest) SetResourceOwnerId(v int64) *DescribeParameterGroupRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeParameterGroupResponseBody struct {
	ParameterGroup []*DescribeParameterGroupResponseBodyParameterGroup `json:"ParameterGroup,omitempty" xml:"ParameterGroup,omitempty" type:"Repeated"`
	RequestId      *string                                             `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeParameterGroupResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupResponseBody) SetParameterGroup(v []*DescribeParameterGroupResponseBodyParameterGroup) *DescribeParameterGroupResponseBody {
	s.ParameterGroup = v
	return s
}

func (s *DescribeParameterGroupResponseBody) SetRequestId(v string) *DescribeParameterGroupResponseBody {
	s.RequestId = &v
	return s
}

type DescribeParameterGroupResponseBodyParameterGroup struct {
	CreateTime         *string                                                            `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	DBType             *string                                                            `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion          *string                                                            `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	ForceRestart       *string                                                            `json:"ForceRestart,omitempty" xml:"ForceRestart,omitempty"`
	ParameterCounts    *int32                                                             `json:"ParameterCounts,omitempty" xml:"ParameterCounts,omitempty"`
	ParameterDetail    []*DescribeParameterGroupResponseBodyParameterGroupParameterDetail `json:"ParameterDetail,omitempty" xml:"ParameterDetail,omitempty" type:"Repeated"`
	ParameterGroupDesc *string                                                            `json:"ParameterGroupDesc,omitempty" xml:"ParameterGroupDesc,omitempty"`
	ParameterGroupId   *string                                                            `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	ParameterGroupName *string                                                            `json:"ParameterGroupName,omitempty" xml:"ParameterGroupName,omitempty"`
	ParameterGroupType *string                                                            `json:"ParameterGroupType,omitempty" xml:"ParameterGroupType,omitempty"`
}

func (s DescribeParameterGroupResponseBodyParameterGroup) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupResponseBodyParameterGroup) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetCreateTime(v string) *DescribeParameterGroupResponseBodyParameterGroup {
	s.CreateTime = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetDBType(v string) *DescribeParameterGroupResponseBodyParameterGroup {
	s.DBType = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetDBVersion(v string) *DescribeParameterGroupResponseBodyParameterGroup {
	s.DBVersion = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetForceRestart(v string) *DescribeParameterGroupResponseBodyParameterGroup {
	s.ForceRestart = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetParameterCounts(v int32) *DescribeParameterGroupResponseBodyParameterGroup {
	s.ParameterCounts = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetParameterDetail(v []*DescribeParameterGroupResponseBodyParameterGroupParameterDetail) *DescribeParameterGroupResponseBodyParameterGroup {
	s.ParameterDetail = v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetParameterGroupDesc(v string) *DescribeParameterGroupResponseBodyParameterGroup {
	s.ParameterGroupDesc = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetParameterGroupId(v string) *DescribeParameterGroupResponseBodyParameterGroup {
	s.ParameterGroupId = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetParameterGroupName(v string) *DescribeParameterGroupResponseBodyParameterGroup {
	s.ParameterGroupName = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroup) SetParameterGroupType(v string) *DescribeParameterGroupResponseBodyParameterGroup {
	s.ParameterGroupType = &v
	return s
}

type DescribeParameterGroupResponseBodyParameterGroupParameterDetail struct {
	ParamName  *string `json:"ParamName,omitempty" xml:"ParamName,omitempty"`
	ParamValue *string `json:"ParamValue,omitempty" xml:"ParamValue,omitempty"`
}

func (s DescribeParameterGroupResponseBodyParameterGroupParameterDetail) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupResponseBodyParameterGroupParameterDetail) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupResponseBodyParameterGroupParameterDetail) SetParamName(v string) *DescribeParameterGroupResponseBodyParameterGroupParameterDetail {
	s.ParamName = &v
	return s
}

func (s *DescribeParameterGroupResponseBodyParameterGroupParameterDetail) SetParamValue(v string) *DescribeParameterGroupResponseBodyParameterGroupParameterDetail {
	s.ParamValue = &v
	return s
}

type DescribeParameterGroupResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeParameterGroupResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeParameterGroupResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupResponse) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupResponse) SetHeaders(v map[string]*string) *DescribeParameterGroupResponse {
	s.Headers = v
	return s
}

func (s *DescribeParameterGroupResponse) SetStatusCode(v int32) *DescribeParameterGroupResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeParameterGroupResponse) SetBody(v *DescribeParameterGroupResponseBody) *DescribeParameterGroupResponse {
	s.Body = v
	return s
}

type DescribeParameterGroupsRequest struct {
	DBType               *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion            *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeParameterGroupsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupsRequest) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupsRequest) SetDBType(v string) *DescribeParameterGroupsRequest {
	s.DBType = &v
	return s
}

func (s *DescribeParameterGroupsRequest) SetDBVersion(v string) *DescribeParameterGroupsRequest {
	s.DBVersion = &v
	return s
}

func (s *DescribeParameterGroupsRequest) SetOwnerAccount(v string) *DescribeParameterGroupsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeParameterGroupsRequest) SetOwnerId(v int64) *DescribeParameterGroupsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeParameterGroupsRequest) SetRegionId(v string) *DescribeParameterGroupsRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeParameterGroupsRequest) SetResourceOwnerAccount(v string) *DescribeParameterGroupsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeParameterGroupsRequest) SetResourceOwnerId(v int64) *DescribeParameterGroupsRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeParameterGroupsResponseBody struct {
	ParameterGroups []*DescribeParameterGroupsResponseBodyParameterGroups `json:"ParameterGroups,omitempty" xml:"ParameterGroups,omitempty" type:"Repeated"`
	RequestId       *string                                               `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeParameterGroupsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupsResponseBody) SetParameterGroups(v []*DescribeParameterGroupsResponseBodyParameterGroups) *DescribeParameterGroupsResponseBody {
	s.ParameterGroups = v
	return s
}

func (s *DescribeParameterGroupsResponseBody) SetRequestId(v string) *DescribeParameterGroupsResponseBody {
	s.RequestId = &v
	return s
}

type DescribeParameterGroupsResponseBodyParameterGroups struct {
	CreateTime         *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	DBType             *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion          *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	ForceRestart       *string `json:"ForceRestart,omitempty" xml:"ForceRestart,omitempty"`
	ParameterCounts    *int64  `json:"ParameterCounts,omitempty" xml:"ParameterCounts,omitempty"`
	ParameterGroupDesc *string `json:"ParameterGroupDesc,omitempty" xml:"ParameterGroupDesc,omitempty"`
	ParameterGroupId   *string `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	ParameterGroupName *string `json:"ParameterGroupName,omitempty" xml:"ParameterGroupName,omitempty"`
	ParameterGroupType *string `json:"ParameterGroupType,omitempty" xml:"ParameterGroupType,omitempty"`
}

func (s DescribeParameterGroupsResponseBodyParameterGroups) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupsResponseBodyParameterGroups) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetCreateTime(v string) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.CreateTime = &v
	return s
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetDBType(v string) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.DBType = &v
	return s
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetDBVersion(v string) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.DBVersion = &v
	return s
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetForceRestart(v string) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.ForceRestart = &v
	return s
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetParameterCounts(v int64) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.ParameterCounts = &v
	return s
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetParameterGroupDesc(v string) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.ParameterGroupDesc = &v
	return s
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetParameterGroupId(v string) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.ParameterGroupId = &v
	return s
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetParameterGroupName(v string) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.ParameterGroupName = &v
	return s
}

func (s *DescribeParameterGroupsResponseBodyParameterGroups) SetParameterGroupType(v string) *DescribeParameterGroupsResponseBodyParameterGroups {
	s.ParameterGroupType = &v
	return s
}

type DescribeParameterGroupsResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeParameterGroupsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeParameterGroupsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterGroupsResponse) GoString() string {
	return s.String()
}

func (s *DescribeParameterGroupsResponse) SetHeaders(v map[string]*string) *DescribeParameterGroupsResponse {
	s.Headers = v
	return s
}

func (s *DescribeParameterGroupsResponse) SetStatusCode(v int32) *DescribeParameterGroupsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeParameterGroupsResponse) SetBody(v *DescribeParameterGroupsResponseBody) *DescribeParameterGroupsResponse {
	s.Body = v
	return s
}

type DescribeParameterTemplatesRequest struct {
	DBType               *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion            *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeParameterTemplatesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterTemplatesRequest) GoString() string {
	return s.String()
}

func (s *DescribeParameterTemplatesRequest) SetDBType(v string) *DescribeParameterTemplatesRequest {
	s.DBType = &v
	return s
}

func (s *DescribeParameterTemplatesRequest) SetDBVersion(v string) *DescribeParameterTemplatesRequest {
	s.DBVersion = &v
	return s
}

func (s *DescribeParameterTemplatesRequest) SetOwnerAccount(v string) *DescribeParameterTemplatesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeParameterTemplatesRequest) SetOwnerId(v int64) *DescribeParameterTemplatesRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeParameterTemplatesRequest) SetRegionId(v string) *DescribeParameterTemplatesRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeParameterTemplatesRequest) SetResourceOwnerAccount(v string) *DescribeParameterTemplatesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeParameterTemplatesRequest) SetResourceOwnerId(v int64) *DescribeParameterTemplatesRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeParameterTemplatesResponseBody struct {
	DBType         *string                                           `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion      *string                                           `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	Engine         *string                                           `json:"Engine,omitempty" xml:"Engine,omitempty"`
	ParameterCount *string                                           `json:"ParameterCount,omitempty" xml:"ParameterCount,omitempty"`
	Parameters     *DescribeParameterTemplatesResponseBodyParameters `json:"Parameters,omitempty" xml:"Parameters,omitempty" type:"Struct"`
	RequestId      *string                                           `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeParameterTemplatesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterTemplatesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeParameterTemplatesResponseBody) SetDBType(v string) *DescribeParameterTemplatesResponseBody {
	s.DBType = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBody) SetDBVersion(v string) *DescribeParameterTemplatesResponseBody {
	s.DBVersion = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBody) SetEngine(v string) *DescribeParameterTemplatesResponseBody {
	s.Engine = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBody) SetParameterCount(v string) *DescribeParameterTemplatesResponseBody {
	s.ParameterCount = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBody) SetParameters(v *DescribeParameterTemplatesResponseBodyParameters) *DescribeParameterTemplatesResponseBody {
	s.Parameters = v
	return s
}

func (s *DescribeParameterTemplatesResponseBody) SetRequestId(v string) *DescribeParameterTemplatesResponseBody {
	s.RequestId = &v
	return s
}

type DescribeParameterTemplatesResponseBodyParameters struct {
	TemplateRecord []*DescribeParameterTemplatesResponseBodyParametersTemplateRecord `json:"TemplateRecord,omitempty" xml:"TemplateRecord,omitempty" type:"Repeated"`
}

func (s DescribeParameterTemplatesResponseBodyParameters) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterTemplatesResponseBodyParameters) GoString() string {
	return s.String()
}

func (s *DescribeParameterTemplatesResponseBodyParameters) SetTemplateRecord(v []*DescribeParameterTemplatesResponseBodyParametersTemplateRecord) *DescribeParameterTemplatesResponseBodyParameters {
	s.TemplateRecord = v
	return s
}

type DescribeParameterTemplatesResponseBodyParametersTemplateRecord struct {
	CheckingCode         *string `json:"CheckingCode,omitempty" xml:"CheckingCode,omitempty"`
	ForceModify          *string `json:"ForceModify,omitempty" xml:"ForceModify,omitempty"`
	ForceRestart         *string `json:"ForceRestart,omitempty" xml:"ForceRestart,omitempty"`
	IsNodeAvailable      *string `json:"IsNodeAvailable,omitempty" xml:"IsNodeAvailable,omitempty"`
	ParamRelyRule        *string `json:"ParamRelyRule,omitempty" xml:"ParamRelyRule,omitempty"`
	ParameterDescription *string `json:"ParameterDescription,omitempty" xml:"ParameterDescription,omitempty"`
	ParameterName        *string `json:"ParameterName,omitempty" xml:"ParameterName,omitempty"`
	ParameterValue       *string `json:"ParameterValue,omitempty" xml:"ParameterValue,omitempty"`
}

func (s DescribeParameterTemplatesResponseBodyParametersTemplateRecord) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterTemplatesResponseBodyParametersTemplateRecord) GoString() string {
	return s.String()
}

func (s *DescribeParameterTemplatesResponseBodyParametersTemplateRecord) SetCheckingCode(v string) *DescribeParameterTemplatesResponseBodyParametersTemplateRecord {
	s.CheckingCode = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBodyParametersTemplateRecord) SetForceModify(v string) *DescribeParameterTemplatesResponseBodyParametersTemplateRecord {
	s.ForceModify = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBodyParametersTemplateRecord) SetForceRestart(v string) *DescribeParameterTemplatesResponseBodyParametersTemplateRecord {
	s.ForceRestart = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBodyParametersTemplateRecord) SetIsNodeAvailable(v string) *DescribeParameterTemplatesResponseBodyParametersTemplateRecord {
	s.IsNodeAvailable = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBodyParametersTemplateRecord) SetParamRelyRule(v string) *DescribeParameterTemplatesResponseBodyParametersTemplateRecord {
	s.ParamRelyRule = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBodyParametersTemplateRecord) SetParameterDescription(v string) *DescribeParameterTemplatesResponseBodyParametersTemplateRecord {
	s.ParameterDescription = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBodyParametersTemplateRecord) SetParameterName(v string) *DescribeParameterTemplatesResponseBodyParametersTemplateRecord {
	s.ParameterName = &v
	return s
}

func (s *DescribeParameterTemplatesResponseBodyParametersTemplateRecord) SetParameterValue(v string) *DescribeParameterTemplatesResponseBodyParametersTemplateRecord {
	s.ParameterValue = &v
	return s
}

type DescribeParameterTemplatesResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeParameterTemplatesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeParameterTemplatesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeParameterTemplatesResponse) GoString() string {
	return s.String()
}

func (s *DescribeParameterTemplatesResponse) SetHeaders(v map[string]*string) *DescribeParameterTemplatesResponse {
	s.Headers = v
	return s
}

func (s *DescribeParameterTemplatesResponse) SetStatusCode(v int32) *DescribeParameterTemplatesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeParameterTemplatesResponse) SetBody(v *DescribeParameterTemplatesResponseBody) *DescribeParameterTemplatesResponse {
	s.Body = v
	return s
}

type DescribePendingMaintenanceActionRequest struct {
	IsHistory            *int32  `json:"IsHistory,omitempty" xml:"IsHistory,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	Region               *string `json:"Region,omitempty" xml:"Region,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
	TaskType             *string `json:"TaskType,omitempty" xml:"TaskType,omitempty"`
}

func (s DescribePendingMaintenanceActionRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribePendingMaintenanceActionRequest) GoString() string {
	return s.String()
}

func (s *DescribePendingMaintenanceActionRequest) SetIsHistory(v int32) *DescribePendingMaintenanceActionRequest {
	s.IsHistory = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetOwnerAccount(v string) *DescribePendingMaintenanceActionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetOwnerId(v int64) *DescribePendingMaintenanceActionRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetPageNumber(v int32) *DescribePendingMaintenanceActionRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetPageSize(v int32) *DescribePendingMaintenanceActionRequest {
	s.PageSize = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetRegion(v string) *DescribePendingMaintenanceActionRequest {
	s.Region = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetResourceOwnerAccount(v string) *DescribePendingMaintenanceActionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetResourceOwnerId(v int64) *DescribePendingMaintenanceActionRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetSecurityToken(v string) *DescribePendingMaintenanceActionRequest {
	s.SecurityToken = &v
	return s
}

func (s *DescribePendingMaintenanceActionRequest) SetTaskType(v string) *DescribePendingMaintenanceActionRequest {
	s.TaskType = &v
	return s
}

type DescribePendingMaintenanceActionResponseBody struct {
	Items            []*DescribePendingMaintenanceActionResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Repeated"`
	PageNumber       *int32                                               `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize         *int32                                               `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	RequestId        *string                                              `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *int32                                               `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribePendingMaintenanceActionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribePendingMaintenanceActionResponseBody) GoString() string {
	return s.String()
}

func (s *DescribePendingMaintenanceActionResponseBody) SetItems(v []*DescribePendingMaintenanceActionResponseBodyItems) *DescribePendingMaintenanceActionResponseBody {
	s.Items = v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBody) SetPageNumber(v int32) *DescribePendingMaintenanceActionResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBody) SetPageSize(v int32) *DescribePendingMaintenanceActionResponseBody {
	s.PageSize = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBody) SetRequestId(v string) *DescribePendingMaintenanceActionResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBody) SetTotalRecordCount(v int32) *DescribePendingMaintenanceActionResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribePendingMaintenanceActionResponseBodyItems struct {
	CreatedTime     *string `json:"CreatedTime,omitempty" xml:"CreatedTime,omitempty"`
	DBClusterId     *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBType          *string `json:"DBType,omitempty" xml:"DBType,omitempty"`
	DBVersion       *string `json:"DBVersion,omitempty" xml:"DBVersion,omitempty"`
	Deadline        *string `json:"Deadline,omitempty" xml:"Deadline,omitempty"`
	Id              *int32  `json:"Id,omitempty" xml:"Id,omitempty"`
	ModifiedTime    *string `json:"ModifiedTime,omitempty" xml:"ModifiedTime,omitempty"`
	PrepareInterval *string `json:"PrepareInterval,omitempty" xml:"PrepareInterval,omitempty"`
	Region          *string `json:"Region,omitempty" xml:"Region,omitempty"`
	ResultInfo      *string `json:"ResultInfo,omitempty" xml:"ResultInfo,omitempty"`
	StartTime       *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
	Status          *int32  `json:"Status,omitempty" xml:"Status,omitempty"`
	SwitchTime      *string `json:"SwitchTime,omitempty" xml:"SwitchTime,omitempty"`
	TaskType        *string `json:"TaskType,omitempty" xml:"TaskType,omitempty"`
}

func (s DescribePendingMaintenanceActionResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribePendingMaintenanceActionResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetCreatedTime(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.CreatedTime = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetDBClusterId(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.DBClusterId = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetDBType(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.DBType = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetDBVersion(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.DBVersion = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetDeadline(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.Deadline = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetId(v int32) *DescribePendingMaintenanceActionResponseBodyItems {
	s.Id = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetModifiedTime(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.ModifiedTime = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetPrepareInterval(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.PrepareInterval = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetRegion(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.Region = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetResultInfo(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.ResultInfo = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetStartTime(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.StartTime = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetStatus(v int32) *DescribePendingMaintenanceActionResponseBodyItems {
	s.Status = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetSwitchTime(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.SwitchTime = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponseBodyItems) SetTaskType(v string) *DescribePendingMaintenanceActionResponseBodyItems {
	s.TaskType = &v
	return s
}

type DescribePendingMaintenanceActionResponse struct {
	Headers    map[string]*string                            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribePendingMaintenanceActionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribePendingMaintenanceActionResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribePendingMaintenanceActionResponse) GoString() string {
	return s.String()
}

func (s *DescribePendingMaintenanceActionResponse) SetHeaders(v map[string]*string) *DescribePendingMaintenanceActionResponse {
	s.Headers = v
	return s
}

func (s *DescribePendingMaintenanceActionResponse) SetStatusCode(v int32) *DescribePendingMaintenanceActionResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribePendingMaintenanceActionResponse) SetBody(v *DescribePendingMaintenanceActionResponseBody) *DescribePendingMaintenanceActionResponse {
	s.Body = v
	return s
}

type DescribePendingMaintenanceActionsRequest struct {
	IsHistory            *int32  `json:"IsHistory,omitempty" xml:"IsHistory,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s DescribePendingMaintenanceActionsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribePendingMaintenanceActionsRequest) GoString() string {
	return s.String()
}

func (s *DescribePendingMaintenanceActionsRequest) SetIsHistory(v int32) *DescribePendingMaintenanceActionsRequest {
	s.IsHistory = &v
	return s
}

func (s *DescribePendingMaintenanceActionsRequest) SetOwnerAccount(v string) *DescribePendingMaintenanceActionsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribePendingMaintenanceActionsRequest) SetOwnerId(v int64) *DescribePendingMaintenanceActionsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribePendingMaintenanceActionsRequest) SetRegionId(v string) *DescribePendingMaintenanceActionsRequest {
	s.RegionId = &v
	return s
}

func (s *DescribePendingMaintenanceActionsRequest) SetResourceOwnerAccount(v string) *DescribePendingMaintenanceActionsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribePendingMaintenanceActionsRequest) SetResourceOwnerId(v int64) *DescribePendingMaintenanceActionsRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribePendingMaintenanceActionsRequest) SetSecurityToken(v string) *DescribePendingMaintenanceActionsRequest {
	s.SecurityToken = &v
	return s
}

type DescribePendingMaintenanceActionsResponseBody struct {
	RequestId *string                                                  `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TypeList  []*DescribePendingMaintenanceActionsResponseBodyTypeList `json:"TypeList,omitempty" xml:"TypeList,omitempty" type:"Repeated"`
}

func (s DescribePendingMaintenanceActionsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribePendingMaintenanceActionsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribePendingMaintenanceActionsResponseBody) SetRequestId(v string) *DescribePendingMaintenanceActionsResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribePendingMaintenanceActionsResponseBody) SetTypeList(v []*DescribePendingMaintenanceActionsResponseBodyTypeList) *DescribePendingMaintenanceActionsResponseBody {
	s.TypeList = v
	return s
}

type DescribePendingMaintenanceActionsResponseBodyTypeList struct {
	Count    *int32  `json:"Count,omitempty" xml:"Count,omitempty"`
	TaskType *string `json:"TaskType,omitempty" xml:"TaskType,omitempty"`
}

func (s DescribePendingMaintenanceActionsResponseBodyTypeList) String() string {
	return tea.Prettify(s)
}

func (s DescribePendingMaintenanceActionsResponseBodyTypeList) GoString() string {
	return s.String()
}

func (s *DescribePendingMaintenanceActionsResponseBodyTypeList) SetCount(v int32) *DescribePendingMaintenanceActionsResponseBodyTypeList {
	s.Count = &v
	return s
}

func (s *DescribePendingMaintenanceActionsResponseBodyTypeList) SetTaskType(v string) *DescribePendingMaintenanceActionsResponseBodyTypeList {
	s.TaskType = &v
	return s
}

type DescribePendingMaintenanceActionsResponse struct {
	Headers    map[string]*string                             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribePendingMaintenanceActionsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribePendingMaintenanceActionsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribePendingMaintenanceActionsResponse) GoString() string {
	return s.String()
}

func (s *DescribePendingMaintenanceActionsResponse) SetHeaders(v map[string]*string) *DescribePendingMaintenanceActionsResponse {
	s.Headers = v
	return s
}

func (s *DescribePendingMaintenanceActionsResponse) SetStatusCode(v int32) *DescribePendingMaintenanceActionsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribePendingMaintenanceActionsResponse) SetBody(v *DescribePendingMaintenanceActionsResponseBody) *DescribePendingMaintenanceActionsResponse {
	s.Body = v
	return s
}

type DescribePolarSQLCollectorPolicyRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribePolarSQLCollectorPolicyRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribePolarSQLCollectorPolicyRequest) GoString() string {
	return s.String()
}

func (s *DescribePolarSQLCollectorPolicyRequest) SetDBClusterId(v string) *DescribePolarSQLCollectorPolicyRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribePolarSQLCollectorPolicyRequest) SetOwnerAccount(v string) *DescribePolarSQLCollectorPolicyRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribePolarSQLCollectorPolicyRequest) SetOwnerId(v int64) *DescribePolarSQLCollectorPolicyRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribePolarSQLCollectorPolicyRequest) SetResourceOwnerAccount(v string) *DescribePolarSQLCollectorPolicyRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribePolarSQLCollectorPolicyRequest) SetResourceOwnerId(v int64) *DescribePolarSQLCollectorPolicyRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribePolarSQLCollectorPolicyResponseBody struct {
	DBClusterId        *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	RequestId          *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	SQLCollectorStatus *string `json:"SQLCollectorStatus,omitempty" xml:"SQLCollectorStatus,omitempty"`
}

func (s DescribePolarSQLCollectorPolicyResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribePolarSQLCollectorPolicyResponseBody) GoString() string {
	return s.String()
}

func (s *DescribePolarSQLCollectorPolicyResponseBody) SetDBClusterId(v string) *DescribePolarSQLCollectorPolicyResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribePolarSQLCollectorPolicyResponseBody) SetRequestId(v string) *DescribePolarSQLCollectorPolicyResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribePolarSQLCollectorPolicyResponseBody) SetSQLCollectorStatus(v string) *DescribePolarSQLCollectorPolicyResponseBody {
	s.SQLCollectorStatus = &v
	return s
}

type DescribePolarSQLCollectorPolicyResponse struct {
	Headers    map[string]*string                           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribePolarSQLCollectorPolicyResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribePolarSQLCollectorPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribePolarSQLCollectorPolicyResponse) GoString() string {
	return s.String()
}

func (s *DescribePolarSQLCollectorPolicyResponse) SetHeaders(v map[string]*string) *DescribePolarSQLCollectorPolicyResponse {
	s.Headers = v
	return s
}

func (s *DescribePolarSQLCollectorPolicyResponse) SetStatusCode(v int32) *DescribePolarSQLCollectorPolicyResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribePolarSQLCollectorPolicyResponse) SetBody(v *DescribePolarSQLCollectorPolicyResponseBody) *DescribePolarSQLCollectorPolicyResponse {
	s.Body = v
	return s
}

type DescribeRegionsRequest struct {
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeRegionsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsRequest) GoString() string {
	return s.String()
}

func (s *DescribeRegionsRequest) SetOwnerAccount(v string) *DescribeRegionsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeRegionsRequest) SetOwnerId(v int64) *DescribeRegionsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeRegionsRequest) SetResourceOwnerAccount(v string) *DescribeRegionsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeRegionsRequest) SetResourceOwnerId(v int64) *DescribeRegionsRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeRegionsResponseBody struct {
	Regions   *DescribeRegionsResponseBodyRegions `json:"Regions,omitempty" xml:"Regions,omitempty" type:"Struct"`
	RequestId *string                             `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s DescribeRegionsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponseBody) SetRegions(v *DescribeRegionsResponseBodyRegions) *DescribeRegionsResponseBody {
	s.Regions = v
	return s
}

func (s *DescribeRegionsResponseBody) SetRequestId(v string) *DescribeRegionsResponseBody {
	s.RequestId = &v
	return s
}

type DescribeRegionsResponseBodyRegions struct {
	Region []*DescribeRegionsResponseBodyRegionsRegion `json:"Region,omitempty" xml:"Region,omitempty" type:"Repeated"`
}

func (s DescribeRegionsResponseBodyRegions) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponseBodyRegions) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponseBodyRegions) SetRegion(v []*DescribeRegionsResponseBodyRegionsRegion) *DescribeRegionsResponseBodyRegions {
	s.Region = v
	return s
}

type DescribeRegionsResponseBodyRegionsRegion struct {
	RegionId *string                                        `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	Zones    *DescribeRegionsResponseBodyRegionsRegionZones `json:"Zones,omitempty" xml:"Zones,omitempty" type:"Struct"`
}

func (s DescribeRegionsResponseBodyRegionsRegion) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponseBodyRegionsRegion) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponseBodyRegionsRegion) SetRegionId(v string) *DescribeRegionsResponseBodyRegionsRegion {
	s.RegionId = &v
	return s
}

func (s *DescribeRegionsResponseBodyRegionsRegion) SetZones(v *DescribeRegionsResponseBodyRegionsRegionZones) *DescribeRegionsResponseBodyRegionsRegion {
	s.Zones = v
	return s
}

type DescribeRegionsResponseBodyRegionsRegionZones struct {
	Zone []*DescribeRegionsResponseBodyRegionsRegionZonesZone `json:"Zone,omitempty" xml:"Zone,omitempty" type:"Repeated"`
}

func (s DescribeRegionsResponseBodyRegionsRegionZones) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponseBodyRegionsRegionZones) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponseBodyRegionsRegionZones) SetZone(v []*DescribeRegionsResponseBodyRegionsRegionZonesZone) *DescribeRegionsResponseBodyRegionsRegionZones {
	s.Zone = v
	return s
}

type DescribeRegionsResponseBodyRegionsRegionZonesZone struct {
	VpcEnabled *bool   `json:"VpcEnabled,omitempty" xml:"VpcEnabled,omitempty"`
	ZoneId     *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s DescribeRegionsResponseBodyRegionsRegionZonesZone) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponseBodyRegionsRegionZonesZone) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponseBodyRegionsRegionZonesZone) SetVpcEnabled(v bool) *DescribeRegionsResponseBodyRegionsRegionZonesZone {
	s.VpcEnabled = &v
	return s
}

func (s *DescribeRegionsResponseBodyRegionsRegionZonesZone) SetZoneId(v string) *DescribeRegionsResponseBodyRegionsRegionZonesZone {
	s.ZoneId = &v
	return s
}

type DescribeRegionsResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeRegionsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeRegionsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponse) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponse) SetHeaders(v map[string]*string) *DescribeRegionsResponse {
	s.Headers = v
	return s
}

func (s *DescribeRegionsResponse) SetStatusCode(v int32) *DescribeRegionsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeRegionsResponse) SetBody(v *DescribeRegionsResponseBody) *DescribeRegionsResponse {
	s.Body = v
	return s
}

type DescribeScheduleTasksRequest struct {
	DBClusterDescription *string `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OrderId              *string `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	Status               *string `json:"Status,omitempty" xml:"Status,omitempty"`
	TaskAction           *string `json:"TaskAction,omitempty" xml:"TaskAction,omitempty"`
}

func (s DescribeScheduleTasksRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeScheduleTasksRequest) GoString() string {
	return s.String()
}

func (s *DescribeScheduleTasksRequest) SetDBClusterDescription(v string) *DescribeScheduleTasksRequest {
	s.DBClusterDescription = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetDBClusterId(v string) *DescribeScheduleTasksRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetOrderId(v string) *DescribeScheduleTasksRequest {
	s.OrderId = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetOwnerAccount(v string) *DescribeScheduleTasksRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetOwnerId(v int64) *DescribeScheduleTasksRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetPageNumber(v int32) *DescribeScheduleTasksRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetPageSize(v int32) *DescribeScheduleTasksRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetPlannedEndTime(v string) *DescribeScheduleTasksRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetPlannedStartTime(v string) *DescribeScheduleTasksRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetRegionId(v string) *DescribeScheduleTasksRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetResourceOwnerAccount(v string) *DescribeScheduleTasksRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetResourceOwnerId(v int64) *DescribeScheduleTasksRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetStatus(v string) *DescribeScheduleTasksRequest {
	s.Status = &v
	return s
}

func (s *DescribeScheduleTasksRequest) SetTaskAction(v string) *DescribeScheduleTasksRequest {
	s.TaskAction = &v
	return s
}

type DescribeScheduleTasksResponseBody struct {
	Data      *DescribeScheduleTasksResponseBodyData `json:"Data,omitempty" xml:"Data,omitempty" type:"Struct"`
	Message   *string                                `json:"Message,omitempty" xml:"Message,omitempty"`
	RequestId *string                                `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	Success   *bool                                  `json:"Success,omitempty" xml:"Success,omitempty"`
}

func (s DescribeScheduleTasksResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeScheduleTasksResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeScheduleTasksResponseBody) SetData(v *DescribeScheduleTasksResponseBodyData) *DescribeScheduleTasksResponseBody {
	s.Data = v
	return s
}

func (s *DescribeScheduleTasksResponseBody) SetMessage(v string) *DescribeScheduleTasksResponseBody {
	s.Message = &v
	return s
}

func (s *DescribeScheduleTasksResponseBody) SetRequestId(v string) *DescribeScheduleTasksResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeScheduleTasksResponseBody) SetSuccess(v bool) *DescribeScheduleTasksResponseBody {
	s.Success = &v
	return s
}

type DescribeScheduleTasksResponseBodyData struct {
	PageNumber       *int32                                             `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize         *int32                                             `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	TimerInfos       []*DescribeScheduleTasksResponseBodyDataTimerInfos `json:"TimerInfos,omitempty" xml:"TimerInfos,omitempty" type:"Repeated"`
	TotalRecordCount *int32                                             `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeScheduleTasksResponseBodyData) String() string {
	return tea.Prettify(s)
}

func (s DescribeScheduleTasksResponseBodyData) GoString() string {
	return s.String()
}

func (s *DescribeScheduleTasksResponseBodyData) SetPageNumber(v int32) *DescribeScheduleTasksResponseBodyData {
	s.PageNumber = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyData) SetPageSize(v int32) *DescribeScheduleTasksResponseBodyData {
	s.PageSize = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyData) SetTimerInfos(v []*DescribeScheduleTasksResponseBodyDataTimerInfos) *DescribeScheduleTasksResponseBodyData {
	s.TimerInfos = v
	return s
}

func (s *DescribeScheduleTasksResponseBodyData) SetTotalRecordCount(v int32) *DescribeScheduleTasksResponseBodyData {
	s.TotalRecordCount = &v
	return s
}

type DescribeScheduleTasksResponseBodyDataTimerInfos struct {
	Action               *string `json:"Action,omitempty" xml:"Action,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DbClusterDescription *string `json:"DbClusterDescription,omitempty" xml:"DbClusterDescription,omitempty"`
	DbClusterStatus      *string `json:"DbClusterStatus,omitempty" xml:"DbClusterStatus,omitempty"`
	OrderId              *string `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	PlannedTime          *string `json:"PlannedTime,omitempty" xml:"PlannedTime,omitempty"`
	Region               *string `json:"Region,omitempty" xml:"Region,omitempty"`
	Status               *string `json:"Status,omitempty" xml:"Status,omitempty"`
	TaskCancel           *bool   `json:"TaskCancel,omitempty" xml:"TaskCancel,omitempty"`
	TaskId               *string `json:"TaskId,omitempty" xml:"TaskId,omitempty"`
}

func (s DescribeScheduleTasksResponseBodyDataTimerInfos) String() string {
	return tea.Prettify(s)
}

func (s DescribeScheduleTasksResponseBodyDataTimerInfos) GoString() string {
	return s.String()
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetAction(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.Action = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetDBClusterId(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.DBClusterId = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetDbClusterDescription(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.DbClusterDescription = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetDbClusterStatus(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.DbClusterStatus = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetOrderId(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.OrderId = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetPlannedEndTime(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.PlannedEndTime = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetPlannedStartTime(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.PlannedStartTime = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetPlannedTime(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.PlannedTime = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetRegion(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.Region = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetStatus(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.Status = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetTaskCancel(v bool) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.TaskCancel = &v
	return s
}

func (s *DescribeScheduleTasksResponseBodyDataTimerInfos) SetTaskId(v string) *DescribeScheduleTasksResponseBodyDataTimerInfos {
	s.TaskId = &v
	return s
}

type DescribeScheduleTasksResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeScheduleTasksResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeScheduleTasksResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeScheduleTasksResponse) GoString() string {
	return s.String()
}

func (s *DescribeScheduleTasksResponse) SetHeaders(v map[string]*string) *DescribeScheduleTasksResponse {
	s.Headers = v
	return s
}

func (s *DescribeScheduleTasksResponse) SetStatusCode(v int32) *DescribeScheduleTasksResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeScheduleTasksResponse) SetBody(v *DescribeScheduleTasksResponseBody) *DescribeScheduleTasksResponse {
	s.Body = v
	return s
}

type DescribeSlowLogRecordsRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	EndTime              *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SQLHASH              *string `json:"SQLHASH,omitempty" xml:"SQLHASH,omitempty"`
	StartTime            *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeSlowLogRecordsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogRecordsRequest) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogRecordsRequest) SetDBClusterId(v string) *DescribeSlowLogRecordsRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetDBName(v string) *DescribeSlowLogRecordsRequest {
	s.DBName = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetEndTime(v string) *DescribeSlowLogRecordsRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetOwnerAccount(v string) *DescribeSlowLogRecordsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetOwnerId(v int64) *DescribeSlowLogRecordsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetPageNumber(v int32) *DescribeSlowLogRecordsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetPageSize(v int32) *DescribeSlowLogRecordsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetRegionId(v string) *DescribeSlowLogRecordsRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetResourceOwnerAccount(v string) *DescribeSlowLogRecordsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetResourceOwnerId(v int64) *DescribeSlowLogRecordsRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetSQLHASH(v string) *DescribeSlowLogRecordsRequest {
	s.SQLHASH = &v
	return s
}

func (s *DescribeSlowLogRecordsRequest) SetStartTime(v string) *DescribeSlowLogRecordsRequest {
	s.StartTime = &v
	return s
}

type DescribeSlowLogRecordsResponseBody struct {
	DBClusterId      *string                                  `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	Engine           *string                                  `json:"Engine,omitempty" xml:"Engine,omitempty"`
	Items            *DescribeSlowLogRecordsResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	PageNumber       *int32                                   `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *int32                                   `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                                  `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *int32                                   `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeSlowLogRecordsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogRecordsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogRecordsResponseBody) SetDBClusterId(v string) *DescribeSlowLogRecordsResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBody) SetEngine(v string) *DescribeSlowLogRecordsResponseBody {
	s.Engine = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBody) SetItems(v *DescribeSlowLogRecordsResponseBodyItems) *DescribeSlowLogRecordsResponseBody {
	s.Items = v
	return s
}

func (s *DescribeSlowLogRecordsResponseBody) SetPageNumber(v int32) *DescribeSlowLogRecordsResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBody) SetPageRecordCount(v int32) *DescribeSlowLogRecordsResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBody) SetRequestId(v string) *DescribeSlowLogRecordsResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBody) SetTotalRecordCount(v int32) *DescribeSlowLogRecordsResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeSlowLogRecordsResponseBodyItems struct {
	SQLSlowRecord []*DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord `json:"SQLSlowRecord,omitempty" xml:"SQLSlowRecord,omitempty" type:"Repeated"`
}

func (s DescribeSlowLogRecordsResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogRecordsResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogRecordsResponseBodyItems) SetSQLSlowRecord(v []*DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) *DescribeSlowLogRecordsResponseBodyItems {
	s.SQLSlowRecord = v
	return s
}

type DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord struct {
	DBName             *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	DBNodeId           *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	ExecutionStartTime *string `json:"ExecutionStartTime,omitempty" xml:"ExecutionStartTime,omitempty"`
	HostAddress        *string `json:"HostAddress,omitempty" xml:"HostAddress,omitempty"`
	LockTimes          *int64  `json:"LockTimes,omitempty" xml:"LockTimes,omitempty"`
	ParseRowCounts     *int64  `json:"ParseRowCounts,omitempty" xml:"ParseRowCounts,omitempty"`
	QueryTimeMS        *int64  `json:"QueryTimeMS,omitempty" xml:"QueryTimeMS,omitempty"`
	QueryTimes         *int64  `json:"QueryTimes,omitempty" xml:"QueryTimes,omitempty"`
	ReturnRowCounts    *int64  `json:"ReturnRowCounts,omitempty" xml:"ReturnRowCounts,omitempty"`
	SQLText            *string `json:"SQLText,omitempty" xml:"SQLText,omitempty"`
}

func (s DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetDBName(v string) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.DBName = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetDBNodeId(v string) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.DBNodeId = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetExecutionStartTime(v string) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.ExecutionStartTime = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetHostAddress(v string) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.HostAddress = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetLockTimes(v int64) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.LockTimes = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetParseRowCounts(v int64) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.ParseRowCounts = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetQueryTimeMS(v int64) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.QueryTimeMS = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetQueryTimes(v int64) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.QueryTimes = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetReturnRowCounts(v int64) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.ReturnRowCounts = &v
	return s
}

func (s *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord) SetSQLText(v string) *DescribeSlowLogRecordsResponseBodyItemsSQLSlowRecord {
	s.SQLText = &v
	return s
}

type DescribeSlowLogRecordsResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeSlowLogRecordsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeSlowLogRecordsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogRecordsResponse) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogRecordsResponse) SetHeaders(v map[string]*string) *DescribeSlowLogRecordsResponse {
	s.Headers = v
	return s
}

func (s *DescribeSlowLogRecordsResponse) SetStatusCode(v int32) *DescribeSlowLogRecordsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeSlowLogRecordsResponse) SetBody(v *DescribeSlowLogRecordsResponseBody) *DescribeSlowLogRecordsResponse {
	s.Body = v
	return s
}

type DescribeSlowLogsRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	EndTime              *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	StartTime            *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s DescribeSlowLogsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogsRequest) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogsRequest) SetDBClusterId(v string) *DescribeSlowLogsRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetDBName(v string) *DescribeSlowLogsRequest {
	s.DBName = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetEndTime(v string) *DescribeSlowLogsRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetOwnerAccount(v string) *DescribeSlowLogsRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetOwnerId(v int64) *DescribeSlowLogsRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetPageNumber(v int32) *DescribeSlowLogsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetPageSize(v int32) *DescribeSlowLogsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetRegionId(v string) *DescribeSlowLogsRequest {
	s.RegionId = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetResourceOwnerAccount(v string) *DescribeSlowLogsRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetResourceOwnerId(v int64) *DescribeSlowLogsRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeSlowLogsRequest) SetStartTime(v string) *DescribeSlowLogsRequest {
	s.StartTime = &v
	return s
}

type DescribeSlowLogsResponseBody struct {
	DBClusterId      *string                            `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EndTime          *string                            `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	Engine           *string                            `json:"Engine,omitempty" xml:"Engine,omitempty"`
	Items            *DescribeSlowLogsResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Struct"`
	PageNumber       *int32                             `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *int32                             `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                            `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	StartTime        *string                            `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
	TotalRecordCount *int32                             `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeSlowLogsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogsResponseBody) SetDBClusterId(v string) *DescribeSlowLogsResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeSlowLogsResponseBody) SetEndTime(v string) *DescribeSlowLogsResponseBody {
	s.EndTime = &v
	return s
}

func (s *DescribeSlowLogsResponseBody) SetEngine(v string) *DescribeSlowLogsResponseBody {
	s.Engine = &v
	return s
}

func (s *DescribeSlowLogsResponseBody) SetItems(v *DescribeSlowLogsResponseBodyItems) *DescribeSlowLogsResponseBody {
	s.Items = v
	return s
}

func (s *DescribeSlowLogsResponseBody) SetPageNumber(v int32) *DescribeSlowLogsResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeSlowLogsResponseBody) SetPageRecordCount(v int32) *DescribeSlowLogsResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeSlowLogsResponseBody) SetRequestId(v string) *DescribeSlowLogsResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeSlowLogsResponseBody) SetStartTime(v string) *DescribeSlowLogsResponseBody {
	s.StartTime = &v
	return s
}

func (s *DescribeSlowLogsResponseBody) SetTotalRecordCount(v int32) *DescribeSlowLogsResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeSlowLogsResponseBodyItems struct {
	SQLSlowLog []*DescribeSlowLogsResponseBodyItemsSQLSlowLog `json:"SQLSlowLog,omitempty" xml:"SQLSlowLog,omitempty" type:"Repeated"`
}

func (s DescribeSlowLogsResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogsResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogsResponseBodyItems) SetSQLSlowLog(v []*DescribeSlowLogsResponseBodyItemsSQLSlowLog) *DescribeSlowLogsResponseBodyItems {
	s.SQLSlowLog = v
	return s
}

type DescribeSlowLogsResponseBodyItemsSQLSlowLog struct {
	CreateTime           *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	DBNodeId             *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	MaxExecutionTime     *int64  `json:"MaxExecutionTime,omitempty" xml:"MaxExecutionTime,omitempty"`
	MaxLockTime          *int64  `json:"MaxLockTime,omitempty" xml:"MaxLockTime,omitempty"`
	ParseMaxRowCount     *int64  `json:"ParseMaxRowCount,omitempty" xml:"ParseMaxRowCount,omitempty"`
	ParseTotalRowCounts  *int64  `json:"ParseTotalRowCounts,omitempty" xml:"ParseTotalRowCounts,omitempty"`
	ReturnMaxRowCount    *int64  `json:"ReturnMaxRowCount,omitempty" xml:"ReturnMaxRowCount,omitempty"`
	ReturnTotalRowCounts *int64  `json:"ReturnTotalRowCounts,omitempty" xml:"ReturnTotalRowCounts,omitempty"`
	SQLHASH              *string `json:"SQLHASH,omitempty" xml:"SQLHASH,omitempty"`
	SQLText              *string `json:"SQLText,omitempty" xml:"SQLText,omitempty"`
	TotalExecutionCounts *int64  `json:"TotalExecutionCounts,omitempty" xml:"TotalExecutionCounts,omitempty"`
	TotalExecutionTimes  *int64  `json:"TotalExecutionTimes,omitempty" xml:"TotalExecutionTimes,omitempty"`
	TotalLockTimes       *int64  `json:"TotalLockTimes,omitempty" xml:"TotalLockTimes,omitempty"`
}

func (s DescribeSlowLogsResponseBodyItemsSQLSlowLog) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogsResponseBodyItemsSQLSlowLog) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetCreateTime(v string) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.CreateTime = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetDBName(v string) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.DBName = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetDBNodeId(v string) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.DBNodeId = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetMaxExecutionTime(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.MaxExecutionTime = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetMaxLockTime(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.MaxLockTime = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetParseMaxRowCount(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.ParseMaxRowCount = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetParseTotalRowCounts(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.ParseTotalRowCounts = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetReturnMaxRowCount(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.ReturnMaxRowCount = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetReturnTotalRowCounts(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.ReturnTotalRowCounts = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetSQLHASH(v string) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.SQLHASH = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetSQLText(v string) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.SQLText = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetTotalExecutionCounts(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.TotalExecutionCounts = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetTotalExecutionTimes(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.TotalExecutionTimes = &v
	return s
}

func (s *DescribeSlowLogsResponseBodyItemsSQLSlowLog) SetTotalLockTimes(v int64) *DescribeSlowLogsResponseBodyItemsSQLSlowLog {
	s.TotalLockTimes = &v
	return s
}

type DescribeSlowLogsResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeSlowLogsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeSlowLogsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeSlowLogsResponse) GoString() string {
	return s.String()
}

func (s *DescribeSlowLogsResponse) SetHeaders(v map[string]*string) *DescribeSlowLogsResponse {
	s.Headers = v
	return s
}

func (s *DescribeSlowLogsResponse) SetStatusCode(v int32) *DescribeSlowLogsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeSlowLogsResponse) SetBody(v *DescribeSlowLogsResponseBody) *DescribeSlowLogsResponse {
	s.Body = v
	return s
}

type DescribeStoragePlanRequest struct {
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s DescribeStoragePlanRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeStoragePlanRequest) GoString() string {
	return s.String()
}

func (s *DescribeStoragePlanRequest) SetOwnerAccount(v string) *DescribeStoragePlanRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeStoragePlanRequest) SetOwnerId(v int64) *DescribeStoragePlanRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeStoragePlanRequest) SetPageNumber(v int32) *DescribeStoragePlanRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeStoragePlanRequest) SetPageSize(v int32) *DescribeStoragePlanRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeStoragePlanRequest) SetResourceOwnerAccount(v string) *DescribeStoragePlanRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeStoragePlanRequest) SetResourceOwnerId(v int64) *DescribeStoragePlanRequest {
	s.ResourceOwnerId = &v
	return s
}

type DescribeStoragePlanResponseBody struct {
	Items            []*DescribeStoragePlanResponseBodyItems `json:"Items,omitempty" xml:"Items,omitempty" type:"Repeated"`
	PageNumber       *int64                                  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize         *int64                                  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	RequestId        *string                                 `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TotalRecordCount *int64                                  `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeStoragePlanResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeStoragePlanResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeStoragePlanResponseBody) SetItems(v []*DescribeStoragePlanResponseBodyItems) *DescribeStoragePlanResponseBody {
	s.Items = v
	return s
}

func (s *DescribeStoragePlanResponseBody) SetPageNumber(v int64) *DescribeStoragePlanResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeStoragePlanResponseBody) SetPageSize(v int64) *DescribeStoragePlanResponseBody {
	s.PageSize = &v
	return s
}

func (s *DescribeStoragePlanResponseBody) SetRequestId(v string) *DescribeStoragePlanResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeStoragePlanResponseBody) SetTotalRecordCount(v int64) *DescribeStoragePlanResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeStoragePlanResponseBodyItems struct {
	AliUid                  *string `json:"AliUid,omitempty" xml:"AliUid,omitempty"`
	CommodityCode           *string `json:"CommodityCode,omitempty" xml:"CommodityCode,omitempty"`
	EndTimes                *string `json:"EndTimes,omitempty" xml:"EndTimes,omitempty"`
	InitCapaCityViewUnit    *string `json:"InitCapaCityViewUnit,omitempty" xml:"InitCapaCityViewUnit,omitempty"`
	InitCapacityViewValue   *string `json:"InitCapacityViewValue,omitempty" xml:"InitCapacityViewValue,omitempty"`
	InstanceId              *string `json:"InstanceId,omitempty" xml:"InstanceId,omitempty"`
	PeriodCapaCityViewUnit  *string `json:"PeriodCapaCityViewUnit,omitempty" xml:"PeriodCapaCityViewUnit,omitempty"`
	PeriodCapacityViewValue *string `json:"PeriodCapacityViewValue,omitempty" xml:"PeriodCapacityViewValue,omitempty"`
	PeriodTime              *string `json:"PeriodTime,omitempty" xml:"PeriodTime,omitempty"`
	ProdCode                *string `json:"ProdCode,omitempty" xml:"ProdCode,omitempty"`
	PurchaseTimes           *string `json:"PurchaseTimes,omitempty" xml:"PurchaseTimes,omitempty"`
	StartTimes              *string `json:"StartTimes,omitempty" xml:"StartTimes,omitempty"`
	Status                  *string `json:"Status,omitempty" xml:"Status,omitempty"`
	StorageType             *string `json:"StorageType,omitempty" xml:"StorageType,omitempty"`
	TemplateName            *string `json:"TemplateName,omitempty" xml:"TemplateName,omitempty"`
}

func (s DescribeStoragePlanResponseBodyItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeStoragePlanResponseBodyItems) GoString() string {
	return s.String()
}

func (s *DescribeStoragePlanResponseBodyItems) SetAliUid(v string) *DescribeStoragePlanResponseBodyItems {
	s.AliUid = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetCommodityCode(v string) *DescribeStoragePlanResponseBodyItems {
	s.CommodityCode = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetEndTimes(v string) *DescribeStoragePlanResponseBodyItems {
	s.EndTimes = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetInitCapaCityViewUnit(v string) *DescribeStoragePlanResponseBodyItems {
	s.InitCapaCityViewUnit = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetInitCapacityViewValue(v string) *DescribeStoragePlanResponseBodyItems {
	s.InitCapacityViewValue = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetInstanceId(v string) *DescribeStoragePlanResponseBodyItems {
	s.InstanceId = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetPeriodCapaCityViewUnit(v string) *DescribeStoragePlanResponseBodyItems {
	s.PeriodCapaCityViewUnit = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetPeriodCapacityViewValue(v string) *DescribeStoragePlanResponseBodyItems {
	s.PeriodCapacityViewValue = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetPeriodTime(v string) *DescribeStoragePlanResponseBodyItems {
	s.PeriodTime = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetProdCode(v string) *DescribeStoragePlanResponseBodyItems {
	s.ProdCode = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetPurchaseTimes(v string) *DescribeStoragePlanResponseBodyItems {
	s.PurchaseTimes = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetStartTimes(v string) *DescribeStoragePlanResponseBodyItems {
	s.StartTimes = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetStatus(v string) *DescribeStoragePlanResponseBodyItems {
	s.Status = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetStorageType(v string) *DescribeStoragePlanResponseBodyItems {
	s.StorageType = &v
	return s
}

func (s *DescribeStoragePlanResponseBodyItems) SetTemplateName(v string) *DescribeStoragePlanResponseBodyItems {
	s.TemplateName = &v
	return s
}

type DescribeStoragePlanResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeStoragePlanResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeStoragePlanResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeStoragePlanResponse) GoString() string {
	return s.String()
}

func (s *DescribeStoragePlanResponse) SetHeaders(v map[string]*string) *DescribeStoragePlanResponse {
	s.Headers = v
	return s
}

func (s *DescribeStoragePlanResponse) SetStatusCode(v int32) *DescribeStoragePlanResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeStoragePlanResponse) SetBody(v *DescribeStoragePlanResponseBody) *DescribeStoragePlanResponse {
	s.Body = v
	return s
}

type DescribeTasksRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeId             *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	EndTime              *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PageNumber           *int32  `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageSize             *int32  `json:"PageSize,omitempty" xml:"PageSize,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	StartTime            *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
	Status               *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s DescribeTasksRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeTasksRequest) GoString() string {
	return s.String()
}

func (s *DescribeTasksRequest) SetDBClusterId(v string) *DescribeTasksRequest {
	s.DBClusterId = &v
	return s
}

func (s *DescribeTasksRequest) SetDBNodeId(v string) *DescribeTasksRequest {
	s.DBNodeId = &v
	return s
}

func (s *DescribeTasksRequest) SetEndTime(v string) *DescribeTasksRequest {
	s.EndTime = &v
	return s
}

func (s *DescribeTasksRequest) SetOwnerAccount(v string) *DescribeTasksRequest {
	s.OwnerAccount = &v
	return s
}

func (s *DescribeTasksRequest) SetOwnerId(v int64) *DescribeTasksRequest {
	s.OwnerId = &v
	return s
}

func (s *DescribeTasksRequest) SetPageNumber(v int32) *DescribeTasksRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeTasksRequest) SetPageSize(v int32) *DescribeTasksRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeTasksRequest) SetResourceOwnerAccount(v string) *DescribeTasksRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *DescribeTasksRequest) SetResourceOwnerId(v int64) *DescribeTasksRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *DescribeTasksRequest) SetStartTime(v string) *DescribeTasksRequest {
	s.StartTime = &v
	return s
}

func (s *DescribeTasksRequest) SetStatus(v string) *DescribeTasksRequest {
	s.Status = &v
	return s
}

type DescribeTasksResponseBody struct {
	DBClusterId      *string                         `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EndTime          *string                         `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	PageNumber       *int32                          `json:"PageNumber,omitempty" xml:"PageNumber,omitempty"`
	PageRecordCount  *int32                          `json:"PageRecordCount,omitempty" xml:"PageRecordCount,omitempty"`
	RequestId        *string                         `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	StartTime        *string                         `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
	Tasks            *DescribeTasksResponseBodyTasks `json:"Tasks,omitempty" xml:"Tasks,omitempty" type:"Struct"`
	TotalRecordCount *int32                          `json:"TotalRecordCount,omitempty" xml:"TotalRecordCount,omitempty"`
}

func (s DescribeTasksResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeTasksResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeTasksResponseBody) SetDBClusterId(v string) *DescribeTasksResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *DescribeTasksResponseBody) SetEndTime(v string) *DescribeTasksResponseBody {
	s.EndTime = &v
	return s
}

func (s *DescribeTasksResponseBody) SetPageNumber(v int32) *DescribeTasksResponseBody {
	s.PageNumber = &v
	return s
}

func (s *DescribeTasksResponseBody) SetPageRecordCount(v int32) *DescribeTasksResponseBody {
	s.PageRecordCount = &v
	return s
}

func (s *DescribeTasksResponseBody) SetRequestId(v string) *DescribeTasksResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeTasksResponseBody) SetStartTime(v string) *DescribeTasksResponseBody {
	s.StartTime = &v
	return s
}

func (s *DescribeTasksResponseBody) SetTasks(v *DescribeTasksResponseBodyTasks) *DescribeTasksResponseBody {
	s.Tasks = v
	return s
}

func (s *DescribeTasksResponseBody) SetTotalRecordCount(v int32) *DescribeTasksResponseBody {
	s.TotalRecordCount = &v
	return s
}

type DescribeTasksResponseBodyTasks struct {
	Task []*DescribeTasksResponseBodyTasksTask `json:"Task,omitempty" xml:"Task,omitempty" type:"Repeated"`
}

func (s DescribeTasksResponseBodyTasks) String() string {
	return tea.Prettify(s)
}

func (s DescribeTasksResponseBodyTasks) GoString() string {
	return s.String()
}

func (s *DescribeTasksResponseBodyTasks) SetTask(v []*DescribeTasksResponseBodyTasksTask) *DescribeTasksResponseBodyTasks {
	s.Task = v
	return s
}

type DescribeTasksResponseBodyTasksTask struct {
	BeginTime          *string `json:"BeginTime,omitempty" xml:"BeginTime,omitempty"`
	CurrentStepName    *string `json:"CurrentStepName,omitempty" xml:"CurrentStepName,omitempty"`
	DBName             *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	ExpectedFinishTime *string `json:"ExpectedFinishTime,omitempty" xml:"ExpectedFinishTime,omitempty"`
	FinishTime         *string `json:"FinishTime,omitempty" xml:"FinishTime,omitempty"`
	Progress           *int32  `json:"Progress,omitempty" xml:"Progress,omitempty"`
	ProgressInfo       *string `json:"ProgressInfo,omitempty" xml:"ProgressInfo,omitempty"`
	Remain             *int32  `json:"Remain,omitempty" xml:"Remain,omitempty"`
	StepProgressInfo   *string `json:"StepProgressInfo,omitempty" xml:"StepProgressInfo,omitempty"`
	StepsInfo          *string `json:"StepsInfo,omitempty" xml:"StepsInfo,omitempty"`
	TaskAction         *string `json:"TaskAction,omitempty" xml:"TaskAction,omitempty"`
	TaskErrorCode      *string `json:"TaskErrorCode,omitempty" xml:"TaskErrorCode,omitempty"`
	TaskErrorMessage   *string `json:"TaskErrorMessage,omitempty" xml:"TaskErrorMessage,omitempty"`
	TaskId             *string `json:"TaskId,omitempty" xml:"TaskId,omitempty"`
}

func (s DescribeTasksResponseBodyTasksTask) String() string {
	return tea.Prettify(s)
}

func (s DescribeTasksResponseBodyTasksTask) GoString() string {
	return s.String()
}

func (s *DescribeTasksResponseBodyTasksTask) SetBeginTime(v string) *DescribeTasksResponseBodyTasksTask {
	s.BeginTime = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetCurrentStepName(v string) *DescribeTasksResponseBodyTasksTask {
	s.CurrentStepName = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetDBName(v string) *DescribeTasksResponseBodyTasksTask {
	s.DBName = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetExpectedFinishTime(v string) *DescribeTasksResponseBodyTasksTask {
	s.ExpectedFinishTime = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetFinishTime(v string) *DescribeTasksResponseBodyTasksTask {
	s.FinishTime = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetProgress(v int32) *DescribeTasksResponseBodyTasksTask {
	s.Progress = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetProgressInfo(v string) *DescribeTasksResponseBodyTasksTask {
	s.ProgressInfo = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetRemain(v int32) *DescribeTasksResponseBodyTasksTask {
	s.Remain = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetStepProgressInfo(v string) *DescribeTasksResponseBodyTasksTask {
	s.StepProgressInfo = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetStepsInfo(v string) *DescribeTasksResponseBodyTasksTask {
	s.StepsInfo = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetTaskAction(v string) *DescribeTasksResponseBodyTasksTask {
	s.TaskAction = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetTaskErrorCode(v string) *DescribeTasksResponseBodyTasksTask {
	s.TaskErrorCode = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetTaskErrorMessage(v string) *DescribeTasksResponseBodyTasksTask {
	s.TaskErrorMessage = &v
	return s
}

func (s *DescribeTasksResponseBodyTasksTask) SetTaskId(v string) *DescribeTasksResponseBodyTasksTask {
	s.TaskId = &v
	return s
}

type DescribeTasksResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeTasksResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeTasksResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeTasksResponse) GoString() string {
	return s.String()
}

func (s *DescribeTasksResponse) SetHeaders(v map[string]*string) *DescribeTasksResponse {
	s.Headers = v
	return s
}

func (s *DescribeTasksResponse) SetStatusCode(v int32) *DescribeTasksResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeTasksResponse) SetBody(v *DescribeTasksResponseBody) *DescribeTasksResponse {
	s.Body = v
	return s
}

type EnableFirewallRulesRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	Enable               *bool   `json:"Enable,omitempty" xml:"Enable,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	RuleNameList         *string `json:"RuleNameList,omitempty" xml:"RuleNameList,omitempty"`
}

func (s EnableFirewallRulesRequest) String() string {
	return tea.Prettify(s)
}

func (s EnableFirewallRulesRequest) GoString() string {
	return s.String()
}

func (s *EnableFirewallRulesRequest) SetDBClusterId(v string) *EnableFirewallRulesRequest {
	s.DBClusterId = &v
	return s
}

func (s *EnableFirewallRulesRequest) SetEnable(v bool) *EnableFirewallRulesRequest {
	s.Enable = &v
	return s
}

func (s *EnableFirewallRulesRequest) SetOwnerAccount(v string) *EnableFirewallRulesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *EnableFirewallRulesRequest) SetOwnerId(v int64) *EnableFirewallRulesRequest {
	s.OwnerId = &v
	return s
}

func (s *EnableFirewallRulesRequest) SetResourceOwnerAccount(v string) *EnableFirewallRulesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *EnableFirewallRulesRequest) SetResourceOwnerId(v int64) *EnableFirewallRulesRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *EnableFirewallRulesRequest) SetRuleNameList(v string) *EnableFirewallRulesRequest {
	s.RuleNameList = &v
	return s
}

type EnableFirewallRulesResponseBody struct {
	Message   *string `json:"Message,omitempty" xml:"Message,omitempty"`
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	Success   *bool   `json:"Success,omitempty" xml:"Success,omitempty"`
}

func (s EnableFirewallRulesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s EnableFirewallRulesResponseBody) GoString() string {
	return s.String()
}

func (s *EnableFirewallRulesResponseBody) SetMessage(v string) *EnableFirewallRulesResponseBody {
	s.Message = &v
	return s
}

func (s *EnableFirewallRulesResponseBody) SetRequestId(v string) *EnableFirewallRulesResponseBody {
	s.RequestId = &v
	return s
}

func (s *EnableFirewallRulesResponseBody) SetSuccess(v bool) *EnableFirewallRulesResponseBody {
	s.Success = &v
	return s
}

type EnableFirewallRulesResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *EnableFirewallRulesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s EnableFirewallRulesResponse) String() string {
	return tea.Prettify(s)
}

func (s EnableFirewallRulesResponse) GoString() string {
	return s.String()
}

func (s *EnableFirewallRulesResponse) SetHeaders(v map[string]*string) *EnableFirewallRulesResponse {
	s.Headers = v
	return s
}

func (s *EnableFirewallRulesResponse) SetStatusCode(v int32) *EnableFirewallRulesResponse {
	s.StatusCode = &v
	return s
}

func (s *EnableFirewallRulesResponse) SetBody(v *EnableFirewallRulesResponseBody) *EnableFirewallRulesResponse {
	s.Body = v
	return s
}

type FailoverDBClusterRequest struct {
	ClientToken          *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	TargetDBNodeId       *string `json:"TargetDBNodeId,omitempty" xml:"TargetDBNodeId,omitempty"`
}

func (s FailoverDBClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s FailoverDBClusterRequest) GoString() string {
	return s.String()
}

func (s *FailoverDBClusterRequest) SetClientToken(v string) *FailoverDBClusterRequest {
	s.ClientToken = &v
	return s
}

func (s *FailoverDBClusterRequest) SetDBClusterId(v string) *FailoverDBClusterRequest {
	s.DBClusterId = &v
	return s
}

func (s *FailoverDBClusterRequest) SetOwnerAccount(v string) *FailoverDBClusterRequest {
	s.OwnerAccount = &v
	return s
}

func (s *FailoverDBClusterRequest) SetOwnerId(v int64) *FailoverDBClusterRequest {
	s.OwnerId = &v
	return s
}

func (s *FailoverDBClusterRequest) SetResourceOwnerAccount(v string) *FailoverDBClusterRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *FailoverDBClusterRequest) SetResourceOwnerId(v int64) *FailoverDBClusterRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *FailoverDBClusterRequest) SetTargetDBNodeId(v string) *FailoverDBClusterRequest {
	s.TargetDBNodeId = &v
	return s
}

type FailoverDBClusterResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s FailoverDBClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s FailoverDBClusterResponseBody) GoString() string {
	return s.String()
}

func (s *FailoverDBClusterResponseBody) SetRequestId(v string) *FailoverDBClusterResponseBody {
	s.RequestId = &v
	return s
}

type FailoverDBClusterResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *FailoverDBClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s FailoverDBClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s FailoverDBClusterResponse) GoString() string {
	return s.String()
}

func (s *FailoverDBClusterResponse) SetHeaders(v map[string]*string) *FailoverDBClusterResponse {
	s.Headers = v
	return s
}

func (s *FailoverDBClusterResponse) SetStatusCode(v int32) *FailoverDBClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *FailoverDBClusterResponse) SetBody(v *FailoverDBClusterResponseBody) *FailoverDBClusterResponse {
	s.Body = v
	return s
}

type GrantAccountPrivilegeRequest struct {
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	AccountPrivilege     *string `json:"AccountPrivilege,omitempty" xml:"AccountPrivilege,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s GrantAccountPrivilegeRequest) String() string {
	return tea.Prettify(s)
}

func (s GrantAccountPrivilegeRequest) GoString() string {
	return s.String()
}

func (s *GrantAccountPrivilegeRequest) SetAccountName(v string) *GrantAccountPrivilegeRequest {
	s.AccountName = &v
	return s
}

func (s *GrantAccountPrivilegeRequest) SetAccountPrivilege(v string) *GrantAccountPrivilegeRequest {
	s.AccountPrivilege = &v
	return s
}

func (s *GrantAccountPrivilegeRequest) SetDBClusterId(v string) *GrantAccountPrivilegeRequest {
	s.DBClusterId = &v
	return s
}

func (s *GrantAccountPrivilegeRequest) SetDBName(v string) *GrantAccountPrivilegeRequest {
	s.DBName = &v
	return s
}

func (s *GrantAccountPrivilegeRequest) SetOwnerAccount(v string) *GrantAccountPrivilegeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *GrantAccountPrivilegeRequest) SetOwnerId(v int64) *GrantAccountPrivilegeRequest {
	s.OwnerId = &v
	return s
}

func (s *GrantAccountPrivilegeRequest) SetResourceOwnerAccount(v string) *GrantAccountPrivilegeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *GrantAccountPrivilegeRequest) SetResourceOwnerId(v int64) *GrantAccountPrivilegeRequest {
	s.ResourceOwnerId = &v
	return s
}

type GrantAccountPrivilegeResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s GrantAccountPrivilegeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GrantAccountPrivilegeResponseBody) GoString() string {
	return s.String()
}

func (s *GrantAccountPrivilegeResponseBody) SetRequestId(v string) *GrantAccountPrivilegeResponseBody {
	s.RequestId = &v
	return s
}

type GrantAccountPrivilegeResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GrantAccountPrivilegeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GrantAccountPrivilegeResponse) String() string {
	return tea.Prettify(s)
}

func (s GrantAccountPrivilegeResponse) GoString() string {
	return s.String()
}

func (s *GrantAccountPrivilegeResponse) SetHeaders(v map[string]*string) *GrantAccountPrivilegeResponse {
	s.Headers = v
	return s
}

func (s *GrantAccountPrivilegeResponse) SetStatusCode(v int32) *GrantAccountPrivilegeResponse {
	s.StatusCode = &v
	return s
}

func (s *GrantAccountPrivilegeResponse) SetBody(v *GrantAccountPrivilegeResponseBody) *GrantAccountPrivilegeResponse {
	s.Body = v
	return s
}

type ListTagResourcesRequest struct {
	NextToken            *string                       `json:"NextToken,omitempty" xml:"NextToken,omitempty"`
	OwnerAccount         *string                       `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64                        `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string                       `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceId           []*string                     `json:"ResourceId,omitempty" xml:"ResourceId,omitempty" type:"Repeated"`
	ResourceOwnerAccount *string                       `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64                        `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	ResourceType         *string                       `json:"ResourceType,omitempty" xml:"ResourceType,omitempty"`
	Tag                  []*ListTagResourcesRequestTag `json:"Tag,omitempty" xml:"Tag,omitempty" type:"Repeated"`
}

func (s ListTagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesRequest) GoString() string {
	return s.String()
}

func (s *ListTagResourcesRequest) SetNextToken(v string) *ListTagResourcesRequest {
	s.NextToken = &v
	return s
}

func (s *ListTagResourcesRequest) SetOwnerAccount(v string) *ListTagResourcesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ListTagResourcesRequest) SetOwnerId(v int64) *ListTagResourcesRequest {
	s.OwnerId = &v
	return s
}

func (s *ListTagResourcesRequest) SetRegionId(v string) *ListTagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *ListTagResourcesRequest) SetResourceId(v []*string) *ListTagResourcesRequest {
	s.ResourceId = v
	return s
}

func (s *ListTagResourcesRequest) SetResourceOwnerAccount(v string) *ListTagResourcesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ListTagResourcesRequest) SetResourceOwnerId(v int64) *ListTagResourcesRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ListTagResourcesRequest) SetResourceType(v string) *ListTagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *ListTagResourcesRequest) SetTag(v []*ListTagResourcesRequestTag) *ListTagResourcesRequest {
	s.Tag = v
	return s
}

type ListTagResourcesRequestTag struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s ListTagResourcesRequestTag) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesRequestTag) GoString() string {
	return s.String()
}

func (s *ListTagResourcesRequestTag) SetKey(v string) *ListTagResourcesRequestTag {
	s.Key = &v
	return s
}

func (s *ListTagResourcesRequestTag) SetValue(v string) *ListTagResourcesRequestTag {
	s.Value = &v
	return s
}

type ListTagResourcesResponseBody struct {
	NextToken    *string                                   `json:"NextToken,omitempty" xml:"NextToken,omitempty"`
	RequestId    *string                                   `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TagResources *ListTagResourcesResponseBodyTagResources `json:"TagResources,omitempty" xml:"TagResources,omitempty" type:"Struct"`
}

func (s ListTagResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBody) SetNextToken(v string) *ListTagResourcesResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListTagResourcesResponseBody) SetRequestId(v string) *ListTagResourcesResponseBody {
	s.RequestId = &v
	return s
}

func (s *ListTagResourcesResponseBody) SetTagResources(v *ListTagResourcesResponseBodyTagResources) *ListTagResourcesResponseBody {
	s.TagResources = v
	return s
}

type ListTagResourcesResponseBodyTagResources struct {
	TagResource []*ListTagResourcesResponseBodyTagResourcesTagResource `json:"TagResource,omitempty" xml:"TagResource,omitempty" type:"Repeated"`
}

func (s ListTagResourcesResponseBodyTagResources) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBodyTagResources) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBodyTagResources) SetTagResource(v []*ListTagResourcesResponseBodyTagResourcesTagResource) *ListTagResourcesResponseBodyTagResources {
	s.TagResource = v
	return s
}

type ListTagResourcesResponseBodyTagResourcesTagResource struct {
	ResourceId   *string `json:"ResourceId,omitempty" xml:"ResourceId,omitempty"`
	ResourceType *string `json:"ResourceType,omitempty" xml:"ResourceType,omitempty"`
	TagKey       *string `json:"TagKey,omitempty" xml:"TagKey,omitempty"`
	TagValue     *string `json:"TagValue,omitempty" xml:"TagValue,omitempty"`
}

func (s ListTagResourcesResponseBodyTagResourcesTagResource) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBodyTagResourcesTagResource) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetResourceId(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.ResourceId = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetResourceType(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.ResourceType = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetTagKey(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.TagKey = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetTagValue(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.TagValue = &v
	return s
}

type ListTagResourcesResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListTagResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListTagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponse) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponse) SetHeaders(v map[string]*string) *ListTagResourcesResponse {
	s.Headers = v
	return s
}

func (s *ListTagResourcesResponse) SetStatusCode(v int32) *ListTagResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *ListTagResourcesResponse) SetBody(v *ListTagResourcesResponseBody) *ListTagResourcesResponse {
	s.Body = v
	return s
}

type ModifyAccountDescriptionRequest struct {
	AccountDescription   *string `json:"AccountDescription,omitempty" xml:"AccountDescription,omitempty"`
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyAccountDescriptionRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyAccountDescriptionRequest) GoString() string {
	return s.String()
}

func (s *ModifyAccountDescriptionRequest) SetAccountDescription(v string) *ModifyAccountDescriptionRequest {
	s.AccountDescription = &v
	return s
}

func (s *ModifyAccountDescriptionRequest) SetAccountName(v string) *ModifyAccountDescriptionRequest {
	s.AccountName = &v
	return s
}

func (s *ModifyAccountDescriptionRequest) SetDBClusterId(v string) *ModifyAccountDescriptionRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyAccountDescriptionRequest) SetOwnerAccount(v string) *ModifyAccountDescriptionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyAccountDescriptionRequest) SetOwnerId(v int64) *ModifyAccountDescriptionRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyAccountDescriptionRequest) SetResourceOwnerAccount(v string) *ModifyAccountDescriptionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyAccountDescriptionRequest) SetResourceOwnerId(v int64) *ModifyAccountDescriptionRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyAccountDescriptionResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyAccountDescriptionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyAccountDescriptionResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyAccountDescriptionResponseBody) SetRequestId(v string) *ModifyAccountDescriptionResponseBody {
	s.RequestId = &v
	return s
}

type ModifyAccountDescriptionResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyAccountDescriptionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyAccountDescriptionResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyAccountDescriptionResponse) GoString() string {
	return s.String()
}

func (s *ModifyAccountDescriptionResponse) SetHeaders(v map[string]*string) *ModifyAccountDescriptionResponse {
	s.Headers = v
	return s
}

func (s *ModifyAccountDescriptionResponse) SetStatusCode(v int32) *ModifyAccountDescriptionResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyAccountDescriptionResponse) SetBody(v *ModifyAccountDescriptionResponseBody) *ModifyAccountDescriptionResponse {
	s.Body = v
	return s
}

type ModifyAccountPasswordRequest struct {
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	NewAccountPassword   *string `json:"NewAccountPassword,omitempty" xml:"NewAccountPassword,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyAccountPasswordRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyAccountPasswordRequest) GoString() string {
	return s.String()
}

func (s *ModifyAccountPasswordRequest) SetAccountName(v string) *ModifyAccountPasswordRequest {
	s.AccountName = &v
	return s
}

func (s *ModifyAccountPasswordRequest) SetDBClusterId(v string) *ModifyAccountPasswordRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyAccountPasswordRequest) SetNewAccountPassword(v string) *ModifyAccountPasswordRequest {
	s.NewAccountPassword = &v
	return s
}

func (s *ModifyAccountPasswordRequest) SetOwnerAccount(v string) *ModifyAccountPasswordRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyAccountPasswordRequest) SetOwnerId(v int64) *ModifyAccountPasswordRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyAccountPasswordRequest) SetResourceOwnerAccount(v string) *ModifyAccountPasswordRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyAccountPasswordRequest) SetResourceOwnerId(v int64) *ModifyAccountPasswordRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyAccountPasswordResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyAccountPasswordResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyAccountPasswordResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyAccountPasswordResponseBody) SetRequestId(v string) *ModifyAccountPasswordResponseBody {
	s.RequestId = &v
	return s
}

type ModifyAccountPasswordResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyAccountPasswordResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyAccountPasswordResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyAccountPasswordResponse) GoString() string {
	return s.String()
}

func (s *ModifyAccountPasswordResponse) SetHeaders(v map[string]*string) *ModifyAccountPasswordResponse {
	s.Headers = v
	return s
}

func (s *ModifyAccountPasswordResponse) SetStatusCode(v int32) *ModifyAccountPasswordResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyAccountPasswordResponse) SetBody(v *ModifyAccountPasswordResponseBody) *ModifyAccountPasswordResponse {
	s.Body = v
	return s
}

type ModifyAutoRenewAttributeRequest struct {
	DBClusterIds         *string `json:"DBClusterIds,omitempty" xml:"DBClusterIds,omitempty"`
	Duration             *string `json:"Duration,omitempty" xml:"Duration,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PeriodUnit           *string `json:"PeriodUnit,omitempty" xml:"PeriodUnit,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	RenewalStatus        *string `json:"RenewalStatus,omitempty" xml:"RenewalStatus,omitempty"`
	ResourceGroupId      *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyAutoRenewAttributeRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyAutoRenewAttributeRequest) GoString() string {
	return s.String()
}

func (s *ModifyAutoRenewAttributeRequest) SetDBClusterIds(v string) *ModifyAutoRenewAttributeRequest {
	s.DBClusterIds = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetDuration(v string) *ModifyAutoRenewAttributeRequest {
	s.Duration = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetOwnerAccount(v string) *ModifyAutoRenewAttributeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetOwnerId(v int64) *ModifyAutoRenewAttributeRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetPeriodUnit(v string) *ModifyAutoRenewAttributeRequest {
	s.PeriodUnit = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetRegionId(v string) *ModifyAutoRenewAttributeRequest {
	s.RegionId = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetRenewalStatus(v string) *ModifyAutoRenewAttributeRequest {
	s.RenewalStatus = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetResourceGroupId(v string) *ModifyAutoRenewAttributeRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetResourceOwnerAccount(v string) *ModifyAutoRenewAttributeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyAutoRenewAttributeRequest) SetResourceOwnerId(v int64) *ModifyAutoRenewAttributeRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyAutoRenewAttributeResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyAutoRenewAttributeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyAutoRenewAttributeResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyAutoRenewAttributeResponseBody) SetRequestId(v string) *ModifyAutoRenewAttributeResponseBody {
	s.RequestId = &v
	return s
}

type ModifyAutoRenewAttributeResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyAutoRenewAttributeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyAutoRenewAttributeResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyAutoRenewAttributeResponse) GoString() string {
	return s.String()
}

func (s *ModifyAutoRenewAttributeResponse) SetHeaders(v map[string]*string) *ModifyAutoRenewAttributeResponse {
	s.Headers = v
	return s
}

func (s *ModifyAutoRenewAttributeResponse) SetStatusCode(v int32) *ModifyAutoRenewAttributeResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyAutoRenewAttributeResponse) SetBody(v *ModifyAutoRenewAttributeResponseBody) *ModifyAutoRenewAttributeResponse {
	s.Body = v
	return s
}

type ModifyBackupPolicyRequest struct {
	BackupFrequency                              *string `json:"BackupFrequency,omitempty" xml:"BackupFrequency,omitempty"`
	BackupRetentionPolicyOnClusterDeletion       *string `json:"BackupRetentionPolicyOnClusterDeletion,omitempty" xml:"BackupRetentionPolicyOnClusterDeletion,omitempty"`
	DBClusterId                                  *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DataLevel1BackupFrequency                    *string `json:"DataLevel1BackupFrequency,omitempty" xml:"DataLevel1BackupFrequency,omitempty"`
	DataLevel1BackupPeriod                       *string `json:"DataLevel1BackupPeriod,omitempty" xml:"DataLevel1BackupPeriod,omitempty"`
	DataLevel1BackupRetentionPeriod              *string `json:"DataLevel1BackupRetentionPeriod,omitempty" xml:"DataLevel1BackupRetentionPeriod,omitempty"`
	DataLevel1BackupTime                         *string `json:"DataLevel1BackupTime,omitempty" xml:"DataLevel1BackupTime,omitempty"`
	DataLevel2BackupAnotherRegionRegion          *string `json:"DataLevel2BackupAnotherRegionRegion,omitempty" xml:"DataLevel2BackupAnotherRegionRegion,omitempty"`
	DataLevel2BackupAnotherRegionRetentionPeriod *string `json:"DataLevel2BackupAnotherRegionRetentionPeriod,omitempty" xml:"DataLevel2BackupAnotherRegionRetentionPeriod,omitempty"`
	DataLevel2BackupPeriod                       *string `json:"DataLevel2BackupPeriod,omitempty" xml:"DataLevel2BackupPeriod,omitempty"`
	DataLevel2BackupRetentionPeriod              *string `json:"DataLevel2BackupRetentionPeriod,omitempty" xml:"DataLevel2BackupRetentionPeriod,omitempty"`
	OwnerAccount                                 *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId                                      *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PreferredBackupPeriod                        *string `json:"PreferredBackupPeriod,omitempty" xml:"PreferredBackupPeriod,omitempty"`
	PreferredBackupTime                          *string `json:"PreferredBackupTime,omitempty" xml:"PreferredBackupTime,omitempty"`
	ResourceOwnerAccount                         *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId                              *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyBackupPolicyRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyBackupPolicyRequest) GoString() string {
	return s.String()
}

func (s *ModifyBackupPolicyRequest) SetBackupFrequency(v string) *ModifyBackupPolicyRequest {
	s.BackupFrequency = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetBackupRetentionPolicyOnClusterDeletion(v string) *ModifyBackupPolicyRequest {
	s.BackupRetentionPolicyOnClusterDeletion = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDBClusterId(v string) *ModifyBackupPolicyRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDataLevel1BackupFrequency(v string) *ModifyBackupPolicyRequest {
	s.DataLevel1BackupFrequency = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDataLevel1BackupPeriod(v string) *ModifyBackupPolicyRequest {
	s.DataLevel1BackupPeriod = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDataLevel1BackupRetentionPeriod(v string) *ModifyBackupPolicyRequest {
	s.DataLevel1BackupRetentionPeriod = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDataLevel1BackupTime(v string) *ModifyBackupPolicyRequest {
	s.DataLevel1BackupTime = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDataLevel2BackupAnotherRegionRegion(v string) *ModifyBackupPolicyRequest {
	s.DataLevel2BackupAnotherRegionRegion = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDataLevel2BackupAnotherRegionRetentionPeriod(v string) *ModifyBackupPolicyRequest {
	s.DataLevel2BackupAnotherRegionRetentionPeriod = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDataLevel2BackupPeriod(v string) *ModifyBackupPolicyRequest {
	s.DataLevel2BackupPeriod = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetDataLevel2BackupRetentionPeriod(v string) *ModifyBackupPolicyRequest {
	s.DataLevel2BackupRetentionPeriod = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetOwnerAccount(v string) *ModifyBackupPolicyRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetOwnerId(v int64) *ModifyBackupPolicyRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetPreferredBackupPeriod(v string) *ModifyBackupPolicyRequest {
	s.PreferredBackupPeriod = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetPreferredBackupTime(v string) *ModifyBackupPolicyRequest {
	s.PreferredBackupTime = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetResourceOwnerAccount(v string) *ModifyBackupPolicyRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyBackupPolicyRequest) SetResourceOwnerId(v int64) *ModifyBackupPolicyRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyBackupPolicyResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyBackupPolicyResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyBackupPolicyResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyBackupPolicyResponseBody) SetRequestId(v string) *ModifyBackupPolicyResponseBody {
	s.RequestId = &v
	return s
}

type ModifyBackupPolicyResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyBackupPolicyResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyBackupPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyBackupPolicyResponse) GoString() string {
	return s.String()
}

func (s *ModifyBackupPolicyResponse) SetHeaders(v map[string]*string) *ModifyBackupPolicyResponse {
	s.Headers = v
	return s
}

func (s *ModifyBackupPolicyResponse) SetStatusCode(v int32) *ModifyBackupPolicyResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyBackupPolicyResponse) SetBody(v *ModifyBackupPolicyResponseBody) *ModifyBackupPolicyResponse {
	s.Body = v
	return s
}

type ModifyDBClusterAccessWhitelistRequest struct {
	DBClusterIPArrayAttribute *string `json:"DBClusterIPArrayAttribute,omitempty" xml:"DBClusterIPArrayAttribute,omitempty"`
	DBClusterIPArrayName      *string `json:"DBClusterIPArrayName,omitempty" xml:"DBClusterIPArrayName,omitempty"`
	DBClusterId               *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	ModifyMode                *string `json:"ModifyMode,omitempty" xml:"ModifyMode,omitempty"`
	OwnerAccount              *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId                   *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount      *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId           *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityGroupIds          *string `json:"SecurityGroupIds,omitempty" xml:"SecurityGroupIds,omitempty"`
	SecurityIps               *string `json:"SecurityIps,omitempty" xml:"SecurityIps,omitempty"`
	WhiteListType             *string `json:"WhiteListType,omitempty" xml:"WhiteListType,omitempty"`
}

func (s ModifyDBClusterAccessWhitelistRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAccessWhitelistRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetDBClusterIPArrayAttribute(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.DBClusterIPArrayAttribute = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetDBClusterIPArrayName(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.DBClusterIPArrayName = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetDBClusterId(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetModifyMode(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.ModifyMode = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetOwnerAccount(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetOwnerId(v int64) *ModifyDBClusterAccessWhitelistRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetResourceOwnerId(v int64) *ModifyDBClusterAccessWhitelistRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetSecurityGroupIds(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.SecurityGroupIds = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetSecurityIps(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.SecurityIps = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistRequest) SetWhiteListType(v string) *ModifyDBClusterAccessWhitelistRequest {
	s.WhiteListType = &v
	return s
}

type ModifyDBClusterAccessWhitelistResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterAccessWhitelistResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAccessWhitelistResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAccessWhitelistResponseBody) SetRequestId(v string) *ModifyDBClusterAccessWhitelistResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterAccessWhitelistResponse struct {
	Headers    map[string]*string                          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterAccessWhitelistResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterAccessWhitelistResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAccessWhitelistResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAccessWhitelistResponse) SetHeaders(v map[string]*string) *ModifyDBClusterAccessWhitelistResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterAccessWhitelistResponse) SetStatusCode(v int32) *ModifyDBClusterAccessWhitelistResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterAccessWhitelistResponse) SetBody(v *ModifyDBClusterAccessWhitelistResponseBody) *ModifyDBClusterAccessWhitelistResponse {
	s.Body = v
	return s
}

type ModifyDBClusterAndNodesParametersRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeIds            *string `json:"DBNodeIds,omitempty" xml:"DBNodeIds,omitempty"`
	FromTimeService      *bool   `json:"FromTimeService,omitempty" xml:"FromTimeService,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ParameterGroupId     *string `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	Parameters           *string `json:"Parameters,omitempty" xml:"Parameters,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterAndNodesParametersRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAndNodesParametersRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetDBClusterId(v string) *ModifyDBClusterAndNodesParametersRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetDBNodeIds(v string) *ModifyDBClusterAndNodesParametersRequest {
	s.DBNodeIds = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetFromTimeService(v bool) *ModifyDBClusterAndNodesParametersRequest {
	s.FromTimeService = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetOwnerAccount(v string) *ModifyDBClusterAndNodesParametersRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetOwnerId(v int64) *ModifyDBClusterAndNodesParametersRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetParameterGroupId(v string) *ModifyDBClusterAndNodesParametersRequest {
	s.ParameterGroupId = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetParameters(v string) *ModifyDBClusterAndNodesParametersRequest {
	s.Parameters = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetPlannedEndTime(v string) *ModifyDBClusterAndNodesParametersRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetPlannedStartTime(v string) *ModifyDBClusterAndNodesParametersRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterAndNodesParametersRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersRequest) SetResourceOwnerId(v int64) *ModifyDBClusterAndNodesParametersRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterAndNodesParametersResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterAndNodesParametersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAndNodesParametersResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAndNodesParametersResponseBody) SetRequestId(v string) *ModifyDBClusterAndNodesParametersResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterAndNodesParametersResponse struct {
	Headers    map[string]*string                             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterAndNodesParametersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterAndNodesParametersResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAndNodesParametersResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAndNodesParametersResponse) SetHeaders(v map[string]*string) *ModifyDBClusterAndNodesParametersResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterAndNodesParametersResponse) SetStatusCode(v int32) *ModifyDBClusterAndNodesParametersResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterAndNodesParametersResponse) SetBody(v *ModifyDBClusterAndNodesParametersResponseBody) *ModifyDBClusterAndNodesParametersResponse {
	s.Body = v
	return s
}

type ModifyDBClusterAuditLogCollectorRequest struct {
	CollectorStatus      *string `json:"CollectorStatus,omitempty" xml:"CollectorStatus,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterAuditLogCollectorRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAuditLogCollectorRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAuditLogCollectorRequest) SetCollectorStatus(v string) *ModifyDBClusterAuditLogCollectorRequest {
	s.CollectorStatus = &v
	return s
}

func (s *ModifyDBClusterAuditLogCollectorRequest) SetDBClusterId(v string) *ModifyDBClusterAuditLogCollectorRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterAuditLogCollectorRequest) SetOwnerAccount(v string) *ModifyDBClusterAuditLogCollectorRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterAuditLogCollectorRequest) SetOwnerId(v int64) *ModifyDBClusterAuditLogCollectorRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterAuditLogCollectorRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterAuditLogCollectorRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterAuditLogCollectorRequest) SetResourceOwnerId(v int64) *ModifyDBClusterAuditLogCollectorRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterAuditLogCollectorResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterAuditLogCollectorResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAuditLogCollectorResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAuditLogCollectorResponseBody) SetRequestId(v string) *ModifyDBClusterAuditLogCollectorResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterAuditLogCollectorResponse struct {
	Headers    map[string]*string                            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterAuditLogCollectorResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterAuditLogCollectorResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterAuditLogCollectorResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterAuditLogCollectorResponse) SetHeaders(v map[string]*string) *ModifyDBClusterAuditLogCollectorResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterAuditLogCollectorResponse) SetStatusCode(v int32) *ModifyDBClusterAuditLogCollectorResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterAuditLogCollectorResponse) SetBody(v *ModifyDBClusterAuditLogCollectorResponseBody) *ModifyDBClusterAuditLogCollectorResponse {
	s.Body = v
	return s
}

type ModifyDBClusterDeletionRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	Protection           *bool   `json:"Protection,omitempty" xml:"Protection,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterDeletionRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterDeletionRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterDeletionRequest) SetDBClusterId(v string) *ModifyDBClusterDeletionRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterDeletionRequest) SetOwnerAccount(v string) *ModifyDBClusterDeletionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterDeletionRequest) SetOwnerId(v int64) *ModifyDBClusterDeletionRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterDeletionRequest) SetProtection(v bool) *ModifyDBClusterDeletionRequest {
	s.Protection = &v
	return s
}

func (s *ModifyDBClusterDeletionRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterDeletionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterDeletionRequest) SetResourceOwnerId(v int64) *ModifyDBClusterDeletionRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterDeletionResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterDeletionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterDeletionResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterDeletionResponseBody) SetRequestId(v string) *ModifyDBClusterDeletionResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterDeletionResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterDeletionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterDeletionResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterDeletionResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterDeletionResponse) SetHeaders(v map[string]*string) *ModifyDBClusterDeletionResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterDeletionResponse) SetStatusCode(v int32) *ModifyDBClusterDeletionResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterDeletionResponse) SetBody(v *ModifyDBClusterDeletionResponseBody) *ModifyDBClusterDeletionResponse {
	s.Body = v
	return s
}

type ModifyDBClusterDescriptionRequest struct {
	DBClusterDescription *string `json:"DBClusterDescription,omitempty" xml:"DBClusterDescription,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterDescriptionRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterDescriptionRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterDescriptionRequest) SetDBClusterDescription(v string) *ModifyDBClusterDescriptionRequest {
	s.DBClusterDescription = &v
	return s
}

func (s *ModifyDBClusterDescriptionRequest) SetDBClusterId(v string) *ModifyDBClusterDescriptionRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterDescriptionRequest) SetOwnerAccount(v string) *ModifyDBClusterDescriptionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterDescriptionRequest) SetOwnerId(v int64) *ModifyDBClusterDescriptionRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterDescriptionRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterDescriptionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterDescriptionRequest) SetResourceOwnerId(v int64) *ModifyDBClusterDescriptionRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterDescriptionResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterDescriptionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterDescriptionResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterDescriptionResponseBody) SetRequestId(v string) *ModifyDBClusterDescriptionResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterDescriptionResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterDescriptionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterDescriptionResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterDescriptionResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterDescriptionResponse) SetHeaders(v map[string]*string) *ModifyDBClusterDescriptionResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterDescriptionResponse) SetStatusCode(v int32) *ModifyDBClusterDescriptionResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterDescriptionResponse) SetBody(v *ModifyDBClusterDescriptionResponseBody) *ModifyDBClusterDescriptionResponse {
	s.Body = v
	return s
}

type ModifyDBClusterEndpointRequest struct {
	AutoAddNewNodes       *string `json:"AutoAddNewNodes,omitempty" xml:"AutoAddNewNodes,omitempty"`
	DBClusterId           *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBEndpointDescription *string `json:"DBEndpointDescription,omitempty" xml:"DBEndpointDescription,omitempty"`
	DBEndpointId          *string `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	EndpointConfig        *string `json:"EndpointConfig,omitempty" xml:"EndpointConfig,omitempty"`
	Nodes                 *string `json:"Nodes,omitempty" xml:"Nodes,omitempty"`
	OwnerAccount          *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId               *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ReadWriteMode         *string `json:"ReadWriteMode,omitempty" xml:"ReadWriteMode,omitempty"`
	ResourceOwnerAccount  *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId       *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterEndpointRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterEndpointRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterEndpointRequest) SetAutoAddNewNodes(v string) *ModifyDBClusterEndpointRequest {
	s.AutoAddNewNodes = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetDBClusterId(v string) *ModifyDBClusterEndpointRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetDBEndpointDescription(v string) *ModifyDBClusterEndpointRequest {
	s.DBEndpointDescription = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetDBEndpointId(v string) *ModifyDBClusterEndpointRequest {
	s.DBEndpointId = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetEndpointConfig(v string) *ModifyDBClusterEndpointRequest {
	s.EndpointConfig = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetNodes(v string) *ModifyDBClusterEndpointRequest {
	s.Nodes = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetOwnerAccount(v string) *ModifyDBClusterEndpointRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetOwnerId(v int64) *ModifyDBClusterEndpointRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetReadWriteMode(v string) *ModifyDBClusterEndpointRequest {
	s.ReadWriteMode = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterEndpointRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterEndpointRequest) SetResourceOwnerId(v int64) *ModifyDBClusterEndpointRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterEndpointResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterEndpointResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterEndpointResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterEndpointResponseBody) SetRequestId(v string) *ModifyDBClusterEndpointResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterEndpointResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterEndpointResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterEndpointResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterEndpointResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterEndpointResponse) SetHeaders(v map[string]*string) *ModifyDBClusterEndpointResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterEndpointResponse) SetStatusCode(v int32) *ModifyDBClusterEndpointResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterEndpointResponse) SetBody(v *ModifyDBClusterEndpointResponseBody) *ModifyDBClusterEndpointResponse {
	s.Body = v
	return s
}

type ModifyDBClusterMaintainTimeRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	MaintainTime         *string `json:"MaintainTime,omitempty" xml:"MaintainTime,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterMaintainTimeRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMaintainTimeRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMaintainTimeRequest) SetDBClusterId(v string) *ModifyDBClusterMaintainTimeRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterMaintainTimeRequest) SetMaintainTime(v string) *ModifyDBClusterMaintainTimeRequest {
	s.MaintainTime = &v
	return s
}

func (s *ModifyDBClusterMaintainTimeRequest) SetOwnerAccount(v string) *ModifyDBClusterMaintainTimeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterMaintainTimeRequest) SetOwnerId(v int64) *ModifyDBClusterMaintainTimeRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterMaintainTimeRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterMaintainTimeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterMaintainTimeRequest) SetResourceOwnerId(v int64) *ModifyDBClusterMaintainTimeRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterMaintainTimeResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterMaintainTimeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMaintainTimeResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMaintainTimeResponseBody) SetRequestId(v string) *ModifyDBClusterMaintainTimeResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterMaintainTimeResponse struct {
	Headers    map[string]*string                       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterMaintainTimeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterMaintainTimeResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMaintainTimeResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMaintainTimeResponse) SetHeaders(v map[string]*string) *ModifyDBClusterMaintainTimeResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterMaintainTimeResponse) SetStatusCode(v int32) *ModifyDBClusterMaintainTimeResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterMaintainTimeResponse) SetBody(v *ModifyDBClusterMaintainTimeResponseBody) *ModifyDBClusterMaintainTimeResponse {
	s.Body = v
	return s
}

type ModifyDBClusterMigrationRequest struct {
	ConnectionStrings     *string `json:"ConnectionStrings,omitempty" xml:"ConnectionStrings,omitempty"`
	DBClusterId           *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	NewMasterInstanceId   *string `json:"NewMasterInstanceId,omitempty" xml:"NewMasterInstanceId,omitempty"`
	OwnerAccount          *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId               *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount  *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId       *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken         *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
	SourceRDSDBInstanceId *string `json:"SourceRDSDBInstanceId,omitempty" xml:"SourceRDSDBInstanceId,omitempty"`
	SwapConnectionString  *string `json:"SwapConnectionString,omitempty" xml:"SwapConnectionString,omitempty"`
}

func (s ModifyDBClusterMigrationRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMigrationRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMigrationRequest) SetConnectionStrings(v string) *ModifyDBClusterMigrationRequest {
	s.ConnectionStrings = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetDBClusterId(v string) *ModifyDBClusterMigrationRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetNewMasterInstanceId(v string) *ModifyDBClusterMigrationRequest {
	s.NewMasterInstanceId = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetOwnerAccount(v string) *ModifyDBClusterMigrationRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetOwnerId(v int64) *ModifyDBClusterMigrationRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterMigrationRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetResourceOwnerId(v int64) *ModifyDBClusterMigrationRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetSecurityToken(v string) *ModifyDBClusterMigrationRequest {
	s.SecurityToken = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetSourceRDSDBInstanceId(v string) *ModifyDBClusterMigrationRequest {
	s.SourceRDSDBInstanceId = &v
	return s
}

func (s *ModifyDBClusterMigrationRequest) SetSwapConnectionString(v string) *ModifyDBClusterMigrationRequest {
	s.SwapConnectionString = &v
	return s
}

type ModifyDBClusterMigrationResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterMigrationResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMigrationResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMigrationResponseBody) SetRequestId(v string) *ModifyDBClusterMigrationResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterMigrationResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterMigrationResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterMigrationResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMigrationResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMigrationResponse) SetHeaders(v map[string]*string) *ModifyDBClusterMigrationResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterMigrationResponse) SetStatusCode(v int32) *ModifyDBClusterMigrationResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterMigrationResponse) SetBody(v *ModifyDBClusterMigrationResponseBody) *ModifyDBClusterMigrationResponse {
	s.Body = v
	return s
}

type ModifyDBClusterMonitorRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	Period               *string `json:"Period,omitempty" xml:"Period,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterMonitorRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMonitorRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMonitorRequest) SetDBClusterId(v string) *ModifyDBClusterMonitorRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterMonitorRequest) SetOwnerAccount(v string) *ModifyDBClusterMonitorRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterMonitorRequest) SetOwnerId(v int64) *ModifyDBClusterMonitorRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterMonitorRequest) SetPeriod(v string) *ModifyDBClusterMonitorRequest {
	s.Period = &v
	return s
}

func (s *ModifyDBClusterMonitorRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterMonitorRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterMonitorRequest) SetResourceOwnerId(v int64) *ModifyDBClusterMonitorRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterMonitorResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterMonitorResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMonitorResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMonitorResponseBody) SetRequestId(v string) *ModifyDBClusterMonitorResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterMonitorResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterMonitorResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterMonitorResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterMonitorResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterMonitorResponse) SetHeaders(v map[string]*string) *ModifyDBClusterMonitorResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterMonitorResponse) SetStatusCode(v int32) *ModifyDBClusterMonitorResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterMonitorResponse) SetBody(v *ModifyDBClusterMonitorResponseBody) *ModifyDBClusterMonitorResponse {
	s.Body = v
	return s
}

type ModifyDBClusterParametersRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	FromTimeService      *bool   `json:"FromTimeService,omitempty" xml:"FromTimeService,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ParameterGroupId     *string `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	Parameters           *string `json:"Parameters,omitempty" xml:"Parameters,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterParametersRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterParametersRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterParametersRequest) SetDBClusterId(v string) *ModifyDBClusterParametersRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetFromTimeService(v bool) *ModifyDBClusterParametersRequest {
	s.FromTimeService = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetOwnerAccount(v string) *ModifyDBClusterParametersRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetOwnerId(v int64) *ModifyDBClusterParametersRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetParameterGroupId(v string) *ModifyDBClusterParametersRequest {
	s.ParameterGroupId = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetParameters(v string) *ModifyDBClusterParametersRequest {
	s.Parameters = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetPlannedEndTime(v string) *ModifyDBClusterParametersRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetPlannedStartTime(v string) *ModifyDBClusterParametersRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterParametersRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterParametersRequest) SetResourceOwnerId(v int64) *ModifyDBClusterParametersRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterParametersResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterParametersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterParametersResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterParametersResponseBody) SetRequestId(v string) *ModifyDBClusterParametersResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterParametersResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterParametersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterParametersResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterParametersResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterParametersResponse) SetHeaders(v map[string]*string) *ModifyDBClusterParametersResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterParametersResponse) SetStatusCode(v int32) *ModifyDBClusterParametersResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterParametersResponse) SetBody(v *ModifyDBClusterParametersResponseBody) *ModifyDBClusterParametersResponse {
	s.Body = v
	return s
}

type ModifyDBClusterPrimaryZoneRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	FromTimeService      *bool   `json:"FromTimeService,omitempty" xml:"FromTimeService,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	VSwitchId            *string `json:"VSwitchId,omitempty" xml:"VSwitchId,omitempty"`
	ZoneId               *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s ModifyDBClusterPrimaryZoneRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterPrimaryZoneRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetDBClusterId(v string) *ModifyDBClusterPrimaryZoneRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetFromTimeService(v bool) *ModifyDBClusterPrimaryZoneRequest {
	s.FromTimeService = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetOwnerAccount(v string) *ModifyDBClusterPrimaryZoneRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetOwnerId(v int64) *ModifyDBClusterPrimaryZoneRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetPlannedEndTime(v string) *ModifyDBClusterPrimaryZoneRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetPlannedStartTime(v string) *ModifyDBClusterPrimaryZoneRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterPrimaryZoneRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetResourceOwnerId(v int64) *ModifyDBClusterPrimaryZoneRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetVSwitchId(v string) *ModifyDBClusterPrimaryZoneRequest {
	s.VSwitchId = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneRequest) SetZoneId(v string) *ModifyDBClusterPrimaryZoneRequest {
	s.ZoneId = &v
	return s
}

type ModifyDBClusterPrimaryZoneResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterPrimaryZoneResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterPrimaryZoneResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterPrimaryZoneResponseBody) SetRequestId(v string) *ModifyDBClusterPrimaryZoneResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterPrimaryZoneResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterPrimaryZoneResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterPrimaryZoneResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterPrimaryZoneResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterPrimaryZoneResponse) SetHeaders(v map[string]*string) *ModifyDBClusterPrimaryZoneResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterPrimaryZoneResponse) SetStatusCode(v int32) *ModifyDBClusterPrimaryZoneResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterPrimaryZoneResponse) SetBody(v *ModifyDBClusterPrimaryZoneResponseBody) *ModifyDBClusterPrimaryZoneResponse {
	s.Body = v
	return s
}

type ModifyDBClusterResourceGroupRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	NewResourceGroupId   *string `json:"NewResourceGroupId,omitempty" xml:"NewResourceGroupId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceGroupId      *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBClusterResourceGroupRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterResourceGroupRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterResourceGroupRequest) SetDBClusterId(v string) *ModifyDBClusterResourceGroupRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterResourceGroupRequest) SetNewResourceGroupId(v string) *ModifyDBClusterResourceGroupRequest {
	s.NewResourceGroupId = &v
	return s
}

func (s *ModifyDBClusterResourceGroupRequest) SetOwnerAccount(v string) *ModifyDBClusterResourceGroupRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterResourceGroupRequest) SetOwnerId(v int64) *ModifyDBClusterResourceGroupRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterResourceGroupRequest) SetResourceGroupId(v string) *ModifyDBClusterResourceGroupRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *ModifyDBClusterResourceGroupRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterResourceGroupRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterResourceGroupRequest) SetResourceOwnerId(v int64) *ModifyDBClusterResourceGroupRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBClusterResourceGroupResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterResourceGroupResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterResourceGroupResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterResourceGroupResponseBody) SetRequestId(v string) *ModifyDBClusterResourceGroupResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterResourceGroupResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterResourceGroupResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterResourceGroupResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterResourceGroupResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterResourceGroupResponse) SetHeaders(v map[string]*string) *ModifyDBClusterResourceGroupResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterResourceGroupResponse) SetStatusCode(v int32) *ModifyDBClusterResourceGroupResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterResourceGroupResponse) SetBody(v *ModifyDBClusterResourceGroupResponseBody) *ModifyDBClusterResourceGroupResponse {
	s.Body = v
	return s
}

type ModifyDBClusterSSLRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBEndpointId         *string `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	NetType              *string `json:"NetType,omitempty" xml:"NetType,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SSLAutoRotate        *string `json:"SSLAutoRotate,omitempty" xml:"SSLAutoRotate,omitempty"`
	SSLEnabled           *string `json:"SSLEnabled,omitempty" xml:"SSLEnabled,omitempty"`
}

func (s ModifyDBClusterSSLRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterSSLRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterSSLRequest) SetDBClusterId(v string) *ModifyDBClusterSSLRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterSSLRequest) SetDBEndpointId(v string) *ModifyDBClusterSSLRequest {
	s.DBEndpointId = &v
	return s
}

func (s *ModifyDBClusterSSLRequest) SetNetType(v string) *ModifyDBClusterSSLRequest {
	s.NetType = &v
	return s
}

func (s *ModifyDBClusterSSLRequest) SetOwnerAccount(v string) *ModifyDBClusterSSLRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterSSLRequest) SetOwnerId(v int64) *ModifyDBClusterSSLRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterSSLRequest) SetResourceOwnerAccount(v string) *ModifyDBClusterSSLRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterSSLRequest) SetResourceOwnerId(v int64) *ModifyDBClusterSSLRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyDBClusterSSLRequest) SetSSLAutoRotate(v string) *ModifyDBClusterSSLRequest {
	s.SSLAutoRotate = &v
	return s
}

func (s *ModifyDBClusterSSLRequest) SetSSLEnabled(v string) *ModifyDBClusterSSLRequest {
	s.SSLEnabled = &v
	return s
}

type ModifyDBClusterSSLResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterSSLResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterSSLResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterSSLResponseBody) SetRequestId(v string) *ModifyDBClusterSSLResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterSSLResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterSSLResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterSSLResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterSSLResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterSSLResponse) SetHeaders(v map[string]*string) *ModifyDBClusterSSLResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterSSLResponse) SetStatusCode(v int32) *ModifyDBClusterSSLResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterSSLResponse) SetBody(v *ModifyDBClusterSSLResponseBody) *ModifyDBClusterSSLResponse {
	s.Body = v
	return s
}

type ModifyDBClusterTDERequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	EncryptNewTables     *string `json:"EncryptNewTables,omitempty" xml:"EncryptNewTables,omitempty"`
	EncryptionKey        *string `json:"EncryptionKey,omitempty" xml:"EncryptionKey,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	RoleArn              *string `json:"RoleArn,omitempty" xml:"RoleArn,omitempty"`
	TDEStatus            *string `json:"TDEStatus,omitempty" xml:"TDEStatus,omitempty"`
}

func (s ModifyDBClusterTDERequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterTDERequest) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterTDERequest) SetDBClusterId(v string) *ModifyDBClusterTDERequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBClusterTDERequest) SetEncryptNewTables(v string) *ModifyDBClusterTDERequest {
	s.EncryptNewTables = &v
	return s
}

func (s *ModifyDBClusterTDERequest) SetEncryptionKey(v string) *ModifyDBClusterTDERequest {
	s.EncryptionKey = &v
	return s
}

func (s *ModifyDBClusterTDERequest) SetOwnerAccount(v string) *ModifyDBClusterTDERequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBClusterTDERequest) SetOwnerId(v int64) *ModifyDBClusterTDERequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBClusterTDERequest) SetResourceOwnerAccount(v string) *ModifyDBClusterTDERequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBClusterTDERequest) SetResourceOwnerId(v int64) *ModifyDBClusterTDERequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyDBClusterTDERequest) SetRoleArn(v string) *ModifyDBClusterTDERequest {
	s.RoleArn = &v
	return s
}

func (s *ModifyDBClusterTDERequest) SetTDEStatus(v string) *ModifyDBClusterTDERequest {
	s.TDEStatus = &v
	return s
}

type ModifyDBClusterTDEResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBClusterTDEResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterTDEResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterTDEResponseBody) SetRequestId(v string) *ModifyDBClusterTDEResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBClusterTDEResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBClusterTDEResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBClusterTDEResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBClusterTDEResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBClusterTDEResponse) SetHeaders(v map[string]*string) *ModifyDBClusterTDEResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBClusterTDEResponse) SetStatusCode(v int32) *ModifyDBClusterTDEResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBClusterTDEResponse) SetBody(v *ModifyDBClusterTDEResponseBody) *ModifyDBClusterTDEResponse {
	s.Body = v
	return s
}

type ModifyDBDescriptionRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBDescription        *string `json:"DBDescription,omitempty" xml:"DBDescription,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBDescriptionRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBDescriptionRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBDescriptionRequest) SetDBClusterId(v string) *ModifyDBDescriptionRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBDescriptionRequest) SetDBDescription(v string) *ModifyDBDescriptionRequest {
	s.DBDescription = &v
	return s
}

func (s *ModifyDBDescriptionRequest) SetDBName(v string) *ModifyDBDescriptionRequest {
	s.DBName = &v
	return s
}

func (s *ModifyDBDescriptionRequest) SetOwnerAccount(v string) *ModifyDBDescriptionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBDescriptionRequest) SetOwnerId(v int64) *ModifyDBDescriptionRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBDescriptionRequest) SetResourceOwnerAccount(v string) *ModifyDBDescriptionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBDescriptionRequest) SetResourceOwnerId(v int64) *ModifyDBDescriptionRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBDescriptionResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBDescriptionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBDescriptionResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBDescriptionResponseBody) SetRequestId(v string) *ModifyDBDescriptionResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBDescriptionResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBDescriptionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBDescriptionResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBDescriptionResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBDescriptionResponse) SetHeaders(v map[string]*string) *ModifyDBDescriptionResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBDescriptionResponse) SetStatusCode(v int32) *ModifyDBDescriptionResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBDescriptionResponse) SetBody(v *ModifyDBDescriptionResponseBody) *ModifyDBDescriptionResponse {
	s.Body = v
	return s
}

type ModifyDBEndpointAddressRequest struct {
	ConnectionStringPrefix   *string `json:"ConnectionStringPrefix,omitempty" xml:"ConnectionStringPrefix,omitempty"`
	DBClusterId              *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBEndpointId             *string `json:"DBEndpointId,omitempty" xml:"DBEndpointId,omitempty"`
	NetType                  *string `json:"NetType,omitempty" xml:"NetType,omitempty"`
	OwnerAccount             *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId                  *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	Port                     *string `json:"Port,omitempty" xml:"Port,omitempty"`
	PrivateZoneAddressPrefix *string `json:"PrivateZoneAddressPrefix,omitempty" xml:"PrivateZoneAddressPrefix,omitempty"`
	PrivateZoneName          *string `json:"PrivateZoneName,omitempty" xml:"PrivateZoneName,omitempty"`
	ResourceOwnerAccount     *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId          *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBEndpointAddressRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBEndpointAddressRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBEndpointAddressRequest) SetConnectionStringPrefix(v string) *ModifyDBEndpointAddressRequest {
	s.ConnectionStringPrefix = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetDBClusterId(v string) *ModifyDBEndpointAddressRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetDBEndpointId(v string) *ModifyDBEndpointAddressRequest {
	s.DBEndpointId = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetNetType(v string) *ModifyDBEndpointAddressRequest {
	s.NetType = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetOwnerAccount(v string) *ModifyDBEndpointAddressRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetOwnerId(v int64) *ModifyDBEndpointAddressRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetPort(v string) *ModifyDBEndpointAddressRequest {
	s.Port = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetPrivateZoneAddressPrefix(v string) *ModifyDBEndpointAddressRequest {
	s.PrivateZoneAddressPrefix = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetPrivateZoneName(v string) *ModifyDBEndpointAddressRequest {
	s.PrivateZoneName = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetResourceOwnerAccount(v string) *ModifyDBEndpointAddressRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBEndpointAddressRequest) SetResourceOwnerId(v int64) *ModifyDBEndpointAddressRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBEndpointAddressResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBEndpointAddressResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBEndpointAddressResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBEndpointAddressResponseBody) SetRequestId(v string) *ModifyDBEndpointAddressResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBEndpointAddressResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBEndpointAddressResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBEndpointAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBEndpointAddressResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBEndpointAddressResponse) SetHeaders(v map[string]*string) *ModifyDBEndpointAddressResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBEndpointAddressResponse) SetStatusCode(v int32) *ModifyDBEndpointAddressResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBEndpointAddressResponse) SetBody(v *ModifyDBEndpointAddressResponseBody) *ModifyDBEndpointAddressResponse {
	s.Body = v
	return s
}

type ModifyDBNodeClassRequest struct {
	ClientToken          *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeTargetClass    *string `json:"DBNodeTargetClass,omitempty" xml:"DBNodeTargetClass,omitempty"`
	ModifyType           *string `json:"ModifyType,omitempty" xml:"ModifyType,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SubCategory          *string `json:"SubCategory,omitempty" xml:"SubCategory,omitempty"`
}

func (s ModifyDBNodeClassRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodeClassRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBNodeClassRequest) SetClientToken(v string) *ModifyDBNodeClassRequest {
	s.ClientToken = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetDBClusterId(v string) *ModifyDBNodeClassRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetDBNodeTargetClass(v string) *ModifyDBNodeClassRequest {
	s.DBNodeTargetClass = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetModifyType(v string) *ModifyDBNodeClassRequest {
	s.ModifyType = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetOwnerAccount(v string) *ModifyDBNodeClassRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetOwnerId(v int64) *ModifyDBNodeClassRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetPlannedEndTime(v string) *ModifyDBNodeClassRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetPlannedStartTime(v string) *ModifyDBNodeClassRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetResourceOwnerAccount(v string) *ModifyDBNodeClassRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetResourceOwnerId(v int64) *ModifyDBNodeClassRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyDBNodeClassRequest) SetSubCategory(v string) *ModifyDBNodeClassRequest {
	s.SubCategory = &v
	return s
}

type ModifyDBNodeClassResponseBody struct {
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OrderId     *string `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	RequestId   *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBNodeClassResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodeClassResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBNodeClassResponseBody) SetDBClusterId(v string) *ModifyDBNodeClassResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBNodeClassResponseBody) SetOrderId(v string) *ModifyDBNodeClassResponseBody {
	s.OrderId = &v
	return s
}

func (s *ModifyDBNodeClassResponseBody) SetRequestId(v string) *ModifyDBNodeClassResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBNodeClassResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBNodeClassResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBNodeClassResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodeClassResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBNodeClassResponse) SetHeaders(v map[string]*string) *ModifyDBNodeClassResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBNodeClassResponse) SetStatusCode(v int32) *ModifyDBNodeClassResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBNodeClassResponse) SetBody(v *ModifyDBNodeClassResponseBody) *ModifyDBNodeClassResponse {
	s.Body = v
	return s
}

type ModifyDBNodesClassRequest struct {
	ClientToken          *string                            `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string                            `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNode               []*ModifyDBNodesClassRequestDBNode `json:"DBNode,omitempty" xml:"DBNode,omitempty" type:"Repeated"`
	ModifyType           *string                            `json:"ModifyType,omitempty" xml:"ModifyType,omitempty"`
	OwnerAccount         *string                            `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64                             `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PlannedEndTime       *string                            `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string                            `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string                            `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64                             `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SubCategory          *string                            `json:"SubCategory,omitempty" xml:"SubCategory,omitempty"`
}

func (s ModifyDBNodesClassRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodesClassRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBNodesClassRequest) SetClientToken(v string) *ModifyDBNodesClassRequest {
	s.ClientToken = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetDBClusterId(v string) *ModifyDBNodesClassRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetDBNode(v []*ModifyDBNodesClassRequestDBNode) *ModifyDBNodesClassRequest {
	s.DBNode = v
	return s
}

func (s *ModifyDBNodesClassRequest) SetModifyType(v string) *ModifyDBNodesClassRequest {
	s.ModifyType = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetOwnerAccount(v string) *ModifyDBNodesClassRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetOwnerId(v int64) *ModifyDBNodesClassRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetPlannedEndTime(v string) *ModifyDBNodesClassRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetPlannedStartTime(v string) *ModifyDBNodesClassRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetResourceOwnerAccount(v string) *ModifyDBNodesClassRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetResourceOwnerId(v int64) *ModifyDBNodesClassRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyDBNodesClassRequest) SetSubCategory(v string) *ModifyDBNodesClassRequest {
	s.SubCategory = &v
	return s
}

type ModifyDBNodesClassRequestDBNode struct {
	DBNodeId    *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	TargetClass *string `json:"TargetClass,omitempty" xml:"TargetClass,omitempty"`
}

func (s ModifyDBNodesClassRequestDBNode) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodesClassRequestDBNode) GoString() string {
	return s.String()
}

func (s *ModifyDBNodesClassRequestDBNode) SetDBNodeId(v string) *ModifyDBNodesClassRequestDBNode {
	s.DBNodeId = &v
	return s
}

func (s *ModifyDBNodesClassRequestDBNode) SetTargetClass(v string) *ModifyDBNodesClassRequestDBNode {
	s.TargetClass = &v
	return s
}

type ModifyDBNodesClassResponseBody struct {
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OrderId     *string `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	RequestId   *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBNodesClassResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodesClassResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBNodesClassResponseBody) SetDBClusterId(v string) *ModifyDBNodesClassResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBNodesClassResponseBody) SetOrderId(v string) *ModifyDBNodesClassResponseBody {
	s.OrderId = &v
	return s
}

func (s *ModifyDBNodesClassResponseBody) SetRequestId(v string) *ModifyDBNodesClassResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBNodesClassResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBNodesClassResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBNodesClassResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodesClassResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBNodesClassResponse) SetHeaders(v map[string]*string) *ModifyDBNodesClassResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBNodesClassResponse) SetStatusCode(v int32) *ModifyDBNodesClassResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBNodesClassResponse) SetBody(v *ModifyDBNodesClassResponseBody) *ModifyDBNodesClassResponse {
	s.Body = v
	return s
}

type ModifyDBNodesParametersRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeIds            *string `json:"DBNodeIds,omitempty" xml:"DBNodeIds,omitempty"`
	FromTimeService      *bool   `json:"FromTimeService,omitempty" xml:"FromTimeService,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ParameterGroupId     *string `json:"ParameterGroupId,omitempty" xml:"ParameterGroupId,omitempty"`
	Parameters           *string `json:"Parameters,omitempty" xml:"Parameters,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyDBNodesParametersRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodesParametersRequest) GoString() string {
	return s.String()
}

func (s *ModifyDBNodesParametersRequest) SetDBClusterId(v string) *ModifyDBNodesParametersRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetDBNodeIds(v string) *ModifyDBNodesParametersRequest {
	s.DBNodeIds = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetFromTimeService(v bool) *ModifyDBNodesParametersRequest {
	s.FromTimeService = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetOwnerAccount(v string) *ModifyDBNodesParametersRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetOwnerId(v int64) *ModifyDBNodesParametersRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetParameterGroupId(v string) *ModifyDBNodesParametersRequest {
	s.ParameterGroupId = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetParameters(v string) *ModifyDBNodesParametersRequest {
	s.Parameters = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetPlannedEndTime(v string) *ModifyDBNodesParametersRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetPlannedStartTime(v string) *ModifyDBNodesParametersRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetResourceOwnerAccount(v string) *ModifyDBNodesParametersRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyDBNodesParametersRequest) SetResourceOwnerId(v int64) *ModifyDBNodesParametersRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyDBNodesParametersResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyDBNodesParametersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodesParametersResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyDBNodesParametersResponseBody) SetRequestId(v string) *ModifyDBNodesParametersResponseBody {
	s.RequestId = &v
	return s
}

type ModifyDBNodesParametersResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyDBNodesParametersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyDBNodesParametersResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyDBNodesParametersResponse) GoString() string {
	return s.String()
}

func (s *ModifyDBNodesParametersResponse) SetHeaders(v map[string]*string) *ModifyDBNodesParametersResponse {
	s.Headers = v
	return s
}

func (s *ModifyDBNodesParametersResponse) SetStatusCode(v int32) *ModifyDBNodesParametersResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyDBNodesParametersResponse) SetBody(v *ModifyDBNodesParametersResponseBody) *ModifyDBNodesParametersResponse {
	s.Body = v
	return s
}

type ModifyGlobalDatabaseNetworkRequest struct {
	GDNDescription       *string `json:"GDNDescription,omitempty" xml:"GDNDescription,omitempty"`
	GDNId                *string `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s ModifyGlobalDatabaseNetworkRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyGlobalDatabaseNetworkRequest) GoString() string {
	return s.String()
}

func (s *ModifyGlobalDatabaseNetworkRequest) SetGDNDescription(v string) *ModifyGlobalDatabaseNetworkRequest {
	s.GDNDescription = &v
	return s
}

func (s *ModifyGlobalDatabaseNetworkRequest) SetGDNId(v string) *ModifyGlobalDatabaseNetworkRequest {
	s.GDNId = &v
	return s
}

func (s *ModifyGlobalDatabaseNetworkRequest) SetOwnerAccount(v string) *ModifyGlobalDatabaseNetworkRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyGlobalDatabaseNetworkRequest) SetOwnerId(v int64) *ModifyGlobalDatabaseNetworkRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyGlobalDatabaseNetworkRequest) SetResourceOwnerAccount(v string) *ModifyGlobalDatabaseNetworkRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyGlobalDatabaseNetworkRequest) SetResourceOwnerId(v int64) *ModifyGlobalDatabaseNetworkRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyGlobalDatabaseNetworkRequest) SetSecurityToken(v string) *ModifyGlobalDatabaseNetworkRequest {
	s.SecurityToken = &v
	return s
}

type ModifyGlobalDatabaseNetworkResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyGlobalDatabaseNetworkResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyGlobalDatabaseNetworkResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyGlobalDatabaseNetworkResponseBody) SetRequestId(v string) *ModifyGlobalDatabaseNetworkResponseBody {
	s.RequestId = &v
	return s
}

type ModifyGlobalDatabaseNetworkResponse struct {
	Headers    map[string]*string                       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyGlobalDatabaseNetworkResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyGlobalDatabaseNetworkResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyGlobalDatabaseNetworkResponse) GoString() string {
	return s.String()
}

func (s *ModifyGlobalDatabaseNetworkResponse) SetHeaders(v map[string]*string) *ModifyGlobalDatabaseNetworkResponse {
	s.Headers = v
	return s
}

func (s *ModifyGlobalDatabaseNetworkResponse) SetStatusCode(v int32) *ModifyGlobalDatabaseNetworkResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyGlobalDatabaseNetworkResponse) SetBody(v *ModifyGlobalDatabaseNetworkResponseBody) *ModifyGlobalDatabaseNetworkResponse {
	s.Body = v
	return s
}

type ModifyLogBackupPolicyRequest struct {
	DBClusterId                           *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	LogBackupAnotherRegionRegion          *string `json:"LogBackupAnotherRegionRegion,omitempty" xml:"LogBackupAnotherRegionRegion,omitempty"`
	LogBackupAnotherRegionRetentionPeriod *string `json:"LogBackupAnotherRegionRetentionPeriod,omitempty" xml:"LogBackupAnotherRegionRetentionPeriod,omitempty"`
	LogBackupRetentionPeriod              *string `json:"LogBackupRetentionPeriod,omitempty" xml:"LogBackupRetentionPeriod,omitempty"`
	OwnerAccount                          *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId                               *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount                  *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId                       *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ModifyLogBackupPolicyRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyLogBackupPolicyRequest) GoString() string {
	return s.String()
}

func (s *ModifyLogBackupPolicyRequest) SetDBClusterId(v string) *ModifyLogBackupPolicyRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyLogBackupPolicyRequest) SetLogBackupAnotherRegionRegion(v string) *ModifyLogBackupPolicyRequest {
	s.LogBackupAnotherRegionRegion = &v
	return s
}

func (s *ModifyLogBackupPolicyRequest) SetLogBackupAnotherRegionRetentionPeriod(v string) *ModifyLogBackupPolicyRequest {
	s.LogBackupAnotherRegionRetentionPeriod = &v
	return s
}

func (s *ModifyLogBackupPolicyRequest) SetLogBackupRetentionPeriod(v string) *ModifyLogBackupPolicyRequest {
	s.LogBackupRetentionPeriod = &v
	return s
}

func (s *ModifyLogBackupPolicyRequest) SetOwnerAccount(v string) *ModifyLogBackupPolicyRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyLogBackupPolicyRequest) SetOwnerId(v int64) *ModifyLogBackupPolicyRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyLogBackupPolicyRequest) SetResourceOwnerAccount(v string) *ModifyLogBackupPolicyRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyLogBackupPolicyRequest) SetResourceOwnerId(v int64) *ModifyLogBackupPolicyRequest {
	s.ResourceOwnerId = &v
	return s
}

type ModifyLogBackupPolicyResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyLogBackupPolicyResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyLogBackupPolicyResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyLogBackupPolicyResponseBody) SetRequestId(v string) *ModifyLogBackupPolicyResponseBody {
	s.RequestId = &v
	return s
}

type ModifyLogBackupPolicyResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyLogBackupPolicyResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyLogBackupPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyLogBackupPolicyResponse) GoString() string {
	return s.String()
}

func (s *ModifyLogBackupPolicyResponse) SetHeaders(v map[string]*string) *ModifyLogBackupPolicyResponse {
	s.Headers = v
	return s
}

func (s *ModifyLogBackupPolicyResponse) SetStatusCode(v int32) *ModifyLogBackupPolicyResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyLogBackupPolicyResponse) SetBody(v *ModifyLogBackupPolicyResponseBody) *ModifyLogBackupPolicyResponse {
	s.Body = v
	return s
}

type ModifyMaskingRulesRequest struct {
	DBClusterId  *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	Enable       *string `json:"Enable,omitempty" xml:"Enable,omitempty"`
	RuleConfig   *string `json:"RuleConfig,omitempty" xml:"RuleConfig,omitempty"`
	RuleName     *string `json:"RuleName,omitempty" xml:"RuleName,omitempty"`
	RuleNameList *string `json:"RuleNameList,omitempty" xml:"RuleNameList,omitempty"`
}

func (s ModifyMaskingRulesRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyMaskingRulesRequest) GoString() string {
	return s.String()
}

func (s *ModifyMaskingRulesRequest) SetDBClusterId(v string) *ModifyMaskingRulesRequest {
	s.DBClusterId = &v
	return s
}

func (s *ModifyMaskingRulesRequest) SetEnable(v string) *ModifyMaskingRulesRequest {
	s.Enable = &v
	return s
}

func (s *ModifyMaskingRulesRequest) SetRuleConfig(v string) *ModifyMaskingRulesRequest {
	s.RuleConfig = &v
	return s
}

func (s *ModifyMaskingRulesRequest) SetRuleName(v string) *ModifyMaskingRulesRequest {
	s.RuleName = &v
	return s
}

func (s *ModifyMaskingRulesRequest) SetRuleNameList(v string) *ModifyMaskingRulesRequest {
	s.RuleNameList = &v
	return s
}

type ModifyMaskingRulesResponseBody struct {
	Message   *string `json:"Message,omitempty" xml:"Message,omitempty"`
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	Success   *bool   `json:"Success,omitempty" xml:"Success,omitempty"`
}

func (s ModifyMaskingRulesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyMaskingRulesResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyMaskingRulesResponseBody) SetMessage(v string) *ModifyMaskingRulesResponseBody {
	s.Message = &v
	return s
}

func (s *ModifyMaskingRulesResponseBody) SetRequestId(v string) *ModifyMaskingRulesResponseBody {
	s.RequestId = &v
	return s
}

func (s *ModifyMaskingRulesResponseBody) SetSuccess(v bool) *ModifyMaskingRulesResponseBody {
	s.Success = &v
	return s
}

type ModifyMaskingRulesResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyMaskingRulesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyMaskingRulesResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyMaskingRulesResponse) GoString() string {
	return s.String()
}

func (s *ModifyMaskingRulesResponse) SetHeaders(v map[string]*string) *ModifyMaskingRulesResponse {
	s.Headers = v
	return s
}

func (s *ModifyMaskingRulesResponse) SetStatusCode(v int32) *ModifyMaskingRulesResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyMaskingRulesResponse) SetBody(v *ModifyMaskingRulesResponseBody) *ModifyMaskingRulesResponse {
	s.Body = v
	return s
}

type ModifyPendingMaintenanceActionRequest struct {
	Ids                  *string `json:"Ids,omitempty" xml:"Ids,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
	SwitchTime           *string `json:"SwitchTime,omitempty" xml:"SwitchTime,omitempty"`
}

func (s ModifyPendingMaintenanceActionRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyPendingMaintenanceActionRequest) GoString() string {
	return s.String()
}

func (s *ModifyPendingMaintenanceActionRequest) SetIds(v string) *ModifyPendingMaintenanceActionRequest {
	s.Ids = &v
	return s
}

func (s *ModifyPendingMaintenanceActionRequest) SetOwnerAccount(v string) *ModifyPendingMaintenanceActionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ModifyPendingMaintenanceActionRequest) SetOwnerId(v int64) *ModifyPendingMaintenanceActionRequest {
	s.OwnerId = &v
	return s
}

func (s *ModifyPendingMaintenanceActionRequest) SetRegionId(v string) *ModifyPendingMaintenanceActionRequest {
	s.RegionId = &v
	return s
}

func (s *ModifyPendingMaintenanceActionRequest) SetResourceOwnerAccount(v string) *ModifyPendingMaintenanceActionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ModifyPendingMaintenanceActionRequest) SetResourceOwnerId(v int64) *ModifyPendingMaintenanceActionRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *ModifyPendingMaintenanceActionRequest) SetSecurityToken(v string) *ModifyPendingMaintenanceActionRequest {
	s.SecurityToken = &v
	return s
}

func (s *ModifyPendingMaintenanceActionRequest) SetSwitchTime(v string) *ModifyPendingMaintenanceActionRequest {
	s.SwitchTime = &v
	return s
}

type ModifyPendingMaintenanceActionResponseBody struct {
	Ids       *string `json:"Ids,omitempty" xml:"Ids,omitempty"`
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ModifyPendingMaintenanceActionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyPendingMaintenanceActionResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyPendingMaintenanceActionResponseBody) SetIds(v string) *ModifyPendingMaintenanceActionResponseBody {
	s.Ids = &v
	return s
}

func (s *ModifyPendingMaintenanceActionResponseBody) SetRequestId(v string) *ModifyPendingMaintenanceActionResponseBody {
	s.RequestId = &v
	return s
}

type ModifyPendingMaintenanceActionResponse struct {
	Headers    map[string]*string                          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyPendingMaintenanceActionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyPendingMaintenanceActionResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyPendingMaintenanceActionResponse) GoString() string {
	return s.String()
}

func (s *ModifyPendingMaintenanceActionResponse) SetHeaders(v map[string]*string) *ModifyPendingMaintenanceActionResponse {
	s.Headers = v
	return s
}

func (s *ModifyPendingMaintenanceActionResponse) SetStatusCode(v int32) *ModifyPendingMaintenanceActionResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyPendingMaintenanceActionResponse) SetBody(v *ModifyPendingMaintenanceActionResponseBody) *ModifyPendingMaintenanceActionResponse {
	s.Body = v
	return s
}

type OpenAITaskRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	Password             *string `json:"Password,omitempty" xml:"Password,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	Username             *string `json:"Username,omitempty" xml:"Username,omitempty"`
}

func (s OpenAITaskRequest) String() string {
	return tea.Prettify(s)
}

func (s OpenAITaskRequest) GoString() string {
	return s.String()
}

func (s *OpenAITaskRequest) SetDBClusterId(v string) *OpenAITaskRequest {
	s.DBClusterId = &v
	return s
}

func (s *OpenAITaskRequest) SetOwnerAccount(v string) *OpenAITaskRequest {
	s.OwnerAccount = &v
	return s
}

func (s *OpenAITaskRequest) SetOwnerId(v int64) *OpenAITaskRequest {
	s.OwnerId = &v
	return s
}

func (s *OpenAITaskRequest) SetPassword(v string) *OpenAITaskRequest {
	s.Password = &v
	return s
}

func (s *OpenAITaskRequest) SetRegionId(v string) *OpenAITaskRequest {
	s.RegionId = &v
	return s
}

func (s *OpenAITaskRequest) SetResourceOwnerAccount(v string) *OpenAITaskRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *OpenAITaskRequest) SetResourceOwnerId(v int64) *OpenAITaskRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *OpenAITaskRequest) SetUsername(v string) *OpenAITaskRequest {
	s.Username = &v
	return s
}

type OpenAITaskResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TaskId    *string `json:"TaskId,omitempty" xml:"TaskId,omitempty"`
}

func (s OpenAITaskResponseBody) String() string {
	return tea.Prettify(s)
}

func (s OpenAITaskResponseBody) GoString() string {
	return s.String()
}

func (s *OpenAITaskResponseBody) SetRequestId(v string) *OpenAITaskResponseBody {
	s.RequestId = &v
	return s
}

func (s *OpenAITaskResponseBody) SetTaskId(v string) *OpenAITaskResponseBody {
	s.TaskId = &v
	return s
}

type OpenAITaskResponse struct {
	Headers    map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *OpenAITaskResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s OpenAITaskResponse) String() string {
	return tea.Prettify(s)
}

func (s OpenAITaskResponse) GoString() string {
	return s.String()
}

func (s *OpenAITaskResponse) SetHeaders(v map[string]*string) *OpenAITaskResponse {
	s.Headers = v
	return s
}

func (s *OpenAITaskResponse) SetStatusCode(v int32) *OpenAITaskResponse {
	s.StatusCode = &v
	return s
}

func (s *OpenAITaskResponse) SetBody(v *OpenAITaskResponseBody) *OpenAITaskResponse {
	s.Body = v
	return s
}

type RefreshDBClusterStorageUsageRequest struct {
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SyncRealTime         *bool   `json:"SyncRealTime,omitempty" xml:"SyncRealTime,omitempty"`
}

func (s RefreshDBClusterStorageUsageRequest) String() string {
	return tea.Prettify(s)
}

func (s RefreshDBClusterStorageUsageRequest) GoString() string {
	return s.String()
}

func (s *RefreshDBClusterStorageUsageRequest) SetOwnerAccount(v string) *RefreshDBClusterStorageUsageRequest {
	s.OwnerAccount = &v
	return s
}

func (s *RefreshDBClusterStorageUsageRequest) SetOwnerId(v int64) *RefreshDBClusterStorageUsageRequest {
	s.OwnerId = &v
	return s
}

func (s *RefreshDBClusterStorageUsageRequest) SetResourceOwnerAccount(v string) *RefreshDBClusterStorageUsageRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *RefreshDBClusterStorageUsageRequest) SetResourceOwnerId(v int64) *RefreshDBClusterStorageUsageRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *RefreshDBClusterStorageUsageRequest) SetSyncRealTime(v bool) *RefreshDBClusterStorageUsageRequest {
	s.SyncRealTime = &v
	return s
}

type RefreshDBClusterStorageUsageResponseBody struct {
	DBClusterId         *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	RequestId           *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	UsedStorage         *string `json:"UsedStorage,omitempty" xml:"UsedStorage,omitempty"`
	UsedStorageModified *string `json:"UsedStorageModified,omitempty" xml:"UsedStorageModified,omitempty"`
}

func (s RefreshDBClusterStorageUsageResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RefreshDBClusterStorageUsageResponseBody) GoString() string {
	return s.String()
}

func (s *RefreshDBClusterStorageUsageResponseBody) SetDBClusterId(v string) *RefreshDBClusterStorageUsageResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *RefreshDBClusterStorageUsageResponseBody) SetRequestId(v string) *RefreshDBClusterStorageUsageResponseBody {
	s.RequestId = &v
	return s
}

func (s *RefreshDBClusterStorageUsageResponseBody) SetUsedStorage(v string) *RefreshDBClusterStorageUsageResponseBody {
	s.UsedStorage = &v
	return s
}

func (s *RefreshDBClusterStorageUsageResponseBody) SetUsedStorageModified(v string) *RefreshDBClusterStorageUsageResponseBody {
	s.UsedStorageModified = &v
	return s
}

type RefreshDBClusterStorageUsageResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RefreshDBClusterStorageUsageResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RefreshDBClusterStorageUsageResponse) String() string {
	return tea.Prettify(s)
}

func (s RefreshDBClusterStorageUsageResponse) GoString() string {
	return s.String()
}

func (s *RefreshDBClusterStorageUsageResponse) SetHeaders(v map[string]*string) *RefreshDBClusterStorageUsageResponse {
	s.Headers = v
	return s
}

func (s *RefreshDBClusterStorageUsageResponse) SetStatusCode(v int32) *RefreshDBClusterStorageUsageResponse {
	s.StatusCode = &v
	return s
}

func (s *RefreshDBClusterStorageUsageResponse) SetBody(v *RefreshDBClusterStorageUsageResponseBody) *RefreshDBClusterStorageUsageResponse {
	s.Body = v
	return s
}

type RemoveDBClusterFromGDNRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	GDNId                *string `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s RemoveDBClusterFromGDNRequest) String() string {
	return tea.Prettify(s)
}

func (s RemoveDBClusterFromGDNRequest) GoString() string {
	return s.String()
}

func (s *RemoveDBClusterFromGDNRequest) SetDBClusterId(v string) *RemoveDBClusterFromGDNRequest {
	s.DBClusterId = &v
	return s
}

func (s *RemoveDBClusterFromGDNRequest) SetGDNId(v string) *RemoveDBClusterFromGDNRequest {
	s.GDNId = &v
	return s
}

func (s *RemoveDBClusterFromGDNRequest) SetOwnerAccount(v string) *RemoveDBClusterFromGDNRequest {
	s.OwnerAccount = &v
	return s
}

func (s *RemoveDBClusterFromGDNRequest) SetOwnerId(v int64) *RemoveDBClusterFromGDNRequest {
	s.OwnerId = &v
	return s
}

func (s *RemoveDBClusterFromGDNRequest) SetResourceOwnerAccount(v string) *RemoveDBClusterFromGDNRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *RemoveDBClusterFromGDNRequest) SetResourceOwnerId(v int64) *RemoveDBClusterFromGDNRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *RemoveDBClusterFromGDNRequest) SetSecurityToken(v string) *RemoveDBClusterFromGDNRequest {
	s.SecurityToken = &v
	return s
}

type RemoveDBClusterFromGDNResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s RemoveDBClusterFromGDNResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RemoveDBClusterFromGDNResponseBody) GoString() string {
	return s.String()
}

func (s *RemoveDBClusterFromGDNResponseBody) SetRequestId(v string) *RemoveDBClusterFromGDNResponseBody {
	s.RequestId = &v
	return s
}

type RemoveDBClusterFromGDNResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RemoveDBClusterFromGDNResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RemoveDBClusterFromGDNResponse) String() string {
	return tea.Prettify(s)
}

func (s RemoveDBClusterFromGDNResponse) GoString() string {
	return s.String()
}

func (s *RemoveDBClusterFromGDNResponse) SetHeaders(v map[string]*string) *RemoveDBClusterFromGDNResponse {
	s.Headers = v
	return s
}

func (s *RemoveDBClusterFromGDNResponse) SetStatusCode(v int32) *RemoveDBClusterFromGDNResponse {
	s.StatusCode = &v
	return s
}

func (s *RemoveDBClusterFromGDNResponse) SetBody(v *RemoveDBClusterFromGDNResponseBody) *RemoveDBClusterFromGDNResponse {
	s.Body = v
	return s
}

type ResetAccountRequest struct {
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	AccountPassword      *string `json:"AccountPassword,omitempty" xml:"AccountPassword,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s ResetAccountRequest) String() string {
	return tea.Prettify(s)
}

func (s ResetAccountRequest) GoString() string {
	return s.String()
}

func (s *ResetAccountRequest) SetAccountName(v string) *ResetAccountRequest {
	s.AccountName = &v
	return s
}

func (s *ResetAccountRequest) SetAccountPassword(v string) *ResetAccountRequest {
	s.AccountPassword = &v
	return s
}

func (s *ResetAccountRequest) SetDBClusterId(v string) *ResetAccountRequest {
	s.DBClusterId = &v
	return s
}

func (s *ResetAccountRequest) SetOwnerAccount(v string) *ResetAccountRequest {
	s.OwnerAccount = &v
	return s
}

func (s *ResetAccountRequest) SetOwnerId(v int64) *ResetAccountRequest {
	s.OwnerId = &v
	return s
}

func (s *ResetAccountRequest) SetResourceOwnerAccount(v string) *ResetAccountRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *ResetAccountRequest) SetResourceOwnerId(v int64) *ResetAccountRequest {
	s.ResourceOwnerId = &v
	return s
}

type ResetAccountResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s ResetAccountResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ResetAccountResponseBody) GoString() string {
	return s.String()
}

func (s *ResetAccountResponseBody) SetRequestId(v string) *ResetAccountResponseBody {
	s.RequestId = &v
	return s
}

type ResetAccountResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ResetAccountResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ResetAccountResponse) String() string {
	return tea.Prettify(s)
}

func (s ResetAccountResponse) GoString() string {
	return s.String()
}

func (s *ResetAccountResponse) SetHeaders(v map[string]*string) *ResetAccountResponse {
	s.Headers = v
	return s
}

func (s *ResetAccountResponse) SetStatusCode(v int32) *ResetAccountResponse {
	s.StatusCode = &v
	return s
}

func (s *ResetAccountResponse) SetBody(v *ResetAccountResponseBody) *ResetAccountResponse {
	s.Body = v
	return s
}

type RestartDBNodeRequest struct {
	DBNodeId             *string `json:"DBNodeId,omitempty" xml:"DBNodeId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s RestartDBNodeRequest) String() string {
	return tea.Prettify(s)
}

func (s RestartDBNodeRequest) GoString() string {
	return s.String()
}

func (s *RestartDBNodeRequest) SetDBNodeId(v string) *RestartDBNodeRequest {
	s.DBNodeId = &v
	return s
}

func (s *RestartDBNodeRequest) SetOwnerAccount(v string) *RestartDBNodeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *RestartDBNodeRequest) SetOwnerId(v int64) *RestartDBNodeRequest {
	s.OwnerId = &v
	return s
}

func (s *RestartDBNodeRequest) SetResourceOwnerAccount(v string) *RestartDBNodeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *RestartDBNodeRequest) SetResourceOwnerId(v int64) *RestartDBNodeRequest {
	s.ResourceOwnerId = &v
	return s
}

type RestartDBNodeResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s RestartDBNodeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RestartDBNodeResponseBody) GoString() string {
	return s.String()
}

func (s *RestartDBNodeResponseBody) SetRequestId(v string) *RestartDBNodeResponseBody {
	s.RequestId = &v
	return s
}

type RestartDBNodeResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RestartDBNodeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RestartDBNodeResponse) String() string {
	return tea.Prettify(s)
}

func (s RestartDBNodeResponse) GoString() string {
	return s.String()
}

func (s *RestartDBNodeResponse) SetHeaders(v map[string]*string) *RestartDBNodeResponse {
	s.Headers = v
	return s
}

func (s *RestartDBNodeResponse) SetStatusCode(v int32) *RestartDBNodeResponse {
	s.StatusCode = &v
	return s
}

func (s *RestartDBNodeResponse) SetBody(v *RestartDBNodeResponseBody) *RestartDBNodeResponse {
	s.Body = v
	return s
}

type RestoreTableRequest struct {
	BackupId             *string `json:"BackupId,omitempty" xml:"BackupId,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	RestoreTime          *string `json:"RestoreTime,omitempty" xml:"RestoreTime,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
	TableMeta            *string `json:"TableMeta,omitempty" xml:"TableMeta,omitempty"`
}

func (s RestoreTableRequest) String() string {
	return tea.Prettify(s)
}

func (s RestoreTableRequest) GoString() string {
	return s.String()
}

func (s *RestoreTableRequest) SetBackupId(v string) *RestoreTableRequest {
	s.BackupId = &v
	return s
}

func (s *RestoreTableRequest) SetDBClusterId(v string) *RestoreTableRequest {
	s.DBClusterId = &v
	return s
}

func (s *RestoreTableRequest) SetOwnerAccount(v string) *RestoreTableRequest {
	s.OwnerAccount = &v
	return s
}

func (s *RestoreTableRequest) SetOwnerId(v int64) *RestoreTableRequest {
	s.OwnerId = &v
	return s
}

func (s *RestoreTableRequest) SetResourceOwnerAccount(v string) *RestoreTableRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *RestoreTableRequest) SetResourceOwnerId(v int64) *RestoreTableRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *RestoreTableRequest) SetRestoreTime(v string) *RestoreTableRequest {
	s.RestoreTime = &v
	return s
}

func (s *RestoreTableRequest) SetSecurityToken(v string) *RestoreTableRequest {
	s.SecurityToken = &v
	return s
}

func (s *RestoreTableRequest) SetTableMeta(v string) *RestoreTableRequest {
	s.TableMeta = &v
	return s
}

type RestoreTableResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s RestoreTableResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RestoreTableResponseBody) GoString() string {
	return s.String()
}

func (s *RestoreTableResponseBody) SetRequestId(v string) *RestoreTableResponseBody {
	s.RequestId = &v
	return s
}

type RestoreTableResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RestoreTableResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RestoreTableResponse) String() string {
	return tea.Prettify(s)
}

func (s RestoreTableResponse) GoString() string {
	return s.String()
}

func (s *RestoreTableResponse) SetHeaders(v map[string]*string) *RestoreTableResponse {
	s.Headers = v
	return s
}

func (s *RestoreTableResponse) SetStatusCode(v int32) *RestoreTableResponse {
	s.StatusCode = &v
	return s
}

func (s *RestoreTableResponse) SetBody(v *RestoreTableResponseBody) *RestoreTableResponse {
	s.Body = v
	return s
}

type RevokeAccountPrivilegeRequest struct {
	AccountName          *string `json:"AccountName,omitempty" xml:"AccountName,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBName               *string `json:"DBName,omitempty" xml:"DBName,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s RevokeAccountPrivilegeRequest) String() string {
	return tea.Prettify(s)
}

func (s RevokeAccountPrivilegeRequest) GoString() string {
	return s.String()
}

func (s *RevokeAccountPrivilegeRequest) SetAccountName(v string) *RevokeAccountPrivilegeRequest {
	s.AccountName = &v
	return s
}

func (s *RevokeAccountPrivilegeRequest) SetDBClusterId(v string) *RevokeAccountPrivilegeRequest {
	s.DBClusterId = &v
	return s
}

func (s *RevokeAccountPrivilegeRequest) SetDBName(v string) *RevokeAccountPrivilegeRequest {
	s.DBName = &v
	return s
}

func (s *RevokeAccountPrivilegeRequest) SetOwnerAccount(v string) *RevokeAccountPrivilegeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *RevokeAccountPrivilegeRequest) SetOwnerId(v int64) *RevokeAccountPrivilegeRequest {
	s.OwnerId = &v
	return s
}

func (s *RevokeAccountPrivilegeRequest) SetResourceOwnerAccount(v string) *RevokeAccountPrivilegeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *RevokeAccountPrivilegeRequest) SetResourceOwnerId(v int64) *RevokeAccountPrivilegeRequest {
	s.ResourceOwnerId = &v
	return s
}

type RevokeAccountPrivilegeResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s RevokeAccountPrivilegeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RevokeAccountPrivilegeResponseBody) GoString() string {
	return s.String()
}

func (s *RevokeAccountPrivilegeResponseBody) SetRequestId(v string) *RevokeAccountPrivilegeResponseBody {
	s.RequestId = &v
	return s
}

type RevokeAccountPrivilegeResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RevokeAccountPrivilegeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RevokeAccountPrivilegeResponse) String() string {
	return tea.Prettify(s)
}

func (s RevokeAccountPrivilegeResponse) GoString() string {
	return s.String()
}

func (s *RevokeAccountPrivilegeResponse) SetHeaders(v map[string]*string) *RevokeAccountPrivilegeResponse {
	s.Headers = v
	return s
}

func (s *RevokeAccountPrivilegeResponse) SetStatusCode(v int32) *RevokeAccountPrivilegeResponse {
	s.StatusCode = &v
	return s
}

func (s *RevokeAccountPrivilegeResponse) SetBody(v *RevokeAccountPrivilegeResponseBody) *RevokeAccountPrivilegeResponse {
	s.Body = v
	return s
}

type SwitchOverGlobalDatabaseNetworkRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	GDNId                *string `json:"GDNId,omitempty" xml:"GDNId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	SecurityToken        *string `json:"SecurityToken,omitempty" xml:"SecurityToken,omitempty"`
}

func (s SwitchOverGlobalDatabaseNetworkRequest) String() string {
	return tea.Prettify(s)
}

func (s SwitchOverGlobalDatabaseNetworkRequest) GoString() string {
	return s.String()
}

func (s *SwitchOverGlobalDatabaseNetworkRequest) SetDBClusterId(v string) *SwitchOverGlobalDatabaseNetworkRequest {
	s.DBClusterId = &v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkRequest) SetGDNId(v string) *SwitchOverGlobalDatabaseNetworkRequest {
	s.GDNId = &v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkRequest) SetOwnerAccount(v string) *SwitchOverGlobalDatabaseNetworkRequest {
	s.OwnerAccount = &v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkRequest) SetOwnerId(v int64) *SwitchOverGlobalDatabaseNetworkRequest {
	s.OwnerId = &v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkRequest) SetRegionId(v string) *SwitchOverGlobalDatabaseNetworkRequest {
	s.RegionId = &v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkRequest) SetResourceOwnerAccount(v string) *SwitchOverGlobalDatabaseNetworkRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkRequest) SetResourceOwnerId(v int64) *SwitchOverGlobalDatabaseNetworkRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkRequest) SetSecurityToken(v string) *SwitchOverGlobalDatabaseNetworkRequest {
	s.SecurityToken = &v
	return s
}

type SwitchOverGlobalDatabaseNetworkResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s SwitchOverGlobalDatabaseNetworkResponseBody) String() string {
	return tea.Prettify(s)
}

func (s SwitchOverGlobalDatabaseNetworkResponseBody) GoString() string {
	return s.String()
}

func (s *SwitchOverGlobalDatabaseNetworkResponseBody) SetRequestId(v string) *SwitchOverGlobalDatabaseNetworkResponseBody {
	s.RequestId = &v
	return s
}

type SwitchOverGlobalDatabaseNetworkResponse struct {
	Headers    map[string]*string                           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *SwitchOverGlobalDatabaseNetworkResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s SwitchOverGlobalDatabaseNetworkResponse) String() string {
	return tea.Prettify(s)
}

func (s SwitchOverGlobalDatabaseNetworkResponse) GoString() string {
	return s.String()
}

func (s *SwitchOverGlobalDatabaseNetworkResponse) SetHeaders(v map[string]*string) *SwitchOverGlobalDatabaseNetworkResponse {
	s.Headers = v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkResponse) SetStatusCode(v int32) *SwitchOverGlobalDatabaseNetworkResponse {
	s.StatusCode = &v
	return s
}

func (s *SwitchOverGlobalDatabaseNetworkResponse) SetBody(v *SwitchOverGlobalDatabaseNetworkResponseBody) *SwitchOverGlobalDatabaseNetworkResponse {
	s.Body = v
	return s
}

type TagResourcesRequest struct {
	OwnerAccount         *string                   `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64                    `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string                   `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceId           []*string                 `json:"ResourceId,omitempty" xml:"ResourceId,omitempty" type:"Repeated"`
	ResourceOwnerAccount *string                   `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64                    `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	ResourceType         *string                   `json:"ResourceType,omitempty" xml:"ResourceType,omitempty"`
	Tag                  []*TagResourcesRequestTag `json:"Tag,omitempty" xml:"Tag,omitempty" type:"Repeated"`
}

func (s TagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesRequest) GoString() string {
	return s.String()
}

func (s *TagResourcesRequest) SetOwnerAccount(v string) *TagResourcesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *TagResourcesRequest) SetOwnerId(v int64) *TagResourcesRequest {
	s.OwnerId = &v
	return s
}

func (s *TagResourcesRequest) SetRegionId(v string) *TagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *TagResourcesRequest) SetResourceId(v []*string) *TagResourcesRequest {
	s.ResourceId = v
	return s
}

func (s *TagResourcesRequest) SetResourceOwnerAccount(v string) *TagResourcesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *TagResourcesRequest) SetResourceOwnerId(v int64) *TagResourcesRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *TagResourcesRequest) SetResourceType(v string) *TagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *TagResourcesRequest) SetTag(v []*TagResourcesRequestTag) *TagResourcesRequest {
	s.Tag = v
	return s
}

type TagResourcesRequestTag struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s TagResourcesRequestTag) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesRequestTag) GoString() string {
	return s.String()
}

func (s *TagResourcesRequestTag) SetKey(v string) *TagResourcesRequestTag {
	s.Key = &v
	return s
}

func (s *TagResourcesRequestTag) SetValue(v string) *TagResourcesRequestTag {
	s.Value = &v
	return s
}

type TagResourcesResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s TagResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *TagResourcesResponseBody) SetRequestId(v string) *TagResourcesResponseBody {
	s.RequestId = &v
	return s
}

type TagResourcesResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *TagResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s TagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesResponse) GoString() string {
	return s.String()
}

func (s *TagResourcesResponse) SetHeaders(v map[string]*string) *TagResourcesResponse {
	s.Headers = v
	return s
}

func (s *TagResourcesResponse) SetStatusCode(v int32) *TagResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *TagResourcesResponse) SetBody(v *TagResourcesResponseBody) *TagResourcesResponse {
	s.Body = v
	return s
}

type TempModifyDBNodeRequest struct {
	ClientToken          *string                          `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string                          `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNode               []*TempModifyDBNodeRequestDBNode `json:"DBNode,omitempty" xml:"DBNode,omitempty" type:"Repeated"`
	ModifyType           *string                          `json:"ModifyType,omitempty" xml:"ModifyType,omitempty"`
	OperationType        *string                          `json:"OperationType,omitempty" xml:"OperationType,omitempty"`
	OwnerAccount         *string                          `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64                           `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	ResourceOwnerAccount *string                          `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64                           `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	RestoreTime          *string                          `json:"RestoreTime,omitempty" xml:"RestoreTime,omitempty"`
}

func (s TempModifyDBNodeRequest) String() string {
	return tea.Prettify(s)
}

func (s TempModifyDBNodeRequest) GoString() string {
	return s.String()
}

func (s *TempModifyDBNodeRequest) SetClientToken(v string) *TempModifyDBNodeRequest {
	s.ClientToken = &v
	return s
}

func (s *TempModifyDBNodeRequest) SetDBClusterId(v string) *TempModifyDBNodeRequest {
	s.DBClusterId = &v
	return s
}

func (s *TempModifyDBNodeRequest) SetDBNode(v []*TempModifyDBNodeRequestDBNode) *TempModifyDBNodeRequest {
	s.DBNode = v
	return s
}

func (s *TempModifyDBNodeRequest) SetModifyType(v string) *TempModifyDBNodeRequest {
	s.ModifyType = &v
	return s
}

func (s *TempModifyDBNodeRequest) SetOperationType(v string) *TempModifyDBNodeRequest {
	s.OperationType = &v
	return s
}

func (s *TempModifyDBNodeRequest) SetOwnerAccount(v string) *TempModifyDBNodeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *TempModifyDBNodeRequest) SetOwnerId(v int64) *TempModifyDBNodeRequest {
	s.OwnerId = &v
	return s
}

func (s *TempModifyDBNodeRequest) SetResourceOwnerAccount(v string) *TempModifyDBNodeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *TempModifyDBNodeRequest) SetResourceOwnerId(v int64) *TempModifyDBNodeRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *TempModifyDBNodeRequest) SetRestoreTime(v string) *TempModifyDBNodeRequest {
	s.RestoreTime = &v
	return s
}

type TempModifyDBNodeRequestDBNode struct {
	TargetClass *string `json:"TargetClass,omitempty" xml:"TargetClass,omitempty"`
	ZoneId      *string `json:"ZoneId,omitempty" xml:"ZoneId,omitempty"`
}

func (s TempModifyDBNodeRequestDBNode) String() string {
	return tea.Prettify(s)
}

func (s TempModifyDBNodeRequestDBNode) GoString() string {
	return s.String()
}

func (s *TempModifyDBNodeRequestDBNode) SetTargetClass(v string) *TempModifyDBNodeRequestDBNode {
	s.TargetClass = &v
	return s
}

func (s *TempModifyDBNodeRequestDBNode) SetZoneId(v string) *TempModifyDBNodeRequestDBNode {
	s.ZoneId = &v
	return s
}

type TempModifyDBNodeResponseBody struct {
	DBClusterId *string   `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	DBNodeIds   []*string `json:"DBNodeIds,omitempty" xml:"DBNodeIds,omitempty" type:"Repeated"`
	OrderId     *string   `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	RequestId   *string   `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s TempModifyDBNodeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s TempModifyDBNodeResponseBody) GoString() string {
	return s.String()
}

func (s *TempModifyDBNodeResponseBody) SetDBClusterId(v string) *TempModifyDBNodeResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *TempModifyDBNodeResponseBody) SetDBNodeIds(v []*string) *TempModifyDBNodeResponseBody {
	s.DBNodeIds = v
	return s
}

func (s *TempModifyDBNodeResponseBody) SetOrderId(v string) *TempModifyDBNodeResponseBody {
	s.OrderId = &v
	return s
}

func (s *TempModifyDBNodeResponseBody) SetRequestId(v string) *TempModifyDBNodeResponseBody {
	s.RequestId = &v
	return s
}

type TempModifyDBNodeResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *TempModifyDBNodeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s TempModifyDBNodeResponse) String() string {
	return tea.Prettify(s)
}

func (s TempModifyDBNodeResponse) GoString() string {
	return s.String()
}

func (s *TempModifyDBNodeResponse) SetHeaders(v map[string]*string) *TempModifyDBNodeResponse {
	s.Headers = v
	return s
}

func (s *TempModifyDBNodeResponse) SetStatusCode(v int32) *TempModifyDBNodeResponse {
	s.StatusCode = &v
	return s
}

func (s *TempModifyDBNodeResponse) SetBody(v *TempModifyDBNodeResponseBody) *TempModifyDBNodeResponse {
	s.Body = v
	return s
}

type TransformDBClusterPayTypeRequest struct {
	ClientToken          *string `json:"ClientToken,omitempty" xml:"ClientToken,omitempty"`
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PayType              *string `json:"PayType,omitempty" xml:"PayType,omitempty"`
	Period               *string `json:"Period,omitempty" xml:"Period,omitempty"`
	RegionId             *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceGroupId      *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	UsedTime             *string `json:"UsedTime,omitempty" xml:"UsedTime,omitempty"`
}

func (s TransformDBClusterPayTypeRequest) String() string {
	return tea.Prettify(s)
}

func (s TransformDBClusterPayTypeRequest) GoString() string {
	return s.String()
}

func (s *TransformDBClusterPayTypeRequest) SetClientToken(v string) *TransformDBClusterPayTypeRequest {
	s.ClientToken = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetDBClusterId(v string) *TransformDBClusterPayTypeRequest {
	s.DBClusterId = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetOwnerAccount(v string) *TransformDBClusterPayTypeRequest {
	s.OwnerAccount = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetOwnerId(v int64) *TransformDBClusterPayTypeRequest {
	s.OwnerId = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetPayType(v string) *TransformDBClusterPayTypeRequest {
	s.PayType = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetPeriod(v string) *TransformDBClusterPayTypeRequest {
	s.Period = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetRegionId(v string) *TransformDBClusterPayTypeRequest {
	s.RegionId = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetResourceGroupId(v string) *TransformDBClusterPayTypeRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetResourceOwnerAccount(v string) *TransformDBClusterPayTypeRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetResourceOwnerId(v int64) *TransformDBClusterPayTypeRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *TransformDBClusterPayTypeRequest) SetUsedTime(v string) *TransformDBClusterPayTypeRequest {
	s.UsedTime = &v
	return s
}

type TransformDBClusterPayTypeResponseBody struct {
	ChargeType  *string `json:"ChargeType,omitempty" xml:"ChargeType,omitempty"`
	DBClusterId *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	ExpiredTime *string `json:"ExpiredTime,omitempty" xml:"ExpiredTime,omitempty"`
	OrderId     *string `json:"OrderId,omitempty" xml:"OrderId,omitempty"`
	RequestId   *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s TransformDBClusterPayTypeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s TransformDBClusterPayTypeResponseBody) GoString() string {
	return s.String()
}

func (s *TransformDBClusterPayTypeResponseBody) SetChargeType(v string) *TransformDBClusterPayTypeResponseBody {
	s.ChargeType = &v
	return s
}

func (s *TransformDBClusterPayTypeResponseBody) SetDBClusterId(v string) *TransformDBClusterPayTypeResponseBody {
	s.DBClusterId = &v
	return s
}

func (s *TransformDBClusterPayTypeResponseBody) SetExpiredTime(v string) *TransformDBClusterPayTypeResponseBody {
	s.ExpiredTime = &v
	return s
}

func (s *TransformDBClusterPayTypeResponseBody) SetOrderId(v string) *TransformDBClusterPayTypeResponseBody {
	s.OrderId = &v
	return s
}

func (s *TransformDBClusterPayTypeResponseBody) SetRequestId(v string) *TransformDBClusterPayTypeResponseBody {
	s.RequestId = &v
	return s
}

type TransformDBClusterPayTypeResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *TransformDBClusterPayTypeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s TransformDBClusterPayTypeResponse) String() string {
	return tea.Prettify(s)
}

func (s TransformDBClusterPayTypeResponse) GoString() string {
	return s.String()
}

func (s *TransformDBClusterPayTypeResponse) SetHeaders(v map[string]*string) *TransformDBClusterPayTypeResponse {
	s.Headers = v
	return s
}

func (s *TransformDBClusterPayTypeResponse) SetStatusCode(v int32) *TransformDBClusterPayTypeResponse {
	s.StatusCode = &v
	return s
}

func (s *TransformDBClusterPayTypeResponse) SetBody(v *TransformDBClusterPayTypeResponseBody) *TransformDBClusterPayTypeResponse {
	s.Body = v
	return s
}

type UntagResourcesRequest struct {
	All                  *bool     `json:"All,omitempty" xml:"All,omitempty"`
	OwnerAccount         *string   `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64    `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	RegionId             *string   `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	ResourceId           []*string `json:"ResourceId,omitempty" xml:"ResourceId,omitempty" type:"Repeated"`
	ResourceOwnerAccount *string   `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64    `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	ResourceType         *string   `json:"ResourceType,omitempty" xml:"ResourceType,omitempty"`
	TagKey               []*string `json:"TagKey,omitempty" xml:"TagKey,omitempty" type:"Repeated"`
}

func (s UntagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesRequest) GoString() string {
	return s.String()
}

func (s *UntagResourcesRequest) SetAll(v bool) *UntagResourcesRequest {
	s.All = &v
	return s
}

func (s *UntagResourcesRequest) SetOwnerAccount(v string) *UntagResourcesRequest {
	s.OwnerAccount = &v
	return s
}

func (s *UntagResourcesRequest) SetOwnerId(v int64) *UntagResourcesRequest {
	s.OwnerId = &v
	return s
}

func (s *UntagResourcesRequest) SetRegionId(v string) *UntagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *UntagResourcesRequest) SetResourceId(v []*string) *UntagResourcesRequest {
	s.ResourceId = v
	return s
}

func (s *UntagResourcesRequest) SetResourceOwnerAccount(v string) *UntagResourcesRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *UntagResourcesRequest) SetResourceOwnerId(v int64) *UntagResourcesRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *UntagResourcesRequest) SetResourceType(v string) *UntagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *UntagResourcesRequest) SetTagKey(v []*string) *UntagResourcesRequest {
	s.TagKey = v
	return s
}

type UntagResourcesResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s UntagResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *UntagResourcesResponseBody) SetRequestId(v string) *UntagResourcesResponseBody {
	s.RequestId = &v
	return s
}

type UntagResourcesResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UntagResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UntagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesResponse) GoString() string {
	return s.String()
}

func (s *UntagResourcesResponse) SetHeaders(v map[string]*string) *UntagResourcesResponse {
	s.Headers = v
	return s
}

func (s *UntagResourcesResponse) SetStatusCode(v int32) *UntagResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *UntagResourcesResponse) SetBody(v *UntagResourcesResponseBody) *UntagResourcesResponse {
	s.Body = v
	return s
}

type UpgradeDBClusterMinorVersionRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	FromTimeService      *bool   `json:"FromTimeService,omitempty" xml:"FromTimeService,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
}

func (s UpgradeDBClusterMinorVersionRequest) String() string {
	return tea.Prettify(s)
}

func (s UpgradeDBClusterMinorVersionRequest) GoString() string {
	return s.String()
}

func (s *UpgradeDBClusterMinorVersionRequest) SetDBClusterId(v string) *UpgradeDBClusterMinorVersionRequest {
	s.DBClusterId = &v
	return s
}

func (s *UpgradeDBClusterMinorVersionRequest) SetFromTimeService(v bool) *UpgradeDBClusterMinorVersionRequest {
	s.FromTimeService = &v
	return s
}

func (s *UpgradeDBClusterMinorVersionRequest) SetOwnerAccount(v string) *UpgradeDBClusterMinorVersionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *UpgradeDBClusterMinorVersionRequest) SetOwnerId(v int64) *UpgradeDBClusterMinorVersionRequest {
	s.OwnerId = &v
	return s
}

func (s *UpgradeDBClusterMinorVersionRequest) SetPlannedEndTime(v string) *UpgradeDBClusterMinorVersionRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *UpgradeDBClusterMinorVersionRequest) SetPlannedStartTime(v string) *UpgradeDBClusterMinorVersionRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *UpgradeDBClusterMinorVersionRequest) SetResourceOwnerAccount(v string) *UpgradeDBClusterMinorVersionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *UpgradeDBClusterMinorVersionRequest) SetResourceOwnerId(v int64) *UpgradeDBClusterMinorVersionRequest {
	s.ResourceOwnerId = &v
	return s
}

type UpgradeDBClusterMinorVersionResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s UpgradeDBClusterMinorVersionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpgradeDBClusterMinorVersionResponseBody) GoString() string {
	return s.String()
}

func (s *UpgradeDBClusterMinorVersionResponseBody) SetRequestId(v string) *UpgradeDBClusterMinorVersionResponseBody {
	s.RequestId = &v
	return s
}

type UpgradeDBClusterMinorVersionResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UpgradeDBClusterMinorVersionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpgradeDBClusterMinorVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s UpgradeDBClusterMinorVersionResponse) GoString() string {
	return s.String()
}

func (s *UpgradeDBClusterMinorVersionResponse) SetHeaders(v map[string]*string) *UpgradeDBClusterMinorVersionResponse {
	s.Headers = v
	return s
}

func (s *UpgradeDBClusterMinorVersionResponse) SetStatusCode(v int32) *UpgradeDBClusterMinorVersionResponse {
	s.StatusCode = &v
	return s
}

func (s *UpgradeDBClusterMinorVersionResponse) SetBody(v *UpgradeDBClusterMinorVersionResponseBody) *UpgradeDBClusterMinorVersionResponse {
	s.Body = v
	return s
}

type UpgradeDBClusterVersionRequest struct {
	DBClusterId          *string `json:"DBClusterId,omitempty" xml:"DBClusterId,omitempty"`
	FromTimeService      *bool   `json:"FromTimeService,omitempty" xml:"FromTimeService,omitempty"`
	OwnerAccount         *string `json:"OwnerAccount,omitempty" xml:"OwnerAccount,omitempty"`
	OwnerId              *int64  `json:"OwnerId,omitempty" xml:"OwnerId,omitempty"`
	PlannedEndTime       *string `json:"PlannedEndTime,omitempty" xml:"PlannedEndTime,omitempty"`
	PlannedStartTime     *string `json:"PlannedStartTime,omitempty" xml:"PlannedStartTime,omitempty"`
	ResourceOwnerAccount *string `json:"ResourceOwnerAccount,omitempty" xml:"ResourceOwnerAccount,omitempty"`
	ResourceOwnerId      *int64  `json:"ResourceOwnerId,omitempty" xml:"ResourceOwnerId,omitempty"`
	UpgradeLabel         *string `json:"UpgradeLabel,omitempty" xml:"UpgradeLabel,omitempty"`
	UpgradePolicy        *string `json:"UpgradePolicy,omitempty" xml:"UpgradePolicy,omitempty"`
	UpgradeType          *string `json:"UpgradeType,omitempty" xml:"UpgradeType,omitempty"`
}

func (s UpgradeDBClusterVersionRequest) String() string {
	return tea.Prettify(s)
}

func (s UpgradeDBClusterVersionRequest) GoString() string {
	return s.String()
}

func (s *UpgradeDBClusterVersionRequest) SetDBClusterId(v string) *UpgradeDBClusterVersionRequest {
	s.DBClusterId = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetFromTimeService(v bool) *UpgradeDBClusterVersionRequest {
	s.FromTimeService = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetOwnerAccount(v string) *UpgradeDBClusterVersionRequest {
	s.OwnerAccount = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetOwnerId(v int64) *UpgradeDBClusterVersionRequest {
	s.OwnerId = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetPlannedEndTime(v string) *UpgradeDBClusterVersionRequest {
	s.PlannedEndTime = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetPlannedStartTime(v string) *UpgradeDBClusterVersionRequest {
	s.PlannedStartTime = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetResourceOwnerAccount(v string) *UpgradeDBClusterVersionRequest {
	s.ResourceOwnerAccount = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetResourceOwnerId(v int64) *UpgradeDBClusterVersionRequest {
	s.ResourceOwnerId = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetUpgradeLabel(v string) *UpgradeDBClusterVersionRequest {
	s.UpgradeLabel = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetUpgradePolicy(v string) *UpgradeDBClusterVersionRequest {
	s.UpgradePolicy = &v
	return s
}

func (s *UpgradeDBClusterVersionRequest) SetUpgradeType(v string) *UpgradeDBClusterVersionRequest {
	s.UpgradeType = &v
	return s
}

type UpgradeDBClusterVersionResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s UpgradeDBClusterVersionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpgradeDBClusterVersionResponseBody) GoString() string {
	return s.String()
}

func (s *UpgradeDBClusterVersionResponseBody) SetRequestId(v string) *UpgradeDBClusterVersionResponseBody {
	s.RequestId = &v
	return s
}

type UpgradeDBClusterVersionResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UpgradeDBClusterVersionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpgradeDBClusterVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s UpgradeDBClusterVersionResponse) GoString() string {
	return s.String()
}

func (s *UpgradeDBClusterVersionResponse) SetHeaders(v map[string]*string) *UpgradeDBClusterVersionResponse {
	s.Headers = v
	return s
}

func (s *UpgradeDBClusterVersionResponse) SetStatusCode(v int32) *UpgradeDBClusterVersionResponse {
	s.StatusCode = &v
	return s
}

func (s *UpgradeDBClusterVersionResponse) SetBody(v *UpgradeDBClusterVersionResponseBody) *UpgradeDBClusterVersionResponse {
	s.Body = v
	return s
}

type Client struct {
	openapi.Client
}

func NewClient(config *openapi.Config) (*Client, error) {
	client := new(Client)
	err := client.Init(config)
	return client, err
}

func (client *Client) Init(config *openapi.Config) (_err error) {
	_err = client.Client.Init(config)
	if _err != nil {
		return _err
	}
	client.EndpointRule = tea.String("regional")
	client.EndpointMap = map[string]*string{
		"cn-qingdao":                  tea.String("polardb.aliyuncs.com"),
		"cn-beijing":                  tea.String("polardb.aliyuncs.com"),
		"cn-hangzhou":                 tea.String("polardb.aliyuncs.com"),
		"cn-shanghai":                 tea.String("polardb.aliyuncs.com"),
		"cn-shenzhen":                 tea.String("polardb.aliyuncs.com"),
		"cn-hongkong":                 tea.String("polardb.aliyuncs.com"),
		"ap-southeast-1":              tea.String("polardb.aliyuncs.com"),
		"us-west-1":                   tea.String("polardb.aliyuncs.com"),
		"us-east-1":                   tea.String("polardb.aliyuncs.com"),
		"cn-hangzhou-finance":         tea.String("polardb.aliyuncs.com"),
		"cn-shanghai-finance-1":       tea.String("polardb.aliyuncs.com"),
		"cn-shenzhen-finance-1":       tea.String("polardb.aliyuncs.com"),
		"ap-northeast-2-pop":          tea.String("polardb.aliyuncs.com"),
		"cn-beijing-finance-1":        tea.String("polardb.aliyuncs.com"),
		"cn-beijing-finance-pop":      tea.String("polardb.aliyuncs.com"),
		"cn-beijing-gov-1":            tea.String("polardb.aliyuncs.com"),
		"cn-beijing-nu16-b01":         tea.String("polardb.aliyuncs.com"),
		"cn-edge-1":                   tea.String("polardb.aliyuncs.com"),
		"cn-fujian":                   tea.String("polardb.aliyuncs.com"),
		"cn-haidian-cm12-c01":         tea.String("polardb.aliyuncs.com"),
		"cn-hangzhou-bj-b01":          tea.String("polardb.aliyuncs.com"),
		"cn-hangzhou-internal-prod-1": tea.String("polardb.aliyuncs.com"),
		"cn-hangzhou-internal-test-1": tea.String("polardb.aliyuncs.com"),
		"cn-hangzhou-internal-test-2": tea.String("polardb.aliyuncs.com"),
		"cn-hangzhou-internal-test-3": tea.String("polardb.aliyuncs.com"),
		"cn-hangzhou-test-306":        tea.String("polardb.aliyuncs.com"),
		"cn-hongkong-finance-pop":     tea.String("polardb.aliyuncs.com"),
		"cn-huhehaote-nebula-1":       tea.String("polardb.aliyuncs.com"),
		"cn-north-2-gov-1":            tea.String("polardb.aliyuncs.com"),
		"cn-qingdao-nebula":           tea.String("polardb.aliyuncs.com"),
		"cn-shanghai-et15-b01":        tea.String("polardb.aliyuncs.com"),
		"cn-shanghai-et2-b01":         tea.String("polardb.aliyuncs.com"),
		"cn-shanghai-inner":           tea.String("polardb.aliyuncs.com"),
		"cn-shanghai-internal-test-1": tea.String("polardb.aliyuncs.com"),
		"cn-shenzhen-inner":           tea.String("polardb.aliyuncs.com"),
		"cn-shenzhen-st4-d01":         tea.String("polardb.aliyuncs.com"),
		"cn-shenzhen-su18-b01":        tea.String("polardb.aliyuncs.com"),
		"cn-wuhan":                    tea.String("polardb.aliyuncs.com"),
		"cn-wulanchabu":               tea.String("polardb.aliyuncs.com"),
		"cn-yushanfang":               tea.String("polardb.aliyuncs.com"),
		"cn-zhangbei":                 tea.String("polardb.aliyuncs.com"),
		"cn-zhangbei-na61-b01":        tea.String("polardb.aliyuncs.com"),
		"cn-zhangjiakou-na62-a01":     tea.String("polardb.aliyuncs.com"),
		"cn-zhengzhou-nebula-1":       tea.String("polardb.aliyuncs.com"),
		"eu-west-1-oxs":               tea.String("polardb.aliyuncs.com"),
		"rus-west-1-pop":              tea.String("polardb.aliyuncs.com"),
	}
	_err = client.CheckConfig(config)
	if _err != nil {
		return _err
	}
	client.Endpoint, _err = client.GetEndpoint(tea.String("polardb"), client.RegionId, client.EndpointRule, client.Network, client.Suffix, client.EndpointMap, client.Endpoint)
	if _err != nil {
		return _err
	}

	return nil
}

func (client *Client) GetEndpoint(productId *string, regionId *string, endpointRule *string, network *string, suffix *string, endpointMap map[string]*string, endpoint *string) (_result *string, _err error) {
	if !tea.BoolValue(util.Empty(endpoint)) {
		_result = endpoint
		return _result, _err
	}

	if !tea.BoolValue(util.IsUnset(endpointMap)) && !tea.BoolValue(util.Empty(endpointMap[tea.StringValue(regionId)])) {
		_result = endpointMap[tea.StringValue(regionId)]
		return _result, _err
	}

	_body, _err := endpointutil.GetEndpointRules(productId, regionId, endpointRule, network, suffix)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelScheduleTasksWithOptions(request *CancelScheduleTasksRequest, runtime *util.RuntimeOptions) (_result *CancelScheduleTasksResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.TaskId)) {
		query["TaskId"] = request.TaskId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CancelScheduleTasks"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CancelScheduleTasksResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelScheduleTasks(request *CancelScheduleTasksRequest) (_result *CancelScheduleTasksResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CancelScheduleTasksResponse{}
	_body, _err := client.CancelScheduleTasksWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CheckAccountNameWithOptions(request *CheckAccountNameRequest, runtime *util.RuntimeOptions) (_result *CheckAccountNameResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CheckAccountName"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CheckAccountNameResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CheckAccountName(request *CheckAccountNameRequest) (_result *CheckAccountNameResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CheckAccountNameResponse{}
	_body, _err := client.CheckAccountNameWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CheckDBNameWithOptions(request *CheckDBNameRequest, runtime *util.RuntimeOptions) (_result *CheckDBNameResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CheckDBName"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CheckDBNameResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CheckDBName(request *CheckDBNameRequest) (_result *CheckDBNameResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CheckDBNameResponse{}
	_body, _err := client.CheckDBNameWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CloseAITaskWithOptions(request *CloseAITaskRequest, runtime *util.RuntimeOptions) (_result *CloseAITaskResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CloseAITask"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CloseAITaskResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CloseAITask(request *CloseAITaskRequest) (_result *CloseAITaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CloseAITaskResponse{}
	_body, _err := client.CloseAITaskWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CloseDBClusterMigrationWithOptions(request *CloseDBClusterMigrationRequest, runtime *util.RuntimeOptions) (_result *CloseDBClusterMigrationResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ContinueEnableBinlog)) {
		query["ContinueEnableBinlog"] = request.ContinueEnableBinlog
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CloseDBClusterMigration"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CloseDBClusterMigrationResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CloseDBClusterMigration(request *CloseDBClusterMigrationRequest) (_result *CloseDBClusterMigrationResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CloseDBClusterMigrationResponse{}
	_body, _err := client.CloseDBClusterMigrationWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateAccountWithOptions(request *CreateAccountRequest, runtime *util.RuntimeOptions) (_result *CreateAccountResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountDescription)) {
		query["AccountDescription"] = request.AccountDescription
	}

	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.AccountPassword)) {
		query["AccountPassword"] = request.AccountPassword
	}

	if !tea.BoolValue(util.IsUnset(request.AccountPrivilege)) {
		query["AccountPrivilege"] = request.AccountPrivilege
	}

	if !tea.BoolValue(util.IsUnset(request.AccountType)) {
		query["AccountType"] = request.AccountType
	}

	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateAccount"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateAccountResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateAccount(request *CreateAccountRequest) (_result *CreateAccountResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateAccountResponse{}
	_body, _err := client.CreateAccountWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateBackupWithOptions(request *CreateBackupRequest, runtime *util.RuntimeOptions) (_result *CreateBackupResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateBackup"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateBackupResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateBackup(request *CreateBackupRequest) (_result *CreateBackupResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateBackupResponse{}
	_body, _err := client.CreateBackupWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateDBClusterWithOptions(request *CreateDBClusterRequest, runtime *util.RuntimeOptions) (_result *CreateDBClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AutoRenew)) {
		query["AutoRenew"] = request.AutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.BackupRetentionPolicyOnClusterDeletion)) {
		query["BackupRetentionPolicyOnClusterDeletion"] = request.BackupRetentionPolicyOnClusterDeletion
	}

	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.CloneDataPoint)) {
		query["CloneDataPoint"] = request.CloneDataPoint
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterNetworkType)) {
		query["ClusterNetworkType"] = request.ClusterNetworkType
	}

	if !tea.BoolValue(util.IsUnset(request.CreationCategory)) {
		query["CreationCategory"] = request.CreationCategory
	}

	if !tea.BoolValue(util.IsUnset(request.CreationOption)) {
		query["CreationOption"] = request.CreationOption
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterDescription)) {
		query["DBClusterDescription"] = request.DBClusterDescription
	}

	if !tea.BoolValue(util.IsUnset(request.DBMinorVersion)) {
		query["DBMinorVersion"] = request.DBMinorVersion
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeClass)) {
		query["DBNodeClass"] = request.DBNodeClass
	}

	if !tea.BoolValue(util.IsUnset(request.DBType)) {
		query["DBType"] = request.DBType
	}

	if !tea.BoolValue(util.IsUnset(request.DBVersion)) {
		query["DBVersion"] = request.DBVersion
	}

	if !tea.BoolValue(util.IsUnset(request.DefaultTimeZone)) {
		query["DefaultTimeZone"] = request.DefaultTimeZone
	}

	if !tea.BoolValue(util.IsUnset(request.GDNId)) {
		query["GDNId"] = request.GDNId
	}

	if !tea.BoolValue(util.IsUnset(request.LowerCaseTableNames)) {
		query["LowerCaseTableNames"] = request.LowerCaseTableNames
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ParameterGroupId)) {
		query["ParameterGroupId"] = request.ParameterGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.PayType)) {
		query["PayType"] = request.PayType
	}

	if !tea.BoolValue(util.IsUnset(request.Period)) {
		query["Period"] = request.Period
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		query["ResourceGroupId"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityIPList)) {
		query["SecurityIPList"] = request.SecurityIPList
	}

	if !tea.BoolValue(util.IsUnset(request.SourceResourceId)) {
		query["SourceResourceId"] = request.SourceResourceId
	}

	if !tea.BoolValue(util.IsUnset(request.TDEStatus)) {
		query["TDEStatus"] = request.TDEStatus
	}

	if !tea.BoolValue(util.IsUnset(request.UsedTime)) {
		query["UsedTime"] = request.UsedTime
	}

	if !tea.BoolValue(util.IsUnset(request.VPCId)) {
		query["VPCId"] = request.VPCId
	}

	if !tea.BoolValue(util.IsUnset(request.VSwitchId)) {
		query["VSwitchId"] = request.VSwitchId
	}

	if !tea.BoolValue(util.IsUnset(request.ZoneId)) {
		query["ZoneId"] = request.ZoneId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateDBCluster"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateDBClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateDBCluster(request *CreateDBClusterRequest) (_result *CreateDBClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateDBClusterResponse{}
	_body, _err := client.CreateDBClusterWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateDBClusterEndpointWithOptions(request *CreateDBClusterEndpointRequest, runtime *util.RuntimeOptions) (_result *CreateDBClusterEndpointResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AutoAddNewNodes)) {
		query["AutoAddNewNodes"] = request.AutoAddNewNodes
	}

	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointDescription)) {
		query["DBEndpointDescription"] = request.DBEndpointDescription
	}

	if !tea.BoolValue(util.IsUnset(request.EndpointConfig)) {
		query["EndpointConfig"] = request.EndpointConfig
	}

	if !tea.BoolValue(util.IsUnset(request.EndpointType)) {
		query["EndpointType"] = request.EndpointType
	}

	if !tea.BoolValue(util.IsUnset(request.Nodes)) {
		query["Nodes"] = request.Nodes
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ReadWriteMode)) {
		query["ReadWriteMode"] = request.ReadWriteMode
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateDBClusterEndpoint"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateDBClusterEndpointResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateDBClusterEndpoint(request *CreateDBClusterEndpointRequest) (_result *CreateDBClusterEndpointResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateDBClusterEndpointResponse{}
	_body, _err := client.CreateDBClusterEndpointWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateDBEndpointAddressWithOptions(request *CreateDBEndpointAddressRequest, runtime *util.RuntimeOptions) (_result *CreateDBEndpointAddressResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ConnectionStringPrefix)) {
		query["ConnectionStringPrefix"] = request.ConnectionStringPrefix
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointId)) {
		query["DBEndpointId"] = request.DBEndpointId
	}

	if !tea.BoolValue(util.IsUnset(request.NetType)) {
		query["NetType"] = request.NetType
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateDBEndpointAddress"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateDBEndpointAddressResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateDBEndpointAddress(request *CreateDBEndpointAddressRequest) (_result *CreateDBEndpointAddressResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateDBEndpointAddressResponse{}
	_body, _err := client.CreateDBEndpointAddressWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateDBLinkWithOptions(request *CreateDBLinkRequest, runtime *util.RuntimeOptions) (_result *CreateDBLinkResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBLinkName)) {
		query["DBLinkName"] = request.DBLinkName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SourceDBName)) {
		query["SourceDBName"] = request.SourceDBName
	}

	if !tea.BoolValue(util.IsUnset(request.TargetDBAccount)) {
		query["TargetDBAccount"] = request.TargetDBAccount
	}

	if !tea.BoolValue(util.IsUnset(request.TargetDBInstanceName)) {
		query["TargetDBInstanceName"] = request.TargetDBInstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.TargetDBName)) {
		query["TargetDBName"] = request.TargetDBName
	}

	if !tea.BoolValue(util.IsUnset(request.TargetDBPasswd)) {
		query["TargetDBPasswd"] = request.TargetDBPasswd
	}

	if !tea.BoolValue(util.IsUnset(request.TargetIp)) {
		query["TargetIp"] = request.TargetIp
	}

	if !tea.BoolValue(util.IsUnset(request.TargetPort)) {
		query["TargetPort"] = request.TargetPort
	}

	if !tea.BoolValue(util.IsUnset(request.VpcId)) {
		query["VpcId"] = request.VpcId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateDBLink"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateDBLinkResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateDBLink(request *CreateDBLinkRequest) (_result *CreateDBLinkResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateDBLinkResponse{}
	_body, _err := client.CreateDBLinkWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateDBNodesWithOptions(request *CreateDBNodesRequest, runtime *util.RuntimeOptions) (_result *CreateDBNodesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNode)) {
		query["DBNode"] = request.DBNode
	}

	if !tea.BoolValue(util.IsUnset(request.EndpointBindList)) {
		query["EndpointBindList"] = request.EndpointBindList
	}

	if !tea.BoolValue(util.IsUnset(request.ImciSwitch)) {
		query["ImciSwitch"] = request.ImciSwitch
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateDBNodes"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateDBNodesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateDBNodes(request *CreateDBNodesRequest) (_result *CreateDBNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateDBNodesResponse{}
	_body, _err := client.CreateDBNodesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateDatabaseWithOptions(request *CreateDatabaseRequest, runtime *util.RuntimeOptions) (_result *CreateDatabaseResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.AccountPrivilege)) {
		query["AccountPrivilege"] = request.AccountPrivilege
	}

	if !tea.BoolValue(util.IsUnset(request.CharacterSetName)) {
		query["CharacterSetName"] = request.CharacterSetName
	}

	if !tea.BoolValue(util.IsUnset(request.Collate)) {
		query["Collate"] = request.Collate
	}

	if !tea.BoolValue(util.IsUnset(request.Ctype)) {
		query["Ctype"] = request.Ctype
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBDescription)) {
		query["DBDescription"] = request.DBDescription
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateDatabase"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateDatabaseResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateDatabase(request *CreateDatabaseRequest) (_result *CreateDatabaseResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateDatabaseResponse{}
	_body, _err := client.CreateDatabaseWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateGlobalDatabaseNetworkWithOptions(request *CreateGlobalDatabaseNetworkRequest, runtime *util.RuntimeOptions) (_result *CreateGlobalDatabaseNetworkResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.GDNDescription)) {
		query["GDNDescription"] = request.GDNDescription
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateGlobalDatabaseNetwork"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateGlobalDatabaseNetworkResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateGlobalDatabaseNetwork(request *CreateGlobalDatabaseNetworkRequest) (_result *CreateGlobalDatabaseNetworkResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateGlobalDatabaseNetworkResponse{}
	_body, _err := client.CreateGlobalDatabaseNetworkWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateParameterGroupWithOptions(request *CreateParameterGroupRequest, runtime *util.RuntimeOptions) (_result *CreateParameterGroupResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBType)) {
		query["DBType"] = request.DBType
	}

	if !tea.BoolValue(util.IsUnset(request.DBVersion)) {
		query["DBVersion"] = request.DBVersion
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ParameterGroupDesc)) {
		query["ParameterGroupDesc"] = request.ParameterGroupDesc
	}

	if !tea.BoolValue(util.IsUnset(request.ParameterGroupName)) {
		query["ParameterGroupName"] = request.ParameterGroupName
	}

	if !tea.BoolValue(util.IsUnset(request.Parameters)) {
		query["Parameters"] = request.Parameters
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateParameterGroup"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateParameterGroupResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateParameterGroup(request *CreateParameterGroupRequest) (_result *CreateParameterGroupResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateParameterGroupResponse{}
	_body, _err := client.CreateParameterGroupWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateStoragePlanWithOptions(request *CreateStoragePlanRequest, runtime *util.RuntimeOptions) (_result *CreateStoragePlanResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Period)) {
		query["Period"] = request.Period
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.StorageClass)) {
		query["StorageClass"] = request.StorageClass
	}

	if !tea.BoolValue(util.IsUnset(request.StorageType)) {
		query["StorageType"] = request.StorageType
	}

	if !tea.BoolValue(util.IsUnset(request.UsedTime)) {
		query["UsedTime"] = request.UsedTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateStoragePlan"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateStoragePlanResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateStoragePlan(request *CreateStoragePlanRequest) (_result *CreateStoragePlanResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &CreateStoragePlanResponse{}
	_body, _err := client.CreateStoragePlanWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteAccountWithOptions(request *DeleteAccountRequest, runtime *util.RuntimeOptions) (_result *DeleteAccountResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteAccount"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteAccountResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteAccount(request *DeleteAccountRequest) (_result *DeleteAccountResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteAccountResponse{}
	_body, _err := client.DeleteAccountWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteBackupWithOptions(request *DeleteBackupRequest, runtime *util.RuntimeOptions) (_result *DeleteBackupResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupId)) {
		query["BackupId"] = request.BackupId
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteBackup"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteBackupResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteBackup(request *DeleteBackupRequest) (_result *DeleteBackupResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteBackupResponse{}
	_body, _err := client.DeleteBackupWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteDBClusterWithOptions(request *DeleteDBClusterRequest, runtime *util.RuntimeOptions) (_result *DeleteDBClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupRetentionPolicyOnClusterDeletion)) {
		query["BackupRetentionPolicyOnClusterDeletion"] = request.BackupRetentionPolicyOnClusterDeletion
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteDBCluster"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteDBClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteDBCluster(request *DeleteDBClusterRequest) (_result *DeleteDBClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteDBClusterResponse{}
	_body, _err := client.DeleteDBClusterWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteDBClusterEndpointWithOptions(request *DeleteDBClusterEndpointRequest, runtime *util.RuntimeOptions) (_result *DeleteDBClusterEndpointResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointId)) {
		query["DBEndpointId"] = request.DBEndpointId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteDBClusterEndpoint"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteDBClusterEndpointResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteDBClusterEndpoint(request *DeleteDBClusterEndpointRequest) (_result *DeleteDBClusterEndpointResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteDBClusterEndpointResponse{}
	_body, _err := client.DeleteDBClusterEndpointWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteDBEndpointAddressWithOptions(request *DeleteDBEndpointAddressRequest, runtime *util.RuntimeOptions) (_result *DeleteDBEndpointAddressResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointId)) {
		query["DBEndpointId"] = request.DBEndpointId
	}

	if !tea.BoolValue(util.IsUnset(request.NetType)) {
		query["NetType"] = request.NetType
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteDBEndpointAddress"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteDBEndpointAddressResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteDBEndpointAddress(request *DeleteDBEndpointAddressRequest) (_result *DeleteDBEndpointAddressResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteDBEndpointAddressResponse{}
	_body, _err := client.DeleteDBEndpointAddressWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteDBLinkWithOptions(request *DeleteDBLinkRequest, runtime *util.RuntimeOptions) (_result *DeleteDBLinkResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBLinkName)) {
		query["DBLinkName"] = request.DBLinkName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteDBLink"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteDBLinkResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteDBLink(request *DeleteDBLinkRequest) (_result *DeleteDBLinkResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteDBLinkResponse{}
	_body, _err := client.DeleteDBLinkWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteDBNodesWithOptions(request *DeleteDBNodesRequest, runtime *util.RuntimeOptions) (_result *DeleteDBNodesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeId)) {
		query["DBNodeId"] = request.DBNodeId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteDBNodes"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteDBNodesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteDBNodes(request *DeleteDBNodesRequest) (_result *DeleteDBNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteDBNodesResponse{}
	_body, _err := client.DeleteDBNodesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteDatabaseWithOptions(request *DeleteDatabaseRequest, runtime *util.RuntimeOptions) (_result *DeleteDatabaseResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteDatabase"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteDatabaseResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteDatabase(request *DeleteDatabaseRequest) (_result *DeleteDatabaseResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteDatabaseResponse{}
	_body, _err := client.DeleteDatabaseWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteGlobalDatabaseNetworkWithOptions(request *DeleteGlobalDatabaseNetworkRequest, runtime *util.RuntimeOptions) (_result *DeleteGlobalDatabaseNetworkResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.GDNId)) {
		query["GDNId"] = request.GDNId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteGlobalDatabaseNetwork"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteGlobalDatabaseNetworkResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteGlobalDatabaseNetwork(request *DeleteGlobalDatabaseNetworkRequest) (_result *DeleteGlobalDatabaseNetworkResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteGlobalDatabaseNetworkResponse{}
	_body, _err := client.DeleteGlobalDatabaseNetworkWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteMaskingRulesWithOptions(request *DeleteMaskingRulesRequest, runtime *util.RuntimeOptions) (_result *DeleteMaskingRulesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.RuleNameList)) {
		query["RuleNameList"] = request.RuleNameList
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteMaskingRules"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteMaskingRulesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteMaskingRules(request *DeleteMaskingRulesRequest) (_result *DeleteMaskingRulesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteMaskingRulesResponse{}
	_body, _err := client.DeleteMaskingRulesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteParameterGroupWithOptions(request *DeleteParameterGroupRequest, runtime *util.RuntimeOptions) (_result *DeleteParameterGroupResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ParameterGroupId)) {
		query["ParameterGroupId"] = request.ParameterGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteParameterGroup"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteParameterGroupResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteParameterGroup(request *DeleteParameterGroupRequest) (_result *DeleteParameterGroupResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DeleteParameterGroupResponse{}
	_body, _err := client.DeleteParameterGroupWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeAITaskStatusWithOptions(request *DescribeAITaskStatusRequest, runtime *util.RuntimeOptions) (_result *DescribeAITaskStatusResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := openapiutil.Query(util.ToMap(request))
	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeAITaskStatus"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeAITaskStatusResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeAITaskStatus(request *DescribeAITaskStatusRequest) (_result *DescribeAITaskStatusResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeAITaskStatusResponse{}
	_body, _err := client.DescribeAITaskStatusWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeAccountsWithOptions(request *DescribeAccountsRequest, runtime *util.RuntimeOptions) (_result *DescribeAccountsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeAccounts"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeAccountsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeAccounts(request *DescribeAccountsRequest) (_result *DescribeAccountsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeAccountsResponse{}
	_body, _err := client.DescribeAccountsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeAutoRenewAttributeWithOptions(request *DescribeAutoRenewAttributeRequest, runtime *util.RuntimeOptions) (_result *DescribeAutoRenewAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterIds)) {
		query["DBClusterIds"] = request.DBClusterIds
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		query["ResourceGroupId"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeAutoRenewAttribute"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeAutoRenewAttributeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeAutoRenewAttribute(request *DescribeAutoRenewAttributeRequest) (_result *DescribeAutoRenewAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeAutoRenewAttributeResponse{}
	_body, _err := client.DescribeAutoRenewAttributeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeBackupLogsWithOptions(request *DescribeBackupLogsRequest, runtime *util.RuntimeOptions) (_result *DescribeBackupLogsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupRegion)) {
		query["BackupRegion"] = request.BackupRegion
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeBackupLogs"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeBackupLogsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeBackupLogs(request *DescribeBackupLogsRequest) (_result *DescribeBackupLogsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeBackupLogsResponse{}
	_body, _err := client.DescribeBackupLogsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeBackupPolicyWithOptions(request *DescribeBackupPolicyRequest, runtime *util.RuntimeOptions) (_result *DescribeBackupPolicyResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeBackupPolicy"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeBackupPolicyResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeBackupPolicy(request *DescribeBackupPolicyRequest) (_result *DescribeBackupPolicyResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeBackupPolicyResponse{}
	_body, _err := client.DescribeBackupPolicyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeBackupTasksWithOptions(request *DescribeBackupTasksRequest, runtime *util.RuntimeOptions) (_result *DescribeBackupTasksResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupJobId)) {
		query["BackupJobId"] = request.BackupJobId
	}

	if !tea.BoolValue(util.IsUnset(request.BackupMode)) {
		query["BackupMode"] = request.BackupMode
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeBackupTasks"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeBackupTasksResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeBackupTasks(request *DescribeBackupTasksRequest) (_result *DescribeBackupTasksResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeBackupTasksResponse{}
	_body, _err := client.DescribeBackupTasksWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeBackupsWithOptions(request *DescribeBackupsRequest, runtime *util.RuntimeOptions) (_result *DescribeBackupsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupId)) {
		query["BackupId"] = request.BackupId
	}

	if !tea.BoolValue(util.IsUnset(request.BackupMode)) {
		query["BackupMode"] = request.BackupMode
	}

	if !tea.BoolValue(util.IsUnset(request.BackupRegion)) {
		query["BackupRegion"] = request.BackupRegion
	}

	if !tea.BoolValue(util.IsUnset(request.BackupStatus)) {
		query["BackupStatus"] = request.BackupStatus
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeBackups"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeBackupsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeBackups(request *DescribeBackupsRequest) (_result *DescribeBackupsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeBackupsResponse{}
	_body, _err := client.DescribeBackupsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeCharacterSetNameWithOptions(request *DescribeCharacterSetNameRequest, runtime *util.RuntimeOptions) (_result *DescribeCharacterSetNameResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeCharacterSetName"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeCharacterSetNameResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeCharacterSetName(request *DescribeCharacterSetNameRequest) (_result *DescribeCharacterSetNameResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeCharacterSetNameResponse{}
	_body, _err := client.DescribeCharacterSetNameWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterAccessWhitelistWithOptions(request *DescribeDBClusterAccessWhitelistRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterAccessWhitelistResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterAccessWhitelist"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterAccessWhitelistResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterAccessWhitelist(request *DescribeDBClusterAccessWhitelistRequest) (_result *DescribeDBClusterAccessWhitelistResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterAccessWhitelistResponse{}
	_body, _err := client.DescribeDBClusterAccessWhitelistWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterAttributeWithOptions(request *DescribeDBClusterAttributeRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterAttribute"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterAttributeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterAttribute(request *DescribeDBClusterAttributeRequest) (_result *DescribeDBClusterAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterAttributeResponse{}
	_body, _err := client.DescribeDBClusterAttributeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterAuditLogCollectorWithOptions(request *DescribeDBClusterAuditLogCollectorRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterAuditLogCollectorResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterAuditLogCollector"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterAuditLogCollectorResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterAuditLogCollector(request *DescribeDBClusterAuditLogCollectorRequest) (_result *DescribeDBClusterAuditLogCollectorResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterAuditLogCollectorResponse{}
	_body, _err := client.DescribeDBClusterAuditLogCollectorWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterAvailableResourcesWithOptions(request *DescribeDBClusterAvailableResourcesRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterAvailableResourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBNodeClass)) {
		query["DBNodeClass"] = request.DBNodeClass
	}

	if !tea.BoolValue(util.IsUnset(request.DBType)) {
		query["DBType"] = request.DBType
	}

	if !tea.BoolValue(util.IsUnset(request.DBVersion)) {
		query["DBVersion"] = request.DBVersion
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PayType)) {
		query["PayType"] = request.PayType
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ZoneId)) {
		query["ZoneId"] = request.ZoneId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterAvailableResources"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterAvailableResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterAvailableResources(request *DescribeDBClusterAvailableResourcesRequest) (_result *DescribeDBClusterAvailableResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterAvailableResourcesResponse{}
	_body, _err := client.DescribeDBClusterAvailableResourcesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterEndpointsWithOptions(request *DescribeDBClusterEndpointsRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterEndpointsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointId)) {
		query["DBEndpointId"] = request.DBEndpointId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterEndpoints"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterEndpointsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterEndpoints(request *DescribeDBClusterEndpointsRequest) (_result *DescribeDBClusterEndpointsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterEndpointsResponse{}
	_body, _err := client.DescribeDBClusterEndpointsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterMigrationWithOptions(request *DescribeDBClusterMigrationRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterMigrationResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterMigration"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterMigrationResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterMigration(request *DescribeDBClusterMigrationRequest) (_result *DescribeDBClusterMigrationResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterMigrationResponse{}
	_body, _err := client.DescribeDBClusterMigrationWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterMonitorWithOptions(request *DescribeDBClusterMonitorRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterMonitorResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterMonitor"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterMonitorResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterMonitor(request *DescribeDBClusterMonitorRequest) (_result *DescribeDBClusterMonitorResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterMonitorResponse{}
	_body, _err := client.DescribeDBClusterMonitorWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterParametersWithOptions(request *DescribeDBClusterParametersRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterParametersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterParameters"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterParametersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterParameters(request *DescribeDBClusterParametersRequest) (_result *DescribeDBClusterParametersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterParametersResponse{}
	_body, _err := client.DescribeDBClusterParametersWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterPerformanceWithOptions(request *DescribeDBClusterPerformanceRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterPerformanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.Key)) {
		query["Key"] = request.Key
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterPerformance"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterPerformanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterPerformance(request *DescribeDBClusterPerformanceRequest) (_result *DescribeDBClusterPerformanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterPerformanceResponse{}
	_body, _err := client.DescribeDBClusterPerformanceWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterSSLWithOptions(request *DescribeDBClusterSSLRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterSSLResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterSSL"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterSSLResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterSSL(request *DescribeDBClusterSSLRequest) (_result *DescribeDBClusterSSLResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterSSLResponse{}
	_body, _err := client.DescribeDBClusterSSLWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterTDEWithOptions(request *DescribeDBClusterTDERequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterTDEResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterTDE"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterTDEResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterTDE(request *DescribeDBClusterTDERequest) (_result *DescribeDBClusterTDEResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterTDEResponse{}
	_body, _err := client.DescribeDBClusterTDEWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClusterVersionWithOptions(request *DescribeDBClusterVersionRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClusterVersionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusterVersion"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClusterVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusterVersion(request *DescribeDBClusterVersionRequest) (_result *DescribeDBClusterVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClusterVersionResponse{}
	_body, _err := client.DescribeDBClusterVersionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClustersWithOptions(request *DescribeDBClustersRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClustersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterDescription)) {
		query["DBClusterDescription"] = request.DBClusterDescription
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterIds)) {
		query["DBClusterIds"] = request.DBClusterIds
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterStatus)) {
		query["DBClusterStatus"] = request.DBClusterStatus
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeIds)) {
		query["DBNodeIds"] = request.DBNodeIds
	}

	if !tea.BoolValue(util.IsUnset(request.DBType)) {
		query["DBType"] = request.DBType
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.PayType)) {
		query["PayType"] = request.PayType
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		query["ResourceGroupId"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Tag)) {
		query["Tag"] = request.Tag
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClusters"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClustersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClusters(request *DescribeDBClustersRequest) (_result *DescribeDBClustersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClustersResponse{}
	_body, _err := client.DescribeDBClustersWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBClustersWithBackupsWithOptions(request *DescribeDBClustersWithBackupsRequest, runtime *util.RuntimeOptions) (_result *DescribeDBClustersWithBackupsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterDescription)) {
		query["DBClusterDescription"] = request.DBClusterDescription
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterIds)) {
		query["DBClusterIds"] = request.DBClusterIds
	}

	if !tea.BoolValue(util.IsUnset(request.DBType)) {
		query["DBType"] = request.DBType
	}

	if !tea.BoolValue(util.IsUnset(request.DBVersion)) {
		query["DBVersion"] = request.DBVersion
	}

	if !tea.BoolValue(util.IsUnset(request.IsDeleted)) {
		query["IsDeleted"] = request.IsDeleted
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBClustersWithBackups"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBClustersWithBackupsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBClustersWithBackups(request *DescribeDBClustersWithBackupsRequest) (_result *DescribeDBClustersWithBackupsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBClustersWithBackupsResponse{}
	_body, _err := client.DescribeDBClustersWithBackupsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBInitializeVariableWithOptions(request *DescribeDBInitializeVariableRequest, runtime *util.RuntimeOptions) (_result *DescribeDBInitializeVariableResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBInitializeVariable"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBInitializeVariableResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBInitializeVariable(request *DescribeDBInitializeVariableRequest) (_result *DescribeDBInitializeVariableResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBInitializeVariableResponse{}
	_body, _err := client.DescribeDBInitializeVariableWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBLinksWithOptions(request *DescribeDBLinksRequest, runtime *util.RuntimeOptions) (_result *DescribeDBLinksResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBLinkName)) {
		query["DBLinkName"] = request.DBLinkName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBLinks"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBLinksResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBLinks(request *DescribeDBLinksRequest) (_result *DescribeDBLinksResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBLinksResponse{}
	_body, _err := client.DescribeDBLinksWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBNodePerformanceWithOptions(request *DescribeDBNodePerformanceRequest, runtime *util.RuntimeOptions) (_result *DescribeDBNodePerformanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeId)) {
		query["DBNodeId"] = request.DBNodeId
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.Key)) {
		query["Key"] = request.Key
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBNodePerformance"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBNodePerformanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBNodePerformance(request *DescribeDBNodePerformanceRequest) (_result *DescribeDBNodePerformanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBNodePerformanceResponse{}
	_body, _err := client.DescribeDBNodePerformanceWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBNodesParametersWithOptions(request *DescribeDBNodesParametersRequest, runtime *util.RuntimeOptions) (_result *DescribeDBNodesParametersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeIds)) {
		query["DBNodeIds"] = request.DBNodeIds
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBNodesParameters"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBNodesParametersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBNodesParameters(request *DescribeDBNodesParametersRequest) (_result *DescribeDBNodesParametersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBNodesParametersResponse{}
	_body, _err := client.DescribeDBNodesParametersWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDBProxyPerformanceWithOptions(request *DescribeDBProxyPerformanceRequest, runtime *util.RuntimeOptions) (_result *DescribeDBProxyPerformanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.Key)) {
		query["Key"] = request.Key
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDBProxyPerformance"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDBProxyPerformanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDBProxyPerformance(request *DescribeDBProxyPerformanceRequest) (_result *DescribeDBProxyPerformanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDBProxyPerformanceResponse{}
	_body, _err := client.DescribeDBProxyPerformanceWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDatabasesWithOptions(request *DescribeDatabasesRequest, runtime *util.RuntimeOptions) (_result *DescribeDatabasesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDatabases"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDatabasesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDatabases(request *DescribeDatabasesRequest) (_result *DescribeDatabasesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDatabasesResponse{}
	_body, _err := client.DescribeDatabasesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeDetachedBackupsWithOptions(request *DescribeDetachedBackupsRequest, runtime *util.RuntimeOptions) (_result *DescribeDetachedBackupsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupId)) {
		query["BackupId"] = request.BackupId
	}

	if !tea.BoolValue(util.IsUnset(request.BackupMode)) {
		query["BackupMode"] = request.BackupMode
	}

	if !tea.BoolValue(util.IsUnset(request.BackupRegion)) {
		query["BackupRegion"] = request.BackupRegion
	}

	if !tea.BoolValue(util.IsUnset(request.BackupStatus)) {
		query["BackupStatus"] = request.BackupStatus
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeDetachedBackups"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeDetachedBackupsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeDetachedBackups(request *DescribeDetachedBackupsRequest) (_result *DescribeDetachedBackupsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeDetachedBackupsResponse{}
	_body, _err := client.DescribeDetachedBackupsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeGlobalDatabaseNetworkWithOptions(request *DescribeGlobalDatabaseNetworkRequest, runtime *util.RuntimeOptions) (_result *DescribeGlobalDatabaseNetworkResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.GDNId)) {
		query["GDNId"] = request.GDNId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeGlobalDatabaseNetwork"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeGlobalDatabaseNetworkResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeGlobalDatabaseNetwork(request *DescribeGlobalDatabaseNetworkRequest) (_result *DescribeGlobalDatabaseNetworkResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeGlobalDatabaseNetworkResponse{}
	_body, _err := client.DescribeGlobalDatabaseNetworkWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeGlobalDatabaseNetworksWithOptions(request *DescribeGlobalDatabaseNetworksRequest, runtime *util.RuntimeOptions) (_result *DescribeGlobalDatabaseNetworksResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.GDNDescription)) {
		query["GDNDescription"] = request.GDNDescription
	}

	if !tea.BoolValue(util.IsUnset(request.GDNId)) {
		query["GDNId"] = request.GDNId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeGlobalDatabaseNetworks"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeGlobalDatabaseNetworksResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeGlobalDatabaseNetworks(request *DescribeGlobalDatabaseNetworksRequest) (_result *DescribeGlobalDatabaseNetworksResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeGlobalDatabaseNetworksResponse{}
	_body, _err := client.DescribeGlobalDatabaseNetworksWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeLogBackupPolicyWithOptions(request *DescribeLogBackupPolicyRequest, runtime *util.RuntimeOptions) (_result *DescribeLogBackupPolicyResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeLogBackupPolicy"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeLogBackupPolicyResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeLogBackupPolicy(request *DescribeLogBackupPolicyRequest) (_result *DescribeLogBackupPolicyResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeLogBackupPolicyResponse{}
	_body, _err := client.DescribeLogBackupPolicyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeMaskingRulesWithOptions(request *DescribeMaskingRulesRequest, runtime *util.RuntimeOptions) (_result *DescribeMaskingRulesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.RuleNameList)) {
		query["RuleNameList"] = request.RuleNameList
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeMaskingRules"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeMaskingRulesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeMaskingRules(request *DescribeMaskingRulesRequest) (_result *DescribeMaskingRulesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeMaskingRulesResponse{}
	_body, _err := client.DescribeMaskingRulesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeMetaListWithOptions(request *DescribeMetaListRequest, runtime *util.RuntimeOptions) (_result *DescribeMetaListResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupId)) {
		query["BackupId"] = request.BackupId
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.GetDbName)) {
		query["GetDbName"] = request.GetDbName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RestoreTime)) {
		query["RestoreTime"] = request.RestoreTime
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeMetaList"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeMetaListResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeMetaList(request *DescribeMetaListRequest) (_result *DescribeMetaListResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeMetaListResponse{}
	_body, _err := client.DescribeMetaListWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeParameterGroupWithOptions(request *DescribeParameterGroupRequest, runtime *util.RuntimeOptions) (_result *DescribeParameterGroupResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ParameterGroupId)) {
		query["ParameterGroupId"] = request.ParameterGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeParameterGroup"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeParameterGroupResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeParameterGroup(request *DescribeParameterGroupRequest) (_result *DescribeParameterGroupResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeParameterGroupResponse{}
	_body, _err := client.DescribeParameterGroupWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeParameterGroupsWithOptions(request *DescribeParameterGroupsRequest, runtime *util.RuntimeOptions) (_result *DescribeParameterGroupsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBType)) {
		query["DBType"] = request.DBType
	}

	if !tea.BoolValue(util.IsUnset(request.DBVersion)) {
		query["DBVersion"] = request.DBVersion
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeParameterGroups"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeParameterGroupsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeParameterGroups(request *DescribeParameterGroupsRequest) (_result *DescribeParameterGroupsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeParameterGroupsResponse{}
	_body, _err := client.DescribeParameterGroupsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeParameterTemplatesWithOptions(request *DescribeParameterTemplatesRequest, runtime *util.RuntimeOptions) (_result *DescribeParameterTemplatesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBType)) {
		query["DBType"] = request.DBType
	}

	if !tea.BoolValue(util.IsUnset(request.DBVersion)) {
		query["DBVersion"] = request.DBVersion
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeParameterTemplates"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeParameterTemplatesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeParameterTemplates(request *DescribeParameterTemplatesRequest) (_result *DescribeParameterTemplatesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeParameterTemplatesResponse{}
	_body, _err := client.DescribeParameterTemplatesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribePendingMaintenanceActionWithOptions(request *DescribePendingMaintenanceActionRequest, runtime *util.RuntimeOptions) (_result *DescribePendingMaintenanceActionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.IsHistory)) {
		query["IsHistory"] = request.IsHistory
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.Region)) {
		query["Region"] = request.Region
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	if !tea.BoolValue(util.IsUnset(request.TaskType)) {
		query["TaskType"] = request.TaskType
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribePendingMaintenanceAction"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribePendingMaintenanceActionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribePendingMaintenanceAction(request *DescribePendingMaintenanceActionRequest) (_result *DescribePendingMaintenanceActionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribePendingMaintenanceActionResponse{}
	_body, _err := client.DescribePendingMaintenanceActionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribePendingMaintenanceActionsWithOptions(request *DescribePendingMaintenanceActionsRequest, runtime *util.RuntimeOptions) (_result *DescribePendingMaintenanceActionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.IsHistory)) {
		query["IsHistory"] = request.IsHistory
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribePendingMaintenanceActions"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribePendingMaintenanceActionsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribePendingMaintenanceActions(request *DescribePendingMaintenanceActionsRequest) (_result *DescribePendingMaintenanceActionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribePendingMaintenanceActionsResponse{}
	_body, _err := client.DescribePendingMaintenanceActionsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribePolarSQLCollectorPolicyWithOptions(request *DescribePolarSQLCollectorPolicyRequest, runtime *util.RuntimeOptions) (_result *DescribePolarSQLCollectorPolicyResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := openapiutil.Query(util.ToMap(request))
	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribePolarSQLCollectorPolicy"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribePolarSQLCollectorPolicyResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribePolarSQLCollectorPolicy(request *DescribePolarSQLCollectorPolicyRequest) (_result *DescribePolarSQLCollectorPolicyResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribePolarSQLCollectorPolicyResponse{}
	_body, _err := client.DescribePolarSQLCollectorPolicyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeRegionsWithOptions(request *DescribeRegionsRequest, runtime *util.RuntimeOptions) (_result *DescribeRegionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeRegions"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeRegionsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeRegions(request *DescribeRegionsRequest) (_result *DescribeRegionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeRegionsResponse{}
	_body, _err := client.DescribeRegionsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeScheduleTasksWithOptions(request *DescribeScheduleTasksRequest, runtime *util.RuntimeOptions) (_result *DescribeScheduleTasksResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterDescription)) {
		query["DBClusterDescription"] = request.DBClusterDescription
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OrderId)) {
		query["OrderId"] = request.OrderId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Status)) {
		query["Status"] = request.Status
	}

	if !tea.BoolValue(util.IsUnset(request.TaskAction)) {
		query["TaskAction"] = request.TaskAction
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeScheduleTasks"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeScheduleTasksResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeScheduleTasks(request *DescribeScheduleTasksRequest) (_result *DescribeScheduleTasksResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeScheduleTasksResponse{}
	_body, _err := client.DescribeScheduleTasksWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeSlowLogRecordsWithOptions(request *DescribeSlowLogRecordsRequest, runtime *util.RuntimeOptions) (_result *DescribeSlowLogRecordsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SQLHASH)) {
		query["SQLHASH"] = request.SQLHASH
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeSlowLogRecords"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeSlowLogRecordsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeSlowLogRecords(request *DescribeSlowLogRecordsRequest) (_result *DescribeSlowLogRecordsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeSlowLogRecordsResponse{}
	_body, _err := client.DescribeSlowLogRecordsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeSlowLogsWithOptions(request *DescribeSlowLogsRequest, runtime *util.RuntimeOptions) (_result *DescribeSlowLogsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeSlowLogs"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeSlowLogsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeSlowLogs(request *DescribeSlowLogsRequest) (_result *DescribeSlowLogsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeSlowLogsResponse{}
	_body, _err := client.DescribeSlowLogsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeStoragePlanWithOptions(request *DescribeStoragePlanRequest, runtime *util.RuntimeOptions) (_result *DescribeStoragePlanResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeStoragePlan"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeStoragePlanResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeStoragePlan(request *DescribeStoragePlanRequest) (_result *DescribeStoragePlanResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeStoragePlanResponse{}
	_body, _err := client.DescribeStoragePlanWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeTasksWithOptions(request *DescribeTasksRequest, runtime *util.RuntimeOptions) (_result *DescribeTasksResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeId)) {
		query["DBNodeId"] = request.DBNodeId
	}

	if !tea.BoolValue(util.IsUnset(request.EndTime)) {
		query["EndTime"] = request.EndTime
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["PageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["PageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.StartTime)) {
		query["StartTime"] = request.StartTime
	}

	if !tea.BoolValue(util.IsUnset(request.Status)) {
		query["Status"] = request.Status
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeTasks"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeTasksResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeTasks(request *DescribeTasksRequest) (_result *DescribeTasksResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &DescribeTasksResponse{}
	_body, _err := client.DescribeTasksWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) EnableFirewallRulesWithOptions(request *EnableFirewallRulesRequest, runtime *util.RuntimeOptions) (_result *EnableFirewallRulesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.Enable)) {
		query["Enable"] = request.Enable
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RuleNameList)) {
		query["RuleNameList"] = request.RuleNameList
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("EnableFirewallRules"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &EnableFirewallRulesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) EnableFirewallRules(request *EnableFirewallRulesRequest) (_result *EnableFirewallRulesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &EnableFirewallRulesResponse{}
	_body, _err := client.EnableFirewallRulesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) FailoverDBClusterWithOptions(request *FailoverDBClusterRequest, runtime *util.RuntimeOptions) (_result *FailoverDBClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.TargetDBNodeId)) {
		query["TargetDBNodeId"] = request.TargetDBNodeId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("FailoverDBCluster"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &FailoverDBClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) FailoverDBCluster(request *FailoverDBClusterRequest) (_result *FailoverDBClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &FailoverDBClusterResponse{}
	_body, _err := client.FailoverDBClusterWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GrantAccountPrivilegeWithOptions(request *GrantAccountPrivilegeRequest, runtime *util.RuntimeOptions) (_result *GrantAccountPrivilegeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.AccountPrivilege)) {
		query["AccountPrivilege"] = request.AccountPrivilege
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GrantAccountPrivilege"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &GrantAccountPrivilegeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GrantAccountPrivilege(request *GrantAccountPrivilegeRequest) (_result *GrantAccountPrivilegeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &GrantAccountPrivilegeResponse{}
	_body, _err := client.GrantAccountPrivilegeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListTagResourcesWithOptions(request *ListTagResourcesRequest, runtime *util.RuntimeOptions) (_result *ListTagResourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["NextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceId)) {
		query["ResourceId"] = request.ResourceId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		query["ResourceType"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.Tag)) {
		query["Tag"] = request.Tag
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListTagResources"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ListTagResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListTagResources(request *ListTagResourcesRequest) (_result *ListTagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ListTagResourcesResponse{}
	_body, _err := client.ListTagResourcesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyAccountDescriptionWithOptions(request *ModifyAccountDescriptionRequest, runtime *util.RuntimeOptions) (_result *ModifyAccountDescriptionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountDescription)) {
		query["AccountDescription"] = request.AccountDescription
	}

	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyAccountDescription"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyAccountDescriptionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyAccountDescription(request *ModifyAccountDescriptionRequest) (_result *ModifyAccountDescriptionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyAccountDescriptionResponse{}
	_body, _err := client.ModifyAccountDescriptionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyAccountPasswordWithOptions(request *ModifyAccountPasswordRequest, runtime *util.RuntimeOptions) (_result *ModifyAccountPasswordResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.NewAccountPassword)) {
		query["NewAccountPassword"] = request.NewAccountPassword
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyAccountPassword"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyAccountPasswordResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyAccountPassword(request *ModifyAccountPasswordRequest) (_result *ModifyAccountPasswordResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyAccountPasswordResponse{}
	_body, _err := client.ModifyAccountPasswordWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyAutoRenewAttributeWithOptions(request *ModifyAutoRenewAttributeRequest, runtime *util.RuntimeOptions) (_result *ModifyAutoRenewAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterIds)) {
		query["DBClusterIds"] = request.DBClusterIds
	}

	if !tea.BoolValue(util.IsUnset(request.Duration)) {
		query["Duration"] = request.Duration
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PeriodUnit)) {
		query["PeriodUnit"] = request.PeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.RenewalStatus)) {
		query["RenewalStatus"] = request.RenewalStatus
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		query["ResourceGroupId"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyAutoRenewAttribute"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyAutoRenewAttributeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyAutoRenewAttribute(request *ModifyAutoRenewAttributeRequest) (_result *ModifyAutoRenewAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyAutoRenewAttributeResponse{}
	_body, _err := client.ModifyAutoRenewAttributeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyBackupPolicyWithOptions(request *ModifyBackupPolicyRequest, runtime *util.RuntimeOptions) (_result *ModifyBackupPolicyResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupFrequency)) {
		query["BackupFrequency"] = request.BackupFrequency
	}

	if !tea.BoolValue(util.IsUnset(request.BackupRetentionPolicyOnClusterDeletion)) {
		query["BackupRetentionPolicyOnClusterDeletion"] = request.BackupRetentionPolicyOnClusterDeletion
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DataLevel1BackupFrequency)) {
		query["DataLevel1BackupFrequency"] = request.DataLevel1BackupFrequency
	}

	if !tea.BoolValue(util.IsUnset(request.DataLevel1BackupPeriod)) {
		query["DataLevel1BackupPeriod"] = request.DataLevel1BackupPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.DataLevel1BackupRetentionPeriod)) {
		query["DataLevel1BackupRetentionPeriod"] = request.DataLevel1BackupRetentionPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.DataLevel1BackupTime)) {
		query["DataLevel1BackupTime"] = request.DataLevel1BackupTime
	}

	if !tea.BoolValue(util.IsUnset(request.DataLevel2BackupAnotherRegionRegion)) {
		query["DataLevel2BackupAnotherRegionRegion"] = request.DataLevel2BackupAnotherRegionRegion
	}

	if !tea.BoolValue(util.IsUnset(request.DataLevel2BackupAnotherRegionRetentionPeriod)) {
		query["DataLevel2BackupAnotherRegionRetentionPeriod"] = request.DataLevel2BackupAnotherRegionRetentionPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.DataLevel2BackupPeriod)) {
		query["DataLevel2BackupPeriod"] = request.DataLevel2BackupPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.DataLevel2BackupRetentionPeriod)) {
		query["DataLevel2BackupRetentionPeriod"] = request.DataLevel2BackupRetentionPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PreferredBackupPeriod)) {
		query["PreferredBackupPeriod"] = request.PreferredBackupPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.PreferredBackupTime)) {
		query["PreferredBackupTime"] = request.PreferredBackupTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyBackupPolicy"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyBackupPolicyResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyBackupPolicy(request *ModifyBackupPolicyRequest) (_result *ModifyBackupPolicyResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyBackupPolicyResponse{}
	_body, _err := client.ModifyBackupPolicyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterAccessWhitelistWithOptions(request *ModifyDBClusterAccessWhitelistRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterAccessWhitelistResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterIPArrayAttribute)) {
		query["DBClusterIPArrayAttribute"] = request.DBClusterIPArrayAttribute
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterIPArrayName)) {
		query["DBClusterIPArrayName"] = request.DBClusterIPArrayName
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.ModifyMode)) {
		query["ModifyMode"] = request.ModifyMode
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityGroupIds)) {
		query["SecurityGroupIds"] = request.SecurityGroupIds
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityIps)) {
		query["SecurityIps"] = request.SecurityIps
	}

	if !tea.BoolValue(util.IsUnset(request.WhiteListType)) {
		query["WhiteListType"] = request.WhiteListType
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterAccessWhitelist"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterAccessWhitelistResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterAccessWhitelist(request *ModifyDBClusterAccessWhitelistRequest) (_result *ModifyDBClusterAccessWhitelistResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterAccessWhitelistResponse{}
	_body, _err := client.ModifyDBClusterAccessWhitelistWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterAndNodesParametersWithOptions(request *ModifyDBClusterAndNodesParametersRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterAndNodesParametersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeIds)) {
		query["DBNodeIds"] = request.DBNodeIds
	}

	if !tea.BoolValue(util.IsUnset(request.FromTimeService)) {
		query["FromTimeService"] = request.FromTimeService
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ParameterGroupId)) {
		query["ParameterGroupId"] = request.ParameterGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.Parameters)) {
		query["Parameters"] = request.Parameters
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterAndNodesParameters"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterAndNodesParametersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterAndNodesParameters(request *ModifyDBClusterAndNodesParametersRequest) (_result *ModifyDBClusterAndNodesParametersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterAndNodesParametersResponse{}
	_body, _err := client.ModifyDBClusterAndNodesParametersWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterAuditLogCollectorWithOptions(request *ModifyDBClusterAuditLogCollectorRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterAuditLogCollectorResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CollectorStatus)) {
		query["CollectorStatus"] = request.CollectorStatus
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterAuditLogCollector"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterAuditLogCollectorResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterAuditLogCollector(request *ModifyDBClusterAuditLogCollectorRequest) (_result *ModifyDBClusterAuditLogCollectorResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterAuditLogCollectorResponse{}
	_body, _err := client.ModifyDBClusterAuditLogCollectorWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterDeletionWithOptions(request *ModifyDBClusterDeletionRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterDeletionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Protection)) {
		query["Protection"] = request.Protection
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterDeletion"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterDeletionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterDeletion(request *ModifyDBClusterDeletionRequest) (_result *ModifyDBClusterDeletionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterDeletionResponse{}
	_body, _err := client.ModifyDBClusterDeletionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterDescriptionWithOptions(request *ModifyDBClusterDescriptionRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterDescriptionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterDescription)) {
		query["DBClusterDescription"] = request.DBClusterDescription
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterDescription"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterDescriptionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterDescription(request *ModifyDBClusterDescriptionRequest) (_result *ModifyDBClusterDescriptionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterDescriptionResponse{}
	_body, _err := client.ModifyDBClusterDescriptionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterEndpointWithOptions(request *ModifyDBClusterEndpointRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterEndpointResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AutoAddNewNodes)) {
		query["AutoAddNewNodes"] = request.AutoAddNewNodes
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointDescription)) {
		query["DBEndpointDescription"] = request.DBEndpointDescription
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointId)) {
		query["DBEndpointId"] = request.DBEndpointId
	}

	if !tea.BoolValue(util.IsUnset(request.EndpointConfig)) {
		query["EndpointConfig"] = request.EndpointConfig
	}

	if !tea.BoolValue(util.IsUnset(request.Nodes)) {
		query["Nodes"] = request.Nodes
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ReadWriteMode)) {
		query["ReadWriteMode"] = request.ReadWriteMode
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterEndpoint"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterEndpointResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterEndpoint(request *ModifyDBClusterEndpointRequest) (_result *ModifyDBClusterEndpointResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterEndpointResponse{}
	_body, _err := client.ModifyDBClusterEndpointWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterMaintainTimeWithOptions(request *ModifyDBClusterMaintainTimeRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterMaintainTimeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.MaintainTime)) {
		query["MaintainTime"] = request.MaintainTime
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterMaintainTime"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterMaintainTimeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterMaintainTime(request *ModifyDBClusterMaintainTimeRequest) (_result *ModifyDBClusterMaintainTimeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterMaintainTimeResponse{}
	_body, _err := client.ModifyDBClusterMaintainTimeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterMigrationWithOptions(request *ModifyDBClusterMigrationRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterMigrationResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ConnectionStrings)) {
		query["ConnectionStrings"] = request.ConnectionStrings
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.NewMasterInstanceId)) {
		query["NewMasterInstanceId"] = request.NewMasterInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	if !tea.BoolValue(util.IsUnset(request.SourceRDSDBInstanceId)) {
		query["SourceRDSDBInstanceId"] = request.SourceRDSDBInstanceId
	}

	if !tea.BoolValue(util.IsUnset(request.SwapConnectionString)) {
		query["SwapConnectionString"] = request.SwapConnectionString
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterMigration"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterMigrationResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterMigration(request *ModifyDBClusterMigrationRequest) (_result *ModifyDBClusterMigrationResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterMigrationResponse{}
	_body, _err := client.ModifyDBClusterMigrationWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterMonitorWithOptions(request *ModifyDBClusterMonitorRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterMonitorResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Period)) {
		query["Period"] = request.Period
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterMonitor"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterMonitorResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterMonitor(request *ModifyDBClusterMonitorRequest) (_result *ModifyDBClusterMonitorResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterMonitorResponse{}
	_body, _err := client.ModifyDBClusterMonitorWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterParametersWithOptions(request *ModifyDBClusterParametersRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterParametersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.FromTimeService)) {
		query["FromTimeService"] = request.FromTimeService
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ParameterGroupId)) {
		query["ParameterGroupId"] = request.ParameterGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.Parameters)) {
		query["Parameters"] = request.Parameters
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterParameters"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterParametersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterParameters(request *ModifyDBClusterParametersRequest) (_result *ModifyDBClusterParametersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterParametersResponse{}
	_body, _err := client.ModifyDBClusterParametersWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterPrimaryZoneWithOptions(request *ModifyDBClusterPrimaryZoneRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterPrimaryZoneResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.FromTimeService)) {
		query["FromTimeService"] = request.FromTimeService
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.VSwitchId)) {
		query["VSwitchId"] = request.VSwitchId
	}

	if !tea.BoolValue(util.IsUnset(request.ZoneId)) {
		query["ZoneId"] = request.ZoneId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterPrimaryZone"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterPrimaryZoneResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterPrimaryZone(request *ModifyDBClusterPrimaryZoneRequest) (_result *ModifyDBClusterPrimaryZoneResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterPrimaryZoneResponse{}
	_body, _err := client.ModifyDBClusterPrimaryZoneWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterResourceGroupWithOptions(request *ModifyDBClusterResourceGroupRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterResourceGroupResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.NewResourceGroupId)) {
		query["NewResourceGroupId"] = request.NewResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		query["ResourceGroupId"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterResourceGroup"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterResourceGroupResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterResourceGroup(request *ModifyDBClusterResourceGroupRequest) (_result *ModifyDBClusterResourceGroupResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterResourceGroupResponse{}
	_body, _err := client.ModifyDBClusterResourceGroupWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterSSLWithOptions(request *ModifyDBClusterSSLRequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterSSLResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointId)) {
		query["DBEndpointId"] = request.DBEndpointId
	}

	if !tea.BoolValue(util.IsUnset(request.NetType)) {
		query["NetType"] = request.NetType
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SSLAutoRotate)) {
		query["SSLAutoRotate"] = request.SSLAutoRotate
	}

	if !tea.BoolValue(util.IsUnset(request.SSLEnabled)) {
		query["SSLEnabled"] = request.SSLEnabled
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterSSL"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterSSLResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterSSL(request *ModifyDBClusterSSLRequest) (_result *ModifyDBClusterSSLResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterSSLResponse{}
	_body, _err := client.ModifyDBClusterSSLWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBClusterTDEWithOptions(request *ModifyDBClusterTDERequest, runtime *util.RuntimeOptions) (_result *ModifyDBClusterTDEResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.EncryptNewTables)) {
		query["EncryptNewTables"] = request.EncryptNewTables
	}

	if !tea.BoolValue(util.IsUnset(request.EncryptionKey)) {
		query["EncryptionKey"] = request.EncryptionKey
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RoleArn)) {
		query["RoleArn"] = request.RoleArn
	}

	if !tea.BoolValue(util.IsUnset(request.TDEStatus)) {
		query["TDEStatus"] = request.TDEStatus
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBClusterTDE"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBClusterTDEResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBClusterTDE(request *ModifyDBClusterTDERequest) (_result *ModifyDBClusterTDEResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBClusterTDEResponse{}
	_body, _err := client.ModifyDBClusterTDEWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBDescriptionWithOptions(request *ModifyDBDescriptionRequest, runtime *util.RuntimeOptions) (_result *ModifyDBDescriptionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBDescription)) {
		query["DBDescription"] = request.DBDescription
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBDescription"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBDescriptionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBDescription(request *ModifyDBDescriptionRequest) (_result *ModifyDBDescriptionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBDescriptionResponse{}
	_body, _err := client.ModifyDBDescriptionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBEndpointAddressWithOptions(request *ModifyDBEndpointAddressRequest, runtime *util.RuntimeOptions) (_result *ModifyDBEndpointAddressResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ConnectionStringPrefix)) {
		query["ConnectionStringPrefix"] = request.ConnectionStringPrefix
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBEndpointId)) {
		query["DBEndpointId"] = request.DBEndpointId
	}

	if !tea.BoolValue(util.IsUnset(request.NetType)) {
		query["NetType"] = request.NetType
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Port)) {
		query["Port"] = request.Port
	}

	if !tea.BoolValue(util.IsUnset(request.PrivateZoneAddressPrefix)) {
		query["PrivateZoneAddressPrefix"] = request.PrivateZoneAddressPrefix
	}

	if !tea.BoolValue(util.IsUnset(request.PrivateZoneName)) {
		query["PrivateZoneName"] = request.PrivateZoneName
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBEndpointAddress"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBEndpointAddressResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBEndpointAddress(request *ModifyDBEndpointAddressRequest) (_result *ModifyDBEndpointAddressResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBEndpointAddressResponse{}
	_body, _err := client.ModifyDBEndpointAddressWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBNodeClassWithOptions(request *ModifyDBNodeClassRequest, runtime *util.RuntimeOptions) (_result *ModifyDBNodeClassResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeTargetClass)) {
		query["DBNodeTargetClass"] = request.DBNodeTargetClass
	}

	if !tea.BoolValue(util.IsUnset(request.ModifyType)) {
		query["ModifyType"] = request.ModifyType
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SubCategory)) {
		query["SubCategory"] = request.SubCategory
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBNodeClass"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBNodeClassResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBNodeClass(request *ModifyDBNodeClassRequest) (_result *ModifyDBNodeClassResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBNodeClassResponse{}
	_body, _err := client.ModifyDBNodeClassWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBNodesClassWithOptions(request *ModifyDBNodesClassRequest, runtime *util.RuntimeOptions) (_result *ModifyDBNodesClassResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNode)) {
		query["DBNode"] = request.DBNode
	}

	if !tea.BoolValue(util.IsUnset(request.ModifyType)) {
		query["ModifyType"] = request.ModifyType
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SubCategory)) {
		query["SubCategory"] = request.SubCategory
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBNodesClass"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBNodesClassResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBNodesClass(request *ModifyDBNodesClassRequest) (_result *ModifyDBNodesClassResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBNodesClassResponse{}
	_body, _err := client.ModifyDBNodesClassWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyDBNodesParametersWithOptions(request *ModifyDBNodesParametersRequest, runtime *util.RuntimeOptions) (_result *ModifyDBNodesParametersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNodeIds)) {
		query["DBNodeIds"] = request.DBNodeIds
	}

	if !tea.BoolValue(util.IsUnset(request.FromTimeService)) {
		query["FromTimeService"] = request.FromTimeService
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ParameterGroupId)) {
		query["ParameterGroupId"] = request.ParameterGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.Parameters)) {
		query["Parameters"] = request.Parameters
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyDBNodesParameters"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyDBNodesParametersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyDBNodesParameters(request *ModifyDBNodesParametersRequest) (_result *ModifyDBNodesParametersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyDBNodesParametersResponse{}
	_body, _err := client.ModifyDBNodesParametersWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyGlobalDatabaseNetworkWithOptions(request *ModifyGlobalDatabaseNetworkRequest, runtime *util.RuntimeOptions) (_result *ModifyGlobalDatabaseNetworkResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.GDNDescription)) {
		query["GDNDescription"] = request.GDNDescription
	}

	if !tea.BoolValue(util.IsUnset(request.GDNId)) {
		query["GDNId"] = request.GDNId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyGlobalDatabaseNetwork"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyGlobalDatabaseNetworkResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyGlobalDatabaseNetwork(request *ModifyGlobalDatabaseNetworkRequest) (_result *ModifyGlobalDatabaseNetworkResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyGlobalDatabaseNetworkResponse{}
	_body, _err := client.ModifyGlobalDatabaseNetworkWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyLogBackupPolicyWithOptions(request *ModifyLogBackupPolicyRequest, runtime *util.RuntimeOptions) (_result *ModifyLogBackupPolicyResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.LogBackupAnotherRegionRegion)) {
		query["LogBackupAnotherRegionRegion"] = request.LogBackupAnotherRegionRegion
	}

	if !tea.BoolValue(util.IsUnset(request.LogBackupAnotherRegionRetentionPeriod)) {
		query["LogBackupAnotherRegionRetentionPeriod"] = request.LogBackupAnotherRegionRetentionPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.LogBackupRetentionPeriod)) {
		query["LogBackupRetentionPeriod"] = request.LogBackupRetentionPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyLogBackupPolicy"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyLogBackupPolicyResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyLogBackupPolicy(request *ModifyLogBackupPolicyRequest) (_result *ModifyLogBackupPolicyResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyLogBackupPolicyResponse{}
	_body, _err := client.ModifyLogBackupPolicyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyMaskingRulesWithOptions(request *ModifyMaskingRulesRequest, runtime *util.RuntimeOptions) (_result *ModifyMaskingRulesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.Enable)) {
		query["Enable"] = request.Enable
	}

	if !tea.BoolValue(util.IsUnset(request.RuleConfig)) {
		query["RuleConfig"] = request.RuleConfig
	}

	if !tea.BoolValue(util.IsUnset(request.RuleName)) {
		query["RuleName"] = request.RuleName
	}

	if !tea.BoolValue(util.IsUnset(request.RuleNameList)) {
		query["RuleNameList"] = request.RuleNameList
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyMaskingRules"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyMaskingRulesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyMaskingRules(request *ModifyMaskingRulesRequest) (_result *ModifyMaskingRulesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyMaskingRulesResponse{}
	_body, _err := client.ModifyMaskingRulesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyPendingMaintenanceActionWithOptions(request *ModifyPendingMaintenanceActionRequest, runtime *util.RuntimeOptions) (_result *ModifyPendingMaintenanceActionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Ids)) {
		query["Ids"] = request.Ids
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	if !tea.BoolValue(util.IsUnset(request.SwitchTime)) {
		query["SwitchTime"] = request.SwitchTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyPendingMaintenanceAction"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyPendingMaintenanceActionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyPendingMaintenanceAction(request *ModifyPendingMaintenanceActionRequest) (_result *ModifyPendingMaintenanceActionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ModifyPendingMaintenanceActionResponse{}
	_body, _err := client.ModifyPendingMaintenanceActionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) OpenAITaskWithOptions(request *OpenAITaskRequest, runtime *util.RuntimeOptions) (_result *OpenAITaskResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Password)) {
		query["Password"] = request.Password
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.Username)) {
		query["Username"] = request.Username
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("OpenAITask"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &OpenAITaskResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) OpenAITask(request *OpenAITaskRequest) (_result *OpenAITaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &OpenAITaskResponse{}
	_body, _err := client.OpenAITaskWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RefreshDBClusterStorageUsageWithOptions(request *RefreshDBClusterStorageUsageRequest, runtime *util.RuntimeOptions) (_result *RefreshDBClusterStorageUsageResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SyncRealTime)) {
		query["SyncRealTime"] = request.SyncRealTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("RefreshDBClusterStorageUsage"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &RefreshDBClusterStorageUsageResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RefreshDBClusterStorageUsage(request *RefreshDBClusterStorageUsageRequest) (_result *RefreshDBClusterStorageUsageResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &RefreshDBClusterStorageUsageResponse{}
	_body, _err := client.RefreshDBClusterStorageUsageWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RemoveDBClusterFromGDNWithOptions(request *RemoveDBClusterFromGDNRequest, runtime *util.RuntimeOptions) (_result *RemoveDBClusterFromGDNResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.GDNId)) {
		query["GDNId"] = request.GDNId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("RemoveDBClusterFromGDN"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &RemoveDBClusterFromGDNResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RemoveDBClusterFromGDN(request *RemoveDBClusterFromGDNRequest) (_result *RemoveDBClusterFromGDNResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &RemoveDBClusterFromGDNResponse{}
	_body, _err := client.RemoveDBClusterFromGDNWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ResetAccountWithOptions(request *ResetAccountRequest, runtime *util.RuntimeOptions) (_result *ResetAccountResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.AccountPassword)) {
		query["AccountPassword"] = request.AccountPassword
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ResetAccount"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &ResetAccountResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ResetAccount(request *ResetAccountRequest) (_result *ResetAccountResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &ResetAccountResponse{}
	_body, _err := client.ResetAccountWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RestartDBNodeWithOptions(request *RestartDBNodeRequest, runtime *util.RuntimeOptions) (_result *RestartDBNodeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBNodeId)) {
		query["DBNodeId"] = request.DBNodeId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("RestartDBNode"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &RestartDBNodeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RestartDBNode(request *RestartDBNodeRequest) (_result *RestartDBNodeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &RestartDBNodeResponse{}
	_body, _err := client.RestartDBNodeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RestoreTableWithOptions(request *RestoreTableRequest, runtime *util.RuntimeOptions) (_result *RestoreTableResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.BackupId)) {
		query["BackupId"] = request.BackupId
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RestoreTime)) {
		query["RestoreTime"] = request.RestoreTime
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	if !tea.BoolValue(util.IsUnset(request.TableMeta)) {
		query["TableMeta"] = request.TableMeta
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("RestoreTable"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &RestoreTableResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RestoreTable(request *RestoreTableRequest) (_result *RestoreTableResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &RestoreTableResponse{}
	_body, _err := client.RestoreTableWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RevokeAccountPrivilegeWithOptions(request *RevokeAccountPrivilegeRequest, runtime *util.RuntimeOptions) (_result *RevokeAccountPrivilegeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AccountName)) {
		query["AccountName"] = request.AccountName
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBName)) {
		query["DBName"] = request.DBName
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("RevokeAccountPrivilege"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &RevokeAccountPrivilegeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RevokeAccountPrivilege(request *RevokeAccountPrivilegeRequest) (_result *RevokeAccountPrivilegeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &RevokeAccountPrivilegeResponse{}
	_body, _err := client.RevokeAccountPrivilegeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) SwitchOverGlobalDatabaseNetworkWithOptions(request *SwitchOverGlobalDatabaseNetworkRequest, runtime *util.RuntimeOptions) (_result *SwitchOverGlobalDatabaseNetworkResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.GDNId)) {
		query["GDNId"] = request.GDNId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityToken)) {
		query["SecurityToken"] = request.SecurityToken
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("SwitchOverGlobalDatabaseNetwork"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &SwitchOverGlobalDatabaseNetworkResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) SwitchOverGlobalDatabaseNetwork(request *SwitchOverGlobalDatabaseNetworkRequest) (_result *SwitchOverGlobalDatabaseNetworkResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &SwitchOverGlobalDatabaseNetworkResponse{}
	_body, _err := client.SwitchOverGlobalDatabaseNetworkWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) TagResourcesWithOptions(request *TagResourcesRequest, runtime *util.RuntimeOptions) (_result *TagResourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceId)) {
		query["ResourceId"] = request.ResourceId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		query["ResourceType"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.Tag)) {
		query["Tag"] = request.Tag
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("TagResources"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &TagResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) TagResources(request *TagResourcesRequest) (_result *TagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &TagResourcesResponse{}
	_body, _err := client.TagResourcesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) TempModifyDBNodeWithOptions(request *TempModifyDBNodeRequest, runtime *util.RuntimeOptions) (_result *TempModifyDBNodeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.DBNode)) {
		query["DBNode"] = request.DBNode
	}

	if !tea.BoolValue(util.IsUnset(request.ModifyType)) {
		query["ModifyType"] = request.ModifyType
	}

	if !tea.BoolValue(util.IsUnset(request.OperationType)) {
		query["OperationType"] = request.OperationType
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RestoreTime)) {
		query["RestoreTime"] = request.RestoreTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("TempModifyDBNode"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &TempModifyDBNodeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) TempModifyDBNode(request *TempModifyDBNodeRequest) (_result *TempModifyDBNodeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &TempModifyDBNodeResponse{}
	_body, _err := client.TempModifyDBNodeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) TransformDBClusterPayTypeWithOptions(request *TransformDBClusterPayTypeRequest, runtime *util.RuntimeOptions) (_result *TransformDBClusterPayTypeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClientToken)) {
		query["ClientToken"] = request.ClientToken
	}

	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PayType)) {
		query["PayType"] = request.PayType
	}

	if !tea.BoolValue(util.IsUnset(request.Period)) {
		query["Period"] = request.Period
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		query["ResourceGroupId"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.UsedTime)) {
		query["UsedTime"] = request.UsedTime
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("TransformDBClusterPayType"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &TransformDBClusterPayTypeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) TransformDBClusterPayType(request *TransformDBClusterPayTypeRequest) (_result *TransformDBClusterPayTypeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &TransformDBClusterPayTypeResponse{}
	_body, _err := client.TransformDBClusterPayTypeWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UntagResourcesWithOptions(request *UntagResourcesRequest, runtime *util.RuntimeOptions) (_result *UntagResourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.All)) {
		query["All"] = request.All
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["RegionId"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceId)) {
		query["ResourceId"] = request.ResourceId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		query["ResourceType"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.TagKey)) {
		query["TagKey"] = request.TagKey
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("UntagResources"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &UntagResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UntagResources(request *UntagResourcesRequest) (_result *UntagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &UntagResourcesResponse{}
	_body, _err := client.UntagResourcesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpgradeDBClusterMinorVersionWithOptions(request *UpgradeDBClusterMinorVersionRequest, runtime *util.RuntimeOptions) (_result *UpgradeDBClusterMinorVersionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.FromTimeService)) {
		query["FromTimeService"] = request.FromTimeService
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("UpgradeDBClusterMinorVersion"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &UpgradeDBClusterMinorVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpgradeDBClusterMinorVersion(request *UpgradeDBClusterMinorVersionRequest) (_result *UpgradeDBClusterMinorVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &UpgradeDBClusterMinorVersionResponse{}
	_body, _err := client.UpgradeDBClusterMinorVersionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpgradeDBClusterVersionWithOptions(request *UpgradeDBClusterVersionRequest, runtime *util.RuntimeOptions) (_result *UpgradeDBClusterVersionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DBClusterId)) {
		query["DBClusterId"] = request.DBClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.FromTimeService)) {
		query["FromTimeService"] = request.FromTimeService
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerAccount)) {
		query["OwnerAccount"] = request.OwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.OwnerId)) {
		query["OwnerId"] = request.OwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedEndTime)) {
		query["PlannedEndTime"] = request.PlannedEndTime
	}

	if !tea.BoolValue(util.IsUnset(request.PlannedStartTime)) {
		query["PlannedStartTime"] = request.PlannedStartTime
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerAccount)) {
		query["ResourceOwnerAccount"] = request.ResourceOwnerAccount
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceOwnerId)) {
		query["ResourceOwnerId"] = request.ResourceOwnerId
	}

	if !tea.BoolValue(util.IsUnset(request.UpgradeLabel)) {
		query["UpgradeLabel"] = request.UpgradeLabel
	}

	if !tea.BoolValue(util.IsUnset(request.UpgradePolicy)) {
		query["UpgradePolicy"] = request.UpgradePolicy
	}

	if !tea.BoolValue(util.IsUnset(request.UpgradeType)) {
		query["UpgradeType"] = request.UpgradeType
	}

	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("UpgradeDBClusterVersion"),
		Version:     tea.String("2017-08-01"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = &UpgradeDBClusterVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpgradeDBClusterVersion(request *UpgradeDBClusterVersionRequest) (_result *UpgradeDBClusterVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	_result = &UpgradeDBClusterVersionResponse{}
	_body, _err := client.UpgradeDBClusterVersionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}
