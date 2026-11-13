package temperature

import (
	"errors"
	"math"
	"testing"
)

var sinkFloat float64

// near reports whether got is within tol of want. Tolerance 1e-12 absolute is used for exact relations: the
// constants 1.8, 32 and 273.15 are not all exact in binary, so the last bits may differ.
func near(got, want, tol float64) bool { return math.Abs(got-want) <= tol }

type convFunc func(float64) (float64, error)

func checkErrors(t *testing.T, name string, f convFunc, below float64) {
	t.Helper()
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := f(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("%s(%v) error = %v, want ErrNotFinite", name, v, err)
		}
	}
	if _, err := f(below); !errors.Is(err, ErrBelowAbsoluteZero) {
		t.Errorf("%s(%v) error = %v, want ErrBelowAbsoluteZero", name, below, err)
	}
}

// NIST SP 811 App. B.9: t/°C = (t/°F − 32)/1.8. Fixed points: water boils 212 °F = 100 °C, freezes 32 °F = 0 °C;
// −40 is the same on both scales; absolute zero −459.67 °F = −273.15 °C.
func TestFahrenheitToCelsius_example(t *testing.T) {
	tests := []struct{ f, want float64 }{
		{212, 100}, {32, 0}, {-40, -40}, {98.6, 37}, {-459.67, -273.15}, {0, -160.0 / 9},
	}
	for _, tc := range tests {
		got, err := FahrenheitToCelsius(tc.f)
		if err != nil || !near(got, tc.want, 1e-12) {
			t.Errorf("FahrenheitToCelsius(%v) = (%v, %v), want %v", tc.f, got, err, tc.want)
		}
	}
}

// parity: ledger temp_converter, source cases 1–6. The source shows °C to one decimal; the comparison is at
// that display precision.
func TestFahrenheitToCelsius_parity(t *testing.T) {
	tests := []struct{ f, want float64 }{
		{94.9, 34.9},  // source case 1
		{95, 35.0},    // source case 2
		{100.3, 37.9}, // source case 3
		{100.4, 38.0}, // source case 4
		{60, 15.6},    // source case 5
		{120, 48.9},   // source case 6
	}
	for _, tc := range tests {
		got, err := FahrenheitToCelsius(tc.f)
		if err != nil || math.Round(got*10)/10 != tc.want {
			t.Errorf("FahrenheitToCelsius(%v) = (%v, %v), want %v at 1 dp", tc.f, got, err, tc.want)
		}
	}
}

func TestFahrenheitToCelsius_errors(t *testing.T) {
	checkErrors(t, "FahrenheitToCelsius", FahrenheitToCelsius, -459.68)
	if _, err := FahrenheitToCelsius(-459.67); err != nil {
		t.Errorf("FahrenheitToCelsius(-459.67) error = %v, want nil (absolute zero itself)", err)
	}
}

// Round trip through CelsiusToFahrenheit over a grid, tolerance 1e-9 absolute (two roundings of 1.8).
func TestFahrenheitToCelsius_roundTrip(t *testing.T) {
	for f := -459.0; f <= 1000; f += 0.7 {
		c, err := FahrenheitToCelsius(f)
		if err != nil {
			t.Fatal(err)
		}
		back, err := CelsiusToFahrenheit(c)
		if err != nil || !near(back, f, 1e-9) {
			t.Fatalf("round trip %v → %v → (%v, %v)", f, c, back, err)
		}
	}
}

func TestFahrenheitToCelsius_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = FahrenheitToCelsius(98.6) }); n != 0 {
		t.Errorf("FahrenheitToCelsius allocates %v times, want 0", n)
	}
}

func BenchmarkFahrenheitToCelsius(b *testing.B) {
	b.ReportAllocs()
	f := 98.6
	for b.Loop() {
		sinkFloat, _ = FahrenheitToCelsius(f)
	}
}
