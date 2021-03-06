package seamcarver

import (
	"container/heap"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"
)

// Distance to travel per iteration
const d int = 1

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
		NumCols: img.Bounds().Dx(),
		NumRows: img.Bounds().Dy(),
	}
	lumMatrix.Matrix = make([][]float64, lumMatrix.NumCols)
	for i := range lumMatrix.Matrix {
		lumMatrix.Matrix[i] = make([]float64, lumMatrix.NumRows)
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
	imgGraph := make(ImageGraph, lumMatrix.NumCols)
	for x := range lumMatrix.Matrix {
		imgGraph[x] = make([]Vertex, lumMatrix.NumRows)
		for y := range lumMatrix.Matrix[x] {
			imgGraph[x][y] = Vertex{
				Cost:     math.MaxFloat64,
				Weights:  [3]float64{},
				EndZone:  false,
				Coords:   Point{X: x, Y: y},
				Previous: Point{},
			}

			// Each of the weights is the magnitude of the energy gradient, energy being luminance

			// North-east
			if x < lumMatrix.NumCols-d && y >= d {
				imgGraph[x][y].Weights[0] = math.Abs(lumMatrix.Matrix[x][y]-lumMatrix.Matrix[x+1][y-1]) + 50
			} else {
				imgGraph[x][y].Weights[0] = -1
			}
			// East
			if x < lumMatrix.NumCols-d {
				imgGraph[x][y].Weights[1] = math.Abs(lumMatrix.Matrix[x][y] - lumMatrix.Matrix[x+1][y])
			} else {
				imgGraph[x][y].Weights[1] = -1
			}
			// South-east
			if x < lumMatrix.NumCols-d && y < lumMatrix.NumRows-d {
				imgGraph[x][y].Weights[2] = math.Abs(lumMatrix.Matrix[x][y]-lumMatrix.Matrix[x+1][y+1]) + 50
			} else {
				imgGraph[x][y].Weights[2] = -1
			}
		}
	}
	return imgGraph
}

func Carve(srcImg image.Image, imgGraph ImageGraph, showVisited bool) {
	bounds := srcImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	// Create a new RGBA image to be manipulated
	dstImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, bounds.Min, draw.Src)

	// Set right hand side of page to be endzone
	for j := 0; j < height; j++ {
		imgGraph[width-d][j].EndZone = true
	}

	// Go along y-axis, or side of page
	now := time.Now()
	for j := 0; j < height; j += 1000 {
		imgGraphCopy := make(ImageGraph, width)
		for x := range imgGraphCopy {
			imgGraphCopy[x] = make([]Vertex, height)
			copy(imgGraphCopy[x], imgGraph[x])
		}
		path, visited, avgHeight := ShortestPath(Point{X: 0, Y: j}, imgGraphCopy, showVisited)
		log.Printf("Time taken: %v", time.Since(now))
		for _, point := range visited {
			dstImg.Set(point.X, point.Y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
		for x := 0; x < width; x++ {
			dstImg.Set(x, avgHeight, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
		for _, point := range path {
			dstImg.Set(point.X, point.Y, color.Black)
		}
		log.Printf("Time to paint: %v", time.Since(now))
	}

	// Write result
	if err := WriteImage(dstImg, "woopdedoop.jpg"); err != nil {
		log.Printf("Error writing image: %v", err)
	}
}

// Djikstra's shortest path algorithm
func ShortestPath(start Point, imgGraph ImageGraph, showVisited bool) (Path, []Point, int) {
	log.Printf("ShortestPath called")
	// Setup data structures
	visited := map[Point]bool{}
	imgGraph[start.X][start.Y].Cost = 0
	visitableNodes := PriorityQueue{&imgGraph[start.X][start.Y]}
	heap.Init(&visitableNodes)
	isVisitable := map[Point]bool{start: true}
	var currentNode Point

	for len(visitableNodes) > 0 {
		currentNode = heap.Pop(&visitableNodes).(*Vertex).Coords
		isVisitable[currentNode] = false
		visited[currentNode] = true
		x := currentNode.X
		y := currentNode.Y

		if imgGraph[x][y].EndZone {
			break
		}
		// Evaluate neighbours and set costs/previouses
		// Choose closest neighbour
		// if neighbour.EndZone == true {end}
		// Add this neighbour to visited list

		// Evaluate neighbours and set costs/previouses...
		cost := imgGraph[x][y].Cost
		NE := imgGraph[x][y].Weights[0]
		E := imgGraph[x][y].Weights[1]
		SE := imgGraph[x][y].Weights[2]

		// .. and add to visitable if unvisited

		// North-east
		nextNode := Point{X: x + d, Y: y - d}
		if NE >= 0 && !visited[nextNode] {
			if NE+cost < imgGraph[x+d][y-d].Cost {
				imgGraph[x+d][y-d].Cost = NE + cost
				imgGraph[x+d][y-d].Previous = currentNode
			}
			if !isVisitable[nextNode] {
				isVisitable[nextNode] = true
				heap.Push(&visitableNodes, &imgGraph[x+d][y-d])
			}
		}
		// East
		nextNode = Point{X: x + d, Y: y}
		if E >= 0 && !visited[nextNode] {
			if E+cost < imgGraph[x+d][y].Cost {
				imgGraph[x+d][y].Cost = E + cost
				imgGraph[x+d][y].Previous = currentNode
			}
			if !isVisitable[nextNode] {
				isVisitable[nextNode] = true
				heap.Push(&visitableNodes, &imgGraph[x+d][y])
			}
		}
		// South-east
		nextNode = Point{X: x + d, Y: y + d}
		if SE >= 0 && !visited[nextNode] {
			if SE+cost < imgGraph[x+d][y+d].Cost {
				imgGraph[x+d][y+d].Cost = SE + cost
				imgGraph[x+d][y+d].Previous = currentNode
			}
			if !isVisitable[nextNode] {
				isVisitable[nextNode] = true
				heap.Push(&visitableNodes, &imgGraph[x+d][y+d])
			}
		}
	}
	log.Printf("Shortest path found")

	f, err1 := os.Create("memprofile.mprof")
	if err1 != nil {
		log.Fatal(err1)
	}
	pprof.WriteHeapProfile(f)
	f.Close()

	path := Path{}
	totalHeight := 0
	count := 0
	oldY := currentNode.Y
	for currentNode != start {
		path.Add(currentNode)
		currentNode = imgGraph[currentNode.X][currentNode.Y].Previous
		totalHeight += currentNode.Y - (currentNode.Y-oldY)*1000
		oldY = currentNode.Y
		count++
	}

	avgHeight := totalHeight / count

	visitedNodes := []Point{}
	if showVisited {
		for node, _ := range visited {
			visitedNodes = append(visitedNodes, node)
		}
	}

	return path, visitedNodes, avgHeight
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
