package emulator

import (
	"fmt"
	"os"
)

type restart uint8

const (
	RST0 restart = iota
	RST1
	RST2
	RST3
	RST4
	RST5
	RST6
	RST7
)

var cycles8080 = []int{
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, //0x00..0x0f
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, //0x10..0x1f
	4, 10, 16, 5, 5, 5, 7, 4, 4, 10, 16, 5, 5, 5, 7, 4,
	4, 10, 13, 5, 10, 10, 10, 4, 4, 10, 13, 5, 5, 5, 7, 4,
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, //0x40..0x4f
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5,
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5,
	7, 7, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 5, 7, 5,
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, //0x80..8x4f
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4,
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4,
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4,
	11, 10, 10, 10, 17, 11, 7, 11, 11, 10, 10, 10, 10, 17, 7, 11, //0xc0..0xcf
	11, 10, 10, 10, 17, 11, 7, 11, 11, 10, 10, 10, 10, 17, 7, 11,
	11, 10, 10, 18, 17, 11, 7, 11, 11, 5, 10, 5, 17, 17, 7, 11,
	11, 10, 10, 4, 17, 11, 7, 11, 11, 5, 10, 4, 17, 17, 7, 11,
}

type conditionCodes struct {
	z   uint8
	s   uint8
	p   uint8
	cy  uint8
	ac  uint8
	pad uint8
}

type state8080 struct {
	a          uint8
	b          uint8
	c          uint8
	d          uint8
	e          uint8
	h          uint8
	l          uint8
	sp         uint16
	pc         uint16
	memory     []uint8
	cc         conditionCodes
	int_enable uint8
}

func newState8080() *state8080 {
	cc := conditionCodes{
		z:   0,
		s:   0,
		p:   0,
		cy:  0,
		ac:  0,
		pad: 0,
	}
	return &state8080{
		a:          0,
		b:          0,
		c:          0,
		d:          0,
		e:          0,
		h:          0,
		l:          0,
		sp:         0,
		pc:         0,
		memory:     make([]uint8, 0, 0x10000), // 16K
		cc:         cc,
		int_enable: 0,
	}
}

func loadFileIntoMemoryAt(state *state8080, file string, offset int) {
	bs, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}
	//copy(state.memory[offset:], bs)

	// state.memory = slices.Insert(state.memory, offset, bs...)
	state.memory = append(state.memory[:offset], append(bs, state.memory[offset:]...)...)
}
func LoadSpaceInvaders(state *state8080) {
	loadFileIntoMemoryAt(state, "rom/invaders/invaders.h", 0x0)    // 0000-07FF
	loadFileIntoMemoryAt(state, "rom/invaders/invaders.g", 0x800)  // 0800-0FFF
	loadFileIntoMemoryAt(state, "rom/invaders/invaders.f", 0x1000) // 1000-17FF
	loadFileIntoMemoryAt(state, "rom/invaders/invaders.e", 0x1800) // 1800-1FFF
	for i := 0; i < 8193; i++ {
		state.memory = append(state.memory, 0xFF)
	}

}

func ScreenData(state *state8080) []uint8 {
	encoded := state.memory[0x2400:0x3FFF]
	decoded := make([]uint8, len(encoded)*8, len(encoded)*8)
	i := 0
	for _, e := range encoded {

		decoded[i] = Btoi(0x80 == (e & 0x80))   // bit 7 0x80
		decoded[i+1] = Btoi(0x40 == (e & 0x40)) // bit 6 0x40
		decoded[i+2] = Btoi(0x20 == (e & 0x20)) // bit 5 0x20
		decoded[i+3] = Btoi(0x10 == (e & 0x10)) // bit 4 0x10
		decoded[i+4] = Btoi(0x08 == (e & 0x08)) // bit 3 0x08
		decoded[i+5] = Btoi(0x04 == (e & 0x04)) // bit 2 0x04
		decoded[i+6] = Btoi(0x02 == (e & 0x02)) // bit 1 0x02
		decoded[i+7] = Btoi(0x01 == (e & 0x01)) // bit 0 0x01
		i += 8
	}

	return decoded
}

