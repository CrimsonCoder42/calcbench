package trig

import (
	"math"
	"testing"
)

// Worked example: DLMF Table 4.23.1, arccot 1 = π/4 and arccot(−1) = −π/4. At 0 the table gives ∓π/2;
// this package returns +π/2 for both zeros (the range (−π/2, π/2]).
func TestAcot_dlmfTable(t *testing.T) {
	for _, tc := range []struct{ x, want float64 }{{1, math.Pi / 4}, {-1, -math.Pi / 4}, {0, math.Pi / 2}, {math.Copysign(0, -1), math.Pi / 2}} {
		if got, err := Acot(tc.x); err != nil || got != tc.want {
			t.Errorf("Acot(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestAcot_cases(t *testing.T) {
	tests := []struct {
		name    string
		x, want float64
	}{
		{"√3 gives π/6", math.Sqrt(3), math.Pi / 6},
		{"large positive tends to +0", 1e300, 1e-300},
		{"large negative tends to −0", -1e300, -1e-300},
		{"tiny positive tends to π/2", 1e-300, math.Pi / 2},
		{"tiny negative rounds to −math.Pi/2", -1e-300, -math.Pi / 2},
		{"subnormal", 5e-324, math.Pi / 2},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Acot(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
				t.Errorf("Acot(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
	// Odd away from 0 (DLMF 4.23.15).
	p, _ := Acot(0.3)
	if m, _ := Acot(-0.3); m != -p {
		t.Errorf("Acot(−0.3) = %v, want %v", m, -p)
	}
}

func TestAcot_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Acot(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Acot(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := Acot(math.Inf(-1)); return err }, ErrNotFinite},
	})
}

// Oracle: cot(Acot x) = x and the result is in (−π/2, π/2] with the sign of x, for 0.01 ≤ |x| ≤ 100;
// tolerance 1e-12 relative.
func TestAcot_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := math.Pow(10, rng.Float64()*4-2)
		if rng.IntN(2) == 0 {
			x = -x
		}
		y, err := Acot(x)
		if err != nil || y <= -math.Pi/2 || y > math.Pi/2 || math.Signbit(y) != math.Signbit(x) || !relNear(math.Cos(y)/math.Sin(y), x, 1e-12) {
			t.Fatalf("Acot(%v) = (%v, %v)", x, y, err)
		}
	}
}

func TestAcot_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Acot(2.5) }); n != 0 {
		t.Errorf("Acot allocates %v times, want 0", n)
	}
}

func BenchmarkAcot(b *testing.B) {
	b.ReportAllocs()
	x := 2.5
	for b.Loop() {
		sink, _ = Acot(x)
	}
}
