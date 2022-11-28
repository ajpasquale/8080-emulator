package emulator

import (
	"fmt"
	"testing"
)

type pcomm struct{}

var finished bool

func (p pcomm) PortIn(state *state8080, port uint8) uint8 {
	return 0x00
}
func (p pcomm) PortOut(state *state8080, port uint8, value uint8) {

	if port == 0 {
		finished = true
	} else if port == 1 {
		operation := state.c

		if operation == 2 {
			fmt.Printf("%c", state.e)
		} else if operation == 9 {
			addr := bytesToPair(state.d, state.e)
			for state.memory[addr] != '$' {
				fmt.Printf("%c", state.memory[addr])
				addr++
			}
		}
	}
}

func TestTST8080(t *testing.T) {
	state := newState8080(pcomm{})
	for i := 0; i < 0x100; i++ {
		state.memory = append(state.memory, 0xFF)
	}
	LoadFileIntoMemoryAt(state, "rom/tests/TST8080.COM", 0x100)
	for i := 0; i < 0x10000; i++ {
		state.memory = append(state.memory, 0xFF)
	}
	fmt.Printf("*** TEST: %s\n", "TST8080.COM")

	state.pc = 0x100

	state.memory[0x0000] = 0xD3
	state.memory[0x0001] = 0x00

	state.memory[0x0005] = 0xD3
	state.memory[0x0006] = 0x01
	state.memory[0x0007] = 0xC9

	finished = false

	for !finished {
		Emulate8080(state)
	}
}
