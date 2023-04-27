package fractals

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"

	"github.com/B3zaleel/fractage/src/helpers"
	"github.com/llgcode/draw2d/draw2dimg"
)

// Properties of a Sierpinski carpet image.
type SierpinskiCarpet struct {
	Width           int
	Height          int
	Color           color.RGBA
	UseRandomColors bool
	Iterations      int
	Background      color.RGBA
}

// Writes the Sierpinski carpet image to the given output.
func (props *SierpinskiCarpet) WriteImage(output io.Writer) {
	viewport := image.Rect(0, 0, props.Width, props.Height)
	img := image.NewRGBA(viewport)
	gc := draw2dimg.NewGraphicContext(img)
	minSide := math.Min(float64(props.Width), float64(props.Height))
	x1 := 0 + float64(props.Width)/2 - minSide/2
	x2 := x1 + minSide
	y1 := 0 + float64(props.Height)/2 - minSide/2
	y2 := y1 + minSide
	helpers.FillImage(img, props.Background)
	helpers.DrawRectangle(gc, x1, y1, x2-x1, y2-y1, color.RGBA{0, 0, 0, 255})
	props.render(gc, x1, y1, x2, y2, props.Iterations)
	err := png.Encode(output, img)
	if err != nil {
		panic(err)
	}
}

// Helper function for rendering the Sierpinski carpet.
func (props *SierpinskiCarpet) render(gc *draw2dimg.GraphicContext, x1, y1, x2, y2 float64, level int) {
	if level > 0 {
		x1n := 2*x1/3 + x2/3
		x2n := x1/3 + 2*x2/3
		y1n := 2*y1/3 + y2/3
		y2n := y1/3 + 2*y2/3
		rectColor := props.Color
		if props.UseRandomColors {
			rectColor = helpers.RandomColor()
		}
		helpers.FillRectangle(gc, x1n, y1n, x2n-x1n, y2n-y1n, rectColor)

		props.render(gc, x1, y1, x1n, y1n, level-1)
		props.render(gc, x1n, y1, x2n, y1n, level-1)
		props.render(gc, x2n, y1, x2, y1n, level-1)
		props.render(gc, x1, y1n, x1n, y2n, level-1)
		props.render(gc, x2n, y1n, x2, y2n, level-1)
		props.render(gc, x1, y2n, x1n, y2, level-1)
		props.render(gc, x1n, y2n, x2n, y2, level-1)
		props.render(gc, x2n, y2n, x2, y2, level-1)
	}
}
