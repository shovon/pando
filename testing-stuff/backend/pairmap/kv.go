package pairmap

import "encoding/json"

// KV is a key-value pair that can be marshalled into JSON as an array tuple
type KV[K comparable, V any] struct {
	Key   K
	Value V
}

var _ json.Marshaler = KV[int, any]{}

func (k KV[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{k.Key, k.Value})
}
