package seamcarver

import (
	"image"
	"math"
)

type ImageMatrix struct {
	Matrix  [][]float64
	NumRows int
	NumCols int
}

type Vertex struct {
	ID    string
	Name  string
	Point image.Point
}

// Fix this stupid hash method
func (v Vertex) HashCode() int64 {
	prime := 31
	result := 1
	return int64(result * prime)
}

// Checks for equality of location
func (v Vertex) Equals(w Vertex) bool {
	if v.Point.X == w.Point.X && v.Point.Y == w.Point.Y {
		return true
	}
	return false
}

type Edge struct {
	ID     string
	Source *Vertex
	Dest   *Vertex
	Weight float64
}

// this function is wrong
func (e Edge) CalcWeight() float64 {
	xDist := math.Abs(float64(e.Dest.Point.X - e.Source.Point.X))
	yDist := math.Abs(float64(e.Dest.Point.Y - e.Source.Point.Y))
	return math.Sqrt((xDist * xDist) + (yDist * yDist))
}

func (e *Edge) SetWeight() {
	e.Weight = e.CalcWeight()
}

type Graph struct {
	Vertexes []Vertex
	Edges    []Edge
}
