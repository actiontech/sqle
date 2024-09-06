package state

import "time"

// SystemState - All kinds of system-related information and metrics
type SystemState struct {
	Info         SystemInfo
	Scheduler    Scheduler
	Memory       Memory
	CPUInfo      CPUInformation
	CPUStats     CPUStatisticMap
	NetworkStats NetworkStatsMap

	Disks          DiskMap
	DiskStats      DiskStatsMap
	DiskPartitions DiskPartitionMap

	DataDirectoryPartition string // Partition that the data directory lives on (identified by the partition's mountpoint)
	XlogPartition          string // Partition that the WAL directory lives on
	XlogUsedBytes          uint64
}

// SystemType - Enum that describes which kind of system we're monitoring
type SystemType int

// Treat this list as append-only and never change the order
const (
	SelfHostedSystem SystemType = iota
	AmazonRdsSystem
	HerokuSystem
	GoogleCloudSQLSystem
	AzureDatabaseSystem
	CrunchyBridgeSystem
	AivenSystem
)

type SystemInfo struct {
	Type        SystemType
	SystemScope string
	SystemID    string

	SelfHosted   *SystemInfoSelfHosted
	AmazonRds    *SystemInfoAmazonRds
	ResourceTags map[string]string

	BootTime time.Time
}

// SystemInfoSelfHosted - System information for self-hosted systems (both physical and virtual)
type SystemInfoSelfHosted struct {
	Hostname                 string
	Architecture             string
	OperatingSystem          string
	Platform                 string
	PlatformFamily           string
	PlatformVersion          string
	VirtualizationSystem     string // Name of the virtualization system (only if we're a guest)
	KernelVersion            string
	DatabaseSystemIdentifier string
}

// SystemInfoAmazonRds - System information for Amazon RDS systems
type SystemInfoAmazonRds struct {
	Region                      string
	InstanceClass               string
	InstanceID                  string
	Status                      string
	AvailabilityZone            string
	PubliclyAccessible          bool
	MultiAz                     bool
	SecondaryAvailabilityZone   string
	CaCertificate               string
	AutoMinorVersionUpgrade     bool
	PreferredMaintenanceWindow  string
	PreferredBackupWindow       string
	LatestRestorableTime        time.Time
	BackupRetentionPeriodDays   int32
	MasterUsername              string
	InitialDbName               string
	CreatedAt                   time.Time
	StorageProvisionedIOPS      int32
	StorageAllocatedGigabytes   int32
	StorageEncrypted            bool
	StorageType                 string
	EnhancedMonitoring          bool
	PerformanceInsights         bool
	PostgresLogExport           bool
	IAMAuthentication           bool
	DeletionProtection          bool
	ParameterApplyStatus        string
	ParameterPgssEnabled        bool
	ParameterAutoExplainEnabled bool
	IsAuroraPostgres            bool
}

// Scheduler - Information about the OS scheduler
type Scheduler struct {
	Loadavg1min  float64
	Loadavg5min  float64
	Loadavg15min float64
}

// Memory - Metrics related to system memory
type Memory struct {
	TotalBytes      uint64
	CachedBytes     uint64
	BuffersBytes    uint64
	FreeBytes       uint64
	WritebackBytes  uint64
	DirtyBytes      uint64
	SlabBytes       uint64
	MappedBytes     uint64
	PageTablesBytes uint64
	ActiveBytes     uint64
	InactiveBytes   uint64
	AvailableBytes  uint64
	SwapUsedBytes   uint64
	SwapTotalBytes  uint64

	HugePagesSizeBytes uint64
	HugePagesFree      uint64
	HugePagesTotal     uint64
	HugePagesReserved  uint64
	HugePagesSurplus   uint64

	ApplicationBytes uint64
}

type CPUInformation struct {
	Model             string
	CacheSizeBytes    int32
	SpeedMhz          float64
	SocketCount       int32
	PhysicalCoreCount int32
	LogicalCoreCount  int32
}

// CPUStatisticMap - Map of all CPU statistics (Key = CPU ID)
type CPUStatisticMap map[string]CPUStatistic

// CPUStatistic - Statistics for a single CPU core
type CPUStatistic struct {
	DiffedOnInput bool // True if has already been diffed on input (and we can simply copy the diff)
	DiffedValues  *DiffedSystemCPUStats

	// Seconds (counter values that need to be diff-ed between runs)
	UserSeconds      float64
	SystemSeconds    float64
	IdleSeconds      float64
	NiceSeconds      float64
	IowaitSeconds    float64
	IrqSeconds       float64
	SoftIrqSeconds   float64
	StealSeconds     float64
	GuestSeconds     float64
	GuestNiceSeconds float64
}

