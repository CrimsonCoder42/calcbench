package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestLogReturns_handWorked(t *testing.T) {
	got, err := LogReturns(nil, []float64{100, 110, 100})
	if err != nil || len(got) != 2 || !near(got[0], 0.09531017980432486, 1e-15) || !near(got[1], -0.09531017980432486, 1e-15) {
		t.Errorf("LogReturns = %v, %v; want [ln 1.1, −ln 1.1]", got, err)
	}
}

// ln(1e300/1e-300) = 600 ln 10: the ratio overflows, the difference of logs does not.
func TestLogReturns_extremeRatio(t *testing.T) {
	got, err := LogReturns(nil, []float64{1e-300, 1e300})
	if err != nil || !near(got[0], 600*math.Ln10, 1e-10) {
		t.Errorf("LogReturns = %v, %v; want %v", got, err, 600*math.Ln10)
	}
}

func TestLogReturns_errors(t *testing.T) {
	_, err := LogReturns(nil, []float64{100})
	wantErr(t, "one price", err, ErrTooFew)
	_, err = LogReturns(nil, []float64{100, math.NaN()})
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = LogReturns(nil, []float64{100, 0})
	wantErr(t, "zero", err, ErrDomain)
}

func TestLogReturns_noAlloc(t *testing.T) {
	prices, dst := []float64{100, 110, 99}, make([]float64, 2)
	if a := testing.AllocsPerRun(100, func() { sinkFloats, sinkErr = LogReturns(dst, prices) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkLogReturns(b *testing.B) {
	for _, n := range benchSizes {
		prices, dst := genPrices(n), make([]float64, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloats, sinkErr = LogReturns(dst, prices)
			}
		})
	}
}
