package market

import (
	"math"
	"testing"
)

// Stoikov 2018 §2 gives the formula without a worked number; hand-checked: bid 100, ask 101, Q^b = 300,
// Q^a = 100 → I = 0.75, M_w = (100 + 0.75)·0.01 = 1.0075. Tolerance 1e-15: two roundings.
func TestWeightedMid_handChecked(t *testing.T) {
	got, err := WeightedMid(100, 101, 300, 100, 0.01)
	if err != nil || !near(got, 1.0075, 1e-15) {
		t.Errorf("WeightedMid = %v, %v; want 1.0075", got, err)
	}
}

func TestWeightedMid_cases(t *testing.T) {
	for _, c := range []struct {
		name     string
		bid, ask int64
		qb, qa   float64
		want     float64
	}{
		{"balanced is mid", 100, 102, 5, 5, 101},
		{"only bids → ask", 100, 102, 5, 0, 102},
		{"only asks → bid", 100, 102, 0, 5, 100},
		{"locked", 100, 100, 1, 9, 100},
	} {
		got, err := WeightedMid(c.bid, c.ask, c.qb, c.qa, 1)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestWeightedMid_errors(t *testing.T) {
	_, err := WeightedMid(101, 100, 1, 1, 1)
	wantErr(t, "crossed", err, ErrDomain)
	_, err = WeightedMid(100, 101, 0, 0, 1)
	wantErr(t, "zero sizes", err, ErrDomain)
	_, err = WeightedMid(100, 101, -1, 1, 1)
	wantErr(t, "negative size", err, ErrNegative)
	_, err = WeightedMid(100, 101, math.Inf(1), 1, 1)
	wantErr(t, "inf size", err, ErrNotFinite)
	_, err = WeightedMid(100, 101, 1, 1, -1)
	wantErr(t, "tick < 0", err, ErrOutOfRange)
	_, err = WeightedMid(-5, 101, 1, 1, 1)
	wantErr(t, "negative bid", err, ErrNegative)
	_, err = WeightedMid(1, maxTicks, 1, 1, math.MaxFloat64)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestWeightedMid_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = WeightedMid(100, 101, 300, 100, 0.01) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkWeightedMid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = WeightedMid(100, 101, 300, 100, 0.01)
	}
}
