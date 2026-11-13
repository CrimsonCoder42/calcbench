package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.42.4 (a/sin A = b/sin B): a = 10, A = B = 60° gives b = 10;
// the 3-4-5 triangle, a = 3, A = atan2(3, 4), B = π/2, gives b = 5.
func TestLawOfSinesSide_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ a, angleA, angleB, want float64 }{{10, radOf(60), radOf(60), 10}, {3, math.Atan2(3, 4), math.Pi / 2, 5}} {
		if got, err := LawOfSinesSide(tc.a, tc.angleA, tc.angleB); err != nil || !relNear(got, tc.want, ulps4) {
			t.Errorf("LawOfSinesSide(%v, %v, %v) = (%v, %v), want %v", tc.a, tc.angleA, tc.angleB, got, err, tc.want)
		}
	}
}

func TestLawOfSinesSide_cases(t *testing.T) {
	tests := []struct {
		name                    string
		a, angleA, angleB, want float64
	}{
		{"30-60-90: side opposite 60° is a√3", 1, math.Pi / 6, math.Pi / 3, math.Sqrt(3)},
		{"tiny angle B", 1, math.Pi / 2, 1e-300, 1e-300},
		{"huge side that fits", 1e300, math.Pi / 2, math.Pi / 6, 5e299},
		{"tiny side", 1e-300, math.Pi / 6, math.Pi / 2, 2e-300},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := LawOfSinesSide(tc.a, tc.angleA, tc.angleB); err != nil || !relNear(got, tc.want, ulps4) {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
		})
	}
}

func TestLawOfSinesSide_errors(t *testing.T) {
	call := func(a, angleA, angleB float64) func() error {
		return func() error { _, err := LawOfSinesSide(a, angleA, angleB); return err }
	}
	runErrCases(t, []errCase{
		{"side NaN", call(math.NaN(), 1, 1), ErrNotFinite},
		{"side +Inf", call(math.Inf(1), 1, 1), ErrNotFinite},
		{"angle A NaN", call(1, math.NaN(), 1), ErrNotFinite},
		{"angle B -Inf", call(1, 1, math.Inf(-1)), ErrNotFinite},
		{"side zero", call(0, 1, 1), ErrNonPositive},
		{"side negative", call(-3, 1, 1), ErrNonPositive},
		{"angle A zero", call(1, 0, 1), ErrDomain},
		{"angle B negative", call(1, 1, -0.1), ErrDomain},
		{"angle A π", call(1, math.Pi, 0.1), ErrDomain},
		{"angles sum to π", call(1, math.Pi/2, math.Pi/2), ErrDomain},
		{"angles sum past π", call(1, 2, 2), ErrDomain},
		{"result overflows", call(1e308, 1e-300, math.Pi/2), ErrOverflow},
	})
}

// Oracle: a random triangle built from its angles has sides proportional to the sines (DLMF 4.42.4);
// b must equal k·sin B for a = k·sin A, to a few ulps (1e-15 relative).
func TestLawOfSinesSide_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		angleA := 0.01 + rng.Float64()*3
		angleB := 0.01 + rng.Float64()*(math.Pi-angleA-0.02)
		k := math.Pow(10, rng.Float64()*20-10)
		b, err := LawOfSinesSide(k*math.Sin(angleA), angleA, angleB)
		if err != nil || !relNear(b, k*math.Sin(angleB), 1e-15) {
			t.Fatalf("LawOfSinesSide(%v, %v, %v) = (%v, %v), want %v", k*math.Sin(angleA), angleA, angleB, b, err, k*math.Sin(angleB))
		}
	}
}

func TestLawOfSinesSide_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = LawOfSinesSide(3, 0.6, 1.2) }); n != 0 {
		t.Errorf("LawOfSinesSide allocates %v times, want 0", n)
	}
}

func BenchmarkLawOfSinesSide(b *testing.B) {
	b.ReportAllocs()
	a, angleA, angleB := 3.0, 0.6, 1.2
	for b.Loop() {
		sink, _ = LawOfSinesSide(a, angleA, angleB)
	}
}
