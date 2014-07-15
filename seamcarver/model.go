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
	// The previous node in the path
	Previous image.Point
}

// Fix this stupid hash method
func (v Vertex) HashCode() int64 {
	prime := 31
	result := 1
	return int64(result * prime)
}

type ImageGraph [][]Vertex

type Path []image.Point
