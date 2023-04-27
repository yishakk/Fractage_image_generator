package fractals

import (
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/B3zaleel/fractage/src/helpers"
	"github.com/llgcode/draw2d/draw2dimg"
)

// Properties of a Cantor set image.
type CantorSet struct {
	Width           int
	Height          int
	Color           color.RGBA
	UseRandomColors bool
	LineHeight      float64
	Iterations      int
	Background      color.RGBA
}

// Writes the Cantor set image to the given output.
func (props *CantorSet) WriteImage(output io.Writer) {
	viewport := image.Rect(0, 0, props.Width, props.Height)
	img := image.NewRGBA(viewport)
	gc := draw2dimg.NewGraphicContext(img)
	x := float64(viewport.Min.X)
	y := float64(props.Height)/2 - float64(props.Iterations)*props.LineHeight + props.LineHeight/2
	helpers.FillImage(img, props.Background)
	props.render(gc, x, y, float64(viewport.Dx()), props.Iterations)
	err := png.Encode(output, img)
	if err != nil {
		panic(err)
	}
}

// Helper function for rendering the Cantor set.
func (props *CantorSet) render(gc *draw2dimg.GraphicContext, x, y, width float64, level int) {
	if level > 0 {
		dx := width / 3
		rectColor := props.Color
		if props.UseRandomColors {
			rectColor = helpers.RandomColor()
		}
		helpers.FillRectangle(gc, x, y, width, props.LineHeight, rectColor)
		props.render(gc, x, y+props.LineHeight*2, dx, level-1)
		props.render(gc, x+2*dx, y+props.LineHeight*2, dx, level-1)
	}
}
