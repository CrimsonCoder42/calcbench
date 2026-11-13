package prob

import (
	"slices"
	"testing"
)

// Egidi et al. 2018 eq. (2.1) (basic normalisation): [1.9, 1.9] → [0.5, 0.5]; [1.5, 2.75] → [11/17, 6/17] (the inverse
// odds 2/3 and 4/11 sum to 34/33). Hand-checked with fractions; tolerance 4 ulps.
func TestNormalizedImpliedProbabilities_sourceExample(t *testing.T) {
	got, err := NormalizedImpliedProbabilities(nil, []float64{1.9, 1.9})
	if err != nil || !slices.Equal(got, []float64{0.5, 0.5}) {
		t.Errorf("NormalizedImpliedProbabilities([1.9 1.9]) = (%v, %v), want [0.5 0.5]", got, err)
	}
	got, err = NormalizedImpliedProbabilities(nil, []float64{1.5, 2.75})
	if err != nil || len(got) != 2 || !within(got[0], 11.0/17, 1e-15) || !within(got[1], 6.0/17, 1e-15) {
		t.Errorf("NormalizedImpliedProbabilities([1.5 2.75]) = (%v, %v), want [11/17 6/17]", got, err)
	}
}

func TestNormalizedImpliedProbabilities_appendsToDst(t *testing.T) {
	got, err := NormalizedImpliedProbabilities([]float64{7}, []float64{2, 2})
	if err != nil || !slices.Equal(got, []float64{7, 0.5, 0.5}) {
		t.Errorf("NormalizedImpliedProbabilities(dst=[7]) = (%v, %v), want [7 0.5 0.5]", got, err)
	}
}

func TestNormalizedImpliedProbabilities_errors(t *testing.T) {
	got, err := NormalizedImpliedProbabilities([]float64{7}, []float64{2})
	requireErr(t, "one outcome", err, ErrOutOfRange)
	if len(got) != 1 {
		t.Errorf("on error dst length = %d, want 1", len(got))
	}
	_, err = NormalizedImpliedProbabilities(nil, []float64{2, 0})
	requireErr(t, "d < 1", err, ErrOutOfRange)
	_, err = NormalizedImpliedProbabilities(nil, []float64{inf, 2})
	requireErr(t, "Inf", err, ErrNotFinite)
}

// Oracle: the results sum to 1 within n ulps and keep the order of the inverse odds.
func TestNormalizedImpliedProbabilities_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		n := 2 + rng.IntN(30)
		odds := make([]float64, n)
		for i := range odds {
			odds[i] = 1 + 50*rng.Float64()
		}
		got, err := NormalizedImpliedProbabilities(nil, odds)
		if err != nil {
			t.Fatal(err)
		}
		var s compensatedSum
		for i, p := range got {
			s.add(p)
			if i > 0 && (odds[i] < odds[i-1]) != (p > got[i-1]) && odds[i] != odds[i-1] {
				t.Fatalf("trial %d: order not kept at %d", trial, i)
			}
		}
		if absDiff(s.total(), 1) > float64(n)*0x1p-52 {
			t.Fatalf("trial %d: sum = %v, want 1", trial, s.total())
		}
	}
}

func TestNormalizedImpliedProbabilities_noAllocWithCapacity(t *testing.T) {
	odds := []float64{1.5, 2.75}
	buf := make([]float64, 0, 2)
	requireNoAlloc(t, "NormalizedImpliedProbabilities", func() { sinkSlice, _ = NormalizedImpliedProbabilities(buf[:0], odds) })
}

func BenchmarkNormalizedImpliedProbabilities(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}} {
		odds := make([]float64, size.n)
		for i := range odds {
			odds[i] = float64(size.n)
		}
		buf := make([]float64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkSlice, _ = NormalizedImpliedProbabilities(buf[:0], odds)
			}
		})
	}
}
