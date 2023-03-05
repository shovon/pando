package connectionstate

import "backend/writer"

// Authenticating is the state when the connection is being authenticated
type Authenticating struct{}

// Connected is the state when the connection is connected, and also gives us
// access to methods for sending messages to the client
type Connected struct {
	writer writer.Writer
}

// Disconnected is the state when the connection is disconnected
type Disconnected struct{}

// Connection is just a safe connection object that can be used to send messages
type Connection struct {
	state any
}

var _ writer.Writer = Connection{}

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

// Write writes a message to the connection
func (c Connection) Write(message []byte) error {
	writer, ok := c.state.(Connected)

	if !ok {
		// TODO: determine whether an unconnected state should be a no-op
		return nil
	}

	return writer.writer.Write(message)
}

// WriteJSON writes a JSON message to the connection
func (c Connection) WriteJSON(message interface{}) error {
	writer, ok := c.state.(Connected)

	if !ok {
		// TOOD: determine whether an unconnected state should be a no-op
		return nil
	}

	return writer.writer.WriteJSON(message)
}
