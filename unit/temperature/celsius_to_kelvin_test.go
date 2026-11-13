package temperature

import "testing"

// NIST SP 811 App. B.9: T/K = t/°C + 273.15.
func TestCelsiusToKelvin_example(t *testing.T) {
	tests := []struct{ c, want float64 }{{0, 273.15}, {100, 373.15}, {-273.15, 0}, {37, 310.15}}
	for _, tc := range tests {
		got, err := CelsiusToKelvin(tc.c)
		if err != nil || !near(got, tc.want, 1e-12) {
			t.Errorf("CelsiusToKelvin(%v) = (%v, %v), want %v", tc.c, got, err, tc.want)
		}
	}
}

func TestCelsiusToKelvin_errors(t *testing.T) {
	checkErrors(t, "CelsiusToKelvin", CelsiusToKelvin, -273.16)
}

func TestCelsiusToKelvin_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = CelsiusToKelvin(37) }); n != 0 {
		t.Errorf("CelsiusToKelvin allocates %v times, want 0", n)
	}
}

func BenchmarkCelsiusToKelvin(b *testing.B) {
	b.ReportAllocs()
	c := 37.0
	for b.Loop() {
		sinkFloat, _ = CelsiusToKelvin(c)
	}
}
