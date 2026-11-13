package gastro

import "testing"

// Moon 2019 Table 3 (PMC6361015): each row's points; stool-count edges both sides. HAND-DERIVED totals.
func TestPUCAIScore_table(t *testing.T) {
	cases := []struct {
		name string
		f    PUCAIFindings
		want int
	}{
		{"zero", PUCAIFindings{}, 0},
		{"pain 1", PUCAIFindings{AbdominalPain: 1}, 5},
		{"pain 2", PUCAIFindings{AbdominalPain: 2}, 10},
		{"bleeding 1", PUCAIFindings{RectalBleeding: 1}, 10},
		{"bleeding 2", PUCAIFindings{RectalBleeding: 2}, 20},
		{"bleeding 3", PUCAIFindings{RectalBleeding: 3}, 30},
		{"consistency 1", PUCAIFindings{StoolConsistency: 1}, 5},
		{"consistency 2", PUCAIFindings{StoolConsistency: 2}, 10},
		{"stools 2", PUCAIFindings{StoolsPer24h: 2}, 0},
		{"stools 3", PUCAIFindings{StoolsPer24h: 3}, 5},
		{"stools 5", PUCAIFindings{StoolsPer24h: 5}, 5},
		{"stools 6", PUCAIFindings{StoolsPer24h: 6}, 10},
		{"stools 8", PUCAIFindings{StoolsPer24h: 8}, 10},
		{"stools 9", PUCAIFindings{StoolsPer24h: 9}, 15},
		{"nocturnal", PUCAIFindings{NocturnalStools: true}, 10},
		{"activity 1", PUCAIFindings{ActivityLevel: 1}, 5},
		{"activity 2", PUCAIFindings{ActivityLevel: 2}, 10},
		{"doc example", PUCAIFindings{AbdominalPain: 1, RectalBleeding: 1, StoolConsistency: 1, StoolsPer24h: 4}, 25},
		{"maximum", PUCAIFindings{2, 3, 2, 12, true, 2}, 85},
	}
	for _, c := range cases {
		got, err := PUCAIScore(c.f)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger pucai, source cases 1–7 (pucai_test.go TestPUCAI_MatchesTSX): picker indexes
// (pain, bleeding, consistency, stool-count band, nocturnal, activity) → score; the count band index k is
// given here as a count inside band k (0→1, 1→4, 2→7, 3→10).
func TestPUCAIScore_parity(t *testing.T) {
	count := [4]int{1, 4, 7, 10}
	cases := []struct {
		p, b, c, n, noc, a int
		want               int
	}{
		{0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 5},
		{1, 1, 0, 0, 0, 0, 15},
		{1, 1, 1, 1, 0, 0, 25},
		{1, 2, 1, 1, 0, 0, 35},
		{2, 2, 1, 2, 1, 1, 60},
		{2, 3, 2, 3, 1, 2, 85},
	}
	for i, c := range cases {
		f := PUCAIFindings{c.p, c.b, c.c, count[c.n], c.noc == 1, c.a}
		if got, err := PUCAIScore(f); err != nil || got != c.want {
			t.Errorf("case %d: got %d, %v; want %d", i+1, got, err, c.want)
		}
	}
}

func TestPUCAIScore_errors(t *testing.T) {
	call := func(f PUCAIFindings) func() error { return func() error { _, err := PUCAIScore(f); return err } }
	runErrCases(t, []errCase{
		{"pain -1", call(PUCAIFindings{AbdominalPain: -1}), ErrOutOfRange},
		{"pain 3", call(PUCAIFindings{AbdominalPain: 3}), ErrOutOfRange},
		{"bleeding 4", call(PUCAIFindings{RectalBleeding: 4}), ErrOutOfRange},
		{"consistency 3", call(PUCAIFindings{StoolConsistency: 3}), ErrOutOfRange},
		{"activity 3", call(PUCAIFindings{ActivityLevel: 3}), ErrOutOfRange},
		{"stools -1", call(PUCAIFindings{StoolsPer24h: -1}), ErrOutOfRange},
	})
}

// Kolho & Turner 2013 (PMC3600299): both sides of each edge. The consumer's bands (parity: ledger pucai,
// source band cases) agree: 5 remission, 15/25 mild, 35/60 moderate, 85 severe.

func TestPUCAI_noAlloc(t *testing.T) {
	f := PUCAIFindings{1, 1, 1, 4, false, 0}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = PUCAIScore(f) }); n != 0 {
		t.Errorf("PUCAIScore allocs = %v, want 0", n)
	}
}

func BenchmarkPUCAIScore(b *testing.B) {
	f := PUCAIFindings{1, 1, 1, 4, false, 0}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = PUCAIScore(f)
	}
}
