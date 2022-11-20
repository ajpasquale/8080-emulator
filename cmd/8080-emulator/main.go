package main

import (
	"fmt"
	"time"

	emulator "github.com/ajpasquale/8080-emulator"
)

var shiftCount uint8
var shiftLSB uint8 // right
var shiftMSB uint8 // left

func main() {
	cycles := 0
	state := emulator.InitializeState()

	now := time.Now()
	emulator.LoadSpaceInvaders(state)
	timer := now
	for {

		// RST 1 - middle of the screen interrupt
		if time.Since(timer) > 8000*time.Microsecond && emulator.GetIntEnabled(state) == 1 {
			//Interrupt8080(state, )
			emulator.Restart8080(state, emulator.RST1)
		}
		// RST 2 - end of screen interrupt
		if time.Since(timer) > 16000*time.Microsecond && emulator.GetIntEnabled(state) == 1 {
			// problem near 17cd db 02	 IN	 $02
			emulator.Restart8080(state, emulator.RST2)
			timer = time.Now()
		}
		if time.Since(now) > 1*time.Second {
			fmt.Println("break")
		}

		// INPUT
		emulator.GetInput(state)
		// OUTPUT
		emulator.GetOutput(state)

		cycles += emulator.Emulate8080(state)
	}
}
