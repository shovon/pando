package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"
)

var jwtHS256Key []byte
var currentProcessKey []byte = make([]byte, 16)

func init() {
	key := strings.Trim(os.Getenv("JWT_HS256_KEY"), " ")
	if key == "" {
		log.Fatal("JWT_HS256_KEY environment variable is not set")
	}

	var err error
	jwtHS256Key, err = base64.StdEncoding.DecodeString(key)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := rand.Read(currentProcessKey); err != nil {
		log.Fatal(err)
	}
}

// GetHS256Key returns a copy of the JWT HS256 key
//
// Note: due to data safety and integrity, this function returns a copy of the
// key. This is to prevent the key from being modified by the caller.
//
// For performance reasons, the caller should cache the key returned by this
// function.
func GetHS256Key() []byte {
	k := make([]byte, len(jwtHS256Key))
	copy(k, jwtHS256Key)
	return k
}

// GetCurrentProcessKey returns a copy of the current process key.
//
// The current process key is just a string of bytes that has been initialized
// at random, on startup.
//
// Use it for whatever you want
func GetCurrentProcessKey() []byte {
	k := make([]byte, len(currentProcessKey))
	copy(k, currentProcessKey)
	return k
}

const defaultPort = 3333

// GetPort gets the appropriate port to get the server running on.
// If the PORT environment variable is set, then it will use that. Otherwise, it
// will use the default port
func GetPort() int {
	port := strings.Trim(os.Getenv("PORT"), " ")
	if port == "" {
		return defaultPort
	}

	num, err := strconv.Atoi(port)
	if err != nil {
		return defaultPort
	}

	return num
}
