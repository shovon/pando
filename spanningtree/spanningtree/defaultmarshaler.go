package spanningtree

import "encoding/json"

type defaultMarshaler struct {
	value interface{}
}

var _ json.Marshaler = &defaultMarshaler{}

func (d *defaultMarshaler) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.value)
}
