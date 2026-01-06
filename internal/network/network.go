package network

import (
	"fmt"

	"github.com/oursimon/flowmo/flowmoerrors"
)

type Network struct {
	adj [][]*edge
}

func New() *Network {
	return &Network{
		adj: make([][]*edge, 0),
	}
}

func (net *Network) AddNode() int {
	net.adj = append(net.adj, []*edge{})
	nodeIdx := len(net.adj) - 1

	return nodeIdx
}

func (net *Network) AddEdge(from, to, capacity int) error {
	if capacity < 0 {
		return fmt.Errorf(
			"capacity must be non-negative %d: %w",
			capacity,
			flowmoerrors.ErrInvalidArgument)
	}

	invalidFrom := from < 0 || from >= len(net.adj)
	if invalidFrom {
		return fmt.Errorf(
			"from node %d: %w",
			from,
			flowmoerrors.ErrInvalidArgument)
	}

	invalidTo := to < 0 || to >= len(net.adj)
	if invalidTo {
		return fmt.Errorf(
			"to node %d: %w",
			to,
			flowmoerrors.ErrInvalidArgument,
		)
	}

	net.addEdge(from, to, capacity)
	net.addResidual(from, to)

	return nil
}

func (net *Network) MaxFlow(source, sink int) (int, error) {
	return maxFlow(net, source, sink)
}

func (net *Network) FlowByNode(node int) (int, error) {
	invalidNode := node < 0 || node >= len(net.adj)
	if invalidNode {
		return -1, fmt.Errorf(
			"node %d: %w",
			node,
			flowmoerrors.ErrInvalidArgument,
		)
	}

	out := 0
	edges := net.adj[node]
	for _, e := range edges {
		if e.isResidual() {
			continue
		}

		out += e.initialCapacity - e.capacity
	}

	return out, nil
}

func (net *Network) addEdge(from, to, capacity int) {
	toIndex := len(net.adj[to])
	forwardEdge := &edge{
		to:              to,
		capacity:        capacity,
		initialCapacity: capacity,
		reverse:         toIndex,
	}

	net.adj[from] = append(
		net.adj[from],
		forwardEdge,
	)
}

func (net *Network) addResidual(from, to int) {
	fromIndex := len(net.adj[from])
	residualEdge := &edge{
		to:       from,
		capacity: 0,
		reverse:  fromIndex,
	}
	net.adj[to] = append(
		net.adj[to],
		residualEdge,
	)
}
