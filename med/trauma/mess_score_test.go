package trauma

import "testing"

// Topal 2010 / Okereke 2022 MESS table: edges of age and ischemia time. HAND-DERIVED.
func TestMESSScore_table(t *testing.T) {
	tests := []struct {
		sk, is int
		h      float64
		sh     int
		age    float64
		want   int
	}{
		{1, 1, 0, 0, 20, 2},
		{4, 3, 0, 2, 20, 9},
		{1, 3, 6, 0, 20, 4},
		{1, 3, 6.01, 0, 20, 7},
		{1, 1, 0, 0, 29.9, 2},
		{1, 1, 0, 0, 30, 3},
		{1, 1, 0, 0, 50, 3},
		{1, 1, 0, 0, 50.1, 4},
		{4, 3, 7, 2, 60, 14},
	}
	for _, tc := range tests {
		if got, err := MESSScore(tc.sk, tc.is, tc.h, tc.sh, tc.age); err != nil || got != tc.want {
			t.Errorf("%+v: (%d, %v), want %d", tc, got, err, tc.want)
		}
	}
}

// Doc comment example: 3 + 2×2 + 1 + 1 = 9 (HAND-DERIVED).
func TestMESSScore_example(t *testing.T) {
	if got, err := MESSScore(3, 2, 8, 1, 40); err != nil || got != 9 {
		t.Errorf("got (%d, %v), want 9", got, err)
	}
}

// parity: ledger mess_score, source example 3 + (2×2) + 1 + 1 = 9 agrees. At age exactly 50 the app gives 2 and
// the table 1 (30–50 → 1): asserted 1 here (CONSUMER DEFECT). The app's "No ischemia (0)" option is refused.
func TestMESSScore_parity(t *testing.T) {
	if got, _ := MESSScore(3, 2, 8, 1, 40); got != 9 {
		t.Errorf("got %d, want 9", got)
	}
	if got, _ := MESSScore(2, 1, 0, 0, 50); got != 4 {
		t.Errorf("age 50: got %d, want 4 (app: 5)", got)
	}
	_, err := MESSScore(2, 0, 0, 0, 40)
	checkErr(t, "ischemia 0", err, ErrOutOfRange)
}

func TestMESSScore_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, err := MESSScore(1, 1, v, 0, 20)
		checkErr(t, "hours", err, ErrNotFinite)
		_, err = MESSScore(1, 1, 0, 0, v)
		checkErr(t, "age", err, ErrNotFinite)
	}
	cases := [][5]float64{{0, 1, 0, 0, 20}, {5, 1, 0, 0, 20}, {1, 4, 0, 0, 20}, {1, 1, 0, -1, 20}, {1, 1, 0, 3, 20}, {1, 1, -1, 0, 20}, {1, 1, 0, 0, -1}}
	for _, c := range cases {
		_, err := MESSScore(int(c[0]), int(c[1]), c[2], int(c[3]), c[4])
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestMESSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MESSScore(3, 2, 8, 1, 40) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMESSScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = MESSScore(3, 2, 8, 1, 40)
	}
}
