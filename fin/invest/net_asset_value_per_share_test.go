package invest

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// SEC Fast Answers "Net Asset Value": NAV $100 million over 10,000,000 shares → $10.00 a share; and
// assets $100 million less liabilities $10 million → NAV $90 million ($9.00 over the same shares).
func TestNetAssetValuePerShare_sec(t *testing.T) {
	got, err := NetAssetValuePerShare(10_000_000_000, 0, 10_000_000)
	if err != nil || got != 1000 {
		t.Errorf("NAV/share = %d, %v; want 1000", got, err)
	}
	got, err = NetAssetValuePerShare(10_000_000_000, 1_000_000_000, 10_000_000)
	if err != nil || got != 900 {
		t.Errorf("NAV/share = %d, %v; want 900", got, err)
	}
}

func TestNetAssetValuePerShare_rounding(t *testing.T) {
	for _, c := range []struct {
		name                 string
		assets, liab, shares int64
		want                 int64
	}{
		{"2.5 → 2", 5, 0, 2, 2},
		{"3.5 → 4", 7, 0, 2, 4},
		{"−2.5 → −2", 0, 5, 2, -2},
		{"−3.5 → −4", 0, 7, 2, -4},
		{"1/3 → 0", 1, 0, 3, 0},
		{"2/3 → 1", 2, 0, 3, 1},
		{"−2/3 → −1", 0, 2, 3, -1},
		{"zero NAV", 100, 100, 7, 0},
		{"max", maxCents, 0, 1, maxCents},
		{"min", 0, maxCents, 1, -maxCents},
		{"huge share count", maxCents, 0, math.MaxInt64, 0},
	} {
		got, err := NetAssetValuePerShare(c.assets, c.liab, c.shares)
		if err != nil || got != c.want {
			t.Errorf("%s: %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// Oracle: big.Rat quotient rounded half-even, on 300 random inputs.
func TestNetAssetValuePerShare_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		a := r.Int64N(maxCents)
		l := r.Int64N(maxCents)
		s := r.Int64N(1_000_000) + 1
		got, err := NetAssetValuePerShare(a, l, s)
		noErr(t, "NAV", err)
		q := new(big.Rat).SetFrac(big.NewInt(a-l), big.NewInt(s))
		f, _ := q.Float64() // exact enough: |q| ≤ 2^53, and ties are checked exactly below
		want := int64(math.RoundToEven(f))
		// Exact tie handling: 2·(nav − want·s) = ±s means a tie, which must go to the even neighbour.
		if rem := (a - l) - want*s; 2*rem > s || 2*rem < -s || (2*rem == s || 2*rem == -s) && want%2 != 0 {
			t.Fatalf("oracle self-check failed for %d/%d", a-l, s)
		}
		if got != want {
			t.Errorf("NetAssetValuePerShare(%d, %d, %d) = %d, want %d", a, l, s, got, want)
		}
	}
}

func TestNetAssetValuePerShare_errors(t *testing.T) {
	_, err := NetAssetValuePerShare(100, 0, 0)
	wantErr(t, "shares 0", err, ErrZeroDenominator)
	_, err = NetAssetValuePerShare(100, 0, -1)
	wantErr(t, "shares < 0", err, ErrNegative)
	_, err = NetAssetValuePerShare(-1, 0, 1)
	wantErr(t, "assets < 0", err, ErrNegative)
	_, err = NetAssetValuePerShare(1, -1, 1)
	wantErr(t, "liabilities < 0", err, ErrNegative)
	_, err = NetAssetValuePerShare(maxCents+1, 0, 1)
	wantErr(t, "assets > 2^53", err, ErrOutOfRange)
	_, err = NetAssetValuePerShare(0, maxCents+1, 1)
	wantErr(t, "liabilities > 2^53", err, ErrOutOfRange)
}

func TestNetAssetValuePerShare_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkCents, sinkErr = NetAssetValuePerShare(10_000_000_000, 1_000_000_000, 10_000_000) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkNetAssetValuePerShare(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, sinkErr = NetAssetValuePerShare(10_000_000_000, 1_000_000_000, 10_000_000)
	}
}
