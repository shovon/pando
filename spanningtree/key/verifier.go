package key

type Verifier interface {
	Verify([]byte, []byte) (bool, error)
	IsKeyValid() bool
}
