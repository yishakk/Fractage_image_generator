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
	MANDELBROT_SET_MAX_ITERATIONS        = 500_000
	MANDELBROT_SET_DEFAULT_ITERATIONS    = 700
	MANDELBROT_SET_DEFAULT_COLOR_PALETTE = "orange_blue"
	MANDELBROT_SET_DEFAULT_BAIL_OUT      = 20
	MANDELBROT_SET_DEFAULT_M             = 2
	MANDELBROT_SET_DEFAULT_REGION        = "-2, -1.25, 3.25, 2.5"
)

func GetMandelbrotSet(ctx iris.Context) {
	query := ctx.Request().URL.Query()
	fractal := fractals.MandelbrotSet{
		Width:         DEFAULT_WIDTH,
		Height:        DEFAULT_HEIGHT,
		MaxIterations: MANDELBROT_SET_DEFAULT_ITERATIONS,
		M:             MANDELBROT_SET_DEFAULT_M,
		BailOut:       MANDELBROT_SET_DEFAULT_BAIL_OUT,
		Background:    color.RGBA{255, 255, 255, 255},
	}
	colorPaletteValue := MANDELBROT_SET_DEFAULT_COLOR_PALETTE
	regionValue := MANDELBROT_SET_DEFAULT_REGION
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
	if query.Has("color_palette") {
		colorPaletteValue = query.Get("color_palette")
	}
	if query.Has("iterations") {
		iterations, err := strconv.Atoi(query.Get("iterations"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if iterations < 0 || iterations > MANDELBROT_SET_MAX_ITERATIONS {
			ctx.Text(fmt.Sprintf("Too many iterations. Max: %d\n", MANDELBROT_SET_MAX_ITERATIONS))
			return
		}
		fractal.MaxIterations = iterations
	}
	if query.Has("m") {
		m, err := strconv.ParseFloat(query.Get("m"), 64)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.M = m
	}
	if query.Has("region") {
		regionValue = query.Get("region")
	}
	if query.Has("bail_out") {
		bailOut, err := strconv.ParseFloat(query.Get("bail_out"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.BailOut = bailOut
	}
	if query.Has("background") {
		background, err := helpers.ParseColor(query.Get("background"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Background = background
	}
	region, err := helpers.ParseRect(regionValue)
	if err != nil {
		ctx.Text(err.Error())
		return
	}
	colorPalette, err := helpers.ParseColorPalette(colorPaletteValue)
	if err != nil {
		ctx.Text(err.Error())
		return
	}
	fractal.Region = region
	fractal.ColorPalette = colorPalette
	ctx.ContentType("image/png")
	err = fractal.WriteImage(ctx.ResponseWriter())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
