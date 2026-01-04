# Flowmo

A Go library for solving maximum flow problems in directed networks using Dinic's algorithm.

## Installation

```bash
go get github.com/oursimon/flowmo
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/oursimon/flowmo"
)

func main() {
    // Create a new flow network
    f := flowmo.New()

    // Add directed edges with capacities
    f.AddEdge("source", "a", 10)
    f.AddEdge("source", "b", 5)
    f.AddEdge("a", "sink", 7)
    f.AddEdge("b", "sink", 8)
    f.AddEdge("a", "b", 3)

    // Compute maximum flow
    maxFlow, err := f.MaxFlow("source", "sink")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Maximum flow: %d\n", maxFlow)

    // Query node capacities
    incoming, _ := f.IncomingCapacityByNode("sink")
    fmt.Printf("Incoming capacity to sink: %d\n", incoming)
}
```

## API Reference

### `New() *Flowmo`
Creates a new flow network instance.

### `AddEdge(from, to Node, capacity int) error`
Adds a directed edge from `from` to `to` with the specified capacity. Nodes are created automatically if they don't exist.

### `MaxFlow(source, sink Node) (int, error)`
Computes the maximum flow from the source node to the sink node using Dinic's algorithm.

### `IncomingCapacityByNode(node Node) (int, error)`
Returns the total incoming flow capacity for the specified node.

### `OutgoingCapacityByNode(node Node) (int, error)`
Returns the total outgoing flow capacity for the specified node.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

