package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestCumulativeReturn_cases(t *testing.T) {
	for _, c := range []struct {
		name    string
		returns []float64
		want    float64
	}{
		{"up then down", []float64{0.10, -0.10}, -0.01},
		{"one", []float64{0.05}, 0.05},
		{"zero", []float64{0, 0, 0}, 0},
		{"doubling twice", []float64{1, 1}, 3},
		{"total loss", []float64{0.5, -1, 0.2}, -1},
	} {
		got, err := CumulativeReturn(c.returns)
		if err != nil || !near(got, c.want, 1e-14) {
			t.Errorf("%s: CumulativeReturn = %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

// Oracle: the direct product Π(1 + R) − 1 on fixed-seed series.
func TestCumulativeReturn_oracle(t *testing.T) {
	r := genReturns(300)
	for n := 1; n <= len(r); n += 13 {
		p := 1.0
		for _, x := range r[:n] {
			p *= 1 + x
		}
		got, err := CumulativeReturn(r[:n])
		if err != nil || !near(got, p-1, 1e-13) {
			t.Errorf("n=%d: CumulativeReturn = %v, %v; want %v", n, got, err, p-1)
		}
	}
}

func TestCumulativeReturn_errors(t *testing.T) {
	_, err := CumulativeReturn(nil)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = CumulativeReturn([]float64{0.1, math.NaN()})
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = CumulativeReturn([]float64{math.Inf(1)})
	wantErr(t, "Inf", err, ErrNotFinite)
	_, err = CumulativeReturn([]float64{0.1, -1.5})
	wantErr(t, "below -1", err, ErrRate)
	_, err = CumulativeReturn([]float64{1e300, 1e300, 1e300})
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestCumulativeReturn_noAlloc(t *testing.T) {
	r := []float64{0.1, -0.1}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = CumulativeReturn(r) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkCumulativeReturn(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = CumulativeReturn(r)
			}
		})
	}
}
