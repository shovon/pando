package main

import (
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"backend/roommanager"
	"backend/roommanager/callroom"
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

var rooms = roommanager.NewRoomManager()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func getPort() int {
	port := strings.Trim(os.Getenv("PORT"), " ")
	if port == "" {
		return 8080
	}

	num, err := strconv.Atoi(port)
	if err != nil {
		return 8080
	}

	return num
}

func handleCall(w http.ResponseWriter, r *http.Request) {

	log.Print("Got connection from client")

	params := mux.Vars(r)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer c.Close()

	log.Println("Got connection object")

	{
		err := wskeyid.HandleAuthConnection(r, c)
		if err != nil {
			log.Println("WebSocket authentication failed ", err.Error())
			return
		}
	}

	log.Println("Got connection wrapper")

	clientId := strings.TrimSpace(r.URL.Query().Get("client_id"))

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

	rooms.InsertParticipant(
		roomId,
		clientId,
		callroom.Participant{Connection: c},
	)

	messageChannel := readLoop(c)

	for event := range messageChannel {
		var message clientmessages.Message
		err := json.Unmarshal(event, &message)
		if err != nil {
			continue
		}

		switch message.Type {
		case "MESSAGE_TO_PARTICIPANT":
			toParticipant, err := clientmessages.ParseMessageToParticipant(message.Data)
			if err == nil {
				rooms.SendMessageToParticipant(roomId, clientId, toParticipant)
			} else {
				// TODO: send a better message to participant
				b, err := json.Marshal(servermessages.CreateClientError(servermessages.ErrorResponse{
					Title: "Bad message",
				}))
				if err != nil {
					panic(err)
				}
				writeTextMessage(c, b)
			}
		}
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/room/{id}", handleCall)

	port := getPort()
	log.Printf("Listening on port %d", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
