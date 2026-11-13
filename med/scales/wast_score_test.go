package scales

import "testing"

// HAND-DERIVED from the CDC (Basile 2007, pp. 72–74) recode-and-sum rule: 8 items, 1–3 each.
func TestWASTScore_sum(t *testing.T) {
	cases := []struct {
		items [8]int
		want  int
	}{
		{[8]int{1, 1, 1, 1, 1, 1, 1, 1}, 8},
		{[8]int{2, 2, 2, 1, 1, 1, 2, 1}, 12},
		{[8]int{3, 3, 3, 3, 3, 3, 3, 3}, 24},
	}
	for _, c := range cases {
		if got, err := WASTScore(c.items); err != nil || got != c.want {
			t.Errorf("%v: got %d, %v; want %d", c.items, got, err, c.want)
		}
	}
}

// parity: ledger wast, source case 1: the consumer scores each picker index + 1 (1–3) and sums 8 items; all
// middle answers → 16. (Its 8-item cut-off "≥ 13" and short-form "≥ 5" are not in the CDC document and are not
// asserted.)
func TestWASTScore_parity(t *testing.T) {
	if got, err := WASTScore([8]int{2, 2, 2, 2, 2, 2, 2, 2}); err != nil || got != 16 {
		t.Errorf("got %d, %v; want 16", got, err)
	}
}

func TestWASTScore_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"item 0", func() error { _, err := WASTScore([8]int{0, 1, 1, 1, 1, 1, 1, 1}); return err }, ErrOutOfRange},
		{"item 4", func() error { _, err := WASTScore([8]int{1, 1, 1, 1, 1, 1, 1, 4}); return err }, ErrOutOfRange},
	})
}

func TestWASTScore_noAlloc(t *testing.T) {
	items := [8]int{2, 2, 2, 1, 1, 1, 2, 1}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = WASTScore(items) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkWASTScore(b *testing.B) {
	items := [8]int{2, 2, 2, 1, 1, 1, 2, 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = WASTScore(items)
	}
}
