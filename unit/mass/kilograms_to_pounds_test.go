package mass

import (
	"errors"
	"math"
	"testing"
)

var sinkMass float64

// nearRel reports whether got is within relTol of want, relative to |want| (absolute when want is 0).
// Tolerance 2⁻⁵¹ relative: one division or multiplication by a factor that is itself rounded to binary.
func nearRel(got, want, relTol float64) bool {
	if want == 0 {
		return math.Abs(got) <= relTol
	}
	return math.Abs(got-want) <= relTol*math.Abs(want)
}

const convTol = 0x1p-51

// NIST SP 811 App. B.8: 1 lb = 0.453 592 37 kg exactly, so 0.45359237 kg is 1 lb and 1 kg is
// 1/0.45359237 = 2.204 622 621 848 775 8… lb.
func TestKilogramsToPounds_example(t *testing.T) {
	tests := []struct{ kg, want float64 }{
		{0.45359237, 1},
		{1, 2.2046226218487757},
		{100, 220.46226218487757},
		{0, 0},
		{45.359237, 100},
	}
	for _, tc := range tests {
		got, err := KilogramsToPounds(tc.kg)
		if err != nil || !nearRel(got, tc.want, convTol) {
			t.Errorf("KilogramsToPounds(%v) = (%v, %v), want %v", tc.kg, got, err, tc.want)
		}
	}
}

// parity: ledger C6 (unit_convert.go:13), source cases. The source multiplies by 2.20462 (a rounded
// factor); calcbench uses the exact definition, so results agree to about 1.3e-6 relative. The test
// asserts the exact value and checks the source's value is within that rounding.
func TestKilogramsToPounds_parity(t *testing.T) {
	for _, kg := range []float64{1, 70, 100.5} {
		got, err := KilogramsToPounds(kg)
		if err != nil {
			t.Fatal(err)
		}
		source := kg * 2.20462
		if !nearRel(source, got, 2e-6) {
			t.Errorf("KilogramsToPounds(%v) = %v, source %v differs by more than its rounding", kg, got, source)
		}
	}
}

func TestKilogramsToPounds_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := KilogramsToPounds(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("KilogramsToPounds(%v) error = %v, want ErrNotFinite", v, err)
		}
	}
	for _, v := range []float64{-1e-300, -1} {
		if _, err := KilogramsToPounds(v); !errors.Is(err, ErrNegative) {
			t.Errorf("KilogramsToPounds(%v) error = %v, want ErrNegative", v, err)
		}
	}
	if _, err := KilogramsToPounds(math.MaxFloat64); !errors.Is(err, ErrOverflow) {
		t.Errorf("KilogramsToPounds(MaxFloat64) error = %v, want ErrOverflow", err)
	}
	if got, err := KilogramsToPounds(0.4 * math.MaxFloat64); err != nil || math.IsInf(got, 0) {
		t.Errorf("KilogramsToPounds(0.4·MaxFloat64) = (%v, %v), want finite", got, err)
	}
	if got, err := KilogramsToPounds(math.SmallestNonzeroFloat64); err != nil || got <= 0 {
		t.Errorf("KilogramsToPounds(smallest) = (%v, %v), want > 0", got, err)
	}
}

// Round trip through PoundsToKilograms over a grid, tolerance 2⁻⁵⁰ relative (two roundings).
func TestKilogramsToPounds_roundTrip(t *testing.T) {
	for kg := 0.001; kg < 1e6; kg *= 1.37 {
		lb, err := KilogramsToPounds(kg)
		if err != nil {
			t.Fatal(err)
		}
		back, err := PoundsToKilograms(lb)
		if err != nil || !nearRel(back, kg, 0x1p-50) {
			t.Fatalf("round trip %v kg → %v lb → (%v, %v)", kg, lb, back, err)
		}
	}
}

func TestKilogramsToPounds_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkMass, _ = KilogramsToPounds(72.5) }); a != 0 {
		t.Errorf("KilogramsToPounds allocates %v times, want 0", a)
	}
}

func BenchmarkKilogramsToPounds(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkMass, _ = KilogramsToPounds(72.5)
	}
}
