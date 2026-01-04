package network

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oursimon/flowmo/flowmoerrors"
)

// Test_MaxFlow_singleEdge:
//
//	(0) --10--> (1)
//
// maxFlow = 10
func Test_MaxFlow_singleEdge(t *testing.T) {
	net := New()
	nodeZero := net.AddNode()
	nodeOne := net.AddNode()

	_ = net.AddEdge(nodeZero, nodeOne, 10)
	flow, err := maxFlow(net, nodeZero, nodeOne)
	assert.NoError(t, err)
	assert.Equal(t, 10, flow)
}

// Test_MaxFlow_parallelPaths:
//
//	    5
//	(0) ---> (1)
//	 \         \
//	  \8        \5
//	   \         \
//	    > (2) ----> (3)
//	         6
//
// maxFlow = 5 + 6 = 11
func Test_MaxFlow_parallelPaths(t *testing.T) {
	net := New()
	net.AddNode()
	net.AddNode()
	net.AddNode()
	net.AddNode()

	_ = net.AddEdge(0, 1, 5)
	_ = net.AddEdge(0, 2, 8)
	_ = net.AddEdge(1, 3, 5)
	_ = net.AddEdge(2, 3, 6)

	flow, err := maxFlow(net, 0, 3)
	assert.NoError(t, err)
	assert.Equal(t, 11, flow)
}

// Test_MaxFlow_bottleneck:
//
//	(0) --100--> (1) --1--> (2) --100--> (3)
//
// Bottleneck at edge 1->2
// maxFlow = 1
func Test_MaxFlow_bottleneck(t *testing.T) {
	net := New()
	net.AddNode()
	net.AddNode()
	net.AddNode()
	net.AddNode()

	_ = net.AddEdge(0, 1, 100)
	_ = net.AddEdge(1, 2, 1)
	_ = net.AddEdge(2, 3, 100)

	flow, err := maxFlow(net, 0, 3)
	assert.NoError(t, err)
	assert.Equal(t, 1, flow)
}

// Test_MaxFlow_MultiEdge:
//
//	(0) --5--> (1)
//	(0) --7--> (1)
//
// maxFlow = 5 + 7 = 12
func Test_MaxFlow_multiEdge(t *testing.T) {
	net := New()

	net.AddNode()
	net.AddNode()

	_ = net.AddEdge(0, 1, 5)
	_ = net.AddEdge(0, 1, 7)

	flow, err := maxFlow(net, 0, 1)
	assert.NoError(t, err)
	assert.Equal(t, 12, flow)
}

// Test_MaxFlow_disjointPaths:
//
//	(0) --3--> (1) --3--> (2)
//	 \                     /
//	  \--------1----------/
//
// maxFlow = 4
func Test_MaxFlow_disjointPaths(t *testing.T) {
	net := New()
	net.AddNode()
	net.AddNode()
	net.AddNode()

	_ = net.AddEdge(0, 1, 3)
	_ = net.AddEdge(1, 2, 3)
	_ = net.AddEdge(0, 2, 1)

	flow, err := maxFlow(net, 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 4, flow)
}

// Test_MaxFlow_cycle:
//
//	(0) --5--> (1) --5--> (2) --4--> (3)
//	             ^---5-----|
//	                  cycle
//
// maxFlow = 4
func Test_MaxFlow_cycle(t *testing.T) {
	net := New()
	net.AddNode()
	net.AddNode()
	net.AddNode()
	net.AddNode()

	_ = net.AddEdge(0, 1, 5)
	_ = net.AddEdge(1, 2, 5)
	_ = net.AddEdge(2, 3, 4)
	_ = net.AddEdge(2, 1, 5) // cycle

	flow, err := maxFlow(net, 0, 3)
	assert.NoError(t, err)
	assert.Equal(t, 4, flow)
}

// Test_MaxFlow_sourceEqualsSink:
// maxFlow = 0.
func Test_MaxFlow_sourceEqualsSink(t *testing.T) {
	net := New()
	net.AddNode()

	flow, err := maxFlow(net, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, flow)
}

// Test_MaxFlow_zeroCapacity:
//
//	(0) --0--> (1)
//	(0) --5--> (1)
//
// maxFlow = 5
func Test_MaxFlow_zeroCapacity(t *testing.T) {
	net := New()
	net.AddNode()
	net.AddNode()

	_ = net.AddEdge(0, 1, 0)
	_ = net.AddEdge(0, 1, 5)

	flow, err := maxFlow(net, 0, 1)
	assert.NoError(t, err)
	assert.Equal(t, 5, flow)
}

// Test_MaxFlow_internalSourceAndSink:
//
//	(0) --100--> (1)
//	           /   \
//	        3 /     \ 2
//	         v       v
//	        (2)     (3)
//	         \       /
//	        2 \     / 3
//	          v   v
//	           (4) ------100 ------- (5)
//
// maxFlow from 1 to 4 = 2 + 2 = 4
func Test_MaxFlow_internalSourceAndSink(t *testing.T) {
	net := New()
	net.AddNode() // 0
	net.AddNode() // 1
	net.AddNode() // 2
	net.AddNode() // 3
	net.AddNode() // 4
	net.AddNode() // 5

	// Edges that shouldn't matter.
	_ = net.AddEdge(0, 1, 100)
	_ = net.AddEdge(4, 5, 100)

	// Internal subgraph where we measure the flow from 1 to 4.
	_ = net.AddEdge(1, 2, 3)
	_ = net.AddEdge(1, 3, 2)
	_ = net.AddEdge(2, 4, 2)
	_ = net.AddEdge(3, 4, 3)

	flow, err := maxFlow(net, 1, 4)
	assert.NoError(t, err)
	assert.Equal(t, 4, flow)
}

func Test_newDinic_givenNilAdjacencyList(t *testing.T) {
	_, err := newDinic(nil, 0, 1)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_newDinic_givenInvalidSink(t *testing.T) {
	_, err := newDinic([][]*edge{}, -1, 1)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_newDinic_givenOutOfRangeSink(t *testing.T) {
	net := New()
	_ = net.AddNode()
	_ = net.AddNode()
	_ = net.AddEdge(0, 1, 100)
	// source index 2 is out of range
	_, err := newDinic(net.adj, 1, 2)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_newDinic_givenInvalidSource(t *testing.T) {
	_, err := newDinic([][]*edge{}, -1, 1)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_newDinic_givenOutOfRangeSource(t *testing.T) {
	net := New()
	_ = net.AddNode()
	_ = net.AddNode()
	_ = net.AddEdge(0, 1, 100)
	// source index 2 is out of range
	_, err := newDinic(net.adj, 2, 1)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}
