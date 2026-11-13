package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Yang & Tobin 1991's f/VT: 30 / 0.3 = 100 breaths/min/L. Tolerance 1e-12 relative.
func TestRapidShallowBreathingIndex_handExample(t *testing.T) {
	if got, err := RapidShallowBreathingIndex(30, 0.3); err != nil || math.Abs(got-100) > 1e-12*100 {
		t.Errorf("got (%v, %v), want 100", got, err)
	}
}

func TestRapidShallowBreathingIndex_errors(t *testing.T) {
	tests := []struct {
		name    string
		rr, vt  float64
		wantErr error
	}{
		{"zero RR", 0, 0.4, ErrOutOfRange},
		{"zero VT", 20, 0, ErrOutOfRange},
		{"negative VT", 20, -0.1, ErrOutOfRange},
		{"NaN RR", math.NaN(), 0.4, ErrNotFinite},
		{"+Inf VT", 20, math.Inf(1), ErrNotFinite},
		{"-Inf RR", math.Inf(-1), 0.4, ErrNotFinite},
		{"overflow", math.MaxFloat64, 0.5, ErrDomain},
	}
	for _, tc := range tests {
		if _, err := RapidShallowBreathingIndex(tc.rr, tc.vt); !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestRapidShallowBreathingIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = RapidShallowBreathingIndex(30, 0.3) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRapidShallowBreathingIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = RapidShallowBreathingIndex(30, 0.3)
	}
}
