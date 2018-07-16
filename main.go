package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	PROGRAM_OFFSET = 512
	TICK_TIME      = 17
	x0             = 0
	x1             = 1
	x2             = 2
	x3             = 3
	x4             = 4
	x5             = 5
	x6             = 6
	x7             = 7
	x8             = 8
	x9             = 9
	xA             = 10
	xB             = 11
	xC             = 12
	xD             = 13
	xE             = 14
	xF             = 15
)

type Chip8 struct {
	display    *Display
	memory     []byte
	programPtr uint16
	registers  map[byte]byte
	regI       uint16
}

type Display struct {
	screen [][]byte
}

func NewDisplay() *Display {
	var d [][]byte
	for i := 0; i < 32; i++ {
		line := make([]byte, 8, 8)
		d = append(d, line)
	}
	return &Display{screen: d}
}

func (d *Display) clear() {
	for i, line := range d.screen {
		for j := range line {
			d.screen[i][j] = 0
		}
	}
}

func (d *Display) bitDump() {
	for _, line := range d.screen {
		lineStr := ""
		for _, byt := range line {
			lineStr = fmt.Sprintf("%s%08b", lineStr, byt)
		}
		fmt.Println(lineStr)
	}
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
			c8.display.bitDump()
			fmt.Printf("Tick: %d\n", tick)
		}
	}()

	wg.Wait()
}

func (c8 *Chip8) execInstr() {
	instr := c8.memory[c8.programPtr : c8.programPtr+2]
	fmt.Printf("Instr: %X %v\n", instr, instr)

	switch {
	// 00E0 - CLS - clear the display
	case instr[0] == 0 && instr[1] == 224:
		c8.display.clear()
	// 6xkk - LD Vx, byte - Load the byte value into the register specified by x
	case lNib(instr[0]) == x6:
		c8.registers[rNib(instr[0])] = instr[1]
	// Annn - LD I, addr - Load the int16 addr specified by nnn into the I register
	case lNib(instr[0]) == xA:
		addr := getAddr(instr[0], instr[1])
		c8.regI = addr
	default:
		msg := fmt.Sprintf("Unknown Instruction: %X\n", instr)
		panic(msg)
	}

	c8.programPtr = c8.programPtr + 2
}

func main() {
	fmt.Println("Starting Chip8 Emulator")
	c8 := NewChip8()
	c8.Load("/Users/zach/chip8/ibm_logo.ch8")
	c8.Run()
	fmt.Println("Closing Chip8 Emulator")
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func lNib(b byte) byte {
	return b >> 4
}

func rNib(b byte) byte {
	return b & 15
}

func getAddr(b1, b2 byte) uint16 {
	right := uint16(rNib(b1)) << 8
	left := uint16(b2)
	return right + left
}
