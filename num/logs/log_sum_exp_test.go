package logs

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

func TestLogSumExp_cases(t *testing.T) {
	inf := math.Inf(-1)
	for _, tt := range []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one", []float64{2.5}, 2.5},
		{"log 1..4", []float64{0, math.Log(2), math.Log(3), math.Log(4)}, math.Log(10)},
		{"three large equal", []float64{1000, 1000, 1000}, 1000 + math.Log(3)},
		{"three very negative", []float64{-1000, -1000, -1000}, -1000 + math.Log(3)},
		{"with log-zero", []float64{inf, 0, inf}, 0},
		{"all log-zero", []float64{inf, inf}, inf},
		{"tiny tail kept", []float64{0, -40}, math.Log1p(math.Exp(-40))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15 relative: n roundings in [0, 1] then one log1p.
			got, err := LogSumExp(tt.values)
			if err != nil || !nearRel(got, tt.want, 1e-15) {
				t.Errorf("LogSumExp(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

// Oracle: for moderate values, log(Σ exp x) computed directly agrees with the shifted form.
func TestLogSumExp_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := 1 + r.IntN(50)
		xs := make([]float64, n)
		var direct float64
		for i := range xs {
			xs[i] = r.Float64()*20 - 10
			direct += math.Exp(xs[i])
		}
		got, err := LogSumExp(xs)
		// Tolerance 1e-13 relative: the direct sum has up to n roundings at magnitude e^10.
		if err != nil || !nearRel(got, math.Log(direct), 1e-13) {
			t.Fatalf("LogSumExp(%v) = %v, %v; want %v", xs, got, err, math.Log(direct))
		}
	}
}

func TestLogSumExp_errors(t *testing.T) {
	if _, err := LogSumExp(nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("LogSumExp(nil) err = %v; want ErrEmpty", err)
	}
	for _, xs := range [][]float64{{1, math.NaN()}, {math.Inf(1), 0}} {
		if _, err := LogSumExp(xs); !errors.Is(err, ErrNotFinite) {
			t.Errorf("LogSumExp(%v) err = %v; want ErrNotFinite", xs, err)
		}
	}
}

func TestLogSumExp_noAlloc(t *testing.T) {
	xs := []float64{1, 2, 3}
	if n := testing.AllocsPerRun(100, func() { _, _ = LogSumExp(xs) }); n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func BenchmarkLogSumExp(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		xs := make([]float64, n)
		for i := range xs {
			xs[i] = float64(i%97) - 50
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = LogSumExp(xs)
			}
		})
	}
}
