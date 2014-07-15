package main

import (
	"log"

	"github.com/ffjiang/PhotonReader/seamcarver"
)

func main() {
	log.Printf("Hello, world.")

	lumMatrix := seamcarver.LoadImage("images/sampletext.jpg")
	for _, row := range lumMatrix.Matrix {
		log.Printf("%v", row)
	}
	imgGraph := seamcarver.SetWeights(lumMatrix)
	seamcarver.Carve(imgGraph)
}
