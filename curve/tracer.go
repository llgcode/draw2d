package curve

// LineTracer is an interface that help segmenting curve into small lines
type LineTracer interface {
	LineTo(x, y float64)
}
