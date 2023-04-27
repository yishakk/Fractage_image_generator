package fractals

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/cmplx"

	"github.com/B3zaleel/fractage/src/helpers"
	math_helpers "github.com/B3zaleel/fractage/src/helpers/math"
)

const (
	MAX_DELTA = 1e-14
)

// Properties of a Newton basin image.
type NewtonBasin struct {
	Width            int
	Height           int
	ColorPalette     helpers.ColorPalette
	MaxIterations    int
	Polynomial       math_helpers.CmplxPolynomial
	BailOut          float64
	Region           helpers.Rect
	Background       color.RGBA
	UseDynamicColors bool
}

// Writes the Newton basin image to the given output.
func (props *NewtonBasin) WriteImage(output io.Writer) error {
	viewport := image.Rect(0, 0, props.Width, props.Height)
	img := image.NewRGBA(viewport)
	helpers.FillImage(img, props.Background)
	err := props.render(img)
	if err != nil {
		return err
	}
	err = png.Encode(output, img)
	if err != nil {
		return err
	}
	return nil
}

// Helper function for rendering the Newton basin.
func (props *NewtonBasin) render(img *image.RGBA) error {
	width, height := float64(props.Width), float64(props.Height)
	step := math.Max(props.Region.Width/width, props.Region.Height/height)
	xOffset := props.Region.X - (width*step-props.Region.Width)/2.0
	yOffset := props.Region.Y - (height*step-props.Region.Height)/2.0
	err := props.ColorPalette.TranslateColorTransitions()
	if err != nil {
		return err
	}
	var pixelColor color.RGBA
	var n int
	poly := props.Polynomial
	polyDeriv := props.Polynomial.FirstDerivative()
	for y := 0; y <= int(height); y++ {
		for x := 0; x <= int(width); x++ {
			n = 0
			Z := complex(xOffset+float64(x)*step, yOffset+float64(y)*step)
			delta := Z
			Z1 := Z
			for (n < props.MaxIterations) && (cmplx.Abs(Z) < props.BailOut) && (cmplx.Abs(delta) > MAX_DELTA) {
				Z = Z - poly.Evaluate(Z)/polyDeriv.Evaluate(Z)
				delta = Z1 - Z
				Z1 = Z
				n++
			}
			mag := float64(props.MaxIterations-n) / float64(props.MaxIterations)
			if props.UseDynamicColors {
				var angle float64
				if Z == 0+0i {
					angle = 0
				} else {
					angle = cmplx.Phase(Z)
				}
				pixelColor = color.RGBA{
					R: uint8(255 * mag * (math.Sin(angle)/2 + 0.5)),
					G: uint8(255 * mag * (math.Sin(angle+1*math.Pi/3)/2 + 0.5)),
					B: uint8(255 * mag * (math.Sin(angle+5*math.Pi/3)/2 + 0.5)),
					A: 255,
				}
			} else {
				pixelColor, err = props.ColorPalette.GetColor(mag)
			}
			img.Set(x, y, pixelColor)
		}
	}
	return nil
}
