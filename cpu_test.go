package emulator

import (
	"reflect"
	"testing"
)

func TestCpu(t *testing.T) {

	// cycles := 0
	// state := InitializeState()

	// now := time.Now()
	// LoadSpaceInvaders(state)
	// timer := now
	// for {

	// 	// RST 1 - middle of the screen interrupt
	// 	if time.Since(timer) > 30*8000*time.Microsecond && GetIntEnabled(state) == 1 {
	// 		//Interrupt8080(state, )
	// 		Restart8080(state, RST1)
	// 	}
	// 	// RST 2 - end of screen interrupt
	// 	if time.Since(timer) > 30*16000*time.Microsecond && GetIntEnabled(state) == 1 {
	// 		// problem near 17cd db 02	 IN	 $02
	// 		Restart8080(state, RST2)
	// 		timer = time.Now()
	// 	}

	// 	// if state.pc == 0x93 &&
	// 	// 	state.sp == 0x94 &&
	// 	// 	state.h == 0x3e &&
	// 	// 	state.d == 0x1f {
	// 	// 	fmt.Println("break")
	// 	// }

	// 	// if cycles == 2707677 {
	// 	// 	fmt.Println("break")
	// 	// }

	// 	// INPUT
	// 	GetInput(state)
	// 	// OUTPUT
	// 	GetOutput(state)
	// 	// fmt.Printf("pc: %x, sp: %x, a: %x, h: %x, l: %x, d: %x, e: %x cycle: %d\n",
	// 	// 	state.pc,
	// 	// 	state.sp,
	// 	// 	state.a,
	// 	// 	state.h,
	// 	// 	state.l,
	// 	// 	state.d,
	// 	// 	state.e,
	// 	// 	cycles)
	// 	fmt.Printf("sp: %x\n", state.sp)
	// 	cycles += Emulate8080(state)

	//}
}
func TestInstructionSTAX(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint16
	}{
		//        instr a     b     c               a      memory
		{[]uint8{0x02, 0xDD, 0x00, 0x01}, []uint16{0x00DD, 0x0001}},
		{[]uint8{0x12, 0xDD, 0x00, 0x01}, []uint16{0x00DD, 0x0001}},
	}

	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[1]
		state.memory = append(state.memory, tt.in[0])

		switch tt.in[0] {
		case 0x02:
			state.b = tt.in[2]
			state.c = tt.in[3]
		case 0x12:
			state.d = tt.in[2]
			state.e = tt.in[3]
		default:
			t.Errorf("TestInstructionSTAX missing case: %x", tt.in[0])
		}

		state.memory = append(state.memory, 0x00)
		Emulate8080(state)

		if !reflect.DeepEqual(state.memory[tt.want[1]], uint8(tt.want[0])) {
			t.Errorf("TestInstructionSTAX(%q)\nhave %v \nwant %v", tt.in, state.b, tt.want[0])
		}
	}
}

func TestInstructionINX(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint16
	}{
		// instruction, state 1, state 2, val, before val   want: after val
		{[]uint8{0x03, 0x00, 0x01}, 0x0002},
		{[]uint8{0x13, 0x00, 0x02}, 0x0003},
		{[]uint8{0x23, 0x00, 0x03}, 0x0004},
		{[]uint8{0x33, 0x00, 0x04}, 0x0005},

		{[]uint8{0x03, 0x00, 0xFF}, 0x0100},
		{[]uint8{0x13, 0xFF, 0x00}, 0xFF01},
		{[]uint8{0x23, 0xFF, 0xFF}, 0x0000},
		{[]uint8{0x33, 0xFA, 0xE0}, 0xFAE1},
	}

	for _, tt := range tests {
		state := newState8080()
		instruction := tt.in[0]

		state.memory = append(state.memory, instruction)

		switch instruction {
		case 0x03: // B
			state.b = tt.in[1]
			state.c = tt.in[2]
		case 0x13: // D
			state.d = tt.in[1]
			state.e = tt.in[2]
		case 0x23: // H
			state.h = tt.in[1]
			state.l = tt.in[2]
		case 0x33: // SP
			hi := tt.in[1]
			lo := tt.in[2]
			state.sp = bytesToPair(hi, lo)
		default:
			t.Errorf("TestInstructionINX missing case: %x", instruction)
		}

		Emulate8080(state)

		var reg uint16
		switch instruction {
		case 0x03: // B
			reg = bytesToPair(state.b, state.c)
		case 0x13: // D
			reg = bytesToPair(state.d, state.e)
		case 0x23: // H
			reg = bytesToPair(state.h, state.l)
		case 0x33: // SP
			reg = state.sp
		default:
			t.Errorf("TestInstructionINX missing case: %x", instruction)
		}

		if !reflect.DeepEqual(reg, tt.want) {
			t.Errorf("TestInstructionINX(%q)\nhave %v \nwant %v", tt.in, reg, tt.want)
		}
	}
}

