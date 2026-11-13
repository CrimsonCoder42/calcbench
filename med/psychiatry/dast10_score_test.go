package psychiatry

import (
	"errors"
	"testing"
)

// Lam 2015: yes = 1 except question 3 (no = 1); levels 0, 1–2, 3–5, 6–8, 9–10.
func TestDAST10Score_example(t *testing.T) {
	var allNo, allYes [10]bool
	for i := range allYes {
		allYes[i] = true
	}
	if got := DAST10Score(allYes); got != 9 {
		t.Errorf("all yes = %d, want 9", got)
	}
	if got := DAST10Score(allNo); got != 1 {
		t.Errorf("all no = %d, want 1", got)
	}
	max := allYes
	max[2] = false
	if got := DAST10Score(max); got != 10 {
		t.Errorf("max = %d, want 10", got)
	}
	min := allNo
	min[2] = true
	if got := DAST10Score(min); got != 0 {
		t.Errorf("min = %d, want 0", got)
	}
}

func TestDAST10LevelOf_example(t *testing.T) {
	cases := []struct {
		total int
		want  DAST10Level
	}{
		{0, DAST10NoProblem}, {1, DAST10Low}, {2, DAST10Low}, {3, DAST10Moderate}, {5, DAST10Moderate},
		{6, DAST10Substantial}, {8, DAST10Substantial}, {9, DAST10Severe}, {10, DAST10Severe},
	}
	for _, c := range cases {
		if got, err := DAST10LevelOf(c.total); err != nil || got != c.want {
			t.Errorf("DAST10LevelOf(%d) = (%q, %v), want %q", c.total, got, err, c.want)
		}
	}
	for _, bad := range []int{-1, 11} {
		if _, err := DAST10LevelOf(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("DAST10LevelOf(%d): %v, want ErrOutOfRange", bad, err)
		}
	}
}

// dastAnswers turns the consumer's per-item points (question 3 already reversed) into yes/no answers.
func dastAnswers(points [10]int) [10]bool {
	var a [10]bool
	for i, p := range points {
		a[i] = (p == 1) != (i == 2)
	}
	return a
}

// parity: ledger dast_10, source cases 1–7 (source dast_10_test.go: zero, one, q3 alone, three, six, nine, ten;
// its levels have the same edges, its labels are the consumer's).
func TestDAST10Score_parity(t *testing.T) {
	first := func(n int) [10]int {
		var p [10]int
		for i := 0; i < n; i++ {
			p[i] = 1
		}
		return p
	}
	cases := []struct {
		points [10]int
		want   int
		level  DAST10Level
	}{
		{first(0), 0, DAST10NoProblem},   // parity: ledger dast_10, source case 1
		{first(1), 1, DAST10Low},         // parity: ledger dast_10, source case 2
		{[10]int{2: 1}, 1, DAST10Low},    // parity: ledger dast_10, source case 3 (q3 alone)
		{first(3), 3, DAST10Moderate},    // parity: ledger dast_10, source case 4
		{first(6), 6, DAST10Substantial}, // parity: ledger dast_10, source case 5
		{first(9), 9, DAST10Severe},      // parity: ledger dast_10, source case 6
		{first(10), 10, DAST10Severe},    // parity: ledger dast_10, source case 7
	}
	for i, c := range cases {
		got := DAST10Score(dastAnswers(c.points))
		if got != c.want {
			t.Errorf("case %d = %d, want %d", i+1, got, c.want)
		}
		if l, _ := DAST10LevelOf(got); l != c.level {
			t.Errorf("case %d level = %q, want %q", i+1, l, c.level)
		}
	}
}

func TestDAST10_noAlloc(t *testing.T) {
	in := [10]bool{true}
	if n := testing.AllocsPerRun(100, func() {
		l, _ := DAST10LevelOf(DAST10Score(in))
		sinkString = string(l)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkDAST10Score(b *testing.B) {
	b.ReportAllocs()
	in := [10]bool{true}
	for b.Loop() {
		sinkInt = DAST10Score(in)
	}
}

func BenchmarkDAST10LevelOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		l, _ := DAST10LevelOf(4)
		sinkString = string(l)
	}
}
