package reduce

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"testing"
)

// Our example: running products of 1..5 are the factorials 1, 2, 6, 24, 120.
func TestCumProd_example(t *testing.T) {
	got, err := CumProd(nil, []float64{1, 2, 3, 4, 5})
	if err != nil || !slices.Equal(got, []float64{1, 2, 6, 24, 120}) {
		t.Errorf("CumProd = (%v, %v)", got, err)
	}
}

func TestCumProd_cases(t *testing.T) {
	got, err := CumProd(nil, nil)
	if err != nil || len(got) != 0 {
		t.Errorf("empty: (%v, %v)", got, err)
	}
	got, err = CumProd(nil, []float64{2, 0, 3})
	if err != nil || !slices.Equal(got, []float64{2, 0, 0}) {
		t.Errorf("zero: (%v, %v)", got, err)
	}
	// 1e-200 · 1e-200 underflows to 0 at step 2 (true value 1e-400), and the scaled product recovers
	// 1e-100 at step 3.
	got, err = CumProd(nil, []float64{1e-200, 1e-200, 1e300})
	if err != nil || got[1] != 0 || math.Abs(got[2]/1e-100-1) > 8*0x1p-52 {
		t.Errorf("scaled: (%v, %v)", got, err)
	}
	in := []float64{2, 3}
	got, _ = CumProd(in, in)
	if !slices.Equal(got, []float64{2, 6}) || &got[0] != &in[0] {
		t.Errorf("in place: %v", got)
	}
}

func TestCumProd_errors(t *testing.T) {
	if _, err := CumProd(nil, []float64{math.Inf(1)}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: err = %v", err)
	}
	if _, err := CumProd(nil, []float64{1e200, 1e200}); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: err = %v", err)
	}
}

func BenchmarkCumProd(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchValues(n)
		dst := make([]float64, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = CumProd(dst, values)
			}
		})
	}
}