func TestInstructionINR(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x04, 0x00}, 0x01}, // B
		{[]uint8{0x0c, 0x00}, 0x01}, // C
		{[]uint8{0x14, 0x00}, 0x01}, // D
		{[]uint8{0x1c, 0x00}, 0x01}, // E
		{[]uint8{0x24, 0x00}, 0x01}, // H
		{[]uint8{0x2c, 0x00}, 0x01}, // L
		{[]uint8{0x34, 0x00}, 0x01}, // M
		{[]uint8{0x3c, 0x00}, 0x01}, // A

		{[]uint8{0x04, 0x99}, 0x9A}, // B
		{[]uint8{0x0c, 0x99}, 0x9A}, // C
		{[]uint8{0x14, 0x99}, 0x9A}, // D
		{[]uint8{0x1c, 0x99}, 0x9A}, // E
		{[]uint8{0x24, 0x99}, 0x9A}, // H
		{[]uint8{0x2c, 0x99}, 0x9A}, // L
		{[]uint8{0x34, 0x99}, 0x9A}, // M
		{[]uint8{0x3c, 0x99}, 0x9A}, // A

		{[]uint8{0x04, 0xFF}, 0x00}, // B
		{[]uint8{0x0c, 0xFF}, 0x00}, // C
		{[]uint8{0x14, 0xFF}, 0x00}, // D
		{[]uint8{0x1c, 0xFF}, 0x00}, // E
		{[]uint8{0x24, 0xFF}, 0x00}, // H
		{[]uint8{0x2c, 0xFF}, 0x00}, // L
		{[]uint8{0x34, 0xFF}, 0x00}, // M
		{[]uint8{0x3c, 0xFF}, 0x00}, // A
	}

	var reg uint8

	for _, tt := range tests {
		state := newState8080()
		state.memory = append(state.memory, tt.in[0])
		switch tt.in[0] {
		case 0x04: // INR B
			state.b = tt.in[1]
		case 0x0c: // INR C
			state.c = tt.in[1]
		case 0x14: // INR D
			state.d = tt.in[1]
		case 0x1c: // INR E
			state.e = tt.in[1]
		case 0x24: // INR H
			state.h = tt.in[1]
		case 0x2c: // INR L
			state.l = tt.in[1]
		case 0x34:
			state.h = 0x00
			state.l = 0x01
			state.memory = append(state.memory, tt.in[1])
		case 0x3c: // INR A
			state.a = tt.in[1]
		default:
			t.Errorf("TestInstructionINR invalid case: %x", tt.in[0])
		}

		Emulate8080(state)

		switch tt.in[0] {
		case 0x04: // INR B
			reg = state.b
		case 0x0c: // INR C
			reg = state.c
		case 0x14: // INR D
			reg = state.d
		case 0x1c: // INR E
			reg = state.e
		case 0x24: // INR H
			reg = state.h
		case 0x2c: // INR L
			reg = state.l
		case 0x34:
			reg = state.memory[0x01]
		case 0x3c: // INR A
			reg = state.a
		default:
			t.Errorf("TestInstructionINR invalid case: %x", tt.in[0])
		}

		if !reflect.DeepEqual(reg, tt.want) {
			t.Errorf("TestInstructionINR reg(%x)\nhave %v \nwant %v", tt.in[0], reg, tt.want)
		}
	}
}

