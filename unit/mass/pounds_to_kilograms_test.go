package mass

import (
	"errors"
	"math"
	"testing"
)

// NIST SP 811 App. B.8: 1 lb = 0.453 592 37 kg exactly.
func TestPoundsToKilograms_example(t *testing.T) {
	tests := []struct{ lb, want float64 }{
		{1, 0.45359237},
		{150, 68.0388555},
		{2.2046226218487757, 1},
		{0, 0},
		{100, 45.359237},
	}
	for _, tc := range tests {
		got, err := PoundsToKilograms(tc.lb)
		if err != nil || !nearRel(got, tc.want, convTol) {
			t.Errorf("PoundsToKilograms(%v) = (%v, %v), want %v", tc.lb, got, err, tc.want)
		}
	}
}

// parity: ledger C6 (unit_convert.go:19), source divides by 2.20462; agreement to its rounding (2e-6).
func TestPoundsToKilograms_parity(t *testing.T) {
	for _, lb := range []float64{1, 154, 220.5} {
		got, err := PoundsToKilograms(lb)
		if err != nil {
			t.Fatal(err)
		}
		if source := lb / 2.20462; !nearRel(source, got, 2e-6) {
			t.Errorf("PoundsToKilograms(%v) = %v, source %v differs by more than its rounding", lb, got, source)
		}
	}
}

func TestPoundsToKilograms_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := PoundsToKilograms(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("PoundsToKilograms(%v) error = %v, want ErrNotFinite", v, err)
		}
	}
	if _, err := PoundsToKilograms(-0.5); !errors.Is(err, ErrNegative) {
		t.Errorf("PoundsToKilograms(-0.5) error = %v, want ErrNegative", err)
	}
	if got, err := PoundsToKilograms(math.MaxFloat64); err != nil || math.IsInf(got, 0) {
		t.Errorf("PoundsToKilograms(MaxFloat64) = (%v, %v), want finite", got, err)
	}
}

func TestPoundsToKilograms_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkMass, _ = PoundsToKilograms(160) }); a != 0 {
		t.Errorf("PoundsToKilograms allocates %v times, want 0", a)
	}
}

func BenchmarkPoundsToKilograms(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkMass, _ = PoundsToKilograms(160)
	}
}
