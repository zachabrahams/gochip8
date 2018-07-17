package main

import "fmt"

type Display struct {
	screen []uint64
}

func NewDisplay() *Display {
	s := make([]uint64, 32, 32)
	return &Display{screen: s}
}

func (d *Display) clear() {
	for i := range d.screen {
		d.screen[i] = 0
	}
}

func (d *Display) bitDump() {
	for _, line := range d.screen {
		fmt.Printf("%063b\n", line)
	}
}
