package seamcarver

import (
	//"fmt"
	//"image"
	"image/color"
	//"image/jpeg"
)

func Luminance(color color.Color) float64 {
	r, g, b, _ := color.RGBA()
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}
