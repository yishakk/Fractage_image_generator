package fractals

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/rand"
	"strconv"

	"github.com/B3zaleel/fractage/src/helpers"
)

const (
	IFS_FXN_INDEX_A           = 0
	IFS_FXN_INDEX_B           = 1
	IFS_FXN_INDEX_C           = 2
	IFS_FXN_INDEX_D           = 3
	IFS_FXN_INDEX_E           = 4
	IFS_FXN_INDEX_F           = 5
	IFS_FXN_INDEX_PROBABILITY = 6
	// The number of variables in each set of the iterated function system.
	IFS_FXN_VARIABLES_COUNT = int(7)
)

// Properties of an iterated function system (IFS) image.
type IteratedFunctionSystem struct {
	Width      int
	Height     int
	Colors     []color.RGBA
	Iterations int
	X          float64
	Y          float64
	Scale      float64
	Focus      bool
	Variables  [][IFS_FXN_VARIABLES_COUNT]float64
	Background color.RGBA
}

// Writes the IFS image to the given output.
func (props *IteratedFunctionSystem) WriteImage(output io.Writer) error {
	viewport := image.Rect(0, 0, props.Width, props.Height)
	img := image.NewRGBA(viewport)
	helpers.FillImage(img, props.Background)
	props.render(img)
	err := png.Encode(output, img)
	if err != nil {
		return err
	}
	return nil
}

// Helper function for rendering the IFS.
func (props *IteratedFunctionSystem) render(img *image.RGBA) {
	xMin, yMin, xMax, yMax := 0.0, 0.0, 0.0, 0.0
	var x, y, xn float64
	var ptColor color.RGBA
	for round := 1; round < 3; round++ {
		x, y, xn = 0, 0, 0
		if props.Focus {
			props.Scale = 1
			if round == 2 {
				boundsWidth, boundsHeight := xMax-xMin, yMax-yMin
				xScale, yScale := float64(props.Width)/boundsWidth, float64(props.Height)/boundsHeight
				props.Scale = math.Min(xScale, yScale)
				props.X = -xMin*props.Scale + (float64(props.Width)-boundsWidth*props.Scale)/2
				props.Y = -yMin*props.Scale + (float64(props.Height)-boundsHeight*props.Scale)/2
			}
		}
		if !props.Focus || (props.Focus && round == 2) {
			rand.Seed(0)
		}
		for i := 0; i < props.Iterations; i++ {
			probability := rand.Float64()
			xn = x
			sum := float64(0.0)
			ptColor = props.Colors[0]
			for j := 0; j < len(props.Variables); j++ {
				fxn := props.Variables[j]
				sum += fxn[IFS_FXN_INDEX_PROBABILITY]
				if probability <= sum {
					a := fxn[IFS_FXN_INDEX_A]
					b := fxn[IFS_FXN_INDEX_B]
					c := fxn[IFS_FXN_INDEX_C]
					d := fxn[IFS_FXN_INDEX_D]
					e := fxn[IFS_FXN_INDEX_E]
					f := fxn[IFS_FXN_INDEX_F]
					x = a*xn + b*y + e
					y = c*xn + d*y + f
					ptColor = props.Colors[j]
					if props.Focus && round == 1 {
						xMin = math.Min(xMin, x)
						yMin = math.Min(yMin, y)
						xMax = math.Max(xMax, x)
						yMax = math.Max(yMax, y)
					}
					break
				}
			}
			if !props.Focus || (props.Focus && round == 2) {
				ptX := int(math.Round(props.X + x*props.Scale))
				ptY := int(math.Round(float64(props.Height) - (props.Y + y*props.Scale)))
				img.Set(ptX, ptY, ptColor)
			}
		}
		if !props.Focus {
			break
		}
	}
}

// Retrieves a comma-separated list of the variables for each set of the IFS.
func GetIFSVariables(txt string) ([][IFS_FXN_VARIABLES_COUNT]float64, error) {
	values, err := helpers.GetCSV(txt)
	if err != nil {
		return nil, err
	}
	if len(values)%IFS_FXN_VARIABLES_COUNT != 0 {
		return nil, errors.New("Incomplete IFS variables provided.")
	}
	count := len(values) / IFS_FXN_VARIABLES_COUNT
	functions := make([][IFS_FXN_VARIABLES_COUNT]float64, count)
	for i, j := 0, 0; i < len(values); i += IFS_FXN_VARIABLES_COUNT {
		a, err := strconv.ParseFloat(values[i+IFS_FXN_INDEX_A], 64)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseFloat(values[i+IFS_FXN_INDEX_B], 64)
		if err != nil {
			return nil, err
		}
		c, err := strconv.ParseFloat(values[i+IFS_FXN_INDEX_C], 64)
		if err != nil {
			return nil, err
		}
		d, err := strconv.ParseFloat(values[i+IFS_FXN_INDEX_D], 64)
		if err != nil {
			return nil, err
		}
		e, err := strconv.ParseFloat(values[i+IFS_FXN_INDEX_E], 64)
		if err != nil {
			return nil, err
		}
		f, err := strconv.ParseFloat(values[i+IFS_FXN_INDEX_F], 64)
		if err != nil {
			return nil, err
		}
		probability, err := strconv.ParseFloat(values[i+IFS_FXN_INDEX_PROBABILITY], 64)
		if err != nil {
			return nil, err
		}
		functions[j][IFS_FXN_INDEX_A] = a
		functions[j][IFS_FXN_INDEX_B] = b
		functions[j][IFS_FXN_INDEX_C] = c
		functions[j][IFS_FXN_INDEX_D] = d
		functions[j][IFS_FXN_INDEX_E] = e
		functions[j][IFS_FXN_INDEX_F] = f
		functions[j][IFS_FXN_INDEX_PROBABILITY] = probability
		j++
	}
	return functions, nil
}

// Retrieves a slice of colors for the IFS.
func GetIFSColors(txt string, count int) []color.RGBA {
	values, err := helpers.GetCSV(txt)
	colors := make([]color.RGBA, count)
	for i := 0; i < count; i++ {
		if err != nil || i >= len(values) {
			colors[i] = helpers.RandomColor()
		} else {
			color, err := helpers.ParseColor(values[i])
			if err != nil {
				color = helpers.RandomColor()
			}
			colors[i] = color
		}
	}
	return colors
}
