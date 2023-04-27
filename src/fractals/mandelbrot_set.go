package fractals

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/cmplx"

	"github.com/B3zaleel/fractage/src/helpers"
)

// Properties of a Mandelbrot set image.
type MandelbrotSet struct {
	Width         int
	Height        int
	ColorPalette  helpers.ColorPalette
	MaxIterations int
	M             float64
	BailOut       float64
	Region        helpers.Rect
	Background    color.RGBA
}

// Writes the Mandelbrot set image to the given output.
func (props *MandelbrotSet) WriteImage(output io.Writer) error {
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

// Helper function for rendering the Mandelbrot set.
func (props *MandelbrotSet) render(img *image.RGBA) error {
	width, height := float64(props.Width), float64(props.Height)
	bailOutPow := math.Pow(props.BailOut, props.M)
	step := math.Max(props.Region.Width/width, props.Region.Height/height)
	xOffset := props.Region.X - (width*step-props.Region.Width)/2.0
	yOffset := props.Region.Y - (height*step-props.Region.Height)/2.0
	err := props.ColorPalette.TranslateColorTransitions()
	if err != nil {
		return err
	}
	var pixelColor color.RGBA
	var x2, y2 float64
	var C, Z complex128
	var n int
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			n = 0
			C = complex(xOffset+float64(x)*step, yOffset+float64(y)*step)
			Z = complex(xOffset+float64(x)*step, yOffset+float64(y)*step)
			for n < props.MaxIterations {
				x2 = math.Pow(real(Z), props.M)
				y2 = math.Pow(imag(Z), props.M)
				if x2+y2 > bailOutPow {
					// Z diverges
					break
				}
				Z = cmplx.Pow(Z, complex(props.M, 0)) + C
				n++
			}
			if n < props.MaxIterations {
				// Z escaped
				pixelColor, err = props.ColorPalette.GetColor(float64(n) / float64(props.MaxIterations))
				if err != nil {
					return err
				}
			} else {
				pixelColor, err = props.ColorPalette.GetColor(1)
				if err != nil {
					return err
				}
			}
			img.Set(x, y, pixelColor)
		}
	}
	return nil
}
