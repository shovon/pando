package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

type Message interface {
	MessageWithData | MessageNoData
}

type MessageWithData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type MessageNoData struct {
	Type string `json:"type"`
}

type ChallengeMessage struct {
	Message string `json:"message"`
}

type ChallengeResponse struct {
	Message   string          `json:"message"`
	Signature json.RawMessage `json:"signature"`
}

func createMessage(title string, data interface{}) (MessageWithData, error) {
	msg, err := json.Marshal(data)
	if err != nil {
		return MessageWithData{}, err
	}
	return MessageWithData{title, msg}, nil
}

func createChallenge() (MessageWithData, error) {
	payload := make([]byte, 32)
	rand.Read(payload)
	msg := ChallengeMessage{
		Message: base64.RawStdEncoding.EncodeToString(payload),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return MessageWithData{}, err
	}
	return MessageWithData{
		Type: "CHALLENGE",
		Data: data,
	}, nil
}

func createClientError(title string, payload interface{}) (MessageWithData, error) {
	data, err := createMessage(title, payload)
	if err != nil {
		return MessageWithData{}, err
	}

	return createMessage("CLIENT_ERROR", data)
}

func createServerError(title string, payload interface{}) (MessageWithData, error) {
	data, err := createMessage(title, payload)
	if err != nil {
		return MessageWithData{}, err
	}

	return createMessage("SERVER_ERROR", data)
}

type ErrorResponse struct {
	ID     string      `json:"id,omitempty"`
	Code   string      `json:"code,omitempty"`
	Title  string      `json:"title,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

func createErrorResponse(err ErrorResponse) (MessageWithData, error) {
	data, e := json.Marshal(err)
	if e != nil {
		return MessageWithData{}, e
	}

	return MessageWithData{
		Type: "ERROR",
		Data: data,
	}, nil
}
