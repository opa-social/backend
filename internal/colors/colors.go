package colors

import (
	"fmt"
	"image/color"

	"github.com/muesli/gamut"
)

func colorToString(c color.Color) string {
	r, g, b, _ := c.RGBA()

	// Shift color space.
	r = r >> 8
	g = g >> 8
	b = b >> 8

	return fmt.Sprintf("%02x%02x%02x", r, g, b)
}

// GenerateRandomColors creates a set of n random pastelle colors and converts them
// to their hex string represenatation.
func GenerateRandomColors(n int) []string {
	colors := make([]string, n)

	palette, _ := gamut.Generate(n, gamut.PastelGenerator{})
	for _, c := range palette {
		colors = append(colors, colorToString(c))
	}

	return colors
}
