package signal

import (
	"errors"
	"math"
	"math/rand/v2"
	"strconv"
	"testing"
)

// convolveOracle is the convolution sum written out (Oppenheim & Schafer 3rd ed. §2.3), mode full.
func convolveOracle(a, v []float64) []float64 {
	out := make([]float64, len(a)+len(v)-1)
	for i := range a {
		for j := range v {
			out[i+j] += a[i] * v[j]
		}
	}
	return out
}

func reversed(v []float64) []float64 {
	r := make([]float64, len(v))
	for i := range v {
		r[len(v)-1-i] = v[i]
	}
	return r
}

func TestConvolve_workedExample(t *testing.T) {
	// Hand-checked: (1 + 2x + 3x²)(x + 0.5x²) = x + 2.5x² + 4x³ + 1.5x⁴.
	a, v := []float64{1, 2, 3}, []float64{0, 1, 0.5}
	full, err := Convolve(nil, a, v, ModeFull)
	if err != nil {
		t.Fatal(err)
	}
	equalFloats(t, "full", full, []float64{0, 1, 2.5, 4, 1.5}, 0)
	same, _ := Convolve(nil, a, v, ModeSame)
	equalFloats(t, "same", same, []float64{1, 2.5, 4}, 0)
	valid, _ := Convolve(nil, a, v, ModeValid)
	equalFloats(t, "valid", valid, []float64{2.5}, 0)
}

func TestCorrelate_workedExample(t *testing.T) {
	// c[k] = Σ a[n+k]·v[n] for k = −2…2 by hand: 0.5, 2, 3.5, 3, 0.
	a, v := []float64{1, 2, 3}, []float64{0, 1, 0.5}
	full, err := Correlate(nil, a, v, ModeFull)
	if err != nil {
		t.Fatal(err)
	}
	equalFloats(t, "full", full, []float64{0.5, 2, 3.5, 3, 0}, 0)
	same, _ := Correlate(nil, a, v, ModeSame)
	equalFloats(t, "same", same, []float64{2, 3.5, 3}, 0)
	valid, _ := Correlate(nil, a, v, ModeValid)
	equalFloats(t, "valid", valid, []float64{3.5}, 0)
}

func TestConvolve_modeLengths(t *testing.T) {
	// same/valid offsets: same starts at ⌊(min−1)/2⌋, valid at min−1 (numpy's convention).
	a := []float64{1, 2, 3, 4}
	got, _ := Convolve(nil, a, []float64{1, 1}, ModeSame)
	equalFloats(t, "same 4x2", got, []float64{1, 3, 5, 7}, 0)
	got, _ = Convolve(nil, a, []float64{1, 1, 1}, ModeSame)
	equalFloats(t, "same 4x3", got, []float64{3, 6, 9, 7}, 0)
	got, _ = Convolve(nil, []float64{1, 1, 1}, a, ModeValid) // commutes
	equalFloats(t, "valid 3x4", got, []float64{6, 9}, 0)
	got, _ = Convolve(nil, []float64{2}, []float64{3}, ModeFull)
	equalFloats(t, "1x1", got, []float64{6}, 0)
	got, _ = Convolve(nil, []float64{1, 2}, []float64{1, 2}, ModeValid)
	equalFloats(t, "equal lengths valid", got, []float64{4}, 0)
}

