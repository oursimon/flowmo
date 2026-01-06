package main

import (
	"fmt"

	"github.com/oursimon/flowmo"
)

func main() {
	// Create a new flow network.
	f := flowmo.New()

	// Add directed edges with capacities.
	// Nodes are identified by string labels.
	// If nodes do not exist, they will be created automatically.
	_ = f.AddEdge("a", "b", 1)
	_ = f.AddEdge("a", "c", 1)
	_ = f.AddEdge("c", "b", 1)

	// Compute the maximum flow from source "a" to sink "b".
	flow, err := f.MaxFlow("a", "b")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Max flow from 'a' to 'b': %d\n", flow)

	// Query the outgoing flow for specific nodes.
	outgoing, err := f.FlowByNode("c")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Outgoing flow from 'c': %d\n", outgoing)
}
