package roommanager

import (
	"backend/messages/clientmessages"
	"backend/roommanager/callroom"
	"sync"
)

type RoomManager struct {
	lock  *sync.RWMutex
	rooms map[string]callroom.Room
}

func NewRoomManager() RoomManager {
	return RoomManager{lock: &sync.RWMutex{}, rooms: make(map[string]callroom.Room)}
}

func (r *RoomManager) SendMessageToParticipant(roomId, fromParticipantId string, message clientmessages.MessageToParticipant) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	room, ok := r.rooms[roomId]
	if !ok {
		return
	}

	room.SendMessageToParticipant(message, fromParticipantId)
}

func (r *RoomManager) InsertParticipant(
	roomId, participantId string,
	participant callroom.Participant,
) {
	r.lock.Lock()
	defer r.lock.Unlock()

	room := r.getRoom(roomId)
	room.InsertParticipant(participantId, participant)
}

func (r *RoomManager) RemoveParticipant(roomId, participantId string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	room, ok := r.rooms[roomId]
	if !ok {
		return
	}

	room.RemoveParticipant(participantId)

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
