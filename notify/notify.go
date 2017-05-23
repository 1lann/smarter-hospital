// Package notify handles everything relating to alert notifications, and
// significant events, including audit logs.
package notify

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/1lann/smarter-hospital/store"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/ws"
)

// Notification represents a notification that is logged to the events
// log, and if Alert is true, then it will also push an alert to nurses.
// The ID, Dismissed and Time of the notification will be automatically set.
type Notification struct {
	ID bson.ObjectId `bson:"_id"`

	Alert     bool
	Dismissed bool

	Heading    string
	SubHeading string
	Icon       string
	Link       string

	Time time.Time

	Target string
}

// Server represents a notification server.
type Server struct {
	ws *ws.Server
}

// Client represents a client that receives notifications.
type Client struct {
	ws             *ws.Client
	notifyHandler  func(n Notification)
	dismissHandler func(id string)
}

// NewServer returns a new notification server with the given WebSocket server.
func NewServer(wsServer *ws.Server) *Server {
	return &Server{ws: wsServer}
}

// Push pushes a new notification, and notifies clients appropriately.
func (s *Server) Push(n Notification) error {
	store.C("notify").EnsureIndexKey("dismissed")
	store.C("notify").EnsureIndexKey("time")

	n.Dismissed = !n.Alert
	n.Time = time.Now()
	n.ID = bson.NewObjectIdWithTime(time.Now())

	err := store.C("notify").Insert(n)
	if err != nil {
		return err
	}

	s.ws.Emit("notification", n)

	return nil
}

// Dismiss dismisses a notification across all clients.
func (s *Server) Dismiss(id string) error {
	err := store.C("notify").UpdateId(id, bson.M{"dismissed": true})
	if err != nil {
		return err
	}

	s.ws.Emit("notification_dismiss", id)
	return nil
}

// Notifications returns all the notifications in the server.
func (s *Server) Notifications() ([]Notification, error) {
	var results []Notification
	err := store.C("notify").Find(nil).Sort("-time").All(&results)
	return results, err
}

// NewClient returns a new notification client which receives notifications.
func NewClient(wsClient *ws.Client) *Client {
	return &Client{ws: wsClient}
}

func getNotifications() ([]Notification, error) {
	resp, err := http.Get(views.Address + "/notify/all")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return nil, views.ErrInternalError
	case http.StatusOK:
		break
	default:
		return nil, errors.New("notify: unknown response")
	}

	var notifications []Notification
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&notifications)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// OnNotification is used to provide a handler whenever a notification
// is received.
func (c *Client) OnNotification(handler func(n Notification)) {
	c.notifyHandler = handler
}

// OnDismiss is used to provide a handler whenever a dismissal is received.
func (c *Client) OnDismiss(handler func(id string)) {
	c.dismissHandler = handler
}

// Start returns the list of current notifications, and starts listening for
// notifications.
func (c *Client) Start() ([]Notification, error) {
	c.ws.Subscribe("notification", c.notifyHandler)
	c.ws.Subscribe("notification_dismiss", c.dismissHandler)

	ns, err := getNotifications()
	if err != nil {
		return nil, err
	}

	return ns, nil
}

// Dismiss dismisses a notification, and notifies the server of such dismissal.
// The handler provided to OnDismiss will be called when the server
// successfully receives the dismissal.
func (c *Client) Dismiss(id string) {
	resp, err := http.Get(views.Address + "/notify/dismiss/" + id)
	if err != nil {
		println("notify: dismiss:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		println("notify: dismiss:", resp.Status)
	}
}
