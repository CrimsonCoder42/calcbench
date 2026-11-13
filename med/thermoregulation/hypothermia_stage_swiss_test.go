package thermoregulation

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkStage   HypothermiaStage
	sinkInTable bool
)

// Paal 2016 Table 1, "Core temperature (°C)": I "35–32 °C", II "<32–28 °C", III "<28 °C", IV "Variable".
// Both sides of every edge; above 35 °C the table has no row.
func TestHypothermiaStageSwiss_edges(t *testing.T) {
	tests := []struct {
		c       float64
		want    HypothermiaStage
		inTable bool
	}{
		{37, "", false},
		{35.001, "", false},
		{35, HypothermiaI, true},
		{33.5, HypothermiaI, true},
		{32, HypothermiaI, true},
		{31.999, HypothermiaII, true},
		{28, HypothermiaII, true},
		{27.999, HypothermiaIII, true},
		{13.7, HypothermiaIII, true}, // Table 1 footnote b's lowest temperature with successful rewarming
		{0, HypothermiaIII, true},
		{-273.15, HypothermiaIII, true},
	}
	for _, tc := range tests {
		got, in, err := HypothermiaStageSwiss(tc.c)
		if err != nil || got != tc.want || in != tc.inTable {
			t.Errorf("HypothermiaStageSwiss(%v) = (%q, %v, %v), want (%q, %v)", tc.c, got, in, err, tc.want, tc.inTable)
		}
	}
}

// The °F edges are the exact images of 35, 32 and 28 °C: 95, 89.6, 82.4 °F.
func TestHypothermiaStageSwissF_edges(t *testing.T) {
	tests := []struct {
		f       float64
		want    HypothermiaStage
		inTable bool
	}{
		{98.6, "", false},
		{95.001, "", false},
		{95, HypothermiaI, true},
		{89.6, HypothermiaI, true},
		{89.599, HypothermiaII, true},
		{82.4, HypothermiaII, true},
		{82.399, HypothermiaIII, true},
		{-459.67, HypothermiaIII, true},
	}
	for _, tc := range tests {
		got, in, err := HypothermiaStageSwissF(tc.f)
		if err != nil || got != tc.want || in != tc.inTable {
			t.Errorf("HypothermiaStageSwissF(%v) = (%q, %v, %v), want (%q, %v)", tc.f, got, in, err, tc.want, tc.inTable)
		}
	}
}

// The twins agree on a grid of whole-tenth °C values away from the edges (°F = °C × 9/5 + 32).
func TestHypothermiaStageSwiss_twinsAgree(t *testing.T) {
	for tenth := 200; tenth <= 400; tenth++ {
		c := float64(tenth) / 10
		if c == 35 || c == 32 || c == 28 {
			continue // the edges are tested exactly above; their °F image is not exact after conversion
		}
		sc, ic, _ := HypothermiaStageSwiss(c)
		sf, iF, _ := HypothermiaStageSwissF(c*9/5 + 32)
		if sc != sf || ic != iF {
			t.Errorf("%v °C: Celsius (%q, %v), Fahrenheit (%q, %v)", c, sc, ic, sf, iF)
		}
	}
}