func Initiate8080() *state8080 {
	state := newState8080()
	// May need to add more to this later
	return state
}

func Emulate8080(state *state8080) int {
	opCode := state.memory[state.pc]
	state.pc++

	switch opCode {
	case 0x00:
	case 0x01: // LXI B
		state.c = state.memory[state.pc]
		state.pc++
		state.b = state.memory[state.pc]
		state.pc++
	case 0x02: // STAX B
		bc := bytesToPair(state.b, state.c)
		state.memory[bc] = state.a
	case 0x03: // INX B
		bc := bytesToPair(state.b, state.c)
		bc += 1
		state.b, state.c = pairToBytes(bc)
	case 0x04: // INR B
		state.b = add8(state, state.b, 1)
	case 0x05: // DCR B
		state.b = sub8(state, state.b, 1)
	case 0x06: // MVI B
		state.b = state.memory[state.pc]
		state.pc++
	case 0x07: // RLC
		r, cb := shiftLeft8(state.a, 1)
		state.a = uint8(r + cb)
		state.cc.cy = cb
	case 0x08: // -
	case 0x09: // DAD B (BC+HL) -> HL
		state.cc.cy = 0
		bc := bytesToPair(state.b, state.c)
		hl := bytesToPair(state.h, state.l)
		r, ok := add16(bc, hl)
		if !ok {
			state.cc.cy = 1
		}
		state.h, state.l = pairToBytes(r)
	case 0x0a: // LDAX B -> A
		bc := bytesToPair(state.b, state.c)
		state.a = state.memory[bc]
	case 0x0b: // DCX B
		bc := bytesToPair(state.b, state.c)
		state.b, state.c = pairToBytes(bc - 1)
	case 0x0c: // INR C
		state.c = add8(state, state.c, 1)
	case 0x0d: // DCR C
		state.c = sub8(state, state.c, 1)
	case 0x0e: // MVI C
		state.c = state.memory[state.pc]
		state.pc++
	case 0x0f: // RRC
		r, cy := shiftRight8(state.a, 1)
		state.a = r
		state.cc.cy = cy
	case 0x10: // -
	case 0x11: // LXI D
		state.e = state.memory[state.pc]
		state.pc++
		state.d = state.memory[state.pc]
		state.pc++
	case 0x12: // STAX D
		de := bytesToPair(state.d, state.e)
		state.memory[de] = state.a
	case 0x13: // INX D
		de := bytesToPair(state.d, state.e)
		de += 1
		state.d, state.e = pairToBytes(de)
	case 0x14: // INR D
		state.d = add8(state, state.d, 1)
	case 0x15: // DCR D
		state.d = sub8(state, state.d, 1)
	case 0x16: // MVI D
		state.d = state.memory[state.pc]
		state.pc++
	case 0x17: // "RAL", 1
		r, cb := shiftLeft8(state.a, 1)
		state.a = uint8(r + state.cc.cy)
		state.cc.cy = cb
	case 0x18: // -
	case 0x19: // DAD D (DE+HL) -> HL
		state.cc.cy = 0
		de := bytesToPair(state.d, state.e)
		hl := bytesToPair(state.h, state.l)
		r, ok := add16(de, hl)
		if !ok {
			state.cc.cy = 1
		}
		state.h, state.l = pairToBytes(r)
	case 0x1a: // LDAX D
		de := bytesToPair(state.d, state.e)
		state.a = state.memory[de]
	case 0x1b: // DCX D
		de := bytesToPair(state.d, state.e)
		state.d, state.e = pairToBytes(de - 1)
	case 0x1c: // INR E
		state.e = add8(state, state.e, 1)
	case 0x1d:
		state.e = sub8(state, state.e, 1)
	case 0x1e: // MVI E
		state.e = state.memory[state.pc]
		state.pc++
	case 0x1f: // RAR
		r, lsb := shiftRight8(state.a, 1)
		state.a = uint8(r + (state.cc.cy << 7))
		state.cc.cy = lsb
	case 0x20: // RIM - 8080 NOP
	case 0x21: // LXI H
		state.l = state.memory[state.pc]
		state.pc++
		state.h = state.memory[state.pc]
		state.pc++
	case 0x22: // SHLD
		state.memory[state.pc] = state.l
		state.pc++
		state.memory[state.pc] = state.h
		state.pc++
	case 0x23: // INX H
		hl := bytesToPair(state.h, state.l)
		hl += 1
		state.h, state.l = pairToBytes(hl)
	case 0x24: // INR H
		state.h = add8(state, state.h, 1)
	case 0x25: // DCR H
		state.h = sub8(state, state.h, 1)
	case 0x26: // MVI H
		state.h = state.memory[state.pc]
		state.pc++
	case 0x27: // DAA
		lo := state.a & 0xF
		if lo > 0x09 || state.cc.ac == 1 {
			state.a = add8(state, state.a, 6)
		}

		hi := state.a >> 4
		if hi > 9 || state.cc.cy == 1 {
			state.a = add8(state, state.a, 0x60)
		}
	case 0x28: // -
	case 0x29: // DAD H (HL+HL) -> HL
		state.cc.cy = 0
		hl := bytesToPair(state.h, state.l)
		r, ok := add16(hl, hl)
		if !ok {
			state.cc.cy = 1
		}
		state.h, state.l = pairToBytes(r)
	case 0x2a: // LHLD
		state.l = state.memory[state.pc]
		state.pc++
		state.h = state.memory[state.pc]
		state.pc++
	case 0x2b: // DCX H
		hl := bytesToPair(state.h, state.l)
		state.h, state.l = pairToBytes(hl - 1)
	case 0x2c: // INR L
		state.l = add8(state, state.l, 1)
	case 0x2d: //  DCR L
		state.l = sub8(state, state.l, 1)
	case 0x2e: // MVI L
		state.l = state.memory[state.pc]
		state.pc++
	case 0x2f: // CMA
		state.a = state.a ^ 0xFF
	case 0x30: // SIM - 8080 NOP
	case 0x31: // LXI SP
		p := state.memory[state.pc]
		state.pc++
		s := state.memory[state.pc]
		state.pc++
		state.sp = bytesToPair(s, p)
	case 0x32: // STA
		lo := state.memory[state.pc]
		state.pc++
		hi := state.memory[state.pc]
		state.pc++
		addr := bytesToPair(hi, lo)
		state.memory[addr] = state.a
	case 0x33: // INX SP
		state.sp++
	case 0x34: // INR M
		m := bytesToPair(state.h, state.l)
		state.memory[m] = add8(state, state.memory[m], 1)
	case 0x35: // DCR M
		m := bytesToPair(state.h, state.l)
		state.memory[m] = sub8(state, state.memory[m], 1)
	case 0x36: // MVI M
		m := bytesToPair(state.h, state.l)
		state.memory[m] = state.memory[state.pc]
		state.pc++
	case 0x37: // STC
		state.cc.cy = 1
	case 0x38: // -
	case 0x39: // DAD SP
		state.cc.cy = 0
		hl := bytesToPair(state.h, state.l)
		sp := state.sp
		r, ok := add16(hl, sp)
		if !ok {
			state.cc.cy = 1
		}
		state.h, state.l = pairToBytes(r)
	case 0x3a: // LDA
		lo := state.memory[state.pc]
		state.pc++
		hi := state.memory[state.pc]
		state.pc++
		m := bytesToPair(hi, lo)
		state.a = state.memory[m]
	case 0x3b: // DCX SP
		state.sp--
	case 0x3c: // INR A
		state.a = add8(state, state.a, 1)
	case 0x3d: // DCR A
		state.a = sub8(state, state.a, 1)
	case 0x3e: // MVI A
		state.a = state.memory[state.pc]
		state.pc++
	case 0x3f: // CMC
		state.cc.cy--
	case 0x40: // MOV B,B
	case 0x41: // MOV B,C
		state.b = state.c
	case 0x42: // MOV B, D
		state.b = state.d
	case 0x43: // MOV B, E
		state.b = state.e
	case 0x44: // MOV B,H
		state.b = state.h
	case 0x45: // MOV B,L
		state.b = state.l
	case 0x46: // MOV B,M
		addr := bytesToPair(state.h, state.l)
		state.b = state.memory[addr]
	case 0x47: // MOV B,A
		state.b = state.a
	case 0x48: // MOV C,B
		state.c = state.b
	case 0x49: // MOV C,C
	case 0x4a: // MOV C,D
		state.c = state.d
	case 0x4b: // MOV C,E
		state.c = state.e
	case 0x4c: // MOV C,H
		state.c = state.h
	case 0x4d: // MOV C,L
		state.c = state.l
	case 0x4e: // MOV C,M
		addr := bytesToPair(state.h, state.l)
		state.c = state.memory[addr]
	case 0x4f: // MOV C,A
		state.c = state.a
	case 0x50: // MOV D,B
		state.d = state.b
	case 0x51: // MOV D,C
		state.d = state.c
	case 0x52: // MOV D,D
	case 0x53: // MOV D,E
		state.d = state.e
	case 0x54: // MOV D,H
		state.d = state.h
	case 0x55: // MOV D,L
		state.d = state.l
	case 0x56: // MOV D,M
		addr := bytesToPair(state.h, state.l)
		state.d = state.memory[addr]
	case 0x57: // MOV D,A
		state.d = state.a
	case 0x58: // MOV E,B
		state.e = state.b
	case 0x59: // MOV E,C
		state.e = state.c
	case 0x5a: // MOV E,D
		state.e = state.d
	case 0x5b: // MOV E,E
	case 0x5c: // MOV E,H
		state.e = state.h
	case 0x5d: // MOV E,L
		state.e = state.l
	case 0x5e: // MOV E,M
		addr := bytesToPair(state.h, state.l)
		state.e = state.memory[addr]
	case 0x5f: // MOV E,A
		state.e = state.a
	case 0x60: // MOV H,B
		state.h = state.b
	case 0x61: // MOV H,C
		state.h = state.c
	case 0x62: // MOV H,D
		state.h = state.d
	case 0x63: // MOV H, E
		state.h = state.e
	case 0x64: // MOV H, H
	case 0x65: // MOV H, L
		state.h = state.l
	case 0x66: // MOV H, M
		addr := bytesToPair(state.h, state.l)
		state.h = state.memory[addr]
	case 0x67: // MOV H, A
		state.h = state.a
	case 0x68: // MOV L, B
		state.l = state.b
	case 0x69: // MOV L, C
		state.l = state.c
	case 0x6a: // MOV L, D
		state.l = state.d
	case 0x6b: // MOV L, E
		state.l = state.e
	case 0x6c: // MOV L, H
		state.l = state.h
	case 0x6d: // MOV L, L
	case 0x6e: // MOV L, M
		addr := bytesToPair(state.h, state.l)
		state.l = state.memory[addr]
	case 0x6f: // MOV L,A
		state.l = state.a
	case 0x70: // MOV M,B
		addr := bytesToPair(state.h, state.l)
		state.memory[addr] = state.b
	case 0x71: // MOV M,C
		addr := bytesToPair(state.h, state.l)
		state.memory[addr] = state.c
	case 0x72: // MOV M, D
		addr := bytesToPair(state.h, state.l)
		state.memory[addr] = state.d
	case 0x73: // MOV M, E
		addr := bytesToPair(state.h, state.l)
		state.memory[addr] = state.e
	case 0x74: // MOV M, H
		addr := bytesToPair(state.h, state.l)
		state.memory[addr] = state.h
	case 0x75: // MOV M, L
		addr := bytesToPair(state.h, state.l)
		state.memory[addr] = state.l
	case 0x76: // HLT
		os.Exit(0x76)
	case 0x77: // MOV M, A
		addr := bytesToPair(state.h, state.l)
		state.memory[addr] = state.a
	case 0x78: // MOV A,B
		state.a = state.b
	case 0x79: // MOV A,C
		state.a = state.c
	case 0x7a: // MOV A, D
		state.a = state.d
	case 0x7b: // MOV A, E
		state.a = state.e
	case 0x7c: // MOV A, H
		state.a = state.h
	case 0x7d: // MOV A, L
		state.a = state.l
	case 0x7e: // MOV A, M
		addr := bytesToPair(state.h, state.l)
		state.a = state.memory[addr]
	case 0x7f: // MOV A, A
	case 0x80: // ADD B
		state.a = add8(state, state.a, state.b)
	case 0x81: // ADD C
		state.a = add8(state, state.a, state.c)
	case 0x82: // ADD D
		state.a = add8(state, state.a, state.d)
	case 0x83: // ADD E
		state.a = add8(state, state.a, state.e)
	case 0x84: // ADD H
		state.a = add8(state, state.a, state.h)
	case 0x85: // ADD L
		state.a = add8(state, state.a, state.l)
	case 0x86: // ADD M
		addr := bytesToPair(state.h, state.l)
		m := state.memory[addr]
		state.a = add8(state, state.a, m)
	case 0x87: // ADD A
		state.a = add8(state, state.a, state.a)
	case 0x88: // ADC B
		state.a = add8WithCarry(state, state.a, state.b)
	case 0x89: // ADC C
		state.a = add8WithCarry(state, state.a, state.c)
	case 0x8a: // ADC D
		state.a = add8WithCarry(state, state.a, state.d)
	case 0x8b: // ADC E
		state.a = add8WithCarry(state, state.a, state.e)
	case 0x8c: // ADC H
		state.a = add8WithCarry(state, state.a, state.h)
	case 0x8d: // ADC L
		state.a = add8WithCarry(state, state.a, state.l)
	case 0x8e: // ADC M
		addr := bytesToPair(state.h, state.l)
		m := state.memory[addr]
		state.a = add8WithCarry(state, state.a, m)
	case 0x8f: // ADC A
		state.a = add8WithCarry(state, state.a, state.a)
	case 0x90: // SUB B
		state.a = sub8(state, state.a, state.b)
	case 0x91: // SUB C
		state.a = sub8(state, state.a, state.c)
	case 0x92: // SUB D
		state.a = sub8(state, state.a, state.d)
	case 0x93: // SUB E
		state.a = sub8(state, state.a, state.e)
	case 0x94: // SUB H
		state.a = sub8(state, state.a, state.h)
	case 0x95: // SUB L
		state.a = sub8(state, state.a, state.l)
	case 0x96: // SUB H
		addr := bytesToPair(state.h, state.l)
		m := state.memory[addr]
		state.a = sub8(state, state.a, m)
	case 0x97: // SUB M
		state.a = sub8(state, state.a, state.a)
	case 0x98: // SBB B
		state.a = sub8WithBorrow(state, state.a, state.b)
	case 0x99: // SBB C
		state.a = sub8WithBorrow(state, state.a, state.c)
	case 0x9a: // SBB D
		state.a = sub8WithBorrow(state, state.a, state.d)
	case 0x9b: // SBB E
		state.a = sub8WithBorrow(state, state.a, state.e)
	case 0x9c: // SBB H
		state.a = sub8WithBorrow(state, state.a, state.h)
	case 0x9d: // SBB L
		state.a = sub8WithBorrow(state, state.a, state.l)
	case 0x9e: // SBB M
		addr := bytesToPair(state.h, state.l)
		m := state.memory[addr]
		state.a = sub8WithBorrow(state, state.a, m)
	case 0x9f: // SBB A
		state.a = sub8WithBorrow(state, state.a, state.a)
	case 0xa0: // ANA B
		state.a = state.a & state.b
		setLogicFlags(state)
	case 0xa1: // ANA C
		state.a = state.a & state.c
		setLogicFlags(state)
	case 0xa2: // ANA D
		state.a = state.a & state.d
		setLogicFlags(state)
	case 0xa3: // ANA E
		state.a = state.a & state.e
		setLogicFlags(state)
	case 0xa4: // ANA H
		state.a = state.a & state.h
		setLogicFlags(state)
	case 0xa5: // ANA L
		state.a = state.a & state.l
		setLogicFlags(state)
	case 0xa6: // ANA M
		addr := bytesToPair(state.h, state.l)
		m := state.memory[addr]
		state.a = state.a & m
		setLogicFlags(state)
	case 0xa7: // ANA A
		state.a = state.a & state.a
		setLogicFlags(state)
	case 0xa8: // XRA B
		state.a = state.a ^ state.b
		setLogicFlags(state)
	case 0xa9: // XRA C
		state.a = state.a ^ state.c
		setLogicFlags(state)
	case 0xaa: // XRA D
		state.a = state.a ^ state.d
		setLogicFlags(state)
	case 0xab: // XRA E
		state.a = state.a ^ state.e
		setLogicFlags(state)
	case 0xac: // XRA H
		state.a = state.a ^ state.h
		setLogicFlags(state)
	case 0xad: // XRA L
		state.a = state.a ^ state.l
		setLogicFlags(state)
	case 0xae: // XRA M
		addr := bytesToPair(state.h, state.l)
		m := state.memory[addr]
		state.a = state.a ^ m
		setLogicFlags(state)
	case 0xaf: // XRA A
		state.a = state.a ^ state.a
		setLogicFlags(state)
	case 0xb0: // ORA B
		state.a = state.a | state.b
		setLogicFlags(state)
	case 0xb1: // ORA C
		state.a = state.a | state.c
		setLogicFlags(state)
	case 0xb2: // ORA D
		state.a = state.a | state.d
		setLogicFlags(state)
	case 0xb3: // ORA E
		state.a = state.a | state.e
		setLogicFlags(state)
	case 0xb4: // ORA H
		state.a = state.a | state.h
		setLogicFlags(state)
	case 0xb5: // ORA L
		state.a = state.a | state.l
		setLogicFlags(state)
	case 0xb6: // ORA M
		addr := bytesToPair(state.h, state.l)
		m := state.memory[addr]
		state.a = state.a | m
	case 0xb7: // ORA A
		state.a = state.a | state.a
		setLogicFlags(state)
	case 0xb8: // CMP B
		// if rp > a then zero reset & carry set
		// if rp < a then zero reset & carry reset
		// if rp == a then zero is set & carry reset
		sub8(state, state.a, state.b)
	case 0xb9: // CMP C
		sub8(state, state.a, state.c)
	case 0xba: // CMP D
		sub8(state, state.a, state.d)
	case 0xbb: // CMP E
		sub8(state, state.a, state.e)
	case 0xbc: // CMP H
		sub8(state, state.a, state.h)
	case 0xbd: // CMP L
		sub8(state, state.a, state.l)
	case 0xbe: // CMP M
		addr := bytesToPair(state.h, state.l)
		m := state.memory[addr]
		sub8(state, state.a, m)
	case 0xbf: // CMP A
		sub8(state, state.a, state.a)
	case 0xc0: // RNZ
		if state.cc.z == 0 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xc1: // POP B
		state.c = state.memory[state.sp]
		state.b = state.memory[state.sp+1]
		state.sp += 2
	case 0xc2: // JNZ
		if state.cc.z == 0 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xc3: // JMP
		state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
	case 0xc4: // CNZ
		if state.cc.z == 1 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+1]
			lo := state.memory[state.pc]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xc5: // PUSH B
		state.memory[state.sp-1] = state.b
		state.memory[state.sp-2] = state.c
		state.sp -= 2
	case 0xc6: // ADI
		state.a = add8(state, state.a, state.memory[state.pc])
		state.pc++
	case 0xc7: // RST 0
		// push pc to stack and jump to 0x0000
		state.int_enable = 0
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
		state.sp -= 2
		state.pc = bytesToPair(0x00, 0x00)
	case 0xc8: // RZ
		if state.cc.z == 1 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xc9: // RET
		state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
		state.sp += 2
	case 0xca: // JZ
		if state.cc.z == 1 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xcb: // -
	case 0xcc: // CZ
		if state.cc.z == 0 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+1]
			lo := state.memory[state.pc]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xcd: // CALL
		ret := state.pc + 2
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
		state.sp -= 2
		hi := state.memory[state.pc+1]
		lo := state.memory[state.pc]
		state.pc = bytesToPair(hi, lo)
	case 0xce: // ACI
		state.a = add8WithCarry(state, state.a, state.memory[state.pc])
		state.pc++
	case 0xcf: // RST 1
		// push pc to stack and jump to 0x0008
		state.int_enable = 0
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
		state.sp -= 2
		state.pc = bytesToPair(0x00, 0x08)
	case 0xd0: // RNC
		if state.cc.cy == 0 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xd1: // POP D
		state.e = state.memory[state.sp]
		state.d = state.memory[state.sp+1]
		state.sp += 2
	case 0xd2: // JNC
		if state.cc.cy == 0 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xd3: // OUT
		//fmt.Printf("OUT port: %x  a: %b\n", state.memory[state.pc], state.a)
		state.pc++
	case 0xd4: // CNC
		if state.cc.cy == 0 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+1]
			lo := state.memory[state.pc]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xd5: // PUSH D
		state.memory[state.sp-1] = state.d
		state.memory[state.sp-2] = state.e
		state.sp -= 2
	case 0xd6: // SUI
		state.a = sub8(state, state.a, state.memory[state.pc])
		state.pc++
	case 0xd7: // RST 2
		// push pc to stack and jump to 0x0010
		state.int_enable = 0
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
		state.sp -= 2
		state.pc = bytesToPair(0x00, 0x10)
	case 0xd8: // RC
		if state.cc.cy == 1 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xd9: // -
	case 0xda: // JC
		if state.cc.cy == 1 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xdb: // IN
		//	fmt.Printf("IN port: %x  a: %b\n", state.memory[state.pc], state.a)
		state.pc++
	case 0xdc: // CC
		if state.cc.cy == 1 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+1]
			lo := state.memory[state.pc]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xdd: // -
	case 0xde: // SBI
		state.a = sub8WithBorrow(state, state.a, state.memory[state.pc])
		state.pc++
	case 0xdf: // RST 3
		// push pc to stack and jump to 0x0018
		state.int_enable = 0
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
		state.sp -= 2
		state.pc = bytesToPair(0x00, 0x18)
	case 0xe0: // RPO
		if state.cc.p == 0 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xe1: // POP H
		state.l = state.memory[state.sp]
		state.h = state.memory[state.sp+1]
		state.sp += 2
	case 0xe2: // JPO
		if state.cc.p == 0 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xe3: // XTHL
		hl := bytesToPair(state.h, state.l)
		state.h, state.l = pairToBytes(state.sp)
		state.sp = hl
	case 0xe4: // CPO
		if state.cc.p == 0 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+1]
			lo := state.memory[state.pc]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xe5: // PUSH H
		state.memory[state.sp-1] = state.h
		state.memory[state.sp-2] = state.l
		state.sp -= 2
	case 0xe6: // ANI
		state.a = state.a & state.memory[state.pc]
		setLogicFlags(state)
		state.pc++
	case 0xe7: // RST 4
		// push pc to stack and jump to 0x0020
		state.int_enable = 0
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
		state.sp -= 2
		state.pc = bytesToPair(0x00, 0x20)
	case 0xe8: // RPE
		if state.cc.p == 1 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xe9: // PCHL
		state.pc = bytesToPair(state.h, state.l)
	case 0xea: // JPE
		if state.cc.p == 1 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xeb: // XCHG swap HL with DE
		h, l := state.h, state.l
		state.h, state.l = state.d, state.e
		state.d, state.e = h, l
	case 0xec: // CPE
		if state.cc.p == 1 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+1]
			lo := state.memory[state.pc]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xed: // -
	case 0xee: // XRI
		state.a = state.a ^ state.memory[state.pc]
		setLogicFlags(state)
		state.pc++
	case 0xef: // RST 5
		// push pc to stack and jump to 0x0028
		state.int_enable = 0
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
		state.sp -= 2
		state.pc = bytesToPair(0x00, 0x28)
	case 0xf0: // RP
		if state.cc.s == 0 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xf1: // POP PSW
		state.a = state.memory[state.sp]
		psw := state.memory[state.sp+1]
		setFlagsFromPSW(state, psw)
		state.sp += 2
	case 0xf2: // JP
		if state.cc.s == 0 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xf3: // DI - 01 Disable Interrupts
		state.int_enable = 0
	case 0xf4: // CP
		if state.cc.s == 0 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+1]
			lo := state.memory[state.pc]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xf5: // PUSH PSW
		psw := setPSW(state)
		state.memory[state.sp-1] = state.a
		state.memory[state.sp-2] = psw
		state.sp -= 2
	case 0xf6: // ORI
		state.a = state.a | state.memory[state.pc]
		setLogicFlags(state)
		state.pc++
	case 0xf7: // RST 6
		// push pc to stack and jump to 0x0030
		state.int_enable = 0
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
		state.sp -= 2
		state.pc = bytesToPair(0x00, 0x30)
	case 0xf8: // RM
		if state.cc.s == 1 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xf9: // SPHL
		state.sp = bytesToPair(state.h, state.l)
	case 0xfa: // JP
		if state.cc.s == 1 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xfb: // EI - EI Enable Interrupts
		state.int_enable = 1
	case 0xfc: // CM
		if state.cc.s == 1 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+1]
			lo := state.memory[state.pc]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xfd: // -
	case 0xfe: // CPI
		sub8(state, state.a, state.memory[state.pc])
		state.pc++
	case 0xff: // RST 7
		// push pc to stack and jump to 0x0038
		state.int_enable = 0
		state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
		state.sp -= 2
		state.pc = bytesToPair(0x00, 0x38)
	}

	return cycles8080[opCode]
}

