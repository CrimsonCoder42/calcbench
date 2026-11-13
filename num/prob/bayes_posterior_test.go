package prob

import (
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// Ross §3.3 (Bayes's formula over a partition): priors [0.5, 0.5], likelihoods [0.3, 0.6] → [1/3, 2/3].
// Hand-checked; tolerance 2 ulps.
func TestBayesPosterior_sourceExample(t *testing.T) {
	got, err := BayesPosterior(nil, []float64{0.3, 0.6}, []float64{0.5, 0.5})
	if err != nil || len(got) != 2 || !within(got[0], 1.0/3, 5e-16) || !within(got[1], 2.0/3, 5e-16) {
		t.Errorf("BayesPosterior = (%v, %v), want ([1/3 2/3], nil)", got, err)
	}
}

func TestBayesPosterior_appendsToDst(t *testing.T) {
	dst := []float64{9}
	got, err := BayesPosterior(dst, []float64{1, 0}, []float64{0.25, 0.75})
	if err != nil || !slices.Equal(got, []float64{9, 1, 0}) {
		t.Errorf("BayesPosterior(dst=[9]) = (%v, %v), want ([9 1 0], nil)", got, err)
	}
}

func TestBayesPosterior_errors(t *testing.T) {
	dst := []float64{7}
	got, err := BayesPosterior(dst, []float64{0, 0}, []float64{0.5, 0.5})
	requireErr(t, "every product 0", err, ErrOutOfRange)
	if len(got) != 1 {
		t.Errorf("on error dst length = %d, want 1", len(got))
	}
	_, err = BayesPosterior(nil, []float64{0.5}, []float64{1, 0})
	requireErr(t, "length mismatch", err, ErrOutOfRange)
	_, err = BayesPosterior(nil, nil, nil)
	requireErr(t, "empty", err, ErrOutOfRange)
	_, err = BayesPosterior(nil, []float64{nan}, []float64{1})
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = BayesPosterior(nil, []float64{2}, []float64{1})
	requireErr(t, "likelihood > 1", err, ErrProbability)
	_, err = BayesPosterior(nil, []float64{0.5, 0.5}, []float64{0.6, 0.6})
	requireErr(t, "priors sum 1.2", err, ErrNotDistribution)
}

// Oracle: exact rational posteriors l·p/Σl·p of the float64 inputs. Tolerance 4 ulps (one compensated
// sum, one product and one division per posterior).
func TestBayesPosterior_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(20)
		lik, pri := make([]float64, n), make([]float64, n)
		for i := range n {
			lik[i] = 0.01 + 0.99*rng.Float64()
			pri[i] = 1 / float64(n)
		}
		den := new(big.Rat)
		for i := range n {
			den.Add(den, new(big.Rat).Mul(new(big.Rat).SetFloat64(lik[i]), new(big.Rat).SetFloat64(pri[i])))
		}
		got, err := BayesPosterior(nil, lik, pri)
		if err != nil {
			t.Fatal(err)
		}
		for i := range n {
			num := new(big.Rat).Mul(new(big.Rat).SetFloat64(lik[i]), new(big.Rat).SetFloat64(pri[i]))
			want, _ := num.Quo(num, den).Float64()
			if !within(got[i], want, 1e-15) {
				t.Fatalf("trial %d: posterior[%d] = %v, want %v", trial, i, got[i], want)
			}
		}
	}
}

func TestBayesPosterior_noAllocWithCapacity(t *testing.T) {
	lik, pri := []float64{0.3, 0.6}, []float64{0.5, 0.5}
	buf := make([]float64, 0, 2)
	requireNoAlloc(t, "BayesPosterior", func() { sinkSlice, _ = BayesPosterior(buf[:0], lik, pri) })
}

func BenchmarkBayesPosterior(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}} {
		lik, pri := make([]float64, size.n), make([]float64, size.n)
		for i := range size.n {
			lik[i], pri[i] = 0.5, 1/float64(size.n)
		}
		buf := make([]float64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkSlice, _ = BayesPosterior(buf[:0], lik, pri)
			}
		})
	}
}
