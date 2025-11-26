package pprof

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/actiontech/sqle/sqle/log"
)

const (
	pprofDirName = "pprof"
)

// CollectHeapProfile 采集堆内存 profile 并保存到文件
func CollectHeapProfile(logPath string) error {
	return collectProfile("heap", logPath, func(f *os.File) error {
		runtime.GC()
		return pprof.WriteHeapProfile(f)
	})
}

// CollectGoroutineProfile 采集 goroutine profile 并保存到文件
func CollectGoroutineProfile(logPath string) error {
	return collectProfile("goroutine", logPath, func(f *os.File) error {
		return pprof.Lookup("goroutine").WriteTo(f, 0)
	})
}

// CollectAllocsProfile 采集内存分配 profile 并保存到文件
func CollectAllocsProfile(logPath string) error {
	return collectProfile("allocs", logPath, func(f *os.File) error {
		return pprof.Lookup("allocs").WriteTo(f, 0)
	})
}

// CollectBlockProfile 采集阻塞 profile 并保存到文件
func CollectBlockProfile(logPath string) error {
	return collectProfile("block", logPath, func(f *os.File) error {
		return pprof.Lookup("block").WriteTo(f, 0)
	})
}

// CollectMutexProfile 采集互斥锁 profile 并保存到文件
func CollectMutexProfile(logPath string) error {
	return collectProfile("mutex", logPath, func(f *os.File) error {
		return pprof.Lookup("mutex").WriteTo(f, 0)
	})
}

// CollectCPUProfile 采集 CPU profile 并保存到文件（持续指定秒数）
func CollectCPUProfile(logPath string, duration time.Duration) error {
	return collectProfile("cpu", logPath, func(f *os.File) error {
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}
		time.Sleep(duration)
		pprof.StopCPUProfile()
		return nil
	})
}

// collectProfile 通用的 profile 采集函数
func collectProfile(profileType, logPath string, writeFunc func(*os.File) error) error {
	pprofDir := filepath.Join(logPath, pprofDirName)
	if err := os.MkdirAll(pprofDir, 0755); err != nil {
		return fmt.Errorf("failed to create pprof directory: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.prof", profileType, timestamp)
	filePath := filepath.Join(pprofDir, filename)

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create profile file: %v", err)
	}
	defer f.Close()

	if err := writeFunc(f); err != nil {
		return fmt.Errorf("failed to write profile: %v", err)
	}

	log.Logger().Infof("pprof %s profile saved to: %s", profileType, filePath)
	return nil
}

// CollectAllProfiles 采集所有类型的 profile（除了 CPU，因为 CPU 需要持续时间）
func CollectAllProfiles(logPath string) error {
	profiles := []struct {
		name string
		fn   func(string) error
	}{
		{"heap", CollectHeapProfile},
		{"goroutine", CollectGoroutineProfile},
		{"allocs", CollectAllocsProfile},
		{"block", CollectBlockProfile},
		{"mutex", CollectMutexProfile},
	}

	var lastErr error
	for _, p := range profiles {
		if err := p.fn(logPath); err != nil {
			log.Logger().Errorf("failed to collect %s profile: %v", p.name, err)
			lastErr = err
		}
	}

	return lastErr
}

// StartPeriodicCollection 启动定期自动采集 pprof profile
// interval: 采集间隔时间，如果为 0 则不启用定期采集
func StartPeriodicCollection(logPath string, interval time.Duration) {
	if interval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Logger().Infof("Starting periodic pprof collection, interval: %v", interval)

		for range ticker.C {
			log.Logger().Infof("Periodic pprof collection triggered")
			if err := CollectAllProfiles(logPath); err != nil {
				log.Logger().Errorf("Periodic pprof collection failed: %v", err)
			}
		}
	}()
}
