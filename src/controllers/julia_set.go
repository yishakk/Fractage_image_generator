package controllers

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/B3zaleel/fractage/src/fractals"
	"github.com/B3zaleel/fractage/src/helpers"
	"github.com/kataras/iris/v12"
)

func GetJuliaSet(ctx iris.Context) {
	query := ctx.Request().URL.Query()
	fractal := fractals.JuliaSet{
		Width:         DEFAULT_WIDTH,
		Height:        DEFAULT_HEIGHT,
		C:             fractals.JULIA_SET_DEFAULT_C,
		MaxIterations: fractals.JULIA_SET_DEFAULT_ITERATIONS,
		BailOut:       fractals.JULIA_SET_DEFAULT_BAIL_OUT,
		Background:    color.RGBA{255, 255, 255, 255},
	}
	colorPaletteValue := fractals.JULIA_SET_DEFAULT_COLOR_PALETTE
	regionValue := fractals.JULIA_SET_DEFAULT_REGION
	seriesName := fractals.JULIA_SET_DEFAULT_SERIES_TYPE
	variablesTxt := fractals.JULIA_SET_DEFAULT_VARIABLES_TEXT
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
	if query.Has("c") {
		c, err := strconv.ParseComplex(query.Get("c"), 64)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.C = c
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
		if iterations < 0 || iterations > fractals.JULIA_SET_MAX_ITERATIONS {
			ctx.Text(fmt.Sprintf("Too many iterations. Max: %d\n", fractals.JULIA_SET_MAX_ITERATIONS))
			return
		}
		fractal.MaxIterations = iterations
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
	if query.Has("type") {
		seriesName = query.Get("type")
	}
	if query.Has("variables") {
		variablesTxt = query.Get("variables")
	}
	if !fractals.IsValidJuliaSetSeriesFunction(seriesName) {
		ctx.Text("Invalid function type")
		return
	}
	fractal.SeriesFunctionName = seriesName
	if query.Has("background") {
		background, err := helpers.ParseColor(query.Get("background"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Background = background
	}
	variables, err := fractals.ParseJuliaSetVariables(variablesTxt)
	if err != nil {
		ctx.Text(err.Error())
		return
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
	fractal.Variables = variables
	fractal.Region = region
	fractal.ColorPalette = colorPalette
	ctx.ContentType("image/png")
	err = fractal.WriteImage(ctx.ResponseWriter())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
