package curve

// LineBuilder is an interface that help segmenting curve into small lines
type LineBuilder interface {
	// LineTo a point
	LineTo(x, y float64)
}
