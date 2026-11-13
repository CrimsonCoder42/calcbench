package stream

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

func addAll(t *testing.T, s *RunningStats, values []float64) {
	t.Helper()
	for _, v := range values {
		if err := s.Add(v); err != nil {
			t.Fatal(err)
		}
	}
}

// twoPass is the oracle: mean, then the sum of squared deviations from it.
func twoPass(values []float64) (mean, m2 float64) {
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))
	for _, v := range values {
		m2 += (v - mean) * (v - mean)
	}
	return mean, m2
}

// Hand-derived: 2 4 4 4 5 5 7 9 has mean 5, squared deviations 9 1 1 1 0 0 4 16 summing to 32,
// population variance 32/8 = 4, sample variance 32/7, min 2, max 9.
func TestRunningStats_example(t *testing.T) {
	var s RunningStats
	addAll(t, &s, []float64{2, 4, 4, 4, 5, 5, 7, 9})
	mean, _ := s.Mean()
	pv, _ := s.PopulationVariance()
	sv, _ := s.SampleVariance()
	lo, _ := s.Min()
	hi, _ := s.Max()
	// Tolerance 1e-15 relative: every intermediate is a small dyadic or a quotient by k ≤ 8.
	if s.Count() != 8 || !closeTo(mean, 5, 1e-15) || !closeTo(pv, 4, 1e-15) || !closeTo(sv, 32.0/7, 1e-15) || lo != 2 || hi != 9 {
		t.Errorf("got n=%d mean=%v pv=%v sv=%v min=%v max=%v", s.Count(), mean, pv, sv, lo, hi)
	}
}

// Large offset, small spread: 1e9 + 4, 7, 13, 16 has deviations −6 −3 3 6, so sample variance
// 90/3 = 30 (hand-derived). A one-pass sum-of-squares formula loses every digit here; Welford does not.
func TestRunningStats_largeOffset(t *testing.T) {
	var s RunningStats
	addAll(t, &s, []float64{1e9 + 4, 1e9 + 7, 1e9 + 13, 1e9 + 16})
	// Tolerance 1e-6: the values are exact; the mean's ulp at 1e9 is 1.2e-7, which bounds the error.
	if sv, ok := s.SampleVariance(); !ok || math.Abs(sv-30) > 1e-6 {
		t.Errorf("SampleVariance = (%v, %v), want 30", sv, ok)
	}
}

func TestRunningStats_emptyAndOne(t *testing.T) {
	var s RunningStats
	for name, f := range map[string]func() (float64, bool){"Mean": s.Mean, "PopulationVariance": s.PopulationVariance,
		"SampleVariance": s.SampleVariance, "Min": s.Min, "Max": s.Max} {
		if v, ok := f(); ok || v != 0 {
			t.Errorf("empty %s = (%v, %v), want (0, false)", name, v, ok)
		}
	}
	addAll(t, &s, []float64{-3.5})
	if pv, ok := s.PopulationVariance(); !ok || pv != 0 {
		t.Errorf("one value: PopulationVariance = (%v, %v), want (0, true)", pv, ok)
	}
	if _, ok := s.SampleVariance(); ok {
		t.Error("one value: SampleVariance defined, want (0, false)")
	}
	if m, _ := s.Mean(); m != -3.5 {
		t.Errorf("one value: Mean = %v, want -3.5", m)
	}
}

func TestRunningStats_cases(t *testing.T) {
	tests := []struct {
		name           string
		values         []float64
		mean, pv       float64
		least, largest float64
	}{
		{"duplicates", []float64{3, 3, 3}, 3, 0, 3, 3},
		{"negative", []float64{-1, -3}, -2, 1, -3, -1},
		{"zeros", []float64{0, 0}, 0, 0, 0, 0},
		{"largest finite, same sign", []float64{math.MaxFloat64, math.MaxFloat64}, math.MaxFloat64, 0, math.MaxFloat64, math.MaxFloat64},
		{"smallest subnormal", []float64{5e-324, 5e-324}, 5e-324, 0, 5e-324, 5e-324},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s RunningStats
			addAll(t, &s, tt.values)
			mean, _ := s.Mean()
			pv, _ := s.PopulationVariance()
			lo, _ := s.Min()
			hi, _ := s.Max()
			if mean != tt.mean || pv != tt.pv || lo != tt.least || hi != tt.largest {
				t.Errorf("mean %v pv %v min %v max %v; want %v %v %v %v", mean, pv, lo, hi, tt.mean, tt.pv, tt.least, tt.largest)
			}
		})
	}
}

func TestRunningStats_signedZero(t *testing.T) {
	var s RunningStats
	addAll(t, &s, []float64{0, math.Copysign(0, -1)})
	lo, _ := s.Min()
	hi, _ := s.Max()
	if !math.Signbit(lo) || math.Signbit(hi) {
		t.Errorf("Min, Max = %v (signbit %v), %v (signbit %v); want -0, +0", lo, math.Signbit(lo), hi, math.Signbit(hi))
	}
}

