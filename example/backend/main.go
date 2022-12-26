package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	port := getPort()
	log.Printf("Listening on port %d", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello, World!"))
	})))
}
