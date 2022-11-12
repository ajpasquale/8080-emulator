package emulator

func pairToBytes(p uint16) (uint8, uint8) {
	const n = 8
	lo := uint8(p & 0x00FF)
	// need to shift down to 1 byte otherwise you lose the higher byte
	hi := uint8((p & 0xFF00) >> n)
	return hi, lo
}

func bytesToPair(hi uint8, lo uint8) uint16 {
	const n = 8
	return uint16(hi)<<n | uint16(lo)
}
