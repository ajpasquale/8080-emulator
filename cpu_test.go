package emulator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCpu(t *testing.T) {

	state := newState8080()

	loadFileIntoMemoryAt(state, "rom/invaders/invaders.h", 0)
	loadFileIntoMemoryAt(state, "rom/invaders/invaders.g", 0x800)
	loadFileIntoMemoryAt(state, "rom/invaders/invaders.f", 0x1000)
	loadFileIntoMemoryAt(state, "rom/invaders/invaders.e", 0x1800)
}

func TestParity(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{0, 1},
		{1, 0},
		{2, 0},
		{3, 1},
		{4, 0},
		{5, 1},
		{6, 1},
		{7, 0},
		{8, 0},
		{9, 1},
		{10, 1},
		{11, 0},
		{12, 1},
		{13, 0},
		{14, 0},
		{15, 1},
		{16, 0},
		{17, 1},
		{18, 1},
		{19, 0},
		{20, 1},
		{21, 0},
		{22, 0},
		{23, 1},
		{24, 1},
		{25, 0},
		{26, 0},
		{27, 1},
		{28, 0},
		{29, 1},
		{30, 1},
		{31, 0},
		{32, 0},
		{33, 1},
		{34, 1},
		{35, 0},
		{36, 1},
		{37, 0},
		{38, 0},
		{39, 1},
		{40, 1},
		{41, 0},
		{42, 0},
		{43, 1},
		{44, 0},
		{45, 1},
		{46, 1},
		{47, 0},
		{48, 1},
		{49, 0},
		{50, 0},
		{51, 1},
		{52, 0},
		{53, 1},
		{54, 1},
		{55, 0},
		{56, 0},
		{57, 1},
		{58, 1},
		{59, 0},
		{60, 1},
		{61, 0},
		{62, 0},
		{63, 1},
		{64, 0},
		{65, 1},
		{66, 1},
		{67, 0},
		{68, 1},
		{69, 0},
		{70, 0},
		{71, 1},
		{72, 1},
		{73, 0},
		{74, 0},
		{75, 1},
		{76, 0},
		{77, 1},
		{78, 1},
		{79, 0},
		{80, 1},
		{81, 0},
		{82, 0},
		{83, 1},
		{84, 0},
		{85, 1},
		{86, 1},
		{87, 0},
		{88, 0},
		{89, 1},
		{90, 1},
		{91, 0},
		{92, 1},
		{93, 0},
		{94, 0},
		{95, 1},
		{96, 1},
		{97, 0},
		{98, 0},
		{99, 1},
	}
	for _, tt := range tests {
		have := parity(tt.in, 8)
		if !reflect.DeepEqual(have, tt.want) {
			t.Errorf("parity(%q)\nhave %v \nwant %v", tt.in, have, tt.want)
		}
	}
}
func TestSetArithmeticFlags(t *testing.T) {
	tests := []struct {
		in   uint16
		want []uint8
	}{
		//Cy=0, AC=0, Z=0, P=0, S=0
		{0x00, []uint8{0, 0, 1, 1, 0}},
		{0x01, []uint8{0, 0, 0, 0, 0}},
		{0xF, []uint8{0, 0, 0, 1, 0}},
		{0x10, []uint8{0, 1, 0, 0, 0}},
		{0xB9, []uint8{0, 1, 0, 0, 1}},
		{0xFF, []uint8{0, 1, 0, 1, 1}},
		{0x100, []uint8{1, 1, 1, 1, 0}},
		{0x101, []uint8{1, 1, 0, 0, 0}},
		{0xFFF, []uint8{1, 1, 0, 1, 1}},
		{0xFFFF, []uint8{1, 1, 0, 1, 1}},
	}

	for _, tt := range tests {

		state := newState8080()
		setArthmeticFlags(state, tt.in)

		if !reflect.DeepEqual(state.cc.cy, tt.want[0]) {
			t.Errorf("TestSetArithmeticFlags(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[0])
		}
		if !reflect.DeepEqual(state.cc.ac, tt.want[1]) {
			t.Errorf("TestSetArithmeticFlags(%q)\nhave %v \nwant %v", tt.in, state.cc.ac, tt.want[1])
		}
		if !reflect.DeepEqual(state.cc.z, tt.want[2]) {
			t.Errorf("TestSetArithmeticFlags(%q)\nhave %v \nwant %v", tt.in, state.cc.z, tt.want[2])
		}
		if !reflect.DeepEqual(state.cc.p, tt.want[3]) {
			t.Errorf("TestSetArithmeticFlags(%q)\nhave %v \nwant %v", tt.in, state.cc.p, tt.want[3])
		}
		if !reflect.DeepEqual(state.cc.s, tt.want[4]) {
			t.Errorf("TestSetArithmeticFlags(%q)\nhave %v \nwant %v", tt.in, state.cc.s, tt.want[4])
		}
	}

}

func TestInstructionINXB(t *testing.T) {
	state := newState8080()
	state.b = 0xAB
	state.c = 0xCD
	fmt.Printf("%x%x\n", state.b, state.c)
	state.memory = append(state.memory, 0x03)
	Emulate8080(state)
	fmt.Printf("%x%x\n", state.b, state.c)

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
			t.Errorf("TestInstructionRLC(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want[0])
		}
		if !reflect.DeepEqual(state.cc.cy, uint8(tt.want[1])) {
			t.Errorf("TestInstructionRLC(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[1])
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
			t.Errorf("TestInstructionADDB(%q)\nhave %v \nwant %v", tt.in, state.a, tt.want)
		}
	}
}

func TestSetAuxCarry(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		{[]uint8{0x00, 0x00}, 0},
		{[]uint8{0x16, 0x01}, 0},
		{[]uint8{0x16, 0x08}, 0},
		{[]uint8{0x0F, 0x01}, 1},
		{[]uint8{0x3D, 0x42}, 0},
		{[]uint8{0x3D, 0x43}, 1},
	}
	for _, tt := range tests {
		have := setAuxCarry(tt.in[0], tt.in[1])
		if !reflect.DeepEqual(have, tt.want) {
			t.Errorf("TestSetAuxCarry(%q)\nhave %v \nwant %v", tt.in, have, tt.want)
		}
	}
}
