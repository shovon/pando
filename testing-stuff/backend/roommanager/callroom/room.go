package callroom

import (
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Represents a single participant, not as far as the problem domain, but as a
// client in the call.
type Client struct {
	// The connection associated with the participant
	Connection *websocket.Conn

	// Participant is the metadata associated with the participant
	Participant servermessages.ParticipantState
}

// Room is a room in the call
//
// Please don't initialize this struct directly, use NewRoom instead
type Room struct {
	lock    *sync.RWMutex
	clients map[string]Client
}

// NewRoom creates a new Room instance
func NewRoom() Room {
	return Room{lock: &sync.RWMutex{}, clients: make(map[string]Client)}
}

// InsertClient inserts a new client into the room
func (r *Room) InsertClient(participantId string, participant Client) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.clients[participantId] = participant

	r.signalRoomState()
}

func (r Room) signalRoomState() {
	for _, participant := range r.clients {
		// This is so innefficient, but it needs to be done, for now
		err := participant.Connection.WriteJSON(
			servermessages.CreateRoomStateMessage(r.GetRoomState()),
		)
		if err != nil {
			// TODO: figure out a more robust solution, for the event when something
			//   goes wrong.
			//
			//   Check to see what causes the error. If it's a connection error, then,
			//   just close the connection
			log.Print("Error", err.Error())
		}
	}
}

// RemoveClient removes a client from the room
func (r *Room) RemoveClient(participantId string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.clients, participantId)

	r.signalRoomState()
}

// SendMessageToClient sends a message to the client
func (r Room) SendMessageToClient(
	message clientmessages.MessageToParticipant,
	fromParticipantId string,
) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	participant, ok := r.clients[message.To]
	if !ok {
		return
	}

	participant.Connection.WriteJSON(
		servermessages.CreateMessageToParticipant(fromParticipantId, message.Data),
	)
}

// Size returns the number of clients in the room
func (r Room) Size() int {
	return len(r.clients)
}

func toParticipantStateMap(
	participant map[string]Client,
) {
	m := make(map[string]servermessages.ParticipantState)

	for k, v := range participant {
		m[k] = v.Participant
	}

	return m
}

// GetRoomState returns the room state, which includes all participants and
// their current state
func (r Room) GetRoomState() servermessages.RoomState {
	return servermessages.RoomState{
		Participants: toParticipantStateMap(r.clients),
	}
}
