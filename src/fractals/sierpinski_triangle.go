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

// Properties of a Sierpinski triangle image.
type SierpinskiTriangle struct {
	Width           int
	Height          int
	Color           color.RGBA
	UseRandomColors bool
	Iterations      int
	Background      color.RGBA
}

// Writes the Sierpinski triangle image to the given output.
func (props *SierpinskiTriangle) WriteImage(output io.Writer) {
	viewport := image.Rect(0, 0, props.Width, props.Height)
	img := image.NewRGBA(viewport)
	gc := draw2dimg.NewGraphicContext(img)
	var side, height float64
	if props.Width > props.Height {
		height = float64(props.Height)
		side = float64(props.Height) / (math.Sqrt(3) / 2.0)
	} else {
		side = float64(props.Width)
		height = math.Sqrt(3) / 2.0 * side
	}
	midX, midY := float64(props.Width/2), float64(props.Height/2)
	pt1 := helpers.Point{X: midX, Y: midY - height/2}
	pt2 := helpers.Point{X: midX + side/2, Y: midY + height/2}
	pt3 := helpers.Point{X: midX - side/2, Y: midY + height/2}
	helpers.FillImage(img, props.Background)
	props.render(gc, pt1, pt2, pt3, props.Iterations)
	err := png.Encode(output, img)
	if err != nil {
		panic(err)
	}
}

// Helper function for rendering the Sierpinski triangle.
func (props *SierpinskiTriangle) render(gc *draw2dimg.GraphicContext, pt1, pt2, pt3 helpers.Point, level int) {
	if level > 0 {
		pt1New := helpers.Point{X: (pt1.X + pt2.X) / 2, Y: (pt1.Y + pt2.Y) / 2}
		pt2New := helpers.Point{X: (pt2.X + pt3.X) / 2, Y: (pt2.Y + pt3.Y) / 2}
		pt3New := helpers.Point{X: (pt3.X + pt1.X) / 2, Y: (pt3.Y + pt1.Y) / 2}
		rectColor := props.Color
		if props.UseRandomColors {
			rectColor = helpers.RandomColor()
		}
		helpers.DrawTriangle(gc, pt1, pt2, pt3, rectColor)
		props.render(gc, pt1, pt1New, pt3New, level-1)
		props.render(gc, pt2, pt1New, pt2New, level-1)
		props.render(gc, pt3, pt2New, pt3New, level-1)
	}
}
