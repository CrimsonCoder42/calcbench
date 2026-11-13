package calculus

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

func near(got, want, tol float64) bool {
	return got == want || math.Abs(got-want) <= tol*math.Max(1, math.Abs(want))
}

// Worked example (A&S 25.4.2 applied panel by panel, hand-checked): y = x² at 0, 1, 2, 3 gives
// (0+1)/2 + (1+4)/2 + (4+9)/2 = 9.5.
func TestTrapezoid_quadratic(t *testing.T) {
	got, err := Trapezoid([]float64{0, 1, 4, 9}, []float64{0, 1, 2, 3})
	if err != nil || got != 9.5 {
		t.Errorf("Trapezoid = %v, %v; want 9.5", got, err)
	}
}

func TestTrapezoid_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		y, x []float64
		want float64
	}{
		{"two points", []float64{1, 3}, []float64{0, 2}, 4},
		{"linear exact uneven", []float64{1, 2, 5}, []float64{0, 1, 4}, 12},
		{"descending x negates", []float64{5, 2, 1}, []float64{4, 1, 0}, -12},
		{"zero", []float64{0, 0, 0}, []float64{-1, 0, 1}, 0},
		{"huge y no overflow", []float64{m, m}, []float64{0, 0.5}, m / 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Trapezoid(tt.y, tt.x)
			if err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("Trapezoid(%v, %v) = %v, %v; want %v", tt.y, tt.x, got, err, tt.want)
			}
		})
	}
}

// Oracle: the composite rule converges to ∫₀^π sin = 2 with error ≈ π h²/12·(average |f”|) ≤ π³/(12 n²).
func TestTrapezoid_sinConvergence(t *testing.T) {
	const n = 1000
	x, y := make([]float64, n+1), make([]float64, n+1)
	for i := range x {
		x[i] = math.Pi * float64(i) / n
		y[i] = math.Sin(x[i])
	}
	got, err := Trapezoid(y, x)
	if err != nil || math.Abs(got-2) > math.Pow(math.Pi, 3)/(12*n*n) {
		t.Errorf("Trapezoid(sin) = %v, %v; want 2 within π³/(12n²)", got, err)
	}
	u, err := TrapezoidUniform(y, math.Pi/n)
	// Tolerance 1e-14: same panels, uniform step.
	if err != nil || !near(u, got, 1e-14) {
		t.Errorf("TrapezoidUniform = %v, %v; want %v", u, err, got)
	}
}

func TestTrapezoid_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		y, x []float64
		want error
	}{
		{"mismatch", []float64{1, 2}, []float64{1}, ErrLengthMismatch},
		{"one", []float64{1}, []float64{1}, ErrTooFew},
		{"empty", nil, nil, ErrTooFew},
		{"NaN x", []float64{1, 2}, []float64{0, math.NaN()}, ErrNotFinite},
		{"Inf y", []float64{1, math.Inf(-1)}, []float64{0, 1}, ErrNotFinite},
		{"panel overflow", []float64{m, m}, []float64{0, 4}, ErrOverflow},
		{"sum overflow", []float64{m, m, m}, []float64{0, 1, 2}, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Trapezoid(tt.y, tt.x); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func TestTrapezoidUniform_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		y    []float64
		dx   float64
		want float64
	}{
		{"worked", []float64{1, 2, 3}, 0.5, 2},
		{"negative dx", []float64{1, 2, 3}, -0.5, -2},
		{"two", []float64{4, 6}, 1, 5},
		{"huge sum small dx", []float64{m, m, m}, 0.25, m / 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TrapezoidUniform(tt.y, tt.dx)
			if err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("TrapezoidUniform(%v, %v) = %v, %v; want %v", tt.y, tt.dx, got, err, tt.want)
			}
		})
	}
}

func TestTrapezoidUniform_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		y    []float64
		dx   float64
		want error
	}{
		{"one", []float64{1}, 1, ErrTooFew},
		{"NaN dx", []float64{1, 2}, math.NaN(), ErrNotFinite},
		{"Inf y", []float64{1, math.Inf(1)}, 1, ErrNotFinite},
		{"overflow", []float64{m, m, m}, 4, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := TrapezoidUniform(tt.y, tt.dx); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func TestTrapezoid_oracleUniform(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		n := 2 + r.IntN(60)
		y, x := make([]float64, n), make([]float64, n)
		for i := range y {
			y[i] = r.Float64()*200 - 100
			x[i] = float64(i) * 0.25
		}
		a, err1 := Trapezoid(y, x)
		b, err2 := TrapezoidUniform(y, 0.25)
		// Tolerance 1e-12: two summation orders of up to 60 terms of magnitude ≤ 100.
		if err1 != nil || err2 != nil || math.Abs(a-b) > 1e-12 {
			t.Fatalf("Trapezoid %v (%v) vs TrapezoidUniform %v (%v)", a, err1, b, err2)
		}
	}
}

func TestTrapezoid_noAlloc(t *testing.T) {
	y, x := []float64{1, 2, 3}, []float64{0, 1, 2}
	if n := testing.AllocsPerRun(50, func() { _, _ = Trapezoid(y, x); _, _ = TrapezoidUniform(y, 1) }); n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func benchData(n int) (y, x []float64) {
	y, x = make([]float64, n), make([]float64, n)
	for i := range y {
		x[i] = float64(i) * 0.001
		y[i] = math.Sin(x[i])
	}
	return y, x
}

func BenchmarkTrapezoid(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		y, x := benchData(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Trapezoid(y, x)
			}
		})
	}
}

func BenchmarkTrapezoidUniform(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		y, _ := benchData(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = TrapezoidUniform(y, 0.001)
			}
		})
	}
}
