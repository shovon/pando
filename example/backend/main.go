package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/shovon/gorillawswrapper"
	"github.com/sparkscience/wskeyid-go/v2"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func getPort() int {
	port := strings.Trim(os.Getenv("PORT"), " ")
	if port == "" {
		return 8080
	}

	num, err := strconv.Atoi(port)
	if err != nil {
		return 8080
	}

	return num
}

func handleCall(w http.ResponseWriter, r *http.Request) {

	log.Print("Got connection from client")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer c.Close()

	log.Println("Got connection object")

	conn := gorillawswrapper.NewWrapper(c)

	{
		err := wskeyid.HandleAuthConnection(r, conn)
		if err != nil {
			log.Println("WebSocket authentication failed ", err.Error())
			return
		}
	}

	log.Println("Got connection wrapper")

	defer conn.Stop()
	defer log.Println("Closing connection")

	conn.WriteTextMessage("Cool")

	for msg := range conn.MessagesChannel() {
		fmt.Println(string(msg.Message))
		if err := conn.WriteTextMessage("Got message"); err != nil {
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/room/{id}", handleCall)

	port := getPort()
	log.Printf("Listening on port %d", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
