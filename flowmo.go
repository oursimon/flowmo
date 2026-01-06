package flowmo

import (
	"fmt"

	"github.com/oursimon/flowmo/flowmoerrors"
	"github.com/oursimon/flowmo/internal/network"
)

type Node string

// Flowmo represents a flow network that computes maximum flow using Dinic's algorithm.
// It maps string-based node identifiers to an internal graph representation.
type Flowmo struct {
	indexByNode map[Node]int
	network     *network.Network
}

func New() *Flowmo {
	net := network.New()
	return &Flowmo{
		indexByNode: make(map[Node]int),
		network:     net,
	}
}

// AddEdge adds a directed edge from the "from" node to the "to" node with the specified capacity.
// If either node does not exist, it will be created automatically.
// Multiple edges between the same pair of nodes are allowed.
func (f *Flowmo) AddEdge(from, to Node, capacity int) error {
	fromIndex, err := f.addNode(from)
	if err != nil {
		return err
	}

	toIndex, err := f.addNode(to)
	if err != nil {
		return err
	}

	return f.network.AddEdge(fromIndex, toIndex, capacity)
}

// MaxFlow computes the maximum flow from the source node
// to the sink node using Dinic's algorithm.
// The algorithm runs in O(V²E) time complexity where V is
// the number of nodes and E is the number of edges.
//
// Returns the maximum flow value and nil error on success.
// Returns 0 and an error if:
//   - Source node does not exist (ErrNotFound)
//   - Sink node does not exist (ErrNotFound)
//
// If source equals sink, returns 0 with no error (valid edge case).
// If no path exists from source to sink, returns 0 with no error.
func (f *Flowmo) MaxFlow(source, sink Node) (int, error) {
	sourceIndex, exists := f.indexByNode[source]
	if !exists {
		return 0, fmt.Errorf(
			"source node %q: %w",
			source,
			flowmoerrors.ErrNotFound,
		)
	}

	sinkIndex, exists := f.indexByNode[sink]
	if !exists {
		return 0, fmt.Errorf(
			"sink node %q: %w",
			sink,
			flowmoerrors.ErrNotFound,
		)
	}

	return f.network.MaxFlow(sourceIndex, sinkIndex)
}

// FlowByNode returns the total outgoing flow
// (used capacity) from the specified node.
// This represents the sum of all flow that has been pushed out of
// the node across all outgoing edges.
// This method should be called after MaxFlow has been computed.
//
// Returns the total outgoing flow and nil error on success.
// Returns -1 and an error if:
//   - Node does not exist (ErrNotFound)
func (f *Flowmo) FlowByNode(node Node) (int, error) {
	idx, exists := f.indexByNode[node]
	if !exists {
		return -1, fmt.Errorf(
			"node %q: %w",
			node,
			flowmoerrors.ErrNotFound,
		)
	}

	return f.network.FlowByNode(idx)
}

func (f *Flowmo) addNode(node Node) (int, error) {
	if node == "" {
		return -1, fmt.Errorf(
			"node cannot be empty: %w",
			flowmoerrors.ErrInvalidArgument,
		)
	}

	if idx, exists := f.indexByNode[node]; exists {
		return idx, nil
	}

	idx := f.network.AddNode()
	f.indexByNode[node] = idx

	return idx, nil
}
