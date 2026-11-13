package signal

import (
	"errors"
	"math"
	"math/cmplx"
	"math/rand/v2"
	"strconv"
	"testing"
)

// dftOracle is the direct O(n²) DFT sum, the definition the fast transforms must match (Oppenheim &
// Schafer 3rd ed. §8.2, eq. (8.67)). Angles are reduced with (j·k mod n) so the oracle itself stays exact
// to rounding for every size used here.
func dftOracle(x []complex128, sign float64) []complex128 {
	n := len(x)
	y := make([]complex128, n)
	for k := range n {
		var s complex128
		for j := range n {
			sn, cs := math.Sincos(sign * 2 * math.Pi * float64((j*k)%n) / float64(n))
			s += x[j] * complex(cs, sn)
		}
		y[k] = s
	}
	return y
}

func randComplex(r *rand.Rand, n int) []complex128 {
	x := make([]complex128, n)
	for i := range x {
		x[i] = complex(2*r.Float64()-1, 2*r.Float64()-1)
	}
	return x
}

// maxAbsDiff returns max |a[i] − b[i]|.
func maxAbsDiff(a, b []complex128) float64 {
	d := 0.0
	for i := range a {
		d = math.Max(d, cmplx.Abs(a[i]-b[i]))
	}
	return d
}

func TestFFT_workedExample(t *testing.T) {
	// Direct DFT of [1, 2, 3, 4] by hand: X[k] = Σ x[j]·(−i)^(jk).
	got, err := FFT(nil, []complex128{1, 2, 3, 4})
	if err != nil {
		t.Fatal(err)
	}
	want := []complex128{10, complex(-2, 2), -2, complex(-2, -2)}
	if d := maxAbsDiff(got, want); d > 1e-14 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestFFT_smallCases(t *testing.T) {
	// n = 1 is the identity; an impulse transforms to all ones; a constant to n·δ.
	got, _ := FFT(nil, []complex128{complex(3, -1)})
	if got[0] != complex(3, -1) {
		t.Fatalf("n=1: %v", got)
	}
	for _, n := range []int{2, 3, 5, 8, 12} {
		x := make([]complex128, n)
		x[0] = 1
		got, err := FFT(nil, x)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range got {
			if cmplx.Abs(v-1) > 1e-14 {
				t.Fatalf("impulse n=%d k=%d: %v", n, k, v)
			}
		}
		for i := range x {
			x[i] = 1
		}
		got, _ = FFT(nil, x)
		if cmplx.Abs(got[0]-complex(float64(n), 0)) > 1e-13 {
			t.Fatalf("constant n=%d: X[0] = %v", n, got[0])
		}
		for k := 1; k < n; k++ {
			if cmplx.Abs(got[k]) > 1e-13 {
				t.Fatalf("constant n=%d: X[%d] = %v", n, k, got[k])
			}
		}
	}
}

// TestFFT_oracle compares with the direct DFT for every n from 1 to 64 and a few larger sizes (powers of
// two, primes for Bluestein, composites). Tolerance 1e-12·n: FFT rounding grows like ε·log n·‖x‖, the
// oracle's like ε·n·‖x‖, with |x[j]| ≤ √2.
func TestFFT_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	sizes := []int{}
	for n := 1; n <= 64; n++ {
		sizes = append(sizes, n)
	}
	sizes = append(sizes, 97, 127, 128, 251, 256, 360, 509, 1000, 1009, 1024)
	for _, n := range sizes {
		x := randComplex(r, n)
		got, err := FFT(nil, x)
		if err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}
		want := dftOracle(x, -1)
		if d := maxAbsDiff(got, want); d > 1e-12*float64(n) {
			t.Fatalf("FFT n=%d: max diff %g", n, d)
		}
		gotInv, err := IFFT(nil, x)
		if err != nil {
			t.Fatal(err)
		}
		wantInv := dftOracle(x, +1)
		for k := range wantInv {
			wantInv[k] /= complex(float64(n), 0)
		}
		if d := maxAbsDiff(gotInv, wantInv); d > 1e-12 {
			t.Fatalf("IFFT n=%d: max diff %g", n, d)
		}
	}
}

