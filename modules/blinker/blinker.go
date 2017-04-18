package blinker

import "github.com/1lann/smarter-hospital/comm"

func init() {
	core.RegisterAction(Action{})
}

// Action represents a blinker action.
type Action struct {
	Rate int // 0 = super slow, 100 = super fast
}
