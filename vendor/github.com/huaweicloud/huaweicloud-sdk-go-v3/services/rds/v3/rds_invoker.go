package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
)

type AddPostgresqlHbaConfInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddPostgresqlHbaConfInvoker) Invoke() (*model.AddPostgresqlHbaConfResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddPostgresqlHbaConfResponse), nil
	}
}

type ApplyConfigurationAsyncInvoker struct {
	*invoker.BaseInvoker
}

func (i *ApplyConfigurationAsyncInvoker) Invoke() (*model.ApplyConfigurationAsyncResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ApplyConfigurationAsyncResponse), nil
	}
}

type AttachEipInvoker struct {
	*invoker.BaseInvoker
}

func (i *AttachEipInvoker) Invoke() (*model.AttachEipResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AttachEipResponse), nil
	}
}

type BatchDeleteManualBackupInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteManualBackupInvoker) Invoke() (*model.BatchDeleteManualBackupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteManualBackupResponse), nil
	}
}

type BatchRestoreDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchRestoreDatabaseInvoker) Invoke() (*model.BatchRestoreDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchRestoreDatabaseResponse), nil
	}
}

type BatchRestorePostgreSqlTablesInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchRestorePostgreSqlTablesInvoker) Invoke() (*model.BatchRestorePostgreSqlTablesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchRestorePostgreSqlTablesResponse), nil
	}
}

type BatchTagAddActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchTagAddActionInvoker) Invoke() (*model.BatchTagAddActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchTagAddActionResponse), nil
	}
}

type BatchTagDelActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchTagDelActionInvoker) Invoke() (*model.BatchTagDelActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchTagDelActionResponse), nil
	}
}

type ChangeFailoverModeInvoker struct {
	*invoker.BaseInvoker
}

func (i *ChangeFailoverModeInvoker) Invoke() (*model.ChangeFailoverModeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ChangeFailoverModeResponse), nil
	}
}

type ChangeFailoverStrategyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ChangeFailoverStrategyInvoker) Invoke() (*model.ChangeFailoverStrategyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ChangeFailoverStrategyResponse), nil
	}
}

type ChangeOpsWindowInvoker struct {
	*invoker.BaseInvoker
}

func (i *ChangeOpsWindowInvoker) Invoke() (*model.ChangeOpsWindowResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ChangeOpsWindowResponse), nil
	}
}

type CopyConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *CopyConfigurationInvoker) Invoke() (*model.CopyConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CopyConfigurationResponse), nil
	}
}

type CreateConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateConfigurationInvoker) Invoke() (*model.CreateConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateConfigurationResponse), nil
	}
}

type CreateDnsNameInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDnsNameInvoker) Invoke() (*model.CreateDnsNameResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDnsNameResponse), nil
	}
}

type CreateInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateInstanceInvoker) Invoke() (*model.CreateInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateInstanceResponse), nil
	}
}

type CreateManualBackupInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateManualBackupInvoker) Invoke() (*model.CreateManualBackupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateManualBackupResponse), nil
	}
}

type CreateRestoreInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRestoreInstanceInvoker) Invoke() (*model.CreateRestoreInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRestoreInstanceResponse), nil
	}
}

type CreateXelLogDownloadInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateXelLogDownloadInvoker) Invoke() (*model.CreateXelLogDownloadResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateXelLogDownloadResponse), nil
	}
}

type DeleteConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteConfigurationInvoker) Invoke() (*model.DeleteConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteConfigurationResponse), nil
	}
}

type DeleteInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteInstanceInvoker) Invoke() (*model.DeleteInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteInstanceResponse), nil
	}
}

type DeleteJobInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteJobInvoker) Invoke() (*model.DeleteJobResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteJobResponse), nil
	}
}

type DeleteManualBackupInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteManualBackupInvoker) Invoke() (*model.DeleteManualBackupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteManualBackupResponse), nil
	}
}

type DeletePostgresqlHbaConfInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeletePostgresqlHbaConfInvoker) Invoke() (*model.DeletePostgresqlHbaConfResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeletePostgresqlHbaConfResponse), nil
	}
}

