package fractals

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"strings"

	"github.com/B3zaleel/fractage/src/helpers"
	"github.com/llgcode/draw2d/draw2dimg"
)

type LindenmayerSystem struct {
	Width                 int
	Height                int
	Axiom                 string
	Iterations            int
	RewriteRules          map[rune]string
	DrawSymbols           string
	SkipSymbols           string
	Angle                 float64
	Color                 color.RGBA
	UseRandomColors       bool
	Focus                 bool
	TurningAngle          float64
	Position              string
	LineWidth             float64
	LineLength            float64
	LineWidthIncrement    float64
	LineLengthScaleFactor float64
	TurningAngleIncrement float64
	Background            color.RGBA
}

// Represents a drawing state
type State struct {
	Angle              float64
	LineWidth          float64
	LineLength         float64
	TurningAngle       float64
	SwapTurnDirections bool
	X                  float64
	Y                  float64
}

// Writes the Lindenmayer system image to the given output.
func (props *LindenmayerSystem) WriteImage(output io.Writer) error {
	var x, y float64
	generator := props.BuildGenerator()
	x, y, _ = ParseLSystemPosition(props.Position, float64(props.Width), float64(props.Height))
	viewport := image.Rect(0, 0, props.Width, props.Height)
	img := image.NewRGBA(viewport)
	gc := draw2dimg.NewGraphicContext(img)
	helpers.FillImage(img, props.Background)
	props.render(gc, &generator, x, y)
	err := png.Encode(output, img)
	if err != nil {
		return err
	}
	return nil
}

func (props *LindenmayerSystem) render(gc *draw2dimg.GraphicContext, generator *[]rune, startX, startY float64) {
	drawingStates := make([]State, 1)
	color := props.Color
	if props.UseRandomColors && gc != nil {
		color = helpers.RandomColor()
	}
	n := len(drawingStates)
	bounds := helpers.Rect{
		X:      startX,
		Y:      startY,
		Width:  startX,
		Height: startY,
	}
	for round := 1; round < 3; round++ {
		i := 0
		polygonOpen := false
		if props.Focus && round == 2 {
			// adjust the drawing parameters
			scaleX := float64(props.Width) / bounds.Width
			scaleY := float64(props.Height) / bounds.Height
			step := math.Min(float64(scaleX), float64(scaleY))
			props.LineLength = step
		}
		drawingStates[0] = State{
			Angle:              props.Angle,
			LineWidth:          props.LineWidth,
			LineLength:         props.LineLength,
			TurningAngle:       props.TurningAngle,
			SwapTurnDirections: false,
			X:                  startX,
			Y:                  startY,
		}
		for _, c := range *generator {
			switch c {
			case 'F', 'f':
				{
					x0, y0 := drawingStates[i].X, drawingStates[i].Y
					x1 := x0 + drawingStates[i].LineLength*math.Cos(drawingStates[i].Angle*math.Pi/180)
					y1 := y0 + drawingStates[i].LineLength*math.Sin(drawingStates[i].Angle*math.Pi/180)
					if !props.Focus || (props.Focus && round == 2) {
						if !polygonOpen {
							gc.SetLineWidth(drawingStates[i].LineWidth)
							gc.SetStrokeColor(color)
							gc.BeginPath()
						}
						canDraw := true
						if canDraw {
							gc.MoveTo(x0, y0)
							if c == 'F' {
								gc.LineTo(x1, y1)
							} else {
								gc.MoveTo(x1, y1)
							}
						}
						if !polygonOpen {
							gc.Close()
							gc.Stroke()
						}
					}
					if x1 > bounds.Width {
						bounds.Width = x1
					}
					if x1 < bounds.X {
						bounds.X = x1
					}
					if y1 > bounds.Height {
						bounds.Height = y1
					}
					if y1 < bounds.Y {
						bounds.Y = y1
					}
					drawingStates[i].X = x1
					drawingStates[i].Y = y1
				}
			case '+':
				{
					if drawingStates[i].SwapTurnDirections {
						drawingStates[i].Angle += drawingStates[i].TurningAngle
					} else {
						drawingStates[i].Angle -= drawingStates[i].TurningAngle
					}
				}
			case '-':
				{
					if drawingStates[i].SwapTurnDirections {
						drawingStates[i].Angle -= drawingStates[i].TurningAngle
					} else {
						drawingStates[i].Angle += drawingStates[i].TurningAngle
					}
				}
			case '|':
				{
					if drawingStates[i].Angle >= 180 {
						drawingStates[i].Angle -= 180.0
					} else {
						drawingStates[i].Angle += 180.0
					}
				}
			case '[':
				{
					newDrawingState := State{
						Angle:              drawingStates[i].Angle,
						LineWidth:          drawingStates[i].LineWidth,
						LineLength:         drawingStates[i].LineLength,
						TurningAngle:       drawingStates[i].TurningAngle,
						SwapTurnDirections: drawingStates[i].SwapTurnDirections,
						X:                  drawingStates[i].X,
						Y:                  drawingStates[i].Y,
					}
					if i+1 < n {
						drawingStates[i+1] = newDrawingState
					} else {
						drawingStates = append(drawingStates, newDrawingState)
						n++
					}
					i++
				}
			case ']':
				{
					if i > 0 {
						i--
					}
				}
			case '#':
				drawingStates[i].LineWidth += props.LineWidthIncrement
			case '!':
				drawingStates[i].LineWidth -= props.LineWidthIncrement
			case '@':
				{
					if (!props.Focus || (props.Focus && round == 2)) && !polygonOpen {
						// TODO: Draw a dot
					}
				}
			case '{':
				{
					if !polygonOpen {
						polygonOpen = true
						if !props.Focus || (props.Focus && round == 2) {
							gc.SetFillColor(color)
							gc.SetStrokeColor(color)
							gc.SetLineWidth(drawingStates[i].LineWidth)
							gc.BeginPath()
						}
					}
				}
			case '}':
				{
					if polygonOpen {
						polygonOpen = false
						if (!props.Focus && round == 1) || (props.Focus && round == 2) {
							gc.Close()
							gc.FillStroke()
						}
					}
				}
			case '>':
				drawingStates[i].LineLength *= props.LineLengthScaleFactor
			case '<':
				drawingStates[i].LineLength /= props.LineLengthScaleFactor
			case '&':
				drawingStates[i].SwapTurnDirections = !drawingStates[i].SwapTurnDirections
			case '(':
				drawingStates[i].TurningAngle -= props.TurningAngleIncrement
			case ')':
				drawingStates[i].TurningAngle += props.TurningAngleIncrement
			}
		}
		if !props.Focus || (props.Focus && round == 2) && polygonOpen {
			gc.Close()
			gc.FillStroke()
		}
		if !props.Focus {
			break
		}
	}
}

