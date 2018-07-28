package chip8

type Beeper interface {
	Beep()
	Close()
}
