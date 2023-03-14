package connectionstate

import (
	"backend/writer"
	"io"
)

const (
	ConnectingState = "CONNECTING"
	ConnectedState  = "CONNECTED"
)

// Connecting is the state when the connection is being authenticated
type Connecting struct{}

type CloserWriter interface {
	io.Closer
	writer.Writer
}

// Connected is the state when the connection is connected, and also gives us
// access to methods for sending messages to the client
type Connected struct {
	writer CloserWriter
}

var _ writer.Writer = Connected{}
var _ io.Closer = Connected{}

func (c Connected) Write(message []byte) error {
	return c.writer.Write(message)
}

func (c Connected) WriteJSON(message interface{}) error {
	return c.writer.WriteJSON(message)
}

func (c Connected) Close() error {
	return c.writer.Close()
}

func ConnectionStatus(state any) string {
	switch state.(type) {
	case Connecting:
		return ConnectingState
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

// NewAuthenticatingConnection creates a new connection in the authenticating
// state
func NewAuthenticatingConnection() Connection {
	return Connection{state: Connecting{}}
}

// NewConnectedConnection creates a new connection in the connected state
func NewConnectedConnection(w CloserWriter) Connection {
	return Connection{state: Connected{writer: w}}
}

// State returns the state of the connection
func (c Connection) State() any {
	return c.state
}
