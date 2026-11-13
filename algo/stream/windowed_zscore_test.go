package stream

import (
	"errors"
	"math"
	"testing"
)

// Hand-derived, window 3. {1,2,3}: mean 2, s = 1, z(3) = 1. {2,3,4}: z(4) = 1.
// {3,4,10}: mean 17/3, s = √(129/9), z(10) = (13/3)/√(129/9) = 13/√129.
func TestWindowedZScore_example(t *testing.T) {
	z, err := NewWindowedZScore(3)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := z.ZScore(); ok {
		t.Error("empty: ZScore defined")
	}
	_ = z.Add(1)
	if _, ok := z.ZScore(); ok {
		t.Error("one value: ZScore defined")
	}
	want := map[int]float64{2: 1, 3: 1, 4: 13 / math.Sqrt(129)}
	for i, x := range []float64{2, 3, 4, 10} {
		if err := z.Add(x); err != nil {
			t.Fatal(err)
		}
		g, ok := z.ZScore()
		w, has := want[i+1]
		// Tolerance 1e-14: a handful of operations on small integers.
		if has && (!ok || !closeTo(g, w, 1e-14)) {
			t.Errorf("after %v: ZScore = (%v, %v), want %v", x, g, ok, w)
		}
	}
	if m, _ := z.Mean(); !closeTo(m, 17.0/3, 1e-14) {
		t.Errorf("Mean = %v, want 17/3", m)
	}
	if s, _ := z.Std(); !closeTo(s, math.Sqrt(129.0/9), 1e-14) {
		t.Errorf("Std = %v, want √(129/9)", s)
	}
	if z.Len() != 3 || z.WindowLen() != 3 {
		t.Errorf("Len, WindowLen = %d, %d", z.Len(), z.WindowLen())
	}
	z.Resync()
	if g, _ := z.ZScore(); !closeTo(g, 13/math.Sqrt(129), 1e-14) {
		t.Errorf("after Resync: ZScore = %v", g)
	}
	z.Reset()
	if z.Len() != 0 {
		t.Error("Reset left values")
	}
}

// A window with no spread has no z-score.
func TestWindowedZScore_noSpread(t *testing.T) {
	z, _ := NewWindowedZScore(4)
	for range 5 {
		_ = z.Add(9.5)
	}
	if g, ok := z.ZScore(); ok {
		t.Errorf("constant window: ZScore = %v, want (0, false)", g)
	}
}

// Oracle: the z of the newest value against two passes over the window, 1,000 fixed-seed values.
func TestWindowedZScore_oracle(t *testing.T) {
	values := normalSeries(1000)
	z, _ := NewWindowedZScore(20)
	for i, x := range values {
		_ = z.Add(x)
		win := values[max(0, i-19) : i+1]
		if len(win) < 2 {
			continue
		}
		mean, m2 := twoPass(win)
		want := (x - mean) / math.Sqrt(m2/float64(len(win)-1))
		g, ok := z.ZScore()
		// Tolerance 1e-9 absolute: |z| ≤ √(n−1) ≈ 4.4, and mean and s match to ~1e-12 relative.
		if !ok || math.Abs(g-want) > 1e-9 {
			t.Fatalf("step %d: z = (%v, %v), want %v", i, g, ok, want)
		}
	}
}

func TestWindowedZScore_errors(t *testing.T) {
	if _, err := NewWindowedZScore(0); !errors.Is(err, ErrWindowLen) {
		t.Errorf("NewWindowedZScore(0) error = %v, want ErrWindowLen", err)
	}
	var zero WindowedZScore
	if err := zero.Add(1); !errors.Is(err, ErrUninitialized) {
		t.Errorf("zero.Add error = %v, want ErrUninitialized", err)
	}
	if zero.Len() != 0 || zero.WindowLen() != 0 {
		t.Error("zero value: Len or WindowLen not 0")
	}
	z, _ := NewWindowedZScore(3)
	_ = z.Add(1)
	_ = z.Add(3)
	if err := z.Add(math.NaN()); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Add(NaN) error = %v, want ErrNotFinite", err)
	}
	// Newest is still 3: mean 2, s = √2, z = 1/√2.
	if g, _ := z.ZScore(); !closeTo(g, 1/math.Sqrt2, 1e-15) {
		t.Errorf("after rejected Add: ZScore = %v, want 1/√2", g)
	}
}

func TestWindowedZScore_noAllocation(t *testing.T) {
	z, _ := NewWindowedZScore(50)
	values := benchFloats(1000)
	allocs := testing.AllocsPerRun(20, func() {
		for _, x := range values {
			_ = z.Add(x)
			_, _ = z.ZScore()
		}
	})
	if allocs != 0 {
		t.Errorf("Add and ZScore allocate %v times per 1000 values, want 0", allocs)
	}
}

func BenchmarkWindowedZScore(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			values := benchFloats(size.n)
			z, _ := NewWindowedZScore(50)
			b.ReportAllocs()
			for b.Loop() {
				for _, x := range values {
					_ = z.Add(x)
					_, _ = z.ZScore()
				}
			}
		})
	}
}
