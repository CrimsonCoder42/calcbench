package trig

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example (Itoh 1982, eq. 4, hand-checked): a phase ramp in steps of 3π/4, wrapped to (−π, π].
func TestUnwrapAngles_ramp(t *testing.T) {
	in := []float64{0, 3 * math.Pi / 4, -3 * math.Pi / 4, -math.Pi / 4}
	want := []float64{0, 3 * math.Pi / 4, 5 * math.Pi / 4, 7 * math.Pi / 4}
	got, err := UnwrapAngles(nil, in, 2*math.Pi)
	if err != nil {
		t.Fatal(err)
	}
	for i := range want {
		// Tolerance 1e-14: one remainder and a few adds per step.
		if math.Abs(got[i]-want[i]) > 1e-14 {
			t.Errorf("UnwrapAngles[%d] = %v; want %v", i, got[i], want[i])
		}
	}
}

func TestUnwrapAngles_cases(t *testing.T) {
	for _, tt := range []struct {
		name     string
		in, want []float64
		period   float64
	}{
		{"empty", nil, []float64{}, 2 * math.Pi},
		{"one", []float64{5}, []float64{5}, 2 * math.Pi},
		{"no jump kept exactly", []float64{0.1, 0.2, 0.3}, []float64{0.1, 0.2, 0.3}, 2 * math.Pi},
		{"degrees period 360", []float64{350, 10, 30, 340}, []float64{350, 370, 390, 340}, 360},
		{"downward jump", []float64{10, 350}, []float64{10, -10}, 360},
		{"half period positive kept", []float64{0, 180}, []float64{0, 180}, 360},
		{"half period negative kept", []float64{0, -180}, []float64{0, -180}, 360},
		{"multi-period jump", []float64{0, 725}, []float64{0, 5}, 360},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnwrapAngles(nil, tt.in, tt.period)
			if err != nil || len(got) != len(tt.want) {
				t.Fatalf("UnwrapAngles(%v) = %v, %v; want %v", tt.in, got, err, tt.want)
			}
			for i := range got {
				if math.Abs(got[i]-tt.want[i]) > 1e-12 {
					t.Errorf("UnwrapAngles(%v)[%d] = %v; want %v", tt.in, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// Oracle: a smooth signal with steps under π, wrapped with NormalizeAngleSigned, unwraps back to itself.
func TestUnwrapAngles_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		n := 1 + r.IntN(100)
		truth := make([]float64, n)
		wrapped := make([]float64, n)
		truth[0] = r.Float64()*2*math.Pi - math.Pi
		for i := range truth {
			if i > 0 {
				truth[i] = truth[i-1] + (r.Float64()*2-1)*3
			}
			wrapped[i], _ = NormalizeAngleSigned(truth[i])
		}
		got, err := UnwrapAngles(nil, wrapped, 2*math.Pi)
		if err != nil {
			t.Fatal(err)
		}
		for i := range got {
			// Tolerance 1e-9: rounding accumulates over up to 100 steps of magnitude ≤ 300.
			if math.Abs(got[i]-truth[i]) > 1e-9 {
				t.Fatalf("index %d: got %v want %v", i, got[i], truth[i])
			}
		}
	}
}

func TestUnwrapAngles_inPlace(t *testing.T) {
	xs := []float64{350, 10, 30}
	got, err := UnwrapAngles(xs, xs, 360)
	if err != nil || &got[0] != &xs[0] || got[1] != 370 || got[2] != 390 {
		t.Errorf("in place = %v, %v", got, err)
	}
	dst := make([]float64, 0, 3)
	if n := testing.AllocsPerRun(50, func() { _, _ = UnwrapAngles(dst, []float64{1, 2, 3}, 2*math.Pi) }); n != 0 {
		t.Errorf("allocs with dst = %v; want 0", n)
	}
}

func TestUnwrapAngles_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name   string
		in     []float64
		period float64
		want   error
	}{
		{"NaN angle", []float64{0, math.NaN()}, 360, ErrNotFinite},
		{"Inf angle", []float64{math.Inf(1)}, 360, ErrNotFinite},
		{"NaN period", []float64{0}, math.NaN(), ErrNotFinite},
		{"Inf period", []float64{0}, math.Inf(1), ErrNotFinite},
		{"zero period", []float64{0}, 0, ErrDomain},
		{"negative period", []float64{0}, -1, ErrDomain},
		{"overflow", []float64{0.9 * m, -0.9 * m}, m, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := UnwrapAngles(nil, tt.in, tt.period); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkUnwrapAngles(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		xs := make([]float64, n)
		for i := range xs {
			xs[i], _ = NormalizeAngleSigned(float64(i) * 0.7)
		}
		dst := make([]float64, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = UnwrapAngles(dst, xs, 2*math.Pi)
			}
		})
	}
}
