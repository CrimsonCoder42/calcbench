package pediatrics

import (
	"errors"
	"testing"
)

// Shafi 2020 (PMC7719480) chart: rows 0–3 and "Score 2 extra" items. HAND-DERIVED totals.
func TestPEWSScore_chart(t *testing.T) {
	cases := []struct {
		name string
		f    PEWSFindings
		want int
	}{
		{"zero", PEWSFindings{}, 0},
		{"behaviour 3", PEWSFindings{Behaviour: 3}, 3},
		{"cardiovascular 3", PEWSFindings{Cardiovascular: 3}, 3},
		{"respiratory 3", PEWSFindings{Respiratory: 3}, 3},
		{"nebulisers", PEWSFindings{QuarterHourlyNebulisers: true}, 2},
		{"vomiting", PEWSFindings{PostoperativeVomiting: true}, 2},
		{"doc example", PEWSFindings{Behaviour: 2, Cardiovascular: 1, QuarterHourlyNebulisers: true}, 5},
		{"maximum 13", PEWSFindings{3, 3, 3, true, true}, 13},
	}
	for _, c := range cases {
		got, err := PEWSScore(c.f)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger pews, source case 1: the consumer adds 1 (not 2) for nebulisers and for vomiting; behaviour
// 1, cardiovascular 2, respiratory 1, nebulisers → consumer 5; the sources give 1 + 2 + 1 + 2 = 6 (asserted).
func TestPEWSScore_parity(t *testing.T) {
	if got, err := PEWSScore(PEWSFindings{1, 2, 1, true, false}); err != nil || got != 6 {
		t.Errorf("got %d, %v; want 6", got, err)
	}
}

func TestPEWSScore_errors(t *testing.T) {
	for _, f := range []PEWSFindings{{Behaviour: -1}, {Behaviour: 4}, {Cardiovascular: 4}, {Respiratory: -1}} {
		if _, err := PEWSScore(f); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%+v: err = %v, want ErrOutOfRange", f, err)
		}
	}
}

func TestPEWSScore_noAlloc(t *testing.T) {
	f := PEWSFindings{2, 1, 0, true, false}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = PEWSScore(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPEWSScore(b *testing.B) {
	f := PEWSFindings{2, 1, 0, true, false}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = PEWSScore(f)
	}
}
