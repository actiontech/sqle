//go:build enterprise
// +build enterprise

package v1

import (
	"testing"

	"github.com/actiontech/sqle/sqle/model"
)

type testContent struct {
	inputFilesToSort  []FileToSort
	inputOriginFileID []uint
	results           []uint
}

var testContents = []testContent{
	// 文件不排序
	{
		inputFilesToSort:  []FileToSort{},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{1, 2, 3, 4, 5, 6, 7, 8},
	},
	// 将第一个文件移到末尾
	{
		inputFilesToSort: []FileToSort{
			{FileID: 1, NewIndex: 7},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{2, 3, 4, 5, 6, 7, 8, 1},
	},
	// 文件最终位置不发生改变
	{
		inputFilesToSort: []FileToSort{
			{FileID: 1, NewIndex: 0},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{1, 2, 3, 4, 5, 6, 7, 8},
	},
	// 将最后一个文件移到末尾
	{
		inputFilesToSort: []FileToSort{
			{FileID: 8, NewIndex: 0},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{8, 1, 2, 3, 4, 5, 6, 7},
	},
	// 移动中间的文件到任意为止
	{
		inputFilesToSort: []FileToSort{
			{FileID: 3, NewIndex: 0},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{3, 1, 2, 4, 5, 6, 7, 8},
	},
	// 所有文件都发生拖拽，但是顺序未发生改变
	{
		inputFilesToSort: []FileToSort{
			{FileID: 1, NewIndex: 0},
			{FileID: 2, NewIndex: 1},
			{FileID: 3, NewIndex: 2},
			{FileID: 4, NewIndex: 3},
			{FileID: 5, NewIndex: 4},
			{FileID: 6, NewIndex: 5},
			{FileID: 7, NewIndex: 6},
			{FileID: 8, NewIndex: 7},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{1, 2, 3, 4, 5, 6, 7, 8},
	},
	// 所有文件倒排
	{
		inputFilesToSort: []FileToSort{
			{FileID: 1, NewIndex: 7},
			{FileID: 2, NewIndex: 6},
			{FileID: 3, NewIndex: 5},
			{FileID: 4, NewIndex: 4},
			{FileID: 5, NewIndex: 3},
			{FileID: 6, NewIndex: 2},
			{FileID: 7, NewIndex: 1},
			{FileID: 8, NewIndex: 0},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{8, 7, 6, 5, 4, 3, 2, 1},
	},
	// 改变两个文件的顺序
	{
		inputFilesToSort: []FileToSort{
			{FileID: 1, NewIndex: 3},
			{FileID: 2, NewIndex: 4},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{3, 4, 5, 1, 2, 6, 7, 8},
	},
	{
		inputFilesToSort: []FileToSort{
			{FileID: 2, NewIndex: 5},
			{FileID: 5, NewIndex: 1},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{1, 5, 3, 4, 6, 2, 7, 8},
	},
	// 改变三个文件顺序
	{
		inputFilesToSort: []FileToSort{
			{FileID: 1, NewIndex: 3},
			{FileID: 2, NewIndex: 4},
			{FileID: 5, NewIndex: 1},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{3, 5, 4, 1, 2, 6, 7, 8},
	},
	{
		inputFilesToSort: []FileToSort{
			{FileID: 2, NewIndex: 5},
			{FileID: 5, NewIndex: 1},
			{FileID: 8, NewIndex: 0},
		},
		inputOriginFileID: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		results:           []uint{8, 5, 1, 3, 4, 2, 6, 7},
	},
}

// TODO (temporarily remove test cases)
func TestCheckRedundantIndex(t *testing.T) {
	// originSortedFileId := []uint{90, 91, 92, 93, 94, 95, 96, 97}
	// for _, content := range testContents {
	// 	auditFiles := mockAuditFilesByIds(content.inputOriginFileID)
	// 	sortedAuditFIles := reorderFiles(content.inputFilesToSort, auditFiles)
	// 	sortedFileId := []uint{}
	// 	for _, file := range sortedAuditFIles {
	// 		sortedFileId = append(sortedFileId, file.ID)
	// 	}
	// 	if !equalFileIdOrder(sortedFileId, content.results) {
	// 		t.Errorf("Expected %v, but got %v", content.results, sortedFileId)
	// 	}
	// }
}

func mockAuditFilesByIds(ids []uint) []*model.AuditFile {
	auditFiles := []*model.AuditFile{}
	for ind, id := range ids {
		auditFiles = append(auditFiles, &model.AuditFile{
			Model: model.Model{
				ID: id,
			},
			ExecOrder: uint(ind),
		})
	}
	return auditFiles
}

func equalFileIdOrder(sortedFileId, result []uint) bool {
	if len(sortedFileId) != len(result) {
		return false
	}
	for ind, fileId := range sortedFileId {
		if fileId != result[ind] {
			return false
		}
	}
	return true
}
