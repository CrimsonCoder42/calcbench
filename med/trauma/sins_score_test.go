package trauma

import "testing"

// Fisher 2014 / Kim 2021 Table 1: every printed point value accepted, every other refused. HAND-DERIVED.
func TestSINSScore_items(t *testing.T) {
	printed := [6][]int{{0, 1, 2, 3}, {0, 1, 3}, {0, 1, 2}, {0, 2, 4}, {0, 1, 2, 3}, {0, 1, 3}}
	for item := range 6 {
		for v := -1; v <= 5; v++ {
			var a [6]int
			a[item] = v
			got, err := SINSScore(a[0], a[1], a[2], a[3], a[4], a[5])
			ok := false
			for _, p := range printed[item] {
				ok = ok || p == v
			}
			if ok && (err != nil || got != v) {
				t.Errorf("item %d = %d: (%d, %v)", item, v, got, err)
			}
			if !ok {
				checkErr(t, "unprinted value", err, ErrOutOfRange)
			}
		}
	}
	if got, _ := SINSScore(3, 3, 2, 4, 3, 3); got != 18 {
		t.Errorf("maximum = %d, want 18", got)
	}
}

// Doc comment example: 2+3+2+0+2+1 = 10 (HAND-DERIVED).
func TestSINSScore_example(t *testing.T) {
	if got, err := SINSScore(2, 3, 2, 0, 2, 1); err != nil || got != 10 {
		t.Errorf("got (%d, %v), want 10", got, err)
	}
}

// parity: ledger sins_score, source example total 11 (the app's item points equal the table): mobile 2 +
// mechanical pain 3 + lytic 2 + normal alignment 0 + > 50% collapse 3 + unilateral 1 = 11.
func TestSINSScore_parity(t *testing.T) {
	if got, err := SINSScore(2, 3, 2, 0, 3, 1); err != nil || got != 11 {
		t.Errorf("got (%d, %v), want 11", got, err)
	}
}

func TestSINSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = SINSScore(2, 3, 2, 0, 2, 1) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSINSScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = SINSScore(2, 3, 2, 0, 2, 1)
	}
}
