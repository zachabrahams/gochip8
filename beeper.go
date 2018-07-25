package main

import (
	"fmt"
	"io/ioutil"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

type Beeper interface {
	Beep()
	Close()
}

type SDLBeeper struct {
	data []byte
}

// For now assume that sld is initialized, I'll update that later.
func NewSDLBeeper() *SDLBeeper {
	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 4, 4096); err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("./beep.wav")
	if err != nil {
		panic(err)
	}

	return &SDLBeeper{data: data}
}

func (b *SDLBeeper) Beep() {
	chunk, err := mix.QuickLoadWAV(b.data)
	defer chunk.Free()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	_, err = chunk.Play(-1, 0)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	for mix.Playing(-1) == 1 {
		sdl.Delay(16)
	}
}

func (b *SDLBeeper) Close() {
	mix.CloseAudio()
}
