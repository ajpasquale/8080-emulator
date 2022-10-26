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
	state.memory = slices.Insert(state.memory, offset, bs...)
}

func Emulate8080(state *state8080) {
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
		state.b++
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
		fmt.Println(opCode, "-")
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
		state.c = state.memory[state.pc]
		state.pc++
	case 0x0f: // RRC
		r, cy := shiftRight8(state.a, 1)
		state.a = r
		state.cc.cy = cy
	case 0x10: // -
		fmt.Println(opCode, "-")
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
		state.d++
	case 0x15: // DCR D
		state.d--
	case 0x16: // MVI D
		state.d = state.memory[state.pc]
		state.pc++
	case 0x17: // "RAL", 1
		r, cb := shiftLeft8(state.a, 1)
		state.a = uint8(r + state.cc.cy)
		state.cc.cy = cb
	case 0x18: // -
		fmt.Println(opCode, "-")
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
	case 0x1e: // MVI E
		state.e = state.memory[state.pc]
		state.pc++
	case 0x1f: // RAR
		r, lsb := shiftRight8(state.a, 1)
		state.a = uint8(r + (state.cc.cy << 7))
		state.cc.cy = lsb
	case 0x20: // RIM ??
		fmt.Println(opCode, "RIM")
	case 0x21: // LXI H
		state.l = state.memory[state.pc]
		state.pc++
		state.h = state.memory[state.pc]
		state.pc++
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
		state.h = state.memory[state.pc]
		state.pc++
	case 0x27: // DAA
		fmt.Println(opCode, "DAA")
	case 0x28: // -
		fmt.Println(opCode, "-")
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

		state.l = state.memory[state.pc]
		state.pc++
	case 0x2f: // CMA needs work to update the flags register and check for overflow
		state.a = state.a ^ 0xFF
	case 0x30: // SIM ??
		fmt.Println(opCode, "SIM")
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
	case 0x33: // INX SP NOT CORRECT!
		state.sp++
	case 0x34: // INR M
		m := bytesToPair(state.h, state.l)
		state.memory[m]++
	case 0x35: // DCR M
		m := bytesToPair(state.h, state.l)
		state.memory[m]--
	case 0x36: // MVI M
		m := bytesToPair(state.h, state.l)
		state.memory[m] = state.memory[state.pc]
		state.pc++
	case 0x37: // STC
		state.cc.cy = 1
	case 0x38: // -
		fmt.Println(opCode, "-")
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
		state.a++
	case 0x3d: // DCR A
		state.a--
	case 0x3e: // MVI A
		state.a = state.memory[state.pc]
		state.pc++
	case 0x3f: // CMC
		state.cc.cy--
	case 0x40: // MOV B,B
		fmt.Println(opCode, "MOV B,B")
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
		fmt.Println(opCode, "MOV C,C")
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
		fmt.Println(opCode, "MOV D,D")
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
		fmt.Println(opCode, "MOV E,E")
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
		fmt.Println(opCode, "MOV H, H")
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
		fmt.Println(opCode, "MOV L, L")
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
		fmt.Println(opCode, "HLT")
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
		fmt.Println(opCode, "MOV A, A")
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
			hi := state.memory[state.pc+2]
			lo := state.memory[state.pc+1]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xc5: // PUSH B
		state.memory[state.sp-1] = state.b
		state.memory[state.sp-2] = state.c
		state.sp -= 2
	case 0xc6: // ADI
		state.a = add8(state, state.a, state.memory[state.pc+1])
		state.pc++
	case 0xc7: // RST 0
		fmt.Println("RST 0")
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
			hi := state.memory[state.pc+2]
			lo := state.memory[state.pc+1]
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
		fmt.Println(opCode, "ACI")
	case 0xcf: // RST 1
		fmt.Println(opCode, "RST 1")
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
		fmt.Println(opCode, "OUT")
	case 0xd4: // CNC
		if state.cc.cy == 0 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+2]
			lo := state.memory[state.pc+1]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xd5: // PUSH D
		state.memory[state.sp-1] = state.d
		state.memory[state.sp-2] = state.e
		state.sp -= 2
	case 0xd6: // SUI
		fmt.Println(opCode, "SUI")
	case 0xd7: // RST 2
		fmt.Println(opCode, "RST 2")
	case 0xd8: // RC
		if state.cc.cy == 1 {
			state.pc = bytesToPair(state.memory[state.sp+1], state.memory[state.sp])
			state.sp += 2
		}
	case 0xd9: // -
		fmt.Println(opCode, "-")
	case 0xda: // JC
		if state.cc.cy == 1 {
			state.pc = bytesToPair(state.memory[state.pc+1], state.memory[state.pc])
		} else {
			state.pc += 2
		}
	case 0xdb: // IN
		fmt.Println(opCode, "IN")
	case 0xdc: // CC
		fmt.Println(opCode, "CC")
	case 0xdd: // -
		fmt.Println(opCode, "-")
	case 0xde: // SBI
		fmt.Println(opCode, "SBI")
	case 0xdf: // RST 3
		fmt.Println(opCode, "RST 3")
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
		fmt.Println(opCode, "XTHL")
	case 0xe4: // CPO
		if state.cc.p == 0 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+2]
			lo := state.memory[state.pc+1]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xe5: // PUSH H
		state.memory[state.sp-1] = state.h
		state.memory[state.sp-2] = state.l
		state.sp -= 2
	case 0xe6: // ANI
		state.a = state.a & state.memory[state.pc+1]
		setLogicFlags(state)
		state.pc++
	case 0xe7: // RST 4
		fmt.Println("RST 4")
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
			hi := state.memory[state.pc+2]
			lo := state.memory[state.pc+1]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xed: // -
		fmt.Println(opCode, "-")
	case 0xee: // XRI
		state.a = state.a & state.memory[state.pc+1]
		setLogicFlags(state)
		state.pc++
	case 0xef: // RST 5
		fmt.Println(opCode, "RST 5")
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
		fmt.Println(opCode, "DI")
		state.int_enable = 0
	case 0xf4: // CP
		if state.cc.s == 0 {
			ret := state.pc + 2
			state.memory[state.sp-1], state.memory[state.sp-2] = pairToBytes(ret)
			state.sp -= 2
			hi := state.memory[state.pc+2]
			lo := state.memory[state.pc+1]
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
		state.a = state.a | state.memory[state.pc+1]
		setLogicFlags(state)
		state.pc++
	case 0xf7: // RST 6
		fmt.Println(opCode, "RST 6")
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
			hi := state.memory[state.pc+2]
			lo := state.memory[state.pc+1]
			state.pc = bytesToPair(hi, lo)
		} else {
			state.pc += 2
		}
	case 0xfd: // -
		fmt.Println(opCode, "-")
	case 0xfe: // CPI
		sub8(state, state.a, state.memory[state.pc])
		state.pc++
	case 0xff: // RST 7
		fmt.Println("RST 7")
	}
}