// Builds the image generation string for this Lindenmayer system.
func (props *LindenmayerSystem) BuildGenerator() []rune {
	previousString := []rune(props.Axiom)
	for i := 0; i < props.Iterations; i++ {
		var newString []rune
	str_replacement:
		for _, c := range previousString {
			for variable, replacement := range props.RewriteRules {
				if variable == c {
					newString = append(newString, []rune(replacement)...)
					continue str_replacement
				}
			}
			newString = append(newString, c)
		}
		previousString = newString
	}
	for i := 0; i < len(previousString); i++ {
		c := previousString[i]
		if strings.ContainsRune(props.DrawSymbols, c) {
			previousString[i] = 'F'
		} else if strings.ContainsRune(props.SkipSymbols, c) {
			previousString[i] = 'f'
		}
	}
	return previousString
}

// Converts a comma-separated list of rewrite rules to a map of
// variable and rewrite string.
func ParseLindenmayerRules(txt string) (map[rune]string, error) {
	values, err := helpers.GetCSV(txt)
	if err != nil {
		return nil, err
	}
	rules := make(map[rune]string, len(values))
	for i := 0; i < len(values); i++ {
		rule := strings.Trim(values[i], helpers.WHITESPACE_CUTSET)
		if strings.ContainsRune(rule, '=') {
			before, expression, _ := strings.Cut(rule, "=")
			variable := []rune(strings.Trim(before, helpers.WHITESPACE_CUTSET))
			if len(variable) == 1 {
				rules[variable[0]] = expression
			} else {
				return nil, errors.New("The variable must be a single character")
			}
		} else {
			return nil, errors.New("Invalid rewrite rule")
		}
	}
	return rules, nil
}

// Retrieves the cartesian position of a point given a position.
func ParseLSystemPosition(txt string, width, height float64) (x, y float64, err error) {
	switch txt {
	case "top-left", "left-top":
		return 0, 0, nil
	case "top-center", "center-top", "top":
		return width / 2.0, 0, nil
	case "top-right", "right-top":
		return width, 0, nil
	case "center-left", "left-center", "left":
		return 0, height / 2.0, nil
	case "center-right", "right-center", "right":
		return width, height / 2.0, nil
	case "center-center", "center":
		return width / 2.0, height / 2.0, nil
	case "bottom-left", "left-bottom":
		return 0, height, nil
	case "bottom-center", "center-bottom", "bottom":
		return width / 2.0, height, nil
	case "bottom-right", "right-bottom":
		return width, height, nil
	}
	return 0, 0, errors.New(fmt.Sprintf("Unknown position: %s\n", txt))
}
