package main

import (
	"errors"
	"net/http"
	"spanningtree/spanningtree"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{}

var trees = make(map[string]*spanningtree.Tree)

func getTree(id string) *spanningtree.Tree {
	t, ok := trees[id]
	if !ok {
		t = &spanningtree.Tree{}
		trees[id] = t
	}

	return t
}

func deleteValue(key string) {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tree/{id}/{userid}", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Err(err)
			w.WriteHeader(500)
			w.Write([]byte("An internal server error occurred"))
			return
		}

		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Err(errors.New("Internal server error"))
			w.WriteHeader(500)
			w.Write([]byte("Internal server error. Failed to parse ID"))
		}

		t := getTree(id)
		listener := t.RegisterChangeListener()

		go func() {
			switch (<-listener).(type) {
			case spanningtree.Deleted:

			}
		}()
	})
}
