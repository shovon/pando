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
	Payload string `json:"payload"`
}

func createChallenge() (Message, error) {
	payload := make([]byte, 32)
	rand.Read(payload)
	msg := ChallengeMessage{
		Payload: base64.RawStdEncoding.EncodeToString(payload),
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
			log.Err(fmt.Errorf("failed to parse the user ID as a valid key", userId))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse User ID"))
		}
		if verifier.IsKeyValid() {
			w.WriteHeader(400)
			w.Write([]byte("The user ID must be a valid key format. Remember: tbe user ID will double as a public key"))
		}

		challenge, err := createChallenge()
		if err != nil {
			log.Err(err)
			w.WriteHeader(500)
			w.Write([]byte("Failed to create challenge message. Investigation is needed"))
		}

		id, ok := vars["id"]
		if !ok {
			log.Err(errors.New("the ID was not set, for some reason. This is bad"))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse ID"))
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Err(err)
			w.WriteHeader(500)
			w.Write([]byte("An internal server error occurred"))
			return
		}

		// t := trees.getTree(id)
		// listener := t.RegisterChangeListener(userId)

		c.WriteJSON(challenge)

		// Hm… Should we be conservative, and close the connection if the client
		// does not send a challenge response as the first message, or should we
		// be liberal, and loop until we get a valid response?

		// go func() {
		// 	switch ev := (<-listener).(type) {
		// 	case spanningtree.NodeState:

		// 	}
		// }()
	})
}
