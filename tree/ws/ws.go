package ws

import "github.com/gorilla/websocket"

type Wrapper struct {
}

func NewWrapper(c *websocket.Conn) Wrapper {
	return Wrapper{}
}
