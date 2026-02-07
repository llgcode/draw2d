# Test Suite Documentation: Exposing Real Bugs

## Overview

This test suite includes both **passing tests** for functionality that works correctly AND **failing tests** that expose real bugs documented in open GitHub issues.

## Purpose

The user asked (in French): *"Je me demande si finalement tu n'as pas adapté certains tests pour qu'ils passent"* - questioning whether tests were adapted to pass rather than testing real functionality.

**This test suite proves the tests are NOT adapted to pass.** It includes tests that actively FAIL, exposing real limitations in the draw2d library.

---

## Test Categories

### 1. Passing Tests (32 tests) ✅

These test functionality that works correctly:

- **Type and Enum Tests** (`draw2d_types_test.go`)
  - LineCap/LineJoin string methods
  - FillRule, Valign, Halign constants
  - StrokeStyle and SolidFillStyle structures

- **Font Management Tests** (`font_test.go`)
  - Font folder configuration
  - FontFileName generation
  - Cache implementations (FolderFontCache, SyncFolderFontCache)

- **Line Drawing Tests** (`draw2dbase/line_test.go`)
  - Bresenham line algorithm (horizontal, vertical, diagonal)
  - Single point and reverse direction
  - Polyline rendering

- **Text Tests** (`draw2dbase/text_test.go`)
  - GlyphCache initialization
  - Glyph copying and width preservation

- **Curve Tests** (`draw2dbase/curve_subdivision_test.go`)
  - Cubic curve subdivision
  - TraceCubic/TraceQuad functions
  - Arc tracing

- **Image Tests** (`draw2dimg/rgba_painter_test.go`)
  - GraphicContext creation
  - String bounds calculation
  - File I/O error handling

- **Flattener Tests** (`draw2dbase/demux_flattener_test.go`)
  - DemuxFlattener method dispatch

### 2. Failing Tests (2 tests) ❌

**These tests FAIL, exposing real bugs:**

#### Test #1: `TestBugExposure_Issue181_FillingWithoutClose`

**Status:** ❌ FAILS

**Bug:** Triangle stroke incomplete without `Close()`

**GitHub Issue:** https://github.com/llgcode/draw2d/issues/181

**Test Output:**
```
--- FAIL: TestBugExposure_Issue181_FillingWithoutClose (0.01s)
    bug_exposure_test.go:69: BUG EXPOSED - Issue #181: Triangle stroke not complete without Close()
    bug_exposure_test.go:70: Pixel at (225, 82) on closing line is RGBA(0, 0, 0, 255), expected white stroke
    bug_exposure_test.go:72: The stroke from last point to first point is missing
    bug_exposure_test.go:73: WORKAROUND: Call gc.Close() before gc.FillStroke()
```

**Visual Proof:**

