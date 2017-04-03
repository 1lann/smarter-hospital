package blinker

import "encoding/gob"

func init() {
	gob.RegisterName("blinkact", Action{})
}

// Action represents a blinker action.
type Action struct {
	Rate int // 0 = super slow, 100 = super fast
}
