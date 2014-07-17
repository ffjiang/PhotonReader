package seamcarver

import (
	"image"
)

type LuminanceMatrix struct {
	Matrix  [][]float64
	NumRows int
	NumCols int
}

type Vertex struct {
	Cost float64
	// Energy gradient magnitude in each direction, clockwise starting from north
	Weights [8]float64
	// Whether this Vertex is connected to the virtual node
	EndZone bool
	Coords  Point
	// The previous node in the path
	Previous Point
}

// Fix this stupid hash method
func (v Vertex) HashCode() int64 {
	prime := 31
	result := 1
	return int64(result * prime)
}

type ImageGraph [][]Vertex

type Path []Point

func (path *Path) Add(p Point) {
	*path = append(*path, p)
}

type Point image.Point

// Minimum priority queue built using the heap interface
type PriorityQueue []*Vertex

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	vertex := x.(*Vertex)
	*pq = append(*pq, vertex)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	vertex := old[n-1]
	*pq = old[:n-1]
	return vertex
}
