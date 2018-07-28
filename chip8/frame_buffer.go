package chip8

import "fmt"

type FrameBuffer struct {
	Buffer []uint64
}

func NewFrameBuffer() *FrameBuffer {
	b := make([]uint64, 32, 32)
	return &FrameBuffer{Buffer: b}
}

func (fb *FrameBuffer) clear() {
	for i := range fb.Buffer {
		fb.Buffer[i] = 0
	}
}

func (fb *FrameBuffer) String() string {
	display := ""
	for _, line := range fb.Buffer {
		display = fmt.Sprintf("%s%064b (%X, %d)\n", display, line, line, line)
	}

	return display
}
