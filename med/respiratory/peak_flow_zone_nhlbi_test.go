package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkPeakFlowZone PeakFlowZone

// NHLBI Asthma Action Plan (NIH Pub. No. 20-HL-5251, Feb 2021), zone peak-flow lines: GREEN ZONE "80 percent
// or more of my best peak flow"; YELLOW ZONE "50 to 79 percent"; RED ZONE "less than … 50 percent". Both
// sides of each edge, with a personal best of 500 L/min (exact quotients) and of 437 L/min (whole L/min either side of the edge).
func TestPeakFlowZoneNHLBI_edges(t *testing.T) {
	tests := []struct {
		pef, best float64
		want      PeakFlowZone
		heading   string
	}{
		{600, 500, PeakFlowGreenZone, "DOING WELL"}, // above personal best
		{500, 500, PeakFlowGreenZone, "DOING WELL"},
		{400, 500, PeakFlowGreenZone, "DOING WELL"}, // exactly 80 percent
		{399.9999, 500, PeakFlowYellowZone, "ASTHMA IS GETTING WORSE"},
		{395, 500, PeakFlowYellowZone, "ASTHMA IS GETTING WORSE"}, // 79 percent, the plan's printed upper yellow
		{250, 500, PeakFlowYellowZone, "ASTHMA IS GETTING WORSE"}, // exactly 50 percent
		{249.9999, 500, PeakFlowRedZone, "MEDICAL ALERT!"},
		{0, 500, PeakFlowRedZone, "MEDICAL ALERT!"},
		{350, 437, PeakFlowGreenZone, "DOING WELL"},               // 80.09 percent
		{349, 437, PeakFlowYellowZone, "ASTHMA IS GETTING WORSE"}, // 79.86 percent
		{219, 437, PeakFlowYellowZone, "ASTHMA IS GETTING WORSE"}, // 50.11 percent
		{218, 437, PeakFlowRedZone, "MEDICAL ALERT!"},             // 49.89 percent
	}
	for _, tc := range tests {
		got, err := PeakFlowZoneNHLBI(tc.pef, tc.best)
		if err != nil || got != tc.want || got.Heading() != tc.heading {
			t.Errorf("PeakFlowZoneNHLBI(%v, %v) = (%q %q, %v), want %q %q", tc.pef, tc.best, got, got.Heading(), err, tc.want, tc.heading)
		}
	}
}

// Parity with the reference's zone classifier: pct < 50 red, pct < 80 yellow, else green, pct = pef/best × 100.
// The reference's range table read 80 as yellow and 81 as the first green; the plan prints "80 percent or more"
// as GREEN ZONE, so 80 asserts the plan's zone (reference: yellow).
func TestPeakFlowZoneNHLBI_parity(t *testing.T) {
	tests := []struct {
		pef, best float64
		want      PeakFlowZone
	}{
		{240, 500, PeakFlowRedZone},      // parity: ledger pef_zone.go, source case 1 (pct < 50)
		{300, 500, PeakFlowYellowZone},   // parity: ledger pef_zone.go, source case 2 (pct < 80)
		{450, 500, PeakFlowGreenZone},    // parity: ledger pef_zone.go, source case 3 (default)
		{400, 500, PeakFlowGreenZone},    // parity: ledger pef_zone_personal_best, source case 2 (reference yellow at 80; plan GREEN)
		{252.5, 505, PeakFlowYellowZone}, // parity: ledger pef_zone_personal_best, source case 2 (50 percent, yellow in both)
	}
	for _, tc := range tests {
		got, err := PeakFlowZoneNHLBI(tc.pef, tc.best)
		if err != nil || got != tc.want {
			t.Errorf("PeakFlowZoneNHLBI(%v, %v) = (%q, %v), want %q", tc.pef, tc.best, got, err, tc.want)
		}
	}
}

func TestPeakFlowZone_Heading_unknown(t *testing.T) {
	if got := PeakFlowZone("BLUE ZONE").Heading(); got != "" {
		t.Errorf(`PeakFlowZone("BLUE ZONE").Heading() = %q, want ""`, got)
	}
}

func TestPeakFlowZoneNHLBI_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := PeakFlowZoneNHLBI(v, 500); !errors.Is(err, ErrNotFinite) {
			t.Errorf("pef %v: err = %v, want ErrNotFinite", v, err)
		}
		if _, err := PeakFlowZoneNHLBI(400, v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("best %v: err = %v, want ErrNotFinite", v, err)
		}
	}
	// pef ≥ 0: −0.001 rejected, 0 accepted (above); best > 0: 0 and −1 rejected, the smallest positive accepted.
	if _, err := PeakFlowZoneNHLBI(-0.001, 500); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("pef −0.001: err = %v, want ErrOutOfRange", err)
	}
	for _, best := range []float64{0, -1} {
		if _, err := PeakFlowZoneNHLBI(400, best); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("best %v: err = %v, want ErrOutOfRange", best, err)
		}
	}
	if got, err := PeakFlowZoneNHLBI(400, math.SmallestNonzeroFloat64); err != nil || got != PeakFlowGreenZone {
		t.Errorf("best smallest positive: (%q, %v), want GREEN ZONE", got, err)
	}
}

func TestPeakFlowZoneNHLBI_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkPeakFlowZone, _ = PeakFlowZoneNHLBI(320, 450) }); n != 0 {
		t.Errorf("PeakFlowZoneNHLBI allocates %v times, want 0", n)
	}
}

func BenchmarkPeakFlowZoneNHLBI(b *testing.B) {
	b.ReportAllocs()
	pef, best := 320.0, 450.0
	for b.Loop() {
		sinkPeakFlowZone, _ = PeakFlowZoneNHLBI(pef, best)
	}
}
