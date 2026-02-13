// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 11/02/2026 by Copilot

package draw2dgles2

import "math"

// Point2D represents a 2D point
type Point2D struct {
	X, Y float32
}

// signedArea2 computes twice the signed area of a polygon.
// Positive = clockwise in screen coordinates (Y-down), Negative = counter-clockwise.
func signedArea2(vertices []Point2D) float32 {
	var sum float32
	n := len(vertices)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		sum += vertices[i].X*vertices[j].Y - vertices[j].X*vertices[i].Y
	}
	return sum
}

// Triangulate converts a polygon (list of vertices) into triangles using ear-clipping algorithm.
// Automatically detects polygon winding order (CW or CCW).
// Returns a list of triangle indices.
func Triangulate(vertices []Point2D) []uint16 {
	if len(vertices) < 3 {
		return nil
	}

	// Remove trailing duplicate/near-duplicate vertices that match the first vertex.
	// This is common from path Close operations that add a LineTo back to start.
	for len(vertices) > 3 {
		last := len(vertices) - 1
		dx := vertices[last].X - vertices[0].X
		dy := vertices[last].Y - vertices[0].Y
		if dx*dx+dy*dy < 0.5 {
			vertices = vertices[:last]
		} else {
			break
		}
	}

	if len(vertices) < 3 {
		return nil
	}

	// Detect winding order using signed area
	cw := signedArea2(vertices) > 0

	// Create index list
	indices := make([]int, len(vertices))
	for i := range indices {
		indices[i] = i
	}

	var triangles []uint16
	count := len(indices)

	// Ear clipping algorithm
	for count > 3 {
		earFound := false

		for i := 0; i < count; i++ {
			prev := indices[(i+count-1)%count]
			curr := indices[i]
			next := indices[(i+1)%count]

			if isEar(vertices, indices, count, prev, curr, next, cw) {
				// Add triangle
				triangles = append(triangles, uint16(prev), uint16(curr), uint16(next))

				// Remove ear
				copy(indices[i:], indices[i+1:])
				count--
				earFound = true
				break
			}
		}

		if !earFound {
			// Degenerate polygon — use fan triangulation as fallback
			for i := 1; i < count-1; i++ {
				triangles = append(triangles, uint16(indices[0]), uint16(indices[i]), uint16(indices[i+1]))
			}
			break
		}
	}

	// Add final triangle
	if count == 3 {
		triangles = append(triangles, uint16(indices[0]), uint16(indices[1]), uint16(indices[2]))
	}

	return triangles
}

// isEar checks if the vertex at curr forms an ear.
// The cw parameter indicates whether the polygon has clockwise winding.
func isEar(vertices []Point2D, indices []int, count, prev, curr, next int, cw bool) bool {
	p1 := vertices[prev]
	p2 := vertices[curr]
	p3 := vertices[next]

	// Check if triangle is convex based on polygon winding.
	// For CW polygons: convex vertex has positive cross product.
	// For CCW polygons: convex vertex has negative cross product.
	cross := cross2D(sub2D(p2, p1), sub2D(p3, p2))
	if cw {
		if cross < 0 {
			return false
		}
	} else {
		if cross > 0 {
			return false
		}
	}

	// Check if any other vertex is inside this triangle
	for i := 0; i < count; i++ {
		idx := indices[i]
		if idx == prev || idx == curr || idx == next {
			continue
		}

		if pointInTriangle(vertices[idx], p1, p2, p3) {
			return false
		}
	}

	return true
}

// pointInTriangle checks if point p is inside triangle (a, b, c)
func pointInTriangle(p, a, b, c Point2D) bool {
	v0 := sub2D(c, a)
	v1 := sub2D(b, a)
	v2 := sub2D(p, a)

	dot00 := dot2D(v0, v0)
	dot01 := dot2D(v0, v1)
	dot02 := dot2D(v0, v2)
	dot11 := dot2D(v1, v1)
	dot12 := dot2D(v1, v2)

	invDenom := 1 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	return (u >= 0) && (v >= 0) && (u+v < 1)
}

func sub2D(a, b Point2D) Point2D {
	return Point2D{a.X - b.X, a.Y - b.Y}
}

func dot2D(a, b Point2D) float32 {
	return a.X*b.X + a.Y*b.Y
}

func cross2D(a, b Point2D) float32 {
	return a.X*b.Y - a.Y*b.X
}

// ConvertToFloat32 converts float64 coordinates to float32
func ConvertToFloat32(x, y float64) (float32, float32) {
	return float32(x), float32(y)
}

// distance calculates the distance between two points
func distance(a, b Point2D) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
