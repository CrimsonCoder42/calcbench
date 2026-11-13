package stream

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

func addPairs(t *testing.T, s *OnlineCovariance, xs, ys []float64) {
	t.Helper()
	for i := range xs {
		if err := s.Add(xs[i], ys[i]); err != nil {
			t.Fatal(err)
		}
	}
}

// Hand-derived: (1,2),(2,4),(3,6),(4,5) has means 2.5, 4.25; deviations x −1.5 −0.5 0.5 1.5,
// y −2.25 −0.25 1.75 0.75; C = 3.375 + 0.125 + 0.875 + 1.125 = 5.5; M2x = 5; M2y = 8.75.
// Covariance 5.5/3, population 5.5/4, beta 5.5/8.75, r = 5.5/√43.75.
func TestOnlineCovariance_example(t *testing.T) {
	var s OnlineCovariance
	addPairs(t, &s, []float64{1, 2, 3, 4}, []float64{2, 4, 6, 5})
	mx, _ := s.MeanX()
	my, _ := s.MeanY()
	c, _ := s.Covariance()
	pc, _ := s.PopulationCovariance()
	b, _ := s.Beta()
	r, _ := s.Correlation()
	// Tolerance 1e-15 relative: dyadic intermediates and one division each.
	if s.Count() != 4 || mx != 2.5 || my != 4.25 || !closeTo(c, 5.5/3, 1e-15) || !closeTo(pc, 5.5/4, 1e-15) ||
		!closeTo(b, 5.5/8.75, 1e-15) || !closeTo(r, 5.5/math.Sqrt(43.75), 1e-15) {
		t.Errorf("got n %d means %v %v cov %v pop %v beta %v r %v", s.Count(), mx, my, c, pc, b, r)
	}
}

// Oracle: Add against two passes, and Merge of a random split against one summary of the
// concatenation.
func TestOnlineCovariance_oracleAndMerge(t *testing.T) {
	rng := rand.New(rand.NewPCG(5, 6))
	for range 300 {
		n := 2 + rng.IntN(200)
		xs, ys := make([]float64, n), make([]float64, n)
		for i := range xs {
			ys[i] = 1e3*rng.NormFloat64() + 50
			xs[i] = -0.7*ys[i] + 300*rng.NormFloat64()
		}
		var whole, left, right OnlineCovariance
		addPairs(t, &whole, xs, ys)
		cut := rng.IntN(n + 1)
		addPairs(t, &left, xs[:cut], ys[:cut])
		addPairs(t, &right, xs[cut:], ys[cut:])
		if err := left.Merge(right); err != nil {
			t.Fatal(err)
		}
		m2x, m2y, c := twoPassCo(xs, ys)
		scale := math.Sqrt(m2x*m2y) / float64(n-1)
		for name, s := range map[string]*OnlineCovariance{"Add": &whole, "Merge": &left} {
			gc, _ := s.Covariance()
			gb, _ := s.Beta()
			gr, _ := s.Correlation()
			// Tolerance 1e-10 of each quantity's scale: n ≤ 200, both methods backward stable
			// (Chan, Golub, LeVeque 1983).
			if s.Count() != int64(n) || math.Abs(gc-c/float64(n-1)) > 1e-10*scale ||
				math.Abs(gb-c/m2y) > 1e-10*math.Sqrt(m2x/m2y) || math.Abs(gr-c/math.Sqrt(m2x*m2y)) > 1e-10 {
				t.Fatalf("%s n=%d cut=%d: cov %v beta %v r %v; want %v %v %v", name, n, cut, gc, gb, gr,
					c/float64(n-1), c/m2y, c/math.Sqrt(m2x*m2y))
			}
		}
	}
}

func TestOnlineCovariance_emptyMergeReset(t *testing.T) {
	var s, empty OnlineCovariance
	for name, f := range map[string]func() (float64, bool){"MeanX": s.MeanX, "MeanY": s.MeanY, "Covariance": s.Covariance,
		"PopulationCovariance": s.PopulationCovariance, "Correlation": s.Correlation, "Beta": s.Beta} {
		if g, ok := f(); ok || g != 0 {
			t.Errorf("empty %s = (%v, %v), want (0, false)", name, g, ok)
		}
	}
	addPairs(t, &s, []float64{1, 2}, []float64{3, 3})
	before := s
	if err := s.Merge(empty); err != nil || s != before {
		t.Errorf("s.Merge(empty) changed s: %+v (%v)", s, err)
	}
	if err := empty.Merge(s); err != nil || empty != s {
		t.Errorf("empty.Merge(s) = %+v, want %+v (%v)", empty, s, err)
	}
	if _, ok := s.Beta(); ok {
		t.Error("Beta defined with no spread in y")
	}
	if _, ok := s.Correlation(); ok {
		t.Error("Correlation defined with no spread in y")
	}
	s.Reset()
	if s.Count() != 0 {
		t.Error("Reset left pairs")
	}
}

func TestOnlineCovariance_errors(t *testing.T) {
	var s OnlineCovariance
	for _, p := range [][2]float64{{math.NaN(), 0}, {0, math.Inf(-1)}} {
		if err := s.Add(p[0], p[1]); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Add(%v) error = %v, want ErrNotFinite", p, err)
		}
	}
	_ = s.Add(-1e308, 0)
	if err := s.Add(1e308, 0); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add(1e308, 0) error = %v, want ErrOverflow", err)
	}
	if s.Count() != 1 {
		t.Errorf("state changed on error: Count %d", s.Count())
	}
	full := OnlineCovariance{n: math.MaxInt64}
	if err := full.Add(1, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add at MaxInt64 error = %v, want ErrOverflow", err)
	}
	if err := full.Merge(s); !errors.Is(err, ErrOverflow) {
		t.Errorf("Merge past MaxInt64 error = %v, want ErrOverflow", err)
	}
	var far OnlineCovariance
	_ = far.Add(1e308, 0)
	if err := s.Merge(far); !errors.Is(err, ErrOverflow) {
		t.Errorf("Merge of means 2e308 apart error = %v, want ErrOverflow", err)
	}
}

func TestOnlineCovariance_noAllocation(t *testing.T) {
	xs, ys := pairSeries(1000)
	var s, other OnlineCovariance
	addPairs(t, &other, xs[:10], ys[:10])
	allocs := testing.AllocsPerRun(20, func() {
		for i := range xs {
			_ = s.Add(xs[i], ys[i])
		}
		_ = s.Merge(other)
	})
	if allocs != 0 {
		t.Errorf("Add and Merge allocate %v times, want 0", allocs)
	}
}

func BenchmarkOnlineCovariance(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			xs, ys := pairSeries(size.n)
			b.ReportAllocs()
			for b.Loop() {
				var s OnlineCovariance
				for i := range xs {
					_ = s.Add(xs[i], ys[i])
				}
			}
		})
	}
}

// Merge is O(1): one case.
func BenchmarkOnlineCovarianceMerge(b *testing.B) {
	xs, ys := pairSeries(1000)
	var part OnlineCovariance
	for i := range xs {
		_ = part.Add(xs[i], ys[i])
	}
	b.ReportAllocs()
	for b.Loop() {
		var s OnlineCovariance
		_ = s.Merge(part)
		_ = s.Merge(part)
	}
}
