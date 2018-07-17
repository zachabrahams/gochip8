package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const SCALING_FACTOR = 10

type Screen struct {
	window *sdl.Window
}

func NewScreen() *Screen {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("gochip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 64*SCALING_FACTOR, 32*SCALING_FACTOR, sdl.WINDOW_SHOWN|sdl.WINDOW_BORDERLESS)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	window.UpdateSurface()
	return &Screen{window: window}
}

func (s *Screen) Update(d *Display) {
	rects := []sdl.Rect{}
	for i, line := range d.screen {
		for j := 0; j < 64; j++ {
			var bit uint64 = 1 << uint(63-j)
			if (line & bit) > 0 {
				rect := &sdl.Rect{
					X: int32(j * SCALING_FACTOR),
					Y: int32(i * SCALING_FACTOR),
					W: SCALING_FACTOR,
					H: SCALING_FACTOR,
				}
				rects = append(rects, *rect)
			}
		}

	}
	surface, err := s.window.GetSurface()
	if err != nil {
		panic(err)
	}

	if len(rects) > 0 {
		surface.FillRect(nil, 0)
		surface.FillRects(rects, 0xffffffff)
		s.window.UpdateSurface()
	}

}

func (s *Screen) Close() {
	s.window.Destroy()
	sdl.Quit()
}
