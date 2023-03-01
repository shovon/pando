package roommanager

import (
	"backend/messages/clientmessages"
	"backend/messages/servermessages"
	"backend/roommanager/callroom"
	"backend/ws"
	"sync"
)

type RoomManager struct {
	lock  *sync.RWMutex
	rooms map[string]callroom.Room
}

func NewRoomManager() RoomManager {
	return RoomManager{lock: &sync.RWMutex{}, rooms: make(map[string]callroom.Room)}
}

// SendMessageToRoom is for handling an event where a participant intends to
// send a message to another participant in the room
func (r *RoomManager) SendMessageToParticipant(
	roomId, fromParticipantId string,
	message clientmessages.MessageToParticipant,
) (bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	room, ok := r.rooms[roomId]
	if !ok {
		// TODO: be specific about the error. Notify the client code that the reason
		//   why sending failed is because the room doesn't exist
		return false, nil
	}

	return room.SendMessageToClient(message, fromParticipantId)
}

// InsertParticipant inserts a new participant into the room
func (r *RoomManager) InsertParticipant(
	roomId, participantId string,
	participant struct {
		WebSocketWriter ws.ThreadSafeWriter
		Name            string
	},
) {
	r.lock.Lock()
	defer r.lock.Unlock()

	room := r.getRoom(roomId)
	room.InsertClient(
		participantId,
		callroom.Client{
			WebSocketWriter: participant.WebSocketWriter,
			Participant:     servermessages.ParticipantState{Name: participant.Name},
		},
	)
}

// RemoveParticipant removes a participant from the room
func (r *RoomManager) RemoveParticipant(roomId, participantId string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	room, ok := r.rooms[roomId]
	if !ok {
		return
	}

	room.RemoveClient(participantId)

	if room.Size() < 0 {
		delete(r.rooms, roomId)
	}
}

// getRoom either gets or creates a room.
//
// NOT THREAD SAFE!
func (r *RoomManager) getRoom(roomId string) callroom.Room {
	room, ok := r.rooms[roomId]
	if !ok {
		room = callroom.NewRoom()
		r.rooms[roomId] = room
	}

	return room
}
