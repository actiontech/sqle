package scsi

import (
	"regexp"
	"testing"
)

func TestRsvKeyMatch(t *testing.T) {
	tcs := make(map[string]string)
	tcs[`  PR generation=0x1, Reservation follows:
    Key=0x2f1f
    scope: LU_SCOPE,  type: Exclusive Access, registrants only`] = "0x2f1f"

	tcs[`  PR generation=0x1, Reservation follows:
    Key=0x111
    scope: LU_SCOPE,  type: Exclusive Access, registrants only`] = "0x111"

	tcs[`  PR generation=0x1, Reservation follows:
    Key=1111
    scope: LU_SCOPE,  type: Exclusive Access, registrants only`] = ""

	tcs[`  PR generation=0x1, Reservation follows:
    Key=0x21fz
    scope: LU_SCOPE,  type: Exclusive Access, registrants only`] = ""
	reg := regexp.MustCompile("(?s)Key=(0x[0-9a-fA-F]+)\\s.*type:(.*)")
	for tc, expect := range tcs {
		matches := reg.FindStringSubmatch(tc)
		if len(matches) < 3 {
			if "" != expect {
				t.Logf("\n%v \n not match	", tc)
				t.Fail()
			}
		} else {
			if expect != matches[1] {
				t.Logf("%v \n key is :(%v) ,result is (%v)", tc, expect, matches[1])
				t.Fail()
			}
		}

	}

}
