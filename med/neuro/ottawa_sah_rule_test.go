package neuro

import "testing"

// Perry 2013 abstract: six findings; counts, HAND-DERIVED.
func TestOttawaSAHRuleCount_example(t *testing.T) {
	if got := OttawaSAHRuleCount(OttawaSAHFindings{true, true, true, true, true, true}); got != 6 {
		t.Errorf("all = %d, want 6", got)
	}
	if got := OttawaSAHRuleCount(OttawaSAHFindings{}); got != 0 {
		t.Errorf("none = %d, want 0", got)
	}
	if got := OttawaSAHRuleCount(OttawaSAHFindings{Age40OrOlder: true}); got != 1 {
		t.Errorf("age = %d, want 1", got)
	}
}

// parity: ledger ottawa_sah_rule, source cases 1–3 (consumer counts the same six; its "≥3 multiple" band is
// not in the source).
func TestOttawaSAHRuleCount_parity(t *testing.T) {
	cases := []struct {
		f    OttawaSAHFindings
		want int
	}{
		{OttawaSAHFindings{Age40OrOlder: true, ThunderclapHeadache: true, LimitedNeckFlexion: true}, 3}, // parity: ledger ottawa_sah_rule, source case 1
		{OttawaSAHFindings{OnsetDuringExertion: true}, 1},                                               // parity: ledger ottawa_sah_rule, source case 2
		{OttawaSAHFindings{}, 0}, // parity: ledger ottawa_sah_rule, source case 3
	}
	for i, c := range cases {
		if got := OttawaSAHRuleCount(c.f); got != c.want {
			t.Errorf("case %d = %d, want %d", i+1, got, c.want)
		}
	}
}

func TestOttawaSAHRuleCount_noAlloc(t *testing.T) {
	f := OttawaSAHFindings{Age40OrOlder: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = OttawaSAHRuleCount(f) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkOttawaSAHRuleCount(b *testing.B) {
	b.ReportAllocs()
	f := OttawaSAHFindings{Age40OrOlder: true}
	for b.Loop() {
		sinkInt = OttawaSAHRuleCount(f)
	}
}
