package main

import (
	"fmt"
	"log"
	"os"
	"image/jpeg"

	"github.com/ffjiang/PhotonReader/seamcarver"
)

func main() {
	fmt.Printf("Hello, world.")

	imgMatrix := loadImage("images/sampletext.jpg")
	for i, row := range imgMatrix.Matrix {
		fmt.Printf("%v", row)
	}
}

func loadImage(fileName string) seamcarver.ImageMatrix {
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

	// Initialise the ImageMatrix object...
	imgMatrix := ImageMatrix{
		NumRows: img.Bounds.Max.X - img.Bounds.Min.X + 1,
		NumCols: img.Bounds.Max.Y - img.Bounds.Min.Y + 1,
	}

	// ... with a matrix representing luminance.
	imgMatrix.Matrix := make([][]float64, imgMatrix.NumRows)
	for i, row := range imgMatrix.Matrix {
		row = make([]float64, imgMatrix.NumCols)
		for j, column := range row {
			column = Luminance(img.At(i, j))
		}
	}
	return imgMatrix
}
