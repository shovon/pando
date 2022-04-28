package key

type BadVerifier struct {
}

var _ Verifier = &BadVerifier{}

func (v *BadVerifier) Verify(message, signature []byte) (bool, error) {
	return false, nil
}

func (v *BadVerifier) IsKeyValid() bool {
	return false
}
