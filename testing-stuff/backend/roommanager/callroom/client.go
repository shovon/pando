package callroom

import (
	"backend/ws"
	"encoding/json"
)

const (
	AwaitingConnection = "AWAITING_CONNECTION"
	Connected          = "CONNECTED"
)

// ParticipantState is the state of a participant
type ParticipantState struct {
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
	// The connection associated with the participant
	WebSocketWriter ws.ThreadSafeWriter

	// Participant is the metadata associated with the participant
	Participant ParticipantState
}

var _ json.Marshaler = Client{}

// ConnectionStatus returns the connection status of the participant
func (c Client) ConnectionStatus() string {
	if !c.WebSocketWriter.IsConnected() {
		return AwaitingConnection
	}
	return Connected
}

func (c Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Participant)
}
