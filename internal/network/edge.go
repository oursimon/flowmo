package network

type edge struct {
	to              int
	capacity        int
	initialCapacity int
	reverse         int
}

func (e *edge) isResidual() bool {
	return e.initialCapacity == 0
}