func Restart8080(state *state8080, oper restart) {
	var addr uint8
	switch oper {
	case RST0:
		addr = 0x00
	case RST1:
		addr = 0x08
	case RST2:
		addr = 0x10
	case RST3:
		addr = 0x18
	case RST4:
		addr = 0x20
	case RST5:
		addr = 0x28
	case RST6:
		addr = 0x30
	case RST7:
		addr = 0x38

	}

	// push pc to stack and jump
	state.int_enable = 0
	state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(state.pc)
	state.sp -= 2
	state.pc = bytesToPair(0x00, addr)
}

func GetIntEnabled(state *state8080) uint8 {
	return state.int_enable
}

func GetInput(state *state8080) {
	// INPUT
	if state.memory[state.pc] == 0xdb {
		switch state.memory[state.pc+1] {
		case 0x0: // fire, left, right?
		case 0x1: // credit,start, player 1 shot, left, right
		case 0x2: // dip3,5,6, player 2 shot, left, right
		case 0x3: // shift reg data
			m := uint16(shiftMSB) << 8
			shift := uint16(m | uint16(shiftLSB))
			state.a = uint8((shift >> (8 - shiftCount)) & 0xFF)
		}
	}
}

func GetOutput(state *state8080) {
	// OUTPUT
	if state.memory[state.pc] == 0xd3 {
		switch state.memory[state.pc+1] {
		case 0x02: // shift amount
			shiftCount = state.a & 7
		case 0x03: // discrete sounds
		case 0x04: // shift data (LSB on 1st write, MSB on 2nd)
			shiftLSB = shiftMSB
			shiftMSB = state.a
		case 0x05: // discrete sounds
		case 0x06: // watchdog?
		}

	}
}
