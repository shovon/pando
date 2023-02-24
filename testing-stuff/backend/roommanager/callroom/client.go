package callroom

import (
	"backend/messages/servermessages"
	"encoding/json"

	"github.com/gorilla/websocket"
)

// Represents a single participant, not as far as the problem domain, but as a
// client in the call.
type Client struct {
	// The connection associated with the participant
	Connection *websocket.Conn

	// Participant is the metadata associated with the participant
	Participant servermessages.ParticipantState
}

var _ json.Marshaler = Client{}

func (c Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Participant)
}
