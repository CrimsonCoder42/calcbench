package stream

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchFloats returns n fixed-seed values in [-1000, 1000). Shared by the tests and benchmarks.
func benchFloats(n int) []float64 {
	rng := rand.New(rand.NewPCG(1, 2))
	s := make([]float64, n)
	for i := range s {
		s[i] = 2000*rng.Float64() - 1000
	}
	return s
}

// exactMean is the oracle: the exact rational mean of values, rounded once to float64.
func exactMean(values []float64) float64 {
	var total big.Rat
	for _, v := range values {
		var r big.Rat
		total.Add(&total, r.SetFloat64(v))
	}
	total.Quo(&total, new(big.Rat).SetInt64(int64(len(values))))
	f, _ := total.Float64()
	return f
}

// closeTo compares within tol·max(1, |want|).
func closeTo(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Max(1, math.Abs(want))
}

func newMA(t *testing.T, w int) *MovingAverageStream {
	t.Helper()
	m, err := NewMovingAverageStream(w)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// NIST/SEMATECH e-Handbook §6.4.2.1: the 12 supplier values of §6.4.2 with a moving average of 3.
// The handbook prints the 3-term averages (the window sums 26, 29, … divided by 3); the first two
// values (fewer than 3 seen) are the means of what has arrived, hand-derived.
func TestMovingAverageStream_nistExample(t *testing.T) {
	supplier := []float64{9, 8, 9, 12, 9, 12, 11, 7, 13, 9, 11, 10}
	want := []float64{9, 8.5, 26.0 / 3, 29.0 / 3, 30.0 / 3, 33.0 / 3, 32.0 / 3, 30.0 / 3, 31.0 / 3, 29.0 / 3, 33.0 / 3, 30.0 / 3}
	m := newMA(t, 3)
	for i, v := range supplier {
		got, err := m.Add(v)
		// Tolerance 1e-15: the sums are exact small integers, so only the final division rounds.
		if err != nil || !closeTo(got, want[i], 1e-15) {
			t.Errorf("after supplier[%d] = %v: mean (%v, %v), want %v", i, v, got, err, want[i])
		}
	}
}

func TestMovingAverageStream_cases(t *testing.T) {
	tests := []struct {
		name   string
		w      int
		values []float64
		want   float64
	}{
		{"one value", 3, []float64{4}, 4},
		{"window 1 is the last value", 1, []float64{3, -1, 2}, 2},
		{"negative and zero", 2, []float64{-2, 0, -3}, -1.5},
		{"duplicates", 2, []float64{5, 5, 5}, 5},
		// After 1e16 leaves, an uncompensated sum would have lost the 1s it absorbed.
		{"large value leaves the window", 2, []float64{1e16, 1, 1, 1}, 1},
		{"largest finite value, window 1", 1, []float64{math.MaxFloat64}, math.MaxFloat64},
		{"smallest subnormal", 2, []float64{5e-324, 5e-324}, 5e-324},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMA(t, tt.w)
			var got float64
			var err error
			for _, v := range tt.values {
				if got, err = m.Add(v); err != nil {
					t.Fatal(err)
				}
			}
			if got != tt.want {
				t.Errorf("mean = %v, want %v", got, tt.want)
			}
			if mean, ok := m.Mean(); !ok || mean != got {
				t.Errorf("Mean() = (%v, %v), want (%v, true)", mean, ok, got)
			}
			if m.Len() != min(len(tt.values), tt.w) || m.WindowLen() != tt.w {
				t.Errorf("Len, WindowLen = %d, %d", m.Len(), m.WindowLen())
			}
		})
	}
}

func TestMovingAverageStream_empty(t *testing.T) {
	m := newMA(t, 4)
	if mean, ok := m.Mean(); ok || mean != 0 || m.Len() != 0 {
		t.Errorf("empty: Mean() = (%v, %v), Len %d; want (0, false), 0", mean, ok, m.Len())
	}
}

