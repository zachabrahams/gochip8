package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type Chip8 struct {
	callStack   []uint16
	frameBuffer *FrameBuffer
	delayTimer  *Timer
	keyboard    *Keyboard
	memory      []byte
	programPtr  uint16
	registers   map[byte]byte
	regI        uint16
	step        chan bool
}

func NewChip8(kb *Keyboard) *Chip8 {
	m := make([]byte, 4096, 4096)
	loadBuiltInSprites(m)
	r := map[byte]byte{}
	// Registers 0-15 are for V0-VF.
	for i := 0; i < 16; i++ {
		r[byte(i)] = byte(0)
	}

	return &Chip8{
		callStack:   []uint16{},
		frameBuffer: NewFrameBuffer(),
		delayTimer:  NewTimer(),
		keyboard:    kb,
		memory:      m,
		programPtr:  PROGRAM_OFFSET,
		registers:   r,
		regI:        0,
		step:        make(chan bool),
	}
}

func (c8 *Chip8) String() {
	var msg bytes.Buffer
	// msg.WriteString(hex.Dump(c8.memory))
	c8.frameBuffer.bitDump()
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

func (c8 *Chip8) Run() {
	ticker := time.NewTicker(CLOCK_TICK * time.Millisecond)
	tick := 0
	go func() {
		for _ = range ticker.C {
			c8.execInstr()
			tick++
			// clearScreen()
			c8.String()
			//<-c8.step
		}
	}()
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
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
