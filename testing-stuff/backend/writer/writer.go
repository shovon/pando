package writer

type Writer interface {
	Write([]byte) error
	WriteJSON(interface{}) error
}
