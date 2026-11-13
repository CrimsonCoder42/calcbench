package reduce

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

// Higham §27.8: scaling avoids spurious overflow. 1e300·1e300·1e-300·1e-300 overflows in a plain loop,
// yet the product is 1 (to rounding: four roundings, so within 4u; we allow 8 ulp).
func TestProd_scalingExample(t *testing.T) {
	got, err := Prod([]float64{1e300, 1e300, 1e-300, 1e-300})
	if err != nil || math.Abs(got-1) > 8*0x1p-52 {
		t.Errorf("Prod = (%v, %v), want (≈1, nil)", got, err)
	}
	got, err = Prod([]float64{1e-300, 1e-300, 1e300, 1e300}) // no spurious underflow either
	if err != nil || math.Abs(got-1) > 8*0x1p-52 {
		t.Errorf("Prod underflow order = (%v, %v), want (≈1, nil)", got, err)
	}
}

func TestProd_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"empty", nil, 1},
		{"one", []float64{-2.5}, -2.5},
		{"two", []float64{3, 4}, 12},
		{"factorial 10", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3628800},
		{"zero", []float64{5, 0, 7}, 0},
		{"signs", []float64{-1, -2, -3}, -6},
		{"max", []float64{math.MaxFloat64, 1}, math.MaxFloat64},
		{"underflow to zero", []float64{math.SmallestNonzeroFloat64, 0.25}, 0},
	}
	for _, tc := range tests {
		got, err := Prod(tc.values)
		if err != nil || got != tc.want {
			t.Errorf("%s: Prod = (%v, %v), want (%v, nil)", tc.name, got, err, tc.want)
		}
	}
}

func TestProd_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"NaN", []float64{math.NaN()}, ErrNotFinite},
		{"Inf", []float64{2, math.Inf(-1)}, ErrNotFinite},
		{"overflow", []float64{math.MaxFloat64, 2}, ErrOverflow},
	}
	for _, tc := range tests {
		if _, err := Prod(tc.values); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func TestProd_noAlloc(t *testing.T) {
	values := []float64{1, 2, 3}
	if a := testing.AllocsPerRun(100, func() { _, _ = Prod(values) }); a != 0 {
		t.Errorf("Prod allocates %v times, want 0", a)
	}
}

func BenchmarkProd(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchValues(n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Prod(values)
			}
		})
	}
}
