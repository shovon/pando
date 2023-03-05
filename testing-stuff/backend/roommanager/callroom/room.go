package callroom

import (
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"backend/pairmap"
	"backend/slice"
	"backend/sortedmap"
	"fmt"
	"log"
	"sync"
)

// Room is a room in the call
//
// Please don't initialize this struct directly, use NewRoom instead
type Room struct {
	lock    *sync.RWMutex
	clients sortedmap.SortedMap[string, Client]
}

// NewRoom creates a new Room instance
func NewRoom() Room {
	return Room{lock: &sync.RWMutex{}, clients: sortedmap.New[string, Client]()}
}

// InsertClient inserts a new client into the room
func (r *Room) InsertClient(participantId string, participant Client) {
	r.lock.Lock()
	defer r.lock.Unlock()

	fmt.Println("Inserting client", participantId)

	r.clients.Set(participantId, participant)

	r.signalRoomState()
}

// RemoveClient removes a client from the room
func (r *Room) RemoveClient(participantId string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.clients.Delete(participantId)
	r.signalRoomState()
}

func (r *Room) DisconnectClient(participantId string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	client, ok := r.clients.Get(participantId)

	if !ok {
		return
	}

	client.WebSocketWriter.Close()
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

	participant, ok := r.clients.Get(message.To)
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
	return r.clients.Len()
}

type DetailedParticipantState struct {
	ParticipantState

	ConnectionStatus string `json:"connectionStatus"`
}

func (r Room) getParticipantsState() []DetailedParticipantState {
	r.lock.RLock()
	defer r.lock.RUnlock()

	participants := make([]DetailedParticipantState, r.clients.Len())
	for i, participant := range r.clients.Values() {
		participants[i] = DetailedParticipantState{
			ParticipantState: participant.Participant,
			ConnectionStatus: participant.ConnectionStatus(),
		}
	}

	return participants
}

// GetRoomState returns the room state, which includes all participants and
// their current state
func (r Room) GetRoomState() servermessages.RoomState {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.getRoomState()
}

func (r Room) getRoomState() servermessages.RoomState {
	p := slice.Map(
		r.clients.Pairs(),
		func(kv sortedmap.KV[string, Client]) pairmap.KV[
			string,
			any,
		] {
			// We're gonna need so much more as well
			return pairmap.KV[string, any]{
				Key:   kv.Key,
				Value: kv.Value.Participant,
			}
		},
	)

	return servermessages.RoomState{
		Participants: p,
	}
}

func (r Room) signalRoomState() {
	for _, participant := range r.clients.Values() {
		// This is so innefficient, but it needs to be done, for now
		fmt.Println("Sending room state to", participant.Participant.Name)
		err := participant.WebSocketWriter.WriteJSON(
			servermessages.CreateRoomStateMessage(r.getRoomState()),
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
