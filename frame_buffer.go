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

func (fb *FrameBuffer) bitDump() {
	for _, line := range fb.buffer {
		fmt.Printf("%064b\n", line)
	}
}
