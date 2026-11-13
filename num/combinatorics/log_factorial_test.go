package combinatorics

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// ln 100! = 363.739375555563, the log of the exact 100! (hand-checked with exact integers). Tolerance
// 1e-12 relative: Lgamma is accurate to a few ulps.
func TestLogFactorial_sourceExample(t *testing.T) {
	got, err := LogFactorial(100)
	if err != nil || math.Abs(got-363.739375555563) > 1e-12*363.74 {
		t.Errorf("LogFactorial(100) = (%v, %v), want (363.739375555563, nil)", got, err)
	}
}

func TestLogFactorial_cases(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want float64
	}{
		{"zero", 0, 0},
		{"one", 1, 0},
		{"two", 2, math.Ln2},
		{"twenty, last table entry", 20, math.Log(2432902008176640000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := LogFactorial(tt.n); err != nil || got != tt.want {
				t.Errorf("LogFactorial(%d) = (%v, %v), want (%v, nil)", tt.n, got, err, tt.want)
			}
		})
	}
	// The largest int still has a finite log factorial (about 3.9e20).
	if got, err := LogFactorial(math.MaxInt); err != nil || math.IsInf(got, 0) || got < 3.8e20 || got > 4.0e20 {
		t.Errorf("LogFactorial(MaxInt) = (%v, %v), want finite ≈ 3.9e20", got, err)
	}
}

func TestLogFactorial_errors(t *testing.T) {
	for _, n := range []int{-1, math.MinInt} {
		if got, err := LogFactorial(n); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("LogFactorial(%d) = (%v, %v), want (0, ErrNegative)", n, got, err)
		}
	}
}

// Against ln of the exact n! in math/big. Tolerance 1e-12 relative (a few ulps of Lgamma, plus the
// conversion of the big value).
func TestLogFactorial_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 2 + rng.IntN(3000)
		want := bigLog(bigFactorial(n))
		got, err := LogFactorial(n)
		if err != nil || math.Abs(got-want) > 1e-12*want {
			t.Fatalf("trial %d: LogFactorial(%d) = (%v, %v), want %v", trial, n, got, err, want)
		}
	}
}

func TestLogFactorial_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkF, _ = LogFactorial(1000) }); a != 0 {
		t.Errorf("LogFactorial allocates %v times, want 0", a)
	}
}

func BenchmarkLogFactorial(b *testing.B) {
	b.Run("n=1000", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkF, _ = LogFactorial(1000)
		}
	})
}
