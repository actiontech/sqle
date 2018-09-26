package dry

import "testing"

func TestEndianSafeSplitUint16(t *testing.T) {

	least, most := EndianSafeSplitUint16(1)
	if !(least == 1 && most == 0) {
		t.Fail()
	}

	least, most = EndianSafeSplitUint16(256)
	if !(least == 0 && most == 1) {
		t.Fail()
	}

}