func TestInstructionDCR(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x05, 0x01}, 0x00}, // B
		{[]uint8{0x0d, 0x01}, 0x00}, // C
		{[]uint8{0x15, 0x01}, 0x00}, // D
		{[]uint8{0x1d, 0x01}, 0x00}, // E
		{[]uint8{0x25, 0x01}, 0x00}, // H
		{[]uint8{0x2d, 0x01}, 0x00}, // L
		{[]uint8{0x35, 0x01}, 0x00}, // M
		{[]uint8{0x3d, 0x01}, 0x00}, // A

		{[]uint8{0x05, 0x9A}, 0x99}, // B
		{[]uint8{0x0d, 0x9A}, 0x99}, // C
		{[]uint8{0x15, 0x9A}, 0x99}, // D
		{[]uint8{0x1d, 0x9A}, 0x99}, // E
		{[]uint8{0x25, 0x9A}, 0x99}, // H
		{[]uint8{0x2d, 0x9A}, 0x99}, // L
		{[]uint8{0x35, 0x9A}, 0x99}, // M
		{[]uint8{0x3d, 0x9A}, 0x99}, // A

		{[]uint8{0x05, 0x00}, 0xFF}, // B
		{[]uint8{0x0d, 0x00}, 0xFF}, // C
		{[]uint8{0x15, 0x00}, 0xFF}, // D
		{[]uint8{0x1d, 0x00}, 0xFF}, // E
		{[]uint8{0x25, 0x00}, 0xFF}, // H
		{[]uint8{0x2d, 0x00}, 0xFF}, // L
		{[]uint8{0x35, 0x00}, 0xFF}, // M
		{[]uint8{0x3d, 0x00}, 0xFF}, // A
	}

	var reg uint8

	for _, tt := range tests {
		state := newState8080()
		state.memory = append(state.memory, tt.in[0])
		switch tt.in[0] {
		case 0x05: // DCR B
			state.b = tt.in[1]
		case 0x0d: // DCR C
			state.c = tt.in[1]
		case 0x15: // DCR D
			state.d = tt.in[1]
		case 0x1d: // DCR E
			state.e = tt.in[1]
		case 0x25: // DCR H
			state.h = tt.in[1]
		case 0x2d: // DCR L
			state.l = tt.in[1]
		case 0x35: // DCR M
			state.h = 0x00
			state.l = 0x01
			state.memory = append(state.memory, tt.in[1])
		case 0x3d: // DCR A
			state.a = tt.in[1]
		default:
			t.Errorf("TestInstructionDCR invalid case: %x", tt.in[0])
		}

		Emulate8080(state)

		switch tt.in[0] {
		case 0x05: // DCR B
			reg = state.b
		case 0x0d: // DCR C
			reg = state.c
		case 0x15: // DCR D
			reg = state.d
		case 0x1d: // DCR E
			reg = state.e
		case 0x25: // DCR H
			reg = state.h
		case 0x2d: // DCR L
			reg = state.l
		case 0x35: // DCR M
			reg = state.memory[0x01]
		case 0x3d: // DCR A
			reg = state.a
		default:
			t.Errorf("TestInstructionDCR invalid case: %x", tt.in[0])
		}

		if !reflect.DeepEqual(reg, tt.want) {
			t.Errorf("TestInstructionDCR reg(%x)\nhave %v \nwant %v", tt.in[0], reg, tt.want)
		}
	}
}

