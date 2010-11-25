// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

type MatrixTransform [6]float

const (
	epsilon = 1e-6
)

func (tr MatrixTransform) TransformX(x, y float) float {
	return x*tr[0] + y*tr[2] + tr[4]
}

func (tr MatrixTransform) TransformY(x, y float) float {
	return x*tr[1] + y*tr[3] + tr[5]
}

func (tr MatrixTransform) Determinant() float {
	return tr[0]*tr[3] - tr[1]*tr[2]
}

func (tr MatrixTransform) InverseTransformX(x, y float) float {
	return ((x-tr[4])*tr[3] - (y-tr[5])*tr[2]) / tr.Determinant()
}

func (tr MatrixTransform) InverseTransformY(x, y float) float {
	return ((y-tr[5])*tr[0] - (x-tr[4])*tr[1]) / tr.Determinant()
}

func (tr MatrixTransform) Transform(points ...*float) {
	for i, j := 0, 1; j < len(points); i, j = i+2, j+2 {
		x := *points[i]
		y := *points[j]
		*points[i] = tr.TransformX(x, y)
		*points[j] = tr.TransformY(x, y)
	}
}


func (tr MatrixTransform) InverseTransform(points ...*float) {
	d := tr.Determinant() // matrix determinant
	for i, j := 0, 1; j < len(points); i, j = i+2, j+2 {
		x := *points[i]
		y := *points[j]
		*points[i] = ((x-tr[4])*tr[3] - (y-tr[5])*tr[2]) / d
		*points[j] = ((y-tr[5])*tr[0] - (x-tr[4])*tr[1]) / d
	}
}

// ******************** Vector transformations ********************

func (tr MatrixTransform) VectorTransformX(x, y float) float {
	return x*tr[0] + y*tr[2]
}

func (tr MatrixTransform) VectorTransformY(x, y float) float {
	return x*tr[1] + y*tr[3]
}

func (tr MatrixTransform) VectorInverseTransformX(x, y float) float {
	d := tr.Determinant() // matrix determinant
	return (x*tr[3] - y*tr[2]) / d
}

func (tr MatrixTransform) VectorInverseTransformY(x, y float) float {
	d := tr.Determinant() // matrix determinant
	return (y*tr[0] - x*tr[1]) / d
}

func (tr MatrixTransform) VectorTransform(points ...*float) {
	for i, j := 0, 1; j < len(points); i, j = i+2, j+2 {
		x := *points[i]
		y := *points[j]
		*points[i] = tr.VectorTransformX(x, y)
		*points[j] = tr.VectorTransformY(x, y)
	}
}

// ******************** Transformations creation ********************

/** Creates an identity transformation. */
func NewIdentityMatrix() MatrixTransform {
	return [6]float{1, 0, 0, 1, 0, 0}
}

/**
 * Creates a transformation with a translation, that,
 * transform point1 into point2.
 */
func NewTranslationMatrix(tx, ty float) MatrixTransform {
	return [6]float{1, 0, 0, 1, tx, ty}
}

/**
 * Creates a transformation with a sx, sy scale factor
 */
func NewScaleMatrix(sx, sy float) MatrixTransform {
	return [6]float{sx, 0, 0, sy, 0, 0}
}

/**
 * Creates a rotation transformation.
 */
func NewRotationMatrix(angle float) MatrixTransform {
	c := cos(angle)
	s := sin(angle)
	return [6]float{c, s, -s, c, 0, 0}
}

/**
 * Creates a transformation, combining a scale and a translation, that transform rectangle1 into rectangle2.
 */
func NewMatrixTransform(rectangle1 [4]float, rectangle2 [4]float) MatrixTransform {
	xScale := (rectangle2[2] - rectangle2[0]) / (rectangle1[2] - rectangle1[0])
	yScale := (rectangle2[3] - rectangle2[1]) / (rectangle1[3] - rectangle1[1])
	xOffset := rectangle2[0] - (rectangle1[0] * xScale)
	yOffset := rectangle2[1] - (rectangle1[1] * yScale)
	return [6]float{xScale, 0, 0, yScale, xOffset, yOffset}
}

// ******************** Transformations operations ********************

/**
 * Returns a transformation that is the inverse of the given transformation.
 */
func (tr MatrixTransform) GetInverseTransformation() MatrixTransform {
	d := tr.Determinant() // matrix determinant
	return [6]float{
		tr[3] / d,
		-tr[1] / d,
		-tr[2] / d,
		tr[0] / d,
		(tr[2]*tr[5] - tr[3]*tr[4]) / d,
		(tr[1]*tr[4] - tr[0]*tr[5]) / d}
}

/**
 * Returns a transformation that is the composition (tr2 o tr1) of the given
 * transformations tr2 and tr1.

 * For given point (x, y), the composed transformation is defined by the
 * equation:
 * (tr2 o tr1)(x, y) = tr2(tr1(x, y))
 */
func (tr1 MatrixTransform) GetComposedTransformation(tr2 MatrixTransform) MatrixTransform {
	return [6]float{
		tr1[0]*tr2[0] + tr1[1]*tr2[2],
		tr1[1]*tr2[3] + tr1[0]*tr2[1],
		tr1[2]*tr2[0] + tr1[3]*tr2[2],
		tr1[3]*tr2[3] + tr1[2]*tr2[1],
		tr1[4]*tr2[0] + tr1[5]*tr2[2] + tr2[4],
		tr1[5]*tr2[3] + tr1[4]*tr2[1] + tr2[5]}
}

func (tr1 *MatrixTransform) Compose(tr2 MatrixTransform) (*MatrixTransform){
	tr1[0] = tr2[0]*tr1[0] + tr2[1]*tr1[2]
	tr1[1] = tr2[1]*tr1[3] + tr2[0]*tr1[1]
	tr1[2] = tr2[2]*tr1[0] + tr2[3]*tr1[2]
	tr1[3] = tr2[3]*tr1[3] + tr2[2]*tr1[1]
	tr1[4] = tr2[4]*tr1[0] + tr2[5]*tr1[2] + tr1[4]
	tr1[5] = tr2[5]*tr1[3] + tr2[4]*tr1[1] + tr1[5]
	return tr1
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
func fequals(float1, float2 float) bool {
	return fabs(float1-float2) <= epsilon
}
