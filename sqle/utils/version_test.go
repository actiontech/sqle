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
		}, {
			name:                "ob versionToBeCompared 1",
			version:             "4.2.1.1-101010012023111012",
			versionToBeCompared: "4.2.0",
			expectedResult:      false,
			expectedError:       false,
		}, {
			name:                "ob versionToBeCompared 2",
			version:             "4.2.1.2-101010012023111012",
			versionToBeCompared: "4.2.0.0-100001282023042317",
			expectedResult:      false,
			expectedError:       false,
		}, {
			name:                "ob versionToBeCompared 3",
			version:             "3.1.3-101010012023111012",
			versionToBeCompared: "4.10.0.0-100001282023042317",
			expectedResult:      true,
			expectedError:       false,
		},
		{
			name:                "ob versionToBeCompared 4",
			version:             "3.1.2-101010012023111012",
			versionToBeCompared: "4.10.0.0-100001282023042317",
			expectedResult:      true,
			expectedError:       false,
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
