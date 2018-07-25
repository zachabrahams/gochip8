package main

import "fmt"

type FrameBuffer struct {
	buffer []uint64
}

func NewFrameBuffer() *FrameBuffer {
	b := make([]uint64, 32, 32)
	return &FrameBuffer{buffer: b}
}

func (fb *FrameBuffer) clear() {
	for i := range fb.buffer {
		fb.buffer[i] = 0
	}
}

func (fb *FrameBuffer) String() string {
	display := ""
	for _, line := range fb.buffer {
		display = fmt.Sprintf("%s%064b (%X, %d)\n", display, line, line, line)
	}

	return display
}
