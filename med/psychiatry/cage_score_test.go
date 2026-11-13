package psychiatry

import "testing"

// O'Rourke 2022 / Peláez-Zuberbuhler 2025: four yes/no items, 0–4.
func TestCAGEScore_example(t *testing.T) {
	cases := []struct {
		in   [4]bool
		want int
	}{{[4]bool{}, 0}, {[4]bool{true, true, true, true}, 4}, {[4]bool{true, false, true, false}, 2}}
	for _, c := range cases {
		if got := CAGEScore(c.in); got != c.want {
			t.Errorf("CAGEScore(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}

// parity: ledger cage_alcohol, source cases 1–5 (source cage_alcohol_test.go "score": zero … four; its band
// labels are the consumer's own and are not reproduced).
func TestCAGEScore_parity(t *testing.T) {
	cases := []struct {
		in   [4]bool
		want int
	}{
		{[4]bool{}, 0},                       // parity: ledger cage_alcohol, source case 1
		{[4]bool{true}, 1},                   // parity: ledger cage_alcohol, source case 2
		{[4]bool{true, true}, 2},             // parity: ledger cage_alcohol, source case 3
		{[4]bool{true, true, true}, 3},       // parity: ledger cage_alcohol, source case 4
		{[4]bool{true, true, true, true}, 4}, // parity: ledger cage_alcohol, source case 5
	}
	for i, c := range cases {
		if got := CAGEScore(c.in); got != c.want {
			t.Errorf("case %d = %d, want %d", i+1, got, c.want)
		}
	}
}

func TestCAGEScore_noAlloc(t *testing.T) {
	in := [4]bool{true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = CAGEScore(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCAGEScore(b *testing.B) {
	b.ReportAllocs()
	in := [4]bool{true}
	for b.Loop() {
		sinkInt = CAGEScore(in)
	}
}
