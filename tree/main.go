package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

type connectionData struct {
	clientId string
}

type clientTree struct {
	clientId, treeId string
}

func getIds(r *http.Request, c *websocket.Conn) (clientTree, bool) {
	vars := mux.Vars(r)

	// Grab the client ID
	clientId := strings.TrimSpace(r.URL.Query().Get("client_id"))
	if len(clientId) <= 0 {
		log.Info().Msg("The client did not supply a client ID. Returning without client ID nor tree ID")
		msg, err := createClientError(
			"CLIENT_ID_NOT_SET",
			ErrorResponse{
				Title:  "Client ID has not been set",
				Detail: "Please set the relevant client ID, via the `client_id` query parameter",
			},
		)
		if err != nil {
			log.Panic().Err(err)
		}
		c.WriteJSON(msg)
		return clientTree{}, false
	}

	// Create a verifier
	panic(fmt.Sprintf("Not yet implemented. Waiting to handle clientId. Meanwhile check this out %s", clientId))

	// Grab the tree's ID
	id, ok := vars["id"]
	if !ok {
		log.Error().Err(errors.New("the tree ID was not set, for some reason. This is bad"))
		msg, err := createServerError("TREE_ID_NOT SET", ErrorResponse{
			Title:  "Tree ID not set",
			Detail: "A server error has resulted in the tree ID not being set on the server",
		})
		if err != nil {
			log.Panic().Err(err)
		}
		c.WriteJSON(msg)
		return clientTree{}, false
	}

	return clientTree{clientId: clientId, treeId: id}, false
}

func challengeClient(c *websocket.Conn) {
	payload := make([]byte, 32)
	rand.Read(payload)

	message := ChallengeMessage{
		Message: base64.RawStdEncoding.EncodeToString(payload),
	}

	c.WriteJSON(message)
}

type Sec256r1Key struct {
	X [32]byte
	Y [32]byte
}

func getSec256r1Key(str string) (Sec256r1Key, error) {
	b, err := base64.RawStdEncoding.DecodeString(str)
	if err != nil {
		return Sec256r1Key{}, err
	}
	if len(b) != 65 {
		return Sec256r1Key{}, fmt.Errorf("the key must be a 65-byte buffer, but but %d", len(b))
	}
	if b[0] != 4 {
		return Sec256r1Key{}, fmt.Errorf("the key must be an x and y coordinate key. This is indicated by the first byte of the key. The first byte must be of value 4, but got %d", len(b))
	}
	var x [32]byte
	var y [32]byte
	copy(x[:], b[1:33])
	copy(y[:], b[33:])

	return Sec256r1Key{X: x, Y: x}, nil
}

func handleTree(w http.ResponseWriter, r *http.Request) {

	// Upgarde the WebSocket connection
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err)
		return
	}

	// Get the IDs for the client and the tree
	ct, success := getIds(r, c)
	if !success {
		log.Info().Msg("Failed to get client and tree ID. Closing the connectiong")
		// No message is needed to be sent, since the `getIds` function already did
		// that.
		return
	}

	// Create the challenge
	challenge, err := createChallenge()
	if err != nil {
		log.Error().Err(err)
		return
	}

	c.WriteJSON(challenge)

	var msg MessageWithData

	for {

		err := c.ReadJSON(&msg)

		// This is an error arising from the fact that
		if err != nil {
			title := "Bad JSON message received"
			log.Info().Err(err).Msg(title)
			msg, err := createErrorResponse(ErrorResponse{Title: title})
			if err != nil {
				log.Panic().Err(err)
			}
			c.WriteJSON(msg)
			continue
		}

		if msg.Type != "CHALLENGE_RESPONSE" {
			title := fmt.Sprintf("Expected a message of type CHALLENGE_RESPONSE, but got %s", msg.Type)
			msg, err := createErrorResponse(ErrorResponse{Title: title})
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
			msg, err := createErrorResponse(ErrorResponse{Title: title})
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
			msg, err := createErrorResponse(ErrorResponse{Title: title})
			if err != nil {
				log.Panic().Err(err)
			}
			c.WriteJSON(msg)
			continue
		}

		verified, err := verifier.Verify(message, response.Signature)
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
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tree/{id}", handleTree).Methods("UPGRADE")
}
