package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	PROGRAM_OFFSET = 512
	//CLOCK_TICK     = 10
	CLOCK_TICK = 100
	TIMER_TICK = 17
)

func main() {
	fmt.Println("Starting Chip8 Emulator")
	screen := NewScreen()
	defer screen.Close()
	c8 := NewChip8()
	c8.Load("/Users/zach/chip8/particle.ch8")
	c8.Run()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		screen.Update(c8.frameBuffer)
	}
	fmt.Println("Closing Chip8 Emulator")
}
