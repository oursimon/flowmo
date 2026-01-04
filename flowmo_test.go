package flowmo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oursimon/flowmo/flowmoerrors"
)

func Test_New(t *testing.T) {
	f := New()
	assert.NotNil(t, f)
	assert.NotNil(t, f.network)
	assert.NotNil(t, f.indexByNode)
	assert.Equal(t, 0, len(f.indexByNode))
}

func Test_AddEdge_basicFunctionality(t *testing.T) {
	f := New()

	err := f.AddEdge("a", "b", 10)
	assert.NoError(t, err)

	// Verify nodes were created
	assert.Equal(t, 2, len(f.indexByNode))
	assert.Contains(t, f.indexByNode, Node("a"))
	assert.Contains(t, f.indexByNode, Node("b"))
}

func Test_AddEdge_sameNodesPair(t *testing.T) {
	// Multiple edges between same nodes
	f := New()

	err := f.AddEdge("a", "b", 10)
	assert.NoError(t, err)

	err = f.AddEdge("a", "b", 5)
	assert.NoError(t, err)

	// Should accumulate capacity
	maxFlow, err := f.MaxFlow("a", "b")
	assert.NoError(t, err)
	assert.Equal(t, 15, maxFlow)
}

func Test_AddEdge_negativeCapacity(t *testing.T) {
	f := New()

	err := f.AddEdge("a", "b", -10)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_AddEdge_emptyFromNode(t *testing.T) {
	f := New()

	err := f.AddEdge("", "b", 10)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_AddEdge_emptyToNode(t *testing.T) {
	f := New()

	err := f.AddEdge("a", "", 10)
	assert.ErrorIs(t, err, flowmoerrors.ErrInvalidArgument)
}

func Test_AddEdge_zeroCapacity(t *testing.T) {
	f := New()

	// Zero capacity should be allowed
	err := f.AddEdge("a", "b", 0)
	assert.NoError(t, err)

	maxFlow, err := f.MaxFlow("a", "b")
	assert.NoError(t, err)
	assert.Equal(t, 0, maxFlow)
}

func Test_MaxFlow_simpleNetwork(t *testing.T) {
	//  a --10--> b
	f := New()

	_ = f.AddEdge("a", "b", 10)

	maxFlow, err := f.MaxFlow("a", "b")
	assert.NoError(t, err)
	assert.Equal(t, 10, maxFlow)
}

func Test_MaxFlow_sourceNotFound(t *testing.T) {
	f := New()

	_ = f.AddEdge("a", "b", 10)

	_, err := f.MaxFlow("x", "b")
	assert.ErrorIs(t, err, flowmoerrors.ErrNotFound)
}

func Test_MaxFlow_sinkNotFound(t *testing.T) {
	f := New()

	_ = f.AddEdge("a", "b", 10)

	_, err := f.MaxFlow("a", "x")
	assert.ErrorIs(t, err, flowmoerrors.ErrNotFound)
}

func Test_IncomingCapacityByNode_nodeNotFound(t *testing.T) {
	f := New()

	_ = f.AddEdge("a", "b", 10)

	_, err := f.IncomingCapacityByNode("x")
	assert.ErrorIs(t, err, flowmoerrors.ErrNotFound)
}

func Test_OutgoingCapacityByNode_nodeNotFound(t *testing.T) {
	f := New()

	_ = f.AddEdge("a", "b", 10)

	_, err := f.OutgoingCapacityByNode("x")
	assert.ErrorIs(t, err, flowmoerrors.ErrNotFound)
}

func Test_CompleteWorkflow(t *testing.T) {
	// End-to-end test: create graph, compute max flow, query capacities
	//       20        10
	//   a -----> b -----> d
	//   |        |        ^
	//   |10      |5       |15
	//   v        v        |
	//   c ---------------->
	f := New()

	_ = f.AddEdge("a", "b", 20)
	_ = f.AddEdge("a", "c", 10)
	_ = f.AddEdge("b", "c", 5)
	_ = f.AddEdge("b", "d", 10)
	_ = f.AddEdge("c", "d", 15)

	maxFlow, err := f.MaxFlow("a", "d")
	assert.NoError(t, err)
	assert.Equal(t, 25, maxFlow)

	// Verify incoming capacities
	incomingB, _ := f.IncomingCapacityByNode("b")
	assert.Equal(t, 15, incomingB)

	incomingC, _ := f.IncomingCapacityByNode("c")
	assert.Equal(t, 15, incomingC)

	incomingD, _ := f.IncomingCapacityByNode("d")
	assert.Equal(t, 25, incomingD)

	// Verify outgoing capacities
	outgoingA, _ := f.OutgoingCapacityByNode("a")
	assert.Equal(t, 25, outgoingA)

	outgoingB, _ := f.OutgoingCapacityByNode("b")
	assert.Equal(t, 15, outgoingB)

	outgoingC, _ := f.OutgoingCapacityByNode("c")
	assert.Equal(t, 15, outgoingC)
}
