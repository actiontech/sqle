package state

type CollectorStats struct {
	GoVersion string

	MemoryHeapAllocatedBytes uint64 // Bytes allocated and not yet freed
	MemoryHeapObjects        uint64 // Total number of allocated objects
	MemorySystemBytes        uint64 // Bytes obtained from system (sum of heap and fixed-size structures)
	MemoryRssBytes           uint64 // Memory allocated in bytes as seen by the OS

	ActiveGoroutines int32

	CgoCalls int64
}

type DiffedCollectorStats CollectorStats

func (curr CollectorStats) DiffSince(prev CollectorStats) DiffedCollectorStats {
	return DiffedCollectorStats{
		GoVersion:                curr.GoVersion,
		MemoryHeapAllocatedBytes: curr.MemoryHeapAllocatedBytes,
		MemoryHeapObjects:        curr.MemoryHeapObjects,
		MemorySystemBytes:        curr.MemorySystemBytes,
		MemoryRssBytes:           curr.MemoryRssBytes,
		ActiveGoroutines:         curr.ActiveGoroutines,
		CgoCalls:                 curr.CgoCalls - prev.CgoCalls,
	}
}
