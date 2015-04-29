package draw2d

import (
	"image/color"
)

// FillRule defines the fill rule used when fill
type FillRule int

const (
	// FillRuleEvenOdd determines the "insideness" of a point in the shape
	// by drawing a ray from that point to infinity in any direction
	// and counting the number of path segments from the given shape that the ray crosses.
	// If this number is odd, the point is inside; if even, the point is outside.
	FillRuleEvenOdd FillRule = iota
	// FillRuleWinding determines the "insideness" of a point in the shape
	// by drawing a ray from that point to infinity in any direction
	// and then examining the places where a segment of the shape crosses the ray.
	// Starting with a count of zero, add one each time a path segment crosses
	// the ray from left to right and subtract one each time
	// a path segment crosses the ray from right to left. After counting the crossings,
	// if the result is zero then the point is outside the path. Otherwise, it is inside.
	FillRuleWinding
)

// LineCap is the style of line extremities
type LineCap int

const (
	// RoundCap defines a rounded shape at the end of the line
	RoundCap LineCap = iota
	// ButtCap defines a squared shape exactly at the end of the line
	ButtCap
	// SquareCap defines a squared shape at the end of the line
	SquareCap
)

// LineJoin is the style of segments joint
type LineJoin int

const (
	// BevelJoin represents cut segments joint
	BevelJoin LineJoin = iota
	// RoundJoin represents rounded segments joint
	RoundJoin
	// MiterJoin represents peaker segments joint
	MiterJoin
)

// StrokeStyle keeps stroke style attributes
// that is used by the Stroke method of a Drawer
type StrokeStyle struct {
	// Color defines the color of stroke
	Color color.Color
	// Line width
	Width float64
	// Line cap style rounded, butt or square
	LineCap LineCap
	// Line join style bevel, round or miter
	LineJoin LineJoin
	// offset of the first dash
	dashOffset float64
	// array represented dash length pair values are plain dash and impair are space between dash
	// if empty display plain line
	dash []float64
}

// FillStyle
type FillStyle struct {
}

// SolidFillStyle define style attributes for a solid fill style
type SolidFillStyle struct {
	FillStyle
	// Color defines the line color
	Color color.Color
	// FillRule defines the file rule to used
	FillRule FillRule
}

type VerticalAlign int

type HorizontalAlign int

// TextStyle
type TextStyle struct {
	// Color defines the color of text
	Color color.Color
	// The font to use
	Font FontData
	// Horizontal Alignment of the text
	HorizontalAlign HorizontalAlign
	// Vertical Alignment of the text
	VerticalAlign VerticalAlign
}

// Drawer can fill and stroke a path
type Drawer interface {
	setMatrix(MatrixTransform)
	Fill(FillStyle, Path)
	Stroke(StrokeStyle, Path)
	Text(TextStyle, text string, x, y float64)
}
