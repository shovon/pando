package maybe

// Maybe is a type that can either have a value or not
//
// Go seems to not fucking have this concept of an algebraic data type.
//
// So we have to deal with this shit.
type Maybe[T any] struct {
	hasValue bool
	value    T
}

// Get returns the value and a boolean indicating if the value is present
func (m Maybe[T]) Get() (T, bool) {
	return m.value, m.hasValue
}

// Something returns a Maybe with a value
func Something[T any](value T) Maybe[T] {
	return Maybe[T]{hasValue: true, value: value}
}

// Nothing returns a Maybe with no value
func Nothing[T any]() Maybe[T] {
	return Maybe[T]{hasValue: false}
}
