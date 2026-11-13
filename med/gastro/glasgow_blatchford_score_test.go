package gastro

import (
	"math"
	"testing"
)

var gbsBase = GlasgowBlatchfordFindings{UreaMmolL: 5, HaemoglobinGL: 140, SystolicMmHg: 120, PulseBpm: 80}

// Banister 2018 Table 1 / Orpen-Palmer 2023 Table 1 (PMC6135483, PMC10253886): every band, each side of each
// edge. HAND-DERIVED totals.
func TestGlasgowBlatchfordScore_tableBands(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*GlasgowBlatchfordFindings)
		want int
	}{
		{"all zero", func(f *GlasgowBlatchfordFindings) {}, 0},
		{"urea 6.49", func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = 6.49 }, 0},
		{"urea 6.5", func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = 6.5 }, 2},
		{"urea 7.99", func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = 7.99 }, 2},
		{"urea 8.0", func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = 8 }, 3},
		{"urea 9.99", func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = 9.99 }, 3},
		{"urea 10.0", func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = 10 }, 4},
		{"urea 24.99", func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = 24.99 }, 4},
		{"urea 25.0", func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = 25 }, 6},
		{"man Hb 130", func(f *GlasgowBlatchfordFindings) { f.Male = true; f.HaemoglobinGL = 130 }, 0},
		{"man Hb 129", func(f *GlasgowBlatchfordFindings) { f.Male = true; f.HaemoglobinGL = 129 }, 1},
		{"man Hb 120", func(f *GlasgowBlatchfordFindings) { f.Male = true; f.HaemoglobinGL = 120 }, 1},
		{"man Hb 119", func(f *GlasgowBlatchfordFindings) { f.Male = true; f.HaemoglobinGL = 119 }, 3},
		{"man Hb 100", func(f *GlasgowBlatchfordFindings) { f.Male = true; f.HaemoglobinGL = 100 }, 3},
		{"man Hb 99", func(f *GlasgowBlatchfordFindings) { f.Male = true; f.HaemoglobinGL = 99 }, 6},
		{"woman Hb 125", func(f *GlasgowBlatchfordFindings) { f.HaemoglobinGL = 125 }, 0},
		{"woman Hb 120", func(f *GlasgowBlatchfordFindings) { f.HaemoglobinGL = 120 }, 0},
		{"woman Hb 119", func(f *GlasgowBlatchfordFindings) { f.HaemoglobinGL = 119 }, 1},
		{"woman Hb 100", func(f *GlasgowBlatchfordFindings) { f.HaemoglobinGL = 100 }, 1},
		{"woman Hb 99", func(f *GlasgowBlatchfordFindings) { f.HaemoglobinGL = 99 }, 6},
		{"SBP 110", func(f *GlasgowBlatchfordFindings) { f.SystolicMmHg = 110 }, 0},
		{"SBP 109", func(f *GlasgowBlatchfordFindings) { f.SystolicMmHg = 109 }, 1},
		{"SBP 100", func(f *GlasgowBlatchfordFindings) { f.SystolicMmHg = 100 }, 1},
		{"SBP 99", func(f *GlasgowBlatchfordFindings) { f.SystolicMmHg = 99 }, 2},
		{"SBP 90", func(f *GlasgowBlatchfordFindings) { f.SystolicMmHg = 90 }, 2},
		{"SBP 89", func(f *GlasgowBlatchfordFindings) { f.SystolicMmHg = 89 }, 3},
		{"pulse 99", func(f *GlasgowBlatchfordFindings) { f.PulseBpm = 99 }, 0},
		{"pulse 100", func(f *GlasgowBlatchfordFindings) { f.PulseBpm = 100 }, 1},
		{"melaena", func(f *GlasgowBlatchfordFindings) { f.Melaena = true }, 1},
		{"syncope", func(f *GlasgowBlatchfordFindings) { f.Syncope = true }, 2},
		{"hepatic", func(f *GlasgowBlatchfordFindings) { f.HepaticDisease = true }, 2},
		{"cardiac failure", func(f *GlasgowBlatchfordFindings) { f.CardiacFailure = true }, 2},
		{"maximum", func(f *GlasgowBlatchfordFindings) {
			*f = GlasgowBlatchfordFindings{UreaMmolL: 30, HaemoglobinGL: 80, SystolicMmHg: 80, PulseBpm: 120, Male: true, Melaena: true, Syncope: true, HepaticDisease: true, CardiacFailure: true}
		}, 23},
		{"doc example", func(f *GlasgowBlatchfordFindings) {
			*f = GlasgowBlatchfordFindings{UreaMmolL: 9, HaemoglobinGL: 115, SystolicMmHg: 95, PulseBpm: 105, Male: true, Melaena: true}
		}, 10},
	}
	for _, c := range cases {
		f := gbsBase
		c.mod(&f)
		got, err := GlasgowBlatchfordScore(f)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger glasgow_blatchford, source case 1: the consumer takes BUN mg/dL with bands 18.2 / 22.4 / 28 /
// > 70 (= urea mmol/L × 2.8) and Hb g/dL; a woman with BUN 30 mg/dL (urea 30/2.8 ≈ 10.7 mmol/L), Hb 11 g/dL
// (110 g/L), SBP 105, pulse ≥ 100, syncope → 4 + 1 + 1 + 1 + 2 = 9. The consumer's "> 70" edge differs from the
// sources' urea 25 mmol/L (= 70 mg/dL) start of the 6-point band: at exactly 25 mmol/L this function gives 6.
func TestGlasgowBlatchfordScore_parity(t *testing.T) {
	f := GlasgowBlatchfordFindings{UreaMmolL: 30.0 / 2.8, HaemoglobinGL: 110, SystolicMmHg: 105, PulseBpm: 100, Syncope: true}
	if got, err := GlasgowBlatchfordScore(f); err != nil || got != 9 {
		t.Errorf("got %d, %v; want 9", got, err)
	}
}

func TestGlasgowBlatchfordScore_errors(t *testing.T) {
	bad := func(mod func(*GlasgowBlatchfordFindings)) func() error {
		return func() error { f := gbsBase; mod(&f); _, err := GlasgowBlatchfordScore(f); return err }
	}
	runErrCases(t, []errCase{
		{"urea NaN", bad(func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = math.NaN() }), ErrNotFinite},
		{"Hb +Inf", bad(func(f *GlasgowBlatchfordFindings) { f.HaemoglobinGL = math.Inf(1) }), ErrNotFinite},
		{"SBP -Inf", bad(func(f *GlasgowBlatchfordFindings) { f.SystolicMmHg = math.Inf(-1) }), ErrNotFinite},
		{"pulse NaN", bad(func(f *GlasgowBlatchfordFindings) { f.PulseBpm = math.NaN() }), ErrNotFinite},
		{"urea negative", bad(func(f *GlasgowBlatchfordFindings) { f.UreaMmolL = -1 }), ErrOutOfRange},
		{"Hb negative", bad(func(f *GlasgowBlatchfordFindings) { f.HaemoglobinGL = -0.1 }), ErrOutOfRange},
		{"SBP negative", bad(func(f *GlasgowBlatchfordFindings) { f.SystolicMmHg = -5 }), ErrOutOfRange},
		{"pulse negative", bad(func(f *GlasgowBlatchfordFindings) { f.PulseBpm = -1 }), ErrOutOfRange},
	})
}

func TestGlasgowBlatchfordScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GlasgowBlatchfordScore(gbsBase) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkGlasgowBlatchfordScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = GlasgowBlatchfordScore(gbsBase)
	}
}
