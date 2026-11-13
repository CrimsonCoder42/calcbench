package psychiatry

import (
	"errors"
	"testing"
)

// Wu 2023 / Zhang 2023: five items 0–4, total 0–20.
func TestCoronavirusAnxietyScaleScore_example(t *testing.T) {
	cases := []struct {
		in   [5]int
		want int
	}{{[5]int{}, 0}, {[5]int{4, 4, 4, 4, 4}, 20}, {[5]int{2, 2, 2, 2, 1}, 9}, {[5]int{2, 2, 2, 1, 1}, 8}}
	for _, c := range cases {
		if got, err := CoronavirusAnxietyScaleScore(c.in); err != nil || got != c.want {
			t.Errorf("CoronavirusAnxietyScaleScore(%v) = (%d, %v), want %d", c.in, got, err, c.want)
		}
	}
	for _, b := range [][5]int{{5, 0, 0, 0, 0}, {0, 0, 0, 0, -1}} {
		if _, err := CoronavirusAnxietyScaleScore(b); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%v: %v, want ErrOutOfRange", b, err)
		}
	}
}

// parity: ledger cas_covid_anxiety, source cases 1–6 (source cas_covid_anxiety_test.go "score"; its 5–8
// "Elevated" band is not in the sources).
func TestCoronavirusAnxietyScaleScore_parity(t *testing.T) {
	cases := []struct {
		in   [5]int
		want int
	}{
		{[5]int{0, 0, 0, 0, 0}, 0},  // parity: ledger cas_covid_anxiety, source case 1
		{[5]int{1, 1, 1, 1, 0}, 4},  // parity: ledger cas_covid_anxiety, source case 2
		{[5]int{1, 1, 1, 1, 1}, 5},  // parity: ledger cas_covid_anxiety, source case 3
		{[5]int{2, 2, 2, 2, 0}, 8},  // parity: ledger cas_covid_anxiety, source case 4
		{[5]int{2, 2, 2, 2, 1}, 9},  // parity: ledger cas_covid_anxiety, source case 5
		{[5]int{4, 4, 4, 4, 4}, 20}, // parity: ledger cas_covid_anxiety, source case 6
	}
	for i, c := range cases {
		if got, err := CoronavirusAnxietyScaleScore(c.in); err != nil || got != c.want {
			t.Errorf("case %d = (%d, %v), want %d", i+1, got, err, c.want)
		}
	}
}

func TestCoronavirusAnxietyScaleScore_noAlloc(t *testing.T) {
	in := [5]int{2, 2, 2, 2, 1}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = CoronavirusAnxietyScaleScore(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCoronavirusAnxietyScaleScore(b *testing.B) {
	b.ReportAllocs()
	in := [5]int{2, 2, 2, 2, 1}
	for b.Loop() {
		sinkInt, _ = CoronavirusAnxietyScaleScore(in)
	}
}
