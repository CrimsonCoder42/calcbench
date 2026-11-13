package reduce

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"testing"
)

// Abramowitz & Stegun §25.1: the forward differences of the squares 1, 4, 9, 16, 25 are 3, 5, 7, 9; the
// second differences are constant 2 and the third are 0 (a quadratic's).
func TestDiff_squares(t *testing.T) {
	sq := []float64{1, 4, 9, 16, 25}
	tests := []struct {
		order int
		want  []float64
	}{
		{0, []float64{1, 4, 9, 16, 25}},
		{1, []float64{3, 5, 7, 9}},
		{2, []float64{2, 2, 2}},
		{3, []float64{0, 0}},
		{5, []float64{}},
	}
	for _, tc := range tests {
		got, err := Diff(nil, sq, tc.order)
		if err != nil || !slices.Equal(got, tc.want) {
			t.Errorf("order %d: Diff = (%v, %v), want %v", tc.order, got, err, tc.want)
		}
	}
	if !slices.Equal(sq, []float64{1, 4, 9, 16, 25}) {
		t.Errorf("input modified: %v", sq)
	}
}

func TestDiff_cases(t *testing.T) {
	got, err := Diff(nil, nil, 0)
	if err != nil || len(got) != 0 {
		t.Errorf("empty: (%v, %v)", got, err)
	}
	in := []float64{1, 2, 4, 8}
	got, err = Diff(in, in, 1)
	if err != nil || !slices.Equal(got, []float64{1, 2, 4}) || &got[0] != &in[0] {
		t.Errorf("in place: (%v, %v)", got, err)
	}
}

func TestDiff_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		order  int
		want   error
	}{
		{"negative order", []float64{1, 2}, -1, ErrOutOfRange},
		{"order too large", []float64{1, 2}, 3, ErrOutOfRange},
		{"NaN", []float64{1, math.NaN()}, 1, ErrNotFinite},
		{"overflow", []float64{-math.MaxFloat64, math.MaxFloat64}, 1, ErrOverflow},
		{"overflow then cancels", []float64{-math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64}, 2, ErrOverflow},
	}
	for _, tc := range tests {
		if _, err := Diff(nil, tc.values, tc.order); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func BenchmarkDiff(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchValues(n)
		dst := make([]float64, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Diff(dst, values, 1)
			}
		})
	}
}
