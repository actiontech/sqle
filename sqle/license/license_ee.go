//go:build !release
// +build !release

package license

func GetDMSLicense(content string) (*License, error) {
	return &License{}, nil
}

type License struct{}

func (l *License) CheckSupportKnowledgeBase() (bool, error) {
	return false, nil
}

func (l *License) GetKnowledgeBaseDBTypes() []string {
	return nil
}