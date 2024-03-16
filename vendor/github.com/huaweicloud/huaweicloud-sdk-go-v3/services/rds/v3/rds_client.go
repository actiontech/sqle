package v3

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
)

type RdsClient struct {
	HcClient *http_client.HcHttpClient
}

func NewRdsClient(hcClient *http_client.HcHttpClient) *RdsClient {
	return &RdsClient{HcClient: hcClient}
}

func RdsClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// AddPostgresqlHbaConf 在pg_hba.conf文件最后新增单个或多个配置
//
// 以传入配置全量覆盖当前pg_hba.conf文件内容，入参为空时用默认配置覆盖当前文件内容
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) AddPostgresqlHbaConf(request *model.AddPostgresqlHbaConfRequest) (*model.AddPostgresqlHbaConfResponse, error) {
	requestDef := GenReqDefForAddPostgresqlHbaConf()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AddPostgresqlHbaConfResponse), nil
	}
}

// AddPostgresqlHbaConfInvoker 在pg_hba.conf文件最后新增单个或多个配置
func (c *RdsClient) AddPostgresqlHbaConfInvoker(request *model.AddPostgresqlHbaConfRequest) *AddPostgresqlHbaConfInvoker {
	requestDef := GenReqDefForAddPostgresqlHbaConf()
	return &AddPostgresqlHbaConfInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ApplyConfigurationAsync 应用参数模板
//
// 应用参数模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ApplyConfigurationAsync(request *model.ApplyConfigurationAsyncRequest) (*model.ApplyConfigurationAsyncResponse, error) {
	requestDef := GenReqDefForApplyConfigurationAsync()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ApplyConfigurationAsyncResponse), nil
	}
}

// ApplyConfigurationAsyncInvoker 应用参数模板
func (c *RdsClient) ApplyConfigurationAsyncInvoker(request *model.ApplyConfigurationAsyncRequest) *ApplyConfigurationAsyncInvoker {
	requestDef := GenReqDefForApplyConfigurationAsync()
	return &ApplyConfigurationAsyncInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AttachEip 绑定和解绑弹性公网IP
//
// 绑定和解绑弹性公网IP。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) AttachEip(request *model.AttachEipRequest) (*model.AttachEipResponse, error) {
	requestDef := GenReqDefForAttachEip()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AttachEipResponse), nil
	}
}

// AttachEipInvoker 绑定和解绑弹性公网IP
func (c *RdsClient) AttachEipInvoker(request *model.AttachEipRequest) *AttachEipInvoker {
	requestDef := GenReqDefForAttachEip()
	return &AttachEipInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchDeleteManualBackup 批量删除手动备份
//
// 批量删除手动备份。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) BatchDeleteManualBackup(request *model.BatchDeleteManualBackupRequest) (*model.BatchDeleteManualBackupResponse, error) {
	requestDef := GenReqDefForBatchDeleteManualBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchDeleteManualBackupResponse), nil
	}
}

// BatchDeleteManualBackupInvoker 批量删除手动备份
func (c *RdsClient) BatchDeleteManualBackupInvoker(request *model.BatchDeleteManualBackupRequest) *BatchDeleteManualBackupInvoker {
	requestDef := GenReqDefForBatchDeleteManualBackup()
	return &BatchDeleteManualBackupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchRestoreDatabase 库级时间点恢复
//
// 库级时间点恢复
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) BatchRestoreDatabase(request *model.BatchRestoreDatabaseRequest) (*model.BatchRestoreDatabaseResponse, error) {
	requestDef := GenReqDefForBatchRestoreDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchRestoreDatabaseResponse), nil
	}
}

// BatchRestoreDatabaseInvoker 库级时间点恢复
func (c *RdsClient) BatchRestoreDatabaseInvoker(request *model.BatchRestoreDatabaseRequest) *BatchRestoreDatabaseInvoker {
	requestDef := GenReqDefForBatchRestoreDatabase()
	return &BatchRestoreDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchRestorePostgreSqlTables 表级时间点恢复（PostgreSQL）
//
// 表级时间点恢复（PostgreSQL）
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) BatchRestorePostgreSqlTables(request *model.BatchRestorePostgreSqlTablesRequest) (*model.BatchRestorePostgreSqlTablesResponse, error) {
	requestDef := GenReqDefForBatchRestorePostgreSqlTables()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchRestorePostgreSqlTablesResponse), nil
	}
}

// BatchRestorePostgreSqlTablesInvoker 表级时间点恢复（PostgreSQL）
func (c *RdsClient) BatchRestorePostgreSqlTablesInvoker(request *model.BatchRestorePostgreSqlTablesRequest) *BatchRestorePostgreSqlTablesInvoker {
	requestDef := GenReqDefForBatchRestorePostgreSqlTables()
	return &BatchRestorePostgreSqlTablesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchTagAddAction 批量添加标签
//
// 批量添加标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) BatchTagAddAction(request *model.BatchTagAddActionRequest) (*model.BatchTagAddActionResponse, error) {
	requestDef := GenReqDefForBatchTagAddAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchTagAddActionResponse), nil
	}
}

// BatchTagAddActionInvoker 批量添加标签
func (c *RdsClient) BatchTagAddActionInvoker(request *model.BatchTagAddActionRequest) *BatchTagAddActionInvoker {
	requestDef := GenReqDefForBatchTagAddAction()
	return &BatchTagAddActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchTagDelAction 批量删除标签
//
// 批量删除标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) BatchTagDelAction(request *model.BatchTagDelActionRequest) (*model.BatchTagDelActionResponse, error) {
	requestDef := GenReqDefForBatchTagDelAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchTagDelActionResponse), nil
	}
}

// BatchTagDelActionInvoker 批量删除标签
func (c *RdsClient) BatchTagDelActionInvoker(request *model.BatchTagDelActionRequest) *BatchTagDelActionInvoker {
	requestDef := GenReqDefForBatchTagDelAction()
	return &BatchTagDelActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeFailoverMode 更改主备实例的数据同步方式
//
// 更改主备实例的数据同步方式。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ChangeFailoverMode(request *model.ChangeFailoverModeRequest) (*model.ChangeFailoverModeResponse, error) {
	requestDef := GenReqDefForChangeFailoverMode()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeFailoverModeResponse), nil
	}
}

// ChangeFailoverModeInvoker 更改主备实例的数据同步方式
func (c *RdsClient) ChangeFailoverModeInvoker(request *model.ChangeFailoverModeRequest) *ChangeFailoverModeInvoker {
	requestDef := GenReqDefForChangeFailoverMode()
	return &ChangeFailoverModeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeFailoverStrategy 切换主备实例的倒换策略
//
// 切换主备实例的倒换策略.
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ChangeFailoverStrategy(request *model.ChangeFailoverStrategyRequest) (*model.ChangeFailoverStrategyResponse, error) {
	requestDef := GenReqDefForChangeFailoverStrategy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeFailoverStrategyResponse), nil
	}
}

// ChangeFailoverStrategyInvoker 切换主备实例的倒换策略
func (c *RdsClient) ChangeFailoverStrategyInvoker(request *model.ChangeFailoverStrategyRequest) *ChangeFailoverStrategyInvoker {
	requestDef := GenReqDefForChangeFailoverStrategy()
	return &ChangeFailoverStrategyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeOpsWindow 设置可维护时间段
//
// 设置可维护时间段
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ChangeOpsWindow(request *model.ChangeOpsWindowRequest) (*model.ChangeOpsWindowResponse, error) {
	requestDef := GenReqDefForChangeOpsWindow()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeOpsWindowResponse), nil
	}
}

// ChangeOpsWindowInvoker 设置可维护时间段
func (c *RdsClient) ChangeOpsWindowInvoker(request *model.ChangeOpsWindowRequest) *ChangeOpsWindowInvoker {
	requestDef := GenReqDefForChangeOpsWindow()
	return &ChangeOpsWindowInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CopyConfiguration 复制参数模板
//
// 复制参数模板
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CopyConfiguration(request *model.CopyConfigurationRequest) (*model.CopyConfigurationResponse, error) {
	requestDef := GenReqDefForCopyConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CopyConfigurationResponse), nil
	}
}

// CopyConfigurationInvoker 复制参数模板
func (c *RdsClient) CopyConfigurationInvoker(request *model.CopyConfigurationRequest) *CopyConfigurationInvoker {
	requestDef := GenReqDefForCopyConfiguration()
	return &CopyConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateConfiguration 创建参数模板
//
// 创建参数模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateConfiguration(request *model.CreateConfigurationRequest) (*model.CreateConfigurationResponse, error) {
	requestDef := GenReqDefForCreateConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateConfigurationResponse), nil
	}
}

// CreateConfigurationInvoker 创建参数模板
func (c *RdsClient) CreateConfigurationInvoker(request *model.CreateConfigurationRequest) *CreateConfigurationInvoker {
	requestDef := GenReqDefForCreateConfiguration()
	return &CreateConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDnsName 申请域名
//
// 申请域名
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateDnsName(request *model.CreateDnsNameRequest) (*model.CreateDnsNameResponse, error) {
	requestDef := GenReqDefForCreateDnsName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDnsNameResponse), nil
	}
}

// CreateDnsNameInvoker 申请域名
func (c *RdsClient) CreateDnsNameInvoker(request *model.CreateDnsNameRequest) *CreateDnsNameInvoker {
	requestDef := GenReqDefForCreateDnsName()
	return &CreateDnsNameInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateInstance 创建数据库实例
//
// 创建数据库实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateInstance(request *model.CreateInstanceRequest) (*model.CreateInstanceResponse, error) {
	requestDef := GenReqDefForCreateInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateInstanceResponse), nil
	}
}

// CreateInstanceInvoker 创建数据库实例
func (c *RdsClient) CreateInstanceInvoker(request *model.CreateInstanceRequest) *CreateInstanceInvoker {
	requestDef := GenReqDefForCreateInstance()
	return &CreateInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateManualBackup 创建手动备份
//
// 创建手动备份。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateManualBackup(request *model.CreateManualBackupRequest) (*model.CreateManualBackupResponse, error) {
	requestDef := GenReqDefForCreateManualBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateManualBackupResponse), nil
	}
}

// CreateManualBackupInvoker 创建手动备份
func (c *RdsClient) CreateManualBackupInvoker(request *model.CreateManualBackupRequest) *CreateManualBackupInvoker {
	requestDef := GenReqDefForCreateManualBackup()
	return &CreateManualBackupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateRestoreInstance 恢复到新实例
//
// 恢复到新实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateRestoreInstance(request *model.CreateRestoreInstanceRequest) (*model.CreateRestoreInstanceResponse, error) {
	requestDef := GenReqDefForCreateRestoreInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateRestoreInstanceResponse), nil
	}
}

// CreateRestoreInstanceInvoker 恢复到新实例
func (c *RdsClient) CreateRestoreInstanceInvoker(request *model.CreateRestoreInstanceRequest) *CreateRestoreInstanceInvoker {
	requestDef := GenReqDefForCreateRestoreInstance()
	return &CreateRestoreInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateXelLogDownload 获取扩展日志下载信息
//
// 获取扩展日志下载信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateXelLogDownload(request *model.CreateXelLogDownloadRequest) (*model.CreateXelLogDownloadResponse, error) {
	requestDef := GenReqDefForCreateXelLogDownload()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateXelLogDownloadResponse), nil
	}
}

// CreateXelLogDownloadInvoker 获取扩展日志下载信息
func (c *RdsClient) CreateXelLogDownloadInvoker(request *model.CreateXelLogDownloadRequest) *CreateXelLogDownloadInvoker {
	requestDef := GenReqDefForCreateXelLogDownload()
	return &CreateXelLogDownloadInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteConfiguration 删除参数模板
//
// 删除参数模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteConfiguration(request *model.DeleteConfigurationRequest) (*model.DeleteConfigurationResponse, error) {
	requestDef := GenReqDefForDeleteConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteConfigurationResponse), nil
	}
}

// DeleteConfigurationInvoker 删除参数模板
func (c *RdsClient) DeleteConfigurationInvoker(request *model.DeleteConfigurationRequest) *DeleteConfigurationInvoker {
	requestDef := GenReqDefForDeleteConfiguration()
	return &DeleteConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteInstance 删除数据库实例
//
// 删除数据库实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteInstance(request *model.DeleteInstanceRequest) (*model.DeleteInstanceResponse, error) {
	requestDef := GenReqDefForDeleteInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteInstanceResponse), nil
	}
}

// DeleteInstanceInvoker 删除数据库实例
func (c *RdsClient) DeleteInstanceInvoker(request *model.DeleteInstanceRequest) *DeleteInstanceInvoker {
	requestDef := GenReqDefForDeleteInstance()
	return &DeleteInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteJob 删除即时任务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteJob(request *model.DeleteJobRequest) (*model.DeleteJobResponse, error) {
	requestDef := GenReqDefForDeleteJob()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteJobResponse), nil
	}
}

// DeleteJobInvoker 删除即时任务
func (c *RdsClient) DeleteJobInvoker(request *model.DeleteJobRequest) *DeleteJobInvoker {
	requestDef := GenReqDefForDeleteJob()
	return &DeleteJobInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteManualBackup 删除手动备份
//
// 删除手动备份。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteManualBackup(request *model.DeleteManualBackupRequest) (*model.DeleteManualBackupResponse, error) {
	requestDef := GenReqDefForDeleteManualBackup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteManualBackupResponse), nil
	}
}

// DeleteManualBackupInvoker 删除手动备份
func (c *RdsClient) DeleteManualBackupInvoker(request *model.DeleteManualBackupRequest) *DeleteManualBackupInvoker {
	requestDef := GenReqDefForDeleteManualBackup()
	return &DeleteManualBackupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeletePostgresqlHbaConf 删除pg_hba.conf文件的单个或多个配置
//
// 删除pg_hba.conf文件的单个或多个配置，以priority做唯一标识
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeletePostgresqlHbaConf(request *model.DeletePostgresqlHbaConfRequest) (*model.DeletePostgresqlHbaConfResponse, error) {
	requestDef := GenReqDefForDeletePostgresqlHbaConf()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeletePostgresqlHbaConfResponse), nil
	}
}

// DeletePostgresqlHbaConfInvoker 删除pg_hba.conf文件的单个或多个配置
func (c *RdsClient) DeletePostgresqlHbaConfInvoker(request *model.DeletePostgresqlHbaConfRequest) *DeletePostgresqlHbaConfInvoker {
	requestDef := GenReqDefForDeletePostgresqlHbaConf()
	return &DeletePostgresqlHbaConfInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DownloadSlowlog 获取慢日志下载链接
//
// 获取慢日志下载链接。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DownloadSlowlog(request *model.DownloadSlowlogRequest) (*model.DownloadSlowlogResponse, error) {
	requestDef := GenReqDefForDownloadSlowlog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DownloadSlowlogResponse), nil
	}
}

