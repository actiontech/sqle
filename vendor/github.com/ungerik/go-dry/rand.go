package dry

import (
	cryptorand "crypto/rand"
	"fmt"
	mathrand "math/rand"
	"time"
)

// RandSeedWithTime calls rand.Seed() with the current time.
func RandSeedWithTime() {
	mathrand.Seed(time.Now().UTC().UnixNano())
}

func getRandomHexString(length int, formatStr string) string {
	var buffer []byte
	if length%2 == 0 {
		buffer = make([]byte, length/2)
	} else {
		buffer = make([]byte, (length+1)/2)
	}
	_, err := cryptorand.Read(buffer)
	if err != nil {
		return ""
	}
	hexString := fmt.Sprintf(formatStr, buffer)
	return hexString[:length]
}

// RandomHexString returns a random lower case hex string with length.
func RandomHexString(length int) string {
	return getRandomHexString(length, "%x")
}

// RandomHEXString returns a random upper case hex string with length.
func RandomHEXString(length int) string {
	return getRandomHexString(length, "%X")
}
