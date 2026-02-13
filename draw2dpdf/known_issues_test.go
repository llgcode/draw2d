// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

// This file contains tests for known issues specific to the PDF backend

package draw2dpdf_test

import (
	"testing"

	"github.com/llgcode/draw2d/draw2dpdf"
)

// TestIssue139_YAxisFlipDoesNotWork tests that Y-axis flipping doesn't work with PDF.
// Issue: https://github.com/llgcode/draw2d/issues/139
// Expected: Scale(1, -1) should flip the Y axis for PDF context just like it does for image context
// Actual: The transformation silently fails with draw2dpdf.GraphicContext
//
// This test demonstrates that while draw2dimg.GraphicContext properly handles
// negative scaling for Y-axis flipping, draw2dpdf.GraphicContext does not.
func TestIssue139_YAxisFlipDoesNotWork(t *testing.T) {
	t.Skip("Known issue #139: Flipping Y axis doesn't work with draw2dpdf.GraphicContext")

	// Create a PDF graphic context
	pdf := draw2dpdf.NewPdf("P", "mm", "A4")
	gc := draw2dpdf.NewGraphicContext(pdf)

	// Get initial transformation matrix
	initialMatrix := gc.GetMatrixTransform()

	// Try to flip Y axis (this should work but doesn't)
	height := 297.0 // A4 height in mm
	gc.Translate(0, height)
	gc.Scale(1, -1)

	// Get transformed matrix
	transformedMatrix := gc.GetMatrixTransform()

	// Check if transformation was actually applied
	// The Y scale component should be negative
	_, scaleY := transformedMatrix.GetScaling()
	if scaleY >= 0 {
		_, initialScaleY := initialMatrix.GetScaling()
		t.Errorf("Bug confirmed: Y-axis flip not applied. Initial ScaleY: %v, After flip ScaleY: %v", initialScaleY, scaleY)
	}

	// Even if the matrix is set, rendering might not respect it
	// The underlying gofpdf library has TransformScale but may not be called
	t.Logf("Known issue: PDF backend doesn't properly handle negative scaling for Y-axis flip")
	t.Logf("Initial matrix: %+v", initialMatrix)
	t.Logf("Transformed matrix: %+v", transformedMatrix)
}

// TestPDFTransformationsAvailable documents that gofpdf has transformation functions.
// Issue: https://github.com/llgcode/draw2d/issues/139
// This test documents that the underlying gofpdf library has the necessary functions,
// but they may not be properly integrated with draw2dpdf.GraphicContext.
func TestPDFTransformationsAvailable(t *testing.T) {
	t.Skip("Reference test documenting available gofpdf transformation functions")

	// The gofpdf package provides these transformation functions:
	// - Transform(tm TransformMatrix)
	// - TransformBegin()
	// - TransformEnd()
	// - TransformScale(scaleWd, scaleHt, x, y float64)
	// - TransformMirrorVertical(y float64)
	// - TransformMirrorHorizontal(x float64)
	// etc.
	//
	// However, draw2dpdf.GraphicContext.Scale() may not properly call these functions
	// when dealing with negative scale values.

	t.Logf("Reference: gofpdf provides TransformScale() and TransformMirrorVertical()")
	t.Logf("Issue: draw2dpdf.GraphicContext doesn't properly integrate these for Y-axis flip")
}
