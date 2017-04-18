package main

import "log"

func authHandler(id string) bool {
	log.Println(id, "just connected to the server!")
	return true
}
