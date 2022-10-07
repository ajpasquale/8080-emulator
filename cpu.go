package emulator

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

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
		memory:     make([]uint8, 0, 0x4000),
		cc:         cc,
		int_enable: 0,
	}
}

func parity(x int, size int) int {
	p := 0
	x = (x & ((1 << size) - 1))

	for i := 0; i < size; i++ {
		if x&0x1 == 1 {
			p++
		}
		x = x >> 1
	}
	if 0 == (p & 0x1) {
		return 1
	}
	return 0

}

func loadFileIntoMemoryAt(state *state8080, file string, offset int) {
	bs, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}
	state.memory = slices.Insert(state.memory, offset, bs...)
}

func Emulate8080(state *state8080) {
	opCode := state.memory[state.pc]

	state.pc++

	switch opCode {
	case 0x0:
		break
	case 0x01: // LXI B
		state.pc++
		state.c = state.memory[state.pc]
		state.pc++
		state.b = state.memory[state.pc]
	case 0x02: // STAX B
		bc := bytesToPair(state.b, state.c)
		state.memory[bc] = state.a
	case 0x03: // INX B
		bc := bytesToPair(state.b, state.c)
		bc += 1
		state.b, state.c = pairToBytes(bc)
	case 0x04: // INR B
		state.b++
	case 0x05: // DCR B
		state.b--
	case 0x06: // MVI B
		state.pc++
		state.b = state.memory[state.pc]
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
		state.c++
	case 0x0d: // DCR C
		state.c--
	case 0x0e: // MVI C
		state.pc++
		state.c = state.memory[state.pc]
	case 0x0f: // RRC
		r, cy := shiftRight8(state.a, 1)
		state.a = r
		state.cc.cy = cy
	case 0x10: // -
	case 0x11: // LXI D
		state.pc++
		state.e = state.memory[state.pc]
		state.pc++
		state.d = state.memory[state.pc]
	case 0x12: // STAX D
		de := bytesToPair(state.d, state.e)
		state.memory[de] = state.a
	case 0x13: // INX D
		de := bytesToPair(state.d, state.e)
		de += 1
		state.d, state.e = pairToBytes(de)
	case 0x14: // INR D
		state.d++
	case 0x15: // DCR D
		state.d--
	case 0x16: // "MVI D,", 2
		state.pc++
		state.d = state.memory[state.pc]
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
	case 0x1a: //LDAX D
		de := bytesToPair(state.d, state.e)
		state.a = state.memory[de]
	case 0x1b: // DCX D
		de := bytesToPair(state.d, state.e)
		state.d, state.e = pairToBytes(de - 1)
	case 0x1c: // INR E
		state.e++
	case 0x1d:
		state.e--
	case 0x1e:
		state.pc++
		state.e = state.memory[state.pc]
	case 0x1f: // RAR
		r, lsb := shiftRight8(state.a, 1)
		state.a = uint8(r + (state.cc.cy << 7))
		state.cc.cy = lsb
	case 0x20: // RIM ??
	case 0x21: // LXI H
		state.pc++
		state.l = state.memory[state.pc]
		state.pc++
		state.h = state.memory[state.pc]
	case 0x22: // SHLD
		state.pc++
		state.memory[state.pc] = state.l
		state.pc++
		state.memory[state.pc] = state.h
	case 0x23: // INX H
		hl := bytesToPair(state.h, state.l)
		hl += 1
		state.h, state.l = pairToBytes(hl)
	case 0x24: // INR H
		state.h++
	case 0x25:
		state.h--
	case 0x26: // MVI H
		state.pc++
		state.d = state.memory[state.pc]
	case 0x27: // DAA
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
		state.pc++
		state.l = state.memory[state.pc]
		state.pc++
		state.h = state.memory[state.pc]
	case 0x2b: // DCX H
		hl := bytesToPair(state.h, state.l)
		state.h, state.l = pairToBytes(hl - 1)
	case 0x2c: // INR L
		state.l++
	case 0x2d: //  DCR L
		state.l--
	case 0x2e: // MVI L
		state.pc++
		state.l = state.memory[state.pc]
	case 0x2f: // CMA needs work to update the flags register and check for overflow
		state.a = state.a ^ 0xFF
	case 0x30: // SIM ??
	case 0x31: // LXI SP
		state.pc++
		p := state.memory[state.pc]
		state.pc++
		s := state.memory[state.pc]
		state.sp = bytesToPair(s, p)
	case 0x32: // STA
		state.pc++
		lo := state.memory[state.pc]
		state.pc++
		hi := state.memory[state.pc]
		addr := bytesToPair(hi, lo)
		state.memory[addr] = state.a
	case 0x33: // INX SP NOT CORRECT!
		state.sp++
	case 0x34: // INR M
		m := bytesToPair(state.h, state.l)
		state.memory[m]++
	case 0x35: // DCR M
		m := bytesToPair(state.h, state.l)
		state.memory[m]--
	case 0x36: // MVI M
		state.pc++
		m := bytesToPair(state.h, state.l)
		state.memory[m] = state.memory[state.pc]
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
		state.pc++
		lo := state.memory[state.pc]
		state.pc++
		hi := state.memory[state.pc]
		m := bytesToPair(hi, lo)
		state.a = state.memory[m]
	case 0x3b: // DCX SP
		state.sp--
	case 0x3c: // INR A
		state.a++
	case 0x3d: // DCR A
		state.a--
	case 0x3e: // MVI A
		state.pc++
		state.a = state.memory[state.pc]
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
	case 0x58:
	case 0x59:
	case 0x5a:
	case 0x5b:
	case 0x5c:
	case 0x5d:
	case 0x5e:
	case 0x5f:
	case 0x60:
	case 0x61:
	case 0x62:
	case 0x63:
	case 0x64:
	case 0x65:
	case 0x66:
	case 0x67:
	case 0x68:
	case 0x69:
	case 0x6a:
	case 0x6b:
	case 0x6c:
	case 0x6d:
	case 0x6e:
	case 0x6f:
	case 0x70:
	case 0x71:
	case 0x72:
	case 0x73:
	case 0x74:
	case 0x75:
	case 0x76:
	case 0x77:
	case 0x78:
	case 0x79:
	case 0x7a:
	case 0x7b:
	case 0x7c:
	case 0x7d:
	case 0x7e:
	case 0x7f:
	case 0x80:
	case 0x81:
	case 0x82:
	case 0x83:
	case 0x84:
	case 0x85:
	case 0x86:
	case 0x87:
	case 0x88:
	case 0x89:
	case 0x8a:
	case 0x8b:
	case 0x8c:
	case 0x8d:
	case 0x8e:
	case 0x8f:
	case 0x90:
	case 0x91:
	case 0x92:
	case 0x93:
	case 0x94:
	case 0x95:
	case 0x96:
	case 0x97:
	case 0x98:
	case 0x99:
	case 0x9a:
	case 0x9b:
	case 0x9c:
	case 0x9d:
	case 0x9e:
	case 0x9f:
	case 0xa0:
	case 0xa1:
	case 0xa2:
	case 0xa3:
	case 0xa4:
	case 0xa5:
	case 0xa6:
	case 0xa7:
	case 0xa8:
	case 0xa9:
	case 0xaa:
	case 0xab:
	case 0xac:
	case 0xad:
	case 0xae:
	case 0xaf:
	case 0xb0:
	case 0xb1:
	case 0xb2:
	case 0xb3:
	case 0xb4:
	case 0xb5:
	case 0xb6:
	case 0xb7:
	case 0xb8:
	case 0xb9:
	case 0xba:
	case 0xbb:
	case 0xbc:
	case 0xbd:
	case 0xbe:
	case 0xbf:
	case 0xc0:
	case 0xc1:
	case 0xc2:
	case 0xc3:
	case 0xc4:
	case 0xc5:
	case 0xc6:
	case 0xc7:
	case 0xc8:
	case 0xc9:
	case 0xca:
	case 0xcb:
	case 0xcc:
	case 0xcd:
	case 0xce:
	case 0xcf:
	case 0xd0:
	case 0xd1:
	case 0xd2:
	case 0xd3:
	case 0xd4:
	case 0xd5:
	case 0xd6:
	case 0xd7:
	case 0xd8:
	case 0xd9:
	case 0xda:
	case 0xdb:
	case 0xdc:
	case 0xdd:
	case 0xde:
	case 0xdf:
	case 0xe0:
	case 0xe1:
	case 0xe2:
	case 0xe3:
	case 0xe4:
	case 0xe5:
	case 0xe6:
	case 0xe7:
	case 0xe8:
	case 0xe9:
	case 0xea:
	case 0xeb:
	case 0xec:
	case 0xed:
	case 0xee:
	case 0xef:
	case 0xf0:
	case 0xf1:
	case 0xf2:
	case 0xf3:
	case 0xf4:
	case 0xf5:
	case 0xf6:
	case 0xf7:
	case 0xf8:
	case 0xf9:
	case 0xfa:
	case 0xfb:
	case 0xfc:
	case 0xfd:
	case 0xfe:
	case 0xff:
	}
}
