package v1

import (
	"bytes"
	e "errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	xmlParser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/bodgit/sevenzip"
	"github.com/labstack/echo/v4"
)

// getSqlsFrom7z 从 7z 文件中提取 SQL 语句。
// 使用 github.com/bodgit/sevenzip 库解压 7z 文件，遍历内部文件并调用 processArchiveEntry 处理。
// 函数签名与 getSqlsFromZip / getSqlsFromRar 保持一致。
// 注意：sevenzip 需要 io.ReaderAt + int64 size（不同于 RAR 的 io.Reader），需先将上传文件读入 bytes.Reader。
func getSqlsFrom7z(c echo.Context) (sqlsFromSQLFile []SQLsFromSQLFile, sqlsFromXML []SQLFromXML, skippedCount int, exist bool, err error) {
	file, err := c.FormFile(InputZipFileName)
	if err == http.ErrMissingFile {
		return nil, nil, 0, false, nil
	}
	if err != nil {
		return nil, nil, 0, false, err
	}

	f, err := file.Open()
	if err != nil {
		return nil, nil, 0, false, err
	}
	defer f.Close()

	// 使用 archiveConfig 进行压缩包总大小限制检查（上传文件大小预检）
	if err := defaultArchiveConfig.checkSize(0, file.Size); err != nil {
		return nil, nil, 0, false, err
	}

	// sevenzip 需要 io.ReaderAt 接口，将上传文件内容读入 bytes.Reader
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, 0, false, fmt.Errorf("read 7z file into memory failed: %v", err)
	}

	sqlsFromSQLFile, sqlsFromXML, skippedCount, err = process7zContent(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return nil, nil, 0, false, err
	}

	return sqlsFromSQLFile, sqlsFromXML, skippedCount, true, nil
}

// process7zContent 从 io.ReaderAt 中读取 7z 内容，遍历 entry 并提取 SQL。
// 该函数封装了 7z 解压的核心逻辑，独立于 echo.Context，便于单元测试。
func process7zContent(r io.ReaderAt, size int64) (sqlsFromSQLFile []SQLsFromSQLFile, sqlsFromXML []SQLFromXML, skippedCount int, err error) {
	// 使用 sevenzip.NewReader 打开 7z 文件
	szr, err := sevenzip.NewReader(r, size)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("open 7z file failed: %v", err)
	}

	var xmlContents []xmlParser.XmlFile
	var totalSize int64
	var fileCount int

	for _, f := range szr.File {
		// 跳过目录
		if f.FileInfo().IsDir() {
			continue
		}

		// 文件数量限制检查
		fileCount++
		if err := defaultArchiveConfig.checkFileCount(fileCount); err != nil {
			return nil, nil, 0, err
		}

		// 嵌套压缩包检查：depth=1 时跳过内层压缩包
		ext := strings.ToLower(filepath.Ext(f.Name))
		if supportedArchiveExts[ext] {
			continue
		}

		// 打开并读取文件内容
		rc, err := f.Open()
		if err != nil {
			return nil, nil, 0, fmt.Errorf("open 7z entry %q failed: %v", f.Name, err)
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, nil, 0, fmt.Errorf("read 7z entry %q content failed: %v", f.Name, err)
		}

		// 累计大小限制检查
		totalSize += int64(len(content))
		if err := defaultArchiveConfig.checkSize(0, totalSize); err != nil {
			return nil, nil, 0, err
		}

		// 委托 processArchiveEntry 按扩展名分发处理
		sqlContent, xmlContent, isSupported, err := processArchiveEntry(f.Name, content)
		if err != nil {
			if e.Is(err, utils.ErrUnknownEncoding) {
				log.NewEntry().WithField("convert_to_utf8", f.Name).Errorf("convert to utf8 failed: %v", err)
				continue
			}
			return nil, nil, 0, err
		}
		if !isSupported {
			skippedCount++
			continue
		}

		if xmlContent != nil {
			xmlContents = append(xmlContents, *xmlContent)
		} else if sqlContent != "" {
			sqlsFromSQLFile = append(sqlsFromSQLFile, SQLsFromSQLFile{
				FilePath: f.Name,
				SQLs:     sqlContent,
			})
		}
	}

	// parse xml content
	// xml文件需要把所有文件内容同时解析，否则会无法解析跨namespace引用的SQL
	{
		sqlsFromXmls, err := parseXMLsWithFilePath(xmlContents)
		if err != nil {
			return nil, nil, 0, err
		}
		sqlsFromXML = append(sqlsFromXML, sqlsFromXmls...)
	}

	// 按文件名自然排序，确保SQL按文件顺序执行
	sort.Slice(sqlsFromSQLFile, func(i, j int) bool {
		return utils.CompareNatural(sqlsFromSQLFile[i].FilePath, sqlsFromSQLFile[j].FilePath)
	})
	sort.Slice(sqlsFromXML, func(i, j int) bool {
		return utils.CompareNatural(sqlsFromXML[i].FilePath, sqlsFromXML[j].FilePath)
	})

	return sqlsFromSQLFile, sqlsFromXML, skippedCount, nil
}
