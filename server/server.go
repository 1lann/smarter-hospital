package main

import (
	"fmt"
	"log"

	"github.com/1lann/smarter-hospital/comm"
)

var server *comm.Server

func authHandler(id string) bool {
	log.Println(id, "just connected to the server!")
	return true
}

func main() {
	var err error
	server, err = comm.NewServer("0.0.0.0:5000", authHandler, make(map[string]interface{}))
	if err != nil {
		panic(err)
	}

	log.Println("Server is running!")

	for {
		var action string
		fmt.Scanf("%s", &action)
		if action == "on" {
			server.Do("led", "on", true)
		} else if action == "off" {
			server.Do("led", "off", true)
		} else {
			fmt.Println("That ain't a valid action!")
		}
	}
}
