package signal

import (
	"errors"
	"math"
	"math/rand/v2"
	"strconv"
	"testing"
)

func randReal(r *rand.Rand, n int) []float64 {
	x := make([]float64, n)
	for i := range x {
		x[i] = 2*r.Float64() - 1
	}
	return x
}

func toComplex(x []float64) []complex128 {
	z := make([]complex128, len(x))
	for i, v := range x {
		z[i] = complex(v, 0)
	}
	return z
}

func TestRFFT_workedExample(t *testing.T) {
	// The first n/2+1 bins of the hand-checked DFT of [1, 2, 3, 4].
	got, err := RFFT(nil, []float64{1, 2, 3, 4})
	if err != nil {
		t.Fatal(err)
	}
	want := []complex128{10, complex(-2, 2), -2}
	if len(got) != 3 || maxAbsDiff(got, want) > 1e-14 {
		t.Fatalf("got %v, want %v", got, want)
	}
	back, err := IRFFT(nil, got, 4)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range []float64{1, 2, 3, 4} {
		if math.Abs(back[i]-v) > 1e-14 {
			t.Fatalf("IRFFT got %v", back)
		}
	}
}

// TestRFFT_oracle compares with the direct DFT for n = 1…64 and larger even, odd and power-of-two sizes;
// tolerance 1e-12·n as for FFT. IRFFT(RFFT(x), n) = x to 1e-12.
func TestRFFT_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	sizes := []int{}
	for n := 1; n <= 64; n++ {
		sizes = append(sizes, n)
	}
	sizes = append(sizes, 98, 101, 256, 1000, 1001, 1024, 2*1009)
	for _, n := range sizes {
		x := randReal(r, n)
		got, err := RFFT(nil, x)
		if err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}
		want := dftOracle(toComplex(x), -1)[:n/2+1]
		if d := maxAbsDiff(got, want); d > 1e-12*float64(n) {
			t.Fatalf("RFFT n=%d: max diff %g", n, d)
		}
		back, err := IRFFT(nil, got, n)
		if err != nil {
			t.Fatalf("IRFFT n=%d: %v", n, err)
		}
		for i := range x {
			if math.Abs(back[i]-x[i]) > 1e-12 {
				t.Fatalf("round trip n=%d i=%d: %v vs %v", n, i, back[i], x[i])
			}
		}
	}
}

// TestIRFFT_matchesIFFT: on a Hermitian spectrum IRFFT equals the real part of IFFT of the full spectrum.
func TestIRFFT_matchesIFFT(t *testing.T) {
	r := rand.New(rand.NewPCG(9, 10))
	for _, n := range []int{2, 5, 16, 33} {
		x := randReal(r, n)
		full, _ := FFT(nil, toComplex(x))
		half := append([]complex128(nil), full[:n/2+1]...)
		half[0] += complex(0, 7) // imaginary part of bin 0 is ignored
		got, _ := IRFFT(nil, half, n)
		for i := range x {
			if math.Abs(got[i]-x[i]) > 1e-12 {
				t.Fatalf("n=%d: %v vs %v", n, got, x)
			}
		}
	}
}

func TestRFFT_errors(t *testing.T) {
	if _, err := RFFT(nil, nil); !errors.Is(err, ErrEmpty) {
		t.Fatalf("empty: %v", err)
	}
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := RFFT(nil, []float64{1, 2, bad}); !errors.Is(err, ErrNotFinite) {
			t.Fatalf("%v: %v", bad, err)
		}
	}
	if _, err := RFFT(nil, []float64{math.MaxFloat64, math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Fatalf("overflow: %v", err)
	}
	if _, err := RFFT(nil, []float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Fatalf("overflow odd: %v", err)
	}
}

func TestIRFFT_errors(t *testing.T) {
	for _, n := range []int{0, -1} {
		if _, err := IRFFT(nil, []complex128{1}, n); !errors.Is(err, ErrDimension) {
			t.Fatalf("n=%d: %v", n, err)
		}
	}
	if _, err := IRFFT(nil, []complex128{1, 2}, 4); !errors.Is(err, ErrDimension) {
		t.Fatalf("short spectrum: %v", err)
	}
	if _, err := IRFFT(nil, []complex128{1, complex(math.NaN(), 0)}, 3); !errors.Is(err, ErrNotFinite) {
		t.Fatalf("NaN: %v", err)
	}
	big := complex(math.MaxFloat64, 0)
	if _, err := IRFFT(nil, []complex128{big, big, big}, 4); !errors.Is(err, ErrOverflow) {
		t.Fatalf("overflow: %v", err)
	}
}

func BenchmarkRFFT(b *testing.B) {
	for _, n := range []int{1 << 10, 1 << 13, 1 << 17, 1000, 10000, 100000} {
		x := randReal(rand.New(rand.NewPCG(1, 2)), n)
		dst := make([]complex128, n/2+1)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = RFFT(dst, x)
			}
		})
	}
}

func BenchmarkIRFFT(b *testing.B) {
	for _, n := range []int{1 << 10, 1 << 13, 1 << 17, 100000} {
		spec, _ := RFFT(nil, randReal(rand.New(rand.NewPCG(1, 2)), n))
		dst := make([]float64, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = IRFFT(dst, spec, n)
			}
		})
	}
}