// DownloadSlowlogInvoker 获取慢日志下载链接
func (c *RdsClient) DownloadSlowlogInvoker(request *model.DownloadSlowlogRequest) *DownloadSlowlogInvoker {
	requestDef := GenReqDefForDownloadSlowlog()
	return &DownloadSlowlogInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// EnableConfiguration 应用参数模板
//
// 应用参数模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) EnableConfiguration(request *model.EnableConfigurationRequest) (*model.EnableConfigurationResponse, error) {
	requestDef := GenReqDefForEnableConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.EnableConfigurationResponse), nil
	}
}

// EnableConfigurationInvoker 应用参数模板
func (c *RdsClient) EnableConfigurationInvoker(request *model.EnableConfigurationRequest) *EnableConfigurationInvoker {
	requestDef := GenReqDefForEnableConfiguration()
	return &EnableConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAuditlogs 获取审计日志列表
//
// 获取审计日志列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListAuditlogs(request *model.ListAuditlogsRequest) (*model.ListAuditlogsResponse, error) {
	requestDef := GenReqDefForListAuditlogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuditlogsResponse), nil
	}
}

// ListAuditlogsInvoker 获取审计日志列表
func (c *RdsClient) ListAuditlogsInvoker(request *model.ListAuditlogsRequest) *ListAuditlogsInvoker {
	requestDef := GenReqDefForListAuditlogs()
	return &ListAuditlogsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListBackups 获取备份列表
//
// 获取备份列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListBackups(request *model.ListBackupsRequest) (*model.ListBackupsResponse, error) {
	requestDef := GenReqDefForListBackups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListBackupsResponse), nil
	}
}

// ListBackupsInvoker 获取备份列表
func (c *RdsClient) ListBackupsInvoker(request *model.ListBackupsRequest) *ListBackupsInvoker {
	requestDef := GenReqDefForListBackups()
	return &ListBackupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListCollations 查询SQLServer可用字符集
//
// 查询SQLServer可用字符集
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListCollations(request *model.ListCollationsRequest) (*model.ListCollationsResponse, error) {
	requestDef := GenReqDefForListCollations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListCollationsResponse), nil
	}
}

// ListCollationsInvoker 查询SQLServer可用字符集
func (c *RdsClient) ListCollationsInvoker(request *model.ListCollationsRequest) *ListCollationsInvoker {
	requestDef := GenReqDefForListCollations()
	return &ListCollationsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListConfigurations 获取参数模板列表
//
// 获取参数模板列表，包括所有数据库的默认参数模板和用户创建的参数模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListConfigurations(request *model.ListConfigurationsRequest) (*model.ListConfigurationsResponse, error) {
	requestDef := GenReqDefForListConfigurations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListConfigurationsResponse), nil
	}
}

// ListConfigurationsInvoker 获取参数模板列表
func (c *RdsClient) ListConfigurationsInvoker(request *model.ListConfigurationsRequest) *ListConfigurationsInvoker {
	requestDef := GenReqDefForListConfigurations()
	return &ListConfigurationsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDatastores 查询数据库引擎的版本
//
// 查询数据库引擎的版本。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListDatastores(request *model.ListDatastoresRequest) (*model.ListDatastoresResponse, error) {
	requestDef := GenReqDefForListDatastores()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDatastoresResponse), nil
	}
}

// ListDatastoresInvoker 查询数据库引擎的版本
func (c *RdsClient) ListDatastoresInvoker(request *model.ListDatastoresRequest) *ListDatastoresInvoker {
	requestDef := GenReqDefForListDatastores()
	return &ListDatastoresInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDrRelations
//
// 批量查询容灾实例信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListDrRelations(request *model.ListDrRelationsRequest) (*model.ListDrRelationsResponse, error) {
	requestDef := GenReqDefForListDrRelations()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDrRelationsResponse), nil
	}
}

// ListDrRelationsInvoker
func (c *RdsClient) ListDrRelationsInvoker(request *model.ListDrRelationsRequest) *ListDrRelationsInvoker {
	requestDef := GenReqDefForListDrRelations()
	return &ListDrRelationsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListEngineFlavors 查询实例可变更规格
//
// 查询实例可变更规格
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListEngineFlavors(request *model.ListEngineFlavorsRequest) (*model.ListEngineFlavorsResponse, error) {
	requestDef := GenReqDefForListEngineFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListEngineFlavorsResponse), nil
	}
}

// ListEngineFlavorsInvoker 查询实例可变更规格
func (c *RdsClient) ListEngineFlavorsInvoker(request *model.ListEngineFlavorsRequest) *ListEngineFlavorsInvoker {
	requestDef := GenReqDefForListEngineFlavors()
	return &ListEngineFlavorsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListErrorLogs 查询数据库错误日志
//
// 查询数据库错误日志。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListErrorLogs(request *model.ListErrorLogsRequest) (*model.ListErrorLogsResponse, error) {
	requestDef := GenReqDefForListErrorLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListErrorLogsResponse), nil
	}
}

// ListErrorLogsInvoker 查询数据库错误日志
func (c *RdsClient) ListErrorLogsInvoker(request *model.ListErrorLogsRequest) *ListErrorLogsInvoker {
	requestDef := GenReqDefForListErrorLogs()
	return &ListErrorLogsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListErrorLogsNew 查询数据库错误日志
//
// 查询数据库错误日志。(与原v3接口相比修改offset,符合华为云服务开放 API遵从性规范3.0)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListErrorLogsNew(request *model.ListErrorLogsNewRequest) (*model.ListErrorLogsNewResponse, error) {
	requestDef := GenReqDefForListErrorLogsNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListErrorLogsNewResponse), nil
	}
}

// ListErrorLogsNewInvoker 查询数据库错误日志
func (c *RdsClient) ListErrorLogsNewInvoker(request *model.ListErrorLogsNewRequest) *ListErrorLogsNewInvoker {
	requestDef := GenReqDefForListErrorLogsNew()
	return &ListErrorLogsNewInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListErrorlogForLts
//
// 查询实例的错误日志数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListErrorlogForLts(request *model.ListErrorlogForLtsRequest) (*model.ListErrorlogForLtsResponse, error) {
	requestDef := GenReqDefForListErrorlogForLts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListErrorlogForLtsResponse), nil
	}
}

// ListErrorlogForLtsInvoker
func (c *RdsClient) ListErrorlogForLtsInvoker(request *model.ListErrorlogForLtsRequest) *ListErrorlogForLtsInvoker {
	requestDef := GenReqDefForListErrorlogForLts()
	return &ListErrorlogForLtsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListFlavors 查询数据库规格
//
// 查询数据库规格。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListFlavors(request *model.ListFlavorsRequest) (*model.ListFlavorsResponse, error) {
	requestDef := GenReqDefForListFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListFlavorsResponse), nil
	}
}

// ListFlavorsInvoker 查询数据库规格
func (c *RdsClient) ListFlavorsInvoker(request *model.ListFlavorsRequest) *ListFlavorsInvoker {
	requestDef := GenReqDefForListFlavors()
	return &ListFlavorsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListHistoryDatabase 查询指定时间点可恢复的库
//
// 查询指定时间点可恢复的库
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListHistoryDatabase(request *model.ListHistoryDatabaseRequest) (*model.ListHistoryDatabaseResponse, error) {
	requestDef := GenReqDefForListHistoryDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListHistoryDatabaseResponse), nil
	}
}

// ListHistoryDatabaseInvoker 查询指定时间点可恢复的库
func (c *RdsClient) ListHistoryDatabaseInvoker(request *model.ListHistoryDatabaseRequest) *ListHistoryDatabaseInvoker {
	requestDef := GenReqDefForListHistoryDatabase()
	return &ListHistoryDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInspectionHistories
//
// 查询实例大版本升级检查历史。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListInspectionHistories(request *model.ListInspectionHistoriesRequest) (*model.ListInspectionHistoriesResponse, error) {
	requestDef := GenReqDefForListInspectionHistories()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInspectionHistoriesResponse), nil
	}
}

// ListInspectionHistoriesInvoker
func (c *RdsClient) ListInspectionHistoriesInvoker(request *model.ListInspectionHistoriesRequest) *ListInspectionHistoriesInvoker {
	requestDef := GenReqDefForListInspectionHistories()
	return &ListInspectionHistoriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstanceDiagnosis 获取诊断后的实例数量
//
// 获取诊断后的实例数量
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListInstanceDiagnosis(request *model.ListInstanceDiagnosisRequest) (*model.ListInstanceDiagnosisResponse, error) {
	requestDef := GenReqDefForListInstanceDiagnosis()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstanceDiagnosisResponse), nil
	}
}

// ListInstanceDiagnosisInvoker 获取诊断后的实例数量
func (c *RdsClient) ListInstanceDiagnosisInvoker(request *model.ListInstanceDiagnosisRequest) *ListInstanceDiagnosisInvoker {
	requestDef := GenReqDefForListInstanceDiagnosis()
	return &ListInstanceDiagnosisInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstanceParamHistories 查询实例参数修改历史
//
// 实例参数修改历史。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListInstanceParamHistories(request *model.ListInstanceParamHistoriesRequest) (*model.ListInstanceParamHistoriesResponse, error) {
	requestDef := GenReqDefForListInstanceParamHistories()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstanceParamHistoriesResponse), nil
	}
}

// ListInstanceParamHistoriesInvoker 查询实例参数修改历史
func (c *RdsClient) ListInstanceParamHistoriesInvoker(request *model.ListInstanceParamHistoriesRequest) *ListInstanceParamHistoriesInvoker {
	requestDef := GenReqDefForListInstanceParamHistories()
	return &ListInstanceParamHistoriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstanceTags 查询实例标签
//
// 查询实例标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListInstanceTags(request *model.ListInstanceTagsRequest) (*model.ListInstanceTagsResponse, error) {
	requestDef := GenReqDefForListInstanceTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstanceTagsResponse), nil
	}
}

// ListInstanceTagsInvoker 查询实例标签
func (c *RdsClient) ListInstanceTagsInvoker(request *model.ListInstanceTagsRequest) *ListInstanceTagsInvoker {
	requestDef := GenReqDefForListInstanceTags()
	return &ListInstanceTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstances 查询数据库实例列表
//
// 查询数据库实例列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListInstances(request *model.ListInstancesRequest) (*model.ListInstancesResponse, error) {
	requestDef := GenReqDefForListInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstancesResponse), nil
	}
}

// ListInstancesInvoker 查询数据库实例列表
func (c *RdsClient) ListInstancesInvoker(request *model.ListInstancesRequest) *ListInstancesInvoker {
	requestDef := GenReqDefForListInstances()
	return &ListInstancesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstancesInfoDiagnosis 获取指定诊断项的诊断结果
//
// 获取指定诊断项的诊断结果
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListInstancesInfoDiagnosis(request *model.ListInstancesInfoDiagnosisRequest) (*model.ListInstancesInfoDiagnosisResponse, error) {
	requestDef := GenReqDefForListInstancesInfoDiagnosis()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstancesInfoDiagnosisResponse), nil
	}
}

// ListInstancesInfoDiagnosisInvoker 获取指定诊断项的诊断结果
func (c *RdsClient) ListInstancesInfoDiagnosisInvoker(request *model.ListInstancesInfoDiagnosisRequest) *ListInstancesInfoDiagnosisInvoker {
	requestDef := GenReqDefForListInstancesInfoDiagnosis()
	return &ListInstancesInfoDiagnosisInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListInstancesSupportFastRestore 获取实例是否能使用极速恢复
//
// 批量获取实例是否能在库表恢复时使用极速恢复。
//
// - 调用接口前，您需要了解API 认证鉴权。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListInstancesSupportFastRestore(request *model.ListInstancesSupportFastRestoreRequest) (*model.ListInstancesSupportFastRestoreResponse, error) {
	requestDef := GenReqDefForListInstancesSupportFastRestore()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListInstancesSupportFastRestoreResponse), nil
	}
}

// ListInstancesSupportFastRestoreInvoker 获取实例是否能使用极速恢复
func (c *RdsClient) ListInstancesSupportFastRestoreInvoker(request *model.ListInstancesSupportFastRestoreRequest) *ListInstancesSupportFastRestoreInvoker {
	requestDef := GenReqDefForListInstancesSupportFastRestore()
	return &ListInstancesSupportFastRestoreInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListJobInfo 获取指定ID的任务信息
//
// 获取指定ID的任务信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListJobInfo(request *model.ListJobInfoRequest) (*model.ListJobInfoResponse, error) {
	requestDef := GenReqDefForListJobInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListJobInfoResponse), nil
	}
}

// ListJobInfoInvoker 获取指定ID的任务信息
func (c *RdsClient) ListJobInfoInvoker(request *model.ListJobInfoRequest) *ListJobInfoInvoker {
	requestDef := GenReqDefForListJobInfo()
	return &ListJobInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListJobInfoDetail 获取指定实例和时间范围的任务信息（SQL Server）
//
// 获取指定实例和时间范围的任务信息（SQL Server）。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListJobInfoDetail(request *model.ListJobInfoDetailRequest) (*model.ListJobInfoDetailResponse, error) {
	requestDef := GenReqDefForListJobInfoDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListJobInfoDetailResponse), nil
	}
}

// ListJobInfoDetailInvoker 获取指定实例和时间范围的任务信息（SQL Server）
func (c *RdsClient) ListJobInfoDetailInvoker(request *model.ListJobInfoDetailRequest) *ListJobInfoDetailInvoker {
	requestDef := GenReqDefForListJobInfoDetail()
	return &ListJobInfoDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListOffSiteBackups 查询跨区域备份列表
//
// 查询跨区域备份列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListOffSiteBackups(request *model.ListOffSiteBackupsRequest) (*model.ListOffSiteBackupsResponse, error) {
	requestDef := GenReqDefForListOffSiteBackups()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteBackupsResponse), nil
	}
}

// ListOffSiteBackupsInvoker 查询跨区域备份列表
func (c *RdsClient) ListOffSiteBackupsInvoker(request *model.ListOffSiteBackupsRequest) *ListOffSiteBackupsInvoker {
	requestDef := GenReqDefForListOffSiteBackups()
	return &ListOffSiteBackupsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListOffSiteInstances 查询跨区域备份实例列表
//
// 查询跨区域备份实例列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListOffSiteInstances(request *model.ListOffSiteInstancesRequest) (*model.ListOffSiteInstancesResponse, error) {
	requestDef := GenReqDefForListOffSiteInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteInstancesResponse), nil
	}
}

