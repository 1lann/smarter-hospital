package views

import (
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

// Message represents a message from a WebSocket.
type Message struct {
	Type string
	Data []byte
}

// Page represents a page.
type Page interface {
	OnLoad()
	OnUnload(client *ws.Client)
	OnConnect(client *ws.Client)
	OnDisconnect()
	Title() string
}

// User represents information about the current user.
type User struct {
	*js.Object
	FirstName string `js:"firstName"`
	LastName  string `js:"lastName"`
}

var registeredPages = make(map[string]Page)

// RegisterPage registers a page to the pages system.
func RegisterPage(path string, page Page) {
	if _, found := registeredPages[path]; found {
		panic("pages: page already exists for path: " + path)
	}

	registeredPages[path] = page
}

// GetTitle returns the title of a registered page.
func GetTitle(path string) string {
	page, found := registeredPages[path]
	if !found {
		return ""
	}

	return page.Title()
}

// AllPages returns a map of all the pages. Do not modify/write to the map.
func AllPages() map[string]Page {
	return registeredPages
}

// Run starts the page handling system.
func Run() {
	path := js.Global.Get("location").Get("pathname").String()
	page, found := registeredPages[path]
	if !found {
		panic("page not found!")
	}

	client := ws.NewClient()
	var scheme string
	if js.Global.Get("location").Get("protocol").String() == "https:" {
		scheme = "wss://"
	} else {
		scheme = "ws://"
	}

	client.HandleConnect(func() {
		page.OnConnect(client)
	})
	client.HandleDisconnect(func() {
		page.OnDisconnect()
	})
	client.Connect(scheme + js.Global.Get("location").Get("host").String() +
		"/ws")

	page.OnLoad()

	select {}
}

// GetUser returns the user information from the page.
func GetUser() User {
	var user User
	user.Object = js.Global.Get("user")
	return user
}
