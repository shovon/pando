package servermessages

import (
	"backend/keyvalue"
)

// TODO: having all server messages in here is just stupid.
//
//   Refactoring is needed to standardize the concept of a type/data message,
//   while still empowering domain models to send messages without relying on
//   this module. Again inversion of control is quite important

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

type ParticipantDoesNotExist struct {
	ParticipantID string `json:"participantId"`
}

type ParticipantAuthenticating struct {
	ParticipantID string `json:"participantId"`
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

// RoomState is the state of a room
type RoomState struct {
	Participants []keyvalue.KV[string, any] `json:"participants"`
}

// CreateRoomStateMessage creates a message containing the room state
func CreateRoomStateMessage(room RoomState) MessageWithData {
	return MessageWithData{
		Type: "ROOM_STATE",
		Data: room,
	}
}

func CreateSessionTokenMessage(token string) MessageWithData {
	return MessageWithData{
		Type: "SESSION_TOKEN",
		Data: token,
	}
}
