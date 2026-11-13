package invest

import (
	"fmt"
	"math"
	"testing"
)

// Hand-worked: mean 0.015; shortfalls below 0 are −0.05, −0.01 → DD = √((0.0025 + 0.0001)/4).
func TestSortinoRatio_handWorked(t *testing.T) {
	got, err := SortinoRatio([]float64{0.10, -0.05, 0.02, -0.01}, 0)
	want := 0.015 / math.Sqrt(0.0026/4)
	if err != nil || !near(got, want, 1e-14) || !near(got, 0.58834, 1e-5) {
		t.Errorf("SortinoRatio = %v, %v; want %v", got, err, want)
	}
}

// MAR 0.02: shortfalls 0.10 → none, −0.07, 0 (not below), −0.03 → DD = √((0.0049 + 0.0009)/4); mean − MAR = −0.005.
func TestSortinoRatio_nonZeroMAR(t *testing.T) {
	got, err := SortinoRatio([]float64{0.10, -0.05, 0.02, -0.01}, 0.02)
	want := -0.005 / math.Sqrt(0.0058/4)
	if err != nil || !near(got, want, 1e-14) {
		t.Errorf("SortinoRatio = %v, %v; want %v", got, err, want)
	}
}

func TestSortinoRatio_errors(t *testing.T) {
	_, err := SortinoRatio(nil, 0)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = SortinoRatio([]float64{0.01, 0.02}, 0)
	wantErr(t, "none below MAR", err, ErrZeroDenominator)
	_, err = SortinoRatio([]float64{0.01, math.NaN()}, 0)
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = SortinoRatio([]float64{0.01, -1.5}, 0)
	wantErr(t, "below -1", err, ErrRate)
	for _, bad := range []float64{math.NaN(), math.Inf(1), -1} {
		_, err = SortinoRatio([]float64{0.01, -0.02}, bad)
		wantErr(t, "bad MAR", err, ErrRate)
	}
	_, err = SortinoRatio([]float64{-0.5, 1e300}, 1e200)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestSortinoRatio_noAlloc(t *testing.T) {
	r := []float64{0.10, -0.05, 0.02, -0.01}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = SortinoRatio(r, 0) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkSortinoRatio(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = SortinoRatio(r, 0)
			}
		})
	}
}
