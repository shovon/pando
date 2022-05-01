package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"tree/keyid"

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

func sendChallenge() {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tree/{id}/{clientid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		clientId, ok := vars["clientid"]
		if !ok {
			log.Err(errors.New("the user ID was not set, for some reason. This is bad"))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse User ID"))
		}

		verifier, err := keyid.CreateVerifier(clientId)
		if err != nil {
			log.Error().Err(fmt.Errorf("failed to parse the user ID as a valid key", clientId))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse User ID"))
		}
		if verifier.IsKeyValid() {
			w.WriteHeader(400)
			w.Write([]byte("The user ID must be a valid key format. Remember: tbe user ID will double as a public key"))
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

		challenge, err := createChallenge()
		if err != nil {
			log.Error().Err(err)
			w.WriteHeader(500)
			w.Write([]byte("Failed to create challenge message. Investigation is needed"))
		}

		c.WriteJSON(challenge)

		var msg Message

		for {

			err := c.ReadJSON(&msg)

			if err != nil {
				title := "Bad JSON message received"
				log.Info().Err(err).Msg(title)
				msg, err := createErrorResponse(ErrorResponse{Title: &title})
				if err != nil {
					log.Panic().Err(err)
				}
				c.WriteJSON(msg)
				continue
			}

			if msg.Type != "CHALLENGE_RESPONSE" {
				title := fmt.Sprintf("Expected a message of type CHALLENGE_RESPONSE, but got %s", msg.Type)
				msg, err := createErrorResponse(ErrorResponse{Title: &title})
				if err != nil {
					log.Panic().Err(err)
				}
				c.WriteJSON(msg)
				continue
			}

			var response ChallengeResponse

			err = json.Unmarshal(msg.Data, &response)
			if err != nil {
				title := "Bad challenge response payload given"
				log.Info().Err(err).Msg(title)
				msg, err := createErrorResponse(ErrorResponse{Title: &title})
				if err != nil {
					log.Panic().Err(err)
				}
				c.WriteJSON(msg)
				continue
			}

			message, err := base64.RawStdEncoding.DecodeString(response.Message)
			if err != nil {
				title := "Failed to parse base64-encoded message"
				log.Info().Err(err).Msg(title)
				msg, err := createErrorResponse(ErrorResponse{Title: &title})
				if err != nil {
					log.Panic().Err(err)
				}
				c.WriteJSON(msg)
				continue
			}

			signature, err := base64.RawStdEncoding.DecodeString(response.Signature)
			if err != nil {
				title := "Failed to parse base64-encoded signature"
				log.Info().Err(err).Msg(title)
				msg, err := createErrorResponse(ErrorResponse{Title: &title})
				if err != nil {
					log.Panic().Err(err)
				}
				c.WriteJSON(msg)
				continue
			}

			verified, err := verifier.Verify(message, signature)
			if err != nil {
				title := "An internal error occurred while attempting to verify the challenge response and signature"
				log.Info().Err(err).Msg(title)
				msg, err := createErrorResponse(ErrorResponse{Title: &title})
				if err != nil {
					log.Panic().Err(err)
					c.WriteJSON(msg)
					continue
				}
			}

			if !verified {
				title := "The signature has not been deemed as authentic"
				log.Info().Err(err).Msg(title)
				msg, err := createErrorResponse(ErrorResponse{Title: &title})
				if err != nil {
					log.Panic().Err(err)
					c.WriteJSON(msg)
					continue
				}
			}

		}

		// go func() {
		// 	switch ev := (<-listener).(type) {
		// 	case spanningtree.NodeState:

		// 	}
		// }()
	})
}
