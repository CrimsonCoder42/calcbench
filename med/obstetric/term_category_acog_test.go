package obstetric

import (
	"errors"
	"testing"
)

// sinkCategory keeps benchmark results observable.
var sinkCategory TermCategory

// ACOG CO 579 (PMID 24150030): both sides of every category edge, compared exactly.
func TestTermCategoryACOG_edges(t *testing.T) {
	tests := []struct {
		weeks, days int
		want        TermCategory
	}{
		{37, 0, TermEarly},
		{38, 6, TermEarly},
		{39, 0, TermFull},
		{40, 6, TermFull},
		{41, 0, TermLate},
		{41, 6, TermLate},
		{42, 0, TermPostterm},
		{45, 3, TermPostterm},
	}
	for _, tc := range tests {
		if got, err := TermCategoryACOG(tc.weeks, tc.days); err != nil || got != tc.want {
			t.Errorf("%d %d/7: got (%q, %v), want %q", tc.weeks, tc.days, got, err, tc.want)
		}
	}
}

// The labels are the source's words, verbatim.
func TestTermCategoryACOG_labels(t *testing.T) {
	if TermEarly != "early term" || TermFull != "full term" || TermLate != "late term" || TermPostterm != "postterm" {
		t.Error("labels differ from ACOG CO 579")
	}
}

// parity: ledger pregnancy_dates, source cases 6–10 (ega_weeks 37, 40, 41, 42, 42 with ega_days 0, 0, 0, 0, 6).
// The source's bands ("Term Pregnancy" 37–40, "Late-Term" 41, "Post-Term Pregnancy" ≥ 42) are the app's
// labels; the primary's labels are asserted. Source cases 1–5 (weeks 0, 10, 13, 14, 28: trimester bands) have
// no CO 579 category and return ErrOutOfRange.
func TestTermCategoryACOG_parity(t *testing.T) {
	tests := []struct {
		name        string
		weeks, days int
		want        TermCategory
	}{
		{"source case 6 (source: Term Pregnancy)", 37, 0, TermEarly},
		{"source case 7 (source: Term Pregnancy)", 40, 0, TermFull},
		{"source case 8 (source: Late-Term)", 41, 0, TermLate},
		{"source case 9 (source: Post-Term Pregnancy)", 42, 0, TermPostterm},
		{"source case 10 (source: Post-Term Pregnancy)", 42, 6, TermPostterm},
	}
	for _, tc := range tests {
		if got, err := TermCategoryACOG(tc.weeks, tc.days); err != nil || got != tc.want {
			t.Errorf("%s: got (%q, %v), want %q", tc.name, got, err, tc.want)
		}
	}
	for i, w := range []struct{ weeks, days int }{{0, 1}, {10, 0}, {13, 6}, {14, 0}, {28, 0}} {
		if _, err := TermCategoryACOG(w.weeks, w.days); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("source case %d (%d %d/7): err = %v, want ErrOutOfRange", i+1, w.weeks, w.days, err)
		}
	}
}

func TestTermCategoryACOG_errors(t *testing.T) {
	tests := []struct {
		name        string
		weeks, days int
	}{
		{"36 6/7", 36, 6},
		{"negative weeks", -1, 0},
		{"days -1", 39, -1},
		{"days 7", 39, 7},
	}
	for _, tc := range tests {
		if _, err := TermCategoryACOG(tc.weeks, tc.days); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%s: err = %v, want ErrOutOfRange", tc.name, err)
		}
	}
}

func TestTermCategoryACOG_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkCategory, _ = TermCategoryACOG(39, 2) }); n != 0 {
		t.Errorf("TermCategoryACOG allocates %v times, want 0", n)
	}
}

func BenchmarkTermCategoryACOG(b *testing.B) {
	b.ReportAllocs()
	w, d := 39, 2
	for b.Loop() {
		sinkCategory, _ = TermCategoryACOG(w, d)
	}
}
