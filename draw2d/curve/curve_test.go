package curve

import (
	"testing"
	"log"
)


var (
	cf64Test1 = NewCubicCurveFloat64(100, 100, 200, 100, 100, 200, 200, 200)
)

func BenchmarkCubicCurveCasteljauTest1(b *testing.B) {
	var s []float64
	for i := 0; i < b.N; i++ {
        s = cf64Test1.SegmentCasteljau()
    }
	log.Printf("Num of points: %d\n", len(s));
}