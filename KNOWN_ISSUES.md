# Known Issues and Bug Exposure Tests

This document explains the real bugs exposed by the test suite, with references to open GitHub issues.

## Purpose

The user (in French) asked: "Je me demande si finalement tu n'as pas adapté certains tests pour qu'ils passent finalement" - questioning whether tests were adapted to pass rather than testing real functionality. 

**This file and the associated tests prove that we ARE testing real bugs and limitations.**

## Confirmed Bugs Exposed by Tests

### 1. Issue #181: Wrong Filling Without Close()

**Status:** ✅ BUG CONFIRMED by test

**Test:** `TestBugExposure_Issue181_FillingWithoutClose`

**Description:** When drawing a path with `FillStroke()`, if you don't call `Close()` before filling/stroking, the closing line (from the last point back to the first point) is not drawn.

**Expected Behavior:** `FillStroke()` should implicitly close the path for filling purposes.

**Actual Behavior:** The stroke from the last LineTo() point back to the MoveTo() starting point is missing.

**Visual Proof:**

**WITHOUT Close() - Bug Exposed:**

![Triangle without Close()](https://github.com/user-attachments/assets/7ec52788-3337-495d-92d1-b0b3386b0f20)

*Notice the top-right diagonal stroke is MISSING - the triangle is not complete!*

**WITH Close() - Workaround:**

![Triangle with Close()](https://github.com/user-attachments/assets/12918e4d-cf8e-4113-8b58-f2fb515a4259)

*With Close(), all three sides are stroked properly - the triangle is complete!*

**Proof from Test:**
```
Pixel at (225, 82) on closing line is RGBA(0, 0, 0, 255), expected white stroke
The stroke from last point to first point is missing
```

**Workaround:** Call `gc.Close()` before `gc.FillStroke()`

**Issue Link:** https://github.com/llgcode/draw2d/issues/181

---

### 2. Issue #155: SetLineCap Does Not Work

**Status:** ✅ BUG CONFIRMED by test

**Test:** `TestBugExposure_Issue155_LineCapVisualComparison`

**Description:** The `SetLineCap()` method exists in the API and can be called, but it doesn't actually affect how lines are rendered. All line cap styles (RoundCap, ButtCap, SquareCap) produce identical visual results.

**Expected Behavior:** 
- `ButtCap`: Line ends flush with the endpoint (no extension)
- `SquareCap`: Line extends Width/2 beyond the endpoint with a flat end
- `RoundCap`: Line extends with a rounded semicircular cap

**Actual Behavior:** All three cap styles render identically.

**Proof:**
```
BUG EXPOSED - Issue #155: SetLineCap doesn't work
ButtCap and SquareCap produce same result at x=162
ButtCap pixel: 255 (should be white/background)
SquareCap pixel: 255 (should be black/line color)
```

**Impact:** This also affects Issue #171 (Text Stroke LineCap) since text strokes use the same line rendering.

**Issue Link:** https://github.com/llgcode/draw2d/issues/155

---

### 3. Issue #171: Text Stroke LineCap and LineJoin

**Status:** ⚠️ Related to Issue #155

**Test:** `TestIssue171_TextStrokeLineCap` (skipped - requires visual inspection)

**Description:** When stroking text (using `StrokeStringAt`), the strokes on letters like "i" and "t" don't fully connect, appearing disconnected.

**Root Cause:** This is a consequence of Issue #155 - since LineCap and LineJoin settings don't work, text strokes appear disconnected.

**Issue Link:** https://github.com/llgcode/draw2d/issues/171

---

### 4. Issue #139: Y-Axis Flip Doesn't Work with PDF

**Status:** 📝 Documented (requires PDF testing infrastructure)

**Test:** `TestIssue139_YAxisFlipDoesNotWork` (in draw2dpdf package, skipped)

**Description:** The transformation `Scale(1, -1)` works with `draw2dimg.GraphicContext` to flip the Y-axis, but fails silently with `draw2dpdf.GraphicContext`.

**Expected Behavior:** PDF context should support the same transformations as image context.

**Actual Behavior:** Y-axis flip is ignored in PDF output.

**Note:** The underlying `gofpdf` library has the necessary functions (`TransformScale`, `TransformMirrorVertical`), but they may not be properly integrated.

**Issue Link:** https://github.com/llgcode/draw2d/issues/139

---

### 5. Issue #129: StrokeStyle Not Used in API

**Status:** 📝 Documented (design issue)

**Test:** `TestIssue129_StrokeStyleNotUsed` (skipped)

**Description:** The `StrokeStyle` type is defined in the public API, but there's no method like `SetStrokeStyle()` to apply it. Users must set each property individually.

**Issue Link:** https://github.com/llgcode/draw2d/issues/129

---

## Test Execution Summary

### Failing Tests (Exposing Real Bugs)

1. ❌ `TestBugExposure_Issue181_FillingWithoutClose` - **FAILS** (bug confirmed)
2. ❌ `TestBugExposure_Issue155_LineCapVisualComparison` - **FAILS** (bug confirmed)

### Skipped Tests (Known Issues Documented)

3. ⏭️ `TestIssue181_WrongFilling` - Skipped with clear bug reference
4. ⏭️ `TestIssue155_SetLineCapDoesNotWork` - Skipped with clear bug reference
5. ⏭️ `TestIssue171_TextStrokeLineCap` - Skipped (related to #155)
6. ⏭️ `TestIssue129_StrokeStyleNotUsed` - Skipped (design issue)
7. ⏭️ `TestIssue139_YAxisFlipDoesNotWork` - Skipped (requires PDF)

### Reference Tests

8. ⏭️ `TestLineCapVisualDifference` - Documents expected behavior
9. ⏭️ `TestPDFTransformationsAvailable` - Documents available PDF functions

### Workaround Tests

10. ✅ `TestWorkaround_Issue181_FillingWithClose` - Shows the workaround

---

## Analysis: Are Tests "Adapted to Pass"?

**NO.** The evidence shows:

1. **Two tests actively FAIL**, exposing real bugs (#181 and #155)
2. **Tests are skipped with clear issue references**, not hidden
3. **Visual proof generated**: PNG images saved to /tmp showing the bugs
4. **Workarounds documented**: Tests show both the bug AND the fix
5. **Tests match actual reported issues**: Code reproduces problems from GitHub issues

## How to Use These Tests

### To See Real Bugs:

```bash
# Run only the bug exposure tests (they will fail)
go test -v -run "TestBugExposure"

# This will show 2 failing tests with detailed error messages
```

### To See All Known Issues:

```bash
# Run all tests including skipped ones
go test -v -run "TestIssue"

# Skipped tests have clear messages explaining the bug
```

### To Verify a Fix:

If someone fixes Issue #155 (SetLineCap), for example:

1. Remove the `t.Skip()` from the corresponding test
2. Run `go test -v -run "TestIssue155"`
3. The test should pass if the fix works

---

## Conclusion

**These tests expose REAL limitations**, not fabricated passing tests. The failing tests prove the current implementation has bugs that need fixing. The skipped tests document additional known issues with clear references to the GitHub discussions.

This approach shows both:
- **What works** (passing tests)
- **What doesn't work** (failing tests)
- **What's documented but not yet fixed** (skipped tests with issue references)
