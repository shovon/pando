package key

import (
	"encoding/base64"
	"errors"
	"strings"
)

var BadKeyFormatError = errors.New("The supplied key format is bad")

func Parse(key string) error {
	parts := strings.Split(key, ".")
	version, remainder := parts[0], parts[1:]

	switch version {
	case "v1":
		return parseV1(strings.Join(remainder, "."))
	default:
		return errors.New("not a valid key")
	}
}

// V1 is a base64-encoded 65-byte ES256 public key, with three values
// concatenated:
//
// - literally the number 0x04, as a single byte
// - a 32-byte x coordinate
// - a 32-byte y coordinate
func parseV1(key string) error {
	v, err := base64.RawStdEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	header := v[0]
	if header != 0x04 {
		return BadKeyFormatError
	}
	x := v[1:32]
	y := v[32:]
	return errors.New("not yet implemented")
}
