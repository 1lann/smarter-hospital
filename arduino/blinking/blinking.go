package main

import (
	"os"
	"time"

	"github.com/1lann/smarter-hospital/comm"
)

func handleAction(action comm.Action) error {
	if action.Name == "off" {
		os.Stdout.Write([]byte{'0'})
	} else if action.Name == "on" {
		os.Stdout.Write([]byte{'1'})
	}
	return nil
}

func handlePing() bool {
	return true
}

func main() {
	comm.Connect("149.171.143.175:5000", "led", handleAction, handlePing)
	time.Sleep(time.Second * 1000)
}
