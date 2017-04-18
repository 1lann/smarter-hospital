package core

import (
	"errors"
	"log"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/cenkalti/rpc2"
)

// ErrHandshakeFail is the error returned by Connect when the handshake fails.
var ErrHandshakeFail = errors.New("comm: handshake failure")

// ActionHandler represents a handler for actions.
type ActionHandler func(action Action) error

// PingHandler represents a handler that pings the ATMega to check if it's
// responsive.
type PingHandler func() bool

// Client represents a connected client to the system.
type Client struct {
	client      *rpc2.Client
	commLock    *sync.Mutex
	pingHandler PingHandler
}

// Emit sends an event to the system.
func (c *Client) Emit(name string, val interface{}) {
	if c.client == nil {
		log.Println("comm: emit failed: client not ready")
		return
	}

	err := c.client.Notify(EventMsg, Event{
		Name:  name,
		Value: val,
	})
	if err != nil {
		log.Println("comm: emit failed:", err)
	}
}

func createActionHandler(actionHandler ActionHandler) func(client *rpc2.Client,
	action *Action, result *Result) error {
	return func(client *rpc2.Client, action *Action, result *Result) error {
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

	var moduleIDs []string

	for moduleID := range setupModules {
		moduleIDs = append(moduleIDs, moduleID)
	}

	var handshakeResp HandshakeResponse
	err = c.client.Call(AuthMsg, &HandshakeRequest{
		ModuleIDs: moduleIDs,
	}, &handshakeResp)
	if err != nil {
		return err
	}

	if !handshakeResp.Successful {
		return ErrHandshakeFail
	}

	for moduleID, settings := range handshakeResp.ModuleSettings {
		module, found := setupModules[moduleID]
		if !found {
			log.Println("comm: handshake module settings: unknown module ID:",
				moduleID)
			continue
		}

		if !reflect.TypeOf(settings).AssignableTo(
			module.module.FieldByName("Settings").Type()) {
			log.Println("comm: handshake module settings: settings is not "+
				"assignable for:", moduleID)
			return ErrHandshakeFail
		}

		module.module.FieldByName("Settings").Set(reflect.ValueOf(settings))
	}

	return nil
}

// Connect creates a client connection to the system, where addr is the address
// of the server, id is the id of the client for authentication, actionHandler
// is handler when an action request is made from the server.
//
// Connect is non-blocking and will return a client connection to emit events
// to the server.
func Connect(addr string, id string, actionHandler ActionHandler,
	pingHandler PingHandler) *Client {
	c := &Client{
		commLock:    new(sync.Mutex),
		pingHandler: pingHandler,
	}

	go func() {
		for {
			err := c.reconnect(addr, id)
			if err != nil {
				log.Println("comm: reconnect error:", err)
				time.Sleep(time.Second * 5)
				continue
			}

			log.Println("comm: re-connected")

			c.client.Handle(HealthCheckMsg, c.healthCheckHandler)
			if actionHandler != nil {
				c.client.Handle(ActionMsg, createActionHandler(actionHandler))
			}

			<-c.client.DisconnectNotify()
			log.Println("comm: disconnected, reconnecting...")
			c.client.Close()
		}
	}()

	return c
}

func (c *Client) healthCheckHandler(client *rpc2.Client, _ *bool,
	result *Result) error {
	result.Successful = c.pingHandler()

	return nil
}
