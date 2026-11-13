package bloom

import (
	"errors"
	"math"
	"testing"
)

// Broder–Mitzenmacher §2.1–2.2, hand-computed: n = 1000, p = 0.01: −n·ln p/(ln 2)² = 9585.06 → m = 9586;
// (m/n)·ln 2 = 6.6446 → k = 7.
func TestOptimalParams_example(t *testing.T) {
	m, k, err := OptimalParams(1000, 0.01)
	if err != nil || m != 9586 || k != 7 {
		t.Errorf("OptimalParams(1000, 0.01) = (%d, %d, %v), want (9586, 7, nil)", m, k, err)
	}
}

// For each n and p: m is the least integer at or above the bound, k is the rounded optimum, and the
// predicted rate is within 1% of p (a Python sweep of the same grid found at most +0.72%, at p = 0.1).
func TestOptimalParams_grid(t *testing.T) {
	for _, p := range []float64{0.5, 0.1, 0.01, 1e-3, 1e-6, 1e-9, 1e-300} {
		for _, n := range []uint64{1, 2, 3, 10, 1000, 1_000_000} {
			m, k, err := OptimalParams(n, p)
			if errors.Is(err, ErrTooLarge) {
				continue // 1e6 items at 1e-300 need more than MaxBits
			}
			if err != nil {
				t.Fatalf("OptimalParams(%d, %v): %v", n, p, err)
			}
			bound := -float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)
			if float64(m) < bound || float64(m-1) >= bound {
				t.Errorf("OptimalParams(%d, %v): m = %d, bound %v", n, p, m, bound)
			}
			if want := max(1, int(math.Round(float64(m)/float64(n)*math.Ln2))); k != want {
				t.Errorf("OptimalParams(%d, %v): k = %d, want %d", n, p, k, want)
			}
			if f, _ := FalsePositiveRate(m, n, k); f > 1.01*p {
				t.Errorf("OptimalParams(%d, %v) = (%d, %d): predicted rate %v > 1.01·p", n, p, m, k, f)
			}
		}
	}
}

func TestOptimalParams_errors(t *testing.T) {
	tests := []struct {
		name    string
		n       uint64
		p       float64
		wantErr error
	}{
		{"no items", 0, 0.01, ErrItems},
		{"p zero", 10, 0, ErrProbability},
		{"p one", 10, 1, ErrProbability},
		{"p negative", 10, -0.1, ErrProbability},
		{"p above one", 10, 1.5, ErrProbability},
		{"p NaN", 10, math.NaN(), ErrProbability},
		{"p +Inf", 10, math.Inf(1), ErrProbability},
		{"p -Inf", 10, math.Inf(-1), ErrProbability},
		{"too many bits", 1_000_000_000, 0.01, ErrTooLarge},
		{"maximum of uint64", math.MaxUint64, 0.5, ErrTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, k, err := OptimalParams(tt.n, tt.p)
			if m != 0 || k != 0 || !errors.Is(err, tt.wantErr) {
				t.Errorf("OptimalParams(%d, %v) = (%d, %d, %v), want (0, 0, %v)", tt.n, tt.p, m, k, err, tt.wantErr)
			}
		})
	}
}

func BenchmarkOptimalParams(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			OptimalParams(1000, 0.01)
		}
	})
}
