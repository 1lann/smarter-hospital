package arclient

import (
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/1lann/smarter-hospital/proto"
	"github.com/cenkalti/rpc2"
)

// ErrAuthFail is the error returned by Connect when authentication fails.
var ErrAuthFail = errors.New("arclient: authentication failure")

// ActionHandler represents a handler for actions.
type ActionHandler func(action proto.Action) error

// The stream used for input and output to communicate with the ATMega.
var (
	Input  io.Reader
	Output io.Writer
)

// Client represents a connected client to the system.
type Client struct {
	client *rpc2.Client
}

// Emit sends an event to the system.
func (c *Client) Emit(name string, val interface{}) {
	err := c.client.Notify(proto.EventMsg, proto.Event{
		Name:  name,
		Value: val,
	})
	if err != nil {
		log.Println("arclient: emit fail:", err)
	}
}

func createActionHandler(actionHandler ActionHandler) func(client *rpc2.Client,
	action *proto.Action, result *proto.Result) error {
	return func(client *rpc2.Client, action *proto.Action, result *proto.Result) error {
		err := actionHandler(*action)
		if err != nil {
			result.Successful = false
			result.Message = err.Error()
			return nil
		}

		result.Successful = true
		return nil
	}
}

func (c *Client) reconnect(addr string, id string) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		return err
	}

	c.client = rpc2.NewClient(conn)
	go c.client.Run()

	authResponse := proto.AuthResponse{}
	err = c.client.Call(proto.AuthMsg, &proto.AuthRequest{ID: id}, &authResponse)
	if err != nil {
		return err
	}

	if !authResponse.Authenticated {
		return ErrAuthFail
	}

	return nil
}

// Connect creates a client connection to the system, where addr is the address
// of the server, id is the id of the client for authentication, actionHandler
// is handler when an action request is made from the server.
//
// Connect is non-blocking, and will return a client connection to emit
// events to the server.
func Connect(addr string, id string, actionHandler ActionHandler) (*Client, error) {
	if Input == nil {
		panic("arclient: connect: arclient.Input must be set")
	}

	if Output == nil {
		panic("arclient: connect: arclient.Output must be set")
	}

	c := &Client{}
	err := c.reconnect(addr, id)
	if err != nil {
		return nil, err
	}

	c.client.Handle(proto.HealthCheckMsg, healthCheckHandler)
	if actionHandler != nil {
		c.client.Handle(proto.ActionMsg, createActionHandler(actionHandler))
	}

	go func() {
		<-c.client.DisconnectNotify()
		log.Println("arclient: disconnected, reconnecting...")
		c.client.Close()

		for {
			err := c.reconnect(addr, id)
			if err != nil {
				log.Println("arclient: reconnect error:", err)
				time.Sleep(time.Second * 5)
				continue
			}

			log.Println("arclient: connection re-established")

			c.client.Handle(proto.HealthCheckMsg, healthCheckHandler)
			if actionHandler != nil {
				c.client.Handle(proto.ActionMsg, createActionHandler(actionHandler))
			}

			<-c.client.DisconnectNotify()
			log.Println("arclient: disconnected, reconnecting...")
			c.client.Close()
		}
	}()

	return c, nil
}

func healthCheckHandler(client *rpc2.Client, _ interface{}, result *proto.Result) error {
	result.Successful = true
	return nil
}
