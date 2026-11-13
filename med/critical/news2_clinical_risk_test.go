package critical

import (
	"errors"
	"testing"
)

// RCP 2017 Chart 2: 0–4 "Low"; red score "Low–medium"; 5–6 "Medium"; 7 or more "High". Each edge both sides.
func TestNEWS2ClinicalRisk_example(t *testing.T) {
	tests := []struct {
		agg  int
		red  bool
		want NEWS2Risk
	}{
		{0, false, NEWS2RiskLow}, {4, false, NEWS2RiskLow}, {3, true, NEWS2RiskLowMedium}, {4, true, NEWS2RiskLowMedium},
		{5, false, NEWS2RiskMedium}, {5, true, NEWS2RiskMedium}, {6, false, NEWS2RiskMedium},
		{7, false, NEWS2RiskHigh}, {7, true, NEWS2RiskHigh}, {20, true, NEWS2RiskHigh},
	}
	for _, tc := range tests {
		if got, err := NEWS2ClinicalRisk(tc.agg, tc.red); err != nil || got != tc.want {
			t.Errorf("NEWS2ClinicalRisk(%d, %v) = (%q, %v), want %q", tc.agg, tc.red, got, err, tc.want)
		}
	}
}

func TestNEWS2ClinicalRisk_errors(t *testing.T) {
	for _, agg := range []int{-1, 21} {
		if _, err := NEWS2ClinicalRisk(agg, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("NEWS2ClinicalRisk(%d) error = %v, want ErrOutOfRange", agg, err)
		}
	}
}

func TestNEWS2ClinicalRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { r, _ := NEWS2ClinicalRisk(5, true); sinkString = string(r) }); n != 0 {
		t.Errorf("NEWS2ClinicalRisk allocates %v times, want 0", n)
	}
}

func BenchmarkNEWS2ClinicalRisk(b *testing.B) {
	b.ReportAllocs()
	agg := 5
	for b.Loop() {
		r, _ := NEWS2ClinicalRisk(agg, true)
		sinkString = string(r)
	}
}
