package psychiatry

import "testing"

// Tian 2023 / Islam 2024: yes = 1, except items 1, 5, 7, 11, 13 (no = 1).
func TestGDS15Score_example(t *testing.T) {
	var allNo, allYes [15]bool
	for i := range allYes {
		allYes[i] = true
	}
	if got := GDS15Score(allYes); got != 10 {
		t.Errorf("all yes = %d, want 10", got)
	}
	if got := GDS15Score(allNo); got != 5 {
		t.Errorf("all no = %d, want 5", got)
	}
	max := allYes
	for _, i := range []int{0, 4, 6, 10, 12} {
		max[i] = false
	}
	if got := GDS15Score(max); got != 15 {
		t.Errorf("max = %d, want 15", got)
	}
	min := allNo
	for _, i := range []int{0, 4, 6, 10, 12} {
		min[i] = true
	}
	if got := GDS15Score(min); got != 0 {
		t.Errorf("min = %d, want 0", got)
	}
}

// gdsAnswers turns per-item points (the consumer's already-reversed 0/1 values) into yes/no answers.
func gdsAnswers(points [15]int) [15]bool {
	var a [15]bool
	for i, p := range points {
		reverse := i == 0 || i == 4 || i == 6 || i == 10 || i == 12
		a[i] = (p == 1) != reverse
	}
	return a
}

// parity: ledger gds15, source cases 1–8 (source gds15_test.go sums per-item points q1..q15 whose reverse
// scoring lives in the consumer's pickers; each case's points are turned back into answers here).
func TestGDS15Score_parity(t *testing.T) {
	firstN := func(n int) [15]int {
		var p [15]int
		for i := 0; i < n; i++ {
			p[i] = 1
		}
		return p
	}
	for i, n := range []int{0, 4, 5, 7, 8, 9, 10, 15} { // parity: ledger gds15, source cases 1–8 (zero … max)
		if got := GDS15Score(gdsAnswers(firstN(n))); got != n {
			t.Errorf("case %d = %d, want %d", i+1, got, n)
		}
	}
}

func TestGDS15Score_noAlloc(t *testing.T) {
	in := [15]bool{true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = GDS15Score(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkGDS15Score(b *testing.B) {
	b.ReportAllocs()
	in := [15]bool{true}
	for b.Loop() {
		sinkInt = GDS15Score(in)
	}
}
