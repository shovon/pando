package connectionstate

import "backend/writer"

const (
	AuthenticatingState = "AUTHENTICATING"
	ConnectedState      = "CONNECTED"
	DisconnectedState   = "DISCONNECTED"
)

// Authenticating is the state when the connection is being authenticated
type Authenticating struct{}

// Connected is the state when the connection is connected, and also gives us
// access to methods for sending messages to the client
type Connected struct {
	writer writer.Writer
}

// Disconnected is the state when the connection is disconnected
type Disconnected struct{}

// TODO: perhaps letting the Connection object be a writer is a bad idea

// Connection is just a safe connection object that can be used to send messages
type Connection struct {
	state any
}

// NewAuthenticatingConnection creates a new connection in the authenticating
// state
func NewAuthenticatingConnection() Connection {
	return Connection{state: Authenticating{}}
}

// NewConnectedConnection creates a new connection in the connected state
func NewConnectedConnection(w writer.Writer) Connection {
	return Connection{state: Connected{writer: w}}
}

// NewDisconnectedConnection creates a new connection in the disconnected state
func NewDisconnectedConnection() Connection {
	return Connection{state: Disconnected{}}
}

// State returns the state of the connection
func (c Connection) State() string {
	switch c.state.(type) {
	case Authenticating:
		return AuthenticatingState
	case Connected:
		return ConnectedState
	case Disconnected:
		return DisconnectedState
	}

	panic("Unknown connection state")
}

func (c *Connection) Disconnect() {
	c.state = NewDisconnectedConnection()
}
