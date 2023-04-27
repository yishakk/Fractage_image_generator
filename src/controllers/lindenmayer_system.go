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
	LSYSTEM_MAX_ITERATIONS                   = 500_000
	LSYSTEM_DEFAULT_ITERATIONS               = 6
	LSYSTEM_DEFAULT_AXIOM                    = "X"
	LSYSTEM_DEFAULT_RULES                    = "F=FF,X=F-[[X]+X]+F[+FX]-X"
	LSYSTEM_DEFAULT_TURNING_ANGLE            = 22.5
	LSYSTEM_DEFAULT_POSITION                 = "bottom-center"
	LSYSTEM_DEFAULT_LINE_WIDTH               = 0.6
	LSYSTEM_DEFAULT_LINE_LENGTH              = 5
	LSYSTEM_DEFAULT_LINE_WIDTH_INCREMENT     = 0.5
	LSYSTEM_DEFAULT_TURNING_ANGLE_INCREMENT  = 5
	LSYSTEM_DEFAULT_LINE_LENGTH_SCALE_FACTOR = 0.125
	LSYSTEM_DEFAULT_ANGLE                    = -90.0
	LSYSTEM_DEFAULT_DRAW_SYMBOLS             = "AB"
	LSYSTEM_DEFAULT_SKIP_SYMBOLS             = ""
)

func GetLindenmayerSystem(ctx iris.Context) {
	query := ctx.Request().URL.Query()
	fractal := fractals.LindenmayerSystem{
		Width:                 DEFAULT_WIDTH,
		Height:                DEFAULT_HEIGHT,
		Axiom:                 LSYSTEM_DEFAULT_AXIOM,
		Iterations:            LSYSTEM_DEFAULT_ITERATIONS,
		DrawSymbols:           LSYSTEM_DEFAULT_DRAW_SYMBOLS,
		SkipSymbols:           LSYSTEM_DEFAULT_SKIP_SYMBOLS,
		Angle:                 LSYSTEM_DEFAULT_ANGLE,
		UseRandomColors:       true,
		Focus:                 false,
		TurningAngle:          LSYSTEM_DEFAULT_TURNING_ANGLE,
		Position:              LSYSTEM_DEFAULT_POSITION,
		LineWidth:             LSYSTEM_DEFAULT_LINE_WIDTH,
		LineLength:            LSYSTEM_DEFAULT_LINE_LENGTH,
		LineWidthIncrement:    LSYSTEM_DEFAULT_LINE_WIDTH_INCREMENT,
		LineLengthScaleFactor: LSYSTEM_DEFAULT_LINE_LENGTH_SCALE_FACTOR,
		TurningAngleIncrement: LSYSTEM_DEFAULT_TURNING_ANGLE_INCREMENT,
		Background:            color.RGBA{255, 255, 255, 255},
	}
	rulesTxt := LSYSTEM_DEFAULT_RULES
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
	if query.Has("axiom") {
		fractal.Axiom = query.Get("axiom")
	}
	if query.Has("rules") {
		rulesTxt = query.Get("rules")
	}
	if query.Has("iterations") {
		iterations, err := strconv.Atoi(query.Get("iterations"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		if iterations < 0 || iterations > LSYSTEM_MAX_ITERATIONS {
			ctx.Text(fmt.Sprintf("Too many iterations. Max: %d\n", LSYSTEM_MAX_ITERATIONS))
			return
		}
		fractal.Iterations = iterations
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
	if query.Has("draw_symbols") {
		fractal.DrawSymbols = query.Get("draw_symbols")
	}
	if query.Has("skip_symbols") {
		fractal.SkipSymbols = query.Get("skip_symbols")
	}
	if query.Has("angle") {
		angle, err := strconv.ParseFloat(query.Get("angle"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Angle = angle
	}
	if query.Has("turning_angle") {
		turningAngle, err := strconv.ParseFloat(query.Get("turning_angle"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.TurningAngle = turningAngle
	}
	if query.Has("position") {
		fractal.Position = query.Get("position")
	}
	if query.Has("focus") {
		focus, err := strconv.ParseBool(query.Get("focus"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Focus = focus
	}
	if query.Has("line_width") {
		lineWidth, err := strconv.ParseFloat(query.Get("line_width"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.LineWidth = lineWidth
	}
	if query.Has("line_length") {
		lineLength, err := strconv.ParseFloat(query.Get("line_length"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.LineLength = lineLength
		fractal.Focus = false
	}
	if query.Has("line_length_scale") {
		lineLengthScaleFactor, err := strconv.ParseFloat(query.Get("line_length_scale"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.LineLengthScaleFactor = lineLengthScaleFactor
	}
	if query.Has("line_width_step") {
		lineWidthStep, err := strconv.ParseFloat(query.Get("line_width_step"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.LineWidthIncrement = lineWidthStep
	}
	if query.Has("turning_angle_step") {
		turningAngleStep, err := strconv.ParseFloat(query.Get("turning_angle_step"), 32)
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.TurningAngleIncrement = turningAngleStep
	}
	if query.Has("background") {
		background, err := helpers.ParseColor(query.Get("background"))
		if err != nil {
			ctx.Text(err.Error())
			return
		}
		fractal.Background = background
	}
	rules, err := fractals.ParseLindenmayerRules(rulesTxt)
	if err != nil {
		ctx.Text(err.Error())
		return
	}
	_, _, err = fractals.ParseLSystemPosition(fractal.Position, float64(fractal.Width), float64(fractal.Height))
	if err != nil {
		ctx.Text(err.Error())
		return
	}
	fractal.RewriteRules = rules
	ctx.ContentType("image/png")
	err = fractal.WriteImage(ctx.ResponseWriter())
	if err != nil {
		ctx.Text(err.Error())
	}
}
