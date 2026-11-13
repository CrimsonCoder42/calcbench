package gastro

import "testing"

// HAND-DERIVED from the weights in Prabhu 2024 (PMC11087879, Table 1) and Jose 2021 (PMC8289675):
// tenderness 2, leukocytosis 2, the other six 1 each, maximum 10.
func TestAlvaradoScore_weights(t *testing.T) {
	cases := []struct {
		name string
		f    AlvaradoFindings
		want int
	}{
		{"none", AlvaradoFindings{}, 0},
		{"migration", AlvaradoFindings{MigrationOfPain: true}, 1},
		{"anorexia", AlvaradoFindings{Anorexia: true}, 1},
		{"nausea", AlvaradoFindings{NauseaOrVomiting: true}, 1},
		{"tenderness", AlvaradoFindings{TendernessRLQ: true}, 2},
		{"rebound", AlvaradoFindings{ReboundPain: true}, 1},
		{"temperature", AlvaradoFindings{ElevatedTemperature: true}, 1},
		{"leukocytosis", AlvaradoFindings{Leukocytosis: true}, 2},
		{"shift", AlvaradoFindings{ShiftToLeft: true}, 1},
		{"doc example", AlvaradoFindings{TendernessRLQ: true, Leukocytosis: true, MigrationOfPain: true}, 5},
		{"all", AlvaradoFindings{true, true, true, true, true, true, true, true}, 10},
	}
	for _, c := range cases {
		if got := AlvaradoScore(c.f); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}

// parity: ledger alvarado_appendicitis, source case 1: the consumer sums migration 1, anorexia 1, nausea 1,
// RLQ tenderness 2, rebound 1, temperature 1, leukocytosis 2, left shift 1; all but rebound and shift → 8.
func TestAlvaradoScore_parity(t *testing.T) {
	f := AlvaradoFindings{MigrationOfPain: true, Anorexia: true, NauseaOrVomiting: true, TendernessRLQ: true, ElevatedTemperature: true, Leukocytosis: true}
	if got := AlvaradoScore(f); got != 8 {
		t.Errorf("got %d, want 8", got)
	}
}

func TestAlvaradoScore_noAlloc(t *testing.T) {
	f := AlvaradoFindings{TendernessRLQ: true, Leukocytosis: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = AlvaradoScore(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkAlvaradoScore(b *testing.B) {
	f := AlvaradoFindings{TendernessRLQ: true, Leukocytosis: true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt = AlvaradoScore(f)
	}
}