type DownloadSlowlogInvoker struct {
	*invoker.BaseInvoker
}

func (i *DownloadSlowlogInvoker) Invoke() (*model.DownloadSlowlogResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DownloadSlowlogResponse), nil
	}
}

type EnableConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *EnableConfigurationInvoker) Invoke() (*model.EnableConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.EnableConfigurationResponse), nil
	}
}

type ListAuditlogsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAuditlogsInvoker) Invoke() (*model.ListAuditlogsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAuditlogsResponse), nil
	}
}

type ListBackupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListBackupsInvoker) Invoke() (*model.ListBackupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListBackupsResponse), nil
	}
}

type ListCollationsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListCollationsInvoker) Invoke() (*model.ListCollationsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListCollationsResponse), nil
	}
}

type ListConfigurationsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListConfigurationsInvoker) Invoke() (*model.ListConfigurationsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListConfigurationsResponse), nil
	}
}

type ListDatastoresInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDatastoresInvoker) Invoke() (*model.ListDatastoresResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDatastoresResponse), nil
	}
}

type ListDrRelationsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDrRelationsInvoker) Invoke() (*model.ListDrRelationsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDrRelationsResponse), nil
	}
}

type ListEngineFlavorsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListEngineFlavorsInvoker) Invoke() (*model.ListEngineFlavorsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListEngineFlavorsResponse), nil
	}
}

type ListErrorLogsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListErrorLogsInvoker) Invoke() (*model.ListErrorLogsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListErrorLogsResponse), nil
	}
}

type ListErrorLogsNewInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListErrorLogsNewInvoker) Invoke() (*model.ListErrorLogsNewResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListErrorLogsNewResponse), nil
	}
}

type ListErrorlogForLtsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListErrorlogForLtsInvoker) Invoke() (*model.ListErrorlogForLtsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListErrorlogForLtsResponse), nil
	}
}

type ListFlavorsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListFlavorsInvoker) Invoke() (*model.ListFlavorsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListFlavorsResponse), nil
	}
}

type ListHistoryDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListHistoryDatabaseInvoker) Invoke() (*model.ListHistoryDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListHistoryDatabaseResponse), nil
	}
}

type ListInspectionHistoriesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInspectionHistoriesInvoker) Invoke() (*model.ListInspectionHistoriesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInspectionHistoriesResponse), nil
	}
}

type ListInstanceDiagnosisInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstanceDiagnosisInvoker) Invoke() (*model.ListInstanceDiagnosisResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstanceDiagnosisResponse), nil
	}
}

type ListInstanceParamHistoriesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstanceParamHistoriesInvoker) Invoke() (*model.ListInstanceParamHistoriesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstanceParamHistoriesResponse), nil
	}
}

type ListInstanceTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstanceTagsInvoker) Invoke() (*model.ListInstanceTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstanceTagsResponse), nil
	}
}

type ListInstancesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstancesInvoker) Invoke() (*model.ListInstancesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstancesResponse), nil
	}
}

type ListInstancesInfoDiagnosisInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstancesInfoDiagnosisInvoker) Invoke() (*model.ListInstancesInfoDiagnosisResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstancesInfoDiagnosisResponse), nil
	}
}

type ListInstancesSupportFastRestoreInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstancesSupportFastRestoreInvoker) Invoke() (*model.ListInstancesSupportFastRestoreResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstancesSupportFastRestoreResponse), nil
	}
}

type ListJobInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListJobInfoInvoker) Invoke() (*model.ListJobInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListJobInfoResponse), nil
	}
}

type ListJobInfoDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListJobInfoDetailInvoker) Invoke() (*model.ListJobInfoDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListJobInfoDetailResponse), nil
	}
}

type ListOffSiteBackupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListOffSiteBackupsInvoker) Invoke() (*model.ListOffSiteBackupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListOffSiteBackupsResponse), nil
	}
}

type ListOffSiteInstancesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListOffSiteInstancesInvoker) Invoke() (*model.ListOffSiteInstancesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListOffSiteInstancesResponse), nil
	}
}

type ListOffSiteRestoreTimesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListOffSiteRestoreTimesInvoker) Invoke() (*model.ListOffSiteRestoreTimesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListOffSiteRestoreTimesResponse), nil
	}
}

type ListPostgresqlHbaInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPostgresqlHbaInfoInvoker) Invoke() (*model.ListPostgresqlHbaInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPostgresqlHbaInfoResponse), nil
	}
}

type ListPostgresqlHbaInfoHistoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPostgresqlHbaInfoHistoryInvoker) Invoke() (*model.ListPostgresqlHbaInfoHistoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPostgresqlHbaInfoHistoryResponse), nil
	}
}

type ListPostgresqlListHistoryTablesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPostgresqlListHistoryTablesInvoker) Invoke() (*model.ListPostgresqlListHistoryTablesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPostgresqlListHistoryTablesResponse), nil
	}
}

type ListPredefinedTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPredefinedTagInvoker) Invoke() (*model.ListPredefinedTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPredefinedTagResponse), nil
	}
}

type ListProjectTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListProjectTagsInvoker) Invoke() (*model.ListProjectTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListProjectTagsResponse), nil
	}
}

type ListRecycleInstancesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRecycleInstancesInvoker) Invoke() (*model.ListRecycleInstancesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRecycleInstancesResponse), nil
	}
}

type ListRestoreTimesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRestoreTimesInvoker) Invoke() (*model.ListRestoreTimesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRestoreTimesResponse), nil
	}
}

type ListSimplifiedInstancesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSimplifiedInstancesInvoker) Invoke() (*model.ListSimplifiedInstancesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSimplifiedInstancesResponse), nil
	}
}

type ListSlowLogFileInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSlowLogFileInvoker) Invoke() (*model.ListSlowLogFileResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSlowLogFileResponse), nil
	}
}

type ListSlowLogStatisticsForLtsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSlowLogStatisticsForLtsInvoker) Invoke() (*model.ListSlowLogStatisticsForLtsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSlowLogStatisticsForLtsResponse), nil
	}
}

type ListSlowLogsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSlowLogsInvoker) Invoke() (*model.ListSlowLogsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSlowLogsResponse), nil
	}
}

type ListSlowLogsNewInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSlowLogsNewInvoker) Invoke() (*model.ListSlowLogsNewResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSlowLogsNewResponse), nil
	}
}

type ListSlowlogForLtsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSlowlogForLtsInvoker) Invoke() (*model.ListSlowlogForLtsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSlowlogForLtsResponse), nil
	}
}

type ListSlowlogStatisticsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSlowlogStatisticsInvoker) Invoke() (*model.ListSlowlogStatisticsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSlowlogStatisticsResponse), nil
	}
}

type ListSslCertDownloadLinkInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSslCertDownloadLinkInvoker) Invoke() (*model.ListSslCertDownloadLinkResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSslCertDownloadLinkResponse), nil
	}
}

type ListStorageTypesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListStorageTypesInvoker) Invoke() (*model.ListStorageTypesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListStorageTypesResponse), nil
	}
}

type ListUpgradeHistoriesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListUpgradeHistoriesInvoker) Invoke() (*model.ListUpgradeHistoriesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListUpgradeHistoriesResponse), nil
	}
}

type ListXellogFilesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListXellogFilesInvoker) Invoke() (*model.ListXellogFilesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListXellogFilesResponse), nil
	}
}

type MigrateFollowerInvoker struct {
	*invoker.BaseInvoker
}

func (i *MigrateFollowerInvoker) Invoke() (*model.MigrateFollowerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.MigrateFollowerResponse), nil
	}
}

type ModifyPostgresqlHbaConfInvoker struct {
	*invoker.BaseInvoker
}

func (i *ModifyPostgresqlHbaConfInvoker) Invoke() (*model.ModifyPostgresqlHbaConfResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ModifyPostgresqlHbaConfResponse), nil
	}
}

type RestoreExistInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *RestoreExistInstanceInvoker) Invoke() (*model.RestoreExistInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RestoreExistInstanceResponse), nil
	}
}

type RestoreTablesInvoker struct {
	*invoker.BaseInvoker
}

