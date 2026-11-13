package temperature

import (
	"errors"
	"testing"
)

// Inverse of NIST SP 811 App. B.9: 100 °C = 212 °F, 0 °C = 32 °F, −40 = −40, 37 °C = 98.6 °F, −273.15 °C = −459.67 °F.
func TestCelsiusToFahrenheit_example(t *testing.T) {
	tests := []struct{ c, want float64 }{{100, 212}, {0, 32}, {-40, -40}, {37, 98.6}, {-273.15, -459.67}}
	for _, tc := range tests {
		got, err := CelsiusToFahrenheit(tc.c)
		if err != nil || !near(got, tc.want, 1e-12) {
			t.Errorf("CelsiusToFahrenheit(%v) = (%v, %v), want %v", tc.c, got, err, tc.want)
		}
	}
}

func TestCelsiusToFahrenheit_errors(t *testing.T) {
	checkErrors(t, "CelsiusToFahrenheit", CelsiusToFahrenheit, -273.16)
}

func TestCelsiusToFahrenheit_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = CelsiusToFahrenheit(37) }); n != 0 {
		t.Errorf("CelsiusToFahrenheit allocates %v times, want 0", n)
	}
}

func BenchmarkCelsiusToFahrenheit(b *testing.B) {
	b.ReportAllocs()
	c := 37.0
	for b.Loop() {
		sinkFloat, _ = CelsiusToFahrenheit(c)
	}
}

// A finite Celsius value whose Fahrenheit value is beyond float64 is ErrOverflow (review finding).
func TestCelsiusToFahrenheit_overflow(t *testing.T) {
	if _, err := CelsiusToFahrenheit(1.7976931348623157e308); !errors.Is(err, ErrOverflow) {
		t.Errorf("CelsiusToFahrenheit(MaxFloat64) err = %v, want ErrOverflow", err)
	}
	// 9e307 × 1.8 = 1.62e308 still fits.
	if _, err := CelsiusToFahrenheit(9e307); err != nil {
		t.Errorf("CelsiusToFahrenheit(9e307) err = %v, want nil", err)
	}
}
