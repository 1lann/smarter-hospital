package core

import (
	"errors"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/cenkalti/rpc2"
)

// ErrHandshakeFail is the error returned by Connect when the handshake fails.
var ErrHandshakeFail = errors.New("core: handshake failure")

// ActionHandler represents a handler for actions.
type ActionHandler func(action Action) error

// PingHandler represents a handler that pings the ATMega to check if it's
// responsive.
type PingHandler func() bool

// Client represents a connected client to the system.
type Client struct {
	client      *rpc2.Client
	pingHandler PingHandler
}

// Emit sends an event to the system.
func (c *Client) Emit(moduleID string, val interface{}) {
	if c.client == nil {
		log.Println("core: emit failed: client not ready")
		return
	}

	if !reflect.TypeOf(val).
		AssignableTo(setupModules[moduleID].registration.eventType) {
		log.Println("core: refusing to emit: event value not assignable to " +
			"module's registered event type")
		return
	}

	err := c.client.Notify(EventMsg, Event{
		ModuleID: moduleID,
		Value:    val,
	})
	if err != nil {
		log.Println("core: emit failed:", err)
	}
}

func (c *Client) reconnect(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		return err
	}

	c.client = rpc2.NewClient(conn)

	c.client.Handle(HealthCheckMsg, c.healthCheckHandler)
	c.client.Handle(ActionMsg, c.actionHandler)

	go c.client.Run()

	var moduleIDs []string

	for moduleID := range setupModules {
		moduleIDs = append(moduleIDs, moduleID)
	}

	var handshakeResp HandshakeResponse
	err = c.client.Call(HandshakeMsg, &HandshakeRequest{
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
			log.Println("core: handshake module settings: unknown module ID:",
				moduleID)
			continue
		}

		if !reflect.TypeOf(settings).AssignableTo(
			module.module.Elem().FieldByName("Settings").Type()) {
			log.Println("core: handshake module settings: settings is not "+
				"assignable for:", moduleID)
			return ErrHandshakeFail
		}

		module.module.Elem().FieldByName("Settings").
			Set(reflect.ValueOf(settings))
	}

	return nil
}

// Connect creates a client connection to the system, where addr is the address
// of the server, id is the id of the client for authentication, actionHandler
// is handler when an action request is made from the server.
//
// Connect is non-blocking and will return a client connection to emit events
// to the server.
func Connect(addr string, pingHandler PingHandler) *Client {
	c := &Client{
		pingHandler: pingHandler,
	}

	go func() {
		for {
			err := c.reconnect(addr)
			if err != nil {
				log.Println("core: reconnect error:", err)
				time.Sleep(time.Second * 5)
				continue
			}

			log.Println("core: re-connected")

			<-c.client.DisconnectNotify()
			log.Println("core: disconnected, reconnecting...")
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

func (c *Client) actionHandler(client *rpc2.Client, action *Action,
	result *Result) error {
	module, found := setupModules[action.ModuleID]
	if !found {
		result.Message = "Module not found"
		result.Successful = false
		return nil
	}

	if module.registration.actionType == nil {
		result.Message = "Module does not support actions"
		result.Successful = false
		return nil
	}

	results := module.module.MethodByName("HandleAction").Call([]reflect.Value{
		reflect.ValueOf(c),
		reflect.ValueOf(action.Value),
	})

	if results[0].IsNil() {
		return nil
	}

	err := results[0].Interface().(error)
	result.Message = err.Error()
	result.Successful = false
	return nil
}
