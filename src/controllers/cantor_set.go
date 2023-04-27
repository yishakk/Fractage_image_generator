package controllers

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/B3zaleel/fractage/src/fractals"
	"github.com/B3zaleel/fractage/src/helpers"
	"github.com/kataras/iris/v12"
)

const (
	MAX_LINE_HEIGHT     = 30
	DEFAULT_LINE_HEIGHT = 5
)

func GetCantorSet(ctx iris.Context) {
	query := ctx.Request().URL.Query()
	fractal := fractals.CantorSet{
		Width:           DEFAULT_WIDTH,
		Height:          DEFAULT_HEIGHT,
		UseRandomColors: true,
		Iterations:      DEFAULT_ITERATIONS,
		LineHeight:      DEFAULT_LINE_HEIGHT,
		Background:      color.RGBA{255, 255, 255, 255},
	}
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
	if query.Has("color") {
		color, err := helpers.ParseColor(query.Get("color"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Color = color
		fractal.UseRandomColors = false
	}
	if query.Has("iterations") {
		iterations, err := strconv.Atoi(query.Get("iterations"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if iterations < 0 || iterations > MAX_ITERATIONS {
			ctx.Text(fmt.Sprintf("Too many iterations. Max: %d\n", MAX_ITERATIONS))
			return
		}
		fractal.Iterations = iterations
	}
	if query.Has("line_height") {
		lineHeight, err := strconv.ParseFloat(query.Get("line_height"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if lineHeight < 0 || lineHeight > MAX_ITERATIONS {
			ctx.Text(fmt.Sprintf("Height is too large. Max: %d\n", MAX_LINE_HEIGHT))
			return
		}
		fractal.LineHeight = lineHeight
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
