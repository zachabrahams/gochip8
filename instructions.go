package main

import (
	"crypto/rand"
	"fmt"
)

func (c8 *Chip8) execInstr() {
	nextInstr := c8.programPtr + 2
	instr := c8.memory[c8.programPtr:nextInstr]

	highI := instr[0]
	lowI := instr[1]
	rHighI := rNib(highI)
	lHighI := lNib(highI)
	rLowI := rNib(lowI)
	lLowI := lNib(lowI)

	switch {
	// 00E0 - CLS - clear the frame buffer
	case highI == 0x00 && lowI == 0xE0:
		c8.frameBuffer.clear()
	// 00EE RET returns from a subroutine
	case highI == 0x00 && lowI == 0xEE:
		cs := c8.callStack
		nextInstr, c8.callStack = cs[len(cs)-1], cs[:len(cs)-1]
	// 1nnn - JP addr
	case lHighI == 0x1:
		addr := getAddr(instr)
		nextInstr = addr
	// 2nnn - JP addr - pushes program counter +2 to the call stack and makes program counter = nnn
	case lHighI == 0x2:
		c8.callStack = append(c8.callStack, c8.programPtr+2)
		addr := getAddr(instr)
		nextInstr = addr
	// 3xkk - SE Vx, byte - Skip next instruction if Vx = kk
	case lHighI == 0x3:
		if c8.registers[rHighI] == lowI {
			nextInstr += 2
		}
	// 4xkk SNE Vx, byte - Skip next instruction if Vx != kk
	case lHighI == 0x4:
		if c8.registers[rHighI] != lowI {
			nextInstr += 2
		}
	// 5xy0 = SE Vx, Vy = Skip next instruction if Vx =  Vy.
	case lHighI == 0x5:
		if c8.registers[rHighI] == c8.registers[lLowI] {
			nextInstr += 2
		}
	// 6xkk - LD Vx, byte - Load the byte value into the register specified by x
	case lHighI == 0x6:
		c8.registers[rHighI] = lowI
	// 7xkk - ADD Vx, byte
	case lHighI == 0x7:
		c8.registers[rHighI] += lowI
	// 8xy0 - LD Vx, Vy - Set Vx to Vy
	case lHighI == 0x8 && rLowI == 0x0:
		c8.registers[rHighI] = c8.registers[lLowI]
	// 8xy2 - AND Vx, Vy Sets Vx to Vx & Vy
	case lHighI == 0x8 && rLowI == 0x2:
		x := c8.registers[rHighI]
		y := c8.registers[lLowI]
		c8.registers[rHighI] = x & y
	// 8xy4 - ADD Vx, Vy - Sets Vx to Vx +  Vy and sets VF to 1 if there is an overflow, 0 otherwise.
	case lHighI == 0x8 && rLowI == 0x4:
		x := uint16(c8.registers[rHighI])
		y := uint16(c8.registers[lLowI])
		sum := x + y
		if sum > 255 {
			c8.registers[0xF] = 1
			sum = sum & 255
		} else {
			c8.registers[0xF] = 0
		}

		c8.registers[rHighI] = byte(sum)
	// 8xy5 - SUB Vx, Vy - set Vx to Vx - Vy and VF = 1 iff Vx > Vy
	case lHighI == 0x8 && rLowI == 0x5:
		x := c8.registers[rHighI]
		y := c8.registers[lLowI]
		if x > y {
			c8.registers[0xF] = 1
		} else {
			c8.registers[0xF] = 0
		}

		c8.registers[rHighI] = x - y
	// 8xy6 - SHR Vx, {, Vy}
	case lHighI == 0x8 && rLowI == 0x6:
		if (c8.registers[rHighI] & 0x01) > 0 {
			c8.registers[0xF] = 1
		} else {
			c8.registers[0xF] = 0
		}

		c8.registers[rHighI] = c8.registers[rHighI] >> 1
	// 8xyE - SHL Vx {, Vy}
	case lHighI == 0x8 && rLowI == 0xE:
		if (c8.registers[rHighI] & 0x80) > 0 {
			c8.registers[0xF] = 1
		} else {
			c8.registers[0xF] = 0
		}

		c8.registers[rHighI] = c8.registers[rHighI] << 1
	// Annn - LD I, addr - Load the int16 addr specified by nnn into the I register
	case lHighI == 0xA:
		addr := getAddr(instr)
		c8.regI = addr
	// Cxkk - RND Vx, byte - generates a random byte, bitwise ANDs it with byte and
	// stores the result in Vx
	case lHighI == 0xC:
		if lowI == 0xFB {
			panic("asdf")
		}
		randBytes := make([]byte, 1, 1)
		_, err := rand.Read(randBytes)
		if err != nil {
			panic(err)
		}
		fmt.Printf("rand: %v\n", randBytes)

		c8.registers[rHighI] = (lowI & randBytes[0])
	// Dxyn - DRW Vx, Vy, nibble - grab an nibble length byte from I and draw it at the
	// values of Vx and Vy. If at least one pixel is erased set VF to 1 otherwise to 0
	// if a part of the sprite is located off screen - wrap it.
	case lHighI == 0xD:
		xOffset := 56 - int(c8.registers[rHighI])
		yOffset := c8.registers[lLowI]
		length := rLowI

		sprite := c8.memory[c8.regI : c8.regI+uint16(length)]
		//for _, line := range sprite {
		//	fmt.Printf("%08b\n", line)
		//}

		c8.registers[0xF] = 0
		for i := 0; i < int(length); i++ {
			var spriteRow uint64
			// if we need to wrap
			if xOffset < 0 {
				unWrappedOffset := xOffset * -1
				unWrapped := uint64(sprite[i]) >> (uint(unWrappedOffset))
				wrapped := uint64(sprite[i]) << uint(64+xOffset)
				spriteRow = unWrapped ^ wrapped

			} else {
				spriteRow = uint64(sprite[i]) << uint(xOffset)
			}
			y := (yOffset + byte(i)) % 32
			currentRow := c8.frameBuffer.buffer[y]

			collisionFree := currentRow | spriteRow
			c8.frameBuffer.buffer[y] = currentRow ^ spriteRow
			if collisionFree^c8.frameBuffer.buffer[y] > 0 && c8.registers[0xF] == 0 {
				c8.registers[0xF] = 1
			}
		}
	// ExA1 - SKNP Vx - Skips the next instruction if the key with Vxs value is not pressed
	case lHighI == 0xE && lowI == 0xA1:
		key := c8.registers[rHighI]
		pressed := c8.keyboard.isPressed(key)
		if !pressed {
			nextInstr += 2
		}
	// Fx07 - LD Vx, DT - Set Vx to be the value of the delay timer
	case lHighI == 0xF && lowI == 0x07:
		c8.registers[rHighI] = c8.delayTimer.Read()
	// Fx15 - LD DT, Vx - Set the delay timer the the value of Vx
	case lHighI == 0xF && lowI == 0x15:
		c8.delayTimer.Set(c8.registers[rHighI])
	// Fx18 - LD ST, Vx - set sound time to Vx's value
	case lHighI == 0xF && lowI == 0x18:
		//not implemented yet
	// Fx1E - ADD I, VX - Add Vx to I and store in I
	case lHighI == 0xF && lowI == 0x1E:
		c8.regI += uint16(c8.registers[rHighI])
	// Fx33 - LD B, Vx - Load BCD - Store 100's digit of B Vx value at I, 10s digit at I +1 and
	// ones at I + 2
	case lHighI == 0xF && lowI == 0x33:
		ones := c8.registers[rHighI] % 10
		tens := (c8.registers[rHighI] % 100) / 10
		hundreds := c8.registers[rHighI] / 100

		c8.memory[c8.regI] = hundreds
		c8.memory[c8.regI+1] = tens
		c8.memory[c8.regI+2] = ones
		// Fx29 - LD F, Vx - Set I to he logcation of the built in sprite for Vx's value
	case lHighI == 0xF && lowI == 0x29:
		// the built in sprites are stored at memory location 0, in order, with 5 bytes to a sprite.
		c8.regI = uint16(c8.registers[rHighI] * 5)
	// Fx55 - LD [I], Vx Load values from Vx into memory starting at I
	case lHighI == 0xF && lowI == 0x55:
		cursor := c8.regI
		var i byte
		for i = 0; i <= rHighI; i++ {
			c8.memory[cursor+uint16(i)] = c8.registers[i]
		}
	// Fx65 - LD Vx, [I] Load values from I into registers V0 to Vx
	case lHighI == 0xF && lowI == 0x65:
		cursor := c8.regI
		var i byte
		for i = 0; i <= rHighI; i++ {
			c8.registers[i] = c8.memory[cursor+uint16(i)]
		}
	default:
		msg := fmt.Sprintf("Unknown Instruction: %X\n", instr)
		panic(msg)
	}

	c8.programPtr = nextInstr
}

func lNib(b byte) byte {
	return b >> 4
}

func rNib(b byte) byte {
	return b & 15
}

func getAddr(bs []byte) uint16 {
	right := uint16(rNib(bs[0])) << 8
	left := uint16(bs[1])
	return right + left
}
