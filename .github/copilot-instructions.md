# Copilot Instructions for draw2d

## Project Overview

draw2d is a Go 2D vector graphics library with multiple backends:
- `draw2d` — Core package: interfaces (`GraphicContext`, `PathBuilder`), types (`Matrix`, `Path`), and font management
- `draw2dbase` — Base implementations shared across backends (`StackGraphicContext`, flattener, stroker, dasher)
- `draw2dimg` — Raster image backend (using freetype-go)
- `draw2dpdf` — PDF backend (using gofpdf)
- `draw2dsvg` — SVG backend
- `draw2dgl` — OpenGL backend
- `draw2dkit` — Drawing helpers (`Rectangle`, `Circle`, `Ellipse`, `RoundedRectangle`)
- `samples/` — Example drawings used as integration tests

## Language and Conventions

- All code, comments, commit messages, and documentation must be written in **English**.
- The project uses **Go 1.20+** (see `go.mod`).

## Code Style

### File Headers

Source files include a copyright header and creation date:
```go
// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff
```

New files should follow the same pattern with the current date and author name.

### Comments

- Exported types, functions, and methods must have GoDoc comments.
- Comments should start with the name of the thing being documented:
  ```go
  // Rectangle draws a rectangle using a path between (x1,y1) and (x2,y2)
  func Rectangle(path draw2d.PathBuilder, x1, y1, x2, y2 float64) {
  ```
- Package comments go in the main source file or a `doc.go` file.

### Naming

- Follow standard Go naming conventions (camelCase for unexported, PascalCase for exported).
- Backend packages are named `draw2d<backend>` (e.g., `draw2dimg`, `draw2dpdf`).
- The `GraphicContext` struct in each backend embeds `*draw2dbase.StackGraphicContext`.

### Error Handling

- Functions that can fail return `error` as the last return value.
- Do not silently ignore errors — log or return them.

## Testing

### Structure

- **Unit tests** go alongside the source file they test (e.g., `matrix_test.go` tests `matrix.go`).
- **Integration/sample tests** live in `samples_test.go`, `draw2dpdf/samples_test.go`, etc.
- Test output files go in the `output/` directory (generated, not committed).

### Writing Tests

- Use the standard `testing` package only — no external test frameworks.
- Use table-driven tests where multiple inputs share the same logic:
  ```go
  tests := []struct {
      name string
      // ...
  }{
      {"case1", ...},
      {"case2", ...},
  }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) { ... })
  }
  ```
- Tests must not depend on external resources (fonts, network) unless testing that specific integration.
- For image-based tests, use `image.NewRGBA(image.Rect(0, 0, w, h))` as the canvas.
- Use `t.TempDir()` for any file output in tests.
- Reference GitHub issue numbers in regression test comments:
  ```go
  // Test related to issue #95: DashVertexConverter state preservation
  ```

### Running Tests

```bash
go test ./...
go test -cover ./... | grep -v "no test"
```

### Test Coverage Goals

- Every exported function and method should have at least one unit test.
- Core types (`Matrix`, `Path`, `StackGraphicContext`) should have thorough coverage.
- Backend-specific operations (`Stroke`, `Fill`, `FillStroke`, `Clear`) should verify pixel output where possible.
- Known bugs in the issue tracker should have corresponding regression tests.

## Documentation

- When adding or changing public API, update the GoDoc comments accordingly.
- When fixing a bug, add a comment referencing the issue number.
- If a change affects behavior described in `README.md` or package READMEs, update them.
- The `samples/` directory serves as living documentation — keep samples working after changes.

## Architecture Notes

- All backends implement the `draw2d.GraphicContext` interface defined in `gc.go`.
- `draw2dbase.StackGraphicContext` provides the common state management (colors, transforms, font, path). Backends embed it and override rendering methods (`Stroke`, `Fill`, `FillStroke`, string drawing, etc.).
- The `draw2dkit` helpers operate on `draw2d.PathBuilder`, not `GraphicContext`, making them backend-agnostic.
- `Matrix` is a `[6]float64` affine transformation matrix. Coordinate system follows the HTML Canvas 2D Context conventions.
