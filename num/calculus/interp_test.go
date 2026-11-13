package calculus

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example (NR §3.3 eq. 3.3.1, hand-checked): table (1,3), (2,2), (3,0) at 2.5 → 1.
func TestInterp_cases(t *testing.T) {
	xp, fp := []float64{1, 2, 3}, []float64{3, 2, 0}
	for _, tt := range []struct {
		name    string
		x, want float64
	}{
		{"worked", 2.5, 1},
		{"first node", 1, 3},
		{"interior node", 2, 2},
		{"last node", 3, 0},
		{"quarter", 1.25, 2.75},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Interp(tt.x, xp, fp)
			if err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("Interp(%v) = %v, %v; want %v", tt.x, got, err, tt.want)
			}
		})
	}
	m := math.MaxFloat64
	if got, err := Interp(0, []float64{-m, m}, []float64{-m, m}); err != nil || got != 0 {
		t.Errorf("Interp over full range = %v, %v; want 0", got, err)
	}
}

// Oracle: a linear scan for the interval gives the same value.
func TestInterp_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := 2 + r.IntN(30)
		xp, fp := make([]float64, n), make([]float64, n)
		for i := range xp {
			xp[i] = float64(i) + r.Float64()*0.5
			fp[i] = r.Float64()*10 - 5
		}
		x := xp[0] + r.Float64()*(xp[n-1]-xp[0])
		j := 0
		for j < n-2 && xp[j+1] <= x {
			j++
		}
		want := fp[j] + (x-xp[j])/(xp[j+1]-xp[j])*(fp[j+1]-fp[j])
		got, err := Interp(x, xp, fp)
		// Tolerance 1e-12: different operation order on values ≤ 30.
		if err != nil || math.Abs(got-want) > 1e-12 {
			t.Fatalf("Interp(%v) = %v, %v; want %v", x, got, err, want)
		}
	}
}

func TestInterp_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		x      float64
		xp, fp []float64
		want   error
	}{
		{"mismatch", 1, []float64{1, 2}, []float64{1}, ErrLengthMismatch},
		{"one point", 1, []float64{1}, []float64{1}, ErrTooFew},
		{"NaN x", math.NaN(), []float64{1, 2}, []float64{1, 2}, ErrNotFinite},
		{"Inf xp", 1, []float64{1, math.Inf(1)}, []float64{1, 2}, ErrNotFinite},
		{"NaN fp", 1, []float64{1, 2}, []float64{1, math.NaN()}, ErrNotFinite},
		{"not ascending", 1, []float64{1, 1}, []float64{1, 2}, ErrNotAscending},
		{"below", 0.5, []float64{1, 2}, []float64{1, 2}, ErrOutOfRange},
		{"above", 2.5, []float64{1, 2}, []float64{1, 2}, ErrOutOfRange},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Interp(tt.x, tt.xp, tt.fp); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func TestInterp_noAlloc(t *testing.T) {
	xp, fp := []float64{1, 2, 3}, []float64{3, 2, 0}
	if n := testing.AllocsPerRun(50, func() { _, _ = Interp(2.5, xp, fp) }); n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func BenchmarkInterp(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		fp, xp := benchData(n)
		x := xp[n/3] + 0.0004
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Interp(x, xp, fp)
			}
		})
	}
}
