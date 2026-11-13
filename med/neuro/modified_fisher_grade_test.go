package neuro

import (
	"errors"
	"testing"
)

// Patel 2021 Table 6 / Zhou 2022 §2.3: grades 0–4.
func TestModifiedFisherGrade_example(t *testing.T) {
	cases := []struct {
		sah  int
		ivh  bool
		want int
	}{{SAHNone, false, 0}, {SAHThin, false, 1}, {SAHThin, true, 2}, {SAHThick, false, 3}, {SAHThick, true, 4}}
	for _, c := range cases {
		if got, err := ModifiedFisherGrade(c.sah, c.ivh); err != nil || got != c.want {
			t.Errorf("ModifiedFisherGrade(%d, %v) = (%d, %v), want %d", c.sah, c.ivh, got, err, c.want)
		}
	}
}

func TestModifiedFisherGrade_errors(t *testing.T) {
	for _, s := range []int{-1, 3} {
		if _, err := ModifiedFisherGrade(s, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("sah %d: %v, want ErrOutOfRange", s, err)
		}
	}
	if _, err := ModifiedFisherGrade(SAHNone, true); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("IVH without SAH: %v, want ErrOutOfRange", err)
	}
}

// parity: ledger modified_fisher_sah, source cases 1–4. Source case 5 (no SAH with IVH → consumer grade 0) is
// not graded by the sources; this function refuses it (ledger finding).
func TestModifiedFisherGrade_parity(t *testing.T) {
	cases := []struct {
		sah  int
		ivh  bool
		want int
	}{
		{2, true, 4},  // parity: ledger modified_fisher_sah, source case 1
		{2, false, 3}, // parity: ledger modified_fisher_sah, source case 2
		{1, true, 2},  // parity: ledger modified_fisher_sah, source case 3
		{1, false, 1}, // parity: ledger modified_fisher_sah, source case 4
	}
	for i, c := range cases {
		if got, _ := ModifiedFisherGrade(c.sah, c.ivh); got != c.want {
			t.Errorf("case %d = %d, want %d", i+1, got, c.want)
		}
	}
	if _, err := ModifiedFisherGrade(0, true); err == nil { // parity: ledger modified_fisher_sah, source case 5 (consumer 0; primary side: not graded)
		t.Error("case 5: want an error")
	}
}

func TestModifiedFisherGrade_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = ModifiedFisherGrade(SAHThick, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkModifiedFisherGrade(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = ModifiedFisherGrade(SAHThick, true)
	}
}
