package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/ffjiang/PhotonReader/seamcarver"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	log.Printf("Hello, world.")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	img := seamcarver.LoadImage("images/sampletext.jpg")

	lumMatrix := seamcarver.CreateLumMatrix(img)

	imgGraph := seamcarver.SetWeights(lumMatrix)

	seamcarver.Carve(img, imgGraph)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}
