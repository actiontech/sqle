package dry

import (
	"unsafe"
)

func EndianIsLittle() bool {
	var word uint16 = 1
	littlePtr := (*uint8)(unsafe.Pointer(&word))
	return (*littlePtr) == 1
}

func EndianIsBig() bool {
	return !EndianIsLittle()
}

func EndianSafeSplitUint16(value uint16) (leastSignificant, mostSignificant uint8) {
	bytes := (*[2]uint8)(unsafe.Pointer(&value))
	if EndianIsLittle() {
		return bytes[0], bytes[1]
	}
	return bytes[1], bytes[0]
}
