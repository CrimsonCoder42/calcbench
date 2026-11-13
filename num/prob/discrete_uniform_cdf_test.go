package prob

import (
	"math"
	"testing"
)

// Casella–Berger §3.2 (discrete uniform): a fair die, k = 4 → 4/6 ≈ 0.6666666666666666. Hand-checked.
func TestDiscreteUniformCDF_sourceExample(t *testing.T) {
	if got, err := DiscreteUniformCDF(4, 1, 6); err != nil || got != 4.0/6 {
		t.Errorf("DiscreteUniformCDF(4, 1, 6) = (%v, %v), want (4/6, nil)", got, err)
	}
}

func TestDiscreteUniformCDF_cases(t *testing.T) {
	tests := []struct {
		name    string
		k, a, b int
		want    float64
	}{
		{"below", 0, 1, 6, 0},
		{"at a", 1, 1, 6, 1.0 / 6},
		{"at b", 6, 1, 6, 1},
		{"above", math.MaxInt, 1, 6, 1},
		{"negative range", -2, -3, 0, 0.5},
		{"wide range middle", 0, -(1 << 61), 1<<61 - 1, 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := DiscreteUniformCDF(tt.k, tt.a, tt.b); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("DiscreteUniformCDF(%d, %d, %d) = (%v, %v), want %v", tt.k, tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

func TestDiscreteUniformCDF_errors(t *testing.T) {
	_, err := DiscreteUniformCDF(1, 2, 1)
	requireErr(t, "a > b", err, ErrOutOfRange)
}

// Ranges wider than int: P(X ≤ −1) on MinInt…MaxInt is 2⁶³/2⁶⁴ = ½ exactly; P(X ≤ MinInt) = 2⁻⁶⁴. On
// MinInt…0 (m = 2⁶³ + 1), P(X ≤ −1) = 2⁶³/(2⁶³ + 1) = 1 − 1.08e-19, which correctly rounds to 1; k ≥ b is 1.
// Hand-checked. Before the fix these were ErrOverflow.
func TestDiscreteUniformCDF_fullIntRange(t *testing.T) {
	for _, c := range []struct {
		k, a, b int
		want    float64
	}{{-1, math.MinInt, math.MaxInt, 0.5}, {math.MinInt, math.MinInt, math.MaxInt, 0x1p-64}, {-1, math.MinInt, 0, 1}, {0, math.MinInt, 0, 1}} {
		if got, err := DiscreteUniformCDF(c.k, c.a, c.b); err != nil || got != c.want {
			t.Errorf("DiscreteUniformCDF(%d, %d, %d) = (%v, %v), want (%v, nil)", c.k, c.a, c.b, got, err, c.want)
		}
	}
}

// Oracle: the CDF equals the running sum of the PMF for small ranges. Tolerance 1e-15 (exact quotients).
func TestDiscreteUniformCDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 200 {
		a := rng.IntN(100) - 50
		b := a + rng.IntN(60)
		k := a - 2 + rng.IntN(b-a+5)
		count := max(0, min(k, b)-a+1)
		want := float64(count) / float64(b-a+1)
		if got, err := DiscreteUniformCDF(k, a, b); err != nil || !within(got, want, 1e-15) {
			t.Fatalf("trial %d: DiscreteUniformCDF(%d, %d, %d) = (%v, %v), want %v", trial, k, a, b, got, err, want)
		}
	}
}

func TestDiscreteUniformCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "DiscreteUniformCDF", func() { sinkFloat, _ = DiscreteUniformCDF(4, 1, 6) })
}

func BenchmarkDiscreteUniformCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DiscreteUniformCDF(4, 1, 6)
	}
}
