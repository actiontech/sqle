package utils

import (
	"fmt"

	goVersion "github.com/hashicorp/go-version"
)
// Check if version is less than version to be compared, for example: input (3.1.1,3.1.2) returns true, input (3.1.1,3.1.1) returns false, input (3.1.1,3.1.0) returns false
func IsVersionLessThan(version, versionToBeCompared string) (bool, error) {
	version1, err := goVersion.NewVersion(version)
	if err != nil {
		return false, fmt.Errorf("input version is invalid:%w", err)
	}
	version2, err := goVersion.NewVersion(versionToBeCompared)
	if err != nil {
		return false, fmt.Errorf("input version to be compared is invalid:%w", err)
	}
	return version1.LessThan(version2), nil
}
