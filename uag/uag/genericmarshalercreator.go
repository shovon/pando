package uag

import "encoding/json"

type genericMarshalerCreator func(interface{}) json.Marshaler
