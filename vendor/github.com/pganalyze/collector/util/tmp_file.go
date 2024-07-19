package util

import (
	"os"
)

// Clean up a created temporary file, reporting any errors to the log
func CleanUpTmpFile(tmpFile *os.File, logger *Logger) {
	err := tmpFile.Close()
	if err != nil {
		logger.PrintError("Failed to close temporary file \"%s\": %s", tmpFile.Name(), err)
	}

	err = os.Remove(tmpFile.Name())
	if err != nil {
		logger.PrintError("Failed to delete temporary file \"%s\": %s", tmpFile.Name(), err)
	}
}
