package signal

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

const besselI0At1 = 1.2660658777520082 // A&S Table 9.8: e⁻¹·I₀(1) = 0.4657596077

func TestKaiser_workedExample(t *testing.T) {
	// M = 3: the ends have r = ±1, so w = I₀(0)/I₀(β) = 1/I₀(β); the centre is 1.
	w, err := Kaiser(nil, 3, 1)
	if err != nil {
		t.Fatal(err)
	}
	e := 1 / besselI0At1
	equalFloats(t, "Kaiser(3,1)", w, []float64{e, 1, e}, 1e-15)
}

func TestKaiser_cases(t *testing.T) {
	// β = 0 is the rectangular window (Harris 1978 §V.H).
	w, _ := Kaiser(nil, 6, 0)
	equalFloats(t, "beta=0", w, []float64{1, 1, 1, 1, 1, 1}, 0)
	// I₀ is even: −β gives the same window.
	p, _ := Kaiser(nil, 9, 5)
	n, _ := Kaiser(nil, 9, -5)
	equalFloats(t, "even in beta", n, p, 0)
	// The window rises strictly to the centre and is exactly symmetric.
	for i := 1; i <= 4; i++ {
		if !(p[i] > p[i-1]) {
			t.Fatalf("not rising: %v", p)
		}
	}
	for i := range p {
		if p[i] != p[8-i] {
			t.Fatalf("not symmetric: %v", p)
		}
	}
	w, _ = Kaiser(nil, 0, 3)
	if len(w) != 0 {
		t.Fatalf("m=0: %v", w)
	}
	w, _ = Kaiser(nil, 1, 3)
	equalFloats(t, "m=1", w, []float64{1}, 0)
	// Large β (asymptotic branch of I₀) stays finite and in (0, 1].
	w, err := Kaiser(nil, 11, 700)
	if err != nil {
		t.Fatal(err)
	}
	if w[5] != 1 || w[0] <= 0 || w[0] > 1e-300 {
		t.Fatalf("beta=700: %v", w)
	}
}

func TestKaiser_errors(t *testing.T) {
	if _, err := Kaiser(nil, -1, 1); !errors.Is(err, ErrDimension) {
		t.Fatalf("m<0: %v", err)
	}
	for _, b := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := Kaiser(nil, 4, b); !errors.Is(err, ErrNotFinite) {
			t.Fatalf("beta=%v: %v", b, err)
		}
	}
	if _, err := Kaiser(nil, 4, 800); !errors.Is(err, ErrOverflow) {
		t.Fatalf("beta=800: %v", err)
	}
}

func BenchmarkKaiser(b *testing.B) {
	for _, m := range []int{1000, 10000, 100000} {
		dst := make([]float64, m)
		b.Run(strconv.Itoa(m), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Kaiser(dst, m, 8.6)
			}
		})
	}
}
