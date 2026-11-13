package endocrine

import "testing"

var sinkDDS DiabetesDistressCategory

func fill17(v int) (a [17]int) {
	for i := range a {
		a[i] = v
	}
	return a
}

// HAND-DERIVED: the mean of 17 ratings. Tolerance 1e-15 relative (one integer sum, one divide).
func TestDiabetesDistressScore17_example(t *testing.T) {
	tests := []struct {
		items [17]int
		want  float64
	}{
		{fill17(1), 1},
		{fill17(2), 2},
		{fill17(6), 6},
		{[17]int{6, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 22.0 / 17},
	}
	for _, tc := range tests {
		if got, err := DiabetesDistressScore17(tc.items); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("DiabetesDistressScore17(%v) = (%v, %v), want %v", tc.items, got, err, tc.want)
		}
	}
}

// parity: ledger dds17, source cases 1–7. The app's pickers are 0–5 (index); item = index + 1, as the source
// formula says. Means and labels below are the source's expected values (its labels are the app's wording;
// the categories here carry Fisher 2012's labels).
func TestDiabetesDistressScore17_parity(t *testing.T) {
	cases := []struct {
		idx  [17]int
		mean float64
		want DiabetesDistressCategory
	}{
		{[17]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 17.0 / 17, DiabetesDistressLittleOrNone}, // source case 1
		{[17]int{1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0}, 27.0 / 17, DiabetesDistressLittleOrNone}, // source case 2
		{[17]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 34.0 / 17, DiabetesDistressModerate},     // source case 3
		{[17]int{4, 4, 3, 3, 2, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0}, 38.0 / 17, DiabetesDistressModerate},     // source case 4
		{[17]int{2, 2, 2, 2, 2, 1, 1, 1, 1, 2, 2, 1, 1, 1, 1, 1, 1}, 41.0 / 17, DiabetesDistressModerate},     // source case 5
		{[17]int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}, 51.0 / 17, DiabetesDistressHigh},         // source case 6
		{[17]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, 102.0 / 17, DiabetesDistressHigh},        // source case 7
	}
	for i, tc := range cases {
		var items [17]int
		for j, v := range tc.idx {
			items[j] = v + 1
		}
		got, err := DiabetesDistressScore17(items)
		if err != nil || !relNear(got, tc.mean, 1e-15) {
			t.Errorf("source case %d: mean (%v, %v), want %v", i+1, got, err, tc.mean)
			continue
		}
		if c, err := DiabetesDistressCategoryFisher(got); err != nil || c != tc.want {
			t.Errorf("source case %d: category (%q, %v), want %q", i+1, c, err, tc.want)
		}
	}
}

// Both sides of the 2.0 and 3.0 edges, the 2.9–3.0 gap, and the 1 and 6 range ends.
func TestDiabetesDistressCategoryFisher_edges(t *testing.T) {
	tests := []struct {
		v    float64
		want DiabetesDistressCategory
	}{
		{1, DiabetesDistressLittleOrNone},
		{1.9999, DiabetesDistressLittleOrNone},
		{2.0, DiabetesDistressModerate},
		{2.9, DiabetesDistressModerate},
		{50.0 / 17, DiabetesDistressModerate},
		{2.9999, DiabetesDistressModerate},
		{3.0, DiabetesDistressHigh},
		{6, DiabetesDistressHigh},
	}
	for _, tc := range tests {
		if got, err := DiabetesDistressCategoryFisher(tc.v); err != nil || got != tc.want {
			t.Errorf("DiabetesDistressCategoryFisher(%v) = (%q, %v), want %q", tc.v, got, err, tc.want)
		}
	}
}

func TestDiabetesDistress_errors(t *testing.T) {
	bad0 := fill17(3)
	bad0[0] = 0
	bad7 := fill17(3)
	bad7[16] = 7
	score := func(a [17]int) func() error {
		return func() error { _, err := DiabetesDistressScore17(a); return err }
	}
	cat := func(v float64) func() error {
		return func() error { _, err := DiabetesDistressCategoryFisher(v); return err }
	}
	runErrCases(t, []errCase{
		{"item 0", score(bad0), ErrOutOfRange},
		{"item 7", score(bad7), ErrOutOfRange},
		{"all zero", score([17]int{}), ErrOutOfRange},
		{"cat NaN", cat(nan), ErrNotFinite},
		{"cat +Inf", cat(pInf), ErrNotFinite},
		{"cat -Inf", cat(nInf), ErrNotFinite},
		{"cat 0.99", cat(0.99), ErrOutOfRange},
		{"cat 6.01", cat(6.01), ErrOutOfRange},
	})
}

func TestDiabetesDistress_noAllocation(t *testing.T) {
	items := fill17(3)
	if n := testing.AllocsPerRun(100, func() { sink, _ = DiabetesDistressScore17(items) }); n != 0 {
		t.Errorf("DiabetesDistressScore17 allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkDDS, _ = DiabetesDistressCategoryFisher(2.5) }); n != 0 {
		t.Errorf("DiabetesDistressCategoryFisher allocates %v times, want 0", n)
	}
}

func BenchmarkDiabetesDistressScore17(b *testing.B) {
	b.ReportAllocs()
	items := fill17(3)
	for b.Loop() {
		sink, _ = DiabetesDistressScore17(items)
	}
}

func BenchmarkDiabetesDistressCategoryFisher(b *testing.B) {
	b.ReportAllocs()
	v := 2.5
	for b.Loop() {
		sinkDDS, _ = DiabetesDistressCategoryFisher(v)
	}
}
