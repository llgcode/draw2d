// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

import (
	"math"
)

type MatrixTransform [6]float64

const (
	epsilon = 1e-6
)

// Determinant compute the determinant of the matrix
func (tr MatrixTransform) Determinant() float64 {
	return tr[0]*tr[3] - tr[1]*tr[2]
}

// Transform apply the transformation matrix to points. It modify the points passed in parameter.
func (tr MatrixTransform) Transform(points []float64) {
	for i, j := 0, 1; j < len(points); i, j = i+2, j+2 {
		x := points[i]
		y := points[j]
		points[i] = x*tr[0] + y*tr[2] + tr[4]
		points[j] = x*tr[1] + y*tr[3] + tr[5]
	}
}

// TransformPoint apply the transformation matrix to point. It returns the point the transformed point.
func (tr MatrixTransform) TransformPoint(x, y float64) (xres, yres float64) {
	xres = x*tr[0] + y*tr[2] + tr[4]
	yres = x*tr[1] + y*tr[3] + tr[5]
	return xres, yres
}

func minMax(x, y float64) (min, max float64) {
	if x > y {
		return y, x
	}
	return x, y
}

// Transform apply the transformation matrix to the rectangle represented by the min and the max point of the rectangle
func (tr MatrixTransform) TransformRectangle(x0, y0, x2, y2 float64) (nx0, ny0, nx2, ny2 float64) {
	points := []float64{x0, y0, x2, y0, x2, y2, x0, y2}
	tr.Transform(points)
	points[0], points[2] = minMax(points[0], points[2])
	points[4], points[6] = minMax(points[4], points[6])
	points[1], points[3] = minMax(points[1], points[3])
	points[5], points[7] = minMax(points[5], points[7])

	nx0 = math.Min(points[0], points[4])
	ny0 = math.Min(points[1], points[5])
	nx2 = math.Max(points[2], points[6])
	ny2 = math.Max(points[3], points[7])
	return nx0, ny0, nx2, ny2
}

// InverseTransform apply the transformation inverse matrix to the rectangle represented by the min and the max point of the rectangle
func (tr MatrixTransform) InverseTransform(points []float64) {
	d := tr.Determinant() // matrix determinant
	for i, j := 0, 1; j < len(points); i, j = i+2, j+2 {
		x := points[i]
		y := points[j]
		points[i] = ((x-tr[4])*tr[3] - (y-tr[5])*tr[2]) / d
		points[j] = ((y-tr[5])*tr[0] - (x-tr[4])*tr[1]) / d
	}
}

// InverseTransformPoint apply the transformation inverse matrix to point. It returns the point the transformed point.
func (tr MatrixTransform) InverseTransformPoint(x, y float64) (xres, yres float64) {
	d := tr.Determinant() // matrix determinant
	xres = ((x-tr[4])*tr[3] - (y-tr[5])*tr[2]) / d
	yres = ((y-tr[5])*tr[0] - (x-tr[4])*tr[1]) / d
	return xres, yres
}

// VectorTransform apply the transformation matrix to points without using the translation parameter of the affine matrix.
// It modify the points passed in parameter.
func (tr MatrixTransform) VectorTransform(points []float64) {
	for i, j := 0, 1; j < len(points); i, j = i+2, j+2 {
		x := points[i]
		y := points[j]
		points[i] = x*tr[0] + y*tr[2]
		points[j] = x*tr[1] + y*tr[3]
	}
}

// NewIdentityMatrix creates an identity transformation matrix.
func NewIdentityMatrix() MatrixTransform {
	return [6]float64{1, 0, 0, 1, 0, 0}
}

// NewTranslationMatrix creates a transformation matrix with a translation tx and ty translation parameter
func NewTranslationMatrix(tx, ty float64) MatrixTransform {
	return [6]float64{1, 0, 0, 1, tx, ty}
}

// NewScaleMatrix creates a transformation matrix with a sx, sy scale factor
func NewScaleMatrix(sx, sy float64) MatrixTransform {
	return [6]float64{sx, 0, 0, sy, 0, 0}
}

// NewRotationMatrix creates a rotation transformation matrix. angle is in radian
func NewRotationMatrix(angle float64) MatrixTransform {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return [6]float64{c, s, -s, c, 0, 0}
}

// NewMatrixTransform creates a transformation matrix, combining a scale and a translation, that transform rectangle1 into rectangle2.
func NewMatrixFromRects(rectangle1, rectangle2 [4]float64) MatrixTransform {
	xScale := (rectangle2[2] - rectangle2[0]) / (rectangle1[2] - rectangle1[0])
	yScale := (rectangle2[3] - rectangle2[1]) / (rectangle1[3] - rectangle1[1])
	xOffset := rectangle2[0] - (rectangle1[0] * xScale)
	yOffset := rectangle2[1] - (rectangle1[1] * yScale)
	return [6]float64{xScale, 0, 0, yScale, xOffset, yOffset}
}

