package main

import (
	"fmt"
)

func (c8 *Chip8) execInstr() {
	instr := c8.memory[c8.programPtr : c8.programPtr+2]
	fmt.Printf("Instr: %X %v\n", instr, instr)

	switch {
	// 00E0 - CLS - clear the display
	case instr[0] == 0 && instr[1] == 224:
		c8.display.clear()
	// 6xkk - LD Vx, byte - Load the byte value into the register specified by x
	case lNib(instr[0]) == x6:
		c8.registers[rNib(instr[0])] = instr[1]
	// Annn - LD I, addr - Load the int16 addr specified by nnn into the I register
	case lNib(instr[0]) == xA:
		addr := getAddr(instr[0], instr[1])
		c8.regI = addr
	// Dxyn - DRW Vx, Vy, nibble - grab an nibble length byte from I and draw it at the
	// values of Vx and Vy.
	// TODO:
	// a. If a pixel is erased set VF to 1, otherwise 0.
	// b. If part of a sprite is offscreen it should be wrapped
	case lNib(instr[0]) == xD:
		xOffset := 55 - rNib(instr[0])
		yOffset := lNib(instr[0])
		length := rNib(instr[1])

		sprite := c8.memory[c8.regI : c8.regI+uint16(length)]
		//for _, line := range sprite {
		//	fmt.Printf("%08b\n", line)
		//}

		// worry about horizontal wrapping later
		c8.registers[xF] = 0
		for i := 0; i < int(length); i++ {
			var spriteRow uint64
			spriteRow = uint64(sprite[i]) << xOffset
			y := (yOffset + byte(i)) % 32
			currentRow := c8.display.screen[y]

			collisionFree := currentRow | spriteRow
			c8.display.screen[y] = currentRow ^ spriteRow
			if collisionFree^c8.display.screen[y] > 0 && c8.registers[xF] == 0 {
				c8.registers[xF] = 1
			}
		}
	default:
		msg := fmt.Sprintf("Unknown Instruction: %X\n", instr)
		panic(msg)
	}

	c8.programPtr = c8.programPtr + 2
}

func lNib(b byte) byte {
	return b >> 4
}

func rNib(b byte) byte {
	return b & 15
}

func getAddr(b1, b2 byte) uint16 {
	right := uint16(rNib(b1)) << 8
	left := uint16(b2)
	return right + left
}
