package callroom

import (
	"backend/maputils"
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"backend/pairmap"
	"log"
	"sync"
)

// Room is a room in the call
//
// Please don't initialize this struct directly, use NewRoom instead
type Room struct {
	lock    *sync.RWMutex
	clients pairmap.PairMap[string, Client]
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

// RemoveClient removes a client from the room
func (r *Room) RemoveClient(participantId string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.clients, participantId)

	r.signalRoomState()
}

// SendMessageToClient is intended to handle the event when a participant
// intends to send a direct message to another participant.
//
// A boolean is returned to indicate whether the message was sent successfully
func (r Room) SendMessageToClient(
	message clientmessages.MessageToParticipant,
	fromParticipantId string,
) (bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	participant, ok := r.clients[message.To]
	if !ok {
		return false, nil
	}

	err := participant.WebSocketWriter.WriteJSON(
		servermessages.CreateMessageToParticipant(fromParticipantId, message.Data),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Size returns the number of clients in the room
func (r Room) Size() int {
	return len(r.clients)
}

// GetRoomState returns the room state, which includes all participants and
// their current state
func (r Room) GetRoomState() servermessages.RoomState {
	return servermessages.RoomState{
		Participants: maputils.Map(
			r.clients,
			func(key string, c Client) (string, servermessages.ParticipantState) {
				return key, c.Participant
			},
		),
	}
}

func (r Room) signalRoomState() {
	for _, participant := range r.clients {
		// This is so innefficient, but it needs to be done, for now
		err := participant.WebSocketWriter.WriteJSON(
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