// ListOffSiteInstancesInvoker 查询跨区域备份实例列表
func (c *RdsClient) ListOffSiteInstancesInvoker(request *model.ListOffSiteInstancesRequest) *ListOffSiteInstancesInvoker {
	requestDef := GenReqDefForListOffSiteInstances()
	return &ListOffSiteInstancesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListOffSiteRestoreTimes 查询跨区域备份可恢复时间段
//
// 查询跨区域备份可恢复时间段。
// 如果您备份策略中的保存天数设置较长，建议您传入查询日期“date”。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListOffSiteRestoreTimes(request *model.ListOffSiteRestoreTimesRequest) (*model.ListOffSiteRestoreTimesResponse, error) {
	requestDef := GenReqDefForListOffSiteRestoreTimes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListOffSiteRestoreTimesResponse), nil
	}
}

// ListOffSiteRestoreTimesInvoker 查询跨区域备份可恢复时间段
func (c *RdsClient) ListOffSiteRestoreTimesInvoker(request *model.ListOffSiteRestoreTimesRequest) *ListOffSiteRestoreTimesInvoker {
	requestDef := GenReqDefForListOffSiteRestoreTimes()
	return &ListOffSiteRestoreTimesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPostgresqlHbaInfo 查询实例的pg_hba.conf文件配置
//
// 查询实例的pg_hba.conf文件配置
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlHbaInfo(request *model.ListPostgresqlHbaInfoRequest) (*model.ListPostgresqlHbaInfoResponse, error) {
	requestDef := GenReqDefForListPostgresqlHbaInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlHbaInfoResponse), nil
	}
}

// ListPostgresqlHbaInfoInvoker 查询实例的pg_hba.conf文件配置
func (c *RdsClient) ListPostgresqlHbaInfoInvoker(request *model.ListPostgresqlHbaInfoRequest) *ListPostgresqlHbaInfoInvoker {
	requestDef := GenReqDefForListPostgresqlHbaInfo()
	return &ListPostgresqlHbaInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPostgresqlHbaInfoHistory 查询实例的pg_hba.conf文件修改历史
//
// 查询实例的pg_hba.conf文件修改历史
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlHbaInfoHistory(request *model.ListPostgresqlHbaInfoHistoryRequest) (*model.ListPostgresqlHbaInfoHistoryResponse, error) {
	requestDef := GenReqDefForListPostgresqlHbaInfoHistory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlHbaInfoHistoryResponse), nil
	}
}

// ListPostgresqlHbaInfoHistoryInvoker 查询实例的pg_hba.conf文件修改历史
func (c *RdsClient) ListPostgresqlHbaInfoHistoryInvoker(request *model.ListPostgresqlHbaInfoHistoryRequest) *ListPostgresqlHbaInfoHistoryInvoker {
	requestDef := GenReqDefForListPostgresqlHbaInfoHistory()
	return &ListPostgresqlHbaInfoHistoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPostgresqlListHistoryTables 查询指定时间点可恢复的表(PostgreSQL)
//
// 查询指定时间点可恢复的表(PostgreSQL)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlListHistoryTables(request *model.ListPostgresqlListHistoryTablesRequest) (*model.ListPostgresqlListHistoryTablesResponse, error) {
	requestDef := GenReqDefForListPostgresqlListHistoryTables()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlListHistoryTablesResponse), nil
	}
}

// ListPostgresqlListHistoryTablesInvoker 查询指定时间点可恢复的表(PostgreSQL)
func (c *RdsClient) ListPostgresqlListHistoryTablesInvoker(request *model.ListPostgresqlListHistoryTablesRequest) *ListPostgresqlListHistoryTablesInvoker {
	requestDef := GenReqDefForListPostgresqlListHistoryTables()
	return &ListPostgresqlListHistoryTablesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPredefinedTag
//
// 查询预定义标签
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListPredefinedTag(request *model.ListPredefinedTagRequest) (*model.ListPredefinedTagResponse, error) {
	requestDef := GenReqDefForListPredefinedTag()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPredefinedTagResponse), nil
	}
}

// ListPredefinedTagInvoker
func (c *RdsClient) ListPredefinedTagInvoker(request *model.ListPredefinedTagRequest) *ListPredefinedTagInvoker {
	requestDef := GenReqDefForListPredefinedTag()
	return &ListPredefinedTagInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProjectTags 查询项目标签
//
// 查询项目标签。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListProjectTags(request *model.ListProjectTagsRequest) (*model.ListProjectTagsResponse, error) {
	requestDef := GenReqDefForListProjectTags()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProjectTagsResponse), nil
	}
}

// ListProjectTagsInvoker 查询项目标签
func (c *RdsClient) ListProjectTagsInvoker(request *model.ListProjectTagsRequest) *ListProjectTagsInvoker {
	requestDef := GenReqDefForListProjectTags()
	return &ListProjectTagsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRecycleInstances 查询回收站
//
// 查询回收站实例信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListRecycleInstances(request *model.ListRecycleInstancesRequest) (*model.ListRecycleInstancesResponse, error) {
	requestDef := GenReqDefForListRecycleInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRecycleInstancesResponse), nil
	}
}

// ListRecycleInstancesInvoker 查询回收站
func (c *RdsClient) ListRecycleInstancesInvoker(request *model.ListRecycleInstancesRequest) *ListRecycleInstancesInvoker {
	requestDef := GenReqDefForListRecycleInstances()
	return &ListRecycleInstancesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListRestoreTimes 查询可恢复时间段
//
// 查询可恢复时间段。
// 如果您备份策略中的保存天数设置较长，建议您传入查询日期“date”。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListRestoreTimes(request *model.ListRestoreTimesRequest) (*model.ListRestoreTimesResponse, error) {
	requestDef := GenReqDefForListRestoreTimes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListRestoreTimesResponse), nil
	}
}

// ListRestoreTimesInvoker 查询可恢复时间段
func (c *RdsClient) ListRestoreTimesInvoker(request *model.ListRestoreTimesRequest) *ListRestoreTimesInvoker {
	requestDef := GenReqDefForListRestoreTimes()
	return &ListRestoreTimesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSimplifiedInstances
//
// 获取指定实例详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSimplifiedInstances(request *model.ListSimplifiedInstancesRequest) (*model.ListSimplifiedInstancesResponse, error) {
	requestDef := GenReqDefForListSimplifiedInstances()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSimplifiedInstancesResponse), nil
	}
}

// ListSimplifiedInstancesInvoker
func (c *RdsClient) ListSimplifiedInstancesInvoker(request *model.ListSimplifiedInstancesRequest) *ListSimplifiedInstancesInvoker {
	requestDef := GenReqDefForListSimplifiedInstances()
	return &ListSimplifiedInstancesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSlowLogFile 查询慢日志文件列表
//
// 查询慢日志文件列表。
// 调用该接口取到慢日志文件名后，可以调用接口/v3/{project_id}/instances/{instance_id}/slowlog-download 获取慢日志文件下载链接
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSlowLogFile(request *model.ListSlowLogFileRequest) (*model.ListSlowLogFileResponse, error) {
	requestDef := GenReqDefForListSlowLogFile()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogFileResponse), nil
	}
}

// ListSlowLogFileInvoker 查询慢日志文件列表
func (c *RdsClient) ListSlowLogFileInvoker(request *model.ListSlowLogFileRequest) *ListSlowLogFileInvoker {
	requestDef := GenReqDefForListSlowLogFile()
	return &ListSlowLogFileInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSlowLogStatisticsForLts
//
// 查询实例慢日志的统计数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSlowLogStatisticsForLts(request *model.ListSlowLogStatisticsForLtsRequest) (*model.ListSlowLogStatisticsForLtsResponse, error) {
	requestDef := GenReqDefForListSlowLogStatisticsForLts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogStatisticsForLtsResponse), nil
	}
}

// ListSlowLogStatisticsForLtsInvoker
func (c *RdsClient) ListSlowLogStatisticsForLtsInvoker(request *model.ListSlowLogStatisticsForLtsRequest) *ListSlowLogStatisticsForLtsInvoker {
	requestDef := GenReqDefForListSlowLogStatisticsForLts()
	return &ListSlowLogStatisticsForLtsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSlowLogs 查询数据库慢日志
//
// 查询数据库慢日志。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSlowLogs(request *model.ListSlowLogsRequest) (*model.ListSlowLogsResponse, error) {
	requestDef := GenReqDefForListSlowLogs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogsResponse), nil
	}
}

// ListSlowLogsInvoker 查询数据库慢日志
func (c *RdsClient) ListSlowLogsInvoker(request *model.ListSlowLogsRequest) *ListSlowLogsInvoker {
	requestDef := GenReqDefForListSlowLogs()
	return &ListSlowLogsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSlowLogsNew 查询数据库慢日志
//
// 查询数据库慢日志。(与原v3接口相比修改offset,符合华为云服务开放 API遵从性规范3.0)
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSlowLogsNew(request *model.ListSlowLogsNewRequest) (*model.ListSlowLogsNewResponse, error) {
	requestDef := GenReqDefForListSlowLogsNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowLogsNewResponse), nil
	}
}

// ListSlowLogsNewInvoker 查询数据库慢日志
func (c *RdsClient) ListSlowLogsNewInvoker(request *model.ListSlowLogsNewRequest) *ListSlowLogsNewInvoker {
	requestDef := GenReqDefForListSlowLogsNew()
	return &ListSlowLogsNewInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSlowlogForLts
//
// 查询实例的慢日志数据。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSlowlogForLts(request *model.ListSlowlogForLtsRequest) (*model.ListSlowlogForLtsResponse, error) {
	requestDef := GenReqDefForListSlowlogForLts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowlogForLtsResponse), nil
	}
}

// ListSlowlogForLtsInvoker
func (c *RdsClient) ListSlowlogForLtsInvoker(request *model.ListSlowlogForLtsRequest) *ListSlowlogForLtsInvoker {
	requestDef := GenReqDefForListSlowlogForLts()
	return &ListSlowlogForLtsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSlowlogStatistics 获取慢日志统计信息
//
// 获取慢日志统计信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSlowlogStatistics(request *model.ListSlowlogStatisticsRequest) (*model.ListSlowlogStatisticsResponse, error) {
	requestDef := GenReqDefForListSlowlogStatistics()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSlowlogStatisticsResponse), nil
	}
}

// ListSlowlogStatisticsInvoker 获取慢日志统计信息
func (c *RdsClient) ListSlowlogStatisticsInvoker(request *model.ListSlowlogStatisticsRequest) *ListSlowlogStatisticsInvoker {
	requestDef := GenReqDefForListSlowlogStatistics()
	return &ListSlowlogStatisticsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSslCertDownloadLink 获取SSL证书下载地址
//
// 获取SSL证书下载地址
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSslCertDownloadLink(request *model.ListSslCertDownloadLinkRequest) (*model.ListSslCertDownloadLinkResponse, error) {
	requestDef := GenReqDefForListSslCertDownloadLink()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSslCertDownloadLinkResponse), nil
	}
}

// ListSslCertDownloadLinkInvoker 获取SSL证书下载地址
func (c *RdsClient) ListSslCertDownloadLinkInvoker(request *model.ListSslCertDownloadLinkRequest) *ListSslCertDownloadLinkInvoker {
	requestDef := GenReqDefForListSslCertDownloadLink()
	return &ListSslCertDownloadLinkInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListStorageTypes 查询数据库磁盘类型
//
// 查询数据库磁盘类型。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListStorageTypes(request *model.ListStorageTypesRequest) (*model.ListStorageTypesResponse, error) {
	requestDef := GenReqDefForListStorageTypes()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListStorageTypesResponse), nil
	}
}

// ListStorageTypesInvoker 查询数据库磁盘类型
func (c *RdsClient) ListStorageTypesInvoker(request *model.ListStorageTypesRequest) *ListStorageTypesInvoker {
	requestDef := GenReqDefForListStorageTypes()
	return &ListStorageTypesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListUpgradeHistories
//
// 查询实例大版本升级历史信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListUpgradeHistories(request *model.ListUpgradeHistoriesRequest) (*model.ListUpgradeHistoriesResponse, error) {
	requestDef := GenReqDefForListUpgradeHistories()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListUpgradeHistoriesResponse), nil
	}
}

// ListUpgradeHistoriesInvoker
func (c *RdsClient) ListUpgradeHistoriesInvoker(request *model.ListUpgradeHistoriesRequest) *ListUpgradeHistoriesInvoker {
	requestDef := GenReqDefForListUpgradeHistories()
	return &ListUpgradeHistoriesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListXellogFiles 查询扩展日志文件列表
//
// 查询扩展日志文件列表。
// 查询扩展日志文件列表，可以调用接口/v3/{project_id}/instances/{instance_id}/xellog-download 获取扩展日志下载链接
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListXellogFiles(request *model.ListXellogFilesRequest) (*model.ListXellogFilesResponse, error) {
	requestDef := GenReqDefForListXellogFiles()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListXellogFilesResponse), nil
	}
}

// ListXellogFilesInvoker 查询扩展日志文件列表
func (c *RdsClient) ListXellogFilesInvoker(request *model.ListXellogFilesRequest) *ListXellogFilesInvoker {
	requestDef := GenReqDefForListXellogFiles()
	return &ListXellogFilesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// MigrateFollower 迁移主备实例的备机
//
// 迁移主备实例的备机
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) MigrateFollower(request *model.MigrateFollowerRequest) (*model.MigrateFollowerResponse, error) {
	requestDef := GenReqDefForMigrateFollower()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.MigrateFollowerResponse), nil
	}
}

// MigrateFollowerInvoker 迁移主备实例的备机
func (c *RdsClient) MigrateFollowerInvoker(request *model.MigrateFollowerRequest) *MigrateFollowerInvoker {
	requestDef := GenReqDefForMigrateFollower()
	return &MigrateFollowerInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ModifyPostgresqlHbaConf 修改pg_hba.conf文件的单个或多个配置
//
// 修改/新增pg_hba.conf文件的单个或多个配置，以priority做唯一标识，priority不存在的新增，存在的修改
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ModifyPostgresqlHbaConf(request *model.ModifyPostgresqlHbaConfRequest) (*model.ModifyPostgresqlHbaConfResponse, error) {
	requestDef := GenReqDefForModifyPostgresqlHbaConf()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ModifyPostgresqlHbaConfResponse), nil
	}
}

// ModifyPostgresqlHbaConfInvoker 修改pg_hba.conf文件的单个或多个配置
func (c *RdsClient) ModifyPostgresqlHbaConfInvoker(request *model.ModifyPostgresqlHbaConfRequest) *ModifyPostgresqlHbaConfInvoker {
	requestDef := GenReqDefForModifyPostgresqlHbaConf()
	return &ModifyPostgresqlHbaConfInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RestoreExistInstance 恢复到已有实例
//
// 恢复到已有实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) RestoreExistInstance(request *model.RestoreExistInstanceRequest) (*model.RestoreExistInstanceResponse, error) {
	requestDef := GenReqDefForRestoreExistInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreExistInstanceResponse), nil
	}
}