// Parity with the reference's °C and °F oral-temperature tables, hypothermia rows: severe < 28 °C (< 82.4 °F),
// moderate 28 to < 32 °C, mild 32 to < 35 °C, and 35 °C itself read as above the hypothermia rows. Table 1
// prints "35–32 °C" for stage I, so 35 °C asserts the table's stage (reference: not hypothermia).
func TestHypothermiaStageSwiss_parity(t *testing.T) {
	tests := []struct {
		c       float64
		want    HypothermiaStage
		inTable bool
	}{
		{27, HypothermiaIII, true}, // parity: ledger temperature_oral_c_adult_11_to_65, source case 1
		{30, HypothermiaII, true},  // parity: ledger temperature_oral_c_adult_11_to_65, source case 2
		{34, HypothermiaI, true},   // parity: ledger temperature_oral_c_adult_11_to_65, source case 3
		{35, HypothermiaI, true},   // parity: ledger temperature_oral_c_adult_11_to_65, source case 4 (reference "Below Normal")
		{36.8, "", false},          // parity: ledger temperature_oral_c_adult_11_to_65, source case 5 (above the table)
	}
	for _, tc := range tests {
		got, in, err := HypothermiaStageSwiss(tc.c)
		if err != nil || got != tc.want || in != tc.inTable {
			t.Errorf("HypothermiaStageSwiss(%v) = (%q, %v, %v), want (%q, %v)", tc.c, got, in, err, tc.want, tc.inTable)
		}
	}
	testsF := []struct {
		f    float64
		want HypothermiaStage
	}{
		{82, HypothermiaIII}, // parity: ledger temperature_oral_adult_11_to_65, source case 1
		{86, HypothermiaII},  // parity: ledger temperature_oral_adult_11_to_65, source case 2
		{94, HypothermiaI},   // parity: ledger temperature_oral_adult_11_to_65, source case 3
	}
	// The reference's °F classifier uses 86 and 90 °F as its severe and moderate edges, not 82.4 and 89.6:
	// the cases assert Table 1's stage, with the classifier's label in the comment.
	testsF = append(testsF, []struct {
		f    float64
		want HypothermiaStage
	}{
		{85, HypothermiaII},  // parity: ledger temp_classifier.go, source case 1 (reference "Severe Hypothermia"; 29.4 °C)
		{89, HypothermiaII},  // parity: ledger temp_classifier.go, source case 2 (reference "Moderate Hypothermia")
		{89.8, HypothermiaI}, // parity: ledger temp_classifier.go, source case 2 (reference "Moderate Hypothermia"; 32.1 °C)
		{94.9, HypothermiaI}, // parity: ledger temp_classifier.go, source case 3 (reference "Mild Hypothermia")
	}...)
	for _, tc := range testsF {
		got, in, err := HypothermiaStageSwissF(tc.f)
		if err != nil || got != tc.want || !in {
			t.Errorf("HypothermiaStageSwissF(%v) = (%q, %v, %v), want %q", tc.f, got, in, err, tc.want)
		}
	}
}

func TestHypothermiaStageSwiss_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, _, err := HypothermiaStageSwiss(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("HypothermiaStageSwiss(%v): err = %v, want ErrNotFinite", v, err)
		}
		if _, _, err := HypothermiaStageSwissF(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("HypothermiaStageSwissF(%v): err = %v, want ErrNotFinite", v, err)
		}
	}
	if _, _, err := HypothermiaStageSwiss(-273.16); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("HypothermiaStageSwiss(−273.16): err = %v, want ErrOutOfRange", err)
	}
	if _, _, err := HypothermiaStageSwissF(-459.68); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("HypothermiaStageSwissF(−459.68): err = %v, want ErrOutOfRange", err)
	}
}

func TestHypothermiaStageSwiss_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkStage, sinkInTable, _ = HypothermiaStageSwiss(30.5) }); n != 0 {
		t.Errorf("HypothermiaStageSwiss allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkStage, sinkInTable, _ = HypothermiaStageSwissF(87) }); n != 0 {
		t.Errorf("HypothermiaStageSwissF allocates %v times, want 0", n)
	}
}

func BenchmarkHypothermiaStageSwiss(b *testing.B) {
	b.ReportAllocs()
	c := 30.5
	for b.Loop() {
		sinkStage, sinkInTable, _ = HypothermiaStageSwiss(c)
	}
}

func BenchmarkHypothermiaStageSwissF(b *testing.B) {
	b.ReportAllocs()
	f := 87.0
	for b.Loop() {
		sinkStage, sinkInTable, _ = HypothermiaStageSwissF(f)
	}
}
