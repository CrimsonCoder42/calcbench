package sample

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

type point struct{ x, y float64 }

func proposeSquare(rng *rand.Rand) point { return point{2*rng.Float64() - 1, 2*rng.Float64() - 1} }

func inDisc(p point, _ float64) bool { return p.x*p.x+p.y*p.y < 1 }

// von Neumann 1951 / Devroye 1986, ch. II, §3: uniform points in the unit disc by rejection from the
// square [−1,1]²; the acceptance probability is π/4, so with maxTries = 1 the success rate over 100 000
// calls is within 5σ of π/4.
func TestRejectionSample_discExample(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const trials = 100_000
	ok := 0
	for range trials {
		p, err := RejectionSample(rng, proposeSquare, inDisc, 1)
		switch {
		case err == nil:
			if !inDisc(p, 0) {
				t.Fatalf("accepted %v outside the disc", p)
			}
			ok++
		case !errors.Is(err, ErrMaxTries):
			t.Fatal(err)
		}
	}
	if !within5Sigma(ok, trials, math.Pi/4) {
		t.Errorf("accepted %d of %d, want about %.0f", ok, trials, trials*math.Pi/4)
	}
}

// The density-ratio form: sample f(x) = 2x on [0, 1) from uniform g with c = 2, accept u < f(x)/(c·g(x))
// = x. The mean of f is 2/3; over 50 000 draws the sample mean is within 5 standard errors
// (σ² = 1/18) of it.
func TestRejectionSample_densityRatio(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const n = 50_000
	sum := 0.0
	for range n {
		x, err := RejectionSample(rng, func(r *rand.Rand) float64 { return r.Float64() },
			func(x, u float64) bool { return u < x }, 1000)
		if err != nil {
			t.Fatal(err)
		}
		sum += x
	}
	if mean, se := sum/n, math.Sqrt(1.0/18/n); math.Abs(mean-2.0/3) > 5*se {
		t.Errorf("mean %v, want 2/3 ± %v", mean, 5*se)
	}
}

func TestRejectionSample_errors(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	never := func(point, float64) bool { return false }
	tests := []struct {
		name     string
		rng      *rand.Rand
		propose  func(*rand.Rand) point
		accept   func(point, float64) bool
		maxTries int
		want     error
	}{
		{"nil rng", nil, proposeSquare, inDisc, 5, ErrNilRand},
		{"nil propose", rng, nil, inDisc, 5, ErrNilFunc},
		{"nil accept", rng, proposeSquare, nil, 5, ErrNilFunc},
		{"zero tries", rng, proposeSquare, inDisc, 0, ErrNonPositive},
		{"negative tries", rng, proposeSquare, inDisc, -3, ErrNonPositive},
		{"never accepted", rng, proposeSquare, never, 50, ErrMaxTries},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := RejectionSample(tt.rng, tt.propose, tt.accept, tt.maxTries); !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

// TestRejectionSample_oracle: the number of tries used is exactly the index of the first accepted
// proposal, checked against a replay of the same stream.
func TestRejectionSample_oracle(t *testing.T) {
	for trial := range 300 {
		a := rand.New(rand.NewPCG(uint64(trial), 3))
		b := rand.New(rand.NewPCG(uint64(trial), 3))
		tries := 0
		got, err := RejectionSample(a, func(r *rand.Rand) int { tries++; return r.IntN(10) },
			func(x int, _ float64) bool { return x == 0 }, 200)
		// Replay: propose IntN(10), then one Float64 per try.
		want, wantErr := -1, error(ErrMaxTries)
		for i := range 200 {
			x := b.IntN(10)
			b.Float64()
			if x == 0 {
				want, wantErr = i+1, nil
				break
			}
		}
		if (err == nil) != (wantErr == nil) || (err == nil && (got != 0 || tries != want)) {
			t.Fatalf("trial %d: got %d after %d tries (%v), want 0 after %d (%v)", trial, got, tries, err, want, wantErr)
		}
	}
}

func BenchmarkRejectionSample(b *testing.B) {
	for _, maxTries := range []int{1, 10, 100} {
		rng := rand.New(rand.NewPCG(1, 2))
		b.Run(fmt.Sprint(maxTries), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				if _, err := RejectionSample(rng, proposeSquare, inDisc, maxTries); err == nil {
					sink++
				}
			}
			benchSink = sink
		})
	}
}
