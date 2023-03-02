package config

import (
	"encoding/base64"
	"log"
	"os"
	"strings"
)

var jwtHS256Key []byte

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
}

// GetJWTKey returns a copy of the JWT HS256 key
//
// Note: due to data safety and integrity, this function returns a copy of the
// key. This is to prevent the key from being modified by the caller.
//
// For performance reasons, the caller should cache the key returned by this
// function.
func GetJWTKey() []byte {
	k := make([]byte, len(jwtHS256Key))
	copy(k, jwtHS256Key)
	return k
}
