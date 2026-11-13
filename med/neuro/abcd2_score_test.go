package neuro

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Johnston 2007 abstract points; each edge on both sides.
func TestABCD2Score_example(t *testing.T) {
	tests := []struct {
		name             string
		age              float64
		bp, weak, speech bool
		dur              float64
		diabetes         bool
		want             int
	}{
		{"none", 59.9, false, false, false, 9.9, false, 0},
		{"age 60", 60, false, false, false, 0, false, 1},
		{"bp", 30, true, false, false, 0, false, 1},
		{"weakness 2", 30, false, true, false, 0, false, 2},
		{"speech 1", 30, false, false, true, 0, false, 1},
		{"weakness wins over speech", 30, false, true, true, 0, false, 2},
		{"duration 10", 30, false, false, false, 10, false, 1},
		{"duration 59.9", 30, false, false, false, 59.9, false, 1},
		{"duration 60", 30, false, false, false, 60, false, 2},
		{"diabetes", 30, false, false, false, 0, true, 1},
		{"worked: 70, BP, weakness, 45 min", 70, true, true, false, 45, false, 5},
		{"max", 80, true, true, true, 120, true, 7},
	}
	for _, tc := range tests {
		if got, err := ABCD2Score(tc.age, tc.bp, tc.weak, tc.speech, tc.dur, tc.diabetes); err != nil || got != tc.want {
			t.Errorf("%s: ABCD2Score = (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

func TestABCD2Score_errors(t *testing.T) {
	for _, tc := range []struct {
		age, dur float64
		want     error
	}{
		{math.NaN(), 0, ErrNotFinite}, {math.Inf(1), 0, ErrNotFinite}, {50, math.Inf(-1), ErrNotFinite},
		{-1, 0, ErrOutOfRange}, {50, -0.1, ErrOutOfRange},
	} {
		if _, err := ABCD2Score(tc.age, false, false, false, tc.dur, false); !errors.Is(err, tc.want) {
			t.Errorf("ABCD2Score(%v, …, %v) error = %v, want %v", tc.age, tc.dur, err, tc.want)
		}
	}
}

func TestABCD2Score_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = ABCD2Score(70, true, true, false, 45, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkABCD2Score(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = ABCD2Score(70, true, true, false, 45, false)
	}
}

// Johnston 2007 abstract strata, each edge on both sides.
func TestABCD2RiskGroupOf_example(t *testing.T) {
	tests := []struct {
		total int
		want  ABCD2RiskGroup
	}{{0, ABCD2LowRisk}, {3, ABCD2LowRisk}, {4, ABCD2ModerateRisk}, {5, ABCD2ModerateRisk}, {6, ABCD2HighRisk}, {7, ABCD2HighRisk}}
	for _, tc := range tests {
		if got, err := ABCD2RiskGroupOf(tc.total); err != nil || got != tc.want {
			t.Errorf("ABCD2RiskGroupOf(%d) = (%q, %v), want %q", tc.total, got, err, tc.want)
		}
	}
	for _, bad := range []int{-1, 8} {
		if _, err := ABCD2RiskGroupOf(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("ABCD2RiskGroupOf(%d) error = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestABCD2RiskGroupOf_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { g, _ := ABCD2RiskGroupOf(5); sinkString = string(g) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkABCD2RiskGroupOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := ABCD2RiskGroupOf(5)
		sinkString = string(g)
	}
}
