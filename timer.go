package main

import (
	"time"
)

type Timer struct {
	getVal     chan bool
	setVal     chan byte
	receiveVal chan byte
}

func (t *Timer) Read() byte {
	t.getVal <- true
	val := <-t.receiveVal
	return val
}

func (t *Timer) Set(newVal byte) {
	t.setVal <- newVal
}

func NewTimer(cb func()) *Timer {
	getVal := make(chan bool)
	setVal := make(chan byte)
	receiveVal := make(chan byte)

	ticker := time.NewTicker(TIMER_TICK * time.Millisecond)
	var val byte
	val = 0

	go func() {
		for {
			select {
			case newVal := <-setVal:
				val = newVal
			case <-getVal:
				receiveVal <- val
			case <-ticker.C:
				if val > 0 {
					val--
					if val == 0 {
						cb()
					}
				}
			}
		}
	}()

	return &Timer{
		getVal:     getVal,
		setVal:     setVal,
		receiveVal: receiveVal,
	}
}
