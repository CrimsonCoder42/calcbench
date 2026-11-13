package market

import (
	"math"
	"testing"
)

// Hand-checked from Stoikov 2018 §2's definition: Q^b = 300, Q^a = 100 → 0.75 (exact in binary).
func TestLevel1Imbalance_cases(t *testing.T) {
	for _, c := range []struct {
		name         string
		qb, qa, want float64
	}{
		{"example", 300, 100, 0.75},
		{"balanced", 4, 4, 0.5},
		{"only bid", 1, 0, 1},
		{"only ask", 0, 1, 0},
		{"tiny", 5e-324, 5e-324, 0.5},
	} {
		got, err := Level1Imbalance(c.qb, c.qa)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestLevel1Imbalance_errors(t *testing.T) {
	_, err := Level1Imbalance(0, 0)
	wantErr(t, "both zero", err, ErrDomain)
	_, err = Level1Imbalance(-1, 2)
	wantErr(t, "negative", err, ErrNegative)
	_, err = Level1Imbalance(math.NaN(), 2)
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = Level1Imbalance(math.MaxFloat64, math.MaxFloat64)
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkLevel1Imbalance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = Level1Imbalance(300, 100)
	}
}
