package utils

import (
	"testing"
)

func TestIsVersionLessThan(t *testing.T) {
	testCases := []struct {
		name                string
		version             string
		versionToBeCompared string
		expectedResult      bool
		expectedError       bool
	}{
		{
			name:                "valid versions, version is less than versionToBeCompared",
			version:             "1.2.3.1",
			versionToBeCompared: "1.2.4",
			expectedResult:      true,
			expectedError:       false,
		},
		{
			name:                "valid versions, version is greater than versionToBeCompared",
			version:             "1.2.4",
			versionToBeCompared: "1.2.3.1",
			expectedResult:      false,
			expectedError:       false,
		},
		{
			name:                "valid versions, version is equal to versionToBeCompared",
			version:             "1.2.3.2",
			versionToBeCompared: "1.2.3.2",
			expectedResult:      false,
			expectedError:       false,
		},
		{
			name:                "invalid version",
			version:             "invalid",
			versionToBeCompared: "1.2.3",
			expectedResult:      false,
			expectedError:       true,
		},
		{
			name:                "invalid versionToBeCompared",
			version:             "1.2.3",
			versionToBeCompared: "invalid",
			expectedResult:      false,
			expectedError:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := IsVersionLessThan(tc.version, tc.versionToBeCompared)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, but got: %v", tc.expectedError, err != nil)
			} else if result != tc.expectedResult {
				t.Errorf("expected result %v, but got %v", tc.expectedResult, result)
			}
		})
	}
}
