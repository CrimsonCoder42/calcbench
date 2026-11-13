package business

import (
	"errors"
	"math/rand/v2"
	"testing"
)

// Principles of Accounting Vol. 2 §3.1, Hicks Manufacturing, April: $50,000 − $10,000 = $40,000. Exact cents.
func TestContributionMargin_sourceExample(t *testing.T) {
	got, err := ContributionMargin(5000000, 1000000)
	if err != nil || got != 4000000 {
		t.Fatalf("got %v, %v; want 4000000", got, err)
	}
}

func TestContributionMargin_table(t *testing.T) {
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
		got, err := ContributionMargin(c.a, c.b)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestContributionMargin_errors(t *testing.T) {
	cases := []struct {
		name string
		a, b int64
		want error
	}{
		{"salesCents above 2^53", maxExactCents + 1, 0, ErrOutOfRange},
		{"variableCostsCents above 2^53", 0, maxExactCents + 1, ErrOutOfRange},
		{"negative salesCents", -1, 0, ErrNegative},
		{"negative variableCostsCents", 0, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := ContributionMargin(c.a, c.b); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: the integer difference on 500 fixed-seed inputs.
func TestContributionMargin_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		a := int64(rng.Uint64N(maxExactCents + 1))
		b := int64(rng.Uint64N(maxExactCents + 1))
		if got, err := ContributionMargin(a, b); err != nil || got != a-b {
			t.Fatalf("(%d, %d) = %v, %v; want %v", a, b, got, err, a-b)
		}
	}
}

func TestContributionMargin_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = ContributionMargin(5000000, 1000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkContributionMargin(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = ContributionMargin(5000000, 1000000)
	}
}
