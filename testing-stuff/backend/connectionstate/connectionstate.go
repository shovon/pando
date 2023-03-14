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

// NewConnectingStatus creates a new connection in the authenticating
// state
func NewConnectingStatus() Connection {
	return Connection{state: Connecting{}}
}

// NewConnectedStatus creates a new connection in the connected state
func NewConnectedStatus(w writer.Writer) Connection {
	return Connection{state: Connected{writer: w}}
}

// State returns the state of the connection
func (c Connection) State() any {
	return c.state
}
