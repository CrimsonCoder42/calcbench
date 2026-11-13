package invest

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// bigBondPrice is the oracle: Σ c·vᵏ + F·vⁿ in 256-bit big.Float, v = 1/(1 + y/m), c = F·cr/m, from the
// exact binary values of the float64 inputs.
func bigBondPrice(faceCents int64, cr, y float64, n, m int) *big.Float {
	const prec = 256
	nf := func(x float64) *big.Float { return new(big.Float).SetPrec(prec).SetFloat64(x) }
	mm := new(big.Float).SetPrec(prec).SetInt64(int64(m))
	face := new(big.Float).SetPrec(prec).SetInt64(faceCents)
	c := new(big.Float).SetPrec(prec).Mul(face, nf(cr))
	c.Quo(c, mm)
	i := new(big.Float).SetPrec(prec).Quo(nf(y), mm)
	v := new(big.Float).SetPrec(prec).Add(nf(1), i)
	v.Quo(nf(1), v)
	sum := new(big.Float).SetPrec(prec)
	vk := nf(1)
	for k := 1; k <= n; k++ {
		vk.Mul(vk, v)
		term := new(big.Float).SetPrec(prec).Mul(c, vk)
		sum.Add(sum, term)
	}
	sum.Add(sum, new(big.Float).SetPrec(prec).Mul(face, vk))
	return sum
}

// OpenStax Principles of Finance §10.2 "Bond Valuation": the 3M bond ($1,000, 2.25 % annual coupon,
// 15 years, YTM 1.24 %) is priced at $1,137.47, and the calculator example (4 % semiannual coupon,
// 30 periods, 5 % yield) at $895.35.
func TestBondPrice_openStax(t *testing.T) {
	got, err := BondPrice(100000, 0.0225, 0.0124, 15, 1)
	noErr(t, "3M", err)
	if got != 113747 {
		t.Errorf("3M bond = %d cents, want 113747", got)
	}
	got, err = BondPrice(100000, 0.04, 0.05, 30, 2)
	noErr(t, "semiannual", err)
	if got != 89535 {
		t.Errorf("semiannual bond = %d cents, want 89535", got)
	}
}

func TestBondPrice_cases(t *testing.T) {
	cases := []struct {
		name  string
		face  int64
		cr, y float64
		n, m  int
		want  int64
	}{
		{"at par", 100000, 0.06, 0.06, 20, 2, 100000},
		{"zero yield", 100000, 0.05, 0, 10, 1, 150000}, // c·n + F
		{"zero coupon zero yield", 100000, 0, 0, 10, 1, 100000},
		{"one period", 100000, 0.10, 0.10, 1, 1, 100000}, // 110000/1.1
		{"zero face", 0, 0.05, 0.05, 10, 1, 0},
		{"zero coupon", 100000, 0, 0.10, 1, 1, 90909},     // 100000/1.1 = 90909.09
		{"negative yield", 100000, 0, -0.5, 1, 1, 200000}, // 100000/0.5
	}
	for _, c := range cases {
		got, err := BondPrice(c.face, c.cr, c.y, c.n, c.m)
		if err != nil || got != c.want {
			t.Errorf("%s: BondPrice = %d, %v; want %d", c.name, got, err, c.want)
		}
	}
	// 2^53 face is accepted; the price at par stays 2^53.
	got, err := BondPrice(maxCents, 0, 0, 1, 1)
	if err != nil || got != maxCents {
		t.Errorf("max face: %d, %v", got, err)
	}
}

func TestBondPrice_errors(t *testing.T) {
	_, err := BondPrice(maxCents+1, 0.05, 0.05, 10, 1)
	wantErr(t, "face > 2^53", err, ErrOutOfRange)
	_, err = BondPrice(-1, 0.05, 0.05, 10, 1)
	wantErr(t, "face < 0", err, ErrNegative)
	_, err = BondPrice(100000, -0.01, 0.05, 10, 1)
	wantErr(t, "coupon < 0", err, ErrNegative)
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err = BondPrice(100000, bad, 0.05, 10, 1)
		wantErr(t, "coupon not finite", err, ErrRate)
		_, err = BondPrice(100000, 0.05, bad, 10, 1)
		wantErr(t, "yield not finite", err, ErrRate)
	}
	_, err = BondPrice(100000, 0.05, -2, 10, 2) // periodic −1
	wantErr(t, "yield/m = −1", err, ErrRate)
	_, err = BondPrice(100000, 0.05, 0.05, 0, 1)
	wantErr(t, "periods 0", err, ErrPeriods)
	_, err = BondPrice(100000, 0.05, 0.05, 10, 0)
	wantErr(t, "periodsPerYear 0", err, ErrPeriods)
	_, err = BondPrice(maxCents, 0, -0.99, 10, 1) // 2^53 · 100^10
	wantErr(t, "overflow", err, ErrOverflow)
}

// Oracle: the price is the 256-bit exact value rounded half-even (F-5 contract; prices ≤ about 2³⁶ cents).
func TestBondPrice_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	ms := []int{1, 2, 4, 12}
	for range 300 {
		face := r.Int64N(100_000_000_00) + 1
		cr := float64(r.IntN(1500)) / 10000
		y := float64(r.IntN(2000)-100) / 10000
		m := ms[r.IntN(len(ms))]
		n := r.IntN(120) + 1
		got, err := BondPrice(face, cr, y, n, m)
		if err != nil {
			t.Fatalf("BondPrice(%d, %v, %v, %d, %d): %v", face, cr, y, n, m, err)
		}
		exact, _ := bigBondPrice(face, cr, y, n, m).Float64()
		checkRounded(t, "BondPrice", got, exact)
	}
}

func TestBondPrice_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkCents, sinkErr = BondPrice(100000, 0.04, 0.05, 30, 2) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkBondPrice(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, sinkErr = BondPrice(100000, 0.04, 0.05, 30, 2)
	}
}
