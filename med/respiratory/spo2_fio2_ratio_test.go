package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED: SpO₂ 94 / FiO₂ 0.4 = 235 (the S/F the Rice 2007 abstract pairs with P/F 200). Tolerance 1e-12 relative.
func TestSpO2FiO2Ratio_handExample(t *testing.T) {
	if got, err := SpO2FiO2Ratio(94, 0.4); err != nil || math.Abs(got-235) > 1e-12*235 {
		t.Errorf("got (%v, %v), want 235", got, err)
	}
}

func TestSpO2FiO2Ratio_edges(t *testing.T) {
	if _, err := SpO2FiO2Ratio(97, 0.21); err != nil {
		t.Errorf("SpO2 97: err = %v, want nil", err)
	}
	tests := []struct {
		name       string
		spo2, fio2 float64
		wantErr    error
	}{
		{"SpO2 97.1 (above the source's ≤ 97 %)", 97.1, 0.4, ErrOutOfRange},
		{"SpO2 0", 0, 0.4, ErrOutOfRange},
		{"FiO2 0.2", 94, 0.2, ErrOutOfRange},
		{"FiO2 1.01", 94, 1.01, ErrOutOfRange},
		{"NaN SpO2", math.NaN(), 0.4, ErrNotFinite},
		{"+Inf FiO2", 94, math.Inf(1), ErrNotFinite},
		{"-Inf SpO2", math.Inf(-1), 0.4, ErrNotFinite},
	}
	for _, tc := range tests {
		if _, err := SpO2FiO2Ratio(tc.spo2, tc.fio2); !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

// Rice 2007 abstract: "S/F = 64 + 0.84 x (P/F)". 200 → 232, 300 → 316; the abstract's fitted pairs (235, 315)
// are within 4 of the line, as a regression's rounded report is. Tolerance 1e-9 (two exact-ish operations).
func TestSpO2FiO2FromPaO2FiO2Rice_equation(t *testing.T) {
	for _, c := range []struct{ pf, want float64 }{{200, 232}, {300, 316}, {100, 148}} {
		if got, err := SpO2FiO2FromPaO2FiO2Rice(c.pf); err != nil || math.Abs(got-c.want) > 1e-9 {
			t.Errorf("%v: got (%v, %v), want %v", c.pf, got, err, c.want)
		}
	}
	for _, bad := range []struct {
		pf      float64
		wantErr error
	}{{0, ErrOutOfRange}, {-1, ErrOutOfRange}, {math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}} {
		if _, err := SpO2FiO2FromPaO2FiO2Rice(bad.pf); !errors.Is(err, bad.wantErr) {
			t.Errorf("%v: err = %v, want %v", bad.pf, err, bad.wantErr)
		}
	}
	if got, err := SpO2FiO2FromPaO2FiO2Rice(math.MaxFloat64); err != nil || math.IsInf(got, 0) {
		t.Errorf("MaxFloat64: got (%v, %v), want finite", got, err)
	}
}

func TestSpO2FiO2Ratio_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = SpO2FiO2Ratio(94, 0.4) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = SpO2FiO2FromPaO2FiO2Rice(200) }); n != 0 {
		t.Errorf("Rice allocs = %v, want 0", n)
	}
}

func BenchmarkSpO2FiO2Ratio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SpO2FiO2Ratio(94, 0.4)
	}
}

func BenchmarkSpO2FiO2FromPaO2FiO2Rice(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SpO2FiO2FromPaO2FiO2Rice(200)
	}
}
