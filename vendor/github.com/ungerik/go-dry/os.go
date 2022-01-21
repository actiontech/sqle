package dry

import (
	"os"
	"strings"
)

// EnvironMap returns the current environment variables as a map.
func EnvironMap() map[string]string {
	return environToMap(os.Environ())
}

func environToMap(environ []string) map[string]string {
	ret := make(map[string]string)

	for _, v := range environ {
		parts := strings.SplitN(v, "=", 2)

		ret[parts[0]] = parts[1]
	}

	return ret
}

// GetenvDefault retrieves the value of the environment variable
// named by the key. It returns the given defaultValue if the
// variable is not present.
func GetenvDefault(key, defaultValue string) string {
	ret := os.Getenv(key)
	if ret == "" {
		return defaultValue
	}

	return ret
}
