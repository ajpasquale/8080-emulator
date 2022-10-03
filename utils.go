package emulator

import "math/bits"

func shiftLeft8(a uint8, k int) (uint8, uint8) {
	msb := a >> 7
	r := uint16(a) << k
	return uint8(r), msb
}

func shiftRight8(a uint8, k int) (uint8, uint8) {
	lsb := a & 0x1
	r := a >> k
	return r, lsb

}

func add16(a uint16, b uint16) (uint16, bool) {
	const n = 16
	r := uint32(a + b)
	overflow := (bits.Len32(r) > n)
	return uint16(r), overflow
}
