package shopping

import (
	"errors"
	"math"
	"testing"
)

// Cases HAND-DERIVED: tip = bill × bp / 10 000, half to even.
// 5500 × 1500 = 8 250 000 → 825; 1010 × 1500 = 1 515 000 → 151.5 → 152 (151 odd → up);
// 1030 × 1500 = 1 545 000 → 154.5 → 154 (even, stays); 1003 × 1500 → 150.45 → 150; 1007 × 1500 → 151.05 → 151;
// 0 bill → 0; 0 rate → 0; 12345 × 1800 → 2222.1 → 2222; 9999 × 2000 → 1999.8 → 2000.
func TestTipCents_examples(t *testing.T) {
	tests := []struct{ bill, bp, want int64 }{
		{5500, 1500, 825},
		{1010, 1500, 152},
		{1030, 1500, 154},
		{1003, 1500, 150},
		{1007, 1500, 151},
		{0, 1500, 0},
		{5500, 0, 0},
		{12345, 1800, 2222},
		{9999, 2000, 2000},
	}
	for _, tc := range tests {
		if got, err := TipCents(tc.bill, tc.bp); err != nil || got != tc.want {
			t.Errorf("TipCents(%d, %d) = (%d, %v), want %d", tc.bill, tc.bp, got, err, tc.want)
		}
	}
}

// Oracle: against float rounding on bills where the exact value is never a tie (bp = 1700, bill odd multiples avoided by check).
func TestTipCents_oracle(t *testing.T) {
	for bill := int64(0); bill < 20000; bill += 7 {
		for _, bp := range []int64{1000, 1500, 1800, 2000, 2250} {
			p := bill * bp
			if p%10000 == 5000 {
				continue
			}
			want := int64(math.Round(float64(p) / 10000))
			if got, err := TipCents(bill, bp); err != nil || got != want {
				t.Fatalf("TipCents(%d, %d) = (%d, %v), want %d", bill, bp, got, err, want)
			}
		}
	}
}

func TestTipCents_errors(t *testing.T) {
	cases := []struct {
		name     string
		bill, bp int64
		want     error
	}{
		{"negative bill", -1, 1500, ErrOutOfRange},
		{"negative rate", 100, -1, ErrOutOfRange},
		{"overflow", math.MaxInt64, 2, ErrOverflow},
		{"overflow edge", math.MaxInt64/1500 + 1, 1500, ErrOverflow},
	}
	for _, tc := range cases {
		if _, err := TipCents(tc.bill, tc.bp); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
	if _, err := TipCents(math.MaxInt64/1500, 1500); err != nil {
		t.Errorf("largest bill at 15%%: err = %v, want nil", err)
	}
}

func TestTipCents_noAllocation(t *testing.T) {
	var s int64
	if n := testing.AllocsPerRun(100, func() { s, _ = TipCents(5500, 1500) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	_ = s
}

func BenchmarkTipCents(b *testing.B) {
	b.ReportAllocs()
	var s int64
	for b.Loop() {
		s, _ = TipCents(5500, 1500)
	}
	_ = s
}