// RestoreExistInstanceInvoker 恢复到已有实例
func (c *RdsClient) RestoreExistInstanceInvoker(request *model.RestoreExistInstanceRequest) *RestoreExistInstanceInvoker {
	requestDef := GenReqDefForRestoreExistInstance()
	return &RestoreExistInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RestoreTables 表级时间点恢复(MySQL)
//
// 表级时间点恢复(MySQL)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) RestoreTables(request *model.RestoreTablesRequest) (*model.RestoreTablesResponse, error) {
	requestDef := GenReqDefForRestoreTables()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreTablesResponse), nil
	}
}

// RestoreTablesInvoker 表级时间点恢复(MySQL)
func (c *RdsClient) RestoreTablesInvoker(request *model.RestoreTablesRequest) *RestoreTablesInvoker {
	requestDef := GenReqDefForRestoreTables()
	return &RestoreTablesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RestoreTablesNew 表级时间点恢复(MySQL)
//
// 表级时间点恢复(MySQL)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) RestoreTablesNew(request *model.RestoreTablesNewRequest) (*model.RestoreTablesNewResponse, error) {
	requestDef := GenReqDefForRestoreTablesNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreTablesNewResponse), nil
	}
}

// RestoreTablesNewInvoker 表级时间点恢复(MySQL)
func (c *RdsClient) RestoreTablesNewInvoker(request *model.RestoreTablesNewRequest) *RestoreTablesNewInvoker {
	requestDef := GenReqDefForRestoreTablesNew()
	return &RestoreTablesNewInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RestoreToExistingInstance 恢复到已有实例
//
// 恢复到已有实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) RestoreToExistingInstance(request *model.RestoreToExistingInstanceRequest) (*model.RestoreToExistingInstanceResponse, error) {
	requestDef := GenReqDefForRestoreToExistingInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RestoreToExistingInstanceResponse), nil
	}
}

// RestoreToExistingInstanceInvoker 恢复到已有实例
func (c *RdsClient) RestoreToExistingInstanceInvoker(request *model.RestoreToExistingInstanceRequest) *RestoreToExistingInstanceInvoker {
	requestDef := GenReqDefForRestoreToExistingInstance()
	return &RestoreToExistingInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetAuditlogPolicy 设置审计日志策略
//
// 设置审计日志策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetAuditlogPolicy(request *model.SetAuditlogPolicyRequest) (*model.SetAuditlogPolicyResponse, error) {
	requestDef := GenReqDefForSetAuditlogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetAuditlogPolicyResponse), nil
	}
}

// SetAuditlogPolicyInvoker 设置审计日志策略
func (c *RdsClient) SetAuditlogPolicyInvoker(request *model.SetAuditlogPolicyRequest) *SetAuditlogPolicyInvoker {
	requestDef := GenReqDefForSetAuditlogPolicy()
	return &SetAuditlogPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetAutoEnlargePolicy 设置自动扩容策略
//
// 设置实例存储空间自动扩容策略，按扩容量扣除存储费用。
// 可用存储空间小于设置值或者10GB时，自动扩容当前存储空间的15%（非10倍数向上取整，账户余额不足，会导致自动扩容失败）。
// 设置只读实例自动扩容与主实例自动扩容互不影响，对只读实例设置自动扩容时，可选择大于或等于主实例的存储空间。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetAutoEnlargePolicy(request *model.SetAutoEnlargePolicyRequest) (*model.SetAutoEnlargePolicyResponse, error) {
	requestDef := GenReqDefForSetAutoEnlargePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetAutoEnlargePolicyResponse), nil
	}
}

// SetAutoEnlargePolicyInvoker 设置自动扩容策略
func (c *RdsClient) SetAutoEnlargePolicyInvoker(request *model.SetAutoEnlargePolicyRequest) *SetAutoEnlargePolicyInvoker {
	requestDef := GenReqDefForSetAutoEnlargePolicy()
	return &SetAutoEnlargePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetBackupPolicy 设置自动备份策略
//
// 设置自动备份策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetBackupPolicy(request *model.SetBackupPolicyRequest) (*model.SetBackupPolicyResponse, error) {
	requestDef := GenReqDefForSetBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetBackupPolicyResponse), nil
	}
}

// SetBackupPolicyInvoker 设置自动备份策略
func (c *RdsClient) SetBackupPolicyInvoker(request *model.SetBackupPolicyRequest) *SetBackupPolicyInvoker {
	requestDef := GenReqDefForSetBackupPolicy()
	return &SetBackupPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetBinlogClearPolicy 设置binlog本地保留时长
//
// 修改指定实例的binlog本地保留时长。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetBinlogClearPolicy(request *model.SetBinlogClearPolicyRequest) (*model.SetBinlogClearPolicyResponse, error) {
	requestDef := GenReqDefForSetBinlogClearPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetBinlogClearPolicyResponse), nil
	}
}

// SetBinlogClearPolicyInvoker 设置binlog本地保留时长
func (c *RdsClient) SetBinlogClearPolicyInvoker(request *model.SetBinlogClearPolicyRequest) *SetBinlogClearPolicyInvoker {
	requestDef := GenReqDefForSetBinlogClearPolicy()
	return &SetBinlogClearPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetOffSiteBackupPolicy 设置跨区域备份策略
//
// 设置跨区域备份策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetOffSiteBackupPolicy(request *model.SetOffSiteBackupPolicyRequest) (*model.SetOffSiteBackupPolicyResponse, error) {
	requestDef := GenReqDefForSetOffSiteBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetOffSiteBackupPolicyResponse), nil
	}
}

// SetOffSiteBackupPolicyInvoker 设置跨区域备份策略
func (c *RdsClient) SetOffSiteBackupPolicyInvoker(request *model.SetOffSiteBackupPolicyRequest) *SetOffSiteBackupPolicyInvoker {
	requestDef := GenReqDefForSetOffSiteBackupPolicy()
	return &SetOffSiteBackupPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetSecondLevelMonitor 设置秒级监控策略
//
// 设置实例秒级监控策略，约五分钟后生效，对于1秒监控和5秒监控，计费方式为按需付费（每小时扣费一次）。
// 设置只读实例秒级监控与主实例互不影响。
// 规格变更到4U以下的实例，秒级监控功能会自动关闭。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetSecondLevelMonitor(request *model.SetSecondLevelMonitorRequest) (*model.SetSecondLevelMonitorResponse, error) {
	requestDef := GenReqDefForSetSecondLevelMonitor()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetSecondLevelMonitorResponse), nil
	}
}

// SetSecondLevelMonitorInvoker 设置秒级监控策略
func (c *RdsClient) SetSecondLevelMonitorInvoker(request *model.SetSecondLevelMonitorRequest) *SetSecondLevelMonitorInvoker {
	requestDef := GenReqDefForSetSecondLevelMonitor()
	return &SetSecondLevelMonitorInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetSecurityGroup 修改安全组
//
// 修改安全组
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetSecurityGroup(request *model.SetSecurityGroupRequest) (*model.SetSecurityGroupResponse, error) {
	requestDef := GenReqDefForSetSecurityGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetSecurityGroupResponse), nil
	}
}

// SetSecurityGroupInvoker 修改安全组
func (c *RdsClient) SetSecurityGroupInvoker(request *model.SetSecurityGroupRequest) *SetSecurityGroupInvoker {
	requestDef := GenReqDefForSetSecurityGroup()
	return &SetSecurityGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetSensitiveSlowLog 慢日志敏感信息的开关
//
// # V3慢日志敏感信息的开关
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetSensitiveSlowLog(request *model.SetSensitiveSlowLogRequest) (*model.SetSensitiveSlowLogResponse, error) {
	requestDef := GenReqDefForSetSensitiveSlowLog()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetSensitiveSlowLogResponse), nil
	}
}

// SetSensitiveSlowLogInvoker 慢日志敏感信息的开关
func (c *RdsClient) SetSensitiveSlowLogInvoker(request *model.SetSensitiveSlowLogRequest) *SetSensitiveSlowLogInvoker {
	requestDef := GenReqDefForSetSensitiveSlowLog()
	return &SetSensitiveSlowLogInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAuditlogDownloadLink 生成审计日志下载链接
//
// 生成审计日志下载链接。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowAuditlogDownloadLink(request *model.ShowAuditlogDownloadLinkRequest) (*model.ShowAuditlogDownloadLinkResponse, error) {
	requestDef := GenReqDefForShowAuditlogDownloadLink()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAuditlogDownloadLinkResponse), nil
	}
}

// ShowAuditlogDownloadLinkInvoker 生成审计日志下载链接
func (c *RdsClient) ShowAuditlogDownloadLinkInvoker(request *model.ShowAuditlogDownloadLinkRequest) *ShowAuditlogDownloadLinkInvoker {
	requestDef := GenReqDefForShowAuditlogDownloadLink()
	return &ShowAuditlogDownloadLinkInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAuditlogPolicy 查询审计日志策略
//
// 查询审计日志策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowAuditlogPolicy(request *model.ShowAuditlogPolicyRequest) (*model.ShowAuditlogPolicyResponse, error) {
	requestDef := GenReqDefForShowAuditlogPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAuditlogPolicyResponse), nil
	}
}

// ShowAuditlogPolicyInvoker 查询审计日志策略
func (c *RdsClient) ShowAuditlogPolicyInvoker(request *model.ShowAuditlogPolicyRequest) *ShowAuditlogPolicyInvoker {
	requestDef := GenReqDefForShowAuditlogPolicy()
	return &ShowAuditlogPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAutoEnlargePolicy 查询自动扩容策略
//
// 查询实例存储空间自动扩容策略
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowAutoEnlargePolicy(request *model.ShowAutoEnlargePolicyRequest) (*model.ShowAutoEnlargePolicyResponse, error) {
	requestDef := GenReqDefForShowAutoEnlargePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAutoEnlargePolicyResponse), nil
	}
}

// ShowAutoEnlargePolicyInvoker 查询自动扩容策略
func (c *RdsClient) ShowAutoEnlargePolicyInvoker(request *model.ShowAutoEnlargePolicyRequest) *ShowAutoEnlargePolicyInvoker {
	requestDef := GenReqDefForShowAutoEnlargePolicy()
	return &ShowAutoEnlargePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAvailableVersion
//
// 查询实例可升级的目标版本
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowAvailableVersion(request *model.ShowAvailableVersionRequest) (*model.ShowAvailableVersionResponse, error) {
	requestDef := GenReqDefForShowAvailableVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAvailableVersionResponse), nil
	}
}

// ShowAvailableVersionInvoker
func (c *RdsClient) ShowAvailableVersionInvoker(request *model.ShowAvailableVersionRequest) *ShowAvailableVersionInvoker {
	requestDef := GenReqDefForShowAvailableVersion()
	return &ShowAvailableVersionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowBackupDownloadLink 获取备份下载链接
//
// 获取备份下载链接。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowBackupDownloadLink(request *model.ShowBackupDownloadLinkRequest) (*model.ShowBackupDownloadLinkResponse, error) {
	requestDef := GenReqDefForShowBackupDownloadLink()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBackupDownloadLinkResponse), nil
	}
}

// ShowBackupDownloadLinkInvoker 获取备份下载链接
func (c *RdsClient) ShowBackupDownloadLinkInvoker(request *model.ShowBackupDownloadLinkRequest) *ShowBackupDownloadLinkInvoker {
	requestDef := GenReqDefForShowBackupDownloadLink()
	return &ShowBackupDownloadLinkInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowBackupPolicy 查询自动备份策略
//
// 查询自动备份策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowBackupPolicy(request *model.ShowBackupPolicyRequest) (*model.ShowBackupPolicyResponse, error) {
	requestDef := GenReqDefForShowBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBackupPolicyResponse), nil
	}
}

// ShowBackupPolicyInvoker 查询自动备份策略
func (c *RdsClient) ShowBackupPolicyInvoker(request *model.ShowBackupPolicyRequest) *ShowBackupPolicyInvoker {
	requestDef := GenReqDefForShowBackupPolicy()
	return &ShowBackupPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowBinlogClearPolicy 获取binlog本地保留时长
//
// 查寻指定实例的binlog本地保留时长。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowBinlogClearPolicy(request *model.ShowBinlogClearPolicyRequest) (*model.ShowBinlogClearPolicyResponse, error) {
	requestDef := GenReqDefForShowBinlogClearPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBinlogClearPolicyResponse), nil
	}
}

// ShowBinlogClearPolicyInvoker 获取binlog本地保留时长
func (c *RdsClient) ShowBinlogClearPolicyInvoker(request *model.ShowBinlogClearPolicyRequest) *ShowBinlogClearPolicyInvoker {
	requestDef := GenReqDefForShowBinlogClearPolicy()
	return &ShowBinlogClearPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowConfiguration 获取指定参数模板的参数
//
// 获取指定参数模板的参数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowConfiguration(request *model.ShowConfigurationRequest) (*model.ShowConfigurationResponse, error) {
	requestDef := GenReqDefForShowConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowConfigurationResponse), nil
	}
}

// ShowConfigurationInvoker 获取指定参数模板的参数
func (c *RdsClient) ShowConfigurationInvoker(request *model.ShowConfigurationRequest) *ShowConfigurationInvoker {
	requestDef := GenReqDefForShowConfiguration()
	return &ShowConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDnsName 查询实例ipv6域名。
//
// 查询实例ipv6域名。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowDnsName(request *model.ShowDnsNameRequest) (*model.ShowDnsNameResponse, error) {
	requestDef := GenReqDefForShowDnsName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDnsNameResponse), nil
	}
}

// ShowDnsNameInvoker 查询实例ipv6域名。
func (c *RdsClient) ShowDnsNameInvoker(request *model.ShowDnsNameRequest) *ShowDnsNameInvoker {
	requestDef := GenReqDefForShowDnsName()
	return &ShowDnsNameInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDomainName
//
// 查询实例ipv4域名
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowDomainName(request *model.ShowDomainNameRequest) (*model.ShowDomainNameResponse, error) {
	requestDef := GenReqDefForShowDomainName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDomainNameResponse), nil
	}
}

// ShowDomainNameInvoker
func (c *RdsClient) ShowDomainNameInvoker(request *model.ShowDomainNameRequest) *ShowDomainNameInvoker {
	requestDef := GenReqDefForShowDomainName()
	return &ShowDomainNameInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowDrReplicaStatus 查询跨云容灾复制状态
//
// 建立跨云容灾关系后，查询主实例和灾备实例间的复制状态及延迟。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowDrReplicaStatus(request *model.ShowDrReplicaStatusRequest) (*model.ShowDrReplicaStatusResponse, error) {
	requestDef := GenReqDefForShowDrReplicaStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowDrReplicaStatusResponse), nil
	}
}

