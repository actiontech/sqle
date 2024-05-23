//go:build enterprise
// +build enterprise

package v1

import (
	"testing"

	"github.com/actiontech/sqle/sqle/model"
)

type testContent struct {
	inputFilesToSort []FileToSort
	results          []uint
}

var testContents = []testContent{
	{
		inputFilesToSort: []FileToSort{
			{
				FileID:   91,
				NewIndex: 7,
			},
		},
		results: []uint{90, 92, 93, 94, 95, 96, 97, 91},
	},
	{
		inputFilesToSort: []FileToSort{
			{
				FileID:   97,
				NewIndex: 1,
			},
		},
		results: []uint{90, 97, 91, 92, 93, 94, 95, 96},
	},
	{
		inputFilesToSort: []FileToSort{
			{
				FileID:   91,
				NewIndex: 5,
			},
			{
				FileID:   93,
				NewIndex: 1,
			},
		},
		results: []uint{90, 93, 92, 94, 95, 91, 96, 97},
	},
	{
		inputFilesToSort: []FileToSort{
			{
				FileID:   91,
				NewIndex: 3,
			},
			{
				FileID:   92,
				NewIndex: 4,
			},
			{
				FileID:   93,
				NewIndex: 5,
			},
		},
		results: []uint{90, 94, 95, 91, 92, 93, 96, 97},
	},
	{
		inputFilesToSort: []FileToSort{
			{
				FileID:   97,
				NewIndex: 1,
			},
			{
				FileID:   96,
				NewIndex: 2,
			},
			{
				FileID:   95,
				NewIndex: 3,
			},
		},
		results: []uint{90, 97, 96, 95, 91, 92, 93, 94},
	},
}

func TestCheckRedundantIndex(t *testing.T) {
	originSortedFileId := []uint{90, 91, 92, 93, 94, 95, 96, 97}
	auditFiles := mockAuditFilesByIds(originSortedFileId)

	for _, content := range testContents {
		sortedAuditFIles := reorderFiles(content.inputFilesToSort, auditFiles)
		sortedFileId := []uint{}
		for _, file := range sortedAuditFIles {
			sortedFileId = append(sortedFileId, file.ID)
		}
		if !equalFileIdOrder(sortedFileId, content.results) {
			t.Errorf("Expected %v, but got %v", content.results, sortedFileId)
		}
	}
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
