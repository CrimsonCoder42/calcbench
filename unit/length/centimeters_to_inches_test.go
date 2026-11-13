package length

import (
	"errors"
	"math"
	"testing"
)

var sinkLength float64

// nearRel reports whether got is within relTol of want, relative to |want| (absolute when want is 0).
func nearRel(got, want, relTol float64) bool {
	if want == 0 {
		return math.Abs(got) <= relTol
	}
	return math.Abs(got-want) <= relTol*math.Abs(want)
}

// convTol is 2⁻⁵¹ relative: one correctly rounded operation by 2.54, which is itself rounded to binary.
const convTol = 0x1p-51

// NIST SP 811 App. B.8: 1 in = 2.54 cm exactly, so 30.48 cm = 12 in and 254 cm = 100 in (hand-checked).
func TestCentimetersToInches_example(t *testing.T) {
	tests := []struct{ cm, want float64 }{
		{2.54, 1}, {30.48, 12}, {180, 70.86614173228347}, {0, 0}, {254, 100},
	}
	for _, tc := range tests {
		got, err := CentimetersToInches(tc.cm)
		if err != nil || !nearRel(got, tc.want, convTol) {
			t.Errorf("CentimetersToInches(%v) = (%v, %v), want %v", tc.cm, got, err, tc.want)
		}
	}
}

// parity: ledger C6 (unit_convert.go:25), source cm / 2.54: the same formula, so the same result exactly.
func TestCentimetersToInches_parity(t *testing.T) {
	for _, cm := range []float64{100, 165.1, 182.88} {
		got, err := CentimetersToInches(cm)
		if err != nil || got != cm/2.54 {
			t.Errorf("CentimetersToInches(%v) = (%v, %v), want %v", cm, got, err, cm/2.54)
		}
	}
}

func TestCentimetersToInches_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := CentimetersToInches(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("CentimetersToInches(%v) error = %v, want ErrNotFinite", v, err)
		}
	}
	if _, err := CentimetersToInches(-1); !errors.Is(err, ErrNegative) {
		t.Errorf("CentimetersToInches(-1) error = %v, want ErrNegative", err)
	}
	if got, err := CentimetersToInches(math.MaxFloat64); err != nil || math.IsInf(got, 0) {
		t.Errorf("CentimetersToInches(MaxFloat64) = (%v, %v), want finite", got, err)
	}
}

// Round trip through InchesToCentimeters, tolerance 2⁻⁵⁰ relative (two roundings).
func TestCentimetersToInches_roundTrip(t *testing.T) {
	for cm := 0.01; cm < 1e7; cm *= 1.41 {
		in, err := CentimetersToInches(cm)
		if err != nil {
			t.Fatal(err)
		}
		back, err := InchesToCentimeters(in)
		if err != nil || !nearRel(back, cm, 0x1p-50) {
			t.Fatalf("round trip %v cm → %v in → (%v, %v)", cm, in, back, err)
		}
	}
}

func TestCentimetersToInches_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkLength, _ = CentimetersToInches(172) }); a != 0 {
		t.Errorf("CentimetersToInches allocates %v times, want 0", a)
	}
}

func BenchmarkCentimetersToInches(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkLength, _ = CentimetersToInches(172)
	}
}