// TestConvolve_oracle compares both paths (direct below 64, FFT at and above) with the written-out sum
// for every mode. Tolerance 1e-12·√(n·m): FFT rounding is about ε·log L·‖a‖₂‖v‖₂ with |a|,|v| ≤ 1.
func TestConvolve_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	pairs := [][2]int{{1, 1}, {1, 9}, {9, 1}, {5, 8}, {63, 63}, {63, 200}, {64, 64}, {64, 65}, {100, 300}, {517, 129}, {1000, 1000}}
	for _, p := range pairs {
		a, v := randReal(r, p[0]), randReal(r, p[1])
		tol := 1e-12 * math.Sqrt(float64(p[0]*p[1]))
		for _, corr := range []bool{false, true} {
			vv := v
			f := Convolve
			if corr {
				vv = reversed(v)
				f = Correlate
			}
			full := convolveOracle(a, vv)
			short, long := min(p[0], p[1]), max(p[0], p[1])
			for _, mc := range []struct {
				mode        ConvolutionMode
				start, size int
			}{{ModeFull, 0, len(full)}, {ModeSame, (short - 1) / 2, long}, {ModeValid, short - 1, long - short + 1}} {
				got, err := f(nil, a, v, mc.mode)
				if err != nil {
					t.Fatal(err)
				}
				equalFloats(t, "pair "+strconv.Itoa(p[0])+"x"+strconv.Itoa(p[1]), got, full[mc.start:mc.start+mc.size], tol)
			}
		}
	}
}

func TestConvolve_errors(t *testing.T) {
	for _, f := range []func([]float64, []float64, []float64, ConvolutionMode) ([]float64, error){Convolve, Correlate} {
		if _, err := f(nil, nil, []float64{1}, ModeFull); !errors.Is(err, ErrEmpty) {
			t.Fatalf("empty a: %v", err)
		}
		if _, err := f(nil, []float64{1}, []float64{}, ModeFull); !errors.Is(err, ErrEmpty) {
			t.Fatalf("empty v: %v", err)
		}
		if _, err := f(nil, []float64{1}, []float64{1}, ConvolutionMode(7)); !errors.Is(err, ErrMode) {
			t.Fatalf("mode: %v", err)
		}
		if _, err := f(nil, []float64{1}, []float64{1}, ConvolutionMode(-1)); !errors.Is(err, ErrMode) {
			t.Fatalf("mode -1: %v", err)
		}
		for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			if _, err := f(nil, []float64{1, bad}, []float64{1}, ModeFull); !errors.Is(err, ErrNotFinite) {
				t.Fatalf("a %v: %v", bad, err)
			}
			if _, err := f(nil, []float64{1}, []float64{bad, 1}, ModeSame); !errors.Is(err, ErrNotFinite) {
				t.Fatalf("v %v: %v", bad, err)
			}
		}
		big := math.MaxFloat64
		if _, err := f(nil, []float64{big, big}, []float64{2, 2}, ModeFull); !errors.Is(err, ErrOverflow) {
			t.Fatalf("overflow direct: %v", err)
		}
		long := make([]float64, 70)
		for i := range long {
			long[i] = big
		}
		if _, err := f(nil, long, long, ModeValid); !errors.Is(err, ErrOverflow) {
			t.Fatalf("overflow fft: %v", err)
		}
	}
}

func TestConvolve_noAllocDirectWithDst(t *testing.T) {
	a, v := make([]float64, 100), make([]float64, 10)
	dst := make([]float64, 109)
	if n := testing.AllocsPerRun(50, func() { _, _ = Convolve(dst, a, v, ModeFull) }); n != 0 {
		t.Fatalf("allocs = %v", n)
	}
}

// BenchmarkConvolve: a short kernel (direct path, O(n·32)) and an equal-length pair (FFT path,
// O(n log n)), at 1e3/1e4/1e5 and a power of two.
func BenchmarkConvolve(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	for _, n := range []int{1000, 1 << 13, 10000, 100000} {
		a := randReal(r, n)
		k := randReal(r, 32)
		eq := randReal(r, n)
		dst := make([]float64, 2*n)
		b.Run("kernel32/"+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Convolve(dst, a, k, ModeFull)
			}
		})
		b.Run("equal/"+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Convolve(dst, a, eq, ModeFull)
			}
		})
	}
}

func BenchmarkCorrelate(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	for _, n := range []int{1000, 10000, 100000} {
		a, v := randReal(r, n), randReal(r, n)
		dst := make([]float64, 2*n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Correlate(dst, a, v, ModeSame)
			}
		})
	}
}
