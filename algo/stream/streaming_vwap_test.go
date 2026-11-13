package stream

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// tradeSeries returns n fixed-seed trades: price 100 + N(0, 1), integer volume in [0, 1000].
func tradeSeries(n int) (prices, volumes []float64) {
	rng := rand.New(rand.NewPCG(7, 8))
	prices, volumes = make([]float64, n), make([]float64, n)
	for i := range prices {
		prices[i] = 100 + rng.NormFloat64()
		volumes[i] = float64(rng.IntN(1001))
	}
	return prices, volumes
}

// directVWAP is the oracle: Σ p·v / Σ v in one plain pass.
func directVWAP(prices, volumes []float64) (float64, bool) {
	var pv, v float64
	for i := range prices {
		pv += prices[i] * volumes[i]
		v += volumes[i]
	}
	if v == 0 {
		return 0, false
	}
	return pv / v, true
}

// Hand-derived: (10, 100), (11, 300), (9.5, 0): (1000 + 3300)/400 = 10.75.
func TestStreamingVWAP_example(t *testing.T) {
	var s StreamingVWAP
	if _, ok := s.VWAP(); ok {
		t.Error("empty: VWAP defined")
	}
	for _, tr := range [][2]float64{{10, 100}, {11, 300}, {9.5, 0}} {
		if err := s.Add(tr[0], tr[1]); err != nil {
			t.Fatal(err)
		}
	}
	if p, ok := s.VWAP(); !ok || p != 10.75 || s.Count() != 3 || s.TotalVolume() != 400 {
		t.Errorf("VWAP (%v, %v) count %d volume %v, want 10.75 3 400", p, ok, s.Count(), s.TotalVolume())
	}
	s.Reset()
	_ = s.Add(5, 0)
	if _, ok := s.VWAP(); ok {
		t.Error("only zero-volume trades: VWAP defined")
	}
}

// Merge of a random split equals one summary of the concatenation, and both match the direct sum.
func TestStreamingVWAP_mergeOracle(t *testing.T) {
	prices, volumes := tradeSeries(1000)
	rng := rand.New(rand.NewPCG(9, 10))
	for range 50 {
		cut := rng.IntN(1001)
		var whole, left, right StreamingVWAP
		for i := range prices {
			_ = whole.Add(prices[i], volumes[i])
			if i < cut {
				_ = left.Add(prices[i], volumes[i])
			} else {
				_ = right.Add(prices[i], volumes[i])
			}
		}
		if err := left.Merge(right); err != nil {
			t.Fatal(err)
		}
		want, _ := directVWAP(prices, volumes)
		gw, _ := whole.VWAP()
		gl, _ := left.VWAP()
		// Tolerance 1e-13 relative: compensated sums are within a few ulps; the plain oracle sum of 1,000
		// positive terms is within 1000·u ≈ 1.1e-13 relative.
		if !relClose(gw, want, 1e-13) || !relClose(gl, gw, 1e-14) || left.Count() != 1000 || left.TotalVolume() != whole.TotalVolume() {
			t.Fatalf("cut %d: whole %v merged %v, want %v", cut, gw, gl, want)
		}
	}
}

func TestStreamingVWAP_errors(t *testing.T) {
	var s StreamingVWAP
	for _, tr := range [][2]float64{{math.NaN(), 1}, {1, math.Inf(1)}} {
		if err := s.Add(tr[0], tr[1]); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Add(%v) error = %v, want ErrNotFinite", tr, err)
		}
	}
	if err := s.Add(10, -1); !errors.Is(err, ErrVolume) {
		t.Errorf("Add(10, -1) error = %v, want ErrVolume", err)
	}
	if err := s.Add(1e308, 10); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add(1e308, 10) error = %v, want ErrOverflow", err)
	}
	if s.Count() != 0 || s.TotalVolume() != 0 {
		t.Error("state changed on error")
	}
	full := StreamingVWAP{n: math.MaxInt64}
	if err := full.Add(1, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add at MaxInt64 error = %v, want ErrOverflow", err)
	}
	_ = s.Add(1, 1)
	if err := full.Merge(s); !errors.Is(err, ErrOverflow) {
		t.Errorf("Merge past MaxInt64 count error = %v, want ErrOverflow", err)
	}
	var a, b StreamingVWAP
	_ = a.Add(1e308, 1.5)
	_ = b.Add(1e308, 1.5)
	if err := a.Merge(b); !errors.Is(err, ErrOverflow) {
		t.Errorf("Merge of sums 3e308 error = %v, want ErrOverflow", err)
	}
}

func TestStreamingVWAP_noAllocation(t *testing.T) {
	prices, volumes := tradeSeries(1000)
	var s, other StreamingVWAP
	_ = other.Add(1, 1)
	allocs := testing.AllocsPerRun(20, func() {
		for i := range prices {
			_ = s.Add(prices[i], volumes[i])
		}
		_ = s.Merge(other)
	})
	if allocs != 0 {
		t.Errorf("Add and Merge allocate %v times, want 0", allocs)
	}
}

func BenchmarkStreamingVWAP(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			prices, volumes := tradeSeries(size.n)
			b.ReportAllocs()
			for b.Loop() {
				var s StreamingVWAP
				for i := range prices {
					_ = s.Add(prices[i], volumes[i])
				}
			}
		})
	}
}
