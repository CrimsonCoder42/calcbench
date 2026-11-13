package psychiatry

import (
	"errors"
	"testing"
)

// Webster 2005 abstract bands, each edge on both sides.
func TestOpioidRiskToolCategory_example(t *testing.T) {
	tests := []struct {
		total int
		want  ORTCategory
	}{
		{0, ORTLowRisk}, {3, ORTLowRisk}, {4, ORTModerateRisk}, {7, ORTModerateRisk}, {8, ORTHighRisk}, {26, ORTHighRisk},
	}
	for _, tc := range tests {
		if got, err := OpioidRiskToolCategory(tc.total); err != nil || got != tc.want {
			t.Errorf("OpioidRiskToolCategory(%d) = (%q, %v), want %q", tc.total, got, err, tc.want)
		}
	}
	if _, err := OpioidRiskToolCategory(-1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("OpioidRiskToolCategory(-1) error = %v, want ErrOutOfRange", err)
	}
}

// parity: ledger opioid_risk_tool, source cases 1–7: the source's totals and bands (lo/mid/hi). The totals
// themselves come from item weights not read in the primary, so only the band of each total is checked.
func TestOpioidRiskToolCategory_parity(t *testing.T) {
	tests := []struct {
		total int
		want  ORTCategory
	}{
		{0, ORTLowRisk},      // source case 1 "zero"
		{4, ORTModerateRisk}, // source case 2 "female four"
		{4, ORTModerateRisk}, // source case 3 "male four"
		{11, ORTHighRisk},    // source case 4 "female eleven"
		{27, ORTHighRisk},    // source case 5 "male high"
		{26, ORTHighRisk},    // source case 6 "female all"
		{3, ORTLowRisk},      // source case 7 "male three"
	}
	for i, tc := range tests {
		if got, err := OpioidRiskToolCategory(tc.total); err != nil || got != tc.want {
			t.Errorf("source case %d: OpioidRiskToolCategory(%d) = (%q, %v), want %q", i+1, tc.total, got, err, tc.want)
		}
	}
}

func TestOpioidRiskToolCategory_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { c, _ := OpioidRiskToolCategory(5); sinkString = string(c) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkOpioidRiskToolCategory(b *testing.B) {
	b.ReportAllocs()
	total := 5
	for b.Loop() {
		c, _ := OpioidRiskToolCategory(total)
		sinkString = string(c)
	}
}
