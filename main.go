package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	PROGRAM_OFFSET = 512
	CLOCK_TICK     = 3
	// CLOCK_TICK = 500
	TIMER_TICK = 17
)

func main() {
	fmt.Println("Starting Chip8 Emulator")
	screen := NewScreen()
	defer screen.Close()
	kb := NewKeyboard()
	c8 := NewChip8(kb)
	c8.Load("/Users/zach/chip8/space_intercept-1978.ch8")
	c8.Run()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				kevent := event.(*sdl.KeyboardEvent)
				if kevent.Type == sdl.KEYUP && kevent.Keysym.Sym == sdl.K_PERIOD {
					//	c8.step <- true
				}
			}
		}
		kbState := sdl.GetKeyboardState()
		newKBState := parseKbState(kbState)

		kb.newKeyboardState <- newKBState
		screen.Update(c8.frameBuffer)
	}
	fmt.Println("Closing Chip8 Emulator")
}

func parseKbState(kbState []uint8) uint16 {
	var keys uint16
	keys = 0
	if kbState[sdl.SCANCODE_1] == 1 {
		keys = keys | (0x1 << 0)
	}
	if kbState[sdl.SCANCODE_2] == 1 {
		keys = keys | (0x1 << 1)
	}
	if kbState[sdl.SCANCODE_3] == 1 {
		keys = keys | (0x1 << 2)
	}
	if kbState[sdl.SCANCODE_4] == 1 {
		keys = keys | (0x1 << 3)
	}
	if kbState[sdl.SCANCODE_Q] == 1 {
		keys = keys | (0x1 << 4)
	}
	if kbState[sdl.SCANCODE_W] == 1 {
		keys = keys | (0x1 << 5)
	}
	if kbState[sdl.SCANCODE_E] == 1 {
		keys = keys | (0x1 << 6)
	}
	if kbState[sdl.SCANCODE_R] == 1 {
		keys = keys | (0x1 << 7)
	}
	if kbState[sdl.SCANCODE_A] == 1 {
		keys = keys | (0x1 << 8)
	}
	if kbState[sdl.SCANCODE_S] == 1 {
		keys = keys | (0x1 << 9)
	}
	if kbState[sdl.SCANCODE_D] == 1 {
		keys = keys | (0x1 << 10)
	}
	if kbState[sdl.SCANCODE_F] == 1 {
		keys = keys | (0x1 << 11)
	}
	if kbState[sdl.SCANCODE_Z] == 1 {
		keys = keys | (0x1 << 12)
	}
	if kbState[sdl.SCANCODE_X] == 1 {
		keys = keys | (0x1 << 13)
	}
	if kbState[sdl.SCANCODE_C] == 1 {
		keys = keys | (0x1 << 14)
	}
	if kbState[sdl.SCANCODE_V] == 1 {
		keys = keys | (0x1 << 15)
	}
	return keys
}
