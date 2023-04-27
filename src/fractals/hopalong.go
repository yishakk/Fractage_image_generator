package fractals

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"

	"github.com/B3zaleel/fractage/src/helpers"
)

var (
	HOPALONG_TYPES = map[string]func(props *Hopalong, xIn, yIn float64) (xOut, yOut float64){
		"classic_bm":      classic_barry_martin_fractal,
		"positive_bm":     positive_barry_martin_fractal,
		"additive_bm":     additive_barry_martin_fractal,
		"gingerbread_man": gingerbread_man_fractal,
	}
)

// Properties of a Hopalong image.
type Hopalong struct {
	Width           int
	Height          int
	Color           color.RGBA
	UseRandomColors bool
	A               float64
	B               float64
	C               float64
	D               float64
	X               float64
	Y               float64
	Type            string
	Scale           float64
	Resolution      int
	Background      color.RGBA
}

// Writes the Hopalong image to the given output.
func (props *Hopalong) WriteImage(output io.Writer) {
	viewport := image.Rect(0, 0, props.Width, props.Height)
	img := image.NewRGBA(viewport)
	helpers.FillImage(img, props.Background)
	props.render(img)
	err := png.Encode(output, img)
	if err != nil {
		panic(err)
	}
}

// Helper function for rendering the Hopalong.
func (props *Hopalong) render(img *image.RGBA) {
	x, y := props.X, props.Y
	midX, midY := float64(props.Width)/2.0, float64(props.Height)/2.0
	ptColor := props.Color
	if props.UseRandomColors {
		ptColor = helpers.RandomColor()
	}
	hopalong_fxn := HOPALONG_TYPES[props.Type]
	for i := 0; i < props.Width; i++ {
		for j := 0; j < props.Height; j++ {
			for k := 0; k < props.Resolution; k++ {
				x, y = hopalong_fxn(props, x, y)
				if props.UseRandomColors && i%50 == 0 {
					ptColor = helpers.RandomColor()
				}
				img.Set(int(midX+x*props.Scale), int(midY-y*props.Scale), ptColor)
			}
		}
	}
}

func classic_barry_martin_fractal(props *Hopalong, xIn, yIn float64) (xOut, yOut float64) {
	xSign := 0
	if xIn < 0 {
		xSign = -1
	} else if xIn > 0 {
		xSign = 1
	}
	xOut = yIn - float64(xSign)*math.Sqrt(math.Abs(float64(props.B)*xIn-float64(props.C)))
	yOut = float64(props.A) - xIn
	return
}

func positive_barry_martin_fractal(props *Hopalong, xIn, yIn float64) (xOut, yOut float64) {
	xSign := 0
	if xIn < 0 {
		xSign = -1
	} else if xIn > 0 {
		xSign = 1
	}
	xOut = yIn + float64(xSign)*math.Sqrt(math.Abs(float64(props.B)*xIn-float64(props.C)))
	yOut = float64(props.A) - xIn
	return
}

func additive_barry_martin_fractal(props *Hopalong, xIn, yIn float64) (xOut, yOut float64) {
	xOut = yIn + math.Sqrt(math.Abs(float64(props.B)*xIn-float64(props.C)))
	yOut = float64(props.A) - xIn
	return
}

func gingerbread_man_fractal(props *Hopalong, xIn, yIn float64) (xOut, yOut float64) {
	xOut = yIn + math.Abs(float64(props.B)*xIn)
	yOut = float64(props.A) - xIn
	return
}
