package main

import (
	"errors"
	"net/http"
	"tree/tree"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{}
var trees = make(map[string]*tree.Tree)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tree/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		tree := &tree.Tree{}
		trees[id] = tree
		defer delete(trees, id)
	})
}
