package respiratory

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

var (
	sinkPEFFloat float64
	sinkPEFBool  bool
)

// Tolerance 1e-9 absolute on percents: the formula is three exact-input float operations, so results agree
// with the hand arithmetic to a few ulps; the hand values below are carried to 13+ digits.
const pefTol = 1e-9

// GINA 2023 Box 1-2 footnote *: "(day's highest minus day's lowest) divided by (mean of day's highest and
// lowest)"; text: "( [Day's highest – day's lowest] / mean of day's highest and lowest) x 100". HAND-DERIVED.
func TestPEFAmplitudePercentMean_gina(t *testing.T) {
	tests := []struct {
		hi, lo, want float64
	}{
		{440, 400, 40.0 / 420 * 100}, // 9.523809523809524
		{460, 380, 80.0 / 420 * 100}, // 19.047619047619047
		{500, 500, 0},                // no variation
		{300, 0, 200},                // lowest 0: the arithmetic maximum
		{1, 0.5, 0.5 / 0.75 * 100},   // 66.666…
		{math.MaxFloat64, 0, 200},    // no overflow at the top of float64
		{math.MaxFloat64, math.MaxFloat64 / 2, (math.MaxFloat64 / 2) / (math.MaxFloat64 * 0.75) * 100},
	}
	for _, tc := range tests {
		got, err := PEFAmplitudePercentMean(tc.hi, tc.lo)
		if err != nil || math.Abs(got-tc.want) > pefTol {
			t.Errorf("PEFAmplitudePercentMean(%v, %v) = (%v, %v), want %v", tc.hi, tc.lo, got, err, tc.want)
		}
	}
}

func TestPEFAmplitudePercentMean_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := PEFAmplitudePercentMean(v, 400); !errors.Is(err, ErrNotFinite) {
			t.Errorf("highest %v: err = %v, want ErrNotFinite", v, err)
		}
		if _, err := PEFAmplitudePercentMean(400, v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("lowest %v: err = %v, want ErrNotFinite", v, err)
		}
	}
	for _, tc := range []struct{ hi, lo float64 }{
		{400, -0.001}, // lowest < 0
		{0, 0},        // highest not > 0
		{-1, 0},       // highest < 0
		{399, 400},    // highest < lowest
	} {
		if _, err := PEFAmplitudePercentMean(tc.hi, tc.lo); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("PEFAmplitudePercentMean(%v, %v): err = %v, want ErrOutOfRange", tc.hi, tc.lo, err)
		}
	}
	// Edges accepted: lowest 0, highest = lowest.
	if _, err := PEFAmplitudePercentMean(400, 0); err != nil {
		t.Errorf("lowest 0: err = %v", err)
	}
	if _, err := PEFAmplitudePercentMean(400, 400); err != nil {
		t.Errorf("highest = lowest: err = %v", err)
	}
}

// GINA 2023: each day's value "averaged over one week" (Box 1-2) / "over 1–2 weeks" (text). HAND-DERIVED.
func TestPEFDiurnalVariabilityPercent_gina(t *testing.T) {
	hi := []float64{440, 460}
	lo := []float64{400, 380}
	want := (40.0/420*100 + 80.0/420*100) / 2 // 14.285714285714286
	got, err := PEFDiurnalVariabilityPercent(hi, lo)
	if err != nil || math.Abs(got-want) > pefTol {
		t.Errorf("PEFDiurnalVariabilityPercent = (%v, %v), want %v", got, err, want)
	}
	// One day equals MED-RESP-050.
	got, err = PEFDiurnalVariabilityPercent([]float64{440}, []float64{400})
	if err != nil || math.Abs(got-40.0/420*100) > pefTol {
		t.Errorf("one day = (%v, %v), want %v", got, err, 40.0/420*100)
	}
	// Fourteen identical days average to that day's value (compensated sum; no drift).
	var h14, l14 [14]float64
	for i := range h14 {
		h14[i], l14[i] = 430, 390
	}
	got, err = PEFDiurnalVariabilityPercent(h14[:], l14[:])
	if err != nil || math.Abs(got-40.0/410*100) > pefTol {
		t.Errorf("14 equal days = (%v, %v), want %v", got, err, 40.0/410*100)
	}
	// The inputs are not modified.
	if hi[0] != 440 || hi[1] != 460 || lo[0] != 400 || lo[1] != 380 {
		t.Errorf("inputs modified: %v %v", hi, lo)
	}
}

