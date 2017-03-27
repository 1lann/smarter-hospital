package comm

import (
	"errors"
	"log"
	"net"
	"reflect"
	"strconv"
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
	handlers     map[string]interface{}
}

// RemoteClient represents a client connected to the RPC2 server.
type RemoteClient struct {
	ID    string
	State *rpc2.State
}

// GetAllClients returns a list of all the clients.
func (s *Server) GetAllClients() []RemoteClient {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	clients := make([]RemoteClient, len(s.clients))

	i := 0
	for id, client := range s.clients {
		clients[i] = RemoteClient{
			ID:    id,
			State: client.State,
		}
		i++
	}

	return clients
}

// Do performs an action to the provided client ID.
func (s *Server) Do(id string, name string, val interface{}) (string, error) {
	s.clientsMutex.Lock()
	client, found := s.clients[id]
	s.clientsMutex.Unlock()
	if !found {
		return "The device that performs that action is disconnected " +
			"from the system!", ErrClientNotFound
	}

	result := Result{}
	err := client.Call(ActionMsg, Action{Name: name, Value: val}, &result)
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
		for id, client := range s.clients {
			go s.checkHealth(id, client)
		}
		s.clientsMutex.Unlock()
	}
}

func (s *Server) disconnect(id string) {
	s.clientsMutex.Lock()
	client, ok := s.clients[id]
	if !ok {
		s.clientsMutex.Unlock()
		return
	}
	delete(s.clients, id)
	s.clientsMutex.Unlock()

	client.Close()
}

func (s *Server) checkHealth(id string, client *rpc2.Client) {
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
			log.Println("comm: client health check timeout, disconnecting:", id)
			s.disconnect(id)
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
		log.Println("comm: client failed health check, disconnecting:", id,
			"reason:", err)
		s.disconnect(id)
		return
	}

	if !result.Successful {
		log.Println("comm: client sent incorrect response for health check, "+
			"disconnecting:", id)
		s.disconnect(id)
		return
	}

	log.Println("comm: debug: health check for \""+id+"\" successful "+
		"with latency:", latency.Nanoseconds()/1000000, "ms")
}

// NewServer starts the rpc2 server and listens asynchronously on the provided
// address.
func NewServer(addr string, authHandler AuthHandler,
	handlers map[string]interface{}) (*Server, error) {
	for eventName, handler := range handlers {
		handlerType := reflect.TypeOf(handler)
		if handlerType.NumIn() != 2 {
			return nil, errors.New("comm: expected 2 arguments for handler \"" +
				eventName + "\", instead got " + strconv.Itoa(handlerType.NumIn()))
		}

		if handlerType.In(0) != reflect.TypeOf(Client{}) {
			return nil, errors.New("comm: first argument to event handler \"" +
				eventName + "\" must be of type Client")
		}

		if handlerType.NumOut() != 1 {
			return nil, errors.New("comm: expected 1 return argument for handler \"" +
				eventName + "\", instead got " + strconv.Itoa(handlerType.NumOut()))
		}

		if !handlerType.Out(0).Implements(reflect.TypeOf(errors.New(""))) {
			return nil, errors.New("comm: return argument for handler \"" +
				eventName + "\" must be an error")
		}
	}

	s := &Server{
		server:       nil,
		clients:      make(map[string]*rpc2.Client),
		clientsMutex: new(sync.Mutex),
		handlers:     handlers,
	}

	s.server = rpc2.NewServer()
	s.server.OnConnect(func(client *rpc2.Client) {
		client.State = rpc2.NewState()

		go func() {
			time.Sleep(time.Second * 10)
			if authed, ok := client.State.Get("authenticated"); ok && authed.(bool) {
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

	s.server.Handle(AuthMsg, func(client *rpc2.Client, req *AuthRequest,
		resp *AuthResponse) error {
		resp.Authenticated = authHandler(req.ID)
		if !resp.Authenticated {
			client.Close()
		}

		s.clientsMutex.Lock()
		if oldClient, found := s.clients[req.ID]; found {
			log.Println("comm: warning: old client for " + req.ID + " found")
			oldClient.Close()
		}
		s.clients[req.ID] = client
		s.clientsMutex.Unlock()

		client.State.Set("authenticated", true)
		client.State.Set("id", req.ID)

		return nil
	})

	s.server.Handle(EventMsg, func(client *rpc2.Client,
		event *Event, result *Result) error {
		id, ok := client.State.Get("id")
		if !ok {
			result.Successful = false
			result.Message = "Refusing to handle event from unauthenticated client"
			log.Println("comm: refusing to handle event from unauthenticated client")
			return nil
		}

		fn, ok := s.handlers[event.Name]
		if !ok {
			result.Successful = false
			result.Message = "Event handler \"" + event.Name + "\" does not exist"
			log.Println("comm: event handler for \"" + event.Name + "\" does not exist")
			return nil
		}

		if reflect.TypeOf(event.Value) != reflect.TypeOf(fn).In(1) {
			result.Successful = false
			result.Message = "Mis-matched types for event: " + event.Name
			log.Println("comm: mis-matched types for event:", event.Name)
			return nil
		}

		results := reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(RemoteClient{
				ID:    id.(string),
				State: client.State,
			}),
			reflect.ValueOf(event.Value),
		})

		err := results[0].Interface().(error)
		if err != nil {
			result.Successful = false
			result.Message = "Server error while handling event"
			log.Println("comm: error while handling event \""+event.Name+
				"\":", err)
			return nil
		}

		result.Successful = true
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

// AuthHandler represents an authentication request handler.
type AuthHandler func(id string) bool
