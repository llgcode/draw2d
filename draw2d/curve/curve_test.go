package curve

import (
	"testing"
	"log"
)


var (
	cf64Test1 = NewCubicCurveFloat64(100, 100, 200, 100, 100, 200, 200, 200)
	cf64Test2 = NewCubicCurveFloat64(100, 100, 300, 200, 200, 200, 200, 100)
)

func TestCubicCurveCasteljauRecTest1(t *testing.T) {
	var s []float64
	d := cf64Test1.EstimateDistance()
	log.Printf("Distance estimation: %f\n", d)
	numSegments := int(d * 0.25)
	log.Printf("Max segments estimation: %d\n", numSegments)
	s = make([]float64, 0, numSegments)
	s = cf64Test1.SegmentCasteljauRec(s)
	log.Printf("points: %v\n", s)
	log.Printf("Num of points: %d\n", len(s))
}

func TestCubicCurveCasteljauTest1(t *testing.T) {
	var s []float64
	d := cf64Test1.EstimateDistance()
	log.Printf("Distance estimation: %f\n", d)
	numSegments := int(d * 0.25)
	log.Printf("Max segments estimation: %d\n", numSegments)
	s = make([]float64, 0, numSegments)
	s = cf64Test1.SegmentCasteljau(s)
	log.Printf("points: %v\n", s)
	log.Printf("Num of points: %d\n", len(s))
}


func BenchmarkCubicCurveCasteljauRecTest1(b *testing.B) {
	var s []float64
	d := cf64Test1.EstimateDistance()
	log.Printf("Distance estimation: %f\n", d)
	numSegments := int(d * 0.25)
	log.Printf("Max segments estimation: %d\n", numSegments)
	for i := 0; i < b.N; i++ {
		s = make([]float64, 0, numSegments)
		s = cf64Test1.SegmentCasteljauRec(s)
	}
	log.Printf("Num of points: %d\n", len(s))
}

func BenchmarkCubicCurveCasteljauRecTest2(b *testing.B) {
	var s []float64
	d := cf64Test1.EstimateDistance()
	log.Printf("Distance estimation: %f\n", d)
	numSegments := int(d * 0.25)
	log.Printf("Max segments estimation: %d\n", numSegments)
	for i := 0; i < b.N; i++ {
		s = make([]float64, 0, numSegments)
		s = cf64Test2.SegmentCasteljauRec(s)
	}
	log.Printf("Num of points: %d\n", len(s))
}
func BenchmarkCubicCurveCasteljauTest1(b *testing.B) {
	var s []float64
	d := cf64Test1.EstimateDistance()
	log.Printf("Distance estimation: %f\n", d)
	numSegments := int(d * 0.25)
	log.Printf("Max segments estimation: %d\n", numSegments)
	for i := 0; i < b.N; i++ {
		s = make([]float64, 0, numSegments)
		s = cf64Test1.SegmentCasteljau(s)
	}
	log.Printf("Num of points: %d\n", len(s))
}

func BenchmarkCubicCurveCasteljauTest2(b *testing.B) {
	var s []float64
	d := cf64Test1.EstimateDistance()
	log.Printf("Distance estimation: %f\n", d)
	numSegments := int(d * 0.25)
	log.Printf("Max segments estimation: %d\n", numSegments)
	for i := 0; i < b.N; i++ {
		s = make([]float64, 0, numSegments)
		s = cf64Test2.SegmentCasteljau(s)
	}
	log.Printf("Num of points: %d\n", len(s))
}
