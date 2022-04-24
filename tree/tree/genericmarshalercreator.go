package tree

import "encoding/json"

type genericMarshalerCreator func(interface{}) json.Marshaler