func (i *RestoreTablesInvoker) Invoke() (*model.RestoreTablesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RestoreTablesResponse), nil
	}
}

type RestoreTablesNewInvoker struct {
	*invoker.BaseInvoker
}

func (i *RestoreTablesNewInvoker) Invoke() (*model.RestoreTablesNewResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RestoreTablesNewResponse), nil
	}
}

type RestoreToExistingInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *RestoreToExistingInstanceInvoker) Invoke() (*model.RestoreToExistingInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RestoreToExistingInstanceResponse), nil
	}
}

type SetAuditlogPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetAuditlogPolicyInvoker) Invoke() (*model.SetAuditlogPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetAuditlogPolicyResponse), nil
	}
}

type SetAutoEnlargePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetAutoEnlargePolicyInvoker) Invoke() (*model.SetAutoEnlargePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetAutoEnlargePolicyResponse), nil
	}
}

type SetBackupPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetBackupPolicyInvoker) Invoke() (*model.SetBackupPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetBackupPolicyResponse), nil
	}
}

type SetBinlogClearPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetBinlogClearPolicyInvoker) Invoke() (*model.SetBinlogClearPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetBinlogClearPolicyResponse), nil
	}
}

type SetOffSiteBackupPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetOffSiteBackupPolicyInvoker) Invoke() (*model.SetOffSiteBackupPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetOffSiteBackupPolicyResponse), nil
	}
}

type SetSecondLevelMonitorInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetSecondLevelMonitorInvoker) Invoke() (*model.SetSecondLevelMonitorResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetSecondLevelMonitorResponse), nil
	}
}

type SetSecurityGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetSecurityGroupInvoker) Invoke() (*model.SetSecurityGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetSecurityGroupResponse), nil
	}
}

type SetSensitiveSlowLogInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetSensitiveSlowLogInvoker) Invoke() (*model.SetSensitiveSlowLogResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetSensitiveSlowLogResponse), nil
	}
}

type ShowAuditlogDownloadLinkInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAuditlogDownloadLinkInvoker) Invoke() (*model.ShowAuditlogDownloadLinkResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAuditlogDownloadLinkResponse), nil
	}
}

type ShowAuditlogPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAuditlogPolicyInvoker) Invoke() (*model.ShowAuditlogPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAuditlogPolicyResponse), nil
	}
}

type ShowAutoEnlargePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutoEnlargePolicyInvoker) Invoke() (*model.ShowAutoEnlargePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutoEnlargePolicyResponse), nil
	}
}

type ShowAvailableVersionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAvailableVersionInvoker) Invoke() (*model.ShowAvailableVersionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAvailableVersionResponse), nil
	}
}

type ShowBackupDownloadLinkInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowBackupDownloadLinkInvoker) Invoke() (*model.ShowBackupDownloadLinkResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBackupDownloadLinkResponse), nil
	}
}

type ShowBackupPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowBackupPolicyInvoker) Invoke() (*model.ShowBackupPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBackupPolicyResponse), nil
	}
}

type ShowBinlogClearPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowBinlogClearPolicyInvoker) Invoke() (*model.ShowBinlogClearPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBinlogClearPolicyResponse), nil
	}
}

type ShowConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowConfigurationInvoker) Invoke() (*model.ShowConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowConfigurationResponse), nil
	}
}

type ShowDnsNameInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDnsNameInvoker) Invoke() (*model.ShowDnsNameResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDnsNameResponse), nil
	}
}

type ShowDomainNameInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDomainNameInvoker) Invoke() (*model.ShowDomainNameResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDomainNameResponse), nil
	}
}

type ShowDrReplicaStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDrReplicaStatusInvoker) Invoke() (*model.ShowDrReplicaStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDrReplicaStatusResponse), nil
	}
}

type ShowInstanceConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowInstanceConfigurationInvoker) Invoke() (*model.ShowInstanceConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowInstanceConfigurationResponse), nil
	}
}

type ShowOffSiteBackupPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowOffSiteBackupPolicyInvoker) Invoke() (*model.ShowOffSiteBackupPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowOffSiteBackupPolicyResponse), nil
	}
}

type ShowQuotasInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowQuotasInvoker) Invoke() (*model.ShowQuotasResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowQuotasResponse), nil
	}
}

type ShowRecyclePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRecyclePolicyInvoker) Invoke() (*model.ShowRecyclePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRecyclePolicyResponse), nil
	}
}

type ShowReplicationStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowReplicationStatusInvoker) Invoke() (*model.ShowReplicationStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowReplicationStatusResponse), nil
	}
}

type ShowSecondLevelMonitoringInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowSecondLevelMonitoringInvoker) Invoke() (*model.ShowSecondLevelMonitoringResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowSecondLevelMonitoringResponse), nil
	}
}

type ShowTdeStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTdeStatusInvoker) Invoke() (*model.ShowTdeStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTdeStatusResponse), nil
	}
}

type ShowUpgradeDbMajorVersionStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowUpgradeDbMajorVersionStatusInvoker) Invoke() (*model.ShowUpgradeDbMajorVersionStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUpgradeDbMajorVersionStatusResponse), nil
	}
}

type StartFailoverInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartFailoverInvoker) Invoke() (*model.StartFailoverResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartFailoverResponse), nil
	}
}

type StartInstanceEnlargeVolumeActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartInstanceEnlargeVolumeActionInvoker) Invoke() (*model.StartInstanceEnlargeVolumeActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartInstanceEnlargeVolumeActionResponse), nil
	}
}

type StartInstanceRestartActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartInstanceRestartActionInvoker) Invoke() (*model.StartInstanceRestartActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartInstanceRestartActionResponse), nil
	}
}

type StartInstanceSingleToHaActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartInstanceSingleToHaActionInvoker) Invoke() (*model.StartInstanceSingleToHaActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartInstanceSingleToHaActionResponse), nil
	}
}

type StartRecyclePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartRecyclePolicyInvoker) Invoke() (*model.StartRecyclePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartRecyclePolicyResponse), nil
	}
}

type StartResizeFlavorActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartResizeFlavorActionInvoker) Invoke() (*model.StartResizeFlavorActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartResizeFlavorActionResponse), nil
	}
}

type StartupInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartupInstanceInvoker) Invoke() (*model.StartupInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartupInstanceResponse), nil
	}
}

type StopInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopInstanceInvoker) Invoke() (*model.StopInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopInstanceResponse), nil
	}
}

type SwitchSslInvoker struct {
	*invoker.BaseInvoker
}

func (i *SwitchSslInvoker) Invoke() (*model.SwitchSslResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SwitchSslResponse), nil
	}
}

type UpdateConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateConfigurationInvoker) Invoke() (*model.UpdateConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateConfigurationResponse), nil
	}
}

type UpdateDataIpInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDataIpInvoker) Invoke() (*model.UpdateDataIpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDataIpResponse), nil
	}
}

type UpdateDnsNameInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDnsNameInvoker) Invoke() (*model.UpdateDnsNameResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDnsNameResponse), nil
	}
}

type UpdateInstanceConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateInstanceConfigurationInvoker) Invoke() (*model.UpdateInstanceConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateInstanceConfigurationResponse), nil
	}
}

type UpdateInstanceConfigurationAsyncInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateInstanceConfigurationAsyncInvoker) Invoke() (*model.UpdateInstanceConfigurationAsyncResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateInstanceConfigurationAsyncResponse), nil
	}
}

type UpdateInstanceNameInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateInstanceNameInvoker) Invoke() (*model.UpdateInstanceNameResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateInstanceNameResponse), nil
	}
}

type UpdatePortInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePortInvoker) Invoke() (*model.UpdatePortResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePortResponse), nil
	}
}

type UpdatePostgresqlInstanceAliasInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePostgresqlInstanceAliasInvoker) Invoke() (*model.UpdatePostgresqlInstanceAliasResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePostgresqlInstanceAliasResponse), nil
	}
}

type UpdateTdeStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTdeStatusInvoker) Invoke() (*model.UpdateTdeStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTdeStatusResponse), nil
	}
}

type UpgradeDbMajorVersionInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeDbMajorVersionInvoker) Invoke() (*model.UpgradeDbMajorVersionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeDbMajorVersionResponse), nil
	}
}

type UpgradeDbMajorVersionPreCheckInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeDbMajorVersionPreCheckInvoker) Invoke() (*model.UpgradeDbMajorVersionPreCheckResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeDbMajorVersionPreCheckResponse), nil
	}
}

type UpgradeDbVersionInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeDbVersionInvoker) Invoke() (*model.UpgradeDbVersionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeDbVersionResponse), nil
	}
}

type UpgradeDbVersionNewInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeDbVersionNewInvoker) Invoke() (*model.UpgradeDbVersionNewResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeDbVersionNewResponse), nil
	}
}

type ListApiVersionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListApiVersionInvoker) Invoke() (*model.ListApiVersionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListApiVersionResponse), nil
	}
}

type ListApiVersionNewInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListApiVersionNewInvoker) Invoke() (*model.ListApiVersionNewResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListApiVersionNewResponse), nil
	}
}

type ShowApiVersionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowApiVersionInvoker) Invoke() (*model.ShowApiVersionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowApiVersionResponse), nil
	}
}

type AllowDbUserPrivilegeInvoker struct {
	*invoker.BaseInvoker
}

func (i *AllowDbUserPrivilegeInvoker) Invoke() (*model.AllowDbUserPrivilegeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AllowDbUserPrivilegeResponse), nil
	}
}

type CreateDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDatabaseInvoker) Invoke() (*model.CreateDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDatabaseResponse), nil
	}
}

type CreateDbUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDbUserInvoker) Invoke() (*model.CreateDbUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDbUserResponse), nil
	}
}

type DeleteDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDatabaseInvoker) Invoke() (*model.DeleteDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDatabaseResponse), nil
	}
}

type DeleteDbUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDbUserInvoker) Invoke() (*model.DeleteDbUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDbUserResponse), nil
	}
}

type ListAuthorizedDatabasesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAuthorizedDatabasesInvoker) Invoke() (*model.ListAuthorizedDatabasesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAuthorizedDatabasesResponse), nil
	}
}

type ListAuthorizedDbUsersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAuthorizedDbUsersInvoker) Invoke() (*model.ListAuthorizedDbUsersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAuthorizedDbUsersResponse), nil
	}
}

type ListDatabasesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDatabasesInvoker) Invoke() (*model.ListDatabasesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDatabasesResponse), nil
	}
}

type ListDbUsersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDbUsersInvoker) Invoke() (*model.ListDbUsersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDbUsersResponse), nil
	}
}

type ResetPwdInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetPwdInvoker) Invoke() (*model.ResetPwdResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetPwdResponse), nil
	}
}

type RevokeInvoker struct {
	*invoker.BaseInvoker
}

func (i *RevokeInvoker) Invoke() (*model.RevokeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RevokeResponse), nil
	}
}

type SetDbUserPwdInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetDbUserPwdInvoker) Invoke() (*model.SetDbUserPwdResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetDbUserPwdResponse), nil
	}
}

type SetReadOnlySwitchInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetReadOnlySwitchInvoker) Invoke() (*model.SetReadOnlySwitchResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetReadOnlySwitchResponse), nil
	}
}

type UpdateDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDatabaseInvoker) Invoke() (*model.UpdateDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDatabaseResponse), nil
	}
}

type UpdateDbUserCommentInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDbUserCommentInvoker) Invoke() (*model.UpdateDbUserCommentResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDbUserCommentResponse), nil
	}
}

type AllowDbPrivilegeInvoker struct {
	*invoker.BaseInvoker
}

func (i *AllowDbPrivilegeInvoker) Invoke() (*model.AllowDbPrivilegeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AllowDbPrivilegeResponse), nil
	}
}

type ChangeProxyScaleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ChangeProxyScaleInvoker) Invoke() (*model.ChangeProxyScaleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ChangeProxyScaleResponse), nil
	}
}

type ChangeTheDelayThresholdInvoker struct {
	*invoker.BaseInvoker
}

