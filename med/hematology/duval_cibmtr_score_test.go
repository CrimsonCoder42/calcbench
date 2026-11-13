package hematology

import (
	"errors"
	"testing"
)

// Duval 2010, Table 5 (AML): CR1 < 6 months 1; poor cytogenetics 1; mismatched unrelated 1, related other
// than HLA-identical sibling 2; circulating blasts 1; Karnofsky or Lansky < 90 1.
func TestDuvalCIBMTRScore_table(t *testing.T) {
	tests := []struct {
		name    string
		d       DuvalDiseaseGroup
		cyto    bool
		donor   DuvalDonorGroup
		bl, prf bool
		want    int
	}{
		{"all zero", DuvalPIFOrCR1Over6Months, false, DuvalDonorSiblingOrMatchedUnrelated, false, false, 0},
		{"CR1 under 6 months", DuvalCR1Under6Months, false, DuvalDonorSiblingOrMatchedUnrelated, false, false, 1},
		{"poor cytogenetics", DuvalPIFOrCR1Over6Months, true, DuvalDonorSiblingOrMatchedUnrelated, false, false, 1},
		{"mismatched unrelated", DuvalPIFOrCR1Over6Months, false, DuvalDonorMismatchedUnrelated, false, false, 1},
		{"related not identical sibling", DuvalPIFOrCR1Over6Months, false, DuvalDonorRelatedNotIdenticalSibling, false, false, 2},
		{"blasts", DuvalPIFOrCR1Over6Months, false, DuvalDonorSiblingOrMatchedUnrelated, true, false, 1},
		{"performance below 90", DuvalPIFOrCR1Over6Months, false, DuvalDonorSiblingOrMatchedUnrelated, false, true, 1},
		{"maximum", DuvalCR1Under6Months, true, DuvalDonorRelatedNotIdenticalSibling, true, true, 6},
	}
	for _, tc := range tests {
		got, err := DuvalCIBMTRScore(tc.d, tc.cyto, tc.donor, tc.bl, tc.prf)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger duval_cibmtr. The source entry has no _test.go and its JSON no example; case 1 is the
// consumer's own worked example (PIF, intermediate cytogenetics, HLA-identical sibling, blasts, Karnofsky 80 → 2).
func TestDuvalCIBMTRScore_parity(t *testing.T) {
	got, err := DuvalCIBMTRScore(DuvalPIFOrCR1Over6Months, false, DuvalDonorSiblingOrMatchedUnrelated, true, true)
	if err != nil || got != 2 {
		t.Errorf("source case 1: got (%d, %v), want 2", got, err)
	}
}

func TestDuvalCIBMTRScore_errors(t *testing.T) {
	for _, d := range []DuvalDiseaseGroup{-1, 2} {
		if _, err := DuvalCIBMTRScore(d, false, DuvalDonorSiblingOrMatchedUnrelated, false, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("disease %d: error = %v, want ErrOutOfRange", d, err)
		}
	}
	for _, g := range []DuvalDonorGroup{-1, 3} {
		if _, err := DuvalCIBMTRScore(DuvalPIFOrCR1Over6Months, false, g, false, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("donor %d: error = %v, want ErrOutOfRange", g, err)
		}
	}
}

func TestDuvalCIBMTRScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = DuvalCIBMTRScore(DuvalCR1Under6Months, true, DuvalDonorMismatchedUnrelated, false, true)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkDuvalCIBMTRScore(b *testing.B) {
	b.ReportAllocs()
	d, g := DuvalCR1Under6Months, DuvalDonorMismatchedUnrelated
	for b.Loop() {
		sinkInt, _ = DuvalCIBMTRScore(d, true, g, false, true)
	}
}
