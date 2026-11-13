package stats

import (
	"errors"
	"math"
	"testing"
)

// Worked example: e-Handbook §1.3.5.3, AUTO83B.DAT: N1 = 249, s1 = 6.41470; N2 = 79, s2 = 6.10771 →
// printed "s_p = 6.34260".
func TestPooledStandardDeviation_eHandbook(t *testing.T) {
	// Tolerance 5e-6 absolute: the printed inputs and result have 5 decimals.
	if got, err := PooledStandardDeviation(249, 6.41470, 79, 6.10771); err != nil || math.Abs(got-6.34260) > 5e-6 {
		t.Errorf("PooledStandardDeviation(AUTO83B) = %v, %v; want 6.34260", got, err)
	}
}

func TestPooledStandardDeviation_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		n1   int
		sd1  float64
		n2   int
		sd2  float64
		want float64
	}{
		{"equal SDs", 10, 3, 20, 3, 3},
		{"one sample of size 1", 1, 0, 3, 2, 2},
		{"zero SDs", 5, 0, 5, 0, 0},
		{"hand-checked", 3, 1, 3, 3, math.Sqrt(5)},
		{"squares would overflow", 2, m, 2, m, m},
		{"largest counts", math.MaxInt, 1, math.MaxInt, 1, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := PooledStandardDeviation(tt.n1, tt.sd1, tt.n2, tt.sd2); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("PooledStandardDeviation(%d, %v, %d, %v) = %v, %v; want %v", tt.n1, tt.sd1, tt.n2, tt.sd2, got, err, tt.want)
			}
		})
	}
}

func TestPooledStandardDeviation_errors(t *testing.T) {
	for _, tt := range []struct {
		name string
		n1   int
		sd1  float64
		n2   int
		sd2  float64
		want error
	}{
		{"n1 zero", 0, 1, 5, 1, ErrTooFew},
		{"n2 negative", 5, 1, -1, 1, ErrTooFew},
		{"total 2", 1, 1, 1, 1, ErrTooFew},
		{"NaN sd1", 3, math.NaN(), 3, 1, ErrNotFinite},
		{"Inf sd2", 3, 1, 3, math.Inf(1), ErrNotFinite},
		{"negative sd1", 3, -1, 3, 1, ErrNegative},
		{"negative sd2", 3, 1, 3, -1, ErrNegative},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PooledStandardDeviation(tt.n1, tt.sd1, tt.n2, tt.sd2); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("PooledStandardDeviation = %v, %v; want 0, %v", got, err, tt.want)
			}
		})
	}
}

func TestPooledStandardDeviation_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = PooledStandardDeviation(249, 6.41470, 79, 6.10771) }); a != 0 {
		t.Errorf("PooledStandardDeviation allocates %v times per call, want 0", a)
	}
}

func BenchmarkPooledStandardDeviation(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchSink, _ = PooledStandardDeviation(249, 6.41470, 79, 6.10771)
	}
}

// Tiny standard deviations: equal sd1 = sd2 = 3·2^-1000 pool to the same value (hand: a weighted RMS of
// equal numbers); with n1 = n2 = 3, sd1 = 3·2^-1000, sd2 = 4·2^-1000 → √((2·9 + 2·16)/4)·2^-1000 =
// √12.5·2^-1000 (hand). No square is formed (math.Hypot), so nothing underflows.
func TestPooledStandardDeviation_tinySD(t *testing.T) {
	// Tolerance 4e-16 relative: Hypot and the two weight roots each round once.
	if got, err := PooledStandardDeviation(3, 0x3p-1000, 5, 0x3p-1000); err != nil || !near(got, 0x3p-1000, 4e-16) {
		t.Errorf("PooledStandardDeviation(3, 3·2^-1000, 5, 3·2^-1000) = %v, %v; want 3·2^-1000", got, err)
	}
	if got, err := PooledStandardDeviation(3, 0x3p-1000, 3, 0x4p-1000); err != nil || !near(got, math.Sqrt(12.5)*0x1p-1000, 4e-16) {
		t.Errorf("PooledStandardDeviation(3, 3·2^-1000, 3, 4·2^-1000) = %v, %v; want √12.5·2^-1000", got, err)
	}
}
