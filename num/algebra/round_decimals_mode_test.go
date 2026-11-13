package algebra

import (
	"errors"
	"math"
	"testing"
)

// Hand computations of the IEEE 754-2019 §4.3 and §9.5 rules on exact binary ties (x.5, x.25, x.125).
func TestRoundDecimalsMode_cases(t *testing.T) {
	E, A, Z := RoundTiesToEven, RoundTiesToAway, RoundTiesTowardZero
	P, N, T := RoundTowardPositive, RoundTowardNegative, RoundTowardZero
	tests := []struct {
		x    float64
		d    int
		mode RoundMode
		want float64
	}{
		{2.5, 0, E, 2}, {2.5, 0, A, 3}, {2.5, 0, Z, 2}, {2.5, 0, P, 3}, {2.5, 0, N, 2}, {2.5, 0, T, 2},
		{-2.5, 0, E, -2}, {-2.5, 0, A, -3}, {-2.5, 0, Z, -2}, {-2.5, 0, P, -2}, {-2.5, 0, N, -3}, {-2.5, 0, T, -2},
		{3.5, 0, E, 4}, {3.5, 0, Z, 3},
		{2.6, 0, Z, 3}, {-2.6, 0, Z, -3}, {2.4, 0, A, 2},
		{1.125, 2, E, 1.12}, {1.125, 2, A, 1.13}, {1.125, 2, Z, 1.12},
		{1.21, 1, P, 1.3}, {-1.21, 1, N, -1.3}, {1.29, 1, T, 1.2},
		{1250, -2, E, 1200}, {1250, -2, A, 1300}, {1250, -2, Z, 1200}, {1201, -2, P, 1300},
		{0, 3, A, 0}, {1e300, 2, A, 1e300},
	}
	for _, tc := range tests {
		got, err := RoundDecimalsMode(tc.x, tc.d, tc.mode)
		if err != nil || math.Abs(got-tc.want) > 1e-12*math.Max(1, math.Abs(tc.want)) {
			t.Errorf("RoundDecimalsMode(%v, %d, %d) = %v, %v; want %v", tc.x, tc.d, tc.mode, got, err, tc.want)
		}
	}
}

// RoundTiesToEven must agree with RoundDecimals (NUM-ALGEB-022), the same rule.
func TestRoundDecimalsMode_agreesWithRoundDecimals(t *testing.T) {
	for _, x := range []float64{0.5, 1.5, 2.5, -0.5, 1234.5678, 2.675, -7.125} {
		for _, d := range []int{-2, 0, 1, 2, 3} {
			a, _ := RoundDecimals(x, d)
			b, _ := RoundDecimalsMode(x, d, RoundTiesToEven)
			if a != b {
				t.Errorf("x = %v, d = %d: RoundDecimals %v, RoundDecimalsMode %v", x, d, a, b)
			}
		}
	}
}

func TestRoundDecimalsMode_errors(t *testing.T) {
	if _, err := RoundDecimalsMode(math.NaN(), 0, RoundTiesToEven); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	if _, err := RoundDecimalsMode(1, 309, RoundTiesToEven); !errors.Is(err, ErrDegenerate) {
		t.Errorf("decimals 309: %v", err)
	}
	if _, err := RoundDecimalsMode(1, 0, RoundMode(99)); !errors.Is(err, ErrDegenerate) {
		t.Errorf("bad mode: %v", err)
	}
	if _, err := RoundDecimalsMode(1, 0, RoundMode(-1)); !errors.Is(err, ErrDegenerate) {
		t.Errorf("negative mode: %v", err)
	}
	if a := testing.AllocsPerRun(100, func() { _, _ = RoundDecimalsMode(1.125, 2, RoundTiesTowardZero) }); a != 0 {
		t.Errorf("allocates %v", a)
	}
}

func BenchmarkRoundDecimalsMode(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = RoundDecimalsMode(1234.5678, 2, RoundTiesToAway)
	}
}
