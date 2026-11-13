package length

import (
	"errors"
	"math"
	"testing"
)

// NIST SP 811 App. B.8: 1 in = 2.54 cm exactly.
func TestInchesToCentimeters_example(t *testing.T) {
	tests := []struct{ in, want float64 }{
		{1, 2.54}, {12, 30.48}, {70, 177.8}, {0, 0}, {100, 254},
	}
	for _, tc := range tests {
		got, err := InchesToCentimeters(tc.in)
		if err != nil || !nearRel(got, tc.want, convTol) {
			t.Errorf("InchesToCentimeters(%v) = (%v, %v), want %v", tc.in, got, err, tc.want)
		}
	}
}

// parity: ledger C6 (unit_convert.go:31), source in × 2.54: the same formula, so the same result exactly.
func TestInchesToCentimeters_parity(t *testing.T) {
	for _, in := range []float64{60, 65.5, 72} {
		got, err := InchesToCentimeters(in)
		if err != nil || got != in*2.54 {
			t.Errorf("InchesToCentimeters(%v) = (%v, %v), want %v", in, got, err, in*2.54)
		}
	}
}

func TestInchesToCentimeters_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := InchesToCentimeters(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("InchesToCentimeters(%v) error = %v, want ErrNotFinite", v, err)
		}
	}
	if _, err := InchesToCentimeters(-2); !errors.Is(err, ErrNegative) {
		t.Errorf("InchesToCentimeters(-2) error = %v, want ErrNegative", err)
	}
	if _, err := InchesToCentimeters(math.MaxFloat64); !errors.Is(err, ErrOverflow) {
		t.Errorf("InchesToCentimeters(MaxFloat64) error = %v, want ErrOverflow", err)
	}
	if got, err := InchesToCentimeters(math.MaxFloat64 / 3); err != nil || math.IsInf(got, 0) {
		t.Errorf("InchesToCentimeters(MaxFloat64/3) = (%v, %v), want finite", got, err)
	}
}

func TestInchesToCentimeters_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkLength, _ = InchesToCentimeters(68) }); a != 0 {
		t.Errorf("InchesToCentimeters allocates %v times, want 0", a)
	}
}

func BenchmarkInchesToCentimeters(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkLength, _ = InchesToCentimeters(68)
	}
}
