package main

import (
	"os"
	"time"

	"github.com/1lann/smarter-hospital/comm"
	"github.com/1lann/smarter-hospital/comm/blinker"
)

func handleAction(action core.Action) error {
	if action.Name == "blink" {
		blinkerAction := action.Value.(blinker.Action)
		os.Stdout.Write([]byte{byte(blinkerAction.Rate)})
	}

	return nil
}

func handlePing() bool {
	return true
}

func main() {
	core.Connect(os.Getenv("ADDR"), "arduino", handleAction, handlePing)

	for {
		time.Sleep(time.Hour * 10)
	}
}
