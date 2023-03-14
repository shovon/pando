package callroom

import (
	"backend/connectionstate"
	"encoding/json"
)

// ParticipantState is the state of a participant
type ParticipantState struct {
	// TODO: a lot of this stuff can easily be soft-coded

	// Name is the name of the participant
	Name string `json:"name"`

	// HasVideo is true if the participant has video enabled
	HasVideo bool `json:"hasVideo"`

	// HasAudio is true if the participant has audio enabled
	HasAudio bool `json:"hasAudio"`
}

// Represents a single participant, not as far as the problem domain, but as a
// client in the call.
type Client struct {
	// Connection is the WebSocket writer is used to send messages to the client.
	// This exists to prevent race conditions when sending messages to the client.
	// A race condition will either result in the
	Connection connectionstate.Connection

	// SessionToken is the token for some asynchronous/stateless authenticated
	// stuff
	SessionToken string

	// Participant is the metadata associated with the participant
	Participant ParticipantState
}

var _ json.Marshaler = Client{}

func (c *Client) Close() error {
	con, ok := c.Connection.State().(connectionstate.Connected)
	if !ok {
		// If already disconnected, it's a no-op
		return nil
	}
	return con.Close()
}

// ConnectionState returns the connection status of the participant
//
// This function has been created because not all participants on the call are
// guaranteed to be connected to the server. For example, if the server crashes
// and then the room is re-created, the participants that were in the room
// before the crash will be re-inserted into the room, but they will not be
// connected to the server.
func (c Client) ConnectionState() any {
	return c.Connection.State()
}

func (c Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Participant)
}
