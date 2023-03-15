package callroom

import (
	"backend/connectionstate"
	"backend/keyvalue"
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"backend/slice"
	"backend/sortedmap"
	"errors"
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

	client.Close()
}

func idempotentSend(
	conn connectionstate.Connection,
	message interface{},
) error {
	writer, ok := conn.State().(connectionstate.Connected)
	if ok {
		return writer.WriteJSON(message)
	}

	return nil
}

func createFailedToDeliverMessage(messageID string) servermessages.MessageWithData {
	return servermessages.MessageWithData{
		Type: "FAILED_TO_DELIVER_MESSAGE",
		Data: map[string]interface{}{
			"messageId": messageID,
		},
	}
}

func CreateParticipantDoesNotExist(participantID string) servermessages.MessageWithData {
	return servermessages.MessageWithData{
		Type: "PARTICIPANT_DOES_NOT_EXIST",
		Data: map[string]interface{}{
			"participantId": participantID,
		},
	}
}

// TODO: perhaps return a more detailed message to the original sender as to
// why their message was not sent

// SendMessageToClient is intended to handle the event when a participant
// intends to send a direct message to another participant.
//
// A boolean is returned to indicate whether the message was sent successfully
func (r Room) SendMessageToClient(
	message clientmessages.MessageToParticipant,
	fromParticipantId string,
) error {
	r.lock.RLock()
	defer r.lock.RUnlock()

	participant, participantExists := r.clients.Get(message.To)

	sender, ok := r.clients.Get(fromParticipantId)

	if !ok {
		return errors.New("sender does not exist")
	}

	if !participantExists {
		return idempotentSend(sender.Connection, CreateParticipantDoesNotExist(message.To))
	}

	switch v := participant.Connection.State().(type) {
	case connectionstate.Disconnected:
		if ok {
			return idempotentSend(
				sender.Connection,
				createFailedToDeliverMessage(
					message.ID,
				),
			)
		}
	case connectionstate.Connected:
		// TODO: this is stupid
		return v.WriteJSON(
			servermessages.CreateMessageToParticipant(fromParticipantId, message.Data),
		)

	}

	return errors.New("unknown connection state")
}

// Size returns the number of clients in the room
func (r Room) Size() int {
	return r.clients.Len()
}

type DetailedParticipantState struct {
	ParticipantData

	ConnectionStatus string `json:"connectionStatus"`
}

func (r Room) getParticipantsState() []DetailedParticipantState {
	r.lock.RLock()
	defer r.lock.RUnlock()

	participants := make([]DetailedParticipantState, r.clients.Len())
	for i, participant := range r.clients.Values() {
		participants[i] = DetailedParticipantState{
			ParticipantData:  participant.Participant,
			ConnectionStatus: connectionstate.ConnectionStatus(participant.Connection.State()),
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
		func(kv keyvalue.KV[string, Client]) keyvalue.KV[
			string,
			any,
		] {
			// We're gonna need so much more as well
			return keyvalue.KV[string, any]{
				Key:   kv.Key,
				Value: kv.Value,
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
		err := idempotentSend(participant.Connection,
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