func (i *ChangeTheDelayThresholdInvoker) Invoke() (*model.ChangeTheDelayThresholdResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ChangeTheDelayThresholdResponse), nil
	}
}

type CreatePostgresqlDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePostgresqlDatabaseInvoker) Invoke() (*model.CreatePostgresqlDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePostgresqlDatabaseResponse), nil
	}
}

type CreatePostgresqlDatabaseSchemaInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePostgresqlDatabaseSchemaInvoker) Invoke() (*model.CreatePostgresqlDatabaseSchemaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePostgresqlDatabaseSchemaResponse), nil
	}
}

type CreatePostgresqlDbUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePostgresqlDbUserInvoker) Invoke() (*model.CreatePostgresqlDbUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePostgresqlDbUserResponse), nil
	}
}

type CreatePostgresqlExtensionInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePostgresqlExtensionInvoker) Invoke() (*model.CreatePostgresqlExtensionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePostgresqlExtensionResponse), nil
	}
}

type DeletePostgresqlDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeletePostgresqlDatabaseInvoker) Invoke() (*model.DeletePostgresqlDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeletePostgresqlDatabaseResponse), nil
	}
}

type DeletePostgresqlDbUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeletePostgresqlDbUserInvoker) Invoke() (*model.DeletePostgresqlDbUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeletePostgresqlDbUserResponse), nil
	}
}

type DeletePostgresqlExtensionInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeletePostgresqlExtensionInvoker) Invoke() (*model.DeletePostgresqlExtensionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeletePostgresqlExtensionResponse), nil
	}
}

type ListPostgresqlDatabaseSchemasInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPostgresqlDatabaseSchemasInvoker) Invoke() (*model.ListPostgresqlDatabaseSchemasResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPostgresqlDatabaseSchemasResponse), nil
	}
}

type ListPostgresqlDatabasesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPostgresqlDatabasesInvoker) Invoke() (*model.ListPostgresqlDatabasesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPostgresqlDatabasesResponse), nil
	}
}

type ListPostgresqlDbUserPaginatedInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPostgresqlDbUserPaginatedInvoker) Invoke() (*model.ListPostgresqlDbUserPaginatedResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPostgresqlDbUserPaginatedResponse), nil
	}
}

type ListPostgresqlExtensionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPostgresqlExtensionInvoker) Invoke() (*model.ListPostgresqlExtensionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPostgresqlExtensionResponse), nil
	}
}

type RevokePostgresqlDbPrivilegeInvoker struct {
	*invoker.BaseInvoker
}

func (i *RevokePostgresqlDbPrivilegeInvoker) Invoke() (*model.RevokePostgresqlDbPrivilegeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RevokePostgresqlDbPrivilegeResponse), nil
	}
}

type SearchQueryScaleComputeFlavorsInvoker struct {
	*invoker.BaseInvoker
}

func (i *SearchQueryScaleComputeFlavorsInvoker) Invoke() (*model.SearchQueryScaleComputeFlavorsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SearchQueryScaleComputeFlavorsResponse), nil
	}
}

type SearchQueryScaleFlavorsInvoker struct {
	*invoker.BaseInvoker
}

func (i *SearchQueryScaleFlavorsInvoker) Invoke() (*model.SearchQueryScaleFlavorsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SearchQueryScaleFlavorsResponse), nil
	}
}

type SetDatabaseUserPrivilegeInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetDatabaseUserPrivilegeInvoker) Invoke() (*model.SetDatabaseUserPrivilegeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetDatabaseUserPrivilegeResponse), nil
	}
}

type SetPostgresqlDbUserPwdInvoker struct {
	*invoker.BaseInvoker
}

func (i *SetPostgresqlDbUserPwdInvoker) Invoke() (*model.SetPostgresqlDbUserPwdResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SetPostgresqlDbUserPwdResponse), nil
	}
}

type ShowInformationAboutDatabaseProxyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowInformationAboutDatabaseProxyInvoker) Invoke() (*model.ShowInformationAboutDatabaseProxyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowInformationAboutDatabaseProxyResponse), nil
	}
}

type ShowPostgresqlParamValueInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPostgresqlParamValueInvoker) Invoke() (*model.ShowPostgresqlParamValueResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPostgresqlParamValueResponse), nil
	}
}

type StartDatabaseProxyInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartDatabaseProxyInvoker) Invoke() (*model.StartDatabaseProxyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartDatabaseProxyResponse), nil
	}
}

type StopDatabaseProxyInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopDatabaseProxyInvoker) Invoke() (*model.StopDatabaseProxyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopDatabaseProxyResponse), nil
	}
}

type UpdateDbUserPrivilegeInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDbUserPrivilegeInvoker) Invoke() (*model.UpdateDbUserPrivilegeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDbUserPrivilegeResponse), nil
	}
}

type UpdatePostgresqlDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePostgresqlDatabaseInvoker) Invoke() (*model.UpdatePostgresqlDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePostgresqlDatabaseResponse), nil
	}
}

type UpdatePostgresqlDbUserCommentInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePostgresqlDbUserCommentInvoker) Invoke() (*model.UpdatePostgresqlDbUserCommentResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePostgresqlDbUserCommentResponse), nil
	}
}

type UpdatePostgresqlParameterValueInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePostgresqlParameterValueInvoker) Invoke() (*model.UpdatePostgresqlParameterValueResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePostgresqlParameterValueResponse), nil
	}
}

type UpdateReadWeightInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateReadWeightInvoker) Invoke() (*model.UpdateReadWeightResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateReadWeightResponse), nil
	}
}

type AllowSqlserverDbUserPrivilegeInvoker struct {
	*invoker.BaseInvoker
}

func (i *AllowSqlserverDbUserPrivilegeInvoker) Invoke() (*model.AllowSqlserverDbUserPrivilegeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AllowSqlserverDbUserPrivilegeResponse), nil
	}
}

type BatchAddMsdtcsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchAddMsdtcsInvoker) Invoke() (*model.BatchAddMsdtcsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchAddMsdtcsResponse), nil
	}
}

type CreateSqlserverDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSqlserverDatabaseInvoker) Invoke() (*model.CreateSqlserverDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSqlserverDatabaseResponse), nil
	}
}

type CreateSqlserverDbUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSqlserverDbUserInvoker) Invoke() (*model.CreateSqlserverDbUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSqlserverDbUserResponse), nil
	}
}

type DeleteSqlserverDatabaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSqlserverDatabaseInvoker) Invoke() (*model.DeleteSqlserverDatabaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSqlserverDatabaseResponse), nil
	}
}

type DeleteSqlserverDatabaseExInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSqlserverDatabaseExInvoker) Invoke() (*model.DeleteSqlserverDatabaseExResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSqlserverDatabaseExResponse), nil
	}
}

type DeleteSqlserverDbUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSqlserverDbUserInvoker) Invoke() (*model.DeleteSqlserverDbUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSqlserverDbUserResponse), nil
	}
}

type ListAuthorizedSqlserverDbUsersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAuthorizedSqlserverDbUsersInvoker) Invoke() (*model.ListAuthorizedSqlserverDbUsersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAuthorizedSqlserverDbUsersResponse), nil
	}
}

type ListMsdtcHostsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListMsdtcHostsInvoker) Invoke() (*model.ListMsdtcHostsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListMsdtcHostsResponse), nil
	}
}

type ListSqlserverDatabasesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSqlserverDatabasesInvoker) Invoke() (*model.ListSqlserverDatabasesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSqlserverDatabasesResponse), nil
	}
}

type ListSqlserverDbUsersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSqlserverDbUsersInvoker) Invoke() (*model.ListSqlserverDbUsersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSqlserverDbUsersResponse), nil
	}
}

type ModifyCollationInvoker struct {
	*invoker.BaseInvoker
}

func (i *ModifyCollationInvoker) Invoke() (*model.ModifyCollationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ModifyCollationResponse), nil
	}
}

type RevokeSqlserverDbUserPrivilegeInvoker struct {
	*invoker.BaseInvoker
}

func (i *RevokeSqlserverDbUserPrivilegeInvoker) Invoke() (*model.RevokeSqlserverDbUserPrivilegeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RevokeSqlserverDbUserPrivilegeResponse), nil
	}
}
