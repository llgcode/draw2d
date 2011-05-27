// Copyright 2011 The draw2d Authors. All rights reserved.
// created: 27/05/2011 by Laurent Le Goff
package raster

import (
	"math"
)

const (
	POLYGON_CLIP_NONE = iota
	POLYGON_CLIP_LEFT
	POLYGON_CLIP_RIGHT
	POLYGON_CLIP_TOP
	POLYGON_CLIP_BOTTOM
)

type Polygon []float64


type PolygonEdge struct {
	X, Slope            float64
	FirstLine, LastLine int
	Winding             int16
}


//! Calculates the edges of the polygon with transformation and clipping to edges array.
/*! \param startIndex the index for the first vertex.
 *  \param vertexCount the amount of vertices to convert.
 *  \param edges the array for result edges. This should be able to contain 2*aVertexCount edges.
 *  \param tr the transformation matrix for the polygon.
 *  \param aClipRectangle the clip rectangle.
 *  \return the amount of edges in the result.
 */
func (p Polygon) getEdges(startIndex, vertexCount int, edges []PolygonEdge, tr [6]float64, clipBound [4]float64) int {
	startIndex = startIndex * 2
	endIndex := startIndex + (vertexCount * 2)
	if endIndex > len(p) {
		endIndex = len(p)
	}

	x := p[startIndex]
	y := p[startIndex+1]
	// inline transformation
	prevX := x*tr[0] + y*tr[2] + tr[4]
	prevY := x*tr[1] + y*tr[3] + tr[5]

	//! Calculates the clip flags for a point.
	prevClipFlags := POLYGON_CLIP_NONE
	if prevX < clipBound[0] {
		prevClipFlags |= POLYGON_CLIP_LEFT
	} else if prevX >= clipBound[2] {
		prevClipFlags |= POLYGON_CLIP_RIGHT
	}

	if prevY < clipBound[1] {
		prevClipFlags |= POLYGON_CLIP_TOP
	} else if prevY >= clipBound[3] {
		prevClipFlags |= POLYGON_CLIP_BOTTOM
	}

	edgeCount := 0
	var k, clipFlags, clipSum, clipUnion int
	var xleft, yleft, xright, yright, oldY, maxX, minX float64
	var swapWinding int16
	for n := startIndex; n < endIndex; n = n + 2 {
		k = (n + 2) % len(p)
		x = p[k]*tr[0] + p[k+1]*tr[2] + tr[4]
		y = p[k]*tr[1] + p[k+1]*tr[3] + tr[5]

		//! Calculates the clip flags for a point.
		clipFlags = POLYGON_CLIP_NONE
		if prevX < clipBound[0] {
			clipFlags |= POLYGON_CLIP_LEFT
		} else if prevX >= clipBound[2] {
			clipFlags |= POLYGON_CLIP_RIGHT
		}
		if prevY < clipBound[1] {
			clipFlags |= POLYGON_CLIP_TOP
		} else if prevY >= clipBound[3] {
			clipFlags |= POLYGON_CLIP_BOTTOM
		}

		clipSum = prevClipFlags | clipFlags
		clipUnion = prevClipFlags & clipFlags

		// Skip all edges that are either completely outside at the top or at the bottom.
		if (clipUnion & (POLYGON_CLIP_TOP | POLYGON_CLIP_BOTTOM)) == 0 {
			if (clipUnion & POLYGON_CLIP_RIGHT) != 0 {
				// Both clip to right, edge is a vertical line on the right side
				if getVerticalEdge(prevY, y, clipBound[2], &(edges[edgeCount]), clipBound) {
					edgeCount++
				}
			} else if (clipUnion & POLYGON_CLIP_LEFT) != 0 {
				// Both clip to left, edge is a vertical line on the left side
				if getVerticalEdge(prevY, y, clipBound[0], &(edges[edgeCount]), clipBound) {
					edgeCount++
				}
			} else if (clipSum & (POLYGON_CLIP_RIGHT | POLYGON_CLIP_LEFT)) == 0 {
				// No clipping in the horizontal direction
				if getEdge(prevX, prevY, x, y, &(edges[edgeCount]), clipBound) {
					edgeCount++
				}
			} else {
				// Clips to left or right or both.

				if x < prevX {
					xleft, yleft = x, y
					xright, yright = prevX, prevY
					swapWinding = -1
				} else {
					xleft, yleft = prevX, prevY
					xright, yright = x, y
					swapWinding = 1
				}

				slope := (yright - yleft) / (xright - xleft)

				if (clipSum & POLYGON_CLIP_RIGHT) != 0 {
					// calculate new position for the right vertex
					oldY = yright
					maxX = clipBound[2]

					yright = yleft + (maxX-xleft)*slope
					xright = maxX

					// add vertical edge for the overflowing part
					if getVerticalEdge(yright, oldY, maxX, &(edges[edgeCount]), clipBound) {
						edges[edgeCount].Winding *= swapWinding
						edgeCount++
					}
				}

				if (clipSum & POLYGON_CLIP_LEFT) != 0 {
					// calculate new position for the left vertex
					oldY = yleft
					minX = clipBound[0]

					yleft = yleft + (minX-xleft)*slope
					xleft = minX

					// add vertical edge for the overflowing part
					if getVerticalEdge(oldY, yleft, minX, &(edges[edgeCount]), clipBound) {
						edges[edgeCount].Winding *= swapWinding
						edgeCount++
					}
				}

				if getEdge(xleft, yleft, xright, yright, &(edges[edgeCount]), clipBound) {
					edges[edgeCount].Winding *= swapWinding
					edgeCount++
				}
			}
		}

		prevClipFlags = clipFlags
		prevX = x
		prevY = y
	}

	return edgeCount
}


