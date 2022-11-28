package emulator

import (
	"math/bits"
	"reflect"
	"testing"
)

func TestParity(t *testing.T) {
	tests := []struct {
		in   uint8
		want uint8
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
		have := Btoi(bits.OnesCount8(uint8(tt.in&0xFF))%2 == 0)
		if !reflect.DeepEqual(have, tt.want) {
			t.Errorf("TestParity(%q)\nhave %v \nwant %v", tt.in, have, tt.want)
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

		c := comm{}
		state := newState8080(c)
		setArthmeticFlags(state, tt.in)

		if !reflect.DeepEqual(state.cc.cy, tt.want[0]) {
			t.Errorf("TestSetArithmeticFlags(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[0])
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

func TestSetPSW(t *testing.T) {
	tests := []struct {
		in   []uint8
		want uint8
	}{
		// S Z K A - P V C
		{[]uint8{0, 0, 0, 0, 0}, 0x00},
		{[]uint8{0, 0, 0, 0, 1}, 0x01},
		{[]uint8{0, 0, 0, 1, 0}, 0x04},
		{[]uint8{0, 0, 1, 0, 0}, 0x10},
		{[]uint8{0, 1, 0, 0, 0}, 0x40},
		{[]uint8{1, 0, 0, 0, 0}, 0x80},
		{[]uint8{1, 1, 0, 0, 1}, 0xC1},
		{[]uint8{1, 1, 1, 1, 1}, 0xD5},
	}

	for _, tt := range tests {

		c := comm{}
		state := newState8080(c)
		state.cc.s = tt.in[0]
		state.cc.z = tt.in[1]
		state.cc.ac = tt.in[2]
		state.cc.p = tt.in[3]
		state.cc.cy = tt.in[4]
		have := setPSW(state)

		if !reflect.DeepEqual(have, tt.want) {
			t.Errorf("TestSetPSW(%q)\nhave %v \nwant %v", tt.in, have, tt.want)
		}
	}
}

func TestSetFlagsFromPSW(t *testing.T) {
	tests := []struct {
		in   uint8
		want []uint8
	}{
		// S Z K A - P V C
		{0x00, []uint8{0, 0, 0, 0, 0}},
		{0x01, []uint8{0, 0, 0, 0, 1}},
		{0x04, []uint8{0, 0, 0, 1, 0}},
		{0x10, []uint8{0, 0, 1, 0, 0}},
		{0x40, []uint8{0, 1, 0, 0, 0}},
		{0x80, []uint8{1, 0, 0, 0, 0}},
		{0xC1, []uint8{1, 1, 0, 0, 1}},
		{0xD5, []uint8{1, 1, 1, 1, 1}},
	}

	for _, tt := range tests {

		c := comm{}
		state := newState8080(c)
		psw := uint8(tt.in)
		setFlagsFromPSW(state, psw)

		if !reflect.DeepEqual(state.cc.s, tt.want[0]) {
			t.Errorf("TestSetFlagsFromPSW(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[0])
		}
		if !reflect.DeepEqual(state.cc.z, tt.want[1]) {
			t.Errorf("TestSetFlagsFromPSW(%q)\nhave %v \nwant %v", tt.in, state.cc.z, tt.want[1])
		}
		if !reflect.DeepEqual(state.cc.ac, tt.want[2]) {
			t.Errorf("TestSetFlagsFromPSW(%q)\nhave %v \nwant %v", tt.in, state.cc.p, tt.want[2])
		}
		if !reflect.DeepEqual(state.cc.p, tt.want[3]) {
			t.Errorf("TestSetFlagsFromPSW(%q)\nhave %v \nwant %v", tt.in, state.cc.s, tt.want[3])
		}
		if !reflect.DeepEqual(state.cc.cy, tt.want[4]) {
			t.Errorf("TestSetFlagsFromPSW(%q)\nhave %v \nwant %v", tt.in, state.cc.s, tt.want[4])
		}
	}
}

func TestAdd16(t *testing.T) {
	tests := []struct {
		in   []uint16
		want []uint16
	}{
		//         a        b              (a+b) carry flag
		{[]uint16{0x339F, 0xA17B}, []uint16{0xD51A, 0}},
		{[]uint16{0x0000, 0x0000}, []uint16{0x0000, 0}},
		{[]uint16{0xFFFF, 0x0001}, []uint16{0x0000, 1}},
	}
	for _, tt := range tests {
		have, cb := add16(tt.in[0], tt.in[1])

		if !reflect.DeepEqual(have, tt.want[0]) {
			t.Errorf("TestAdd16 value(%q)\nhave %v \nwant %v", tt.in, have, tt.want[0])
		}
		carry := Btoi(cb)
		if !reflect.DeepEqual(uint16(carry), tt.want[1]) {
			t.Errorf("TestAdd16 carry(%x)\nhave %v \nwant %v", tt.in, carry, tt.want[1])
		}
	}
}

func TestSub8Flags(t *testing.T) {
	// tests := []struct {
	// 	in   []uint8
	// 	want []uint8
	// }{
	// 	// S Z A P C
	// 	// {[]uint8{0x3E, 0x3E}, []uint8{0, 1, 1, 1, 0}},
	// 	// {[]uint8{0x4A, 0x40, 0}, []uint8{0, 0, 0, 1, 0}},

	// 	//;A=F5H,C=0,Z=0
	// 	{[]uint8{0xF5, 0x00, 0}, []uint8{1, 0, 0, 1, 0}},
	// }
	// for _, tt := range tests {
	// 	c := comm{}
	// 	state := newState8080(c)
	// 	state.a = tt.in[0]
	// 	state.b = tt.in[1]

	// 	sub8(state, state.a, state.b)
	// 	if !reflect.DeepEqual(state.cc.s, tt.want[0]) {
	// 		t.Errorf("TestSub8Flags(%q)\nhave %v \nwant %v", tt.in, state.cc.s, tt.want[0])
	// 	}
	// 	if !reflect.DeepEqual(state.cc.z, tt.want[1]) {
	// 		t.Errorf("TestSub8Flags(%q)\nhave %v \nwant %v", tt.in, state.cc.z, tt.want[1])
	// 	}
	// 	if !reflect.DeepEqual(state.cc.ac, tt.want[2]) {
	// 		t.Errorf("TestSub8Flags(%q)\nhave %v \nwant %v", tt.in, state.cc.ac, tt.want[2])
	// 	}
	// 	if !reflect.DeepEqual(state.cc.p, tt.want[3]) {
	// 		t.Errorf("TestSub8Flags(%q)\nhave %v \nwant %v", tt.in, state.cc.p, tt.want[3])
	// 	}
	// 	if !reflect.DeepEqual(state.cc.cy, tt.want[4]) {
	// 		t.Errorf("TestSub8Flags(%q)\nhave %v \nwant %v", tt.in, state.cc.cy, tt.want[4])
	// 	}
	//}

}
