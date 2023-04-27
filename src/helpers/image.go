package helpers

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	LINE_WIDTH = 0.2
)

// Represents a rectangular region.
type Rect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Fills an image with color.
//  *image*: The image to fill.
//  *color*: The color to fill the image with.
func FillImage(image *image.RGBA, color color.RGBA) {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()
	for i := 0; i <= width; i++ {
		for j := 0; j <= height; j++ {
			image.SetRGBA(i, j, color)
		}
	}
}

// Draws a rectangle in an image.
//  *image*: The image to draw the rectangle in.
//  *x*: The horizontal offset of the rectangle.
//  *y*: The vertical offset of the rectangle.
//  *width*: The width of the rectangle.
//  *height*: The height of the rectangle.
//  *color*: The color to stroke the rectangle with.
func DrawRectangle(gc *draw2dimg.GraphicContext, x, y, width, height float64, color color.RGBA) {
	gc.SetStrokeColor(color)
	gc.SetLineWidth(LINE_WIDTH)
	gc.BeginPath()
	gc.MoveTo(x, y)
	gc.LineTo(x+width, y)
	gc.LineTo(x+width, y+height)
	gc.LineTo(x, y+height)
	gc.LineTo(x, y)
	gc.Close()
	gc.Stroke()
}

// Draws a filled rectangle in an image.
//  *image*: The image to draw the rectangle in.
//  *x*: The horizontal offset of the rectangle.
//  *y*: The vertical offset of the rectangle.
//  *width*: The width of the rectangle.
//  *height*: The height of the rectangle.
//  *color*: The color to fill the rectangle with.
func FillRectangle(gc *draw2dimg.GraphicContext, x, y, width, height float64, color color.RGBA) {
	gc.SetFillColor(color)
	gc.SetStrokeColor(color)
	gc.SetLineWidth(LINE_WIDTH)
	gc.BeginPath()
	gc.MoveTo(x, y)
	gc.LineTo(x+width, y)
	gc.LineTo(x+width, y+height)
	gc.LineTo(x, y+height)
	gc.LineTo(x, y)
	gc.Close()
	gc.FillStroke()
}

// Draws a line in an image.
//  *pt1*: The starting point of the line.
//  *pt2*: The ending point of the line.
//  *color*: The color to stroke the line with.
func DrawLine(gc *draw2dimg.GraphicContext, pt1, pt2 Point, color color.RGBA) {
	gc.SetStrokeColor(color)
	gc.SetLineWidth(LINE_WIDTH)
	gc.BeginPath()
	gc.MoveTo(pt1.X, pt1.Y)
	gc.LineTo(pt2.X, pt2.Y)
	gc.Close()
	gc.Stroke()
}

// Draws a triangle in an image.
//  *pt1*: The first point of the triangle.
//  *pt2*: The second point of the triangle.
//  *pt3*: The third point of the triangle.
//  *strokeColor*: The color to stroke the triangle with.
func DrawTriangle(gc *draw2dimg.GraphicContext, pt1, pt2, pt3 Point, strokeColor color.RGBA) {
	gc.SetStrokeColor(strokeColor)
	gc.SetLineWidth(LINE_WIDTH)
	gc.BeginPath()
	gc.MoveTo(pt1.X, pt1.Y)
	gc.LineTo(pt2.X, pt2.Y)
	gc.LineTo(pt3.X, pt3.Y)
	gc.LineTo(pt1.X, pt1.Y)
	gc.Close()
	gc.Stroke()
}

// Draws a filled triangle in an image.
//  *pt1*: The first point of the triangle.
//  *pt2*: The second point of the triangle.
//  *pt3*: The third point of the triangle.
//  *strokeColor*: The color to stroke the triangle with.
//  *fillColor*: The color to fill the triangle with.
func DrawFilledTriangle(gc *draw2dimg.GraphicContext, pt1, pt2, pt3 Point, strokeColor, fillColor color.RGBA) {
	gc.SetStrokeColor(strokeColor)
	gc.SetFillColor(fillColor)
	gc.SetLineWidth(LINE_WIDTH)
	gc.BeginPath()
	gc.MoveTo(pt1.X, pt1.Y)
	gc.LineTo(pt2.X, pt2.Y)
	gc.LineTo(pt3.X, pt3.Y)
	gc.LineTo(pt1.X, pt1.Y)
	gc.Close()
	gc.FillStroke()
}

// Adds a pixel to an image.
//  *x*: The horizontal position of the pixel.
//  *y*: The vertical position of the pixel.
func PutPixel(gc *draw2dimg.GraphicContext, x, y float64, color color.RGBA) {
	gc.SetFillColor(color)
	gc.SetStrokeColor(color)
	gc.SetLineWidth(LINE_WIDTH)
	gc.SetLineJoin(draw2d.MiterJoin)
	gc.BeginPath()
	gc.MoveTo(x, y)
	gc.LineTo(x+LINE_WIDTH/2, y)
	gc.LineTo(x+LINE_WIDTH/2, y+LINE_WIDTH/2)
	gc.LineTo(x, y+LINE_WIDTH/2)
	gc.LineTo(x, y)
	gc.Close()
	gc.FillStroke()
}
