package critical

import (
	"errors"
	"testing"
)

// CPOT (Front Pain Res 2024 Table 1; Asriyanto 2024): four indicators 0–2, total 0–8.
func TestCPOTScore_table(t *testing.T) {
	tests := []struct {
		f, b, m, v int
		want       int
	}{
		{0, 0, 0, 0, 0}, {2, 1, 1, 1, 5}, {2, 2, 2, 2, 8}, {1, 0, 0, 0, 1}, {0, 0, 0, 2, 2}, {1, 1, 1, 0, 3},
	}
	for _, tc := range tests {
		if got, err := CPOTScore(tc.f, tc.b, tc.m, tc.v); err != nil || got != tc.want {
			t.Errorf("CPOTScore(%d,%d,%d,%d) = (%d, %v), want %d", tc.f, tc.b, tc.m, tc.v, got, err, tc.want)
		}
	}
}

// parity: ledger cpot_pain. The source entry has no _test.go and its JSON no example; the consumer sums the
// four 0–2 items (its "intubated" input selects the fourth item's wording only), so case 1 (2,2,2,2) is 8.
func TestCPOTScore_parity(t *testing.T) {
	if got, err := CPOTScore(2, 2, 2, 2); err != nil || got != 8 {
		t.Errorf("source case 1: got (%d, %v), want 8", got, err)
	}
}

func TestCPOTScore_errors(t *testing.T) {
	bad := [][4]int{{-1, 0, 0, 0}, {3, 0, 0, 0}, {0, -1, 0, 0}, {0, 3, 0, 0}, {0, 0, -1, 0}, {0, 0, 3, 0}, {0, 0, 0, -1}, {0, 0, 0, 3}}
	for _, x := range bad {
		if _, err := CPOTScore(x[0], x[1], x[2], x[3]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("CPOTScore(%v) error = %v, want ErrOutOfRange", x, err)
		}
	}
}

func TestCPOTScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = CPOTScore(2, 1, 1, 1) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCPOTScore(b *testing.B) {
	b.ReportAllocs()
	f := 2
	for b.Loop() {
		sinkInt, _ = CPOTScore(f, 1, 1, 1)
	}
}
