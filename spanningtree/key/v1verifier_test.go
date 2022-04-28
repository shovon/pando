package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestExample(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.FailNow()
	}

	message := []byte("Hello, World!")

	hash := sha256.Sum256(message)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.FailNow()
	}

	publicKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: privateKey.X, Y: privateKey.Y}

	ecdsa.Verify(publicKey, hash[:], r, s)
}

func TestVerify(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.FailNow()
	}

	message := []byte("Hello, World!")

	hash := sha256.Sum256(message)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.FailNow()
	}

	signature := append(r.Bytes(), s.Bytes()...)

	key := append([]byte{0x04}, privateKey.X.Bytes()...)
	key = append(key, privateKey.Y.Bytes()...)
	keyStr := base64.RawStdEncoding.EncodeToString(key)

	verifier := V1Verifier{keyStr}

	ok, err := verifier.Verify(message, signature)

	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Error("Failed to verify signature")
	}
}