// Parity with the reference's daily-amplitude routine: readings grouped by calendar day, each day with ≥ 2
// readings gives (highest − lowest) / mean(highest, lowest) × 100, averaged; it rounds each day to 0.1 and the
// average to 0.1. calcbench does not round (the primary does not), so the unrounded values are asserted.
func TestPEFDiurnalVariabilityPercent_parity(t *testing.T) {
	// parity: ledger pef_diurnal.go, source case 1 — day 1 readings 400, 440, 420 (highest 440, lowest 400);
	// day 2 readings 380, 460; reference days 9.5 and 19.0, average 14.3 (rounded).
	got, err := PEFDiurnalVariabilityPercent([]float64{440, 460}, []float64{400, 380})
	if err != nil || math.Abs(got-14.285714285714286) > pefTol {
		t.Errorf("source case 1 = (%v, %v), want 14.285714285714286 (reference 14.3 after rounding)", got, err)
	}
	// parity: ledger pef_diurnal.go, source case 2 — one day 300 and 300: variability 0.
	got, err = PEFDiurnalVariabilityPercent([]float64{300}, []float64{300})
	if err != nil || got != 0 {
		t.Errorf("source case 2 = (%v, %v), want 0", got, err)
	}
	// parity: ledger X2 diurnal.go:114 — the candidate used max(PM) − min(AM); GINA uses the day's highest
	// and lowest of all its readings. AM 400, 410; PM 380, 430: GINA day = (430 − 380)/405 × 100, the
	// candidate's (430 − 400)/415 × 100 = 7.23 (not the target).
	got, err = PEFAmplitudePercentMean(430, 380)
	if err != nil || math.Abs(got-50.0/405*100) > pefTol {
		t.Errorf("X2 case = (%v, %v), want %v (candidate 7.23)", got, err, 50.0/405*100)
	}
}

func TestPEFDiurnalVariabilityPercent_errors(t *testing.T) {
	if _, err := PEFDiurnalVariabilityPercent(nil, nil); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("no day: err = %v, want ErrOutOfRange", err)
	}
	if _, err := PEFDiurnalVariabilityPercent([]float64{440, 460}, []float64{400}); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("length mismatch: err = %v, want ErrOutOfRange", err)
	}
	if _, err := PEFDiurnalVariabilityPercent([]float64{440, math.NaN()}, []float64{400, 380}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN highest: err = %v, want ErrNotFinite", err)
	}
	if _, err := PEFDiurnalVariabilityPercent([]float64{440, 460}, []float64{400, math.Inf(-1)}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("-Inf lowest: err = %v, want ErrNotFinite", err)
	}
	if _, err := PEFDiurnalVariabilityPercent([]float64{440, 370}, []float64{400, 380}); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("highest < lowest on day 1: err = %v, want ErrOutOfRange", err)
	}
}

// GINA 2023 Box 1-2: "Adults: average daily diurnal PEF variability >10%"; "Children: … >13%"; the box
// covers "adults, adolescents, and children 6–11 years". Both sides of each threshold and of the age edges.
func TestPEFDiurnalVariabilityExcessiveGINA_edges(t *testing.T) {
	tests := []struct {
		v, age float64
		want   bool
	}{
		{10, 40, false}, {10.000001, 40, true}, {9.99, 40, false}, {0, 40, false}, {200, 40, true},
		{13, 8, false}, {13.000001, 8, true}, {12, 8, false},
		{12, 6, false}, {13.5, 6, true}, // age 6: the box's youngest group
		{12, 11.99, false}, {12, 12, true}, // 12 years: adolescents take the adults' > 10 %
		{10.5, 17, true}, {10.5, 18, true},
	}
	for _, tc := range tests {
		got, err := PEFDiurnalVariabilityExcessiveGINA(tc.v, tc.age)
		if err != nil || got != tc.want {
			t.Errorf("PEFDiurnalVariabilityExcessiveGINA(%v, %v) = (%v, %v), want %v", tc.v, tc.age, got, err, tc.want)
		}
	}
}

