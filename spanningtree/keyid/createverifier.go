package key

import (
	"errors"
	"strings"
)

var ErrBadKeyFormat = errors.New("the supplied key format is bad")

const separator = "$"

func CreateVerifier(key string) (Verifier, error) {
	parts := strings.Split(key, separator)
	version, remainder := parts[0], parts[1:]

	switch version {
	case "v1":
		return parseV1(strings.Join(remainder, separator))
	default:
		return &BadVerifier{}, errors.New("not a valid key")
	}
}

// V1 is a base64-encoded 65-byte ES256 public key, with three values
// concatenated:
//
// - literally the number 0x04, as a single byte
// - a 32-byte x coordinate
// - a 32-byte y coordinate
func parseV1(key string) (Verifier, error) {
	return &V1Verifier{key}, nil
}
