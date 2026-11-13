package prob

import (
	"math"
	"testing"
)

// Casella–Berger §3.2 (discrete uniform on 1…N, P = 1/N): a fair die (1, 6), k = 4 → 1/6. Hand-checked.
func TestDiscreteUniformPMF_sourceExample(t *testing.T) {
	if got, err := DiscreteUniformPMF(4, 1, 6); err != nil || got != 1.0/6 {
		t.Errorf("DiscreteUniformPMF(4, 1, 6) = (%v, %v), want (1/6, nil)", got, err)
	}
}

func TestDiscreteUniformPMF_cases(t *testing.T) {
	tests := []struct {
		name    string
		k, a, b int
		want    float64
	}{
		{"one point", 3, 3, 3, 1},
		{"below", 0, 1, 6, 0},
		{"above", 7, 1, 6, 0},
		{"negative range", -2, -3, 0, 0.25},
		{"widest allowed", 5, 1, math.MaxInt, 1 / float64(math.MaxInt)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := DiscreteUniformPMF(tt.k, tt.a, tt.b); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("DiscreteUniformPMF(%d, %d, %d) = (%v, %v), want %v", tt.k, tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

func TestDiscreteUniformPMF_errors(t *testing.T) {
	_, err := DiscreteUniformPMF(1, 6, 1)
	requireErr(t, "a > b", err, ErrOutOfRange)
}

// The full int range has m = 2⁶⁴ points (beyond int and uint64, exact in float64): P = 2⁻⁶⁴. Before the fix
// this was ErrOverflow. Hand-checked: [MinInt, 0] has 2⁶³ + 1 points, rounded in float64 to 2⁶³.
func TestDiscreteUniformPMF_fullIntRange(t *testing.T) {
	for _, c := range []struct {
		k, a, b int
		want    float64
	}{{0, math.MinInt, math.MaxInt, 0x1p-64}, {math.MaxInt, math.MinInt, math.MaxInt, 0x1p-64}, {0, math.MinInt, 0, 0x1p-63}} {
		if got, err := DiscreteUniformPMF(c.k, c.a, c.b); err != nil || got != c.want {
			t.Errorf("DiscreteUniformPMF(%d, %d, %d) = (%v, %v), want (%v, nil)", c.k, c.a, c.b, got, err, c.want)
		}
	}
}

func TestDiscreteUniformPMF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "DiscreteUniformPMF", func() { sinkFloat, _ = DiscreteUniformPMF(4, 1, 6) })
}

func BenchmarkDiscreteUniformPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DiscreteUniformPMF(4, 1, 6)
	}
}