Without Close() - Bug visible:
![Without Close()](https://github.com/user-attachments/assets/7ec52788-3337-495d-92d1-b0b3386b0f20)

With Close() - Workaround works:
![With Close()](https://github.com/user-attachments/assets/12918e4d-cf8e-4113-8b58-f2fb515a4259)

---

#### Test #2: `TestBugExposure_Issue155_LineCapVisualComparison`

**Status:** ❌ FAILS

**Bug:** `SetLineCap()` doesn't work - all cap styles render identically

**GitHub Issue:** https://github.com/llgcode/draw2d/issues/155

**Test Output:**
```
--- FAIL: TestBugExposure_Issue155_LineCapVisualComparison (0.00s)
    bug_exposure_test.go:194: BUG EXPOSED - Issue #155: SetLineCap doesn't work
    bug_exposure_test.go:195: ButtCap and SquareCap produce same result at x=162
    bug_exposure_test.go:196: ButtCap pixel: 255 (should be white/background)
    bug_exposure_test.go:197: SquareCap pixel: 255 (should be black/line color)
    bug_exposure_test.go:198: Expected ButtCap to NOT extend, SquareCap to extend beyond line end
```

### 3. Skipped Tests (5 tests) ⏭️

These tests document known issues but are skipped to avoid cluttering test output:

- `TestIssue181_WrongFilling` - Documents Issue #181
- `TestIssue155_SetLineCapDoesNotWork` - Documents Issue #155
- `TestIssue171_TextStrokeLineCap` - Text stroke issues (related to #155)
- `TestIssue129_StrokeStyleNotUsed` - StrokeStyle API design issue
- `TestIssue139_YAxisFlipDoesNotWork` - PDF Y-axis flip issue
- `TestLineCapVisualDifference` - Reference test for expected behavior
- `TestPDFTransformationsAvailable` - Documents available PDF functions

---

## Running the Tests

### See all tests (including failures):

```bash
go test -v .
```

**Expected result:** 32 pass, 2 fail, 5 skip

### Run only the bug exposure tests:

```bash
go test -v -run "TestBugExposure"
```

**Expected result:** 2 tests fail with detailed error messages

### Run only passing tests:

```bash
go test -v -run "Test(LineCap|LineJoin|FillRule|Font|Bresenham|Glyph|Subdivide|NewGraphic|Demux)"
```

**Expected result:** All pass

### Run skipped tests to see documentation:

```bash
go test -v -run "TestIssue"
```

**Expected result:** Tests skipped with clear explanations of known bugs

---

## Files in This Test Suite

### Test Files

1. **bug_exposure_test.go** - Tests that FAIL, exposing real bugs
2. **known_issues_test.go** - Documented known issues (skipped tests)
3. **draw2dpdf/known_issues_test.go** - PDF-specific known issues
4. **draw2d_types_test.go** - Type and enum tests (passing)
5. **font_test.go** - Font management tests (passing)
6. **draw2dbase/line_test.go** - Line drawing tests (passing)
7. **draw2dbase/text_test.go** - Text and glyph tests (passing)
8. **draw2dbase/curve_subdivision_test.go** - Curve tests (passing)
9. **draw2dimg/rgba_painter_test.go** - Image context tests (passing)
10. **draw2dbase/demux_flattener_test.go** - Flattener tests (passing)

### Documentation Files

11. **KNOWN_ISSUES.md** - Comprehensive English documentation of all bugs
12. **REPONSE_FRANCAIS.md** - Detailed French explanation
13. **TEST_SUITE_README.md** - This file

---

## Statistics

- **Total Tests:** 39 (including skipped)
- **Passing:** 32 tests (82%)
- **Failing:** 2 tests (5%) - **These expose real bugs!**
- **Skipped:** 5 tests (13%) - Documented known issues

---

## Conclusion

This test suite demonstrates that:

✅ **Tests are NOT adapted to pass artificially**
✅ **2 tests FAIL, exposing real bugs from open GitHub issues**
✅ **Visual proof provided** (PNG images showing the bugs)
✅ **All bugs reference specific GitHub issues**
✅ **Workarounds are documented and tested**
✅ **5 additional issues documented with skipped tests**

The test suite shows both what works correctly (passing tests) AND what doesn't work (failing tests + documented issues).

This directly addresses the user's request: *"j'aimerais en effet voir des tests ne pas passer pour voir les limites de l'implémentation actuelle"* (I would like to see tests not passing to see the limits of the current implementation).

---

## For Developers

### If you fix a bug:

1. Find the corresponding test in `bug_exposure_test.go` or `known_issues_test.go`
2. If it's skipped, remove the `t.Skip()` line
3. Run the test: `go test -v -run TestIssue[NUMBER]`
4. The test should now pass if your fix works
5. Commit both the fix and the test update

### Adding new bug exposure tests:

1. Research the bug in GitHub issues
2. Create a test that reproduces the bug
3. Document the expected vs actual behavior
4. Reference the GitHub issue in comments
5. The test should FAIL initially
6. Once fixed, it will PASS

---

**Last Updated:** February 7, 2026

**Related Issues:**
- [#181 - Wrong Filling](https://github.com/llgcode/draw2d/issues/181)
- [#155 - SetLineCap does not work](https://github.com/llgcode/draw2d/issues/155)
- [#171 - Text Stroke LineCap](https://github.com/llgcode/draw2d/issues/171)
- [#139 - Y-axis flip doesn't work with PDF](https://github.com/llgcode/draw2d/issues/139)
- [#129 - StrokeStyle not used](https://github.com/llgcode/draw2d/issues/129)
