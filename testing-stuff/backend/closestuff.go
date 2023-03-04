package main

import (
	"backend/config"
	"backend/nextint"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
	"unsafe"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWT(clientId, roomId string) (string, error) {
	i := nextint.NextInt()
	iat := time.Now().Unix()

	key := config.GetCurrentProcessKey()

	b := make(
		[]byte,
		int(unsafe.Sizeof(i))+int(unsafe.Sizeof(iat))+len(key),
	)

	binary.LittleEndian.PutUint64(b, uint64(i))
	binary.LittleEndian.PutUint64(b[8:], uint64(iat))
	copy(b[16:], key)

	hash := sha256.Sum256(b)

	hashString := string(hash[:])

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// TODO: consider softcoding the iss field
		"iss":      "demo-room-backend",
		"clientId": clientId,
		"roomId":   roomId,
		"jti":      hashString,
		"iat":      iat,
	})

	return token.SignedString(config.GetHS256Key())
}

func verifyJWT(clientId, roomId, j string) (bool, error) {
	token, err := jwt.Parse(j, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return config.GetHS256Key(), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["clientId"] != clientId {
			return false, nil
		}

		if claims["roomId"] != roomId {
			return false, nil
		}

		return true, nil
	} else {
		return false, nil
	}
}
