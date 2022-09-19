package emulator

type conditionCodes struct {
	uint8 z
	uint8 s
	uint8 p
	uint8 cy
	uint8 ac
	uint8 pad
}

type state8080 {
	uint8 a
	uint8 b
	uint8 c
	uint8 d
	uint8 e
	uint8 h
	uint8 l
	uint16 sp
	uint16 pc
	uint8 []memory
	conditionCodes cc
	uint8 int_enable
}


//
func initState8080(){

}
func loadFileToMemoryAt(state state8080, filename string, ){

}