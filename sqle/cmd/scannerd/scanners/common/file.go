package common

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	mybatisParser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

func GetSQLFromPath(pathName string, skipErrorQuery, skipErrorFile bool, fileSuffix string) (allSQL []driverV2.Node, err error) {
	if !path.IsAbs(pathName) {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		pathName = path.Join(pwd, pathName)
	}

	fileInfos, err := os.ReadDir(pathName)
	if err != nil {
		return nil, err
	}
	for _, fi := range fileInfos {
		var sqlList []driverV2.Node
		pathJoin := path.Join(pathName, fi.Name())

		if fi.IsDir() {
			sqlList, err = GetSQLFromPath(pathJoin, skipErrorQuery, skipErrorFile, fileSuffix)
		} else if strings.HasSuffix(fi.Name(), fileSuffix) {
			sqlList, err = GetSQLFromFile(pathJoin, skipErrorQuery, fileSuffix)
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

func GetSQLFromFile(file string, skipErrorQuery bool, fileSuffix string) (r []driverV2.Node, err error) {
	content, err := ReadFileContent(file)
	if err != nil {
		return nil, err
	}
	switch fileSuffix {
	case scanners.MybatisFileSuffix:
		sqls, err := mybatisParser.ParseXMLQuery(content, skipErrorQuery)
		if err != nil {
			return nil, err
		}
		for _, sql := range sqls {
			n, err := Parse(context.TODO(), sql)
			if err != nil {
				return nil, err
			}
			r = append(r, n...)
		}
	case scanners.SQLFileSuffix:
		n, err := Parse(context.TODO(), content)
		if err != nil {
			return nil, err
		}
		r = append(r, n...)
	}
	return r, nil
}

func ReadFileContent(file string) (content string, err error) {
	data, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return "", err
	}
	return string(data), err
}
