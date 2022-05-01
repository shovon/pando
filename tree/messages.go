package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ChallengeMessage struct {
	Message string `json:"message"`
}

type ChallengeResponse struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func createChallenge() (Message, error) {
	payload := make([]byte, 32)
	rand.Read(payload)
	msg := ChallengeMessage{
		Message: base64.RawStdEncoding.EncodeToString(payload),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return Message{}, err
	}
	return Message{
		Type: "CHALLENGE",
		Data: data,
	}, nil
}

type ErrorResponse struct {
	ID     *string     `json:"id,omitempty"`
	Code   *string     `json:"code,omitempty"`
	Title  *string     `json:"title,omitempty"`
	Detail *string     `json:"detail,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

func createErrorResponse(err ErrorResponse) (Message, error) {
	data, e := json.Marshal(err)
	if e != nil {
		return Message{}, e
	}

	return Message{
		Type: "ERROR",
		Data: data,
	}, nil
}
