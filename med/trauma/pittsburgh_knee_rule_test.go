package trauma

import "testing"

// Seaberg 1994 abstract: fall or blunt trauma with inability to ambulate or age <12 or >50. HAND-DERIVED.
func TestPittsburghKneeRule_source(t *testing.T) {
	tests := []struct {
		age            float64
		trauma, unable bool
		want           bool
	}{
		{30, true, false, false},
		{30, true, true, true},
		{11.9, true, false, true},
		{12, true, false, false},
		{50, true, false, false},
		{50.1, true, false, true},
		{70, false, true, false},
	}
	for i, tc := range tests {
		if got, err := PittsburghKneeRule(tc.age, tc.trauma, tc.unable); err != nil || got != tc.want {
			t.Errorf("case %d: (%v, %v), want %v", i, got, err, tc.want)
		}
	}
}

// parity: ledger pittsburgh_knee, source case 1 (app: no mechanism → "Cannot Apply"; here false).
func TestPittsburghKneeRule_parity(t *testing.T) {
	if got, _ := PittsburghKneeRule(70, false, false); got {
		t.Error("no mechanism → want false")
	}
}

func TestPittsburghKneeRule_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, err := PittsburghKneeRule(v, true, false)
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	_, err := PittsburghKneeRule(-0.5, true, false)
	checkErr(t, "range", err, ErrOutOfRange)
}

func TestPittsburghKneeRule_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = PittsburghKneeRule(60, true, false) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPittsburghKneeRule(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, _ = PittsburghKneeRule(60, true, false)
	}
}
