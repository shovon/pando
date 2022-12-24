package reactforcegraph

import (
	"encoding/json"
	"errors"
	"tree/adjacencylist"
)

type ReactForceGraphMarshaler adjacencylist.AdjacencyList

var _ json.Marshaler = &ReactForceGraphMarshaler{}

type ForceGraphNode struct {
	ID string `json:"id"`
}

type ForceGraphLink struct {
	Target string `json:"target"`
	Source string `json:"source"`
}

type forceGraph struct {
	Nodes []ForceGraphNode `json:"nodes"`
	Links []ForceGraphLink `json:"links"`
}

func (s ReactForceGraphMarshaler) MarshalJSON() ([]byte, error) {
	nodes := []ForceGraphNode{}
	links := []ForceGraphLink{}
	for k, list := range s {
		s, ok := k.(string)
		if !ok {
			return nil, errors.New("unable to cast key from adjacency list to string")
		}
		nodes = append(nodes, ForceGraphNode{s})

		for link := range list.Neighbors {
			target, ok := link.(string)
			if !ok {
				return nil, errors.New("unable to cast target from adjacency list to string")
			}
			links = append(links, ForceGraphLink{s, target})
		}
	}

	graph := forceGraph{nodes, links}
	return json.Marshal(graph)
}
