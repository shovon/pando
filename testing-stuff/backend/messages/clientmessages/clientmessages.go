package clientmessages

import "encoding/json"

// Message is a message sent from a client to the server
type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type MessageToParticipant struct {
	To   string          `json:"type"`
	Data json.RawMessage `json:"data"`
	ID   string          `json:"id"`
}

type SessionToken string

type UnknownMessage json.RawMessage

func ParseMessageToParticipant(message json.RawMessage) (MessageToParticipant, error) {
	var m MessageToParticipant
	err := json.Unmarshal(message, &m)
	return m, err
}

func ParseSessionToken(message json.RawMessage) (SessionToken, error) {
	var m SessionToken
	err := json.Unmarshal(message, &m)
	return m, err
}

func ParseParticipantName(message json.RawMessage) (string, error) {
	var name string
	err := json.Unmarshal(message, &name)
	return name, err
}

const (
	MessageToParticipantType = "MESSAGE_TO_PARTICIPANT"
	SetNameType              = "SET_NAME"
)

func ParseMessage(message Message) (any, error) {
	switch message.Type {
	case "MESSAGE_TO_PARTICIPANT":
		return ParseMessageToParticipant(message.Data)
	case "SET_NAME":
		return ParseParticipantName(message.Data)
	case "BROADCAST_MESSAGE":
	case "ENABLE_VIDEO":
	case "DISABLE_VIDEO":
	case "ENABLE_AUDIO":
	case "DISABLE_AUDIO":
	case "CLOSE_CONNECTION":
		return ParseSessionToken(message.Data)
	case "RESTORE_STATE":
	}

	return UnknownMessage(message.Data), nil
}
