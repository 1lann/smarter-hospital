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

// Some errors that the Do can possibly return
var (
	ErrClientNotFound = errors.New("comm: no such client")
	ErrDeviceError    = errors.New("comm: device reported error")
)

// Server represents a running RPC2 server.
type Server struct {
	server       *rpc2.Server
	clients      map[string]*rpc2.Client
	clientsMutex *sync.Mutex
}

// Do performs an action to the provided client ID.
func (s *Server) Do(moduleID string, val interface{}) (string, error) {
	s.clientsMutex.Lock()
	client, found := s.clients[moduleID]
	s.clientsMutex.Unlock()
	if !found {
		return "The device that performs that action is disconnected " +
			"from the system!", ErrClientNotFound
	}

	result := Result{}
	err := client.Call(ActionMsg, Action{ModuleID: moduleID, Value: val},
		&result)
	if err != nil {
		return "There was an error communicating to the device which " +
			"performs that action!", err
	}

	if !result.Successful {
		if result.Message == "" {
			return "The device reported that the action could not be " +
				"completed!", ErrDeviceError
		}
		return result.Message, ErrDeviceError
	}

	return "", nil
}

func (s *Server) healthCheck() {
	ticker := time.Tick(time.Second * 5)

	for _ = range ticker {
		s.clientsMutex.Lock()
		for moduleID, client := range s.clients {
			go s.checkHealth(moduleID, client)
		}
		s.clientsMutex.Unlock()
	}
}

func (s *Server) disconnect(moduleID string) {
	s.clientsMutex.Lock()
	client, ok := s.clients[moduleID]
	if !ok {
		s.clientsMutex.Unlock()
		return
	}
	delete(s.clients, moduleID)
	s.clientsMutex.Unlock()

	client.Close()
}

func (s *Server) checkHealth(moduleID string, client *rpc2.Client) {
	result := Result{}

	seq := 0
	rawSeq, ok := client.State.Get("health_check_ongoing")
	if ok {
		seq = rawSeq.(int)
	}

	go func(expectedSeq int) {
		time.Sleep(time.Second * 5)
		seq, _ := client.State.Get("health_check_ongoing")
		if seq.(int) == expectedSeq {
			log.Println("comm: client health check timeout, disconnecting:",
				moduleID)
			s.disconnect(moduleID)
		}
	}(seq)

	start := time.Now()
	client.State.Set("health_check_ongoing", seq)
	err := client.Call(HealthCheckMsg, true, &result)
	latency := time.Since(start)
	client.State.Set("health_check_latency", latency)

	seq = seq + 1
	seq = seq % 65535

	client.State.Set("health_check_ongoing", seq)

	if err != nil {
		log.Println("comm: client failed health check, disconnecting:",
			moduleID, "reason:", err)
		s.disconnect(moduleID)
		return
	}

	if !result.Successful {
		log.Println("comm: client sent incorrect response for health check, "+
			"disconnecting:", moduleID)
		s.disconnect(moduleID)
		return
	}

	log.Println("comm: debug: health check for \""+moduleID+"\" successful "+
		"with latency:", latency.Nanoseconds()/1000000, "ms")
}

// NewServer starts the rpc2 server and listens asynchronously on the provided
// address.
func NewServer(addr string) (*Server, error) {
	s := &Server{
		server:       nil,
		clients:      make(map[string]*rpc2.Client),
		clientsMutex: new(sync.Mutex),
	}

	s.server = rpc2.NewServer()
	s.server.OnConnect(func(client *rpc2.Client) {
		client.State = rpc2.NewState()

		go func() {
			time.Sleep(time.Second * 10)
			if shaken, ok := client.State.Get("handshaken"); ok &&
				shaken.(bool) {
				return
			}

			log.Println("comm: client took too long to authenticate")
			client.Close()
		}()
	})

	s.server.OnDisconnect(func(client *rpc2.Client) {
		id, ok := client.State.Get("id")
		if !ok {
			return
		}

		s.clientsMutex.Lock()
		delete(s.clients, id.(string))
		s.clientsMutex.Unlock()

		log.Println("comm: debug: client disconnected:", id.(string))
	})

	s.server.Handle(HandshakeMsg, func(client *rpc2.Client,
		req *HandshakeRequest, resp *HandshakeResponse) error {
		s.clientsMutex.Lock()
		for _, moduleID := range req.ModuleIDs {
			if oldClient, found := s.clients[moduleID]; found {
				log.Println("comm: warning: old client for " + moduleID +
					" found")
				go oldClient.Close()
			}
			s.clients[moduleID] = client
		}
		s.clientsMutex.Unlock()

		client.State.Set("handshaken", true)

		return nil
	})

	s.server.Handle(EventMsg, func(client *rpc2.Client,
		event *Event, result *Result) error {
		go s.handleMessage(client, event, result)
		return nil
	})

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	go s.server.Accept(listener)
	go s.healthCheck()

	return s, nil
}

func (s *Server) handleMessage(client *rpc2.Client, event *Event,
	result *Result) {
	shaken, ok := client.State.Get("handshaken")
	if !ok || !shaken.(bool) {
		result.Successful = false
		result.Message =
			"Refusing to handle event from unauthenticated client"
		log.Println(
			"comm: refusing to handle event from unauthenticated client")
		return
	}

	module, found := setupModules[event.ModuleID]
	if !found {
		result.Successful = false
		result.Message = "Module \"" + event.ModuleID +
			"\" does not exist"
		log.Println("comm: module \"" + event.ModuleID +
			"\" does not exist")
		return
	}

	if !module.registration.hasEventHandler {
		result.Successful = false
		result.Message = "Module \"" + event.ModuleID +
			"\" does not support events"
		log.Println("comm: module \"" + event.ModuleID +
			"\" does not support events")
		return
	}

	if !reflect.TypeOf(event.Value).
		AssignableTo(module.registration.eventType) {
		result.Successful = false
		result.Message = "Unassignable types for event: " + event.ModuleID
		log.Println("comm: unassignable types for event:", event.ModuleID)
		return
	}

	results := module.module.MethodByName("HandleEvent").
		Call([]reflect.Value{reflect.ValueOf(event.Value)})

	err := results[0].Interface()
	if err != nil {
		result.Successful = false
		result.Message = "Server error while handling event"
		log.Println("comm: error while handling event \""+event.ModuleID+
			"\":", err.(error))
		return
	}

	result.Successful = true
}
