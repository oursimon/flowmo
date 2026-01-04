package network

import (
	"fmt"

	"github.com/oursimon/flowmo/flowmoerrors"
)

const infinity = 1 << 50

type dinic struct {
	adj       [][]*edge
	nodeLevel []int
	iterator  []int
	source    int
	sink      int
}

func maxFlow(net *Network, source, sink int) (int, error) {
	if net == nil {
		return 0, fmt.Errorf(
			"network is nil: %w",
			flowmoerrors.ErrInvalidArgument,
		)
	}

	if source == sink {
		return 0, nil
	}

	d, err := newDinic(net.adj, source, sink)
	if err != nil {
		return 0, err
	}

	return d.run()
}

func newDinic(adj [][]*edge, source, sink int) (*dinic, error) {
	nrOfNodes := len(adj)
	invalidSource := source < 0 || source >= nrOfNodes
	if invalidSource {
		return nil, fmt.Errorf(
			"source node %d: %w",
			source,
			flowmoerrors.ErrInvalidArgument,
		)
	}

	invalidSink := sink < 0 || sink >= nrOfNodes
	if invalidSink {
		return nil, fmt.Errorf(
			"sink node %d: %w",
			sink,
			flowmoerrors.ErrInvalidArgument,
		)
	}

	return &dinic{
		adj:       adj,
		nodeLevel: make([]int, nrOfNodes),
		iterator:  make([]int, nrOfNodes),
		source:    source,
		sink:      sink,
	}, nil
}

func (d *dinic) run() (int, error) {
	total := 0
	for d.buildLevelGraph() {
		d.resetIterator()
		for {
			pushed := d.sendFlow(d.source, infinity)
			if pushed == 0 {
				break
			}

			total += pushed
		}
	}

	return total, nil
}

func (d *dinic) buildLevelGraph() bool {
	d.resetLevel()

	queue := []int{d.source}
	d.nodeLevel[d.source] = 0
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		for _, e := range d.adj[currentNode] {
			hasCapacity := e.capacity > 0
			isUnvisited := d.nodeLevel[e.to] < 0

			if !hasCapacity || !isUnvisited {
				continue
			}

			d.nodeLevel[e.to] = d.nodeLevel[currentNode] + 1
			queue = append(queue, e.to)
		}
	}

	return d.nodeLevel[d.sink] >= 0
}

// nolint:gocyclo
func (d *dinic) sendFlow(currentNode, incomingFlow int) int {
	// Base case: we reached the sink, so we can push everything we carried here.
	if currentNode == d.sink {
		return incomingFlow
	}

	edges := d.adj[currentNode]
	for d.iterator[currentNode] < len(edges) {
		edgeIndex := d.iterator[currentNode]
		e := edges[edgeIndex]

		d.iterator[currentNode]++

		if e.capacity == 0 {
			continue
		}

		isOnNextLevel := d.nodeLevel[e.to] == d.nodeLevel[currentNode]+1
		if !isOnNextLevel {
			continue
		}

		flowLimit := incomingFlow
		if e.capacity < flowLimit {
			flowLimit = e.capacity
		}

		// push it further down the path (dfs)
		pushed := d.sendFlow(e.to, flowLimit)
		if pushed == 0 {
			continue
		}

		e.capacity -= pushed
		reverseEdge := d.adj[e.to][e.reverse]
		reverseEdge.capacity += pushed

		return pushed
	}

	// No augmenting path from this node.
	return 0
}

func (d *dinic) resetLevel() {
	for i := range d.nodeLevel {
		d.nodeLevel[i] = -1
	}
}

func (d *dinic) resetIterator() {
	for i := range d.iterator {
		d.iterator[i] = 0
	}
}
