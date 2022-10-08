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

func setArthmeticFlags(state *state8080, result uint16) {
	state.cc.ac = Btoi(result > 0xF)
	state.cc.cy = Btoi(result > 0xFF)
	state.cc.z = Btoi((result & 0xFF) == 0) // checking the first 8 bits is zero not the entire result.
	state.cc.s = Btoi((result & 0x80) == 0x80)
	state.cc.p = Btoi(bits.OnesCount8(uint8(result&0xFF))%2 == 0)
}

func Btoi(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
