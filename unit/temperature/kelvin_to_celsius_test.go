package temperature

import "testing"

// NIST SP 811 App. B.9: t/°C = T/K − 273.15.
func TestKelvinToCelsius_example(t *testing.T) {
	tests := []struct{ k, want float64 }{{273.15, 0}, {373.15, 100}, {0, -273.15}, {310.15, 37}}
	for _, tc := range tests {
		got, err := KelvinToCelsius(tc.k)
		if err != nil || !near(got, tc.want, 1e-12) {
			t.Errorf("KelvinToCelsius(%v) = (%v, %v), want %v", tc.k, got, err, tc.want)
		}
	}
}

func TestKelvinToCelsius_errors(t *testing.T) {
	checkErrors(t, "KelvinToCelsius", KelvinToCelsius, -1e-9)
}

func TestKelvinToCelsius_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = KelvinToCelsius(310.15) }); n != 0 {
		t.Errorf("KelvinToCelsius allocates %v times, want 0", n)
	}
}

func BenchmarkKelvinToCelsius(b *testing.B) {
	b.ReportAllocs()
	k := 310.15
	for b.Loop() {
		sinkFloat, _ = KelvinToCelsius(k)
	}
}
