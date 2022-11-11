package main

func main() {
	// cycles := 0
	// state := emulator.newState8080()
	// LoadSpaceInvaders(state)

	// now := time.Now()

	// timer := now
	// for {

	// 	// RST 1 - middle of the screen interrupt
	// 	if time.Since(timer) > 8000*time.Microsecond && state.int_enable == 1 {
	// 		//Interrupt8080(state, )
	// 		Restart8080(state, RST1)
	// 	}
	// 	// RST 2 - end of screen interrupt
	// 	if time.Since(timer) > 16000*time.Microsecond && state.int_enable == 1 {
	// 		// problem near 17cd db 02	 IN	 $02
	// 		Restart8080(state, RST2)
	// 		timer = time.Now()
	// 	}
	// 	if time.Since(now) > 1*time.Second || state.pc == 0x024b {
	// 		fmt.Println("break")
	// 	}

	// 	// INPUT
	// 	if state.memory[state.pc] == 0xdb {
	// 		switch state.memory[state.pc+1] {
	// 		case 0x0: // fire, left, right?
	// 		case 0x1: // credit,start, player 1 shot, left, right
	// 		case 0x2: // dip3,5,6, player 2 shot, left, right
	// 		case 0x3: // shift reg data
	// 			m := uint16(shiftMSB) << 8
	// 			shift := uint16(m | uint16(shiftLSB))
	// 			state.a = uint8((shift >> (8 - shiftCount)) & 0xFF)
	// 		}

	// 	}
	// 	// OUTPUT
	// 	if state.memory[state.pc] == 0xd3 {
	// 		switch state.memory[state.pc+1] {
	// 		case 0x02: // shift amount
	// 			shiftCount = state.a & 7
	// 		case 0x03: // discrete sounds
	// 		case 0x04: // shift data (LSB on 1st write, MSB on 2nd)
	// 			shiftLSB = shiftMSB
	// 			shiftMSB = state.a
	// 		case 0x05: // discrete sounds
	// 		case 0x06: // watchdog?
	// 		}

	// 	}
	// 	fmt.Printf("pc: %x, a: %x, h: %x, l: %x\n",
	// 		state.pc,
	// 		state.a,
	// 		state.h,
	// 		state.l)
	// 	cycles += Emulate8080(state)
	// }
}
