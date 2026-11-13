package loan

import (
	"errors"
	"testing"
)

// Hand-checked from the CFPB formula (principal + interest + mortgage insurance + escrow of homeowners
// insurance and tax): $1,500 + $50 + ($3,600 + $1,200)/12 = $1,950.00.
func TestPITIPayment_example(t *testing.T) {
	got, err := PITIPayment(150_000, 360_000, 120_000, 5_000)
	if err != nil || got != 195_000 {
		t.Fatalf("got %d, %v; want 195000", got, err)
	}
}

// The escrow twelfth is rounded half-even exactly: 18/12 = 1.5 → 2, 30/12 = 2.5 → 2, 7/12 → 1, 5/12 → 0.
func TestPITIPayment_halfEven(t *testing.T) {
	cases := []struct{ annual, want int64 }{{18, 2}, {30, 2}, {7, 1}, {5, 0}, {6, 0}, {66, 6}, {78, 6}, {90, 8}, {0, 0}}
	for _, c := range cases {
		got, err := PITIPayment(0, c.annual, 0, 0)
		if err != nil || got != c.want {
			t.Errorf("annual %d: got %d, %v; want %d", c.annual, got, err, c.want)
		}
		got, err = PITIPayment(0, 0, c.annual, 0)
		if err != nil || got != c.want {
			t.Errorf("insurance %d: got %d, %v; want %d", c.annual, got, err, c.want)
		}
	}
}

// Largest inputs: P&I of 2⁵³ alone fits; tax and insurance of 2⁵³ each give an escrow of
// 2⁵⁴/12 = 1501199875790165 remainder 4 (rounds down), computed without float64.
func TestPITIPayment_extremes(t *testing.T) {
	if got, err := PITIPayment(maxExactCents, 0, 0, 0); err != nil || got != maxExactCents {
		t.Errorf("max P&I: got %d, %v", got, err)
	}
	if got, err := PITIPayment(0, maxExactCents, maxExactCents, 0); err != nil || got != 1_501_199_875_790_165 {
		t.Errorf("max escrow: got %d, %v; want 1501199875790165", got, err)
	}
}

func TestPITIPayment_errors(t *testing.T) {
	cases := []struct {
		name             string
		pi, tax, ins, mi int64
		want             error
	}{
		{"piAboveMax", maxExactCents + 1, 0, 0, 0, ErrOutOfRange},
		{"taxAboveMax", 0, maxExactCents + 1, 0, 0, ErrOutOfRange},
		{"insAboveMax", 0, 0, maxExactCents + 1, 0, ErrOutOfRange},
		{"miAboveMax", 0, 0, 0, maxExactCents + 1, ErrOutOfRange},
		{"piNegative", -1, 0, 0, 0, ErrNegative},
		{"taxNegative", 0, -1, 0, 0, ErrNegative},
		{"insNegative", 0, 0, -1, 0, ErrNegative},
		{"miNegative", 0, 0, 0, -1, ErrNegative},
		{"totalOverflow", maxExactCents, 0, 0, 1, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := PITIPayment(c.pi, c.tax, c.ins, c.mi); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestPITIPayment_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkInt, _ = PITIPayment(150_000, 360_000, 120_000, 5_000) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkPITIPayment(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = PITIPayment(150_000, 360_000, 120_000, 5_000)
	}
}