func TestMovingAverageStream_errors(t *testing.T) {
	for _, w := range []int{0, -1, math.MinInt, maxWindowLen + 1, math.MaxInt} {
		if m, err := NewMovingAverageStream(w); !errors.Is(err, ErrWindowLen) || m != nil {
			t.Errorf("NewMovingAverageStream(%d) = (%v, %v), want (nil, ErrWindowLen)", w, m, err)
		}
	}
	m := newMA(t, 2)
	if _, err := m.Add(1); err != nil {
		t.Fatal(err)
	}
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := m.Add(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Add(%v) error = %v, want ErrNotFinite", v, err)
		}
	}
	if mean, _ := m.Mean(); mean != 1 || m.Len() != 1 {
		t.Errorf("a rejected value changed the state: mean %v, Len %d", mean, m.Len())
	}
	var zero MovingAverageStream
	if _, err := zero.Add(1); !errors.Is(err, ErrUninitialized) {
		t.Errorf("zero value Add error = %v, want ErrUninitialized", err)
	}
	if _, ok := zero.Mean(); ok || zero.Len() != 0 || zero.WindowLen() != 0 {
		t.Error("zero value reports values")
	}
}

func TestMovingAverageStream_overflowLeavesStateUnchanged(t *testing.T) {
	m := newMA(t, 2)
	if _, err := m.Add(math.MaxFloat64); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Add(math.MaxFloat64); !errors.Is(err, ErrOverflow) {
		t.Errorf("sum beyond MaxFloat64: error = %v, want ErrOverflow", err)
	}
	if mean, _ := m.Mean(); mean != math.MaxFloat64 || m.Len() != 1 {
		t.Errorf("after the rejected Add: mean %v, Len %d; want MaxFloat64, 1", mean, m.Len())
	}
	// Once MaxFloat64 has left the window, the same value fits again.
	if _, err := m.Add(-math.MaxFloat64); err != nil {
		t.Fatal(err)
	}
	if got, err := m.Add(1); err != nil || got != (-math.MaxFloat64+1)/2 {
		t.Errorf("mean = (%v, %v), want (-MaxFloat64/2, nil)", got, err)
	}
}

func TestMovingAverageStream_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		w := 1 + rng.IntN(10)
		n := rng.IntN(60)
		m := newMA(t, w)
		var seen []float64
		for range n {
			v := math.Ldexp(rng.Float64()-0.5, rng.IntN(60)-20) // magnitudes over 18 decades
			seen = append(seen, v)
			got, err := m.Add(v)
			if err != nil {
				t.Fatal(err)
			}
			window := seen[max(0, len(seen)-w):]
			absSum := 0.0
			for _, x := range seen[max(0, len(seen)-2*w):] {
				absSum += math.Abs(x)
			}
			// Tolerance: since the last recomputation at most the last 2w values have entered the
			// compensated sum; its error is a small multiple of u² times their absolute sum (Higham
			// §4.3), far below 2⁻⁴⁹ of it, scaled by 1/len for the mean; plus the two final roundings.
			tol := 0x1p-49 * absSum / float64(len(window))
			if want := exactMean(window); math.Abs(got-want) > tol+math.Abs(want)*0x1p-52 {
				t.Fatalf("w=%d, values %v: mean %v, want %v", w, seen, got, want)
			}
		}
	}
}

// The recomputation every windowLen values keeps the error bounded over a long stream.
func TestMovingAverageStream_longStreamNoDrift(t *testing.T) {
	values := benchFloats(100_000)
	m := newMA(t, 7)
	var got float64
	for _, v := range values {
		got, _ = m.Add(v)
	}
	want := exactMean(values[len(values)-7:])
	// Tolerance 1e-12 absolute: values are below 1000, so a few ulps of 7000 bound the window's error.
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("after 1e5 values: mean %v, want %v", got, want)
	}
}

func TestMovingAverageStream_noAllocation(t *testing.T) {
	m := newMA(t, 50)
	values := benchFloats(1000)
	allocs := testing.AllocsPerRun(20, func() {
		for _, v := range values {
			_, _ = m.Add(v)
		}
	})
	if allocs != 0 {
		t.Errorf("Add allocates %v times per 1000 values, want 0", allocs)
	}
}

func BenchmarkMovingAverageStream(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			values := benchFloats(size.n)
			m, _ := NewMovingAverageStream(50)
			b.ReportAllocs()
			for b.Loop() {
				for _, v := range values {
					_, _ = m.Add(v)
				}
			}
		})
	}
}
