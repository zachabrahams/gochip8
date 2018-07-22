package main

type Keyboard struct {
	checkKey   chan byte
	anyKey     chan bool
	keyPressed chan bool

	newKeyboardState chan uint16
}

func NewKeyboard() *Keyboard {
	checkKey := make(chan byte)
	keyPressed := make(chan bool)
	anyKey := make(chan bool)

	newKeyboardState := make(chan uint16)

	var keys uint16
	keys = 0
	go func() {
		for {
			select {
			case keys = <-newKeyboardState:
			case keyToCheck := <-checkKey:
				if keys&(0x1<<keyToCheck) == uint16(0x1<<keyToCheck) {
					keyPressed <- true
				} else {
					keyPressed <- false
				}
			case <-anyKey:
				if keys != 0 {
					keyPressed <- true
				} else {
					keyPressed <- false
				}
			}
		}
	}()

	return &Keyboard{
		checkKey:   checkKey,
		anyKey:     anyKey,
		keyPressed: keyPressed,

		newKeyboardState: newKeyboardState,
	}
}

func (k *Keyboard) isPressed(key byte) bool {
	k.checkKey <- key
	pressed := <-k.keyPressed
	return pressed
}
