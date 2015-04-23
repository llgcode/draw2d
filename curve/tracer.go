package curve

// LineTracer is an interface that help segmenting curve into small lines
type LineTracer interface {
	// AddPoint a point
	AddPoint(x, y float64)
}
