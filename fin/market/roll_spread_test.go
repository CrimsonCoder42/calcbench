package market

import (
	"math"
	"testing"
)

// Roll 1984 §I: s = 2√(−cov). Hand-checked: 100, 101, 100, 101, 100 → Δp = 1, −1, 1, −1; pairs
// (−1, 1), (1, −1), (−1, 1); sample cov = −(24/9)/2 = −4/3; s = 2√(4/3) = 4/√3. Tolerance 1e-14.
func TestRollSpread_handChecked(t *testing.T) {
	got, err := RollSpread([]int64{100, 101, 100, 101, 100}, 1)
	if err != nil || !near(got, 4/math.Sqrt(3), 1e-14) {
		t.Errorf("RollSpread = %v, %v; want %v", got, err, 4/math.Sqrt(3))
	}
	// tick size scales the result
	got, err = RollSpread([]int64{100, 101, 100, 101, 100}, 0.01)
	if err != nil || !near(got, 0.04/math.Sqrt(3), 1e-16) {
		t.Errorf("tick 0.01: %v, %v", got, err)
	}
}

func TestRollSpread_errors(t *testing.T) {
	_, err := RollSpread([]int64{1, 2, 1}, 1)
	wantErr(t, "three prices", err, ErrTooFew)
	_, err = RollSpread([]int64{1, 2, 3, 4, 5}, 1)
	wantErr(t, "trend: cov = 0", err, ErrDomain)
	_, err = RollSpread([]int64{1, 2, 4, 7, 11}, 1)
	wantErr(t, "positive cov", err, ErrDomain)
	_, err = RollSpread([]int64{1, 2, -1, 2}, 1)
	wantErr(t, "negative price", err, ErrNegative)
	_, err = RollSpread([]int64{1, 2, 1, 2}, math.Inf(1))
	wantErr(t, "tick Inf", err, ErrNotFinite)
}

func BenchmarkRollSpread(b *testing.B) {
	p := make([]int64, 1000)
	for i := range p {
		p[i] = int64(100 + i%2 + (i/7)%3)
	}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = RollSpread(p, 0.01)
	}
}
