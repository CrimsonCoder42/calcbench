package stats

import (
	"errors"
	"math"
	"testing"
)

// Hand computations: z₀.₉₇₅ = 1.959963984540054, so x̄ 20.6, σ 3.2, n 50 gives margin
// 1.959963984540054·3.2/√50 = 0.88698… and limits 20.6 ∓ margin; t₀.₉₇₅ with 9 df = 2.262157163 (NIST
// e-Handbook §1.3.6.7.2 t table), x̄ 10, s 2, n 10 → margin 2.262157163·2/√10 = 1.430714…
func TestMeanConfidenceInterval_cases(t *testing.T) {
	const z = 1.959963984540054
	wantM := z * 3.2 / math.Sqrt(50)
	lo, hi, m, err := MeanConfidenceIntervalZ(20.6, 3.2, 50, 0.95)
	if err != nil || math.Abs(m-wantM) > 1e-12 || math.Abs(lo-(20.6-wantM)) > 1e-12 || math.Abs(hi-(20.6+wantM)) > 1e-12 {
		t.Errorf("Z: %v %v %v %v; want margin %v", lo, hi, m, err, wantM)
	}
	_, _, m, err = MeanConfidenceIntervalT(10, 2, 10, 0.95)
	if want := 2.262157163 * 2 / math.Sqrt(10); err != nil || math.Abs(m-want) > 1e-8 {
		t.Errorf("T: margin %v, %v; want %v", m, err, want)
	}
	// sd 0: a degenerate interval at the mean.
	lo, hi, m, err = MeanConfidenceIntervalZ(5, 0, 3, 0.9)
	if err != nil || lo != 5 || hi != 5 || m != 0 {
		t.Errorf("sd 0: %v %v %v %v", lo, hi, m, err)
	}
	// t is wider than z for small n and approaches it for large n.
	_, _, mz, _ := MeanConfidenceIntervalZ(0, 1, 5, 0.95)
	_, _, mt, _ := MeanConfidenceIntervalT(0, 1, 5, 0.95)
	if !(mt > mz) {
		t.Errorf("t margin %v not wider than z %v at n = 5", mt, mz)
	}
}

func TestMeanConfidenceInterval_errors(t *testing.T) {
	e := func(_, _, _ float64, err error) error { return err }
	tests := []struct {
		name      string
		err, want error
	}{
		{"Z mean NaN", e(MeanConfidenceIntervalZ(math.NaN(), 1, 5, 0.95)), ErrNotFinite},
		{"Z sd Inf", e(MeanConfidenceIntervalZ(0, math.Inf(1), 5, 0.95)), ErrNotFinite},
		{"Z sd negative", e(MeanConfidenceIntervalZ(0, -1, 5, 0.95)), ErrNegative},
		{"Z n 0", e(MeanConfidenceIntervalZ(0, 1, 0, 0.95)), ErrOutOfRange},
		{"Z confidence 1", e(MeanConfidenceIntervalZ(0, 1, 5, 1)), ErrOutOfRange},
		{"Z overflow", e(MeanConfidenceIntervalZ(math.MaxFloat64, math.MaxFloat64, 1, 0.95)), ErrOverflow},
		{"T n 1", e(MeanConfidenceIntervalT(0, 1, 1, 0.95)), ErrOutOfRange},
		{"T confidence NaN", e(MeanConfidenceIntervalT(0, 1, 5, math.NaN())), ErrNotFinite},
		{"T sd negative", e(MeanConfidenceIntervalT(0, -1, 5, 0.95)), ErrNegative},
	}
	for _, tc := range tests {
		if !errors.Is(tc.err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, tc.err, tc.want)
		}
	}
}

func BenchmarkMeanConfidenceIntervalZ(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _, _ = MeanConfidenceIntervalZ(20.6, 3.2, 50, 0.95)
	}
}

func BenchmarkMeanConfidenceIntervalT(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _, _ = MeanConfidenceIntervalT(20.6, 3.2, 50, 0.95)
	}
}
