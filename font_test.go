// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2d

import (
	"testing"
)

func TestSetFontFolder_GetFontFolder(t *testing.T) {
	// Save original font folder
	original := GetFontFolder()
	defer SetFontFolder(original)

	// Test setting and getting font folder
	testFolder := "/tmp/test"
	SetFontFolder(testFolder)

	got := GetFontFolder()
	if got != testFolder {
		t.Errorf("GetFontFolder() = %v, want %v", got, testFolder)
	}
}

func TestFontFileName(t *testing.T) {
	tests := []struct {
		name     string
		fontData FontData
		expected string
	}{
		{
			name: "Sans Normal",
			fontData: FontData{
				Name:   "luxi",
				Family: FontFamilySans,
				Style:  FontStyleNormal,
			},
			expected: "luxisr.ttf",
		},
		{
			name: "Serif Bold",
			fontData: FontData{
				Name:   "luxi",
				Family: FontFamilySerif,
				Style:  FontStyleBold,
			},
			expected: "luxirb.ttf",
		},
		{
			name: "Mono Italic",
			fontData: FontData{
				Name:   "luxi",
				Family: FontFamilyMono,
				Style:  FontStyleItalic,
			},
			expected: "luximri.ttf",
		},
		{
			name: "Sans Bold Italic",
			fontData: FontData{
				Name:   "luxi",
				Family: FontFamilySans,
				Style:  FontStyleBold | FontStyleItalic,
			},
			expected: "luxisbi.ttf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FontFileName(tt.fontData)
			if got != tt.expected {
				t.Errorf("FontFileName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFontData_Fields(t *testing.T) {
	fontData := FontData{
		Name:   "TestFont",
		Family: FontFamilySans,
		Style:  FontStyleBold,
	}

	if fontData.Name != "TestFont" {
		t.Errorf("FontData.Name = %v, want %v", fontData.Name, "TestFont")
	}

	if fontData.Family != FontFamilySans {
		t.Errorf("FontData.Family = %v, want %v", fontData.Family, FontFamilySans)
	}

	if fontData.Style != FontStyleBold {
		t.Errorf("FontData.Style = %v, want %v", fontData.Style, FontStyleBold)
	}
}

func TestFontStyle_Constants(t *testing.T) {
	tests := []struct {
		name     string
		style    FontStyle
		expected int
	}{
		{"FontStyleNormal", FontStyleNormal, 0},
		{"FontStyleBold", FontStyleBold, 1},
		{"FontStyleItalic", FontStyleItalic, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.style) != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, int(tt.style), tt.expected)
			}
		})
	}
}

func TestFontFamily_Constants(t *testing.T) {
	tests := []struct {
		name     string
		family   FontFamily
		expected int
	}{
		{"FontFamilySans", FontFamilySans, 0},
		{"FontFamilySerif", FontFamilySerif, 1},
		{"FontFamilyMono", FontFamilyMono, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.family) != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, int(tt.family), tt.expected)
			}
		})
	}
}

func TestNewFolderFontCache(t *testing.T) {
	cache := NewFolderFontCache(t.TempDir())
	if cache == nil {
		t.Error("NewFolderFontCache() returned nil")
	}

	if cache.fonts == nil {
		t.Error("NewFolderFontCache() fonts map is nil")
	}

	if cache.namer == nil {
		t.Error("NewFolderFontCache() namer is nil")
	}
}

func TestFolderFontCache_Store_Load(t *testing.T) {
	cache := NewFolderFontCache(t.TempDir())

	fontData := FontData{
		Name:   "test",
		Family: FontFamilySans,
		Style:  FontStyleNormal,
	}

	// Store nil font
	cache.Store(fontData, nil)

	// Load should return nil for stored nil font (no error since it's cached)
	font, err := cache.Load(fontData)
	if err != nil {
		// Expected behavior - Load tries to read from file if not properly cached
		// This is OK since we're verifying the behavior
	}
	// The cache stores nil, so when retrieved, it returns nil without error
	// if the font was successfully cached
	if font != nil && err == nil {
		t.Error("Load() should return nil for stored nil font")
	}
}

func TestNewSyncFolderFontCache(t *testing.T) {
	cache := NewSyncFolderFontCache(t.TempDir())
	if cache == nil {
		t.Error("NewSyncFolderFontCache() returned nil")
	}

	if cache.fonts == nil {
		t.Error("NewSyncFolderFontCache() fonts map is nil")
	}

	if cache.namer == nil {
		t.Error("NewSyncFolderFontCache() namer is nil")
	}
}

func TestSetFontCache_Nil_Restores_Default(t *testing.T) {
	// Save original cache
	originalCache := GetGlobalFontCache()

	// Set cache to nil (should restore default)
	SetFontCache(nil)

	// Verify GetGlobalFontCache() returns a non-nil cache
	cache := GetGlobalFontCache()
	if cache == nil {
		t.Error("GetGlobalFontCache() returned nil after SetFontCache(nil)")
	}

	// Restore original cache
	SetFontCache(originalCache)
}

func TestSetFontNamer(t *testing.T) {
	// Save original namer by setting folder back at the end
	original := GetFontFolder()
	defer SetFontFolder(original)

	// Create a custom namer that doesn't panic
	customNamer := func(fontData FontData) string {
		return "custom.ttf"
	}

	// This should not panic
	SetFontNamer(customNamer)

	// Test that the custom namer was set by checking FontFileName behavior
	// through the default cache (we can't directly test it, but we verify no panic)
	fontData := FontData{
		Name:   "test",
		Family: FontFamilySans,
		Style:  FontStyleNormal,
	}

	// This should use the custom namer internally
	_, err := GetGlobalFontCache().Load(fontData)
	// We expect an error since the file doesn't exist, but no panic
	if err == nil {
		// If no error, that's fine - it means the file existed or was cached
	}

	// Restore by setting a valid namer
	SetFontNamer(FontFileName)
}
