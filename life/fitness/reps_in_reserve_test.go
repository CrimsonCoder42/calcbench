package fitness

import "testing"

// Zourdos 2016 abstract: "RPE-10 = 0-RIR, RPE-9 = 1-RIR, and so forth"; Helms 2016: "an RIR-based RPE of
// 8-10 (RIR 0-2)". Exact integer comparison.
func TestRepsInReserve_example(t *testing.T) {
	for rpe, want := range map[int]int{10: 0, 9: 1, 8: 2, 7: 3} {
		if got, err := RepsInReserve(rpe); err != nil || got != want {
			t.Errorf("RepsInReserve(%d) = (%d, %v), want %d", rpe, got, err, want)
		}
	}
}

// parity: ledger rpe_rir, source cases 1–5.
func TestRepsInReserve_parity(t *testing.T) {
	cases := []struct{ rpe, rir int }{
		{10, 0}, // parity: ledger rpe_rir, source case 1
		{9, 1},  // parity: ledger rpe_rir, source case 2
		{8, 2},  // parity: ledger rpe_rir, source case 3
		{7, 3},  // parity: ledger rpe_rir, source case 4
		{10, 0}, // parity: ledger rpe_rir, source case 5
	}
	for _, c := range cases {
		if got, err := RepsInReserve(c.rpe); err != nil || got != c.rir {
			t.Errorf("RepsInReserve(%d) = (%d, %v), want %d", c.rpe, got, err, c.rir)
		}
	}
}

// parity: ledger rpe_rir, source refusal cases 6, 5, 1, 11 (9.5 cannot be passed: rpe is int; the unit
// refusal has no counterpart). Edges: 7 accepted (example test), 6 refused; 10 accepted, 11 refused.
func TestRepsInReserve_errors(t *testing.T) {
	call := func(r int) func() error {
		return func() error { _, err := RepsInReserve(r); return err }
	}
	runErrCases(t, []errCase{
		{"rpe 6 (parity)", call(6), ErrOutOfRange},
		{"rpe 5 (parity)", call(5), ErrOutOfRange},
		{"rpe 1 (parity)", call(1), ErrOutOfRange},
		{"rpe 11 (parity)", call(11), ErrOutOfRange},
		{"rpe 0", call(0), ErrOutOfRange},
		{"rpe negative", call(-10), ErrOutOfRange},
	})
}

func TestRepsInReserve_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkI, _ = RepsInReserve(8) }); n != 0 {
		t.Errorf("RepsInReserve allocates %v times, want 0", n)
	}
}

func BenchmarkRepsInReserve(b *testing.B) {
	b.ReportAllocs()
	r := 8
	for b.Loop() {
		sinkI, _ = RepsInReserve(r)
	}
}
