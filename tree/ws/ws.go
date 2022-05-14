package ws

import (
	"net/http"
	"sync"
	"time"

	"tree/stoppable"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait = 60 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10
)

// A wrapper for Gorilla's WebSocket
type Wrapper struct {
	stoppable.Stoppable
	mut *sync.Mutex
	c   *websocket.Conn
}

func NewWrapper(c *websocket.Conn) Wrapper {
	var mut sync.Mutex

	wrapper := Wrapper{stoppable.NewStoppable(), &mut, c}

	return wrapper
}

func UpgradeWebSocket(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

func (w *Wrapper) Loop() {
loop:
	for {
		select {
		case <-time.After(pingPeriod):
			w.c.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Info().Err(err).Msg("The end host probably closed the connection")
			}
		case <-w.OnStopped():
			break loop
		}
	}
}

func (w *Wrapper) WriteMessage(messageType int, data []byte) error {
	w.mut.Lock()
	defer w.mut.Unlock()
	w.c.SetWriteDeadline(time.Now().Add(writeWait))
	return w.c.WriteMessage(messageType, data)
}

func (w *Wrapper) WriteJSON(v interface{}) error {
	w.mut.Lock()
	defer w.mut.Unlock()
	return w.c.WriteJSON(v)
}