// ShowDrReplicaStatusInvoker 查询跨云容灾复制状态
func (c *RdsClient) ShowDrReplicaStatusInvoker(request *model.ShowDrReplicaStatusRequest) *ShowDrReplicaStatusInvoker {
	requestDef := GenReqDefForShowDrReplicaStatus()
	return &ShowDrReplicaStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowInstanceConfiguration 获取指定实例的参数模板
//
// 获取指定实例的参数模板。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowInstanceConfiguration(request *model.ShowInstanceConfigurationRequest) (*model.ShowInstanceConfigurationResponse, error) {
	requestDef := GenReqDefForShowInstanceConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowInstanceConfigurationResponse), nil
	}
}

// ShowInstanceConfigurationInvoker 获取指定实例的参数模板
func (c *RdsClient) ShowInstanceConfigurationInvoker(request *model.ShowInstanceConfigurationRequest) *ShowInstanceConfigurationInvoker {
	requestDef := GenReqDefForShowInstanceConfiguration()
	return &ShowInstanceConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowOffSiteBackupPolicy 查询跨区域备份策略
//
// 查询跨区域备份策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowOffSiteBackupPolicy(request *model.ShowOffSiteBackupPolicyRequest) (*model.ShowOffSiteBackupPolicyResponse, error) {
	requestDef := GenReqDefForShowOffSiteBackupPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowOffSiteBackupPolicyResponse), nil
	}
}

// ShowOffSiteBackupPolicyInvoker 查询跨区域备份策略
func (c *RdsClient) ShowOffSiteBackupPolicyInvoker(request *model.ShowOffSiteBackupPolicyRequest) *ShowOffSiteBackupPolicyInvoker {
	requestDef := GenReqDefForShowOffSiteBackupPolicy()
	return &ShowOffSiteBackupPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowQuotas 查询配额
//
// 查询当前项目下资源配额情况。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowQuotas(request *model.ShowQuotasRequest) (*model.ShowQuotasResponse, error) {
	requestDef := GenReqDefForShowQuotas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowQuotasResponse), nil
	}
}

// ShowQuotasInvoker 查询配额
func (c *RdsClient) ShowQuotasInvoker(request *model.ShowQuotasRequest) *ShowQuotasInvoker {
	requestDef := GenReqDefForShowQuotas()
	return &ShowQuotasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRecyclePolicy 查询回收站的回收策略。
//
// 查询回收站的回收策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowRecyclePolicy(request *model.ShowRecyclePolicyRequest) (*model.ShowRecyclePolicyResponse, error) {
	requestDef := GenReqDefForShowRecyclePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRecyclePolicyResponse), nil
	}
}

// ShowRecyclePolicyInvoker 查询回收站的回收策略。
func (c *RdsClient) ShowRecyclePolicyInvoker(request *model.ShowRecyclePolicyRequest) *ShowRecyclePolicyInvoker {
	requestDef := GenReqDefForShowRecyclePolicy()
	return &ShowRecyclePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowReplicationStatus 获取实例的复制状态。
//
// 获取实例的复制状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowReplicationStatus(request *model.ShowReplicationStatusRequest) (*model.ShowReplicationStatusResponse, error) {
	requestDef := GenReqDefForShowReplicationStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowReplicationStatusResponse), nil
	}
}

// ShowReplicationStatusInvoker 获取实例的复制状态。
func (c *RdsClient) ShowReplicationStatusInvoker(request *model.ShowReplicationStatusRequest) *ShowReplicationStatusInvoker {
	requestDef := GenReqDefForShowReplicationStatus()
	return &ShowReplicationStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSecondLevelMonitoring 查询秒级监控策略
//
// 查询实例秒级监控策略
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowSecondLevelMonitoring(request *model.ShowSecondLevelMonitoringRequest) (*model.ShowSecondLevelMonitoringResponse, error) {
	requestDef := GenReqDefForShowSecondLevelMonitoring()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSecondLevelMonitoringResponse), nil
	}
}

// ShowSecondLevelMonitoringInvoker 查询秒级监控策略
func (c *RdsClient) ShowSecondLevelMonitoringInvoker(request *model.ShowSecondLevelMonitoringRequest) *ShowSecondLevelMonitoringInvoker {
	requestDef := GenReqDefForShowSecondLevelMonitoring()
	return &ShowSecondLevelMonitoringInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTdeStatus 根据实例id查询sqlserver TDE状态
//
// 根据实例id查询sqlserver TDE状态
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowTdeStatus(request *model.ShowTdeStatusRequest) (*model.ShowTdeStatusResponse, error) {
	requestDef := GenReqDefForShowTdeStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTdeStatusResponse), nil
	}
}

// ShowTdeStatusInvoker 根据实例id查询sqlserver TDE状态
func (c *RdsClient) ShowTdeStatusInvoker(request *model.ShowTdeStatusRequest) *ShowTdeStatusInvoker {
	requestDef := GenReqDefForShowTdeStatus()
	return &ShowTdeStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowUpgradeDbMajorVersionStatus
//
// 查询大版本检查状态或升级状态。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowUpgradeDbMajorVersionStatus(request *model.ShowUpgradeDbMajorVersionStatusRequest) (*model.ShowUpgradeDbMajorVersionStatusResponse, error) {
	requestDef := GenReqDefForShowUpgradeDbMajorVersionStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowUpgradeDbMajorVersionStatusResponse), nil
	}
}

// ShowUpgradeDbMajorVersionStatusInvoker
func (c *RdsClient) ShowUpgradeDbMajorVersionStatusInvoker(request *model.ShowUpgradeDbMajorVersionStatusRequest) *ShowUpgradeDbMajorVersionStatusInvoker {
	requestDef := GenReqDefForShowUpgradeDbMajorVersionStatus()
	return &ShowUpgradeDbMajorVersionStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartFailover 手动倒换主备
//
// 手动倒换主备.
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StartFailover(request *model.StartFailoverRequest) (*model.StartFailoverResponse, error) {
	requestDef := GenReqDefForStartFailover()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartFailoverResponse), nil
	}
}

// StartFailoverInvoker 手动倒换主备
func (c *RdsClient) StartFailoverInvoker(request *model.StartFailoverRequest) *StartFailoverInvoker {
	requestDef := GenReqDefForStartFailover()
	return &StartFailoverInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartInstanceEnlargeVolumeAction 扩容数据库实例的磁盘空间
//
// 扩容数据库实例的磁盘空间。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StartInstanceEnlargeVolumeAction(request *model.StartInstanceEnlargeVolumeActionRequest) (*model.StartInstanceEnlargeVolumeActionResponse, error) {
	requestDef := GenReqDefForStartInstanceEnlargeVolumeAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceEnlargeVolumeActionResponse), nil
	}
}

// StartInstanceEnlargeVolumeActionInvoker 扩容数据库实例的磁盘空间
func (c *RdsClient) StartInstanceEnlargeVolumeActionInvoker(request *model.StartInstanceEnlargeVolumeActionRequest) *StartInstanceEnlargeVolumeActionInvoker {
	requestDef := GenReqDefForStartInstanceEnlargeVolumeAction()
	return &StartInstanceEnlargeVolumeActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartInstanceRestartAction 重启数据库实例
//
// 重启数据库实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StartInstanceRestartAction(request *model.StartInstanceRestartActionRequest) (*model.StartInstanceRestartActionResponse, error) {
	requestDef := GenReqDefForStartInstanceRestartAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceRestartActionResponse), nil
	}
}

// StartInstanceRestartActionInvoker 重启数据库实例
func (c *RdsClient) StartInstanceRestartActionInvoker(request *model.StartInstanceRestartActionRequest) *StartInstanceRestartActionInvoker {
	requestDef := GenReqDefForStartInstanceRestartAction()
	return &StartInstanceRestartActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartInstanceSingleToHaAction 单机转主备实例
//
// 单机转主备实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StartInstanceSingleToHaAction(request *model.StartInstanceSingleToHaActionRequest) (*model.StartInstanceSingleToHaActionResponse, error) {
	requestDef := GenReqDefForStartInstanceSingleToHaAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartInstanceSingleToHaActionResponse), nil
	}
}

// StartInstanceSingleToHaActionInvoker 单机转主备实例
func (c *RdsClient) StartInstanceSingleToHaActionInvoker(request *model.StartInstanceSingleToHaActionRequest) *StartInstanceSingleToHaActionInvoker {
	requestDef := GenReqDefForStartInstanceSingleToHaAction()
	return &StartInstanceSingleToHaActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartRecyclePolicy 设置回收站策略
//
// 设置回收站策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StartRecyclePolicy(request *model.StartRecyclePolicyRequest) (*model.StartRecyclePolicyResponse, error) {
	requestDef := GenReqDefForStartRecyclePolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartRecyclePolicyResponse), nil
	}
}

// StartRecyclePolicyInvoker 设置回收站策略
func (c *RdsClient) StartRecyclePolicyInvoker(request *model.StartRecyclePolicyRequest) *StartRecyclePolicyInvoker {
	requestDef := GenReqDefForStartRecyclePolicy()
	return &StartRecyclePolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartResizeFlavorAction 变更数据库实例的规格
//
// 变更数据库实例的规格。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StartResizeFlavorAction(request *model.StartResizeFlavorActionRequest) (*model.StartResizeFlavorActionResponse, error) {
	requestDef := GenReqDefForStartResizeFlavorAction()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartResizeFlavorActionResponse), nil
	}
}

// StartResizeFlavorActionInvoker 变更数据库实例的规格
func (c *RdsClient) StartResizeFlavorActionInvoker(request *model.StartResizeFlavorActionRequest) *StartResizeFlavorActionInvoker {
	requestDef := GenReqDefForStartResizeFlavorAction()
	return &StartResizeFlavorActionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartupInstance 开启实例
//
// 停止实例以节省费用，在停止数据库实例后，支持手动重新开启实例。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StartupInstance(request *model.StartupInstanceRequest) (*model.StartupInstanceResponse, error) {
	requestDef := GenReqDefForStartupInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartupInstanceResponse), nil
	}
}

// StartupInstanceInvoker 开启实例
func (c *RdsClient) StartupInstanceInvoker(request *model.StartupInstanceRequest) *StartupInstanceInvoker {
	requestDef := GenReqDefForStartupInstance()
	return &StartupInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopInstance 停止实例
//
// 实例进行关机，通过暂时停止按需实例以节省费用，实例默认停止七天。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StopInstance(request *model.StopInstanceRequest) (*model.StopInstanceResponse, error) {
	requestDef := GenReqDefForStopInstance()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopInstanceResponse), nil
	}
}

// StopInstanceInvoker 停止实例
func (c *RdsClient) StopInstanceInvoker(request *model.StopInstanceRequest) *StopInstanceInvoker {
	requestDef := GenReqDefForStopInstance()
	return &StopInstanceInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SwitchSsl 设置SSL数据加密
//
// 设置SSL数据加密。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SwitchSsl(request *model.SwitchSslRequest) (*model.SwitchSslResponse, error) {
	requestDef := GenReqDefForSwitchSsl()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SwitchSslResponse), nil
	}
}

// SwitchSslInvoker 设置SSL数据加密
func (c *RdsClient) SwitchSslInvoker(request *model.SwitchSslRequest) *SwitchSslInvoker {
	requestDef := GenReqDefForSwitchSsl()
	return &SwitchSslInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateConfiguration 修改参数模板参数
//
// 修改参数模板参数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateConfiguration(request *model.UpdateConfigurationRequest) (*model.UpdateConfigurationResponse, error) {
	requestDef := GenReqDefForUpdateConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateConfigurationResponse), nil
	}
}

// UpdateConfigurationInvoker 修改参数模板参数
func (c *RdsClient) UpdateConfigurationInvoker(request *model.UpdateConfigurationRequest) *UpdateConfigurationInvoker {
	requestDef := GenReqDefForUpdateConfiguration()
	return &UpdateConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDataIp 修改内网地址
//
// 修改内网地址
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateDataIp(request *model.UpdateDataIpRequest) (*model.UpdateDataIpResponse, error) {
	requestDef := GenReqDefForUpdateDataIp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDataIpResponse), nil
	}
}

// UpdateDataIpInvoker 修改内网地址
func (c *RdsClient) UpdateDataIpInvoker(request *model.UpdateDataIpRequest) *UpdateDataIpInvoker {
	requestDef := GenReqDefForUpdateDataIp()
	return &UpdateDataIpInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDnsName 修改域名
//
// 修改域名
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateDnsName(request *model.UpdateDnsNameRequest) (*model.UpdateDnsNameResponse, error) {
	requestDef := GenReqDefForUpdateDnsName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDnsNameResponse), nil
	}
}

// UpdateDnsNameInvoker 修改域名
func (c *RdsClient) UpdateDnsNameInvoker(request *model.UpdateDnsNameRequest) *UpdateDnsNameInvoker {
	requestDef := GenReqDefForUpdateDnsName()
	return &UpdateDnsNameInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateInstanceConfiguration 修改指定实例的参数
//
// 修改指定实例的参数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateInstanceConfiguration(request *model.UpdateInstanceConfigurationRequest) (*model.UpdateInstanceConfigurationResponse, error) {
	requestDef := GenReqDefForUpdateInstanceConfiguration()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceConfigurationResponse), nil
	}
}

// UpdateInstanceConfigurationInvoker 修改指定实例的参数
func (c *RdsClient) UpdateInstanceConfigurationInvoker(request *model.UpdateInstanceConfigurationRequest) *UpdateInstanceConfigurationInvoker {
	requestDef := GenReqDefForUpdateInstanceConfiguration()
	return &UpdateInstanceConfigurationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateInstanceConfigurationAsync 修改指定实例的参数
//
// 修改指定实例的参数。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateInstanceConfigurationAsync(request *model.UpdateInstanceConfigurationAsyncRequest) (*model.UpdateInstanceConfigurationAsyncResponse, error) {
	requestDef := GenReqDefForUpdateInstanceConfigurationAsync()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceConfigurationAsyncResponse), nil
	}
}

// UpdateInstanceConfigurationAsyncInvoker 修改指定实例的参数
func (c *RdsClient) UpdateInstanceConfigurationAsyncInvoker(request *model.UpdateInstanceConfigurationAsyncRequest) *UpdateInstanceConfigurationAsyncInvoker {
	requestDef := GenReqDefForUpdateInstanceConfigurationAsync()
	return &UpdateInstanceConfigurationAsyncInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateInstanceName 修改实例名称
//
// 修改实例名称。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateInstanceName(request *model.UpdateInstanceNameRequest) (*model.UpdateInstanceNameResponse, error) {
	requestDef := GenReqDefForUpdateInstanceName()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateInstanceNameResponse), nil
	}
}

