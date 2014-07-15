package seamcarver

import (
	"log"
	"os"
	//"image"
	"image/color"
	"image/jpeg"
	"math"
)

func LoadImage(fileName string) LuminanceMatrix {
	// Load the image file
	imgFile, err := os.Open(fileName)
	if err != nil {
		log.Printf("Cannot open image file: %v", err)
	}
	defer imgFile.Close()

	// Decode the image file into an Image object
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		log.Printf("Cannot decode jpeg: %v", err)
	}

	log.Printf("%T,%v", img, img)

	// Initialise the LuminanceMatrix object...
	lumMatrix := LuminanceMatrix{
		NumRows: img.Bounds().Max.X - img.Bounds().Min.X,
		NumCols: img.Bounds().Max.Y - img.Bounds().Min.Y,
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
	log.Printf("%v, %T", rgbaColour, rgbaColour)
	r, g, b := rgbaColour.R, rgbaColour.G, rgbaColour.B
	log.Printf("%v, %v, %v rgb", r, g, b)
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}

func SetWeights(lumMatrix LuminanceMatrix) ImageGraph {
	imgGraph := make(ImageGraph, lumMatrix.NumRows)
	for i := range lumMatrix.Matrix {
		imgGraph[i] = make([]Vertex, lumMatrix.NumCols)
		for j := range lumMatrix.Matrix[i] {
			imgGraph[i][j] = Vertex{
				Cost:    math.MaxFloat64,
				Weights: [8]float64{},
			}
			// North
			if i > 0 {
				imgGraph[i][j].Weights[0] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i-1][j])
			} else {
				imgGraph[i][j].Weights[0] = -1
			}
			// North-east
			if i > 0 && j < lumMatrix.NumCols {
				imgGraph[i][j].Weights[1] = math.Abs(lumMatrix.Matrix[i][j]-lumMatrix.Matrix[i-1][j+1]) / math.Sqrt2
			} else {
				imgGraph[i][j].Weights[1] = -1
			}
			// East
			if j < lumMatrix.NumCols {
				imgGraph[i][j].Weights[2] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i][j+1])
			} else {
				imgGraph[i][j].Weights[2] = -1
			}
			// South-east
			if i < lumMatrix.NumRows && j < lumMatrix.NumCols {
				imgGraph[i][j].Weights[3] = math.Abs(lumMatrix.Matrix[i][j]-lumMatrix.Matrix[i+1][j+1]) / math.Sqrt2
			} else {
				imgGraph[i][j].Weights[3] = -1
			}
			// South
			if i < lumMatrix.NumRows {
				imgGraph[i][j].Weights[4] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i+1][j])
			} else {
				imgGraph[i][j].Weights[4] = -1
			}
			// South-west
			if i < lumMatrix.NumRows && j > 0 {
				imgGraph[i][j].Weights[5] = math.Abs(lumMatrix.Matrix[i][j]-lumMatrix.Matrix[i+1][j-1]) / math.Sqrt2
			} else {
				imgGraph[i][j].Weights[5] = 0
			}
			// West
			if j > 0 {
				imgGraph[i][j].Weights[6] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i][j-1])
			} else {
				imgGraph[i][j].Weights[6] = -1
			}
			// North-west
			if i > 0 && j > 0 {
				imgGraph[i][j].Weights[7] = math.Abs(lumMatrix.Matrix[i][j]-lumMatrix.Matrix[i-1][j-1]) / math.Sqrt2
			} else {
				imgGraph[i][j].Weights[7] = -1
			}
		}
	}
	return imgGraph
}

func Carve(imgGraph ImageGraph) {
	graphLength := len(imgGraph)
	for i := 0; i < graphLength; i += 10 {
		log.Printf("hi")
	}
}

// Djikstra's shortest path algorithm
func ShortestPath() {

}
