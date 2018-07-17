package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Chip8 struct {
	display    *Display
	memory     []byte
	programPtr uint16
	registers  map[byte]byte
	regI       uint16
}

func NewChip8() *Chip8 {
	m := make([]byte, 4096, 4096)
	r := map[byte]byte{}
	// Registers 0-15 are for V0-VF.
	for i := 0; i < 16; i++ {
		r[byte(i)] = byte(0)
	}

	return &Chip8{
		display:    NewDisplay(),
		memory:     m,
		programPtr: PROGRAM_OFFSET,
		registers:  r,
		regI:       0,
	}
}

func (c8 *Chip8) String() {
	var msg bytes.Buffer
	// msg.WriteString(hex.dump(c8.memory)
	c8.display.bitDump()
	msg.WriteString(fmt.Sprintf("Program Counter: %X (%d)\n", c8.programPtr, c8.programPtr))

	instr := c8.memory[c8.programPtr : c8.programPtr+2]
	msg.WriteString(fmt.Sprintf("Instr: %X %v\n", instr, instr))
	msg.WriteString("Registers:\n")
	for i := 0; i < 16; i++ {
		msg.WriteString(fmt.Sprintf("V%X: %02X (%d)\n", i, c8.registers[byte(i)], c8.registers[byte(i)]))
	}
	msg.WriteString(fmt.Sprintf("I: %03X (%d)\n", c8.regI, c8.regI))
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
	var wg sync.WaitGroup
	wg.Add(1)
	ticker := time.NewTicker(TICK_TIME * time.Millisecond)
	tick := 0
	go func() {
		for _ = range ticker.C {
			c8.execInstr()
			tick++
			clearScreen()
			//c8.String()
			c8.display.bitDump()
		}
	}()

	wg.Wait()
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
