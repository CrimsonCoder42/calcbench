package signal

import (
	"slices"
	"strconv"
	"testing"
)

func TestFFTShift_cases(t *testing.T) {
	cases := []struct {
		name      string
		in, shift []int
	}{
		{"empty", []int{}, []int{}},
		{"one", []int{7}, []int{7}},
		{"two", []int{0, 1}, []int{1, 0}},
		{"odd", []int{0, 1, 2, 3, 4}, []int{3, 4, 0, 1, 2}},
		{"even", []int{0, 1, 2, 3, 4, 5}, []int{3, 4, 5, 0, 1, 2}},
	}
	for _, c := range cases {
		got := FFTShift(nil, c.in)
		if !slices.Equal(got, c.shift) {
			t.Fatalf("%s: FFTShift = %v, want %v", c.name, got, c.shift)
		}
		back := IFFTShift(nil, got)
		if !slices.Equal(back, c.in) {
			t.Fatalf("%s: IFFTShift = %v, want %v", c.name, back, c.in)
		}
	}
}

// TestFFTShift_freqCentred: FFTShift of FFTFreq is ascending (zero frequency in the centre).
func TestFFTShift_freqCentred(t *testing.T) {
	for _, n := range []int{1, 2, 7, 8} {
		f, _ := FFTFreq(nil, n, 1)
		if s := FFTShift(nil, f); !slices.IsSorted(s) {
			t.Fatalf("n=%d: %v", n, s)
		}
	}
}

func TestFFTShift_inPlace(t *testing.T) {
	for n := 0; n <= 9; n++ {
		x := make([]int, n)
		for i := range x {
			x[i] = i
		}
		want := FFTShift(nil, x)
		got := FFTShift(x, x)
		if !slices.Equal(got, want) {
			t.Fatalf("n=%d FFTShift in place: %v, want %v", n, got, want)
		}
		wantInv := IFFTShift(nil, got)
		gotInv := IFFTShift(got, got)
		if !slices.Equal(gotInv, wantInv) {
			t.Fatalf("n=%d IFFTShift in place: %v, want %v", n, gotInv, wantInv)
		}
	}
}

func TestFFTShift_noAllocWithDst(t *testing.T) {
	x := make([]complex128, 100)
	dst := make([]complex128, 100)
	if a := testing.AllocsPerRun(100, func() { _ = FFTShift(dst, x) }); a != 0 {
		t.Fatalf("allocs = %v", a)
	}
}

func BenchmarkFFTShift(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		x := make([]complex128, n)
		dst := make([]complex128, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = FFTShift(dst, x)
			}
		})
	}
}

func BenchmarkIFFTShift(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		x := make([]complex128, n)
		dst := make([]complex128, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = IFFTShift(dst, x)
			}
		})
	}
}
