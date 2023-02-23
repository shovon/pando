package main

import (
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"backend/roommanager"
	"backend/roommanager/callroom"
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

const defaultPort = 8080

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

	// Insert the participant into the room
	rooms.InsertParticipant(
		roomId,
		clientId,
		callroom.Client{Connection: c},
	)

	defer rooms.RemoveParticipant(roomId, clientId)

	// Creates a message channel, to read from
	messageChannel := ws.ReadLoop(c)

	// Just something to ensure that there are not thread safety issues.
	//
	// TODO: ensure there is absolutely no way to do anything that is
	//   thread-unsafe
	writer := ws.NewThreadSafeWriter(c)

	for event := range messageChannel {
		var message clientmessages.Message
		err := json.Unmarshal(event, &message)
		if err != nil {
			continue
		}

		switch message.Type {
		case "MESSAGE_TO_PARTICIPANT":
			// Recieved a message from participant, inteded for another participant

			toParticipant, err := clientmessages.ParseMessageToParticipant(message.Data)
			if err == nil {
				rooms.SendMessageToParticipant(roomId, clientId, toParticipant)
			} else {
				// TODO: send a better message to participant
				b, err := json.Marshal(servermessages.CreateClientError(servermessages.ErrorResponse{
					Title: "Bad message",
				}))
				// TODO: handle the event when an error occurred attempting to marshal
				//   JSON
				if err != nil {
					log.Println("Was not able to marshal error message to be sent to client")
					return
				}
				writer.Write(b)
			}
		case "BROADCAST_MESSAGE":
			b, err := json.Marshal(servermessages.CreateServerError(servermessages.ErrorResponse{
				Title: "Not yet implemented",
			}))
			// TODO: handle the event when an error occurred attempting to marshal
			//   JSON
			if err != nil {
				log.Println("Was not able to marshal error message to be sent to client")
				return
			}
			writer.Write(b)
		case "ENABLE_VIDEO":
			b, err := json.Marshal(servermessages.CreateServerError(servermessages.ErrorResponse{
				Title: "Not yet implemented",
			}))
			// TODO: handle the event when an error occurred attempting to marshal
			//   JSON
			if err != nil {
				log.Println("Was not able to marshal error message to be sent to client")
				return
			}
			writer.Write(b)
		case "DISABLE_VIDEO":
			b, err := json.Marshal(servermessages.CreateServerError(servermessages.ErrorResponse{
				Title: "Not yet implemented",
			}))
			// TODO: handle the event when an error occurred attempting to marshal
			//   JSON
			if err != nil {
				log.Println("Was not able to marshal error message to be sent to client")
				return
			}
			writer.Write(b)
		case "ENABLE_AUDIO":
			b, err := json.Marshal(servermessages.CreateServerError(servermessages.ErrorResponse{
				Title: "Not yet implemented",
			}))
			// TODO: handle the event when an error occurred attempting to marshal
			//   JSON
			if err != nil {
				log.Println("Was not able to marshal error message to be sent to client")
				return
			}
			writer.Write(b)
		case "DISABLE_AUDIO":
			b, err := json.Marshal(servermessages.CreateServerError(servermessages.ErrorResponse{
				Title: "Not yet implemented",
			}))
			// TODO: handle the event when an error occurred attempting to marshal
			//   JSON
			if err != nil {
				log.Println("Was not able to marshal error message to be sent to client")
				return
			}
			writer.Write(b)
		}

	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/room/{id}", handleRoom)

	port := getPort()
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
