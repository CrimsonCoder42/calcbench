package calculus

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example (Fornberg 1988 three-node stencil, hand-checked): y = x² at 0, 1, 3 → 0, 2, 6.
func TestGradient_quadratic(t *testing.T) {
	got, err := Gradient(nil, []float64{0, 1, 9}, []float64{0, 1, 3})
	want := []float64{0, 2, 6}
	for i := range want {
		if err != nil || math.Abs(got[i]-want[i]) > 1e-14 {
			t.Fatalf("Gradient = %v, %v; want %v", got, err, want)
		}
	}
}

// Uniform grid: interior weights −1/2, 0, 1/2 over h and end weights −3/2, 2, −1/2 (Fornberg 1988, Table 1).
func TestGradient_uniformTable(t *testing.T) {
	got, err := Gradient(nil, []float64{1, 2, 4, 7, 11}, []float64{0, 1, 2, 3, 4})
	want := []float64{0.5, 1.5, 2.5, 3.5, 4.5}
	for i := range want {
		if err != nil || math.Abs(got[i]-want[i]) > 1e-14 {
			t.Fatalf("Gradient = %v, %v; want %v", got, err, want)
		}
	}
}

// Oracle: exact (to rounding) for random quadratics on random strictly ascending grids.
func TestGradient_quadraticOracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := 3 + r.IntN(40)
		a, b, c := r.Float64()*4-2, r.Float64()*4-2, r.Float64()*4-2
		x, y := make([]float64, n), make([]float64, n)
		for i := range x {
			x[i] = float64(i) + r.Float64()*0.9
			y[i] = a*x[i]*x[i] + b*x[i] + c
		}
		got, err := Gradient(nil, y, x)
		if err != nil {
			t.Fatal(err)
		}
		for i := range x {
			want := 2*a*x[i] + b
			// Tolerance 1e-9 absolute: values up to ~3000 differenced over steps ≥ 0.1.
			if math.Abs(got[i]-want) > 1e-9 {
				t.Fatalf("i=%d got %v want %v", i, got[i], want)
			}
		}
	}
}

func TestGradient_twoPoints(t *testing.T) {
	got, err := Gradient(nil, []float64{1, 4}, []float64{0, 2})
	if err != nil || got[0] != 1.5 || got[1] != 1.5 {
		t.Errorf("Gradient = %v, %v; want [1.5 1.5]", got, err)
	}
}

func TestGradient_dst(t *testing.T) {
	dst := make([]float64, 0, 4)
	y, x := []float64{0, 1, 4, 9}, []float64{0, 1, 2, 3}
	got, err := Gradient(dst, y, x)
	if err != nil || &got[:1][0] != &dst[:1][0] {
		t.Errorf("Gradient did not reuse dst: %v", err)
	}
	if n := testing.AllocsPerRun(50, func() { _, _ = Gradient(dst, y, x) }); n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func TestGradient_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		y, x []float64
		want error
	}{
		{"mismatch", []float64{1, 2}, []float64{1}, ErrLengthMismatch},
		{"one", []float64{1}, []float64{1}, ErrTooFew},
		{"NaN y", []float64{1, math.NaN()}, []float64{0, 1}, ErrNotFinite},
		{"Inf x", []float64{1, 2}, []float64{0, math.Inf(1)}, ErrNotFinite},
		{"not ascending", []float64{1, 2, 3}, []float64{0, 2, 1}, ErrNotAscending},
		{"repeated x", []float64{1, 2}, []float64{1, 1}, ErrNotAscending},
		{"step overflow", []float64{1, 2}, []float64{-m, m}, ErrOverflow},
		{"derivative overflow", []float64{-m, m}, []float64{0, 0.5}, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Gradient(nil, tt.y, tt.x); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkGradient(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		y, x := benchData(n)
		dst := make([]float64, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Gradient(dst, y, x)
			}
		})
	}
}
