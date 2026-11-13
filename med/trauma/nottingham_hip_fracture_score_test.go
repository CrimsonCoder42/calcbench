package trauma

import (
	"errors"
	"testing"
)

var sinkNHFS int

// Lisk 2020 (PMC7497295) and Karres 2023 (PMC10293368) Table 2, both citing Maxwell 2008: age < 66 = 0,
// 66–85 = 3, ≥ 86 = 4, each other variable 1, maximum 10. Both sides of every edge.
func TestNottinghamHipFractureScore_points(t *testing.T) {
	ages := []struct {
		age  float64
		want int
	}{{0, 0}, {65.9, 0}, {66, 3}, {85, 3}, {85.9, 3}, {86, 4}, {130, 4}}
	for _, tc := range ages {
		got, err := NottinghamHipFractureScore(tc.age, false, 12, 10, false, 0, false)
		if err != nil || got != tc.want {
			t.Errorf("age %v → %d (%v), want %d", tc.age, got, err, tc.want)
		}
	}
	// haemoglobin ≤ 10 g/dL scores; 10.01 does not.
	for _, tc := range []struct {
		hb   float64
		want int
	}{{10, 1}, {10.01, 0}, {3, 1}} {
		got, _ := NottinghamHipFractureScore(50, false, tc.hb, 10, false, 0, false)
		if got != tc.want {
			t.Errorf("Hb %v → %d, want %d", tc.hb, got, tc.want)
		}
	}
	// AMTS ≤ 6 scores; 7 does not. Comorbidities ≥ 2 score; 1 does not.
	for _, tc := range []struct{ amts, com, want int }{{6, 0, 1}, {7, 0, 0}, {0, 0, 1}, {10, 1, 0}, {10, 2, 1}, {10, 9, 1}} {
		got, _ := NottinghamHipFractureScore(50, false, 12, tc.amts, false, tc.com, false)
		if got != tc.want {
			t.Errorf("AMTS %d comorbidities %d → %d, want %d", tc.amts, tc.com, got, tc.want)
		}
	}
	if got, _ := NottinghamHipFractureScore(90, true, 8, 3, true, 3, true); got != 10 {
		t.Errorf("all items → %d, want 10 (the printed maximum)", got)
	}
	if got, _ := NottinghamHipFractureScore(50, true, 12, 10, false, 0, false); got != 1 {
		t.Errorf("male only → %d, want 1", got)
	}
	if got, _ := NottinghamHipFractureScore(50, false, 12, 10, true, 0, true); got != 2 {
		t.Errorf("institution + malignancy → %d, want 2", got)
	}
}

// HAND-DERIVED from the points: age 90 (4) + Hb 9.5 g/dL (1) + 2 comorbidities (1) = 6;
// SI: age 70 (3) + male + Hb 100 g/L + AMTS 4 + institution + malignancy = 8.
func TestNottinghamHipFractureScore_source(t *testing.T) {
	if got, err := NottinghamHipFractureScore(90, false, 9.5, 8, false, 2, false); err != nil || got != 6 {
		t.Errorf("got %d (%v), want 6", got, err)
	}
	if got, err := NottinghamHipFractureScoreSI(70, true, 100, 4, true, 1, true); err != nil || got != 8 {
		t.Errorf("SI got %d (%v), want 8", got, err)
	}
}

// The two unit systems agree away from the threshold and at it (10 g/dL and 100 g/L are each printed by a source).
func TestNottinghamHipFractureScoreSI_agrees(t *testing.T) {
	for _, hb := range []float64{8, 10, 10.5, 14} {
		a, _ := NottinghamHipFractureScore(75, true, hb, 9, false, 1, false)
		b, _ := NottinghamHipFractureScoreSI(75, true, hb*10, 9, false, 1, false)
		if a != b {
			t.Errorf("Hb %v g/dL → %d, %v g/L → %d", hb, a, hb*10, b)
		}
	}
}

// parity: ledger nhfs, source case 1. No _test.go or JSON example; the entry takes pre-scored items
// (age 0–2 index of the options "<66 (0)", "66–85 (3)", "≥86 (4)"). Its options sum to the same points:
// age ≥ 86, male, Hb ≤ 10, AMTS ≤ 6, institution, ≥ 2 comorbidities, malignancy → 10.
func TestNottinghamHipFractureScore_parity(t *testing.T) {
	if got, _ := NottinghamHipFractureScore(88, true, 9, 5, true, 2, true); got != 10 {
		t.Errorf("got %d, want 10", got)
	}
}

func TestNottinghamHipFractureScore_errors(t *testing.T) {
	for _, v := range nonFinite {
		if _, err := NottinghamHipFractureScore(v, false, 12, 10, false, 0, false); !errors.Is(err, ErrNotFinite) {
			t.Errorf("age %v: err = %v, want ErrNotFinite", v, err)
		}
		if _, err := NottinghamHipFractureScore(70, false, v, 10, false, 0, false); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Hb %v: err = %v, want ErrNotFinite", v, err)
		}
		if _, err := NottinghamHipFractureScoreSI(70, false, v, 10, false, 0, false); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Hb SI %v: err = %v, want ErrNotFinite", v, err)
		}
	}
	bad := []func() error{
		func() error { _, e := NottinghamHipFractureScore(-0.1, false, 12, 10, false, 0, false); return e },
		func() error { _, e := NottinghamHipFractureScore(130.1, false, 12, 10, false, 0, false); return e },
		func() error { _, e := NottinghamHipFractureScore(70, false, 0, 10, false, 0, false); return e },
		func() error { _, e := NottinghamHipFractureScoreSI(70, false, -1, 10, false, 0, false); return e },
		func() error { _, e := NottinghamHipFractureScore(70, false, 12, -1, false, 0, false); return e },
		func() error { _, e := NottinghamHipFractureScore(70, false, 12, 11, false, 0, false); return e },
		func() error { _, e := NottinghamHipFractureScore(70, false, 12, 10, false, -1, false); return e },
	}
	for i, f := range bad {
		if err := f(); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("case %d: err = %v, want ErrOutOfRange", i, err)
		}
	}
}

func TestNottinghamHipFractureScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkNHFS, _ = NottinghamHipFractureScore(90, true, 9, 5, true, 2, true) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkNHFS, _ = NottinghamHipFractureScoreSI(90, true, 90, 5, true, 2, true) }); n != 0 {
		t.Errorf("SI allocs = %v, want 0", n)
	}
}

func BenchmarkNottinghamHipFractureScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkNHFS, _ = NottinghamHipFractureScore(90, true, 9, 5, true, 2, true)
	}
}

func BenchmarkNottinghamHipFractureScoreSI(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkNHFS, _ = NottinghamHipFractureScoreSI(90, true, 90, 5, true, 2, true)
	}
}