//! Creates a polygon edge between two vectors.
/*! Clips the edge vertically to the clip rectangle. Returns true for edges that
 *  should be rendered, false for others.
 */
func getEdge(x0, y0, x1, y1 float64, edge *PolygonEdge, clipBound [4]float64) bool {
	var startX, startY, endX, endY float64
	var winding int16

	if y0 <= y1 {
		startX = x0
		startY = y0
		endX = x1
		endY = y1
		winding = 1
	} else {
		startX = x1
		startY = y1
		endX = x0
		endY = y0
		winding = -1
	}

	// Essentially, firstLine is floor(startY + 1) and lastLine is floor(endY).
	// These are refactored to integer casts in order to avoid function
	// calls. The difference with integer cast is that numbers are always
	// rounded towards zero. Since values smaller than zero get clipped away,
	// only coordinates between 0 and -1 require greater attention as they
	// also round to zero. The problems in this range can be avoided by
	// adding one to the values before conversion and subtracting after it.

	firstLine := int(math.Floor(startY)) + 1
	lastLine := int(math.Floor(endY))

	minClip := int(clipBound[1])
	maxClip := int(clipBound[3])

	// If start and end are on the same line, the edge doesn't cross
	// any lines and thus can be ignored.
	// If the end is smaller than the first line, edge is out.
	// If the start is larger than the last line, edge is out.
	if firstLine > lastLine || lastLine < minClip || firstLine >= maxClip {
		return false
	}

	// Adjust the start based on the target.
	if firstLine < minClip {
		firstLine = minClip
	}

	if lastLine >= maxClip {
		lastLine = maxClip - 1
	}
	edge.Slope = (endX - startX) / (endY - startY)
	edge.X = startX + (float64(firstLine)-startY)*edge.Slope
	edge.Winding = winding
	edge.FirstLine = firstLine
	edge.LastLine = lastLine

	return true
}


//! Creates a vertical polygon edge between two y values.
/*! Clips the edge vertically to the clip rectangle. Returns true for edges that
 *  should be rendered, false for others.
 */
func getVerticalEdge(startY, endY, x float64, edge *PolygonEdge, clipBound [4]float64) bool {
	var start, end float64
	var winding int16
	if startY < endY {
		start = startY
		end = endY
		winding = 1
	} else {
		start = endY
		end = startY
		winding = -1
	}

	firstLine := int(math.Floor(start)) + 1
	lastLine := int(math.Floor(end))

	minClip := int(clipBound[1])
	maxClip := int(clipBound[3])

	// If start and end are on the same line, the edge doesn't cross
	// any lines and thus can be ignored.
	// If the end is smaller than the first line, edge is out.
	// If the start is larger than the last line, edge is out.
	if firstLine > lastLine || lastLine < minClip || firstLine >= maxClip {
		return false
	}

	// Adjust the start based on the clip rect.
	if firstLine < minClip {
		firstLine = minClip
	}
	if lastLine >= maxClip {
		lastLine = maxClip - 1
	}

	edge.Slope = 0
	edge.X = x
	edge.Winding = winding
	edge.FirstLine = firstLine
	edge.LastLine = lastLine

	return true
}
