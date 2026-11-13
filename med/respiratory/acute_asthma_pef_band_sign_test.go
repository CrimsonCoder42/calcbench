package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkAcuteAsthmaPEFBand AcuteAsthmaPEFBand

// SIGN 158 (Nov 2024, corrections May 2025) Table 15: "PEF >50–75% best or predicted" (Moderate acute
// asthma), "PEF 33–50% best or predicted" (Acute severe asthma), "PEF <33% best or predicted"
// (Life-threatening asthma). Both sides of 33, 50 and 75; the cell text and its row verbatim.
func TestAcuteAsthmaPEFBandSIGN_table15(t *testing.T) {
	tests := []struct {
		p     float64
		want  AcuteAsthmaPEFBand
		level string
	}{
		{0, AcuteAsthmaPEFLifeThreatening, "Life-threatening asthma"},
		{32.999, AcuteAsthmaPEFLifeThreatening, "Life-threatening asthma"},
		{33, AcuteAsthmaPEFAcuteSevere, "Acute severe asthma"},
		{50, AcuteAsthmaPEFAcuteSevere, "Acute severe asthma"},
		{50.001, AcuteAsthmaPEFModerate, "Moderate acute asthma"},
		{75, AcuteAsthmaPEFModerate, "Moderate acute asthma"},
		{75.001, "", ""},
		{120, "", ""},
	}
	for _, tc := range tests {
		got, err := AcuteAsthmaPEFBandSIGN(tc.p, 40)
		if err != nil || got != tc.want || got.Level() != tc.level {
			t.Errorf("AcuteAsthmaPEFBandSIGN(%v, 40) = (%q %q, %v), want %q %q", tc.p, got, got.Level(), err, tc.want, tc.level)
		}
	}
}

// Parity with the reference's bands (< 33 "Life-Threatening Exacerbation", 33–< 50 "Acute Severe
// Exacerbation", 50–< 75 "Moderate Exacerbation", ≥ 75 "Mild Exacerbation"), applied alike to % personal best
// and % predicted. Table 15 puts 50 in "33–50%", 75 in ">50–75%", and has no cell above 75 (no "mild" level):
// those cases assert the table.
func TestAcuteAsthmaPEFBandSIGN_parity(t *testing.T) {
	tests := []struct {
		p    float64
		want AcuteAsthmaPEFBand
	}{
		{32.9, AcuteAsthmaPEFLifeThreatening}, // parity: ledger pef_exacerbation_percent_personal_best, source case 1
		{33, AcuteAsthmaPEFAcuteSevere},       // parity: ledger pef_exacerbation_percent_personal_best, source case 2
		{49.9, AcuteAsthmaPEFAcuteSevere},     // parity: ledger pef_exacerbation_percent_personal_best, source case 3
		{50, AcuteAsthmaPEFAcuteSevere},       // parity: ledger pef_exacerbation_percent_personal_best, source case 4 (reference Moderate; table "33–50%")
		{74.9, AcuteAsthmaPEFModerate},        // parity: ledger pef_exacerbation_percent_personal_best, source case 5
		{75, AcuteAsthmaPEFModerate},          // parity: ledger pef_exacerbation_percent_personal_best, source case 6 (reference Mild; table ">50–75%")
		{80, ""},                              // parity: ledger pef_exacerbation_percent_personal_best, source case 7 (reference Mild; no table cell)
		{40, AcuteAsthmaPEFAcuteSevere},       // parity: ledger pef_exacerbation_percent_predicted, source case 1 (same bands on % predicted)
		{60, AcuteAsthmaPEFModerate},          // parity: ledger pef_exacerbation_percent_predicted, source case 2
	}
	for _, tc := range tests {
		got, err := AcuteAsthmaPEFBandSIGN(tc.p, 30)
		if err != nil || got != tc.want {
			t.Errorf("AcuteAsthmaPEFBandSIGN(%v, 30) = (%q, %v), want %q", tc.p, got, err, tc.want)
		}
	}
}

func TestAcuteAsthmaPEFBand_Level_unknown(t *testing.T) {
	if got := AcuteAsthmaPEFBand("PEF >75%").Level(); got != "" {
		t.Errorf(`Level() = %q, want ""`, got)
	}
}

func TestAcuteAsthmaPEFBandSIGN_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := AcuteAsthmaPEFBandSIGN(v, 30); !errors.Is(err, ErrNotFinite) {
			t.Errorf("pef %v: err = %v, want ErrNotFinite", v, err)
		}
		if _, err := AcuteAsthmaPEFBandSIGN(40, v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("age %v: err = %v, want ErrNotFinite", v, err)
		}
	}
	if _, err := AcuteAsthmaPEFBandSIGN(-0.001, 30); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("pef < 0: err = %v, want ErrOutOfRange", err)
	}
	if _, err := AcuteAsthmaPEFBandSIGN(40, 17.99); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("age 17.99: err = %v, want ErrOutOfRange", err)
	}
	if _, err := AcuteAsthmaPEFBandSIGN(40, 18); err != nil {
		t.Errorf("age 18: err = %v, want nil", err)
	}
}

func TestAcuteAsthmaPEFBandSIGN_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkAcuteAsthmaPEFBand, _ = AcuteAsthmaPEFBandSIGN(45, 30) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkAcuteAsthmaPEFBandSIGN(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkAcuteAsthmaPEFBand, _ = AcuteAsthmaPEFBandSIGN(45, 30)
	}
}
