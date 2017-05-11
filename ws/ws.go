package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"

	jsws "github.com/gopherjs/websocket"
	"github.com/gorilla/websocket"
)

type encodeMessage struct {
	Event string
	Value interface{}
}

type decodeMessage struct {
	Event string
	Value json.RawMessage
}

type registeredHandler struct {
	valueType reflect.Type
	handler   reflect.Value
}

// Server is the object for a WebSocket server.
type Server struct {
	upgrader       websocket.Upgrader
	wsClients      map[string]*websocket.Conn
	wsClientsMutex *sync.RWMutex
}

// Client is the object for the WebSocket client connected to a WebSocket
// server.
type Client struct {
	registeredHandlers map[string]registeredHandler
	connectHandler     func()
	disconnectHandler  func()
}

// NewServer returns a new server object, to be used on a server.
func NewServer() *Server {
	return &Server{
		upgrader:       websocket.Upgrader{},
		wsClients:      make(map[string]*websocket.Conn),
		wsClientsMutex: new(sync.RWMutex),
	}
}

// NewClient returns a new client object, to be used on a client. It is
// not ever used for a server.
func NewClient() *Client {
	return &Client{
		registeredHandlers: make(map[string]registeredHandler),
		connectHandler:     func() {},
		disconnectHandler:  func() {},
	}
}

// Handle is the handler to be called by the HTTP server to handle
// for websocket connections.
func (s *Server) Handle(r *http.Request, wr http.ResponseWriter) {
	conn, err := s.upgrader.Upgrade(wr, r, nil)
	if err != nil {
		log.Println("ws: upgrade:", err)
		return
	}

	s.wsClientsMutex.Lock()
	if _, exists := s.wsClients[r.RemoteAddr]; exists {
		log.Println("ws: already existing connection:",
			r.RemoteAddr)
		s.wsClients[r.RemoteAddr].Close()
	}
	s.wsClients[r.RemoteAddr] = conn
	s.wsClientsMutex.Unlock()
}

// Emit sends a WebSocket message to all connected WebSocket clients.
// Ensure that gob.Register is called on the message object on init.
// gob.Register does not need to be called when using a client, as it is
// automatically performed by HandleEvent.
func (s *Server) Emit(event string, msg interface{}) {
	type connection struct {
		id   string
		conn *websocket.Conn
	}

	var flaggedForRemoval []string
	var connections []connection

	s.wsClientsMutex.RLock()
	for id, conn := range s.wsClients {
		connections = append(connections, connection{
			id:   id,
			conn: conn,
		})
	}
	s.wsClientsMutex.RUnlock()

	for _, connPair := range connections {
		data, err := json.Marshal(encodeMessage{
			Event: event,
			Value: msg,
		})
		if err != nil {
			log.Println("ws: json encode error:", err)
			continue
		}

		err = connPair.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			flaggedForRemoval = append(flaggedForRemoval, connPair.id)
		}
	}

	s.wsClientsMutex.Lock()
	for _, id := range flaggedForRemoval {
		if _, found := s.wsClients[id]; found {
			delete(s.wsClients, id)
		}
	}
	s.wsClientsMutex.Unlock()
}

// HandleEvent registers an event handler to handle an incoming WebSocket
// message from the server. To be only used on the client.
func (c *Client) HandleEvent(event string, handler interface{}) {
	panicPrefix := "core: register handler: "

	if _, found := c.registeredHandlers[event]; found {
		panic(panicPrefix + "handler receiver for event \"" + event +
			"\" already registered")
	}

	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Func {
		panic(panicPrefix + "handler must be of type func")
	}

	if handlerType.NumIn() != 1 {
		panic(panicPrefix + "expected 1 input argument, " +
			"instead got " + strconv.Itoa(handlerType.NumIn()))
	}

	c.registeredHandlers[event] = registeredHandler{
		valueType: handlerType.In(0),
		handler:   reflect.ValueOf(handler),
	}

	if handlerType.NumOut() != 0 {
		panic(panicPrefix + "expected no return arguments, " +
			"instead got " + strconv.Itoa(handlerType.NumOut()))
	}
}

// HandleConnect calls the provided handler when the client successfully connects
// to the server. To be only used on the client.
func (c *Client) HandleConnect(handler func()) {
	c.connectHandler = handler
}

// HandleDisconnect calls the provided handler when the client disconnects
// from the server. There is no need to attempt to re-establish a connection
// as Connect() automatically does this. To be only used on the client.
func (c *Client) HandleDisconnect(handler func()) {
	c.disconnectHandler = handler
}

// Connect creates a WebSocket connection to the server to receive realtime
// events from. Connect runs a goroutine that automatically maintains the
// connection, and is non-blocking.
func (c *Client) Connect(url string) {
	for {
		conn, err := jsws.Dial(url)
		if err != nil {
			println("ws connect:", err)
			time.Sleep(time.Second * 3)
			continue
		}

		buffer := make([]byte, 100000)

		for {
			n, err := conn.Read(buffer)
			if err != nil {
				println("ws: error reading:", err)
				break
			}

			var msg decodeMessage
			err = json.Unmarshal(buffer[:n], &msg)
			if err != nil {
				println("ws: could not decode json:", err)
				continue
			}

			handler, found := c.registeredHandlers[msg.Event]
			if !found {
				println("ws: could not find handler for event:", msg.Event)
				continue
			}

			value := reflect.New(handler.valueType)
			err = json.Unmarshal([]byte(msg.Value), value.Interface())
			if err != nil {
				println("ws: could not decode value json for event \""+
					msg.Event+"\":", err)
				continue
			}

			go handler.handler.Call([]reflect.Value{value.Elem()})
		}
		conn.Close()
	}

}
