// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 26/06/2015 by Stani Michiels

package pdf2d

import "github.com/stanim/draw2d"

// VertexMatrixTransform implements Vectorizer and applies the Matrix
// transformation tr. It is normally wrapped around gofpdf Fpdf.
type VertexMatrixTransform struct {
	tr   draw2d.MatrixTransform
	Next Vectorizer
}

func NewVertexMatrixTransform(tr draw2d.MatrixTransform,
	vectorizer Vectorizer) *VertexMatrixTransform {
	return &VertexMatrixTransform{tr, vectorizer}
}

// MoveTo creates a new subpath that start at the specified point
func (vmt *VertexMatrixTransform) MoveTo(x, y float64) {
	vmt.tr.Transform(&x, &y)
	vmt.Next.MoveTo(x, y)
}

// LineTo adds a line to the current subpath
func (vmt *VertexMatrixTransform) LineTo(x, y float64) {
	vmt.tr.Transform(&x, &y)
	vmt.Next.LineTo(x, y)
}

// CurveTo adds a quadratic bezier curve to the current subpath
func (vmt *VertexMatrixTransform) CurveTo(cx, cy, x, y float64) {
	vmt.tr.Transform(&cx, &cy, &x, &y)
	vmt.Next.CurveTo(cx, cy, x, y)

}

// CurveTo adds a cubic bezier curve to the current subpath
func (vmt *VertexMatrixTransform) CurveBezierCubicTo(cx1, cy1,
	cx2, cy2, x, y float64) {
	vmt.tr.Transform(&cx1, &cy1, &cx2, &cy2, &x, &y)
	vmt.Next.CurveBezierCubicTo(cx1, cy1, cx2, cy2, x, y)
}

// ArcTo adds an arc to the current subpath
func (vmt *VertexMatrixTransform) ArcTo(x, y, rx, ry, degRotate, degStart, degEnd float64) {
	vmt.tr.Transform(&x, &y)
	vmt.Next.ArcTo(x, y, rx, ry, degRotate, degStart, degEnd)
}

// ClosePath closes the subpath
func (vmt *VertexMatrixTransform) ClosePath() {
	vmt.Next.ClosePath()
}
