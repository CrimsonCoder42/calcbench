package economics

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// DOJ/FTC Merger Guidelines 2023, §2.1, n. 13: five equal firms → 2,000 (exact: 20² is exact in binary);
// six equal firms at 16.67 % → 6 × 16.67² = 1667.3334, printed as 1,667 (tolerance 1e-9: 16.67 is not
// exact in binary).
func TestHerfindahlHirschmanIndex_mergerGuidelinesExample(t *testing.T) {
	got, err := HerfindahlHirschmanIndex([]float64{20, 20, 20, 20, 20})
	if err != nil || got != 2000 {
		t.Fatalf("five equal: got %v, %v; want 2000", got, err)
	}
	six := []float64{16.67, 16.67, 16.67, 16.67, 16.67, 16.67}
	got, err = HerfindahlHirschmanIndex(six)
	if err != nil || math.Abs(got-1667.3334) > 1e-9 {
		t.Fatalf("six equal: got %v, %v; want 1667.3334", got, err)
	}
	if math.Round(got) != 1667 {
		t.Fatalf("six equal: %v does not round to the Guidelines' 1,667", got)
	}
}

// Merger Guidelines 2023, §2.1 footnote: merging firms with shares a and b changes the HHI by 2ab;
// 2 × 20 × 5 = 200.
func TestHerfindahlHirschmanIndex_mergerChange(t *testing.T) {
	before, _ := HerfindahlHirschmanIndex([]float64{20, 5, 75})
	after, _ := HerfindahlHirschmanIndex([]float64{25, 75})
	if after-before != 200 {
		t.Fatalf("ΔHHI = %v, want 200", after-before)
	}
}

func TestHerfindahlHirschmanIndex_cases(t *testing.T) {
	cases := []struct {
		name   string
		shares []float64
		want   float64
	}{
		{"monopoly", []float64{100}, 10000},
		{"oneZero", []float64{0}, 0},
		{"two", []float64{50, 50}, 5000},
		{"duplicates", []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10}, 1000},
		{"withZeros", []float64{0, 60, 0, 40}, 5200},
		{"sumBelow100", []float64{30}, 900},
	}
	for _, c := range cases {
		got, err := HerfindahlHirschmanIndex(c.shares)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestHerfindahlHirschmanIndex_errors(t *testing.T) {
	cases := []struct {
		name   string
		shares []float64
		want   error
	}{
		{"nil", nil, ErrDomain},
		{"empty", []float64{}, ErrDomain},
		{"NaN", []float64{50, math.NaN()}, ErrNotFinite},
		{"Inf", []float64{math.Inf(1)}, ErrNotFinite},
		{"negative", []float64{-1, 50}, ErrOutOfRange},
		{"above100", []float64{100.0001}, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := HerfindahlHirschmanIndex(c.shares); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: the exact sum, in 256-bit big.Float, of the float64 squares on 300 random markets (fixed
// seed); the compensated float64 sum must be within 1 ulp of it (Neumaier's bound for n ≤ 2000 terms).
func TestHerfindahlHirschmanIndex_bigOracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := 0; trial < 300; trial++ {
		n := 1 + rng.IntN(2000)
		shares := make([]float64, n)
		for k := range shares {
			shares[k] = rng.Float64() * 100 / float64(n)
		}
		got, err := HerfindahlHirschmanIndex(shares)
		if err != nil {
			t.Fatal(err)
		}
		acc := new(big.Float).SetPrec(256)
		for _, s := range shares {
			acc.Add(acc, new(big.Float).SetPrec(256).SetFloat64(s*s))
		}
		want, _ := acc.Float64()
		if math.Abs(got-want) > math.Nextafter(want, math.Inf(1))-want {
			t.Fatalf("trial %d: got %v, want %v", trial, got, want)
		}
	}
}

func TestHerfindahlHirschmanIndex_noAlloc(t *testing.T) {
	shares := []float64{30, 30, 20, 20}
	if a := testing.AllocsPerRun(100, func() { _, _ = HerfindahlHirschmanIndex(shares) }); a != 0 {
		t.Fatalf("allocs = %v, want 0", a)
	}
}

func BenchmarkHerfindahlHirschmanIndex(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		shares := make([]float64, n)
		for k := range shares {
			shares[k] = 100 / float64(n)
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = HerfindahlHirschmanIndex(shares)
			}
		})
	}
}
