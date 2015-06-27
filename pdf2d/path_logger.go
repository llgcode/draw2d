// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 26/06/2015 by Stani Michiels

package pdf2d

import (
	"bytes"
	"log"
	"strconv"
)

func ftoas(xs ...float64) string {
	var buffer bytes.Buffer
	for i, x := range xs {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(strconv.FormatFloat(x, 'f', 2, 64))
	}
	return buffer.String()
}

// PathLogger implements Vectorizer and applies the Matrix
// transformation tr. It is used as debugging middleware.
// It should wrap gofpdf.Fpdf directly.
type PathLogger struct {
	logger *log.Logger
	Next   Vectorizer
}

// NewPathLogger constructs a new PathLogger
func NewPathLogger(logger *log.Logger,
	vectorizer Vectorizer) *PathLogger {
	return &PathLogger{logger, vectorizer}
}

// MoveTo creates a new subpath that start at the specified point
func (pl *PathLogger) MoveTo(x, y float64) {
	pl.logger.Printf("MoveTo(x=%.2f, y=%.2f)", x, y)
	pl.Next.MoveTo(x, y)
}

// LineTo adds a line to the current subpath
func (pl *PathLogger) LineTo(x, y float64) {
	pl.logger.Printf("LineTo(x=%.2f, y=%.2f)", x, y)
	pl.Next.LineTo(x, y)
}

// CurveTo adds a quadratic bezier curve to the current subpath
func (pl *PathLogger) CurveTo(cx, cy, x, y float64) {
	pl.logger.Printf("CurveTo(cx=%.2f, cy=%.2f, x=%.2f, y=%.2f)", cx, cy, x, y)
	pl.Next.CurveTo(cx, cy, x, y)

}

// CurveBezierCubicTo adds a cubic bezier curve to the current subpath
func (pl *PathLogger) CurveBezierCubicTo(cx1, cy1,
	cx2, cy2, x, y float64) {
	pl.logger.Printf("CurveBezierCubicTo(cx1=%.2f, cy1=%.2f, cx2=%.2f, cy2=%.2f, x=%.2f, y=%.2f)", cx1, cy1, cx2, cy2, x, y)
	pl.Next.CurveBezierCubicTo(cx1, cy1, cx2, cy2, x, y)
}

// ArcTo adds an arc to the current subpath
func (pl *PathLogger) ArcTo(x, y, rx, ry, degRotate, degStart, degEnd float64) {
	pl.logger.Printf("ArcTo(x=%.2f, y=%.2f, rx=%.2f, ry=%.2f, degRotate=%.2f, degStart=%.2f, degEnd=%.2f)", x, y, rx, ry, degRotate, degStart, degEnd)
	pl.Next.ArcTo(x, y, rx, ry, degRotate, degStart, degEnd)
}

// ClosePath closes the subpath
func (pl *PathLogger) ClosePath() {
	pl.Next.ClosePath()
}
