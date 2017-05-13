package store

import "gopkg.in/mgo.v2"

// ConnectOpts are the parameters for connecting to the database.
type ConnectOpts struct {
	Address  string
	Database string
	Username string
	Password string
}

// Session represents a connected database session.
type Session struct {
	session *mgo.Session
	db      *mgo.Database
}

var session *Session

// Connect attempts a connection to the database with the provided connection
// options, and returns a session if successful.
func Connect(options ConnectOpts) error {
	s, err := mgo.Dial(options.Address)
	if err != nil {
		return err
	}

	// Optional. Switch the session to a monotonic behavior.
	s.SetMode(mgo.Monotonic, true)

	db := s.DB(options.Database)

	session = &Session{
		session: s,
		db:      db,
	}

	return nil
}

// C returns a collection of the connected database session.
func C(name string) *mgo.Collection {
	return session.db.C(name)
}