func TestInstructionRLC(t *testing.T) {
	tests := []struct {
		in   uint8
		want []uint8
	}{
		{0x00, []uint8{0x00, 0}},
		{0x35, []uint8{0x6a, 0}},
		{0x95, []uint8{0x2b, 1}},
		{0xFF, []uint8{0xFF, 1}},
	}

	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in
		state.memory = append(state.memory, 0x07)
		Emulate8080(state)
		if !reflect.DeepEqual(state.a, tt.want[0]) {
			t.Errorf("TestInstructionRLC(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want[0])
		}
		if !reflect.DeepEqual(state.cc.cy, tt.want[1]) {
			t.Errorf("TestInstructionRLC(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[1])
		}
	}
}

func TestInstructionDAD(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint16
	}{
		// instruction state1,state2, state3,state4   value, carry flag
		// HL + HL only uses state3 and state4
		{[]uint8{0x09, 0x00, 0x00, 0x00, 0x00}, []uint16{0x0000, 0}},
		{[]uint8{0x19, 0x00, 0x00, 0x00, 0x00}, []uint16{0x0000, 0}},
		{[]uint8{0x39, 0x00, 0x00, 0x00, 0x00}, []uint16{0x0000, 0}},

		{[]uint8{0x09, 0x80, 0x00, 0x00, 0x08}, []uint16{0x8008, 0}},
		{[]uint8{0x19, 0x33, 0x9F, 0xA1, 0x7B}, []uint16{0xD51A, 0}},
		{[]uint8{0x39, 0x00, 0x40, 0x00, 0x40}, []uint16{0x0080, 0}},

		{[]uint8{0x09, 0xFF, 0x00, 0x00, 0xFF}, []uint16{0xFFFF, 0}},
		{[]uint8{0x19, 0x33, 0x9F, 0xA1, 0x7B}, []uint16{0xD51A, 0}},
		{[]uint8{0x39, 0x00, 0x40, 0x00, 0x40}, []uint16{0x0080, 0}},

		{[]uint8{0x29, 0x00, 0x00, 0xFF, 0xFF}, []uint16{0xFFFE, 1}},
		{[]uint8{0x29, 0x00, 0x00, 0x00, 0xFF}, []uint16{0x01FE, 0}},
	}
	for _, tt := range tests {
		state := newState8080()
		instruction := tt.in[0]
		state.memory = append(state.memory, instruction)

		switch instruction {
		case 0x09: // DAD B (BC+HL) -> HL
			state.b = tt.in[1]
			state.c = tt.in[2]
		case 0x19: // DAD D (DE+HL) -> HL
			state.d = tt.in[1]
			state.e = tt.in[2]
		case 0x29: // DAD H (HL+HL) -> HL
		case 0x39: // DAD SP (SP+HL) -> HL
			hi := tt.in[1]
			lo := tt.in[2]
			state.sp = bytesToPair(hi, lo)
		default:
			t.Errorf("TestInstructionDAD missing case: %x", instruction)
		}

		state.h = tt.in[3]
		state.l = tt.in[4]

		Emulate8080(state)

		if !reflect.DeepEqual(bytesToPair(state.h, state.l), tt.want[0]) {
			t.Errorf("TestInstructionDAD(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want[0])
		}
		if !reflect.DeepEqual(state.cc.cy, uint8(tt.want[1])) {
			t.Errorf("TestInstructionDAD carry(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[1])
		}
	}
}

func TestInstructionLDA(t *testing.T) {
	ptests := []struct {
		in   []uint8
		want uint8
	}{
		//instruction, lo, hi, value
		{[]uint8{0x3a, 0x05, 0x00, 0x00}, 0x00},
		{[]uint8{0x3a, 0x05, 0x00, 0x01}, 0x01},
		{[]uint8{0x3a, 0x05, 0x00, 0xFF}, 0xFF},
	}

	for _, tt := range ptests {
		state := newState8080()
		instruction := tt.in[0]

		state.memory = append(state.memory, instruction)
		state.memory = append(state.memory, tt.in[1]) //lo addr
		state.memory = append(state.memory, tt.in[2]) //hi addr
		state.memory = append(state.memory, 0x00)     // buffer
		state.memory = append(state.memory, 0x00)     // buffer
		state.memory = append(state.memory, tt.in[3]) // value to be moved to state.a

		Emulate8080(state)

		have := state.a

		if !reflect.DeepEqual(have, tt.want) {
			t.Errorf("TestInstructionLDA(%q)\nhave %v \nwant %v", tt.in, have, tt.want)
		}
	}

	ntests := []struct {
		in   []uint8
		want uint8
	}{
		//instruction, lo, hi, value
		{[]uint8{0x3a, 0x05, 0x00, 0x01}, 0x00},
		{[]uint8{0x3a, 0x05, 0x00, 0x00}, 0x01},
		{[]uint8{0x3a, 0x05, 0x00, 0xFA}, 0xFF},
	}

	for _, tt := range ntests {
		state := newState8080()
		instruction := tt.in[0]

		state.memory = append(state.memory, instruction)
		state.memory = append(state.memory, tt.in[1]) //lo addr
		state.memory = append(state.memory, tt.in[2]) //hi addr
		state.memory = append(state.memory, 0x00)     // buffer
		state.memory = append(state.memory, 0x00)     // buffer
		state.memory = append(state.memory, tt.in[3]) // value to be moved to state.a

		Emulate8080(state)

		have := state.a

		if reflect.DeepEqual(have, tt.want) {
			t.Errorf("TestInstructionLDA(%q)\nhave %v \nwant %v", tt.in, have, tt.want)
		}
	}

}
func TestInstructionLHLD(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint8
	}{
		//    instruction, val1, val2
		{[]uint8{0x2a, 0x00, 0x00}, []uint8{0x00, 0x00}},
		{[]uint8{0x2a, 0x01, 0x00}, []uint8{0x01, 0x00}},
		{[]uint8{0x2a, 0x01, 0x01}, []uint8{0x01, 0x01}},
		{[]uint8{0x2a, 0x02, 0x03}, []uint8{0x02, 0x03}},
		{[]uint8{0x2a, 0x05, 0x03}, []uint8{0x05, 0x03}},
		{[]uint8{0x2a, 0xFF, 0xFF}, []uint8{0xFF, 0xFF}},
	}

	for _, tt := range tests {
		state := newState8080()
		instruction := tt.in[0]

		state.memory = append(state.memory, instruction)
		state.memory = append(state.memory, tt.in[1]) // state.l
		state.memory = append(state.memory, tt.in[2]) // state.h

		Emulate8080(state)

		if !reflect.DeepEqual(state.l, tt.want[0]) {
			t.Errorf("TestInstructionLHLD(%q)\nhave %v \nwant %v", tt.in, state.l, tt.want[0])
		}
		if !reflect.DeepEqual(state.h, tt.want[1]) {
			t.Errorf("TestInstructionLHLD(%q)\nhave %v \nwant %v", tt.in, state.l, tt.want[1])
		}
	}
}

func TestInstructionCMC(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		// instruction, carry, expected carry
		{[]uint8{0x3f, 0x00}, 0x01},
		{[]uint8{0x3f, 0x01}, 0x00},
	}

	for _, tt := range tests {
		state := newState8080()
		instruction := tt.in[0]
		state.cc.cy = tt.in[1]

		state.memory = append(state.memory, instruction)

		Emulate8080(state)

		if !reflect.DeepEqual(state.cc.cy, tt.want) {
			t.Errorf("TestInstructionCMC(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want)
		}
	}
}

func TestInstructionSHLD(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint8
	}{
		//instruction   l     h              l      h
		{[]uint8{0x22, 0x00, 0x00}, []uint8{0x00, 0x00}},
		{[]uint8{0x22, 0x01, 0x01}, []uint8{0x01, 0x01}},
		{[]uint8{0x22, 0xFF, 0xFF}, []uint8{0xFF, 0xFF}},
		{[]uint8{0x22, 0x00, 0xFF}, []uint8{0x00, 0xFF}},
		{[]uint8{0x22, 0xFF, 0x00}, []uint8{0xFF, 0x00}},
	}

	for _, tt := range tests {
		state := newState8080()
		instruction := tt.in[0]

		state.memory = append(state.memory, instruction)

		state.memory = append(state.memory, 0x00) // state.l
		state.memory = append(state.memory, 0x00) // state.h

		state.l = tt.in[1]
		state.h = tt.in[2]

		Emulate8080(state)

		if !reflect.DeepEqual(state.memory[1], tt.want[0]) {
			t.Errorf("TestInstructionLHLD(%q)\nhave %v \nwant %v", tt.in, state.l, tt.want[0])
		}
		if !reflect.DeepEqual(state.memory[2], tt.want[1]) {
			t.Errorf("TestInstructionLHLD(%q)\nhave %v \nwant %v", tt.in, state.l, tt.want[1])
		}
	}
}
func TestInstructionSTA(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		//instruction, lo, hi, state.a
		{[]uint8{0x32, 0x05, 0x00, 0x01}, 0x01},
		{[]uint8{0x32, 0x05, 0x00, 0x0F}, 0x0F},
		{[]uint8{0x32, 0x05, 0x00, 0xFF}, 0xFF},
	}

	for _, tt := range tests {
		state := newState8080()
		instruction := tt.in[0]

		state.memory = append(state.memory, instruction)
		state.memory = append(state.memory, tt.in[1]) // lo addr
		state.memory = append(state.memory, tt.in[2]) // hi addr
		state.memory = append(state.memory, 0x00)     // buffer
		state.memory = append(state.memory, 0x00)     // buffer
		state.memory = append(state.memory, 0x00)     // value pointed to by lo/hi addr

		state.a = tt.in[3]

		Emulate8080(state)

		addr := bytesToPair(tt.in[2], tt.in[1])

		if !reflect.DeepEqual(state.memory[addr], tt.want) {
			t.Errorf("TestInstructionSTA(%q)\nhave %v \nwant %v", tt.in, state.memory[addr], tt.want)
		}
	}
}

