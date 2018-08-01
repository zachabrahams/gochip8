package chip8

import (
	"fmt"
	"strings"
)

func translateOpCode(instr []byte) string {
	var out string

	high := instr[0]
	low := instr[1]
	first := lNib(high)
	fourth := rNib(low)

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
		out = fmtVxByte("SE", instr)
	// 4xkk
	case first == 0x4:
		out = fmtVxByte("SNE", instr)
	// 5xy0
	case first == 0x5:
		out = fmtVxVy("SE", instr)
	// 6xkk
	case first == 0x6:
		out = fmtVxByte("LD", instr)
	// 7xkk
	case first == 0x7:
		out = fmtVxByte("ADD", instr)
	// 8xy0
	case first == 0x8 && fourth == 0x0:
		out = fmtVxVy("LD", instr)
	// 8xy1
	case first == 0x8 && fourth == 0x1:
		out = fmtVxVy("OR", instr)
	// 8xy2
	case first == 0x8 && fourth == 0x2:
		out = fmtVxVy("AND", instr)
	// 8xy3
	case first == 0x8 && fourth == 0x3:
		out = fmtVxVy("XOR", instr)
	// 8xy4
	case first == 0x8 && fourth == 0x4:
		out = fmtVxVy("ADD", instr)
	// 8xy5
	case first == 0x8 && fourth == 0x5:
		out = fmtVxVy("SUB", instr)
	// 8xy6
	case first == 0x8 && fourth == 0x6:
		x, y := rNib(high), lNib(low)
		out = fmt.Sprintf("SHR V%X, {, V%X}", x, y)
	// 8xy7
	case first == 0x8 && fourth == 0x7:
		out = fmtVxVy("SUBN", instr)
	// 8xyE
	case first == 0x8 && fourth == 0xE:
		x, y := rNib(high), lNib(low)
		out = fmt.Sprintf("SHL V%X, {, V%X}", x, y)
	// 9xy0
	case first == 0x9:
		out = fmtVxVy("SNE", instr)
	// Annn
	case first == 0xA:
		addr := getAddr([]byte{high, low})
		out = fmt.Sprintf("LD I, 0x%03X", addr)
	// Bnnn
	case first == 0xB:
		addr := getAddr([]byte{high, low})
		out = fmt.Sprintf("JP V0, 0x%03X", addr)
	// CxKK
	case first == 0xC:
		out = fmtVxByte("RND", instr)
	// Dxyn
	case first == 0xD:
		x, y := rNib(high), lNib(low)
		out = fmt.Sprintf("DRW V%X, V%X, 0x%X", x, y, fourth)
	// Ex9E
	case first == 0xE && low == 0x9E:
		x := rNib(high)
		out = fmt.Sprintf("SKP V%X", x)
	// ExA1
	case first == 0xE && low == 0xA1:
		x := rNib(high)
		out = fmt.Sprintf("SKNP V%X", x)
	// Fx07
	case first == 0xF && low == 0x07:
		x := rNib(high)
		out = fmt.Sprintf("LD V%X, DT", x)
	// Fx0A
	case first == 0xF && low == 0X0A:
		x := rNib(high)
		out = fmt.Sprintf("LD V%X, K", x)
	// Fx15
	case first == 0xF && low == 0x15:
		x := rNib(high)
		out = fmt.Sprintf("LD DT, V%X", x)
	// Fx18
	case first == 0xF && low == 0x18:
		x := rNib(high)
		out = fmt.Sprintf("LD ST, V%X", x)
	// Fx1E
	case first == 0xF && low == 0x1E:
		x := rNib(high)
		out = fmt.Sprintf("ADD I, V%X", x)
	// Fx29
	case first == 0xF && low == 0x29:
		x := rNib(high)
		out = fmt.Sprintf("LD F, V%X", x)
	// Fx33
	case first == 0xF && low == 0x33:
		x := rNib(high)
		out = fmt.Sprintf("LD B, V%X", x)
	//Fx55
	case first == 0xF && low == 0x55:
		x := rNib(high)
		out = fmt.Sprintf("LD [I], V%X", x)
	// Fx65
	case first == 0xF && low == 0x65:
		x := rNib(high)
		out = fmt.Sprintf("LD V%X, [I]", x)
	default:
		out = "BAD INSTR"
	}
	return out
}

func fmtVxVy(op string, instr []byte) string {
	x, y := rNib(instr[0]), lNib(instr[1])
	return fmt.Sprintf("%s V%X, V%X", op, x, y)
}

func fmtVxByte(op string, instr []byte) string {
	x := rNib(instr[0])
	return fmt.Sprintf("%s V%X, 0x%02X", op, x, instr[1])
}

func Disassemble(opCodes []byte, offset uint16) strings.Builder {
	var out strings.Builder

	// IF there's an odd number of bytes passed in we ignore the
	// last byte.  To do so we compare i to one less then the
	// lenght of the opCodes slice.
	for i := 0; i < len(opCodes)-1; i += 2 {
		high, low := opCodes[i], opCodes[i+1]
		out.WriteString(fmt.Sprintf("0x%03X   %02X%02X", offset+uint16(i), high, low) + "   " + translateOpCode([]byte{high, low}) + "\n")
	}
	return out
}
