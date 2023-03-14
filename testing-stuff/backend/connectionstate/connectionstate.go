package connectionstate

import (
	"backend/writer"
	"io"
)

const (
	DisconnectedState = "DISCONNECTED"
	ConnectedState    = "CONNECTED"
)

// Disconnected is the state when the connection is being authenticated
type Disconnected struct{}

type CloserWriter struct {
	io.Closer
	writer.Writer
}

// Connected is the state when the connection is connected, and also gives us
// access to methods for sending messages to the client
type Connected struct {
	writer CloserWriter
}

var _ io.Closer = Connected{}
var _ writer.Writer = Connected{}

func (c Connected) Close() error {
	return c.writer.Close()
}

func (c Connected) Write(message []byte) error {
	return c.writer.Write(message)
}

func (c Connected) WriteJSON(message interface{}) error {
	return c.writer.WriteJSON(message)
}

// ConnectionStatus infers the connection status given the passed-in state
// object
func ConnectionStatus(state any) string {
	switch state.(type) {
	case Disconnected:
		return DisconnectedState
	case Connected:
		return ConnectedState
	default:
		return "UNKNOWN"
	}
}

// Connection is just a safe connection object that can be used to send messages
type Connection struct {
	state any
}

// NewDisconnectedStatus creates a new connection in the authenticating
// state
func NewDisconnectedStatus() Connection {
	return Connection{state: Disconnected{}}
}

// NewConnectedStatus creates a new connection in the connected state
func NewConnectedStatus(w CloserWriter) Connection {
	return Connection{state: Connected{writer: w}}
}

// State returns the state of the connection
func (c Connection) State() any {
	return c.state
}
