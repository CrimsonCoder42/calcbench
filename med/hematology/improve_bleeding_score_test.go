package hematology

import (
	"errors"
	"math"
	"testing"
)

// IMPROVE bleeding points (Arakaki 2024 Table 1; Chindamo 2021 Table 2). Sums of halves are exact, so the
// comparisons are exact.
func TestIMPROVEBleedingScore_table(t *testing.T) {
	type in struct {
		age                                         float64
		male                                        bool
		renal                                       IMPROVERenal
		cancer, rheum, cvc, icu, liver, plt, bl, ul bool
	}
	tests := []struct {
		name string
		in   in
		want float64
	}{
		{"woman 39, nothing", in{39.99, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, false, false}, 0},
		{"age 40 scores 1.5", in{40, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, false, false}, 1.5},
		{"age 84.99 scores 1.5", in{84.99, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, false, false}, 1.5},
		{"age 85 scores 3.5", in{85, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, false, false}, 3.5},
		{"male", in{30, true, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, false, false}, 1},
		{"GFR 30–59", in{30, false, IMPROVERenalGFR30To59, false, false, false, false, false, false, false, false}, 1},
		{"GFR < 30", in{30, false, IMPROVERenalGFRUnder30, false, false, false, false, false, false, false, false}, 2.5},
		{"cancer", in{30, false, IMPROVERenalGFR60OrMore, true, false, false, false, false, false, false, false}, 2},
		{"rheumatic", in{30, false, IMPROVERenalGFR60OrMore, false, true, false, false, false, false, false, false}, 2},
		{"CVC", in{30, false, IMPROVERenalGFR60OrMore, false, false, true, false, false, false, false, false}, 2},
		{"ICU/CCU", in{30, false, IMPROVERenalGFR60OrMore, false, false, false, true, false, false, false, false}, 2.5},
		{"hepatic failure", in{30, false, IMPROVERenalGFR60OrMore, false, false, false, false, true, false, false, false}, 2.5},
		{"platelets < 50", in{30, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, true, false, false}, 4},
		{"bleeding in 3 months", in{30, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, true, false}, 4},
		{"active ulcer", in{30, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, false, true}, 4.5},
		{"example: man 86, GFR 25, ICU, ulcer", in{86, true, IMPROVERenalGFRUnder30, false, false, false, true, false, false, false, true}, 14},
		{"maximum 30.5", in{90, true, IMPROVERenalGFRUnder30, true, true, true, true, true, true, true, true}, 30.5},
		{"≥7 edge: age 50, ICU, CVC, male", in{50, true, IMPROVERenalGFR60OrMore, false, false, true, true, false, false, false, false}, 7},
		{"just under 7", in{50, false, IMPROVERenalGFR60OrMore, false, false, true, true, false, false, false, false}, 6},
	}
	for _, tc := range tests {
		x := tc.in
		got, err := IMPROVEBleedingScore(x.age, x.male, x.renal, x.cancer, x.rheum, x.cvc, x.icu, x.liver, x.plt, x.bl, x.ul)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger improve_bleeding. The source entry has no _test.go and its JSON no example; the consumer sums
// the same weights from yes/no inputs and a 3-level GFR input, so case 1 (all present, age ≥ 85, GFR < 30) is 30.5.
func TestIMPROVEBleedingScore_parity(t *testing.T) {
	got, err := IMPROVEBleedingScore(85, true, IMPROVERenalGFRUnder30, true, true, true, true, true, true, true, true)
	if err != nil || got != 30.5 {
		t.Errorf("source case 1: got (%v, %v), want 30.5", got, err)
	}
}

func TestIMPROVEBleedingScore_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := IMPROVEBleedingScore(v, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, false, false); !errors.Is(err, ErrNotFinite) {
			t.Errorf("ageYears %v: error = %v, want ErrNotFinite", v, err)
		}
	}
	for _, v := range []float64{-0.01, 130.01} {
		if _, err := IMPROVEBleedingScore(v, false, IMPROVERenalGFR60OrMore, false, false, false, false, false, false, false, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("ageYears %v: error = %v, want ErrOutOfRange", v, err)
		}
	}
	for _, r := range []IMPROVERenal{-1, 3} {
		if _, err := IMPROVEBleedingScore(50, false, r, false, false, false, false, false, false, false, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("renal %d: error = %v, want ErrOutOfRange", r, err)
		}
	}
}

func TestIMPROVEBleedingScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkFloat, _ = IMPROVEBleedingScore(86, true, IMPROVERenalGFRUnder30, false, false, false, true, false, false, false, true)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkIMPROVEBleedingScore(b *testing.B) {
	b.ReportAllocs()
	age, r := 86.0, IMPROVERenalGFRUnder30
	for b.Loop() {
		sinkFloat, _ = IMPROVEBleedingScore(age, true, r, false, false, false, true, false, false, false, true)
	}
}
