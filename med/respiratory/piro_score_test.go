package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkPIRO PIRORisk

// HAND-DERIVED from Rello 2009 (PMID 19114916), one point per feature: age 75 + shock + bacteremia = 3;
// all eight = 8; none (age 70, not > 70) = 0.
func TestPIROScore_handExample(t *testing.T) {
	if got, err := PIROScore(PIROFindings{AgeYears: 75, Shock: true, Bacteremia: true}); err != nil || got != 3 {
		t.Errorf("got (%d, %v), want 3", got, err)
	}
	all := PIROFindings{true, 80, true, true, true, true, true, true}
	if got, err := PIROScore(all); err != nil || got != 8 {
		t.Errorf("all: got (%d, %v), want 8", got, err)
	}
	if got, err := PIROScore(PIROFindings{AgeYears: 70}); err != nil || got != 0 {
		t.Errorf("age 70: got (%d, %v), want 0", got, err)
	}
	if got, err := PIROScore(PIROFindings{AgeYears: 70.01}); err != nil || got != 1 {
		t.Errorf("age 70.01: got (%d, %v), want 1", got, err)
	}
}

// Levels: "Low, 0-2 points; Mild, 3 points; high, 4 points; Very high, 5-8 points" — both sides of each edge.
func TestPIRORiskLevel_edges(t *testing.T) {
	cases := []struct {
		s    int
		want PIRORisk
	}{{0, PIROLow}, {2, PIROLow}, {3, PIROMild}, {4, PIROHigh}, {5, PIROVeryHigh}, {8, PIROVeryHigh}}
	for _, c := range cases {
		if got, err := PIRORiskLevel(c.s); err != nil || got != c.want {
			t.Errorf("%d: got (%q, %v), want %q", c.s, got, err, c.want)
		}
	}
	for _, bad := range []int{-1, 9} {
		if _, err := PIRORiskLevel(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%d: err = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestPIROScore_errors(t *testing.T) {
	for _, bad := range []struct {
		age     float64
		wantErr error
	}{{-1, ErrOutOfRange}, {math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}, {math.Inf(-1), ErrNotFinite}} {
		if _, err := PIROScore(PIROFindings{AgeYears: bad.age}); !errors.Is(err, bad.wantErr) {
			t.Errorf("age %v: err = %v, want %v", bad.age, err, bad.wantErr)
		}
	}
}

func TestPIROScore_noAlloc(t *testing.T) {
	f := PIROFindings{AgeYears: 75, Shock: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = PIROScore(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkPIRO, _ = PIRORiskLevel(4) }); n != 0 {
		t.Errorf("level allocs = %v, want 0", n)
	}
}

func BenchmarkPIROScore(b *testing.B) {
	b.ReportAllocs()
	f := PIROFindings{AgeYears: 75, Shock: true, Bacteremia: true}
	for b.Loop() {
		sinkInt, _ = PIROScore(f)
	}
}

func BenchmarkPIRORiskLevel(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkPIRO, _ = PIRORiskLevel(4)
	}
}
