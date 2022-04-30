package keyid

type Verifier interface {
	Verify([]byte, []byte) (bool, error)
	IsKeyValid() bool
}
