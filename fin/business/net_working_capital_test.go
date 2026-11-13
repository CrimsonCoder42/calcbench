package business

import (
	"errors"
	"math/rand/v2"
	"testing"
)

// Principles of Accounting Vol. 1, Appendix A, Banyan Goods: $200,000 − $100,000 = $100,000. Exact cents.
func TestNetWorkingCapital_sourceExample(t *testing.T) {
	got, err := NetWorkingCapital(20000000, 10000000)
	if err != nil || got != 10000000 {
		t.Fatalf("got %v, %v; want 10000000", got, err)
	}
}

func TestNetWorkingCapital_table(t *testing.T) {
	cases := []struct {
		name string
		a, b int64
		want int64
	}{
		{"zero", 0, 0, 0},
		{"negative result", 10, 30, -20},
		{"max", maxExactCents, 0, maxExactCents},
		{"min", 0, maxExactCents, -maxExactCents},
		{"equal extremes", maxExactCents, maxExactCents, 0},
	}
	for _, c := range cases {
		got, err := NetWorkingCapital(c.a, c.b)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestNetWorkingCapital_errors(t *testing.T) {
	cases := []struct {
		name string
		a, b int64
		want error
	}{
		{"currentAssetsCents above 2^53", maxExactCents + 1, 0, ErrOutOfRange},
		{"currentLiabilitiesCents above 2^53", 0, maxExactCents + 1, ErrOutOfRange},
		{"negative currentAssetsCents", -1, 0, ErrNegative},
		{"negative currentLiabilitiesCents", 0, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := NetWorkingCapital(c.a, c.b); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: the integer difference on 500 fixed-seed inputs.
func TestNetWorkingCapital_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		a := int64(rng.Uint64N(maxExactCents + 1))
		b := int64(rng.Uint64N(maxExactCents + 1))
		if got, err := NetWorkingCapital(a, b); err != nil || got != a-b {
			t.Fatalf("(%d, %d) = %v, %v; want %v", a, b, got, err, a-b)
		}
	}
}

func TestNetWorkingCapital_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = NetWorkingCapital(20000000, 10000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkNetWorkingCapital(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = NetWorkingCapital(20000000, 10000000)
	}
}
