package node

type AdjacencyNode struct {
	KeyValue Pair          `json:"keyValue"`
	Links    []interface{} `json:"links"`
}
