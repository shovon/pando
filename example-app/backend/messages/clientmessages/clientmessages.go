package clientmessages

import "encoding/json"

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type MessageToParticipant struct {
	To   string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func ParseMessageToParticipant(message json.RawMessage) (MessageToParticipant, error) {
	var m MessageToParticipant
	err := json.Unmarshal(message, &m)
	return m, err
}
