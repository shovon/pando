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

type ParticipantState struct{}

type RoomState struct {
	Participants map[string]ParticipantState `json:"participants"`
}

func CreateRoomStateMessage(room RoomState) MessageWithData {
	return MessageWithData{
		Type: "ROOM_STATE",
		Data: room,
	}
}
