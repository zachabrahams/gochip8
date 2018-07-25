package main

type Keyboard struct {
	checkKey   chan byte
	anyKey     chan bool
	nextKey    chan byte
	keyPressed chan bool

	newKeyboardState chan uint16
}

func NewKeyboard() *Keyboard {
	checkKey := make(chan byte)
	keyPressed := make(chan bool)

	anyKey := make(chan bool)
	nextKey := make(chan byte)

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
				pressed := false
				for i := 0; i < 16; i++ {
					if (keys & (0x1 << uint16(i))) > 0 {
						nextKey <- byte(i)
						pressed = true
						break
					}
				}

				if !pressed {
					nextKey <- byte(255)
				}
			}
		}
	}()

	return &Keyboard{
		checkKey:   checkKey,
		anyKey:     anyKey,
		keyPressed: keyPressed,

		nextKey:          nextKey,
		newKeyboardState: newKeyboardState,
	}
}

func (k *Keyboard) isPressed(key byte) bool {
	k.checkKey <- key
	pressed := <-k.keyPressed
	return pressed
}

func (k *Keyboard) nextPress() byte {
	k.anyKey <- true
	pressed := <-k.nextKey
	return pressed
}
