package main

import (
	"fmt"
	"log"

	"github.com/1lann/smarter-hospital/comm"
	"github.com/1lann/smarter-hospital/comm/heartrate"
)

var server *comm.Server

var handlers = map[string]interface{}{
	"test": testHandler,
}

func testHandler(c comm.RemoteClient, data heartrate.Event) error {
	fmt.Println(data)
	return nil
}

func authHandler(id string) bool {
	log.Println(id, "just connected to the server!")
	return true
}
