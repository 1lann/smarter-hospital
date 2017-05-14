package main

import (
	"time"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"

	_ "github.com/1lann/smarter-hospital/modules/lights"
	_ "github.com/1lann/smarter-hospital/modules/ultrasonic"
)

func main() {
	core.SetupModule("lights", "light1")
	core.SetupModule("ultrasonic", "ultrasonic1")

	arduino.Connect("/dev/ttyATH0")

	core.Connect("192.168.240.232:5000")
	for {
		time.Sleep(time.Minute)
	}
}