// DiffedSystemCPUStatsMap - Map of all CPU statistics (Key = CPU ID)
type DiffedSystemCPUStatsMap map[string]DiffedSystemCPUStats

// DiffedSystemCPUStats - CPU statistics as percentages
type DiffedSystemCPUStats struct {
	UserPercent      float64
	SystemPercent    float64
	IdlePercent      float64
	NicePercent      float64
	IowaitPercent    float64
	IrqPercent       float64
	SoftIrqPercent   float64
	StealPercent     float64
	GuestPercent     float64
	GuestNicePercent float64
}

// NetworkStatsMap - Map of all network statistics (Key = Interface Name)
type NetworkStatsMap map[string]NetworkStats

// NetworkStats - Information about the network activity on a single interface
type NetworkStats struct {
	DiffedOnInput bool // True if has already been diffed on input (and we can simply copy the diff)
	DiffedValues  *DiffedNetworkStats

	ReceiveThroughputBytes  uint64
	TransmitThroughputBytes uint64
}

// DiffedNetworkStats - Network statistics for a single interface as a diff
type DiffedNetworkStats struct {
	ReceiveThroughputBytesPerSecond  uint64
	TransmitThroughputBytesPerSecond uint64
}

// DiffedNetworkStatsMap - Map of network statistics as a diff (Key = Interface Name)
type DiffedNetworkStatsMap map[string]DiffedNetworkStats

// Disk - Information about an individual disk device in the system
type Disk struct {
	DiskType        string // Disk type (hdd/sdd/io1/gp2)
	Scheduler       string // Linux Scheduler (noop/anticipatory/deadline/cfq)
	ProvisionedIOPS uint32 // If applicable, how many IOPS are provisioned for this device
	Encrypted       bool   // If applicable, is this device encrypted? (default false)

	ComponentDisks []string // Identifiers for component disks (e.g. for a software RAID)
}

// DiskStats - Statistics about an individual disk device in the system
type DiskStats struct {
	DiffedOnInput bool // True if has already been diffed on input (and we can simply copy the diff)
	DiffedValues  *DiffedDiskStats

	// Counter values
	ReadsCompleted  uint64 // /proc/diskstats 4 - reads completed successfully
	ReadsMerged     uint64 // /proc/diskstats 5 - reads merged
	BytesRead       uint64 // /proc/diskstat 6 - sectors read, multiplied by sector size
	ReadTimeMs      uint64 // /proc/diskstat 7 - time spent reading (ms)
	WritesCompleted uint64 // /proc/diskstats 8 - writes completed
	WritesMerged    uint64 // /proc/diskstats 9 - writes merged
	BytesWritten    uint64 // /proc/diskstat 10 - sectors written, multiplied by sector size
	WriteTimeMs     uint64 // /proc/diskstat 11 - time spent writing (ms)
	AvgQueueSize    int32  // /proc/diskstat 12 - I/Os currently in progress
	IoTime          uint64 // /proc/diskstat 13 - time spent doing I/Os (ms)
}

type DiffedDiskStats struct {
	ReadOperationsPerSecond float64 // The average number of read requests that were issued to the device per second
	ReadsMergedPerSecond    float64 // The average number of read requests merged per second that were queued to the device
	BytesReadPerSecond      float64 // The average number of bytes read from the device per second
	AvgReadLatency          float64 // The average time (in milliseconds) for read requests issued to the device to be served

	WriteOperationsPerSecond float64 // The average number of write requests that were issued to the device per second
	WritesMergedPerSecond    float64 // The average number of write requests merged per second that were queued to the device
	BytesWrittenPerSecond    float64 // The average number of bytes written to the device per second
	AvgWriteLatency          float64 // The average time (in milliseconds) for write requests issued to the device to be served

	AvgQueueSize       int32   // Average I/O operations in flight at the same time (waiting or worked on by the device)
	UtilizationPercent float64 // Percentage of CPU time during which I/O requests were issued to the device (bandwidth utilization for the device)
}

// DiskMap - Map of all disks (key = device name)
type DiskMap map[string]Disk

// DiskStatsMap - Map of all disk statistics (key = device name)
type DiskStatsMap map[string]DiskStats

// DiffedDiskStatsMap - Map of all diffed disk statistics (key = device name)
type DiffedDiskStatsMap map[string]DiffedDiskStats

