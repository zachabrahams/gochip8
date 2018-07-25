package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Chip8 is the struct that represents a full Chip8 VM
// The attribues are:
//
// beepTimer: A timer that counts down at 60hz and beeps when it reaches 0
//
// callStack: A stack of addresses to return to from subroutines
//
// deplayTimer: A timer that counts down at 60 hz
//
// frameBuffer: A representation of the current state of the screen
//
// keyboard: A representation of the current state of the keyboard
//
// memory: a 4kb byte slice reprsenting the memory available to the system
//
// programPtr: the register that points to the next instruction to run
//
// regI: the 16 bit I register - used for storing the location of sprites
//
// registers: An array of the 16 8-bit registesr used by the CPU.
// They are named V0-VF.
//
// step: a channel for doing hacky debugging - should be refactored away.
type Chip8 struct {
	beepTimer   *Timer
	callStack   []uint16
	delayTimer  *Timer
	frameBuffer *FrameBuffer
	keyboard    *Keyboard
	memory      []byte
	programPtr  uint16
	regI        uint16
	registers   map[byte]byte
	step        chan bool
}

// NewChip8 accepts a keyboard and a beeper and returns a pointer to a full
// Chip8.
func NewChip8(kb *Keyboard, b Beeper) *Chip8 {
	m := make([]byte, 4096, 4096)
	loadBuiltInSprites(m)
	r := map[byte]byte{}
	for i := 0; i < 16; i++ {
		r[byte(i)] = byte(0)
	}

	return &Chip8{
		beepTimer:   NewTimer(func() { b.Beep() }),
		callStack:   []uint16{},
		delayTimer:  NewTimer(func() {}),
		frameBuffer: NewFrameBuffer(),
		keyboard:    kb,
		memory:      m,
		programPtr:  PROGRAM_OFFSET,
		regI:        0,
		registers:   r,
		step:        make(chan bool),
	}
}

// String provides a text representation fof the current state of the Chip8.
func (c8 *Chip8) String() {
	var msg bytes.Buffer
	// Uncomment the following to get a hex dump of the entire memory stack
	// msg.WriteString(hex.Dump(c8.memory))
	msg.WriteString(c8.frameBuffer.String())
	msg.WriteString(fmt.Sprintf("Program Counter: %X (%d)\n", c8.programPtr, c8.programPtr))

	instr := c8.memory[c8.programPtr : c8.programPtr+2]
	msg.WriteString(fmt.Sprintf("Instr: %X %v\n", instr, instr))
	msg.WriteString("Registers:\n")
	for i := 0; i < 16; i++ {
		msg.WriteString(fmt.Sprintf("V%X: %02X (%d)\n", i, c8.registers[byte(i)], c8.registers[byte(i)]))
	}
	msg.WriteString(fmt.Sprintf("I: %03X (%d)\n", c8.regI, c8.regI))
	msg.WriteString(fmt.Sprintf("Call Stack: %v", c8.callStack))
	fmt.Println(msg.String())
}

// Load is used to load a Chip8 program into memory from a system file.
func (c8 *Chip8) Load(filename string) {
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

// Run creates a ticker using the CLOCK_TICK variable and executes an instruction on every tick.
func (c8 *Chip8) Run() {
	ticker := time.NewTicker(CLOCK_TICK * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			c8.execInstr()
			//c8.String()
			//<-c8.step
		}
	}()
}

func loadBuiltInSprites(m []byte) {
	sprites := [][]byte{
		[]byte{0xF0, 0x90, 0x90, 0x90, 0xF0}, // 0
		[]byte{0x20, 0x60, 0x20, 0x20, 0x70}, // 1
		[]byte{0xF0, 0x10, 0xF0, 0x80, 0xF0}, // 2
		[]byte{0xF0, 0x10, 0xF0, 0x10, 0xF0}, // 3
		[]byte{0x90, 0x90, 0xF0, 0x10, 0x10}, // 4
		[]byte{0xF0, 0x80, 0xF0, 0x10, 0xF0}, // 5
		[]byte{0xF0, 0x80, 0xF0, 0x90, 0xF0}, // 6
		[]byte{0xF0, 0x10, 0x20, 0x40, 0x40}, // 7
		[]byte{0xF0, 0x90, 0xF0, 0x90, 0xF0}, // 8
		[]byte{0xF0, 0x90, 0xF0, 0x10, 0xF0}, // 9
		[]byte{0xF0, 0x90, 0xF0, 0x90, 0x90}, // A
		[]byte{0xE0, 0x90, 0xE0, 0x90, 0xE0}, // B
		[]byte{0xF0, 0x80, 0x80, 0x80, 0xF0}, // C
		[]byte{0xE0, 0x90, 0x90, 0x90, 0xE0}, // D
		[]byte{0xF0, 0x80, 0xF0, 0x80, 0xF0}, // E
		[]byte{0xF0, 0x80, 0xF0, 0x80, 0x80}, // F
	}

	// We load the sprites starting at memory 0x000, for easy fetching
	for i, sprite := range sprites {
		for j, line := range sprite {
			m[(i*5)+j] = line
		}
	}
}