// Parity with the reference's range bands (adult: ≤ 10 lowest band, > 10 "Diagnostic-Level
// Variability (GINA)", > 20 "Elevated Variability — Poorly Controlled"; child 5–17: 13 in place of 10). The
// > 20 band has no GINA source: calcbench returns only the Box 1-2 criterion. The reference's child split at
// 18 and from age 5 disagrees with the box (children 6–11); those cases assert the box.
func TestPEFDiurnalVariabilityExcessiveGINA_parity(t *testing.T) {
	tests := []struct {
		v, age float64
		want   bool
	}{
		{10, 30, false},  // parity: ledger pef_variability_adult, source case 1 (10 → lowest band)
		{10.1, 30, true}, // parity: ledger pef_variability_adult, source case 2 (> 10 → diagnostic-level)
		{20, 30, true},   // parity: ledger pef_variability_adult, source case 3 (20 → diagnostic-level)
		{25, 30, true},   // parity: ledger pef_variability_adult, source case 4 (> 20 band: no primary; criterion met)
		{13, 8, false},   // parity: ledger pef_variability_child, source case 1 (13 → lowest band)
		{13.1, 8, true},  // parity: ledger pef_variability_child, source case 2 (> 13 → diagnostic-level)
		{12, 15, true},   // parity: ledger pef_variability_child, source case 3 (reference lowest band at 15 y; box: adolescents > 10 %)
	}
	for _, tc := range tests {
		got, err := PEFDiurnalVariabilityExcessiveGINA(tc.v, tc.age)
		if err != nil || got != tc.want {
			t.Errorf("PEFDiurnalVariabilityExcessiveGINA(%v, %v) = (%v, %v), want %v", tc.v, tc.age, got, err, tc.want)
		}
	}
	// parity: ledger pef_variability_child, source case 4 — age 5 is in the reference's child band; the box
	// starts at 6, so it is refused.
	if _, err := PEFDiurnalVariabilityExcessiveGINA(14, 5); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("age 5: err = %v, want ErrOutOfRange", err)
	}
}

func TestPEFDiurnalVariabilityExcessiveGINA_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := PEFDiurnalVariabilityExcessiveGINA(v, 30); !errors.Is(err, ErrNotFinite) {
			t.Errorf("variability %v: err = %v, want ErrNotFinite", v, err)
		}
		if _, err := PEFDiurnalVariabilityExcessiveGINA(5, v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("age %v: err = %v, want ErrNotFinite", v, err)
		}
	}
	for _, tc := range []struct{ v, age float64 }{{-0.001, 30}, {200.001, 30}, {5, 5.99}} {
		if _, err := PEFDiurnalVariabilityExcessiveGINA(tc.v, tc.age); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("(%v, %v): err = %v, want ErrOutOfRange", tc.v, tc.age, err)
		}
	}
}

func TestPEFDiurnalVariability_noAlloc(t *testing.T) {
	hi := []float64{440, 460, 450, 470, 430, 445, 455}
	lo := []float64{400, 380, 410, 405, 395, 400, 390}
	if n := testing.AllocsPerRun(100, func() { sinkPEFFloat, _ = PEFAmplitudePercentMean(440, 400) }); n != 0 {
		t.Errorf("PEFAmplitudePercentMean allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkPEFFloat, _ = PEFDiurnalVariabilityPercent(hi, lo) }); n != 0 {
		t.Errorf("PEFDiurnalVariabilityPercent allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkPEFBool, _ = PEFDiurnalVariabilityExcessiveGINA(12, 30) }); n != 0 {
		t.Errorf("PEFDiurnalVariabilityExcessiveGINA allocs = %v, want 0", n)
	}
}

func BenchmarkPEFAmplitudePercentMean(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkPEFFloat, _ = PEFAmplitudePercentMean(440, 400)
	}
}

func BenchmarkPEFDiurnalVariabilityPercent(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		hi := make([]float64, n)
		lo := make([]float64, n)
		for i := range hi {
			hi[i], lo[i] = 400+float64(i%60), 380+float64(i%20)
		}
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sinkPEFFloat, _ = PEFDiurnalVariabilityPercent(hi, lo)
			}
		})
	}
}

func BenchmarkPEFDiurnalVariabilityExcessiveGINA(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkPEFBool, _ = PEFDiurnalVariabilityExcessiveGINA(12, 30)
	}
}
