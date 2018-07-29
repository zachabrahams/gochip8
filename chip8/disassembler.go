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
	// 2nnn
	case first == 0x2:
		addr := getAddr([]byte{high, low})
		out = fmt.Sprintf("CALL 0x%03X", addr)
	// 3xkk
	case first == 0x3:
		out = "SE" + fmtVxByte(second, low)
	// 4xkk
	case first == 0x4:
		out = "SNE" + fmtVxByte(second, low)
	// 5xy0
	case first == 0x5:
		out = "SE" + fmtVxVy(second, third)
	// 6xkk
	case first == 0x6:
		out = "LD" + fmtVxByte(second, low)
	default:
		out = "BAD INSTR"
	}
	return out
}

func fmtVxVy(x, y byte) string {
	return fmt.Sprintf(" V%X, V%X", x, y)
}

func fmtVxByte(x, b byte) string {
	return fmt.Sprintf(" V%X, 0x%02X", x, b)
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
