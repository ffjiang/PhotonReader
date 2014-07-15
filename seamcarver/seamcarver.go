package seamcarver

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func LoadImage(fileName string) image.Image {
	// Load the image file
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Cannot open image file: %v", err)
	}
	defer file.Close()

	// Decode the image file into an Image object
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Printf("Cannot decode jpeg: %v", err)
	}

	return img
}

func CreateLumMatrix(img image.Image) LuminanceMatrix {
	// Initialise the LuminanceMatrix object...
	lumMatrix := LuminanceMatrix{
		NumRows: img.Bounds().Dx(),
		NumCols: img.Bounds().Dy(),
	}
	lumMatrix.Matrix = make([][]float64, lumMatrix.NumRows)
	for i := range lumMatrix.Matrix {
		lumMatrix.Matrix[i] = make([]float64, lumMatrix.NumCols)
		for j := range lumMatrix.Matrix[i] {
			lumMatrix.Matrix[i][j] = Luminance(img.At(i, j))
		}
	}
	return lumMatrix
}

func Luminance(colour color.Color) float64 {
	rgbaColour := color.RGBAModel.Convert(colour).(color.RGBA)
	r, g, b := rgbaColour.R, rgbaColour.G, rgbaColour.B
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}

func SetWeights(lumMatrix LuminanceMatrix) ImageGraph {
	imgGraph := make(ImageGraph, lumMatrix.NumRows)
	for x := range lumMatrix.Matrix {
		imgGraph[x] = make([]Vertex, lumMatrix.NumCols)
		for y := range lumMatrix.Matrix[x] {
			imgGraph[x][y] = Vertex{
				Cost:    math.MaxFloat64,
				Weights: [8]float64{},
				EndZone: false,
			}

			// Each of the weights is the magnitude of the energy gradient, energy being luminance

			// North
			if y > 0 {
				imgGraph[x][y].Weights[0] = math.Abs(lumMatrix.Matrix[x][y] - lumMatrix.Matrix[x][y-1])
			} else {
				imgGraph[x][y].Weights[0] = -1
			}
			// North-east
			if x < lumMatrix.NumCols-1 && y > 0 {
				imgGraph[x][y].Weights[1] = math.Abs(lumMatrix.Matrix[x][y]-lumMatrix.Matrix[x+1][y-1]) * math.Sqrt2
			} else {
				imgGraph[x][y].Weights[1] = -1
			}
			// East
			if x < lumMatrix.NumCols-1 {
				imgGraph[x][y].Weights[2] = math.Abs(lumMatrix.Matrix[x][y] - lumMatrix.Matrix[x+1][y])
			} else {
				imgGraph[x][y].Weights[2] = -1
			}
			// South-east
			if x < lumMatrix.NumCols-1 && y < lumMatrix.NumRows-1 {
				imgGraph[x][y].Weights[3] = math.Abs(lumMatrix.Matrix[x][y]-lumMatrix.Matrix[x+1][y+1]) * math.Sqrt2
			} else {
				imgGraph[x][y].Weights[3] = -1
			}
			// South
			if y < lumMatrix.NumRows-1 {
				imgGraph[x][y].Weights[4] = math.Abs(lumMatrix.Matrix[x][y] - lumMatrix.Matrix[x][y+1])
			} else {
				imgGraph[x][y].Weights[4] = -1
			}
			// South-west
			if x > 0 && y < lumMatrix.NumRows-1 {
				imgGraph[x][y].Weights[5] = math.Abs(lumMatrix.Matrix[x][y]-lumMatrix.Matrix[x-1][y+1]) * math.Sqrt2
			} else {
				imgGraph[x][y].Weights[5] = 0
			}
			// West
			if x > 0 {
				imgGraph[x][y].Weights[6] = math.Abs(lumMatrix.Matrix[x][y] - lumMatrix.Matrix[x-1][y])
			} else {
				imgGraph[x][y].Weights[6] = -1
			}
			// North-west
			if x > 0 && jy > 0 {
				imgGraph[x][y].Weights[7] = math.Abs(lumMatrix.Matrix[x][y]-lumMatrix.Matrix[x-1][y-1]) * math.Sqrt2
			} else {
				imgGraph[x][y].Weights[7] = -1
			}
		}
	}
	return imgGraph
}

func Carve(srcImg image.Image, imgGraph ImageGraph) {
	bounds := srcImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	// Create a new RGBA image to be manipulated
	dstImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, bounds.Min, draw.Src)

	// Set right hand side of page to be endzone
	for j := 0; j < height; j++ {
		imgGraph[width-1][j].EndZone = true
	}
	// Go along y-axis, or side of page
	for j := 0; j < height; j += 10 {
		ShortestPath(image.Point{X: 0, Y: j}, imgGraph)
	}

	// Write result
	if err := WriteImage(dstImg, "woopdedoop.jpg"); err != nil {
		log.Printf("Error writing image: %v", err)
	}
}

// Djikstra's shortest path algorithm
func ShortestPath(start image.Point, imgGraph ImageGraph) Path {
	visited := map[image.Point]bool{start: true}
	// Coords of current node
	currentNode := start
	imgGraph[start.X][start.Y].Cost = 0
	nodesUnvisited := len(imgGraph)*len(imgGraph[0]) - 1

	for nodesUnvisited > 0 {
		x := currentNode.X
		y := currentNode.Y
		// Evaluate neighbours and set costs/previouses
		// Choose closest neighbour
		// if neighbour.EndZone == true {end}
		// Add this neighbour to visited list

		// Evaluate neighbours and set costs/previouses
		cost := imgGraph[x][y].Cost
		N := imgGraph[x][y].Weights[0]
		NE := imgGraph[x][y].Weights[1]
		E := imgGraph[x][y].Weights[2]
		SE := imgGraph[x][y].Weights[3]
		S := imgGraph[x][y].Weights[4]
		SW := imgGraph[x][y].Weights[5]
		W := imgGraph[x][y].Weights[6]
		NW := imgGraph[x][y].Weights[7]

		minCostNode := image.Point{}
		minCost := math.MaxFloat64
		// North
		if N > 0 {
			if N+cost < imgGraph[x][y-1].Cost {
				imgGraph[x][y-1].Cost = N + cost
			}
			if imgGraph[x][y-1].Cost < minCost {
				minCost = imgGraph[x][y-1].Cost
				minCostNode.X, minCostNode.Y = x, y-1
			}
		}
		// North-east
		if NE > 0 {
			if NE+cost < imgGraph[x+1][y-1].Cost {
				imgGraph[x+1][y-1].Cost = NE + cost
			}
			if imgGraph[x+1][y-1].Cost < minCost {
				minCost = imgGraph[x+1][y-1].Cost
				minCostNode.X, minCostNode.Y = x+1, y-1
			}
		}
		// East
		if E > 0 {
			if E+cost < imgGraph[x+1][y].Cost {
				imgGraph[x+1][y].Cost = E + cost
			}
			if imgGraph[x+1][y].Cost < minCost {
				minCost = imgGraph[x+1][y].Cost
				minCostNode.X, minCostNode.Y = x+1, y
			}
		}
		// South-east
		if SE > 0 {
			if SE+cost < imgGraph[x+1][y+1].Cost {
				imgGraph[x+1][y+1].Cost = SE + cost
			}
			if imgGraph[x+1][y+1].Cost < minCost {
				minCost = imgGraph[x+1][y+1].Cost
				minCostNode.X, minCostNode.y = x+1, y+1
			}
		}
		// South

		nodesUnvisited--
	}

}

func WriteImage(img image.Image, filename string) error {
	// Open a file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	return nil
}
