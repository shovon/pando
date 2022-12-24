package set

import "tree/iterable"

// Set is for representing a set of objects, irrespective insertion order.
// additionally, duplicate insertion of the same key into a Set will result in
// subsequent `Add` invocations to effectively be a no-op
type Set map[interface{}]bool

var _ iterable.Iterable = Set{}

func New(items ...interface{}) Set {
	result := Set{}

	for _, item := range items {
		result.Add(item)
	}

	return result
}

// FromSlice creates a new Set from an existing slice
func FromSlice(s []interface{}) Set {
	newSet := Set{}
	for _, v := range s {
		newSet.Add(v)
	}

	return newSet
}

// Add adds a key to the set
func (s Set) Add(value interface{}) {
	s[value] = true
}

// Has checks for existence of a key in the set
func (s Set) Has(value interface{}) bool {
	return s[value]
}

// Iterate creates a channel purely for iteration purposes
func (s Set) Iterate() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for k := range s {
			c <- k
		}
		close(c)
	}()

	return c
}

// Equal checks two set equality.
//
// What's nice about sets is that insertion order does not matter; only the
// cardinality of the set, and also wither all existing keys match each other's
// keys will be checked against
func (s Set) Equals(s1 Set) bool {
	if len(s) != len(s1) {
		return false
	}

	for k := range s {
		if !s1.Has(k) {
			return false
		}
	}

	return true
}

func (s Set) IsSubsetTo(s1 Set) bool {
	if len(s) > len(s1) {
		return false
	}

	for k := range s {
		if !s1.Has(k) {
			return false
		}
	}

	return true
}

// Union performs a union of two sets. You can think of this method as if it
// were to concatenate two sets together
func (s Set) Union(s1 Set) Set {
	result := Set{}
	for k, ok := range s {
		if ok {
			result.Add(k)
		}
	}
	for k, ok := range s1 {
		if ok {
			result.Add(k)
		}
	}
	return result
}