// UpdateInstanceNameInvoker 修改实例名称
func (c *RdsClient) UpdateInstanceNameInvoker(request *model.UpdateInstanceNameRequest) *UpdateInstanceNameInvoker {
	requestDef := GenReqDefForUpdateInstanceName()
	return &UpdateInstanceNameInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePort 修改数据库端口
//
// 修改数据库端口
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdatePort(request *model.UpdatePortRequest) (*model.UpdatePortResponse, error) {
	requestDef := GenReqDefForUpdatePort()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePortResponse), nil
	}
}

// UpdatePortInvoker 修改数据库端口
func (c *RdsClient) UpdatePortInvoker(request *model.UpdatePortRequest) *UpdatePortInvoker {
	requestDef := GenReqDefForUpdatePort()
	return &UpdatePortInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePostgresqlInstanceAlias 修改实例备注信息
//
// 修改指定数据库实例的备注信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdatePostgresqlInstanceAlias(request *model.UpdatePostgresqlInstanceAliasRequest) (*model.UpdatePostgresqlInstanceAliasResponse, error) {
	requestDef := GenReqDefForUpdatePostgresqlInstanceAlias()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePostgresqlInstanceAliasResponse), nil
	}
}

// UpdatePostgresqlInstanceAliasInvoker 修改实例备注信息
func (c *RdsClient) UpdatePostgresqlInstanceAliasInvoker(request *model.UpdatePostgresqlInstanceAliasRequest) *UpdatePostgresqlInstanceAliasInvoker {
	requestDef := GenReqDefForUpdatePostgresqlInstanceAlias()
	return &UpdatePostgresqlInstanceAliasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTdeStatus sqlserverTDE开关
//
// sqlserverTDE开关。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateTdeStatus(request *model.UpdateTdeStatusRequest) (*model.UpdateTdeStatusResponse, error) {
	requestDef := GenReqDefForUpdateTdeStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTdeStatusResponse), nil
	}
}

// UpdateTdeStatusInvoker sqlserverTDE开关
func (c *RdsClient) UpdateTdeStatusInvoker(request *model.UpdateTdeStatusRequest) *UpdateTdeStatusInvoker {
	requestDef := GenReqDefForUpdateTdeStatus()
	return &UpdateTdeStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeDbMajorVersion
//
// PostgreSQL数据库升级大版本。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpgradeDbMajorVersion(request *model.UpgradeDbMajorVersionRequest) (*model.UpgradeDbMajorVersionResponse, error) {
	requestDef := GenReqDefForUpgradeDbMajorVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeDbMajorVersionResponse), nil
	}
}

// UpgradeDbMajorVersionInvoker
func (c *RdsClient) UpgradeDbMajorVersionInvoker(request *model.UpgradeDbMajorVersionRequest) *UpgradeDbMajorVersionInvoker {
	requestDef := GenReqDefForUpgradeDbMajorVersion()
	return &UpgradeDbMajorVersionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeDbMajorVersionPreCheck
//
// 大版本升级前进行升级检查。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpgradeDbMajorVersionPreCheck(request *model.UpgradeDbMajorVersionPreCheckRequest) (*model.UpgradeDbMajorVersionPreCheckResponse, error) {
	requestDef := GenReqDefForUpgradeDbMajorVersionPreCheck()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeDbMajorVersionPreCheckResponse), nil
	}
}

// UpgradeDbMajorVersionPreCheckInvoker
func (c *RdsClient) UpgradeDbMajorVersionPreCheckInvoker(request *model.UpgradeDbMajorVersionPreCheckRequest) *UpgradeDbMajorVersionPreCheckInvoker {
	requestDef := GenReqDefForUpgradeDbMajorVersionPreCheck()
	return &UpgradeDbMajorVersionPreCheckInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeDbVersion 升级内核小版本
//
// 对实例进行小版本升级。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpgradeDbVersion(request *model.UpgradeDbVersionRequest) (*model.UpgradeDbVersionResponse, error) {
	requestDef := GenReqDefForUpgradeDbVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeDbVersionResponse), nil
	}
}

// UpgradeDbVersionInvoker 升级内核小版本
func (c *RdsClient) UpgradeDbVersionInvoker(request *model.UpgradeDbVersionRequest) *UpgradeDbVersionInvoker {
	requestDef := GenReqDefForUpgradeDbVersion()
	return &UpgradeDbVersionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpgradeDbVersionNew 升级内核小版本
//
// 对实例进行小版本升级。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpgradeDbVersionNew(request *model.UpgradeDbVersionNewRequest) (*model.UpgradeDbVersionNewResponse, error) {
	requestDef := GenReqDefForUpgradeDbVersionNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpgradeDbVersionNewResponse), nil
	}
}

// UpgradeDbVersionNewInvoker 升级内核小版本
func (c *RdsClient) UpgradeDbVersionNewInvoker(request *model.UpgradeDbVersionNewRequest) *UpgradeDbVersionNewInvoker {
	requestDef := GenReqDefForUpgradeDbVersionNew()
	return &UpgradeDbVersionNewInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListApiVersion 查询API版本列表
//
// 查询API版本列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListApiVersion(request *model.ListApiVersionRequest) (*model.ListApiVersionResponse, error) {
	requestDef := GenReqDefForListApiVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListApiVersionResponse), nil
	}
}

// ListApiVersionInvoker 查询API版本列表
func (c *RdsClient) ListApiVersionInvoker(request *model.ListApiVersionRequest) *ListApiVersionInvoker {
	requestDef := GenReqDefForListApiVersion()
	return &ListApiVersionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListApiVersionNew 查询API版本列表
//
// 查询API版本列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListApiVersionNew(request *model.ListApiVersionNewRequest) (*model.ListApiVersionNewResponse, error) {
	requestDef := GenReqDefForListApiVersionNew()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListApiVersionNewResponse), nil
	}
}

// ListApiVersionNewInvoker 查询API版本列表
func (c *RdsClient) ListApiVersionNewInvoker(request *model.ListApiVersionNewRequest) *ListApiVersionNewInvoker {
	requestDef := GenReqDefForListApiVersionNew()
	return &ListApiVersionNewInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowApiVersion 查询指定的API版本信息
//
// 查询指定的API版本信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowApiVersion(request *model.ShowApiVersionRequest) (*model.ShowApiVersionResponse, error) {
	requestDef := GenReqDefForShowApiVersion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowApiVersionResponse), nil
	}
}

// ShowApiVersionInvoker 查询指定的API版本信息
func (c *RdsClient) ShowApiVersionInvoker(request *model.ShowApiVersionRequest) *ShowApiVersionInvoker {
	requestDef := GenReqDefForShowApiVersion()
	return &ShowApiVersionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AllowDbUserPrivilege 授权数据库帐号
//
// 授权数据库帐号。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) AllowDbUserPrivilege(request *model.AllowDbUserPrivilegeRequest) (*model.AllowDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowDbUserPrivilegeResponse), nil
	}
}

// AllowDbUserPrivilegeInvoker 授权数据库帐号
func (c *RdsClient) AllowDbUserPrivilegeInvoker(request *model.AllowDbUserPrivilegeRequest) *AllowDbUserPrivilegeInvoker {
	requestDef := GenReqDefForAllowDbUserPrivilege()
	return &AllowDbUserPrivilegeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDatabase 创建数据库
//
// 创建数据库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateDatabase(request *model.CreateDatabaseRequest) (*model.CreateDatabaseResponse, error) {
	requestDef := GenReqDefForCreateDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDatabaseResponse), nil
	}
}

// CreateDatabaseInvoker 创建数据库
func (c *RdsClient) CreateDatabaseInvoker(request *model.CreateDatabaseRequest) *CreateDatabaseInvoker {
	requestDef := GenReqDefForCreateDatabase()
	return &CreateDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDbUser 创建数据库用户
//
// 在指定实例中创建数据库帐号。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateDbUser(request *model.CreateDbUserRequest) (*model.CreateDbUserResponse, error) {
	requestDef := GenReqDefForCreateDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDbUserResponse), nil
	}
}

// CreateDbUserInvoker 创建数据库用户
func (c *RdsClient) CreateDbUserInvoker(request *model.CreateDbUserRequest) *CreateDbUserInvoker {
	requestDef := GenReqDefForCreateDbUser()
	return &CreateDbUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDatabase 删除数据库
//
// 删除数据库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteDatabase(request *model.DeleteDatabaseRequest) (*model.DeleteDatabaseResponse, error) {
	requestDef := GenReqDefForDeleteDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDatabaseResponse), nil
	}
}

// DeleteDatabaseInvoker 删除数据库
func (c *RdsClient) DeleteDatabaseInvoker(request *model.DeleteDatabaseRequest) *DeleteDatabaseInvoker {
	requestDef := GenReqDefForDeleteDatabase()
	return &DeleteDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDbUser 删除数据库用户
//
// 删除数据库用户。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteDbUser(request *model.DeleteDbUserRequest) (*model.DeleteDbUserResponse, error) {
	requestDef := GenReqDefForDeleteDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDbUserResponse), nil
	}
}

// DeleteDbUserInvoker 删除数据库用户
func (c *RdsClient) DeleteDbUserInvoker(request *model.DeleteDbUserRequest) *DeleteDbUserInvoker {
	requestDef := GenReqDefForDeleteDbUser()
	return &DeleteDbUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAuthorizedDatabases 查询指定用户的已授权数据库
//
// 查询指定用户的已授权数据库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListAuthorizedDatabases(request *model.ListAuthorizedDatabasesRequest) (*model.ListAuthorizedDatabasesResponse, error) {
	requestDef := GenReqDefForListAuthorizedDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedDatabasesResponse), nil
	}
}

// ListAuthorizedDatabasesInvoker 查询指定用户的已授权数据库
func (c *RdsClient) ListAuthorizedDatabasesInvoker(request *model.ListAuthorizedDatabasesRequest) *ListAuthorizedDatabasesInvoker {
	requestDef := GenReqDefForListAuthorizedDatabases()
	return &ListAuthorizedDatabasesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAuthorizedDbUsers 查询指定数据库的已授权用户
//
// 查询指定数据库的已授权用户。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListAuthorizedDbUsers(request *model.ListAuthorizedDbUsersRequest) (*model.ListAuthorizedDbUsersResponse, error) {
	requestDef := GenReqDefForListAuthorizedDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedDbUsersResponse), nil
	}
}

// ListAuthorizedDbUsersInvoker 查询指定数据库的已授权用户
func (c *RdsClient) ListAuthorizedDbUsersInvoker(request *model.ListAuthorizedDbUsersRequest) *ListAuthorizedDbUsersInvoker {
	requestDef := GenReqDefForListAuthorizedDbUsers()
	return &ListAuthorizedDbUsersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDatabases 查询数据库列表
//
// 查询数据库列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListDatabases(request *model.ListDatabasesRequest) (*model.ListDatabasesResponse, error) {
	requestDef := GenReqDefForListDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDatabasesResponse), nil
	}
}

// ListDatabasesInvoker 查询数据库列表
func (c *RdsClient) ListDatabasesInvoker(request *model.ListDatabasesRequest) *ListDatabasesInvoker {
	requestDef := GenReqDefForListDatabases()
	return &ListDatabasesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListDbUsers 查询数据库用户列表
//
// 查询数据库用户列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListDbUsers(request *model.ListDbUsersRequest) (*model.ListDbUsersResponse, error) {
	requestDef := GenReqDefForListDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListDbUsersResponse), nil
	}
}

// ListDbUsersInvoker 查询数据库用户列表
func (c *RdsClient) ListDbUsersInvoker(request *model.ListDbUsersRequest) *ListDbUsersInvoker {
	requestDef := GenReqDefForListDbUsers()
	return &ListDbUsersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ResetPwd 重置数据库密码
//
// 重置数据库密码.
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ResetPwd(request *model.ResetPwdRequest) (*model.ResetPwdResponse, error) {
	requestDef := GenReqDefForResetPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ResetPwdResponse), nil
	}
}

// ResetPwdInvoker 重置数据库密码
func (c *RdsClient) ResetPwdInvoker(request *model.ResetPwdRequest) *ResetPwdInvoker {
	requestDef := GenReqDefForResetPwd()
	return &ResetPwdInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// Revoke 解除数据库帐号权限
//
// 解除数据库帐号权限。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) Revoke(request *model.RevokeRequest) (*model.RevokeResponse, error) {
	requestDef := GenReqDefForRevoke()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeResponse), nil
	}
}

// RevokeInvoker 解除数据库帐号权限
func (c *RdsClient) RevokeInvoker(request *model.RevokeRequest) *RevokeInvoker {
	requestDef := GenReqDefForRevoke()
	return &RevokeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetDbUserPwd 设置数据库账号密码
//
// 设置数据库账号密码
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetDbUserPwd(request *model.SetDbUserPwdRequest) (*model.SetDbUserPwdResponse, error) {
	requestDef := GenReqDefForSetDbUserPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetDbUserPwdResponse), nil
	}
}

// SetDbUserPwdInvoker 设置数据库账号密码
func (c *RdsClient) SetDbUserPwdInvoker(request *model.SetDbUserPwdRequest) *SetDbUserPwdInvoker {
	requestDef := GenReqDefForSetDbUserPwd()
	return &SetDbUserPwdInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetReadOnlySwitch 设置数据库用户只读参数
//
// 根据业务需求，设置数据库用户只读
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetReadOnlySwitch(request *model.SetReadOnlySwitchRequest) (*model.SetReadOnlySwitchResponse, error) {
	requestDef := GenReqDefForSetReadOnlySwitch()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetReadOnlySwitchResponse), nil
	}
}

// SetReadOnlySwitchInvoker 设置数据库用户只读参数
func (c *RdsClient) SetReadOnlySwitchInvoker(request *model.SetReadOnlySwitchRequest) *SetReadOnlySwitchInvoker {
	requestDef := GenReqDefForSetReadOnlySwitch()
	return &SetReadOnlySwitchInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDatabase 修改指定实例的数据库备注
//
// 修改指定实例中的数据库备注。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateDatabase(request *model.UpdateDatabaseRequest) (*model.UpdateDatabaseResponse, error) {
	requestDef := GenReqDefForUpdateDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDatabaseResponse), nil
	}
}

// UpdateDatabaseInvoker 修改指定实例的数据库备注
func (c *RdsClient) UpdateDatabaseInvoker(request *model.UpdateDatabaseRequest) *UpdateDatabaseInvoker {
	requestDef := GenReqDefForUpdateDatabase()
	return &UpdateDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDbUserComment 修改数据库用户名备注
//
// 修改数据库用户名备注
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateDbUserComment(request *model.UpdateDbUserCommentRequest) (*model.UpdateDbUserCommentResponse, error) {
	requestDef := GenReqDefForUpdateDbUserComment()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDbUserCommentResponse), nil
	}
}

