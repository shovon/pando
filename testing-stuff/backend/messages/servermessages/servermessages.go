package servermessages

import (
	"backend/pairmap"
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

// CreateClientError creates a client error message, in order to notify the
// client that it has sent a message that the server cannot understand
func CreateClientError(err ErrorResponse) MessageWithData {
	return MessageWithData{
		Type: "CLIENT_ERROR",
		Data: err,
	}
}

// CreateServerError creates a server error message
func CreateServerError(err ErrorResponse) MessageWithData {
	return MessageWithData{
		Type: "SERVER_ERROR",
		Data: err,
	}
}

// CreateMessageToParticipant creates a message to be sent to a participant
func CreateMessageToParticipant(from string, message json.RawMessage) MessageWithData {
	return MessageWithData{
		Type: "MESSAGE_FROM_PARTICIPANT",
		Data: MessageToParticipant{
			From: from,
			Data: message,
		},
	}
}

// ParticipantState is the state of a participant
type ParticipantState struct {
	// Name is the name of the participant
	Name string `json:"name"`

	// HasVideo is true if the participant has video enabled
	HasVideo bool `json:"hasVideo"`

	// HasAudio is true if the participant has audio enabled
	HasAudio bool `json:"hasAudio"`
}

// ROomState is the state of a room
type RoomState struct {
	Participants []pairmap.KV[string, ParticipantState] `json:"participants"`
}

func CreateRoomStateMessage(room RoomState) MessageWithData {
	return MessageWithData{
		Type: "ROOM_STATE",
		Data: room,
	}
}
