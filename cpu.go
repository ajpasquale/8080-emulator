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


func newState8080(){
	cc := &conditionCodes{
		z:0, 
		s:0,
		p:0,
		cy:0,
		ac:0,
		pad:0,

	}
	
	return &state8080{
		a: 0,
		b: 0,
		c: 0,
		d: 0,
		e: 0,
		h: 0,
		l: 0,
		sp: 0,
		pc: 0,
		memory: make(uint8[], 16384, 16384),
		cc: cc,
		int_enable:0
	}
}