package network

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oursimon/flowmo/flowmoerrors"
)

func Test_AddEdge_invalidCapacity_shouldReturnError(t *testing.T) {
	net := New()
	err := net.AddEdge(0, 1, -1)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_AddEdge_outOfRangeFrom_shouldReturnError(t *testing.T) {
	net := New()
	net.AddNode()
	err := net.AddEdge(1, 0, 10)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_AddEdge_outOfRangeTo_shouldReturnError(t *testing.T) {
	net := New()
	net.AddNode()
	err := net.AddEdge(0, 1, 10)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_AddEdge_invalidFrom_shouldReturnError(t *testing.T) {
	net := New()
	net.AddNode()
	err := net.AddEdge(-1, 0, 10)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_AddEdge_invalidTo_shouldReturnError(t *testing.T) {
	net := New()
	net.AddNode()
	err := net.AddEdge(0, -1, 10)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_AddEdge_happyCase(t *testing.T) {
	net := New()
	from := net.AddNode()
	to := net.AddNode()
	capacity := 15

	err := net.AddEdge(from, to, capacity)
	assert.NoError(t, err)

	assert.Len(t, net.adj, 2)
	assert.Len(t, net.adj[from], 1)
	assert.Len(t, net.adj[to], 1)

	forwardEdge := net.adj[from][0]
	assert.Equal(t, to, forwardEdge.to)
	assert.Equal(t, capacity, forwardEdge.capacity)
	assert.Equal(t, capacity, forwardEdge.initialCapacity)

	residualEdge := net.adj[to][0]
	assert.Equal(t, from, residualEdge.to)
	assert.Equal(t, 0, residualEdge.capacity)
	assert.Equal(t, 0, residualEdge.initialCapacity)
}

func Test_IncomingFlowByNode_beforeMaxFlow_shouldReturnZero(t *testing.T) {
	net := New()
	from := net.AddNode()
	to := net.AddNode()
	capacity := 10

	_ = net.AddEdge(from, to, capacity)

	incoming, _ := net.IncomingFlowByNode(to)
	assert.Equal(t, 0, incoming)
}

func Test_IncomingFlowByNode_afterMaxFlow_shouldReturnFlow(t *testing.T) {
	net := New()
	from := net.AddNode()
	to := net.AddNode()
	capacity := 10

	_ = net.AddEdge(from, to, capacity)

	_, _ = net.MaxFlow(from, to)

	incoming, _ := net.IncomingFlowByNode(to)
	assert.Equal(t, capacity, incoming)
}

func Test_OutgoingFlowByNode_beforeMaxFlow_shouldReturnZero(t *testing.T) {
	net := New()
	from := net.AddNode()
	to := net.AddNode()
	capacity := 10

	_ = net.AddEdge(from, to, capacity)

	outgoing, _ := net.OutgoingFlowByNode(from)
	assert.Equal(t, 0, outgoing)
}

func Test_OutgoingFlowByNode_afterMaxFlow_shouldReturnFlow(t *testing.T) {
	net := New()
	from := net.AddNode()
	to := net.AddNode()
	capacity := 10

	_ = net.AddEdge(from, to, capacity)

	_, _ = net.MaxFlow(from, to)

	outgoing, _ := net.OutgoingFlowByNode(from)
	assert.Equal(t, capacity, outgoing)
}

func Test_AddEdge_selfLoop_shouldWork(t *testing.T) {
	// Self-loop: edge from node to itself
	net := New()
	node := net.AddNode()

	err := net.AddEdge(node, node, 10)
	assert.NoError(t, err)

	// Self-loop should contribute 0 to max flow
	flow, err := net.MaxFlow(node, node)
	assert.NoError(t, err)
	assert.Equal(t, 0, flow)
}

func Test_IncomingFlowByNode_invalidNode_shouldReturnError(t *testing.T) {
	net := New()
	_ = net.AddNode()

	incoming, err := net.IncomingFlowByNode(5)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
	assert.Equal(t, -1, incoming)
}

func Test_OutgoingFlowByNode_invalidNode_shouldReturnError(t *testing.T) {
	net := New()
	_ = net.AddNode()

	outgoing, err := net.OutgoingFlowByNode(5)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
	assert.Equal(t, -1, outgoing)
}
