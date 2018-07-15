package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	PROGRAM_OFFSET = 512
)

type Chip8 struct {
	memory    []byte
	registers map[byte]byte
	regI      uint16
}

func NewChip8() *Chip8 {
	m := make([]byte, 4096, 4096)
	r := map[byte]byte{}
	// Registers 0-15 are for V0-VF.
	for i := 0; i < 16; i++ {
		r[byte(i)] = byte(0)
	}

	return &Chip8{
		memory:    m,
		registers: r,
		regI:      0,
	}
}

func (c8 *Chip8) LoadProgram(filename string) {
	fmt.Printf("Loading Program From File: %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("could not load program file: %v", err)
	}
	defer file.Close()
	binData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("error reading data from file: %v")
	}

	for i := 0; i < len(binData); i++ {
		c8.memory[PROGRAM_OFFSET+i] = binData[i]
	}
	fmt.Printf("Finshed loading program. Loaded %d bytes\n", len(binData))
}

func main() {
	fmt.Println("Starting Chip8 Emulator")
	chip8 := NewChip8()
	chip8.LoadProgram("/Users/zach/chip8/ibm_logo.ch8")
	fmt.Println("Closing Chip8 Emulator")
}
