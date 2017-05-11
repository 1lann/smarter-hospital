package core

import (
	"errors"
	"log"
	"net"
	"reflect"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/1lann/smarter-hospital/store"
	"github.com/cenkalti/rpc2"
)

// Some errors that the Do can possibly return
var (
	ErrClientNotFound = errors.New("core: no such client")
	ErrDeviceError    = errors.New("core: device reported error")
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

	for moduleID, connectedClient := range s.clients {
		if client == connectedClient {
			dcHandler, found := registeredDisconnect[moduleID]
			if found {
				go dcHandler()
			}

			delete(s.clients, moduleID)
		}
	}
	s.clientsMutex.Unlock()

	go client.Close()
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
			log.Println("core: client health check timeout, disconnecting:",
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
		log.Println("core: client failed health check, disconnecting:",
			moduleID, "reason:", err)
		s.disconnect(moduleID)
		return
	}

	if !result.Successful {
		log.Println("core: client sent incorrect response for health check, "+
			"disconnecting:", moduleID)
		s.disconnect(moduleID)
		return
	}

	log.Println("core: debug: health check for \""+moduleID+"\" successful "+
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

			log.Println("core: client took too long to handshake")
			client.Close()
		}()
	})

	s.server.OnDisconnect(func(client *rpc2.Client) {
		s.clientsMutex.Lock()
		for moduleID, connectedClient := range s.clients {
			if client == connectedClient {
				dcHandler, found := registeredDisconnect[moduleID]
				if found {
					go dcHandler()
				}

				delete(s.clients, moduleID)
			}
		}
		s.clientsMutex.Unlock()

		log.Println("core: debug: client disconnected")
	})

	s.server.Handle(HandshakeMsg, s.handleHandshake)

	s.server.Handle(EventMsg, func(client *rpc2.Client, event *Event,
		result *bool) error {
		go s.handleEvent(client, event)
		*result = true
		return nil
	})

	s.server.Handle(ErrorMsg, func(client *rpc2.Client, errorMsg *string,
		result *bool) error {
		go s.handleError(client, *errorMsg)
		*result = true
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

func (s *Server) handleHandshake(client *rpc2.Client, req *HandshakeRequest,
	resp *HandshakeResponse) error {
	settings := make(map[string]interface{})

	s.clientsMutex.Lock()
	for _, moduleID := range req.ModuleIDs {
		module, found := setupModules[moduleID]
		if !found {
			resp.Successful = false
			log.Println("core: handshake: client attempted to handshake "+
				"with non existent module ID:", moduleID)
			go client.Close()
			s.clientsMutex.Unlock()
			return nil
		}

		connectHandler, found := registeredConnect[moduleID]
		if found {
			go connectHandler()
		}

		if module.registration.settingsType != nil {
			settings[moduleID] = module.module.Elem().
				FieldByName("Settings").Interface()
		}

		if oldClient, found := s.clients[moduleID]; found {
			log.Println("core: warning: old client for " + moduleID +
				" found")
			go oldClient.Close()
		}
		s.clients[moduleID] = client
	}
	s.clientsMutex.Unlock()

	client.State.Set("handshaken", true)

	resp.Successful = true
	resp.ModuleSettings = settings

	return nil
}

func (s *Server) handleEvent(client *rpc2.Client, event *Event) {
	s.clientsMutex.Lock()
	expectedClient, found := s.clients[event.ModuleID]
	s.clientsMutex.Unlock()

	if !found {
		log.Println("core: refusing to handle event from non-existent module, "+
			"disconnected client or incomplete handshake:", event.ModuleID)
		return
	}

	if expectedClient != client {
		log.Println("core: refusing to handle event from client that does "+
			"own module:", event.ModuleID)
		return
	}

	shaken, ok := client.State.Get("handshaken")
	if !ok || !shaken.(bool) {
		log.Println(
			"core: refusing to handle event from client that has not " +
				"completed handshake (you should never see this)")
		return
	}

	module, found := setupModules[event.ModuleID]
	if !found {
		log.Println("core: module \"" + event.ModuleID + "\" does not exist")
		return
	}

	if !module.registration.hasEventHandler {
		log.Println("core: module \"" + event.ModuleID +
			"\" does not support events")
		return
	}

	if !reflect.TypeOf(event.Value).
		AssignableTo(module.registration.eventType) {
		log.Println("core: unable to handle event: unassignable type for "+
			"event:", event.ModuleID)
		return
	}

	module.module.MethodByName("HandleEvent").
		Call([]reflect.Value{reflect.ValueOf(event.Value)})

	logics := moduleToLogic[event.ModuleID]
	for _, logic := range logics {
		go registeredLogic[logic].trigger(s)
	}
}

func (s *Server) handleError(client *rpc2.Client, errorMessage string) {
	s.clientsMutex.Lock()
	found := false
	var moduleID string
	for connectedModuleID, connectedClient := range s.clients {
		if connectedClient == client {
			found = true
			moduleID = connectedModuleID
			break
		}
	}
	s.clientsMutex.Unlock()

	if !found {
		log.Println("core: refusing to handle error from disconnected client " +
			"or incomplete handshake")
		return
	}

	log.Println("core: module \"" + moduleID + "\" reported error")

	store.C("module_errors").EnsureIndexKey("Time")
	store.C("module_errors").EnsureIndexKey("ModuleID")
	err := store.C("module_errors").Insert(bson.M{
		"moduleid": moduleID,
		"error":    errorMessage,
		"time":     time.Now(),
	})
	if err != nil {
		log.Println("core: failed to store error in store:", err)
	}
}
