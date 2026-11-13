package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestAnnualizedReturn_cases(t *testing.T) {
	for _, c := range []struct {
		name    string
		returns []float64
		ppy     float64
		want    float64
	}{
		{"two quarters of 10 %, 4 per year", []float64{0.10, 0.10}, 4, math.Pow(1.1, 4) - 1},
		{"12 months of 1 %, 12 per year", []float64{0.01, 0.01, 0.01, 0.01, 0.01, 0.01, 0.01, 0.01, 0.01, 0.01, 0.01, 0.01}, 12, math.Pow(1.01, 12) - 1},
		{"two years of 21 % total, 1 per year", []float64{0.10, 0.10}, 1, 0.10},
		{"zero", []float64{0, 0}, 252, 0},
		{"total loss", []float64{0.1, -1}, 252, -1},
	} {
		got, err := AnnualizedReturn(c.returns, c.ppy)
		if err != nil || !near(got, c.want, 1e-13) {
			t.Errorf("%s: AnnualizedReturn = %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestAnnualizedReturn_errors(t *testing.T) {
	_, err := AnnualizedReturn(nil, 252)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = AnnualizedReturn([]float64{math.NaN()}, 252)
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = AnnualizedReturn([]float64{-2}, 252)
	wantErr(t, "below -1", err, ErrRate)
	for _, bad := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		_, err = AnnualizedReturn([]float64{0.01}, bad)
		wantErr(t, "bad ppy", err, ErrPeriods)
	}
	_, err = AnnualizedReturn([]float64{10}, 1e6)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestAnnualizedReturn_noAlloc(t *testing.T) {
	r := []float64{0.1, 0.1}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = AnnualizedReturn(r, 4) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkAnnualizedReturn(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = AnnualizedReturn(r, 252)
			}
		})
	}
}
