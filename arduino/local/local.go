package main

import (
	"github.com/1lann/smarter-hospital/core"
	_ "github.com/1lann/smarter-hospital/modules/ping"
)

func main() {
	core.SetupModule("ping", "ping1")
	// err := arduino.Connect("/dev/tty.usbmodem1411")
	// if err != nil {
	// 	panic(err)
	// }
	core.Connect("127.0.0.1:5000")
	select {}
}
