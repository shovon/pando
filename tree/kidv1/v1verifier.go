package kidv1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	keyiderrors "tree/keyid/errors"
	"tree/keyid/verifier"
)

type V1Verifier struct {
	Key string
}

var _ verifier.Verifier = &V1Verifier{}

func toBigInt(b []byte) (i *big.Int) {
	i = &big.Int{}
	i.SetBytes(b)
	return i
}

func (v *V1Verifier) Verify(message, signature []byte) (bool, error) {
	value, err := base64.RawStdEncoding.DecodeString(v.Key)
	if err != nil {
		return false, err
	}
	header := value[0]
	if header != 0x04 {
		return false, keyiderrors.ErrBadKeyFormat
	}
	if len(signature) != 64 {
		return false, nil
	}
	x := toBigInt(value[1:33])
	y := toBigInt(value[33:])

	r := toBigInt(signature[:32])
	s := toBigInt(signature[32:])

	pubKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	messageHash := sha256.Sum256(message)

	result := ecdsa.Verify(pubKey, messageHash[:], r, s)
	return result, nil
}

func (v *V1Verifier) IsKeyValid() bool {
	value, err := base64.RawStdEncoding.DecodeString(v.Key)
	if err != nil {
		return false
	}
	if len(value) != 65 {
		return false
	}
	return true
}
