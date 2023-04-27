package controllers

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/B3zaleel/fractage/src/fractals"
	"github.com/B3zaleel/fractage/src/helpers"
	math_helper "github.com/B3zaleel/fractage/src/helpers/math"
	"github.com/kataras/iris/v12"
)

const (
	NEWTON_BASIN_MAX_ITERATIONS        = 500_000
	NEWTON_BASIN_DEFAULT_ITERATIONS    = 32
	NEWTON_BASIN_DEFAULT_COLOR_PALETTE = "orange_blue"
	NEWTON_BASIN_DEFAULT_BAIL_OUT      = 1e15
	NEWTON_BASIN_DEFAULT_POLYNOMIAL    = "-1+x^5"
	NEWTON_BASIN_DEFAULT_REGION        = "-2, -1.5, 4, 3"
)

func GetNewtonBasin(ctx iris.Context) {
	query := ctx.Request().URL.Query()
	fractal := fractals.NewtonBasin{
		Width:            DEFAULT_WIDTH,
		Height:           DEFAULT_HEIGHT,
		MaxIterations:    NEWTON_BASIN_DEFAULT_ITERATIONS,
		BailOut:          NEWTON_BASIN_DEFAULT_BAIL_OUT,
		UseDynamicColors: true,
		Background:       color.RGBA{255, 255, 255, 255},
	}
	colorPaletteValue := NEWTON_BASIN_DEFAULT_COLOR_PALETTE
	regionValue := NEWTON_BASIN_DEFAULT_REGION
	polynomialValue := NEWTON_BASIN_DEFAULT_POLYNOMIAL
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
	if query.Has("polynomial") {
		polynomialValue = query.Get("polynomial")
	}
	if query.Has("color_palette") {
		colorPaletteValue = query.Get("color_palette")
		fractal.UseDynamicColors = false
	}
	if query.Has("iterations") {
		iterations, err := strconv.Atoi(query.Get("iterations"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if iterations < 0 || iterations > NEWTON_BASIN_MAX_ITERATIONS {
			ctx.Text(fmt.Sprintf("Too many iterations. Max: %d\n", NEWTON_BASIN_MAX_ITERATIONS))
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
	polynomial, err := math_helper.ParseCmplxPolynomial(polynomialValue)
	if err != nil {
		ctx.Text(err.Error())
		return
	}
	fractal.Polynomial = polynomial
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
