package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"reflect"
	"runtime/debug"
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

type subscription struct {
	valueType reflect.Type
	handler   reflect.Value
}

// Server is the object for a WebSocket server.
type Server struct {
	upgrader          websocket.Upgrader
	wsClients         map[string]*websocket.Conn
	subscribedClients map[string][]string
	wsClientsMutex    *sync.RWMutex
	subscribedMutex   *sync.RWMutex
}

// Client is the object for the WebSocket client connected to a WebSocket
// server.
type Client struct {
	subscriptions     map[string]subscription
	connectHandler    func()
	disconnectHandler func()
	conn              net.Conn
}

// SubscriptionMessage represents a subscription request from the client
// to the server.
type SubscriptionMessage struct {
	Event     string
	Subscribe bool
}

// NewServer returns a new server object, to be used on a server.
func NewServer() *Server {
	return &Server{
		upgrader:          websocket.Upgrader{},
		wsClients:         make(map[string]*websocket.Conn),
		subscribedClients: make(map[string][]string),
		wsClientsMutex:    new(sync.RWMutex),
		subscribedMutex:   new(sync.RWMutex),
	}
}

// NewClient returns a new client object, to be used on a client. It is
// not ever used for a server.
func NewClient() *Client {
	return &Client{
		subscriptions:     make(map[string]subscription),
		connectHandler:    func() {},
		disconnectHandler: func() {},
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

	go func(id string, conn *websocket.Conn) {
		for {
			var sub SubscriptionMessage
			err := conn.ReadJSON(&sub)
			if err != nil {
				s.subscribedMutex.Lock()
				for k, clients := range s.subscribedClients {
					var i int
					found := false
					for i = 0; i < len(clients); i++ {
						if clients[i] == id {
							found = true
							break
						}
					}

					if found {
						clients[i] = clients[len(clients)-1]
						clients = clients[:len(clients)-1]
						s.subscribedClients[k] = clients
					}
				}
				s.subscribedMutex.Unlock()

				s.wsClientsMutex.Lock()
				delete(s.wsClients, id)
				s.wsClientsMutex.Unlock()
				conn.Close()
				return
			}

			s.subscribedMutex.Lock()
			if sub.Subscribe {
				alreadySubscribed := false
				for _, client := range s.subscribedClients[sub.Event] {
					if client == id {
						alreadySubscribed = true
						break
					}
				}

				if !alreadySubscribed {
					s.subscribedClients[sub.Event] = append(
						s.subscribedClients[sub.Event], id)
				}
			} else {
				var i int
				found := false
				clients := s.subscribedClients[sub.Event]
				for i = 0; i < len(clients); i++ {
					if clients[i] == id {
						found = true
						break
					}
				}

				if found {
					clients[i] = clients[len(clients)-1]
					clients = clients[:len(clients)-1]
					s.subscribedClients[sub.Event] = clients
				}
			}
			s.subscribedMutex.Unlock()
		}
	}(r.RemoteAddr, conn)
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

	var connections []connection

	s.subscribedMutex.RLock()
	s.wsClientsMutex.RLock()
	clients := s.subscribedClients[event]
	for _, client := range clients {
		connections = append(connections, connection{
			id:   client,
			conn: s.wsClients[client],
		})
	}
	s.wsClientsMutex.RUnlock()
	s.subscribedMutex.RUnlock()

	for _, connPair := range connections {
		err := connPair.conn.WriteJSON(encodeMessage{
			Event: event,
			Value: msg,
		})
		if err != nil {
			log.Println("ws: write error:", err)
			continue
		}
	}
}

// Subscribe subscribes to and registers an event handler to handle an incoming
// WebSocket message from the server. All subscriptions are reset when the
// client disconnects, thus this should only be used in HandleConnect.
// To be only used on the client.
func (c *Client) Subscribe(event string, handler interface{}) {
	panicPrefix := "core: subscribe: "

	if _, found := c.subscriptions[event]; found {
		panic(panicPrefix + "subscription for event \"" + event +
			"\" already subscribed")
	}

	handlerType := reflect.TypeOf(handler)

	if handlerType.Kind() != reflect.Func {
		panic(panicPrefix + "handler must be of type func")
	}

	if handlerType.NumIn() != 1 {
		panic(panicPrefix + "expected 1 input argument, " +
			"instead got " + strconv.Itoa(handlerType.NumIn()))
	}

	c.subscriptions[event] = subscription{
		valueType: handlerType.In(0),
		handler:   reflect.ValueOf(handler),
	}

	if handlerType.NumOut() != 0 {
		panic(panicPrefix + "expected no return arguments, " +
			"instead got " + strconv.Itoa(handlerType.NumOut()))
	}

	data, _ := json.Marshal(SubscriptionMessage{
		Event:     event,
		Subscribe: true,
	})
	c.conn.Write(data)
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
	go func() {
		for {
			conn, err := jsws.Dial(url)
			c.conn = conn
			if err != nil {
				println("ws connect:", err)
				time.Sleep(time.Second * 3)
				continue
			}

			go func(c *Client) {
				defer func() {
					if r := recover(); r != nil {
						log.Println("ws: HandleConnect panic: " +
							fmt.Sprint(r) + "\n" + string(debug.Stack()))
					}
				}()

				c.connectHandler()
			}(c)
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

				handler, found := c.subscriptions[msg.Event]
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

				go func(handler subscription, value reflect.Value) {
					defer func() {
						if r := recover(); r != nil {
							log.Println("ws: subscription handler for for \"" +
								msg.Event + "\" panic: " + fmt.Sprint(r) +
								"\n" + string(debug.Stack()))
						}
					}()
				}(handler, value)

				handler.handler.Call([]reflect.Value{value.Elem()})
			}
			conn.Close()
			c.subscriptions = make(map[string]subscription)

			func(c *Client) {
				defer func() {
					if r := recover(); r != nil {
						log.Println("ws: HandleDisconnect panic: " +
							fmt.Sprint(r) + "\n" + string(debug.Stack()))
					}
				}()

				c.disconnectHandler()
			}(c)
		}
	}()
}
