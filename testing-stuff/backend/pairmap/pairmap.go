package pairmap

import (
	"backend/keyvalue"
	"encoding/json"
)

// PairMap is a map that can be marshalled into JSON as a tuple list of
// key-value pairs.
type PairMap[K comparable, V any] map[K]V

var _ json.Marshaler = PairMap[int, any]{}

func (p PairMap[K, V]) MarshalJSON() ([]byte, error) {
	m := []keyvalue.KV[K, V]{}

	for key, value := range p {
		tuple := keyvalue.KV[K, V]{Key: key, Value: value}
		m = append(m, tuple)
	}

	return json.Marshal(m)
}