// TestFFT_bluesteinPrime checks primes (never powers of two, so always Bluestein) against the direct sum.
func TestFFT_bluesteinPrime(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))
	for _, n := range []int{3, 7, 13, 31, 101, 257, 1021, 2003} {
		x := randComplex(r, n)
		got, _ := FFT(nil, x)
		if d := maxAbsDiff(got, dftOracle(x, -1)); d > 1e-12*float64(n) {
			t.Fatalf("prime n=%d: max diff %g", n, d)
		}
	}
}

// TestFFT_roundTrip: IFFT(FFT(x)) = x to 1e-12 (relative to |x| ≤ √2), sizes up to 2¹⁴ + 1.
func TestFFT_roundTrip(t *testing.T) {
	r := rand.New(rand.NewPCG(5, 6))
	for _, n := range []int{1, 2, 6, 17, 100, 1 << 10, 1000, 1<<14 + 1} {
		x := randComplex(r, n)
		y, _ := FFT(nil, x)
		z, err := IFFT(nil, y)
		if err != nil {
			t.Fatal(err)
		}
		if d := maxAbsDiff(x, z); d > 1e-12 {
			t.Fatalf("n=%d: round trip diff %g", n, d)
		}
	}
}

// TestFFT_parseval: Σ|x|² = (1/n)·Σ|X|² (Oppenheim & Schafer 3rd ed. §8.6, Parseval's relation).
func TestFFT_parseval(t *testing.T) {
	r := rand.New(rand.NewPCG(7, 8))
	for _, n := range []int{8, 15, 64, 99, 4096} {
		x := randComplex(r, n)
		y, _ := FFT(nil, x)
		var ex, ey float64
		for i := range x {
			ex += real(x[i])*real(x[i]) + imag(x[i])*imag(x[i])
			ey += real(y[i])*real(y[i]) + imag(y[i])*imag(y[i])
		}
		ey /= float64(n)
		if math.Abs(ex-ey) > 1e-11*ex {
			t.Fatalf("n=%d: %g vs %g", n, ex, ey)
		}
	}
}

func TestFFT_dstAliasAndNoModify(t *testing.T) {
	x := []complex128{1, 2, 3, 4, 5}
	orig := append([]complex128(nil), x...)
	y, _ := FFT(nil, x)
	for i := range x {
		if x[i] != orig[i] {
			t.Fatal("input modified")
		}
	}
	z, _ := FFT(x, x) // in place
	if maxAbsDiff(y, z) > 1e-14 || &z[0] != &x[0] {
		t.Fatal("in-place result differs or not in dst")
	}
}

func TestFFT_errors(t *testing.T) {
	if _, err := FFT(nil, nil); !errors.Is(err, ErrEmpty) {
		t.Fatalf("empty: %v", err)
	}
	if _, err := IFFT(nil, []complex128{}); !errors.Is(err, ErrEmpty) {
		t.Fatalf("empty: %v", err)
	}
	for _, bad := range []complex128{complex(math.NaN(), 0), complex(0, math.Inf(1)), complex(math.Inf(-1), 0)} {
		if _, err := FFT(nil, []complex128{1, bad}); !errors.Is(err, ErrNotFinite) {
			t.Fatalf("FFT %v: %v", bad, err)
		}
		if _, err := IFFT(nil, []complex128{bad, 1, 2}); !errors.Is(err, ErrNotFinite) {
			t.Fatalf("IFFT %v: %v", bad, err)
		}
	}
	big := math.MaxFloat64
	if _, err := FFT(nil, []complex128{complex(big, 0), complex(big, 0)}); !errors.Is(err, ErrOverflow) {
		t.Fatalf("overflow: %v", err)
	}
	if _, err := IFFT(nil, []complex128{complex(big, 0), complex(big, 0), complex(big, 0)}); !errors.Is(err, ErrOverflow) {
		t.Fatalf("overflow: %v", err)
	}
}

func BenchmarkFFT(b *testing.B) {
	for _, n := range []int{1 << 10, 1 << 13, 1 << 17, 1000, 10007, 100000} {
		x := randComplex(rand.New(rand.NewPCG(1, 2)), n)
		dst := make([]complex128, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = FFT(dst, x)
			}
		})
	}
}

func BenchmarkIFFT(b *testing.B) {
	for _, n := range []int{1 << 10, 1 << 13, 1 << 17, 100000} {
		x := randComplex(rand.New(rand.NewPCG(1, 2)), n)
		dst := make([]complex128, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = IFFT(dst, x)
			}
		})
	}
}
