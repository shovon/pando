package sortedmap

import (
	"backend/pairmap"
	"backend/slice"
	"encoding/json"
	"errors"
	"sort"
)

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

type IntValue[V any] struct {
	value V
	order int
}

type SortedMap[K comparable, V any] struct {
	m           map[K]IntValue[V]
	latestOrder int
}

var _ json.Marshaler = SortedMap[int, any]{}
var _ json.Unmarshaler = &SortedMap[int, any]{}

func (s SortedMap[K, V]) MarshalJSON() ([]byte, error) {
	p := slice.Map(
		s.Pairs(),
		func(kv KV[K, V]) pairmap.KV[
			K,
			V,
		] {
			// We're gonna need so much more as well
			return pairmap.KV[K, V]{
				Key:   kv.Key,
				Value: kv.Value,
			}
		},
	)

	return json.Marshal(p)
}

func (s *SortedMap[K, V]) UnmarshalJSON(data []byte) error {
	return errors.New("not yet implemented")
}

func New[K comparable, V any]() SortedMap[K, V] {
	return SortedMap[K, V]{
		m: map[K]IntValue[V]{},
	}
}

func (s *SortedMap[K, V]) Set(key K, value V) {
	s.latestOrder++
	s.m[key] = IntValue[V]{
		value: value,
		order: s.latestOrder,
	}
}

func (s SortedMap[K, V]) Get(key K) (V, bool) {
	value, ok := s.m[key]
	return value.value, ok
}

func (s *SortedMap[K, V]) Delete(key K) {
	delete(s.m, key)
}

func (s SortedMap[K, V]) Len() int {
	return len(s.m)
}

func (s SortedMap[K, V]) Keys() []K {
	return slice.Map(s.Pairs(), func(kv KV[K, V]) K {
		return kv.Key
	})
}

func (s SortedMap[K, V]) Values() []V {
	return slice.Map(s.Pairs(), func(kv KV[K, V]) V {
		return kv.Value
	})
}

func (s SortedMap[K, V]) Pairs() []KV[K, V] {
	pairs := make([]KV[K, V], 0, len(s.m))
	for key, value := range s.m {
		pairs = append(pairs, KV[K, V]{
			Key:   key,
			Value: value.value,
		})
	}
	sort.SliceStable(pairs, func(i int, j int) bool {
		return s.m[pairs[i].Key].order < s.m[pairs[j].Key].order
	})
	return pairs
}

func (s SortedMap[K, V]) Has(key K) bool {
	_, ok := s.m[key]
	return ok
}

func (s *SortedMap[K, V]) Clear() {
	s.m = map[K]IntValue[V]{}
}

func (s SortedMap[K, V]) Copy() SortedMap[K, V] {
	return SortedMap[K, V]{
		m:           s.m,
		latestOrder: s.latestOrder,
	}
}

func (s SortedMap[K, V]) Clone() SortedMap[K, V] {
	return s.Copy()
}

func (s SortedMap[K, V]) Merge(other SortedMap[K, V]) SortedMap[K, V] {
	merged := s.Copy()
	for key, value := range other.m {
		merged.m[key] = value
	}
	return merged
}
