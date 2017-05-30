// +build linux

package main

import (
	"log"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/pi/drivers"

	_ "github.com/1lann/smarter-hospital/modules/climate"
	_ "github.com/1lann/smarter-hospital/modules/heartrate"
	_ "github.com/1lann/smarter-hospital/modules/lights"
	_ "github.com/1lann/smarter-hospital/modules/proximity"
	_ "github.com/1lann/smarter-hospital/modules/thermometer"
	_ "github.com/1lann/smarter-hospital/modules/ultrasonic"
	_ "github.com/kidoman/embd/host/rpi"
)

func main() {
	core.SetupModule("ultrasonic", "ultrasonic1")
	core.SetupModule("thermometer", "thermometer1")
	core.SetupModule("heartrate", "heartrate1")
	core.SetupModule("proximity", "proximity1")
	core.SetupModule("climate", "climate1")
	core.SetupModule("lights", "lights1")
	err := drivers.Connect()
	if err != nil {
		panic(err)
	}

	log.Println("Starting up...")
	core.Connect("192.168.8.232:5000")
	select {}
}
