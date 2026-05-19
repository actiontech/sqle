package v1

import (
	"os"
	"path/filepath"
	"strings"

	javaParser "github.com/actiontech/java-sql-extractor/parser"
	xmlParser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/utils"
)

// processArchiveEntry 处理压缩包内的单个文件，根据扩展名分发处理。
// 返回值：
//   - sqlContent: 从文件中提取的 SQL 内容（.sql/.txt 为文件原文，.java 为提取后的 SQL 语句）
//   - xmlContent: XML 文件内容（仅 .xml 文件非空，用于后续跨 namespace 解析）
//   - isSupported: 该文件格式是否受支持
//   - err: 处理过程中的错误
func processArchiveEntry(filename string, content []byte) (sqlContent string, xmlContent *xmlParser.XmlFile, isSupported bool, err error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".sql", ".txt":
		utf8Content, err := utils.ConvertToUtf8(content)
		if err != nil {
			return "", nil, true, err
		}
		return string(utf8Content), nil, true, nil

	case ".xml":
		utf8Content, err := utils.ConvertToUtf8(content)
		if err != nil {
			return "", nil, true, err
		}
		return "", &xmlParser.XmlFile{
			FilePath: filename,
			Content:  string(utf8Content),
		}, true, nil

	case ".java":
		utf8Content, err := utils.ConvertToUtf8(content)
		if err != nil {
			return "", nil, true, err
		}
		sqls, err := getSqlFromJavaContent(string(utf8Content))
		if err != nil {
			return "", nil, true, err
		}
		if len(sqls) == 0 {
			return "", nil, true, nil
		}
		return strings.Join(sqls, ";\n"), nil, true, nil

	default:
		return "", nil, false, nil
	}
}

// getSqlFromJavaContent 将 Java 源码内容写入临时文件，然后调用 javaParser.GetSqlFromJavaFile 提取 SQL。
// javaParser.GetSqlFromJavaFile 依赖 antlr 的 NewFileStream，需要从文件路径读取，因此需要临时文件。
func getSqlFromJavaContent(javaContent string) ([]string, error) {
	tmpFile, err := os.CreateTemp("", "sqle-java-*.java")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(javaContent); err != nil {
		return nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	return javaParser.GetSqlFromJavaFile(tmpFile.Name())
}
