package controllers

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/B3zaleel/fractage/src/fractals"
	"github.com/B3zaleel/fractage/src/helpers"
	"github.com/kataras/iris/v12"
)

const (
	IFS_MAX_ITERATIONS           = 5_000_000_000
	IFS_DEFAULT_ITERATIONS       = 500_000
	IFS_DEFAULT_X                = 0
	IFS_DEFAULT_Y                = 0
	IFS_DEFAULT_SCALE            = 1
	IFS_MIN_SCALE                = 0
	IFS_MAX_SCALE                = 50_000
	IFS_DEFAULT_SYSTEM_VARIABLES = "0.0,0.0,0.0,0.16,0.0,0.0,0.01, 0.2,-0.26,0.23,0.22,0.0,1.6,0.07, -0.15,0.28,0.26,0.24,0.0,0.44,0.07, 0.85,0.04,-0.04,0.85,0.0,1.6,0.85"
	IFS_DEFAULT_SYSTEM_COLORS    = "mahogany, mahogany, mahogany, mahogany"
	IFS_DEFAULT_FOCUS            = true
)

func GetIFS(ctx iris.Context) {
	query := ctx.Request().URL.Query()
	fractal := fractals.IteratedFunctionSystem{
		Width:      DEFAULT_WIDTH,
		Height:     DEFAULT_HEIGHT,
		Iterations: IFS_DEFAULT_ITERATIONS,
		X:          IFS_DEFAULT_X,
		Y:          IFS_DEFAULT_Y,
		Focus:      IFS_DEFAULT_FOCUS,
		Scale:      IFS_DEFAULT_SCALE,
		Background: color.RGBA{255, 255, 255, 255},
	}
	ifsVariables := IFS_DEFAULT_SYSTEM_VARIABLES
	ifsColors := IFS_DEFAULT_SYSTEM_COLORS
	if query.Has("width") {
		width, err := strconv.Atoi(query.Get("width"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Width = width
	}
	if query.Has("height") {
		height, err := strconv.Atoi(query.Get("height"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Height = height
	}
	if query.Has("variables") {
		ifsVariables = query.Get("variables")
	}
	variables, err := fractals.GetIFSVariables(ifsVariables)
	if err != nil {
		ctx.Text(err.Error())
		return
	}
	fractal.Variables = variables
	if query.Has("color") {
		colorStr := fmt.Sprintf("\"%s\", ", query.Get("color"))
		ifsColors = strings.Repeat(colorStr, len(variables))
	}
	if query.Has("colors") {
		ifsColors = query.Get("colors")
	}
	fractal.Colors = fractals.GetIFSColors(ifsColors, len(variables))
	if query.Has("iterations") {
		iterations, err := strconv.Atoi(query.Get("iterations"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if iterations < 0 || iterations > IFS_MAX_ITERATIONS {
			ctx.Text(fmt.Sprintf("Iterations is too high. Max: %d\n", IFS_MAX_ITERATIONS))
			return
		}
		fractal.Iterations = iterations
	}
	if query.Has("x") {
		x, err := strconv.ParseFloat(query.Get("x"), 64)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.X = x
		fractal.Focus = false
	}
	if query.Has("y") {
		y, err := strconv.ParseFloat(query.Get("y"), 64)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Y = y
		fractal.Focus = false
	}
	if query.Has("scale") {
		scale, err := strconv.ParseFloat(query.Get("scale"), 64)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if scale < IFS_MIN_SCALE || scale > IFS_MAX_SCALE {
			ctx.Text(fmt.Sprintf("scale must be between %d and %d\n", IFS_MIN_SCALE, IFS_MAX_SCALE))
			return
		}
		fractal.Scale = scale
		fractal.Focus = false
	}
	if query.Has("focus") {
		focus, err := strconv.ParseBool(query.Get("focus"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Focus = focus
	}
	if query.Has("background") {
		background, err := helpers.ParseColor(query.Get("background"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Background = background
	}
	ctx.ContentType("image/png")
	fractal.WriteImage(ctx.ResponseWriter())
}
