package chip8

import (
	"fmt"
	"strings"
)

func translateOpCode(high, low byte) string {
	var out string

	first := lNib(high)
	second := rNib(high)
	third := lNib(low)
	fourth := rNib(low)

	_ = second
	_ = third
	_ = fourth

	switch {
	// 00E0
	case high == 0x00 && low == 0xE0:
		out = "CLS"
	// OOEE
	case high == 0x00 && low == 0xEE:
		out = "RET"
	// 1nnn
	case first == 0x1:
		addr := getAddr([]byte{high, low})
		out = fmt.Sprintf("JP 0x%03X", addr)
	default:
		out = "BAD INSTR"
	}
	return out
}

func disassemble(opCodes []byte, offset uint16) strings.Builder {
	var out strings.Builder

	// IF there's an odd number of bytes passed in we ignore the
	// last byte.  To do so we compare i to one less then the
	// lenght of the opCodes slice.
	for i := 0; i < len(opCodes)-1; i += 2 {
		high, low := opCodes[i], opCodes[i+1]
		out.WriteString(fmt.Sprintf("0x%03X   %02X%02X", offset+uint16(i), high, low) + "   " + translateOpCode(high, low) + "\n")
	}
	return out
}
