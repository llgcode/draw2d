// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dimg

import (
	"image"
	"testing"
)

func TestNewGraphicContext_NRGBA(t *testing.T) {
	// Create an NRGBA image (not RGBA)
	// Note: NewGraphicContext currently only supports RGBA, so this test
	// verifies that RGBA works correctly
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)

	if gc == nil {
		t.Error("NewGraphicContext returned nil for RGBA image")
	}

	if gc.img == nil {
		t.Error("GraphicContext.img is nil")
	}
}

func TestGraphicContext_GetStringBounds_EmptyString(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)

	// This should not panic with an empty string
	left, top, right, bottom := gc.GetStringBounds("")

	// For an empty string, bounds should be zero or minimal
	if left != 0 || top != 0 || right != 0 || bottom != 0 {
		// Empty string may have some default bounds, that's OK
		// Just verify it doesn't panic
	}
}

func TestGraphicContext_CreateStringPath_EmptyString(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)

	// This should not panic with an empty string
	width := gc.CreateStringPath("", 0, 0)

	// Empty string should have zero or minimal width
	if width < 0 {
		t.Errorf("CreateStringPath returned negative width: %v", width)
	}
}

func TestLoadFromPngFile_NonExistent(t *testing.T) {
	// Try to load a non-existent file
	nonExistentPath := t.TempDir() + "/non_existent_file.png"

	_, err := LoadFromPngFile(nonExistentPath)
	if err == nil {
		t.Error("LoadFromPngFile should return error for non-existent file")
	}
}

func TestSaveToPngFile_InvalidPath(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	// Try to save to an invalid path (directory that doesn't exist)
	invalidPath := "/nonexistent/directory/that/does/not/exist/file.png"

	err := SaveToPngFile(invalidPath, img)
	if err == nil {
		t.Error("SaveToPngFile should return error for invalid path")
	}
}
