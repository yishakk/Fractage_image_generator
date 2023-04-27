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
	HOPALONG_MAX_RESOLUTION     = 5_000
	HOPALONG_DEFAULT_RESOLUTION = 5
	HOPALONG_DEFAULT_A          = 5
	HOPALONG_DEFAULT_B          = 1
	HOPALONG_DEFAULT_C          = 5
	HOPALONG_DEFAULT_D          = 0
	HOPALONG_DEFAULT_X          = -1
	HOPALONG_DEFAULT_Y          = 0
	HOPALONG_DEFAULT_Scale      = 5
	HOPALONG_DEFAULT_FXN_TYPE   = "classic_bm"
)

func GetHopalong(ctx iris.Context) {
	query := ctx.Request().URL.Query()
	fractal := fractals.Hopalong{
		Width:           DEFAULT_WIDTH,
		Height:          DEFAULT_HEIGHT,
		A:               HOPALONG_DEFAULT_A,
		B:               HOPALONG_DEFAULT_B,
		C:               HOPALONG_DEFAULT_C,
		D:               HOPALONG_DEFAULT_D,
		X:               HOPALONG_DEFAULT_X,
		Y:               HOPALONG_DEFAULT_Y,
		Type:            HOPALONG_DEFAULT_FXN_TYPE,
		Scale:           HOPALONG_DEFAULT_Scale,
		UseRandomColors: true,
		Resolution:      HOPALONG_DEFAULT_RESOLUTION,
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
	if query.Has("resolution") {
		resolution, err := strconv.Atoi(query.Get("resolution"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if resolution < 0 || resolution > HOPALONG_MAX_RESOLUTION {
			ctx.Text(fmt.Sprintf("Resolution is too high. Max: %d\n", HOPALONG_MAX_RESOLUTION))
			return
		}
		fractal.Resolution = resolution
	}
	if query.Has("a") {
		a, err := strconv.ParseFloat(query.Get("a"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.A = a
	}
	if query.Has("b") {
		b, err := strconv.ParseFloat(query.Get("b"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.B = b
	}
	if query.Has("c") {
		c, err := strconv.ParseFloat(query.Get("c"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.C = c
	}
	if query.Has("d") {
		d, err := strconv.ParseFloat(query.Get("d"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.D = d
	}
	if query.Has("x") {
		x, err := strconv.ParseFloat(query.Get("x"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.X = x
	}
	if query.Has("y") {
		y, err := strconv.ParseFloat(query.Get("y"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Y = y
	}
	if query.Has("scale") {
		scale, err := strconv.ParseFloat(query.Get("scale"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Scale = scale
	}
	if query.Has("type") {
		fxnType := query.Get("type")
		validFxn := false
		for k := range fractals.HOPALONG_TYPES {
			if k == fxnType {
				validFxn = true
				break
			}
		}
		if !validFxn {
			ctx.Text("Invalid function type")
			return
		}
		fractal.Type = fxnType
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
