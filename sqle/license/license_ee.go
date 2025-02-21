//go:build !release
// +build !release

package license

import "fmt"


func CheckKnowledgeBaseLicense(license string) error {
	return fmt.Errorf("knowledge base license is not supported in this edition")
}
