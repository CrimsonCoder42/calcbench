package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from PALICC-2 Table 4 footnote (PMID 36661420): OI = 15 × 60 / 60 = 15; OSI = 12 × 50 / 92.
// Tolerance 1e-12 relative (one multiply, one divide).
func TestOxygenationIndex_handExample(t *testing.T) {
	if got, err := OxygenationIndex(15, 60, 60); err != nil || math.Abs(got-15) > 1e-12*15 {
		t.Errorf("OI got (%v, %v), want 15", got, err)
	}
	// The PALICC-2 severity edge OI 16: MAP 16, FiO₂ 100 %, PaO₂ 100.
	if got, err := OxygenationIndex(16, 100, 100); err != nil || math.Abs(got-16) > 1e-12*16 {
		t.Errorf("OI got (%v, %v), want 16", got, err)
	}
	want := 12.0 * 50 / 92
	if got, err := OxygenSaturationIndex(12, 50, 92); err != nil || math.Abs(got-want) > 1e-12*want {
		t.Errorf("OSI got (%v, %v), want %v", got, err, want)
	}
}

func TestOxygenationIndex_errors(t *testing.T) {
	tests := []struct {
		name       string
		osi        bool
		mapP, f, d float64
		wantErr    error
	}{
		{"OI zero MAP", false, 0, 60, 60, ErrOutOfRange},
		{"OI FiO2 20.9", false, 15, 20.9, 60, ErrOutOfRange},
		{"OI FiO2 100.1", false, 15, 100.1, 60, ErrOutOfRange},
		{"OI zero PaO2", false, 15, 60, 0, ErrOutOfRange},
		{"OI NaN MAP", false, math.NaN(), 60, 60, ErrNotFinite},
		{"OI +Inf PaO2", false, 15, 60, math.Inf(1), ErrNotFinite},
		{"OI overflow", false, math.MaxFloat64, 100, 60, ErrDomain},
		{"OSI SpO2 97.1", true, 12, 50, 97.1, ErrOutOfRange},
		{"OSI SpO2 0", true, 12, 50, 0, ErrOutOfRange},
		{"OSI -Inf FiO2", true, 12, math.Inf(-1), 92, ErrNotFinite},
		{"OSI NaN SpO2", true, 12, 50, math.NaN(), ErrNotFinite},
		{"OSI zero MAP", true, 0, 50, 92, ErrOutOfRange},
		{"OSI overflow", true, math.MaxFloat64, 100, 50, ErrDomain},
	}
	for _, tc := range tests {
		var err error
		if tc.osi {
			_, err = OxygenSaturationIndex(tc.mapP, tc.f, tc.d)
		} else {
			_, err = OxygenationIndex(tc.mapP, tc.f, tc.d)
		}
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
	if _, err := OxygenSaturationIndex(12, 21, 97); err != nil {
		t.Errorf("OSI at SpO2 97, FiO2 21: err = %v, want nil", err)
	}
}

func TestOxygenationIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = OxygenationIndex(15, 60, 60) }); n != 0 {
		t.Errorf("OI allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = OxygenSaturationIndex(12, 50, 92) }); n != 0 {
		t.Errorf("OSI allocs = %v, want 0", n)
	}
}

func BenchmarkOxygenationIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = OxygenationIndex(15, 60, 60)
	}
}

func BenchmarkOxygenSaturationIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = OxygenSaturationIndex(12, 50, 92)
	}
}
