package common

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	mybatisParser "github.com/actiontech/mybatis-mapper-2-sql"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/utils"
)

func GetSQLFromPath(pathName string, skipErrorQuery, skipErrorFile bool, fileSuffix string, showFileContent bool) (allSQL []driverV2.Node, err error) {
	if !path.IsAbs(pathName) {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		pathName = path.Join(pwd, pathName)
	}

	fileInfos, err := ioutil.ReadDir(pathName)
	if err != nil {
		return nil, err
	}
	for _, fi := range fileInfos {
		var sqlList []driverV2.Node
		pathJoin := path.Join(pathName, fi.Name())

		if fi.IsDir() {
			sqlList, err = GetSQLFromPath(pathJoin, skipErrorQuery, skipErrorFile, fileSuffix, showFileContent)
		} else if strings.HasSuffix(fi.Name(), fileSuffix) {
			sqlList, err = GetSQLFromFile(pathJoin, skipErrorQuery, fileSuffix, showFileContent)
		}

		if err != nil {
			if skipErrorFile {
				fmt.Printf("[parse %s file error] parse file %s error: %v\n", fileSuffix, pathJoin, err)
			} else {
				return nil, fmt.Errorf("parse file %s error: %v", pathJoin, err)
			}
		}
		allSQL = append(allSQL, sqlList...)
	}
	return allSQL, err
}

func GetSQLFromFile(file string, skipErrorQuery bool, fileSuffix string, showFileContent bool) (r []driverV2.Node, err error) {
	content, err := ReadFileContent(file)
	if err != nil {
		return nil, err
	}
	switch fileSuffix {
	case utils.MybatisFileSuffix:
		var sqls []string
		var err error
		if skipErrorQuery {
			sqls, err = mybatisParser.ParseXMLQuery(content, mybatisParser.SkipErrorQuery, mybatisParser.RestoreOriginSql)
		} else {
			sqls, err = mybatisParser.ParseXMLQuery(content, mybatisParser.RestoreOriginSql)
		}
		if err != nil {
			if showFileContent {
				fmt.Printf("failed to parse xml file content: %s", content)
			}
			return nil, err
		}
		for _, sql := range sqls {
			n, err := Parse(context.TODO(), sql)
			if err != nil {
				if showFileContent {
					fmt.Printf("failed to parse xml file content: %s", content)
				}
				return nil, err
			}
			r = append(r, n...)
		}
	case utils.SQLFileSuffix:
		n, err := Parse(context.TODO(), content)
		if err != nil {
			if showFileContent {
				fmt.Printf("failed to parse sql file content: %s", content)
			}
			return nil, err
		}
		r = append(r, n...)
	}
	return r, nil
}

func ReadFileContent(file string) (content string, err error) {
	data, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		return "", err
	}
	return string(data), err
}
