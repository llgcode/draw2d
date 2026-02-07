# Known Issues and Failing Tests

This document describes the failing tests that demonstrate real bugs and limitations in draw2d's current implementation.

## Purpose

These tests are **expected to fail** and serve to:
1. Document known bugs tracked in GitHub issues
2. Provide reproducible test cases for each bug
3. Help developers understand the limits of the current implementation
4. Serve as validation when bugs are fixed (the tests should pass once fixed)

## Failing Tests

### Issue #155: SetLineCap Does Not Work

**Tests:**
- `TestIssue155_SetLineCapButtCap` (FAILS)
- `TestIssue155_SetLineCapSquareCap` (FAILS)
- `draw2dimg/TestIssue155_LineCapVisualDifference` (FAILS)

**Problem:** Different line caps (ButtCap, RoundCap, SquareCap) all render identically. The line end appearance doesn't change when calling `SetLineCap()` with different values.

**Expected:** Each line cap should produce visually different line endings:
- ButtCap: Line ends exactly at the endpoint
- RoundCap: Line extends with a rounded cap
- SquareCap: Line extends with a square cap (by half the line width)

**Actual:** All line caps render the same way, showing no visual difference.

**Issue:** https://github.com/llgcode/draw2d/issues/155

### Issue #139: PDF Y-Axis Flipping Doesn't Work

**Test:** `TestIssue139_PDFVerticalFlip` (FAILS)

**Problem:** Calling `Scale(1, -1)` to flip the Y-axis doesn't work properly with the PDF backend (`draw2dpdf.GraphicContext`), while it works fine with the image backend.

**Expected:** The transformation matrix should have Y scale = -1 after calling `Scale(1, -1)`.

**Actual:** The matrix Y scale remains 1, indicating the transformation wasn't applied properly.

**Issue:** https://github.com/llgcode/draw2d/issues/139

### Issue #171: Text Stroke Disconnections

**Test:** `TestIssue171_TextStrokeDisconnected` (SKIPPED - requires visual inspection)

**Problem:** When drawing text with stroke, the stroke has gaps and disconnections, especially visible in letters like 'i' and 't'. This is related to issue #155 (SetLineCap not working).

**Issue:** https://github.com/llgcode/draw2d/issues/171

### Issue #181: Triangle Filling Without Close

**Test:** `TestIssue181_TriangleFillingWithoutClose` (PASSES - bug may be fixed or misunderstood)

**Note:** This test currently passes, indicating that FillStroke() does work without an explicit Close() call. The original issue may have been about a different aspect or may have been fixed.

**Workaround:** Calling `Close()` explicitly before `FillStroke()` ensures proper filling (verified in `TestIssue181_TriangleFillingWithClose`).

**Issue:** https://github.com/llgcode/draw2d/issues/181

### Issue #143: Unsupported Image Types

**Test:** `draw2dimg/TestIssue143_UnsupportedImageTypesDocumented`

**Problem:** Only `*image.RGBA` is supported. Other image types like `image.Paletted`, `image.Gray`, etc. cause panics.

**Expected:** Support for more image types or graceful error handling.

**Actual:** Panics with "Image type not supported".

**Issue:** https://github.com/llgcode/draw2d/issues/143

### Issue #147: Performance (~10-30x slower than Cairo)

**Test:** `TestPerformanceNote` (Informational)

**Benchmarks:** Run with `go test -bench=. -benchmem`

**Problem:** draw2d is significantly slower than Cairo for similar operations.

**Example Results:**
```
BenchmarkFillStrokeRectangle  - measures FillStroke performance
BenchmarkStrokeSimpleLine     - measures simple line drawing
BenchmarkFillCircle           - measures circle filling
```

**Issue:** https://github.com/llgcode/draw2d/issues/147

## Running the Tests

To see the failing tests:

```bash
# Run all known issue tests
go test -v -run "TestIssue"

# Run specific issue tests
go test -v -run "TestIssue155"
go test -v -run "TestIssue139"

# Run benchmark tests
go test -bench=. -benchmem

# Run tests in draw2dimg
go test -v ./draw2dimg -run "TestIssue"
```

## When These Tests Pass

When a bug is fixed, the corresponding test should start passing. This indicates that:
1. The bug has been successfully addressed
2. The test can be moved to the regular test suite
3. The issue can be closed on GitHub

## Contributing

If you're working on fixing one of these issues:
1. Make your changes
2. Run the corresponding test to verify it now passes
3. Ensure all other tests still pass
4. Update this README to note the fix
5. Reference the test in your pull request

## Analysis

### Why Were Original Tests Passing?

The original test suite only tested working functionality, not edge cases or known bugs. This gave a false sense of completeness. These failing tests reveal:

1. **Real Implementation Limits:** The line cap/join rendering code may not be fully implemented
2. **Backend Differences:** PDF backend doesn't properly support all transformations
3. **Performance Characteristics:** The pure-Go implementation has inherent performance trade-offs

### Lessons Learned

- Tests should include **negative test cases** that verify bugs are tracked
- Tests should validate **expected failures** for known limitations  
- Performance benchmarks help quantify trade-offs
- Visual rendering tests are challenging to automate but critical for graphics libraries
