package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/cenkalti/rpc2"
)

func main() {
	var lastClient *rpc2.Client

	s := rpc2.NewServer()
	s.OnConnect(func(client *rpc2.Client) {
		lastClient = client
	})

	s.Handle("message", func(client *rpc2.Client, msg *Message, result *Result) error {
		fmt.Print(">", msg.Text)
		result.OK = true
		return nil
	})

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	go s.Accept(listener)

	rd := bufio.NewReader(os.Stdin)
	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			panic(err)
		}

		var result Result
		err = lastClient.Call("message", Message{str}, &result)
		if err != nil {
			fmt.Println("! error while sending message: ", err)
		} else if !result.OK {
			fmt.Println("! non OK response when sending message")
		}
	}
}
