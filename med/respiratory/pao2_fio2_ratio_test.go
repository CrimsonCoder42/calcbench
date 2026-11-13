package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkBerlin BerlinARDSOxygenation

// HAND-DERIVED: PaO₂ 80 / FiO₂ 0.4 = 200 mmHg; 60 / 0.21 = 285.714…; tolerance 1e-12 relative (one division).
func TestPaO2FiO2Ratio_handExample(t *testing.T) {
	if got, err := PaO2FiO2Ratio(80, 0.4); err != nil || math.Abs(got-200) > 1e-12*200 {
		t.Errorf("got (%v, %v), want 200", got, err)
	}
	if got, err := PaO2FiO2Ratio(60, 0.21); err != nil || math.Abs(got-60/0.21) > 1e-12*got {
		t.Errorf("got (%v, %v), want %v", got, err, 60/0.21)
	}
}

func TestPaO2FiO2Ratio_errors(t *testing.T) {
	tests := []struct {
		name     string
		pa, fio2 float64
		wantErr  error
	}{
		{"zero PaO2", 0, 0.4, ErrOutOfRange},
		{"FiO2 0.2", 80, 0.2, ErrOutOfRange},
		{"FiO2 1.01", 80, 1.01, ErrOutOfRange},
		{"NaN PaO2", math.NaN(), 0.4, ErrNotFinite},
		{"+Inf PaO2", math.Inf(1), 0.4, ErrNotFinite},
		{"-Inf FiO2", 80, math.Inf(-1), ErrNotFinite},
	}
	for _, tc := range tests {
		if _, err := PaO2FiO2Ratio(tc.pa, tc.fio2); !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
	for _, f := range []float64{0.21, 1.0} {
		if _, err := PaO2FiO2Ratio(80, f); err != nil {
			t.Errorf("FiO2 %v: err = %v, want nil", f, err)
		}
	}
}

// Berlin bounds (Ranieri 2012 abstract): severe ≤ 100, moderate (100, 200], mild (200, 300]; both sides of each.
func TestBerlinARDSOxygenationCategory_edges(t *testing.T) {
	cases := []struct {
		pf   float64
		want BerlinARDSOxygenation
	}{
		{1, BerlinSevere}, {100, BerlinSevere}, {100.01, BerlinModerate}, {200, BerlinModerate},
		{200.01, BerlinMild}, {300, BerlinMild}, {300.01, BerlinAbove300}, {600, BerlinAbove300},
	}
	for _, c := range cases {
		if got, err := BerlinARDSOxygenationCategory(c.pf); err != nil || got != c.want {
			t.Errorf("%v: got (%q, %v), want %q", c.pf, got, err, c.want)
		}
	}
	for _, bad := range []struct {
		pf      float64
		wantErr error
	}{{0, ErrOutOfRange}, {-5, ErrOutOfRange}, {math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}} {
		if _, err := BerlinARDSOxygenationCategory(bad.pf); !errors.Is(err, bad.wantErr) {
			t.Errorf("%v: err = %v, want %v", bad.pf, err, bad.wantErr)
		}
	}
}

func TestPaO2FiO2Ratio_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = PaO2FiO2Ratio(80, 0.4) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkBerlin, _ = BerlinARDSOxygenationCategory(150) }); n != 0 {
		t.Errorf("category allocs = %v, want 0", n)
	}
}

func BenchmarkPaO2FiO2Ratio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PaO2FiO2Ratio(80, 0.4)
	}
}

func BenchmarkBerlinARDSOxygenationCategory(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBerlin, _ = BerlinARDSOxygenationCategory(150)
	}
}
