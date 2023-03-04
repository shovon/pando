package main

import (
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"backend/roommanager"
	"backend/ws"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sparkscience/wskeyid-golang"
)

const defaultPort = 3333

var rooms = roommanager.NewRoomManager()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Gets the appropriate port to get the server running on.
// If the PORT environment variable is set, then it will use that. Otherwise, it
// will use the default port
func getPort() int {
	port := strings.Trim(os.Getenv("PORT"), " ")
	if port == "" {
		return defaultPort
	}

	num, err := strconv.Atoi(port)
	if err != nil {
		return defaultPort
	}

	return num
}

// handleRoom is the event handler for the room endpoint.
//
// This is where when there is a connection to a room, the participant will
// interact with the room
func handleRoom(w http.ResponseWriter, r *http.Request) {
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

	for {
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
			} else {
				name = n
				log.Println("Got name:", name)
				break
			}

		}
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
			sent, err := rooms.SendMessageToParticipant(roomId, clientId, v)
			if !sent {
				log.Println("Attempted to send message to a participant that does not exist")
			}
			// TODO: this seems like a serious bug. Please investigate this
			if err != nil {
				log.Println("Error sending message to participant: ", err.Error())
			}
		}
	}

	log.Println("Participant left the room")
	rooms.RemoveParticipant(roomId, clientId)
}

func handleLeaveRoom(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/room/{id}", handleRoom)

	// TODO: ensure that this works exclusively with POST requests
	r.HandleFunc("/leave-room/{roomId}/{participantId}", handleLeaveRoom)

	port := getPort()
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
