package seamcarver

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
)

func Luminance(color color.Color) float64 {
	r, g, b, _ := color.RGBA()
	return 0.2126*r + 0.7152*g + 0.0722*b
}