// UpdateDbUserCommentInvoker 修改数据库用户名备注
func (c *RdsClient) UpdateDbUserCommentInvoker(request *model.UpdateDbUserCommentRequest) *UpdateDbUserCommentInvoker {
	requestDef := GenReqDefForUpdateDbUserComment()
	return &UpdateDbUserCommentInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AllowDbPrivilege 授权数据库帐号
//
// 在指定实例的数据库中, 设置帐号的权限。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) AllowDbPrivilege(request *model.AllowDbPrivilegeRequest) (*model.AllowDbPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowDbPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowDbPrivilegeResponse), nil
	}
}

// AllowDbPrivilegeInvoker 授权数据库帐号
func (c *RdsClient) AllowDbPrivilegeInvoker(request *model.AllowDbPrivilegeRequest) *AllowDbPrivilegeInvoker {
	requestDef := GenReqDefForAllowDbPrivilege()
	return &AllowDbPrivilegeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeProxyScale 数据库代理规格变更
//
// 数据库代理实例进行规格变更。
//
// - 调用接口前，您需要了解API 认证鉴权。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ChangeProxyScale(request *model.ChangeProxyScaleRequest) (*model.ChangeProxyScaleResponse, error) {
	requestDef := GenReqDefForChangeProxyScale()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeProxyScaleResponse), nil
	}
}

// ChangeProxyScaleInvoker 数据库代理规格变更
func (c *RdsClient) ChangeProxyScaleInvoker(request *model.ChangeProxyScaleRequest) *ChangeProxyScaleInvoker {
	requestDef := GenReqDefForChangeProxyScale()
	return &ChangeProxyScaleInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ChangeTheDelayThreshold 修改读写分离阈值
//
// 修改指定实例的读写分离延时阈值。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ChangeTheDelayThreshold(request *model.ChangeTheDelayThresholdRequest) (*model.ChangeTheDelayThresholdResponse, error) {
	requestDef := GenReqDefForChangeTheDelayThreshold()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ChangeTheDelayThresholdResponse), nil
	}
}

// ChangeTheDelayThresholdInvoker 修改读写分离阈值
func (c *RdsClient) ChangeTheDelayThresholdInvoker(request *model.ChangeTheDelayThresholdRequest) *ChangeTheDelayThresholdInvoker {
	requestDef := GenReqDefForChangeTheDelayThreshold()
	return &ChangeTheDelayThresholdInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePostgresqlDatabase 创建数据库
//
// 在指定实例中创建数据库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreatePostgresqlDatabase(request *model.CreatePostgresqlDatabaseRequest) (*model.CreatePostgresqlDatabaseResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDatabaseResponse), nil
	}
}

// CreatePostgresqlDatabaseInvoker 创建数据库
func (c *RdsClient) CreatePostgresqlDatabaseInvoker(request *model.CreatePostgresqlDatabaseRequest) *CreatePostgresqlDatabaseInvoker {
	requestDef := GenReqDefForCreatePostgresqlDatabase()
	return &CreatePostgresqlDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePostgresqlDatabaseSchema 创建数据库SCHEMA
//
// 在指定实例的数据库中, 创建数据库schema。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreatePostgresqlDatabaseSchema(request *model.CreatePostgresqlDatabaseSchemaRequest) (*model.CreatePostgresqlDatabaseSchemaResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDatabaseSchema()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDatabaseSchemaResponse), nil
	}
}

// CreatePostgresqlDatabaseSchemaInvoker 创建数据库SCHEMA
func (c *RdsClient) CreatePostgresqlDatabaseSchemaInvoker(request *model.CreatePostgresqlDatabaseSchemaRequest) *CreatePostgresqlDatabaseSchemaInvoker {
	requestDef := GenReqDefForCreatePostgresqlDatabaseSchema()
	return &CreatePostgresqlDatabaseSchemaInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePostgresqlDbUser 创建数据库用户
//
// 在指定实例中创建数据库用户。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreatePostgresqlDbUser(request *model.CreatePostgresqlDbUserRequest) (*model.CreatePostgresqlDbUserResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlDbUserResponse), nil
	}
}

// CreatePostgresqlDbUserInvoker 创建数据库用户
func (c *RdsClient) CreatePostgresqlDbUserInvoker(request *model.CreatePostgresqlDbUserRequest) *CreatePostgresqlDbUserInvoker {
	requestDef := GenReqDefForCreatePostgresqlDbUser()
	return &CreatePostgresqlDbUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreatePostgresqlExtension 创建插件
//
// 在指定数据库上创建插件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreatePostgresqlExtension(request *model.CreatePostgresqlExtensionRequest) (*model.CreatePostgresqlExtensionResponse, error) {
	requestDef := GenReqDefForCreatePostgresqlExtension()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreatePostgresqlExtensionResponse), nil
	}
}

// CreatePostgresqlExtensionInvoker 创建插件
func (c *RdsClient) CreatePostgresqlExtensionInvoker(request *model.CreatePostgresqlExtensionRequest) *CreatePostgresqlExtensionInvoker {
	requestDef := GenReqDefForCreatePostgresqlExtension()
	return &CreatePostgresqlExtensionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeletePostgresqlDatabase 删除数据库
//
// 删除数据库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeletePostgresqlDatabase(request *model.DeletePostgresqlDatabaseRequest) (*model.DeletePostgresqlDatabaseResponse, error) {
	requestDef := GenReqDefForDeletePostgresqlDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeletePostgresqlDatabaseResponse), nil
	}
}

// DeletePostgresqlDatabaseInvoker 删除数据库
func (c *RdsClient) DeletePostgresqlDatabaseInvoker(request *model.DeletePostgresqlDatabaseRequest) *DeletePostgresqlDatabaseInvoker {
	requestDef := GenReqDefForDeletePostgresqlDatabase()
	return &DeletePostgresqlDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeletePostgresqlDbUser 删除数据库用户
//
// 删除数据库用户。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeletePostgresqlDbUser(request *model.DeletePostgresqlDbUserRequest) (*model.DeletePostgresqlDbUserResponse, error) {
	requestDef := GenReqDefForDeletePostgresqlDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeletePostgresqlDbUserResponse), nil
	}
}

// DeletePostgresqlDbUserInvoker 删除数据库用户
func (c *RdsClient) DeletePostgresqlDbUserInvoker(request *model.DeletePostgresqlDbUserRequest) *DeletePostgresqlDbUserInvoker {
	requestDef := GenReqDefForDeletePostgresqlDbUser()
	return &DeletePostgresqlDbUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeletePostgresqlExtension 删除插件
//
// 在指定数据库上删除插件。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeletePostgresqlExtension(request *model.DeletePostgresqlExtensionRequest) (*model.DeletePostgresqlExtensionResponse, error) {
	requestDef := GenReqDefForDeletePostgresqlExtension()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeletePostgresqlExtensionResponse), nil
	}
}

// DeletePostgresqlExtensionInvoker 删除插件
func (c *RdsClient) DeletePostgresqlExtensionInvoker(request *model.DeletePostgresqlExtensionRequest) *DeletePostgresqlExtensionInvoker {
	requestDef := GenReqDefForDeletePostgresqlExtension()
	return &DeletePostgresqlExtensionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPostgresqlDatabaseSchemas 查询数据库SCHEMA列表
//
// 查询指定实例的数据库SCHEMA列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlDatabaseSchemas(request *model.ListPostgresqlDatabaseSchemasRequest) (*model.ListPostgresqlDatabaseSchemasResponse, error) {
	requestDef := GenReqDefForListPostgresqlDatabaseSchemas()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDatabaseSchemasResponse), nil
	}
}

// ListPostgresqlDatabaseSchemasInvoker 查询数据库SCHEMA列表
func (c *RdsClient) ListPostgresqlDatabaseSchemasInvoker(request *model.ListPostgresqlDatabaseSchemasRequest) *ListPostgresqlDatabaseSchemasInvoker {
	requestDef := GenReqDefForListPostgresqlDatabaseSchemas()
	return &ListPostgresqlDatabaseSchemasInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPostgresqlDatabases 查询数据库列表
//
// 查询指定实例中的数据库列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlDatabases(request *model.ListPostgresqlDatabasesRequest) (*model.ListPostgresqlDatabasesResponse, error) {
	requestDef := GenReqDefForListPostgresqlDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDatabasesResponse), nil
	}
}

// ListPostgresqlDatabasesInvoker 查询数据库列表
func (c *RdsClient) ListPostgresqlDatabasesInvoker(request *model.ListPostgresqlDatabasesRequest) *ListPostgresqlDatabasesInvoker {
	requestDef := GenReqDefForListPostgresqlDatabases()
	return &ListPostgresqlDatabasesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPostgresqlDbUserPaginated 查询数据库用户列表
//
// 在指定实例中查询数据库用户列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlDbUserPaginated(request *model.ListPostgresqlDbUserPaginatedRequest) (*model.ListPostgresqlDbUserPaginatedResponse, error) {
	requestDef := GenReqDefForListPostgresqlDbUserPaginated()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlDbUserPaginatedResponse), nil
	}
}

// ListPostgresqlDbUserPaginatedInvoker 查询数据库用户列表
func (c *RdsClient) ListPostgresqlDbUserPaginatedInvoker(request *model.ListPostgresqlDbUserPaginatedRequest) *ListPostgresqlDbUserPaginatedInvoker {
	requestDef := GenReqDefForListPostgresqlDbUserPaginated()
	return &ListPostgresqlDbUserPaginatedInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListPostgresqlExtension 查询插件
//
// 获取指定数据库的插件信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListPostgresqlExtension(request *model.ListPostgresqlExtensionRequest) (*model.ListPostgresqlExtensionResponse, error) {
	requestDef := GenReqDefForListPostgresqlExtension()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListPostgresqlExtensionResponse), nil
	}
}

// ListPostgresqlExtensionInvoker 查询插件
func (c *RdsClient) ListPostgresqlExtensionInvoker(request *model.ListPostgresqlExtensionRequest) *ListPostgresqlExtensionInvoker {
	requestDef := GenReqDefForListPostgresqlExtension()
	return &ListPostgresqlExtensionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RevokePostgresqlDbPrivilege 解除数据库帐号权限
//
// 解除数据库帐号权限
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) RevokePostgresqlDbPrivilege(request *model.RevokePostgresqlDbPrivilegeRequest) (*model.RevokePostgresqlDbPrivilegeResponse, error) {
	requestDef := GenReqDefForRevokePostgresqlDbPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokePostgresqlDbPrivilegeResponse), nil
	}
}

// RevokePostgresqlDbPrivilegeInvoker 解除数据库帐号权限
func (c *RdsClient) RevokePostgresqlDbPrivilegeInvoker(request *model.RevokePostgresqlDbPrivilegeRequest) *RevokePostgresqlDbPrivilegeInvoker {
	requestDef := GenReqDefForRevokePostgresqlDbPrivilege()
	return &RevokePostgresqlDbPrivilegeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SearchQueryScaleComputeFlavors 查询数据库代理可变更的规格
//
// 查询数据库代理可变更的规格信息。
//
// - 调用接口前，您需要了解API 认证鉴权。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SearchQueryScaleComputeFlavors(request *model.SearchQueryScaleComputeFlavorsRequest) (*model.SearchQueryScaleComputeFlavorsResponse, error) {
	requestDef := GenReqDefForSearchQueryScaleComputeFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SearchQueryScaleComputeFlavorsResponse), nil
	}
}

// SearchQueryScaleComputeFlavorsInvoker 查询数据库代理可变更的规格
func (c *RdsClient) SearchQueryScaleComputeFlavorsInvoker(request *model.SearchQueryScaleComputeFlavorsRequest) *SearchQueryScaleComputeFlavorsInvoker {
	requestDef := GenReqDefForSearchQueryScaleComputeFlavors()
	return &SearchQueryScaleComputeFlavorsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SearchQueryScaleFlavors 查询数据库代理可变更的规格
//
// 查询数据库代理可变更的规格信息。
//
// - 调用接口前，您需要了解API 认证鉴权。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SearchQueryScaleFlavors(request *model.SearchQueryScaleFlavorsRequest) (*model.SearchQueryScaleFlavorsResponse, error) {
	requestDef := GenReqDefForSearchQueryScaleFlavors()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SearchQueryScaleFlavorsResponse), nil
	}
}

// SearchQueryScaleFlavorsInvoker 查询数据库代理可变更的规格
func (c *RdsClient) SearchQueryScaleFlavorsInvoker(request *model.SearchQueryScaleFlavorsRequest) *SearchQueryScaleFlavorsInvoker {
	requestDef := GenReqDefForSearchQueryScaleFlavors()
	return &SearchQueryScaleFlavorsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetDatabaseUserPrivilege 设置数据库用户权限
//
// 设置数据库用户权限：只读或可读写。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetDatabaseUserPrivilege(request *model.SetDatabaseUserPrivilegeRequest) (*model.SetDatabaseUserPrivilegeResponse, error) {
	requestDef := GenReqDefForSetDatabaseUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetDatabaseUserPrivilegeResponse), nil
	}
}

// SetDatabaseUserPrivilegeInvoker 设置数据库用户权限
func (c *RdsClient) SetDatabaseUserPrivilegeInvoker(request *model.SetDatabaseUserPrivilegeRequest) *SetDatabaseUserPrivilegeInvoker {
	requestDef := GenReqDefForSetDatabaseUserPrivilege()
	return &SetDatabaseUserPrivilegeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// SetPostgresqlDbUserPwd 重置数据库帐号密码
//
// 重置指定数据库帐号的密码。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) SetPostgresqlDbUserPwd(request *model.SetPostgresqlDbUserPwdRequest) (*model.SetPostgresqlDbUserPwdResponse, error) {
	requestDef := GenReqDefForSetPostgresqlDbUserPwd()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.SetPostgresqlDbUserPwdResponse), nil
	}
}

// SetPostgresqlDbUserPwdInvoker 重置数据库帐号密码
func (c *RdsClient) SetPostgresqlDbUserPwdInvoker(request *model.SetPostgresqlDbUserPwdRequest) *SetPostgresqlDbUserPwdInvoker {
	requestDef := GenReqDefForSetPostgresqlDbUserPwd()
	return &SetPostgresqlDbUserPwdInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowInformationAboutDatabaseProxy 查询数据库代理信息
//
// 查询指定实例的数据库代理详细信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowInformationAboutDatabaseProxy(request *model.ShowInformationAboutDatabaseProxyRequest) (*model.ShowInformationAboutDatabaseProxyResponse, error) {
	requestDef := GenReqDefForShowInformationAboutDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowInformationAboutDatabaseProxyResponse), nil
	}
}

