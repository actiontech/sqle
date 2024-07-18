//go:build enterprise
// +build enterprise

package tbase_audit_log

import (
	"sort"
	"time"
	"os"
	"path/filepath"
)

type FileInfo struct {
    Name    string
    ModTime time.Time
}

func sortFilesByModTime(fileInfos []FileInfo) []FileInfo {
    sort.Slice(fileInfos, func(i, j int) bool {
        return fileInfos[i].ModTime.After(fileInfos[j].ModTime)
    })

    return fileInfos
}

func getSortedFilesByModTime(filePaths []string) ([]string, error) {
    var fileInfos []FileInfo
    for _, filePath := range filePaths {
        info, err := os.Stat(filePath)
        if err != nil {
			return nil, err
        }
        fileInfos = append(fileInfos, FileInfo{
            Name:    filePath,
            ModTime: info.ModTime(),
        })
    }

	sortedFileInfos := sortFilesByModTime(fileInfos)
	sortedFilePaths := []string{}
	for _, fileInfo := range sortedFileInfos {
		sortedFilePaths = append(sortedFilePaths, fileInfo.Name)
	}

	return sortedFilePaths, nil
}

func getSortedFilesByFolderPath(folderPath string) ([]string, error) {
	filePaths, err := filepath.Glob(folderPath)
    if err != nil {
        return nil, err
    }
	
	return getSortedFilesByModTime(filePaths)
}
