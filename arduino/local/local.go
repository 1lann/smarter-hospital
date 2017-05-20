package main

import (
	"log"

	"github.com/1lann/smarter-hospital/core"

	_ "github.com/1lann/smarter-hospital/modules/heartrate"
	_ "github.com/1lann/smarter-hospital/modules/lights"
	_ "github.com/1lann/smarter-hospital/modules/ultrasonic"
)

func main() {
	core.SetupModule("lights", "light1")
	core.SetupModule("ultrasonic", "ultrasonic1")
	core.SetupModule("heartrate", "heartrate1")

	// err := arduino.Connect("/dev/tty.usbmodem1411")
	// if err != nil {
	// 	panic(err)
	// }
	log.Println("connecting...")
	core.Connect("127.0.0.1:5000")
	select {}
}