// ShowInformationAboutDatabaseProxyInvoker 查询数据库代理信息
func (c *RdsClient) ShowInformationAboutDatabaseProxyInvoker(request *model.ShowInformationAboutDatabaseProxyRequest) *ShowInformationAboutDatabaseProxyInvoker {
	requestDef := GenReqDefForShowInformationAboutDatabaseProxy()
	return &ShowInformationAboutDatabaseProxyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowPostgresqlParamValue 获取实例指定参数的值
//
// 获取实例指定参数的值。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ShowPostgresqlParamValue(request *model.ShowPostgresqlParamValueRequest) (*model.ShowPostgresqlParamValueResponse, error) {
	requestDef := GenReqDefForShowPostgresqlParamValue()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowPostgresqlParamValueResponse), nil
	}
}

// ShowPostgresqlParamValueInvoker 获取实例指定参数的值
func (c *RdsClient) ShowPostgresqlParamValueInvoker(request *model.ShowPostgresqlParamValueRequest) *ShowPostgresqlParamValueInvoker {
	requestDef := GenReqDefForShowPostgresqlParamValue()
	return &ShowPostgresqlParamValueInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartDatabaseProxy 开启数据库代理
//
// 为指定实例开启数据库代理。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StartDatabaseProxy(request *model.StartDatabaseProxyRequest) (*model.StartDatabaseProxyResponse, error) {
	requestDef := GenReqDefForStartDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartDatabaseProxyResponse), nil
	}
}

// StartDatabaseProxyInvoker 开启数据库代理
func (c *RdsClient) StartDatabaseProxyInvoker(request *model.StartDatabaseProxyRequest) *StartDatabaseProxyInvoker {
	requestDef := GenReqDefForStartDatabaseProxy()
	return &StartDatabaseProxyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopDatabaseProxy 关闭数据库代理
//
// 为指定实例关闭数据库代理。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) StopDatabaseProxy(request *model.StopDatabaseProxyRequest) (*model.StopDatabaseProxyResponse, error) {
	requestDef := GenReqDefForStopDatabaseProxy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopDatabaseProxyResponse), nil
	}
}

// StopDatabaseProxyInvoker 关闭数据库代理
func (c *RdsClient) StopDatabaseProxyInvoker(request *model.StopDatabaseProxyRequest) *StopDatabaseProxyInvoker {
	requestDef := GenReqDefForStopDatabaseProxy()
	return &StopDatabaseProxyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDbUserPrivilege
//
// 数据库帐号授权。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateDbUserPrivilege(request *model.UpdateDbUserPrivilegeRequest) (*model.UpdateDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForUpdateDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDbUserPrivilegeResponse), nil
	}
}

// UpdateDbUserPrivilegeInvoker
func (c *RdsClient) UpdateDbUserPrivilegeInvoker(request *model.UpdateDbUserPrivilegeRequest) *UpdateDbUserPrivilegeInvoker {
	requestDef := GenReqDefForUpdateDbUserPrivilege()
	return &UpdateDbUserPrivilegeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePostgresqlDatabase 修改指定实例的数据库备注
//
// 修改指定实例中的数据库备注。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdatePostgresqlDatabase(request *model.UpdatePostgresqlDatabaseRequest) (*model.UpdatePostgresqlDatabaseResponse, error) {
	requestDef := GenReqDefForUpdatePostgresqlDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePostgresqlDatabaseResponse), nil
	}
}

// UpdatePostgresqlDatabaseInvoker 修改指定实例的数据库备注
func (c *RdsClient) UpdatePostgresqlDatabaseInvoker(request *model.UpdatePostgresqlDatabaseRequest) *UpdatePostgresqlDatabaseInvoker {
	requestDef := GenReqDefForUpdatePostgresqlDatabase()
	return &UpdatePostgresqlDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePostgresqlDbUserComment 修改数据库用户名备注
//
// 修改数据库用户名备注
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdatePostgresqlDbUserComment(request *model.UpdatePostgresqlDbUserCommentRequest) (*model.UpdatePostgresqlDbUserCommentResponse, error) {
	requestDef := GenReqDefForUpdatePostgresqlDbUserComment()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePostgresqlDbUserCommentResponse), nil
	}
}

// UpdatePostgresqlDbUserCommentInvoker 修改数据库用户名备注
func (c *RdsClient) UpdatePostgresqlDbUserCommentInvoker(request *model.UpdatePostgresqlDbUserCommentRequest) *UpdatePostgresqlDbUserCommentInvoker {
	requestDef := GenReqDefForUpdatePostgresqlDbUserComment()
	return &UpdatePostgresqlDbUserCommentInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdatePostgresqlParameterValue 修改实例指定参数的值
//
// 修改实例指定参数的值。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdatePostgresqlParameterValue(request *model.UpdatePostgresqlParameterValueRequest) (*model.UpdatePostgresqlParameterValueResponse, error) {
	requestDef := GenReqDefForUpdatePostgresqlParameterValue()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdatePostgresqlParameterValueResponse), nil
	}
}

// UpdatePostgresqlParameterValueInvoker 修改实例指定参数的值
func (c *RdsClient) UpdatePostgresqlParameterValueInvoker(request *model.UpdatePostgresqlParameterValueRequest) *UpdatePostgresqlParameterValueInvoker {
	requestDef := GenReqDefForUpdatePostgresqlParameterValue()
	return &UpdatePostgresqlParameterValueInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateReadWeight 修改读写分离权重
//
// 修改指定实例的读写分离权重。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) UpdateReadWeight(request *model.UpdateReadWeightRequest) (*model.UpdateReadWeightResponse, error) {
	requestDef := GenReqDefForUpdateReadWeight()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateReadWeightResponse), nil
	}
}

// UpdateReadWeightInvoker 修改读写分离权重
func (c *RdsClient) UpdateReadWeightInvoker(request *model.UpdateReadWeightRequest) *UpdateReadWeightInvoker {
	requestDef := GenReqDefForUpdateReadWeight()
	return &UpdateReadWeightInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// AllowSqlserverDbUserPrivilege 授权数据库帐号
//
// 授权数据库帐号。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) AllowSqlserverDbUserPrivilege(request *model.AllowSqlserverDbUserPrivilegeRequest) (*model.AllowSqlserverDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForAllowSqlserverDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.AllowSqlserverDbUserPrivilegeResponse), nil
	}
}

// AllowSqlserverDbUserPrivilegeInvoker 授权数据库帐号
func (c *RdsClient) AllowSqlserverDbUserPrivilegeInvoker(request *model.AllowSqlserverDbUserPrivilegeRequest) *AllowSqlserverDbUserPrivilegeInvoker {
	requestDef := GenReqDefForAllowSqlserverDbUserPrivilege()
	return &AllowSqlserverDbUserPrivilegeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// BatchAddMsdtcs 添加MSDTC
//
// 添加MSDTC相关主机host地址
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) BatchAddMsdtcs(request *model.BatchAddMsdtcsRequest) (*model.BatchAddMsdtcsResponse, error) {
	requestDef := GenReqDefForBatchAddMsdtcs()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchAddMsdtcsResponse), nil
	}
}

// BatchAddMsdtcsInvoker 添加MSDTC
func (c *RdsClient) BatchAddMsdtcsInvoker(request *model.BatchAddMsdtcsRequest) *BatchAddMsdtcsInvoker {
	requestDef := GenReqDefForBatchAddMsdtcs()
	return &BatchAddMsdtcsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSqlserverDatabase 创建数据库
//
// 创建数据库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateSqlserverDatabase(request *model.CreateSqlserverDatabaseRequest) (*model.CreateSqlserverDatabaseResponse, error) {
	requestDef := GenReqDefForCreateSqlserverDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSqlserverDatabaseResponse), nil
	}
}

// CreateSqlserverDatabaseInvoker 创建数据库
func (c *RdsClient) CreateSqlserverDatabaseInvoker(request *model.CreateSqlserverDatabaseRequest) *CreateSqlserverDatabaseInvoker {
	requestDef := GenReqDefForCreateSqlserverDatabase()
	return &CreateSqlserverDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSqlserverDbUser 创建数据库用户
//
// 创建数据库用户。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) CreateSqlserverDbUser(request *model.CreateSqlserverDbUserRequest) (*model.CreateSqlserverDbUserResponse, error) {
	requestDef := GenReqDefForCreateSqlserverDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSqlserverDbUserResponse), nil
	}
}

// CreateSqlserverDbUserInvoker 创建数据库用户
func (c *RdsClient) CreateSqlserverDbUserInvoker(request *model.CreateSqlserverDbUserRequest) *CreateSqlserverDbUserInvoker {
	requestDef := GenReqDefForCreateSqlserverDbUser()
	return &CreateSqlserverDbUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSqlserverDatabase 删除数据库
//
// 删除数据库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteSqlserverDatabase(request *model.DeleteSqlserverDatabaseRequest) (*model.DeleteSqlserverDatabaseResponse, error) {
	requestDef := GenReqDefForDeleteSqlserverDatabase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSqlserverDatabaseResponse), nil
	}
}

// DeleteSqlserverDatabaseInvoker 删除数据库
func (c *RdsClient) DeleteSqlserverDatabaseInvoker(request *model.DeleteSqlserverDatabaseRequest) *DeleteSqlserverDatabaseInvoker {
	requestDef := GenReqDefForDeleteSqlserverDatabase()
	return &DeleteSqlserverDatabaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSqlserverDatabaseEx 删除数据库
//
// 删除数据库。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteSqlserverDatabaseEx(request *model.DeleteSqlserverDatabaseExRequest) (*model.DeleteSqlserverDatabaseExResponse, error) {
	requestDef := GenReqDefForDeleteSqlserverDatabaseEx()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSqlserverDatabaseExResponse), nil
	}
}

// DeleteSqlserverDatabaseExInvoker 删除数据库
func (c *RdsClient) DeleteSqlserverDatabaseExInvoker(request *model.DeleteSqlserverDatabaseExRequest) *DeleteSqlserverDatabaseExInvoker {
	requestDef := GenReqDefForDeleteSqlserverDatabaseEx()
	return &DeleteSqlserverDatabaseExInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteSqlserverDbUser 删除数据库用户
//
// 删除数据库用户。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) DeleteSqlserverDbUser(request *model.DeleteSqlserverDbUserRequest) (*model.DeleteSqlserverDbUserResponse, error) {
	requestDef := GenReqDefForDeleteSqlserverDbUser()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSqlserverDbUserResponse), nil
	}
}

// DeleteSqlserverDbUserInvoker 删除数据库用户
func (c *RdsClient) DeleteSqlserverDbUserInvoker(request *model.DeleteSqlserverDbUserRequest) *DeleteSqlserverDbUserInvoker {
	requestDef := GenReqDefForDeleteSqlserverDbUser()
	return &DeleteSqlserverDbUserInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListAuthorizedSqlserverDbUsers 查询指定数据库的已授权用户
//
// 查询指定数据库的已授权用户。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListAuthorizedSqlserverDbUsers(request *model.ListAuthorizedSqlserverDbUsersRequest) (*model.ListAuthorizedSqlserverDbUsersResponse, error) {
	requestDef := GenReqDefForListAuthorizedSqlserverDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListAuthorizedSqlserverDbUsersResponse), nil
	}
}

// ListAuthorizedSqlserverDbUsersInvoker 查询指定数据库的已授权用户
func (c *RdsClient) ListAuthorizedSqlserverDbUsersInvoker(request *model.ListAuthorizedSqlserverDbUsersRequest) *ListAuthorizedSqlserverDbUsersInvoker {
	requestDef := GenReqDefForListAuthorizedSqlserverDbUsers()
	return &ListAuthorizedSqlserverDbUsersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListMsdtcHosts 查询MSDTC的hosts信息
//
// 查询MSDTC的hosts信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListMsdtcHosts(request *model.ListMsdtcHostsRequest) (*model.ListMsdtcHostsResponse, error) {
	requestDef := GenReqDefForListMsdtcHosts()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListMsdtcHostsResponse), nil
	}
}

// ListMsdtcHostsInvoker 查询MSDTC的hosts信息
func (c *RdsClient) ListMsdtcHostsInvoker(request *model.ListMsdtcHostsRequest) *ListMsdtcHostsInvoker {
	requestDef := GenReqDefForListMsdtcHosts()
	return &ListMsdtcHostsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSqlserverDatabases 查询数据库列表
//
// 查询数据库列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSqlserverDatabases(request *model.ListSqlserverDatabasesRequest) (*model.ListSqlserverDatabasesResponse, error) {
	requestDef := GenReqDefForListSqlserverDatabases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSqlserverDatabasesResponse), nil
	}
}

// ListSqlserverDatabasesInvoker 查询数据库列表
func (c *RdsClient) ListSqlserverDatabasesInvoker(request *model.ListSqlserverDatabasesRequest) *ListSqlserverDatabasesInvoker {
	requestDef := GenReqDefForListSqlserverDatabases()
	return &ListSqlserverDatabasesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSqlserverDbUsers 查询数据库用户列表
//
// 查询数据库用户列表。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ListSqlserverDbUsers(request *model.ListSqlserverDbUsersRequest) (*model.ListSqlserverDbUsersResponse, error) {
	requestDef := GenReqDefForListSqlserverDbUsers()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSqlserverDbUsersResponse), nil
	}
}

// ListSqlserverDbUsersInvoker 查询数据库用户列表
func (c *RdsClient) ListSqlserverDbUsersInvoker(request *model.ListSqlserverDbUsersRequest) *ListSqlserverDbUsersInvoker {
	requestDef := GenReqDefForListSqlserverDbUsers()
	return &ListSqlserverDbUsersInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ModifyCollation 修改实例字符集
//
// 修改实例字符集。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) ModifyCollation(request *model.ModifyCollationRequest) (*model.ModifyCollationResponse, error) {
	requestDef := GenReqDefForModifyCollation()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ModifyCollationResponse), nil
	}
}

// ModifyCollationInvoker 修改实例字符集
func (c *RdsClient) ModifyCollationInvoker(request *model.ModifyCollationRequest) *ModifyCollationInvoker {
	requestDef := GenReqDefForModifyCollation()
	return &ModifyCollationInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RevokeSqlserverDbUserPrivilege 解除数据库帐号权限
//
// 解除数据库帐号权限。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *RdsClient) RevokeSqlserverDbUserPrivilege(request *model.RevokeSqlserverDbUserPrivilegeRequest) (*model.RevokeSqlserverDbUserPrivilegeResponse, error) {
	requestDef := GenReqDefForRevokeSqlserverDbUserPrivilege()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RevokeSqlserverDbUserPrivilegeResponse), nil
	}
}

// RevokeSqlserverDbUserPrivilegeInvoker 解除数据库帐号权限
func (c *RdsClient) RevokeSqlserverDbUserPrivilegeInvoker(request *model.RevokeSqlserverDbUserPrivilegeRequest) *RevokeSqlserverDbUserPrivilegeInvoker {
	requestDef := GenReqDefForRevokeSqlserverDbUserPrivilege()
	return &RevokeSqlserverDbUserPrivilegeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
