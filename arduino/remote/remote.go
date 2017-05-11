package main

import (
	"time"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
	_ "github.com/1lann/smarter-hospital/modules/dummy"
)

func main() {
	core.SetupModule("dummy", "dummy1")
	arduino.Connect("/dev/ttyATH0")
	core.Connect("127.0.0.1:5000")
	for {
		time.Sleep(time.Minute)
	}
}
