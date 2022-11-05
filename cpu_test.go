package emulator

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCpu(t *testing.T) {
	state := newState8080()
	LoadSpaceInvaders(state)

	now := time.Now()

	timer := now
	for {

		// RST 1 - middle of the screen interrupt
		if time.Since(timer) > 8000*time.Microsecond && state.int_enable == 1 {
			//Interrupt8080(state, )
			Restart8080(state, RST1)
		}
		// RST 2 - end of screen interrupt
		if time.Since(timer) > 16000*time.Microsecond && state.int_enable == 1 {
			// problem near 17cd db 02	 IN	 $02
			Restart8080(state, RST2)
			timer = time.Now()
		}
		if time.Since(now) > 1*time.Second || state.pc == 0x024b {
			fmt.Println("break")
		}
		fmt.Printf("pc: %x, a: %x, h: %x, l: %x\n",
			state.pc,
			state.a,
			state.h,
			state.l)
		Emulate8080(state)
	}
}

func TestInstructionINXB(t *testing.T) {
	tests := []struct {
		in   []uint8
		want []uint8
	}{
		{[]uint8{0x00, 0x00}, []uint8{0x00, 0x01}},
		{[]uint8{0xAB, 0xCD}, []uint8{0xAB, 0xCE}},
		{[]uint8{0xFF, 0xFF}, []uint8{0x00, 0x00}},
	}

	for _, tt := range tests {
		state := newState8080()
		state.b = tt.in[0]
		state.c = tt.in[1]
		state.memory = append(state.memory, 0x03)
		Emulate8080(state)
		if !reflect.DeepEqual(state.b, tt.want[0]) {
			t.Errorf("TestInstructionINXB(%q)\nhave %v \nwant %v", tt.in, state.b, tt.want[0])
		}
		if !reflect.DeepEqual(state.c, tt.want[1]) {
			t.Errorf("TestInstructionINXB(%q)\nhave %v \nwant %v", tt.in, state.c, tt.want[1])
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

func TestInstructionDADB(t *testing.T) {
	tests := []struct {
		in   []uint16
		want []uint16
	}{
		//         H,L   +   B,C   =   HL
		{[]uint16{0x2061, 0x4050}, []uint16{0x60B1, 1}},
	}
	for _, tt := range tests {
		state := newState8080()
		state.h, state.l = pairToBytes(tt.in[0])
		state.b, state.c = pairToBytes(tt.in[1])
		state.memory = append(state.memory, 0x09)
		Emulate8080(state)

		if !reflect.DeepEqual(bytesToPair(state.h, state.l), tt.want[0]) {
			t.Errorf("TestInstructionDADB(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want[0])
		}
		if !reflect.DeepEqual(state.cc.cy, uint8(tt.want[1])) {
			t.Errorf("TestInstructionDADB(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[1])
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
