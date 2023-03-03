package ws

type Writer interface {
	Write(message []byte) error
	WriteJSON(message interface{}) error
}
