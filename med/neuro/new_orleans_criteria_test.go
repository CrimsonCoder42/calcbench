package neuro

import "testing"

// Haydel 2000 abstract: seven findings; count each, HAND-DERIVED.
func TestNewOrleansCriteriaCount_example(t *testing.T) {
	all := NewOrleansFindings{true, true, true, true, true, true, true}
	if got := NewOrleansCriteriaCount(all); got != 7 {
		t.Errorf("all = %d, want 7", got)
	}
	if got := NewOrleansCriteriaCount(NewOrleansFindings{}); got != 0 {
		t.Errorf("none = %d, want 0", got)
	}
	if got := NewOrleansCriteriaCount(NewOrleansFindings{Headache: true, Vomiting: true}); got != 2 {
		t.Errorf("two = %d, want 2", got)
	}
	if got := NewOrleansCriteriaCount(NewOrleansFindings{Seizure: true}); got != 1 {
		t.Errorf("seizure = %d, want 1", got)
	}
}

// parity: ledger new_orleans_head_ct, source cases 1–3 (consumer counts the same seven; its "≥3 multiple"
// band is not in the source and is not reproduced).
func TestNewOrleansCriteriaCount_parity(t *testing.T) {
	cases := []struct {
		f    NewOrleansFindings
		want int
	}{
		{NewOrleansFindings{Headache: true, AgeOver60: true, Seizure: true}, 3}, // parity: ledger new_orleans_head_ct, source case 1
		{NewOrleansFindings{TraumaAboveClavicles: true}, 1},                     // parity: ledger new_orleans_head_ct, source case 2
		{NewOrleansFindings{}, 0},                                               // parity: ledger new_orleans_head_ct, source case 3
	}
	for i, c := range cases {
		if got := NewOrleansCriteriaCount(c.f); got != c.want {
			t.Errorf("case %d = %d, want %d", i+1, got, c.want)
		}
	}
}

func TestNewOrleansCriteriaCount_noAlloc(t *testing.T) {
	f := NewOrleansFindings{Headache: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = NewOrleansCriteriaCount(f) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkNewOrleansCriteriaCount(b *testing.B) {
	b.ReportAllocs()
	f := NewOrleansFindings{Headache: true}
	for b.Loop() {
		sinkInt = NewOrleansCriteriaCount(f)
	}
}
