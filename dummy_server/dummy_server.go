package main

import (
	"time"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/modules/dummy"
)

func main() {
	core.SetupModule("dummy", "dummy1", dummy.Settings{
		HelloWorld: "Boop",
	})
	server, err := core.NewServer("127.0.0.1:5000")
	if err != nil {
		panic(err)
	}

	core.RegisterConnect("dummy1", func() {
		time.Sleep(time.Second)

		server.Do("dummy1", dummy.Action{
			SetState: true,
		})
	})

	for {
		time.Sleep(time.Minute)
	}
}
