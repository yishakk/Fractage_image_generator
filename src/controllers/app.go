package controllers

import (
	"strconv"

	"github.com/B3zaleel/fractage/src/helpers"
	"github.com/kataras/iris/v12"
)

const (
	PALETTE_DEFAULT_WIDTH     = 200
	PALETTE_DEFAULT_HEIGHT    = 50
	PALETTE_DEFAULT_DIVISIONS = 5
	PALETTE_DEFAULT_VALUE     = "orange_blue"
)

func GetPalette(ctx iris.Context) {
	query := ctx.Request().URL.Query()
	width, height := int(PALETTE_DEFAULT_WIDTH), int(PALETTE_DEFAULT_HEIGHT)
	divisions := int(PALETTE_DEFAULT_DIVISIONS)
	colorPaletteValue := PALETTE_DEFAULT_VALUE
	var err error
	if query.Has("width") {
		width, err = strconv.Atoi(query.Get("width"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if width <= 0 {
			ctx.Text("width must be greater than 0")
			return
		}
	}
	if query.Has("height") {
		height, err = strconv.Atoi(query.Get("height"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if height <= 0 {
			ctx.Text("height must be greater than 0")
			return
		}
	}
	if query.Has("divisions") {
		divisions, err = strconv.Atoi(query.Get("divisions"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if divisions <= 0 {
			ctx.Text("divisions must be greater than 0")
			return
		}
	}
	if query.Has("value") {
		colorPaletteValue = query.Get("value")
	}
	colorPalette, err := helpers.ParseColorPalette(colorPaletteValue)
	if err != nil {
		ctx.Text(err.Error())
		return
	}
	step := 0.0
	if colorPalette.Transitions != nil {
		step = float64(width) / float64((len(colorPalette.Transitions) - 1) * divisions)
	}
	ctx.ContentType("image/png")
	err = colorPalette.Render(ctx.ResponseWriter(), width, height, step)
	if err != nil {
		ctx.Text(err.Error())
	}
}
