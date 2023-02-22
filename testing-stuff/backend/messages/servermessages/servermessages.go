package servermessages

import (
	"encoding/json"
)

type MessageWithData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	ID     string      `json:"id,omitempty"`
	Code   string      `json:"code,omitempty"`
	Title  string      `json:"title,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

type MessageToParticipant struct {
	From string          `json:"from"`
	Data json.RawMessage `json:"data"`
}

func CreateClientError(err ErrorResponse) MessageWithData {
	return MessageWithData{
		Type: "CLIENT_ERROR",
		Data: err,
	}
}

func CreateServerError(err ErrorResponse) MessageWithData {
	return MessageWithData{
		Type: "SERVER_ERROR",
		Data: err,
	}
}

func CreateMessageToParticipant(from string, message json.RawMessage) MessageWithData {
	return MessageWithData{
		Type: "MESSAGE_FROM_PARTICIPANT",
		Data: MessageToParticipant{
			From: from,
			Data: message,
		},
	}
}

type ParticipantState struct {

	// HasVideo is true if the participant has video enabled
	HasVideo bool

	// HasAudio is true if the participant has audio enabled
	HasAudio bool
}

type ParticipantKeyValuePair []struct {
	Key   string
	Value ParticipantState
}

var _ json.Marshaler = ParticipantKeyValuePair{}

func (p ParticipantKeyValuePair) MarshalJSON() ([]byte, error) {
	m := [][]interface{}{}

	for _, keyValue := range p {
		tuple := []interface{}{keyValue.Key, keyValue.Value}
		m = append(m, tuple)
	}

	return json.Marshal(m)
}

type RoomState struct {
	Participants map[string]ParticipantState `json:"participants"`
}

func CreateRoomStateMessage(room RoomState) MessageWithData {
	return MessageWithData{
		Type: "ROOM_STATE",
		Data: room,
	}
}
