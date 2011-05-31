package main

import "draw2d.googlecode.com/hg/draw2d/curve"
import "testing"
import __os__ "os"
import __regexp__ "regexp"

var tests = []testing.InternalTest{
	{"curve.TestCubicCurveRec", curve.TestCubicCurveRec},
	{"curve.TestCubicCurve", curve.TestCubicCurve},
	{"curve.TestCubicCurveAdaptiveRec", curve.TestCubicCurveAdaptiveRec},
	{"curve.TestCubicCurveAdaptive", curve.TestCubicCurveAdaptive},
	{"curve.TestCubicCurveParabolic", curve.TestCubicCurveParabolic},
	{"curve.TestQuadCurve", curve.TestQuadCurve},
}

var benchmarks = []testing.InternalBenchmark{	{"curve.BenchmarkCubicCurveRec", curve.BenchmarkCubicCurveRec},
	{"curve.BenchmarkCubicCurve", curve.BenchmarkCubicCurve},
	{"curve.BenchmarkCubicCurveAdaptiveRec", curve.BenchmarkCubicCurveAdaptiveRec},
	{"curve.BenchmarkCubicCurveAdaptive", curve.BenchmarkCubicCurveAdaptive},
	{"curve.BenchmarkCubicCurveParabolic", curve.BenchmarkCubicCurveParabolic},
	{"curve.BenchmarkQuadCurve", curve.BenchmarkQuadCurve},
}

var matchPat string
var matchRe *__regexp__.Regexp

func matchString(pat, str string) (result bool, err __os__.Error) {
	if matchRe == nil || matchPat != pat {
		matchPat = pat
		matchRe, err = __regexp__.Compile(matchPat)
		if err != nil {
			return
		}
	}
	return matchRe.MatchString(str), nil
}

func main() {
	testing.Main(matchString, tests, benchmarks)
}