// Inverse returns a matrix that is the inverse of the given matrix.
func (tr MatrixTransform) Inverse() MatrixTransform {
	d := tr.Determinant() // matrix determinant
	return [6]float64{
		tr[3] / d,
		-tr[1] / d,
		-tr[2] / d,
		tr[0] / d,
		(tr[2]*tr[5] - tr[3]*tr[4]) / d,
		(tr[1]*tr[4] - tr[0]*tr[5]) / d}
}

// Multiply Compose Matrix tr1 with tr2 returns the resulting matrix
func (tr1 MatrixTransform) Multiply(tr2 MatrixTransform) MatrixTransform {
	return [6]float64{
		tr1[0]*tr2[0] + tr1[1]*tr2[2],
		tr1[1]*tr2[3] + tr1[0]*tr2[1],
		tr1[2]*tr2[0] + tr1[3]*tr2[2],
		tr1[3]*tr2[3] + tr1[2]*tr2[1],
		tr1[4]*tr2[0] + tr1[5]*tr2[2] + tr2[4],
		tr1[5]*tr2[3] + tr1[4]*tr2[1] + tr2[5]}
}

// Scale add a scale to the matrix
func (tr *MatrixTransform) Scale(sx, sy float64) *MatrixTransform {
	tr[0] = sx * tr[0]
	tr[1] = sx * tr[1]
	tr[2] = sy * tr[2]
	tr[3] = sy * tr[3]
	return tr
}

// Translate add a translation to the matrix
func (tr *MatrixTransform) Translate(tx, ty float64) *MatrixTransform {
	tr[4] = tx*tr[0] + ty*tr[2] + tr[4]
	tr[5] = ty*tr[3] + tx*tr[1] + tr[5]
	return tr
}

// Rotate add a rotation to the matrix. angle is in radian
func (tr *MatrixTransform) Rotate(angle float64) *MatrixTransform {
	c := math.Cos(angle)
	s := math.Sin(angle)
	t0 := c*tr[0] + s*tr[2]
	t1 := s*tr[3] + c*tr[1]
	t2 := c*tr[2] - s*tr[0]
	t3 := c*tr[3] - s*tr[1]
	tr[0] = t0
	tr[1] = t1
	tr[2] = t2
	tr[3] = t3
	return tr
}

// GetTranslation
func (tr MatrixTransform) GetTranslation() (x, y float64) {
	return tr[4], tr[5]
}

// GetScaling
func (tr MatrixTransform) GetScaling() (x, y float64) {
	return tr[0], tr[3]
}

// GetScale computes the scale of the matrix
func (tr MatrixTransform) GetScale() float64 {
	x := 0.707106781*tr[0] + 0.707106781*tr[1]
	y := 0.707106781*tr[2] + 0.707106781*tr[3]
	return math.Sqrt(x*x + y*y)
}

func (tr MatrixTransform) GetMaxAbsScaling() (s float64) {
	sx := math.Abs(tr[0])
	sy := math.Abs(tr[3])
	if sx > sy {
		return sx
	}
	return sy
}

func (tr MatrixTransform) GetMinAbsScaling() (s float64) {
	sx := math.Abs(tr[0])
	sy := math.Abs(tr[3])
	if sx > sy {
		return sy
	}
	return sx
}

// ******************** Testing ********************

/**
 * Tests if a two transformation are equal. A tolerance is applied when
 * comparing matrix elements.
 */
func (tr1 MatrixTransform) Equals(tr2 MatrixTransform) bool {
	for i := 0; i < 6; i = i + 1 {
		if !fequals(tr1[i], tr2[i]) {
			return false
		}
	}
	return true
}

/**
 * Tests if a transformation is the identity transformation. A tolerance
 * is applied when comparing matrix elements.
 */
func (tr MatrixTransform) IsIdentity() bool {
	return fequals(tr[4], 0) && fequals(tr[5], 0) && tr.IsTranslation()
}

/**
 * Tests if a transformation is is a pure translation. A tolerance
 * is applied when comparing matrix elements.
 */
func (tr MatrixTransform) IsTranslation() bool {
	return fequals(tr[0], 1) && fequals(tr[1], 0) && fequals(tr[2], 0) && fequals(tr[3], 1)
}

/**
 * Compares two floats.
 * return true if the distance between the two floats is less than epsilon, false otherwise
 */
func fequals(float1, float2 float64) bool {
	return math.Abs(float1-float2) <= epsilon
}

// Transformer apply the Matrix transformation tr
type Transformer struct {
	Tr        MatrixTransform
	Flattener Flattener
}

func (t Transformer) MoveTo(x, y float64) {
	u := x*t.Tr[0] + y*t.Tr[2] + t.Tr[4]
	v := x*t.Tr[1] + y*t.Tr[3] + t.Tr[5]
	t.Flattener.MoveTo(u, v)
}

func (t Transformer) LineTo(x, y float64) {
	u := x*t.Tr[0] + y*t.Tr[2] + t.Tr[4]
	v := x*t.Tr[1] + y*t.Tr[3] + t.Tr[5]
	t.Flattener.LineTo(u, v)
}

func (t Transformer) LineJoin() {
	t.Flattener.LineJoin()
}

func (t Transformer) Close() {
	t.Flattener.Close()
}

func (t Transformer) End() {
	t.Flattener.End()
}
