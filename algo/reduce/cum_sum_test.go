package reduce

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"testing"
)

// Our example: running sums of 1, 2, 3, 4 are 1, 3, 6, 10; and every prefix of Neumaier's 1, 1e100, 1,
// -1e100 is summed exactly (1, 1e100, 1e100, 2).
func TestCumSum_example(t *testing.T) {
	got, err := CumSum(nil, []float64{1, 2, 3, 4})
	if err != nil || !slices.Equal(got, []float64{1, 3, 6, 10}) {
		t.Errorf("CumSum = (%v, %v)", got, err)
	}
	got, err = CumSum(nil, []float64{1, 1e100, 1, -1e100})
	if err != nil || !slices.Equal(got, []float64{1, 1e100, 1e100, 2}) {
		t.Errorf("CumSum Neumaier = (%v, %v), want [1 1e100 1e100 2]", got, err)
	}
}

func TestCumSum_cases(t *testing.T) {
	got, err := CumSum(nil, nil)
	if err != nil || len(got) != 0 {
		t.Errorf("empty: (%v, %v)", got, err)
	}
	in := []float64{5, -5, 2}
	got, err = CumSum(in, in) // in place
	if err != nil || !slices.Equal(got, []float64{5, 0, 2}) || &got[0] != &in[0] {
		t.Errorf("in place: (%v, %v)", got, err)
	}
	dst := make([]float64, 0, 8)
	got, _ = CumSum(dst, []float64{1, 1})
	if &got[0] != &dst[:1][0] {
		t.Errorf("dst with capacity was not reused")
	}
}

func TestCumSum_errors(t *testing.T) {
	if _, err := CumSum(nil, []float64{1, math.NaN()}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: err = %v", err)
	}
	if _, err := CumSum(nil, []float64{math.MaxFloat64, math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: err = %v", err)
	}
}

// Oracle: each out[i] equals Sum(values[:i+1]).
func TestCumSum_matchesSum(t *testing.T) {
	values := benchValues(300)
	got, err := CumSum(nil, values)
	if err != nil {
		t.Fatal(err)
	}
	for i := range values {
		want, _ := Sum(values[:i+1])
		if got[i] != want {
			t.Fatalf("out[%d] = %v, Sum = %v", i, got[i], want)
		}
	}
}

func BenchmarkCumSum(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchValues(n)
		dst := make([]float64, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = CumSum(dst, values)
			}
		})
	}
}
