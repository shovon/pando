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

func CreateFailedToDeliverMessage(m MessageWithData, messageID string) MessageWithData {
	return MessageWithData{
		Type: "FAILED_TO_DELIVER_MESSAGE",
	}
}

type ParticipantDoesNotExist struct {
	ParticipantID string `json:"participantId"`
}

func CreateParticipantDoesNotExist(participantID string) MessageWithData {
	return MessageWithData{
		Type: "PARTICIPANT_DOES_NOT_EXIST",
		Data: ParticipantDoesNotExist{
			ParticipantID: participantID,
		},
	}
}

type ParticipantAuthenticating struct {
	ParticipantID string `json:"participantId"`
}

func CreateParticipantConnecting(participantID string) MessageWithData {
	return MessageWithData{
		Type: "PARTICIPANT_CONNECTING",
		Data: ParticipantAuthenticating{
			ParticipantID: participantID,
		},
	}
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

// RoomState is the state of a room
type RoomState struct {
	Participants []pairmap.KV[string, any] `json:"participants"`
}

// CreateRoomStateMessage creates a message containing the room state
func CreateRoomStateMessage(room RoomState) MessageWithData {
	return MessageWithData{
		Type: "ROOM_STATE",
		Data: room,
	}
}
