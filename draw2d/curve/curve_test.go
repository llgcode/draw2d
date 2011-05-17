package curve

import (
	"testing"
	"log"
)


var (

	cf64Test1 = NewCubicCurveFloat64(0, 0, 20000, 0, 0, 20000, 20000, 20000)
)

func BenchmarkCubicCurveCasteljauFloat64(b *testing.B) {
	var s []float64
	for i := 0; i < b.N; i++ {
        s = cf64Test1.SegmentCasteljau()
    }
	log.Printf("Num of points: %d\n", len(s));
}