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
	display   *Display
	memory    []byte
	registers map[byte]byte
	regI      uint16
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
		display:   NewDisplay(),
		memory:    m,
		registers: r,
		regI:      0,
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
			tick++
			clearScreen()
			c8.display.bitDump()
			fmt.Println(tick)
			if tick >= 300 {
				wg.Done()
				ticker.Stop()
			}
		}
	}()

	wg.Wait()
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
