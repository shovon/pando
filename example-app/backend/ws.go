package main

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 60 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func wrapper() {

}

func setWriteDeadline(c *websocket.Conn) error {
	return c.SetWriteDeadline(time.Now().Add(writeWait))
}

func readLoop(c *websocket.Conn) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		for {
			t, b, err := c.ReadMessage()
			c.SetReadDeadline(time.Now().Add(pongWait))

			if err != nil {
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
