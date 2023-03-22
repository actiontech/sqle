package v1

import (
	e "errors"
	"io/fs"
	"io/ioutil"
	"strings"
)

const (
	// LogoUrlBase sqle static 服务接口的url前缀
	LogoUrlBase = "/static/media"

	// LogoDir sqle logo 的本地目录
	LogoDir = "./ui/static/media"
)

func getLogoFileInfo() (fs.FileInfo, error) {
	fileInfos, err := ioutil.ReadDir(LogoDir)
	if err != nil {
		return nil, e.New("read logo dir failed")
	}

	var hasLogoFile bool
	var logoFileInfo fs.FileInfo
	for _, fileInfo := range fileInfos {
		if strings.HasPrefix(fileInfo.Name(), "logo.") {
			hasLogoFile = true
			logoFileInfo = fileInfo
			break
		}
	}
	if !hasLogoFile {
		return nil, e.New("no logo file")
	}

	return logoFileInfo, nil
}
