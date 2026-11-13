package invest

import (
	"fmt"
	"math"
	"testing"
)

// periodsPerYear 252 is a stated parameter value (trading days per year), not a constant of the code.
func TestAnnualizedVolatility_handWorked(t *testing.T) {
	got, err := AnnualizedVolatility([]float64{0.01, -0.01}, 252)
	if err != nil || !near(got, math.Sqrt(0.0002*252), 1e-15) || !near(got, 0.2244994432, 1e-10) {
		t.Errorf("AnnualizedVolatility = %v, %v; want 0.2244994432", got, err)
	}
}

func TestAnnualizedVolatility_cases(t *testing.T) {
	got, err := AnnualizedVolatility([]float64{0.02, 0.02, 0.02}, 252)
	if err != nil || !near(got, 0, 1e-17) {
		t.Errorf("constant: %v, %v; want 0", got, err)
	}
	// s of 1, 2, 3, 4 is √(5/3); ppy 1 leaves it unchanged.
	got, err = AnnualizedVolatility([]float64{1, 2, 3, 4}, 1)
	if err != nil || !near(got, math.Sqrt(5.0/3), 1e-15) {
		t.Errorf("1..4: %v, %v; want √(5/3)", got, err)
	}
}

func TestAnnualizedVolatility_errors(t *testing.T) {
	_, err := AnnualizedVolatility([]float64{0.01}, 252)
	wantErr(t, "one", err, ErrTooFew)
	_, err = AnnualizedVolatility([]float64{0.01, math.NaN()}, 252)
	wantErr(t, "NaN", err, ErrNotFinite)
	for _, bad := range []float64{0, -4, math.Inf(1)} {
		_, err = AnnualizedVolatility([]float64{0.01, 0.02}, bad)
		wantErr(t, "bad ppy", err, ErrPeriods)
	}
	_, err = AnnualizedVolatility([]float64{-1e300, 1e300}, 1e300)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestAnnualizedVolatility_noAlloc(t *testing.T) {
	r := []float64{0.01, -0.01}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = AnnualizedVolatility(r, 252) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkAnnualizedVolatility(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = AnnualizedVolatility(r, 252)
			}
		})
	}
}
