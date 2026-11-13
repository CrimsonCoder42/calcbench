package signal

import (
	"errors"
	"strconv"
	"testing"
)

type windowFunc func([]float64, int) ([]float64, error)

// TestWindows_knownPoints: values at n = 0, (M−1)/4, (M−1)/2 follow from Harris 1978 eqs. (20), (27a),
// (30a), (32) with cos 0 = 1, cos π/2 = 0, cos π = −1. Tolerance 1e-15 (a few ulps of the cosines).
func TestWindows_knownPoints(t *testing.T) {
	cases := []struct {
		name string
		f    windowFunc
		want []float64 // M = 5
	}{
		{"Hanning", Hanning, []float64{0, 0.5, 1, 0.5, 0}},
		{"Hamming", Hamming, []float64{0.08, 0.54, 1, 0.54, 0.08}},
		{"Blackman", Blackman, []float64{0, 0.34, 1, 0.34, 0}},
		{"Bartlett", Bartlett, []float64{0, 0.5, 1, 0.5, 0}},
	}
	for _, c := range cases {
		got, err := c.f(nil, 5)
		if err != nil {
			t.Fatal(err)
		}
		equalFloats(t, c.name, got, c.want, 1e-15)
	}
}

// TestWindows_evenLength checks M = 4 against the formulas evaluated by hand: cos(2π/3) = −1/2.
func TestWindows_evenLength(t *testing.T) {
	h, _ := Hanning(nil, 4)
	equalFloats(t, "Hanning(4)", h, []float64{0, 0.75, 0.75, 0}, 1e-15)
	hm, _ := Hamming(nil, 4)
	equalFloats(t, "Hamming(4)", hm, []float64{0.08, 0.77, 0.77, 0.08}, 1e-15)
	// Blackman: 0.42 + 0.25 + 0.08·cos(4π/3) = 0.67 − 0.04 = 0.63.
	b, _ := Blackman(nil, 4)
	equalFloats(t, "Blackman(4)", b, []float64{0, 0.63, 0.63, 0}, 1e-15)
	bt, _ := Bartlett(nil, 4)
	equalFloats(t, "Bartlett(4)", bt, []float64{0, 2.0 / 3, 2.0 / 3, 0}, 1e-15)
}

func TestWindows_smallAndSymmetric(t *testing.T) {
	for _, f := range []windowFunc{Hanning, Hamming, Blackman, Bartlett} {
		w, err := f(nil, 0)
		if err != nil || len(w) != 0 {
			t.Fatalf("m=0: %v %v", w, err)
		}
		w, err = f(nil, 1)
		if err != nil || len(w) != 1 || w[0] != 1 {
			t.Fatalf("m=1: %v %v", w, err)
		}
		for _, m := range []int{2, 3, 10, 101, 1000} {
			w, _ := f(nil, m)
			for i := range w {
				if w[i] != w[m-1-i] || w[i] > 1+1e-15 || w[i] < -1e-15 {
					t.Fatalf("m=%d i=%d: %v", m, i, w[i])
				}
			}
		}
		if _, err := f(nil, -1); !errors.Is(err, ErrDimension) {
			t.Fatalf("m=-1: %v", err)
		}
	}
}

func TestWindows_noAllocWithDst(t *testing.T) {
	dst := make([]float64, 64)
	for _, f := range []windowFunc{Hanning, Hamming, Blackman, Bartlett} {
		if a := testing.AllocsPerRun(50, func() { _, _ = f(dst, 64) }); a != 0 {
			t.Fatalf("allocs = %v", a)
		}
	}
}

func benchWindow(b *testing.B, f windowFunc) {
	for _, m := range []int{1000, 10000, 100000} {
		dst := make([]float64, m)
		b.Run(strconv.Itoa(m), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = f(dst, m)
			}
		})
	}
}

func BenchmarkHanning(b *testing.B)  { benchWindow(b, Hanning) }
func BenchmarkHamming(b *testing.B)  { benchWindow(b, Hamming) }
func BenchmarkBlackman(b *testing.B) { benchWindow(b, Blackman) }
