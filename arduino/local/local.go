package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"

	// _ "github.com/1lann/smarter-hospital/modules/heartrate"
	_ "github.com/1lann/smarter-hospital/modules/lights"
	_ "github.com/1lann/smarter-hospital/modules/ultrasonic"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	core.SetupModule("lights", "light1")
	core.SetupModule("ultrasonic", "ultrasonic1")
	// core.SetupModule("heartrate", "heartrate1")

	err := arduino.Connect("/dev/tty.usbmodem1411")
	if err != nil {
		panic(err)
	}
	log.Println("connecting...")
	core.Connect("127.0.0.1:5000")
	select {}
}
