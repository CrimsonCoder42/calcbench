package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// ln C(100,50) = 66.783841652017 (ln of the exact C(100,50), hand-checked with exact integers), and the
// accuracy case of the catalog: ln C(10⁹, 3) = 60.378038038611, where three Lgamma values differ by
// 1.5e-6. Tolerance 1e-13 relative: both are sums of a few correctly rounded logs.
func TestLogBinomial_sourceExample(t *testing.T) {
	for _, c := range []struct {
		n, k int
		want float64
	}{{100, 50, 66.783841652017429}, {1_000_000_000, 3, 60.378038038611}} {
		got, err := LogBinomial(c.n, c.k)
		if err != nil || math.Abs(got-c.want) > 1e-13*c.want {
			t.Errorf("LogBinomial(%d, %d) = (%v, %v), want (%v, nil)", c.n, c.k, got, err, c.want)
		}
	}
}

func TestLogBinomial_cases(t *testing.T) {
	tests := []struct {
		name string
		n, k int
		want float64
	}{
		{"zero of zero", 0, 0, 0},
		{"k zero", 9, 0, 0},
		{"k equals n", 9, 9, 0},
		{"five choose two", 5, 2, math.Log(10)},
		{"max int choose 1", math.MaxInt, 1, math.Log(math.MaxInt)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := LogBinomial(tt.n, tt.k); err != nil || math.Abs(got-tt.want) > 1e-15*math.Max(1, tt.want) {
				t.Errorf("LogBinomial(%d, %d) = (%v, %v), want (%v, nil)", tt.n, tt.k, got, err, tt.want)
			}
		})
	}
	// Huge arguments stay finite.
	if got, err := LogBinomial(math.MaxInt, math.MaxInt/2); err != nil || math.IsInf(got, 0) || math.IsNaN(got) || got <= 0 {
		t.Errorf("LogBinomial(MaxInt, MaxInt/2) = (%v, %v), want finite positive", got, err)
	}
}

func TestLogBinomial_symmetric(t *testing.T) {
	for _, c := range [][2]int{{200, 64}, {200, 65}, {1000, 3}, {5000, 2400}} {
		a, _ := LogBinomial(c[0], c[1])
		b, _ := LogBinomial(c[0], c[0]-c[1])
		if a != b {
			t.Errorf("LogBinomial(%d, %d) = %v, LogBinomial(%d, %d) = %v, want equal", c[0], c[1], a, c[0], c[0]-c[1], b)
		}
	}
}

func TestLogBinomial_errors(t *testing.T) {
	for _, c := range [][2]int{{-1, 0}, {0, -1}, {math.MinInt, 0}} {
		if got, err := LogBinomial(c[0], c[1]); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("LogBinomial(%d, %d) = (%v, %v), want (0, ErrNegative)", c[0], c[1], got, err)
		}
	}
	for _, c := range [][2]int{{5, 6}, {0, 1}, {0, math.MaxInt}} {
		if got, err := LogBinomial(c[0], c[1]); !errors.Is(err, ErrOutOfRange) || got != 0 {
			t.Errorf("LogBinomial(%d, %d) = (%v, %v), want (0, ErrOutOfRange)", c[0], c[1], got, err)
		}
	}
}

// Against ln of the exact C(n,k) in math/big, on both sides of the switch at min(k, n−k) = 64.
// Tolerance 1e-12 relative (at most 64 rounded logs, or the series truncated below 1e-16).
func TestLogBinomial_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		n := 1 + rng.IntN(3000)
		k := rng.IntN(n + 1)
		exact := new(big.Int).Binomial(int64(n), int64(k))
		got, err := LogBinomial(n, k)
		if exact.Cmp(big.NewInt(1)) == 0 {
			if err != nil || got != 0 {
				t.Fatalf("trial %d: LogBinomial(%d, %d) = (%v, %v), want 0", trial, n, k, got, err)
			}
			continue
		}
		want := bigLog(exact)
		if err != nil || math.Abs(got-want) > 1e-12*want {
			t.Fatalf("trial %d: LogBinomial(%d, %d) = (%v, %v), want %v", trial, n, k, got, err, want)
		}
	}
}

func TestLogBinomial_noAlloc(t *testing.T) {
	for _, c := range [][2]int{{1000, 30}, {1000, 500}} {
		if a := testing.AllocsPerRun(100, func() { sinkF, _ = LogBinomial(c[0], c[1]) }); a != 0 {
			t.Errorf("LogBinomial(%d, %d) allocates %v times, want 0", c[0], c[1], a)
		}
	}
}

func BenchmarkLogBinomial(b *testing.B) {
	for _, c := range []struct {
		name string
		n, k int
	}{{"n=1000,k=64", 1000, 64}, {"n=1000,k=500", 1000, 500}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = LogBinomial(c.n, c.k)
			}
		})
	}
}
