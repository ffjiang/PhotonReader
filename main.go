package main

import (
	"log"

	"github.com/ffjiang/PhotonReader/seamcarver"
)

func main() {
	log.Printf("Hello, world.")

	img := seamcarver.LoadImage("images/sampletext.jpg")

	lumMatrix := seamcarver.CreateLumMatrix(img)

	imgGraph := seamcarver.SetWeights(lumMatrix)

	seamcarver.Carve(img, imgGraph)
}