func TestRunningStats_errors(t *testing.T) {
	var s RunningStats
	addAll(t, &s, []float64{1})
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if err := s.Add(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Add(%v) error = %v, want ErrNotFinite", v, err)
		}
	}
	if s.Count() != 1 {
		t.Errorf("a rejected value changed the count: %d", s.Count())
	}

	var big RunningStats
	addAll(t, &big, []float64{math.MaxFloat64})
	before := big
	if err := big.Add(-math.MaxFloat64); !errors.Is(err, ErrOverflow) || big != before {
		t.Errorf("deviation beyond the float64 range: error %v, state changed %v", err, big != before)
	}

	full := RunningStats{n: math.MaxInt64}
	if err := full.Add(1); !errors.Is(err, ErrOverflow) {
		t.Errorf("count beyond MaxInt64: error = %v, want ErrOverflow", err)
	}
	if err := full.Merge(s); !errors.Is(err, ErrOverflow) || full.n != math.MaxInt64 {
		t.Errorf("Merge count beyond MaxInt64: error = %v", err)
	}

	var neg RunningStats
	addAll(t, &neg, []float64{-math.MaxFloat64})
	before = big
	if err := big.Merge(neg); !errors.Is(err, ErrOverflow) || big != before {
		t.Errorf("Merge with means MaxFloat64 apart: error %v, state changed %v", err, big != before)
	}
}

func TestRunningStats_mergeEmpty(t *testing.T) {
	var a, empty RunningStats
	addAll(t, &a, []float64{1, 2, 3})
	before := a
	if err := a.Merge(empty); err != nil || a != before {
		t.Errorf("Merge(empty) changed the summary or failed: %v", err)
	}
	var b RunningStats
	if err := b.Merge(a); err != nil || b != a {
		t.Errorf("empty.Merge(a) = %+v, want %+v (%v)", b, a, err)
	}
}

// Oracle: Welford against two passes; Merge of a random split against adding everything to one summary.
func TestRunningStats_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := 2 + rng.IntN(200)
		values := make([]float64, n)
		for i := range values {
			values[i] = 1e3*rng.NormFloat64() + 50
		}
		var whole, left, right RunningStats
		addAll(t, &whole, values)
		cut := rng.IntN(n + 1)
		addAll(t, &left, values[:cut])
		addAll(t, &right, values[cut:])
		if err := left.Merge(right); err != nil {
			t.Fatal(err)
		}
		mean, m2 := twoPass(values)
		lo, hi := values[0], values[0]
		for _, v := range values {
			lo, hi = min(lo, v), max(hi, v)
		}
		// Tolerance 1e-10 relative: n ≤ 200 values of size ~1e3; both methods are backward stable, with
		// error a small multiple of n·u ≈ 4.4e-14 relative to the spread (Chan, Golub, LeVeque 1983, Table 1).
		for name, s := range map[string]*RunningStats{"Welford": &whole, "merged": &left} {
			gm, _ := s.Mean()
			gv, _ := s.PopulationVariance()
			gs, _ := s.SampleVariance()
			gl, _ := s.Min()
			gh, _ := s.Max()
			if s.Count() != int64(n) || !closeTo(gm, mean, 1e-10) || !closeTo(gv, m2/float64(n), 1e-10) ||
				!closeTo(gs, m2/float64(n-1), 1e-10) || gl != lo || gh != hi {
				t.Fatalf("%s n=%d cut=%d: mean %v var %v; want %v %v", name, n, cut, gm, gv, mean, m2/float64(n))
			}
		}
	}
}

func TestRunningStats_noAllocation(t *testing.T) {
	values := benchFloats(1000)
	var s, other RunningStats
	addAll(t, &other, values[:10])
	allocs := testing.AllocsPerRun(20, func() {
		for _, v := range values {
			_ = s.Add(v)
		}
		_ = s.Merge(other)
	})
	if allocs != 0 {
		t.Errorf("Add and Merge allocate %v times, want 0", allocs)
	}
}

func BenchmarkRunningStats(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			values := benchFloats(size.n)
			b.ReportAllocs()
			for b.Loop() {
				var s RunningStats
				for _, v := range values {
					_ = s.Add(v)
				}
			}
		})
	}
}

// Merge is O(1): one case.
func BenchmarkRunningStatsMerge(b *testing.B) {
	values := benchFloats(1000)
	var part RunningStats
	for _, v := range values {
		_ = part.Add(v)
	}
	b.ReportAllocs()
	for b.Loop() {
		var s RunningStats
		_ = s.Merge(part)
		_ = s.Merge(part)
	}
}
