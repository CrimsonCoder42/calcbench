package neuro

import (
	"errors"
	"testing"
)

// Kwon 2026 / Ruiz-Marrara 2026: domain maxima 5,3,6,3,2,5,6 (sum 30), +1 for ≤ 12 years, range 0–30.
func TestMoCATotal_example(t *testing.T) {
	maxD := MoCADomains{5, 3, 6, 3, 2, 5, 6}
	if got, err := MoCATotal(maxD, false); err != nil || got != 30 {
		t.Errorf("max = (%d, %v), want 30", got, err)
	}
	if got, err := MoCATotal(maxD, true); err != nil || got != 30 {
		t.Errorf("max + education = (%d, %v), want 30 (capped)", got, err)
	}
	d25 := MoCADomains{4, 3, 5, 2, 2, 3, 6}
	if got, _ := MoCATotal(d25, true); got != 26 {
		t.Errorf("25 + education = %d, want 26", got)
	}
	if got, _ := MoCATotal(d25, false); got != 25 {
		t.Errorf("25 = %d, want 25", got)
	}
	if got, _ := MoCATotal(MoCADomains{}, true); got != 1 {
		t.Errorf("0 + education = %d, want 1", got)
	}
}

func TestMoCATotal_errors(t *testing.T) {
	bad := []MoCADomains{
		{VisuospatialExecutive: 6}, {Naming: 4}, {Attention: 7}, {Language: 4}, {Abstraction: 3},
		{DelayedRecall: 6}, {Orientation: 7}, {Naming: -1},
	}
	for _, b := range bad {
		if _, err := MoCATotal(b, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("MoCATotal(%+v): %v, want ErrOutOfRange", b, err)
		}
	}
}

// parity: ledger moca_score, source cases 1–2 (consumer: min(30, sum + education bonus)).
func TestMoCATotal_parity(t *testing.T) {
	if got, _ := MoCATotal(MoCADomains{5, 3, 6, 3, 2, 5, 5}, true); got != 30 { // parity: ledger moca_score, source case 1
		t.Errorf("case 1 = %d, want 30", got)
	}
	if got, _ := MoCATotal(MoCADomains{3, 2, 4, 2, 1, 2, 5}, false); got != 19 { // parity: ledger moca_score, source case 2
		t.Errorf("case 2 = %d, want 19", got)
	}
}

func TestMoCATotal_noAlloc(t *testing.T) {
	d := MoCADomains{4, 3, 5, 2, 2, 3, 6}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MoCATotal(d, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkMoCATotal(b *testing.B) {
	b.ReportAllocs()
	d := MoCADomains{4, 3, 5, 2, 2, 3, 6}
	for b.Loop() {
		sinkInt, _ = MoCATotal(d, true)
	}
}
