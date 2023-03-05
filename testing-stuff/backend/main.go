package main

import (
	"backend/config"
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"backend/roommanager"
	"backend/ws"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sparkscience/wskeyid-golang"
)

var rooms = roommanager.NewRoomManager()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// handleRoom is the event handler for the room endpoint.
//
// This is where when there is a connection to a room, the participant will
// interact with the room
func handleRoom(w http.ResponseWriter, r *http.Request) {
	// Several notes about a connection.
	//
	// - can be assigned (usually upon first connection)
	// - can be reassigned (when a participant jumps from one Internet connection
	//   to another)
	// - can be unassigned (when the server crashes and restarts, and the last
	//   room state is loaded)
	//   - in this situation, we give the participant 60 seconds to connect before
	//     officially kicking them out of the room
	//   - (Note: there might also be other situations where a connection is
	//     unassigned)

	log.Print("Got connection from client")

	// Grab the list of parameters
	params := mux.Vars(r)

	// HTTP -> WebSocket
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer c.Close()
	defer fmt.Println("Connection ended")

	log.Println("Got connection object")

	// Authenticate
	{
		err := wskeyid.HandleAuthConnection(r, c)
		if err != nil {
			log.Println("WebSocket authentication failed ", err.Error())
			return
		}
	}

	log.Println("Got connection wrapper")

	// Get the client ID from the URL
	clientId := strings.TrimSpace(r.URL.Query().Get("client_id"))

	// Get the room ID from the URL
	roomId, ok := params["id"]
	if !ok {
		// This should have technically not been possible at all. Thus closing the
		// connection, while also notifying the client that something went wrong.
		c.WriteJSON(
			servermessages.
				CreateServerError(
					servermessages.ErrorResponse{Title: "An internal server error"},
				),
		)
		return
	}

	// Wait until the participant provides a name
	messageChannel := ws.ReadLoop(c)

	log.Println("Waiting for name")

	var name string

	attempts := 0

	for {
		if attempts >= 10 {
			c.WriteJSON(servermessages.CreateClientError(servermessages.ErrorResponse{
				Title: "Too many failed attempts at providing a name closing",
			}))
			return
		}

		event, ok := <-messageChannel

		if !ok {
			return
		}

		var message clientmessages.Message
		err := json.Unmarshal(event, &message)
		if err != nil {
			continue
		}

		if message.Type == "SET_NAME" {
			// Got the participant name
			n, err := clientmessages.ParseParticipantName(message.Data)
			if err != nil {
				log.Println("Error parsing participant name: ", err.Error())
				c.WriteJSON(servermessages.CreateClientError(servermessages.ErrorResponse{
					Title: "Expected a string for the name",
				}))
			} else {
				name = n
				log.Println("Got name:", name)
				break
			}
		}

		attempts++
	}

	// Just something to ensure that there are not thread safety issues
	writer := ws.NewThreadSafeWriter(c)

	// Insert the participant into the room
	rooms.InsertParticipant(
		roomId,
		clientId,
		struct {
			WebSocketWriter ws.ThreadSafeWriter
			Name            string
		}{WebSocketWriter: writer, Name: name},
	)

	log.Println("Participant left the room")
	defer rooms.RemoveParticipant(roomId, clientId)

	for {
		log.Println("Waiting for next message")
		event, ok := <-messageChannel
		if !ok {
			break
		}
		var message clientmessages.Message
		err := json.Unmarshal(event, &message)
		if err != nil {
			// TODO: this seems like a serious mistake. Please investigate this
			log.Println("Error parsing message: ", err.Error())
			continue
		}

		m, err := clientmessages.ParseMessage(message)
		if err != nil {
			b, err := json.Marshal(servermessages.CreateClientError(servermessages.ErrorResponse{
				Title: "Failed to parse message",
			}))

			if err != nil {
				log.Println("Was not able to marshal error message to be sent to client")

				// TODO: this is a serious error. Please investigate this
				writer.Write([]byte("Bad message body, and also failed to send error message. This is a serious server error"))
				return
			} else {
				writer.Write(b)
			}
		}

		switch v := m.(type) {
		case clientmessages.MessageToParticipant:
			err := rooms.SendMessageToParticipant(roomId, clientId, v)
			if err != nil {
				log.Println("Error sending message to participant: ", err.Error())
				b, err := json.Marshal(servermessages.CreateClientError(servermessages.ErrorResponse{
					Title: "Failed to send message to participant",
				}))
				if err != nil {
					log.Println("Was not able to marshal error message to be sent to client")
				}
				writer.Write(b)
			}
		}
	}
}

func handleLeaveRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	roomId, ok := params["roomId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	participantId, ok := params["participantId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rooms.RemoveParticipant(roomId, participantId)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/room/{id}", handleRoom)
	r.HandleFunc("/leave-room/{roomId}/{participantId}", handleLeaveRoom).Methods(http.MethodPost)

	port := config.GetPort()
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
