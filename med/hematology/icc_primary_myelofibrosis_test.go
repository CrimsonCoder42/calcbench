package hematology

import (
	"errors"
	"testing"
)

// Arber 2022, Table 5: all 3 major criteria and at least 1 minor; leukoerythroblastosis is a minor criterion
// for overt PMF only.
func TestICCPrimaryMyelofibrosis_table(t *testing.T) {
	type in struct {
		m                            ICCPMFMarrow
		clonal, excl                 bool
		anemia, wbc, spleen, ldh, le bool
	}
	tests := []struct {
		name string
		in   in
		want ICCPMFCategory
	}{
		{"pre-PMF with anemia", in{ICCPMFMarrowPrefibrotic, true, true, true, false, false, false, false}, ICCPMFPrefibrotic},
		{"pre-PMF with leukocytosis", in{ICCPMFMarrowPrefibrotic, true, true, false, true, false, false, false}, ICCPMFPrefibrotic},
		{"pre-PMF with splenomegaly", in{ICCPMFMarrowPrefibrotic, true, true, false, false, true, false, false}, ICCPMFPrefibrotic},
		{"pre-PMF with LDH", in{ICCPMFMarrowPrefibrotic, true, true, false, false, false, true, false}, ICCPMFPrefibrotic},
		{"pre-PMF: leukoerythroblastosis does not count", in{ICCPMFMarrowPrefibrotic, true, true, false, false, false, false, true}, ICCPMFNotMet},
		{"overt with leukoerythroblastosis only", in{ICCPMFMarrowOvert, true, true, false, false, false, false, true}, ICCPMFOvert},
		{"overt with anemia", in{ICCPMFMarrowOvert, true, true, true, false, false, false, false}, ICCPMFOvert},
		{"overt with no minor", in{ICCPMFMarrowOvert, true, true, false, false, false, false, false}, ICCPMFNotMet},
		{"no clonal marker", in{ICCPMFMarrowOvert, false, true, true, true, true, true, true}, ICCPMFNotMet},
		{"exclusion not met", in{ICCPMFMarrowPrefibrotic, true, false, true, true, true, true, true}, ICCPMFNotMet},
		{"no marrow finding", in{ICCPMFMarrowNeither, true, true, true, true, true, true, true}, ICCPMFNotMet},
	}
	for _, tc := range tests {
		x := tc.in
		got, err := ICCPrimaryMyelofibrosis(x.m, x.clonal, x.excl, x.anemia, x.wbc, x.spleen, x.ldh, x.le)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%q, %v), want %q", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger icc_pmf. The source entry has no _test.go and its JSON no example; cases follow the consumer's
// rule (pre-PMF: prefibrotic marrow + clonal + exclusion + ≥1 of 4 minors; overt: overt marrow + clonal +
// exclusion + ≥1 of 5 minors), which agrees with Table 5.
func TestICCPrimaryMyelofibrosis_parity(t *testing.T) {
	if got, _ := ICCPrimaryMyelofibrosis(ICCPMFMarrowPrefibrotic, true, true, false, false, true, false, false); got != ICCPMFPrefibrotic {
		t.Errorf("source case 1: got %q, want pre-PMF", got)
	}
	if got, _ := ICCPrimaryMyelofibrosis(ICCPMFMarrowOvert, true, true, false, false, false, false, true); got != ICCPMFOvert {
		t.Errorf("source case 2: got %q, want overt PMF", got)
	}
}

func TestICCPrimaryMyelofibrosis_errors(t *testing.T) {
	for _, m := range []ICCPMFMarrow{-1, 3} {
		if _, err := ICCPrimaryMyelofibrosis(m, true, true, true, true, true, true, true); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("marrow %d: error = %v, want ErrOutOfRange", m, err)
		}
	}
}

func TestICCPrimaryMyelofibrosis_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		c, _ := ICCPrimaryMyelofibrosis(ICCPMFMarrowOvert, true, true, false, false, false, false, true)
		sinkString = string(c)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkICCPrimaryMyelofibrosis(b *testing.B) {
	b.ReportAllocs()
	m := ICCPMFMarrowOvert
	for b.Loop() {
		c, _ := ICCPrimaryMyelofibrosis(m, true, true, false, false, false, false, true)
		sinkString = string(c)
	}
}
