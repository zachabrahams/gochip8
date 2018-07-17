package main

import (
	"fmt"
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

func main() {
	fmt.Println("Starting Chip8 Emulator")
	screen := NewScreen()
	defer screen.Close()
	c8 := NewChip8(screen)
	c8.Load("/Users/zach/chip8/ibm_logo.ch8")
	c8.Run()
	fmt.Println("Closing Chip8 Emulator")
}
