package callroom

import (
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"log"
	"sync"

	"github.com/shovon/gorillawswrapper"
)

type Participant struct {
	Connection gorillawswrapper.Wrapper
}

type Room struct {
	lock         *sync.RWMutex
	participants map[string]Participant
}

func NewRoom() Room {
	return Room{lock: &sync.RWMutex{}, participants: make(map[string]Participant)}
}

func (r *Room) InsertParticipant(participantId string, participant Participant) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.participants[participantId] = participant

	r.signalRoomState()
}

func (r Room) signalRoomState() {
	for _, participant := range r.participants {
		// This is so innefficient, but it needs to be done, for now
		err := participant.Connection.WriteJSON(
			servermessages.CreateRoomStateMessage(r.GetRoomState()),
		)
		if err != nil {
			log.Print("Error", err.Error())
		}
	}
}

func (r *Room) RemoveParticipant(participantId string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.participants, participantId)

	r.signalRoomState()
}

func (r Room) SendMessageToParticipant(message clientmessages.MessageToParticipant, fromParticipantId string) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	participant, ok := r.participants[message.To]
	if !ok {
		return
	}

	participant.Connection.WriteJSON(
		servermessages.CreateMessageToParticipant(fromParticipantId, message.Data),
	)
}

func (r Room) Size() int {
	return len(r.participants)
}

func toParticipantStateMap(
	participant map[string]Participant,
) map[string]servermessages.ParticipantState {
	m := make(map[string]servermessages.ParticipantState)
	for key := range participant {
		m[key] = servermessages.ParticipantState{}
	}

	return m
}

func (r Room) GetRoomState() servermessages.RoomState {
	return servermessages.RoomState{
		Participants: toParticipantStateMap(r.participants),
	}
}
