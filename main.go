package main

import (
    "fmt"
)

type Chip8 struct {
    memory []byte
	registers map[byte]byte
}

func NewChip8() *Chip8 {
	m := make([]byte, 4096, 4096)
	r := map[byte]byte{}
	// Registers 0-15 are for V0-VF. 16 is I.
	for i := 0; i < 17; i++ {
		r[byte(i)] = byte(0)
	}

	return &Chip8{
		memory: m,
		registers: r,
	}
}

func (*Chip8) LoadProgram (filename string) {
	fmt.Printf("Loading Program From File: %s\n", filename)
}

func main() {
	fmt.Println("Starting Chip8 Emulator")
	chip8 := NewChip8()
	chip8.LoadProgram("test.com")
	fmt.Println("Closing Chip8 Emulator")
}

