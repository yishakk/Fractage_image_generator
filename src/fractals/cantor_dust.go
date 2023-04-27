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

// Properties of a Cantor set image.
type CantorDust struct {
	Width           int
	Height          int
	Color           color.RGBA
	UseRandomColors bool
	Iterations      int
	Background      color.RGBA
}

// Writes the Cantor dust image to the given output.
func (props *CantorDust) WriteImage(output io.Writer) {
	viewport := image.Rect(0, 0, props.Width, props.Height)
	img := image.NewRGBA(viewport)
	gc := draw2dimg.NewGraphicContext(img)
	length := math.Min(float64(props.Width), float64(props.Height))
	x := float64(viewport.Min.X) + float64(props.Width)/2 - length/2
	y := float64(viewport.Min.Y) + float64(props.Height)/2 - length/2
	helpers.FillImage(img, props.Background)
	props.render(gc, x, y, length, length, props.Iterations)
	err := png.Encode(output, img)
	if err != nil {
		panic(err)
	}
}

// Helper function for rendering the Cantor dust.
func (props *CantorDust) render(gc *draw2dimg.GraphicContext, x, y, width, height float64, level int) {
	if level > 0 {
		dx, dy := width/3, height/3
		props.render(gc, x, y, dx, dy, level-1)
		props.render(gc, x+2*dx, y, dx, dy, level-1)
		props.render(gc, x, y+2*dy, dx, dy, level-1)
		props.render(gc, x+2*dx, y+2*dy, dx, dy, level-1)
	} else {
		rectColor := props.Color
		if props.UseRandomColors {
			rectColor = helpers.RandomColor()
		}
		helpers.FillRectangle(gc, x, y, width, height, rectColor)
	}
}
