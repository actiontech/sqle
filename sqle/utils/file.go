package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"
)

// Grant file owners read and write execution permissions, group and other users read-only permissions
const OwnerPrivilegedAccessMode fs.FileMode = 0740

func EnsureFilePathWithPermission(filePath string, perm os.FileMode) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err = MkdirWithPermission(filePath, perm)
		if err != nil {
			return err
		}
	} else {
		err = EnsureFilePermission(filePath, perm)
		if err != nil {
			return err
		}
	}
	return nil
}

func MkdirWithPermission(filePath string, perm os.FileMode) error {
	err := os.Mkdir(filePath, os.ModeDir)
	if err != nil {
		return err
	}
	err = os.Chmod(filePath, perm)
	if err != nil {
		return err
	}
	return nil
}

func EnsureFilePermission(filePath string, perm os.FileMode) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	currentPerm := info.Mode().Perm()
	// 使用按位与来检查每一个权限位是否满足要求
	if currentPerm&perm != perm {
		return os.Chmod(filePath, perm)
	}
	return nil
}

// SaveFile 从 io.ReadSeeker 读取内容并保存到指定路径的文件中。
func SaveFile(file io.ReadSeeker, targetPath string) (err error) {
	// 从文件头开始读取
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// 创建目标文件
	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	// 保存文件
	_, err = io.Copy(target, file)
	if err != nil {
		return err
	}
	return nil
}

// ExportDataResult 导出数据的结果
type ExportDataResult struct {
	Content     []byte
	ContentType string
	FileName    string
}

// ExportDataAsExcel 将数据导出为 Excel 格式
// header: 表头字符串数组
// rows: 数据行，二维字符串数组
// fileNamePrefix: 文件名前缀，会自动添加时间戳和 .xlsx 扩展名
// prependRows: 可选的前置行，会在表头之前写入（可以为 nil）
func ExportDataAsExcel(header []string, rows [][]string, fileNamePrefix string, prependRows ...[][]string) (*ExportDataResult, error) {
	excelBuilder, err := NewExcelBuilder()
	if err != nil {
		return nil, fmt.Errorf("create excel builder failed: %v", err)
	}
	defer excelBuilder.Close()

	// 如果有前置行，先写入前置行
	if len(prependRows) > 0 && prependRows[0] != nil {
		if err = excelBuilder.WriteRows(prependRows[0]); err != nil {
			return nil, fmt.Errorf("write excel prepend rows failed: %v", err)
		}
	}

	if err = excelBuilder.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("write excel header failed: %v", err)
	}
	if err = excelBuilder.WriteRows(rows); err != nil {
		return nil, fmt.Errorf("write excel rows failed: %v", err)
	}

	fileBytes, err := excelBuilder.FlushAndGetBuffer()
	if err != nil {
		return nil, fmt.Errorf("flush excel buffer failed: %v", err)
	}

	fileName := fmt.Sprintf("%s_%s.xlsx", fileNamePrefix, time.Now().Format("20060102150405"))
	return &ExportDataResult{
		Content:     fileBytes,
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		FileName:    fileName,
	}, nil
}

// ExportDataAsCSV 将数据导出为 CSV 格式
// header: 表头字符串数组
// rows: 数据行，二维字符串数组
// fileNamePrefix: 文件名前缀，会自动添加时间戳和 .csv 扩展名
// prependRows: 可选的前置行，会在表头之前写入（可以为 nil）
func ExportDataAsCSV(header []string, rows [][]string, fileNamePrefix string, prependRows ...[][]string) (*ExportDataResult, error) {
	csvBuilder := NewCSVBuilder()

	// 如果有前置行，先写入前置行
	if len(prependRows) > 0 && prependRows[0] != nil {
		if err := csvBuilder.WriteRows(prependRows[0]); err != nil {
			return nil, fmt.Errorf("write csv prepend rows failed: %v", err)
		}
	}

	if err := csvBuilder.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("write csv header failed: %v", err)
	}
	if err := csvBuilder.WriteRows(rows); err != nil {
		return nil, fmt.Errorf("write csv rows failed: %v", err)
	}

	fileBytes := csvBuilder.FlushAndGetBuffer().Bytes()
	fileName := fmt.Sprintf("%s_%s.csv", fileNamePrefix, time.Now().Format("20060102150405"))
	return &ExportDataResult{
		Content:     fileBytes,
		ContentType: "text/csv",
		FileName:    fileName,
	}, nil
}

// NormalizeExportFormat 规范化导出格式，支持 string 和 *string 类型
// 如果格式不是 "csv" 或 "excel"，则默认返回 "csv"
// 参数 format 可以是 string 或 *string 类型
func NormalizeExportFormat(format interface{}) string {
	var formatStr string

	switch v := format.(type) {
	case string:
		formatStr = v
	case *string:
		if v != nil {
			formatStr = *v
		}
	default:
		return "csv"
	}

	if formatStr == "" {
		return "csv"
	}

	formatStr = strings.ToLower(formatStr)
	if formatStr != "csv" && formatStr != "excel" {
		return "csv"
	}

	return formatStr
}
