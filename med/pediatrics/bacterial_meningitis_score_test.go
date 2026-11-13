package pediatrics

import (
	"errors"
	"math"
	"testing"
)

// Nigrovic 2002 abstract: 2 points for a positive Gram stain, 1 point for each of CSF ANC ≥ 1000 cells/mm³,
// CSF protein ≥ 80 mg/dL, peripheral ANC ≥ 10,000 cells/mm³ and seizure at or before presentation (range
// 0–6, Lee 2015 PMC4543184). Each threshold is tested on both sides.
func TestBacterialMeningitisScore_points(t *testing.T) {
	cases := []struct {
		gram                bool
		csfANC, prot, blood float64
		seizure             bool
		want                int
	}{
		{false, 0, 0, 0, false, 0},
		{false, 999.9, 79.9, 9999.9, false, 0},
		{true, 0, 0, 0, false, 2},
		{false, 1000, 0, 0, false, 1},
		{false, 0, 80, 0, false, 1},
		{false, 0, 0, 10000, false, 1},
		{false, 0, 0, 0, true, 1},
		{true, 1000, 79.9, 9999, false, 3},
		{true, 1000, 80, 10000, true, 6},
		{false, 5000, 150, 20000, true, 4},
	}
	for _, c := range cases {
		got, err := BacterialMeningitisScore(c.gram, c.csfANC, c.prot, c.blood, c.seizure)
		if err != nil || got != c.want {
			t.Errorf("BacterialMeningitisScore(%+v) = (%d, %v), want %d", c, got, err, c.want)
		}
	}
}

// parity: ledger X3 bacterial_meningitis_score (W2_B6 row); the entry's rule is "very low risk" when no
// criterion is present, i.e. score 0 here.
func TestBacterialMeningitisScore_parity(t *testing.T) {
	if got, _ := BacterialMeningitisScore(false, 200, 45, 8000, false); got != 0 {
		t.Errorf("no criterion present: score %d, want 0", got)
	}
	if got, _ := BacterialMeningitisScore(false, 200, 45, 8000, true); got == 0 {
		t.Errorf("seizure present: score 0, want > 0")
	}
}

func TestBacterialMeningitisScore_errors(t *testing.T) {
	cases := []struct {
		csfANC, prot, blood float64
		want                error
	}{
		{math.NaN(), 0, 0, ErrNotFinite}, {0, math.Inf(1), 0, ErrNotFinite}, {0, 0, math.Inf(-1), ErrNotFinite},
		{-1, 0, 0, ErrOutOfRange}, {1e6 + 1, 0, 0, ErrOutOfRange}, {0, -0.1, 0, ErrOutOfRange},
		{0, 5000.1, 0, ErrOutOfRange}, {0, 0, -1, ErrOutOfRange}, {0, 0, 1e6 + 1, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := BacterialMeningitisScore(false, c.csfANC, c.prot, c.blood, false); !errors.Is(err, c.want) {
			t.Errorf("BacterialMeningitisScore(%+v) error = %v, want %v", c, err, c.want)
		}
	}
	if got, err := BacterialMeningitisScore(true, 1e6, 5000, 1e6, true); err != nil || got != 6 {
		t.Errorf("upper range edge = (%d, %v), want (6, nil)", got, err)
	}
}

func TestBacterialMeningitisScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = BacterialMeningitisScore(true, 1200, 90, 12000, false)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkBacterialMeningitisScore(b *testing.B) {
	b.ReportAllocs()
	csf, prot, blood := 1200.0, 90.0, 12000.0
	for b.Loop() {
		sinkInt, _ = BacterialMeningitisScore(true, csf, prot, blood, false)
	}
}
