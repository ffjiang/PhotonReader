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
func (v Vertex) HashCode() {
	prime := 31
	result := 1
	return result
}

// Checks for equality of location
func (v Vertex) Equals(w Vertex) bool {
	if v.X == w.X && v.Y == w.Y {
		return true
	}
	return false
}

type Edge struct {
	ID     string
	Source *Vertex
	Dest   *Vertex
	Weight int
}

func (e Edge) CalcWeight() {
	xDist := math.Abs(float64(e.Dest.X - e.Source.X))
	yDist := math.Abs(float64(e.Dest.Y - e.Source.Y))
	return math.Sqrt((xDist * xDist) + (yDist * yDist))
}

func (e *Edge) SetWeight() {
	e.Weight = e.CalcWeight()
}

type Graph struct {
	Vertexes []Vertex
	Edges    []Edge
}
