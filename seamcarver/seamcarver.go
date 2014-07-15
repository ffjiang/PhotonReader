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

func Carve(lumMatrix LuminanceMatrix) {
	imgGraph := make(Graph, lumMatrix.NumRows)
	for i := range lumMatrix.Matrix {
		imgGraph[i] = make([]Vertex, lumMatrix.NumCols)
		for j := range lumMatrix.Matrix[i] {
			imgGraph[i][j] = Vertex{
				Cost:   math.MaxFloat64,
				Weight: [8]float64{},
			}
			// North
			if i > 0 {
				imgGraph[i][j].Weight[0] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i-1][j])
			} else {
				imgGraph[i][j].Weight[0] = -1
			}
			// North-east
			if i > 0 && j < lumMatrix.NumCols {
				imgGraph[i][j].Weight[1] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i-1][j+1]) / math.Sqrt2
			} else {
				imgGraph[i][j].Weight[1] = -1
			}
			// East
			if j < lumMatrix.NumCols {
				imgGraph[i][j].Weight[2] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i][j+1])
			} else {
				imgGraph[i][j].Weight[2] = -1
			}
			// South-east
			if i < lumMatrix.Matrix.NumRows && j < lumMatrix.Matrix.NumCols {
				imgGraph[i][j].Weight[3] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i+1][j+1]) / math.Sqrt2
			} else {
				imgGraph[i][j].Weight[3] = -1
			}
			// South
			if i < lumMatrix.Matrix.NumRows {
				imgGraph[i][j].Weight[4] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i+1][j])
			} else {
				imgGraph[i][j].Weight[4] = -1
			}
			// South-west
			if i < lumMatrix.Matrix.NumRows && j > 0 {
				imgGraph[i][j].Weight[5] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i+1][j-1]) / math.Sqrt2
			} else {
				imgGraph[i][j].Weight[5] = 0
			}
			// West
			if j > 0 {
				imgGraph[i][j].Weight[6] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i][j-1])
			} else {
				imgGraph[i][j].Weight[6] = -1
			}
			// North-west
			if i > 0 && j > 0 {
				imgGraph[i][j].Weight[7] = math.Abs(lumMatrix.Matrix[i][j] - lumMatrix.Matrix[i-1][j-1]) / math.Sqrt2
			} else {
				imgGraph[i][j].Weight[7] = -1
			}
		}
		}
	}
}

// Djikstra's shortest path algorithm
func ShortestPath() {

}