// DiskPartition - Information and statistics about one of the disk partitions in the system
type DiskPartition struct {
	DiskName       string // Name of the base device disk that this partition resides on (e.g. /dev/sda)
	PartitionName  string // Platform-specific name of the partition (e.g. /dev/sda1)
	FilesystemType string
	FilesystemOpts string

	UsedBytes  uint64
	TotalBytes uint64
}

// DiskPartitionMap - Map of all disk partitions (key = mountpoint)
type DiskPartitionMap map[string]DiskPartition

// ---

// DiffSince - Calculate the diff between two CPU stats runs
func (curr CPUStatistic) DiffSince(prev CPUStatistic) DiffedSystemCPUStats {
	userSecs := curr.UserSeconds - prev.UserSeconds
	systemSecs := curr.SystemSeconds - prev.SystemSeconds
	idleSecs := curr.IdleSeconds - prev.IdleSeconds
	niceSecs := curr.NiceSeconds - prev.NiceSeconds
	iowaitSecs := curr.IowaitSeconds - prev.IowaitSeconds
	irqSecs := curr.IrqSeconds - prev.IrqSeconds
	softIrqSecs := curr.SoftIrqSeconds - prev.SoftIrqSeconds
	stealSecs := curr.StealSeconds - prev.StealSeconds
	guestSecs := curr.GuestSeconds - prev.GuestSeconds
	guestNiceSecs := curr.GuestNiceSeconds - prev.GuestNiceSeconds
	totalSecs := userSecs + systemSecs + idleSecs + niceSecs + iowaitSecs + irqSecs + softIrqSecs + stealSecs + guestSecs + guestNiceSecs

	if totalSecs == 0 {
		return DiffedSystemCPUStats{}
	}

	return DiffedSystemCPUStats{
		UserPercent:      userSecs / totalSecs * 100,
		SystemPercent:    systemSecs / totalSecs * 100,
		IdlePercent:      idleSecs / totalSecs * 100,
		NicePercent:      niceSecs / totalSecs * 100,
		IowaitPercent:    iowaitSecs / totalSecs * 100,
		IrqPercent:       irqSecs / totalSecs * 100,
		SoftIrqPercent:   softIrqSecs / totalSecs * 100,
		StealPercent:     stealSecs / totalSecs * 100,
		GuestPercent:     guestSecs / totalSecs * 100,
		GuestNicePercent: guestNiceSecs / totalSecs * 100,
	}
}

// DiffSince - Calculate the diff between two network stats runs
func (curr NetworkStats) DiffSince(prev NetworkStats, collectedIntervalSecs uint32) DiffedNetworkStats {
	return DiffedNetworkStats{
		ReceiveThroughputBytesPerSecond:  (curr.ReceiveThroughputBytes - prev.ReceiveThroughputBytes) / uint64(collectedIntervalSecs),
		TransmitThroughputBytesPerSecond: (curr.TransmitThroughputBytes - prev.TransmitThroughputBytes) / uint64(collectedIntervalSecs),
	}
}

// DiffSince - Calculate the diff between two disk stats runs
func (curr DiskStats) DiffSince(prev DiskStats, collectedIntervalSecs uint32) DiffedDiskStats {
	reads := float64(curr.ReadsCompleted - prev.ReadsCompleted)
	writes := float64(curr.WritesCompleted - prev.WritesCompleted)

	diffed := DiffedDiskStats{
		ReadOperationsPerSecond:  reads / float64(collectedIntervalSecs),
		ReadsMergedPerSecond:     float64(curr.ReadsMerged-prev.ReadsMerged) / float64(collectedIntervalSecs),
		BytesReadPerSecond:       float64(curr.BytesRead-prev.BytesRead) / float64(collectedIntervalSecs),
		WriteOperationsPerSecond: writes / float64(collectedIntervalSecs),
		WritesMergedPerSecond:    float64(curr.WritesMerged-prev.WritesMerged) / float64(collectedIntervalSecs),
		BytesWrittenPerSecond:    float64(curr.BytesWritten-prev.BytesWritten) / float64(collectedIntervalSecs),
		AvgQueueSize:             curr.AvgQueueSize,
		UtilizationPercent:       100 * float64(curr.IoTime-prev.IoTime) / float64(1000*collectedIntervalSecs),
	}

	if reads > 0 {
		diffed.AvgReadLatency = float64(curr.ReadTimeMs-prev.ReadTimeMs) / reads
	}

	if writes > 0 {
		diffed.AvgWriteLatency = float64(curr.WriteTimeMs-prev.WriteTimeMs) / writes
	}

	return diffed
}
