package callroom

import (
	"backend/messages/servermessages"
	"backend/ws"
	"encoding/json"
)

const (
	AwaitingConnection = "AWAITING_CONNECTION"
	Connected          = "CONNECTED"
)

// Represents a single participant, not as far as the problem domain, but as a
// client in the call.
type Client struct {
	// The connection associated with the participant
	WebSocketWriter ws.ThreadSafeWriter

	// Participant is the metadata associated with the participant
	Participant servermessages.ParticipantState
}

var _ json.Marshaler = Client{}

func (c Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Participant)
}
