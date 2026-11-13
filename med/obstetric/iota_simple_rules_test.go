package obstetric

import (
	"errors"
	"testing"
)

// sinkIOTA keeps benchmark results observable.
var sinkIOTA IOTASimpleRulesResult

// Timmerman 2010 BMJ (PMC3001703), Methods, rules 1–3, one case each plus both sides of every edge.
func TestIOTASimpleRules_example(t *testing.T) {
	tests := []struct {
		name string
		m, b int
		want IOTASimpleRulesResult
	}{
		{"rule 1: one M, no B", 1, 0, IOTAMalignant},
		{"rule 1: five M, no B", 5, 0, IOTAMalignant},
		{"rule 2: one B, no M", 0, 1, IOTABenign},
		{"rule 2: five B, no M", 0, 5, IOTABenign},
		{"rule 3: both present", 1, 1, IOTAInconclusive},
		{"rule 3: both present, all", 5, 5, IOTAInconclusive},
		{"rule 3: none present", 0, 0, IOTAInconclusive},
		{"edge: M present, B appears", 2, 1, IOTAInconclusive},
		{"edge: B present, M appears", 1, 3, IOTAInconclusive},
	}
	for _, tc := range tests {
		if got, err := IOTASimpleRules(tc.m, tc.b); err != nil || got != tc.want {
			t.Errorf("%s: IOTASimpleRules(%d, %d) = (%q, %v), want %q", tc.name, tc.m, tc.b, got, err, tc.want)
		}
	}
}

// Every one of the 36 count pairs agrees with the rule as stated (exhaustive, 6 × 6).
func TestIOTASimpleRules_exhaustive(t *testing.T) {
	for m := 0; m <= 5; m++ {
		for b := 0; b <= 5; b++ {
			want := IOTAInconclusive
			if m > 0 && b == 0 {
				want = IOTAMalignant
			} else if b > 0 && m == 0 {
				want = IOTABenign
			}
			if got, err := IOTASimpleRules(m, b); err != nil || got != want {
				t.Errorf("IOTASimpleRules(%d, %d) = (%q, %v), want %q", m, b, got, err, want)
			}
		}
	}
}

// The labels are the source's words, verbatim.
func TestIOTASimpleRules_labels(t *testing.T) {
	if IOTAMalignant != "malignant" || IOTABenign != "benign" || IOTAInconclusive != "inconclusive" {
		t.Error("labels differ from Timmerman 2010")
	}
}

// parity: ledger iota_simple_rules — the source entry is cite-only (no formula, no test file, no example) and
// its output is the 2016 risk model's percentage, not this classification; no source cases exist.
func TestIOTASimpleRules_parity(t *testing.T) {
	if _, err := IOTASimpleRules(0, 0); err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestIOTASimpleRules_errors(t *testing.T) {
	tests := []struct {
		name string
		m, b int
	}{
		{"M negative", -1, 0},
		{"M above 5", 6, 0},
		{"B negative", 0, -1},
		{"B above 5", 0, 6},
	}
	for _, tc := range tests {
		if _, err := IOTASimpleRules(tc.m, tc.b); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%s: err = %v, want ErrOutOfRange", tc.name, err)
		}
	}
}

func TestIOTASimpleRules_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkIOTA, _ = IOTASimpleRules(1, 0) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkIOTASimpleRules(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkIOTA, _ = IOTASimpleRules(1, 0)
	}
}
