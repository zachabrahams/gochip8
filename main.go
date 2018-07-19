package main

import (
	"fmt"
)

const (
	PROGRAM_OFFSET = 512
	TICK_TIME      = 17
)

func main() {
	fmt.Println("Starting Chip8 Emulator")
	screen := NewScreen()
	defer screen.Close()
	c8 := NewChip8(screen)
	c8.Load("/Users/zach/chip8/jumping_x_o.ch8")
	c8.Run()
	fmt.Println("Closing Chip8 Emulator")
}
