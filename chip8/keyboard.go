package chip8

import "sync"

type Keyboard struct {
	mutex *sync.Mutex
	state uint16
}

func NewKeyboard() *Keyboard {

	mutex := &sync.Mutex{}
	var state uint16 = 0

	return &Keyboard{
		mutex: mutex,
		state: state,
	}
}

func (k *Keyboard) Update(newState uint16) {
	k.mutex.Lock()
	k.state = newState
	k.mutex.Unlock()
}

func (k *Keyboard) isPressed(key byte) bool {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	if k.state&(0x1<<key) == uint16(0x1<<key) {
		return true
	}

	return false
}

func (k *Keyboard) nextPress() byte {
	k.mutex.Lock()
	keys := k.state
	k.mutex.Unlock()
	for i := 0; i < 16; i++ {
		if (keys & (0x1 << uint16(i))) > 0 {
			return byte(i)
		}
	}

	return byte(255)
}
