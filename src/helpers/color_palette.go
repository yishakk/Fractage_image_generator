package helpers

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/llgcode/draw2d/draw2dimg"
	"gopkg.in/yaml.v3"
)

var (
	NIL_COLOR = color.RGBA{0, 0, 0, 0}
)

// Represents a color palette.
type ColorPalette struct {
	Name        string       `yaml:"name"`
	Transitions []Transition `yaml:"transitions"`
}

// Represents a color transition.
type Transition struct {
	Color    string `yaml:"color"`
	_Color   *color.RGBA
	Position float32 `yaml:"position"`
}

// Gets the color value of a given position in this ColorPalette.
func (palette *ColorPalette) GetColor(pos float64) (color.RGBA, error) {
	if palette.Transitions == nil || len(palette.Transitions) == 0 {
		return NIL_COLOR, errors.New("ColorPalette has no color transitions")
	}
	value := math.Max(0, math.Min(1.0, pos))
	idx := 0
	for (idx < len(palette.Transitions)) && (value >= float64(palette.Transitions[idx].Position)) {
		idx++
	}
	if idx > 0 {
		idx--
	}
	curTransition := palette.Transitions[idx]
	if idx >= len(palette.Transitions)-1 {
		if palette.Transitions[idx]._Color != nil {
			return *palette.Transitions[idx]._Color, nil
		}
		return ParseColor(palette.Transitions[idx].Color)
	}
	nextTransition := palette.Transitions[idx+1]
	var err error
	var curColor, nextColor color.RGBA
	if curTransition._Color != nil {
		curColor = *curTransition._Color
	} else {
		curColor, err = ParseColor(curTransition.Color)
		if err != nil {
			return NIL_COLOR, err
		}
	}
	if nextTransition._Color != nil {
		nextColor = *nextTransition._Color
	} else {
		nextColor, err = ParseColor(nextTransition.Color)
		if err != nil {
			return NIL_COLOR, err
		}
	}
	grad := (value - float64(curTransition.Position))
	grad /= (float64(nextTransition.Position) - float64(curTransition.Position))
	posColor := color.RGBA{
		R: curColor.R + uint8(grad*float64(nextColor.R-curColor.R)),
		G: curColor.G + uint8(grad*float64(nextColor.G-curColor.G)),
		B: curColor.B + uint8(grad*float64(nextColor.B-curColor.B)),
		A: 255,
	}
	return posColor, nil
}

// Translate the value of the color transitions for this color palette.
func (palette *ColorPalette) TranslateColorTransitions() error {
	for i := 0; i < len(palette.Transitions); i++ {
		transitionColor, err := ParseColor(palette.Transitions[i].Color)
		if err != nil {
			return err
		}
		palette.Transitions[i]._Color = &transitionColor
	}
	return nil
}

// Displays an image of this palette to the given output.
func (palette *ColorPalette) Render(output io.Writer, width, height int, step float64) error {
	viewport := image.Rect(0, 0, width, height)
	img := image.NewRGBA(viewport)
	gc := draw2dimg.NewGraphicContext(img)
	err := palette.TranslateColorTransitions()
	if err != nil {
		return err
	}
	if step <= 0 {
		return errors.New("steps between transitions must be greater than 0")
	}
	for x := 0.0; x <= float64(width); x += step {
		pos := x / float64(width)
		curColor, err := palette.GetColor(pos)
		if err != nil {
			return err
		}
		FillRectangle(gc, x, 0.0, step, float64(height), curColor)
	}
	err = png.Encode(output, img)
	if err != nil {
		return err
	}
	return nil
}

// Parses a string into a color palette.
func ParseColorPalette(txt string) (ColorPalette, error) {
	text := strings.Trim(txt, WHITESPACE_CUTSET)
	palette, err := ParseNameColorPalette(text)
	if err == nil {
		return palette, err
	}
	nilPalette := ColorPalette{
		Name:        "",
		Transitions: nil,
	}
	values, err := GetCSV(text)
	if len(values)%2 != 0 {
		return nilPalette, errors.New("Invalid color palette")
	}
	transitions := make([]Transition, len(values)/2)
	for i, j := 0, 0; i < len(values); i += 2 {
		_, err := ParseColor(values[i])
		if err != nil {
			return nilPalette, err
		}
		position, err := strconv.ParseFloat(values[i+1], 32)
		if err != nil {
			return nilPalette, err
		}
		if j == 0 && position != 0.0 {
			return nilPalette, errors.New("The first position must be 0")
		}
		if j == len(values)/2-1 && position != 1.0 {
			return nilPalette, errors.New("The last position must be 1")
		}
		transitions[j] = Transition{
			Color:    values[i],
			Position: float32(position),
		}
		j++
	}
	return ColorPalette{Name: "custom_palette", Transitions: transitions}, nil
}

// Returns the color value of a predetermined color palette that
// matches the given name.
func ParseNameColorPalette(name string) (ColorPalette, error) {
	var palettes []ColorPalette
	nilPalette := ColorPalette{
		Name:        "",
		Transitions: nil,
	}
	file, err := os.ReadFile("src/data/color_palettes.yaml")
	if err != nil {
		return nilPalette, err
	}
	err = yaml.Unmarshal(file, &palettes)
	if err != nil {
		return nilPalette, err
	}
	for _, palette := range palettes {
		if palette.Name == name {
			return palette, nil
		}
	}
	return nilPalette, errors.New("Palette not found")
}
