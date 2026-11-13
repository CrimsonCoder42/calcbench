package oncology

import (
	"errors"
	"testing"
)

var sinkInt int

// Lee 2019, Table 2: grade = 1 with fever alone; otherwise the more severe of the hypotension and hypoxia rows.
func TestASTCTCRSGrade_table(t *testing.T) {
	tests := []struct {
		name         string
		fever, after bool
		hypo         CRSHypotension
		ox           CRSHypoxia
		want         int
	}{
		{"fever alone is grade 1", true, false, CRSHypotensionNone, CRSHypoxiaNone, 1},
		{"no vasopressors is grade 2", true, false, CRSHypotensionNoVasopressors, CRSHypoxiaNone, 2},
		{"low-flow is grade 2", true, false, CRSHypotensionNone, CRSHypoxiaLowFlow, 2},
		{"one vasopressor is grade 3", true, false, CRSHypotensionOneVasopressor, CRSHypoxiaNone, 3},
		{"high-flow or mask is grade 3", true, false, CRSHypotensionNone, CRSHypoxiaHighFlowOrMask, 3},
		{"multiple vasopressors is grade 4", true, false, CRSHypotensionMultipleVasopressors, CRSHypoxiaNone, 4},
		{"positive pressure is grade 4", true, false, CRSHypotensionNone, CRSHypoxiaPositivePressure, 4},
		{"more severe event wins: hypotension", true, false, CRSHypotensionOneVasopressor, CRSHypoxiaLowFlow, 3},
		{"more severe event wins: hypoxia", true, false, CRSHypotensionNoVasopressors, CRSHypoxiaPositivePressure, 4},
		{"fever not required after antipyretic or anti-cytokine therapy", false, true, CRSHypotensionOneVasopressor, CRSHypoxiaNone, 3},
		{"after therapy, no hypotension or hypoxia is grade 1", false, true, CRSHypotensionNone, CRSHypoxiaNone, 1},
		{"fever condition unmet is 0", false, false, CRSHypotensionMultipleVasopressors, CRSHypoxiaPositivePressure, 0},
	}
	for _, tc := range tests {
		got, err := ASTCTCRSGrade(tc.fever, tc.after, tc.hypo, tc.ox)
		if err != nil || got != tc.want {
			t.Errorf("%s: ASTCTCRSGrade = (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger crs_grading. The source entry has no _test.go and its JSON no example; these cases are the
// consumer's rule (no fever → no grade; grade = 1 + max(hypotension level, hypoxia level); organ toxicity
// ignored), which agrees with Table 2.
func TestASTCTCRSGrade_parity(t *testing.T) {
	tests := []struct {
		fever    bool
		hypo, ox int
		want     int
	}{
		{false, 2, 2, 0}, // source case 1: no fever → N/A
		{true, 0, 0, 1},  // source case 2
		{true, 1, 2, 3},  // source case 3
		{true, 3, 0, 4},  // source case 4
	}
	for i, tc := range tests {
		got, err := ASTCTCRSGrade(tc.fever, false, CRSHypotension(tc.hypo), CRSHypoxia(tc.ox))
		if err != nil || got != tc.want {
			t.Errorf("case %d: got (%d, %v), want %d", i+1, got, err, tc.want)
		}
	}
}

func TestASTCTCRSGrade_errors(t *testing.T) {
	for _, h := range []CRSHypotension{-1, 4} {
		if _, err := ASTCTCRSGrade(true, false, h, CRSHypoxiaNone); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("hypotension %d: error = %v, want ErrOutOfRange", h, err)
		}
	}
	for _, o := range []CRSHypoxia{-1, 4} {
		if _, err := ASTCTCRSGrade(true, false, CRSHypotensionNone, o); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("hypoxia %d: error = %v, want ErrOutOfRange", o, err)
		}
	}
}

func TestASTCTCRSGrade_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = ASTCTCRSGrade(true, false, CRSHypotensionOneVasopressor, CRSHypoxiaLowFlow)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkASTCTCRSGrade(b *testing.B) {
	b.ReportAllocs()
	h, o := CRSHypotensionOneVasopressor, CRSHypoxiaLowFlow
	for b.Loop() {
		sinkInt, _ = ASTCTCRSGrade(true, false, h, o)
	}
}
