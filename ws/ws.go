package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var wsConnected = make(map[string]*websocket.Conn)
var wsConnectedMutex = new(sync.RWMutex)

// Message represents a message that can be sent over WebSockets.
type Message struct {
	Type        string
	Name        string      `json:",omitempty"` // For names of displays, and events (or short title)
	Value       interface{} `json:",omitempty"` // Generic value
	Description string      `json:",omitempty"` // Longer description of the message
	Error       error       `json:",omitempty"` // Error value
}

// Handle is the handler to be called by the HTTP server to handle
// for websocket connections.
func Handle(r *http.Request, wr http.ResponseWriter) {
	conn, err := upgrader.Upgrade(wr, r, nil)
	if err != nil {
		log.Println("handleWS upgrade:", err)
		return
	}

	wsConnectedMutex.Lock()
	if _, exists := wsConnected[r.RemoteAddr]; exists {
		log.Println("server: handleWS: already existing connection:",
			r.RemoteAddr)
		wsConnected[r.RemoteAddr].Close()
	}
	wsConnected[r.RemoteAddr] = conn
	wsConnectedMutex.Unlock()
}

// Emit sends a websocket message to all connected websocket clients.
func Emit(msg Message) {
	type connection struct {
		id   string
		conn *websocket.Conn
	}

	var flaggedForRemoval []string
	var connections []connection

	wsConnectedMutex.RLock()
	for id, conn := range wsConnected {
		connections = append(connections, connection{
			id:   id,
			conn: conn,
		})
	}
	wsConnectedMutex.RUnlock()

	for _, connPair := range connections {
		err := connPair.conn.WriteJSON(msg)
		if err != nil {
			flaggedForRemoval = append(flaggedForRemoval, connPair.id)
		}
	}

	wsConnectedMutex.Lock()
	for _, id := range flaggedForRemoval {
		if _, found := wsConnected[id]; found {
			delete(wsConnected, id)
		}
	}
	wsConnectedMutex.Unlock()
}
