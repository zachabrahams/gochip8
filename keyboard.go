package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Keyboard struct {
	checkKey   chan int
	keyPressed chan bool

	getKey     chan bool
	receiveKey chan int
}

func NewKeyboard() *Keyboard {
	var checkKey chan int
	var keyPressed chan bool

	var getKey chan bool
	var receiveKey chan int

	go func() {
		for {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					println("Quit")
					break
				}
			}
		}
	}()
	return &Keyboard{
		checkKey:   checkKey,
		keyPressed: keyPressed,
		getKey:     getKey,
		receiveKey: receiveKey,
	}
}