func TestInstructionLDAX(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x0a, 0x00, 0x01, 0x02}, 0x02},
		{[]uint8{0x1a, 0x00, 0x01, 0x04}, 0x04},
	}

	for _, tt := range tests {
		state := newState8080()
		instruction := tt.in[0]

		state.memory = append(state.memory, instruction)

		// add one memory location for the
		state.memory = append(state.memory, tt.in[3])

		switch instruction {
		case 0x0a: // LDAX B -> A
			state.b = tt.in[1]
			state.c = tt.in[2]
		case 0x1a: // LDAX D -> A
			state.d = tt.in[1]
			state.e = tt.in[2]
		default:
			t.Errorf("TestInstructionLDAX missing case: %x", instruction)
		}

		Emulate8080(state)

		if !reflect.DeepEqual(state.a, uint8(tt.want)) {
			t.Errorf("TestInstructionLDAX(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionDCX(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint16
	}{
		// instruction, state 1, state 2, val, before val   want: after val
		{[]uint8{0x0b, 0x00, 0x02}, 0x0001},
		{[]uint8{0x1b, 0x00, 0x02}, 0x0001},
		{[]uint8{0x2b, 0x00, 0x02}, 0x0001},
		{[]uint8{0x3b, 0x00, 0x02}, 0x0001},

		{[]uint8{0x0b, 0x00, 0xFF}, 0x00FE},
		{[]uint8{0x1b, 0x00, 0x00}, 0xFFFF},
		{[]uint8{0x2b, 0xFF, 0x00}, 0xFEFF},
		{[]uint8{0x3b, 0xAE, 0x8C}, 0xAE8B},
	}

	for _, tt := range tests {
		state := newState8080()
		instruction := tt.in[0]

		state.memory = append(state.memory, instruction)

		switch instruction {
		case 0x0b: //  B
			//bc
			state.b = tt.in[1]
			state.c = tt.in[2]
		case 0x1b: //  D
			//de
			state.d = tt.in[1]
			state.e = tt.in[2]
		case 0x2b: //  H
			//hl
			state.h = tt.in[1]
			state.l = tt.in[2]
		case 0x3b: //  SP
			//sp
			hi := tt.in[1]
			lo := tt.in[2]
			state.sp = bytesToPair(hi, lo)
		default:
			t.Errorf("TestInstructionDCX missing case: %x", instruction)
		}

		Emulate8080(state)

		var reg uint16
		switch instruction {
		case 0x0b: //  B
			reg = bytesToPair(state.b, state.c)
		case 0x1b: //  D
			reg = bytesToPair(state.d, state.e)
		case 0x2b: //  H
			reg = bytesToPair(state.h, state.l)
		case 0x3b: //  SP
			reg = state.sp
		default:
			t.Errorf("TestInstructionDCX missing case: %x", instruction)
		}

		if !reflect.DeepEqual(reg, tt.want) {
			t.Errorf("TestInstructionDCX(%q)\nhave %v \nwant %v", tt.in, reg, tt.want)
		}
	}
}

func TestInstructionRRC(t *testing.T) {
	tests := []struct {
		in   uint8
		want []uint8
	}{
		{0x00, []uint8{0x0, 0}},
		{0x8A, []uint8{0x45, 0}},
		{0x81, []uint8{0x40, 1}},
		{0xFF, []uint8{0x7F, 1}},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in
		state.memory = append(state.memory, 0x0f)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want[0]) {
			t.Errorf("TestInstructionRLC(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
		if !reflect.DeepEqual(state.cc.cy, uint8(tt.want[1])) {
			t.Errorf("TestInstructionRLC(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want)
		}
	}
}

func TestInstructionRAL(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint8
	}{
		{[]uint8{0x00, 0}, []uint8{0x0, 0}},
		{[]uint8{0x00, 1}, []uint8{0x1, 0}},
		{[]uint8{0x35, 0}, []uint8{0x6a, 0}},
		{[]uint8{0x35, 1}, []uint8{0x6b, 0}},
		{[]uint8{0x95, 0}, []uint8{0x2a, 1}},
		{[]uint8{0x95, 1}, []uint8{0x2b, 1}},
		{[]uint8{0xFF, 0}, []uint8{0xFE, 1}},
		{[]uint8{0xFF, 1}, []uint8{0xFF, 1}},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.cc.cy = tt.in[1]
		state.memory = append(state.memory, 0x17)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want[0]) {
			t.Errorf("TestInstructionRAL(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
		if !reflect.DeepEqual(state.cc.cy, uint8(tt.want[1])) {
			t.Errorf("TestInstructionRAL(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want)
		}
	}
}

func TestInstructionRAR(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint8
	}{
		{[]uint8{0x00, 0}, []uint8{0x0, 0}},
		{[]uint8{0x00, 1}, []uint8{0x80, 0}},
		{[]uint8{0x35, 0}, []uint8{0x1A, 1}},
		{[]uint8{0x35, 1}, []uint8{0x9A, 1}},
		{[]uint8{0x95, 0}, []uint8{0x4A, 1}},
		{[]uint8{0x95, 1}, []uint8{0xCA, 1}},
		{[]uint8{0xFF, 0}, []uint8{0x7F, 1}},
		{[]uint8{0xFF, 1}, []uint8{0xFF, 1}},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.cc.cy = tt.in[1]
		state.memory = append(state.memory, 0x1f)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want[0]) {
			t.Errorf("TestInstructionRAR(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want[0])
		}
		if !reflect.DeepEqual(state.cc.cy, uint8(tt.want[1])) {
			t.Errorf("TestInstructionRAR(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[1])
		}
	}
}

func TestInstructionCMA(t *testing.T) {
	tests := []struct {
		in   uint8
		want uint8
	}{
		{0x00, 0xFF},
		{0xFF, 0x00},
		{0xAA, 0x55},
		{0x51, 0xAE},
		{0x89, 0x76},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in
		state.memory = append(state.memory, 0x2f)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionCMA(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionADDB(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x00, 0x00}, 0x00},
		{[]uint8{0x00, 0x01}, 0x01},
		{[]uint8{0x01, 0x01}, 0x02},
		{[]uint8{0x01, 0x0F}, 0x10},
		{[]uint8{0xFF, 0x01}, 0x00},
		{[]uint8{0xFF, 0x02}, 0x01},
		{[]uint8{0x2E, 0x6C}, 0x9A},
		{[]uint8{0xFF, 0x0F}, 0x0E},
		{[]uint8{0xFF, 0xFF}, 0xFE},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.b = tt.in[1]
		state.memory = append(state.memory, 0x80)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionADDB(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionADCB(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x00, 0x00, 0}, 0x00},
		{[]uint8{0x00, 0x00, 1}, 0x01},
		{[]uint8{0x3D, 0x42, 0}, 0x7F},
		{[]uint8{0x3E, 0xC1, 1}, 0x00},
		{[]uint8{0x3D, 0x42, 1}, 0x80},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.b = tt.in[1]
		state.cc.cy = tt.in[2]
		state.memory = append(state.memory, 0x88)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionADCB(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionSUBB(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x0A, 0x05, 0}, 0x05},
		{[]uint8{0x02, 0x05, 0}, 0xFD},
		{[]uint8{0xE5, 0x05, 0}, 0xE0},
		{[]uint8{0x3E, 0x3E, 0}, 0x00},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.b = tt.in[1]
		state.memory = append(state.memory, 0x90)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionSUBB(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionSBBB(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x04, 0x02, 1}, 0x01},
		{[]uint8{0x3E, 0x3E, 0}, 0x00},
		{[]uint8{0x04, 0x02, 0}, 0x02},
		{[]uint8{0x3E, 0x3E, 1}, 0xFF},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.b = tt.in[1]
		state.cc.cy = tt.in[2]
		state.memory = append(state.memory, 0x98)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionSBBB(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionCMPB(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint8
	}{
		{[]uint8{0x0A, 0x05}, []uint8{0x0A, 0x05}},
		{[]uint8{0x02, 0x05}, []uint8{0x02, 0x05}},
		{[]uint8{0xE5, 0x05}, []uint8{0xE5, 0x05}},
		{[]uint8{0x3E, 0x3E}, []uint8{0x3E, 0x3E}},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.b = tt.in[1]
		state.memory = append(state.memory, 0xb8)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want[0]) {
			t.Errorf("TestInstructionCMPB(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want[0])
		}
		if !reflect.DeepEqual(state.b, tt.want[1]) {
			t.Errorf("TestInstructionCMPB(%q)\nhave %v \nwant %v", tt.in, state.b, tt.want[1])
		}
	}
}

func TestInstructionLXI(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint8
	}{
		//       instr, op1, op2           hi byte, lo byte
		{[]uint8{0x01, 0xAD, 0xDE}, []uint8{0xDE, 0xAD}},
		{[]uint8{0x11, 0xAD, 0xDE}, []uint8{0xDE, 0xAD}},
		{[]uint8{0x21, 0xAD, 0xDE}, []uint8{0xDE, 0xAD}},
		{[]uint8{0x31, 0xAD, 0xDE}, []uint8{0xDE, 0xAD}},
	}
	for _, tt := range tests {
		var hi uint8
		var lo uint8
		state := newState8080()
		state.memory = append(state.memory, tt.in[0])
		state.memory = append(state.memory, tt.in[1])
		state.memory = append(state.memory, tt.in[2])
		Emulate8080(state)

		switch tt.in[0] {
		case 0x01:
			hi = state.b
			lo = state.c
		case 0x11:
			hi = state.d
			lo = state.e
		case 0x21:
			hi = state.h
			lo = state.l
		case 0x31:
			hi, lo = pairToBytes(state.sp)
		default:
			t.Errorf("TestInstructionLXI missing case: %x", tt.in[0])
		}

		if !reflect.DeepEqual(hi, tt.want[0]) {
			t.Errorf("TestInstructionLXI Hi(%q)\nhave %v \nwant %v", tt.in, hi, tt.want[0])
		}
		if !reflect.DeepEqual(lo, tt.want[1]) {
			t.Errorf("TestInstructionLXI Lo(%q)\nhave %v \nwant %v", tt.in, lo, tt.want[1])
		}
	}
}

func TestInstructionMVI(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint8
	}{
		//       instr, op1
		{[]uint8{0x06, 0xAD}, []uint8{0xAD}},
		{[]uint8{0x0E, 0xAD}, []uint8{0xAD}},
		{[]uint8{0x16, 0xAD}, []uint8{0xAD}},
		{[]uint8{0x1E, 0xAD}, []uint8{0xAD}},
		{[]uint8{0x26, 0xAD}, []uint8{0xAD}},
		{[]uint8{0x2E, 0xAD}, []uint8{0xAD}},
		{[]uint8{0x3E, 0xAD}, []uint8{0xAD}},
	}
	for _, tt := range tests {
		var hi uint8

		state := newState8080()

		state.memory = append(state.memory, tt.in[0])
		state.memory = append(state.memory, tt.in[1])

		Emulate8080(state)

		switch tt.in[0] {
		case 0x06: // MVI B
			hi = state.b
		case 0x0E: // MVI C
			hi = state.c
		case 0x16: // MVI D
			hi = state.d
		case 0x1E: // MVI E
			hi = state.e
		case 0x26: // MVI H
			hi = state.h
		case 0x2E: // MVI L
			hi = state.l
		case 0x3E: // MVI A
			hi = state.a
		default:
			t.Errorf("TestInstructionMVI missing case: %x", tt.in[0])
		}

		if !reflect.DeepEqual(hi, tt.want[0]) {
			t.Errorf("TestInstructionMVI Hi(%x)\nhave %v \nwant %v", tt.in[0], hi, tt.want[0])
		}
	}
}

func TestInstructionDAA(t *testing.T) {
	tests := []struct {
		in   uint8
		want uint8
	}{
		{0x9B, 0x01},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in
		state.memory = append(state.memory, 0x27)
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionDAA(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionANI(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x40, 0x40}, 0x40},
		{[]uint8{0x40, 0x00}, 0x00},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.memory = append(state.memory, 0xe6)
		state.memory = append(state.memory, tt.in[1])
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionANI(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionXRI(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0xFF, 0xFF}, 0x00},
		{[]uint8{0xFF, 0xDD}, 0x22},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.memory = append(state.memory, 0xee)
		state.memory = append(state.memory, tt.in[1])
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionXRI(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionSBI(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x40, 0x40}, 0x00},
		{[]uint8{0x40, 0x20}, 0x20},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.memory = append(state.memory, 0xde)
		state.memory = append(state.memory, tt.in[1])
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionSBI(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionACI(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x00, 0x00}, 0x00},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.memory = append(state.memory, 0xce)
		state.memory = append(state.memory, tt.in[1])
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionACI(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionADI(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x00, 0x00}, 0x00},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.memory = append(state.memory, 0xc6)
		state.memory = append(state.memory, tt.in[1])
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionADI(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestInstructionSUI(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x00, 0x00}, 0x00},
	}
	for _, tt := range tests {
		state := newState8080()
		state.a = tt.in[0]
		state.memory = append(state.memory, 0xd6)
		state.memory = append(state.memory, tt.in[1])
		Emulate8080(state)

		if !reflect.DeepEqual(state.a, tt.want) {
			t.Errorf("TestInstructionSUI(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}
func TestInstructionRST0(t *testing.T) {
	// state := newState8080()
	// state.memory = append(state.memory, 0x00) // NOP
	// state.memory = append(state.memory, 0x03) // INX B
	// state.memory = append(state.memory, 0xc9) // RET
	// state.memory = append(state.memory, 0xc7) // RST 0
	// state.memory = append(state.memory, 0x00) // NOP
	// state.memory = append(state.memory, 0x76) // HLT
	// state.memory = append(state.memory, 0x00) // Stack
	// state.memory = append(state.memory, 0x00) // Stack
	// state.pc = 3                              // Start at RST 0
	// state.sp = 8                              // End of stack
	// Emulate8080(state)
	// Emulate8080(state)
	// Emulate8080(state)
	// Emulate8080(state)
	// Emulate8080(state)
	// Emulate8080(state)
}

func TestInstructionCall(t *testing.T) {
	state := newState8080()
	state.pc = 0
	state.sp = 8                              // end of the stack
	state.memory = append(state.memory, 0xcd) // CALL
	state.memory = append(state.memory, 0x04) // LO ADDR
	state.memory = append(state.memory, 0x00) // HI ADDR
	state.memory = append(state.memory, 0x76) // HLT
	state.memory = append(state.memory, 0x3c) // INR A
	state.memory = append(state.memory, 0xc9) // RET
	state.memory = append(state.memory, 0xFF) // sp
	state.memory = append(state.memory, 0xFF) // sp
	Emulate8080(state)
	Emulate8080(state)
	Emulate8080(state)

	if !reflect.DeepEqual(state.a, uint8(0x01)) {
		t.Errorf("TestInstructionCall(%q)\nhave %v \nwant %v", 0x01, state.a, 0x01)
	}

	if !reflect.DeepEqual(state.memory[state.pc], uint8(0x76)) {
		t.Errorf("TestInstructionCall(%q)\nhave %v \nwant %v", 0x76, state.a, 0x76)
	}

	if !reflect.DeepEqual(state.sp, uint16(8)) {
		t.Errorf("TestInstructionCall(%q)\nhave %v \nwant %v", 0x76, state.a, 0x76)
	}
}
