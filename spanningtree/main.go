package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"spanningtree/spanningtree"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{}

var trees = newTreeManager()

type participant struct {
	conn *websocket.Conn
	meta json.RawMessage
}

var _ json.Marshaler = &participant{}

func (p *participant) MarshalJSON() ([]byte, error) {
	return p.meta, nil
}

func parseKey(id string) error {
	bytes := make([]byte)
	parsed, err := base64.StdEncoding.Decode(bytes, []byte(id))



	return err
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tree/{id}/{userid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		userId, ok := vars["userid"]
		if !ok {
			log.Err(errors.New("The user ID was not set, for some reason. This is bad"))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse User ID"))
		}

		id, ok := vars["id"]
		if !ok {
			log.Err(errors.New("The ID was not set, for some reason. This is bad"))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse ID"))
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Err(err)
			w.WriteHeader(500)
			w.Write([]byte("An internal server error occurred"))
			return
		}

		c.WriteJSON(v interface{})

		t := trees.getTree(id)
		listener := t.RegisterChangeListener(userId)

		go func() {
			switch ev := (<-listener).(type) {
			case spanningtree.NodeState:
				
			}
		}()
	})
}
