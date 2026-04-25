package v1

import (
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
	"github.com/labstack/echo/v4"
	"github.com/nwaples/rardecode"
)

// getSqlsFromRar 从 RAR 文件中提取 SQL 语句。
// 使用 github.com/nwaples/rardecode 库解压 RAR 文件，遍历内部文件并调用 processArchiveEntry 处理。
// 函数签名与 getSqlsFromZip 保持一致。
func getSqlsFromRar(c echo.Context) (sqlsFromSQLFile []SQLsFromSQLFile, sqlsFromXML []SQLFromXML, skippedCount int, exist bool, err error) {
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

	sqlsFromSQLFile, sqlsFromXML, skippedCount, err = processRarContent(f)
	if err != nil {
		return nil, nil, 0, false, err
	}

	return sqlsFromSQLFile, sqlsFromXML, skippedCount, true, nil
}

// processRarContent 从 io.Reader 中读取 RAR 内容，遍历 entry 并提取 SQL。
// 该函数封装了 RAR 解压的核心逻辑，独立于 echo.Context，便于单元测试。
func processRarContent(r io.Reader) (sqlsFromSQLFile []SQLsFromSQLFile, sqlsFromXML []SQLFromXML, skippedCount int, err error) {
	// 使用 rardecode.NewReader 打开 RAR 文件，密码参数为空字符串（不支持加密 RAR）
	rr, err := rardecode.NewReader(r, "")
	if err != nil {
		return nil, nil, 0, fmt.Errorf("open rar file failed: %v", err)
	}

	var xmlContents []xmlParser.XmlFile
	var totalSize int64
	var fileCount int

	for {
		header, err := rr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, 0, fmt.Errorf("read rar entry failed: %v", err)
		}

		// 跳过目录
		if header.IsDir {
			continue
		}

		// 文件数量限制检查
		fileCount++
		if err := defaultArchiveConfig.checkFileCount(fileCount); err != nil {
			return nil, nil, 0, err
		}

		// 嵌套压缩包检查：depth=1 时跳过内层压缩包
		ext := strings.ToLower(filepath.Ext(header.Name))
		if supportedArchiveExts[ext] {
			continue
		}

		// 读取文件内容
		content, err := io.ReadAll(rr)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("read rar entry content failed: %v", err)
		}

		// 累计大小限制检查
		totalSize += int64(len(content))
		if err := defaultArchiveConfig.checkSize(0, totalSize); err != nil {
			return nil, nil, 0, err
		}

		// 委托 processArchiveEntry 按扩展名分发处理
		sqlContent, xmlContent, isSupported, err := processArchiveEntry(header.Name, content)
		if err != nil {
			if e.Is(err, utils.ErrUnknownEncoding) {
				log.NewEntry().WithField("convert_to_utf8", header.Name).Errorf("convert to utf8 failed: %v", err)
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
				FilePath: header.Name,
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
