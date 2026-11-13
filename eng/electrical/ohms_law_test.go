package electrical

import (
	"errors"
	"math"
	"testing"
)

// Hand computations from V = IR and P = IV = I²R = V²/R (OpenStax University Physics Vol. 2, §9.4-9.5):
// the 12 V, 2 A, 6 Ω, 24 W circuit from every pair; 120 V, 60 W → 0.5 A, 240 Ω.
// Tolerance 1e-12 relative: at most a square root and two operations.
func near(a, b float64) bool { return math.Abs(a-b) <= 1e-12*math.Max(1, math.Abs(b)) }

func sameCircuit(c, w Circuit) bool {
	return near(c.Volts, w.Volts) && near(c.Amperes, w.Amperes) && near(c.Ohms, w.Ohms) && near(c.Watts, w.Watts)
}

func TestOhmsLaw_everyPair(t *testing.T) {
	want := Circuit{12, 2, 6, 24}
	type pairFn func(a, b float64) (Circuit, error)
	cases := []struct {
		name string
		f    pairFn
		a, b float64
	}{
		{"V,I", OhmsLawFromVoltageCurrent, 12, 2},
		{"V,R", OhmsLawFromVoltageResistance, 12, 6},
		{"V,P", OhmsLawFromVoltagePower, 12, 24},
		{"I,R", OhmsLawFromCurrentResistance, 2, 6},
		{"I,P", OhmsLawFromCurrentPower, 2, 24},
		{"R,P", OhmsLawFromResistancePower, 6, 24},
	}
	for _, c := range cases {
		got, err := c.f(c.a, c.b)
		if err != nil || !sameCircuit(got, want) {
			t.Errorf("%s: got %+v %v, want %+v", c.name, got, err, want)
		}
	}
	if got, err := OhmsLawFromVoltagePower(120, 60); err != nil || !sameCircuit(got, Circuit{120, 0.5, 240, 60}) {
		t.Errorf("120 V 60 W: %+v %v", got, err)
	}
}

func TestOhmsLawFromVoltageCurrent_errors(t *testing.T) {
	if _, err := OhmsLawFromVoltageCurrent(12, 0); !errors.Is(err, ErrZero) {
		t.Errorf("zero current: %v", err)
	}
	if _, err := OhmsLawFromVoltageCurrent(-1, 1); !errors.Is(err, ErrNegative) {
		t.Errorf("negative: %v", err)
	}
	if _, err := OhmsLawFromVoltageCurrent(math.NaN(), 1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	if _, err := OhmsLawFromVoltageCurrent(math.MaxFloat64, 2); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
}

func TestOhmsLaw_zeroDivisors(t *testing.T) {
	if _, err := OhmsLawFromVoltageResistance(12, 0); !errors.Is(err, ErrZero) {
		t.Errorf("V,R: %v", err)
	}
	if _, err := OhmsLawFromVoltagePower(0, 24); !errors.Is(err, ErrZero) {
		t.Errorf("V,P V=0: %v", err)
	}
	if _, err := OhmsLawFromVoltagePower(12, 0); !errors.Is(err, ErrZero) {
		t.Errorf("V,P P=0: %v", err)
	}
	if _, err := OhmsLawFromCurrentPower(0, 24); !errors.Is(err, ErrZero) {
		t.Errorf("I,P: %v", err)
	}
	if _, err := OhmsLawFromResistancePower(0, 24); !errors.Is(err, ErrZero) {
		t.Errorf("R,P: %v", err)
	}
	if _, err := OhmsLawFromCurrentResistance(math.Inf(1), 1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("I,R Inf: %v", err)
	}
	if got, err := OhmsLawFromCurrentResistance(0, 6); err != nil || got != (Circuit{0, 0, 6, 0}) {
		t.Errorf("I,R zero current: %+v %v", got, err)
	}
	if _, err := OhmsLawFromCurrentResistance(math.MaxFloat64, 2); !errors.Is(err, ErrOverflow) {
		t.Errorf("I,R overflow: %v", err)
	}
}

func BenchmarkOhmsLawFromVoltageCurrent(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = OhmsLawFromVoltageCurrent(12, 2)
	}
}

func BenchmarkOhmsLawFromVoltageResistance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = OhmsLawFromVoltageResistance(12, 6)
	}
}

func BenchmarkOhmsLawFromVoltagePower(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = OhmsLawFromVoltagePower(12, 24)
	}
}

func BenchmarkOhmsLawFromCurrentResistance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = OhmsLawFromCurrentResistance(2, 6)
	}
}

func BenchmarkOhmsLawFromCurrentPower(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = OhmsLawFromCurrentPower(2, 24)
	}
}

func BenchmarkOhmsLawFromResistancePower(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = OhmsLawFromResistancePower(6, 24)
	}
}
