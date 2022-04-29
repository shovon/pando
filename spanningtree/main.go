package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"spanningtree/key"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{}

var trees = newTreeManager()

type participant struct {
	conn *websocket.Conn
	meta json.RawMessage
}

var _ json.Marshaler = &participant{}

func (p *participant) MarshalJSON() ([]byte, error) {
	return p.meta, nil
}

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
	ID     *string     `json:"id",omitempty"`
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tree/{id}/{userid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		userId, ok := vars["userid"]
		if !ok {
			log.Err(errors.New("the user ID was not set, for some reason. This is bad"))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse User ID"))
		}

		verifier, err := key.CreateVerifier(userId)
		if err != nil {
			log.Error().Err(fmt.Errorf("failed to parse the user ID as a valid key", userId))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse User ID"))
		}
		if verifier.IsKeyValid() {
			w.WriteHeader(400)
			w.Write([]byte("The user ID must be a valid key format. Remember: tbe user ID will double as a public key"))
		}

		challenge, err := createChallenge()
		if err != nil {
			log.Error().Err(err)
			w.WriteHeader(500)
			w.Write([]byte("Failed to create challenge message. Investigation is needed"))
		}

		id, ok := vars["id"]
		if !ok {
			log.Error().Err(errors.New("the ID was not set, for some reason. This is bad"))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse ID"))
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err)
			w.WriteHeader(500)
			w.Write([]byte("An internal server error occurred"))
			return
		}

		// t := trees.getTree(id)
		// listener := t.RegisterChangeListener(userId)

		c.WriteJSON(challenge)

		var msg Message

		for {

			err := c.ReadJSON(&msg)

			if err != nil {
				log.Info().Err(err).Msg("Bad JSON message received")
			}

			if msg.Type != "CHALLENGE_RESPONSE" {
				title := fmt.Sprintf("Expected a message of type CHALLENGE_RESPONSE, but got %s", msg.Type)
				msg, err := createErrorResponse(ErrorResponse{Title: &title})
				if err != nil {
					log.Panic().Err(err)
				}
				c.WriteJSON(msg)
				c.Close()
				return
			}

			var response ChallengeResponse

			err := json.Unmarshal(msg.Data, &response)
			if err != nil {
				log.Info().Err(err).Msg("Bad challenge response payload given")
				msg, err := createErrorResponse(ErrorResponse{})

			}

		}

		// go func() {
		// 	switch ev := (<-listener).(type) {
		// 	case spanningtree.NodeState:

		// 	}
		// }()
	})
}
