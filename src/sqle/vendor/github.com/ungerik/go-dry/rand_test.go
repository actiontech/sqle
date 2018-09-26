package dry

import (
	"testing"
)

func randomHexStringTestHelper(
	t *testing.T,
	testMethod func(int) string,
	upperCase bool) {
	testFn := func(n int) string {
		random_bytes := testMethod(n)
		if len(random_bytes) != n {
			t.FailNow()
		}
		for i := range []byte(random_bytes) {
			if i > 9 {
				if upperCase {
					if !(i >= 'A' && i <= 'Z') {
						t.FailNow()
					}
				} else {
					if !(i >= 'a' && i <= 'z') {
						t.FailNow()
					}
				}
			}
		}
		return random_bytes
	}
	random_bytes1 := testFn(3)
	testFn(10)
	random_bytes2 := testFn(3)
	if random_bytes1 == random_bytes2 {
		t.FailNow()
	}
}

func TestRandomHexString(t *testing.T) {
	randomHexStringTestHelper(t, RandomHexString, false)
}

func TestRandomHEXString(t *testing.T) {
	randomHexStringTestHelper(t, RandomHEXString, true)
}
