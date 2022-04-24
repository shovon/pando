package spanningtree

import "encoding/json"

type genericMarshalerCreator func(interface{}) json.Marshaler
