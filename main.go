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
var showVisited = flag.Bool("showvisited", false, "whether or not to paint visited nodes")

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

	imgPath := "images/sampletext.jpg"
	if len(os.Args) > 1 && os.Args[1][0] != '-' {
		imgPath = os.Args[1]
	}

	img := seamcarver.LoadImage(imgPath)

	lumMatrix := seamcarver.CreateLumMatrix(img)

	imgGraph := seamcarver.SetWeights(lumMatrix)

	seamcarver.Carve(img, imgGraph, *showVisited)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}
