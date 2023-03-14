package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func setWriteDeadline(c *websocket.Conn) error {
	return c.SetWriteDeadline(time.Now().Add(writeWait))
}

// ReadLoop gets a channel that will receive messages from the websocket
func ReadLoop(c *websocket.Conn) <-chan []byte {
	ch := make(chan []byte)

	once := &sync.Once{}

	go func() {
		for {
			t, b, err := c.ReadMessage()
			c.SetReadDeadline(time.Now().Add(pongWait))

			if err != nil {
				once.Do(func() {
					close(ch)
				})
				return
			}
			if t == websocket.TextMessage || t == websocket.BinaryMessage {
				ch <- b
			}
		}
	}()

	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		for {
			<-time.After(pingPeriod)
			setWriteDeadline(c)
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				once.Do(func() {
					close(ch)
				})
				return
			}
		}
	}()

	return ch
}

func writeTextMessage(c *websocket.Conn, m []byte) error {

	err := setWriteDeadline(c)
	if err != nil {
		return err
	}

	return c.WriteMessage(websocket.TextMessage, m)
}

func writeJSONMessage(c *websocket.Conn, m any) error {
	err := setWriteDeadline(c)
	if err != nil {
		return err
	}

	return c.WriteJSON(m)
}
