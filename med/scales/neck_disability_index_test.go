package scales

import "testing"

var sinkNDI NDICategory

// HAND-DERIVED: 10 sections × 0–5 (Rodríguez-Sanz 2021 PMC8627131).
func TestNeckDisabilityIndex_sum(t *testing.T) {
	cases := []struct {
		s    [10]int
		want int
	}{
		{[10]int{}, 0},
		{[10]int{2, 1, 2, 1, 0, 1, 2, 1, 1, 0}, 11},
		{[10]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, 50},
	}
	for _, c := range cases {
		if got, err := NeckDisabilityIndex(c.s); err != nil || got != c.want {
			t.Errorf("%v: got %d, %v; want %d", c.s, got, err, c.want)
		}
	}
}

// parity: ledger ndi, source cases 1–8 (ndi_test.go TestNDI_MatchesTSX): totals and bands. The consumer's
// labels are title-case ("No Disability" …); the bands agree at every case.
func TestNeckDisabilityIndex_parity(t *testing.T) {
	cases := []struct {
		s    [10]int
		want int
		cat  NDICategory
	}{
		{[10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 0, NDINoDisability},
		{[10]int{1, 1, 1, 1, 0, 0, 0, 0, 0, 0}, 4, NDINoDisability},
		{[10]int{1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, 5, NDIMildDisability},
		{[10]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, 10, NDIMildDisability},
		{[10]int{2, 2, 2, 2, 2, 2, 1, 1, 1, 0}, 15, NDIModerateDisability},
		{[10]int{3, 3, 3, 3, 3, 2, 2, 2, 2, 2}, 25, NDISevereDisability},
		{[10]int{4, 4, 4, 4, 4, 4, 4, 3, 2, 2}, 35, NDICompleteDisability},
		{[10]int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, 50, NDICompleteDisability},
	}
	for i, c := range cases {
		got, err := NeckDisabilityIndex(c.s)
		if err != nil || got != c.want {
			t.Errorf("case %d: got %d, %v; want %d", i+1, got, err, c.want)
			continue
		}
		if cat, _ := NDIDisability(got); cat != c.cat {
			t.Errorf("case %d: band %q, want %q", i+1, cat, c.cat)
		}
	}
}

func TestNeckDisabilityIndex_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"section 6", func() error { _, err := NeckDisabilityIndex([10]int{6}); return err }, ErrOutOfRange},
		{"section -1", func() error { _, err := NeckDisabilityIndex([10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, -1}); return err }, ErrOutOfRange},
		{"total -1", func() error { _, err := NDIDisability(-1); return err }, ErrOutOfRange},
		{"total 51", func() error { _, err := NDIDisability(51); return err }, ErrOutOfRange},
	})
}

// Both sides of each band edge.
func TestNDIDisability_edges(t *testing.T) {
	cases := []struct {
		total int
		want  NDICategory
	}{
		{0, NDINoDisability}, {4, NDINoDisability}, {5, NDIMildDisability}, {14, NDIMildDisability},
		{15, NDIModerateDisability}, {24, NDIModerateDisability}, {25, NDISevereDisability}, {34, NDISevereDisability},
		{35, NDICompleteDisability}, {50, NDICompleteDisability},
	}
	for _, c := range cases {
		if got, err := NDIDisability(c.total); err != nil || got != c.want {
			t.Errorf("%d: got %q, %v; want %q", c.total, got, err, c.want)
		}
	}
}

func TestNDI_noAlloc(t *testing.T) {
	s := [10]int{2, 1, 2, 1, 0, 1, 2, 1, 1, 0}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = NeckDisabilityIndex(s) }); n != 0 {
		t.Errorf("NeckDisabilityIndex allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkNDI, _ = NDIDisability(20) }); n != 0 {
		t.Errorf("NDIDisability allocs = %v, want 0", n)
	}
}

func BenchmarkNeckDisabilityIndex(b *testing.B) {
	s := [10]int{2, 1, 2, 1, 0, 1, 2, 1, 1, 0}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = NeckDisabilityIndex(s)
	}
}

func BenchmarkNDIDisability(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkNDI, _ = NDIDisability(i % 51)
	}
}
