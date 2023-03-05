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

var _ writer.Writer = Connected{}

func (c Connected) Write(message []byte) error {
	return c.writer.Write(message)
}

func (c Connected) WriteJSON(message interface{}) error {
	return c.writer.WriteJSON(message)
}

// Disconnected is the state when the connection is disconnected
type Disconnected struct{}

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
func (c Connection) State() any {
	return c.state
}

func (c *Connection) Disconnect() {
	c.state = NewDisconnectedConnection()
}
